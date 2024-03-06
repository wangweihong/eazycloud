package interceptorcli

import (
	"context"
	"net/http"

	"github.com/wangweihong/eazycloud/pkg/code"
	"github.com/wangweihong/eazycloud/pkg/errors"
	"github.com/wangweihong/eazycloud/pkg/httpcli"
	"github.com/wangweihong/eazycloud/pkg/log"
	"github.com/wangweihong/eazycloud/pkg/skipper"
)

func StatusCodeInterceptor(name string, skipperFunc ...skipper.SkipperFunc) httpcli.Interceptor {
	return func(ctx context.Context, req *httpcli.HttpRequest, arg, reply interface{}, cc *httpcli.Client,
		invoker httpcli.Invoker, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {

		if skipper.Skip(req.GetPath(), skipperFunc...) {
			log.F(ctx).Debugf("skip interceptor %s for rawrurl %s", name, req.GetPath())

			return invoker(ctx, req, arg, reply, cc, opts...)
		}
		rawResp, err := invoker(ctx, req, arg, reply, cc, opts...)
		if err != nil {
			return rawResp, errors.UpdateStack(err)
		}

		if rawResp.GetStatusCode() != http.StatusOK {
			return rawResp, errors.Wrap(code.ErrHTTPError, "response code is not 200")
		}

		return rawResp, nil
	}
}
