package iapiserver

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	monitorv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/wangweihong/eazycloud/apis/iprometheus"
	"github.com/wangweihong/eazycloud/internal/pkg/libprometheus"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type (
	AlertManagerRequest struct {
		ResourceRequest
		Resource *monitorv1.Alertmanager `json:"resource" binding:"required"`
	}

	AlertManagerListRequest struct {
		ResourceListRequest
	}

	AlertManagerGetRequest struct {
		ResourceGetRequest
	}

	AlertManagerResponse struct {
		*ResourceInfo[*monitorv1.Alertmanager]
	}
	AlertManagerListResponse struct {
		TotalCount int                 `json:"total_count"`
		List       []*AlertManagerInfo `json:"list"`
	}

	AlertManagerInfo struct {
		*ResourceInfo[*monitorv1.Alertmanager]
	}
)

func NewAlertManagerInfo(meta *monitorv1.Alertmanager, cluster *Cluster) *AlertManagerInfo {
	return &AlertManagerInfo{
		ResourceInfo: NewResourceInfo(meta, cluster),
	}
}

type (
	ServiceMonitorRequest struct {
		ResourceRequest
		Resource *monitorv1.ServiceMonitor `json:"resource"`
	}

	ServiceMonitorGetRequest struct {
		ResourceGetRequest
	}

	ServiceMonitorListRequest struct {
		ResourceListRequest
	}

	ServiceMonitorResponse struct {
		Info *ServiceMonitorInfo `json:"info"`
	}

	ServiceMonitorListResponse struct {
		TotalCount int                   `json:"total_count"`
		List       []*ServiceMonitorInfo `json:"list"`
	}

	ServiceMonitorInfo struct {
		*ResourceInfo[*monitorv1.ServiceMonitor]
	}
)

func NewServiceMonitorInfo(meta *monitorv1.ServiceMonitor, cluster *Cluster) *ServiceMonitorInfo {
	return &ServiceMonitorInfo{
		ResourceInfo: NewResourceInfo(meta, cluster),
	}
}

type (
	PodMonitorRequest struct {
		ResourceRequest
		Resource *monitorv1.PodMonitor `json:"resource"`
	}

	PodMonitorListRequest struct {
		ResourceListRequest
	}

	PodMonitorGetRequest struct {
		ResourceGetRequest
	}

	PodMonitorResponse struct {
		Info *PodMonitorInfo `json:"info"`
	}

	PodMonitorListResponse struct {
		TotalCount int               `json:"total_count"`
		List       []*PodMonitorInfo `json:"list"`
	}

	PodMonitorInfo struct {
		*ResourceInfo[*monitorv1.PodMonitor]
	}
)

func NewPodMonitorInfo(meta *monitorv1.PodMonitor, cluster *Cluster) *PodMonitorInfo {
	return &PodMonitorInfo{
		ResourceInfo: NewResourceInfo(meta, cluster),
	}
}

type (
	MonitoringTimeRange struct {
		StartTime int64 `json:"start_time" form:"start_time"` //起始时间
		EndTime   int64 `json:"end_time"   form:"end_time"`   //结束时间
		Duration  int64 `json:"duration"   form:"duration"`   //采集间隔。秒为单位
	}

	Monitoring struct {
		MonitoringTimeRange
		Level string `json:"level" form:"level"`
		Time  int64  `json:"time" form:"time"`

		Metrics  []string `json:"metrics" form:"metrics"`
		Type     string   `json:"type" form:"type"` //resourceType, such deployment/statefulset/daemonset
		Name     string   `json:"name" form:"name"` //resourceName, such as deployment Name
		NodeName string   `json:"node" form:"node"`
	}

	MonitoringParam struct {
		MonitoringTimeRange
		Level    string `json:"level" form:"level"`
		Time     int64  `json:"time" form:"time"` //某个时间点
		Metrics  string `json:"metrics" form:"metrics"`
		Type     string `json:"type" form:"type"` //resourceType, such deployment/statefulset/daemonset
		Name     string `json:"name" form:"name"` //resourceName, such as deployment Name
		NodeName string `json:"node" form:"node"`
	}
)

