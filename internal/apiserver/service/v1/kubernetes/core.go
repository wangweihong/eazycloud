package kubernetes

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/wangweihong/gotoolbox/pkg/compareutil"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/log"
	"github.com/wangweihong/gotoolbox/pkg/mathutil"
	"github.com/wangweihong/gotoolbox/pkg/paging"
	"github.com/wangweihong/gotoolbox/pkg/sets"
	"github.com/wangweihong/gotoolbox/pkg/waitgroup"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubectl/pkg/util/resource"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/internal/apiserver/store"
	"github.com/wangweihong/eazycloud/internal/pkg/clientset"
	"github.com/wangweihong/eazycloud/internal/pkg/libkubernetes"
)

func (k *kubernetesService) PodCreate(ctx context.Context, req *iapiserver.PodRequest) (*iapiserver.PodInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.PodCreate(ctx, cluster, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sPodToApi(meta, cluster, nil), nil
}

func (k *kubernetesService) PodDelete(ctx context.Context, req *iapiserver.PodRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := clientset.PodDelete(ctx, cluster, req.Resource.Namespace, req.Resource, req.DeleteOpts); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) PodLogList(ctx context.Context, req *iapiserver.PodLogRequest) (*iapiserver.PodLogResponse, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	ret, err := clientset.PodLogList(ctx, cluster, req.Namespace, req.Name, req.ToLogOption())
	if err != nil {
		return nil, err
	}

	resp := &iapiserver.PodLogResponse{}
	resp.LogList = ret
	sort.Slice(resp.LogList, func(i, j int) bool {
		return LogToTimeInt64(resp.LogList[i].Msg) < LogToTimeInt64(resp.LogList[j].Msg)
	})

	if len(resp.LogList) != 0 {
		resp.EndTime = LogToTimeInt64(resp.LogList[len(resp.LogList)-1].Msg)
	}

	return resp, nil
}

func (k *kubernetesService) PodLogStream(ctx context.Context, req *iapiserver.PodLogRequest) (io.ReadCloser, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	ret, err := clientset.PodLogStream(ctx, cluster, req.Namespace, req.Name, req.ToLogOption())
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func LogToTimeInt64(msg string) int64 {
	is := strings.Split(msg, " ")
	if len(is) == 0 {
		return time.Now().UnixNano()
	}
	msg = is[0]

	t, err := time.Parse(time.RFC3339, msg)
	if err != nil {
		return time.Now().UnixNano()
	}

	return t.In(time.Local).UnixNano()
}

func (k *kubernetesService) PodList(ctx context.Context, req *iapiserver.PodListRequest) (*iapiserver.PodListResponse, error) {
	resp := &iapiserver.PodListResponse{}

	clusters, err := getVisitScope(ctx, k.store, req.ResourceListRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	childParent, resourceInfo, err := getAllControllers(ctx, k.store, req.ResourceListRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	wg := waitgroup.RunGenericConcurrently[*iapiserver.Cluster, iapiserver.EachResourceRangeListState[*iapiserver.PodInfo]](
		ctx, clusters, func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.PodInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.PodInfo](cluster.ID, cluster.Name)
			clusterListOne.ClusterName = cluster.Name
			clusterListOne.ClusterUUID = cluster.ID
			resList, err := clientset.PodList(ctx, cluster, req.Namespace, req.ToListOpts())
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			addressNode := make(map[string]*v1.ObjectReference)
			nodeList, _ := clientset.NodeList(ctx, cluster, metav1.ListOptions{})
			if nodeList != nil {
				for _, v := range nodeList.Items {
					for _, addr := range v.Status.Addresses {
						if addr.Type == v1.NodeExternalIP || addr.Type == v1.NodeInternalIP {
							addressNode[addr.Address] = &v1.ObjectReference{
								Kind:            v.Kind,
								Namespace:       v.Namespace,
								Name:            v.Name,
								UID:             v.UID,
								APIVersion:      v.APIVersion,
								ResourceVersion: v.ResourceVersion,
							}
						}
					}
				}
			} else {
				log.Errorf("nodeList error:%v", err.Error())
			}

			var resInfos []*iapiserver.PodInfo
			for i := range resList.Items {
				podInfo := convertK8sPodToApi(&resList.Items[i], cluster, addressNode)
				if filterPod(podInfo, req.Fuzzy, req.FilterVolumeName, req.FilterConfigMap, req.FilterSecret) {
					continue
				}
				podInfo.Controller = getResourceController(string(podInfo.Resource.UID), childParent, resourceInfo)

				resInfos = append(resInfos, podInfo)
			}

			// if req.ShowMonitorData && clientset.IsMonitorServiceReady(cluster.UUID) {
			// 	podHistoryList, err := GetTopkeManager().MonitorPodStateList(ctx, cluster, &iapiserver.Monitoring{MonitoringTimeRange: req.MonitoringTimeRange})
			// 	if err != nil {
			// 		log.Errorf("")
			// 	} else {
			// 		translate := func(namespace, name, metric string) iapiserver.Metric {
			// 			me := podHistoryList[namespace][name]
			// 			me.MetricName = metric
			// 			for _, v := range me.MetricValues {
			// 				if v.Name == metric {
			// 					me.MetricValues = []iapiserver.MetricValue{v}
			// 					return me
			// 				}
			// 			}
			// 			return me
			// 		}
			// 		for _, v := range resInfos {
			// 			v.Metric = []iapiserver.Metric{translate(v.Namespace, v.Name, namespace_pod_cpu_rate.Name), translate(v.Namespace, v.Name, namespace_pod_memory_size.Name)}
			// 		}
			// 	}
			// }

			clusterListOne.TotalCount = len(resInfos)
			clusterListOne.List = resInfos

			return waitgroup.NewGenericResult(clusterListOne, err)
		}, 20*time.Second)

	resp.EachRangeListState = wg.GetSuccessResultList()
	resp.TotalCount = CutPagingSliceResourceList[*iapiserver.PodInfo](resp.EachRangeListState, &resp.List, req.PageNum, req.PageSize, func(i, j int) bool {
		switch req.SortBy {
		case iapiserver.PodSortByState:
			if resp.List[i].PodStatus.Status != resp.List[j].PodStatus.Status { // if equal, compare create time instead
				return compareutil.Compare(resp.List[i].PodStatus.Status, resp.List[j].PodStatus.Status, req.SortDesc)
			}
		case iapiserver.PodSortByRestart:
			if resp.List[i].MaxRestarts != nil && resp.List[j].MaxRestarts != nil && resp.List[i].MaxRestarts != resp.List[j].MaxRestarts { // if equal, compare create time instead
				return compareutil.Compare(*resp.List[i].MaxRestarts, *resp.List[j].MaxRestarts, req.SortDesc)
			}
		case iapiserver.PodSortByHostIP:
			if resp.List[i].Resource.Status.HostIP != resp.List[j].Resource.Status.HostIP { // if equal, compare create time instead
				return compareutil.Compare(resp.List[i].Resource.Status.HostIP, resp.List[j].Resource.Status.HostIP, req.SortDesc)
			}
		case iapiserver.PodSortByCpuRequest:
			request1, _ := resource.PodRequestsAndLimits(resp.List[i].Resource)
			request2, _ := resource.PodRequestsAndLimits(resp.List[j].Resource)

			if request1.Cpu().Size() != request2.Cpu().Size() { // if equal, compare create time instead
				return compareutil.Compare(request1.Cpu().Size(), request2.Cpu().Size(), req.SortDesc)
			}
		case iapiserver.PodSortByPodIP:
			return compareutil.Compare(resp.List[i].Resource.Status.PodIP, resp.List[j].Resource.Status.PodIP, req.SortDesc)
		case iapiserver.PodSortByMemoryRequest:
			request1, _ := resource.PodRequestsAndLimits(resp.List[i].Resource)
			request2, _ := resource.PodRequestsAndLimits(resp.List[j].Resource)

			if request1.Memory().Size() != request2.Memory().Size() { // if equal, compare create time instead
				return compareutil.Compare(request1.Memory().Size(), request2.Memory().Size(), req.SortDesc)
			}
		case iapiserver.PodSortByPhysicalMemory:
		case iapiserver.PodSortByPhysicalCpu:
		}

		return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
	})
	return resp, nil
}

func (k *kubernetesService) PodGet(ctx context.Context, req *iapiserver.ResourceGetRequest) (*iapiserver.PodInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	meta, err := clientset.PodGet(ctx, cluster, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, err
	}

	info := convertK8sPodToApi(meta, cluster, nil)

	//FIXME: better way. why not just get controller from owner reference?
	childParent, resourceInfo, err := getAllControllers(ctx, k.store, iapiserver.ResourceListRequest{
		Cluster:   req.Cluster,
		Namespace: req.Namespace,
	})
	if err != nil {
		return nil, err
	}

	info.Controller = getResourceController(string(meta.UID), childParent, resourceInfo)

	return info, nil
}

func getResourceController(uuid string, childParent map[string]map[string]struct{}, resourceInfo map[string]*iapiserver.ObjectTypeMeta) *iapiserver.ObjectTypeMeta {
	if uuid == "" || childParent == nil || resourceInfo == nil {
		return nil
	}

	parent := ""
	parents := childParent[uuid]
	for len(parents) != 0 {
		for parent = range parents {
			break
		}
		parents = childParent[parent]
	}

	return resourceInfo[parent]
}

// TODO: 全局变量, 定时更新
func getAllControllers(ctx context.Context, st store.Factory, req iapiserver.ResourceListRequest) (map[string]map[string]struct{}, map[string]*iapiserver.ObjectTypeMeta, error) {
	childParent := map[string]map[string]struct{}{}
	resourceInfo := map[string]*iapiserver.ObjectTypeMeta{}

	kinds := map[string]bool{"Pod": true, "Deployment": true, "StatefulSet": true, "Secret": true, "Job": true, "CronJob": true, "ReplicaSet": true}
	clusters, err := st.Kubernetes().List(ctx)
	if err != nil {
		return nil, nil, err
	}

	wg := waitgroup.NewGenericGroup[[]iapiserver.ObjectTypeMeta](ctx)
	wg.WithGlobalTimeout(10 * time.Second)
	for _, cluster := range clusters {
		cluster := cluster
		wg.Start(waitgroup.NewGenericFunc("", func(context.Context) waitgroup.GenericResult[[]iapiserver.ObjectTypeMeta] {
			resourceList, err := clientset.DeploymentList(ctx, cluster, req.Namespace, metav1.ListOptions{})
			if err != nil {
				return waitgroup.NewGenericResult[[]iapiserver.ObjectTypeMeta](nil, err)
			}
			var objects []iapiserver.ObjectTypeMeta
			for _, resource := range resourceList.Items {
				objects = append(objects, iapiserver.ObjectTypeMeta{ObjectMeta: resource.ObjectMeta, TypeMeta: resource.TypeMeta})
			}
			return waitgroup.NewGenericResult[[]iapiserver.ObjectTypeMeta](objects, nil)
		}))

		wg.Start(waitgroup.NewGenericFunc("", func(context.Context) waitgroup.GenericResult[[]iapiserver.ObjectTypeMeta] {
			resourceList, err := clientset.StatefulSetList(ctx, cluster, req.Namespace, metav1.ListOptions{})
			if err != nil {
				return waitgroup.NewGenericResult[[]iapiserver.ObjectTypeMeta](nil, err)
			}
			var objects []iapiserver.ObjectTypeMeta
			for _, resource := range resourceList.Items {
				objects = append(objects, iapiserver.ObjectTypeMeta{ObjectMeta: resource.ObjectMeta, TypeMeta: resource.TypeMeta})
			}
			return waitgroup.NewGenericResult(objects, nil)
		}))
		wg.Start(waitgroup.NewGenericFunc("", func(context.Context) waitgroup.GenericResult[[]iapiserver.ObjectTypeMeta] {
			resourceList, err := clientset.ReplicaSetList(ctx, cluster, req.Namespace, metav1.ListOptions{})
			if err != nil {
				return waitgroup.NewGenericResult[[]iapiserver.ObjectTypeMeta](nil, err)
			}
			var objects []iapiserver.ObjectTypeMeta
			for _, resource := range resourceList.Items {
				objects = append(objects, iapiserver.ObjectTypeMeta{ObjectMeta: resource.ObjectMeta, TypeMeta: resource.TypeMeta})
			}
			return waitgroup.NewGenericResult[[]iapiserver.ObjectTypeMeta](objects, nil)
		}))
		wg.Start(waitgroup.NewGenericFunc("", func(context.Context) waitgroup.GenericResult[[]iapiserver.ObjectTypeMeta] {
			resourceList, err := clientset.JobList(ctx, cluster, req.Namespace, metav1.ListOptions{})
			if err != nil {
				return waitgroup.NewGenericResult[[]iapiserver.ObjectTypeMeta](nil, err)
			}
			var objects []iapiserver.ObjectTypeMeta
			for _, resource := range resourceList.Items {
				objects = append(objects, iapiserver.ObjectTypeMeta{ObjectMeta: resource.ObjectMeta, TypeMeta: resource.TypeMeta})
			}
			return waitgroup.NewGenericResult[[]iapiserver.ObjectTypeMeta](objects, nil)
		}))

		wg.Start(waitgroup.NewGenericFunc("", func(context.Context) waitgroup.GenericResult[[]iapiserver.ObjectTypeMeta] {
			resourceList, err := clientset.CronJobList(ctx, cluster, req.Namespace, metav1.ListOptions{})
			if err != nil {
				return waitgroup.NewGenericResult[[]iapiserver.ObjectTypeMeta](nil, err)
			}
			var objects []iapiserver.ObjectTypeMeta
			for _, resource := range resourceList.Items {
				objects = append(objects, iapiserver.ObjectTypeMeta{ObjectMeta: resource.ObjectMeta, TypeMeta: resource.TypeMeta})
			}
			return waitgroup.NewGenericResult[[]iapiserver.ObjectTypeMeta](objects, nil)
		}))
		wg.Start(waitgroup.NewGenericFunc("", func(context.Context) waitgroup.GenericResult[[]iapiserver.ObjectTypeMeta] {
			resourceList, err := clientset.PodList(ctx, cluster, req.Namespace, metav1.ListOptions{})
			if err != nil {
				return waitgroup.NewGenericResult[[]iapiserver.ObjectTypeMeta](nil, err)
			}
			var objects []iapiserver.ObjectTypeMeta
			for _, resource := range resourceList.Items {
				objects = append(objects, iapiserver.ObjectTypeMeta{ObjectMeta: resource.ObjectMeta, TypeMeta: resource.TypeMeta})
			}
			return waitgroup.NewGenericResult[[]iapiserver.ObjectTypeMeta](objects, nil)
		}))
	}
	wg.Wait()

	for _, ret := range wg.GetResults() {
		for i, obj := range ret.Data {
			resourceInfo[string(obj.ObjectMeta.UID)] = &ret.Data[i]
			for _, parent := range obj.ObjectMeta.OwnerReferences {
				if !kinds[parent.Kind] {
					continue
				}
				parents, ok := childParent[string(obj.ObjectMeta.UID)]
				if !ok {
					parents = map[string]struct{}{}
					childParent[string(obj.ObjectMeta.UID)] = parents
				}
				parents[string(parent.UID)] = struct{}{}
			}
		}

	}

	return childParent, resourceInfo, nil
}

func (k *kubernetesService) GetComponentPod(ctx context.Context, req *iapiserver.PodListRequest) (*iapiserver.PodListResponse, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	switch req.Component {
	case iapiserver.ComponentKubectl:
		resp := &iapiserver.PodListResponse{}
		resList, err := clientset.PodList(ctx, cluster, req.Namespace, metav1.ListOptions{LabelSelector: "app=kubectl"})
		if err != nil {
			return nil, errors.WithStack(err)
		}
		var resInfos []*iapiserver.PodInfo
		for i := range resList.Items {
			podInfo := convertK8sPodToApi(&resList.Items[i], cluster, nil)
			if NewObjectCommonFieldFilter(podInfo.Resource).AddField(podInfo.Resource.Status.HostIP).AddField(podInfo.PodStatus.Status).AddField(podInfo.Resource.Status.PodIP).
				Filter(req.Fuzzy) {
				continue
			}
			resInfos = append(resInfos, podInfo)
		}
		resp.List = append(resp.List, resInfos...)
		resp.TotalCount = len(resp.List)
		return resp, nil
	default:
		return nil, errors.Errorf("invalid component:%v", req.Component)
	}
}

func (k *kubernetesService) PodPatch(ctx context.Context, req *iapiserver.PodRequest) (*iapiserver.PodInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	meta, err := libkubernetes.PodPatch(ctx, cluster.Config, req.Resource.Namespace, req.Resource.Name, req.Pt, req.Data, req.PatchOpts, req.SubResources...)
	if err != nil {
		return nil, err
	}

	return convertK8sPodToApi(meta, cluster, nil), nil
}

func (k *kubernetesService) PodEvict(ctx context.Context, req *iapiserver.PodRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return err
	}

	if err := clientset.PodEvict(ctx, cluster, req.Resource.Namespace, req.Resource, req.DeleteOpts); err != nil {
		return err
	}

	return nil
}

func (k *kubernetesService) ServiceUpdate(ctx context.Context, req *iapiserver.ServiceRequest) (*iapiserver.ServiceInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	meta, err := clientset.ServiceUpdate(ctx, cluster, req.Resource.Namespace, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, err
	}

	return iapiserver.NewServiceInfo(meta, cluster), nil
}

func (k *kubernetesService) ServiceCreate(ctx context.Context, req *iapiserver.ServiceRequest) (*iapiserver.ServiceInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	meta, err := clientset.ServiceCreate(ctx, cluster, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, err
	}

	return iapiserver.NewServiceInfo(meta, cluster), nil
}

func (k *kubernetesService) ServiceDelete(ctx context.Context, req *iapiserver.ServiceRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return err
	}

	if err := clientset.ServiceDelete(ctx, cluster, req.Resource.Namespace, req.Resource.Name, req.DeleteOpts); err != nil {
		return err
	}

	return nil
}

