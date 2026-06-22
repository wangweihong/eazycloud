package postgresql

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
)

type kubernetes struct {
	ds *datastore
}

func newKubernetes(ds *datastore) *kubernetes {
	return &kubernetes{ds}
}

func (s *kubernetes) List(ctx context.Context) ([]*iapiserver.Cluster, error) {
	return nil, nil
}

func (s *kubernetes) Get(ctx context.Context, id string) (*iapiserver.Cluster, error) {
	return nil, nil
}

func (s *kubernetes) Add(ctx context.Context, data *iapiserver.Cluster) (*iapiserver.Cluster, error) {
	return data, nil
}

func (s *kubernetes) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *kubernetes) Update(ctx context.Context, data *iapiserver.Cluster) error {
	return nil
}

func (s *kubernetes) Sync(ctx context.Context, pages []*iapiserver.Cluster) error {
	if len(pages) == 0 {
		return nil
	}

	return nil
}
