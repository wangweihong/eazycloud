package huaweicloud_test

import (
	"context"
	"testing"

	"github.com/wangweihong/eazycloud/internal/pkg/libs/huaweicloud"
	"github.com/wangweihong/gotoolbox/pkg/json"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestImageList(t *testing.T) {
	Convey("image list", t, func() {
		Convey("1", func() {
			ac, err := huaweicloud.NewAccessConfig(ak, sk, region)
			So(err, ShouldBeNil)

			c, err := huaweicloud.NewImage(ac)
			So(err, ShouldBeNil)
			resp, err := c.ImageList(context.Background(), &model.ListImagesRequest{})
			So(err, ShouldBeNil)
			json.PrintStructObject(resp)
		})
	})
}
