package ctyun

import (
	"context"

	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"

	"github.com/wangweihong/eazycloud/apis/ictyun"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
)

// RegionList 资源池列表
func (p *Ecs) RegionList(
	ctx context.Context,
	req *ictyun.RegionListRequest,
	resp *ictyun.RegionListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/region/list-regions").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	if err != nil {
		return httpResp, errors.WithStack(err)
	}

	if req.RegionName != nil {
		list := resp.ReturnObj.RegionList
		resp.ReturnObj.RegionList = nil
		for _, v := range list {
			if v.RegionName == *req.RegionName {
				resp.ReturnObj.RegionList = []ictyun.RegionListEntry{v}
				break
			}
		}
	}

	return httpResp, errors.WithStack(err)
}

// RegionInfo 资源池信息
func (p *Ecs) RegionInfo(
	ctx context.Context,
	req *ictyun.RegionInfoRequest,
	resp *ictyun.RegionInfoResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/ecs/describe-regions").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	if err != nil {
		return httpResp, errors.WithStack(err)
	}

	if req.RegionName != nil {
		list := resp.ReturnObj.RegionList
		resp.ReturnObj.RegionList = nil
		for _, v := range list {
			if v.Name == *req.RegionName {
				resp.ReturnObj.RegionList = []ictyun.RegionInfoEntry{v}
				break
			}
		}
	}

	return httpResp, nil
}

// RegionProduct 资源池支持产品信息
func (p *Ecs) RegionProduct(
	ctx context.Context,
	req *ictyun.RegionInfoRequest,
	resp *ictyun.RegionInfoResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/region/get-products").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// RegionZoneList 资源池下可用区列表
func (p *Ecs) RegionZoneList(
	ctx context.Context,
	req *ictyun.RegionZoneListRequest,
	resp *ictyun.RegionZoneListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/region/get-zones").
		AddQueryParamByObject(req).
		Build()
		// 有些资源池是多可用区,有些资源池是单可用区。单可用区资源池调用这个接口不会返回可用区信息
	// 部分接口如创建弹性云主机，可用区是必传(单可用区需要传default作为参数)。
	//if len(resp.ReturnObj.ZoneList) == 0 {
	//	resp.ReturnObj.ZoneList = append(resp.ReturnObj.ZoneList, ictyun.ZoneEntry{Name: "default", AzDisplayName: "默认"})
	//}
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// RegionResourceSummary 资源池用户资源统计
func (p *Ecs) RegionResourceSummary(
	ctx context.Context,
	req *ictyun.RegionResourceSummaryRequest,
	resp *ictyun.RegionResourceSummaryResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/region/customer-resources").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}
