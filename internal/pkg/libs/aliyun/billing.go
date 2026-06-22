package aliyun

import (
	"context"

	"github.com/wangweihong/gotoolbox/pkg/errors"

	bssopenapi20171214 "github.com/alibabacloud-go/bssopenapi-20171214/v3/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type Billing struct {
	c *bssopenapi20171214.Client
}

func NewBilling(ak, sk string, config *openapi.Config) (*Billing, error) {
	if config == nil {
		config = &openapi.Config{
			AccessKeyId:     tea.String(ak),
			AccessKeySecret: tea.String(sk),
		}
	}

	config.Endpoint = tea.String("business.aliyuncs.com")
	c, err := bssopenapi20171214.NewClient(config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Billing{c: c}, nil
}

// QueryAccountBalance 查询用户余额
func (c *Billing) QueryAccountBalance(
	ctx context.Context,
) (*bssopenapi20171214.QueryAccountBalanceResponse, error) {
	ret, err := c.c.QueryAccountBalanceWithOptions(&util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// QueryAccountTransactionDetails 查询账户交易明细
func (c *Billing) QueryAccountTransactionDetails(
	ctx context.Context,
	req *bssopenapi20171214.QueryAccountTransactionDetailsRequest,
) (*bssopenapi20171214.QueryAccountTransactionDetailsResponse, error) {
	ret, err := c.c.QueryAccountTransactionDetailsWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// QueryAccountTransactions 查询用户账户流水信息
func (c *Billing) QueryAccountTransactions(
	ctx context.Context,
	req *bssopenapi20171214.QueryAccountTransactionsRequest,
) (*bssopenapi20171214.QueryAccountTransactionsResponse, error) {
	ret, err := c.c.QueryAccountTransactionsWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// QueryBillOverview 查询用户某个账期内结算账单概览
func (c *Billing) QueryBillOverview(
	ctx context.Context,
	req *bssopenapi20171214.QueryBillOverviewRequest,
) (*bssopenapi20171214.QueryBillOverviewResponse, error) {
	ret, err := c.c.QueryBillOverviewWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// QueryBill 查询用户某个账期内结算账单
func (c *Billing) QueryBill(
	ctx context.Context,
	req *bssopenapi20171214.QueryBillRequest,
) (*bssopenapi20171214.QueryBillResponse, error) {
	ret, err := c.c.QueryBillWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// QueryAccountBill 查询用户某个账期内的消费，并以资源所有者的维度进行汇总
func (c *Billing) QueryAccountBill(
	ctx context.Context,
	req *bssopenapi20171214.QueryAccountBillRequest,
) (*bssopenapi20171214.QueryAccountBillResponse, error) {
	ret, err := c.c.QueryAccountBillWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// QuerySettleBill 查询用户某个账期内结算账单。支持账单条目超过50000条的国内账号
func (c *Billing) QuerySettleBill(
	ctx context.Context,
	req *bssopenapi20171214.QuerySettleBillRequest,
) (*bssopenapi20171214.QuerySettleBillResponse, error) {
	ret, err := c.c.QuerySettleBillWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// DescribeInstanceBill 查询用户某个账期内所有商品实例或计费项的消费汇总
func (c *Billing) DescribeInstanceBill(
	ctx context.Context,
	req *bssopenapi20171214.DescribeInstanceBillRequest,
) (*bssopenapi20171214.DescribeInstanceBillResponse, error) {
	ret, err := c.c.DescribeInstanceBillWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}
