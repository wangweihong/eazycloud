package httpapi

import (
	"context"

	"github.com/wangweihong/eazycloud/examples/httpcli/example"
	"github.com/wangweihong/eazycloud/pkg/code"
	"github.com/wangweihong/eazycloud/pkg/errors"
	"github.com/wangweihong/eazycloud/pkg/httpcli"
)

type user struct {
	c *client
}

func newUser(c *client) *user {
	return &user{
		c: c,
	}
}

func (p *user) Create(
	ctx context.Context,
	req *example.UserRequest,
	opts ...httpcli.CallOption,
) (*example.UserResponse, error) {
	if req == nil {
		return nil, errors.Wrap(code.ErrValidation, "UserListReq is empty")
	}

	resp := &example.UserResponse{}
	_, err := p.c.Invoke(ctx, "POST", "/user/create", req, resp, opts...)
	if err != nil {
		return nil, errors.WrapError(code.ErrHTTPError, err)
	}
	return resp, nil
}
