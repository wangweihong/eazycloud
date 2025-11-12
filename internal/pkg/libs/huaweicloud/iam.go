package huaweicloud

import (
	"context"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
	"github.com/wangweihong/gotoolbox/pkg/errors"

	iam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
)

type Iam struct {
	c   *iam.IamClient
	err error
}

func NewIam(cfg *AccessConfig) *Iam {
	c, err := cfg.HcIamClient()
	return &Iam{c: c, err: errors.WithStack(err)}

}

// RegionList 查看全部区域
func (p *Iam) RegionList(
	ctx context.Context,
	req *model.KeystoneListRegionsRequest,
) (*model.KeystoneListRegionsResponse, error) {
	resp, err := libs.InvokeApi(ctx, p.err, req, p.c.KeystoneListRegions)
	return resp, errors.WithStack(err)
}

// ProjectList 查看全部区域
func (p *Iam) KeystoneListProjects(
	ctx context.Context,
	req *model.KeystoneListProjectsRequest,
) (*model.KeystoneListProjectsResponse, error) {
	resp, err := libs.InvokeApi(ctx, p.err, req, p.c.KeystoneListProjects)
	return resp, errors.WithStack(err)
}