func (p MonitoringParam) ToMonitoring() *Monitoring {
	var metrics []string
	if p.Metrics != "" {
		metrics = strings.Split(p.Metrics, ",")
	}

	return &Monitoring{
		MonitoringTimeRange: p.MonitoringTimeRange,
		Level:               p.Level,
		Time:                p.Time,
		Metrics:             metrics,
		Type:                p.Type,
		Name:                p.Name,
		NodeName:            p.NodeName,
	}
}

//	v1.Range要求(endTime-startTime)/step < 11000, 防止数据点过大
//
// duration是纳秒!
func (m MonitoringTimeRange) ToRange() v1.Range {
	if m.StartTime == 0 && m.EndTime == 0 {
		d, _ := time.ParseDuration("-30m")
		m.EndTime = time.Now().Unix()
		m.StartTime = time.Now().Add(d).Unix()
	}

	//1 second
	if m.Duration == 0 {
		m.Duration = 1
	}

	return v1.Range{Start: time.Unix(m.StartTime, 0), End: time.Unix(m.EndTime, 0), Step: time.Duration(m.Duration) * time.Second}
}

type MonitoringRequest struct {
	ResourceRequest
	FilterMaster    bool        `json:"filter_master"`
	Monitoring      *Monitoring `json:"monitoring"`
	AlertDeleteList []*v1.Alert `json:"alert_delete_list" binding:"required"`
	Namespace       string      `json:"namespace"`
}

type MonitoringGetRequest struct {
	ResourceRequest
	MonitoringTimeRange

	Namespace string `json:"namespace" form:"namespace"` // only use in list action
	Name      string `json:"name" form:"name" binding:"required"`
}

type MonitoringListRequest struct {
	ResourceListRequest
	MonitoringParam
	FilterMaster bool `json:"filter_master" form:"filter_master"`
}

func (req *MonitoringListRequest) Validate() error {
	if req == nil {
		return errors.Errorf("metrics is empty")
	}
	if req.StartTime == 0 && req.EndTime == 0 {
		d, _ := time.ParseDuration("-30m")
		req.EndTime = time.Now().Unix()
		req.StartTime = time.Now().Add(d).Unix()
	}

	if req.Duration == 0 {
		req.Duration = 1
	}

	switch req.MonitoringParam.Level {
	case libprometheus.NamespaceLevel,
		libprometheus.NodeLevel,
		libprometheus.ClusterLevel,
		libprometheus.PodLevel,
		libprometheus.AlertLevel:
	default:
		return errors.Errorf("monitoring.level[%v] not support", req.MonitoringParam.Level)
	}

	return nil
}

type MonitoringResponse struct {
	TotalCount    int                  `json:"total_count"`
	Metrics       []iprometheus.Metric `json:"metrics,omitempty"`
	Metric        *iprometheus.Metric  `json:"metric,omitempty"`
	NodeUsedTotal *NodeOverview        `json:"node_used_total,omitempty"`
	MonitorDatas  []*MonitorData       `json:"namespace_top,omitempty"` // 命名空间使用量排行
	NodeList      []*NodeOverview      `json:"node_list,omitempty"`     // 节点使用量排行列表
	PodList       []*PodOverview       `json:"pod_list,omitempty"`      // pod资源使用列表
	AlertList     []*Alert             `json:"alert_list,omitempty"`
	MetricList    []string             `json:"metric_list,omitempty"`
	MonitorStatus *MonitorStatus       `json:"monitor_status,omitempty"` //grafana/prometheus组件的状态
	//EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
}

type MonitoringListResponse struct {
	TotalCount    int                  `json:"total_count"`
	Metrics       []iprometheus.Metric `json:"metrics,omitempty"`
	Metric        *iprometheus.Metric  `json:"metric,omitempty"`
	NodeUsedTotal *NodeOverview        `json:"node_used_total,omitempty"`
	MonitorDatas  []*MonitorData       `json:"namespace_top,omitempty"` // 命名空间使用量排行
	NodeList      []*NodeOverview      `json:"node_list,omitempty"`     // 节点使用量排行列表
	PodList       []*PodOverview       `json:"pod_list,omitempty"`      // pod资源使用列表
	AlertList     []*Alert             `json:"alert_list,omitempty"`
	MetricList    []string             `json:"metric_list,omitempty"`
	MonitorStatus *MonitorStatus       `json:"monitor_status,omitempty"` //grafana/prometheus组件的状态
	//EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
}

