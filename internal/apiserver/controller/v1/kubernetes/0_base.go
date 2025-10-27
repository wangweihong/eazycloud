package kubernetes

import (
	"github.com/gin-gonic/gin"

	srvv1 "github.com/wangweihong/eazycloud/internal/apiserver/service/v1"
	"github.com/wangweihong/eazycloud/internal/apiserver/store"
	"github.com/wangweihong/eazycloud/pkg/core"
)

type KubernetesController struct {
	srv srvv1.Service
}

// NewKubernetesController creates a registry service handler.
func NewKubernetesController(store store.Factory) *KubernetesController {
	return &KubernetesController{
		srv: srvv1.NewService(store),
	}
}

func run[T any](c *gin.Context, req T, action func(r T) (any, error)) {
	if err := core.DecodeParameter(c, req); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	ret, err := action(req)
	core.WriteResponse(c, err, ret)
}

func runKubernetes[T any](c *gin.Context, req T, action func(r T) (any, error)) {
	if err := core.DecodeParameter(c, req); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	ret, err := action(req)
	core.WriteResponse(c, err, ret)
}
