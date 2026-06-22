package iharbor

import (
	"net/url"
	"time"

	"github.com/wangweihong/gotoolbox/pkg/errors"
)

type PagingParam struct {
	PageNum  int `json:"page" form:"page"`           //不能为0，会报400错误
	PageSize int `json:"page_size" form:"page_size"` //默认10条，不能为0, 会报400错误
}

func (p *PagingParam) Valiate() error {
	if p.PageNum <= 0 {
		p.PageNum = 1
	}

	// project实际测试中每页最多返回500条,多于500条在下页返回
	if p.PageSize <= 0 {
		p.PageSize = 500
	}
	return nil
}

type (
	Artifact struct {
		Digest            string                             `json:"digest"              form:"digest"`
		ExtraAttrs        ArtifactExtraAttr                  `json:"extra_attrs"`
		Icon              string                             `json:"icon"`
		ID                int                                `json:"id"`
		Labels            any                                `json:"labels"`
		ManifestMediaType string                             `json:"manifest_media_type"`
		MediaType         string                             `json:"media_type"`
		ProjectID         int                                `json:"project_id"`
		PullTime          *time.Time                         `json:"pull_time"`
		PushTime          *time.Time                         `json:"push_time"`
		References        any                                `json:"references"`
		RepositoryID      int                                `json:"repository_id"`
		ScanOverview      map[string]ArtifactTagScanOverview `json:"scan_overview"` // 为NIL表示未扫描
		Size              int                                `json:"size"`
		Tags              []ArtifactTag                      `json:"tags"`
		Type              string                             `json:"type"`
	}

	ArtifactAdditionLink struct {
		BuildHistory struct {
			Absolute bool   `json:"absolute"`
			Href     string `json:"href"`
		} `json:"build_history"`
		Vulnerabilities struct {
			Absolute bool   `json:"absolute"`
			Href     string `json:"href"`
		} `json:"vulnerabilities"`
	}

	ArtifactExtraAttr struct {
		Architecture string     `json:"architecture"`
		Author       any        `json:"author"`
		Created      *time.Time `json:"created"`
		Os           string     `json:"os"`
	}

	ArtifactTagScanOverview struct {
		ReportID   string                 `json:"report_id"`
		ScanStatus string                 `json:"scan_status"` //Success/Running/Pending/Error 扫描状态（结束/运行）
		Severity   string                 `json:"severity"`    //None
		Duration   int                    `json:"duration"`    //扫描时间
		Summary    ArtifactTagScanSummary `json:"summary"`     //扫描结束后才有数据
		StartTime  *time.Time             `json:"start_time"`
		EndTime    *time.Time             `json:"end_time"`
		Scanner    *ArtifactTagScanner    `json:"scanner"` //扫描结束后才有数据
	}

	ArtifactTagScanner struct {
		Name    string `json:"name"`
		Vendor  string `json:"vendor"`
		Version string `json:"version"`
	}

	ArtifactTagScanSummary struct {
		Total   int                 `json:"total"`   //总的漏洞数？
		Fixable int                 `json:"fixable"` //可修复漏洞数
		Summary ArtifactScanSummary `json:"summary"`
	}

	ArtifactTagSignature struct {
	}

	// 扫描出来的漏洞总结
	ArtifactScanSummary struct {
		Low        int `json:"Low"`        //严重级别为低
		Medium     int `json:"Medium"`     //严重级别为中等
		Negligible int `json:"Negligible"` //严重界别为忽略不计
		High       int `json:"High"`       //严重界别为高
		Critical   int `json:"Critical"`   //严重级别为严重
		Unknown    int `json:"Unknown"`    //严重级别未知
	}

	ArtifactTag struct {
		ArtifactID   int        `json:"artifact_id"`
		ID           int        `json:"id"`
		Immutable    bool       `json:"immutable"`
		Name         string     `json:"name"           form:"name"`
		PullTime     *time.Time `json:"pull_time"`
		PushTime     *time.Time `json:"push_time"`
		RepositoryID int        `json:"repository_id"`
		Signed       bool       `json:"signed"`
	}

	ArtifactLabel struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Description  string `json:"description"`
		Color        string `json:"color"`
		Scope        string `json:"scope"`
		ProjectID    int    `json:"project_id"`
		CreationTime string `json:"creation_time"`
		UpdateTime   string `json:"update_time"`
		Deleted      bool   `json:"deleted"`
	}

	ArtifactListRequest struct {
		PagingParam
		ProjectName      string `path:"project"`
		RepoName         string `path:"repository"`
		WithTag          bool   `form:"with_tag"`
		WithScanOverview bool   `form:"with_scan_overview"`
		WithLabel        bool   `form:"with_label"`
	}

	ArtifactListResponse struct {
		List       []*Artifact `json:"list"`
		TotalCount int         `json:"total_count"`
	}
)

