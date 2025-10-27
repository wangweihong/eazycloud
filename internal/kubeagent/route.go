package kubeagent

import (
	"github.com/gin-gonic/gin"

	"github.com/wangweihong/eazycloud/internal/apiserver/store/postgresql"
	"github.com/wangweihong/eazycloud/internal/kubeagent/controller/v1/deployment"
	"github.com/wangweihong/eazycloud/pkg/httpsvr/genericmiddleware"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installApis(g)
}

func installMiddleware(g *gin.Engine) {
	g.Use(genericmiddleware.RequestID())
	g.Use(genericmiddleware.Context())
	g.Use(genericmiddleware.LoggerMiddleware())
}

func installApis(g *gin.Engine) *gin.Engine {

	//g.NoRoute(func(c *gin.Context) {
	//	core.WriteResponse(c, errors.WithCode(code.ErrPageNotFound, "Page not found."), nil)
	//})
	storeIns, _ := postgresql.GetPostgresSQLFactoryOr(nil)
	v1 := g.Group("/v1")
	{
		deploymentv1 := v1.Group("/deployment")
		{
			deploymentController := deployment.NewDeploymentController(storeIns)

			deploymentv1.POST("/deploy-master", deploymentController.DeployMaster)
			deploymentv1.POST("/join-cluster", deploymentController.JoinCluster)
			deploymentv1.POST("/reset-dpeloy", deploymentController.ResetDeploy)
			deploymentv1.GET("/join-command", deploymentController.GetJoinCommand)
			deploymentv1.GET("/kube-config", deploymentController.GetKubeConfig)
			deploymentv1.GET("/deploy-log", deploymentController.GetDeployLog)
			deploymentv1.GET("/deploy-state", deploymentController.GetDeploymentState)
			deploymentv1.GET("/check-node-dep", deploymentController.CheckNodeDependency)

		}
	}

	return g
}
