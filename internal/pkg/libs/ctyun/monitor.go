package ctyun

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/ictyun"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"
)

type Monitor struct {
	c           *Client
	serviceType string
}

// MonitorEcsList 根据筛选条件查询资源池下云主机的列表。 返回对应的设备ID
func (p *Monitor) MonitorEcsList(
	ctx context.Context,
	req *ictyun.MonitorEcsListRequest,
	resp *ictyun.MonitorEcsListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {

	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/monitor/query-ecs-list").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// MonitorDiskList  根据筛选条件查询资源池下云硬盘的列表。 返回对应的设备ID
func (p *Monitor) MonitorDiskList(
	ctx context.Context,
	req *ictyun.MonitorDiskListRequest,
	resp *ictyun.MonitorDiskListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {

	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/monitor/query-evs-list").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// MonitorItemList  监控对象：查询各设备类型支持的监控项列表
func (p *Monitor) MonitorItemList(
	ctx context.Context,
	req *ictyun.MonitorItemListRequest,
	resp *ictyun.MonitorItemListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {

	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/monitor/query-monitor-items").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// MonitorHistoryList  监控对象：查询各设备类型支持的监控项列表
func (p *Monitor) MonitorHistoryList(
	ctx context.Context,
	req *ictyun.MonitorHistoryListRequest,
	resp *ictyun.MonitorHistoryListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {

	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4.1/monitor/query-vm-historymetricdata").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}
