package ikubeagent

import (
	"regexp"

	"github.com/wangweihong/gotoolbox/pkg/errors"
)

const (
	KubernetesDeployStateUninitialized = "uninitialized"
	KubernetesDeployStateDeploying     = "deploying"
	KubernetesDeployStateSuccess       = "success"
	KubernetesDeployStateError         = "error"

	NetworkPluginCalico = "calico"

	CalicoV3_14 = "v3.14"

	KubernetesVersionV1_18_0 = "v1.18.0"
	KubernetesVersionV1_30_0 = "v1.30.0"

	KubeadmBinary                = "/usr/bin/kubeadm"
	KubeletBinary                = "/usr/bin/kubelet"
	KubeletPreRunScript          = "/usr/bin/kubelet-pre-start.sh"
	KubeletServicePath           = "/lib/systemd/system/kubelet.service"
	KubeletServiceConfigDir      = "/etc/systemd/system/kubelet.service.d"
	KubeletServiceConfigFilePath = "/etc/systemd/system/kubelet.service.d/10-kubeadm.conf"
	KubeadmLogPath               = "/var/log/kubegent/kubeadm.log"
	SystemCtlBinary              = "systemctl"
	DockerDeamonJsonPath         = "/etc/docker/daemon.json"
	ConntrackBinary              = "/usr/bin/conntrack"
	KubeConfigPath               = "/etc/kubernetes/admin.conf"
	KubectlBinary                = "/usr/bin/kubectl"
	KubeletRootDir               = "/var/lib/kubelet"
	KubeletRegistryConfigPath    = "/var/lib/kubelet/config.json"
	KubeadmConfigYamlPath        = "/var/lib/kubelet/kubeadm.yaml"
	ManifestsDir                 = "/etc/kubernetes/manifests"
	EtcdRootDir                  = "/var/lib/etcd"
	TemplateDirPath              = "/var/lib/kubeagent/template"
	AddonDir                     = "/etc/kubernetes/addons"
	PatchDir                     = "/etc/kubernetes/patch"
	KustomizeBinary              = "/usr/bin/kustomize"
	DockerCertDir                = "/etc/docker/certs.d"
	EtcdCaPath                   = "/etc/kubernetes/pki/etcd/"

	CalicoDefaultPodSubnet    = "192.168.0.0/16"
	DefaultServiceSubnet      = "10.96.0.0/12"
	DefaultRepository         = "k8s.gcr.io"
	PriRegistrySecretFilePath = "/etc/kubernetes/manifests/private-registry-secret.yaml"
)

const DNS1123LabelMaxLength int = 63

const dns1123LabelFmt string = "[a-z0-9]([-a-z0-9]*[a-z0-9])?"
const dns1123SubdomainFmt string = dns1123LabelFmt + "(\\." + dns1123LabelFmt + ")*"

var dns1123LabelRegexp = regexp.MustCompile("^" + dns1123LabelFmt + "$")

// IsDNS1123Label check node name is ok
func IsDNS1123Label(value string) error {
	if value == "" {
		return errors.New("name is empty")
	}

	if len(value) > DNS1123LabelMaxLength {
		return errors.New("name too long")
	}
	if !dns1123LabelRegexp.MatchString(value) {
		return errors.New("not match name pattern: " + dns1123LabelRegexp.String())
	}
	return nil
}
