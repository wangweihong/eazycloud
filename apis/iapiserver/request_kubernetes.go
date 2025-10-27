package iapiserver

import (
	"fmt"
	"strconv"

	"strings"

	jsoniter "github.com/json-iterator/go"
	snapshotv1beta1 "github.com/kubernetes-csi/external-snapshotter/client/v3/apis/volumesnapshot/v1beta1"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	pkgtypes "k8s.io/apimachinery/pkg/types"

	"github.com/wangweihong/eazycloud/apis/ikubernetes"
	"github.com/wangweihong/eazycloud/apis/imachinery"
)

const (
	StorageClassDriverNfs       = "nfs.topke.io"
	StorageClassDriverGlusterFs = "kubernetes.io/glusterfs"
)

type (
	ResourceRequest struct {
		Cluster          string               `json:"cluster"           binding:"required"`
		CreateOpts       metav1.CreateOptions `json:"create_opts"`
		GetOpts          metav1.GetOptions    `json:"get_opts"`
		DeleteOpts       metav1.DeleteOptions `json:"delete_opts"`
		UpdateOpts       metav1.UpdateOptions `json:"update_opts"`
		ListOpts         metav1.ListOptions   `json:"list_opts"`
		PatchOpts        metav1.PatchOptions  `json:"patch_opts"`
		Pt               pkgtypes.PatchType   `json:"pt"`
		Data             []byte               `json:"data"`
		SubResources     []string             `json:"sub_resources"`
		DeleteCollection bool                 `json:"delete_collection"`

		Yaml bool `json:"yaml" form:"yaml"`
	}

	ResourceListRequest struct {
		imachinery.PagingParams
		Cluster   string `json:"cluster"   form:"cluster"`
		Namespace string `json:"namespace" form:"namespace"`
		Fuzzy     string `json:"fuzzy"     form:"fuzzy"`
		SortBy    string `json:"sort_by"   form:"sort_by"`
		SortDesc  bool   `json:"sort_desc" form:"sort_desc"`
		Yaml      bool   `json:"yaml"      form:"yaml"`

		FieldSelector string `json:"field_selector" form:"field_selector"`
		LabelSelector string `json:"label_selector" form:"label_selector"`
	}

	ResourceInfo[T any] struct {
		Resource T        `json:"resource"`
		Cluster  *Cluster `json:"cluster"`
	}
)

func NewResourceInfo[T any](resource T, cluster *Cluster) *ResourceInfo[T] {
	return &ResourceInfo[T]{
		Resource: resource,
		Cluster:  cluster,
	}
}

// sel := labels.NewSelector()
// req, err := labels.NewRequirement("mykey", selection.Exists, []string{})
// if err != nil {
// ....
// }
// sel.Add(*req)
//
// deployOpts := []client.ListOption{
// client.MatchingLabelsSelector{Selector: sel},
// ...
// }
//
// b.cl.List(context.Background(), deployments, deployOpts...)
func (r ResourceListRequest) ToListOpts() metav1.ListOptions {
	opt := metav1.ListOptions{
		FieldSelector: r.FieldSelector,
		LabelSelector: r.LabelSelector,
	}

	if r.LabelSelector != "" {
		rs := strings.Split(r.LabelSelector, "=")
		//only key without value
		if len(rs) == 1 {
			sel := labels.NewSelector()
			if req, err := labels.NewRequirement(r.LabelSelector, selection.Exists, []string{}); err == nil {
				opt.LabelSelector = sel.Add(*req).String()
			}
		}
	}

	return opt
}

type ResourceGetRequest struct {
	Cluster   string `json:"cluster"   form:"cluster" binding:"required"`
	Namespace string `json:"namespace" form:"namespace"`
	Name      string `json:"name"      form:"name" binding:"required"`
	Yaml      bool   `json:"yaml"      form:"yaml"`
}

func (r ResourceGetRequest) ToGetOpts() metav1.GetOptions {
	return metav1.GetOptions{}
}

const (

	// common sort type
	KubernetesResourceSortByCreateTime  = "KubernetesResourceSortByCreateTime"
	KubernetesResourceSortByName        = "KubernetesResourceSortByName"
	KubernetesResourceSortByNamespace   = "KubernetesResourceSortByNamespace"
	KubernetesResourceSortByClusterName = "KubernetesResourceSortByClusterName"
	KubernetesResourceSortByCpuUsed     = "KubernetesResourceSortByCpuUsed"

	// pod special sort type
	KubernetesResourcePodSortByState         = "KubernetesResourcePodSortByState"
	KubernetesResourcePodSortByRestart       = "KubernetesResourcePodSortByRestart"
	KubernetesResourcePodSortByHostIP        = "KubernetesResourcePodSortByHostIP"
	KubernetesResourcePodSortByCpuRequest    = "KubernetesResourcePodSortByCpuRequest"
	KubernetesResourcePodSortByMemoryRequest = "KubernetesResourcePodSortByMemoryRequest"
	KubernetesResourcePodSortByPodIP         = "KubernetesResourcePodSortByPodIP"
	//only work in Overview
	KubernetesResourcePodSortByPhysicalCpu          = "KubernetesResourcePodSortByPhysicalCpu"
	KubernetesResourcePodSortByPhysicalMemory       = "KubernetesResourcePodSortByPhysicalMemory"
	KubernetesResourcePodSortByNetworkTransferRatio = "KubernetesResourcePodSortByNetworkTransferRatio"
	KubernetesResourcePodSortByNetworkReceiveRatio  = "KubernetesResourcePodSortByNetworkReceiveRatio"
	// deployment special sort type
	KubernetesResourceDeploymentSortByReplicas = "KubernetesResourceDeploymentSortByReplicas"
	// statefulSet special sort type
	KubernetesResourceStatefulSetSortByReplicas = "KubernetesResourceStatefulSetSortByReplicas"
	// job special sort type
	KubernetesResourceJobSortByStartTime = "KubernetesResourceJobSortByStartTime"
	KubernetesResourceJobSortByEndTime   = "KubernetesResourceJobSortByEndTime"
	KubernetesResourceJobSortByDuration  = "KubernetesResourceJobSortByDuration"
	//cronjob special sort type
	KubernetesResourceCronJobSortByLastJobTime = "KubernetesResourceCronJobSortByLastJobTime"
	KubernetesResourceCronJobSortByJobNum      = "KubernetesResourceCronJobSortByJobNum"
	//service special sort type
	KubernetesResourceServiceSortByType = "KubernetesResourceServiceSortByType"
	//
	KubernetesResourceSecretSortByType = "KubernetesResourceSecretSortByType"
	//namespace special sort type
	KubernetesResourceNodeSortByState = "KubernetesResourceNodeSortByState"
	//cluster
	KubernetesResourceClusterSortByState = "ClusterSortByState"
	//node
	KubernetesResourceNodeSortBYCpuResourceLimitRatio      = "CpuResourceLimitRatio"
	KubernetesResourceNodeSortBYCpuResourceRequestRatio    = "CpuResourceRequestRatio"
	KubernetesResourceNodeSortBYMemoryResourceRequestRatio = "MemoryResourceRequestRatio"
	KubernetesResourceNodeSortBYMemoryResourceLimitRatio   = "MemoryResourceLimitRatio"
	KubernetesResourceNodeSortBYPodUsedRatio               = "PodUsedRatio"
)

