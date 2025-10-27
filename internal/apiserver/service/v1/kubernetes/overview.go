package kubernetes

import (
	"context"
	"strings"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/internal/pkg/clientset"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/mathutil"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *kubernetesService) getClusterNodesOverview(ctx context.Context, clusterInfo *iapiserver.Cluster) (map[string]*iapiserver.NodeOverview, error) {
	resp := map[string]*iapiserver.NodeOverview{}

	nodeList, err := clientset.NodeList(ctx, clusterInfo, metav1.ListOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rets, err := k.MonitorNodeStateList(ctx, clusterInfo)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, v := range nodeList.Items {
		no := &iapiserver.NodeOverview{NodeName: v.Name}
		resp[no.NodeName] = no
		no.NodeAddr = getNodeAddr(&v)
		no.Role = getNodeRole(&v)
		no.Status = nodeStatusHealth
		no.OsImage, no.ContainerRuntimeVersion = getNodeImageVersion(&v)
		if !isNodeReady(&v) {
			no.Status = nodeStatusUnhealth
		}
		if no.Role == nodeRoleGateway {
			no.IsGateway = true
		}

		nodeState, ok := rets[no.NodeName]
		if !ok {
			continue
		}
		no.ClusterUUID = clusterInfo.ID
		no.ClusterName = clusterInfo.Name
		no.CreateTime = v.CreationTimestamp.Unix()
		no.CpuUsed = nodeState.CpuUsed
		no.CpuCapacity = nodeState.CpuCapacity
		no.CpuUsedRatio = nodeState.CpuUsedRatio
		no.CpuResourceRequest = nodeState.CpuResourceRequest
		no.CpuResourceLimit = nodeState.CpuResourceLimit
		no.CpuResourceLimitRatio = mathutil.Divide(float64(no.CpuResourceLimit), no.CpuCapacity) * 100
		no.CpuResourceRequestRatio = mathutil.Divide(float64(no.CpuResourceRequest), no.CpuCapacity) * 100
		no.MemoryUsed = nodeState.MemoryUsed
		no.MemoryCapacity = nodeState.MemoryCapacity
		no.MemoryUsedRatio = nodeState.MemoryUsedRatio
		no.MemoryResourceRequest = nodeState.MemoryResourceRequest
		no.MemoryResourceLimit = nodeState.MemoryResourceLimit
		no.MemoryResourceLimitRatio = mathutil.Divide(float64(no.MemoryResourceLimit), no.MemoryCapacity) * 100
		no.MemoryResourceRequestRatio = mathutil.Divide(float64(no.MemoryResourceRequest), no.MemoryCapacity) * 100
		no.DiskUsedRatio = nodeState.DiskUsedRatio * 100
		no.DiskUsed = nodeState.DiskUsed
		no.DiskCapacity = nodeState.DiskCapacity
		no.PodUsed = nodeState.PodUsed
		no.PodCapacity = nodeState.PodCapacity
		no.PodUsedRatio = nodeState.PodUsedRatio * 100
	}

	return resp, nil
}

func KVStringToSelector(kvs []string) string {
	sep := ","
	equal := "="
	resp := ""

	count := len(kvs) / 2
	for i := 0; i < count; i++ {
		resp += kvs[i*2] + equal + kvs[i*2+1] + sep
	}

	return strings.Trim(resp, sep)
}

func getEtcdInfo(ctx context.Context, clusterInfo *iapiserver.Cluster) (*iapiserver.ComponentInfo, error) {
	resp := &iapiserver.ComponentInfo{Namespace: iapiserver.NamespaceKubeSystem, Name: componentEtcd, Status: "unhealth"}

	podList, err := selectorToPodList(ctx, clusterInfo, iapiserver.NamespaceKubeSystem, KVStringToSelector([]string{"component", "etcd"}))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(podList) == 0 {
		return nil, errors.Errorf("etcd not find")
	}

	half := 0
	for _, v := range podList {
		if getPodStatus(&v) == podStatusRunning {
			half++
			if half > len(podList)/2 {
				resp.Status = podStatusRunning
				break
			}
		}
	}

	return resp, nil
}

func getSchedulerInfo(ctx context.Context, clusterInfo *iapiserver.Cluster) (*iapiserver.ComponentInfo, error) {
	resp := &iapiserver.ComponentInfo{}

	pod, err := selectorToPod(ctx, clusterInfo, iapiserver.NamespaceKubeSystem, KVStringToSelector([]string{"component", "kube-scheduler"}))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.Namespace = pod.Namespace
	resp.Name = pod.Name
	resp.Status = getPodStatus(pod)
	return resp, nil
}

func getControllerManagerInfo(ctx context.Context, clusterInfo *iapiserver.Cluster) (*iapiserver.ComponentInfo, error) {
	resp := &iapiserver.ComponentInfo{}

	pod, err := selectorToPod(ctx, clusterInfo, iapiserver.NamespaceKubeSystem, KVStringToSelector([]string{"component", "kube-controller-manager"}))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.Namespace = pod.Namespace
	resp.Name = pod.Name
	resp.Status = getPodStatus(pod)
	return resp, nil
}

func getApiserverInfo(ctx context.Context, clusterInfo *iapiserver.Cluster) (*iapiserver.ComponentInfo, error) {
	resp := &iapiserver.ComponentInfo{}

	pod, err := selectorToPod(ctx, clusterInfo, iapiserver.NamespaceKubeSystem, KVStringToSelector([]string{"component", "kube-apiserver"}))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.Namespace = pod.Namespace
	resp.Name = pod.Name
	resp.Status = getPodStatus(pod)
	return resp, nil
}

// 需要保证selector能够唯一确定一个pod，当有多个pod时这里只会返回第一pod
func selectorToPod(ctx context.Context, clusterInfo *iapiserver.Cluster, namespace string, selector string) (*corev1.Pod, error) {
	podList, err := selectorToPodList(ctx, clusterInfo, namespace, selector)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(podList) == 0 {
		return nil, errors.Errorf("namespace[%v],selector[%v] pod list is empty", namespace, selector)
	}

	return &podList[0], nil
}

// 需要保证selector能够唯一确定一个pod，当有多个pod时这里只会返回第一pod
func selectorToPodList(ctx context.Context, clusterInfo *iapiserver.Cluster, namespace string, selector string) ([]corev1.Pod, error) {
	if clusterInfo == nil || namespace == "" || selector == "" {
		return nil, errors.Errorf("namespace[%v],selector[%v] is empty", namespace, selector)
	}

	listOption := metav1.ListOptions{}
	listOption.LabelSelector = selector
	podList, err := clientset.PodList(ctx, clusterInfo, namespace, listOption)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return podList.Items, nil
}

const (
	// status
	unknown = "unknown"
)

var (
	podStatusRunning = "Running"
	nodeRoleMaster   = "master"
	nodeRoleWorker   = "worker"
	nodeRoleGateway  = "gateway"
)

func getPodStatus(pod *v1.Pod) string {
	if pod == nil {
		return unknown
	}
	return string(pod.Status.Phase)
}
