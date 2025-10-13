package cluster

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/internal/apiserver/store"
)

type ClusterSrv interface {
	List(ctx context.Context, req *iapiserver.ClusterListRequest) (*iapiserver.ClusterListResponse, error)
	Get(ctx context.Context, req *iapiserver.ClusterGetRequest) (*iapiserver.ClusterGetResponse, error)
	Delete(ctx context.Context, req *iapiserver.PodRequest) error
	Add(ctx context.Context, req *iapiserver.ClusterAddRequest) error
	Update(ctx context.Context, req *iapiserver.ClusterUpdateRequest) error
	Stop(ctx context.Context, req *iapiserver.ClusterStartRequest) error
	Start(ctx context.Context, req *iapiserver.ClusterStopResponse) error
}

type clusterService struct {
	store store.Factory
}

func NewService(str store.Factory) *clusterService {
	return &clusterService{store: str}
}
