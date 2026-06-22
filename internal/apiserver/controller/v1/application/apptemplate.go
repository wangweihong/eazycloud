package application

import (
	"github.com/gin-gonic/gin"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/pkg/core"
	"github.com/wangweihong/eazycloud/pkg/httpform"
)

func (rc *ApplicationController) ApplicationTemplateValidate(c *gin.Context) {
	_, buf, err := httpform.FormUpload(c, httpform.FileLimitSize)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	if err := rc.srv.Applications().
		ApplicationTemplateValidate(c, &iapiserver.ApplicationTemplateValidateRequest{PackageData: buf.Bytes()}); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, nil)
}

func (rc *ApplicationController) ApplicationTemplateAdd(c *gin.Context) {
	req := &iapiserver.ApplicationTemplateAddRequest{}
	if err := c.ShouldBind(req); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	_, buf, err := httpform.FormUpload(c, httpform.FileLimitSize)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	req.PackageData = buf.Bytes()

	if _, err := rc.srv.Applications().ApplicationTemplateAdd(c, req); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, nil)
}

func (rc *ApplicationController) ApplicationTemplateList(c *gin.Context) {
	core.Run(c, &iapiserver.ApplicationTemplateListRequest{}, func(r *iapiserver.ApplicationTemplateListRequest) (any, error) {
		ret, err := rc.srv.Applications().ApplicationTemplateList(c, r)
		return ret, err
	})
}

func (rc *ApplicationController) ApplicationTemplateGet(c *gin.Context) {
	core.Run(c, &iapiserver.ApplicationTemplateGetRequest{}, func(r *iapiserver.ApplicationTemplateGetRequest) (any, error) {
		ret, err := rc.srv.Applications().ApplicationTemplateGet(c, r)
		return ret, err
	})
}

func (rc *ApplicationController) ApplicationTemplateDelete(c *gin.Context) {
	core.Run(c, &iapiserver.ApplicationTemplateDeleteRequest{}, func(r *iapiserver.ApplicationTemplateDeleteRequest) (any, error) {
		err := rc.srv.Applications().ApplicationTemplateDelete(c, r)
		return nil, err
	})
}

func (rc *ApplicationController) ApplicationTemplateUpdate(c *gin.Context) {
	core.Run(c, &iapiserver.ApplicationTemplateUpdateRequest{}, func(r *iapiserver.ApplicationTemplateUpdateRequest) (any, error) {
		err := rc.srv.Applications().ApplicationTemplateUpdate(c, r)
		return nil, err
	})
}

func (rc *ApplicationController) ApplicationTemplateVersionList(c *gin.Context) {
	core.Run(c, &iapiserver.ApplicationTemplateVersionListRequest{}, func(r *iapiserver.ApplicationTemplateVersionListRequest) (any, error) {
		ret, err := rc.srv.Applications().ApplicationTemplateVersionList(c, r)
		return ret, err
	})
}

func (rc *ApplicationController) ApplicationTemplateVersionGet(c *gin.Context) {
	core.Run(c, &iapiserver.ApplicationTemplateVersionGetRequest{}, func(r *iapiserver.ApplicationTemplateVersionGetRequest) (any, error) {
		ret, err := rc.srv.Applications().ApplicationTemplateVersionGet(c, r)
		return ret, err
	})
}

func (rc *ApplicationController) ApplicationTemplateVersionDelete(c *gin.Context) {
	core.Run(c, &iapiserver.ApplicationTemplateVersionDeleteRequest{}, func(r *iapiserver.ApplicationTemplateVersionDeleteRequest) (any, error) {
		err := rc.srv.Applications().ApplicationTemplateVersionDelete(c, r)
		return nil, err
	})
}

func (rc *ApplicationController) ApplicationTemplateVersionUpdate(c *gin.Context) {
	core.Run(c, &iapiserver.ApplicationTemplateVersionUpdateRequest{}, func(r *iapiserver.ApplicationTemplateVersionUpdateRequest) (any, error) {
		err := rc.srv.Applications().ApplicationTemplateVersionUpdate(c, r)
		return nil, err
	})
}

func (rc *ApplicationController) ApplicationTemplateVersionAdd(c *gin.Context) {
	core.Run(c, &iapiserver.ApplicationTemplateVersionAddRequest{}, func(r *iapiserver.ApplicationTemplateVersionAddRequest) (any, error) {
		ret, err := rc.srv.Applications().ApplicationTemplateVersionAdd(c, r)
		return ret, err
	})
}
