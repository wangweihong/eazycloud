package version130

import (
	"encoding/base64"
	"os"
	osexec "os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/executil"
	"github.com/wangweihong/gotoolbox/pkg/json"
	"github.com/wangweihong/gotoolbox/pkg/log"
	"github.com/wangweihong/gotoolbox/pkg/stringutil"
	"github.com/wangweihong/gotoolbox/pkg/systemctl"

	"github.com/wangweihong/eazycloud/apis/ikubeagent"
	"github.com/wangweihong/eazycloud/internal/pkg/run"
)

// TODO: 扩展kubeagent功能, 能够调用nvidia-smi查询gpu功能

// master and worker node prepare kubelet/containerd service
func (i *impl) deployPrepareCommon(config *ikubeagent.KubernetesDeployConfig) error {
	sc, err := initSetupContext(config, false)
	if err != nil {
		return err
	}

	if err := i.GenerateTemplate(); err != nil {
		return err
	}

	if err := updateKubeletService(sc); err != nil {
		return err
	}

	if err := updateContainerdService(config); err != nil {
		return err
	}

	return nil
}

// func updateKubeletService() error {
// 	if err := os.MkdirAll(ikubeagent.KubeletRootDir, 0755); err != nil {
// 		log.Errorf("MkdirAll %v  fail:%v", ikubeagent.KubeletRootDir, err)
// 		return err
// 	}

// 	if err := generateKubeletServiceConfig(); err != nil {
// 		log.Errorf("generateKubeletServiceConfig fail:%v", err)
// 		return err
// 	}

// 	if err := generateKubeletPreRunScript(); err != nil {
// 		log.Errorf("generateKubeletPreRunScript fail:%v", err)
// 		return err
// 	}

// 	ctx := map[string]any{
// 		"KubeletBinaryPath":       ikubeagent.KubeletBinary,
// 		"KubeletPreRunScriptPath": ikubeagent.KubeletPreRunScript,
// 	}

// 	t := KubeletServiceTemplate
// 	if err := t.SetContexts(ctx).LocateToDisk().Error(); err != nil {
// 		log.Errorf("parseTemplate %v fail:%v", t.Name(), err)
// 		return err
// 	}

// 	if err := systemctl.NewCommand().Restart("kubelet", true); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func generateKubeletServiceConfig() error {
// 	if err := KubeletServiceConfigTemplate.LocateToDisk().Error(); err != nil {
// 		return err
// 	}
// 	return nil
// }

func generateKubeletPreRunScript() error {
	ctx := map[string]any{}
	t := KubeletPreRunScriptTemplate
	if err := t.SetContexts(ctx).SetFileMode(0755).LocateToDisk().Error(); err != nil {
		return err
	}
	return nil
}

func genenrateKubeletService(ctx map[string]any) error {
	t := KubeletServiceTemplate
	if err := t.SetContexts(ctx).LocateToDisk().Error(); err != nil {
		log.Errorf("parseTemplate %v fail:%v", t.Name(), err)
		return err
	}
	return nil
}

func updateContainerdService(config *ikubeagent.KubernetesDeployConfig) error {
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

	registry := config.RegistryConfig
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
	if err := ContainerdConfigTemplate.SetContexts(ctx).LocateToDisk().Error(); err != nil {
		return err
	}

	if err := systemctl.NewCommand().Restart("containerd", false); err != nil {
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

		if err := systemctl.NewCommand().Restart("containerd", false); err != nil {
			return err
		}
	} else {
		log.Infof("nvidia-ctk not found, ignore")
	}

	return nil
}

func genRegistryTlsConfig(
	certDir string,
	registry *ikubeagent.Registry,
) (cacertPath string, keyPath string, certPath string, err error) {
	if registry.TlsConfig != nil {
		address := TrimScheme(registry.Address)
		dirPath := filepath.Join(certDir, address)
		if err = os.MkdirAll(dirPath, 0755); err != nil {
			return
		}

		cacertPath = filepath.Join(dirPath, "ca.crt")
		if err = os.WriteFile(cacertPath, []byte(registry.TlsConfig.CaData), 0644); err != nil {
			return
		}
		if registry.TlsConfig.ClientKey != "" {
			keyPath = filepath.Join(dirPath, address+"."+"key")
			if err = os.WriteFile(keyPath, []byte(registry.TlsConfig.ClientKey), 0644); err != nil {
				return
			}
		}
		if registry.TlsConfig.ClientCert != "" {
			certPath = filepath.Join(dirPath, address+"."+"cert")
			if err = os.WriteFile(certPath, []byte(registry.TlsConfig.ClientCert), 0644); err != nil {
				return
			}
		}
		return
	}
	return
}

