package iapiserver

import (
	"github.com/wangweihong/eazycloud/apis/imachinery"
)

type (
	ClusterListRequest struct {
		imachinery.PagingParams
		imachinery.ListOptions

		FilterState string `json:"filter_state" form:"filter_state"`
		FilterType  string `json:"filter_type" form:"filter_type"`
	}

	ClusterListResponse struct {
		imachinery.ListRet
		List []*Cluster `json:"list"`
	}
)

type (
	ClusterAddRequest struct {
	}

	ClusterAddResponse struct {
	}
)

type (
	ClusterGetRequest struct {
	}

	ClusterGetResponse struct {
	}
)

type (
	ClusterUpdateRequest struct {
	}

	ClusterUpdateResponse struct {
	}
)

type (
	ClusterDeleteRequest struct {
	}

	ClusterDeleteResponse struct {
	}
)

type (
	ClusterStopRequest struct {
	}

	ClusterStopResponse struct {
	}
)

type (
	ClusterStartRequest struct {
	}

	ClusterStartResponse struct {
	}
)
