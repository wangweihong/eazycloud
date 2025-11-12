package deployment

import (
	"github.com/wangweihong/eazycloud/internal/kubeagent/service/v1/deployment"
	deploySrv "github.com/wangweihong/eazycloud/internal/kubeagent/service/v1/deployment"

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
		srv: deployment.NewService(store),
	}
}
