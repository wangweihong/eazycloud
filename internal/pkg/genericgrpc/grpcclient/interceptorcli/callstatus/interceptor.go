package callstatus

import (
	"context"
	"reflect"

	"github.com/wangweihong/eazycloud/pkg/skipper"

	"google.golang.org/grpc"

	"github.com/wangweihong/eazycloud/internal/pkg/code"
	"github.com/wangweihong/eazycloud/internal/pkg/genericgrpc/grpcproto/apis/callstatus"
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

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			log.F(ctx).Error("invoker fail", log.Err(err))
			return errors.UpdateStack(err)
		}

		cs, exist := fetchCallStatusField(reply)
		if !exist {
			return errors.Wrap(code.ErrGRPCResponseDataParseError, "`CallStatus` field not exist in response")
		}

		if cs == nil {
			return errors.Wrap(code.ErrGRPCResponseDataParseError, "CallStatus is nil")
		}

		return callstatus.ToError(cs)
	}
}

func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return nil, nil
	}
}

func fetchCallStatusField(v interface{}) (*callstatus.CallStatus, bool) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return nil, false
	}
	field := rv.FieldByName("CallStatus")
	if !field.IsValid() {
		return nil, false
	}
	st, ok := field.Interface().(*callstatus.CallStatus)
	if !ok {
		return nil, false
	}
	return st, true
}
