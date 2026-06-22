package ctyun_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wangweihong/eazycloud/apis/ictyun"
	"github.com/wangweihong/eazycloud/internal/pkg/libs/ctyun"
)

const (
	vpc = ""
)

func TestVpc_VpcList(t *testing.T) {
	Convey("VpcList", t, func() {
		_, err := ctyun.NewClient(ak, sk).Vpcs().VpcList(context.Background(), &ictyun.VpcListRequest{
			RegionID: region,
		}, nil)
		So(err, ShouldBeNil)

	})
}

func TestVpc_VpcSubnetList(t *testing.T) {
	Convey("VpcSubnetList", t, func() {
		_, err := ctyun.NewClient(ak, sk).Vpcs().VpcSubnetList(context.Background(), &ictyun.VpcSubnetListRequest{
			RegionID: region,
			VpcID:    vpc,
		}, nil)
		So(err, ShouldBeNil)
	})
}
