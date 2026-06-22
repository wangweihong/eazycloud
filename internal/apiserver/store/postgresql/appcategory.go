package postgresql

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"gorm.io/gorm"
)

type applicationCategory struct {
	ds *datastore
}

func newApplicationCategory(ds *datastore) *applicationCategory {
	return &applicationCategory{ds}
}

func (s *applicationCategory) List(ctx context.Context, param *iapiserver.ApplicationCategoryListRequest) ([]*iapiserver.ApplicationCategory, int64, error) {
	var meta []*iapiserver.ApplicationCategory
	var total int64
	resourceSpecificFilter := func(q *gorm.DB) *gorm.DB {
		if param.AppStoreID != "" {
			q = q.Where("app_store_id = ?", param.AppStoreID)
		}
		return q
	}
	err := param.ToQuery(ctx, s.ds.db, resourceSpecificFilter).
		Model(&iapiserver.ApplicationCategory{}).
		Preload("Templates", func(db *gorm.DB) *gorm.DB {
			// 只加载需要的字段
			return db.Select("id", "name")
		}).
		Find(&meta).Count(&total).Error

	return meta, total, err
}

func (s *applicationCategory) Get(ctx context.Context, id string) (*iapiserver.ApplicationCategory, error) {
	var meta iapiserver.ApplicationCategory
	meta.ID = id
	err := s.ds.db.WithContext(ctx).
		Model(&iapiserver.ApplicationCategory{}).
		Find(&meta).Error
	return &meta, err
}

func (s *applicationCategory) GetByName(ctx context.Context, store, name string) (*iapiserver.ApplicationCategory, error) {
	var meta iapiserver.ApplicationCategory
	err := s.ds.db.WithContext(ctx).Model(&iapiserver.ApplicationCategory{}).
		Where("app_store_id = ?", store).
		Where("name = ?", name).
		First(&meta).Error
	return &meta, err
}

func (s *applicationCategory) Add(ctx context.Context, data *iapiserver.ApplicationCategory) (*iapiserver.ApplicationCategory, error) {
	err := s.ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var store iapiserver.AppStore
		store.ID = data.AppStoreID
		if err := tx.Model(&store).Find(&store).Error; err != nil {
			return errors.WithStack(err)
		}

		switch store.AppSource {
		case iapiserver.AppSourceExternal:
			return errors.Errorf("external app store doesn't support add catetory, please add template in app store")
		}

		if CheckExists(tx, &iapiserver.ApplicationCategory{}, map[string]any{
			"name": data.Name,
		}) {
			return errors.Errorf("exists name with %v", data.Name)
		}

		// 创建分类（忽略关联）
		if err := tx.Omit("Templates").Create(data).Error; err != nil {
			return err
		}

		if data.Templates != nil {
			if err := tx.Model(data).Association("Templates").Append(data.Templates); err != nil {
				return errors.WithStack(err)
			}
		}

		return nil
	})

	return data, err
}

func (s *applicationCategory) Delete(ctx context.Context, id string) error {
	return s.ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("application_template_id = ?", id).
			Delete(&iapiserver.ApplicationCategory{}).Error; err != nil {
			return err
		}

		if err := tx.Where("application_template_id = ?", id).Delete(&iapiserver.ApplicationTemplateCategory{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&iapiserver.ApplicationCategory{}, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *applicationCategory) Update(ctx context.Context, data *iapiserver.ApplicationCategory) error {
	return s.ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&iapiserver.ApplicationCategory{}).Where("id = ?", data.ID).Updates(data).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *applicationCategory) AddCategoryTemplates(ctx context.Context, categoryID string, templates []*iapiserver.ApplicationTemplate) error {
	return s.ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var category iapiserver.ApplicationCategory
		if err := tx.First(&category, "id = ?", categoryID).Error; err != nil {
			return err
		}

		var templateIDs []string
		var newTemplates []iapiserver.ApplicationTemplate
		for _, template := range templates {
			templateIDs = append(templateIDs, template.ID)
		}
		if err := tx.Where("id IN ?", templateIDs).Find(&newTemplates).Error; err != nil {
			return err
		}

		if err := tx.Model(&category).Association("Templates").Append(newTemplates); err != nil {
			return err
		}

		return nil
	})
}

func (s *applicationCategory) DeleteCategoryTemplates(ctx context.Context, categoryID string, templates []*iapiserver.ApplicationTemplate) error {
	return s.ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var category iapiserver.ApplicationCategory
		if err := tx.First(&category, "id = ?", categoryID).Error; err != nil {
			return err
		}

		var templateIDs []string
		var newTemplates []iapiserver.ApplicationTemplate
		for _, template := range templates {
			templateIDs = append(templateIDs, template.ID)
		}
		if err := tx.Where("id IN ?", templateIDs).Find(&newTemplates).Error; err != nil {
			return err
		}

		if err := tx.Model(&category).Association("Templates").Delete(newTemplates); err != nil {
			return err
		}

		return nil
	})
}

func (s *applicationCategory) ReplaceCategoryTemplates(ctx context.Context, categoryID string, templates []*iapiserver.ApplicationTemplate) error {
	return s.ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var category iapiserver.ApplicationCategory
		if err := tx.First(&category, "id = ?", categoryID).Error; err != nil {
			return err
		}

		var templateIDs []string
		var newTemplates []iapiserver.ApplicationTemplate
		for _, template := range templates {
			templateIDs = append(templateIDs, template.ID)
		}
		if err := tx.Where("id IN ?", templateIDs).Find(&newTemplates).Error; err != nil {
			return err
		}

		if err := tx.Model(&category).Association("Templates").Replace(newTemplates); err != nil {
			return err
		}

		return nil
	})
}

func (s *applicationCategory) Sync(ctx context.Context, pages []*iapiserver.ApplicationCategory) error {
	if len(pages) == 0 {
		return nil
	}

	return nil
}
