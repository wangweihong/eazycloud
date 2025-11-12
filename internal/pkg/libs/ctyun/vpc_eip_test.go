package ctyun_test

//
//import (
//	"context"
//	"testing"
//
//	. "github.com/smartystreets/goconvey/convey"
//	"github.com/wangweihong/eazycloud/internal/pkg/multicloud/ctyun"
//	"github.com/wangweihong/eazycloud/internal/pkg/multicloud/ctyun/ictyun"
//	"github.com/wangweihong/gotoolbox/pkg/json"
//	"github.com/wangweihong/gotoolbox/log"
//)
//
//func TestEip_EipList(t *testing.T) {
//	Convey("EipList", t, func() {
//		c, err := ctyun.NewClient(endpoint, ak, sk)
//		So(err, ShouldBeNil)
//
//		Convey("列表", func() {
//			resp, err := ctyun.NewEip(c).EipList(context.Background(), &ictyun.EipListRequest{
//				RegionID:    region,
//				ClientToken: clientToken,
//			})
//			So(err, ShouldBeNil)
//			log.Infof("len:%v", len(resp.ReturnObj.Eips))
//		})
//	})
//}
//
//func TestEip_EipCreate(t *testing.T) {
//	Convey("EipCreate", t, func() {
//		c, err := ctyun.NewClient(endpoint, ak, sk)
//		So(err, ShouldBeNil)
//
//		Convey("创建", func() {
//			resp, err := ctyun.NewEip(c).EipCreate(context.Background(), &ictyun.EipCreateRequest{
//				RegionID:          region,
//				ClientToken:       clientToken,
//				CycleType:         "on_demand",
//				Name:              "test1",
//				CycleCount:        nil,
//				Bandwidth:         nil,
//				BandwidthID:       nil,
//				DemandBillingType: nil,
//				ProjectID:         nil,
//			})
//			So(err, ShouldBeNil)
//			json.PrintStructObject(resp)
//		})
//	})
//}
