package huaweicloud

import (
	"context"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2/model"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
	"github.com/wangweihong/gotoolbox/pkg/errors"

	ims "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2"
)

type Image struct {
	c   *ims.ImsClient
	err error
}

func NewImage(cfg *AccessConfig) (*Image, error) {
	c, err := cfg.HcImsClient(cfg.Region)
	return &Image{c: c, err: errors.WithStack(err)}, nil
}

// ImageList 查看镜像列表
func (p *Image) ImageList(
	ctx context.Context,
	req *model.ListImagesRequest,
) (*model.ListImagesResponse, error) {
	resp, err := libs.InvokeApi(ctx, p.err, req, p.c.ListImages)
	return resp, errors.WithStack(err)
}
