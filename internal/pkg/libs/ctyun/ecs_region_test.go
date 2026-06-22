package ctyun_test

// import (
// 	"context"
// 	"testing"

// 	"github.com/wangweihong/gotoolbox/pkg/typeutil"

// 	. "github.com/smartystreets/goconvey/convey"
// 	"github.com/wangweihong/eazycloud/internal/pkg/multicloud/ctyun"
// 	"github.com/wangweihong/eazycloud/internal/pkg/multicloud/ctyun/ictyun"
// )

// func TestRegion_RegionList(t *testing.T) {
// 	Convey("RegionList", t, func() {
// 		Convey("列表", func() {
// 			resp, err := ctyun.NewRegion(tc).RegionList(context.Background(), &ictyun.RegionListRequest{})
// 			So(err, ShouldBeNil)
// 			So(len(resp.ReturnObj.RegionList), ShouldNotEqual, 0)
// 		})
// 		Convey("根据名字知道指定的资源池信息", func() {
// 			resp, err := ctyun.NewRegion(tc).RegionList(context.Background(), &ictyun.RegionListRequest{RegionName: typeutil.String(rgName)})
// 			So(err, ShouldBeNil)
// 			So(resp.ReturnObj.RegionList, ShouldNotBeNil)
// 			So(len(resp.ReturnObj.RegionList), ShouldNotEqual, 0)
// 		})
// 	})
// }

// func TestRegion_RegionInfo(t *testing.T) {
// 	Convey("RegionInfo", t, func() {
// 		_, err := ctyun.NewRegion(tc).RegionInfo(context.Background(), &ictyun.RegionInfoRequest{})
// 		So(err, ShouldBeNil)
// 	})
// }

// func TestRegion_RegionResourceSummary(t *testing.T) {
// 	Convey("RegionResourceSummary", t, func() {
// 		_, err := ctyun.NewRegion(tc).RegionResourceSummary(context.Background(), &ictyun.RegionResourceSummaryRequest{
// 			RegionID: region,
// 		})
// 		So(err, ShouldBeNil)

// 	})
// }

// func TestZone_RegionZoneList(t *testing.T) {
// 	Convey("ZoneList", t, func() {
// 		_, err := ctyun.NewRegion(tc).RegionZoneList(context.Background(), &ictyun.RegionZoneListRequest{
// 			RegionID: region,
// 		})
// 		So(err, ShouldBeNil)
// 	})
// }

// func TestZone_ZoneList(t *testing.T) {
// 	Convey("ZoneList", t, func() {
// 		_, err := ctyun.NewRegion(tc).RegionResourceSummary(context.Background(), &ictyun.RegionResourceSummaryRequest{
// 			RegionID: region,
// 		})
// 		So(err, ShouldBeNil)
// 	})
// }
