package ginx

import (
	"net/http"

	"github.com/wangweihong/eazycloud/internal/pkg/code"
	"github.com/wangweihong/eazycloud/pkg/errors"
	"github.com/wangweihong/eazycloud/pkg/log"

	"github.com/gin-gonic/gin"
)

// ErrResponse defines the return messages when an error occurred.
// Reference will be omitted if it does not exist.
// swagger:model
type Response struct {
	// Status contains the detail of this request.
	// Caller should check code to determine this request is success or not.
	Status *CallStatus `json:"status"`

	// Data contains dta
	Data interface{} `json:"data,omitempty" `
}

type CallStatus struct {
	Code        int64             `json:"code,omitempty"`
	Message     map[string]string `json:"message,omitempty"`
	Stack       []string          `json:"stack,omitempty"`
	Description string            `json:"description,omitempty"`
}

// ToError convert  call status to err.
func ToError(cs *CallStatus) error {
	if cs == nil || cs.Code == int64(code.ErrSuccess) {
		return nil
	}

	return errors.WrapStack(int(cs.Code), cs.Description, cs.Stack)
}

// FromError convert err to grpc call status.
func FromError(err error) *CallStatus {
	e := errors.Wrap(code.ErrSuccess, "")
	var stack []string
	if err != nil {
		e = errors.FromError(err)
		stack = e.Stack()
	}

	cs := &CallStatus{
		Code:        int64(e.Code()),
		Message:     e.Message(),
		Stack:       stack,
		Description: e.Description(),
	}

	return cs
}

// WriteResponse write an error or the response data into http response body.
// If err is nil, return a success code to tell request is ok.
func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		e := errors.FromError(err)
		log.F(c).Errorf("%#+v", e)
		c.JSON(e.HTTPStatus(), Response{
			Status: FromError(err),
			Data:   data,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Status: FromError(err),
		Data:   data,
	})
}
