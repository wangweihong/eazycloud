package tecentcloud

import (
	"context"

	tcommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

type VpcClient struct {
	c   *vpc.Client
	err error
}

func (c *VpcClient) NewVpcClient(securityId, securityKey string, region string) *VpcClient {
	credential := tcommon.NewCredential(securityId, securityKey)
	clientProfile := profile.NewClientProfile()
	//clientProfile.HttpProfile.Endpoint = endpoint
	client, err := vpc.NewClient(credential, region, clientProfile)
	return &VpcClient{
		c:   client,
		err: err,
	}
}

func (c *VpcClient) DescribeVpcs(ctx context.Context, req *vpc.DescribeVpcsRequest) (*vpc.DescribeVpcsResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeVpcs)
	return resp, errors.WithStack(err)
}

func (c *VpcClient) DescribeSubnets(ctx context.Context, req *vpc.DescribeSubnetsRequest) (*vpc.DescribeSubnetsResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeSubnets)
	return resp, errors.WithStack(err)
}

func (c *VpcClient) DescribeRouteTables(ctx context.Context, req *vpc.DescribeRouteTablesRequest) (*vpc.DescribeRouteTablesResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeRouteTables)
	return resp, errors.WithStack(err)
}

func (c *VpcClient) DescribeAddresses(ctx context.Context, clusterUUID string, region string, req *vpc.DescribeAddressesRequest) (*vpc.DescribeAddressesResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeAddresses)
	return resp, errors.WithStack(err)
}