func configContainerdRegistryAuth(registry *ikubeagent.Registry) error {
	if registry.TlsConfig != nil {
		caPath, keyPath, certPath, err := genRegistryTlsConfig(certDir, registry)
		if err != nil {
			return err
		}
		if registry.User != "" && registry.Password != "" {
			authConfig := ikubeagent.DockerRegistryAccount{
				ServerAddress: registry.Address,
				User:          registry.User,
				Password:      registry.Password,
			}

			ctx := map[string]any{
				"RegistryAddress":      registry.Address,
				"RegistryCAPath":       caPath,
				"RegistryEnableVerify": registry.SkipTlsVerify,
				"RegistryCertPath":     certPath,
				"RegistryKeyPath":      keyPath,
				"RegistryAuth":         encodeAuth(&authConfig),
			}
			t := ContainerdRegistryConfigTemplate
			t.FilePath = filepath.Join(certDir, TrimScheme(registry.Address), "hosts.toml")
			if err := t.SetContexts(ctx).SetFileMode(0755).LocateToDisk().Error(); err != nil {
				log.Errorf("generate Template %v fail:%v", t.Name(), err)
				return err
			}

			return nil
		}
	}

	log.Info("registry User and Password is empty, maybe using public project")
	return nil
}

func encodeAuth(authConfig *ikubeagent.DockerRegistryAccount) string {
	if authConfig.User == "" && authConfig.Password == "" {
		return ""
	}

	authStr := authConfig.User + ":" + authConfig.Password
	msg := []byte(authStr)
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(msg)))
	base64.StdEncoding.Encode(encoded, msg)
	return string(encoded)
}

func encodeRegistryAccount(authConfig *ikubeagent.DockerRegistryAccount) string {
	b, _ := json.Marshal(authConfig)
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
	base64.StdEncoding.Encode(encoded, b)
	return string(encoded)
}

func getImageRepository(config *ikubeagent.KubernetesDeployConfig, PluginType string) (repo string) {
	if config != nil && config.RegistryConfig != nil {
		repo = TrimScheme(config.RegistryConfig.Address)
		if config.RegistryConfig.Project != "" {
			repo = repo + "/" + config.RegistryConfig.Project
		}
	}
	if config != nil {
		switch PluginType {
		case "network":
		case "monitor":
		case "application":
		case "kubectl":
		}
	}
	return
}

func TrimScheme(str string) string {
	return stringutil.TrimAnyPrefix(str, "https://", "http://")
}

func LogOperation(name string, fn func(c *setupContext) error) func(c *setupContext) error {
	return func(c *setupContext) (err error) {
		start := time.Now()
		log.Infof("start to run [%s]\n", name)

		defer func() {
			duration := time.Since(start)
			if r := recover(); r != nil {
				err = errors.Errorf(" caught panic: %v (cost: %v)", r, duration.Round(time.Millisecond))
			}
			if err != nil {
				log.Infof("[%s] run fail : %v (cost: %v)\n", name, err, duration.Round(time.Millisecond))
			} else {
				log.Infof("[%s] run success (cost: %v)\n", name, duration.Round(time.Millisecond))
			}
		}()

		err = fn(c)
		return
	}

}

// func LogOperation2(name string,fn func(c *setupContext) error) func(func(c *setupContext) error) func(c *setupContext) error {
// 	return func(fn func(c *setupContext) error) func(c *setupContext) error {
// 		return func(c *setupContext) (err error) {
// 			start := time.Now()
// 			log.Infof("🚀 开始执行 [%s]\n", name)

// 			defer func() {
// 				duration := time.Since(start)
// 				if r := recover(); r != nil {
// 					err = errors.Errorf("❌ 执行崩溃: %v (耗时: %v)", r, duration.Round(time.Millisecond))
// 				}
// 				if err != nil {
// 					log.Infof("🔥 [%s] 执行失败: %v (耗时: %v)\n", name, err, duration.Round(time.Millisecond))
// 				} else {
// 					log.Infof("✅ [%s] 执行成功 (耗时: %v)\n", name, duration.Round(time.Millisecond))
// 				}
// 			}()

// 			err = fn(c)
// 			return
// 		}
// 	}
// }
