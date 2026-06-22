package ctyun

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/ictyun"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"
)

type EbsBackup struct {
	c           *Client
	serviceType string
}

// EbsBackupList 查询云硬盘备份
func (p *EbsBackup) EbsBackupList(
	ctx context.Context,
	req *ictyun.EbsBackupListRequest,
	resp *ictyun.EbsBackupListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/ebs-backup/list").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupCreate 创建云硬盘备份
func (p *EbsBackup) EbsBackupCreate(
	ctx context.Context,
	req *ictyun.EbsBackupCreateRequest,
	resp *ictyun.EbsBackupCreateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs-backup/create").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupRestore 恢复云硬盘备份
func (p *EbsBackup) EbsBackupRestore(
	ctx context.Context,
	req *ictyun.EbsBackupRestoreRequest,
	resp *ictyun.EbsBackupRestoreResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs-backup/restore").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupDelete 删除云硬盘备份
func (p *EbsBackup) EbsBackupDelete(
	ctx context.Context,
	req *ictyun.EbsBackupDeleteRequest,
	resp *ictyun.EbsBackupDeleteResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs-backup/delete").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupRepoList 查询云硬盘备份库
func (p *EbsBackup) EbsBackupRepoList(
	ctx context.Context,
	req *ictyun.EbsBackupRepoListRequest,
	resp *ictyun.EbsBackupRepoListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/ebs-backup/repo/list").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupRepoCreate 创建云硬盘备份库
func (p *EbsBackup) EbsBackupRepoCreate(
	ctx context.Context,
	req *ictyun.EbsBackupRepoCreateRequest,
	resp *ictyun.EbsBackupRepoCreateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs-backup/repo/create").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupRepoDelete 删除云硬盘备份库
func (p *EbsBackup) EbsBackupRepoDelete(
	ctx context.Context,
	req *ictyun.EbsBackupRepoDeleteRequest,
	resp *ictyun.EbsBackupRepoDeleteResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs-backup/repo/delete").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupRepoResize 扩容云硬盘备份库
func (p *EbsBackup) EbsBackupRepoResize(
	ctx context.Context,
	req *ictyun.EbsBackupRepoResizeRequest,
	resp *ictyun.EbsBackupRepoResizeResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs-backup/repo/resize").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupRepoRenew 扩容云硬盘备份库
func (p *EbsBackup) EbsBackupRepoRenew(
	ctx context.Context,
	req *ictyun.EbsBackupRepoRenewRequest,
	resp *ictyun.EbsBackupRepoRenewResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs-backup/repo/renew").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)

}

// EbsBackupPolicyList 查询云硬盘备份策略
func (p *EbsBackup) EbsBackupPolicyList(
	ctx context.Context,
	req *ictyun.EbsBackupPolicyListRequest,
	resp *ictyun.EbsBackupPolicyListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/ebs-backup/policy/list").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupPolicyDiskList 查询云硬盘备份策略绑定的云硬盘列表
func (p *EbsBackup) EbsBackupPolicyDiskList(
	ctx context.Context,
	req *ictyun.EbsBackupPolicyDiskListRequest,
	resp *ictyun.EbsBackupPolicyDiskListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/ebs-backup/policy/list-disks").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupPolicyCreate 创建云硬盘备份策略
func (p *EbsBackup) EbsBackupPolicyCreate(
	ctx context.Context,
	req *ictyun.EbsBackupPolicyCreateRequest,
	resp *ictyun.EbsBackupPolicyCreateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs-backup/policy/create").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupPolicyDelete 删除云硬盘备份策略
func (p *EbsBackup) EbsBackupPolicyDelete(
	ctx context.Context,
	req *ictyun.EbsBackupPolicyDeleteRequest,
	resp *ictyun.EbsBackupPolicyDeleteResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {

	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs-backup/policy/delete").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupPolicyExecute 执行云硬盘备份策略
func (p *EbsBackup) EbsBackupPolicyExecute(
	ctx context.Context,
	req *ictyun.EbsBackupPolicyExecuteRequest,
	resp *ictyun.EbsBackupPolicyExecuteResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs-backup/policy/execute").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupPolicyBindDisk 云硬盘备份策略绑定云硬盘
func (p *EbsBackup) EbsBackupPolicyBindDisk(
	ctx context.Context,
	req *ictyun.EbsBackupPolicyBindDiskRequest,
	resp *ictyun.EbsBackupPolicyBindDiskResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs-backup/policy/bind-volumes").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupPolicyUnbindDisk 云硬盘备份策略解绑云硬盘
func (p *EbsBackup) EbsBackupPolicyUnbindDisk(
	ctx context.Context,
	req *ictyun.EbsBackupPolicyUnbindDiskRequest,
	resp *ictyun.EbsBackupPolicyUnbindDiskResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs-backup/policy/unbind-volumes").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupPolicyEnable 启用云硬盘备份策略
func (p *EbsBackup) EbsBackupPolicyEnable(
	ctx context.Context,
	req *ictyun.EbsBackupPolicyEnableRequest,
	resp *ictyun.EbsBackupPolicyEnableResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs-backup/policy/enable").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupPolicyDisable 禁用云硬盘备份策略
func (p *EbsBackup) EbsBackupPolicyDisable(
	ctx context.Context,
	req *ictyun.EbsBackupPolicyDisableRequest,
	resp *ictyun.EbsBackupPolicyDisableResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs-backup/policy/disable").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupPolicyRestore 恢复云硬盘备份策略
func (p *EbsBackup) EbsBackupPolicyRestore(
	ctx context.Context,
	req *ictyun.EbsBackupPolicyRestoreRequest,
	resp *ictyun.EbsBackupPolicyRestoreResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs-backup/policy/restore").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupPolicyUpdate 恢复云硬盘备份策略
func (p *EbsBackup) EbsBackupPolicyUpdate(
	ctx context.Context,
	req *ictyun.EbsBackupPolicyUpdateRequest,
	resp *ictyun.EbsBackupPolicyUpdateResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs-backup/policy/update").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupPolicyBindRepo 云硬盘备份策略绑定备份库
func (p *EbsBackup) EbsBackupPolicyBindRepo(
	ctx context.Context,
	req *ictyun.EbsBackupPolicyBindRepoRequest,
	resp *ictyun.EbsBackupPolicyBindRepoResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs-backup/policy/bind-repo").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupPolicyUnbindRepo 云硬盘备份策略解绑备份库
func (p *EbsBackup) EbsBackupPolicyUnbindRepo(
	ctx context.Context,
	req *ictyun.EbsBackupPolicyUnbindRepoRequest,
	resp *ictyun.EbsBackupPolicyUnbindRepoResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		POST().
		WithPath("/v4/ebs-backup/policy/unbind-repo").
		WithBody("", req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

// EbsBackupPolicyTaskList 查询云硬盘备份策略任务列表
func (p *EbsBackup) EbsBackupPolicyTaskList(
	ctx context.Context,
	req *ictyun.EbsBackupPolicyTaskListRequest,
	resp *ictyun.EbsBackupPolicyTaskListResponse,
	opts ...httpcli.CallOption,
) (*httpcli.HttpResponse, error) {
	r := httpcli.NewHttpRequestBuilder().
		WithEndpoint(p.c.serviceEndpoint(p.serviceType)).
		GET().
		WithPath("/v4/ebs-backup/policy/list-tasks").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := libs.Invoke(ctx, p.c.err, p.c.cc, r, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}
