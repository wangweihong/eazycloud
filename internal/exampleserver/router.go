package example_server

import (
	"github.com/gin-gonic/gin"

	"github.com/wangweihong/eazycloud/internal/pkg/genericserver/genericmiddleware"
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
	return g
}
