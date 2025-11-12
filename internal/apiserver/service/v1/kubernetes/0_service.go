package kubernetes

import (
	"context"
	"io"
	"time"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/apis/imachinery"
	"github.com/wangweihong/eazycloud/internal/apiserver/store"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/sliceutil"
	"github.com/wangweihong/gotoolbox/pkg/waitgroup"
)

type KubernetesSrv interface {
	PodList(ctx context.Context, req *iapiserver.PodListRequest) (*iapiserver.PodListResponse, error)
	PodGet(ctx context.Context, req *iapiserver.ResourceGetRequest) (*iapiserver.PodInfo, error)
	PodAdd(ctx context.Context, req *iapiserver.PodRequest) (*iapiserver.PodInfo, error)
	PodDelete(ctx context.Context, req *iapiserver.PodRequest) error
	PodBatchDelete(ctx context.Context, req *iapiserver.PodBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.PodRequest]
	PodLogList(ctx context.Context, req *iapiserver.PodLogRequest) (*iapiserver.PodLogResponse, error)
	PodLogStream(ctx context.Context, req *iapiserver.PodLogRequest) (io.ReadCloser, error)
	PodPatch(ctx context.Context, req *iapiserver.PodRequest) (*iapiserver.PodInfo, error)
	PodEvict(ctx context.Context, req *iapiserver.PodRequest) error
	GetComponentPod(ctx context.Context, req *iapiserver.PodListRequest) (*iapiserver.PodListResponse, error)

	ServiceList(ctx context.Context, req *iapiserver.ServiceListRequest) (*iapiserver.ServiceListResponse, error)
	ServiceGet(ctx context.Context, req *iapiserver.ServiceGetRequest) (*iapiserver.ServiceInfo, error)
	ServiceAdd(ctx context.Context, req *iapiserver.ServiceRequest) (*iapiserver.ServiceInfo, error)
	ServiceDelete(ctx context.Context, req *iapiserver.ServiceRequest) error
	ServiceBatchDelete(ctx context.Context, req *iapiserver.ServiceBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.ServiceRequest]
	ServiceUpdate(ctx context.Context, req *iapiserver.ServiceRequest) (*iapiserver.ServiceInfo, error)

	NodePortList(ctx context.Context, req *iapiserver.NodePortListRequest) (*iapiserver.NodePortListResponse, error)
	IsNodePortUsed(ctx context.Context, req *iapiserver.NodePortRequest) error

	ResourceQuotaAdd(ctx context.Context, req *iapiserver.ResourceQuotaRequest) (*iapiserver.ResourceQuotaInfo, error)
	ResourceQuotaDelete(ctx context.Context, req *iapiserver.ResourceQuotaRequest) error
	ResourceQuotaGet(ctx context.Context, req *iapiserver.ResourceQuotaGetRequest) (*iapiserver.ResourceQuotaInfo, error)
	ResourceQuotaList(ctx context.Context, req *iapiserver.ResourceQuotaListRequest) (*iapiserver.ResourceQuotaListResponse, error)
	ResourceQuotaUpdate(ctx context.Context, req *iapiserver.ResourceQuotaRequest) (*iapiserver.ResourceQuotaInfo, error)

	LimitRangeAdd(ctx context.Context, req *iapiserver.LimitRangeRequest) (*iapiserver.LimitRangeInfo, error)
	LimitRangeDelete(ctx context.Context, req *iapiserver.LimitRangeRequest) error
	LimitRangeGet(ctx context.Context, req *iapiserver.LimitRangeGetRequest) (*iapiserver.LimitRangeInfo, error)
	LimitRangeList(ctx context.Context, req *iapiserver.LimitRangeListRequest) (*iapiserver.LimitRangeListResponse, error)
	LimitRangeUpdate(ctx context.Context, req *iapiserver.LimitRangeRequest) (*iapiserver.LimitRangeInfo, error)

	ConfigMapAdd(ctx context.Context, req *iapiserver.ConfigMapRequest) (*iapiserver.ConfigMapInfo, error)
	ConfigMapUpdate(ctx context.Context, req *iapiserver.ConfigMapRequest) (*iapiserver.ConfigMapInfo, error)
	ConfigMapDelete(ctx context.Context, req *iapiserver.ConfigMapRequest) error
	ConfigMapBatchDelete(ctx context.Context, req *iapiserver.ConfigMapBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.ConfigMapRequest]
	ConfigMapGet(ctx context.Context, req *iapiserver.ConfigMapGetRequest) (*iapiserver.ConfigMapInfo, error)
	ConfigMapList(ctx context.Context, req *iapiserver.ConfigMapListRequest) (*iapiserver.ConfigMapListResponse, error)

	SecretAdd(ctx context.Context, req *iapiserver.SecretRequest) (*iapiserver.SecretInfo, error)
	SecretUpdate(ctx context.Context, req *iapiserver.SecretRequest) (*iapiserver.SecretInfo, error)
	SecretDelete(ctx context.Context, req *iapiserver.SecretRequest) error
	SecretBatchDelete(ctx context.Context, req *iapiserver.SecretBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.SecretRequest]
	SecretGet(ctx context.Context, req *iapiserver.SecretGetRequest) (*iapiserver.SecretInfo, error)
	SecretList(ctx context.Context, req *iapiserver.SecretListRequest) (*iapiserver.SecretListResponse, error)

	NodeGatewayUpdate(ctx context.Context, req *iapiserver.NodeRequest) error
	NodeTaintUpdate(ctx context.Context, req *iapiserver.NodeRequest) (*iapiserver.NodeInfo, error)
	NodeUpdateLabel(ctx context.Context, req *iapiserver.NodeRequest) (*iapiserver.NodeInfo, error)
	NodeAdd(ctx context.Context, req *iapiserver.NodeRequest) (*iapiserver.NodeInfo, error)
	NodeUpdate(ctx context.Context, req *iapiserver.NodeRequest) (*iapiserver.NodeInfo, error)
	NodeDelete(ctx context.Context, req *iapiserver.NodeRequest) error
	NodeGet(ctx context.Context, req *iapiserver.NodeGetRequest) (*iapiserver.NodeInfo, error)
	NodeList(ctx context.Context, req *iapiserver.NodeListRequest) (*iapiserver.NodeListResponse, error)
	NodeCordon(ctx context.Context, req *iapiserver.NodeRequest) error
	NodeUncordon(ctx context.Context, req *iapiserver.NodeRequest) error
	NodeDrain(ctx context.Context, req *iapiserver.NodeDrainRequest) (*imachinery.BatchOutput, error)
	NodeEvent(ctx context.Context, req *iapiserver.NodeEventRequest) ([]*iapiserver.EventInfo, error)
	//NodeState(clusterInfo *iapiserver.Cluster, nodeName string, timeout int64) (*iapiserver.MonitorData, error)

	EventGet(ctx context.Context, req *iapiserver.EventGetRequest) (*iapiserver.EventInfo, error)
	EventList(ctx context.Context, req *iapiserver.EventListRequest) (*iapiserver.EventListResponse, error)

	DaemonSetVersionUpdate(ctx context.Context, req *iapiserver.DaemonSetRequest) (*iapiserver.DaemonSetInfo, error)
	DaemonSetVersionList(ctx context.Context, req *iapiserver.DaemonSetVersionListRequest) (*iapiserver.DaemonSetVersionListResponse, error)
	DaemonSetAdd(ctx context.Context, req *iapiserver.DaemonSetRequest) (*iapiserver.DaemonSetInfo, error)
	DaemonSetRecreate(ctx context.Context, req *iapiserver.DaemonSetRequest) (*iapiserver.DaemonSetInfo, error)
	DaemonSetUpdate(ctx context.Context, req *iapiserver.DaemonSetRequest) (*iapiserver.DaemonSetInfo, error)
	DaemonSetDelete(ctx context.Context, req *iapiserver.DaemonSetRequest) error
	DaemonSetBatchDelete(ctx context.Context, req *iapiserver.DaemonSetBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.DaemonSetRequest]
	DaemonSetGet(ctx context.Context, req *iapiserver.DaemonSetGetRequest) (*iapiserver.DaemonSetInfo, error)
	DaemonSetList(ctx context.Context, req *iapiserver.DaemonSetListRequest) (*iapiserver.DaemonSetListResponse, error)

	ReplicaSetAdd(ctx context.Context, req *iapiserver.ReplicaSetRequest) (*iapiserver.ReplicaSetInfo, error)
	ReplicaSetDelete(ctx context.Context, req *iapiserver.ReplicaSetRequest) error
	ReplicaSetUpdate(ctx context.Context, req *iapiserver.ReplicaSetRequest) (*iapiserver.ReplicaSetInfo, error)
	ReplicaSetGet(ctx context.Context, req *iapiserver.ReplicaSetGetRequest) (*iapiserver.ReplicaSetInfo, error)
	ReplicaSetList(ctx context.Context, req *iapiserver.ReplicaSetListRequest) (*iapiserver.ReplicaSetListResponse, error)

	HpaAdd(ctx context.Context, req *iapiserver.HpaRequest) (*iapiserver.HpaInfo, error)
	HpaDelete(ctx context.Context, req *iapiserver.HpaRequest) error
	HpaUpdate(ctx context.Context, req *iapiserver.HpaRequest) (*iapiserver.HpaInfo, error)
	HpaGet(ctx context.Context, req *iapiserver.HpaGetRequest) (*iapiserver.HpaInfo, error)
	HpaList(ctx context.Context, req *iapiserver.HpaListRequest) (*iapiserver.HpaListResponse, error)

	StatefulSetAdd(ctx context.Context, req *iapiserver.StatefulSetRequest) (*iapiserver.StatefulSetInfo, error)
	StatefulSetVersionUpdate(ctx context.Context, req *iapiserver.StatefulSetRequest) (*iapiserver.StatefulSetInfo, error)
	StatefulSetVersionList(ctx context.Context, req *iapiserver.StatefulSetVersionListRequest) (*iapiserver.StatefulSetVersionListResponse, error)
	StatefulSetRecreate(ctx context.Context, req *iapiserver.StatefulSetRequest) (*iapiserver.StatefulSetInfo, error)
	StatefulSetUpdate(ctx context.Context, req *iapiserver.StatefulSetRequest) (*iapiserver.StatefulSetInfo, error)
	StatefulSetDelete(ctx context.Context, req *iapiserver.StatefulSetRequest) error
	StatefulSetBatchDelete(ctx context.Context, req *iapiserver.StatefulSetBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.StatefulSetRequest]
	StatefulSetGet(ctx context.Context, req *iapiserver.StatefulSetGetRequest) (*iapiserver.StatefulSetInfo, error)
	StatefulSetList(ctx context.Context, req *iapiserver.StatefulSetListRequest) (*iapiserver.StatefulSetListResponse, error)

	DeploymentUpdate(ctx context.Context, req *iapiserver.DeploymentRequest) (*iapiserver.DeploymentInfo, error)
	DeploymentAdd(ctx context.Context, req *iapiserver.DeploymentRequest) (*iapiserver.DeploymentInfo, error)
	DeploymentDelete(ctx context.Context, req *iapiserver.DeploymentRequest) error
	DeploymentBatchDelete(ctx context.Context, req *iapiserver.DeploymentBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.DeploymentRequest]
	DeploymentRecreate(ctx context.Context, req *iapiserver.DeploymentRequest) (*iapiserver.DeploymentInfo, error)
	DeploymentGet(ctx context.Context, req *iapiserver.DeploymentGetRequest) (*iapiserver.DeploymentInfo, error)
	DeploymentList(ctx context.Context, req *iapiserver.DeploymentListRequest) (*iapiserver.DeploymentListResponse, error)
	DeploymentVersionUpdate(ctx context.Context, req *iapiserver.DeploymentRequest) error
	DeploymentVersionList(ctx context.Context, req *iapiserver.DeploymentVersionListRequest) (*iapiserver.DeploymentVersionListResponse, error)

	JobAdd(ctx context.Context, req *iapiserver.JobRequest) (*iapiserver.JobInfo, error)
	JobDelete(ctx context.Context, req *iapiserver.JobRequest) error
	JobBatchDelete(ctx context.Context, req *iapiserver.JobBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.JobRequest]
	JobGet(ctx context.Context, req *iapiserver.JobGetRequest) (*iapiserver.JobResponse, error)
	JobList(ctx context.Context, req *iapiserver.JobListRequest) (*iapiserver.JobListResponse, error)

	CronJobAdd(ctx context.Context, req *iapiserver.CronJobRequest) (*iapiserver.CronJobInfo, error)
	CronJobUpdate(ctx context.Context, req *iapiserver.CronJobRequest) (*iapiserver.CronJobInfo, error)
	CronJobUpdateSuspend(ctx context.Context, req *iapiserver.CronJobRequest) (*iapiserver.CronJobInfo, error)
	CronJobUpdateSchedule(ctx context.Context, req *iapiserver.CronJobRequest) (*iapiserver.CronJobInfo, error)
	CronJobDelete(ctx context.Context, req *iapiserver.CronJobRequest) error
	CronJobBatchDelete(ctx context.Context, req *iapiserver.CronJobBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.CronJobRequest]
	CronJobGet(ctx context.Context, req *iapiserver.CronJobGetRequest) (*iapiserver.CronJobInfo, error)
	CronJobList(ctx context.Context, req *iapiserver.CronJobListRequest) (*iapiserver.CronJobListResponse, error)
	CronJobTrigger(ctx context.Context, req *iapiserver.CronJobRequest) (*iapiserver.JobInfo, error)

	StorageClassAdd(ctx context.Context, req *iapiserver.StorageClassRequest) (*iapiserver.StorageClassInfo, error)
	StorageClassUpdate(ctx context.Context, req *iapiserver.StorageClassRequest) (*iapiserver.StorageClassInfo, error)
	StorageClassDelete(ctx context.Context, req *iapiserver.StorageClassRequest) error
	StorageClassBatchDelete(ctx context.Context, req *iapiserver.StorageClassBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.StorageClassRequest]
	StorageClassGet(ctx context.Context, req *iapiserver.StorageClassGetRequest) (*iapiserver.StorageClassInfo, error)
	StorageClassList(ctx context.Context, req *iapiserver.StorageClassListRequest) (*iapiserver.StorageClassListResponse, error)

	PersistentVolumeAdd(ctx context.Context, req *iapiserver.PersistentVolumeRequest) (*iapiserver.PersistentVolumeInfo, error)
	PersistentVolumeUpdate(ctx context.Context, req *iapiserver.PersistentVolumeRequest) (*iapiserver.PersistentVolumeInfo, error)
	PersistentVolumeDelete(ctx context.Context, req *iapiserver.PersistentVolumeRequest) error
	PersistentVolumeBatchDelete(ctx context.Context, req *iapiserver.PersistentVolumeBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.PersistentVolumeRequest]
	PersistentVolumeGet(ctx context.Context, req *iapiserver.PersistentVolumeGetRequest) (*iapiserver.PersistentVolumeInfo, error)
	PersistentVolumeList(ctx context.Context, req *iapiserver.PersistentVolumeListRequest) (*iapiserver.PersistentVolumeListResponse, error)

	PersistentVolumeClaimAdd(ctx context.Context, req *iapiserver.PersistentVolumeClaimRequest) (*iapiserver.PersistentVolumeClaimInfo, error)
	PersistentVolumeClaimUpdate(ctx context.Context, req *iapiserver.PersistentVolumeClaimRequest) (*iapiserver.PersistentVolumeClaimInfo, error)
	PersistentVolumeClaimDelete(ctx context.Context, req *iapiserver.PersistentVolumeClaimRequest) error
	PersistentVolumeClaimBatchDelete(ctx context.Context, req *iapiserver.PersistentVolumeClaimBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.PersistentVolumeClaimRequest]
	PersistentVolumeClaimGet(ctx context.Context, req *iapiserver.PersistentVolumeClaimGetRequest) (*iapiserver.PersistentVolumeClaimInfo, error)
	PersistentVolumeClaimList(ctx context.Context, req *iapiserver.PersistentVolumeClaimListRequest) (*iapiserver.PersistentVolumeClaimListResponse, error)
	PersistentVolumeClaimExpand(ctx context.Context, req *iapiserver.PersistentVolumeClaimRequest) (*iapiserver.PersistentVolumeClaimInfo, error)

	VolumeSnapshotClassAdd(ctx context.Context, req *iapiserver.VolumeSnapshotClassRequest) (*iapiserver.VolumeSnapshotClassInfo, error)
	VolumeSnapshotClassUpdate(ctx context.Context, req *iapiserver.VolumeSnapshotClassRequest) (*iapiserver.VolumeSnapshotClassInfo, error)
	VolumeSnapshotClassDelete(ctx context.Context, req *iapiserver.VolumeSnapshotClassRequest) error
	VolumeSnapshotClassBatchDelete(ctx context.Context, req *iapiserver.VolumeSnapshotClassBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.VolumeSnapshotClassRequest]
	VolumeSnapshotClassGet(ctx context.Context, req *iapiserver.VolumeSnapshotClassGetRequest) (*iapiserver.VolumeSnapshotClassInfo, error)
	VolumeSnapshotClassList(ctx context.Context, req *iapiserver.VolumeSnapshotClassListRequest) (*iapiserver.VolumeSnapshotClassListResponse, error)

	VolumeSnapshotContentAdd(ctx context.Context, req *iapiserver.VolumeSnapshotContentRequest) (*iapiserver.VolumeSnapshotContentInfo, error)
	VolumeSnapshotContentUpdate(ctx context.Context, req *iapiserver.VolumeSnapshotContentRequest) (*iapiserver.VolumeSnapshotContentInfo, error)
	VolumeSnapshotContentDelete(ctx context.Context, req *iapiserver.VolumeSnapshotContentRequest) error
	VolumeSnapshotContentBatchDelete(ctx context.Context, req *iapiserver.VolumeSnapshotContentBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.VolumeSnapshotContentRequest]
	VolumeSnapshotContentGet(ctx context.Context, req *iapiserver.VolumeSnapshotContentGetRequest) (*iapiserver.VolumeSnapshotContentInfo, error)
	VolumeSnapshotContentList(ctx context.Context, req *iapiserver.VolumeSnapshotContentListRequest) (*iapiserver.VolumeSnapshotContentListResponse, error)

	VolumeSnapshotAdd(ctx context.Context, req *iapiserver.VolumeSnapshotRequest) (*iapiserver.VolumeSnapshotInfo, error)
	VolumeSnapshotUpdate(ctx context.Context, req *iapiserver.VolumeSnapshotRequest) (*iapiserver.VolumeSnapshotInfo, error)
	VolumeSnapshotDelete(ctx context.Context, req *iapiserver.VolumeSnapshotRequest) error
	VolumeSnapshotBatchDelete(ctx context.Context, req *iapiserver.VolumeSnapshotBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.VolumeSnapshotRequest]
	VolumeSnapshotGet(ctx context.Context, req *iapiserver.VolumeSnapshotGetRequest) (*iapiserver.VolumeSnapshotInfo, error)
	VolumeSnapshotList(ctx context.Context, req *iapiserver.VolumeSnapshotListRequest) (*iapiserver.VolumeSnapshotListResponse, error)

	NetworkPolicyAdd(ctx context.Context, req *iapiserver.NetworkPolicyRequest) (*iapiserver.NetworkPolicyInfo, error)
	NetworkPolicyUpdate(ctx context.Context, req *iapiserver.NetworkPolicyRequest) (*iapiserver.NetworkPolicyInfo, error)
	NetworkPolicyDelete(ctx context.Context, req *iapiserver.NetworkPolicyRequest) error
	NetworkPolicyBatchDelete(ctx context.Context, req *iapiserver.NetworkPolicyBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.NetworkPolicyRequest]
	NetworkPolicyGet(ctx context.Context, req *iapiserver.NetworkPolicyGetRequest) (*iapiserver.NetworkPolicyInfo, error)
	NetworkPolicyList(ctx context.Context, req *iapiserver.NetworkPolicyListRequest) (*iapiserver.NetworkPolicyListResponse, error)

	IngressAdd(ctx context.Context, req *iapiserver.IngressRequest) (*iapiserver.IngressInfo, error)
	IngressUpdate(ctx context.Context, req *iapiserver.IngressRequest) (*iapiserver.IngressInfo, error)
	IngressDelete(ctx context.Context, req *iapiserver.IngressRequest) error
	IngressBatchDelete(ctx context.Context, req *iapiserver.IngressBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.IngressRequest]
	IngressGet(ctx context.Context, req *iapiserver.IngressGetRequest) (*iapiserver.IngressInfo, error)
	IngressList(ctx context.Context, req *iapiserver.IngressListRequest) (*iapiserver.IngressListResponse, error)

	NamespaceAdd(ctx context.Context, req *iapiserver.NamespaceRequest) (*iapiserver.NamespaceInfo, error)
	NamespaceDelete(ctx context.Context, req *iapiserver.NamespaceRequest) error
	NamespaceGet(ctx context.Context, req *iapiserver.NamespaceGetRequest) (*iapiserver.NamespaceInfo, error)
	NamespaceGatewayGet(ctx context.Context, req *iapiserver.PodListRequest) (*iapiserver.NamespaceGatewayResponse, error)
	NamespaceList(ctx context.Context, req *iapiserver.NamespaceListRequest) (*iapiserver.NamespaceListResponse, error)
	NamespaceUpdate(ctx context.Context, req *iapiserver.NamespaceRequest) (*iapiserver.NamespaceInfo, error)
}

