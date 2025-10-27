package deployment

import (
	"github.com/gin-gonic/gin"

	"github.com/wangweihong/eazycloud/apis/ikubeagent"
	"github.com/wangweihong/eazycloud/pkg/core"
)

func run[T any](c *gin.Context, req T, action func(r T) (any, error)) {
	if err := core.DecodeParameter(c, req); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	ret, err := action(req)
	core.WriteResponse(c, err, ret)
}

func (d *DeploymentController) DeployMaster(c *gin.Context) {
	// run(c, &iregistry.ListRegistryParam{}, func(r *iregistry.ListRegistryParam) (any, error) {
	// 	ret, err := rc.srv.Registries().List(c, imachinery.ListOptions{})
	// 	return ret, err
	// })
}

func (d *DeploymentController) CheckDependency(c *gin.Context) {
	run(c, &ikubeagent.CheckDependencyReq{}, func(r *ikubeagent.CheckDependencyReq) (any, error) {
		return nil, d.srv.CheckDependency(r.Version)
	})

}

func (d *DeploymentController) JoinCluster(c *gin.Context) {

}

func (d *DeploymentController) ResetDeploy(c *gin.Context) {

}

func (d *DeploymentController) GetJoinCommand(c *gin.Context) {

}

func (d *DeploymentController) GetKubeConfig(c *gin.Context) {

}

func (d *DeploymentController) GetDeployLog(c *gin.Context) {

}
