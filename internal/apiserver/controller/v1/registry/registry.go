package registry

import (
	"github.com/gin-gonic/gin"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/apis/imachinery"
	"github.com/wangweihong/eazycloud/apis/iregistry"
	"github.com/wangweihong/eazycloud/pkg/core"
)

func (rc *RegistryController) List(c *gin.Context) {
	run(c, &iapiserver.RegistryListRequest{}, func(r *iapiserver.RegistryListRequest) (any, error) {
		ret, err := rc.srv.Registries().List(c, r)
		return ret, err
	})
}

func (rc *RegistryController) Get(c *gin.Context) {
	run(c, &iregistry.GetRegistryParam{}, func(r *iregistry.GetRegistryParam) (any, error) {
		ret, err := rc.srv.Registries().Get(c, r.ID, imachinery.GetOptions{})
		return ret, err
	})
}

func (rc *RegistryController) Add(c *gin.Context) {
	run(c, &iregistry.Registry{}, func(r *iregistry.Registry) (any, error) {
		ret, err := rc.srv.Registries().Add(c, r, imachinery.CreateOptions{})
		return ret, err
	})
}

func (rc *RegistryController) Delete(c *gin.Context) {
	run(c, &iregistry.DeleteRegistryParam{}, func(r *iregistry.DeleteRegistryParam) (any, error) {
		err := rc.srv.Registries().Delete(c, r.ID, imachinery.DeleteOptions{})
		return nil, err
	})
}

func (rc *RegistryController) Update(c *gin.Context) {
	run(c, &iregistry.Registry{}, func(r *iregistry.Registry) (any, error) {
		ret, err := rc.srv.Registries().Update(c, r, imachinery.UpdateOptions{})
		return ret, err
	})
}

func (rc *RegistryController) Search(c *gin.Context) {
	var r iregistry.SearchRegistryParam
	if err := core.DecodeParameter(c, &r); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	// if err := rc.srv.Registries().Search(c, r.ID, imachinery.DeleteOptions{}); err != nil {
	// 	core.WriteResponse(c, err, nil)
	// 	return
	// }
}
