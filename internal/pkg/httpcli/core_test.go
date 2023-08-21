package httpcli_test

import (
	"context"
	"testing"
	"time"

	"github.com/wangweihong/eazycloud/internal/pkg/genericserver"
	"github.com/wangweihong/eazycloud/internal/pkg/httpcli"
	"github.com/wangweihong/eazycloud/pkg/log"
	"github.com/wangweihong/eazycloud/pkg/version"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	opts := log.NewOptions()
	opts.OutputPaths = nil
	opts.ErrorOutputPaths = nil
	opts.Level = "debug"
	log.Init(opts)
}

func installServer(conf *genericserver.Config) *genericserver.GenericHTTPServer {
	s, err := conf.Complete().New()
	So(err, ShouldBeNil)
	go func() {
		s.Run()
	}()
	// Wait for the server to start (you can use a more sophisticated wait mechanism)
	time.Sleep(5 * time.Second)
	return s
}

func TestClient_Invoke(t *testing.T) {
	Convey("客户端调用", t, func() {
		conf := genericserver.NewConfig()
		conf.Healthz = true
		conf.Version = true
		conf.EnableMetrics = false
		conf.InsecureServing = &genericserver.InsecureServingInfo{
			Address:  "0.0.0.0:57217",
			Required: true,
		}
		s := installServer(conf)
		defer s.Close()

		Convey("GET请求", func() {
			c, err := httpcli.NewClient("http://0.0.0.0:57217")
			So(err, ShouldBeNil)
			ctx := context.Background()

			vi := version.Info{}
			resp, err := c.Invoke(ctx, "GET", "/version", nil, &vi)
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, 200)
			So(vi.Platform, ShouldNotBeEmpty)
		})

	})
}

func TestClient_Interceptor(t *testing.T) {
	Convey("拦截器", t, func() {
		conf := genericserver.NewConfig()
		conf.Healthz = true
		conf.Version = true
		conf.EnableMetrics = false
		conf.InsecureServing = &genericserver.InsecureServingInfo{
			Address:  "0.0.0.0:57218",
			Required: true,
		}
		s := installServer(conf)
		defer s.Close()
		Convey("拦截器", func() {
			inter1 := func(ctx context.Context, method string, rawURL string, arg, reply interface{}, cc *httpcli.Client, invoker httpcli.Invoker, opts ...httpcli.CallOption) (*httpcli.RawResponse, error) {
				opt := httpcli.QueryCallOption(map[string]string{
					"inter1": "b",
				})
				opts = append(opts, opt)

				resp, err := invoker(ctx, method, rawURL, arg, reply, cc, opts...)
				if err != nil {
					return resp, err
				}
				resp.Header.Set("inter1", "bbbb")
				return resp, err
			}
			inter2 := func(ctx context.Context, method string, rawURL string, arg, reply interface{}, cc *httpcli.Client, invoker httpcli.Invoker, opts ...httpcli.CallOption) (*httpcli.RawResponse, error) {
				opt := httpcli.QueryCallOption(map[string]string{
					"inter2": "bbbb",
				})
				opts = append(opts, opt)

				ctx = context.WithValue(ctx, "inter2", "bbbb")
				resp, err := invoker(ctx, method, rawURL, arg, reply, cc, opts...)
				if err != nil {
					return resp, err
				}
				resp.Header.Set("inter2", "bbbb")
				return resp, err
			}

			c, err := httpcli.NewClient("http://0.0.0.0:57218",
				httpcli.WithIntercepts(inter1, inter2))
			So(err, ShouldBeNil)
			ctx := context.Background()

			vi := version.Info{}
			resp, err := c.Invoke(ctx, "GET", "/version", nil, &vi)
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, 200)
			So(vi.Platform, ShouldNotBeEmpty)

			So(resp.Header.Get("inter1"), ShouldEqual, "bbbb")
			So(resp.Header.Get("inter2"), ShouldEqual, "bbbb")
		})
	})
}
