package ctyun

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/ictyun"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"
)

type Billing struct {
	c           *Client
	serviceType string
}


// BillCycleFeeList 查询包周期流水账单
func (p *Billing) BillCycleFeeList(ctx context.Context, req *ictyun.BillCycleFeeListRequest, resp *ictyun.BillCycleFeeListResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/queryBillCycleFee").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// BillingList 查询按需流水账单
func (p *Billing) BillOnDemandFeeList(ctx context.Context, req *ictyun.BillOnDemandFeeListRequest, resp *ictyun.BillOnDemandFeeListResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/queryBillOnDemandFee").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// BillCycleDetailProdCycleId  账单详情 - 统计维度(产品) - 统计周期（按账期) - 计费模式 (包周期）
func (p *Billing) BillCycleDetailProdCycleId(
	ctx context.Context,
	req *ictyun.BillCycleBillDetailProdCycleIdListRequest,
	resp *ictyun.BillCycleBillDetailProdCycleIdListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/qryCycleBillDetail_Prod_CycleId").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// BillOnDemandBillDetailResDetail  账单详情 - 统计维度(产品) - 统计周期（按账期) - 计费模式 (包周期）
func (p *Billing) BillOnDemandBillDetailResDetail(
	ctx context.Context,
	req *ictyun.BillOnDemandBillDetailResDetailListRequest,
	resp *ictyun.BillOnDemandBillDetailResDetailListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/bill_qryOnDemandBillDetail_Res_Detail").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// BillOnDemandBillDetailProductDetail 账单详情 - 统计维度(产品) - 统计周期（账期) - 计费模式 (按需）
func (p *Billing) BillOnDemandBillDetailProductDetail(
	ctx context.Context,
	req *ictyun.BillOnDemandBillDetailProductDetailListRequest,
	resp *ictyun.BillOnDemandBillDetailProductDetailListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/qryOnDemandBillDetail_Prod_CycleId").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// BillOnDemandBillDetailResCycleId  账单详情 - 统计维度(资源) - 统计周期（账期) - 计费模式 (按需）
func (p *Billing) BillOnDemandBillDetailResCycleId(
	ctx context.Context,
	req *ictyun.BillOnDemandBillDetailResCycleIdListRequest,
	resp *ictyun.BillOnDemandBillDetailResCycleIdListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/qryOnDemandBillDetail_Res_CycleId").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// BillOnDemandBillDetailResCycleId  账单详情 - 统计维度(使用量) - 统计周期（账期) - 计费模式 (按需）
func (p *Billing) BillOnDemandBillDetailUsageCycleId(
	ctx context.Context,
	req *ictyun.BillOnDemandBillDetailUsageCycleIdListRequest,
	resp *ictyun.BillOnDemandBillDetailUsageCycleIdListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/qryOnDemandBillDetail_Usage_CycleId").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// BillOnDemandBillDetailResCycleId 统计维度(使用量) - 统计周期（账期) - 计费模式 (按需）
func (p *Billing) BillOnDemandBillDetailUsageDetail(
	ctx context.Context,
	req *ictyun.BillOnDemandBillDetailUsageDetailListRequest,
	resp *ictyun.BillOnDemandBillDetailUsageDetailListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/qryOnDemandBillDetail_Usage_Detail").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}
