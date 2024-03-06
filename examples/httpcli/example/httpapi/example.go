package httpapi

import (
	"context"
	"net/http"

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
	arg := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.address).
		WithMethod("POST").
		WithPath("/user/create").
		WithBody("", req).Build()

	reply, err := p.c.Invoke(ctx, arg, req, resp, opts...)
	if err != nil {
		return nil, errors.WrapError(code.ErrHTTPError, err)
	}

	if reply.GetStatusCode() != http.StatusOK {
		return nil, errors.Wrap(code.ErrHTTPError, "status code not 200")
	}

	if err = reply.Decode(resp); err != nil {
		return nil, err
	}
	return resp, nil
}
