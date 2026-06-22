package ctyun

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/ictyun"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"
)

type Image struct {
	c           *Client
	serviceType string
}

// ImageList 查询镜像
func (p *Image) ImageList(
	ctx context.Context,
	req *ictyun.ImageListRequest,
	resp *ictyun.ImageListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {

	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/image/list").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// ImageGet 查询镜像详情
func (p *Image) ImageGet(
	ctx context.Context,
	req *ictyun.ImageGetRequest,
	resp *ictyun.ImageGetResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/image/detail").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// ImageDataDiskCreate 使用指定的云主机数据盘来创建一份私有镜像。
func (p *Image) ImageDataDiskCreate(
	ctx context.Context,
	req *ictyun.ImageDataDiskCreateRequest,
	resp *ictyun.ImageDataDiskCreateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/image/create-from-data-disk").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// ImageSystemDiskCreate 使用指定的云主机系统盘来创建一份私有镜像。
func (p *Image) ImageSystemDiskCreate(
	ctx context.Context,
	req *ictyun.ImageSystemDiskCreateRequest,
	resp *ictyun.ImageSystemDiskCreateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/image/create").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// ImageDelete 删除私有镜像
func (p *Image) ImageDelete(
	ctx context.Context,
	req *ictyun.ImageDeleteRequest,
	resp *ictyun.ImageDeleteResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/image/delete").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// ImageImport 使用指定的存在对象存储（原生版）Ⅰ 型的镜像文件来创建一份私有镜像
func (p *Image) ImageImport(
	ctx context.Context,
	req *ictyun.ImageImportRequest,
	resp *ictyun.ImageImportResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/image/import").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// ImageExport 导出一份私有镜像到指定的对象存储（原生版）Ⅰ 型的桶
func (p *Image) ImageExport(
	ctx context.Context,
	req *ictyun.ImageExportRequest,
	resp *ictyun.ImageExportResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/image/export").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}
