package iharbor

import (
	"time"

	"github.com/wangweihong/eazycloud/pkg/validator"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

type ProjectSummary struct {
	RepoCount         int          `json:"repo_count"`
	ChartCount        int          `json:"chart_count"`
	ProjectAdminCount int          `json:"project_admin_count"`
	MasterCount       int          `json:"master_count"`
	DeveloperCount    int          `json:"developer_count"`
	GuestCount        int          `json:"guest_count"`
	LimitedGuestCount int          `json:"limited_guest_count"`
	Quota             ProjectQuota `json:"quota"`
}

type ProjectQuota struct {
	Hard Quota `json:"hard"`
	Used Quota `json:"used"`
}

type Quota struct {
	Count   int `json:"count"`
	Storage int `json:"storage"`
}

type Project struct {
	ProjectID          int          `json:"project_id" form:"project_id"`
	OwnerID            int          `json:"owner_id"`
	Name               string       `json:"name"`
	CreationTime       *time.Time   `json:"creation_time,omitempty"`
	UpdateTime         *time.Time   `json:"update_time,omitempty"`
	Deleted            bool         `json:"deleted"`
	OwnerName          string       `json:"owner_name"`
	CurrentUserRoleID  int          `json:"current_user_role_id"`
	CurrentUserRoleIds []int        `json:"current_user_role_ids,omitempty"`
	RepoCount          int          `json:"repo_count"`
	ChartCount         int          `json:"chart_count"`
	Metadata           Metadata     `json:"metadata"`
	CveWhitelist       CveWhiteList `json:"cve_whitelist,omitempty"` //忽略指定漏洞，详情参考https://cloud.tencent.com/developer/article/1533716
}

func setFalseIfEmpty(field *string) {
	if *field == "" {
		*field = "false"
	}
}

func (p *Project) Convert() *Project {
	setFalseIfEmpty(&p.Metadata.PreventVul)
	setFalseIfEmpty(&p.Metadata.AutoScan)
	setFalseIfEmpty(&p.Metadata.EnableContentTrust)
	setFalseIfEmpty(&p.Metadata.Public)
	return p
}

type CveWhiteList struct {
	ID           int        `json:"id"`
	ProjectID    int        `json:"project_id"`
	Items        []CveItem  `json:"items"`
	CreationTime *time.Time `json:"creation_time"`
	UpdateTime   *time.Time `json:"update_time"`
	ExpiresAt    *int64     `json:"expires_at,omitempty"` //过期时间，
}

type CveItem struct {
	ID string `json:"cve_id"`
}

type Metadata struct {
	AutoScan             string `json:"auto_scan,omitempty"`            //true/false 镜像上传后自动扫描
	EnableContentTrust   string `json:"enable_content_trust,omitempty"` //true/false
	PreventVul           string `json:"prevent_vul,omitempty"`          //true/false 阻止潜在漏洞镜像，需要设置Severity严重级别。
	Public               string `json:"public,omitempty"`               //true/false 创建时必须为true/false，不能为空. 公有/私有
	Severity             string `json:"severity,omitempty"`             //none(无)/low(较低)/"medium"(中等)/high(严重)/critical(危急)
	ReuseSysCveWhiteList string `json:"reuse_sys_cve_whitelist"`        //使用系统缺陷白名单
	//	ReuseSysCveAllowList string `json:"reuse_sys_cve_allowlist"`        //使用系统缺陷白名单2.0
}

func (m Metadata) Validate() error {
	var err error
	checkBool := func(field string) {
		if err != nil {
			return
		}
		if field != "true" && field != "false" {
			err = errors.Errorf("bool string field value %v is not true|false", field)
			return
		}
	}
	checkBool(m.AutoScan)
	checkBool(m.EnableContentTrust)
	checkBool(m.PreventVul)
	checkBool(m.Public)
	checkBool(m.ReuseSysCveWhiteList)
	return err
}

type ProjectListRequest struct {
	PagingParam
	Name   string `json:"name" form:"name"`
	Public string `json:"public" form:"public"` //public=""(公有/私有) 和public="0"(私有),public="公有"三种行为不同,不为以上值会报400错误
}

func (r *ProjectListRequest) Validate() error {
	r.PagingParam.Valiate()
	if r.Public != "" && r.Public != "1" && r.Public != "0" {
		return errors.Errorf("invalid public param")
	}
	return nil
}

type ProjectRequest struct {
	Project
}

func (r *ProjectRequest) Validate() error {
	if r.ProjectID == 0 {
		return errors.Errorf("prject id is empty")
	}
	return nil
}

type ProjectUpdateRequest struct {
	Project
}

func (r *ProjectUpdateRequest) Validate() error {
	if r.ProjectID == 0 {
		return errors.Errorf("prject id is empty")
	}
	if err := r.Metadata.Validate(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

type ProjectListResponse struct {
	List       []Project `json:"list"`
	TotalCount int       `json:"total_count"`
}

type ProjectIdentifier struct {
	ProjectID   int    `json:"project_id" form:"project_id"`     //harbor 1.10 use this in repoproject
	ProjectName string `json:"project_name" form:"project_name"` //harbor 2.1.0 use this as repo  project
}

type ProjectGetResponse struct {
	Project Project `json:"project"`
}

type ProjectCreateRequest struct {
	ProjectCreateQuota
	RegistryID  *int                  `json:"registry_id"`
	ProjectName string                `json:"project_name"`
	Metadata    ProjectCreateMetadata `json:"metadata"` //这里不能直接使用Metadata结构()，会报400错误
}

func (r *ProjectCreateRequest) Validate() error {
	if r.ProjectName == "" {
		return errors.Errorf("project name is empty")
	}
	if r.Metadata.Public != "true" && r.Metadata.Public != "false" {
		return errors.Errorf("invalid public param")
	}

	if r.StorageLimit == 0 {
		return errors.Errorf("invalid quota param")
	}
	return nil
}

type ProjectCreateMetadata struct {
	Public string `json:"public"` //true/false  不能不传，也不能为空
}

type ProjectCreateQuota struct {
	CountLimit   int   `json:"count_limit"`   // -1表示无受限,  harbor 2.0 doesn't support count limit any more.
	StorageLimit int64 `json:"storage_limit"` // -1表示无受限, 单位为bytes
}

func NewUnlimitedProjectQuota() ProjectCreateQuota {
	return ProjectCreateQuota{
		CountLimit:   -1,
		StorageLimit: -1,
	}
}

type ProjectQuotaListRequest struct {
	PagingParam
	Reference string `json:"reference"`
}

func (r *ProjectQuotaListRequest) Validate() error {
	r.PagingParam.Valiate()
	r.Reference = "project"
	return nil
}

type ProjectQuotaListResponse struct {
	TotalCount int                `json:"total_count"`
	List       []ProjectQuotaData `json:"list"`
}

type ProjectQuotaUpdateRequest struct {
	ID           int   `json:"id" binding:"required,ne=0"`
	StorageLimit int64 `json:"storage_limit" binding:"required,ne=0"` // -1表示无限制，单位为bytes
}

func (r *ProjectQuotaUpdateRequest) Valiate() error {
	return validator.ValidateAll(r)
}

type ProjectQuotaUpdateResponse struct {
}

type ProjectQuotaData struct {
	ID  int      `json:"id"` // quota id
	Ref struct { //  project info
		ID        int    `json:"id"`
		Name      string `json:"name"`
		OwnerName string `json:"owner_name"`
	} `json:"ref"`
	CreationTime time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
	Hard         struct {
		Storage int64 `json:"storage"`
	} `json:"hard"`
	Used struct {
		Storage int `json:"storage"`
	} `json:"used"`
}
