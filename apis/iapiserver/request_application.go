package iapiserver

import (
	"strings"

	"github.com/wangweihong/eazycloud/apis/imachinery"
)

type (
	AppStoreListRequest struct {
		imachinery.BasicQueryParam
		AppSource   *int   `json:"app_source"   form:"app_source"`
		FilterState string `json:"filter_state" form:"filter_state"`
	}

	AppStoreListResponse struct {
		imachinery.ListRet
		List []*AppStore `json:"list"`
	}
)

type (
	AppStoreAddRequest struct {
		AppStore
		//local regsitry As store
		// 使用已添加的registry, 作为仓库源
		//UseManagedRegistryProject *RegistryProject `json:"use_managed_registry_project"
	}

	AppStoreAddResponse struct {
		AppStore
	}
)

type (
	AppStoreGetRequest struct {
		AppStore
	}

	AppStoreGetResponse struct {
		AppStore
	}
)

type (
	AppStoreUpdateRequest struct {
		AppStore
	}

	AppStoreUpdateResponse struct {
	}
)

type (
	AppStoreSyncRequest struct {
		AppStore
	}

	AppStoreSyncResponse struct {
	}
)

type (
	AppStoreDeleteRequest struct {
		AppStore
	}

	AppStoreDeleteResponse struct {
	}
)

type (
	AppStoreProbeRequest struct {
		RemoteConfig *AppStoreRemoteConfig `json:"remote_config"`
	}

	AppStoreProbeResponse struct {
		Healthy bool   `json:"healthy"`
		Error   string `json:"error,omitempty"`
	}
)

type (
	ApplicationTemplateValidateRequest struct {
		PackageData []byte `json:"package_data"`
	}

	ApplicationTemplateValidateResponse struct {
	}
)

type (
	ApplicationTemplateEntry struct {
		ApplicationTemplate
		// CategorySet   []*Category                          `json:"category_set,omitempty" description:"分组"` //
		// category id ,overwrite embeded struct field
		LatestVersion *ApplicationTemplateVersion `json:"latest_version"`
		//VersionNum    int                                  `json:"version_num" description:"版本数"`
		//InstanceNum   int                                  `json:"instance_num" description:"实例数"`
		//AppStoreName  string                               `json:"app_store_name"  description:"应用商店名"`
	}

	ApplicationTemplateListRequest struct {
		imachinery.BasicQueryParam
		FilterState string `json:"filter_state"            form:"filter_state"`
		FilterStore string `json:"filter_store"            form:"filter_store"`

		//查询多个分类, 通过","来分割
		Categories       string   `json:"categories" form:"categories"`
		CategoriesShadow []string `json:"-"`
		// 是否在指定所有分类中
		AllCategories bool `json:"all_categories" form:"all_categories"`
	}

	ApplicationTemplateListResponse struct {
		imachinery.ListRet
		List []*ApplicationTemplate `json:"list"`
	}
)

func (r ApplicationTemplateListRequest) PostBind() error {
	if r.Categories != "" {
		r.CategoriesShadow = strings.Split(r.Categories, ",")
	} else {
		r.CategoriesShadow = []string{}
	}
	return nil
}

type (
	ApplicationTemplateAddRequest struct {
		ApplicationTemplate
		PackageData []byte `json:"-"`
	}

	ApplicationTemplateAddResponse struct {
		ApplicationTemplate
	}
)

type (
	ApplicationTemplateGetRequest struct {
		ApplicationTemplate
	}

	ApplicationTemplateGetResponse struct {
		ApplicationTemplate
	}
)

type (
	ApplicationTemplateUpdateRequest struct {
		ApplicationTemplate
	}

	ApplicationTemplateUpdateResponse struct {
		ApplicationTemplate
	}
)

type (
	ApplicationTemplateDeleteRequest struct {
		ApplicationTemplate
	}

	ApplicationTemplateDeleteResponse struct {
		ApplicationTemplate
	}
)

