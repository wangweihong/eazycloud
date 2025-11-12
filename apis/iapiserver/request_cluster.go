package iapiserver

import (
	"net"

	//"github.com/wangweihong/eazycloud/apis/ikubeagent"
	"github.com/wangweihong/eazycloud/apis/ikubeagent"
	"github.com/wangweihong/eazycloud/apis/imachinery"
	"github.com/wangweihong/eazycloud/pkg/validator"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/sets"
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

type (
	ClusterNode struct {
		IP       string `json:"ip"`
		NodeName string `json:"node_name"`
	}
	ClusterInstallRequest struct {
		imachinery.ObjectMeta
		KubeConfig  *ikubeagent.InstallMasterRequest `json:"config" binding:"required"`
		MasterNodes []*ClusterNode                   `json:"master_nodes" binding:"required"`
		WorkerNodes []*ClusterNode                   `json:"worker_nodes"`
		NodeScale   string                           `json:"node_scale"`
	}

	ClusterInstallResponse struct {
	}
)

func (r *ClusterInstallRequest) Validate() error {
	if err := r.KubeConfig.Validate(); err != nil {
		return errors.WithStack(err)
	}

	nodeIPMap, nodeNameMap := sets.NewString(), sets.NewString()
	for _, v := range r.MasterNodes {
		if err := v.Validate(); err != nil {
			return errors.WithStack(err)
		}
		if nodeIPMap.Has(v.IP) || nodeNameMap.Has(v.NodeName) {
			return errors.Errorf("duplicate ip '%v' or name '%v'", v.IP, v.NodeName)
		}
		nodeIPMap.Insert(v.IP)
		nodeNameMap.Insert(v.NodeName)
	}
	return nil
}

func (r *ClusterNode) Validate() error {
	if r.IP == "" || net.ParseIP(r.IP) == nil {
		return errors.Errorf("ip is invalid")
	}
	if err := validator.IsDNS1123Label(r.NodeName); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
