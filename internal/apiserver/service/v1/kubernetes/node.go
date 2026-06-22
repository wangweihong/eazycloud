package kubernetes

import (
	"context"
	"fmt"
	"time"

	//"github.com/wangweihong/gotoolbox/pkg/async"
	"github.com/wangweihong/gotoolbox/pkg/compareutil"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/log"
	"github.com/wangweihong/gotoolbox/pkg/mathutil"
	"github.com/wangweihong/gotoolbox/pkg/paging"
	"github.com/wangweihong/gotoolbox/pkg/waitgroup"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubetypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/reference"
	"k8s.io/kubectl/pkg/scheme"
	"k8s.io/kubectl/pkg/util/resource"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/apis/imachinery"
	"github.com/wangweihong/eazycloud/internal/pkg/clientset"
)

const (
	// master label in kubernetes version 1.18
	nodeMasterLabelKey = "node-role.kubernetes.io/master"
	// master label in kubernetes version 1.30
	nodeMasterLabelKeyNew = "node-role.kubernetes.io/control-plane"
	nodeStatusHealth      = "healthy"
	nodeStatusUnhealth    = "unhealthy"

	actionGatewayCancel = "gateway-cancel"
	actionGatewaySet    = "gateway-set"
)

var (
	nodeStatusReady = corev1.NodeConditionType("Ready")
)

func (k *kubernetesService) NodeGatewayUpdate(ctx context.Context, req *iapiserver.NodeRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	node, err := clientset.NodeGet(ctx, cluster, req.Node.Name, req.GetOpts)
	if err != nil {
		return errors.WithStack(err)
	}

	if getNodeRole(node) == iapiserver.NodeRoleMaster {
		return errors.Errorf("master node cannot be gateway node")
	}

	if err := k.nodeActionUpdate(req.Action, node); err != nil {
		return errors.WithStack(err)
	}

	node, err = clientset.NodeUpdate(ctx, cluster, node, req.UpdateOpts)
	if err != nil {
		return errors.WithStack(err)
	}

	// 删掉本节点和nginx相关的pod
	if req.Action == actionGatewayCancel {
		go func() {
			podList, err := clientset.PodList(ctx, cluster, iapiserver.IngressControllerNamespace, req.ListOpts)
			if err != nil {
				log.Errorf("list pod in gateway node fail:%v", err)
				return
			}
			for _, v := range podList.Items {
				if v.Spec.NodeName != req.Node.Name {
					continue
				}
				go func(pod v1.Pod) {
					if err := clientset.PodDelete(ctx, cluster, iapiserver.IngressControllerNamespace, &pod, req.DeleteOpts); err != nil {
						log.Errorf("delete pod in gateway node fail:%v", err)
					}
				}(v)
			}
		}()
	}

	return nil
}

func (k *kubernetesService) IsGatewayNode(node *v1.Node) bool {
	if node == nil || len(node.Labels) == 0 {
		return false
	}

	gateway, _ := node.Labels["type"]
	if gateway != iapiserver.NodeRoleGateway {
		return false
	}

	for _, v := range node.Spec.Taints {
		if v.Key == "type" && v.Value == iapiserver.NodeRoleGateway && v.Effect == "NoExecute" {
			return true
		}
	}

	return false
}

func (k *kubernetesService) nodeActionUpdate(action string, node *v1.Node) error {
	if node == nil {
		return nil
	}
	if node.Labels == nil {
		node.Labels = map[string]string{}
	}
	switch action {
	case actionGatewaySet:
		node.Labels["type"] = iapiserver.NodeRoleGateway
		for k, v := range node.Spec.Taints {
			if v.Key == "type" && v.Value == iapiserver.NodeRoleGateway {
				v.Effect = "NoExecute"
				node.Spec.Taints[k] = v
				return nil
			}
		}
		node.Spec.Taints = append(node.Spec.Taints, v1.Taint{Key: "type", Value: iapiserver.NodeRoleGateway, Effect: "NoExecute"})
	case actionGatewayCancel:
		delete(node.Labels, "type")
		for k, v := range node.Spec.Taints {
			if v.Key == "type" && v.Value == iapiserver.NodeRoleGateway {
				node.Spec.Taints = append(node.Spec.Taints[0:k], node.Spec.Taints[k+1:]...)
				return nil
			}
		}
	default:
		return errors.Errorf("action[%v] not support", action)
	}

	return nil
}

