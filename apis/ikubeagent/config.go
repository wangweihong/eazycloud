package ikubeagent

type KubernetesDeployConfig struct {
	RegistryConfig    *Registry
	KubernetesVersion string
	// master配置, 传递了则认为是master
	ControlPlaneConfig *KubernetesControlPlaneConfig
	// worker配置, 传递了则认为是worker
	WorkerConfig *KubernetesWorkerPlaneConfig
	NodeConfig   *KubernetesNodeConfig
	// etcd快照策略
	DataBaseConfig *EtcdDatabaseConfig
	// control plane也作为工作节点
	MasterAlsoWorker bool
}

type KubernetesControlPlaneConfig struct {
	ConfigYaml          string
	NetworkPluginConfig *KubernetesNetworkPluginConfig
	HAConfig            *KubernetesHAConfig
	MonitorPluginConfig *KubernetesMonitorPluginConfig
	KubectlPluginConfig *KubernetesKubectlPluginConfig
	KubeProxyConfig     *KubeProxyConfig
	StoragePluginConfig *KubernetesStoragePluginConfig
	GpuConfig           *GpuConfig
	NamespaceConfig     *NamespaceConfig
	ServiceCIDR         string
	PodCIDR             string
	ServiceDnsDomain    string
}

type GpuConfig struct{}

type Registry struct {
	Address       string
	Project       string
	User          string
	Password      string
	TlsConfig     *TlsConfig
	SkipTlsVerify bool
	MutualTls     bool
}

type TlsConfig struct {
	CaData     string
	ClientCert string
	ClientKey  string
}

type EtcdDatabaseBackupPolicy struct {
	BackupScheduleCron string
	BackupHostPath     string
}

type KubernetesHAConfig struct {
	// 集群虚拟IP
	VIP string
	// 集群对外接口
	ExternalPort int32
	// 控制平面节点
	ControlPlaneNodes []*KubernetesNodeConfig
	// keepalived路由
	RouteID           int32
	KubeadmConfigYaml string
}

type KubernetesNodeConfig struct {
	IP       string
	NodeName string
	// control plane优先级, 仅对control plane节点有效
	Priority int32
}

type KubernetesNetworkPluginConfig struct {
	PluginType   string
	CalicoConfig *CalicoConfig
}

type CalicoConfig struct {
	Version string
}

type KubernetesMonitorPluginConfig struct {
}

type KubernetesKubectlPluginConfig struct {
}

type KubeProxyConfig struct {
	Mode string
}

type KubernetesStoragePluginConfig struct {
	LocalStorage           *LocalStoragePluginConfig
	NfsStoragePluginConfig *NfsStoragePluginConfig
	// 允许部署集群时创建一个共用的NFS PV卷
	NFSVolumeConfig *NFSVolumeConfig
}

type LocalStoragePluginConfig struct{}
type NfsStoragePluginConfig struct{}
type NFSVolumeConfig struct {
	ServerIP  string
	MountPath string
}

type KubernetesWorkerPlaneConfig struct {
	ConfigYaml  string
	JoinCommand string
	Sync        bool
}

type EtcdDatabaseConfig struct {
	Local        *LocalDatabaseConfig
	External     *ExternalDatabaseConfig
	BackupPolicy *EtcdDatabaseBackupPolicy
}

type LocalDatabaseConfig struct {
	HostPath string
}

type ExternalDatabaseConfig struct {
	Endpoints []string
	TlsConfig *TlsConfig
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
