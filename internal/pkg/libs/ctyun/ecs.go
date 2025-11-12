package ctyun

import (
	"context"
	"strings"

	"github.com/wangweihong/eazycloud/apis/ictyun"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"
)

type Ecs struct {
	c           *Client
	serviceType string
}

// EcsList 查询弹性云主机列表
func (p *Ecs) EcsList(
	ctx context.Context,
	req *ictyun.EcsListRequest,
	resp *ictyun.EcsListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/list-instances").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EcsGetStatusList 查询弹性云主机状态列表
func (p *Ecs) EcsGetStatusList(
	ctx context.Context,
	req *ictyun.EcsGetStatusListRequest,
	resp *ictyun.EcsGetStatusListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {

	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/instance-status-list").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EcsGetStatistics 获取弹性云主机统计
func (p *Ecs) EcsGetStatistics(
	ctx context.Context,
	req *ictyun.EcsGetStatisticsRequest,
	resp *ictyun.EcsGetStatisticsResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {

	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/statistics-instance").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EcsCreate 创建弹性云主机
func (p *Ecs) EcsCreate(
	ctx context.Context,
	req *ictyun.EcsCreateRequest,
	resp *ictyun.EcsCreateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {

	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/create-instance").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EcsUpdate 更新弹性云主机信息
func (p *Ecs) EcsUpdate(
	ctx context.Context,
	req *ictyun.EcsUpdateRequest,
	resp *ictyun.EcsUpdateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {

	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/update-instance").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// EcsReboot 重启弹性云主机
func (p *Ecs) EcsReboot(
	ctx context.Context,
	req *ictyun.EcsRebootRequest,
	resp *ictyun.EcsRebootResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/reboot-instance").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EcsStart 重启弹性云主机
func (p *Ecs) EcsStart(
	ctx context.Context,
	req *ictyun.EcsStartRequest,
	resp *ictyun.EcsStartResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/start-instance").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EcsStop 停止弹性云主机
func (p *Ecs) EcsStop(
	ctx context.Context,
	req *ictyun.EcsStopRequest,
	resp *ictyun.EcsStopResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/stop-instance").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EcsDelete 删除/退订弹性云主机
func (p *Ecs) EcsDelete(
	ctx context.Context,
	req *ictyun.EcsDeleteRequest,
	resp *ictyun.EcsDeleteResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/unsubscribe-instance").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EcsResetPassword 重置弹性云主机密码
func (p *Ecs) EcsResetPassword(
	ctx context.Context,
	req *ictyun.EcsResetPasswordRequest,
	resp *ictyun.EcsResetPasswordResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/reset-password").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EcsBatchCreate 批量创建弹性云主机
func (p *Ecs) EcsBatchCreate(
	ctx context.Context,
	req *ictyun.EcsBatchCreateRequest,
	resp *ictyun.EcsBatchCreateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {

	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/batch-create-instances").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// EcsBatchReboot 批量重启弹性云主机
func (p *Ecs) EcsBatchReboot(
	ctx context.Context,
	req *ictyun.EcsBatchRebootRequest,
	resp *ictyun.EcsBatchRebootResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {

	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/batch-reboot-instances").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EcsBatchStart 批量启动弹性云主机
func (p *Ecs) EcsBatchStart(
	ctx context.Context,
	req *ictyun.EcsBatchStartRequest,
	resp *ictyun.EcsBatchStartResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {

	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/batch-start-instances").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EcsBatchStop 批量停止弹性云主机
func (p *Ecs) EcsBatchStop(
	ctx context.Context,
	req *ictyun.EcsBatchStopRequest,
	resp *ictyun.EcsBatchStopResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {

	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/batch-stop-instances").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EcsBatchDelete 批量删除弹性云主机
func (p *Ecs) EcsBatchDelete(
	ctx context.Context,
	req *ictyun.EcsBatchDeleteRequest,
	resp *ictyun.EcsBatchDeleteResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/batch-delete-instances").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EcsBatchResetPassword 批量删除弹性云主机
func (p *Ecs) EcsBatchResetPassword(
	ctx context.Context,
	req *ictyun.EcsBatchResetPasswordRequest,
	resp *ictyun.EcsBatchResetPasswordResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/batch-reset-password").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EcsGetFlavorFamily 查询弹性云主机规格族列表
func (p *Ecs) EcsGetFlavorFamily(
	ctx context.Context,
	req *ictyun.EcsGetFlavorFamilyRequest,
	resp *ictyun.EcsGetFlavorFamilyResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/flavor-families/list").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// EcsFlavorList 查询弹性云主机规格列表
func (p *Ecs) EcsFlavorList(
	ctx context.Context,
	req *ictyun.EcsFlavorListRequest,
	resp *ictyun.EcsFlavorListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {

	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/flavor/list").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if req.Available != nil || req.Fuzzy != nil {
		var fl []ictyun.EcsFlavorInfo
		for _, v := range resp.ReturnObj.FlavorList {
			// 根据available返回已售罄/未售罄规格
			if req.Available != nil && v.Available != *req.Available {
				continue
			}
			// 支持规格名进行模糊匹配
			if req.Fuzzy != nil && strings.Contains(v.FlavorName, *req.Fuzzy) {
				continue
			}
			fl = append(fl, v)
		}

		resp.ReturnObj.FlavorList = fl
	}

	return httpResp, nil
}

// EcsGetOrderID 根据masterOrderID查询云主机ID
func (p *Ecs) EcsGetOrderID(
	ctx context.Context,
	req *ictyun.EcsGetFlavorFamilyRequest,
	resp *ictyun.EcsGetFlavorFamilyResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/ecs/order/query-uuid").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}
