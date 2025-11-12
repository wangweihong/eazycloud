package store

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/ikubeagent"
)

type InstallStateStore interface {
	GetByName(ctx context.Context, name string) (*ikubeagent.InstallState, error)
	Upsert(ctx context.Context, data *ikubeagent.InstallState) (*ikubeagent.InstallState, error)
}
