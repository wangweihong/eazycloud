package interceptor

import (
	"google.golang.org/grpc"

	"github.com/wangweihong/eazycloud/internal/pkg/grpcserver/interceptor/context"
	"github.com/wangweihong/eazycloud/internal/pkg/grpcserver/interceptor/logging"
	"github.com/wangweihong/eazycloud/internal/pkg/grpcserver/interceptor/recovery"
	"github.com/wangweihong/eazycloud/internal/pkg/grpcserver/interceptor/requestid"
)

const (
	InterceptorNameContext   = "context"
	InterceptorNameRequestID = "requestid"
	InterceptorNameRecovery  = "recovery"
	InterceptorNameLogger    = "logger"
)

var (
	UnaryServerInterceptorList  = defaultUnaryServerInterceptorList()
	UnaryServerInterceptorNames = defaultInterceptorListNames()
)

func defaultUnaryServerInterceptorList() map[string]grpc.UnaryServerInterceptor {
	return map[string]grpc.UnaryServerInterceptor{
		InterceptorNameContext:   context.UnaryServerInterceptor(),
		InterceptorNameRequestID: requestid.UnaryServerInterceptor(),
		InterceptorNameRecovery: recovery.UnaryServerInterceptor(
			recovery.WithRecoveryHandlerContext(recovery.CustomPanicHandler),
		),
		InterceptorNameLogger: logging.UnaryServerInterceptor(),
	}
}

func defaultInterceptorListNames() []string {
	names := make([]string, 0, len(defaultUnaryServerInterceptorList()))
	for name := range defaultUnaryServerInterceptorList() {
		names = append(names, name)
	}
	return names
}
