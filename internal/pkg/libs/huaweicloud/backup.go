package huaweicloud

import (
	"context"

	cbr "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cbr/v1"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cbr/v1/model"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

type Backup struct {
	c   *cbr.CbrClient
	err error
}

func NewBackup(cfg *AccessConfig) *Backup {
	c, err := cfg.HcCbrClient(cfg.Region)
	return &Backup{c: c, err: err}
}

// BackupList 查看备份列表
func (p *Backup) BackupList(
	ctx context.Context,
	req *model.ListBackupsRequest,
) (*model.ListBackupsResponse, error) {
	resp, err := libs.InvokeApi(ctx, p.err, req, p.c.ListBackups)
	return resp, errors.WithStack(err)
}