type EachResourceRangeListState[T any] struct {
	ClusterUUID string `json:"cluster_uuid"`
	ClusterName string `json:"cluster_name"`
	Result      error  `json:"result"`
	TotalCount  int    `json:"total_count"`
	List        []T    `json:"list,omitempty"`
}

func NewEachResourceRangeListState[T any](uuid, name string) EachResourceRangeListState[T] {
	return EachResourceRangeListState[T]{
		ClusterUUID: uuid,
		ClusterName: name,
	}
}

type NodeOverview struct {
	// 节点名
	NodeName string `json:"node_name"`
	// 节点地址
	NodeAddr string `json:"node_addr"`
	// 总CPU
	CpuCapacity float64 `json:"cpu_capacity"`
	// 已使用CPU
	CpuUsed float64 `json:"cpu_used"`
	// CPU使用率
	CpuUsedRatio float64 `json:"cpu_used_ratio"`
	// CPU资源请求
	CpuResourceRequest int64 `json:"cpu_resource_request"`
	// CPU资源请求率
	CpuResourceRequestRatio float64 `json:"cpu_resource_request_ratio"`
	// CPU资源限制
	CpuResourceLimit int64 `json:"cpu_resource_limit"`
	// CPU资源限制率ListRequestG
	CpuResourceLimitRatio float64 `json:"cpu_resource_limit_ratio"`
	// 总内存
	MemoryCapacity float64 `json:"memory_capacity"`
	// 已使用内存
	MemoryUsed float64 `json:"memory_used"`
	// 内存使用率
	MemoryUsedRatio float64 `json:"memory_used_ratio"`
	// 内存资源请求
	MemoryResourceRequest int64 `json:"memory_resource_request"`
	// 内存资源请求率
	MemoryResourceRequestRatio float64 `json:"memory_resource_request_ratio"`
	// 内存资源限制
	MemoryResourceLimit int64 `json:"memory_resource_limit"`
	// 内存资源限制
	MemoryResourceLimitRatio float64 `json:"memory_resource_limit_ratio"`
	// 总存储容量
	DiskCapacity float64 `json:"disk_capacity"`
	// 已使用存储容量
	DiskUsed float64 `json:"disk_used"`
	// 存储容量使用率
	DiskUsedRatio float64 `json:"disk_used_ratio"`
	// Pod总数
	PodCapacity int `json:"pod_capacity"`
	// 已使用Pod总数
	PodUsed int `json:"pod_used"`
	// Pod使用率
	PodUsedRatio float64 `json:"pod_used_ratio"`
	// 是否网关节点
	IsGateway bool `json:"is_gateway"`
	// 运行状态
	Status string `json:"status"`
	// 节点角色
	Role string `json:"role"`
	// 系统镜像
	OsImage string `json:"os_image"`
	// 容器运行时
	ContainerRuntimeVersion string `json:"container_runtime_version"`
	// 错误信息
	ErrorMsg string `json:"error_msg"`
	// 创建时间
	CreateTime int64             `json:"create_time"`
	Extra      map[string]string `json:"extra"`

	ClusterName string `json:"cluster_name"`
	ClusterUUID string `json:"cluster_uuid"`
}

// cluster service request
type (
	ServiceRequest struct {
		ResourceRequest
		Resource *v1.Service `json:"service" binding:"required,namespaced"`
	}

	ServiceListRequest struct {
		ResourceListRequest
	}

	ServiceGetRequest struct {
		ResourceGetRequest
	}
	// ServiceInfo struct {
	// 	Resource *v1.Service `json:"resource"`
	// }
	ServiceResponse struct {
		Info *ResourceInfo[*v1.Service] `json:"info"`
	}

	ServiceListResponse struct {
		TotalCount         int                                                      `json:"total_count"`
		List               []*ResourceInfo[*v1.Service]                             `json:"list"`
		EachRangeListState []EachResourceRangeListState[*ResourceInfo[*v1.Service]] `json:"each_range_list_state"`
	}

	ServiceBatchRequest struct {
		Resources []*ServiceRequest `json:"resources" binding:"required,dive"`
	}
)

// cluster endpoints request
type EndpointsRequest struct {
	ResourceRequest
	FilterName string        `json:"filter_name" form:"filter_name"`
	Endpoints  *v1.Endpoints `json:"endpoints"` // endpoints name is equal to service name
}

type EndpointsGetRequest struct {
	ResourceGetRequest
}

// cluster endpoints request
type EndpointListRequest struct {
	ResourceListRequest
}

type EndpointGetRequest struct {
	ResourceGetRequest
	FilterName string `json:"filter_name" form:"filter_name"`
}

// cluster endpoints info
type EndpointsInfo struct {
	Resource *v1.Endpoints `json:"resource"`
}

// cluster endpoints response
type EndpointsResponse struct {
	Info *EndpointsInfo `json:"info,omitempty"`
}

// cluster endpoints response
type EndpointsListResponse struct {
	TotalCount         int                                          `json:"total_count"`
	List               []*EndpointsInfo                             `json:"list"`
	EachRangeListState []EachResourceRangeListState[*EndpointsInfo] `json:"each_range_list_state,omitempty"`
}

const (
	NodeRoleMaster  = "master"
	NodeRoleGateway = "gateway"
	NodeRoleWorker  = "worker"

	IngressControllerNamespace = "system-router"
	IngressControllerPrefix    = "router-"
)