type NamespaceObjectDelete func(cluster *iapiserver.Cluster, namespace, name string, options metav1.DeleteOptions) error

func (k *kubernetesService) NamespaceResourceDelete(ctx context.Context, clusterUUID string, obj metav1.Object, opt metav1.DeleteOptions, delFunc NamespaceObjectDelete) (*v1.Service, error) {
	if err := libkubernetes.ValidateNamespacedObjectParameters(obj); err != nil {
		return nil, err
	}

	cluster, err := k.store.Kubernetes().Get(ctx, clusterUUID)
	if err != nil {
		return nil, err
	}

	if err := delFunc(cluster, obj.GetNamespace(), obj.GetName(), opt); err != nil {
		return nil, err
	}

	return nil, err
}

func (k *kubernetesService) ServiceBatchDelete(ctx context.Context, req *iapiserver.ServiceBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.ServiceRequest] {
	wg := waitgroup.RunGenericConcurrently[*iapiserver.ServiceRequest, *iapiserver.ServiceRequest](ctx, req.Resources, func(ctx context.Context, res *iapiserver.ServiceRequest) waitgroup.GenericResult[*iapiserver.ServiceRequest] {
		cluster, err := k.store.Clusters().Get(ctx, res.Cluster)
		if err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		if err := clientset.ServiceDelete(ctx, cluster, res.Resource.Namespace, res.Resource.Name, res.DeleteOpts); err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		return waitgroup.NewGenericResult(res, nil)
	})
	return wg.BatchGenericOutput()
}

