package ctyun

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/ictyun"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"
)

type Oss struct {
	c           *Client
	serviceType string
}

// OssBucketList 查询存储桶
func (p *Oss) OssBucketList(
	ctx context.Context,
	req *ictyun.OssBucketListRequest,
	resp *ictyun.OssBucketListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/oss/list-buckets").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// OssKeyList 查询存储桶awk访问密钥
func (p *Oss) OssKeyList(
	ctx context.Context,
	req *ictyun.OssKeyRequest,
	resp *ictyun.OssKeyResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/oss/get-keys").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// OssEndpointList 查询存储桶awk访问端点
func (p *Oss) OssEndpointList(
	ctx context.Context,
	req *ictyun.OssEndpointRequest,
	resp *ictyun.OssEndpointResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/oss/get-endpoint").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// OssBucketGet 查询存储桶
func (p *Oss) OssBucketGet(
	ctx context.Context,
	req *ictyun.OssBucketGetRequest,
	resp *ictyun.OssBucketGetResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/oss/get-bucket-info").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// OssBucketGetACL 查询存储桶ACL
func (p *Oss) OssBucketGetACL(
	ctx context.Context,
	req *ictyun.OssBucketGetACLRequest,
	resp *ictyun.OssBucketGetACLResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/oss/get-bucket-acl").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// OssBucketCreate 创建存储桶
func (p *Oss) OssBucketCreate(
	ctx context.Context,
	req *ictyun.OssBucketCreateRequest,
	resp *ictyun.OssBucketCreateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/oss/create-bucket").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// OssBucketDelete 删除存储桶
func (p *Oss) OssBucketDelete(
	ctx context.Context,
	req *ictyun.OssBucketDeleteRequest,
	resp *ictyun.OssBucketDeleteResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/oss/delete-bucket").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// OssDirectoryCreate 创建存储桶目录
func (p *Oss) OssDirectoryCreate(
	ctx context.Context,
	req *ictyun.OssDirectoryCreateRequest,
	resp *ictyun.OssDirectoryCreateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/oss/create-directory").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// OssDirectoryDelete 删除存储桶目录
func (p *Oss) OssDirectoryDelete(
	ctx context.Context,
	req *ictyun.OssDirectoryDeleteRequest,
	resp *ictyun.OssDirectoryDeleteResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/oss/delete-directory").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// OssObjectList 查询存储桶对象列表
func (p *Oss) OssObjectList(
	ctx context.Context,
	req *ictyun.OssObjectListRequest,
	resp *ictyun.OssObjectListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/oss/list-objects").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// OssObjectNum 查询存储桶对象数量
func (p *Oss) OssObjectNum(
	ctx context.Context,
	req *ictyun.OssObjectNumRequest,
	resp *ictyun.OssObjectNumResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/oss/get-object-num").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// OssObjectDelete 删除存储桶对象
func (p *Oss) OssObjectDelete(
	ctx context.Context,
	req *ictyun.OssObjectDeleteRequest,
	resp *ictyun.OssObjectDeleteResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/oss/delete-object").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// OssObjectDownloadLink 获取存储桶对象下载链接
func (p *Oss) OssObjectDownloadLink(
	ctx context.Context,
	req *ictyun.OssObjectDownloadLinkRequest,
	resp *ictyun.OssObjectDownloadLinkResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/oss/generate-object-download-link").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}
