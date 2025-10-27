package version130

import (
	"os"
	osexec "os/exec"
	"path/filepath"
	"strings"

	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/executil"
	"github.com/wangweihong/gotoolbox/pkg/log"

	"github.com/wangweihong/eazycloud/apis/ikubeagent"
	"github.com/wangweihong/eazycloud/internal/kubeagent/service/v1/deployment"
	"github.com/wangweihong/eazycloud/internal/kubeagent/store"
	"github.com/wangweihong/eazycloud/internal/pkg/run"
)

func NewDeployService(version string, store store.Factory) deployment.AgentService {
	return &impl{
		version: version,
		deps: []string{
			"kubeadm",
			"kubelet",
			"conntrack",
			"kubectl",
			"containerd",
			"crictl",
			"containerd.service",
			// 这个非必须得
			//"systemd-resolved.service",
			//TODO: nvidia-ctl?
		},
	}
}

type NamedOperation struct {
	Name string
	Op   func(c *setupContext) error
}

type SetupContext struct {
	dc        *ikubeagent.KubernetesDeployConfig
	isMaster0 bool
}

type impl struct {
	version string
	deps    []string
	store   store.Factory
}

func (i *impl) Version() string {
	return i.version
}

var (
	initMasterProcess = []setupFunc{
		prepareDataDir,
	}
)

func (i *impl) InitMaster(config *ikubeagent.KubernetesDeployConfig, isMaster0 bool) (err error) {
	log.Infof("start to init master, isMaster0:%v", isMaster0)

	if config.ControlPlaneConfig == nil {
		log.Errorf("init master must parsing k8s cluster control plane param")
		return errors.Errorf("missing ControlPlaneConfig")
	}

	c, err := initSetupContext(config, isMaster0)
	if err != nil {
		return err
	}

	Ops := []NamedOperation{
		{Name: "prepare data dir", Op: prepareDataDir},
		{Name: "prepare component template ", Op: prepareComponentTemplate},
		{Name: "update kubelet service", Op: updateKubeletService},
		{Name: "update containerd service", Op: updateContainerdService2},
		{Name: "prepare ha", Op: deployPrepareHA},
		{Name: "setup k8s cluster", Op: execKubeadmInit},
		{Name: "install system namespace", Op: installSystemNamespace},
		{Name: "install network plugin", Op: installNetworkPlugin},
		{Name: "install monitor plugin", Op: installMonitorPlugin},
		{Name: "install storage plugin", Op: installStoragePlugin},
		{Name: "install kubectl plugin", Op: installKubectlPlugin},
		{Name: "set master as worker", Op: setMasterAsWorker},
		{Name: "install gpu plugin", Op: installGpuPlugin},
	}

	for _, op := range Ops {
		setupDep := LogOperation(op.Name, op.Op)
		if err := setupDep(c); err != nil {
			return err
		}
	}
	//p := postInstall{}
	// p.err = p.createSystemNamespace()
	// p.err = p.deployNetworkPlugin(config)
	//p.err = p.deployMonitorPlugin(config)
	//p.err = p.deployStoragePlugin(config)
	//p.err = p.deployNamespaceControllerPlugin(config)
	//p.err = p.deployKubectlPlugin(config)
	//p.err = p.createSecretOfEtcd(config)
	//p.err = p.setMasterAsWorker(config)
	// p.err = p.deployGpuPlugin(config)
	// if p.err != nil {
	// 	return p.err
	// }
	return nil
}

func (i *impl) JoinCluster(config *ikubeagent.KubernetesDeployConfig, joinCmd string, isControlPlane bool) (err error) {
	log.Info("start to JoinCluster ")

	c, err := initSetupContext(config, false)
	if err != nil {
		return err
	}

	Ops := []NamedOperation{
		{Name: "prepare data dir", Op: prepareDataDir},
		{Name: "prepare component template ", Op: prepareComponentTemplate},
		{Name: "update kubelet service", Op: updateKubeletService},
		{Name: "update containerd service", Op: updateContainerdService2},
		{Name: "prepare ha", Op: deployPrepareHA},
	}

	for _, op := range Ops {
		setupDep := LogOperation(op.Name, op.Op)
		if err := setupDep(c); err != nil {
			return err
		}
	}

	log.Info("Join Cluster run runKubeadmJoin")
	if err := runKubeadmJoin(joinCmd, isControlPlane, config.NodeConfig.NodeName); err != nil {
		log.Errorf("runKubeadmJoin fail:%v", err)
		return err
	}
	log.Info("JoinCluster success")
	return nil
}