// node request
type (
	NodeRequest struct {
		ResourceRequest
		// 动作,  1，设置为网关虚拟机 2，取消网关虚拟机
		Action string `json:"action"`
		// 过滤名字
		FilterName string `json:"filter_name"`
		// 过滤状态
		FilterStatus string `json:"filter_status"`
		// 过滤角色
		FilterRole string `json:"filter_role"`
		// 节点
		Node *v1.Node `json:"node"`
		// 显示使用率
		ShowResourceUsage bool `json:"show_resource_usage"`
		// 节点IP
		NodeIP string `json:"node_ip"`
	}

	NodeEventRequest struct {
		ResourceGetRequest
		imachinery.PagingParams
		Fuzzy string `json:"fuzzy" form:"fuzzy"`
	}

	NodeListRequest struct {
		ResourceListRequest
		ShowResourceUsage bool `json:"show_resource_usage" form:"show_resource_usage"`
	}

	NodeGetRequest struct {
		ResourceGetRequest
		ShowResourceUsage bool `json:"show_resource_usage" form:"show_resource_usage"`
	}

	NodeInfo struct {
		Resource      *v1.Node              `json:"resource"`
		ResourceUsage *NodeRealtimeResource `json:"resource_usage,omitempty"`
		NodeStatus    string                `json:"node_status,omitempty"`
		NodeAddr      string                `json:"node_addr,omitempty"`
		Role          string                `json:"role,omitempty"`
	}

	// node response
	NodeResponse struct {
		Info *NodeInfo `json:"info"`
	}

	NodeListResponse struct {
		NodeRoleCount
		TotalCount         int                                     `json:"total_count"`
		List               []*NodeInfo                             `json:"list"`
		EachRangeListState []EachResourceRangeListState[*NodeInfo] `json:"each_range_list_state,omitempty"`
	}

	NodeRoleCount struct {
		NodeMasterCount  int `json:"node_master"`
		NodeWorkCount    int `json:"node_work"`
		NodeGatewayCount int `json:"node_gateway"`
	}

	// node request
	NodeDrainRequest struct {
		NodeRequest
		DisableEviction bool `json:"disable_eviction"` // if disable use delete instead of evict
	}

	// node request
	NodeDrainResponse struct {
		//*BatchOutput
	}
	NodeRealtimeResource struct {
		*NodeRealTimePhysicalResourceUsage
		*NodeRealTimePodUsage
		*NodePodListResourceRequestAndLimit
	}

	NodeRealTimePhysicalResourceUsage struct {
		CpuCapacity     float64 `json:"cpu_capacity"` // m unit(1core=1000m)
		CpuUsed         float64 `json:"cpu_used"`
		CpuUsedRatio    float64 `json:"cpu_used_ratio"`
		MemoryCapacity  float64 `json:"memory_capacity"` //bytes
		MemoryUsed      float64 `json:"memory_used"`     //bytes
		MemoryUsedRatio float64 `json:"memory_used_ratio"`
		DiskCapacity    float64 `json:"disk_capacity"` //bytes
		DiskUsed        float64 `json:"disk_used"`     //bytes
		DiskUsedRatio   float64 `json:"disk_used_ratio"`
	}

	NodeRealTimePodUsage struct {
		PodCapacity  int     `json:"pod_capacity"`
		PodUsed      int     `json:"pod_used"`
		PodUsedRatio float64 `json:"pod_used_ratio"`
	}
)

type NodePodListResourceRequestAndLimit struct {
	NodeCpuCapacity    float64 `json:"node_cpu_capacity"`
	NodeMemoryCapacity float64 `json:"node_memory_capacity"`

	*PodListResourceRequest
	*PodListResourceLimit
}

type PodListResourceRequest struct {
	PodResourceRequest
	CpuResourceRequestRatio    float64 `json:"cpu_resource_request_ratio"`
	MemoryResourceRequestRatio float64 `json:"memory_resource_request_ratio"`
}

type PodListResourceLimit struct {
	PodResourceLimit
	CpuResourceLimitRatio    float64 `json:"cpu_resource_limit_ratio"`
	MemoryResourceLimitRatio float64 `json:"memory_resource_limit_ratio"`
}

type PodResourceRequest struct {
	CpuResourceRequest    int64 `json:"cpu_resource_request"`
	MemoryResourceRequest int64 `json:"memory_resource_request"`
}

type PodResourceLimit struct {
	CpuResourceLimit    int64 `json:"cpu_resource_limit"`
	MemoryResourceLimit int64 `json:"memory_resource_limit"`
}

type (
	SecretRequest struct {
		ResourceRequest
		Resource *v1.Secret `json:"resource" binding:"required,namespaced"`
	}

	SecretGetRequest struct {
		ResourceGetRequest
	}

	SecretListRequest struct {
		ResourceListRequest
	}
	SecretListResponse struct {
		TotalCount         int                                                     `json:"total_count"`
		List               []*ResourceInfo[*v1.Secret]                             `json:"list"`
		EachRangeListState []EachResourceRangeListState[*ResourceInfo[*v1.Secret]] `json:"each_range_list_state"`
	}
	SecretBatchRequest struct {
		Resources []*SecretRequest `json:"resources" binding:"dive"`
	}
)

type RouterRequest struct {
	ResourceRequest
	// 容器镜像
	Image     string `json:"image"`
	Namespace string `json:"namespace" binding:"required"`
}

type RouterGetRequest struct {
	ResourceGetRequest
}

var (
	KubernetesSystemNamespaces = []string{"kube-system"}
)

