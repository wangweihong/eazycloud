package iapiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/wangweihong/eazycloud/pkg/httpform"
	"github.com/wangweihong/gotoolbox/pkg/maputil"
)

type (
	IdentityProviderMetadataUpsetRequest struct {
		// idp的访问端点
		Endpoint          string `json:"endpoint" form:"endpoint" binding:"required"`
		AuthnNameIDFormat string `json:"authn_name_id_format" form:"authn_name_id_format"`
		// 告知sp在sso时应该跳转到idp哪个uri
		RedirectSSOFrontendURL string `json:"redirect_sso_frontend_url" form:"redirect_sso_frontend_url" binding:"required"`
		KeyEncode              []byte `json:"-"`
		CertEncode             []byte `json:"-"`
	}

	IdentityProviderMetadataUpsetResponse struct {
	}
)

func (r *IdentityProviderMetadataUpsetRequest) Decode(c *gin.Context) error {
	if err := c.ShouldBind(r); err != nil {
		return err
	}

	_, keybuf, err := httpform.FormUploadFileKey(c, httpform.FileLimitSize, httpform.KeyFileFormKey)
	if err != nil {
		return err
	}

	_, certbuf, err := httpform.FormUploadFileKey(c, httpform.FileLimitSize, httpform.CertFileFormKey)
	if err != nil {
		return err
	}

	r.KeyEncode = keybuf.Bytes()
	r.CertEncode = certbuf.Bytes()
	return nil
}

type (
	IdentityProviderMetadataGetResponse struct {
		Setting    *Setting `json:"setting"`
		DecodeKey  string   `json:"decode_key"`
		DecodeCert string   `json:"decode_cert"`
		XML        string   `json:"xml"`
	}
)

type (
	ServiceProviderMetadataUpsetRequest struct {
		KeyEncode  []byte `json:"-"`
		CertEncode []byte `json:"-"`
	}
)

type (
	ServiceProviderMetadataGetResponse struct {
		Setting    *Setting `json:"setting"`
		DecodeKey  string   `json:"decode_key"`
		DecodeCert string   `json:"decode_cert"`
	}
)

type (
	IdentityProviderMetadata struct {
		*Setting
	}
)

func (m IdentityProviderMetadata) GetKey() string {
	return maputil.TypedGet[string, string](m.Extend, "key")
}

func (m IdentityProviderMetadata) GetCert() string {
	return maputil.TypedGet[string, string](m.Extend, "cert")
}

func (m IdentityProviderMetadata) GetEndpoint() string {
	return maputil.TypedGet[string, string](m.Extend, "endpoint")
}
