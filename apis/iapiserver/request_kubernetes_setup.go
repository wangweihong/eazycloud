package iapiserver

import (
	"github.com/wangweihong/eazycloud/apis/ikubeagent"
	"github.com/wangweihong/eazycloud/apis/imachinery"
)

type (
	HostScanRequest struct {
		imachinery.PagingParams
		Begin string `json:"begin" binding:"required,ip"`
		End   string `json:"end" binding:"required,ip"`
	}

	HostScanResponse struct {
		HostsInfo  []*HostScanInfo `json:"hosts_info"`
		TotalCount int             `json:"total_count"`
	}

	HostScanInfo struct {
		IP string `json:"ip"`
		//	HostName string                           `json:"host_name"`
		Version  string               `json:"version"`
		HostInfo *ikubeagent.HostInfo `json:"host_info,omitempty"`
	}
)

type (
	VipUsedTestRequest struct {
		Vip             string   `json:"vip" binding:"required,ip"`
		ExternalPort    int      `json:"external_port" binding:"required,port"`
		NodeList        []string `json:"node_list" binding:"required,ips"`
		LocalAdvisePort int      `json:"local_advise_port" binding:"required,port"`
	}
)

type (
	SetupClusterRequest struct {
		MasterList []*ikubeagent.KubernetesNodeConfig `json:"master_list"`
		NodeList   []*ikubeagent.KubernetesNodeConfig `json:"node_list"`
		Name       string                             `json:"name" binding:"required"`
		NodeScale  string                             `json:"node_scale"`

		Config *ikubeagent.InstallMasterRequest

		ikubeagent.InstallMasterRequest
		//	UseManagedPegistry
	}

	SetupClusterResponse struct {
		ID string `json:"id"`
	}
)

type (
	GetClusterNodeScaleRequirementResponse struct {
		Requirements []ClusterNodeScaleRequirement `json:"requirements"`
	}

	ClusterNodeScaleRequirement struct {
		Name             string `json:"name"  description:"要求名"`
		MinNodeNum       int    `json:"min_node_num" description:"最小节点数"`
		MaxNodeNum       int    `json:"max_node_num" description:"最大节点数"`
		LowestCpuCores   int    `json:"lowest_cpu_cores" description:"最低cpu核数"`
		LowestMemoryUnit int    `json:"lowest_memory_unit" description:"最低内存单元"`
	}
)

type (
	KubeMasterNodeFitClusterScaleCheckRequest struct {
		ClusterNodeScale int                                `json:"cluster_node_scale" binding:"required"` //集群工作节点规模，部署前检测
		MasterList       []*ikubeagent.KubernetesNodeConfig `json:"master_list" binding:"required"`
	}
)
