package postgresql

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"gorm.io/gorm"
)

type serviceProvider struct {
	ds *datastore
}

func newServiceProvider(ds *datastore) *serviceProvider {
	return &serviceProvider{ds}
}

func (s *serviceProvider) List(ctx context.Context, param *iapiserver.ServiceProviderListRequest) ([]*iapiserver.ServiceProvider, int64, error) {
	var meta []*iapiserver.ServiceProvider
	var total int64

	resourceSpecificFilter := func(q *gorm.DB) *gorm.DB {
		if param.Protocol != "" {
			q = q.Where("protocol = ?", param.Protocol)
		}
		return q
	}

	err := param.ToQuery(ctx, s.ds.db, resourceSpecificFilter).
		Find(&meta).Count(&total).Error

	return meta, total, err
}

func (s *serviceProvider) Get(ctx context.Context, id string) (*iapiserver.ServiceProvider, error) {
	var meta iapiserver.ServiceProvider
	meta.ID = id
	err := s.ds.db.
		Model(&iapiserver.ServiceProvider{}).
		Find(&meta).Error
	return &meta, err
}

func (s *serviceProvider) GetByName(ctx context.Context, name string) (*iapiserver.ServiceProvider, error) {
	return nil, nil
}

func (s *serviceProvider) Add(ctx context.Context, data *iapiserver.ServiceProvider) (*iapiserver.ServiceProvider, error) {
	return data, nil
}

func (s *serviceProvider) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *serviceProvider) Update(ctx context.Context, data *iapiserver.ServiceProvider) error {
	return nil
}

func (s *serviceProvider) Sync(ctx context.Context, pages []*iapiserver.ServiceProvider) error {
	if len(pages) == 0 {
		return nil
	}

	return nil
}
