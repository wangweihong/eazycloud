package postgresql

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type setting struct {
	ds *datastore
}

func newSetting(ds *datastore) *setting {
	return &setting{ds}
}

func (s *setting) List(ctx context.Context) ([]*iapiserver.Setting, error) {
	var meta []*iapiserver.Setting
	err := s.ds.db.WithContext(ctx).Model(&iapiserver.Setting{}).
		Find(&meta).Error
	return meta, err
}

func (s *setting) Get(ctx context.Context, id string) (*iapiserver.Setting, error) {
	var meta iapiserver.Setting
	meta.ID = id
	err := s.ds.db.
		Model(&iapiserver.Setting{}).
		Find(&meta).Error
	return &meta, err
}

func (s *setting) GetMultiByNames(ctx context.Context, names ...string) ([]*iapiserver.Setting, error) {
	var meta []*iapiserver.Setting

	err := s.ds.db.WithContext(ctx).Model(&iapiserver.Setting{}).
		Where("name IN ?", names).
		First(&meta).Error
	return meta, err
}

func (s *setting) GetByName(ctx context.Context, name string) (*iapiserver.Setting, error) {
	var meta *iapiserver.Setting

	err := s.ds.db.WithContext(ctx).Model(&iapiserver.Setting{}).
		Where("name = ?", name).
		First(&meta).Error
	return meta, err
}

func (s *setting) Add(ctx context.Context, data *iapiserver.Setting) (*iapiserver.Setting, error) {
	err := s.ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if CheckExists(tx, &iapiserver.Setting{}, map[string]any{
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

func (s *setting) Delete(ctx context.Context, id string) error {
	return s.ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Where("id = ?", id).Delete(&iapiserver.Setting{}).Error
	})
}

// 创建或更新
func (s *setting) Upsert(ctx context.Context, data *iapiserver.Setting) error {
	return s.ds.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}}, // 冲突检测字段
		UpdateAll: true,                            // 更新全部字段
		Where: clause.Where{Exprs: []clause.Expression{
			clause.Lt{Column: "resource_version", Value: data.ResourceVersion}, // 仅当新版本更高时更新
		}},
	}).Create(data).Error
}
