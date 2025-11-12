package huaweicloud_test

import (
	"context"
	"testing"

	"github.com/wangweihong/eazycloud/internal/pkg/libs/huaweicloud"
	"github.com/wangweihong/gotoolbox/pkg/json"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/bss/v2/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewBusiness_ShowCustomerAccountBalances(t *testing.T) {
	Convey("", t, func() {
		cfg, err := huaweicloud.NewAccessConfig(ak, sk, "")
		So(err, ShouldBeNil)

		c, err := huaweicloud.NewBusiness(cfg)
		So(err, ShouldBeNil)
		resp, err := c.ShowCustomerAccountBalances(context.Background(), &model.ShowCustomerAccountBalancesRequest{})
		So(err, ShouldBeNil)
		json.PrintStructObject(resp)
	})
}

func TestNewBusiness_ShowCustomerMonthlySum(t *testing.T) {
	Convey("", t, func() {
		cfg, err := huaweicloud.NewAccessConfig(ak, sk, "")
		So(err, ShouldBeNil)

		c, err := huaweicloud.NewBusiness(cfg)
		So(err, ShouldBeNil)
		resp, err := c.ShowCustomerMonthlySum(context.Background(), &model.ShowCustomerMonthlySumRequest{
			BillCycle: "2024-01",
		})
		So(err, ShouldBeNil)
		json.PrintStructObject(resp)
	})
}

func TestNewBusiness_ListPartnerAccountChangeRecordsResponse(t *testing.T) {
	Convey("", t, func() {
		cfg, err := huaweicloud.NewAccessConfig(ak, sk, "")
		So(err, ShouldBeNil)

		c, err := huaweicloud.NewBusiness(cfg)
		So(err, ShouldBeNil)
		resp, err := c.ListCustomerAccountChangeRecords(context.Background(), &model.ListCustomerAccountChangeRecordsRequest{
			BalanceType: "BALANCE_TYPE_DEBIT",
		})
		So(err, ShouldBeNil)
		json.PrintStructObject(resp)
	})
}

func TestNewBusiness_ListFreeResourcesUsageRecords(t *testing.T) {
	Convey("", t, func() {
		cfg, err := huaweicloud.NewAccessConfig(ak, sk, "")
		So(err, ShouldBeNil)

		c, err := huaweicloud.NewBusiness(cfg)
		So(err, ShouldBeNil)
		resp, err := c.ListFreeResourcesUsageRecords(context.Background(), &model.ListFreeResourcesUsageRecordsRequest{
			FreeResourceId:   nil,
			ProductId:        nil,
			ResourceTypeCode: nil,
			DeductTimeBegin:  "2024-01-01",
			DeductTimeEnd:    "2024-01-11",
			Offset:           nil,
			Limit:            nil,
		})
		So(err, ShouldBeNil)
		json.PrintStructObject(resp)
	})
}