type MonitoringAlertListResponse struct {
	TotalCount         int                                  `json:"total_count"`
	AlertList          []*Alert                             `json:"alert_list,omitempty"`
	EachRangeListState []EachResourceRangeListState[*Alert] `json:"each_range_list_state,omitempty"`
}

type MonitoringNodeListResponse struct {
	TotalCount         int                                         `json:"total_count"`
	List               []*NodeOverview                             `json:"list,omitempty"`
	EachRangeListState []EachResourceRangeListState[*NodeOverview] `json:"each_range_list_state,omitempty"`
}

type Alert struct {
	ClusterUuid string    `json:"cluster_uuid"`
	ClusterName string    `json:"cluster_name"`
	Namespace   string    `json:"namespace"`
	Name        string    `json:"name"`
	Pod         string    `json:"pod"`
	Container   string    `json:"container"`
	Alertstate  string    `json:"alertstate"`
	Severity    string    `json:"severity"`
	ActiveTime  int64     `json:"active_time"`
	Message     string    `json:"message"`
	Value       string    `json:"value"`
	Metadata    *v1.Alert `json:"metadata"`
}

type MonitoringApiserverResponse struct {
	TimeRange  *v1.Range             `json:"time_range,omitempty"`
	CpuUsed    iprometheus.MetricOne `json:"cpu_used"`
	MemoryUsed iprometheus.MetricOne `json:"memory_used"`
	Goroutine  iprometheus.MetricOne `json:"goroutine"`
}

type MonitoringControllerManagerResponse struct {
	TimeRange        *v1.Range             `json:"time_range,omitempty"`
	UpCount          int                   `json:"up_count"`
	WorkQueueAddRate iprometheus.MetricOne `json:"work_queue_add_rate"`
	CpuUsed          iprometheus.MetricOne `json:"cpu_used"`
	MemoryUsed       iprometheus.MetricOne `json:"memory_used"`
	Goroutine        iprometheus.MetricOne `json:"goroutine"`
}

type MonitoringKubeletResponse struct {
	TimeRange  *v1.Range             `json:"time_range,omitempty"`
	UpCount    int                   `json:"up_count"`
	PodCount   int                   `json:"pod_count"`
	RpcRate    iprometheus.MetricOne `json:"rpc_rate"`
	CpuUsed    iprometheus.MetricOne `json:"cpu_used"`
	MemoryUsed iprometheus.MetricOne `json:"memory_used"`
	Goroutine  iprometheus.MetricOne `json:"goroutine"`
}

type MonitoringSchedulerResponse struct {
	TimeRange     *v1.Range             `json:"time_range"`
	UpCount       int                   `json:"up_count"`
	ScheduleCount iprometheus.MetricOne `json:"schedule_count"`
	CpuUsed       iprometheus.MetricOne `json:"cpu_used"`
	MemoryUsed    iprometheus.MetricOne `json:"memory_used"`
	Goroutine     iprometheus.MetricOne `json:"goroutine"`
}

type MonitoringEtcdResponse struct {
	TimeRange            *v1.Range             `json:"time_range,omitempty"`
	HasLeader            bool                  `json:"has_leader"`
	LeaderChangeCount    int                   `json:"leader_change_count"`
	NetReadByte          iprometheus.MetricOne `json:"net_read_byte"`
	NetWriteByte         iprometheus.MetricOne `json:"net_write_byte"`
	MemorySize           iprometheus.MetricOne `json:"memory_size"`
	RaftProposalApplied  iprometheus.MetricOne `json:"raft_proposal_applied"`
	RaftProposalCommited iprometheus.MetricOne `json:"raft_proposal_commited"`
	RaftProposalFailed   iprometheus.MetricOne `json:"raft_proposal_failed"`
	RaftProposalPending  iprometheus.MetricOne `json:"raft_proposal_pending"`
}

