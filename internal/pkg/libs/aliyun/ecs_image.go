package aliyun

import (
	"context"

	"github.com/wangweihong/gotoolbox/pkg/errors"

	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
)

// ImageList 查询镜像列表
func (c *Ecs) ImageList(
	ctx context.Context,
	req *ecs20140526.DescribeImagesRequest,
) (*ecs20140526.DescribeImagesResponse, error) {
	ret, err := c.DescribeImagesWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// ImageDelete 删除镜像
func (c *Ecs) ImageDelete(
	ctx context.Context,
	req *ecs20140526.DeleteImageRequest,
) (*ecs20140526.DeleteImageResponse, error) {
	ret, err := c.DeleteImageWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// ImageCreate 创建镜像
func (c *Ecs) ImageCreate(
	ctx context.Context,
	req *ecs20140526.CreateImageRequest,
) (*ecs20140526.CreateImageResponse, error) {
	ret, err := c.CreateImageWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}
