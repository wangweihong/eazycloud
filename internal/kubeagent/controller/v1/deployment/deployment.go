package deployment

import (
	"os"
	"strings"

	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/executil"

	deploySrv "github.com/wangweihong/eazycloud/internal/kubeagent/service/v1/deployment"
	"github.com/wangweihong/eazycloud/internal/kubeagent/service/v1/deployment/version130"
	"github.com/wangweihong/eazycloud/internal/kubeagent/store"
)

// 除了部署kubernetes集群
// kubeagent还可以用来提供部署算力模型。
// 部署镜像仓库
// 推送离线包
// 创建预设置AI模型
// 填充预设置AI数据集
// ebedding模型
// 下载特定的数据模型（从
// 提供模型列表，数据集列表

type DeploymentController struct {
	srv deploySrv.AgentService
}

// NewDeploymentController creates a deployment service handler.
func NewDeploymentController(store store.Factory) *DeploymentController {
	return &DeploymentController{
		srv: newAgentService(store),
	}
}

func newAgentService(store store.Factory) deploySrv.AgentService {
	kubedamPath := "/usr/bin/kubeadm"

	if _, err := os.Stat(kubedamPath); err != nil && os.IsNotExist(err) {
		return deploySrv.NewInvalidKubeadm(err)
	}

	output, err := executil.ExecuteTimeout("kubeadm", []string{"version", "-o", "short"}, 10)
	if err != nil {
		return deploySrv.NewInvalidKubeadm(err)
	}
	version := strings.TrimSuffix(output, "\n")
	switch version {
	case "v1.30.0":
		return version130.NewDeployService(version, store)
		// case "v1.18.0":
		// 	return version118.NewDeployService(version)
	}
	err = errors.Errorf("unsupport kubeadm version:%v", version)
	return deploySrv.NewInvalidKubeadm(err)
}
