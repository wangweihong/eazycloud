package tracectx_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/wangweihong/eazycloud/pkg/tracectx"
)

type (
	TraceIDKey struct{} // store traceID in context
)

func TestFromTraceIDContext(t *testing.T) {
	Convey("FromTraceIDContext", t, func() {
		ctx := tracectx.NewTraceIDContext(context.Background(), "12345")
		ctx = context.WithValue(ctx, TraceIDKey{}, "54321")
		So(tracectx.FromTraceIDContext(ctx), ShouldNotEqual, ctx.Value(TraceIDKey{}))
	})
}
