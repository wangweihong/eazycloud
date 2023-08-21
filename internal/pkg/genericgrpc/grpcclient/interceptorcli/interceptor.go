package interceptorcli

import (
	"google.golang.org/grpc"

	"github.com/wangweihong/eazycloud/pkg/skipper"

	"github.com/wangweihong/eazycloud/internal/pkg/genericgrpc/grpcclient/interceptorcli/callstatus"
	"github.com/wangweihong/eazycloud/internal/pkg/genericgrpc/grpcclient/interceptorcli/logging"
)

const (
	InterceptorNameLogger     = "logger"
	InterceptorNameCallStatus = "callstatus"
)

var (
	UnaryClientInterceptorList  = defaultUnaryClientInterceptorList()
	UnaryClientInterceptorNames = defaultInterceptorListNames()
)

func defaultUnaryClientInterceptorList(skipperFunc ...skipper.SkipperFunc) map[string]grpc.UnaryClientInterceptor {
	return map[string]grpc.UnaryClientInterceptor{
		InterceptorNameLogger:     logging.UnaryClientInterceptor(skipperFunc...),
		InterceptorNameCallStatus: callstatus.UnaryClientInterceptor(skipperFunc...),
	}
}

func defaultInterceptorListNames() []string {
	names := make([]string, 0, len(defaultUnaryClientInterceptorList()))
	for name := range defaultUnaryClientInterceptorList() {
		names = append(names, name)
	}
	return names
}

// 生成带有跳过条件的拦截器列表.
func GetUnaryClientInterceptorWithSkippers(skipperFunc ...skipper.SkipperFunc) map[string]grpc.UnaryClientInterceptor {
	return defaultUnaryClientInterceptorList(skipperFunc...)
}