type (
	ResourceQuotaRequest struct {
		ResourceRequest
		ResourceQuota  *v1.ResourceQuota `json:"resource_quota" binding:"required,namespaced"`
		UpdateIfExists bool              `json:"update_if_exists" description:"存在则更新"`
	}

	ResourceQuotaListRequest struct {
		ResourceListRequest
	}

	ResourceQuotaGetRequest struct {
		ResourceGetRequest
	}

	ResourceQuotaResponse struct {
		Info *ResourceQuotaInfo `json:"info"`
	}

	ResourceQuotaListResponse struct {
		TotalCount         int                                              `json:"total_count"`
		List               []*ResourceQuotaInfo                             `json:"list"`
		EachRangeListState []EachResourceRangeListState[*ResourceQuotaInfo] `json:"each_range_list_state,omitempty"`
	}

	ResourceQuotaInfo struct {
		Resource *v1.ResourceQuota     `json:"resource"`
		Convert  *ResourceQuotaConvert `json:"convert,omitempty"`
	}

	ResourceQuotaConvert struct {
		Spec   ResourceQuotaSpec   `json:"spec,omitempty"`
		Status ResourceQuotaStatus `json:"status,omitempty"`
	}

	ResourceQuotaSpec struct {
		Hard map[string]int64 `json:"hard,omitempty"`
	}

	ResourceQuotaStatus struct {
		Hard map[string]int64 `json:"hard,omitempty"`
		Used map[string]int64 `json:"used,omitempty"`
	}
)
type (
	JobRequest struct {
		Resource *batchv1.Job `json:"resource" binding:"required,namespaced"`
		ResourceRequest
	}

	JobListRequest struct {
		ResourceListRequest
	}

	JobGetRequest struct {
		ResourceGetRequest
	}

	JobBatchRequest struct {
		Resources []*JobRequest `json:"resources" binding:"dive"`
	}

	JobResponse struct {
		Info *JobInfo `json:"info"`
	}

	JobListResponse struct {
		TotalCount         int                                    `json:"total_count"`
		List               []*JobInfo                             `json:"list"`
		EachRangeListState []EachResourceRangeListState[*JobInfo] `json:"each_range_list_state,omitempty"`
	}

	JobInfo struct {
		Resource        *batchv1.Job       `json:"resource"`
		ResourceConvert []*ResourceConvert `json:"resource_convert,omitempty"`
	}
)

type MonitorData struct {
	Namespace    string            `json:"namespace"`
	Name         string            `json:"name"`
	Addr         string            `json:"addr"`
	Type         string            `json:"type"`
	Cpu          float64           `json:"cpu"`
	CpuUsed      float64           `json:"cpu_used"`
	CpuTotal     float64           `json:"cpu_total"`
	Memory       float64           `json:"memory"`
	MemoryUsed   float64           `json:"memory_used"`
	MemoryTotal  float64           `json:"memory_total"`
	DiskRatio    float64           `json:"disk_ratio"`
	DiskUsed     float64           `json:"disk_used"`
	DiskTotal    float64           `json:"disk_total"`
	DiskRead     float64           `json:"disk_read"`
	DiskWrite    float64           `json:"disk_write"`
	PodRatio     float64           `json:"pod_ratio"`
	PodCount     float64           `json:"pod_count"` // 当前有多少的pod
	PodTotal     float64           `json:"pod_total"` // 集群总共可以创建多少个pod
	NetworkRead  float64           `json:"network_read"`
	NetworkWrite float64           `json:"network_write"`
	Extra        map[string]string `json:"extra"`
}

type MonitorDataResponse struct {
	Status       error             `json:"status"`
	Namespace    string            `json:"namespace,omitempty"`
	Name         string            `json:"name"`
	Addr         string            `json:"addr"`
	Type         string            `json:"type"`
	Cpu          float64           `json:"cpu"`
	CpuUsed      float64           `json:"cpu_used"`
	CpuTotal     float64           `json:"cpu_total"`
	Memory       float64           `json:"memory"`
	MemoryUsed   float64           `json:"memory_used"`
	MemoryTotal  float64           `json:"memory_total"`
	DiskRatio    float64           `json:"disk_ratio"`
	DiskUsed     float64           `json:"disk_used"`
	DiskTotal    float64           `json:"disk_total"`
	DiskRead     float64           `json:"disk_read"`
	DiskWrite    float64           `json:"disk_write"`
	PodRatio     float64           `json:"pod_ratio"`
	PodCount     float64           `json:"pod_count"`
	PodTotal     float64           `json:"pod_total"`
	NetworkRead  float64           `json:"network_read"`
	NetworkWrite float64           `json:"network_write"`
	Extra        map[string]string `json:"extra"`
}

const (
	PodSortByState         = "PodSortByState"
	PodSortByRestart       = "PodSortByRestart"
	PodSortByHostIP        = "PodSortByHostIP"
	PodSortByCpuRequest    = "PodSortByCpuRequest"
	PodSortByMemoryRequest = "PodSortByMemoryRequest"
	PodSortByPodIP         = "PodSortByPodIP"
	//only work in Overview
	PodSortByPhysicalCpu          = "PodSortByPhysicalCpu"
	PodSortByPhysicalMemory       = "PodSortByPhysicalMemory"
	PodSortByNetworkTransferRatio = "PodSortByNetworkTransferRatio"
	PodSortByNetworkReceiveRatio  = "PodSortByNetworkReceiveRatio"
)

const (
	ComponentKubectl = "kubectl"
)

type (
	PodRequest struct {
		ResourceRequest
		// MonitoringTimeRange
		// ShowMonitorData bool             `json:"show_monitor_data" form:"show_monitor_data"`
		PodLogOptions v1.PodLogOptions `json:"pod_log_options"`
		Resource      *v1.Pod          `json:"resource"        binding:"required,namespaced"`
	}

	PodListRequest struct {
		ResourceListRequest
		MonitoringTimeRange
		ShowMonitorData  bool   `json:"show_monitor_data"  form:"show_monitor_data"`
		FilterName       string `json:"filter_name"        form:"filter_name"        description:"过滤名"`
		FilterHostIp     string `json:"filter_host_ip"     form:"filter_host_ip"` // 字段选择器不支持hostip
		FilterVolumeName string `json:"filter_volume_name" form:"filter_volume_name"`
		FilterConfigMap  string `json:"filter_config_map"  form:"filter_config_map"`
		FilterSecret     string `json:"filter_secret"      form:"filter_secret"`
		Component        string `json:"component"          form:"component"` // return specific component pod info, such as `kubectl`
	}

	PodLogRequest struct {
		ResourceGetRequest           `binding:"namespaced"`
		Container                    string `json:"container"                    form:"container"`
		TailLines                    int64  `json:"tail_lines"                   form:"tail_lines"`
		Follow                       bool   `json:"follow"                       form:"follow"`
		TimeStamps                   bool   `json:"timestamps"                   form:"timestamps"`
		InsecureSkipTLSVerifyBackend bool   `json:"insecureSkipTLSVerifyBackend" form:"insecureSkipTLSVerifyBackend"`
		SinceSeconds                 int64  `json:"sinceSeconds"                 form:"sinceSeconds"`
	}

	PodGetRequest struct {
		ResourceGetRequest
	}

	PodBatchRequest struct {
		Resources []*PodRequest `json:"resources" binding:"dive"`
	}

	ObjectTypeMeta struct {
		ObjectMeta metav1.ObjectMeta `json:"object_meta"`
		TypeMeta   metav1.TypeMeta   `json:"type_meta"`
	}

	PodInfo struct {
		Resource        *v1.Pod             `json:"resource"`
		Metric          []Metric            `json:"metric,omitempty"`
		ResourceRequest v1.ResourceList     `json:"resource_request,omitempty"`
		ResourceLimit   v1.ResourceList     `json:"resource_limit,omitempty"`
		PodStatus       *PodStatus          `json:"pod_status,omitempty"`
		MaxRestarts     *int                `json:"max_restarts,omitempty"`
		Controller      *ObjectTypeMeta     `json:"controller,omitempty"`
		NodeInfo        *v1.ObjectReference `json:"node_info,omitempty"`
		ResourceConvert []*ResourceConvert  `json:"resource_convert,omitempty"`
		Cluster         *Cluster            `json:"cluster"`
	}
)

