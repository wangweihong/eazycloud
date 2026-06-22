package iregistry

import (
	"time"

	"github.com/wangweihong/eazycloud/apis/imachinery"
)

type ListRegistryParam struct {
	imachinery.PagingParams
}

type ListRegistryResult struct {
}

type DeleteRegistryParam struct {
	ID string `json:"id" query:"id" binding:"required"`
}

type GetRegistryParam struct {
	ID string `json:"id" query:"id" binding:"required"`
}

type SearchRegistryParam struct {
	ID    string `json:"id"    query:"id"    binding:"required"`
	Query string `json:"query" query:"query" binding:"required"`
}

type (
	ProjectListRequest struct {
		imachinery.PagingParams
		Name   string `json:"name"   form:"name"`
		Public string `json:"public" form:"public"` //public=""(公有/私有) 和public="0"(私有),public="公有"三种行为不同,不为以上值会报400错误
	}

	ProjectListResponse struct {
		List       []Project `json:"list"`
		TotalCount int64     `json:"total_count"`
	}
	Project struct {
		ProjectID          int          `json:"project_id"                      form:"project_id"`
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
	//	Metadata           Metadata     `json:"metadata"`
	//	CveWhitelist       CveWhiteList `json:"cve_whitelist,omitempty"` //忽略指定漏洞，详情参考https://cloud.tencent.com/developer/article/1533716
	}
)
