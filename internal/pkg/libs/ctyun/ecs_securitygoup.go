package ctyun

import (
	"context"

	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"

	"github.com/wangweihong/eazycloud/apis/ictyun"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
)

type SecurityGroup struct {
	c           *Client
	serviceType string
}

func NewSecurityGroup(c *Client) *SecurityGroup {
	return &SecurityGroup{
		c:           c,
		serviceType: SERVICE_ECS,
	}
}

// SecurityGroupList 查询安全组列表
func (p *SecurityGroup) SecurityGroupList(
	ctx context.Context,
	req *ictyun.SecurityGroupListRequest,
	resp *ictyun.SecurityGroupListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/ecs/vpc/query-security-groups").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// SecurityGroupGet 查询安全组详情
func (p *SecurityGroup) SecurityGroupGet(
	ctx context.Context,
	req *ictyun.SecurityGroupGetRequest,
	resp *ictyun.SecurityGroupGetResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/ecs/vpc/describe-security-group-attribute").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// SecurityGroupCreate 创建安全组
func (p *SecurityGroup) SecurityGroupCreate(
	ctx context.Context,
	req *ictyun.SecurityGroupCreateRequest,
	resp *ictyun.SecurityGroupCreateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/vpc/create-security-group").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// SecurityGroupDelete 删除安全组
func (p *SecurityGroup) SecurityGroupDelete(
	ctx context.Context,
	req *ictyun.SecurityGroupDeleteRequest,
	resp *ictyun.SecurityGroupDeleteResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/vpc/delete-security-group").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// SecurityGroupLeave 将弹性云主机移除出安全组
func (p *SecurityGroup) SecurityGroupLeave(
	ctx context.Context,
	req *ictyun.SecurityGroupLeaveRequest,
	resp *ictyun.SecurityGroupLeaveResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/vpc/leave-security-group").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// SecurityGroupJoin 将弹性云主机加入安全组
func (p *SecurityGroup) SecurityGroupJoin(
	ctx context.Context,
	req *ictyun.SecurityGroupJoinRequest,
	resp *ictyun.SecurityGroupJoinResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/vpc/join-security-group").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// SecurityGroupIngressCreate 安全组创建入口规则
func (p *SecurityGroup) SecurityGroupIngressCreate(
	ctx context.Context,
	req *ictyun.SecurityGroupIngressCreateRequest,
	resp *ictyun.SecurityGroupIngressCreateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/vpc/create-security-group-ingress").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// SecurityGroupIngressUpdate 安全组更改入口规则
func (p *SecurityGroup) SecurityGroupIngressUpdate(
	ctx context.Context,
	req *ictyun.SecurityGroupIngressUpdateRequest,
	resp *ictyun.SecurityGroupIngressUpdateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/vpc/modify-security-group-ingress").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// SecurityGroupIngressCreate 安全组删除入口规则
func (p *SecurityGroup) SecurityGroupIngressDelete(
	ctx context.Context,
	req *ictyun.SecurityGroupIngressDeleteRequest,
	resp *ictyun.SecurityGroupIngressDeleteResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/vpc/revoke-security-group-ingress").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// SecurityGroupEgressCreate 安全组创建出口规则
func (p *SecurityGroup) SecurityGroupEgressCreate(
	ctx context.Context,
	req *ictyun.SecurityGroupEgressCreateRequest,
	resp *ictyun.SecurityGroupEgressCreateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/vpc/create-security-group-egress").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// SecurityGroupEgressUpdate 安全组更改出口规则
func (p *SecurityGroup) SecurityGroupEgressUpdate(
	ctx context.Context,
	req *ictyun.SecurityGroupEgressUpdateRequest,
	resp *ictyun.SecurityGroupEgressUpdateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/vpc/modify-security-group-egress").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// SecurityGroupEgressCreate 安全组删除出口规则
func (p *SecurityGroup) SecurityGroupEgressDelete(
	ctx context.Context,
	req *ictyun.SecurityGroupEgressDeleteRequest,
	resp *ictyun.SecurityGroupEgressDeleteResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ecs/vpc/revoke-security-group-egress").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}