func (k *kubernetesService) ServiceGet(ctx context.Context, req *iapiserver.ServiceGetRequest) (*iapiserver.ServiceInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	meta, err := clientset.ServiceGet(ctx, cluster, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, err
	}
	return iapiserver.NewServiceInfo(meta, cluster), nil
}

func (k *kubernetesService) ServiceList(ctx context.Context, req *iapiserver.ServiceListRequest) (*iapiserver.ServiceListResponse, error) {
	resp := &iapiserver.ServiceListResponse{}
	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList[*iapiserver.ServiceInfo](ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.ServiceInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.ServiceInfo](cluster.ID, cluster.Name)
			resList, err := clientset.ServiceList(ctx, cluster, req.Namespace, req.ToListOpts())
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.ServiceInfo
			for i := range resList.Items {
				resInfo := iapiserver.NewServiceInfo(&resList.Items[i], cluster)
				if NewObjectCommonFieldFilter(resInfo.Resource).Filter(req.Fuzzy) {
					continue
				}
				resInfos = append(resInfos, resInfo)
			}
			clusterListOne.TotalCount = len(resInfos)
			clusterListOne.List = resInfos
			return waitgroup.NewGenericResult(clusterListOne, nil)
		}, func(i, j int) bool {
			switch req.SortBy {
			case iapiserver.KubernetesResourceServiceSortByType:
				if string(resp.List[i].Resource.Spec.Type) != string(resp.List[j].Resource.Spec.Type) {
					return compareutil.Compare(string(resp.List[i].Resource.Spec.Type), string(resp.List[j].Resource.Spec.Type), req.SortDesc)
				}
			}
			return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
		}, 10*time.Second)
	return resp, err
}

