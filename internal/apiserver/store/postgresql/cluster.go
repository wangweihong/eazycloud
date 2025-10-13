package postgresql

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
)

type cluster struct {
	ds *datastore
}

func newCluster(ds *datastore) *cluster {
	return &cluster{ds}
}

func (s *cluster) List(ctx context.Context) ([]*iapiserver.Cluster, error) {
	return nil, nil
}

func (s *cluster) Get(ctx context.Context, id string) (*iapiserver.Cluster, error) {
	return nil, nil
}

func (s *cluster) Add(ctx context.Context, data *iapiserver.Cluster) (*iapiserver.Cluster, error) {
	return data, nil
}

func (s *cluster) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *cluster) Update(ctx context.Context, data *iapiserver.Cluster) error {
	return nil
}

func (s *cluster) Sync(ctx context.Context, pages []*iapiserver.Cluster) error {
	if len(pages) == 0 {
		return nil
	}

	return nil
}
