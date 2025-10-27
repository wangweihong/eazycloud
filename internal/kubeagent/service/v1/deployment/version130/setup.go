package version130

import (
	"os"
	osexec "os/exec"
	"path/filepath"
	"strings"

	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/executil"
	"github.com/wangweihong/gotoolbox/pkg/log"
	"github.com/wangweihong/gotoolbox/pkg/systemctl"
	"github.com/wangweihong/gotoolbox/pkg/template"

	"github.com/wangweihong/eazycloud/apis/ikubeagent"
	"github.com/wangweihong/eazycloud/internal/pkg/run"
)

type setupContext struct {
	templateDir           string
	kubeRootDir           string
	KubeadmConfigYamlPath string

	config    *ikubeagent.KubernetesDeployConfig
	isMaster0 bool
}

type setupFunc func(c *setupContext) error

func initSetupContext(dc *ikubeagent.KubernetesDeployConfig, isMaster bool) (*setupContext, error) {
	templateDir := ikubeagent.TemplateDirPath
	kubeRootDir := ikubeagent.KubeletRootDir
	KubeadmConfigYamlPath := ikubeagent.KubeadmConfigYamlPath

	return &setupContext{
		config:                dc,
		isMaster0:             isMaster,
		templateDir:           templateDir,
		kubeRootDir:           kubeRootDir,
		KubeadmConfigYamlPath: KubeadmConfigYamlPath,
	}, nil
}

func writeFile(tp string, data []byte) error {
	err := os.MkdirAll(filepath.Dir(tp), 0755)
	if err != nil {
		return errors.Errorf("mkdir dir %v fail:%v", filepath.Dir(tp), err.Error())
	}
	if err := os.WriteFile(tp, data, 0644); err != nil {
		return errors.Errorf("write file %v fail:%v", tp, err.Error())
	}
	return nil
}

func prepareDataDir(c *setupContext) error {
	log.Infof("prepare data dir")

	remainCleanDir := []string{
		c.templateDir,
	}

	dirNeedToCreate := []string{
		c.kubeRootDir,
	}

	for _, dir := range remainCleanDir {
		log.Infof("rmdir %v", dir)
		_ = os.RemoveAll(dir)
	}

	for _, dir := range dirNeedToCreate {
		log.Infof("mkdir %v", dir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Errorf("mkdir  %v  fail:%v", dir, err)
			return err
		}
	}

	log.Infof("prepare data dir success")

	return nil
}

func prepareComponentTemplate(c *setupContext) error {
	log.Infof("prepare component template under %v ..", c.templateDir)

	for _, v := range fps {
		tp := filepath.Join(c.templateDir, v.Path)
		log.Infof("write file %v start", tp)
		if err := writeFile(tp, []byte(v.Data)); err != nil {
			log.Errorf("write file  %v error: %v", tp, err)
			return err
		}
	}

	log.Infof("prepare component template success")
	return nil
}

type templateFileObject struct {
	tp       template.FileProcessor
	ctx      map[string]any
	fileMode os.FileMode
}

func updateKubeletService(c *setupContext) error {
	ctx := map[string]any{
		"KubeletBinaryPath":       ikubeagent.KubeletBinary,
		"KubeletPreRunScriptPath": ikubeagent.KubeletPreRunScript,
	}
	var tfos = []templateFileObject{
		{KubeadmConfigTemplate, nil, 0644},
		{KubeletPreRunScriptTemplate, nil, 0755},
		{KubeletServiceTemplate, ctx, 0755},
	}

	for _, v := range tfos {
		if err := run.WriteTemplate(v.tp, v.ctx, v.fileMode); err != nil {
			return err
		}
	}

	if err := restartKubeletService(c); err != nil {
		return err
	}
	return nil
}

func restartKubeletService(c *setupContext) error {
	return restartSystemdService("kubelet")
}

func restartSystemdService(service string) error {
	log.Infof("restart %v service", service)

	// TODO： check service state
	// if noactive, use start command
	if err := systemctl.NewCommand().Restart(service, true); err != nil {
		return err
	}

	log.Infof("restart %v service success", service)

	return nil
}

func prepareHA(c *setupContext) error {
	log.Infof("prepare HA")

	// ignore non-control-plane and ha
	if c.config.ControlPlaneConfig == nil || c.config.ControlPlaneConfig.HAConfig == nil {
		log.Infof("ignore ha Prepare phase form no ha config")
		return nil
	}

	for _, setup := range []setupFunc{generateKeepalivedConfig, generateKeepalivedCheckerScript, generateKeepalivedService} {
		if err := setup(c); err != nil {
			return nil
		}
	}

	if err := generateHaproxyService(c.config); err != nil {
		return err
	}

	_ = os.MkdirAll(filepath.Dir(ikubeagent.KubeadmConfigYamlPath), 0755)
	if err := os.WriteFile(ikubeagent.KubeadmConfigYamlPath, []byte(c.config.ControlPlaneConfig.HAConfig.KubeadmConfigYaml), 0755); err != nil {
		return err
	}

	log.Infof("run prepare work for HA complete")
	return nil
}

