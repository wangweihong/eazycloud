package huaweicloud

import (
	"context"

	bss "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/bss/v2"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
	"github.com/wangweihong/gotoolbox/pkg/errors"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/bss/v2/model"
)

type Business struct {
	c   *bss.BssClient
	err error
}

func NewBusiness(cfg *AccessConfig) (*Business, error) {
	c, err := cfg.HcBssClient()
	return &Business{c: c, err: errors.WithStack(err)}, nil
}

// ShowCustomerAccountBalances 查询账户余额
func (p *Business) ShowCustomerAccountBalances(
	ctx context.Context,
	req *model.ShowCustomerAccountBalancesRequest,
) (*model.ShowCustomerAccountBalancesResponse, error) {
	resp, err := libs.InvokeApi(ctx, p.err, req, p.c.ShowCustomerAccountBalances)
	return resp, errors.WithStack(err)
}

// ShowCustomerAccountBalances 查询汇总账单
func (p *Business) ShowCustomerMonthlySum(
	ctx context.Context,
	req *model.ShowCustomerMonthlySumRequest,
) (*model.ShowCustomerMonthlySumResponse, error) {
	resp, err := libs.InvokeApi(ctx, p.err, req, p.c.ShowCustomerMonthlySum)
	return resp, errors.WithStack(err)
}

// ListCustomerAccountChangeRecords 查询收支明细
func (p *Business) ListCustomerAccountChangeRecords(
	ctx context.Context,
	req *model.ListCustomerAccountChangeRecordsRequest,
) (*model.ListCustomerAccountChangeRecordsResponse, error) {
	resp, err := libs.InvokeApi(ctx, p.err, req, p.c.ListCustomerAccountChangeRecords)
	return resp, errors.WithStack(err)
}

// ListResourceUsage 查询95计费资源用量明细
// 当前仅支持查询CDN、OBS、IEC和VPC四种云服务类型的资源用量明细，仅针对95计费场景。
func (p *Business) ListResourceUsage(
	ctx context.Context,
	req *model.ListResourceUsageRequest,
) (*model.ListResourceUsageResponse, error) {
	resp, err := libs.InvokeApi(ctx, p.err, req, p.c.ListResourceUsage)
	return resp, errors.WithStack(err)

}

// ListResourceUsage 查询资源包使用明细
func (p *Business) ListFreeResourcesUsageRecords(
	ctx context.Context,
	req *model.ListFreeResourcesUsageRecordsRequest,
) (*model.ListFreeResourcesUsageRecordsResponse, error) {
	resp, err := libs.InvokeApi(ctx, p.err, req, p.c.ListFreeResourcesUsageRecords)
	return resp, errors.WithStack(err)
}

// ListCustomerBillsFeeRecords 查询流水账单
func (p *Business) ListCustomerBillsFeeRecords(
	ctx context.Context,
	req *model.ListCustomerBillsFeeRecordsRequest,
) (*model.ListCustomerBillsFeeRecordsResponse, error) {
	resp, err := libs.InvokeApi(ctx, p.err, req, p.c.ListCustomerBillsFeeRecords)
	return resp, errors.WithStack(err)
}

// ListCustomerBillsFeeRecords 查询账单明细
func (p *Business) ListCustomerselfResourceRecordDetails(
	ctx context.Context,
	req *model.ListCustomerselfResourceRecordDetailsRequest,
) (*model.ListCustomerselfResourceRecordDetailsResponse, error) {
	resp, err := libs.InvokeApi(ctx, p.err, req, p.c.ListCustomerselfResourceRecordDetails)
	return resp, errors.WithStack(err)

}
