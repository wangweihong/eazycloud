package log_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/wangweihong/eazycloud/pkg/log"
)

func TestWithFieldPair(t *testing.T) {
	defer log.Flush()

	var ctx = context.Background()
	Convey("TestWithFieldPair", t, func() {
		Convey("", func() {
			ctx = log.WithFieldPair(ctx, "aaa", "bbb")
			v := ctx.Value(log.FieldKeyCtx{})
			fields, ok := v.(map[string]interface{})
			So(ok, ShouldBeTrue)
			So(fields, ShouldNotBeNil)

			d := fields["aaa"].(string)
			So(d, ShouldEqual, "bbb")
		})
	})
}