func (k *kubernetesService) SecretCreate(ctx context.Context, req *iapiserver.SecretRequest) (*iapiserver.SecretInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.SecretCreate(ctx, cluster, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return iapiserver.NewSecretInfo(meta, cluster), nil
}

func (k *kubernetesService) SecretUpdate(ctx context.Context, req *iapiserver.SecretRequest) (*iapiserver.SecretInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.SecretUpdate(ctx, cluster, req.Resource.Namespace, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return iapiserver.NewSecretInfo(meta, cluster), nil
}

func (k *kubernetesService) SecretDelete(ctx context.Context, req *iapiserver.SecretRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := clientset.SecretDelete(ctx, cluster, req.Resource.Namespace, req.Resource, req.DeleteOpts); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) SecretBatchDelete(ctx context.Context, req *iapiserver.SecretBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.SecretRequest] {
	wg := waitgroup.RunGenericConcurrently[*iapiserver.SecretRequest, *iapiserver.SecretRequest](ctx, req.Resources, func(ctx context.Context, res *iapiserver.SecretRequest) waitgroup.GenericResult[*iapiserver.SecretRequest] {
		cluster, err := k.store.Clusters().Get(ctx, res.Cluster)
		if err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		if err := clientset.SecretDelete(ctx, cluster, res.Resource.Namespace, res.Resource, res.DeleteOpts); err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		return waitgroup.NewGenericResult(res, nil)
	})
	return wg.BatchGenericOutput()
}

func (k *kubernetesService) SecretGet(ctx context.Context, req *iapiserver.SecretGetRequest) (*iapiserver.SecretInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	meta, err := clientset.SecretGet(ctx, cluster, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, err
	}
	return iapiserver.NewSecretInfo(meta, cluster), nil
}

func (k *kubernetesService) SecretListAll(ctx context.Context, req *iapiserver.SecretListRequest) (*iapiserver.SecretListResponse, error) {
	resp := &iapiserver.SecretListResponse{}
	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList(ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.SecretInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.SecretInfo](cluster.ID, cluster.Name)
			resList, err := clientset.SecretList(ctx, cluster, req.Namespace, req.ToListOpts())
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.SecretInfo
			for i := range resList.Items {
				if NewObjectCommonFieldFilter(&resList.Items[i]).Filter(req.Fuzzy) {
					continue
				}
				resInfos = append(resInfos, iapiserver.NewSecretInfo(&resList.Items[i], cluster))
			}
			clusterListOne.TotalCount = len(resInfos)
			clusterListOne.List = resInfos
			return waitgroup.NewGenericResult(clusterListOne, nil)
		}, func(i, j int) bool {
			switch req.SortBy {
			case iapiserver.KubernetesResourceSecretSortByType:
				if string(resp.List[i].Resource.Type) != string(resp.List[j].Resource.Type) {
					return compareutil.Compare(string(resp.List[i].Resource.Type), string(resp.List[j].Resource.Type), req.SortDesc)
				}
			}
			return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
		}, 10*time.Second)
	return resp, err
}

