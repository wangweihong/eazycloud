package ikubeagent

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/robfig/cron"
	"github.com/wangweihong/eazycloud/pkg/validator"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/netutil"
	"github.com/wangweihong/gotoolbox/pkg/sets"
	"github.com/wangweihong/gotoolbox/pkg/stringutil"
	"github.com/wangweihong/gotoolbox/pkg/validation"
)

type KubernetesDeployConfig struct {
	RegistryConfig    *Registry `json:"registry_config" binding:"required,dive"`
	KubernetesVersion string    `json:"kubernetes_version"`
	// master配置, 传递了则认为是master
	ControlPlaneConfig *KubernetesControlPlaneConfig `json:"control_plane_config" binding:"required"`
	// worker配置, 传递了则认为是worker
	WorkerConfig *KubernetesWorkerConfig `json:"worker_config"`
	NodeConfig   *KubernetesNodeConfig   `json:"node_config" binding:"required"`
	// etcd快照策略
	DataBaseConfig *EtcdDatabaseConfig `json:"data_base_config"`
	// control plane也作为工作节点, 一体机环境设置
	MasterAlsoWorker bool `json:"master_also_worker"`
}

func (r *KubernetesDeployConfig) Validate() error {
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

type KubernetesControlPlaneConfig struct {
	NetworkPlugin    *KubernetesNetworkPluginConfig `json:"network_plugin"`
	HAConfig         *KubernetesHAConfig            `json:"ha_config"`
	MonitorPlugin    *KubernetesMonitorPluginConfig `json:"monitor_plugin"`
	KubectlPlugin    *KubernetesKubectlPluginConfig `json:"kubectl_plugin"`
	KubeProxy        *KubeProxyConfig               `json:"kube_proxy"`
	StoragePlugin    *KubernetesStoragePluginConfig `json:"storage_plugin"`
	Gpu              *GpuConfig                     `json:"gpu"`
	Npu              *NpuConfig                     `json:"npu"`
	Namespace        *NamespaceConfig               `json:"namespace"`
	ServiceCIDR      string                         `json:"service_cidr" binding:"omitempty,cidr"`
	PodCIDR          string                         `json:"pod_cidr" binding:"omitempty,cidr"`
	ServiceDnsDomain string                         `json:"service_dns_domain"`
}

func (r *KubernetesControlPlaneConfig) Validate() error {
	if err := validator.ValidateAll(r); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

type GpuConfig struct{}

type NpuConfig struct{}

type DcuConfig struct{}

type Registry struct {
	Address       string     `json:"address" binding:"required"`
	Project       string     `json:"project"`
	User          string     `json:"user"`
	Password      string     `json:"password"`
	TlsConfig     *TlsConfig `json:"tls_config"`
	SkipTlsVerify bool       `json:"skip_tls_verify"`
	MutualTls     bool       `json:"mutual_tls"`
}

func (r *Registry) Validate() error {
	if r == nil {
		return nil
	}

	if r.Address == "" {
		return errors.Errorf("registry address is empty")
	}

	if !r.SkipTlsVerify {
		if r.TlsConfig == nil || r.TlsConfig.CaData == "" {
			return errors.Errorf("ca is empty when SkipTlsVerify == true")
		}
	}

	return nil
}

type TlsConfig struct {
	CaData     string `json:"ca_data"`
	ClientCert string `json:"client_cert"`
	ClientKey  string `json:"client_key"`
}

type EtcdDatabaseBackupPolicy struct {
	BackupScheduleCron string `json:"backup_schedule_cron"`
	BackupHostPath     string `json:"backup_host_path"`
}

type KubernetesHAConfig struct {
	// 集群虚拟IP
	VIP string `json:"vip" binding:"required"`
	// 集群对外接口
	ExternalPort int32 `json:"external_port" binding:"required,port"`
	// 控制平面节点
	ControlPlaneNodes []*KubernetesNodeConfig `json:"control_plane_nodes"`
	// keepalived路由
	RouteID int32 `json:"route_id"`

	KubeadmConfigYaml string `json:"kubeadm_config_yaml" binding:"required"`
}

func (r *KubernetesHAConfig) Validate() error {
	if r.VIP == "" || net.ParseIP(r.VIP) == nil {
		return errors.Errorf("invalid vip '%v'", r.VIP)
	}

	if netutil.PingV2(r.VIP) {
		return errors.Errorf("vip %v has beed used", r.VIP)
	}

	if !validation.IsPortValid(int(r.ExternalPort)) || validation.IsPortUsed(int(r.ExternalPort)) {
		return errors.Errorf("external port  %v invalid or used", r.ExternalPort)
	}

	localIPs, err := netutil.GetLocalIPsV2(false, func(iface net.Interface) bool {
		return stringutil.HasAnyPrefix(iface.Name, "e")
	})
	if err != nil {
		return errors.WithStack(err)
	}

	if len(localIPs) == 0 {
		return errors.Errorf("cannot get any localIP")
	}

	if r.ControlPlaneNodes != nil {
		var isFound bool
		var iplist []string

		for _, v := range r.ControlPlaneNodes {
			if v.Priority <= 0 {
				return errors.Errorf("ha node priority must > 0 ")
			}

			iplist = append(iplist, v.IP)
			if sets.NewString(localIPs...).Has(v.IP) {
				isFound = true
				break
			}
		}
		if !isFound {
			return errors.Errorf("current ip:%v execlude in deploy ip lists [%v]", localIPs, iplist)
		}
	} else {
		r.ControlPlaneNodes = append(r.ControlPlaneNodes, &KubernetesNodeConfig{
			IP: localIPs[0],
		})
	}

	return nil
}

type KubernetesJoinMasterHAConfig struct {
	ControlPlaneNodes []*KubernetesNodeConfig `json:"control_plane_nodes"`
	KubeadmConfigYaml string                  `json:"kubeadm_config_yaml" binding:"required"`
}

func (r *KubernetesJoinMasterHAConfig) Validate() error {
	localIPs, err := netutil.GetLocalIPsV2(false, func(iface net.Interface) bool {
		return stringutil.HasAnyPrefix(iface.Name, "e")
	})
	if err != nil {
		return errors.WithStack(err)
	}

	if len(localIPs) == 0 {
		return errors.Errorf("cannot get any localIP")
	}

	if r.ControlPlaneNodes != nil {
		var isFound bool
		var iplist []string

		for _, v := range r.ControlPlaneNodes {
			if v.Priority <= 0 {
				return errors.Errorf("ha node priority must > 0 ")
			}

			iplist = append(iplist, v.IP)
			if sets.NewString(localIPs...).Has(v.IP) {
				isFound = true
				break
			}
		}
		if !isFound {
			return errors.Errorf("current ip:%v execlude in deploy ip lists [%v]", localIPs, iplist)
		}
	} else {
		r.ControlPlaneNodes = append(r.ControlPlaneNodes, &KubernetesNodeConfig{
			IP: localIPs[0],
		})
	}

	return nil
}

type KubernetesNodeConfig struct {
	IP       string `json:"ip"`
	NodeName string `json:"node_name"`
	// control plane优先级, 仅对control plane节点有效
	Priority int32 `json:"priority"`
}

func (r *KubernetesNodeConfig) Validate() error {
	if r == nil {
		r = &KubernetesNodeConfig{}
	}
	if r.NodeName == "" {
		r.NodeName, _ = os.Hostname()
	}

	if err := validator.IsDNS1123Label(r.NodeName); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

type KubernetesNetworkPluginConfig struct {
	PluginType   string        `json:"plugin_type"`
	CalicoConfig *CalicoConfig `json:"calico_config"`
}

func (r *KubernetesNetworkPluginConfig) Validate() error {
	if r == nil {
		return fmt.Errorf("network plugin is empty")
	}
	return nil
}

type CalicoConfig struct {
	Version string `json:"version"`
}

type KubernetesMonitorPluginConfig struct {
}

func (r *KubernetesMonitorPluginConfig) Validate() error {
	return nil
}

type KubernetesKubectlPluginConfig struct {
}

func (r *KubernetesKubectlPluginConfig) Validate() error {
	return nil
}

type KubeProxyConfig struct {
	Mode string `json:"mode"`
}

func (r *KubeProxyConfig) Validate() error {
	switch r.Mode {
	case "iptables":
	case "ipvs":
	case "":
	default:
		return errors.Errorf("unsupport kube-proxy mode %v", r.Mode)
	}
	return nil
}

type KubernetesStoragePluginConfig struct {
	Local *LocalStoragePluginConfig `json:"local"`
	Nfs   *NfsStoragePluginConfig   `json:"nfs"`
	// 允许部署集群时创建一个共用的NFS PV卷
	NFSVolume *NFSVolumeConfig `json:"nfs_volume"`
}

func (r *KubernetesStoragePluginConfig) Validate() error {
	return nil
}

type LocalStoragePluginConfig struct{}
type NfsStoragePluginConfig struct{}
type NFSVolumeConfig struct {
	ServerIP  string `json:"server_ip"`
	MountPath string `json:"mount_path"`
}

type KubernetesWorkerConfig struct {
	JoinCommand    string `json:"join_command" binding:"required"`
	IsControlPlane bool   `json:"-"`
}

func (r *KubernetesWorkerConfig) Validate() error {
	if r.JoinCommand == "" {
		return errors.Errorf("Join command is empty")
	}

	if strings.Contains(r.JoinCommand, "--control-plane") && !strings.Contains(r.JoinCommand, "--certificate-key") {
		return errors.Errorf("invalid join Command, --control-plane must with --certificate-key")
	}

	return nil
}

type EtcdDatabaseConfig struct {
	Local        *LocalDatabaseConfig      `json:"local"`
	External     *ExternalDatabaseConfig   `json:"external"`
	BackupPolicy *EtcdDatabaseBackupPolicy `json:"backup_policy"`
}

func (r *EtcdDatabaseConfig) Validate() error {
	if r.Local != nil {
		if r.Local.PeerPort != nil {
			if validation.IsPortValid(*r.Local.PeerPort) || validation.IsPortUsed(*r.Local.PeerPort) {
				return errors.Errorf("peer port '%v' is invalid or used", *r.Local.PeerPort)
			}
		}

		if r.Local.ListenPort != nil {
			if validation.IsPortValid(*r.Local.ListenPort) || validation.IsPortUsed(*r.Local.ListenPort) {
				return errors.Errorf("listen port '%v' is invalid or used", *r.Local.ListenPort)
			}
		}
	}

	if r.BackupPolicy != nil {
		if r.BackupPolicy.BackupHostPath == "" || !filepath.IsAbs(r.BackupPolicy.BackupHostPath) {
			return errors.Errorf("backup host path '%v' is no abs path", r.BackupPolicy.BackupHostPath)
		}

		specParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		if _, err := specParser.Parse(r.BackupPolicy.BackupScheduleCron); err != nil {
			return errors.Errorf("parse schedule cron fail:%v", err.Error())
		}
	}

	return nil
}

type LocalDatabaseConfig struct {
	HostPath   string `json:"host_path"`
	ListenPort *int   `json:"listen_port"`
	PeerPort   *int   `json:"peer_port"`
}

type ExternalDatabaseConfig struct {
	Endpoints []string   `json:"endpoints" `
	TlsConfig *TlsConfig `json:"tls_config"`
}

type NamespaceConfig struct{}

type DockerRegistryAccount struct {
	User     string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	Auth     string `json:"auth,omitempty"`

	ServerAddress string `json:"serveraddress,omitempty"`
}

type DockerRegistryAccounts struct {
	Auths map[string]DockerRegistryAccount `json:"auths"`
}