// func updateContainerdService(config *ikubeagent.KubernetesDeployConfig) error {
// 	ctx := make(map[string]any)

// 	// 根据系统的实际cgroup更新containerd cgroup配置
// 	systemdCgroup := "false"
// 	args := []string{"-fcH`", "%T", "sys/fs/cgroup/"}
// 	stdout, stderr, err := executil.ExecuteCmdSplitStdoutStderr("stat", args, 0)
// 	if err != nil || stderr != "" {
// 		log.Errorf("run command [%v:%v] fail:%v,stderr:%v", ikubeagent.KubectlBinary, args, run.TrimError(err), stderr)
// 		return run.TrimError(err)
// 	}

// 	log.Infof("cgroupfs:%v", stdout)
// 	if strings.Contains(stdout, "cgroup2fs") {
// 		systemdCgroup = "true"
// 	}
// 	ctx["SystemdCgroup"] = systemdCgroup

// 	registry := config.RegistryConfig
// 	if registry != nil {
// 		repo := TrimScheme(registry.Address)
// 		if registry.Project != "" {
// 			repo = repo + "/" + registry.Project
// 		}
// 		if err := configContainerdRegistryAuth(registry); err != nil {
// 			log.Errorf("configContainerdRegistryAuth fail:%v", err)
// 			return err
// 		}
// 		ctx["ImageRepository"] = repo
// 	}
// 	if err := ContainerdConfigTemplate.SetContexts(ctx).LocateToDisk().Error(); err != nil {
// 		return err
// 	}

// 	if err := systemctl.NewCommand().Restart("containerd", false); err != nil {
// 		return err
// 	}

// 	// FIXME? 当集群需要gpu环境, 又缺少nvidia-ctk, 是否应该报错? 同时需要兼容其他可能性 如华为npu以及混合N/GPU场景
// 	if _, err := osexec.LookPath("nvidia-ctk"); err == nil {
// 		log.Infof("nvidia-ctk found, prepare to configure")
// 		args := []string{"runtime", "configure", "--runtime=containerd"}
// 		if _, err := executil.Execute("nvidia-ctk", args); err != nil {
// 			log.Errorf("run command [%v:%v] fail:%v", "nvidia-ctk", args, run.TrimError(err))
// 			return run.TrimError(err)
// 		}

// 		if err := systemctl.NewCommand().Restart("containerd", false); err != nil {
// 			return err
// 		}
// 	} else {
// 		log.Infof("nvidia-ctk not found, ignore")
// 	}

// 	return nil
// }

// FIXME: 兼容container v3 config
func updateContainerdService2(c *setupContext) error {
	ctx := make(map[string]any)

	// 根据系统的实际cgroup更新containerd cgroup配置
	systemdCgroup := "false"
	args := []string{"-fc", "%T", "sys/fs/cgroup/"}
	stdout, stderr, err := executil.ExecuteCmdSplitStdoutStderr("stat", args, 0)
	if err != nil || stderr != "" {
		log.Errorf("run command [%v:%v] fail:%v,stderr:%v", ikubeagent.KubectlBinary, args, run.TrimError(err), stderr)
		return run.TrimError(err)
	}

	log.Infof("cgroupfs:%v", stdout)
	if strings.Contains(stdout, "cgroup2fs") {
		systemdCgroup = "true"
	}
	ctx["SystemdCgroup"] = systemdCgroup

	registry := c.config.RegistryConfig
	if registry != nil {
		repo := TrimScheme(registry.Address)
		if registry.Project != "" {
			repo = repo + "/" + registry.Project
		}
		if err := configContainerdRegistryAuth(registry); err != nil {
			log.Errorf("configContainerdRegistryAuth fail:%v", err)
			return err
		}
		ctx["ImageRepository"] = repo
	}
	if err := run.WriteTemplate(ContainerdConfigTemplate, ctx, 0644); err != nil {
		return err
	}

	// FIXME? 当集群需要gpu环境, 又缺少nvidia-ctk, 是否应该报错? 同时需要兼容其他可能性 如华为npu以及混合N/GPU场景
	if _, err := osexec.LookPath("nvidia-ctk"); err == nil {
		log.Infof("nvidia-ctk found, prepare to configure")
		args := []string{"runtime", "configure", "--runtime=containerd"}
		if _, err := executil.Execute("nvidia-ctk", args); err != nil {
			log.Errorf("run command [%v:%v] fail:%v", "nvidia-ctk", args, run.TrimError(err))
			return run.TrimError(err)
		}
	} else {
		log.Infof("nvidia-ctk not found, ignore")
	}
	if err := restartSystemdService("containerd"); err != nil {
		return err
	}
	return nil
}