type kubernetesService struct {
	store store.Factory
}

func NewService(st store.Factory) *kubernetesService {
	return &kubernetesService{store: st}
}

func getVisitScope(ctx context.Context, st store.Factory, req iapiserver.ResourceListRequest) ([]*iapiserver.Cluster, error) {
	if req.Cluster != "" {
		cluster, err := st.Kubernetes().Get(ctx, req.Cluster)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return []*iapiserver.Cluster{cluster}, nil
	}
	metas, err := st.Kubernetes().List(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return sliceutil.CopyIf(metas, func(o *iapiserver.Cluster) bool {
		return !o.IsStop
	}), nil
}

func multiClusterResourceList[T any](
	ctx context.Context, st store.Factory, list *[]T, req iapiserver.ResourceListRequest,
	concurrentFunc func(context.Context, *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[T]],
	sortFunc func(i, j int) bool, timeouts ...time.Duration,
) ([]iapiserver.EachResourceRangeListState[T], int, error) {
	clusters, err := getVisitScope(ctx, st, req)
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}
	wg := waitgroup.RunGenericConcurrently(ctx, clusters, concurrentFunc, timeouts...)
	eachClusterResult := wg.GetSuccessResultList()

	totalCount := CutPagingSliceResourceList(eachClusterResult, list, req.PageNum, req.PageSize, sortFunc)
	return eachClusterResult, totalCount, nil
}

func multiClusterResourceList2[T any](
	ctx context.Context, st store.Factory, list *[]T, req iapiserver.ResourceListRequest,
	concurrentFunc func(context.Context, *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[T]],
	sortfields string, sortAsc bool, timeouts ...time.Duration,
) ([]iapiserver.EachResourceRangeListState[T], int, error) {
	clusters, err := getVisitScope(ctx, st, req)
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}
	wg := waitgroup.RunGenericConcurrently(ctx, clusters, concurrentFunc, timeouts...)
	eachClusterResult := wg.GetSuccessResultList()

	totalCount := CutPagingSliceResourceList2(eachClusterResult, list, req.PageNum, req.PageSize, sortfields, sortAsc)
	return eachClusterResult, totalCount, nil
}
