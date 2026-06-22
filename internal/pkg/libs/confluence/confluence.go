package libs

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/iconfluence"
	"github.com/wangweihong/eazycloud/apis/imachinery"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"
)

type ConfluenceClient struct {
	endpoint       string
	user, password string
	c              *httpcli.Client
	err            error
}

func NewConfluenceClient(endpoint string, user, password string, opts ...httpcli.Option) *ConfluenceClient {
	callOpts := libs.DefaultCallOptions()
	if opts != nil {
		callOpts = opts
	}

	c, err := httpcli.NewClient(nil, callOpts...)
	return &ConfluenceClient{
		endpoint: endpoint,
		c:        c,
		err:      err,
		user:     user,
		password: password,
	}
}

func (c *ConfluenceClient) DownloadResource(ctx context.Context, url string, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(url).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		Build()
	httpResp, err := libs.Invoke(ctx, c.err, c.c, httpReq, &imachinery.Empty{}, &imachinery.Empty{}, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *ConfluenceClient) ContentAllPages(ctx context.Context, req *imachinery.Empty, resp *iconfluence.ContentAllPagesResponse, opts ...httpcli.CallOption) error {
	start := 0
	limit := 1000
	resp = &iconfluence.ContentAllPagesResponse{}
	for {
		oneResp := &iconfluence.ContentAllPagesResponse{}
		httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
			GET().
			AddBasicAuthHeaderParam(c.user, c.password).
			AddHeaderParam("Content-Type", "application/json").
			AddHeaderParam("X-Atlassian-Token", "no-check").
			WithPath("/rest/api/content").
			AddQueryParam("expand", "body.view,history,history.lastUpdated,ancestors,metadata,space,space.homepage").
			AddQueryParam("start", start).
			AddQueryParam("limit", limit).
			Build()
		_, err := libs.Invoke(ctx, c.err, c.c, httpReq, req, oneResp, opts...)
		if err != nil {
			return errors.WithStack(err)
		}

		limit := oneResp.Limit
		resp.Results = append(resp.Results, oneResp.Results...)
		if oneResp.Size < oneResp.Limit || len(oneResp.Results) == 0 {
			break
		}
		start += limit
	}
	resp.Size = len(resp.Results)
	return nil
}

func (c *ConfluenceClient) SpaceList(ctx context.Context, req *imachinery.Empty, resp *iconfluence.SpaceListResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		AddHeaderParam("Content-Type", "application/json").
		AddHeaderParam("X-Atlassian-Token", "no-check").
		WithPath("/rest/api/space").
		Build()
	httpResp, err := libs.Invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}
