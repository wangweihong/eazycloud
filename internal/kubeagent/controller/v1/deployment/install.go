package deployment

import (
	"github.com/gin-gonic/gin"
	"github.com/wangweihong/eazycloud/apis/ikubeagent"
	"github.com/wangweihong/eazycloud/apis/imachinery"
	"github.com/wangweihong/eazycloud/pkg/core"
)

func (d *DeploymentController) InstallMaster(c *gin.Context) {
	core.Run(c, &ikubeagent.InstallMasterRequest{}, func(r *ikubeagent.InstallMasterRequest) (any, error) {
		return nil, d.srv.InitMaster(c, r)
	})
}

func (d *DeploymentController) CheckDependency(c *gin.Context) {
	core.Run(c, &ikubeagent.CheckDependencyRequest{}, func(r *ikubeagent.CheckDependencyRequest) (any, error) {
		return nil, d.srv.CheckDependency(c, r.Version)
	})
}

func (d *DeploymentController) JoinCluster(c *gin.Context) {
	core.Run(c, &ikubeagent.JoinClusterRequest{}, func(r *ikubeagent.JoinClusterRequest) (any, error) {
		return nil, d.srv.JoinCluster(c, r)
	})
}

func (d *DeploymentController) ResetDeploy(c *gin.Context) {
	core.Run(c, &imachinery.Empty{}, func(r *imachinery.Empty) (any, error) {
		return nil, d.srv.Reset(c)
	})
}

func (d *DeploymentController) GetJoinCommand(c *gin.Context) {
	core.Run(c, &ikubeagent.GetJoinCommandRequest{}, func(r *ikubeagent.GetJoinCommandRequest) (any, error) {
		joinCommand, kubeadmYamlData, err := d.srv.GetJoinCommand(c, r)
		if err != nil {
			return nil, err
		}
		return &ikubeagent.GetJoinCommandResponse{
			JoinCommand:     joinCommand,
			KubeadmYamlData: kubeadmYamlData,
		}, nil
	})
}

func (d *DeploymentController) GetKubeConfig(c *gin.Context) {
	core.Run(c, &ikubeagent.GetKubeConfigRequest{}, func(r *ikubeagent.GetKubeConfigRequest) (any, error) {
		data, err := d.srv.GetKubeConfig(c)
		if err != nil {
			return nil, err
		}
		return &ikubeagent.GetKubeConfigResponse{Data: data}, nil
	})
}

func (d *DeploymentController) GetInstallLog(c *gin.Context) {
	core.Run(c, &ikubeagent.GetInstallLogRequest{}, func(r *ikubeagent.GetInstallLogRequest) (any, error) {
		data, err := d.srv.GetDeployLog(c)
		if err != nil {
			return nil, err
		}
		return &ikubeagent.GetInstallLogResp{Data: data}, nil
	})
}

func (d *DeploymentController) GetInstallState(c *gin.Context) {
	core.Run(c, &ikubeagent.InstallStateRequest{}, func(r *ikubeagent.InstallStateRequest) (any, error) {
		return d.srv.GetInstallState(c)
	})
}

func (d *DeploymentController) SetOwner(c *gin.Context) {
	core.Run(c, &ikubeagent.InstallStateRequest{}, func(r *ikubeagent.InstallStateRequest) (any, error) {
		return d.srv.GetInstallState(c)
	})
}