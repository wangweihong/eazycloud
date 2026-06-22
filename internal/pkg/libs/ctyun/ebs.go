package ctyun

import (
	"context"

	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"

	"github.com/wangweihong/eazycloud/apis/ictyun"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
)

type Ebs struct {
	c           *Client
	serviceType string
}

// EbsList 查询云硬盘列表
func (p *Ebs) EbsList(
	ctx context.Context,
	req *ictyun.EbsListRequest,
	resp *ictyun.EbsListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/ebs/list-ebs").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsGetByName 基于名称查询云硬盘列表
func (p *Ebs) EbsGetByName(
	ctx context.Context,
	req *ictyun.EbsGetByNameRequest,
	resp *ictyun.EbsGetResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/ebs/info-by-name-ebs").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsGetByID 基于磁盘ID查询云硬盘列表
func (p *Ebs) EbsGetByID(
	ctx context.Context,
	req *ictyun.EbsGetByIDRequest,
	resp *ictyun.EbsGetResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/ebs/info-ebs").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsCreate 创建云硬盘
func (p *Ebs) EbsCreate(
	ctx context.Context,
	req *ictyun.EbsCreateRequest,
	resp *ictyun.EbsCreateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs/new-ebs").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// EbsRefund 退订/删除云硬盘
func (p *Ebs) EbsRefund(
	ctx context.Context,
	req *ictyun.EbsRefundRequest,
	resp *ictyun.EbsRefundResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs/refund-ebs").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// EbsResize 扩容云硬盘
func (p *Ebs) EbsResize(
	ctx context.Context,
	req *ictyun.EbsResizeRequest,
	resp *ictyun.EbsResizeResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs/resize-ebs").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsRenew  续订云硬盘
func (p *Ebs) EbsRenew(
	ctx context.Context,
	req *ictyun.EbsRenewRequest,
	resp *ictyun.EbsRenewResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs/renew-ebs").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsUpdate 更新云硬盘属性
func (p *Ebs) EbsUpdate(
	ctx context.Context,
	req *ictyun.EbsUpdateRequest,
	resp *ictyun.EbsUpdateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {

	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs/update-attr-ebs").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsAttach 挂载云硬盘
func (p *Ebs) EbsAttach(
	ctx context.Context,
	req *ictyun.EbsAttachRequest,
	resp *ictyun.EbsAttachResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {

	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs/attach-ebs").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsDetach 卸载云硬盘
func (p *Ebs) EbsDetach(
	ctx context.Context,
	req *ictyun.EbsDetachRequest,
	resp *ictyun.EbsDetachResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {

	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs/detach-ebs").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}
