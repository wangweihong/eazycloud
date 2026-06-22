package iharbor

import (
	"strings"

	"github.com/wangweihong/gotoolbox/pkg/errors"
)

type PreheatInstance struct {
	Insecure bool `json:"insecure"`

	AuthMode       string                   `json:"auth_mode"` //BASIC,NONE,OAUTH
	AuthInfo       *PreheatInstanceAuthInfo `json:"auth_info,omitempty"`
	Description    string                   `json:"description,omitempty"`
	Endpoint       string                   `json:"endpoint"`
	ID             int                      `json:"id"`
	Name           string                   `json:"name"`
	SetupTimestamp int                      `json:"setup_timestamp"`
	Status         string                   `json:"status"`
	StatusMessage  string                   `json:"status_message"`
	Vendor         string                   `json:"vendor"`
	Enabled        bool                     `json:"enabled"`
}

type PreheatInstanceAuthInfo struct {
	Token    string `json:"token"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

type PreheatInstanceProvider struct {
	Icon        string   `json:"icon"`
	ID          string   `json:"id"`
	Maintainers []string `json:"maintainers"`
	Name        string   `json:"name"`
	Source      string   `json:"source"`
	Version     string   `json:"version"`
}

type PreheatInstanceProviderListRequest struct {
	PagingParam
}

func (r *PreheatInstanceProviderListRequest) Validate() error {
	r.PagingParam.Valiate()
	return nil
}

type PreheatInstanceProviderListResponse struct {
	List []PreheatInstanceProvider `json:"list"`
}

type PreheatInstanceRequest struct {
	PreheatInstance
}

func (r *PreheatInstanceRequest) Validate() error {
	if r.Name == "" {
		return errors.Errorf("preheat instance name is empty")
	}
	//provider返回产商是大写字母开头, 但创建和更新vendor不能直接用, 必须要小写。
	r.Vendor = strings.ToLower(r.Vendor)
	return nil
}

type PreheatInstanceResponse struct {
	List       []PreheatInstance `json:"list"`
	TotalCount int               `json:"total_count"`
}
