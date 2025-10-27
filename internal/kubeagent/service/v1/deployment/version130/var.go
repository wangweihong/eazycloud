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
		TemplateSubPath: filepath.Join("systemd", "kubelet_service_config.template"),
		TemplateText:    nil,
		Context:         map[string]any{},
		FilePath:        ikubeagent.KubeletServiceConfigFilePath,
	}

	// KubeletPreRunScriptTemplate
	KubeletPreRunScriptTemplate = template.FileProcessor{
		TemplateName: "kubelet_pre_run_script.template",
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, "systemd", "kubelet_pre_run_script.template"),
		TemplateText: nil,
		Context:      map[string]any{},
		FilePath:     ikubeagent.KubeletPreRunScript,
	}

	KubeletServiceTemplate = template.FileProcessor{
		TemplateName: "kubelet_service.template",
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, "systemd", "kubelet_service.template"),
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
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, "misc", "kubeadm_config.template"),
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
			"KubeProxyMode": "",
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
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, "systemd", "containerd_config_toml.template"),
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
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, "systemd", "containerd_registry_hosts.template"),
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
		FilePath: "/etc/kubernetes/addons/namespace.yaml",
	}
)

var (
	CalicoManifest = template.DirectoryProcessor{
		TemplateName: "calico.template",
		TemplateDir:  filepath.Join(ikubeagent.TemplateDirPath, "network", "calico"),
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
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, "metrics_server.template"),
		Context: map[string]any{
			"ImageRepository": "",
		},
		FilePath: filepath.Join(ikubeagent.AddonDir, "monitor", "metrics_server.yaml"),
	}

	MonitorStackTemplate = template.DirectoryProcessor{
		TemplateName: "monitor-stack",
		TemplateDir:  filepath.Join(ikubeagent.TemplateDirPath, "kube-prometheus-0.6.0", "monitor-stack"),
		Context: map[string]any{
			"ImageRepository": "",
		},
		LocateParsedDir: filepath.Join(ikubeagent.AddonDir, "monitor", "monitor-stack"),
	}
	//monitorOperators Must deploy before monitor stack
	MonitorOperatorsTemplate = template.DirectoryProcessor{
		TemplateName: "monitor-operators",
		TemplateDir:  filepath.Join(ikubeagent.TemplateDirPath, "kube-prometheus-0.6.0", "monitor-operator"),
		Context: map[string]any{
			"ImageRepository": "",
		},
		LocateParsedDir: filepath.Join(ikubeagent.AddonDir, "monitor", "monitor-operator"),
	}
)

var (
	LocalStorageAllInOneTemplate = template.DirectoryProcessor{
		TemplateName: "local-storage-controller",
		TemplateDir:  filepath.Join(ikubeagent.TemplateDirPath, "storage", "local"),
		Context: map[string]any{
			"ImageRepository": "",
			"Namespace":       "",
		},
		LocateParsedDir: filepath.Join(ikubeagent.AddonDir, "storage", "local"),
	}

	NFSStorageClassTemplate = template.DirectoryProcessor{
		TemplateName:    "nfs-storage-class",
		TemplateDir:     filepath.Join(ikubeagent.TemplateDirPath, "storage", "nfs"),
		Context:         map[string]any{},
		LocateParsedDir: filepath.Join(ikubeagent.AddonDir, "storage", "nfs"),
	}

	NFSPVTemplate = template.FileProcessor{
		TemplateName: "nfs-pv-pvc",
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, "nfs-volume.template"),
		Context: map[string]any{
			"NFSIP":       "",
			"NFSHostPath": "/var/lib/deploy/nfs",
		},
		FilePath: filepath.Join(ikubeagent.AddonDir, "nfs-volume.yaml"),
	}
)

var (
	// namespacecontroller do some work when receive namespace events, such as create podpreset to keep timezone sync
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
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, "kubectl.yaml.template"),
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
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, "ha", "keepalived_config.template"),
		Context: map[string]any{
			"STATE":                "",
			"INTERFACE":            "",
			"ROUTER_ID":            "50",
			"AUTH_PASS":            "42",
			"ikubeagentSERVER_VIP": "",
			"PRIORITY":             "",
		},
		FilePath: "/etc/keepalived/keepalived.conf",
	}

	KeepalivedCheckApiserverScript = template.FileProcessor{
		TemplateName: "keepalived_vrrp_script.template",
		TemplateText: nil,
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, "ha", "keepalived_vrrp_script.template"),
		Context: map[string]any{
			"EXTERNAL_PORT":        "",
			"ikubeagentSERVER_VIP": "",
		},
		FilePath: "/etc/keepalived/check_ikubeagentserver.sh",
	}
	KeepalivedYamlTemplate = template.FileProcessor{
		TemplateName: "keepalived_pod.template",
		TemplateText: nil,
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, "ha", "keepalived_pod.template"),
		Context: map[string]any{
			"ImageRepository": "",
		},
		FilePath: "/etc/kubernetes/manifests/keepalived.yaml",
	}

	HaproxyYamlTemplate = template.FileProcessor{
		TemplateName: "haproxy_pod.template",
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, "ha", "haproxy_pod.template"),
		TemplateText: nil,
		Context: map[string]any{
			"ImageRepository": "library",
			"ExternalPort":    "8443",
		},
		FilePath: "/etc/kubernetes/manifests/haproxy.yaml",
	}

	HaproxyConfigTemplate = template.FileProcessor{
		TemplateName: "haproxy_cfg.template",
		TemplatePath: filepath.Join(ikubeagent.TemplateDirPath, "ha", "haproxy_cfg.template"),
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
		TemplateDir:     filepath.Join(ikubeagent.TemplateDirPath, "gpu"),
		LocateParsedDir: filepath.Join(ikubeagent.AddonDir, "gpu"),
	}
)