func (r PodLogRequest) ToLogOption() v1.PodLogOptions {
	opt := v1.PodLogOptions{
		Container:                    r.Container,
		Follow:                       r.Follow,
		Timestamps:                   r.TimeStamps,
		InsecureSkipTLSVerifyBackend: r.InsecureSkipTLSVerifyBackend,
	}

	if r.SinceSeconds != 0 {
		opt.SinceSeconds = &r.SinceSeconds
	}

	if r.TailLines != 0 {
		opt.TailLines = &r.TailLines
	}

	return opt
}

type ResourceConvert struct {
	Container string `json:"container"`
	// 资源限制数
	Limits map[string]int64 `json:"limits,omitempty"`
	// 资源请求数
	Request map[string]int64 `json:"requests,omitempty"`
}

const (
	PodStatusScheduleFail            = "ScheduleFailed"
	PodInitContainerFailStatusPrefix = "Init:"
	PodStatusTerminating             = "Terminating"
)

type PodStatus struct {
	Status          string `json:"status"`
	Message         string `json:"message"`
	Reason          string `json:"reason"`
	ReadyContainers int    `json:"ready_containers"`
	TotalContainers int    `json:"total_containers"`
}

type PodGetComponentResponse struct {
	TotalCount         int                                    `json:"total_count"`
	List               []*PodInfo                             `json:"list"`
	EachRangeListState []EachResourceRangeListState[*PodInfo] `json:"each_range_list_state,omitempty"`
}

type PodListResponse struct {
	TotalCount         int                                    `json:"total_count"`
	List               []*PodInfo                             `json:"list"`
	EachRangeListState []EachResourceRangeListState[*PodInfo] `json:"each_range_list_state,omitempty"`
}

type (
	PodLogResponse struct {
		LogList    []ikubernetes.PodLog `json:"log_list"`
		TotalCount int                  `json:"total_count"`
		EndTime    int64                `json:"end_time"`
	}
)

type Metadata struct {
	Metric string `json:"metric"`
	Type   string `json:"type"`
	Help   string `json:"help"`
}

type Metric struct {
	MetricData
	MetricName string `json:"metric_name"`
	Error      string `json:"error"`
}

type MetricOne struct {
	MetricName string  `json:"metric_name"`
	Series     []Point `json:"series"`
	Sample     *Point  `json:"sample"`
	Error      string  `json:"error"`
}

type MetricData struct {
	MetricType   string        `json:"metric_type"`
	MetricValues []MetricValue `json:"metric_values"`
}

type MetricValue struct {
	Name     string            `json:"name"`
	Max      float64           `json:"max"`
	Metadata map[string]string `json:"metadata"`
	Sample   *Point            `json:"sample"`
	Series   []Point           `json:"series"`
}
type Point [2]float64

func (p Point) Timestamp() float64 {
	return p[0]
}

func (p Point) Value() float64 {
	return p[1]
}

func (p Point) MarshalJSON() ([]byte, error) {
	t, err := jsoniter.Marshal(p.Timestamp())
	if err != nil {
		return nil, err
	}
	v, err := jsoniter.Marshal(strconv.FormatFloat(p.Value(), 'f', -1, 64))
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("[%s,%s]", t, v)), nil
}

func (p *Point) UnmarshalJSON(b []byte) error {
	var v []interface{}
	if err := jsoniter.Unmarshal(b, &v); err != nil {
		return err
	}

	if v == nil {
		return nil
	}

	if len(v) != 2 {
		return errors.Errorf("unsupported array length")
	}

	ts, ok := v[0].(float64)
	if !ok {
		return errors.Errorf("failed to unmarshal [timestamp]")
	}
	valstr, ok := v[1].(string)
	if !ok {
		return errors.Errorf("failed to unmarshal [value]")
	}
	valf, err := strconv.ParseFloat(valstr, 64)
	if err != nil {
		return err
	}

	p[0] = ts
	p[1] = valf
	return nil
}

type (
	ConfigMapRequest struct {
		ResourceRequest
		Resource *v1.ConfigMap `json:"resource" binding:"required,namespaced"`
	}

	ConfigMapListRequest struct {
		ResourceListRequest
	}

	ConfigMapListResponse struct {
		TotalCount         int
		List               []*ConfigMapInfo
		EachRangeListState []EachResourceRangeListState[*ConfigMapInfo]
	}

	ConfigMapGetRequest struct {
		ResourceGetRequest `binding:"namespace_resource"`
	}

	ConfigMapBatchRequest struct {
		Resources []*ConfigMapRequest `json:"resources" binding:"dive"`
	}

	ConfigMapInfo struct {
		Resource *v1.ConfigMap `json:"resource"`
	}
)

