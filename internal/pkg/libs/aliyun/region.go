package aliyun

import (
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/typeutil"

	cloudapi20160714 "github.com/alibabacloud-go/cloudapi-20160714/v4/client"

	util "github.com/alibabacloud-go/tea-utils/v2/service"
)

type Region struct {
	*cloudapi20160714.Client
}

func (c *Region) NewRegion(ak, sk string, region *string) (*Region, error) {
	if region == nil || *region == "" {
		tmp := defaultRegion
		region = &tmp
	}

	config := &openapi.Config{
		AccessKeyId:     typeutil.String(ak),
		AccessKeySecret: typeutil.String(sk),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/CloudAPI
	config.Endpoint = tea.String(fmt.Sprintf("apigateway.%v.aliyuncs.com", *region))
	rc, err := cloudapi20160714.NewClient(config)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &Region{Client: rc}, nil
}

func (c *Region) ListRegion(
	clusterUUID string,
	req *cloudapi20160714.DescribeRegionsRequest,
) (*cloudapi20160714.DescribeRegionsResponse, error) {
	ret, err := c.DescribeRegionsWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Region) ListZone(
	clusterUUID string,
	req *cloudapi20160714.DescribeZonesRequest,
) (*cloudapi20160714.DescribeZonesResponse, error) {
	ret, err := c.DescribeZonesWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}
