package example

import (
	"context"

	"github.com/wangweihong/eazycloud/pkg/httpcli"
)

type UserAPI interface {
	Create(
		ctx context.Context,
		req *UserRequest,
		opts ...httpcli.CallOption,
	) (*UserResponse, error)
}

type UserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UserResponse struct{}