type (
	CronJobRequest struct {
		ResourceRequest
		Resource *batchv1.CronJob `json:"resource"`
	}

	CronJobListRequest struct {
		ResourceListRequest
	}

	CronJobGetRequest struct {
		ResourceGetRequest `binding:"namespaced"`
	}

	CronJobBatchRequest struct {
		Resources []*CronJobRequest `json:"resources" binding:"dive"`
	}

	CronJobResponse struct {
		Info *CronJobInfo `json:"info,omitempty"`
	}

	CronJobListResponse struct {
		Status             error                                      `json:"status"`
		TotalCount         int                                        `json:"total_count"`
		List               []*CronJobInfo                             `json:"list"`
		EachRangeListState []EachResourceRangeListState[*CronJobInfo] `json:"each_range_list_state,omitempty"`
	}

	CronJobInfo struct {
		Resource        *batchv1.CronJob   `json:"resource"`
		ResourceConvert []*ResourceConvert `json:"resource_convert,omitempty"`
	}
)
type (
	NetworkPolicyRequest struct {
		ResourceRequest
		Resource *networkingv1.NetworkPolicy `json:"resource" binding:"required,namespaced"`
	}

	NetworkPolicyListRequest struct {
		ResourceListRequest
	}

	NetworkPolicyGetRequest struct {
		ResourceGetRequest
	}

	NetworkPolicyBatchRequest struct {
		Resources []*NetworkPolicyRequest `json:"resources" binding:"dive"`
	}

	NetworkPolicyResponse struct {
		Info *NetworkPolicyInfo `json:"info"`
	}

	NetworkPolicyListResponse struct {
		TotalCount         int                                              `json:"total_count"`
		List               []*NetworkPolicyInfo                             `json:"list"`
		EachRangeListState []EachResourceRangeListState[*NetworkPolicyInfo] `json:"each_range_list_state,omitempty"`
	}

	NetworkPolicyInfo struct {
		Resource *networkingv1.NetworkPolicy `json:"resource"`
	}
)
type (
	IngressRequest struct {
		ResourceRequest
		Resource *networkingv1.Ingress `json:"resource" binding:"required,namespaced"`
	}

	IngressGetRequest struct {
		ResourceGetRequest
	}

	IngressListRequest struct {
		ResourceListRequest
	}

	IngressBatchRequest struct {
		Resources []*IngressRequest `json:"resources" binding:"dive"`
	}

	IngressResponse struct {
		Status error        `json:"status"`
		Info   *IngressInfo `json:"info"`
	}

	IngressListResponse struct {
		TotalCount         int                                        `json:"total_count"`
		List               []*IngressInfo                             `json:"list"`
		EachRangeListState []EachResourceRangeListState[*IngressInfo] `json:"each_range_list_state,omitempty"`
	}

	IngressInfo struct {
		Resource *networkingv1.Ingress `json:"resource"`
		// 控制器
		Controller v1.ObjectReference `json:"controller"`
		// 控制器是否健康
		ControllerHealthy bool `json:"controller_healthy"`
	}
)

type EventRequest struct {
	ResourceRequest
	Resource *v1.Event `json:"resource" binding:"required,namespaced"`
}

type EventGetRequest struct {
	ResourceGetRequest
}

type EventListRequest struct {
	ResourceListRequest
}

type EventResponse struct {
	Status error      `json:"status"`
	Info   *EventInfo `json:"info"`
}

type EventListResponse struct {
	TotalCount int          `json:"total_count"`
	List       []*EventInfo `json:"list"`
}

type EventInfo struct {
	Resource *v1.Event `json:"resource"`
}

type (
	DaemonSetRequest struct {
		ResourceRequest
		Resource *appsv1.DaemonSet `json:"resource" binding:"required,namespaced"`
		Version  string            `json:"version"`
	}

	DaemonSetListRequest struct {
		ResourceListRequest
	}

	DaemonSetGetRequest struct {
		ResourceGetRequest `binding:"namespaced"`
	}

	DaemonSetBatchRequest struct {
		Resources []*DaemonSetRequest `json:"resources" binding:"dive"`
	}

	DaemonSetResponse struct {
		Status error          `json:"status"`
		Info   *DaemonSetInfo `json:"info"`
	}

	DaemonSetVersionListRequest struct {
		ResourceGetRequest
		imachinery.PagingParams
	}

	DaemonSetVersionListResponse struct {
		// 版本列表
		VersionList []VersionInfo `json:"version_list"`
		// 当前版本
		CurrentVersion *VersionInfo `json:"current_version"`
		TotalCount     int          `json:"total_count"`
	}

	DaemonSetListResponse struct {
		TotalCount         int                                          `json:"total_count"`
		List               []*DaemonSetInfo                             `json:"list"`
		EachRangeListState []EachResourceRangeListState[*DaemonSetInfo] `json:"each_range_list_state,omitempty"`
	}

	DaemonSetInfo struct {
		*ResourceInfo[*appsv1.DaemonSet]
		ResourceConvert []*ResourceConvert `json:"resource_convert,omitempty"`
	}
)

type VersionInfo struct {
	Namespace   string              `json:"namespace"`
	Name        string              `json:"name"`
	Version     string              `json:"version"`
	CreateTime  int64               `json:"create_time"`
	StatefulSet *appsv1.StatefulSet `json:"statefulset,omitempty"`
	Deployment  *appsv1.Deployment  `json:"deployment,omitempty"`
	DaemonSet   *appsv1.DaemonSet   `json:"daemonset,omitempty"`
}

type (
	ReplicaSetRequest struct {
		ResourceRequest
		Resource *appsv1.ReplicaSet `json:"resource" binding:"required,namespaced"`
	}

	ReplicaSetListRequest struct {
		ResourceListRequest
	}

	ReplicaSetGetRequest struct {
		ResourceGetRequest
	}

	ReplicaSetResponse struct {
		Status error           `json:"status"`
		Info   *ReplicaSetInfo `json:"info"`
	}

	ReplicaSetListResponse struct {
		TotalCount         int                                           `json:"total_count"`
		List               []*ReplicaSetInfo                             `json:"list"`
		EachRangeListState []EachResourceRangeListState[*ReplicaSetInfo] `json:"each_range_list_state,omitempty"`
	}

	ReplicaSetInfo struct {
		Resource *appsv1.ReplicaSet `json:"resource"`
	}
)

type (
	HpaRequest struct {
		ResourceRequest
		HorizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler `json:"horizontal_pod_autoscaler" binding:"required,namespaced"`
	}

	HpaListRequest struct {
		ResourceListRequest
		FilterDeploymentName string `json:"filter_deployment_name"`
	}

	HpaGetRequest struct {
		ResourceGetRequest
	}

	HpaResponse struct {
		Status error    `json:"status"`
		Info   *HpaInfo `json:"info"`
	}

	HpaListResponse struct {
		TotalCount         int                                    `json:"total_count"`
		List               []*HpaInfo                             `json:"list"  `
		EachRangeListState []EachResourceRangeListState[*HpaInfo] `json:"each_range_list_state,omitempty"`
	}

	HpaInfo struct {
		Resource *autoscalingv1.HorizontalPodAutoscaler `json:"resource"`
	}
)

