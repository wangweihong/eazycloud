package iapiserver

// monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
// monitorv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
// v1 "github.com/prometheus/client_golang/api/prometheus/v1"

type MonitoringTimeRange struct {
	StartTime int64 `json:"start_time" form:"start_time"` //起始时间
	EndTime   int64 `json:"end_time"   form:"end_time"`   //结束时间
	Duration  int64 `json:"duration"   form:"duration"`   //采集间隔。秒为单位
}

// type ServiceMonitorRequest struct {
// 	Request
// 	ServiceMonitor *monitorv1.ServiceMonitor `json:"service_monitor" description:"ServiceMonitor资源"`
// }

// type ServiceMonitorGetRequest struct {
// 	GetRequest
// }

// type ServiceMonitorListRequest struct {
// 	ListRequest
// }

// type ServiceMonitorResponse struct {
// 	Status error      `json:"status" description:"状态码"`
// 	Info   *ServiceMonitorInfo `json:"info" description:"详情"`
// }

// type ServiceMonitorListResponse struct {
// 	Status     error        `json:"status" description:"状态码"`
// 	TotalCount int                   `json:"total_count" description:"总数"`
// 	List       []*ServiceMonitorInfo `json:"list" description:"列表"`
// }

// type ServiceMonitorInfo struct {
// 	ServiceMonitor *monitorv1.ServiceMonitor `json:"service_monitor" description:"ServiceMonitor资源"`
// }

// type PodMonitorRequest struct {
// 	Request
// 	PodMonitor *monitorv1.PodMonitor `json:"pod_monitor" description:"PodMonistor资源"`
// }

// type PodMonitorListRequest struct {
// 	ListRequest
// }

// type PodMonitorGetRequest struct {
// 	GetRequest
// }

// type PodMonitorResponse struct {
// 	Status error  `json:"status" description:"状态码"`
// 	Info   *PodMonitorInfo `json:"info" description:"详情"`
// }

// type PodMonitorListResponse struct {
// 	Status     error    `json:"status" description:"状态码"`
// 	TotalCount int               `json:"total_count" description:"总数"`
// 	List       []*PodMonitorInfo `json:"list" description:"列表"`
// }

// type PodMonitorInfo struct {
// 	PodMonitor *monitorv1.PodMonitor `json:"pod_monitor" description:"PodMonistor资源"`
// 	CTime      int64
// }

// const (
// 	MetricTypeMatrix = "matrix"
// 	MetricTypeVector = "vector"
// )

// type ProqlOption struct {
// 	Level     string
// 	Node      string
// 	Namespace string
// 	PodName   string
// 	Metric    string
// 	Paras     []string

// 	Time      time.Time
// 	StartTime time.Time
// 	EndTime   time.Time
// 	Duration  time.Duration
// }

// func (m *ProqlOption) String() string {
// 	if m == nil {
// 		return ""
// 	}
// 	return fmt.Sprintf("level[%v],node[%v],namespace[%v],podname[%v]metric[%v]",
// 		m.Level, m.Node, m.Namespace, m.PodName, m.Metric)
// }

// type PrometheusConfig struct {
// 	Address string `json:"address"`
// }

// //type Range struct {
// //	Start int64
// //	End   int64
// //	Step  int64 // 相邻指标间隔时间
// //}

// type MonitoringParam struct {
// 	MonitoringTimeRange
// 	Level    string `json:"level" form:"level"`
// 	Time     int64  `json:"time" form:"time"` //某个时间点
// 	Metrics  string `json:"metrics" form:"metrics"`
// 	Type     string `json:"type" form:"type"` //resourceType, such deployment/statefulset/daemonset
// 	Name     string `json:"name" form:"name"` //resourceName, such as deployment Name
// 	NodeName string `json:"node" form:"node"`
// }

// func (p MonitoringParam) ToMonitoring() *Monitoring {
// 	var metrics []string
// 	if p.Metrics != "" {
// 		metrics = strings.Split(p.Metrics, ",")
// 	}

// 	return &Monitoring{
// 		MonitoringTimeRange: p.MonitoringTimeRange,
// 		Level:               p.Level,
// 		Time:                p.Time,
// 		Metrics:             metrics,
// 		Type:                p.Type,
// 		Name:                p.Name,
// 		NodeName:            p.NodeName,
// 	}
// }

// type Monitoring struct {
// 	MonitoringTimeRange
// 	Level string `json:"level" form:"level"`
// 	Time  int64  `json:"time" form:"time"`

