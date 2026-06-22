package kubeagent

import (
	"github.com/gin-gonic/gin"

	"github.com/wangweihong/eazycloud/internal/kubeagent/controller/v1/deployment"
	"github.com/wangweihong/eazycloud/internal/kubeagent/store/postgresql"
	"github.com/wangweihong/eazycloud/internal/pkg/code"
	"github.com/wangweihong/eazycloud/pkg/core"
	"github.com/wangweihong/eazycloud/pkg/httpsvr/genericmiddleware"
	"github.com/wangweihong/gotoolbox/pkg/errors"
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

	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errors.WithCode(code.ErrPageNotFound, "Page not found."), nil)
	})
	storeIns, _ := postgresql.GetPostgresSQLFactoryOr(nil)
	v1 := g.Group("/v1")
	{
		deploymentv1 := v1.Group("/deployment")
		{
			deploymentController := deployment.NewDeploymentController(storeIns)

			deploymentv1.POST("/install-master", deploymentController.InstallMaster)
			deploymentv1.POST("/join-cluster", deploymentController.JoinCluster)
			deploymentv1.POST("/reset-dpeloy", deploymentController.ResetDeploy)
			deploymentv1.GET("/join-command", deploymentController.GetJoinCommand)
			deploymentv1.GET("/kubeconfig", deploymentController.GetKubeConfig)
			deploymentv1.GET("/install-log", deploymentController.GetInstallLog)
			deploymentv1.GET("/install-state", deploymentController.GetInstallState)
			deploymentv1.GET("/check-dependency", deploymentController.CheckDependency)
		}
	}

	return g
}
