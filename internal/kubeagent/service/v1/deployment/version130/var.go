package version130

import (
	"path/filepath"

	"github.com/wangweihong/eazycloud/apis/ikubeagent"

	gotemplate "text/template"

	"github.com/lithammer/dedent"
	"github.com/wangweihong/gotoolbox/pkg/template"
)

var (
	KubeletServiceConfigTemplate = template.FileProcessor{
		TemplateName: "kubelet_service_config.template",
		//TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, "systemd", "kubelet_service_config.template"),
		TemplateDir:     ikubeagent.TemplateDirPath,
		TemplateSubPath: systemd_kubelet_service_config_template.Path,
		TemplateText:    nil,
		Context:         map[string]any{},
		FilePath:        ikubeagent.KubeletServiceConfigFilePath,
	}

	// KubeletPreRunScriptTemplate
	KubeletPreRunScriptTemplate = template.FileProcessor{
		TemplateName: "kubelet_pre_run_script.template",
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, systemd_kubelet_pre_run_script_template.Path),
		TemplateText: nil,
		Context:      map[string]any{},
		FilePath:     ikubeagent.KubeletPreRunScript,
	}

	KubeletServiceTemplate = template.FileProcessor{
		TemplateName: "kubelet_service.template",
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, systemd_kubelet_service_template.Path),
		TemplateText: nil,
		Context: map[string]any{
			"KubeletBinaryPath":       "",
			"KubeletPreRunScriptPath": "",
		},
		FilePath: ikubeagent.KubeletServicePath,
	}
)

var (
	KubeadmConfigTemplate = template.FileProcessor{
		TemplateName: "kubeadm_config.template",
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, misc_kubeadm_config_template.Path),
		TemplateText: nil,
		Context: map[string]any{
			// 当前控制面板服务地址
			"LocalAdvertiseAddress": "",
			// 当前控制面板服务端口
			"LocalBindPort":     "6443",
			"ImageRepository":   "registry.k8s.io",
			"KubernetesVersion": "v1.30.0",
			//NodeRegistrationName 节点名。除非用户指定，这里置空，由kubeadm填充当前节点名。
			"NodeRegistrationName": "",
			// 集群控制面板对外端点。仅用于高可用集群主节点群统一对外的服务端点。
			"ControlPlaneEndpoint": "",
			// 集群DNS服务地址, 默认为服务网段第10个ip.
			"ClusterDNS": "10.96.0.10",
			// 服务网段
			"ServiceSubnet": "10.96.0.0/12",
			// 集群Pod网段, 注意不能和已存在网络冲突。还必须和插件保持一致, 如calico
			"PodSubnet": "172.18.0.0/16",
			// KubeProxy模型,,为空默认是iptables, 可选iptables,ipvs
			"KubeProxyMode":  "",
			"EtcdListenPort": "2379",
			"EtcdPeerPort":   "2380",
		},
		FilePath: ikubeagent.KubeadmConfigYamlPath,
	}
)

const (
	certDir = "/etc/containerd/certs.d"
)

var (
	ContainerdConfigTemplate = template.FileProcessor{
		TemplateName: "containerd_config.template",
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, systemd_containerd_config_toml_template.Path),
		TemplateText: nil,
		Context: map[string]any{
			"RegistryConfigPath": certDir,
			"ImageRepository":    "",
			"SystemdCgroup":      "false",
		},
		FilePath: "/etc/containerd/config.toml",
	}

	ContainerdRegistryConfigTemplate = template.FileProcessor{
		TemplateName: "containerd_registry.template",
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, systemd_containerd_registry_hosts_template.Path),
		TemplateText: nil,
		Context: map[string]any{
			"RegistryAddress":      "",
			"RegistryCAPath":       "",
			"RegistryEnableVerify": "false",
			"RegistryAuth":         "",
		},
		FilePath: filepath.Join(certDir, "hosts.toml"),
	}
)
var (
	NamespacesManifests = template.FileProcessor{
		TemplateName: "namespaces.template",
		TemplateData: `
ikubeagentVersion: v1
kind: Namespace
metadata:
  name: kubeagent-namespace
`,
		Context:  map[string]any{},
		FilePath: filepath.Join(ikubeagent.AddonDir, "namespace.yaml"),
	}
)

