package libs

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/igenregistry"
	"github.com/wangweihong/eazycloud/apis/imachinery"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"
)

type GenenralRegistryClient struct {
	endpoint string
	c        *httpcli.Client
	err      error
}

func NewGenenralRegistryClient(endpoint string, opts ...httpcli.Option) *GenenralRegistryClient {
	callOpts := DefaultCallOptions()
	if opts != nil {
		callOpts = opts
	}

	c, err := httpcli.NewClient(nil, callOpts...)
	return &GenenralRegistryClient{
		endpoint: endpoint,
		c:        c,
		err:      err,
	}
}

func (c *GenenralRegistryClient) ListImages(ctx context.Context, req *imachinery.Empty, resp *igenregistry.ListImagesResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		WithPath("/v2/_catalog").
		Build()
	return invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
}

func (c *GenenralRegistryClient) GetImageTags(ctx context.Context, req *igenregistry.GetImageTagsRequest, resp *igenregistry.GetImageTagsResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		WithPath("/v2/{repo}/tags/list").
		AddPathParam("repo", req.Repo).
		Build()
	return invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
}

func (c *GenenralRegistryClient) GetImageManifests(ctx context.Context, req *igenregistry.GetImageManifestsRequest, resp *igenregistry.GetImageManifestsResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		// 要求返回v2版本的镜像信息
		AddHeaderParam("Accept", "application/vnd.docker.distribution.manifest.v2+json").
		WithPath("/v2/{repo}/manifests/{tag}").
		AddPathParam("repo", req.Repo).
		AddPathParam("tag", req.Tag).
		Build()
	return invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
}
