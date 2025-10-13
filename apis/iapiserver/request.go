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
)
type (
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
)

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
	Cluster   string `json:"cluster"   form:"cluster"`
	Namespace string `json:"namespace" form:"namespace"`
	Name      string `json:"name"      form:"name"`
	Yaml      bool   `json:"yaml"      form:"yaml"`
}

func (r ResourceGetRequest) ToGetOpts() metav1.GetOptions {
	return metav1.GetOptions{}
}

type PodRequest struct {
	ResourceRequest
	// MonitoringTimeRange
	// ShowMonitorData bool             `json:"show_monitor_data" form:"show_monitor_data"`
	PodLogOptions v1.PodLogOptions `json:"pod_log_options"`
	Resource      *v1.Pod          `json:"resource"        binding:"required,namespaced"`
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

type PodListRequest struct {
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

type PodLogRequest struct {
	ResourceGetRequest           `       binding:"namespaced"`
	Container                    string `                     json:"container"                    form:"container"`
	TailLines                    int64  `                     json:"tail_lines"                   form:"tail_lines"`
	Follow                       bool   `                     json:"follow"                       form:"follow"`
	TimeStamps                   bool   `                     json:"timestamps"                   form:"timestamps"`
	InsecureSkipTLSVerifyBackend bool   `                     json:"insecureSkipTLSVerifyBackend" form:"insecureSkipTLSVerifyBackend"`
	SinceSeconds                 int64  `                     json:"sinceSeconds"                 form:"sinceSeconds"`
}

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

type PodGetRequest struct {
	ResourceGetRequest
}

type PodBatchRequest struct {
	Resources []*PodRequest `json:"resources" binding:"dive"`
}

type ObjectTypeMeta struct {
	ObjectMeta metav1.ObjectMeta `json:"object_meta"`
	TypeMeta   metav1.TypeMeta   `json:"type_meta"`
}

type PodInfo struct {
	Resource        *v1.Pod             `json:"resource"`
	Metric          []Metric            `json:"metric,omitempty"`
	ResourceRequest v1.ResourceList     `json:"resource_request,omitempty"`
	ResourceLimit   v1.ResourceList     `json:"resource_limit,omitempty"`
	PodStatus       *PodStatus          `json:"pod_status,omitempty"`
	MaxRestarts     *int                `json:"max_restarts,omitempty"`
	Controller      *ObjectTypeMeta     `json:"controller,omitempty"`
	NodeInfo        *v1.ObjectReference `json:"node_info,omitempty"`
	ResourceConvert []*ResourceConvert  `json:"resource_convert,omitempty"`
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
	TotalCount         int                          `json:"total_count"`
	List               []*PodInfo                   `json:"list"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
}

type PodListResponse struct {
	TotalCount         int                          `json:"total_count"`
	List               []*PodInfo                   `json:"list"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
}

type PodLog struct {
	Msg string `json:"msg"`
}

type PodLogResponse struct {
	LogList    []ikubernetes.PodLog `json:"log_list"`
	TotalCount int                  `json:"total_count"`
	EndTime    int64                `json:"end_time"`
}

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

type ConfigMapRequest struct {
	ResourceRequest
	Resource *v1.ConfigMap `json:"resource" binding:"required,namespaced"`
}

type ConfigMapListRequest struct {
	ResourceListRequest
}

type ConfigMapGetRequest struct {
	ResourceGetRequest `binding:"namespace_resource"`
}

type ConfigMapBatchRequest struct {
	ConfigMaps []*ConfigMapRequest `json:"configmaps"`
}

type CronJobRequest struct {
	ResourceRequest
	Resource *batchv1.CronJob `json:"resource"`
}

type CronJobListRequest struct {
	ResourceListRequest
}

type CronJobGetRequest struct {
	ResourceGetRequest `binding:"namespaced"`
}

type CronJobBatchRequest struct {
	Resources []*CronJobRequest `json:"resources" binding:"dive"`
}

type CronJobResponse struct {
	Info *CronJobInfo `json:"info,omitempty"`
}

type CronJobListResponse struct {
	Status             error                        `json:"status"`
	TotalCount         int                          `json:"total_count"`
	List               []*CronJobInfo               `json:"list"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
}

type CronJobInfo struct {
	Resource        *batchv1.CronJob   `json:"resource"`
	ResourceConvert []*ResourceConvert `json:"resource_convert,omitempty"`
}

type NetworkPolicyRequest struct {
	ResourceRequest
	Resource *networkingv1.NetworkPolicy `json:"resource" binding:"required,namespaced"`
}

type NetworkPolicyListRequest struct {
	ResourceListRequest
}

type NetworkPolicyGetRequest struct {
	ResourceGetRequest
}

type NetworkPolicyBatchRequest struct {
	Resources []NetworkPolicyRequest `json:"resources" binding:"dive"`
}

type NetworkPolicyResponse struct {
	Status error              `json:"status"`
	Info   *NetworkPolicyInfo `json:"info"`
}

type NetworkPolicyListResponse struct {
	Status             error                        `json:"status"`
	TotalCount         int                          `json:"total_count"`
	List               []*NetworkPolicyInfo         `json:"list"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
}

type NetworkPolicyInfo struct {
	Resource *networkingv1.NetworkPolicy `json:"resource"`
}

type IngressRequest struct {
	ResourceRequest
	Resource *networkingv1.Ingress `json:"resource" binding:"required,namespaced"`
}

type IngressGetRequest struct {
	ResourceGetRequest
}

type IngressListRequest struct {
	ResourceListRequest
}

type IngressBatchRequest struct {
	Resources []*IngressRequest `json:"resources" binding:"dive"`
}

type IngressResponse struct {
	Status error        `json:"status"`
	Info   *IngressInfo `json:"info"`
}

type IngressListResponse struct {
	Status             error                        `json:"status"`
	TotalCount         int                          `json:"total_count"`
	List               []*IngressInfo               `json:"list"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
}

type IngressInfo struct {
	Resource *networkingv1.Ingress `json:"resource"`
	// 控制器
	Controller v1.ObjectReference `json:"controller"`
	// 控制器是否健康
	ControllerHealthy bool `json:"controller_healthy"`
}

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
	Status             error                        `json:"status"`
	TotalCount         int                          `json:"total_count"`
	List               []*EventInfo                 `json:"list"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
}

type EventInfo struct {
	Resource *v1.Event `json:"resource"`
}

type DaemonSetRequest struct {
	ResourceRequest
	Resource *appsv1.DaemonSet `json:"resource" binding:"required,namespaced"`
	Version  string            `json:"version"`
}

type DaemonSetListRequest struct {
	ResourceListRequest
}

type DaemonSetGetRequest struct {
	ResourceGetRequest `binding:"namespaced"`
}

type DaemonSetBatchRequest struct {
	Resources []*DaemonSetRequest `json:"resources" binding:"dive"`
}

type DaemonSetResponse struct {
	Status error          `json:"status"`
	Info   *DaemonSetInfo `json:"info"`
}

type DaemonSetVersionListRequest struct {
	ResourceGetRequest
	imachinery.PagingParams
}

type DaemonSetVersionListResponse struct {
	// 版本列表
	VersionList []VersionInfo `json:"version_list"`
	// 当前版本
	CurrentVersion *VersionInfo `json:"current_version"`
	TotalCount     int          `json:"total_count"`
}

type DaemonSetListResponse struct {
	Status             error                        `json:"status"`
	TotalCount         int                          `json:"total_count"`
	List               []*DaemonSetInfo             `json:"list"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
}

type DaemonSetInfo struct {
	Resource        *appsv1.DaemonSet  `json:"resource"`
	ResourceConvert []*ResourceConvert `json:"resource_convert,omitempty"`
}

type VersionInfo struct {
	Namespace   string              `json:"namespace"`
	Name        string              `json:"name"`
	Version     string              `json:"version"`
	CreateTime  int64               `json:"create_time"`
	StatefulSet *appsv1.StatefulSet `json:"statefulset,omitempty"`
	Deployment  *appsv1.Deployment  `json:"deployment,omitempty"`
	DaemonSet   *appsv1.DaemonSet   `json:"daemonset,omitempty"`
}

type ReplicaSetRequest struct {
	ResourceRequest
	Resource *appsv1.ReplicaSet `json:"resource" binding:"required,namespaced"`
}

type ReplicaSetListRequest struct {
	ResourceListRequest
}

type ReplicaSetGetRequest struct {
	ResourceGetRequest
}

type ReplicaSetResponse struct {
	Status error           `json:"status"`
	Info   *ReplicaSetInfo `json:"info"`
}

type ReplicaSetListResponse struct {
	Status             error                        `json:"status"`
	TotalCount         int                          `json:"total_count"`
	List               []*ReplicaSetInfo            `json:"list"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
}

type ReplicaSetInfo struct {
	Resource *appsv1.ReplicaSet `json:"resource"`
}

type HpaRequest struct {
	ResourceRequest
	HorizontalPodAutoscaler *autoscalingv1.HorizontalPodAutoscaler `json:"horizontal_pod_autoscaler" binding:"required,namespaced"`
}

type HpaListRequest struct {
	ResourceListRequest
	FilterDeploymentName string `json:"filter_deployment_name" description:"过滤deployment名"`
}

type HpaGetRequest struct {
	ResourceGetRequest
}

type HpaResponse struct {
	Status error    `json:"status"`
	Info   *HpaInfo `json:"info"`
}

// type HpaListResponse struct {
// 	Status             error    `json:"status"`
// 	TotalCount         int `json:"total_count"`
// 	List               []*HpaInfo                   `json:"list"  `
// 	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
// }

type HpaInfo struct {
	Resource *autoscalingv1.HorizontalPodAutoscaler `json:"resource"`
}

type StatefulSetRequest struct {
	ResourceRequest
	Resource *appsv1.StatefulSet `json:"resource" binding:"required,namespaced"`
	Version  string              `json:"version"`
}

type StatefulSetListRequest struct {
	ResourceListRequest
}

type StatefulSetGetRequest struct {
	ResourceGetRequest
}

type StatefulSetVersionListRequest struct {
	imachinery.PagingParams
	ResourceGetRequest
}

type StatefulSetBatchRequest struct {
	Resources []*StatefulSetRequest `json:"resources" binding:"dive"`
}

type StatefulSetListResponse struct {
	Status             error                        `json:"status"`
	TotalCount         int                          `json:"total_count"`
	List               []*StatefulSetInfo           `json:"list"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
}

type StatefulSetVersionListResponse struct {
	Status         error            `json:"status"`
	TotalCount     int              `json:"total_count"`
	VersionList    []VersionInfo    `json:"version_list"              description:"版本列表"`
	CurrentVersion *VersionInfo     `json:"current_version,omitempty" description:"当前版本"`
	Info           *StatefulSetInfo `json:"info,omitempty"`
}

type StatefulSetResponse struct {
	Status error            `json:"status"`
	Info   *StatefulSetInfo `json:"info,omitempty"`
}

type StatefulSetInfo struct {
	Resource        *appsv1.StatefulSet `json:"resource"`
	ResourceConvert []*ResourceConvert  `json:"resource_convert,omitempty"`
}

type PodDisruptionBudgetRequest struct {
	ResourceRequest
	PodDisruptionBudget *policyv1.PodDisruptionBudget `json:"pod_disruption_budget" binding:"required,namespaced"`
}

type PodDisruptionBudgetGetRequest struct {
	ResourceGetRequest
}

type PodDisruptionBudgetResponse struct {
	Status error                    `json:"status"`
	Info   *PodDisruptionBudgetInfo `json:"info,omitempty"`
}

type PodDisruptionBudgetListRequest struct {
	ResourceListRequest
}

type PodDisruptionBudgetListResponse struct {
	Status             error                        `json:"status"`
	TotalCount         int                          `json:"total_count"`
	List               []*PodDisruptionBudgetInfo   `json:"list"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
}

type PodDisruptionBudgetBatchRequest struct {
	Resources []*PodDisruptionBudgetRequest `json:"resources" binding:"dive" description:"应用中断干扰"`
}

type PodDisruptionBudgetInfo struct {
	*policyv1.PodDisruptionBudget
}

type StorageClassRequest struct {
	ResourceRequest
	StorageClass *storagev1.StorageClass `json:"storage_class" binding:"required,clusterd"`
	Secret       *v1.Secret              `json:"secret"` // secret for provisioner storageclass
}

type StorageClassListRequest struct {
	ResourceListRequest
	StorageClass *storagev1.StorageClass `json:"storage_class"`
	Secret       *v1.Secret              `json:"secret"` // secret for provisioner storageclass
}

type StorageClassGetRequest struct {
	ResourceGetRequest
}

type StorageClassBatchRequest struct {
	Resources []*StorageClassRequest `json:"resources" binding:"dive"`
}

// cluster service response
type StorageClassResponse struct {
	Status error             `json:"status"`
	Info   *StorageClassInfo `json:"info,omitempty"`
}

type StorageClassListResponse struct {
	Status             error                        `json:"status"`
	TotalCount         int                          `json:"total_count"`
	List               []*StorageClassInfo          `json:"list"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
}

type StorageClassInfo struct {
	*storagev1.StorageClass
}

type VolumeSnapshotClassRequest struct {
	ResourceRequest
	VolumeSnapshotClass *snapshotv1beta1.VolumeSnapshotClass `json:"volume_snapshot_class" binding:"required,clusterd"`
}

type VolumeSnapshotClassListRequest struct {
	ResourceListRequest
}

type VolumeSnapshotClassGetRequest struct {
	ResourceGetRequest
}

type VolumeSnapshotClassBatchRequest struct {
	Resources []*VolumeSnapshotClassRequest `json:"resources" binding:"dive"`
}

// cluster service response
type VolumeSnapshotClassResponse struct {
	Status error                    `json:"status"`
	Info   *VolumeSnapshotClassInfo `json:"info,omitempty"`
}

type VolumeSnapshotClassListResponse struct {
	Status             error                        `json:"status"`
	TotalCount         int                          `json:"total_count"`
	List               []*VolumeSnapshotClassInfo   `json:"list,omitempty"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
}

type VolumeSnapshotClassInfo struct {
	*snapshotv1beta1.VolumeSnapshotClass `json:"volume_snapshot_class"`
}

type VolumeSnapshotRequest struct {
	ResourceRequest
	VolumeSnapshot                   *snapshotv1beta1.VolumeSnapshot `json:"volume_snapshot"`
	PersistentVolumeClaimName        string                          `json:"persistent_volume_claim_name"`
	PersistentVolumeClaimAccessModes []v1.PersistentVolumeAccessMode `json:"persistent_volume_claim_access_modes"`
}

type VolumeSnapshotListRequest struct {
	ResourceListRequest
	VolumeSnapshot            *snapshotv1beta1.VolumeSnapshot `json:"volume_snapshot"`
	PersistentVolumeClaimName string                          `json:"persistent_volume_claim_name"`
}

type VolumeSnapshotGetRequest struct {
	ResourceGetRequest
}

type VolumeSnapshotBatchRequest struct {
	VolumeSnapshots []*VolumeSnapshotRequest `json:"volume_snapshots"`
}

// cluster service response
type VolumeSnapshotListResponse struct {
	Status             error                        `json:"status"`
	TotalCount         int                          `json:"total_count"`
	List               []*VolumeSnapshotInfo        `json:"list,omitempty"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
}

type VolumeSnapshotResponse struct {
	Status error               `json:"status"`
	Info   *VolumeSnapshotInfo `json:"info,omitempty"`
}

type VolumeSnapshotInfo struct {
	*snapshotv1beta1.VolumeSnapshot
}

type VolumeSnapshotContentRequest struct {
	ResourceRequest
	Resource *snapshotv1beta1.VolumeSnapshotContent `json:"resource" binding:"required,clusterd"`
}

type VolumeSnapshotContentListRequest struct {
	ResourceListRequest
}

type VolumeSnapshotContentGetRequest struct {
	ResourceGetRequest `binding:"clusterd"`
}

type VolumeSnapshotContentBatchRequest struct {
	Resources []*VolumeSnapshotContentRequest `json:"resources" binding:"dive"`
}

type VolumeSnapshotContentResponse struct {
	Status error                      `json:"status"`
	Info   *VolumeSnapshotContentInfo `json:"info,omitempty"`
}

type VolumeSnapshotContentListResponse struct {
	Status             error                        `json:"status"`
	TotalCount         int                          `json:"total_count"`
	List               []*VolumeSnapshotContentInfo `json:"list,omitempty"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
}

type VolumeSnapshotContentInfo struct {
	Resource *snapshotv1beta1.VolumeSnapshotContent `json:"resource"`
}

type PersistentVolumeListRequest struct {
	ResourceListRequest
}

type PersistentVolumeGetRequest struct {
	ResourceGetRequest `binding:"clusterd"`
}

type PersistentVolumeRequest struct {
	ResourceRequest
	Resource *v1.PersistentVolume `json:"resource" binding:"required,clusterd"`
}

type PersistentVolumeBatchRequest struct {
	Resources []*PersistentVolumeRequest `json:"resources" binding:"dive"`
}

type PersistentVolumeResponse struct {
	Status error                 `json:"status"`
	Info   *PersistentVolumeInfo `json:"info"`
}

type PersistentVolumeListResponse struct {
	Status             error                        `json:"status"`
	TotalCount         int                          `json:"total_count"`
	List               []*PersistentVolumeInfo      `json:"list"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
}

type PersistentVolumeInfo struct {
	Resource *v1.PersistentVolume `json:"resource,omitempty"`
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
	Status             error                        `json:"status"`
	TotalCount         int                          `json:"total_count"`
	List               []*PersistentVolumeClaimInfo `json:"list"`
	EachRangeListState []EachResourceRangeListState `json:"each_range_list_state,omitempty"`
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