func (k *kubernetesService) EndpointsAdd(ctx context.Context, req *iapiserver.EndpointsRequest) (*v1.Endpoints, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	meta, err := clientset.EndpointsCreate(ctx, cluster, req.Endpoints.Namespace, req.Endpoints, req.CreateOpts)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

func (k *kubernetesService) EndpointsDelete(ctx context.Context, req *iapiserver.EndpointsRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return err
	}

	if err = clientset.EndpointsDelete(ctx, cluster, req.Endpoints.Namespace, req.Endpoints, req.DeleteOpts); err != nil {
		return err
	}

	return nil
}

func (k *kubernetesService) EndpointsGet(ctx context.Context, req *iapiserver.EndpointsGetRequest) (*v1.Endpoints, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}
	return clientset.EndpointsGet(ctx, cluster, req.Namespace, req.Name, req.ToGetOpts())
}

func (k *kubernetesService) EndpointsListAll(ctx context.Context, req *iapiserver.EndpointListRequest) (*iapiserver.EndpointsListResponse, error) {
	resp := &iapiserver.EndpointsListResponse{}
	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList(ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.EndpointsInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.EndpointsInfo](cluster.ID, cluster.Name)
			resList, err := clientset.EndpointsList(ctx, cluster, req.Namespace, req.ToListOpts())
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.EndpointsInfo
			for i := range resList.Items {
				resInfo := iapiserver.NewEndpointsInfo(&resList.Items[i], cluster)
				if NewObjectCommonFieldFilter(resInfo.Resource).Filter(req.Fuzzy) {
					continue
				}
				resInfos = append(resInfos, resInfo)
			}
			clusterListOne.TotalCount = len(resInfos)
			clusterListOne.List = resInfos
			return waitgroup.NewGenericResult(clusterListOne, nil)
		}, func(i, j int) bool {
			return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
		}, 10*time.Second)
	return resp, err
}

func (k *kubernetesService) ResourceQuotaCreate(ctx context.Context, req *iapiserver.ResourceQuotaRequest) (*iapiserver.ResourceQuotaInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	// 限制系统命名空间创建resource quota
	if sets.NewString(iapiserver.KubernetesSystemNamespaces...).Has(req.Resource.Namespace) {
		return nil, errors.Errorf("system namespace[%v] limit create resource quota", iapiserver.KubernetesSystemNamespaces)
	}

	if req.UpdateIfExists {
		//resource quota has exist
		if meta, _ := clientset.ResourceQuotaGet(ctx, cluster, req.Resource.Namespace, req.Resource.Name, req.GetOpts); meta != nil {
			meta, err = clientset.ResourceQuotaUpdate(ctx, cluster, req.Resource.Namespace, req.Resource, req.UpdateOpts)
			if err != nil {
				return nil, err
			}

			return convertK8sResourceQuotaToApi(meta, cluster), nil
		}
	}
	meta, err := clientset.ResourceQuotaCreate(ctx, cluster, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, err
	}

	return convertK8sResourceQuotaToApi(meta, cluster), nil
}

