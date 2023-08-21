package callstatus

import (
	"github.com/wangweihong/eazycloud/internal/pkg/code"
	"github.com/wangweihong/eazycloud/pkg/errors"
)

// ToError convert grpc call status to err.
func ToError(cs *CallStatus) error {
	if cs == nil || cs.Code == int64(code.ErrSuccess) {
		return nil
	}

	return errors.WrapStack(int(cs.Code), cs.Description, cs.Stack)
}

// FromError convert err to grpc call status.
func FromError(err error) *CallStatus {
	e := errors.FromError(err)
	if e == nil {
		e = errors.Wrap(code.ErrSuccess, "")
	}

	cs := &CallStatus{
		Code:        int64(e.Code()),
		Message:     e.Message(),
		Stack:       e.Stack(),
		Description: e.Description(),
	}

	if errors.IsCode(e, code.ErrSuccess) {
		cs.Stack = nil
	}

	return cs
}
