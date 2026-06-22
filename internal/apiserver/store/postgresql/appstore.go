package postgresql

import (
	"context"

	"gorm.io/gorm"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

type appStore struct {
	ds *datastore
}

func newAppStore(ds *datastore) *appStore {
	return &appStore{ds}
}

func (s *appStore) List(ctx context.Context) ([]*iapiserver.AppStore, error) {
	var meta []*iapiserver.AppStore
	err := s.ds.db.WithContext(ctx).Model(&iapiserver.AppStore{}).
		Preload("ApplicationTemplates").
		Preload("ApplicationCategories").
		Find(&meta).Error
	return meta, err
}

func (s *appStore) Query(ctx context.Context, param *iapiserver.AppStoreListRequest) ([]*iapiserver.AppStore, int64, error) {
	var meta []*iapiserver.AppStore
	var total int64

	resourceSpecificFilter := func(q *gorm.DB) *gorm.DB {
		if param.AppSource != nil {
			q = q.Where("app_source = ?", param.AppSource)
		}

		if param.FilterState != "" {
			q = q.Where("state = ?", param.FilterState)
		}
		return q
	}

	err := param.ToQuery(ctx, s.ds.db, resourceSpecificFilter).
		Model(&iapiserver.AppStore{}).
		Preload("ApplicationTemplates").
		Preload("ApplicationCategories").
		Find(&meta).Count(&total).Error

	return meta, total, err
}

func (s *appStore) Get(ctx context.Context, id string) (*iapiserver.AppStore, error) {
	var meta iapiserver.AppStore
	meta.ID = id
	err := s.ds.db.Model(&iapiserver.AppStore{}).
		Preload("ApplicationTemplates").
		Preload("ApplicationCategories").
		Find(&meta).Error
	return &meta, err
}

func (s *appStore) GetByName(ctx context.Context, name string) (*iapiserver.AppStore, error) {
	var meta *iapiserver.AppStore

	err := s.ds.db.WithContext(ctx).Model(&iapiserver.AppStore{}).
		Preload("ApplicationTemplates").
		Preload("ApplicationCategories").
		Where("name = ?", name).
		First(&meta).Error
	return meta, err
}

func (s *appStore) Add(ctx context.Context, data *iapiserver.AppStore) (*iapiserver.AppStore, error) {
	err := s.ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if CheckExists(tx, &iapiserver.AppStore{}, map[string]any{
			"name": data.Name,
		}) {
			return errors.Errorf("appstore exists name with %v", data.Name)
		}

		categories := data.ApplicationCategories
		data.ApplicationCategories = nil
		if err := tx.Create(data).Error; err != nil {
			return errors.WithStack(err)
		}

		for _, ver := range categories {
			ver.AppStoreID = data.ID
			if err := tx.Create(&ver).Error; err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})

	return data, err
}

func (s *appStore) Delete(ctx context.Context, id string) error {
	return s.ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("app_store_id = ?", id).Delete(&iapiserver.ApplicationInstance{}).Error; err != nil {
			return err
		}

		var templateIDs []string
		if err := tx.Model(&iapiserver.ApplicationTemplate{}).Where("app_store_id = ?", id).Pluck("id", &templateIDs).Error; err != nil {
			return err
		}

		if len(templateIDs) > 0 {
			if err := tx.Where("application_template_id IN (?)", templateIDs).Delete(&iapiserver.ApplicationTemplateCategory{}).Error; err != nil {
				return err
			}
			if err := tx.Where("id IN (?)", templateIDs).Delete(&iapiserver.ApplicationTemplate{}).Error; err != nil {
				return err
			}
		}

		if err := tx.Where("app_store_id = ?", id).Delete(&iapiserver.ApplicationCategory{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id).Delete(&iapiserver.AppStore{}).Error
	})
}

func (s *appStore) Update(ctx context.Context, data *iapiserver.AppStore) error {
	return s.ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&iapiserver.AppStore{}).Where("id = ?", data.ID).Updates(data).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *appStore) Sync(ctx context.Context, pages []*iapiserver.AppStore) error {
	if len(pages) == 0 {
		return nil
	}

	return nil
}