func (k *kubernetesService) ResourceQuotaDelete(ctx context.Context, req *iapiserver.ResourceQuotaRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return err
	}

	if err = clientset.ResourceQuotaDelete(ctx, cluster, req.Resource.Namespace, req.Resource, req.DeleteOpts); err != nil {
		return err
	}

	return nil
}

func (k *kubernetesService) ResourceQuotaUpdate(ctx context.Context, req *iapiserver.ResourceQuotaRequest) (*iapiserver.ResourceQuotaInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	meta, err := clientset.ResourceQuotaUpdate(ctx, cluster, req.Resource.Namespace, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, err
	}

	return convertK8sResourceQuotaToApi(meta, cluster), nil
}

func (k *kubernetesService) ResourceQuotaGet(ctx context.Context, req *iapiserver.ResourceQuotaGetRequest) (*iapiserver.ResourceQuotaInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	meta, err := clientset.ResourceQuotaGet(ctx, cluster, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, err
	}

	return convertK8sResourceQuotaToApi(meta, cluster), nil
}

func (k *kubernetesService) ResourceQuotaListAll(ctx context.Context, req *iapiserver.ResourceQuotaListRequest) (*iapiserver.ResourceQuotaListResponse, error) {
	resp := &iapiserver.ResourceQuotaListResponse{}

	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList[*iapiserver.ResourceQuotaInfo](ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.ResourceQuotaInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.ResourceQuotaInfo](cluster.ID, cluster.Name)
			resList, err := clientset.ResourceQuotaList(ctx, cluster, req.Namespace, req.ToListOpts())
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.ResourceQuotaInfo
			for i := range resList.Items {
				resInfo := convertK8sResourceQuotaToApi(&resList.Items[i], cluster)
				if NewObjectCommonFieldFilter(resInfo.Resource).Filter(req.Fuzzy) {
					continue
				}
				resInfos = append(resInfos, resInfo)
			}
			clusterListOne.TotalCount = len(resInfos)
			clusterListOne.List = resInfos
			return waitgroup.NewGenericResult(clusterListOne, nil)
		}, func(i, j int) bool {
			return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
		}, 10*time.Second)
	return resp, err
}

// k8s会把限制转成human-readable,  如 100000--> 100k , 10000000 ---> 1m
func convertK8sResourceQuotaToApi(meta *v1.ResourceQuota, cluster *iapiserver.Cluster) *iapiserver.ResourceQuotaInfo {
	ri := iapiserver.NewResourceQuotaInfo(meta, cluster)

	// convert to int64 for frontend

	ri.Convert = &iapiserver.ResourceQuotaConvert{}
	covertFunc := func(origin v1.ResourceList, convert map[string]int64) {
		for k, v := range origin {
			if string(k) == "requests.memory" {
				convert[string(k)] = v.Value() / 1024 / 1024 // convert to memory
				continue
			}
			if string(k) == "requests.cpu" { // if use Value(), 0.1 cpu/100m cpu will convert to 1 cpu
				convert[string(k)] = v.MilliValue()
				continue
			}
			convert[string(k)] = v.Value()
		}
	}

	if meta.Spec.Hard != nil {
		ri.Convert.Spec.Hard = make(map[string]int64)
		covertFunc(meta.Spec.Hard, ri.Convert.Spec.Hard)
	}
	if meta.Status.Hard != nil {
		ri.Convert.Status.Hard = make(map[string]int64)
		covertFunc(meta.Status.Hard, ri.Convert.Status.Hard)
	}
	if meta.Status.Used != nil {
		ri.Convert.Status.Used = make(map[string]int64)
		covertFunc(meta.Status.Used, ri.Convert.Status.Used)
	}

	return ri
}

