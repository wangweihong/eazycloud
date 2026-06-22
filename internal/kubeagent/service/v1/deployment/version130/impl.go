package version130

import (
	"context"
	gerrors "errors"
	"os"
	osexec "os/exec"
	"path/filepath"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/executil"
	"github.com/wangweihong/gotoolbox/pkg/log"
	"github.com/wangweihong/gotoolbox/pkg/statemachine"
	"gorm.io/gorm"

	"github.com/wangweihong/eazycloud/apis/ikubeagent"
	"github.com/wangweihong/eazycloud/apis/imachinery"
	"github.com/wangweihong/eazycloud/internal/kubeagent/store"
	"github.com/wangweihong/eazycloud/internal/pkg/run"
)

func NewDeployService(version string, state *ikubeagent.InstallState, store store.Factory) *impl {
	return &impl{
		state:   state,
		sm:      statemachine.New(statemachine.State(state.State)),
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
	sm      *statemachine.StateMachine
	state   *ikubeagent.InstallState
}

func (i *impl) Version() string {
	return i.version
}

func (i *impl) InitMaster(ctx context.Context, param *ikubeagent.InstallMasterRequest) error {
	if i.sm.CurrentState() != ikubeagent.KubernetesDeployStateUninitialized {
		return errors.Errorf("current node has used to deploy")
	}
	var err error

	defer func() {
		updateState(ctx, i.store, err, i.state, ikubeagent.KubernetesDeployStateSuccess)
	}()
	err = i.initMaster(ctx, param)
	return err
}

func (i *impl) initMaster(ctx context.Context, param *ikubeagent.InstallMasterRequest) (err error) {
	log.Infof("start to init master")

	config := param.ToKubernetesInstallConfig()

	if i.sm.CurrentState() != ikubeagent.KubernetesDeployStateUninitialized {
		return errors.Errorf("current node has used to deploy,cuurent state: %v", i.sm.CurrentState())
	}

	c, err := initSetupContext(config, true)
	if err != nil {
		return errors.WithStack(err)
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
		{Name: "install display card plugin", Op: installDisplayCardPlugin},
	}

	for _, op := range Ops {
		setupDep := LogOperation(op.Name, op.Op)
		if err := setupDep(c); err != nil {
			return err
		}
	}

	return nil
}

func updateState(ctx context.Context, store store.Factory, err error, state *ikubeagent.InstallState, targetState statemachine.State) {
	if x := recover(); x != nil {
		err = errors.Errorf("panic")
	}

	switch targetState {
	case ikubeagent.KubernetesDeployStateError, ikubeagent.KubernetesDeployStateSuccess:
		state.EndTime = imachinery.Now()
	case ikubeagent.KubernetesDeployStateUninitialized:
		state = new(ikubeagent.InstallState)
		state.State = string(ikubeagent.KubernetesDeployStateUninitialized)

	case ikubeagent.KubernetesDeployStateDeploying:
		state.StartTime = imachinery.Now()
	}

	if err != nil {
		state.State = string(ikubeagent.KubernetesDeployStateError)
		state.ErrorMessage = err.Error()
	} else {
		state.State = string(ikubeagent.KubernetesDeployStateSuccess)
	}

	if _, err := store.InstallStateStores().Upsert(ctx, state); err != nil {
		log.F(ctx).Errorf("update install state error:%v", err)
	}

}
func (i *impl) JoinCluster(ctx context.Context, param *ikubeagent.JoinClusterRequest) error {
	if i.sm.CurrentState() != ikubeagent.KubernetesDeployStateUninitialized {
		return errors.Errorf("current node has used to deploy")
	}
	var err error

	defer func() {
		updateState(ctx, i.store, err, i.state, ikubeagent.KubernetesDeployStateSuccess)
	}()
	err = i.joinCluster(ctx, param)
	return err
}

func (i *impl) joinCluster(ctx context.Context, param *ikubeagent.JoinClusterRequest) (err error) {
	log.Info("start to JoinCluster")

	config := param.ToKubernetesInstallConfig()

	c, err := initSetupContext(config, false)
	if err != nil {
		return err
	}

	Ops := []NamedOperation{
		{Name: "prepare data dir", Op: prepareDataDir},
		{Name: "prepare component template ", Op: prepareComponentTemplate},
		{Name: "update kubelet service", Op: updateKubeletService},
		{Name: "update containerd service", Op: updateContainerdService2},
		// {Name: "prepare ha", Op: deployPrepareHA},
	}
	if param.WorkerConfig.IsControlPlane {
		Ops = append(Ops, NamedOperation{Name: "prepare ha", Op: deployPrepareHA})
	}

	for _, op := range Ops {
		setupDep := LogOperation(op.Name, op.Op)
		if err := setupDep(c); err != nil {
			return err
		}
	}

	log.Info("Join Cluster run runKubeadmJoin")
	if err := runKubeadmJoin(param.WorkerConfig.JoinCommand, param.WorkerConfig.IsControlPlane, config.NodeConfig.NodeName); err != nil {
		log.Errorf("runKubeadmJoin fail:%v", err)
		return err
	}
	log.Info("JoinCluster success")
	return nil
}

func (i *impl) Reset(ctx context.Context) error {

	if i.sm.CurrentState() != ikubeagent.KubernetesDeployStateDeploying {
		return errors.Errorf("current node is deploying")
	}

	var err error
	defer func() {
		updateState(ctx, i.store, err, i.state, ikubeagent.KubernetesDeployStateUninitialized)
	}()
	err = i.reset(ctx, i.state.HighAvailable)
	return err
}

func (i *impl) reset(ctx context.Context, isHa bool) error {
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

func (i *impl) GetJoinCommand(ctx context.Context, param *ikubeagent.GetJoinCommandRequest) (string, string, error) {
	if i.sm.CurrentState() != ikubeagent.KubernetesDeployStateSuccess {
		return "", "", errors.Errorf("kubernetes cluster doesn't deploy success")
	}

	if !i.state.ControlPlane {
		return "", "", errors.Errorf("current node is not control-plane node")
	}

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

	if param.ShowControlPlane {
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

func (i *impl) GetKubeConfig(ctx context.Context) (string, error) {
	if i.sm.CurrentState() != ikubeagent.KubernetesDeployStateSuccess {
		return "", errors.Errorf("kubernetes cluster doesn't deploy success")
	}

	if !i.state.ControlPlane {
		return "", errors.Errorf("current node is not control-plane node")
	}

	data, err := os.ReadFile(ikubeagent.KubeConfigPath)
	if err != nil {
		return "", errors.Errorf("read file %v fail:%v", ikubeagent.KubeConfigPath, err.Error())
	}

	return string(data), nil
}

func (i *impl) GetDeployLog(ctx context.Context) (string, error) {
	data, err := os.ReadFile(ikubeagent.KubeadmLogPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}

		return "", errors.Errorf("read file %v fail:%v", ikubeagent.KubeadmLogPath, err.Error())
	}

	return string(data), nil
}

func (i *impl) CheckDependency(ctx context.Context, version string) error {
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

func (i *impl) GetInstallState(ctx context.Context) (*ikubeagent.InstallStateResponse, error) {
	state, err := i.store.InstallStateStores().GetByName(ctx, ikubeagent.KubernetesInstallStateUniqueName)
	if err != nil && !gerrors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.WithStack(err)
	}

	if err != nil {
		state = &ikubeagent.InstallState{
			State: string(ikubeagent.KubernetesDeployStateUninitialized),
		}
	}

	hostName, _ := os.Hostname()
	hostInfo, _ := hostInfos()
	return &ikubeagent.InstallStateResponse{
		State:    state,
		HostStat: hostInfo,
		HostName: hostName,
	}, nil
}

func hostInfos() (*ikubeagent.HostInfo, error) {
	cpucores, err := cpu.Counts(false)
	if err != nil {
		return nil, err
	}
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	return &ikubeagent.HostInfo{
		Cpu: ikubeagent.CpuInfo{
			Cores: int64(cpucores),
		},
		Mem: ikubeagent.MemInfo{
			Total: memInfo.Total,
		},
	}, nil
}

func (i *impl) SetOwner(ctx context.Context, owner string) error {
	//  fixme: 只更新单个数据
	// if _, err := store.InstallStateStores().Upsert(ctx, state); err != nil {
	// 	log.F(ctx).Errorf("update install state error:%v", err)
	// }
	return nil
}