type (
	StatefulSetRequest struct {
		ResourceRequest
		Resource *appsv1.StatefulSet `json:"resource" binding:"required,namespaced"`
		Version  string              `json:"version"`
	}

	StatefulSetListRequest struct {
		ResourceListRequest
	}

	StatefulSetGetRequest struct {
		ResourceGetRequest
	}

	StatefulSetVersionListRequest struct {
		imachinery.PagingParams
		ResourceGetRequest
	}

	StatefulSetVersionListResponse struct {
		TotalCount     int              `json:"total_count"`
		VersionList    []VersionInfo    `json:"version_list"`
		CurrentVersion *VersionInfo     `json:"current_version,omitempty"`
		Info           *StatefulSetInfo `json:"info,omitempty"`
	}

	StatefulSetBatchRequest struct {
		Resources []*StatefulSetRequest `json:"resources" binding:"dive"`
	}

	StatefulSetListResponse struct {
		Status             error                                          `json:"status"`
		TotalCount         int                                            `json:"total_count"`
		List               []*StatefulSetInfo                             `json:"list"`
		EachRangeListState []EachResourceRangeListState[*StatefulSetInfo] `json:"each_range_list_state,omitempty"`
	}

	StatefulSetResponse struct {
		Status error            `json:"status"`
		Info   *StatefulSetInfo `json:"info,omitempty"`
	}

	StatefulSetInfo struct {
		Resource        *appsv1.StatefulSet `json:"resource"`
		ResourceConvert []*ResourceConvert  `json:"resource_convert,omitempty"`
	}
)

type (

	//pod request
	DeploymentRequest struct {
		ResourceRequest
		Resource      *appsv1.Deployment `json:"resource" binding:"required,namespaced"`
		MigrateTarget *appsv1.Deployment `json:"migrate_target"`
		Version       string             `json:"version"`
	}

	DeploymentListRequest struct {
		ResourceListRequest
	}

	DeploymentGetRequest struct {
		ResourceGetRequest
	}

	DeploymentBatchRequest struct {
		Resources []*DeploymentRequest `json:"resources" binding:"dive"`
	}

	ReasonMessage struct {
		Message string `json:"message"`
		Reason  string `json:"reason"`
	}

	DeploymentInfo struct {
		Resource        *appsv1.Deployment `json:"resource"`
		CurrentVersion  *VersionInfo       `json:"current_version,omitempty"`
		Warning         *ReasonMessage     `json:"warning,omitempty"` // return warn message if condition fail
		ResourceConvert []*ResourceConvert `json:"resource_convert,omitempty"`
	}

	//deployment response
	DeploymentResponse struct {
		Info *DeploymentInfo `json:"info"`
	}

	DeploymentVersoinListRequest struct {
		ResourceGetRequest
		imachinery.PagingParams
	}

	DeploymentVersionListResponse struct {
		VersionList    []VersionInfo `json:"version_list"`
		CurrentVersion *VersionInfo  `json:"current_version,omitempty"`
		TotalCount     int           `json:"total_count"`
	}

	//deployment response
	DeploymentListResponse struct {
		TotalCount         int                                           `json:"total_count"`
		List               []*DeploymentInfo                             `json:"list"`
		EachRangeListState []EachResourceRangeListState[*DeploymentInfo] `json:"each_range_list_state,omitempty"`
	}
)

type (
	PodDisruptionBudgetRequest struct {
		ResourceRequest
		Resource *policyv1.PodDisruptionBudget `json:"resource" binding:"required,namespaced"`
	}

	PodDisruptionBudgetGetRequest struct {
		ResourceGetRequest
	}

	PodDisruptionBudgetResponse struct {
		Status error                    `json:"status"`
		Info   *PodDisruptionBudgetInfo `json:"info,omitempty"`
	}

	PodDisruptionBudgetListRequest struct {
		ResourceListRequest
	}

	PodDisruptionBudgetListResponse struct {
		TotalCount         int                                                    `json:"total_count"`
		List               []*PodDisruptionBudgetInfo                             `json:"list"`
		EachRangeListState []EachResourceRangeListState[*PodDisruptionBudgetInfo] `json:"each_range_list_state,omitempty"`
	}

	PodDisruptionBudgetBatchRequest struct {
		Resources []*PodDisruptionBudgetRequest `json:"resources" binding:"dive" description:"应用中断干扰"`
	}

	PodDisruptionBudgetInfo struct {
		Resource *policyv1.PodDisruptionBudget `json:"resource"`
	}
)

type (
	StorageClassRequest struct {
		ResourceRequest
		Resource *storagev1.StorageClass `json:"resource" binding:"required,clusterd"`
		Secret   *v1.Secret              `json:"secret"` // secret for provisioner storageclass
	}

	StorageClassListRequest struct {
		ResourceListRequest
		StorageClass *storagev1.StorageClass `json:"storage_class"`
		Secret       *v1.Secret              `json:"secret"` // secret for provisioner storageclass
	}

	StorageClassGetRequest struct {
		ResourceGetRequest
	}

	StorageClassBatchRequest struct {
		Resources []*StorageClassRequest `json:"resources" binding:"dive"`
	}

	// cluster service response
	StorageClassResponse struct {
		Info *StorageClassInfo `json:"info,omitempty"`
	}

	StorageClassListResponse struct {
		TotalCount         int                                             `json:"total_count"`
		List               []*StorageClassInfo                             `json:"list"`
		EachRangeListState []EachResourceRangeListState[*StorageClassInfo] `json:"each_range_list_state,omitempty"`
	}

	StorageClassInfo struct {
		Resource *storagev1.StorageClass `json:"resource"`
	}
)

type (
	VolumeSnapshotClassRequest struct {
		ResourceRequest
		Resource *snapshotv1beta1.VolumeSnapshotClass `json:"resource" binding:"required,clusterd"`
	}

	VolumeSnapshotClassListRequest struct {
		ResourceListRequest
	}

	VolumeSnapshotClassGetRequest struct {
		ResourceGetRequest
	}

	VolumeSnapshotClassBatchRequest struct {
		Resources []*VolumeSnapshotClassRequest `json:"resources" binding:"dive"`
	}

	// cluster service response
	VolumeSnapshotClassResponse struct {
		Status error                    `json:"status"`
		Info   *VolumeSnapshotClassInfo `json:"info,omitempty"`
	}

	VolumeSnapshotClassListResponse struct {
		TotalCount         int                                                    `json:"total_count"`
		List               []*VolumeSnapshotClassInfo                             `json:"list,omitempty"`
		EachRangeListState []EachResourceRangeListState[*VolumeSnapshotClassInfo] `json:"each_range_list_state,omitempty"`
	}

	VolumeSnapshotClassInfo struct {
		Resource *snapshotv1beta1.VolumeSnapshotClass `json:"resource"`
	}
)