// 	Metrics  []string `json:"metrics" form:"metrics"`
// 	Type     string   `json:"type" form:"type"` //resourceType, such deployment/statefulset/daemonset
// 	Name     string   `json:"name" form:"name"` //resourceName, such as deployment Name
// 	NodeName string   `json:"node" form:"node"`
// }

// type MonitoringTimeRange struct {
// 	StartTime int64 `json:"start_time" form:"start_time"` //起始事件
// 	EndTime   int64 `json:"end_time" form:"end_time"`     //结束事件
// 	Duration  int64 `json:"duration" form:"duration"`     //采集间隔。秒为单位
// }

// //  v1.Range要求(endTime-startTime)/step < 11000, 防止数据点过大
// // duration是纳秒!
// func (m MonitoringTimeRange) ToRange() v1.Range {
// 	if m.StartTime == 0 && m.EndTime == 0 {
// 		d, _ := time.ParseDuration("-30m")
// 		m.EndTime = time.Now().Unix()
// 		m.StartTime = time.Now().Add(d).Unix()
// 	}

// 	//1 second
// 	if m.Duration == 0 {
// 		m.Duration = 1
// 	}

// 	return v1.Range{Start: time.Unix(m.StartTime, 0), End: time.Unix(m.EndTime, 0), Step: time.Duration(m.Duration) *
// time.Second}
// }

// type Alert struct {
// 	TenantUuid  string    `json:"tenant_uuid"`
// 	Tenant      string    `json:"tenant"`
// 	ClusterUuid string    `json:"cluster_uuid"`
// 	ClusterName string    `json:"cluster_name"`
// 	Namespace   string    `json:"namespace"`
// 	Name        string    `json:"name"`
// 	Pod         string    `json:"pod"`
// 	Container   string    `json:"container"`
// 	Alertstate  string    `json:"alertstate"`
// 	Severity    string    `json:"severity"`
// 	ActiveTime  int64     `json:"active_time"`
// 	Message     string    `json:"message"`
// 	Value       string    `json:"value"`
// 	Metadata    *v1.Alert `json:"metadata"`
// }

// type MonitoringRequest struct {
// 	Request
// 	FilterMaster    bool        `json:"filter_master"`
// 	Monitoring      *Monitoring `json:"monitoring"`
// 	AlertDeleteList []*v1.Alert `json:"alert_delete_list"`
// 	Namespace       string      `json:"namespace"`
// }

// type MonitoringGetRequest struct {
// 	MonitoringTimeRange
// 	ClusterUUID string `json:"topke_cluster_uuid" form:"topke_cluster_uuid"`
// 	Namespace   string `json:"namespace" form:"namespace"` // only use in list action
// 	Name        string `json:"name" form:"name"`
// }

// type MonitoringListRequest struct {
// 	ListRequest
// 	MonitoringParam
// 	FilterMaster bool `json:"filter_master" form:"filter_master"`
// }

// type MonitoringResponse struct {
// 	Status             error               `json:"status"`
// 	TotalCount         int                          `json:"total_count"`
// 	Metrics            []Metric                     `json:"metrics,omitempty"`
// 	Metric             *Metric                      `json:"metric，omitempty"`
// 	NodeUsedTotal      *NodeOverview                `json:"node_used_total,omitempty"`
// 	MonitorDatas       []*MonitorData               `json:"namespace_top,omitempty"` // 命名空间使用量排行
// 	NodeList           []*NodeOverview              `json:"node_list,omitempty"`     // 节点使用量排行列表
// 	PodList            []*PodOverview               `json:"pod_list,omitempty"`      // pod资源使用列表
// 	AlertList          []*Alert                     `json:"alert_list,omitempty"`
// 	MetricList         []string                     `json:"metric_list,omitempty"`
// 	MonitorStatus      *MonitorStatus               `json:"monitor_status,omitempty"` //grafana/prometheus组件的状态
// 	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
// }

// type MonitoringListResponse struct {
// 	Status             error               `json:"status"`
// 	TotalCount         int                          `json:"total_count"`
// 	Metrics            []Metric                     `json:"metrics,omitempty"`
// 	Metric             *Metric                      `json:"metric，omitempty"`
// 	NodeUsedTotal      *NodeOverview                `json:"node_used_total,omitempty"`
// 	MonitorDatas       []*MonitorData               `json:"namespace_top,omitempty"` // 命名空间使用量排行
// 	NodeList           []*NodeOverview              `json:"node_list,omitempty"`     // 节点使用量排行列表
// 	PodList            []*PodOverview               `json:"pod_list,omitempty"`      // pod资源使用列表
// 	AlertList          []*Alert                     `json:"alert_list,omitempty"`
// 	MetricList         []string                     `json:"metric_list,omitempty"`
// 	MonitorStatus      *MonitorStatus               `json:"monitor_status,omitempty"` //grafana/prometheus组件的状态
// 	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
// }