var (
	CalicoManifest = template.DirectoryProcessor{
		TemplateName: "calico.template",
		TemplateDir:  filepath.Join(ikubeagent.TemplateDirPath, filepath.Dir(network_calico_custom_resource_yaml_template.Path)),
		Context: map[string]any{
			"ImageRepository":             "",
			"PodSubnet":                   ikubeagent.CalicoDefaultPodSubnet,
			"ImageSecretDockerConfigJson": "",
		},
		LocateParsedDir: filepath.Join(ikubeagent.AddonDir, "network", "calico"),
	}
)

var (
	ImageRegistrySecretTemplate = template.FileProcessor{
		TemplateText: gotemplate.Must(gotemplate.New("imageRegistrySecret").Parse(dedent.Dedent(`
			ikubeagentVersion: v1
			data:
			  .dockerconfigjson: {{.ImageSecretDockerConfigJson}}
			kind: Secret
			metadata:
			  name: {{.ImageSecretName}}
			  namespace: {{.ImageSecretNamespace}}
			type: kubernetes.io/dockerconfigjson`))),
		Context: map[string]any{
			"ImageSecretDockerConfigJson": "",
			"ImageSecretName":             "",
			"ImageSecretNamespace":        "kube-system",
		},
		FilePath: ikubeagent.PriRegistrySecretFilePath,
	}
)

func NewImageRegistrySecretTemplate(fp string) template.FileProcessor {
	return template.FileProcessor{
		TemplateText: gotemplate.Must(gotemplate.New("imageRegistrySecret").Parse(dedent.Dedent(`
			ikubeagentVersion: v1
			data:
			  .dockerconfigjson: {{.ImageSecretDockerConfigJson}}
			kind: Secret
			metadata:
			  name: {{.ImageSecretName}}
			  namespace: {{.ImageSecretNamespace}}
			type: kubernetes.io/dockerconfigjson`))),
		Context: map[string]any{
			"ImageSecretDockerConfigJson": "",
			"ImageSecretName":             "",
			"ImageSecretNamespace":        "kube-system",
		},
		FilePath: fp,
	}
}

var (
	MetricsServerTemplate = template.FileProcessor{
		TemplateName: "metrics_server.template",
		TemplateText: nil,
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, monitor_metrics_server_template.Path),
		Context: map[string]any{
			"ImageRepository": "",
		},
		FilePath: filepath.Join(ikubeagent.AddonDir, "monitor", "metrics_server.yaml"),
	}

	MonitorStackTemplate = template.DirectoryProcessor{
		TemplateName: "monitor-stack",
		TemplateDir:  filepath.Join(ikubeagent.TemplateDirPath, filepath.Dir(monitor_kube_prometheus_0_6_0_monitor_stack_alertmanager_secret_yaml.Path)),
		Context: map[string]any{
			"ImageRepository": "",
		},
		LocateParsedDir: filepath.Join(ikubeagent.AddonDir, "monitor", "monitor-stack"),
	}
	//monitorOperators Must deploy before monitor stack
	MonitorOperatorsTemplate = template.DirectoryProcessor{
		TemplateName: "monitor-operators",
		TemplateDir:  filepath.Join(ikubeagent.TemplateDirPath, filepath.Dir(monitor_kube_prometheus_0_6_0_monitor_operator_0namespace_namespace_yaml.Path)),
		Context: map[string]any{
			"ImageRepository": "",
		},
		LocateParsedDir: filepath.Join(ikubeagent.AddonDir, "monitor", "monitor-operator"),
	}
)

var (
	LocalStorageAllInOneTemplate = template.DirectoryProcessor{
		TemplateName: "local-storage-controller",
		TemplateDir:  filepath.Join(ikubeagent.TemplateDirPath, filepath.Dir(storage_local_local_storage_class_yaml.Path)),
		Context: map[string]any{
			"ImageRepository": "",
			"Namespace":       "",
		},
		LocateParsedDir: filepath.Join(ikubeagent.AddonDir, "storage", "local"),
	}

	NFSStorageClassTemplate = template.DirectoryProcessor{
		TemplateName:    "nfs-storage-class",
		TemplateDir:     filepath.Join(ikubeagent.TemplateDirPath, filepath.Dir(storage_nfs_nfs_storage_class_yaml.Path)),
		Context:         map[string]any{},
		LocateParsedDir: filepath.Join(ikubeagent.AddonDir, "storage", "nfs"),
	}

	NFSPVTemplate = template.FileProcessor{
		TemplateName: "nfs-pv-pvc",
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, misc_nfs_volume_template.Path),
		Context: map[string]any{
			"NFSIP":       "",
			"NFSHostPath": "/var/lib/deploy/nfs",
		},
		FilePath: filepath.Join(ikubeagent.AddonDir, "nfs-volume.yaml"),
	}
)