func (i *impl) Reset(isHa bool) error {
	// FiXME: we should skip init etcd preflight check error if we reset ignore etcd. but at sametime we should ignore
	// init phase /var/lib/etcd remain data dir error.
	args := []string{"reset", "-f"}
	if !isHa {
		args = []string{"reset", "-f", "--skip-phases", "remove-etcd-member"}
	}
	stdout, stderr, err := executil.ExecuteCmdSplitStdoutStderr(ikubeagent.KubeadmBinary, args, 600)
	if err != nil {
		log.Errorf("kubeadm reset fail:%v", run.TrimError(err))
		if stdout != "" {
			logdata := stdout + "\n" + stderr
			if err := os.WriteFile(ikubeagent.KubeadmLogPath, []byte(logdata), 0755); err != nil {
				log.Errorf("save kubeadm reset log to path :%v fail:%v", ikubeagent.KubeadmLogPath, err)
			}
		}
		return err
	}
	if !isHa {
		if err := os.RemoveAll(ikubeagent.EtcdRootDir); err != nil {
			log.Errorf("remove dir path %v fail:%v", ikubeagent.EtcdRootDir, err)
		}
	}
	//clean deploy dirs
	_ = os.Remove(ikubeagent.KubeadmLogPath)
	_ = os.RemoveAll(ikubeagent.AddonDir)
	_ = os.RemoveAll(ikubeagent.PatchDir)

	return nil
}

func (i *impl) GetJoinCommand(showControlPlane bool) (string, string, error) {
	cmd := ikubeagent.KubeadmBinary
	args := []string{
		"token",
		"create",
		"--print-join-command",
		"--logtostderr=false",
		"--config",
		ikubeagent.KubeadmConfigYamlPath,
	}
	joinCmd, err := executil.Execute(cmd, args)
	if err != nil {
		log.Errorf("generate kubeadm join command (%v %v) fail:%v", cmd, args, err.Error())
		return "", "", err
	}
	joinCmd = strings.TrimSuffix(joinCmd, "\n")

	if showControlPlane {
		joinCmd = joinCmd + "--control-plane"
		//kubeadm-certs exist TTL is 2 minute, "kubeadm alpha certs certificate-key will fail if kubeadm-certs delete."
		// so generate a new kubeadm-certs if control-plane ha
		args = []string{
			"init",
			"phase",
			"upload-certs",
			"--upload-certs",
			"--config",
			ikubeagent.KubeadmConfigYamlPath,
			"--logtostderr=false",
		}
		output, stderr, err := executil.ExecuteCmdSplitStdoutStderr(ikubeagent.KubeadmBinary, args, 0)
		if err != nil {
			log.Errorf("upload-certs fail:%v", err.Error()+":"+stderr)
			return "", "", err
		}

		/*
			args = []string{"alpha", "certs", "certificate-key", "--logtostderr=false"}
			certkey, err := util.Execute(ikubeagent.KubeadmBinary, args)
			if err != nil {
				log.Errorf("generate kubeadm certificate-key fail:%v", err.Error())
				return "", err
			}
		*/
		log.Info(output)
		linedatas := strings.Split(output, "\n")
		joinCmd = joinCmd + " --certificate-key" + " " + linedatas[len(linedatas)-2]

		kubeadmYamlData, err := os.ReadFile(ikubeagent.KubeadmConfigYamlPath)
		if err != nil {
			log.Errorf("read kubedam yaml fail:%v", err.Error()+":"+stderr)
			return joinCmd, "", err
		}

		return joinCmd, string(kubeadmYamlData), nil
	}

	return joinCmd, "", nil
}

func (i *impl) GetKubeConfig() (string, error) {
	data, err := os.ReadFile(ikubeagent.KubeConfigPath)
	if err != nil {
		return "", errors.Errorf("read file %v fail:%v", ikubeagent.KubeConfigPath, err.Error())
	}

	return string(data), nil
}

func (i *impl) GetDeployLog() (string, error) {
	data, err := os.ReadFile(ikubeagent.KubeadmLogPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}

		return "", errors.Errorf("read file %v fail:%v", ikubeagent.KubeadmLogPath, err.Error())
	}

	return string(data), nil
}

func (i *impl) CheckDependency(version string) error {
	if i.version != version {
		return errors.Errorf("kubeadm version %v no match:%v", i.version, version)
	}
	for _, v := range i.deps {
		if strings.HasSuffix(v, ".service") {
			output, err := executil.ExecuteTimeout("systemctl", []string{"is-active", v}, 10)
			if err != nil {
				return errors.Errorf("check containerd service fail:%v", err.Error())
			}
			if strings.TrimSuffix(output, "\n") != "active" {
				return errors.Errorf("containerd service is inactive: current:%v", output)
			}
		} else {
			_, err := osexec.LookPath(v)
			if err != nil {
				return errors.Errorf("dep %v tools not exist", v)
			}
		}
	}

	return nil
}

func (i *impl) GenerateTemplate() error {
	_ = os.RemoveAll(ikubeagent.TemplateDirPath)

	for _, v := range fps {
		tp := filepath.Join(ikubeagent.TemplateDirPath, v.Path)
		err := os.MkdirAll(filepath.Dir(tp), 0755)
		if err != nil {
			return errors.Errorf("generate template %v fail:%v", tp, err.Error())
		}
		if err := os.WriteFile(tp, []byte(v.Data), 0644); err != nil {
			return errors.Errorf("generate template %v fail:%v", tp, err.Error())
		}
	}
	return nil
}
