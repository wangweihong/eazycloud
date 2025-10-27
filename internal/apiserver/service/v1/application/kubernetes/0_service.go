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
	PodCreate(ctx context.Context, req *iapiserver.PodRequest) (*iapiserver.PodInfo, error)
	PodGet(ctx context.Context, req *iapiserver.ResourceGetRequest) (*iapiserver.PodInfo, error)
	PodDelete(ctx context.Context, req *iapiserver.PodRequest) error
	PodLogList(ctx context.Context, req *iapiserver.PodLogRequest) (*iapiserver.PodLogResponse, error)
	PodLogStream(ctx context.Context, req *iapiserver.PodLogRequest) (io.ReadCloser, error)
	PodPatch(ctx context.Context, req *iapiserver.PodRequest) (*iapiserver.PodInfo, error)
	PodEvict(ctx context.Context, req *iapiserver.PodRequest) error
	PodList(ctx context.Context, req *iapiserver.PodListRequest) (*iapiserver.PodListResponse, error)
	GetComponentPod(ctx context.Context, req *iapiserver.PodListRequest) (*iapiserver.PodListResponse, error)

	ServiceUpdate(ctx context.Context, req *iapiserver.ServiceRequest) (*v1.Service, error) 
	
	NodeGatewayUpdate(ctx context.Context, req *iapiserver.NodeRequest) error
	NodeTaintUpdate(ctx context.Context, req *iapiserver.NodeRequest) (*iapiserver.NodeInfo, error)
	NodeUpdateLabel(ctx context.Context, req *iapiserver.NodeRequest) (*iapiserver.NodeInfo, error)
	NodeCreate(ctx context.Context, req *iapiserver.NodeRequest) (*iapiserver.NodeInfo, error)
	NodeUpdate(ctx context.Context, req *iapiserver.NodeRequest) (*iapiserver.NodeInfo, error)
	NodeDelete(ctx context.Context, req *iapiserver.NodeRequest) error
	NodeGet(ctx context.Context, req *iapiserver.NodeGetRequest) (*iapiserver.NodeInfo, error)
	NodeList(ctx context.Context, req *iapiserver.NodeListRequest) (*iapiserver.NodeListResponse, error)
	NodeCordon(ctx context.Context, req *iapiserver.NodeRequest) error
	NodeUncordon(ctx context.Context, req *iapiserver.NodeRequest) error
	NodeDrain(ctx context.Context, req *iapiserver.NodeDrainRequest) (*imachinery.BatchOutput, error)
	NodeEvent(ctx context.Context, req *iapiserver.NodeEventRequest) ([]*iapiserver.EventInfo, error)
	//NodeState(clusterInfo *iapiserver.Cluster, nodeName string, timeout int64) (*iapiserver.MonitorData, error)

	DaemonSetVersionUpdate(ctx context.Context, req *iapiserver.DaemonSetRequest) (*iapiserver.DaemonSetInfo, error)
	DaemonSetVersionList(ctx context.Context, req *iapiserver.DaemonSetVersionListRequest) (*iapiserver.DaemonSetVersionListResponse, error)
	DaemonSetCreate(ctx context.Context, req *iapiserver.DaemonSetRequest) (*iapiserver.DaemonSetInfo, error)
	DaemonSetRecreate(ctx context.Context, req *iapiserver.DaemonSetRequest) (*iapiserver.DaemonSetInfo, error)
	DaemonSetUpdate(ctx context.Context, req *iapiserver.DaemonSetRequest) (*iapiserver.DaemonSetInfo, error)
	DaemonSetDelete(ctx context.Context, req *iapiserver.DaemonSetRequest) error
	DaemonSetBatchDelete(ctx context.Context, req *iapiserver.DaemonSetBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.DaemonSetRequest]
	DaemonSetGet(ctx context.Context, req *iapiserver.DaemonSetGetRequest) (*iapiserver.DaemonSetInfo, error)
	DaemonSetList(ctx context.Context, req *iapiserver.DaemonSetListRequest) (*iapiserver.DaemonSetListResponse, error)

	ReplicaSetCreate(ctx context.Context, req *iapiserver.ReplicaSetRequest) (*iapiserver.ReplicaSetInfo, error)
	ReplicaSetDelete(ctx context.Context, req *iapiserver.ReplicaSetRequest) error
	ReplicaSetUpdate(ctx context.Context, req *iapiserver.ReplicaSetRequest) (*iapiserver.ReplicaSetInfo, error)
	ReplicaSetGet(ctx context.Context, req *iapiserver.ReplicaSetGetRequest) (*iapiserver.ReplicaSetInfo, error)
	ReplicaSetList(ctx context.Context, req *iapiserver.ReplicaSetListRequest) (*iapiserver.ReplicaSetListResponse, error)

	HpaCreate(ctx context.Context, req *iapiserver.HpaRequest) (*iapiserver.HpaInfo, error)
	HpaDelete(ctx context.Context, req *iapiserver.HpaRequest) error
	HpaUpdate(ctx context.Context, req *iapiserver.HpaRequest) (*iapiserver.HpaInfo, error)
	HpaGet(ctx context.Context, req *iapiserver.HpaGetRequest) (*iapiserver.HpaInfo, error)
	HpaList(ctx context.Context, req *iapiserver.HpaListRequest) (*iapiserver.HpaListResponse, error)

	StatefulSetCreate(ctx context.Context, req *iapiserver.StatefulSetRequest) (*iapiserver.StatefulSetInfo, error)
	StatefulSetVersionUpdate(ctx context.Context, req *iapiserver.StatefulSetRequest) (*iapiserver.StatefulSetInfo, error)
	StatefulSetVersionList(ctx context.Context, req *iapiserver.StatefulSetVersionListRequest) (*iapiserver.StatefulSetVersionListResponse, error)
	StatefulSetRecreate(ctx context.Context, req *iapiserver.StatefulSetRequest) (*iapiserver.StatefulSetInfo, error)
	StatefulSetUpdate(ctx context.Context, req *iapiserver.StatefulSetRequest) (*iapiserver.StatefulSetInfo, error)
	StatefulSetDelete(ctx context.Context, req *iapiserver.StatefulSetRequest) error
	StatefulSetBatchDelete(ctx context.Context, req *iapiserver.StatefulSetBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.StatefulSetRequest]
	StatefulSetGet(ctx context.Context, req *iapiserver.StatefulSetGetRequest) (*iapiserver.StatefulSetInfo, error)
	StatefulSetList(ctx context.Context, req *iapiserver.StatefulSetListRequest) (*iapiserver.StatefulSetListResponse, error)

	DeploymentUpdate(ctx context.Context, req *iapiserver.DeploymentRequest) (*iapiserver.DeploymentResponse, error)
	DeploymentCreate(ctx context.Context, req *iapiserver.DeploymentRequest) (*iapiserver.DeploymentResponse, error)
	DeploymentDelete(ctx context.Context, req *iapiserver.DeploymentRequest) error
	DeploymentBatchDelete(ctx context.Context, req *iapiserver.DeploymentBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.DeploymentRequest]
	DeploymentRecreate(ctx context.Context, req *iapiserver.DeploymentRequest) (*iapiserver.DeploymentResponse, error)
	DeploymentGet(ctx context.Context, req *iapiserver.DeploymentGetRequest) (*iapiserver.DeploymentResponse, error)
	DeploymentList(ctx context.Context, req *iapiserver.DeploymentListRequest) (*iapiserver.DeploymentListResponse, error)
	DeploymentVersionUpdate(ctx context.Context, req *iapiserver.DeploymentRequest) error
	DeploymentVersionList(ctx context.Context, req *iapiserver.DeploymentVersoinListRequest) (*iapiserver.DeploymentVersionListResponse, error)

	JobCreate(ctx context.Context, req *iapiserver.JobRequest) (*iapiserver.JobInfo, error)
	JobDelete(ctx context.Context, req *iapiserver.JobRequest) error
	JobBatchDelete(ctx context.Context, req *iapiserver.JobBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.JobRequest]
	JobGet(ctx context.Context, req *iapiserver.JobGetRequest) (*iapiserver.JobResponse, error)
	JobList(ctx context.Context, req *iapiserver.JobListRequest) (*iapiserver.JobListResponse, error)

	CronJobCreate(ctx context.Context, req *iapiserver.CronJobRequest) (*iapiserver.CronJobInfo, error)
	CronJobUpdate(ctx context.Context, req *iapiserver.CronJobRequest) (*iapiserver.CronJobInfo, error)
	CronJobUpdateSuspend(ctx context.Context, req *iapiserver.CronJobRequest) (*iapiserver.CronJobInfo, error)
	CronJobUpdateSchedule(ctx context.Context, req *iapiserver.CronJobRequest) (*iapiserver.CronJobInfo, error)
	CronJobDelete(ctx context.Context, req *iapiserver.CronJobRequest) error
	CronJobGet(ctx context.Context, req *iapiserver.CronJobGetRequest) (*iapiserver.CronJobInfo, error)
	CronJobList(ctx context.Context, req *iapiserver.CronJobListRequest) (*iapiserver.CronJobListResponse, error)
	CronJobTrigger(ctx context.Context, req *iapiserver.CronJobRequest) (*iapiserver.JobInfo, error)

	ConfigMapCreate(ctx context.Context, req *iapiserver.ConfigMapRequest) (*iapiserver.ConfigMapInfo, error)
	ConfigMapUpdate(ctx context.Context, req *iapiserver.ConfigMapRequest) (*iapiserver.ConfigMapInfo, error)
	ConfigMapDelete(ctx context.Context, req *iapiserver.ConfigMapRequest) error
	ConfigMapBatchDelete(ctx context.Context, req *iapiserver.ConfigMapBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.ConfigMapRequest]
	ConfigMapGet(ctx context.Context, req *iapiserver.ConfigMapGetRequest) (*iapiserver.ConfigMapInfo, error)
	ConfigMapList(ctx context.Context, req *iapiserver.ConfigMapListRequest) (*iapiserver.ConfigMapListResponse, error)
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
		return nil, 0, err
	}
	wg := waitgroup.RunGenericConcurrently(ctx, clusters, concurrentFunc, timeouts...)
	eachClusterResult := wg.GetSuccessResultList()

	totalCount := CutPagingSliceResourceList(eachClusterResult, list, req.PageNum, req.PageSize, sortFunc)
	return eachClusterResult, totalCount, nil
}
