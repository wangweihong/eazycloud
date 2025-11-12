package kubernetes

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/internal/pkg/clientset"
	"github.com/wangweihong/eazycloud/internal/pkg/libkubernetes"
	"github.com/wangweihong/gotoolbox/pkg/compareutil"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/log"
	"github.com/wangweihong/gotoolbox/pkg/paging"
	"github.com/wangweihong/gotoolbox/pkg/typeutil"
	"github.com/wangweihong/gotoolbox/pkg/waitgroup"
	appsv1 "k8s.io/api/apps/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/meta"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubetypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	deploy "k8s.io/kubectl/pkg/util/deployment"

	corev1 "k8s.io/api/core/v1"
)

func (k *kubernetesService) DaemonSetVersionUpdate(ctx context.Context, req *iapiserver.DaemonSetRequest) (*iapiserver.DaemonSetInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	daemonSet, err := clientset.DaemonSetGet(ctx, cluster, req.Resource.Namespace, req.Resource.Name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	selector, err := metav1.LabelSelectorAsSelector(daemonSet.Spec.Selector)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	accessor, err := meta.Accessor(daemonSet)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	listOPt := metav1.ListOptions{LabelSelector: selector.String()}
	historyList, err := clientset.RevisionList(ctx, cluster, req.Resource.Namespace, listOPt)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	isFind := false
	history := appsv1.ControllerRevision{}
	for i := range historyList.Items {
		history = historyList.Items[i]
		if metav1.IsControlledBy(&history, accessor) {
			if strconv.Itoa(int(history.Revision)) == req.Version {
				isFind = true
				break
			}
		}
	}

	if !isFind {
		return nil, errors.Errorf("version[%v] not exist", req.Version)
	}

	daemonSet, err = clientset.DaemonSetPatch(ctx, cluster, req.Resource.Namespace, req.Resource, kubetypes.StrategicMergePatchType, history.Data.Raw, metav1.PatchOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sDaemonSetToApiDaemonSet(daemonSet, cluster), nil
}

// 获取版本列表
func (k *kubernetesService) DaemonSetVersionList(ctx context.Context, req *iapiserver.DaemonSetVersionListRequest) (*iapiserver.DaemonSetVersionListResponse, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sts, err := clientset.DaemonSetGet(ctx, cluster, req.Namespace, req.Name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	selector, err := metav1.LabelSelectorAsSelector(sts.Spec.Selector)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	accessor, err := meta.Accessor(sts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	listOPt := metav1.ListOptions{LabelSelector: selector.String()}
	historyList, err := clientset.RevisionList(ctx, cluster, req.Namespace, listOPt)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp := &iapiserver.DaemonSetVersionListResponse{}
	resp.CurrentVersion = &iapiserver.VersionInfo{Namespace: req.Namespace, Name: req.Name}
	if len(historyList.Items) != 0 {
		sort.Slice(historyList.Items, func(i, j int) bool {
			return historyList.Items[i].Revision > historyList.Items[j].Revision
		})
		resp.CurrentVersion.Version = strconv.Itoa(int(historyList.Items[0].Revision))
	}
	for i := range historyList.Items {
		history := historyList.Items[i]
		if metav1.IsControlledBy(&history, accessor) {
			tmpDae, err := k.versionToDaemonSet(sts, &history)
			if err != nil {
				log.Errorf("%v", err)
				continue
			}
			versionInfo := &iapiserver.VersionInfo{
				Version:    strconv.Itoa(int(history.Revision)),
				CreateTime: history.CreationTimestamp.Unix(),
				DaemonSet:  tmpDae,
			}

			resp.VersionList = append(resp.VersionList, *versionInfo)
		}
	}
	sort.SliceStable(resp.VersionList, func(i, j int) bool {
		return resp.VersionList[j].CreateTime > resp.VersionList[i].CreateTime
	})

	resp.TotalCount = len(resp.VersionList)
	s, e := paging.Index(resp.TotalCount, req.PageNum, req.PageSize)
	resp.VersionList = resp.VersionList[s:e]
	return resp, nil
}

func (k *kubernetesService) versionToDaemonSet(dae *appsv1.DaemonSet, version *appsv1.ControllerRevision) (*appsv1.DaemonSet, error) {
	resp := &appsv1.DaemonSet{}
	if dae == nil || version == nil {
		return resp, nil
	}

	resp.ObjectMeta = version.ObjectMeta
	if err := json.Unmarshal(version.Data.Raw, &resp); err != nil {
		return nil, errors.Errorf("decode error:%v", err)
	}

	return resp, nil
}

func (k *kubernetesService) DaemonSetAdd(ctx context.Context, req *iapiserver.DaemonSetRequest) (*iapiserver.DaemonSetInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.DaemonSetCreate(ctx, cluster, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sDaemonSetToApiDaemonSet(meta, cluster), nil
}

func (k *kubernetesService) DaemonSetRecreate(ctx context.Context, req *iapiserver.DaemonSetRequest) (*iapiserver.DaemonSetInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.DaemonSetGet(ctx, cluster, req.Resource.Namespace, req.Resource.Name, req.GetOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	selector, err := metav1.LabelSelectorAsSelector(meta.Spec.Selector)
	if err != nil {
		return nil, errors.Errorf("selector '%v' parse fail:%v", meta.Spec.Selector, err)
	}

	if err := clientset.PodDeleteCollection(ctx, cluster, req.Resource.Namespace, metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: selector.String()}); err != nil {
		return nil, errors.WithStack(err)
	}
	return convertK8sDaemonSetToApiDaemonSet(meta, cluster), nil
}

func (k *kubernetesService) DaemonSetUpdate(ctx context.Context, req *iapiserver.DaemonSetRequest) (*iapiserver.DaemonSetInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.DaemonSetUpdate(ctx, cluster, req.Resource.Namespace, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return convertK8sDaemonSetToApiDaemonSet(meta, cluster), nil
}

func (k *kubernetesService) DaemonSetDelete(ctx context.Context, req *iapiserver.DaemonSetRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := clientset.DaemonSetDelete(ctx, cluster, req.Resource.Namespace, req.Resource, req.DeleteOpts); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) DaemonSetBatchDelete(ctx context.Context, req *iapiserver.DaemonSetBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.DaemonSetRequest] {
	wg := waitgroup.RunGenericConcurrently[*iapiserver.DaemonSetRequest, *iapiserver.DaemonSetRequest](ctx, req.Resources, func(ctx context.Context, res *iapiserver.DaemonSetRequest) waitgroup.GenericResult[*iapiserver.DaemonSetRequest] {
		cluster, err := k.store.Clusters().Get(ctx, res.Cluster)
		if err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		if err := clientset.DaemonSetDelete(ctx, cluster, res.Resource.Namespace, res.Resource, res.DeleteOpts); err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		return waitgroup.NewGenericResult(res, nil)
	})
	return wg.BatchGenericOutput()
}

func (k *kubernetesService) DaemonSetGet(ctx context.Context, req *iapiserver.DaemonSetGetRequest) (*iapiserver.DaemonSetInfo, error) {
	if err := libkubernetes.ValidateNamespacedScopeParameters(req.Namespace, req.Name); err != nil {
		return nil, errors.WithStack(err)
	}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.DaemonSetGet(ctx, cluster, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sDaemonSetToApiDaemonSet(meta, cluster), nil
}

func (k *kubernetesService) DaemonSetList(ctx context.Context, req *iapiserver.DaemonSetListRequest) (*iapiserver.DaemonSetListResponse, error) {
	resp := &iapiserver.DaemonSetListResponse{}
	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList(ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, c *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.DaemonSetInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.DaemonSetInfo](c.ID, c.Name)
			resList, err := clientset.DaemonSetList(ctx, c, req.Namespace, req.ToListOpts())
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.DaemonSetInfo
			for i := range resList.Items {
				resInfo := convertK8sDaemonSetToApiDaemonSet(&resList.Items[i], c)
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

func convertK8sDaemonSetToApiDaemonSet(meta *appsv1.DaemonSet, cluster *iapiserver.Cluster) *iapiserver.DaemonSetInfo {
	resp := &iapiserver.DaemonSetInfo{
		ResourceInfo: iapiserver.NewResourceInfo(meta, cluster),
	}

	resp.ResourceConvert = make([]*iapiserver.ResourceConvert, 0, len(meta.Spec.Template.Spec.Containers))
	for _, container := range meta.Spec.Template.Spec.Containers {
		resp.ResourceConvert = append(resp.ResourceConvert, convertResourceLimitToPersistentUnit(container.Name, container.Resources.Requests, container.Resources.Limits))
	}

	return resp
}

func (k *kubernetesService) ReplicaSetAdd(ctx context.Context, req *iapiserver.ReplicaSetRequest) (*iapiserver.ReplicaSetInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.ReplicaSetCreate(ctx, cluster, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewReplicaSetInfo(meta, cluster), nil
}

func (k *kubernetesService) ReplicaSetDelete(ctx context.Context, req *iapiserver.ReplicaSetRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	err = clientset.ReplicaSetDelete(ctx, cluster, req.Resource.Namespace, req.Resource, req.DeleteOpts)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) ReplicaSetUpdate(ctx context.Context, req *iapiserver.ReplicaSetRequest) (*iapiserver.ReplicaSetInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.ReplicaSetUpdate(ctx, cluster, req.Resource.Namespace, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewReplicaSetInfo(meta, cluster), nil
}

func (k *kubernetesService) ReplicaSetGet(ctx context.Context, req *iapiserver.ReplicaSetGetRequest) (*iapiserver.ReplicaSetInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.ReplicaSetGet(ctx, cluster, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return iapiserver.NewReplicaSetInfo(meta, cluster), nil
}

func (k *kubernetesService) ReplicaSetList(ctx context.Context, req *iapiserver.ReplicaSetListRequest) (*iapiserver.ReplicaSetListResponse, error) {
	resp := &iapiserver.ReplicaSetListResponse{}

	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList(ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, c *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.ReplicaSetInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.ReplicaSetInfo](c.ID, c.Name)
			resList, err := clientset.ReplicaSetList(ctx, c, req.Namespace, req.ToListOpts())
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.ReplicaSetInfo
			for i := range resList.Items {
				resInfo := iapiserver.NewReplicaSetInfo(&resList.Items[i], c)
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

func (k *kubernetesService) HpaAdd(ctx context.Context, req *iapiserver.HpaRequest) (*iapiserver.HpaInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.HpaCreate(ctx, cluster, req.HorizontalPodAutoscaler.Namespace, req.HorizontalPodAutoscaler, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return iapiserver.NewHpaInfo(meta, cluster), nil
}

func (k *kubernetesService) HpaDelete(ctx context.Context, req *iapiserver.HpaRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = clientset.HpaDelete(ctx, cluster, req.HorizontalPodAutoscaler.Namespace, req.HorizontalPodAutoscaler, req.DeleteOpts); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) HpaUpdate(ctx context.Context, req *iapiserver.HpaRequest) (*iapiserver.HpaInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.HpaUpdate(ctx, cluster, req.HorizontalPodAutoscaler.Namespace, req.HorizontalPodAutoscaler, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewHpaInfo(meta, cluster), nil
}

func (k *kubernetesService) HpaGet(ctx context.Context, req *iapiserver.HpaGetRequest) (*iapiserver.HpaInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.HpaGet(ctx, cluster, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewHpaInfo(meta, cluster), nil
}

func (k *kubernetesService) HpaList(ctx context.Context, req *iapiserver.HpaListRequest) (*iapiserver.HpaListResponse, error) {
	resp := &iapiserver.HpaListResponse{}

	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList(ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, c *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.HpaInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.HpaInfo](c.ID, c.Name)
			resList, err := clientset.HpaList(ctx, c, req.Namespace, req.ToListOpts())
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.HpaInfo
			for i := range resList.Items {
				resInfo := iapiserver.NewHpaInfo(&resList.Items[i], c)
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

func (k *kubernetesService) StatefulSetAdd(ctx context.Context, req *iapiserver.StatefulSetRequest) (*iapiserver.StatefulSetInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.StatefulSetCreate(ctx, cluster, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sStatefulSetToApi(meta, cluster), nil
}

func (k *kubernetesService) StatefulSetVersionUpdate(ctx context.Context, req *iapiserver.StatefulSetRequest) (*iapiserver.StatefulSetInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sts, err := clientset.StatefulSetGet(ctx, cluster, req.Resource.Namespace, req.Resource.Name, metav1.GetOptions{})
	if err != nil {

		return nil, errors.WithStack(err)
	}

	selector, err := metav1.LabelSelectorAsSelector(sts.Spec.Selector)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	accessor, err := apimeta.Accessor(sts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	listOPt := metav1.ListOptions{LabelSelector: selector.String()}
	historyList, err := clientset.RevisionList(ctx, cluster, req.Resource.Namespace, listOPt)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	isFind := false

	history := appsv1.ControllerRevision{}
	for i := range historyList.Items {
		history = historyList.Items[i]
		if metav1.IsControlledBy(&history, accessor) {
			if strconv.Itoa(int(history.Revision)) == req.Version {
				isFind = true
				break
			}
		}
	}

	if !isFind {
		return nil, errors.Errorf("version[%v] not exist", req.Version)
	}

	sts, err = clientset.StatefulSetPatch(ctx, cluster, req.Resource.Namespace, req.Resource, kubetypes.StrategicMergePatchType, history.Data.Raw, metav1.PatchOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return convertK8sStatefulSetToApi(sts, cluster), nil
}

// 获取版本列表
func (k *kubernetesService) StatefulSetVersionList(ctx context.Context, req *iapiserver.StatefulSetVersionListRequest) (*iapiserver.StatefulSetVersionListResponse, error) {
	resp := &iapiserver.StatefulSetVersionListResponse{}
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sts, err := clientset.StatefulSetGet(ctx, cluster, req.Namespace, req.Name, v1.GetOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	selector, err := metav1.LabelSelectorAsSelector(sts.Spec.Selector)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	accessor, err := apimeta.Accessor(sts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	listOPt := metav1.ListOptions{LabelSelector: selector.String()}
	historyList, err := clientset.RevisionList(ctx, cluster, req.Namespace, listOPt)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range historyList.Items {
		history := historyList.Items[i]
		if metav1.IsControlledBy(&history, accessor) {
			tmpSts, err := versionToStatefulSet(sts, &history)
			if err != nil {
				log.Errorf("%v", err)
				continue
			}

			versionInfo := &iapiserver.VersionInfo{
				Version:     strconv.Itoa(int(history.Revision)),
				CreateTime:  history.CreationTimestamp.Unix(),
				StatefulSet: tmpSts,
			}

			resp.VersionList = append(resp.VersionList, *versionInfo)
		}
	}
	sort.SliceStable(resp.VersionList, func(i, j int) bool {
		return resp.VersionList[j].CreateTime > resp.VersionList[i].CreateTime
	})
	resp.CurrentVersion = &iapiserver.VersionInfo{Namespace: req.Namespace, Name: req.Name}
	if len(historyList.Items) != 0 {
		sort.Slice(historyList.Items, func(i, j int) bool {
			return historyList.Items[i].Revision > historyList.Items[j].Revision
		})
		resp.CurrentVersion.Version = strconv.Itoa(int(historyList.Items[0].Revision))
	}

	resp.TotalCount = len(resp.VersionList)
	s, e := paging.Index(resp.TotalCount, req.PageNum, req.PageSize)
	resp.VersionList = resp.VersionList[s:e]

	return nil, errors.WithStack(err)
}

func versionToStatefulSet(sts *appsv1.StatefulSet, version *appsv1.ControllerRevision) (*appsv1.StatefulSet, error) {
	resp := &appsv1.StatefulSet{}

	if sts == nil || version == nil {
		return resp, nil
	}

	resp.Spec = appsv1.StatefulSetSpec{}

	resp.ObjectMeta = version.ObjectMeta
	resp.TypeMeta = version.TypeMeta
	if err := json.Unmarshal(version.Data.Raw, &resp); err != nil {
		return nil, errors.WithStack(err)
	}
	return resp, nil
}

func (k *kubernetesService) versionToStatefulSet(sts *appsv1.StatefulSet, version *appsv1.ControllerRevision) (*appsv1.StatefulSet, error) {
	resp := &appsv1.StatefulSet{}
	if sts == nil || version == nil {
		return resp, nil
	}

	resp.Spec = appsv1.StatefulSetSpec{}

	resp.ObjectMeta = version.ObjectMeta
	resp.TypeMeta = version.TypeMeta
	if err := json.Unmarshal(version.Data.Raw, &resp); err != nil {
		return nil, errors.WithStack(err)
	}
	return resp, nil
}

func (k *kubernetesService) StatefulSetRecreate(ctx context.Context, req *iapiserver.StatefulSetRequest) (*iapiserver.StatefulSetInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.StatefulSetGet(ctx, cluster, req.Resource.Namespace, req.Resource.Name, req.GetOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	selector, err := metav1.LabelSelectorAsSelector(meta.Spec.Selector)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := clientset.PodDeleteCollection(ctx, cluster, req.Resource.Namespace, metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: selector.String()}); err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sStatefulSetToApi(meta, cluster), nil
}

func (k *kubernetesService) StatefulSetUpdate(ctx context.Context, req *iapiserver.StatefulSetRequest) (*iapiserver.StatefulSetInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	//avoid old resource version error
	req.Resource.ResourceVersion = ""
	meta, err := clientset.StatefulSetUpdate(ctx, cluster, req.Resource.Namespace, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return convertK8sStatefulSetToApi(meta, cluster), nil
}

func (k *kubernetesService) StatefulSetDelete(ctx context.Context, req *iapiserver.StatefulSetRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := clientset.StatefulSetDelete(ctx, cluster, req.Resource.Namespace, req.Resource, req.DeleteOpts); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) StatefulSetBatchDelete(ctx context.Context, req *iapiserver.StatefulSetBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.StatefulSetRequest] {
	wg := waitgroup.RunGenericConcurrently(ctx, req.Resources, func(ctx context.Context, res *iapiserver.StatefulSetRequest) waitgroup.GenericResult[*iapiserver.StatefulSetRequest] {
		cluster, err := k.store.Clusters().Get(ctx, res.Cluster)
		if err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		if err := clientset.StatefulSetDelete(ctx, cluster, res.Resource.Namespace, res.Resource, res.DeleteOpts); err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		return waitgroup.NewGenericResult(res, nil)
	})
	return wg.BatchGenericOutput()
}

func (k *kubernetesService) StatefulSetGet(ctx context.Context, req *iapiserver.StatefulSetGetRequest) (*iapiserver.StatefulSetInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.StatefulSetGet(ctx, cluster, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sStatefulSetToApi(meta, cluster), nil
}

func (k *kubernetesService) StatefulSetList(ctx context.Context, req *iapiserver.StatefulSetListRequest) (*iapiserver.StatefulSetListResponse, error) {
	resp := &iapiserver.StatefulSetListResponse{}
	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList[*iapiserver.StatefulSetInfo](ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, c *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.StatefulSetInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.StatefulSetInfo](c.ID, c.Name)
			resList, err := clientset.StatefulSetList(ctx, c, req.Namespace, req.ToListOpts())
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.StatefulSetInfo
			for i := range resList.Items {
				resInfo := convertK8sStatefulSetToApi(&resList.Items[i], c)
				resInfos = append(resInfos, resInfo)
			}
			clusterListOne.TotalCount = len(resInfos)
			clusterListOne.List = resInfos
			return waitgroup.NewGenericResult(clusterListOne, nil)
		}, func(i, j int) bool {
			switch req.SortBy {
			case iapiserver.KubernetesResourceStatefulSetSortByReplicas:
				var replica1, replica2 int32
				if resp.List[i].Resource.Spec.Replicas != nil {
					replica1 = *resp.List[i].Resource.Spec.Replicas
				}
				if resp.List[j].Resource.Spec.Replicas != nil {
					replica2 = *resp.List[j].Resource.Spec.Replicas
				}
				if replica1 != replica2 {
					return compareutil.Compare(replica1, replica2, req.SortDesc)
				}
			}
			return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
		}, 10*time.Second)
	return resp, err
}

func convertK8sStatefulSetToApi(meta *appsv1.StatefulSet, cluster *iapiserver.Cluster) *iapiserver.StatefulSetInfo {
	resp := iapiserver.NewStatefulSetInfo(meta, cluster)

	resp.ResourceConvert = make([]*iapiserver.ResourceConvert, 0)
	for _, container := range meta.Spec.Template.Spec.Containers {
		resp.ResourceConvert = append(resp.ResourceConvert, convertResourceLimitToPersistentUnit(container.Name, container.Resources.Requests, container.Resources.Limits))
	}

	return resp
}

var (
	deploymentVersionFlag = "deployment.kubernetes.io/revision"
)

func (k *kubernetesService) DeploymentUpdate(ctx context.Context, req *iapiserver.DeploymentRequest) (*iapiserver.DeploymentInfo, error) {
	req.Resource.ResourceVersion = ""
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.DeploymentUpdate(ctx, cluster, req.Resource.Namespace, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sDeploymentToApi(meta, cluster), nil

}

func (k *kubernetesService) DeploymentAdd(ctx context.Context, req *iapiserver.DeploymentRequest) (*iapiserver.DeploymentInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.DeploymentCreate(ctx, cluster, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sDeploymentToApi(meta, cluster), nil
}

func (k *kubernetesService) DeploymentDelete(ctx context.Context, req *iapiserver.DeploymentRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	if req.DeleteCollection {
		meta, err := clientset.DeploymentGet(ctx, cluster, req.Resource.Namespace, req.Resource.Name, req.GetOpts)
		if err != nil {
			return errors.WithStack(err)
		}

		selector, err := metav1.LabelSelectorAsSelector(meta.Spec.Selector)
		if err != nil {
			return errors.WithStack(err)
		}

		if err := clientset.PodDeleteCollection(ctx, cluster, req.Resource.Namespace, metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: selector.String()}); err != nil {
			return errors.WithStack(err)
		}
	}

	if err = clientset.DeploymentDelete(ctx, cluster, req.Resource.Namespace, req.Resource.Name, req.DeleteOpts); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) DeploymentBatchDelete(ctx context.Context, req *iapiserver.DeploymentBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.DeploymentRequest] {
	wg := waitgroup.RunGenericConcurrently(ctx, req.Resources, func(ctx context.Context, res *iapiserver.DeploymentRequest) waitgroup.GenericResult[*iapiserver.DeploymentRequest] {
		cluster, err := k.store.Clusters().Get(ctx, res.Cluster)
		if err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		if err := clientset.DeploymentDelete(ctx, cluster, res.Resource.Namespace, res.Resource.Name, res.DeleteOpts); err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		return waitgroup.NewGenericResult(res, nil)
	})
	return wg.BatchGenericOutput()
}

// TODO: rewrite logic
func (k *kubernetesService) DeploymentRecreate(ctx context.Context, req *iapiserver.DeploymentRequest) (*iapiserver.DeploymentInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.DeploymentGet(ctx, cluster, req.Resource.Namespace, req.Resource.Name, req.GetOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	selector, err := metav1.LabelSelectorAsSelector(meta.Spec.Selector)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := clientset.PodDeleteCollection(ctx, cluster, req.Resource.Namespace, metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: selector.String()}); err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sDeploymentToApi(meta, cluster), nil
}

func (k *kubernetesService) DeploymentGet(ctx context.Context, req *iapiserver.DeploymentGetRequest) (*iapiserver.DeploymentInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.DeploymentGet(ctx, cluster, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sDeploymentToApi(meta, cluster), nil
}

func (k *kubernetesService) DeploymentList(ctx context.Context, req *iapiserver.DeploymentListRequest) (*iapiserver.DeploymentListResponse, error) {
	resp := &iapiserver.DeploymentListResponse{}
	clusters, err := getVisitScope(ctx, k.store, req.ResourceListRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	wg := waitgroup.RunGenericConcurrently[*iapiserver.Cluster, iapiserver.EachResourceRangeListState[*iapiserver.DeploymentInfo]](
		ctx, clusters, func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.DeploymentInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.DeploymentInfo](cluster.ID, cluster.Name)
			resList, err := libkubernetes.DeploymentList(ctx, cluster.Config, req.Namespace, req.ToListOpts())
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.DeploymentInfo
			for i := range resList.Items {
				resInfo := convertK8sDeploymentToApi(&resList.Items[i], cluster)
				if NewObjectCommonFieldFilter(resInfo.Resource).Filter(req.Fuzzy) {
					continue
				}
				resInfos = append(resInfos, resInfo)
			}
			clusterListOne.TotalCount = len(resInfos)
			clusterListOne.List = resInfos
			return waitgroup.NewGenericResult(clusterListOne, err)
		})

	resp.EachRangeListState = wg.GetSuccessResultList()
	resp.TotalCount = CutPagingSliceResourceList[*iapiserver.DeploymentInfo](resp.EachRangeListState, &resp.List, req.PageNum, req.PageSize, func(i, j int) bool {
		switch req.SortBy {
		case iapiserver.KubernetesResourceDeploymentSortByReplicas:
			replica1 := typeutil.GenericIndirectValue(resp.List[i].Resource.Spec.Replicas)
			replica2 := typeutil.GenericIndirectValue(resp.List[j].Resource.Spec.Replicas)
			if replica1 != replica2 {
				return compareutil.Compare(replica1, replica2, req.SortDesc)
			}
		}
		//fall back to common object param sort
		return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
	})
	return resp, nil
}

func getDependencyCurrentVersion(ctx context.Context, cluster *iapiserver.Cluster, denpendency *appsv1.Deployment) string {
	matchLabel := map[string]string{}
	if denpendency.Spec.Selector != nil && len(denpendency.Spec.Selector.MatchLabels) != 0 {
		matchLabel = denpendency.Spec.Selector.MatchLabels
	}
	sets := make([]string, 0, len(matchLabel))
	for k, v := range matchLabel {
		sets = append(sets, k+"="+v)
	}

	selector := strings.Join(sets, ",")
	listOption := metav1.ListOptions{LabelSelector: selector}
	rsList, err := clientset.ReplicaSetList(ctx, cluster, denpendency.Namespace, listOption)
	if err != nil {
		log.Errorf("getDependencyCurrentVersion error:%s", err.Error())
		return ""
	}

	for _, v := range rsList.Items {
		if v.Status.AvailableReplicas != 0 {
			version, ok := v.Annotations[deploymentVersionFlag]
			if ok {
				return version
			}
		}
	}

	log.Error("getDependencyCurrentVersion error: not found")
	return ""
}

func (k *kubernetesService) DeploymentVersionUpdate(ctx context.Context, req *iapiserver.DeploymentRequest) error {
	if req.Version == "" {
		return errors.Errorf("version is empty")
	}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	// check deployment state
	deployment, err := clientset.DeploymentGet(ctx, cluster, req.Resource.Namespace, req.Resource.Name, metav1.GetOptions{})
	if err != nil {
		return errors.WithStack(err)
	}
	if deployment.Spec.Paused {
		return errors.Errorf("you cannot rollback a paused deployment; resume it first with 'kubectl rollout resume deployment/%s' and try again", deployment.Name)

	}

	// find to revision
	toRevision, err := deploymentToRevision(ctx, cluster, deployment, req.Version)
	if err != nil {
		return errors.WithStack(err)
	}
	if equalIgnoreHash(&toRevision.Spec.Template, &deployment.Spec.Template) {
		return errors.Errorf("same revision")
	}

	// patch
	annotationsToSkip := getAnnotationsToSkip()
	annotations := map[string]string{}
	for k := range annotationsToSkip {
		if v, ok := deployment.Annotations[k]; ok {
			annotations[k] = v
		}
	}
	for k, v := range toRevision.Annotations {
		if !annotationsToSkip[k] {
			annotations[k] = v
		}
	}
	delete(toRevision.Spec.Template.Labels, appsv1.DefaultDeploymentUniqueLabelKey)
	patch, err := getDeploymentPatch(&toRevision.Spec.Template, annotations)
	if err != nil {
		return errors.WithStack(err)
	}
	if _, err := clientset.DeploymentPatch(ctx, cluster, req.Resource.Namespace, req.Resource, kubetypes.JSONPatchType, patch, metav1.PatchOptions{}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) DeploymentVersionList(ctx context.Context, req *iapiserver.DeploymentVersionListRequest) (*iapiserver.DeploymentVersionListResponse, error) {
	resp := &iapiserver.DeploymentVersionListResponse{}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dep, err := clientset.DeploymentGet(ctx, cluster, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rsList, err := replicasetList(ctx, cluster, dep)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// create tmp deployment for user
	for _, v := range rsList {
		version, ok := v.Annotations[deploymentVersionFlag]
		if ok {
			tmpDep, err := deploymentCopy(dep)
			if err != nil {
				log.Errorf("%s", err)
				continue
			}

			tmpDep.Annotations[deploymentVersionFlag] = version
			tmpDep.Spec.Template = v.Spec.Template

			versionInfo := iapiserver.VersionInfo{
				Version:    version,
				CreateTime: v.CreationTimestamp.Unix(),
				Deployment: tmpDep,
			}
			resp.VersionList = append(resp.VersionList, versionInfo)
		}
	}

	resp.CurrentVersion = &iapiserver.VersionInfo{Namespace: req.Namespace, Name: req.Name, Version: dep.Annotations[deploymentVersionFlag]}
	sort.SliceStable(resp.VersionList, func(i, j int) bool {
		return resp.VersionList[j].CreateTime > resp.VersionList[i].CreateTime
	})
	resp.TotalCount = len(resp.VersionList)
	s, e := paging.Index(resp.TotalCount, req.PageNum, req.PageSize)
	resp.VersionList = resp.VersionList[s:e]

	return resp, nil
}

func getDeploymentPatch(podTemplate *corev1.PodTemplateSpec, annotations map[string]string) ([]byte, error) {
	patch, err := json.Marshal([]interface{}{
		map[string]interface{}{
			"op":    "replace",
			"path":  "/spec/template",
			"value": podTemplate,
		},
		map[string]interface{}{
			"op":    "replace",
			"path":  "/metadata/annotations",
			"value": annotations,
		},
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return patch, nil
}

func deploymentToRevision(ctx context.Context, cluster *iapiserver.Cluster, deployment *appsv1.Deployment, toRevision string) (revision *appsv1.ReplicaSet, err error) {
	rsList, err := replicasetList(ctx, cluster, deployment)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, v := range rsList {
		if v.Annotations[deploymentVersionFlag] == toRevision {
			return v, nil
		}
	}

	return nil, errors.Errorf("revision[%v] not found", toRevision)
}

func replicasetList(ctx context.Context, cluster *iapiserver.Cluster, deployment *appsv1.Deployment) ([]*appsv1.ReplicaSet, error) {
	var resp []*appsv1.ReplicaSet

	namespace := deployment.Namespace
	selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	options := metav1.ListOptions{LabelSelector: selector.String()}
	all, err := clientset.ReplicaSetList(ctx, cluster, namespace, options)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, rs := range all.Items {
		rs := rs
		if metav1.IsControlledBy(&rs, deployment) {
			resp = append(resp, &rs)
		}
	}

	return resp, nil
}

func getAnnotationsToSkip() map[string]bool {
	return map[string]bool{
		corev1.LastAppliedConfigAnnotation: true,
		deploy.RevisionAnnotation:          true,
		deploy.RevisionHistoryAnnotation:   true,
		deploy.DesiredReplicasAnnotation:   true,
		deploy.MaxReplicasAnnotation:       true,
		appsv1.DeprecatedRollbackTo:        true,
	}
}

func equalIgnoreHash(template1, template2 *corev1.PodTemplateSpec) bool {
	t1Copy := template1.DeepCopy()
	t2Copy := template2.DeepCopy()
	// Remove hash labels from template.Labels before comparing
	delete(t1Copy.Labels, appsv1.DefaultDeploymentUniqueLabelKey)
	delete(t2Copy.Labels, appsv1.DefaultDeploymentUniqueLabelKey)

	return apiequality.Semantic.DeepEqual(t1Copy, t2Copy)
}

func deploymentCopy(deploy *appsv1.Deployment) (*appsv1.Deployment, error) {
	resp := &appsv1.Deployment{}

	if deploy == nil {
		return resp, nil
	}

	*resp = *deploy
	resp.Annotations = map[string]string{}
	for k, v := range deploy.Annotations {
		resp.Annotations[k] = v
	}
	data, err := json.Marshal(deploy.Spec)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := json.Unmarshal(data, &resp.Spec); err != nil {
		return nil, errors.WithStack(err)
	}

	return resp, nil
}

func convertK8sDeploymentToApi(meta *appsv1.Deployment, cluster *iapiserver.Cluster) *iapiserver.DeploymentInfo {
	resp := iapiserver.NewDeploymentInfo(meta, cluster)

	var warn *iapiserver.ReasonMessage
	for _, v := range meta.Status.Conditions {
		if v.Type == appsv1.DeploymentReplicaFailure && v.Status == corev1.ConditionTrue {
			warn = &iapiserver.ReasonMessage{
				Message: v.Message,
				Reason:  v.Reason,
			}
		}
	}
	resp.Warning = warn
	resp.ResourceConvert = make([]*iapiserver.ResourceConvert, 0, len(meta.Spec.Template.Spec.Containers))
	for _, container := range meta.Spec.Template.Spec.Containers {
		resp.ResourceConvert = append(resp.ResourceConvert, convertResourceLimitToPersistentUnit(container.Name, container.Resources.Requests, container.Resources.Limits))
	}

	return resp
}
