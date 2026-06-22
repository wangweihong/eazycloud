package ctyun

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/ictyun"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"
)

// NetworkInterfaceList 查询网卡列表
func (p *Ecs) NetworkInterfaceList(
	ctx context.Context,
	req *ictyun.NetworkInterfaceListRequest,
	resp *ictyun.NetworkInterfaceListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	//部分资源池不支持
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/ecs/ports/list").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// NetworkInterfaceGet 查询网卡详情
func (p *Ecs) NetworkInterfaceGet(
	ctx context.Context,
	req *ictyun.NetworkInterfaceGetRequest,
	resp *ictyun.NetworkInterfaceGetResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/ecs/ports/show").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// NetworkInterfaceCreate 创建网卡
func (p *Ecs) NetworkInterfaceCreate(
	ctx context.Context,
	req *ictyun.NetworkInterfaceCreateRequest,
	resp *ictyun.NetworkInterfaceCreateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/ports/create").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// NetworkInterfaceDelete 删除网卡
func (p *Ecs) NetworkInterfaceDelete(
	ctx context.Context,
	req *ictyun.NetworkInterfaceDeleteRequest,
	resp *ictyun.NetworkInterfaceDeleteResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/ports/delete").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// NetworkInterfaceAttach 关联网卡到实例
func (p *Ecs) NetworkInterfaceAttach(
	ctx context.Context,
	req *ictyun.NetworkInterfaceAttachRequest,
	resp *ictyun.NetworkInterfaceAttachResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/ports/attach").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// NetworkInterfaceDetach 移除网卡和实例的关联
func (p *Ecs) NetworkInterfaceDetach(
	ctx context.Context,
	req *ictyun.NetworkInterfaceDetachRequest,
	resp *ictyun.NetworkInterfaceDetachResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/ports/detach").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// NetworkInterfaceUpdate 更改网卡信息
func (p *Ecs) NetworkInterfaceUpdate(
	ctx context.Context,
	req *ictyun.NetworkInterfaceUpdateRequest,
	resp *ictyun.NetworkInterfaceUpdateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/ports/update").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}
