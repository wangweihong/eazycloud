package tecentcloud

import (
	"context"

	billing "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/billing/v20180709"
	tcommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

type BillingClient struct {
	c   *billing.Client
	err error
}

func (c *BillingClient) NewBillingClient(securityId, securityKey string, region string) *BillingClient {
	credential := tcommon.NewCredential(securityId, securityKey)
	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = HttpProfileUserEndpoint
	client, err := billing.NewClient(credential, region, clientProfile)
	return &BillingClient{
		c:   client,
		err: err,
	}
}

func (c *BillingClient) DescribeBillSummaryByProject(ctx context.Context, req *billing.DescribeBillSummaryByProjectRequest) (*billing.DescribeBillSummaryByProjectResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeBillSummaryByProject)
	return resp, errors.WithStack(err)
}

func (c *BillingClient) DescribeBillList(ctx context.Context, req *billing.DescribeBillListRequest) (*billing.DescribeBillListResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeBillList)
	return resp, errors.WithStack(err)
}

func (c *BillingClient) DescribeAccountBalance(ctx context.Context, req *billing.DescribeAccountBalanceRequest) (*billing.DescribeAccountBalanceResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeAccountBalance)
	return resp, errors.WithStack(err)
}

func (c *BillingClient) DescribeCostDetail(ctx context.Context, req *billing.DescribeCostDetailRequest) (*billing.DescribeCostDetailResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeCostDetail)
	return resp, errors.WithStack(err)
}

func (c *BillingClient) DescribeDealsByCond(ctx context.Context, req *billing.DescribeDealsByCondRequest) (*billing.DescribeDealsByCondResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeDealsByCond)
	return resp, errors.WithStack(err)
}

func (c *BillingClient) PayDeals(ctx context.Context, req *billing.PayDealsRequest) (*billing.PayDealsResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.PayDeals)
	return resp, errors.WithStack(err)
}

func (c *BillingClient) DescribeCostSummaryByProduct(ctx context.Context, req *billing.DescribeCostSummaryByProductRequest) (*billing.DescribeCostSummaryByProductResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeCostSummaryByProduct)
	return resp, errors.WithStack(err)
}

func (c *BillingClient) DescribeCostSummaryByRegion(ctx context.Context, req *billing.DescribeCostSummaryByRegionRequest) (*billing.DescribeCostSummaryByRegionResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeCostSummaryByRegion)
	return resp, errors.WithStack(err)
}

func (c *BillingClient) DescribeBillDetail(ctx context.Context, req *billing.DescribeBillDetailRequest) (*billing.DescribeBillDetailResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeBillDetail)
	return resp, errors.WithStack(err)
}

func (c *BillingClient) DescribeBillResourceSummary(ctx context.Context, req *billing.DescribeBillResourceSummaryRequest) (*billing.DescribeBillResourceSummaryResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeBillResourceSummary)
	return resp, errors.WithStack(err)
}

func (c *BillingClient) DescribeBillSummaryByProduct(ctx context.Context, req *billing.DescribeBillSummaryByProductRequest) (*billing.DescribeBillSummaryByProductResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeBillSummaryByProduct)
	return resp, errors.WithStack(err)
}

func (c *BillingClient) DescribeBillSummaryByPayMode(ctx context.Context, req *billing.DescribeBillSummaryByPayModeRequest) (*billing.DescribeBillSummaryByPayModeResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeBillSummaryByPayMode)
	return resp, errors.WithStack(err)
}

func (c *BillingClient) DescribeBillSummaryByRegion(ctx context.Context, req *billing.DescribeBillSummaryByRegionRequest) (*billing.DescribeBillSummaryByRegionResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeBillSummaryByRegion)
	return resp, errors.WithStack(err)
}
