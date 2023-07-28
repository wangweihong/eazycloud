package requestid

import (
	"context"

	"google.golang.org/grpc"

	"github.com/wangweihong/eazycloud/pkg/tracectx"
)

// UnaryServerInterceptor returns a new unary server interceptor for trace.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		// TODO: how to get trace id from  request
		ctx = tracectx.WithTraceIDContext(ctx)
		// 调用下一个拦截器或最终的RPC处理程序
		resp, err := handler(ctx, req)
		return resp, err
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for trace.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		// TODO: how to trace stream request?
		return handler(srv, stream)
	}
}
