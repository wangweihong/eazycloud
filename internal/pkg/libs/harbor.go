package libs

import (
	"context"
	"strconv"

	"github.com/wangweihong/eazycloud/apis/iharbor"
	"github.com/wangweihong/eazycloud/apis/imachinery"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"
	"github.com/wangweihong/gotoolbox/pkg/typeutil"
)

type HarborClient struct {
	endpoint       string
	user, password string
	c              *httpcli.Client
	err            error
}

func NewHarborClient(endpoint string, user, password string, opts ...httpcli.Option) *HarborClient {
	callOpts := DefaultCallOptions()
	if opts != nil {
		callOpts = opts
	}

	c, err := httpcli.NewClient(nil, callOpts...)
	return &HarborClient{
		endpoint: endpoint,
		c:        c,
		err:      err,
		user:     user,
		password: password,
	}
}

func (c *HarborClient) ArtifactList(ctx context.Context, req *iharbor.ArtifactListRequest, resp *iharbor.ArtifactListResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	//https://10.30.100.176/api/v2.0/projects/library/repositories/prepare/artifacts?with_tag=true&with_scan_overview=true&with_label=true&page_size=15&page=1
	//注意如果repo路径带有/, 必须要转义再转义
	//如镜像k8s.gcr.io/k8s.gcr.io/ingress-nginx/controller
	//https://10.30.100.179/api/v2.0/projects/k8s.gcr.io/repositories/k8s.gcr.io%252Fingress-nginx%252Fcontroller/artifacts?with_tag=true&with_scan_overview=true&with_label=true&page_size=15&page=1

	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/v2.0/projects/{project}/repositories/{repository}/artifacts").
		AddPathParamByObject(req).
		AddQueryParamByObject(req).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	if err != nil {
		return httpResp, errors.WithStack(err)
	}

	if resp != nil {
		total, _ := strconv.Atoi(httpResp.GetHeader("X-Total-Count"))
		resp.TotalCount = total
	}
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ArtifactGet(ctx context.Context, req *iharbor.ArtifactGetRequest, resp *iharbor.Artifact, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/v2.0/projects/{project}/repositories/{repository}/artifacts/{digest}").
		AddPathParamByObject(req).
		AddQueryParamByObject(req).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ArtifactScan(ctx context.Context, req *iharbor.ArtifactRequest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		POST().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/v2.0/projects/{project}/repositories/{repository}/artifacts/{digest}/scan").
		AddPathParamByObject(req).
		AddQueryParamByObject(req).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ArtifactScanLog(ctx context.Context, req *iharbor.ArtifactScanLogRequest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/v2.0/projects/{project}/repositories/{repository}/artifacts/{digest}/scan/{scan_id}/log").
		AddPathParamByObject(req).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	if err != nil {
		return httpResp, errors.WithStack(err)
	}

	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ArtifactBuildHistory(ctx context.Context, req *iharbor.ArtifactRequest, resp *iharbor.ArtifactBuildHistory, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/v2.0/projects/{project}/repositories/{repository}/artifacts/{digest}/additions/build_history").
		AddPathParamByObject(req).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	if err != nil {
		return httpResp, errors.WithStack(err)
	}

	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ArtifactVulnerabilities(ctx context.Context, req *iharbor.ArtifactRequest, resp *iharbor.ArtifactScanResults, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/v2.0/projects/{project}/repositories/{repository}/artifacts/{digest}/additions/vulnerabilities").
		AddPathParamByObject(req).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	if err != nil {
		return httpResp, errors.WithStack(err)
	}

	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ArtifactDelete(ctx context.Context, req *iharbor.ArtifactRequest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		DELETE().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/v2.0/projects/{project}/repositories/{repository}/artifacts/{digest}").
		AddPathParamByObject(req).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ArtifactTagList(ctx context.Context, req *iharbor.ArtifactTagListRequest, resp *iharbor.ArtifactTagListResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/v2.0/projects/{project}/repositories/{repository}/artifacts/{digest}/tags").
		AddPathParamByObject(req).
		AddQueryParamByObject(req).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	if err != nil {
		return httpResp, errors.WithStack(err)
	}

	if resp != nil {
		total, _ := strconv.Atoi(httpResp.GetHeader("X-Total-Count"))
		resp.TotalCount = total
	}

	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ArtifactTagDelete(ctx context.Context, req *iharbor.ArtifactTagRequest, resp *iharbor.ArtifactTagListResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		DELETE().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/v2.0/projects/{project}/repositories/{repository}/artifacts/{digest}/tags/{name}").
		AddPathParamByObject(req).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ArtifactTagCreate(ctx context.Context, req *iharbor.ArtifactTagRequest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		POST().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/v2.0/projects/{project}/repositories/{repository}/artifacts/{digest}/tags").
		AddPathParamByObject(req).
		WithBody("", req).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

/* -------------------------------- config -------------------------- */

func (c *HarborClient) Config(ctx context.Context, req *imachinery.Empty, resp *iharbor.ConfigResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/v2.0/configurations").
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

/* -------------------------------- distribute -------------------------- */

func (c *HarborClient) PreheatInstanceProviderList(ctx context.Context, req *iharbor.PreheatInstanceProviderListRequest, resp *iharbor.PreheatInstanceProviderListResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/v2.0/p2p/preheat/providers").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) PreheatInstanceProviderInstances(ctx context.Context, req *imachinery.Empty, resp *iharbor.PreheatInstanceProviderListResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/v2.0/p2p/preheat/instances").
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) PreheatInstanceProviderPing(ctx context.Context, req *iharbor.PreheatInstance, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		POST().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/v2.0/p2p/preheat/instances/ping").
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) PreheatInstanceProviderCreate(ctx context.Context, req *iharbor.PreheatInstanceRequest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		POST().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/v2.0/p2p/preheat/instances").
		WithBody("", req).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) PreheatInstanceProviderUpdate(ctx context.Context, req *iharbor.PreheatInstanceRequest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		PUT().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/v2.0/p2p/preheat/instances").
		WithBody("", req).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) PreheatInstanceProviderDelete(ctx context.Context, req *iharbor.PreheatInstanceRequest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		DELETE().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/v2.0/p2p/preheat/instances").
		WithBody("", req).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

/* -------------------------------- healthy -------------------------- */

func (c *HarborClient) Health(ctx context.Context, req *imachinery.Empty, resp *iharbor.SystemHealth, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/v2.0/p2p/health").
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

/* -------------------------------- helm chart -------------------------- */

func (c *HarborClient) HelmChartAdd(ctx context.Context, req *iharbor.HelmChartAddRequest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {

	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		PUT().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/chartrepo/{project}/charts").
		AddPathParam("project", req.ProjectName).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) HelmChartList(ctx context.Context, req *iharbor.HelmChartListRequest, resp *iharbor.HelmChartListResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/chartrepo/{project}/charts").
		AddPathParam("project", req.ProjectName).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) HelmChartDelete(ctx context.Context, req *iharbor.HelmChartRequest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		DELETE().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/chartrepo/{project}/charts/{chart}").
		AddPathParam("project", req.ProjectName).
		AddPathParam("chart", req.ChartName).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) HelmChartVersionList(ctx context.Context, req *iharbor.HelmChartRequest, resp *iharbor.HelmChartVersionListResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/chartrepo/{project}/charts/{chart}").
		AddPathParam("project", req.ProjectName).
		AddPathParam("chart", req.ChartName).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) HelmChartVersionDelete(ctx context.Context, req *iharbor.HelmChartVersionRequest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		DELETE().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/chartrepo/{project}/charts/{chart}/{version}").
		AddPathParam("project", req.ProjectName).
		AddPathParam("chart", req.ChartName).
		AddPathParam("version", req.Version).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) HelmChartVersionManifests(ctx context.Context, req *iharbor.HelmChartVersionRequest, resp *iharbor.HelmChartVersionListResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/chartrepo/{project}/charts/{chart}/{version}").
		AddPathParam("project", req.ProjectName).
		AddPathParam("chart", req.ChartName).
		AddPathParam("version", req.Version).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

/*
buf := make([]byte, 32*1024) // 32KB缓冲区
_, err = io.CopyBuffer(c.Writer, getResp.Body, buf)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "传输中断"})
	}
*/
func (c *HarborClient) HelmChartVersionDownload(ctx context.Context, req *iharbor.HelmChartVersionRequest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/chartrepo/{project}/charts/{chart}/{version}").
		AddPathParam("project", req.ProjectName).
		AddPathParam("chart", req.ChartName).
		AddPathParam("version", req.Version).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) AuditLog(ctx context.Context, req *iharbor.LogListRequest, resp *iharbor.LogListResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	logs := make([]iharbor.Log, 0)
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/v2.0/audit-logs").
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, logs, opts...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp = &iharbor.LogListResponse{
		List:       logs,
		TotalCount: typeutil.MustAtoi(httpResp.GetHeader("X-Total-Count")),
	}

	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ProjectList(ctx context.Context, req *iharbor.ProjectListRequest, resp *iharbor.ProjectListResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	projects := make([]iharbor.Project, 0)
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/projects").
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, projects, opts...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp = &iharbor.ProjectListResponse{
		List:       projects,
		TotalCount: typeutil.MustAtoi(httpResp.GetHeader("X-Total-Count")),
	}
	return httpResp, nil
}

func (c *HarborClient) ProjectCreate(ctx context.Context, req *iharbor.ProjectCreateRequest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		POST().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/projects").
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ProjectUpdate(ctx context.Context, req *iharbor.ProjectRequest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		PUT().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/projects/{project}").
		AddPathParam("project", strconv.Itoa(req.ProjectID)).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ProjectDelete(ctx context.Context, req *iharbor.ProjectRequest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		DELETE().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/projects/{project}").
		AddPathParam("project", strconv.Itoa(req.ProjectID)).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ProjectGet(ctx context.Context, req *iharbor.ProjectRequest, resp *iharbor.Project, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/projects/{project}").
		AddPathParam("project", strconv.Itoa(req.ProjectID)).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ProjectSummary(ctx context.Context, req *iharbor.ProjectRequest, resp *iharbor.ProjectSummary, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/projects/{project}/summary").
		AddPathParam("project", strconv.Itoa(req.ProjectID)).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ProjectSearch(ctx context.Context, req *iharbor.ProjectListRequest, resp *iharbor.ProjectListResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/search").
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ProjectQuotaList(ctx context.Context, req *iharbor.ProjectQuotaListRequest, resp *iharbor.ProjectQuotaListResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	quotas := make([]iharbor.ProjectQuotaData, 0)
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/quotas").
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, quotas, opts...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp = &iharbor.ProjectQuotaListResponse{
		List:       quotas,
		TotalCount: typeutil.MustAtoi(httpResp.GetHeader("X-Total-Count")),
	}
	return httpResp, nil
}

func (c *HarborClient) ProjectQuotaUpdate(ctx context.Context, req *iharbor.ProjectQuotaUpdateRequest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	data := struct {
		Hard struct {
			Storage int64 `json:"storage"`
		} `json:"hard"`
	}{
		Hard: struct {
			Storage int64 `json:"storage"`
		}{
			Storage: req.StorageLimit,
		},
	}
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		PUT().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/quotas").
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, data, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) RemoteRegistryPing(ctx context.Context, req *iharbor.RegistryPingRquest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		POST().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/registries/ping").
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) RemoteRegistryInfo(ctx context.Context, req *iharbor.Registry, resp *iharbor.RegistryInfo, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/registries/{registry}/info").
		AddPathParam("registry", strconv.Itoa(req.ID)).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) RemoteRegistryAdd(ctx context.Context, req *iharbor.Registry, resp *iharbor.RegistryInfo, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		POST().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/registries").
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) RemoteRegistryDelete(ctx context.Context, req *iharbor.Registry, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		DELETE().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/registries/{registry}").
		AddPathParam("registry", strconv.Itoa(req.ID)).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) RemoteRegistryUpdate(ctx context.Context, req *iharbor.Registry, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		PUT().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/registries/{registry}").
		AddPathParam("registry", strconv.Itoa(req.ID)).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ReplicationAdapterNames(ctx context.Context, req *imachinery.Empty, resp *iharbor.AdapterInfoResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/replication/adapterinfos").
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ReplicationPolicyList(ctx context.Context, req *iharbor.PolicyListRequest, resp *iharbor.PolicyListResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	quotas := make([]iharbor.Policy, 0)
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/replication/policies").
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, quotas, opts...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp = &iharbor.PolicyListResponse{
		List:       quotas,
		TotalCount: typeutil.MustAtoi(httpResp.GetHeader("X-Total-Count")),
	}
	return httpResp, nil
}

func (c *HarborClient) ReplicationPolicyGet(ctx context.Context, req *iharbor.Policy, resp *iharbor.Policy, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/replication/policies/{policy}").
		AddPathParam("policy", strconv.FormatInt(req.ID, 10)).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ReplicationPolicyAdd(ctx context.Context, req *iharbor.Policy, resp *iharbor.Policy, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		POST().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/replication/policies").
		WithBody("", req).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ReplicationPolicyUpdate(ctx context.Context, req *iharbor.Policy, resp *iharbor.Policy, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		PUT().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/replication/policies").
		AddPathParam("policy", strconv.FormatInt(req.ID, 10)).
		WithBody("", req).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ReplicationPolicyDelete(ctx context.Context, req *iharbor.Policy, resp *iharbor.Policy, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		DELETE().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/replication/policies").
		AddPathParam("policy", strconv.FormatInt(req.ID, 10)).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ReplicationPolicyExcute(ctx context.Context, req *iharbor.PolicyExecuteIdentifier, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		POST().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/replication/executions").
		WithBody("", req).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) ReplicationPolicyExecuteResultList(ctx context.Context, req *iharbor.PolicyExecuteListRequest, resp *iharbor.PolicyExecuteListResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	list := make([]iharbor.PolicyExecuteResult, 0)
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		AddQueryParamByObject(req).
		WithPath("/api/2.0/replication/executions").
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, list, opts...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp = &iharbor.PolicyExecuteListResponse{
		List:       list,
		TotalCount: typeutil.MustAtoi(httpResp.GetHeader("X-Total-Count")),
	}
	return httpResp, nil
}

func (c *HarborClient) RepositoryList(ctx context.Context, req *iharbor.RepositoryListRequest, resp *iharbor.RepositoryListResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	list := make([]iharbor.Repository, 0)
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		AddQueryParamByObject(req).
		WithPath("/api/2.0/projects/{project}/repositories").
		AddPathParam("project", req.ProjectName).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, list, opts...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp = &iharbor.RepositoryListResponse{
		List:       list,
		TotalCount: typeutil.MustAtoi64(httpResp.GetHeader("X-Total-Count")),
	}
	return httpResp, nil
}

func (c *HarborClient) RepositoryDelete(ctx context.Context, req *iharbor.RepositoryDeleteRequest, resp *iharbor.RepositoryDeleteResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		DELETE().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/projects/{project}/repositories/{repo}").
		AddPathParam("project", req.ProjectName).
		AddPathParam("repo", req.RepoName).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) RepositoryGet(ctx context.Context, req *iharbor.RepositoryGetRequest, resp *iharbor.RepositoryGetResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/projects/{project}/repositories/{repo}").
		AddPathParam("project", req.ProjectName).
		AddPathParam("repo", req.RepositoryName).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) RepositoryTagList(ctx context.Context, req *iharbor.TagListRequest, resp *iharbor.TagListResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	artifacts := make([]iharbor.Artifact, 0)
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		AddQueryParamByObject(req).
		WithPath("/api/2.0/projects/{project}/repositories/{repo}/artifacts").
		AddPathParam("project", req.ProjectName).
		AddPathParam("repo", req.ProjectName).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, artifacts, opts...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tags := make([]*iharbor.Tag, 0)
	for _, v := range artifacts {
		tags = append(tags, iharbor.ConvertArtifactToTagList(&v)...)
	}
	resp = &iharbor.TagListResponse{
		List:       tags,
		TotalCount: typeutil.MustAtoi64(httpResp.GetHeader("X-Total-Count")),
	}
	return httpResp, nil
}

// func (c *HarborClient) RepositoryTagGet(ctx context.Context, req *iharbor.TagRequest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
// 	var artifact iharbor.Artifact
// 	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
// 		GET().
// 		AddBasicAuthHeaderParam(c.user, c.password).
// 		WithPath("/api/2.0/projects/{project}/repositories/{repo}/artifacts/{artifact}").
// 		AddPathParam("project", req.ProjectName).
// 		AddPathParam("repo", req.RepoName).
// 		AddPathParam("artifact", req.ArtifactDigest).
// 		AddQueryParamByObject(req).
// 		Build()
// 	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, artifact, opts...)
// 	if err != nil {
// 		return nil, errors.WithStack(err)
// 	}

// 	return httpResp, errors.WithStack(err)
// }

func (c *HarborClient) RepositoryTagDelete(ctx context.Context, req *iharbor.TagRequest, resp *imachinery.Empty, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		DELETE().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/projects/{project}/repositories/{repo}/artifacts/{artifact}/tags/{tag}").
		AddPathParam("project", req.ProjectName).
		AddPathParam("repo", req.RepoName).
		AddPathParam("artifact", req.ArtifactDigest).
		AddPathParam("tag", req.TagName).
		AddQueryParamByObject(req).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) RepositorySearch(ctx context.Context, req *iharbor.RepositoryListRequest, resp *iharbor.RepositoryListSearchResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/search").
		AddQueryParamByObject(req).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func (c *HarborClient) Search(ctx context.Context, req *iharbor.SearchRequest, resp *iharbor.SearchResponse, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	httpReq := httpcli.NewHttpRequestBuilder().WithEndpoint(c.endpoint).
		GET().
		AddBasicAuthHeaderParam(c.user, c.password).
		WithPath("/api/2.0/search").
		AddQueryParam("q", req.Query).
		Build()
	httpResp, err := invoke(ctx, c.err, c.c, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}