func (k *kubernetesService) EventGet(ctx context.Context, req *iapiserver.EventGetRequest) (*iapiserver.EventInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.EventGet(ctx, cluster, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sEventTopApi(meta, cluster), nil
}

func (k *kubernetesService) EventList(ctx context.Context, req *iapiserver.EventListRequest) (*iapiserver.EventListResponse, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resList, err := clientset.EventList(ctx, cluster, req.Namespace, req.ToListOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var resInfos []*iapiserver.EventInfo
	for i := range resList.Items {
		resInfo := convertK8sEventTopApi(&resList.Items[i], cluster)
		if NewObjectCommonFieldFilter(resInfo.Resource).
			AddField(resInfo.Resource.Message).
			AddField(resInfo.Resource.Reason).
			AddField(resInfo.Resource.Type).
			Filter(req.Fuzzy) {
			continue
		}
		resInfos = append(resInfos, resInfo)
	}

	resp := &iapiserver.EventListResponse{}
	resp.TotalCount = len(resInfos)
	resp.List = resInfos
	s, e := paging.Index(len(resp.List), req.PageNum, req.PageSize)
	resp.List = resp.List[s:e]
	return resp, nil
}

func convertK8sEventTopApi(meta *v1.Event, cluster *iapiserver.Cluster) *iapiserver.EventInfo {
	annotations := meta.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	if !meta.FirstTimestamp.IsZero() {
		annotations["first_timestamp"] = fmt.Sprintf("%v", meta.FirstTimestamp.Unix())
	}
	if !meta.LastTimestamp.IsZero() {
		annotations["last_timestamp"] = fmt.Sprintf("%v", meta.LastTimestamp.Unix())
	}
	meta.SetAnnotations(annotations)
	return iapiserver.NewEventInfo(meta, cluster)
}

func maxContainerRestarts(pod *v1.Pod) int {
	maxRestarts := 0
	for _, c := range pod.Status.ContainerStatuses {
		maxRestarts = mathutil.Max(maxRestarts, int(c.RestartCount))
	}

	// if pod is running init container
	for _, v := range pod.Status.Conditions {
		if v.Type == v1.PodInitialized && v.Status != v1.ConditionTrue {
			for _, c := range pod.Status.InitContainerStatuses {
				maxRestarts = mathutil.Max(maxRestarts, int(c.RestartCount))
			}
			return maxRestarts
		}
	}

	return maxRestarts
}

func convertK8sPodToApi(pod *v1.Pod, cluster *iapiserver.Cluster, addressNode map[string]*v1.ObjectReference) *iapiserver.PodInfo {
	podInfo := &iapiserver.PodInfo{
		Resource: pod,
	}

	maxRestarts := maxContainerRestarts(pod)
	podStatus := calculatePodStatus(pod)
	request, limit := resource.PodRequestsAndLimits(pod)

	podInfo.PodStatus = &podStatus
	podInfo.MaxRestarts = &maxRestarts
	podInfo.ResourceRequest = request
	podInfo.ResourceLimit = limit

	if addressNode != nil {
		if node, ok := addressNode[podInfo.Resource.Status.HostIP]; ok {
			podInfo.NodeInfo = node
		}
	}

	podInfo.ResourceConvert = make([]*iapiserver.ResourceConvert, 0)
	for _, container := range pod.Spec.Containers {
		podInfo.ResourceConvert = append(podInfo.ResourceConvert, convertResourceLimitToPersistentUnit(container.Name, container.Resources.Requests, container.Resources.Limits))
	}
	podInfo.Cluster = cluster

	return podInfo
}

func calculatePodStatus(pod *v1.Pod) iapiserver.PodStatus {
	st := pod.Status
	var readyContainer, totalContainer int
	totalContainer = len(pod.Spec.Containers)
	for _, c := range st.ContainerStatuses {
		if c.Ready {
			readyContainer++
		}
	}

	switch st.Phase {
	case v1.PodPending:
		if noReadyConditions := getNotSuccessPodConditions(st.Conditions); len(noReadyConditions) != 0 {
			for _, v := range noReadyConditions {
				if v.Type == v1.PodScheduled && v.Status != v1.ConditionTrue {
					return newPodStatus(iapiserver.PodStatusScheduleFail, v.Reason, v.Message, readyContainer, totalContainer)
				}

				if v.Type == v1.PodInitialized && v.Status != v1.ConditionTrue {
					for i, ct := range st.InitContainerStatuses {
						if !ct.Ready {
							if ct.State.Running != nil {
								return newPodStatus(iapiserver.PodInitContainerFailStatusPrefix+fmt.Sprintf("%v/%v", i, len(st.InitContainerStatuses)), "", "", readyContainer, totalContainer)
							}

							if ct.State.Waiting != nil {
								return newPodStatus(iapiserver.PodInitContainerFailStatusPrefix+ct.State.Waiting.Reason, ct.State.Waiting.Reason, ct.State.Waiting.Message, readyContainer, totalContainer)
							}
						}
					}
				}

				if v.Type == v1.ContainersReady && v.Status != v1.ConditionTrue {
					for _, ct := range st.ContainerStatuses {
						if !ct.Ready {
							if ct.State.Waiting != nil {
								return newPodStatus(ct.State.Waiting.Reason, ct.State.Waiting.Reason, ct.State.Waiting.Message, readyContainer, totalContainer)
							}
							if ct.State.Terminated != nil {
								return newPodStatus(ct.State.Terminated.Reason, ct.State.Terminated.Reason, ct.State.Terminated.Message, readyContainer, totalContainer)
							}
						}
					}
				}
				if v.Type == v1.PodReady && v.Status != v1.ConditionTrue {
					//FIXME:
				}
			}
		}

	case v1.PodRunning:
		if pod.DeletionTimestamp != nil {
			return newPodStatus(iapiserver.PodStatusTerminating, "", "DeleteTimeStamp:"+pod.DeletionTimestamp.String(), readyContainer, totalContainer)
		}

		if noReadyConditions := getNotSuccessPodConditions(st.Conditions); len(noReadyConditions) != 0 {
			for _, v := range noReadyConditions {
				if v.Type == v1.ContainersReady && v.Status != v1.ConditionTrue {
					for _, ct := range st.ContainerStatuses {
						if !ct.Ready {
							if ct.State.Waiting != nil {
								return newPodStatus(ct.State.Waiting.Reason, ct.State.Waiting.Reason, ct.State.Waiting.Message, readyContainer, totalContainer)
							}

							if ct.State.Terminated != nil && ct.State.Terminated.ExitCode != 0 {
								return newPodStatus(ct.State.Terminated.Reason, ct.State.Terminated.Reason, ct.State.Terminated.Message, readyContainer, totalContainer)
							}
						}
					}
				}
				if v.Type == v1.PodReady && v.Status != v1.ConditionTrue {

				}

			}
		}
	case v1.PodFailed:
		if pod.Status.Reason == "UnexpectedAdmissionError" {
			return newPodStatus("UnexpectedAdmissionError", pod.Status.Message, "", readyContainer, totalContainer)
		}
	}
	//rollback use phase as pod status.
	return newPodStatus(string(st.Phase), st.Reason, st.Message, readyContainer, totalContainer)
}

func newPodStatus(podState, reason, message string, readyContainer, totalContainer int) iapiserver.PodStatus {
	if podState == reason {
		if message == "" {
			reason = ""
		}
	}

	return iapiserver.PodStatus{
		Status:          podState,
		Message:         message,
		Reason:          reason,
		ReadyContainers: readyContainer,
		TotalContainers: totalContainer,
	}
}

func getNotSuccessPodConditions(conditions []v1.PodCondition) []v1.PodCondition {
	notReadyConditions := make([]v1.PodCondition, 0)
	for _, v := range conditions {
		if v.Status != v1.ConditionTrue {
			notReadyConditions = append(notReadyConditions, v)
		}
	}
	return notReadyConditions
}

func (k *kubernetesService) LimitRangeCreate(ctx context.Context, req *iapiserver.LimitRangeRequest) (*iapiserver.LimitRangeInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	if req.UpdateIfExists {
		if meta, _ := clientset.LimitRangeGet(ctx, cluster, req.Resource.Namespace, req.Resource.Name, req.GetOpts); meta != nil {
			meta, err := clientset.LimitRangeUpdate(ctx, cluster, req.Resource.Namespace, req.Resource, req.UpdateOpts)
			if err != nil {
				return nil, err
			}

			return convertK8sLimitRangeToApi(meta, cluster, req.Yaml), nil
		}
	}
	meta, err := clientset.LimitRangeCreate(ctx, cluster, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, err
	}

	return convertK8sLimitRangeToApi(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) LimitRangeDelete(ctx context.Context, req *iapiserver.LimitRangeRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	err = clientset.LimitRangeDelete(ctx, cluster, req.Resource.Namespace, req.Resource, req.DeleteOpts)
	if err != nil {
		return err
	}

	return nil
}

func (k *kubernetesService) LimitRangeUpdate(ctx context.Context, req *iapiserver.LimitRangeRequest) (*iapiserver.LimitRangeInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	meta, err := clientset.LimitRangeUpdate(ctx, cluster, req.Resource.Namespace, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, err
	}

	return convertK8sLimitRangeToApi(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) LimitRangeGet(ctx context.Context, req *iapiserver.LimitRangeGetRequest) (*iapiserver.LimitRangeInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.LimitRangeGet(ctx, cluster, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sLimitRangeToApi(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) LimitRangeList(ctx context.Context, req *iapiserver.LimitRangeListRequest) (*iapiserver.LimitRangeListResponse, error) {
	resp := &iapiserver.LimitRangeListResponse{}
	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList[*iapiserver.LimitRangeInfo](ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.LimitRangeInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.LimitRangeInfo](cluster.ID, cluster.Name)
			resList, err := clientset.LimitRangeList(ctx, cluster, req.Namespace, req.ToListOpts())
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.LimitRangeInfo
			for i := range resList.Items {
				resInfo := iapiserver.NewLimitRangeInfo(&resList.Items[i], cluster)
				if NewObjectCommonFieldFilter(resInfo.Resource).Filter(req.Fuzzy) {
					continue
				}
				resInfos = append(resInfos, resInfo)
			}
			clusterListOne.TotalCount = len(resInfos)
			clusterListOne.List = resInfos
			return waitgroup.NewGenericResult(clusterListOne, nil)
		}, func(i, j int) bool {
			return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
		}, 10*time.Second)
	return resp, err
}

func convertK8sLimitRangeToApi(meta *v1.LimitRange, cluster *iapiserver.Cluster, yaml bool) *iapiserver.LimitRangeInfo {
	lri := iapiserver.NewLimitRangeInfo(meta, cluster)
	if !yaml {
		lri.Convert = &iapiserver.LimitRangeLimitConvert{}
		covertFunc := func(origin v1.ResourceList, convert map[string]int64) {
			for k, v := range origin {
				if string(k) == "memory" {
					convert[string(k)] = v.Value() / 1024 / 1024 // convert to memory to M unit. frontend will pass *M as create arg.
					continue
				}
				if string(k) == "cpu" {
					convert[string(k)] = v.MilliValue() // if use Value(), 0.1 core to calculate to 1 core. so use milliValue to convert, let frontend to convert.
					continue
				}
				convert[string(k)] = v.Value()
			}
		}

		for _, v := range meta.Spec.Limits {
			var limitconvert iapiserver.LimitRangeConvert
			limitconvert.Type = string(v.Type)

			if v.Max != nil {
				limitconvert.Max = make(map[string]int64)
				covertFunc(v.Max, limitconvert.Max)
			}
			if v.Min != nil {
				limitconvert.Min = make(map[string]int64)
				covertFunc(v.Min, limitconvert.Min)
			}
			if v.Default != nil {
				limitconvert.Default = make(map[string]int64)
				covertFunc(v.Default, limitconvert.Default)
			}
			if v.DefaultRequest != nil {
				limitconvert.DefaultRequest = make(map[string]int64)
				covertFunc(v.DefaultRequest, limitconvert.DefaultRequest)
			}
			if v.MaxLimitRequestRatio != nil {
				limitconvert.MaxLimitRequestRatio = make(map[string]int64)
				covertFunc(v.MaxLimitRequestRatio, limitconvert.MaxLimitRequestRatio)
			}
			lri.Convert.Limits = append(lri.Convert.Limits, limitconvert)
		}
	}
	return lri
}