type (
	ApplicationInstanceListRequest struct {
		imachinery.BasicQueryParam

		FilterState   string `json:"filter_state" form:"filter_state"`
		FilterCluster string `json:"filter_cluster" form:"filter_cluster"`
	}

	ApplicationInstanceListResponse struct {
		imachinery.ListRet
		List []*ApplicationInstance `json:"list"`
	}
)

type (
	ApplicationInstanceAddRequest struct {
	}

	ApplicationInstanceAddResponse struct {
	}
)

type (
	ApplicationInstanceGetRequest struct {
		ApplicationInstance
	}

	ApplicationInstanceGetResponse struct {
		ApplicationInstance
	}
)

type (
	ApplicationInstanceWorkloadRequest struct {
		ApplicationInstance
	}

	ApplicationInstanceWorkloadResponse struct {
		ApplicationInstance
	}
)

type (
	ApplicationInstanceUpdateRequest struct {
		ApplicationInstance
	}

	ApplicationInstanceUpdateResponse struct {
	}
)

type (
	ApplicationInstanceDeleteRequest struct {
		ApplicationInstance
	}

	ApplicationInstanceDeleteResponse struct {
	}
)

func TranslateApplicationTemplateState(state_cn string) string {
	switch state_cn {
	case "开发中", "开发":
		return ApplicationTemplateStatusDeveloping
	case "已上架", "上架":
		return ApplicationTemplateStatusOnShelves
	case "已下架", "下架":
		return ApplicationTemplateStatusOffShelves
	}
	return state_cn
}

func TranslateApplicationTemplateVersionState(state_cn string) string {
	switch state_cn {
	case "待提交", "提交":
		return ApplicationTemplateVersionStatusDeveloping
	case "等待审核", "审核", "等待":
		return ApplicationTemplateVersionStatusApproving
	case "通过":
		return ApplicationTemplateVersionStatusApproved
	case "已上架", "上架":
		return ApplicationTemplateVersionStatusOnShelves
	case "已下架", "下架":
		return ApplicationTemplateVersionStatusOffShelves
	case "已拒绝", "拒绝":
		return ApplicationTemplateVersionStatusApprovalReject
	}
	return state_cn
}

type (
	ApplicationCategoryListRequest struct {
		imachinery.BasicQueryParam
		AppStoreID string `json:"app_store_id"  form:"app_store_id"`
	}

	ApplicationCategoryListResponse struct {
		imachinery.ListRet
		List []*ApplicationCategory `json:"list"`
	}
)

type (
	ApplicationCategoryAddRequest struct {
		ApplicationCategory
	}

	ApplicationCategoryAddResponse struct {
		ApplicationCategory
	}
)

type (
	ApplicationCategoryGetRequest struct {
		ApplicationCategory
	}

	ApplicationCategoryGetResponse struct {
		ApplicationCategory
	}
)

type (
	ApplicationCategoryUpdateRequest struct {
		ApplicationCategory
	}

	ApplicationCategoryUpdateResponse struct {
		ApplicationCategory
	}
)

type (
	ApplicationCategoryDeleteRequest struct {
		ApplicationCategory
	}

	ApplicationCategoryDeleteResponse struct {
		ApplicationCategory
	}
)

type (
	ApplicationTemplateVersionListRequest struct {
		imachinery.PagingParams
		imachinery.ListOptions

		FilterState string `json:"filter_state" form:"filter_state"`
	}

	ApplicationTemplateVersionListResponse struct {
		imachinery.ListRet
		List []*ApplicationTemplateVersion `json:"list"`
	}
)

type (
	ApplicationTemplateVersionAddRequest struct {
	}

	ApplicationTemplateVersionAddResponse struct {
	}
)

type (
	ApplicationTemplateVersionGetRequest struct {
	}

	ApplicationTemplateVersionGetResponse struct {
	}
)

type (
	ApplicationTemplateVersionUpdateRequest struct {
	}

	ApplicationTemplateVersionUpdateResponse struct {
	}
)

type (
	ApplicationTemplateVersionDeleteRequest struct {
	}

	ApplicationTemplateVersionDeleteResponse struct {
	}
)
