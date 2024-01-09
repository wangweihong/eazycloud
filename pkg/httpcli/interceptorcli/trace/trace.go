package trace

import (
	"context"
	"github.com/wangweihong/eazycloud/pkg/httpcli"
	"github.com/wangweihong/eazycloud/pkg/log"
	"github.com/wangweihong/eazycloud/pkg/skipper"
	"github.com/wangweihong/eazycloud/pkg/tracectx"
)

// TraceInterceptor is a Intercept that injects traceID in context
func TraceInterceptor(name string, skipperFunc ...skipper.SkipperFunc) httpcli.Interceptor {
	return func(ctx context.Context, method string, rawURL string, arg, reply interface{}, cc *httpcli.Client, invoker httpcli.Invoker, opts ...httpcli.CallOption) (*httpcli.RawResponse, error) {
		log.F(ctx).Debugf("Interceptor %s Enter", name)
		defer log.F(ctx).Debugf("Interceptor %s Finish", name)
		if skipper.Skip(rawURL, skipperFunc...) {
			log.F(ctx).Debugf("skip interceptor %s for rawrurl %s", name, rawURL)
			return invoker(ctx, method, rawURL, arg, reply, cc, opts...)
		}
		traceID := tracectx.NewTraceID()
		opts = append(opts, httpcli.SetHeaderValueCallOption(tracectx.XRequestIDKey, traceID))

		ctx = context.WithValue(ctx, log.KeyRequestID, traceID)
		rawResp, err := invoker(ctx, method, rawURL, arg, reply, cc, opts...)
		return rawResp, err
	}
}
