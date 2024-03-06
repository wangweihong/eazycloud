package httpapi

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/wangweihong/eazycloud/pkg/httpcli"

	"github.com/wangweihong/eazycloud/pkg/httpcli/interceptorcli"

	"github.com/wangweihong/eazycloud/examples/httpcli/example/options"
	"github.com/wangweihong/eazycloud/pkg/code"
	"github.com/wangweihong/eazycloud/pkg/errors"
	"github.com/wangweihong/eazycloud/pkg/log"
	"github.com/wangweihong/eazycloud/pkg/skipper"

	"github.com/wangweihong/eazycloud/examples/httpcli/example"
)

type client struct {
	*httpcli.Client
	address string
	timeout time.Duration
}

func (c *client) Users() example.UserAPI {
	return newUser(c)
}

var (
	httpApiFactory example.Factory
	once           sync.Once
)

// GetHttpApiFactoryOr create dingtalkapi factory with the given config.
func GetHttpApiFactoryOr(opts *options.BackendOptions) (example.Factory, error) {
	if opts == nil && httpApiFactory == nil {
		return nil, fmt.Errorf("failed to get example api factory")
	}

	var c *httpcli.Client
	var err error
	once.Do(func() {
		// 建立长连接?
		HTTPTransport := &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second, // 连接超时时间
				KeepAlive: 60 * time.Second, // 保持长连接的时间
			}).DialContext, // 设置连接的参数
			MaxIdleConns:          500,              // 最大空闲连接
			IdleConnTimeout:       60 * time.Second, // 空闲连接的超时时间
			ExpectContinueTimeout: 30 * time.Second, // 等待服务第一个响应的超时时间
			MaxIdleConnsPerHost:   100,              // 每个host保持的空闲连接数
		}
		hc := &client{
			address: opts.Address,
			timeout: 30 * time.Second,
		}

		c, err = httpcli.NewClient(
			nil,
			httpcli.WithTransport(HTTPTransport),
			httpcli.WithTimeout(30*time.Second),
			httpcli.WithIntercepts(
				// 注意顺序, 队列也靠后的越早执行调用后
				// TokenInterceptor("TokenInterceptor", hc, skipper.AllowPathPrefixSkipper("/gettoken")),
				// ErrorCodeInterceptor(),
				interceptorcli.StatusCodeInterceptor(""),
				interceptorcli.LoggingInterceptor(""),
			),
		)
		hc.Client = c

		httpApiFactory = hc
	})

	if httpApiFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get dingtalkapi factory, httpApiFactory: %+v, error: %w", httpApiFactory, err)
	}

	return httpApiFactory, nil
}

// 错误码拦截.
func ErrorCodeInterceptor(skipperFunc ...skipper.SkipperFunc) httpcli.Interceptor {
	name := "ErrorCode"
	return func(ctx context.Context, req *httpcli.HttpRequest, arg, reply interface{}, cc *httpcli.Client, invoker httpcli.Invoker, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
		log.F(ctx).Debugf("Interceptor %s Enter", name)
		defer log.F(ctx).Debugf("Interceptor %s Finish", name)

		if skipper.Skip(req.GetPath(), skipperFunc...) {
			log.F(ctx).Debugf("skip interceptor %s for rawrurl %s", name, req.GetHeaderParams())

			return invoker(ctx, req, arg, reply, cc, opts...)
		}

		rawResp, err := invoker(ctx, req, arg, reply, cc, opts...)
		if err != nil {
			return rawResp, errors.UpdateStack(err)
		}

		type ErrorResponse struct {
			ErrorMessage string `json:"errmsg"`  // 返回码描述
			ErrorCode    int64  `json:"errcode"` // 返回码. 0表示成功
		}

		var er ErrorResponse
		if err := rawResp.Decode(&er); err != nil {
			return rawResp, err
		}

		if er.ErrorCode != 0 {
			return rawResp, errors.WrapError(
				code.ErrHTTPError,
				fmt.Errorf("got err code %d,msg:%s", er.ErrorCode, er.ErrorMessage),
			)
		}
		return rawResp, nil
	}
}
