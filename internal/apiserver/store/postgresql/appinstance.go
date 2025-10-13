package postgresql

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"gorm.io/gorm"
)

type applicationInstance struct {
	ds *datastore
}

func newApplicationInstance(ds *datastore) *applicationInstance {
	return &applicationInstance{ds}
}

func (s *applicationInstance) List(ctx context.Context, param *iapiserver.ApplicationInstanceListRequest) ([]*iapiserver.ApplicationInstance, int64, error) {
	var meta []*iapiserver.ApplicationInstance
	var total int64

	resourceSpecificFilter := func(q *gorm.DB) *gorm.DB {
		if param.FilterState != "" {
			q = q.Where("state = ?", param.FilterState)
		}

		if param.FilterCluster != "" {
			q = q.Where("cluster_id = ?", param.FilterCluster)
		}
		return q
	}

	err := param.ToQuery(ctx, s.ds.db, resourceSpecificFilter).
		Find(&meta).Count(&total).Error

	return meta, total, err
}

func (s *applicationInstance) Get(ctx context.Context, id string) (*iapiserver.ApplicationInstance, error) {
	var meta iapiserver.ApplicationInstance
	meta.ID = id
	err := s.ds.db.
		Model(&iapiserver.ApplicationInstance{}).
		Find(&meta).Error
	return &meta, err
}

func (s *applicationInstance) GetByName(ctx context.Context, name string) (*iapiserver.ApplicationInstance, error) {
	return nil, nil
}

func (s *applicationInstance) Add(ctx context.Context, data *iapiserver.ApplicationInstance) (*iapiserver.ApplicationInstance, error) {
	return data, nil
}

func (s *applicationInstance) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *applicationInstance) Update(ctx context.Context, data *iapiserver.ApplicationInstance) error {
	return nil
}

func (s *applicationInstance) Sync(ctx context.Context, pages []*iapiserver.ApplicationInstance) error {
	if len(pages) == 0 {
		return nil
	}

	return nil
}
