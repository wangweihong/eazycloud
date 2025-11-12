package libs

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/ikubeagent"
	"github.com/wangweihong/eazycloud/apis/imachinery"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"
)

type KubeagentClient struct {
	endpoint string
	c        *httpcli.Client
	err      error
}

func NewKubeagentClient(endpoint string, opts ...httpcli.Option) *KubeagentClient {
	callOpts := DefaultCallOptions()
	if opts != nil {
		callOpts = opts
	}

	c, err := httpcli.NewClient(nil, callOpts...)
	return &KubeagentClient{
		endpoint: endpoint,
		c:        c,
		err:      err,
	}
}

func (c *KubeagentClient) InstallMaster(ctx context.Context, r *ikubeagent.InstallMasterRequest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		POST().
		WithPath("/v1/deployment/install-master").
		WithBody("", r).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, r, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *KubeagentClient) JoinCluster(ctx context.Context, r *ikubeagent.JoinClusterRequest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		POST().
		WithPath("/v1/deployment/join-cluster").
		WithBody("", r).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, r, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *KubeagentClient) ResetInstall(ctx context.Context, r *imachinery.Empty, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		POST().
		WithPath("/v1/deployment/reset-deploy").
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, r, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *KubeagentClient) CheckDependency(ctx context.Context, r *ikubeagent.CheckDependencyRequest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddQueryParamByObject(r).
		WithPath("/v1/deployment/check-dependency").
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, r, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *KubeagentClient) GetInstallState(ctx context.Context, r *imachinery.Empty, resp *ikubeagent.InstallStateResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		WithPath("/v1/deployment/install-state").
		Build()

	httpResp, err := invoke(ctx, c.err, c.c, httpReq, r, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *KubeagentClient) GetKubeConfig(ctx context.Context, r *imachinery.Empty, resp *ikubeagent.GetKubeConfigResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		WithPath("/v1/deployment/kubeconfig").
		Build()

	httpResp, err := invoke(ctx, c.err, c.c, httpReq, r, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *KubeagentClient) GetInstallLog(ctx context.Context, r *imachinery.Empty, resp *ikubeagent.GetInstallLogResp, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		WithPath("/v1/deployment/install-log").
		Build()

	httpResp, err := invoke(ctx, c.err, c.c, httpReq, r, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *KubeagentClient) GetJoinCommand(ctx context.Context, req *ikubeagent.GetJoinCommandRequest, resp *ikubeagent.GetJoinCommandResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		WithPath("/v1/deployment/join-command").
		Build()

	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}
