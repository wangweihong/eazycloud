package iapiserver

import (
	"github.com/wangweihong/eazycloud/apis/imachinery"
	"github.com/wangweihong/eazycloud/apis/iregistry"
)

type (
	RegistryListRequest struct {
		imachinery.BasicQueryParam
		FilterType  string `form:"filter_type"`
		FilterState string `form:"filter_state"`
	}

	RegistryListResponse struct {
		imachinery.ListRet
		List []*iregistry.Registry `json:"list"`
	}
)

type (
	RegistryAddRequest struct {
		iregistry.Registry
	}

	RegistryAddResponse struct {
		iregistry.Registry
	}
)

type (
	RegistryGetRequest struct {
		iregistry.Registry
	}

	RegistryGetResponse struct {
		iregistry.Registry
	}
)

type (
	RegistryUpdateRequest struct {
		iregistry.Registry
	}

	RegistryUpdateResponse struct {
	}
)

type (
	RegistryDeleteRequest struct {
		iregistry.Registry
	}

	RegistryDeleteResponse struct {
	}
)
