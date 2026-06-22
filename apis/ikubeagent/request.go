package ikubeagent

import (
	"os"

	"github.com/wangweihong/eazycloud/pkg/validator"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/validation"
)

type (
	InstallStateRequest struct {
	}

	InstallStateResponse struct {
		State    *InstallState `json:"state"`
		HostName string        `json:"host_name"`
		HostStat *HostInfo     `json:"host_stat"`
	}
)
type (
	CheckDependencyRequest struct {
		Version string `json:"version" form:"version"`
	}

	CheckDependencyResponse struct {
	}
)

type (
	InstallMasterRequest struct {
		RegistryConfig    *Registry `json:"registry_config" binding:"required,dive"`
		KubernetesVersion string    `json:"kubernetes_version"`
		// master配置, 传递了则认为是master
		ControlPlaneConfig *KubernetesControlPlaneConfig `json:"control_plane_config" binding:"required"`
		NodeConfig         *KubernetesNodeConfig         `json:"node_config" binding:"required"`
		// etcd快照策略
		DataBaseConfig *EtcdDatabaseConfig `json:"database_config"`
		// control plane也作为工作节点, 一体机环境设置
		MasterAlsoWorker bool `json:"master_also_worker"`
	}
)

func (r *InstallMasterRequest) ToKubernetesInstallConfig() *KubernetesDeployConfig {
	return &KubernetesDeployConfig{
		RegistryConfig:     r.RegistryConfig,
		KubernetesVersion:  r.KubernetesVersion,
		ControlPlaneConfig: r.ControlPlaneConfig,
		NodeConfig:         r.NodeConfig,
		DataBaseConfig:     r.DataBaseConfig,
		MasterAlsoWorker:   r.MasterAlsoWorker,
	}
}

func (r *InstallMasterRequest) Validate() error {
	if r.NodeConfig.NodeName == "" {
		r.NodeConfig.NodeName, _ = os.Hostname()
	}

	if err := validator.IsDNS1123Label(r.NodeConfig.NodeName); err != nil {
		return errors.WithStack(err)
	}

	if err := r.RegistryConfig.Validate(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

type (
	JoinClusterRequest struct {
		RegistryConfig *Registry               `json:"registry_config"`
		WorkerConfig   *KubernetesWorkerConfig `json:"worker_config" binding:"required"`
		HAConfig       *KubernetesHAConfig     `json:"ha_config"`
		Sync           bool                    `json:"sync"` //是否同步安装
		NodeConfig     *KubernetesNodeConfig   `json:"node_config" `
	}
)

func (r *JoinClusterRequest) ToKubernetesInstallConfig() *KubernetesDeployConfig {
	return &KubernetesDeployConfig{
		RegistryConfig: r.RegistryConfig,
		NodeConfig:     r.NodeConfig,
		WorkerConfig:   r.WorkerConfig,
		ControlPlaneConfig: &KubernetesControlPlaneConfig{
			HAConfig: r.HAConfig,
		},
	}
}

func (r *JoinClusterRequest) Validate() error {
	var valList []validation.Validator

	valList = append(valList, r.WorkerConfig)
	if r.NodeConfig != nil {
		valList = append(valList, r.NodeConfig)
	}

	if r.RegistryConfig != nil {
		valList = append(valList, r.RegistryConfig)
	}

	if r.HAConfig != nil {
		if r.HAConfig.KubeadmConfigYaml == "" {
			return errors.Errorf("kubeadm config is empty when ha")
		}
		valList = append(valList, r.HAConfig)
	}

	if err := validator.ValidateList(valList...); err != nil {
		return errors.WithStack(err)
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

type (
	GetJoinCommandRequest struct {
		ShowControlPlane bool `json:"show_control_plane"`
	}

	GetJoinCommandResponse struct {
		JoinCommand     string `json:"join_command"`
		KubeadmYamlData string `json:"kubeadm_yaml_data,omitempty"`
	}
)

type (
	GetKubeConfigRequest struct {
	}

	GetKubeConfigResponse struct {
		Data string `json:"data"`
	}
)

type (
	GetInstallLogRequest struct {
	}

	GetInstallLogResp struct {
		Data string `json:"data"`
	}
)
