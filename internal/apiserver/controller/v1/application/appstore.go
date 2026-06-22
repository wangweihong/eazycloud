package application

import (
	"github.com/gin-gonic/gin"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/pkg/core"
)

func (rc *ApplicationController) AppStoreList(c *gin.Context) {
	core.Run(c, &iapiserver.AppStoreListRequest{}, func(r *iapiserver.AppStoreListRequest) (any, error) {
		ret, err := rc.srv.Applications().AppStoreQuery(c, r)
		return ret, err
	})
}

func (rc *ApplicationController) AppStoreGet(c *gin.Context) {
	core.Run(c, &iapiserver.AppStoreGetRequest{}, func(r *iapiserver.AppStoreGetRequest) (any, error) {
		ret, err := rc.srv.Applications().AppStoreGet(c, r)
		return ret, err
	})
}

func (rc *ApplicationController) AppStoreAdd(c *gin.Context) {
	core.Run(c, &iapiserver.AppStoreAddRequest{}, func(r *iapiserver.AppStoreAddRequest) (any, error) {
		ret, err := rc.srv.Applications().AppStoreAdd(c, r)
		return ret, err
	})
}

func (rc *ApplicationController) AppStoreUpdate(c *gin.Context) {
	core.Run(c, &iapiserver.AppStoreUpdateRequest{}, func(r *iapiserver.AppStoreUpdateRequest) (any, error) {
		err := rc.srv.Applications().AppStoreUpdate(c, r)
		return nil, err
	})
}

func (rc *ApplicationController) AppStoreDelete(c *gin.Context) {
	core.Run(c, &iapiserver.AppStoreDeleteRequest{}, func(r *iapiserver.AppStoreDeleteRequest) (any, error) {
		err := rc.srv.Applications().AppStoreDelete(c, r)
		return nil, err
	})
}

func (rc *ApplicationController) AppStoreSync(c *gin.Context) {
	core.Run(c, &iapiserver.AppStoreSyncRequest{}, func(r *iapiserver.AppStoreSyncRequest) (any, error) {
		err := rc.srv.Applications().AppStoreSync(c, r)
		return nil, err
	})
}
