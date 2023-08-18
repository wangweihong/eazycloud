package statuscode

import (
	"context"
	"net/http"

	"github.com/wangweihong/eazycloud/internal/pkg/code"
	"github.com/wangweihong/eazycloud/pkg/json"

	"github.com/wangweihong/eazycloud/internal/pkg/httpcli"
	"github.com/wangweihong/eazycloud/pkg/errors"
	"github.com/wangweihong/eazycloud/pkg/log"
	"github.com/wangweihong/eazycloud/pkg/skipper"
)

// 非200状态码拦截.
func NoSuccessStatusCodeInterceptor(name string, skipperFunc ...skipper.SkipperFunc) httpcli.Interceptor {
	return func(ctx context.Context, method string, rawURL string, arg, reply interface{}, cc *httpcli.Client, invoker httpcli.Invoker, opts ...httpcli.CallOption) (*httpcli.RawResponse, error) {
		log.F(ctx).Debugf("Intercepttor %s Enter", name)
		defer log.F(ctx).Debugf("Interceptor %s Finish", name)

		if skipper.Skip(rawURL, skipperFunc...) {
			log.F(ctx).Debugf("skip interceptor %s for rawrurl %s", name, rawURL)

			return invoker(ctx, method, rawURL, arg, reply, cc, opts...)
		}

		// tell `invoker` do not parse response data in invoke
		opts = append(opts, httpcli.ResponseNotParseCallOption())

		rawResp, err := invoker(ctx, method, rawURL, arg, reply, cc, opts...)
		if err != nil {
			log.F(ctx).Error("invoker fail", log.Err(err))
			return rawResp, errors.UpdateStack(err)
		}

		if rawResp.StatusCode != http.StatusOK {
			log.F(ctx).Error("invoker fail, no 200")

			return rawResp, errors.Wrap(code.ErrHTTPError, "response code is not 200")
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