type MonitorStatus struct {
	GrafanaEnable    bool   `json:"grafana_enable"`
	GrafanaUrl       string `json:"grafana_url"`
	GrafanaError     string `json:"grafana_error"`
	PrometheusEnable bool   `json:"prometheus_enable"`
	PrometheusUrl    string `json:"prometheus_url"`
	PrometheusError  string `json:"prometheus_error"`
}

type PodOverview struct {
	Namespace  string            `json:"namespace" description:"命名空间"`
	Name       string            `json:"name" description:"名称"`
	CreateTime int64             `json:"create_time" description:"创建时间"`
	NodeName   string            `json:"node_name" description:"节点名"`
	NodeAddr   string            `json:"node_addr" description:"节点地址"`
	Extra      map[string]string `json:"extra" description:"额外数据"`
	HardWareResource
	//use to overview
	ClusterName string `json:"cluster_name" description:"集群名"`
}
type HardWareResource struct {
	Cpu           ResourceUsageInfo `json:"cpu" description:"物理CPU"`
	CpuLimit      ResourceUsageInfo `json:"cpu_limit" description:"CPU限制"`
	CpuRequest    ResourceUsageInfo `json:"cpu_request" description:"CPU请求"`
	Memory        ResourceUsageInfo `json:"memory" description:"物理内存"`
	MemoryLimit   ResourceUsageInfo `json:"memory_limit" description:"内存限制"`
	MemoryRequest ResourceUsageInfo `json:"memory_request" description:"内存请求"`
	Disk          ResourceUsageInfo `json:"disk" description:"物理磁盘"`
	Net           ResourceUsageInfo `json:"net" description:"网络"`
	NetRead       ResourceUsageInfo `json:"net_read" description:"网络读"`
	NetWrite      ResourceUsageInfo `json:"net_write" description:"网络写"`
}

type SoftWareResource struct {
	Deployment  ResourceUsageInfo `json:"deployment" description:"deployment"`
	Namespace   ResourceUsageInfo `json:"namespace" description:"命名空间"`
	Statefulset ResourceUsageInfo `json:"statefulset" description:"statefulset"`
	Daemonset   ResourceUsageInfo `json:"daemonset" description:"daemonset"`
	Service     ResourceUsageInfo `json:"service" description:"service"`
	CronJob     ResourceUsageInfo `json:"cronjob" description:"cronjob"`
	Job         ResourceUsageInfo `json:"job" description:"job"`
	Ingress     ResourceUsageInfo `json:"ingress" description:"ingress"`
	Pod         ResourceUsageInfo `json:"pod" description:"pod"`
	Pvc         ResourceUsageInfo `json:"pvc" description:"pvc"`
}

type ResourceUsageInfo struct {
	Used  float64 `json:"Used" description:"已使用"`
	Total float64 `json:"Total" description:"总数"`
	Usage float64 `json:"Usage" description:"使用率"`
}

type MonitoringMetricsHistoryListResponse struct {
	TimeRange *v1.Range            `json:"time_range"`
	Metrics   []iprometheus.Metric `json:"metrics,omitempty"`
}

type MonistoringHardwareHistoryListRequest struct {
	ResourceRequest
	MonitoringTimeRange
}

type MonistoringHardwareHistoryListResponse struct {
	TimeRange *v1.Range            `json:"time_range"`
	Metrics   []iprometheus.Metric `json:"metrics,omitempty"`
}

type MonitoringNodeGetRequest struct {
	ResourceRequest
	MonitoringTimeRange

	NodeName string `json:"node" form:"node" binding:"required"`
}

type MonitoringNodeGetResponse struct {
	TimeRange *v1.Range            `json:"time_range"`
	Metric    *iprometheus.Metric  `json:"metric,omitempty"`
	Metrics   []iprometheus.Metric `json:"metrics,omitempty"`
}

type MonitoringNamespaceListRequest struct {
	ResourceListRequest
}

type MonitoringNamespaceListResponse struct {
	TotalCount         int                                        `json:"total_count"`
	List               []*MonitorData                             `json:"list,omitempty"` // 命名空间使用量排行
	EachRangeListState []EachResourceRangeListState[*MonitorData] `json:"each_range_list_state,omitempty"`
}

