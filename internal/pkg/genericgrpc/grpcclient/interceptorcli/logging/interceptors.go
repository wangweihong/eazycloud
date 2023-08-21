package logging

import (
	"context"

	"github.com/wangweihong/eazycloud/pkg/skipper"

	"google.golang.org/grpc"

	"github.com/wangweihong/eazycloud/pkg/errors"
	"github.com/wangweihong/eazycloud/pkg/log"
)

// UnaryClientInterceptor returns a new unary client interceptor for logging.
func UnaryClientInterceptor(skipperFunc ...skipper.SkipperFunc) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if skipper.Skip(method, skipperFunc...) {
			log.F(ctx).Debugf("skip intercept method %s", method)
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		log.F(ctx).Debug("request param", log.Every("req", req), log.String("method", method))
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			log.F(ctx).Error("invoker fail", log.Err(err))
			return errors.UpdateStack(err)
		}
		log.F(ctx).Debug("response data", log.Every("out", reply))
		return nil
	}
}

func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return nil, nil
	}
}
