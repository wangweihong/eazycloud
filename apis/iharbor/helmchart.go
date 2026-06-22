package iharbor

import (
	"time"

	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"
)

// type HelmChartAddRequest struct {
// 		ProjectName string    `json:"project_name" form:"project_name" path:"project"`
// 		ChartData   io.Reader `json:"chart_data"`
// 		ProvData    io.Reader `json:"prov_data"`
// }

// func (r *HelmChartAddRequest) Validate() error {
// 	if r.ProjectName == "" {
// 		return errors.Errorf("project name is empty")
// 	}
// 	if r.ChartData == nil {
// 		return errors.Errorf("chart data is empty")
// 	}
// 	return nil
// }

type HelmChartAddRequest struct {
	ProjectName   string `json:"project_name"`
	OriginBuilder *httpcli.HttpRequestBuilder
}

type HelmChartListRequest struct {
	PagingParam
	ProjectName string `json:"project_name" form:"project_name"`
}

type HelmChartListResponse struct {
	List       []HelmChart `json:"list"`
	TotalCount int         `json:"total_count"`
}

type HelmChart struct {
	Name          string     `json:"name"`
	TotalVersions int        `json:"total_versions"`
	LatestVersion string     `json:"latest_version"`
	Created       *time.Time `json:"created"`
	Updated       *time.Time `json:"updated"`
	Icon          string     `json:"icon"`
	Home          string     `json:"home"`
	Deprecated    bool       `json:"deprecated"`
}

type HelmChartVersion struct {
	Sources     []string `json:"sources"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Maintainers []struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"maintainers"`
	Icon       string           `json:"icon"`
	APIVersion string           `json:"apiVersion"`
	AppVersion string           `json:"appVersion"`
	Urls       []string         `json:"urls"`
	Created    *time.Time       `json:"created"`
	Digest     string           `json:"digest"`
	Labels     []HelmChartLabel `json:"labels"`
}

type HelmChartLabel struct {
}

type HelmChartRequest struct {
	ProjectName string `json:"project_name" form:"project_name" path:"project"`
	ChartName   string `json:"chart_name" form:"chart_name" path:"chart"`
}

type HelmChartVersionRequest struct {
	ProjectName string `json:"project_name" form:"project_name" path:"project"`
	ChartName   string `json:"chart_name" form:"chart_name" path:"chart"`
	Version     string `json:"version" form:"version" path:"version"`
}

func (r *HelmChartVersionRequest) Validate() error {
	if r.ProjectName == "" || r.ChartName == "" || r.Version == "" {
		return errors.Errorf("project|chart|verison is nil")
	}
	return nil
}

type HelmChartVersionListResponse struct {
	List       []HelmChartVersion `json:"list"`
	TotalCount int                `json:"total_count"`
}

type HelmChartManifestData struct {
	Metadata struct {
		Name        string   `json:"name"`
		Sources     []string `json:"sources"`
		Version     string   `json:"version"`
		Description string   `json:"description"`
		Maintainers []struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"maintainers"`
		Icon       string    `json:"icon"`
		APIVersion string    `json:"apiVersion"`
		AppVersion string    `json:"appVersion"`
		Urls       []string  `json:"urls"`
		Created    time.Time `json:"created"`
		Digest     string    `json:"digest"`
	} `json:"metadata"`
	Dependencies interface{} `json:"dependencies"`
	Files        map[string]string
	Security     struct {
		Signature struct {
			Signed   bool   `json:"signed"`
			ProvFile string `json:"prov_file"`
		} `json:"signature"`
	} `json:"security"`
	Labels []HelmChartLabel `json:"labels"`
}

type HelmChartVersionDownloadResponse struct {
	Data string `json:"data"`
}
