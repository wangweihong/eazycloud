package callstatus

import (
	"context"
	"reflect"

	"github.com/wangweihong/eazycloud/pkg/skipper"

	"google.golang.org/grpc"

	"github.com/wangweihong/eazycloud/pkg/code"
	"github.com/wangweihong/eazycloud/pkg/errors"
	"github.com/wangweihong/eazycloud/pkg/grpcproto/apis/callstatus"
	"github.com/wangweihong/eazycloud/pkg/log"
)

// UnaryClientInterceptor returns a new unary client interceptor for logging.
func UnaryClientInterceptor(skipperFunc ...skipper.SkipperFunc) grpc.UnaryClientInterceptor {
	name := "callstatus"

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.F(ctx).Debugf("Interceptor %s Enter", name)
		defer log.F(ctx).Debugf("Interceptor %s Finish", name)

		if skipper.Skip(method, skipperFunc...) {
			log.F(ctx).Debugf("skip interceptor %s for %s", name, method)

			return invoker(ctx, method, req, reply, cc, opts...)
		}

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			return errors.UpdateStack(err)
		}

		cs, exist := fetchCallStatusField(reply)
		if !exist {
			log.F(ctx).Errorf("`CallStatus` field not exist in response")
			return errors.Wrap(code.ErrGRPCResponseDataParseError, "`CallStatus` field not exist in response")
		}

		if cs == nil {
			log.F(ctx).Errorf("CallStatus is nil")
			return errors.Wrap(code.ErrGRPCResponseDataParseError, "CallStatus is nil")
		}

		if err := callstatus.ToError(cs); err != nil {
			log.F(ctx).Error(err.Error())
			return errors.UpdateStack(err)
		}

		return nil
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
