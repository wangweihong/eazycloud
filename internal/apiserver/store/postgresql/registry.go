package postgresql

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/apis/iregistry"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"gorm.io/gorm"
)

type registry struct {
	ds *datastore
}

func newRegistry(ds *datastore) *registry {
	return &registry{ds: ds}
}

func (s *registry) List(ctx context.Context, param *iapiserver.RegistryListRequest) ([]*iregistry.Registry, int64, error) {
	var meta []*iregistry.Registry
	var total int64

	resourceSpecificFilter := func(q *gorm.DB) *gorm.DB {
		// if param.FilterType != "" {
		// 	q = q.Where("state = ?", param.FilterState)
		// }
		if param.FilterState != "" {
			q = q.Where("state = ?", param.FilterState)
		}
		return q
	}

	err := param.ToQuery(ctx, s.ds.db, resourceSpecificFilter).Model(&iregistry.Registry{}).
		Find(&meta).Count(&total).Error

	return meta, total, err
}

func (s *registry) Get(ctx context.Context, id string) (*iregistry.Registry, error) {
	return nil, nil
}

func (s *registry) Add(ctx context.Context, data *iregistry.Registry) (*iregistry.Registry, error) {
	err := s.ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if CheckExists(tx, &iregistry.Registry{}, map[string]any{
			"name": data.Name,
		}) {
			return errors.Errorf("exists name with %v", data.Name)
		}

		if err := tx.Create(data).Error; err != nil {
			return errors.WithStack(err)
		}

		return nil
	})

	return data, err
}

func (s *registry) Delete(ctx context.Context, id string) error {
	return s.ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Where("id = ?", id).Delete(&iregistry.Registry{}).Error
	})
}

func (s *registry) Update(ctx context.Context, data *iregistry.Registry) error {
	return s.ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&iregistry.Registry{}).Where("id = ?", data.ID).Updates(data).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *registry) Sync(ctx context.Context, pages []*iregistry.Registry) error {
	if len(pages) == 0 {
		return nil
	}

	return nil
}
