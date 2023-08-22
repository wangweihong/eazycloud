package example

import (
	"context"

	"github.com/wangweihong/eazycloud/pkg/grpccli"

	"google.golang.org/grpc"

	"github.com/wangweihong/eazycloud/pkg/errors"
	"github.com/wangweihong/eazycloud/pkg/grpcproto/apis/debug"
	"github.com/wangweihong/eazycloud/pkg/grpcproto/apis/version"
)

type Backend interface {
	Version(ctx context.Context, in *version.VersionRequest, opt ...grpc.CallOption) (*version.VersionResponse, error)
	Example(ctx context.Context, in *debug.ExampleRequest, opt ...grpc.CallOption) (*debug.ExampleResponse, error)

	Close()
}

func NewBackend(addr string, opt ...grpccli.Option) (Backend, error) {
	c, err := grpccli.NewClient(addr, opt...)
	if err != nil {
		return nil, errors.UpdateStack(err)
	}
	return &example{c: c}, nil
}

type example struct {
	c *grpccli.Client
}

func (c *example) Version(
	ctx context.Context,
	in *version.VersionRequest,
	opt ...grpc.CallOption,
) (*version.VersionResponse, error) {
	out := &version.VersionResponse{}

	if err := c.c.Call(ctx, func(ctx context.Context, conn *grpc.ClientConn) error {
		var e error
		out, e = version.NewVersionServiceClient(conn).Version(ctx, in, opt...)
		if e != nil {
			return errors.UpdateStack(e)
		}
		return nil
	}); err != nil {
		return nil, errors.UpdateStack(err)
	}
	return out, nil
}

func (c *example) Example(
	ctx context.Context,
	in *debug.ExampleRequest,
	opt ...grpc.CallOption,
) (*debug.ExampleResponse, error) {
	out := &debug.ExampleResponse{}

	if err := c.c.Call(ctx, func(ctx context.Context, conn *grpc.ClientConn) error {
		var e error
		out, e = debug.NewDebugServiceClient(conn).Example(ctx, in, opt...)
		if e != nil {
			return errors.UpdateStack(e)
		}
		return nil
	}); err != nil {
		return nil, errors.UpdateStack(err)
	}
	return out, nil
}

func (c *example) Close() {
	c.c.Close()
}
