package kubernetes

import (
	"github.com/gin-gonic/gin"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
)

func (rc *KubernetesController) PodCreate(c *gin.Context) {
	runKubernetes(c, &iapiserver.PodRequest{}, func(r *iapiserver.PodRequest) (any, error) {
		return rc.srv.Kubernetes().PodCreate(c, r)
	})
}

func (rc *KubernetesController) PodGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.ResourceGetRequest{}, func(r *iapiserver.ResourceGetRequest) (any, error) {
		return rc.srv.Kubernetes().PodGet(c, r)
	})
}