// type MonitoringApiserverResponse struct {
// 	Status     error `json:"status"`
// 	TimeRange  *v1.Range      `json:"time_range,omitempty"`
// 	CpuUsed    MetricOne      `json:"cpu_used"`
// 	MemoryUsed MetricOne      `json:"memory_used"`
// 	Goroutine  MetricOne      `json:"goroutine"`
// }

// type MonitoringControllerManagerResponse struct {
// 	Status           error `json:"status"`
// 	TimeRange        *v1.Range      `json:"time_range,omitempty"`
// 	UpCount          int            `json:"up_count"`
// 	WorkQueueAddRate MetricOne      `json:"work_queue_add_rate"`
// 	CpuUsed          MetricOne      `json:"cpu_used"`
// 	MemoryUsed       MetricOne      `json:"memory_used"`
// 	Goroutine        MetricOne      `json:"goroutine"`
// }

// type MonitoringKubeletResponse struct {
// 	Status     error `json:"status"`
// 	TimeRange  *v1.Range      `json:"time_range,omitempty"`
// 	UpCount    int            `json:"up_count"`
// 	PodCount   int            `json:"pod_count"`
// 	RpcRate    MetricOne      `json:"rpc_rate"`
// 	CpuUsed    MetricOne      `json:"cpu_used"`
// 	MemoryUsed MetricOne      `json:"memory_used"`
// 	Goroutine  MetricOne      `json:"goroutine"`
// }

// type MonitoringSchedulerResponse struct {
// 	Status        error `json:"status"`
// 	TimeRange     *v1.Range      `json:"time_range"`
// 	UpCount       int            `json:"up_count"`
// 	ScheduleCount MetricOne      `json:"schedule_count"`
// 	CpuUsed       MetricOne      `json:"cpu_used"`
// 	MemoryUsed    MetricOne      `json:"memory_used"`
// 	Goroutine     MetricOne      `json:"goroutine"`
// }

// type MonitoringEtcdResponse struct {
// 	Status               error `json:"status"`
// 	TimeRange            *v1.Range      `json:"time_range,omitempty"`
// 	HasLeader            bool           `json:"has_leader"`
// 	LeaderChangeCount    int            `json:"leader_change_count"`
// 	NetReadByte          MetricOne      `json:"net_read_byte"`
// 	NetWriteByte         MetricOne      `json:"net_write_byte"`
// 	MemorySize           MetricOne      `json:"memory_size"`
// 	RaftProposalApplied  MetricOne      `json:"raft_proposal_applied"`
// 	RaftProposalCommited MetricOne      `json:"raft_proposal_commited"`
// 	RaftProposalFailed   MetricOne      `json:"raft_proposal_failed"`
// 	RaftProposalPending  MetricOne      `json:"raft_proposal_pending"`
// }

// type MonitorStatus struct {
// 	GrafanaEnable    bool   `json:"grafana_enable"`
// 	GrafanaUrl       string `json:"grafana_url"`
// 	GrafanaError     string `json:"grafana_error"`
// 	PrometheusEnable bool   `json:"prometheus_enable"`
// 	PrometheusUrl    string `json:"prometheus_url"`
// 	PrometheusError  string `json:"prometheus_error"`
// }

// type PodOverview struct {
// 	Namespace  string            `json:"namespace" description:"命名空间"`
// 	Name       string            `json:"name" description:"名称"`
// 	CreateTime int64             `json:"create_time" description:"创建时间"`
// 	NodeName   string            `json:"node_name" description:"节点名"`
// 	NodeAddr   string            `json:"node_addr" description:"节点地址"`
// 	Extra      map[string]string `json:"extra" description:"额外数据"`
// 	HardWareResource
// 	//use to overview
// 	ClusterName string `json:"ClusterName" description:"集群名"`
// }

