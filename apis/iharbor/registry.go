package iharbor

import "time"

const (
	RegistryTypeHarbor = "harbor"
)

type RegistryInfo struct {
	Type                     string `json:"type"`
	Description              string `json:"description"`
	SupportedResourceFilters []struct {
		Type   string   `json:"type"`
		Style  string   `json:"style"`
		Values []string `json:"values,omitempty"`
	} `json:"supported_resource_filters"`
	SupportedTriggers []string `json:"supported_triggers"`
}

//RegistryPing 只需要传type/url/即可
type RegistryPingRquest struct {
	Credential   RegistryCredential `json:"credential"`
	AccessKey    *string            `json:"access_key"`    //账号
	AccessSecret *string            `json:"access_secret"` //密码
	Description  string             `json:"description"`
	Insecure     bool               `json:"insecure"`
	Name         string             `json:"name"`
	Type         string             `json:"type"` //类型 harbor
	URL          string             `json:"url"`
}

func (r *RegistryPingRquest) Validate() error {
	if r.AccessKey == nil || r.AccessSecret == nil {
		r.AccessKey = r.Credential.AccessKey
		r.AccessSecret = r.Credential.AccessSecret
	}
	return nil
}

type Registry struct {
	ID              int                `json:"id" form:"id"`
	Name            string             `json:"name,omitempty"`
	Description     string             `json:"description"`
	Type            string             `json:"type,omitempty"` //从/api/replication/adapters获取中
	URL             string             `json:"url,omitempty"`
	TokenServiceURL string             `json:"token_service_url,omitempty"`
	Credential      RegistryCredential `json:"credential"`
	Insecure        bool               `json:"insecure"`
	Status          string             `json:"status"` //healthy
	CreationTime    *time.Time         `json:"creation_time"`
	UpdateTime      *time.Time         `json:"update_time"`

	AccessKey    *string `json:"access_key,omitempty"`    //账号
	AccessSecret *string `json:"access_secret,omitempty"` //密码
}

type RegistryCredential struct {
	Type         *string `json:"type,omitempty"`          //basic等
	AccessKey    *string `json:"access_key,omitempty"`    //账号
	AccessSecret *string `json:"access_secret,omitempty"` //密码
}

type RegistryListResponse struct {
	List       []Registry `json:"list"`
	TotalCount int        `json:"total_count"`
}
