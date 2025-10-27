package ikubeagent

import "github.com/wangweihong/gotoolbox/pkg/errors"

type (
	InstallStateResponse struct {
		State    *InstallState `json:"state"`
		HostName string        `json:"host_name"`
		HostStat *HostInfo     `json:"host_stat"`
	}
)
type (
	CheckDependencyReq struct {
		Version string `json:"version,omitempty"`
	}

	CheckDependencyResp struct {
	}
)

type (
	JoinClusterRequest struct {
		WorkerConfig *KubernetesWorkerPlaneConfig `json:"worker_config"`
		NodeConfig   *KubernetesNodeConfig        `json:"-"`
	}
)

func (r *JoinClusterRequest) Validate() error {
	if r.WorkerConfig == nil {
		return errors.Errorf("worker config is nil")
	}

	return nil
}

type (
	CpuInfo struct {
		Cores int64 `json:"cores"`
	}

	MemInfo struct {
		Total uint64 `json:"total"`
	}
	HostInfo struct {
		Cpu CpuInfo `json:"cpu"`
		Mem MemInfo `json:"mem"`
	}
)