// type MonitorData struct {
// 	Namespace    string            `json:"namespace"` //
// 	Name         string            `json:"name"`      //
// 	Addr         string            `json:"addr"`
// 	Type         string            `json:"type"`          //
// 	Cpu          float64           `json:"cpu"`           // 使用率
// 	CpuUsed      float64           `json:"cpu_used"`      //
// 	CpuTotal     float64           `json:"cpu_total"`     //
// 	Memory       float64           `json:"memory"`        // 使用率
// 	MemoryUsed   float64           `json:"memory_used"`   //
// 	MemoryTotal  float64           `json:"memory_total"`  //
// 	DiskRatio    float64           `json:"disk_ratio"`    // 使用率
// 	DiskUsed     float64           `json:"disk_used"`     //
// 	DiskTotal    float64           `json:"disk_total"`    //
// 	DiskRead     float64           `json:"disk_read"`     //
// 	DiskWrite    float64           `json:"disk_write"`    //
// 	PodRatio     float64           `json:"pod_ratio"`     // 使用率
// 	PodCount     float64           `json:"pod_count"`     // 当前有多少的pod
// 	PodTotal     float64           `json:"pod_total"`     // 集群总共可以创建多少个pod
// 	NetworkRead  float64           `json:"network_read"`  //
// 	NetworkWrite float64           `json:"network_write"` //
// 	Extra        map[string]string `json:"extra"`
// }

// type MonitorDataResponse struct {
// 	Status       error    `json:"status"`
// 	Namespace    string            `json:"namespace,omitempty"` //
// 	Name         string            `json:"name"`                //
// 	Addr         string            `json:"addr"`
// 	Type         string            `json:"type"`          //
// 	Cpu          float64           `json:"cpu"`           // 使用率
// 	CpuUsed      float64           `json:"cpu_used"`      //
// 	CpuTotal     float64           `json:"cpu_total"`     //
// 	Memory       float64           `json:"memory"`        // 使用率
// 	MemoryUsed   float64           `json:"memory_used"`   //
// 	MemoryTotal  float64           `json:"memory_total"`  //
// 	DiskRatio    float64           `json:"disk_ratio"`    // 使用率
// 	DiskUsed     float64           `json:"disk_used"`     //
// 	DiskTotal    float64           `json:"disk_total"`    //
// 	DiskRead     float64           `json:"disk_read"`     //
// 	DiskWrite    float64           `json:"disk_write"`    //
// 	PodRatio     float64           `json:"pod_ratio"`     // 使用率
// 	PodCount     float64           `json:"pod_count"`     // 当前有多少的pod
// 	PodTotal     float64           `json:"pod_total"`     // 集群总共可以创建多少个pod
// 	NetworkRead  float64           `json:"network_read"`  //
// 	NetworkWrite float64           `json:"network_write"` //
// 	Extra        map[string]string `json:"extra"`
// }

// type Metadata struct {
// 	Metric string `json:"metric"`
// 	Type   string `json:"type"`
// 	Help   string `json:"help"`
// }

// type Metric struct {
// 	MetricData
// 	MetricName string `json:"metric_name"`
// 	Error      string `json:"error"`
// }

// type MetricOne struct {
// 	MetricName string  `json:"metric_name"`
// 	Series     []Point `json:"series"`
// 	Sample     *Point  `json:"sample"`
// 	Error      string  `json:"error"`
// }

// type MetricData struct {
// 	MetricType   string        `json:"metric_type"`
// 	MetricValues []MetricValue `json:"metric_values"`
// }

// type MetricValue struct {
// 	Name     string            `json:"name"`
// 	Max      float64           `json:"max"`
// 	Metadata map[string]string `json:"metadata"`
// 	Sample   *Point            `json:"sample"`
// 	Series   []Point           `json:"series"`
// }
// type Point [2]float64

// func (p Point) Timestamp() float64 {
// 	return p[0]
// }

// func (p Point) Value() float64 {
// 	return p[1]
// }

// func (p Point) MarshalJSON() ([]byte, error) {
// 	t, err := jsoniter.Marshal(p.Timestamp())
// 	if err != nil {
// 		return nil, status.NewStatusDesc(scode.ScodeManagerCommonParameterError, err.Error())
// 	}
// 	v, err := jsoniter.Marshal(strconv.FormatFloat(p.Value(), 'f', -1, 64))
// 	if err != nil {
// 		return nil, status.NewStatusDesc(scode.ScodeManagerCommonParameterError, err.Error())
// 	}
// 	return []byte(fmt.Sprintf("[%s,%s]", t, v)), nil
// }