var (
	// namespacecontroller do some work： when receive namespace events, such as create podpreset to keep timezone sync
	NamespacesControllerManifests = template.DirectoryProcessor{
		TemplateName: "namespace-controllers",
		Context: map[string]any{
			"ImageRepository": "",
		},
		TemplateDir:     filepath.Join(ikubeagent.TemplateDirPath, "namespace-controller"),
		LocateParsedDir: filepath.Join(ikubeagent.AddonDir, "namespace-controller"),
	}
)

var (
	KubectlManifest = template.FileProcessor{
		TemplateName: "kubectl.yaml.template",
		TemplateText: nil,
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, misc_kubectl_yaml_template.Path),
		Context: map[string]any{
			"ImageRepository": "k8s.gcr.io",
			"Namespace":       "",
		},
		FilePath: filepath.Join(ikubeagent.AddonDir, "kubectl.yaml"),
	}
)

var (
	KeepalivedConfigTempl = template.FileProcessor{
		TemplateName: "keepalived_config.template",
		TemplateText: nil,
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, ha_keepalived_config_template.Path),
		Context: map[string]any{
			"STATE":                "",
			"INTERFACE":            "",
			"ROUTER_ID":            "50",
			"AUTH_PASS":            "42",
			"APISERVER_VIP": "",
			"PRIORITY":             "",
		},
		FilePath: "/etc/keepalived/keepalived.conf",
	}

	KeepalivedCheckApiserverScript = template.FileProcessor{
		TemplateName: "keepalived_vrrp_script.template",
		TemplateText: nil,
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, ha_keepalived_vrrp_script_template.Path),
		Context: map[string]any{
			"EXTERNAL_PORT": "",
			"APISERVER_VIP": "",
		},
		FilePath: "/etc/keepalived/check_apiserver.sh",
	}
	KeepalivedYamlTemplate = template.FileProcessor{
		TemplateName: "keepalived_pod.template",
		TemplateText: nil,
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, ha_keepalived_pod_template.Path),
		Context: map[string]any{
			"ImageRepository": "",
		},
		FilePath: "/etc/kubernetes/manifests/keepalived.yaml",
	}

	HaproxyYamlTemplate = template.FileProcessor{
		TemplateName: "haproxy_pod.template",
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, ha_haproxy_pod_template.Path),
		TemplateText: nil,
		Context: map[string]any{
			"ImageRepository": "library",
			"ExternalPort":    "8443",
		},
		FilePath: "/etc/kubernetes/manifests/haproxy.yaml",
	}

	HaproxyConfigTemplate = template.FileProcessor{
		TemplateName: "haproxy_cfg.template",
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, ha_haproxy_cfg_template.Path),
		TemplateText: nil,
		Context: map[string]any{
			"HAPROXY_PORT": "8443",
			"ServerIPs":    []string{},
			"ServerPort":   6443,
		},
		FilePath: "/etc/haproxy/haproxy.cfg",
	}
	Gpu = template.DirectoryProcessor{
		TemplateName: "gpu",
		Context: map[string]any{
			"ImageRepository": "",
		},
		TemplateDir:     filepath.Join(ikubeagent.TemplateDirPath, filepath.Dir(gpu_nfd_crd_yaml.Path)),
		LocateParsedDir: filepath.Join(ikubeagent.AddonDir, "gpu"),
	}

	Npu = template.DirectoryProcessor{
		TemplateName: "npu",
		TemplateDir:  filepath.Join(ikubeagent.TemplateDirPath, filepath.Dir(npu_device_plugin_310p_v6_0_0_yaml_template.Path)),
		Context: map[string]interface{}{
			"ImageRepository": "",
		},
		LocateParsedDir: filepath.Join(ikubeagent.AddonDir, "npu"),
	}
)
