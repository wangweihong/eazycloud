package aliyun

import (
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

func (c *Ecs) ListRegion(
	clusterUUID string,
	req *ecs20140526.DescribeRegionsRequest,
	opts ...*util.RuntimeOptions,
) (*ecs20140526.DescribeRegionsResponse, error) {
	ret, err := c.DescribeRegionsWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) ListZone(
	clusterUUID string,
	req *ecs20140526.DescribeZonesRequest,
	opts ...*util.RuntimeOptions,
) (*ecs20140526.DescribeZonesResponse, error) {
	ret, err := c.DescribeZonesWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) ListAvailableResource(
	clusterUUID string,
	req *ecs20140526.DescribeAvailableResourceRequest,
) (*ecs20140526.DescribeAvailableResourceResponse, error) {
	ret, err := c.DescribeAvailableResourceWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}