// func (p *Point) UnmarshalJSON(b []byte) error {
// 	var v []interface{}
// 	if err := jsoniter.Unmarshal(b, &v); err != nil {
// 		return status.NewStatusDesc(scode.ScodeManagerCommonParameterError, err.Error())
// 	}

// 	if v == nil {
// 		return nil
// 	}

// 	if len(v) != 2 {
// 		return status.NewStatusDesc(scode.ScodeManagerCommonParameterError, "unsupported array length")
// 	}

// 	ts, ok := v[0].(float64)
// 	if !ok {
// 		return status.NewStatusDesc(scode.ScodeManagerCommonParameterError, "failed to unmarshal [timestamp]")
// 	}
// 	valstr, ok := v[1].(string)
// 	if !ok {
// 		return status.NewStatusDesc(scode.ScodeManagerCommonParameterError, "failed to unmarshal [value]")
// 	}
// 	valf, err := strconv.ParseFloat(valstr, 64)
// 	if err != nil {
// 		return status.NewStatusDesc(scode.ScodeManagerCommonParameterError, err.Error())
// 	}

// 	p[0] = ts
// 	p[1] = valf
// 	return nil
// }

// type MonitoringTopGetRequest struct {
// 	imachinery.PagingParams
// 	ClusterUUID string `json:"topke_cluster_uuid" form:"topke_cluster_uuid"`
// 	Namespace   string `json:"namespace" form:"namespace"`
// }

// type MonitoringTopGetResponse struct {
// 	Status       error `json:"status"`
// 	TotalCount   int            `json:"total_count"`
// 	MonitorDatas []*MonitorData `json:"namespace_top,omitempty"` // 命名空间使用量排行
// }

// type MonitoringNamespaceListRequest struct {
// 	ListRequest
// }

// type MonitoringNamespaceListResponse struct {
// 	Status             error               `json:"status"`
// 	TotalCount         int                          `json:"total_count"`
// 	MonitorDatas       []*MonitorData               `json:"namespace_top,omitempty"` // 命名空间使用量排行
// 	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
// }

// type MonistoringPodContainerHistoryListRequest struct {
// 	MonitoringTimeRange
// 	ClusterUUID string `json:"topke_cluster_uuid" form:"topke_cluster_uuid"`
// 	Name        string `json:"name" form:"name"` //pod 名

// }

// type MonistoringPodContainerHistoryListResponse struct {
// 	Status    error `json:"status"`
// 	TimeRange *v1.Range      `json:"time_range,omitempty"`
// 	Metrics   []Metric       `json:"metrics,omitempty"`
// }

// type MonitoringNodeGetRequest struct {
// 	MonitoringTimeRange
// 	ClusterUUID string `json:"topke_cluster_uuid" form:"topke_cluster_uuid"`
// 	NodeName    string `json:"node" form:"node"`
// }

// type MonitoringNodeGetResponse struct {
// 	Status    error `json:"status"`
// 	TimeRange *v1.Range      `json:"time_range"`
// 	Metric    *Metric        `json:"metric,omitempty"`
// 	Metrics   []Metric       `json:"metrics,omitempty"`
// }

// type MonitoringPodListRequest struct {
// 	ListRequest
// }

// type MonitoringPodListResponse struct {
// 	Status             error               `json:"status"`
// 	TotalCount         int                          `json:"total_count"`
// 	PodList            []*PodOverview               `json:"pod_list,omitempty"` // pod资源使用列表
// 	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
// }

// type MonitoringMetricsNameListRequest struct {
// 	ClusterUUID string `json:"topke_cluster_uuid" form:"topke_cluster_uuid"`
// }

// type MonitoringMetricsNameListResponse struct {
// 	Status     error `json:"status"`
// 	TotalCount int            `json:"total_count"`
// 	MetricList []string       `json:"metric_list,omitempty"`
// }

// type MonistoringHardwareHistoryListRequest struct {
// 	MonitoringTimeRange
// 	ClusterUUID string `json:"topke_cluster_uuid" form:"topke_cluster_uuid"`
// }

// type MonistoringHardwareHistoryListResponse struct {
// 	Status    error `json:"status"`
// 	TimeRange *v1.Range      `json:"time_range"`
// 	Metrics   []Metric       `json:"metrics,omitempty"`
// }

// type MonitoringMetricsHistoryListResponse struct {
// 	Status    error `json:"status"`
// 	TimeRange *v1.Range      `json:"time_range"`
// 	Metrics   []Metric       `json:"metrics,omitempty"`
// }
