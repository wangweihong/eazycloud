package iapiserver

import (
	"time"

	"github.com/wangweihong/eazycloud/apis/imachinery"
)

const DefaultCertificateExpireDuration = time.Hour * 24 * 365 * 10

const (
	SettingKindSSOSamlIdpMetadata = "saml_idp"
	SettingKindSSOSamlSpMetadata  = "saml_sp"
)


type Setting struct {
	imachinery.ObjectMeta
}
