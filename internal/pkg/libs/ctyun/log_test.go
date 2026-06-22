package ctyun_test

//import (
//	"context"
//	"testing"
//
//	. "github.com/smartystreets/goconvey/convey"
//	"github.com/wangweihong/eazycloud/internal/pkg/multicloud/ctyun"
//	"github.com/wangweihong/eazycloud/internal/pkg/multicloud/ctyun/ictyun"
//	"github.com/wangweihong/gotoolbox/pkg/json"
//)
//
//func TestLog_LogList(t *testing.T) {
//	Convey("LogList", t, func() {
//		c, err := ctyun.NewClient(endpoint, ak, sk)
//		So(err, ShouldBeNil)
//		resp, err := ctyun.NewLog(c).LogList(context.Background(), &ictyun.LogListRequest{
//			RegionID: region,
//		})
//		So(err, ShouldBeNil)
//		json.PrintStructObject(resp)
//	})
//}
