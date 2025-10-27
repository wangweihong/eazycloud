package cluster

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
)

func (s *clusterService) List(ctx context.Context, req *iapiserver.ClusterListRequest) (*iapiserver.ClusterListResponse, error) {
	clusters, err := s.store.Clusters().List(ctx)
	if err != nil {
		return nil, err
	}

	resp := &iapiserver.ClusterListResponse{}
	resp.List = clusters
	resp.Total = int64(len(resp.List))
	return resp, nil
}

func (s *clusterService) Get(ctx context.Context, req *iapiserver.ClusterGetRequest) (*iapiserver.ClusterGetResponse, error) {
	return nil, nil
}

func (s *clusterService) Delete(ctx context.Context, req *iapiserver.PodRequest) error {
	return nil
}

func (s *clusterService) Add(ctx context.Context, req *iapiserver.ClusterAddRequest) error {
	return nil
}

func (s *clusterService) Update(ctx context.Context, req *iapiserver.ClusterUpdateRequest) error {
	return nil
}

func (s *clusterService) Stop(ctx context.Context, req *iapiserver.ClusterStartRequest) error {
	return nil
}

func (s *clusterService) Start(ctx context.Context, req *iapiserver.ClusterStopResponse) error {
	return nil
}
