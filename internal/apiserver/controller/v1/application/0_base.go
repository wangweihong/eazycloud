package application

import (
	srvv1 "github.com/wangweihong/eazycloud/internal/apiserver/service/v1"
	"github.com/wangweihong/eazycloud/internal/apiserver/store"
)

type ApplicationController struct {
	srv srvv1.Service
}

func NewApplicationController(store store.Factory) *ApplicationController {
	return &ApplicationController{
		srv: srvv1.NewService(store),
	}
}
