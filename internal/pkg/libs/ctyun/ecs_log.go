package ctyun

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/ictyun"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"
)

// LogList 查询操作日志列表
func (p *Ecs) LogList(
	ctx context.Context,
	req *ictyun.LogListRequest,
	resp *ictyun.LogListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/op-log/query").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}