type MonitorData struct {
	Namespace    string  `json:"namespace"`
	Name         string  `json:"name"`
	Addr         string  `json:"addr"`
	Type         string  `json:"type"`
	Cpu          float64 `json:"cpu"`
	CpuUsed      float64 `json:"cpu_used"`
	CpuTotal     float64 `json:"cpu_total"`
	Memory       float64 `json:"memory"`
	MemoryUsed   float64 `json:"memory_used"`
	MemoryTotal  float64 `json:"memory_total"`
	DiskRatio    float64 `json:"disk_ratio"`
	DiskUsed     float64 `json:"disk_used"`
	DiskTotal    float64 `json:"disk_total"`
	DiskRead     float64 `json:"disk_read"`
	DiskWrite    float64 `json:"disk_write"`
	PodRatio     float64 `json:"pod_ratio"`
	PodCount     float64 `json:"pod_count"` // 当前有多少的pod
	PodTotal     float64 `json:"pod_total"` // 集群总共可以创建多少个pod
	NetworkRead  float64 `json:"network_read"`
	NetworkWrite float64 `json:"network_write"`
}

// +gen:sortfields
type MonitorDataResponse struct {
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

type MonitoringPodListRequest struct {
	ResourceListRequest
}

type MonitoringPodListResponse struct {
	TotalCount         int                                        `json:"total_count"`
	PodList            []*PodOverview                             `json:"pod_list,omitempty"` // pod资源使用列表
	EachRangeListState []EachResourceRangeListState[*PodOverview] `json:"each_range_list_state,omitempty"`
}

type (
	PrometheusRequest struct {
		ResourceRequest
		Resource *monitorv1.Prometheus `json:"prometheus"`
	}

	PrometheusListRequest struct {
		ResourceListRequest
	}

	PrometheusGetRequest struct {
		ResourceGetRequest
	}

	PrometheusResponse struct {
		Info *PrometheusInfo `json:"info"`
	}

	PrometheusListResponse struct {
		TotalCount int               `json:"total_count"`
		List       []*PrometheusInfo `json:"list"`
	}

	PrometheusInfo struct {
		Url             string                     `json:"url,omitempty"`         // prometheus
		GrafanaUrl      string                     `json:"grafana_url,omitempty"` // grafana
		Prometheus      *monitorv1.Prometheus      `json:"prometheus"`
		ResourceConvert *PrometheusResourceConvert `json:"resource_convert,omitempty"`
	}

	PrometheusResourceConvert struct {
		Limits  map[string]int64 `json:"limits,omitempty"`
		Request map[string]int64 `json:"requests,omitempty"`
	}
)

type MonistoringPodContainerHistoryListRequest struct {
	ResourceRequest
	MonitoringTimeRange
	Name string `json:"name" form:"name" binding:"required"` //pod 名

}

type MonistoringPodContainerHistoryListResponse struct {
	TimeRange *v1.Range            `json:"time_range,omitempty"`
	Metrics   []iprometheus.Metric `json:"metrics,omitempty"`
}

type MonitoringMetricsNameListRequest struct {
	ResourceRequest
}

type MonitoringMetricsNameListResponse struct {
	TotalCount int      `json:"total_count"`
	MetricList []string `json:"metric_list,omitempty"`
}

type MonitoringTopGetRequest struct {
	ResourceListRequest
	Namespace string `json:"namespace" form:"namespace"`
}

type MonitoringTopGetResponse struct {
	TotalCount   int            `json:"total_count"`
	MonitorDatas []*MonitorData `json:"namespace_top,omitempty"` // 命名空间使用量排行
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
	GrafanaRequest struct {
		ResourceRequest
		SubUrl    string  `json:"sub_url"`
		Namespace *string `json:"namespace"`
		Workload  *string `json:"workload"`
		Pod       *string `json:"pod"`
		Node      *string `json:"node"`
	}

	GrafanaGetRequest struct {
		ResourceGetRequest
		SubUrl   string `json:"sub_url" form:"sub_url" description:"子路径URL"`
		Workload string `json:"workload" form:"workload" description:"工作负载类型"`
		Pod      string `json:"pod" form:"pod" description:"Pod"`
		Node     string `json:"node" form:"node" description:"节点名"`
	}

	GrafanaListRequest struct {
		ResourceListRequest
	}

	GrafanaResponse struct {
		TotalUrl string `json:"total_url" description:"URL"`
	}
)

type (
	PrometheusRule struct {
		metav1.ObjectMeta `json:"metadata"`
		Spec              PrometheusRuleSpec `json:"spec"`
	}
	PrometheusRuleSpec struct {
		monitorv1.Rule
		Operator string `json:"operator"`
		Value    string `json:"value"`
		Level    string `json:"level"`
		Desc     string `json:"desc"`
	}

	PrometheusGroupInfo struct {
		RuleCount int `json:"rule_count"`
		monitorv1.RuleGroup
	}

	PrometheusRuleInfo struct {
		*ResourceInfo[*monitorv1.PrometheusRule]
	}

	PrometheusRuleRequest struct {
		ResourceRequest
		PrometheusRule *monitorv1.PrometheusRule `json:"prometheus_rule" binding:"required"`
	}

	PrometheusRuleGetRequest struct {
		ResourceGetRequest
	}

	PrometheusRuleListRequest struct {
		ResourceListRequest
	}

	PrometheusRuleListResponse struct {
		TotalCount         int                                                    `json:"total_count,omitempty" description:"总数"`
		List               []*PrometheusRuleInfo                                  `json:"list,omitempty" description:"规则详情"`
		EachRangeListState []EachResourceRangeListState[*PrometheusRuleSpecGroup] `json:"each_range_list_state,omitempty" description:"每个集群返回状态"`
	}
)

func NewPrometheusRuleInfo(meta *monitorv1.PrometheusRule, cluster *Cluster) *PrometheusRuleInfo {
	return &PrometheusRuleInfo{
		ResourceInfo: NewResourceInfo(meta, cluster),
	}
}

type (
	PrometheusRuleSpecGroup struct {
		metav1.ObjectMeta `json:"metadata"`
		Spec              PrometheusGroupInfo `json:"spec"`
	}

	BuiltinPrometheusRuleSpecGroupRequest struct {
		ResourceRequest
		GroupList []*PrometheusRuleSpecGroup `json:"group_list" binding:"required,dive"`
	}

	BuiltinPrometheusRuleSpecGroupGetRequest struct {
		ResourceRequest
		GroupName string `json:"group_name" binding:"required"`
	}

	BuiltinPrometheusRuleSpecGroupListRequest struct {
		ResourceListRequest
	}

	BuiltinPrometheusRuleSpecGroupListResponse struct {
		TotalCount         int                                                    `json:"total_count,omitempty" description:"总数"`
		List               []*PrometheusRuleSpecGroup                             `json:"list,omitempty" description:"规则组列表"`
		EachRangeListState []EachResourceRangeListState[*PrometheusRuleSpecGroup] `json:"each_range_list_state,omitempty" description:"每个集群返回状态"`
	}
	BuiltinPrometheusRuleSpecGroupRuleRequest struct {
		ResourceRequest
		GroupName string            `json:"group_name" binding:"required"`
		RuleList  []*PrometheusRule `json:"rule_list" binding:"required,dive"`
	}

	BuiltinPrometheusRuleSpecGroupRuleListRequest struct {
		ResourceListRequest
		GroupName string `json:"group_name" form:"group_name" binding:"required"`
		RuleName  string `json:"rule_name" form:"rule_name" binding:"required"`
	}

	BuiltinPrometheusRuleSpecGroupRuleListResponse struct {
		TotalCount int               `json:"total_count,omitempty" description:"总数"`
		List       []*PrometheusRule `json:"rule_list,omitempty" description:"规则列表"`
	}

	BuiltinPrometheusRuleSpecGroupRuleGetRequest struct {
		ResourceRequest
		GroupName string `json:"group_name" binding:"required"`
		RuleName  string `json:"rule_name" binding:"required"`
	}

	BuiltinPrometheusRuleSpecGroupRuleGetResponse struct {
		TotalCount int               `json:"total_count,omitempty"`
		List       []*PrometheusRule `json:"list,omitempty"`
	}
)

type (
	PrometheusRuleAlertRequest struct {
		ResourceRequest
		PrometheusRule *monitorv1.PrometheusRule `json:"prometheus_rule" binding:"required"`
	}
)
