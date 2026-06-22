package ctyun

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/ictyun"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"
)

type KeyPair struct {
	c           *Client
	serviceType string
}

func NewKeyPair(c *Client) *KeyPair {
	return &KeyPair{
		c:           c,
		serviceType: SERVICE_ECS,
	}
}

// KeyPairList 查询密钥对列表
func (p *Ecs) KeyPairList(
	ctx context.Context,
	req *ictyun.KeyPairListRequest,
	resp *ictyun.KeyPairListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/keypair/describe").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// KeyPairCreate 创建密钥对
func (p *Ecs) KeyPairCreate(
	ctx context.Context,
	req *ictyun.KeyPairCreateRequest,
	resp *ictyun.KeyPairCreateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/keypair/create-keypair").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// KeyPairDelete 删除密钥对
func (p *Ecs) KeyPairDelete(
	ctx context.Context,
	req *ictyun.KeyPairDeleteRequest,
	resp *ictyun.KeyPairDeleteResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/keypair/delete").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// KeyPairImport 导入密钥对
func (p *Ecs) KeyPairImport(
	ctx context.Context,
	req *ictyun.KeyPairImportRequest,
	resp *ictyun.KeyPairImportResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/keypair/import-keypair").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// KeyPairAttachInstance 密钥对关联实例
func (p *Ecs) KeyPairAttachInstance(
	ctx context.Context,
	req *ictyun.KeyPairAttachInstanceRequest,
	resp *ictyun.KeyPairAttachInstanceResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {

	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/keypair/attach-instance").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// KeyPairDetachInstance 密钥对取消关联实例
func (p *Ecs) KeyPairDetachInstance(
	ctx context.Context,
	req *ictyun.KeyPairDetachInstanceRequest,
	resp *ictyun.KeyPairDetachInstanceResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {

	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/keypair/detach-instance").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}
