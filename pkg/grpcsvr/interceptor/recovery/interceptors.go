package recovery

import (
	"context"

	"github.com/wangweihong/eazycloud/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/wangweihong/eazycloud/pkg/log"
)

// RecoveryHandlerFunc is a function that recovers from the panic `p` by returning an `error`.
type RecoveryHandlerFunc func(p interface{}) (_ interface{}, err error)

// RecoveryHandlerFuncContext is a function that recovers from the panic `p` by returning an `error`.
// The context can be used to extract request scoped metadata and context values.
type RecoveryHandlerFuncContext func(ctx context.Context, p interface{}) (_ interface{}, err error)

// UnaryServerInterceptor returns a new unary server interceptor for panic recovery.
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	name := "recovery"

	o := evaluateOptions(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log.F(ctx).Debugf("Interceptor %s Enter", name)
		defer log.F(ctx).Debugf("Interceptor %s Finish", name)

		panicked := true

		defer func() {
			if r := recover(); r != nil || panicked {
				log.F(ctx).Errorf("method `%v` crashed with panic: %v", info.FullMethod, r)
				resp, err = recoverFrom(ctx, r, o.recoveryHandlerFunc)
			}
		}()

		resp, err = handler(ctx, req)
		panicked = false
		return resp, errors.UpdateStack(err)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for panic recovery.
func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	o := evaluateOptions(opts)
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		panicked := true

		defer func() {
			if r := recover(); r != nil || panicked {
				_, err = recoverFrom(stream.Context(), r, o.recoveryHandlerFunc)
			}
		}()

		err = handler(srv, stream)
		panicked = false
		return err
	}
}

func recoverFrom(ctx context.Context, p interface{}, r RecoveryHandlerFuncContext) (interface{}, error) {
	if r == nil {
		return nil, status.Errorf(codes.Internal, "%v", p)
	}
	return r(ctx, p)
}
