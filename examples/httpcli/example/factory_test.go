package example_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/wangweihong/eazycloud/examples/httpcli/example"
	"github.com/wangweihong/eazycloud/examples/httpcli/example/httpapi"
	"github.com/wangweihong/eazycloud/examples/httpcli/example/options"
	"github.com/wangweihong/eazycloud/pkg/log"
	"github.com/wangweihong/eazycloud/pkg/tracectx"
)

func init() {
	opts := log.NewOptions()
	//opts.OutputPaths = nil
	//opts.ErrorOutputPaths = nil
	opts.Level = "info"
	log.Init(opts)
}

func TestFactory(t *testing.T) {
	Convey("Factory", t, func() {
		instance, err := httpapi.GetHttpApiFactoryOr(&options.BackendOptions{
			Address: "http://127.0.0.1:8081",
		})
		So(err, ShouldBeNil)

		ctx := context.Background()
		ctx = log.WithFieldPair(ctx, "routineID", tracectx.NewTraceID())

		_, _ = instance.Users().Create(ctx, &example.UserRequest{
			Name: "Test",
			Age:  1000,
		})
	})
}
