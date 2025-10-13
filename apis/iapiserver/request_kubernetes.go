package iapiserver

import (
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"

	"github.com/wangweihong/eazycloud/apis/imachinery"
)

type EachResourceRangeListState struct {
	ClusterUUID string `json:"cluster_uuid"`
	ClusterName string `json:"cluster_name"`
	Result      error  `json:"result"`
	TotalCount  int    `json:"total_count"`
	List        any    `json:"list,omitempty"`
}

func NewEachResourceRangeListState(uuid, name string) EachResourceRangeListState {
	return EachResourceRangeListState{
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
type ServiceRequest struct {
	ResourceRequest
	Resource *v1.Service `json:"service"`
}

type ServiceListRequest struct {
	ResourceListRequest
}

type ServiceGetRequest struct {
	ResourceGetRequest
}

// cluster service info
type ServiceInfo struct {
	*v1.Service
}

// cluster service response
type ServiceResponse struct {
	Info *ServiceInfo `json:"info" description:"服务详情"`
}

type ServiceListResponse struct {
	TotalCount         int                          `json:"total_count"`
	List               []*ServiceInfo               `json:"list"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state"`
}

type ServiceBatchRequest struct {
	Services []*ServiceRequest `json:"services"`
}

// cluster endpoints request
type EndpointsRequest struct {
	FilterName string        `json:"filter_name" form:"filter_name"`
	Endpoints  *v1.Endpoints `json:"endpoints"` // endpoints name is equal to service name
	ResourceRequest
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
	*v1.Endpoints
}

// cluster endpoints response
type EndpointsResponse struct {
	Info *EndpointsInfo `json:"info,omitempty"`
}

// cluster endpoints response
type EndpointsListResponse struct {
	TotalCount         int                          `json:"total_count"`
	List               []*EndpointsInfo             `json:"list"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
}

const (
	NodeRoleMaster  = "master"
	NodeRoleGateway = "gateway"
	NodeRoleWorker  = "worker"

	IngressControllerNamespace = "system-router"
	IngressControllerPrefix    = "router-"
)

// node request
type NodeRequest struct {
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

type NodeEventRequest struct {
	ResourceGetRequest
	imachinery.PagingParams
	Fuzzy string `json:"fuzzy" form:"fuzzy"`
}

type NodeListRequest struct {
	ResourceListRequest
	ShowResourceUsage bool `json:"show_resource_usage" form:"show_resource_usage"`
}

type NodeGetRequest struct {
	ResourceGetRequest
	ShowResourceUsage bool `json:"show_resource_usage" form:"show_resource_usage"`
}

type NodeInfo struct {
	*v1.Node
	ResourceUsage *NodeRealtimeResource `json:"resource_usage,omitempty"`
	NodeStatus    string                `json:"node_status,omitempty"`
	NodeAddr      string                `json:"node_addr,omitempty"`
	Role          string                `json:"role,omitempty"`
}

// node response
type NodeResponse struct {
	Info *NodeInfo `json:"info"`
}

type NodeListResponse struct {
	NodeRoleCount
	TotalCount         int                          `json:"total_count"`
	List               []*NodeInfo                  `json:"list"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
}

type NodeRoleCount struct {
	NodeMasterCount  int `json:"node_master"`
	NodeWorkCount    int `json:"node_work"`
	NodeGatewayCount int `json:"node_gateway"`
}

// node request
type NodeDrainRequest struct {
	NodeRequest
	DisableEviction bool `json:"disable_eviction"` // if disable use delete instead of evict
}

// node request
type NodeDrainResponse struct {
	//*BatchOutput
}
type NodeRealtimeResource struct {
	*NodeRealTimePhysicalResourceUsage
	*NodeRealTimePodUsage
	*NodePodListResourceRequestAndLimit
}

type NodeRealTimePhysicalResourceUsage struct {
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

type NodeRealTimePodUsage struct {
	PodCapacity  int     `json:"pod_capacity"`
	PodUsed      int     `json:"pod_used"`
	PodUsedRatio float64 `json:"pod_used_ratio"`
}

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

type SecretRequest struct {
	ResourceRequest
	Secret *v1.Secret `json:"secret"`
}

type SecretGetRequest struct {
	ResourceGetRequest
}

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

type ResourceQuotaRequest struct {
	ResourceRequest
	ResourceQuota  *v1.ResourceQuota `json:"resource_quota"`
	UpdateIfExists bool              `json:"update_if_exists" description:"存在则更新"`
}

type ResourceQuotaListRequest struct {
	ResourceListRequest
}

type ResourceQuotaGetRequest struct {
	ResourceGetRequest
}

type ResourceQuotaResponse struct {
	Info *ResourceQuotaInfo `json:"info" description:"详情"`
}

type ResourceQuotaListResponse struct {
	Status             error                        `json:"status"                          description:"状态码"`
	TotalCount         int                          `json:"total_count"                     description:"总数"`
	List               []*ResourceQuotaInfo         `json:"list"                            description:"列表"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty" description:"每个集群返回详情"`
}

type ResourceQuotaInfo struct {
	*v1.ResourceQuota
	Convert *ResourceQuotaConvert `json:"convert,omitempty"`
}

type ResourceQuotaConvert struct {
	Spec   ResourceQuotaSpec   `json:"spec,omitempty"`
	Status ResourceQuotaStatus `json:"status,omitempty"`
}

type ResourceQuotaSpec struct {
	Hard map[string]int64 `json:"hard,omitempty"`
}

type ResourceQuotaStatus struct {
	Hard map[string]int64 `json:"hard,omitempty"`
	Used map[string]int64 `json:"used,omitempty"`
}

type JobRequest struct {
	Job *batchv1.Job `json:"job"`
	ResourceRequest
}

type JobListRequest struct {
	ResourceListRequest
}

type JobGetRequest struct {
	ResourceGetRequest
}

type JobBatchRequest struct {
	Jobs []*JobRequest `json:"jobs"`
}

type JobResponse struct {
	Status error    `json:"status" description:"状态码"`
	Info   *JobInfo `json:"info"   description:"详情"`
}

type JobListResponse struct {
	Status             error                        `json:"status"                          description:"状态码"`
	TotalCount         int                          `json:"total_count"                     description:"总数"`
	List               []*JobInfo                   `json:"list"                            description:"列表"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty" description:"每个集群返回详情"`
}

type JobInfo struct {
	*batchv1.Job
	ResourceConvert []*ResourceConvert `json:"resource_convert,omitempty"`
}

type MonitorData struct {
	Namespace    string            `json:"namespace"` //
	Name         string            `json:"name"`      //
	Addr         string            `json:"addr"`
	Type         string            `json:"type"`          //
	Cpu          float64           `json:"cpu"`           // 使用率
	CpuUsed      float64           `json:"cpu_used"`      //
	CpuTotal     float64           `json:"cpu_total"`     //
	Memory       float64           `json:"memory"`        // 使用率
	MemoryUsed   float64           `json:"memory_used"`   //
	MemoryTotal  float64           `json:"memory_total"`  //
	DiskRatio    float64           `json:"disk_ratio"`    // 使用率
	DiskUsed     float64           `json:"disk_used"`     //
	DiskTotal    float64           `json:"disk_total"`    //
	DiskRead     float64           `json:"disk_read"`     //
	DiskWrite    float64           `json:"disk_write"`    //
	PodRatio     float64           `json:"pod_ratio"`     // 使用率
	PodCount     float64           `json:"pod_count"`     // 当前有多少的pod
	PodTotal     float64           `json:"pod_total"`     // 集群总共可以创建多少个pod
	NetworkRead  float64           `json:"network_read"`  //
	NetworkWrite float64           `json:"network_write"` //
	Extra        map[string]string `json:"extra"`
}

type MonitorDataResponse struct {
	Status       error             `json:"status"`
	Namespace    string            `json:"namespace,omitempty"` //
	Name         string            `json:"name"`                //
	Addr         string            `json:"addr"`
	Type         string            `json:"type"`          //
	Cpu          float64           `json:"cpu"`           // 使用率
	CpuUsed      float64           `json:"cpu_used"`      //
	CpuTotal     float64           `json:"cpu_total"`     //
	Memory       float64           `json:"memory"`        // 使用率
	MemoryUsed   float64           `json:"memory_used"`   //
	MemoryTotal  float64           `json:"memory_total"`  //
	DiskRatio    float64           `json:"disk_ratio"`    // 使用率
	DiskUsed     float64           `json:"disk_used"`     //
	DiskTotal    float64           `json:"disk_total"`    //
	DiskRead     float64           `json:"disk_read"`     //
	DiskWrite    float64           `json:"disk_write"`    //
	PodRatio     float64           `json:"pod_ratio"`     // 使用率
	PodCount     float64           `json:"pod_count"`     // 当前有多少的pod
	PodTotal     float64           `json:"pod_total"`     // 集群总共可以创建多少个pod
	NetworkRead  float64           `json:"network_read"`  //
	NetworkWrite float64           `json:"network_write"` //
	Extra        map[string]string `json:"extra"`
}