func (r *ArtifactListRequest) Validate() error {
	if r.ProjectName == "" || r.RepoName == "" {
		return errors.Errorf("project or repository is empty")
	}
	r.PagingParam.Valiate()

	r.WithLabel = true
	r.WithScanOverview = true
	r.WithLabel = true
	return nil
}

type ArtifactGetRequest struct {
	ProjectName string `path:"project"`
	RepoName    string `path:"repository"`
	Digest      string `path:"digest"`

	WithTag          bool `form:"with_tag"`
	WithScanOverview bool `form:"with_scan_overview"`
	WithLabel        bool `form:"with_label"`
}

func (r *ArtifactGetRequest) Validate() error {
	if r.ProjectName == "" || r.RepoName == "" {
		return errors.Errorf("project or repository is empty")
	}
	r.RepoName = url.PathEscape(url.PathEscape(r.RepoName))

	r.WithLabel = true
	r.WithScanOverview = true
	r.WithLabel = true
	return nil
}

type ArtifactRequest struct {
	ProjectName string `path:"project"`
	RepoName    string `path:"repository"`
	Digest      string `path:"digest"`
}

func (r *ArtifactRequest) Validate() error {
	if r.ProjectName == "" || r.RepoName == "" {
		return errors.Errorf("project or repository is empty")
	}

	if r.Digest == "" {
		return errors.Errorf("digest is empty")
	}
	r.RepoName = url.PathEscape(url.PathEscape(r.RepoName))
	return nil
}

type ArtifactTagListRequest struct {
	ArtifactRequest
	PagingParam
}

func (r *ArtifactTagListRequest) Validate() error {
	r.PagingParam.Valiate()
	return r.ArtifactRequest.Validate()
}

type ArtifactTagListResponse struct {
	List       []*ArtifactTag `json:"list"`
	TotalCount int            `json:"total_count"`
}

type ArtifactTagRequest struct {
	ArtifactRequest
	Tag string `json:"name" path:"name"`
}

func (r *ArtifactTagRequest) Validate() error {
	if r.Tag == "" {
		return errors.Errorf("tag is empty")
	}
	return r.ArtifactRequest.Validate()
}

type ArtifactScanLogRequest struct {
	ArtifactRequest
	ReportID string `path:"report_id"`
}
type ArtifactScanLogResponse struct {
	Log string `json:"log"`
}

type ArtifactBuildHistory struct {
	Historys []*BuildHistory `json:"historys"`
}

type BuildHistory struct {
	Created    *time.Time `json:"created"`
	CreatedBy  string     `json:"created_by"`
	EmptyLayer bool       `json:"empty_layer,omitempty"`
}

type ArtifactScanResults struct {
	Results map[string]ArtifactScanResult `json:"results"`
}

type ArtifactScanResult struct {
	GeneratedAt     time.Time               `json:"generated_at"`
	Scanner         ScannerBase             `json:"scanner"`
	Severity        string                  `json:"severity"` // "Medium"
	Vulnerabilities []ArtifactVulnerability `json:"vulnerabilities"`
}

type ArtifactVulnerability struct {
	ID           string `json:"id"`
	Severity     string `json:"severity"`
	Package      string `json:"package"`
	Version      string `json:"version"`
	Description  string `json:"description"`
	FixedVersion string `json:"fix_version"`
}
