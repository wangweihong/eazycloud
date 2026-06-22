package postgresql

import (
	"context"
	gerrors "errors"

	"github.com/wangweihong/eazycloud/apis/ikubeagent"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type installstate struct {
	ds *datastore
}

func newInstallState(ds *datastore) *installstate {
	return &installstate{ds}
}

func (s *installstate) GetByName(ctx context.Context, name string) (*ikubeagent.InstallState, error) {
	var meta ikubeagent.InstallState
	err := s.ds.db.WithContext(ctx).
		Model(&ikubeagent.InstallState{}).
		Where("name = ?", name).
		Find(&meta).Error
	return &meta, err
}

func (s *installstate) Upsert(ctx context.Context, data *ikubeagent.InstallState) (*ikubeagent.InstallState, error) {
	result := ikubeagent.InstallState{}
	// 1. 尝试锁定并查询现有记录（使用悲观锁）
	err := s.ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("name = ?", data.Name).
			First(&result).Error
		// err := tx.
		// 	Where("name = ?", data.Name).
		// 	First(&result).Error
		if gerrors.Is(err, gorm.ErrRecordNotFound) {
			result = *data
			return errors.WithStack(tx.Create(&result).Error) // 直接插入新记录
		}

		if err != nil {
			return errors.WithStack(err) // 其他查询错误
		}

		result.Extend = data.Extend
		result.ExtendShadow = data.ExtendShadow
		result.Description = data.Description
		result.State = data.State
		result.ControlPlane = data.ControlPlane
		result.HighAvailable = data.HighAvailable
		result.EndTime = data.EndTime
		result.ErrorMessage = data.ErrorMessage
		return tx.Model(&result).
			Where("id = ?", result.ID).
			Updates(&result).Error
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}
