package huaweicloud_test

import (
	"context"
	"testing"

	"github.com/wangweihong/eazycloud/internal/pkg/libs/huaweicloud"
	"github.com/wangweihong/gotoolbox/pkg/typeutil"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"
	"github.com/wangweihong/gotoolbox/pkg/json"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMonitorListMetrics(t *testing.T) {
	Convey("metrics list", t, func() {
		SkipConvey("获取指定服务器的指标列表", func() {
			resp, err := huaweicloud.NewMonitor(ac).MonitorListMetrics(context.Background(), &model.ListMetricsRequest{
				Dim0: typeutil.String("instance_id,3dea2e88-92ef-487f-a08c-12a99c53480a"),
			})
			So(err, ShouldBeNil)
			json.PrintStructObject(resp)
		})
	})
}