func (k *kubernetesService) NodeTaintUpdate(ctx context.Context, req *iapiserver.NodeRequest) (*iapiserver.NodeInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	node, err := clientset.NodeGet(ctx, cluster, req.Node.Name, req.GetOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	node.Spec.Taints = req.Node.Spec.Taints
	meta, err := clientset.NodeUpdate(ctx, cluster, node, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sNodeToNodeInfo(meta, cluster, req.Yaml), nil
}

// 修改标签, 减少前端的参数传递
func (k *kubernetesService) NodeUpdateLabel(ctx context.Context, req *iapiserver.NodeRequest) (*iapiserver.NodeInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	//TODO: restrict label prefix check
	node, err := clientset.NodeGet(ctx, cluster, req.Node.Name, req.GetOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	node.Labels = req.Node.Labels
	node, err = clientset.NodeUpdate(ctx, cluster, node, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sNodeToNodeInfo(node, cluster, req.Yaml), nil
}

func (k *kubernetesService) NodeAdd(ctx context.Context, req *iapiserver.NodeRequest) (*iapiserver.NodeInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	node, err := clientset.NodeCreate(ctx, cluster, req.Node, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sNodeToNodeInfo(node, cluster, req.Yaml), nil
}

func (k *kubernetesService) NodeUpdate(ctx context.Context, req *iapiserver.NodeRequest) (*iapiserver.NodeInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	node, err := clientset.NodeUpdate(ctx, cluster, req.Node, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return convertK8sNodeToNodeInfo(node, cluster, req.Yaml), nil
}

func (k *kubernetesService) NodeDelete(ctx context.Context, req *iapiserver.NodeRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = clientset.NodeDelete(ctx, cluster, req.Node, req.DeleteOpts); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) NodeGet(ctx context.Context, req *iapiserver.NodeGetRequest) (*iapiserver.NodeInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	node, err := clientset.NodeGet(ctx, cluster, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	nodeInfo := convertK8sNodeToNodeInfo(node, cluster, req.Yaml)
	// if req.ShowResourceUsage {
	// 	nodeInfo.ResourceUsage = k.getNodeRealtimeResourceUsage(ctx, cluster, node, 5)
	// 	if nodeInfo.ResourceUsage != nil {
	// 		nodeInfo.ResourceUsage.CpuUsed = nodeInfo.ResourceUsage.CpuUsed / 1000
	// 		nodeInfo.ResourceUsage.CpuUsedRatio = nodeInfo.ResourceUsage.CpuUsedRatio / 1000
	// 	}
	// }

	return nodeInfo, nil
}

func (k *kubernetesService) NodeList(ctx context.Context, req *iapiserver.NodeListRequest) (*iapiserver.NodeListResponse, error) {
	resp := &iapiserver.NodeListResponse{}

	var list []*iapiserver.NodeInfo
	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList[*iapiserver.NodeInfo](ctx, k.store, &list, req.ResourceListRequest,
		func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.NodeInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.NodeInfo](cluster.ID, cluster.Name)

			nodeList, err := clientset.NodeList(ctx, cluster, req.ToListOpts())
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			nodes := make([]*iapiserver.NodeInfo, 0)
			for i := range nodeList.Items {
				nodeInfo := convertK8sNodeToNodeInfo(&nodeList.Items[i], cluster, req.Yaml)
				if filterNode(nodeInfo, req.Fuzzy) {
					continue
				}
				nodes = append(nodes, nodeInfo)
			}
			///	if req.ShowResourceUsage && clientset.IsMonitorServiceReady(cluster.UUID) {
			// if req.ShowResourceUsage {
			// 	wg := waitgroup.NewWaitGroup(ctx)
			// 	for i := range nodes {
			// 		i := i
			// 		wg.Start(waitgroup.NewWaitGroupHandleFunc.NewWaitGroupHandleFunc("", nil, func() utils.WaitGroupResult {
			// 			nodes[i].ResourceUsage = getNodeRealtimeResourceUsage(nil, cluster, nodes[i].Node, 5, "ignoreMonitor")
			// 			return waitgroup.NewResult(nil, nil)
			// 		}))
			// 	}
			// 	wg.Wait()
			// }

			clusterListOne.List = nodes
			clusterListOne.TotalCount = len(nodes)
			return waitgroup.NewGenericResult(clusterListOne, nil)
		}, func(i, j int) bool {
			switch req.SortBy {
			case iapiserver.KubernetesResourceNodeSortBYCpuResourceLimitRatio:
				if resp.List[i].ResourceUsage.CpuResourceLimitRatio != resp.List[j].ResourceUsage.CpuResourceLimitRatio { // if equal, compare create time instead
					return compareutil.Compare(resp.List[i].ResourceUsage.CpuResourceLimitRatio, resp.List[j].ResourceUsage.CpuResourceLimitRatio, req.SortDesc)
				}
			case iapiserver.KubernetesResourceNodeSortBYCpuResourceRequestRatio:
				if resp.List[i].ResourceUsage.CpuResourceRequestRatio != resp.List[j].ResourceUsage.CpuResourceRequestRatio { // if equal, compare create time instead
					return compareutil.Compare(resp.List[i].ResourceUsage.CpuResourceRequestRatio, resp.List[j].ResourceUsage.CpuResourceRequestRatio, req.SortDesc)
				}
			case iapiserver.KubernetesResourceNodeSortBYMemoryResourceRequestRatio:
				if resp.List[i].ResourceUsage.MemoryResourceRequestRatio != resp.List[j].ResourceUsage.MemoryResourceRequestRatio { // if equal, compare create time instead
					return compareutil.Compare(resp.List[i].ResourceUsage.MemoryResourceRequestRatio, resp.List[j].ResourceUsage.MemoryResourceRequestRatio, req.SortDesc)
				}
			case iapiserver.KubernetesResourceNodeSortBYMemoryResourceLimitRatio:
				if resp.List[i].ResourceUsage.MemoryResourceLimitRatio != resp.List[j].ResourceUsage.MemoryResourceLimitRatio { // if equal, compare create time instead
					return compareutil.Compare(resp.List[i].ResourceUsage.MemoryResourceLimitRatio, resp.List[j].ResourceUsage.MemoryResourceLimitRatio, req.SortDesc)
				}
			case iapiserver.KubernetesResourceNodeSortBYPodUsedRatio:
				if resp.List[i].ResourceUsage.PodUsedRatio != resp.List[j].ResourceUsage.PodUsedRatio { // if equal, compare create time instead
					return compareutil.Compare(resp.List[i].ResourceUsage.PodUsedRatio, resp.List[j].ResourceUsage.PodUsedRatio, req.SortDesc)
				}

			}

			return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
		}, 10*time.Second)
	resp.List = list
	return resp, err
}

// func getNodeRealtimeResourceUsage(ctx context.Context, clusterInfo *iapiserver.Cluster, node *v1.Node, timeout int64, ignoreMonitor ...string) *iapiserver.NodeRealtimeResource {
// 	nrr := &iapiserver.NodeRealtimeResource{}
// 	wg := waitgroup.NewWaitGroup(nil)
// 	// only work  when monitor service is ready
// 	if clientset.IsMonitorServiceReady(clusterInfo.ID) {
// 		wg.Start(utils.NewWaitGroupHandleFunc("", nil, func() utils.WaitGroupResult {
// 			done := make(chan utils.WaitGroupResult, 1)
// 			go func() {
// 				defer utils.AsyncPanicRecover(nil, func() {
// 					done <- utils.NewWaitGroupResult(nil, fmt.Errorf("panic"))
// 				})
// 				ret, err := getNodeRealTimeMonitorResourceUsage(ctx, clusterInfo, node, 0)
// 				if err != nil {
// 					logrus.Errorf("getNodeRealTimeMonitorResourceUsage err:%v", err.Error())
// 					done <- utils.NewWaitGroupResult(nil, err)
// 				}
// 				done <- utils.NewWaitGroupResult(ret, nil)
// 			}()

// 			select {
// 			case ret := <-done:
// 				return ret

// 			case <-time.After(time.Duration(timeout) * time.Second):
// 				return utils.NewWaitGroupResult(nil, fmt.Errorf("timeout"))
// 			}
// 		}))
// 	}
// 	wg.Start(utils.NewWaitGroupHandleFunc("", nil, func() utils.WaitGroupResult {
// 		podList, err := nodeNonTerminatedPodsList(nil, clusterInfo, node, timeout)
// 		if err != nil {
// 			logrus.Errorf("PodListTimeout error: %v", err)
// 			return utils.NewWaitGroupResult(nil, err)
// 		}
// 		podsRequestAndLimit := calculateNodesPodsResourceRequestAndLimit(podList, node)
// 		return utils.NewWaitGroupResult(podsRequestAndLimit, nil)
// 	}))
// 	// don't not use monitor to get pod usage. when monitor broken, it will broke pod usage
// 	wg.Start(utils.NewWaitGroupHandleFunc("", nil, func() utils.WaitGroupResult {
// 		podList, err := nodeAllPodsList(nil, clusterInfo, node)
// 		if err != nil {
// 			logrus.Errorf("PodListTimeout error:%v", err)
// 			return utils.NewWaitGroupResult(nil, err)
// 		}
// 		podUsage := calculateNodesPodsUsage(podList, node)
// 		return utils.NewWaitGroupResult(podUsage, nil)
// 	}))
// 	wg.Wait()

// 	for _, v := range wg.GetResults() {
// 		if v.Error != nil {
// 			continue
// 		}
// 		if v.Data != nil {
// 			switch ret := v.Data.(type) {
// 			case *iapiserver.NodeRealTimePhysicalResourceUsage:
// 				nrr.NodeRealTimePhysicalResourceUsage = ret
// 			case *iapiserver.NodePodListResourceRequestAndLimit:
// 				nrr.NodePodListResourceRequestAndLimit = ret
// 			case *iapiserver.NodeRealTimePodUsage:
// 				nrr.NodeRealTimePodUsage = ret
// 			}
// 		}
// 	}

// 	if nrr.NodeRealTimePhysicalResourceUsage == nil {
// 		nrr.NodeRealTimePhysicalResourceUsage = &topke.NodeRealTimePhysicalResourceUsage{}
// 		nrr.CpuCapacity, nrr.MemoryCapacity = getCpuMemoryCapacityFromNode(node)
// 	}

// 	return nrr
// }

// func getNodeRealTimeMonitorResourceUsage(ctx context.Context, clusterInfo *iapiserver.Cluster, node *v1.Node, timeout int64) (*topke.NodeRealTimePhysicalResourceUsage, error) {
// 	resourceUsage := &topke.NodeRealTimePhysicalResourceUsage{}
// 	ret, err := GetTopkeManager().NodeState(clusterInfo, node.Name, timeout)
// 	if err != nil {
// 		logrus.Errorf("MonitorGetNodeState err:%v", err)
// 		return nil, err
// 	}
// 	ret.Addr = getNodeAddr(node)
// 	resourceUsage.CpuCapacity = float64(node.Status.Capacity.Cpu().MilliValue())
// 	resourceUsage.MemoryCapacity = float64(node.Status.Capacity.Memory().Value())
// 	resourceUsage.CpuUsed = ret.CpuUsed * 1000
// 	resourceUsage.MemoryUsed = ret.MemoryUsed
// 	resourceUsage.DiskUsed = ret.DiskUsed
// 	resourceUsage.DiskCapacity = ret.DiskTotal

// 	resourceUsage.CpuUsedRatio = Divide(resourceUsage.CpuUsed, resourceUsage.CpuCapacity) * 100
// 	resourceUsage.MemoryUsedRatio = Divide(resourceUsage.MemoryUsed, resourceUsage.MemoryCapacity) * 100
// 	resourceUsage.DiskUsedRatio = Divide(resourceUsage.DiskUsed, resourceUsage.DiskCapacity) * 100

// 	return resourceUsage, nil
// }

func (k *kubernetesService) NodeCordon(ctx context.Context, req *iapiserver.NodeRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}
	if _, err := clientset.NodeCordon(ctx, cluster, req.Node); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) NodeUncordon(ctx context.Context, req *iapiserver.NodeRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}
	if _, err := clientset.NodeUncordon(ctx, cluster, req.Node); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) NodeDrain(ctx context.Context, req *iapiserver.NodeDrainRequest) (*imachinery.BatchOutput, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	node, err := clientset.NodeGet(ctx, cluster, req.Node.Name, req.GetOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	podList, err := nodeAllPodsList(ctx, cluster, node)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !req.DisableEviction {
		return evictOrDeletePods(cluster, podList, clientset.PodEvict), nil
	}
	return evictOrDeletePods(cluster, podList, clientset.PodDelete), nil
}

func (k *kubernetesService) NodeEvent(ctx context.Context, req *iapiserver.NodeEventRequest) ([]*iapiserver.EventInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	node, err := clientset.NodeGet(ctx, cluster, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ref, err := reference.GetReference(scheme.Scheme, node)
	if err != nil {
		return nil, errors.Errorf("Unable to construct reference to '%#v': %v", node, err)
	}

	ref.UID = kubetypes.UID(ref.Name)
	resList, err := clientset.EventSearch(ctx, cluster, scheme.Scheme, ref)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var resInfos []*iapiserver.EventInfo
	for i := range resList.Items {
		resInfo := convertK8sEventTopApi(&resList.Items[i], cluster)
		if NewObjectCommonFieldFilter(resInfo.Resource).Filter(req.Fuzzy) {
			continue
		}
		resInfos = append(resInfos, resInfo)
	}

	s, e := paging.Index(len(resInfos), req.PageNum, req.PageSize)
	return resInfos[s:e], nil
}

// func (ks *kubernetesService) getNodeRealtimeResourceUsage(ctx context.Context, clusterInfo *iapiserver.Cluster, node *v1.Node, timeout int64, ignoreMonitor ...string) *iapiserver.NodeRealtimeResource {
// 	nrr := &iapiserver.NodeRealtimeResource{}
// 	wg := waitgroup.NewWaitGroup(ctx)
// 	// only work  when monitor service is ready
// 	if clientset.IsMonitorServiceReady(clusterInfo.ID) {
// 		wg.Start(waitgroup.NewWaitGroupHandleFunc(ctx, "", func() waitgroup.Result {
// 			done := make(chan waitgroup.Result, 1)
// 			async.GoRoutineCustomPanicCover(ctx, func() {
// 				done <- waitgroup.NewResult(nil, fmt.Errorf("panic"))
// 			}, func(ctx context.Context) {
// 				ret, err := ks.getNodeRealTimeMonitorResourceUsage(ctx, clusterInfo, node, 0)
// 				if err != nil {
// 					done <- waitgroup.NewResult(nil, err)
// 				}
// 				done <- waitgroup.NewResult(ret, nil)
// 			})
// 			select {
// 			case ret := <-done:
// 				return ret

// 			case <-time.After(time.Duration(timeout) * time.Second):
// 				log.Errorf("getNodeRealTimeMonitorResourceUsage timeout, %v second has pass ", timeout)
// 				return waitgroup.NewResult(nil, fmt.Errorf("timeout"))
// 			}
// 		}))
// 	}
// 	wg.Start(waitgroup.NewWaitGroupHandleFunc(ctx, "", func() waitgroup.Result {
// 		podList, err := nodeNonTerminatedPodsList(ctx, clusterInfo, node, timeout)
// 		if err != nil {
// 			log.Errorf("PodListTimeout error: %v", err)
// 			return waitgroup.NewResult(nil, err)
// 		}
// 		podsRequestAndLimit := calculateNodesPodsResourceRequestAndLimit(podList, node)
// 		return waitgroup.NewResult(podsRequestAndLimit, nil)
// 	}))
// 	// don't not use monitor to get pod usage. when monitor broken, it will broke pod usage
// 	wg.Start(waitgroup.NewWaitGroupHandleFunc(ctx, "", func() waitgroup.Result {
// 		podList, err := nodeAllPodsList(nil, clusterInfo, node)
// 		if err != nil {
// 			log.Errorf("PodListTimeout error:%v", err)
// 			return waitgroup.NewResult(nil, err)
// 		}
// 		podUsage := calculateNodesPodsUsage(podList, node)
// 		return waitgroup.NewResult(podUsage, nil)
// 	}))
// 	wg.Wait()

// 	for _, v := range wg.GetResults() {
// 		if v.Error != nil {
// 			continue
// 		}
// 		if v.Data != nil {
// 			switch ret := v.Data.(type) {
// 			case *iapiserver.NodeRealTimePhysicalResourceUsage:
// 				nrr.NodeRealTimePhysicalResourceUsage = ret
// 			case *iapiserver.NodePodListResourceRequestAndLimit:
// 				nrr.NodePodListResourceRequestAndLimit = ret
// 			case *iapiserver.NodeRealTimePodUsage:
// 				nrr.NodeRealTimePodUsage = ret
// 			}
// 		}
// 	}

// 	if nrr.NodeRealTimePhysicalResourceUsage == nil {
// 		nrr.NodeRealTimePhysicalResourceUsage = &iapiserver.NodeRealTimePhysicalResourceUsage{}
// 		nrr.CpuCapacity, nrr.MemoryCapacity = getCpuMemoryCapacityFromNode(node)
// 	}

// 	return nrr
// }

func getCpuMemoryCapacityFromNode(node *v1.Node) (float64, float64) {
	return float64(node.Status.Capacity.Cpu().MilliValue()), float64(node.Status.Capacity.Memory().Value())
}

// func (k *kubernetesService) getNodeRealTimeMonitorResourceUsage(ctx context.Context, clusterInfo *iapiserver.Cluster, node *v1.Node, timeout int64) (*iapiserver.NodeRealTimePhysicalResourceUsage, error) {
// 	resourceUsage := &iapiserver.NodeRealTimePhysicalResourceUsage{}
// 	ret, err := k.NodeState(clusterInfo, node.Name, timeout)
// 	if err != nil {
// 		log.Errorf("MonitorGetNodeState err:%v", err)
// 		return nil, err
// 	}
// 	ret.Addr = getNodeAddr(node)
// 	resourceUsage.CpuCapacity = float64(node.Status.Capacity.Cpu().MilliValue())
// 	resourceUsage.MemoryCapacity = float64(node.Status.Capacity.Memory().Value())
// 	resourceUsage.CpuUsed = ret.CpuUsed * 1000
// 	resourceUsage.MemoryUsed = ret.MemoryUsed
// 	resourceUsage.DiskUsed = ret.DiskUsed
// 	resourceUsage.DiskCapacity = ret.DiskTotal

// 	resourceUsage.CpuUsedRatio = mathutil.Divide(resourceUsage.CpuUsed, resourceUsage.CpuCapacity) * 100
// 	resourceUsage.MemoryUsedRatio = mathutil.Divide(resourceUsage.MemoryUsed, resourceUsage.MemoryCapacity) * 100
// 	resourceUsage.DiskUsedRatio = mathutil.Divide(resourceUsage.DiskUsed, resourceUsage.DiskCapacity) * 100

// 	return resourceUsage, nil
// }

// func (k *kubernetesService) NodeState(ctx context.Context, clusterInfo *iapiserver.Cluster, nodeName string, timeout int64) (*iapiserver.MonitorData, error) {
// 	resp := &iapiserver.MonitorData{Name: nodeName}

// 	rets, err := getClusterNodesOverview(ctx, clusterInfo)
// 	if err != nil {
// 		return nil, err
// 	}
// 	nodeOverview, ok := rets[nodeName]
// 	if !ok {
// 		return resp, nil
// 	}

// 	resp.Cpu = nodeOverview.CpuUsedRatio
// 	resp.CpuUsed = nodeOverview.CpuUsed
// 	resp.CpuTotal = nodeOverview.CpuCapacity
// 	resp.Memory = nodeOverview.MemoryUsedRatio
// 	resp.MemoryUsed = nodeOverview.MemoryUsed
// 	resp.MemoryTotal = nodeOverview.MemoryCapacity
// 	resp.DiskRatio = nodeOverview.DiskUsedRatio
// 	resp.DiskUsed = nodeOverview.DiskUsed
// 	resp.DiskTotal = nodeOverview.DiskCapacity
// 	resp.PodRatio = nodeOverview.PodUsedRatio
// 	resp.PodCount = float64(nodeOverview.PodUsed)
// 	resp.PodTotal = float64(nodeOverview.PodCapacity)

// 	return resp, nil
// }

// func getClusterNodesOverview(ctx context.Context, clusterInfo *iapiserver.Cluster) (map[string]*iapiserver.NodeOverview, error) {
// 	resp := map[string]*iapiserver.NodeOverview{}

// 	nodeList, err := clientset.NodeList(ctx, clusterInfo, metav1.ListOptions{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	rets, err := k.MonitorNodeStateList(context.Background(), clusterInfo)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, v := range nodeList.Items {
// 		no := &iapiserver.NodeOverview{NodeName: v.Name}
// 		resp[no.NodeName] = no

// 		no.NodeAddr = getNodeAddr(&v)
// 		no.Role = getNodeRole(&v)
// 		no.Status = nodeStatusHealth
// 		no.OsImage, no.ContainerRuntimeVersion = getNodeImageVersion(&v)
// 		if !isNodeReady(&v) {
// 			no.Status = nodeStatusUnhealth
// 		}
// 		if no.Role == nodeRoleGateway {
// 			no.IsGateway = true
// 		}

// 		nodeState, ok := rets[no.NodeName]
// 		if !ok {
// 			log.Errorf("prom nodename[%v]", no.NodeName)
// 			continue
// 		}
// 		no.ClusterUUID = clusterInfo.ID
// 		no.ClusterName = clusterInfo.Name
// 		no.CreateTime = v.CreationTimestamp.Unix()
// 		no.CpuUsed = nodeState.CpuUsed
// 		no.CpuCapacity = nodeState.CpuCapacity
// 		no.CpuUsedRatio = nodeState.CpuUsedRatio
// 		no.CpuResourceRequest = nodeState.CpuResourceRequest
// 		no.CpuResourceLimit = nodeState.CpuResourceLimit
// 		no.CpuResourceLimitRatio = Divide(float64(no.CpuResourceLimit), no.CpuCapacity) * 100
// 		no.CpuResourceRequestRatio = Divide(float64(no.CpuResourceRequest), no.CpuCapacity) * 100
// 		no.MemoryUsed = nodeState.MemoryUsed
// 		no.MemoryCapacity = nodeState.MemoryCapacity
// 		no.MemoryUsedRatio = nodeState.MemoryUsedRatio
// 		no.MemoryResourceRequest = nodeState.MemoryResourceRequest
// 		no.MemoryResourceLimit = nodeState.MemoryResourceLimit
// 		no.MemoryResourceLimitRatio = Divide(float64(no.MemoryResourceLimit), no.MemoryCapacity) * 100
// 		no.MemoryResourceRequestRatio = Divide(float64(no.MemoryResourceRequest), no.MemoryCapacity) * 100
// 		no.DiskUsedRatio = nodeState.DiskUsedRatio * 100
// 		no.DiskUsed = nodeState.DiskUsed
// 		no.DiskCapacity = nodeState.DiskCapacity
// 		no.PodUsed = nodeState.PodUsed
// 		no.PodCapacity = nodeState.PodCapacity
// 		no.PodUsedRatio = nodeState.PodUsedRatio * 100
// 	}

// 	return nil, nil
// }

func nodeNonTerminatedPodsList(ctx context.Context, clusterInfo *iapiserver.Cluster, node *v1.Node, timeout int64) ([]*v1.Pod, error) {
	podList, err := clientset.PodList(ctx, clusterInfo, metav1.NamespaceAll, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	pods := make([]*v1.Pod, 0)
	for i := range podList.Items {
		if podList.Items[i].Spec.NodeName == node.Name && podList.Items[i].Status.Phase != v1.PodSucceeded && podList.Items[i].Status.Phase != v1.PodFailed {
			pods = append(pods, &podList.Items[i])
		}
	}
	return pods, nil
}

func nodeAllPodsList(ctx context.Context, clusterInfo *iapiserver.Cluster, node *v1.Node) ([]*v1.Pod, error) {
	podList, err := clientset.PodList(ctx, clusterInfo, metav1.NamespaceAll, metav1.ListOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	pods := make([]*v1.Pod, 0)
	for i := range podList.Items {
		if podList.Items[i].Spec.NodeName == node.Name {
			pods = append(pods, &podList.Items[i])
		}
	}
	return pods, nil
}

func calculateNodesPodsResourceRequestAndLimit(pods []*v1.Pod, node *v1.Node) *iapiserver.NodePodListResourceRequestAndLimit {
	resourceRequest := &iapiserver.PodListResourceRequest{}
	resourceLimit := &iapiserver.PodListResourceLimit{}

	for _, pod := range pods {
		req, limit := resource.PodRequestsAndLimits(pod)
		cpuReq, cpuLimit, memoryReq, memoryLimit := req[v1.ResourceCPU], limit[v1.ResourceCPU], req[v1.ResourceMemory], limit[v1.ResourceMemory]
		resourceRequest.CpuResourceRequest += cpuReq.MilliValue()
		resourceRequest.MemoryResourceRequest += memoryReq.Value()
		resourceLimit.CpuResourceLimit += cpuLimit.MilliValue()
		resourceLimit.MemoryResourceLimit += memoryLimit.Value()
	}

	cpuCapacity := node.Status.Capacity.Cpu().MilliValue()
	memoryCapacity := node.Status.Capacity.Memory().Value()
	//calculate according to node
	resourceRequest.CpuResourceRequestRatio = mathutil.Divide(float64(resourceRequest.CpuResourceRequest), float64(cpuCapacity)) * 100
	resourceLimit.CpuResourceLimitRatio = mathutil.Divide(float64(resourceLimit.CpuResourceLimit), float64(cpuCapacity)) * 100
	// memory
	resourceRequest.MemoryResourceRequestRatio = mathutil.Divide(float64(resourceRequest.MemoryResourceRequest), float64(memoryCapacity)) * 100
	resourceLimit.MemoryResourceLimitRatio = mathutil.Divide(float64(resourceLimit.MemoryResourceLimit), float64(memoryCapacity)) * 100
	return &iapiserver.NodePodListResourceRequestAndLimit{
		NodeCpuCapacity:        float64(cpuCapacity),
		NodeMemoryCapacity:     float64(memoryCapacity),
		PodListResourceRequest: resourceRequest,
		PodListResourceLimit:   resourceLimit,
	}
}

func calculateNodesPodsUsage(pods []*v1.Pod, node *v1.Node) *iapiserver.NodeRealTimePodUsage {
	podUsage := &iapiserver.NodeRealTimePodUsage{}
	podUsage.PodCapacity = int(node.Status.Capacity.Pods().Value())
	podUsage.PodUsed = len(pods)
	podUsage.PodUsedRatio = mathutil.FloatDivide(float64(podUsage.PodUsed), float64(podUsage.PodCapacity)) * 100
	return podUsage
}

func convertK8sNodeToNodeInfo(node *v1.Node, cluster *iapiserver.Cluster, yaml bool) *iapiserver.NodeInfo {
	nodeInfo := &iapiserver.NodeInfo{
		Resource: node,
	}

	if !yaml {
		nodeInfo.NodeStatus = "healthy"
		if !isNodeReady(node) {
			nodeInfo.NodeStatus = "unhealthy"
		}
		nodeInfo.NodeAddr = getNodeAddr(node)
		nodeInfo.Role = getNodeRole(node)
	}
	return nodeInfo
}

type PodEvictOrDeleteFunc func(ctx context.Context, config *iapiserver.Cluster, namespace string, pod *v1.Pod, opts metav1.DeleteOptions) error

func evictOrDeletePods(cluster *iapiserver.Cluster, pods []*v1.Pod, drainFunc PodEvictOrDeleteFunc) *imachinery.BatchOutput {
	type PodResult struct {
		Pod metav1.ObjectMeta
		Err error
	}
	returnCh := make(chan PodResult, 1)
	// 0 timeout means infinite, we use MaxInt64 to represent it.
	var globalTimeout = 30 * time.Second

	ctx, cancel := context.WithTimeout(context.TODO(), globalTimeout)
	defer cancel()
	for _, pod := range pods {
		go func(pod *v1.Pod, returnCh chan PodResult) {
			//TODO: wait for pod evict or timeout?
			for {
				var pr = PodResult{Pod: pod.ObjectMeta}
				select {
				case <-ctx.Done():
					// return here or we'll leak a goroutine.
					pr.Err = fmt.Errorf("error when evicting pod %q: global timeout reached: %v", pod.Name, globalTimeout)
					returnCh <- pr
					return
				default:
				}
				err := drainFunc(ctx, cluster, pod.Namespace, pod, metav1.DeleteOptions{})
				if err == nil {
					pr.Err = nil
					returnCh <- pr
					return
				} else if apierrors.IsNotFound(err) {
					pr.Err = nil
					returnCh <- pr
					return
				} else if apierrors.IsTooManyRequests(err) {
					log.Errorf("error when deleting pod %q (will retry after 5s): %v", pod.Name, err)
					time.Sleep(5 * time.Second)
				} else {
					pr.Err = fmt.Errorf("error when deleting pod %q: %v", pod.Name, err)
					returnCh <- pr
					return
				}
			}
		}(pod, returnCh)
	}

	output := &imachinery.BatchOutput{}
	doneCount := 0
	output.Total = len(pods)
	for doneCount < output.Total {
		select {
		case pr := <-returnCh:
			doneCount++
			if pr.Err != nil {
				output.Fail += 1
				output.Results = append(output.Results, imachinery.SetOutput(pr.Pod, pr.Err))
			} else {
				output.Success += 1
				output.Results = append(output.Results, imachinery.SetOutput(pr.Pod, nil))
			}

		default:
		}
	}
	return output
}

func translateNodeStatusToState(s string) string {
	switch s {
	case "健康", "healthy":
		return "healthy"
	case "不健康", "unhealthy":
		return "unhealthy"
	}
	return ""
}

func filterNode(node *iapiserver.NodeInfo, fuzzy string) bool {
	if fuzzy == "" {
		return false
	}
	if filterState := translateNodeStatusToState(fuzzy); filterState != "" {
		if node.NodeStatus != filterState {
			return true
		}
	}
	if fuzzy == "可调度" || fuzzy == "不可调度" {
		if fuzzy == "可调度" && !node.Resource.Spec.Unschedulable {
			return false
		}
		if fuzzy == "不可调度" && node.Resource.Spec.Unschedulable {
			return false
		}
		//do not return,check if other filter match
	}

	return NewFieldFilter().
		AddField(node.Resource.Name).
		AddField(node.NodeAddr).
		AddField(node.Role).
		AddObjectField(node.Resource).
		Filter(fuzzy)
}

func getNodeStatus(node *corev1.Node) string {
	if node == nil {
		return ""
	}

	if isNodeReady(node) {
		return "healthy"
	}

	return "unhealthy"
}

func getNodeAddr(node *corev1.Node) string {
	if node != nil {
		for _, v := range node.Status.Addresses {
			if v.Type == "InternalIP" {
				return v.Address
			}
		}
	}

	return ""
}

func getNodeName(node *corev1.Node) string {
	if node != nil {
		return node.Name
	}
	return ""
}

func getNodeImageVersion(node *corev1.Node) (string, string) {
	if node == nil {
		return "", ""
	}
	return node.Status.NodeInfo.OSImage, node.Status.NodeInfo.ContainerRuntimeVersion
}

func getNodeRole(node *corev1.Node) string {
	if node == nil || len(node.Labels) == 0 {
		return ""
	}

	gateway, _ := node.Labels["type"]
	if gateway == iapiserver.NodeRoleGateway {
		return gateway
	}

	_, ok := node.Labels[nodeMasterLabelKey]
	if ok {
		return iapiserver.NodeRoleMaster
	}

	if _, ok := node.Labels[nodeMasterLabelKeyNew]; ok {
		return iapiserver.NodeRoleMaster
	}

	return iapiserver.NodeRoleWorker
}

func calculateByNodeRole(meta []*iapiserver.NodeInfo) iapiserver.NodeRoleCount {
	resp := iapiserver.NodeRoleCount{}
	for i := range meta {
		switch getNodeRole(meta[i].Resource) {
		case iapiserver.NodeRoleGateway:
			resp.NodeGatewayCount++
		case iapiserver.NodeRoleMaster:
			resp.NodeMasterCount++
		default:
			resp.NodeWorkCount++
		}
	}

	return resp
}

func isNodeReady(node *corev1.Node) bool {
	if node != nil {
		for _, v := range node.Status.Conditions {
			if v.Type == nodeStatusReady && v.Status == "True" {
				return true
			}
		}
	}

	return false
}
