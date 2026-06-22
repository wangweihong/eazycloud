package postgresql

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"gorm.io/gorm"
)

type applicationTemplate struct {
	ds *datastore
}

func newApplicationTemplate(ds *datastore) *applicationTemplate {
	return &applicationTemplate{ds}
}

func (s *applicationTemplate) List(ctx context.Context, param *iapiserver.ApplicationTemplateListRequest) ([]*iapiserver.ApplicationTemplate, int64, error) {
	var meta []*iapiserver.ApplicationTemplate
	var total int64

	resourceSpecificFilter := func(q *gorm.DB) *gorm.DB {
		if param.FilterStore != "" {
			q = q.Where("app_store_id = ?", param.FilterStore)
		}

		if param.FilterState != "" {
			q = q.Where("state = ?", param.FilterState)
		}

		if len(param.CategoriesShadow) != 0 {
			if param.AllCategories {
				// AND条件：必须包含所有分类
				subQuery := q.
					Select("template_id").
					Table("application_template_categories").
					Where("category_id IN ?", param.CategoriesShadow).
					Group("template_id").
					Having("COUNT(DISTINCT category_id) = ?", len(param.CategoriesShadow))

				q = q.Where("id IN (?)", subQuery)
			} else {
				// OR条件：包含任一分类
				subQuery := q.
					Select("DISTINCT template_id").
					Table("application_template_categories").
					Where("category_id IN ?", len(param.CategoriesShadow))

				q = q.Where("id IN (?)", subQuery)
			}
		}
		return q
	}

	err := param.ToQuery(ctx, s.ds.db, resourceSpecificFilter).
		Model(&iapiserver.ApplicationTemplate{}).
		Preload("Versions", func(db *gorm.DB) *gorm.DB {
			// 只加载需要的字段
			return db.Select("id", "name", "state")
		}).
		Preload("Instances").
		Preload("Categories").
		Find(&meta).Count(&total).Error

	return meta, total, err
}

func (s *applicationTemplate) Get(ctx context.Context, id string) (*iapiserver.ApplicationTemplate, error) {
	var meta iapiserver.ApplicationTemplate
	meta.ID = id
	err := s.ds.db.Model(&iapiserver.ApplicationTemplate{}).
		Preload("Categories").
		Preload("Versions").
		Preload("Instances").
		Find(&meta).Error
	return &meta, err
}

func (s *applicationTemplate) GetByName(ctx context.Context, store string, name string) (*iapiserver.ApplicationTemplate, error) {
	var meta iapiserver.ApplicationTemplate

	err := s.ds.db.Model(&iapiserver.ApplicationTemplate{}).
		Preload("Categories").
		Preload("Versions").
		Preload("Instances").
		Where("app_store_id = ?", store).
		Where("name = ?", name).
		First(&meta).Error
	return &meta, err
}

func (s *applicationTemplate) Add(ctx context.Context, data *iapiserver.ApplicationTemplate) (*iapiserver.ApplicationTemplate, error) {
	err := s.ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var store iapiserver.AppStore
		store.ID = data.AppStoreID
		if err := tx.Model(&iapiserver.AppStore{}).Find(&store).Error; err != nil {
			return errors.WithStack(err)
		}

		switch store.AppSource {
		case iapiserver.AppSourceExternal:
			return errors.Errorf("external app store doesn't support add template, please add template in app store")
		}

		if CheckExists(tx, &iapiserver.ApplicationTemplate{}, map[string]any{
			"name": data.Name,
		}) {
			return errors.Errorf("exists name with %v", data.Name)
		}

		versions := data.Versions
		data.Versions = nil
		if err := tx.Omit("Categories").Create(data).Error; err != nil {
			return errors.WithStack(err)
		}
		if data.Categories != nil {
			if err := tx.Model(data).Association("Categories").Append(data.Categories); err != nil {
				return errors.WithStack(err)
			}
		}

		for _, ver := range versions {
			ver.ApplicationTemplateID = data.ID
			if err := tx.Create(&ver).Error; err != nil {
				return errors.WithStack(err)
			}
		}

		return nil
	})

	return data, err
}

func (s *applicationTemplate) Delete(ctx context.Context, id string) error {
	return s.ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("application_template_id = ?", id).Delete(&iapiserver.ApplicationTemplateVersion{}).Error; err != nil {
			return err
		}
		//FIXME:实例需要清理k8s的数据
		if err := tx.Where("application_template_id = ?", id).Delete(&iapiserver.ApplicationInstance{}).Error; err != nil {
			return err
		}

		if err := tx.Where("application_template_id = ?", id).Delete(&iapiserver.ApplicationTemplateCategory{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&iapiserver.ApplicationTemplate{}, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *applicationTemplate) Update(ctx context.Context, data *iapiserver.ApplicationTemplate) error {
	return nil
}

func (s *applicationTemplate) Sync(ctx context.Context, pages []*iapiserver.ApplicationTemplate) error {
	if len(pages) == 0 {
		return nil
	}

	return nil
}
