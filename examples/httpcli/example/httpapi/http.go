package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/wangweihong/eazycloud/examples/httpcli/example/options"
	"github.com/wangweihong/eazycloud/pkg/code"
	"github.com/wangweihong/eazycloud/pkg/errors"
	"github.com/wangweihong/eazycloud/pkg/httpcli/interceptorcli/logging"
	"github.com/wangweihong/eazycloud/pkg/httpcli/interceptorcli/statuscode"
	"github.com/wangweihong/eazycloud/pkg/log"
	"github.com/wangweihong/eazycloud/pkg/skipper"

	"github.com/wangweihong/eazycloud/examples/httpcli/example"
	"github.com/wangweihong/eazycloud/pkg/httpcli"
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
			opts.Address,
			httpcli.WithTransport(HTTPTransport),
			httpcli.WithTimeout(30*time.Second),
			httpcli.WithIntercepts(
				// 注意顺序, 队列也靠后的越早执行调用后
				// TokenInterceptor("TokenInterceptor", hc, skipper.AllowPathPrefixSkipper("/gettoken")),
				// ErrorCodeInterceptor(),
				statuscode.NoSuccessStatusCodeInterceptor(),
				logging.LoggingInterceptor(),
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
	return func(ctx context.Context, method string, rawURL string, arg, reply interface{}, cc *httpcli.Client, invoker httpcli.Invoker, opts ...httpcli.CallOption) (*httpcli.RawResponse, error) {
		log.F(ctx).Debugf("Interceptor %s Enter", name)
		defer log.F(ctx).Debugf("Interceptor %s Finish", name)

		if skipper.Skip(rawURL, skipperFunc...) {
			log.F(ctx).Debugf("skip interceptor %s for rawrurl %s", name, rawURL)

			return invoker(ctx, method, rawURL, arg, reply, cc, opts...)
		}

		type ErrorResponse struct {
			ErrorMessage string `json:"errmsg"`  // 返回码描述
			ErrorCode    int64  `json:"errcode"` // 返回码. 0表示成功
		}

		log.F(ctx).Debug("", log.Every("arg", arg))
		var er ErrorResponse
		rawResp, err := invoker(ctx, method, rawURL, arg, &er, cc, opts...)
		if err != nil {
			return rawResp, errors.UpdateStack(err)
		}

		if er.ErrorCode != 0 {
			return rawResp, errors.WrapError(
				code.ErrHTTPError,
				fmt.Errorf("got err code %d,msg:%s", er.ErrorCode, er.ErrorMessage),
			)
		}

		if reply != nil {
			if err := json.Unmarshal(rawResp.Body, reply); err != nil {
				log.F(ctx).Errorf("decode  err:%s", err.Error())
				return rawResp, err
			}
		}

		return rawResp, nil
	}
}
