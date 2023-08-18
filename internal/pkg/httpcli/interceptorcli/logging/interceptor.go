package logging

import (
	"context"
	"time"

	"github.com/wangweihong/eazycloud/internal/pkg/httpcli"
	"github.com/wangweihong/eazycloud/pkg/errors"
	"github.com/wangweihong/eazycloud/pkg/log"
	"github.com/wangweihong/eazycloud/pkg/skipper"
)

func LoggingInterceptor(name string, skipperFunc ...skipper.SkipperFunc) httpcli.Interceptor {
	return func(ctx context.Context, method string, rawURL string, arg, reply interface{}, cc *httpcli.Client, invoker httpcli.Invoker, opts ...httpcli.CallOption) (*httpcli.RawResponse, error) {
		log.F(ctx).Debugf("Intercepttor %s Enter", name)
		defer log.F(ctx).Debugf("Interceptor %s Finish", name)
		if skipper.Skip(rawURL, skipperFunc...) {
			log.F(ctx).Debugf("skip interceptor %s for rawrurl %s", name, rawURL)

			return invoker(ctx, method, rawURL, arg, reply, cc, opts...)
		}
		start := time.Now()
		fields := make(map[string]interface{})
		fields["req_time_begin"] = start.Format("2006-01-02 15:04:05.000000")
		fields["req_raw_url"] = rawURL
		fields["method"] = method

		// no parse response
		opts = append(opts, httpcli.ResponseNotParseCallOption())

		rawResp, err := invoker(ctx, method, rawURL, arg, reply, cc, opts...)

		end := time.Now()
		Latency := time.Since(start)
		if Latency > time.Minute {
			// Truncate in a golang < 1.8 safe way
			Latency -= Latency % time.Second
		}
		fields["req_latency"] = Latency
		fields["req_time_end"] = end.Format("2006-01-02 15:04:05.000000")

		reqURL := rawURL
		reqAddr := cc.GetAddr()
		var statusCode int
		if rawResp != nil {
			reqURL = rawResp.ReqURL
			reqAddr = rawResp.ReqAddr
			statusCode = rawResp.StatusCode
			fields["resp_status"] = rawResp.StatusCode
			fields["resp_length"] = len(rawResp.Body)
			fields["req_url"] = rawResp.ReqURL
			fields["req_media_type"] = rawResp.Header.Get("Content-Type")
			fields["req_addr"] = rawResp.ReqAddr
		}
		log.F(ctx).L(ctx).Infof("%3d - [%s] %v %s  %s", statusCode, reqAddr, Latency, method, reqURL)
		if err != nil {
			return rawResp, errors.UpdateStack(err)
		}
		return rawResp, nil
	}
}