type (
	VolumeSnapshotRequest struct {
		ResourceRequest
		Resource                         *snapshotv1beta1.VolumeSnapshot `json:"resource"`
		PersistentVolumeClaimName        string                          `json:"persistent_volume_claim_name"`
		PersistentVolumeClaimAccessModes []v1.PersistentVolumeAccessMode `json:"persistent_volume_claim_access_modes"`
	}

	VolumeSnapshotListRequest struct {
		ResourceListRequest
		VolumeSnapshot            *snapshotv1beta1.VolumeSnapshot `json:"volume_snapshot"`
		PersistentVolumeClaimName string                          `json:"persistent_volume_claim_name"`
	}

	VolumeSnapshotGetRequest struct {
		ResourceGetRequest
	}

	VolumeSnapshotBatchRequest struct {
		Resources []*VolumeSnapshotRequest `json:"resources" binding:"required,dive"`
	}

	// cluster service response
	VolumeSnapshotListResponse struct {
		Status             error                                             `json:"status"`
		TotalCount         int                                               `json:"total_count"`
		List               []*VolumeSnapshotInfo                             `json:"list,omitempty"`
		EachRangeListState []EachResourceRangeListState[*VolumeSnapshotInfo] `json:"each_range_list_state,omitempty"`
	}

	VolumeSnapshotResponse struct {
		Status error               `json:"status"`
		Info   *VolumeSnapshotInfo `json:"info,omitempty"`
	}

	VolumeSnapshotInfo struct {
		Resource *snapshotv1beta1.VolumeSnapshot `json:"resource"`
	}
)

type (
	VolumeSnapshotContentRequest struct {
		ResourceRequest
		Resource *snapshotv1beta1.VolumeSnapshotContent `json:"resource" binding:"required,clusterd"`
	}

	VolumeSnapshotContentListRequest struct {
		ResourceListRequest
	}

	VolumeSnapshotContentGetRequest struct {
		ResourceGetRequest `binding:"clusterd"`
	}

	VolumeSnapshotContentBatchRequest struct {
		Resources []*VolumeSnapshotContentRequest `json:"resources" binding:"dive"`
	}

	VolumeSnapshotContentResponse struct {
		Info *VolumeSnapshotContentInfo `json:"info,omitempty"`
	}

	VolumeSnapshotContentListResponse struct {
		Status             error                                                    `json:"status"`
		TotalCount         int                                                      `json:"total_count"`
		List               []*VolumeSnapshotContentInfo                             `json:"list,omitempty"`
		EachRangeListState []EachResourceRangeListState[*VolumeSnapshotContentInfo] `json:"each_range_list_state,omitempty"`
	}

	VolumeSnapshotContentInfo struct {
		Resource *snapshotv1beta1.VolumeSnapshotContent `json:"resource"`
	}
)

type (
	PersistentVolumeListRequest struct {
		ResourceListRequest
		// pv不支持storage class field过滤, 提供支持
		StorageClassFilter string `json:"-" form:"-"`
	}

	PersistentVolumeGetRequest struct {
		ResourceGetRequest `binding:"clusterd"`
	}

	PersistentVolumeRequest struct {
		ResourceRequest
		Resource *v1.PersistentVolume `json:"resource" binding:"required,clusterd"`
	}

	PersistentVolumeBatchRequest struct {
		Resources []*PersistentVolumeRequest `json:"resources" binding:"dive"`
	}

	PersistentVolumeResponse struct {
		Status error                 `json:"status"`
		Info   *PersistentVolumeInfo `json:"info"`
	}

	PersistentVolumeListResponse struct {
		Status             error                                               `json:"status"`
		TotalCount         int                                                 `json:"total_count"`
		List               []*PersistentVolumeInfo                             `json:"list"`
		EachRangeListState []EachResourceRangeListState[*PersistentVolumeInfo] `json:"each_range_list_state,omitempty"`
	}

	PersistentVolumeInfo struct {
		Resource *v1.PersistentVolume `json:"resource,omitempty"`
	}
)

var _ imachinery.PostBinder = &PersistentVolumeListRequest{}

func (r *PersistentVolumeListRequest) PostBind() error {
	storageClassFilter := ""
	if r.FieldSelector != "" && strings.Contains(r.FieldSelector, "spec.storageClassName") {
		filters := strings.Split(r.FieldSelector, ",")
		r.FieldSelector = ""
		for _, v := range filters {
			kvs := strings.Split(v, "=")
			if kvs[0] == "spec.storageClassName" {
				if len(kvs) != 1 {
					storageClassFilter = kvs[1]
				}
			} else {
				//remove unsupported field
				r.FieldSelector = r.FieldSelector + "," + v
			}
		}
	}
	r.StorageClassFilter = storageClassFilter
	return nil
}

type PersistentVolumeClaimListRequest struct {
	ResourceListRequest
}

type PersistentVolumeClaimGetRequest struct {
	ResourceGetRequest `binding:"namespaced"`
}

type PersistentVolumeClaimRequest struct {
	ResourceRequest
	Resource *v1.PersistentVolumeClaim `json:"resource" binding:"required,namespaced"`
}

type PersistentVolumeClaimBatchRequest struct {
	Resources []*PersistentVolumeClaimRequest `json:"resources" binding:"dive"`
}

type PersistentVolumeClaimListResponse struct {
	Status             error                                                    `json:"status"`
	TotalCount         int                                                      `json:"total_count"`
	List               []*PersistentVolumeClaimInfo                             `json:"list"`
	EachRangeListState []EachResourceRangeListState[*PersistentVolumeClaimInfo] `json:"each_range_list_state,omitempty"`
}

type PersistentVolumeClaimResponse struct {
	Status error                      `json:"status"`
	Info   *PersistentVolumeClaimInfo `json:"info"`
}

type PersistentVolumeClaimInfo struct {
	Resource *v1.PersistentVolumeClaim `json:"resource"`
	// 支持扩展
	AllowExpansion *bool `json:"allow_expansion,omitempty"`
	// 支持快照
	AllowSnapshot *bool `json:"allow_snapshot,omitempty"`
}
