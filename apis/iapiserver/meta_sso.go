package iapiserver

import (
	"fmt"

	"github.com/crewjam/saml/samlsp"
	"github.com/wangweihong/eazycloud/apis/imachinery"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/json"
	"gorm.io/gorm"
)

const (
	UriSSO         = "/v1/auth/sso/login"
	UriSLO         = "/v1/auth/sso/logout"
	UriSSOMetadata = "/v1/auth/sso/metadata"
)

const (
	SSOProtocolSAML   = "saml"
	SSOProtocolOAUTH2 = "oauth2"
)

type SAMLIdP struct {
	Metadata string `json:"metadata"`
}

type SpSAMLMetadata struct {
	Endpoint          string `json:"endpoint"`
	Key               string `json:"key"`
	Cert              string `json:"cert"`
	AuthnNameIDFormat string `json:"authn_name_id_format"`
}

// type Oauth2Config struct {
// 	// TenantUseLDAP      bool   `json:"tenant_use_ldap" description:"是否用于租户LDAP用户`
// 	// Tenant             string `json:"tenant" description:"租户"`
// 	ServerAddress      string `json:"server_address"`
// 	AuthorizeURIPrefix string `json:"authorize_uri_prefix"`
// 	TokenURIPrefix     string `json:"token_uri_prefix"`
// 	UserInfoURIPrefix  string `json:"user_info_uri_prefix"`
// 	UserNameField      string `json:"user_name_field"`
// 	ClientID           string `json:"client_id"`
// 	ClientSecret       string `json:"client_secret"`
// }

/* 当前服务作为sso sp时，记录信任的idP*/

// sso idp(单点登录身份提供商验证信息)
type IdentityProvider struct {
	imachinery.ObjectMeta
	Enable   bool     `json:"enable"`
	Protocol string   `json:"protocol"`
	Endpoint string   `json:"endpoint"`
	SAML     *SAMLIdP `json:"saml" gorm:"-"`
	// Oauth2
}

func (IdentityProvider) TableName() string {
	return "identity_provider"
}

func (s IdentityProvider) Validate() error {
	switch s.Protocol {
	case SSOProtocolSAML:
		if s.SAML == nil {
			return errors.Errorf("saml is empty when protocol is saml")
		}

		if s.SAML.Metadata == "" {
			return errors.Errorf("idp metadata is empty")
		}

		if _, err := samlsp.ParseMetadata([]byte(s.SAML.Metadata)); err != nil {
			return errors.WithStack(err)
		}

	case SSOProtocolOAUTH2:
	}
	return errors.Errorf("invalid sso protocol")
}

func (s *IdentityProvider) BeforeCreate(tx *gorm.DB) error {
	if s.Extend == nil {
		s.Extend = make(map[string]any)
	}
	if s.SAML != nil {
		s.Extend["saml"] = json.ToString(s.SAML)
	}

	if err := s.ObjectMeta.BeforeCreate(tx); err != nil {
		return errors.Errorf("failed to run `BeforeCreate` hook: %w", err)
	}
	return nil
}

// AfterCreate run after create database record.
func (s *IdentityProvider) AfterCreate(tx *gorm.DB) error {
	return tx.Save(s).Error
}

// BeforeUpdate run before update database record.
func (s *IdentityProvider) BeforeUpdate(tx *gorm.DB) error {
	if s.Extend == nil {
		s.Extend = make(map[string]any)
	}
	if s.SAML != nil {
		s.Extend["saml"] = s.SAML
	}
	if err := s.ObjectMeta.BeforeUpdate(tx); err != nil {
		return fmt.Errorf("failed to run `BeforeUpdate` hook: %w", err)
	}
	return nil
}

func (s *IdentityProvider) AfterFind(tx *gorm.DB) error {
	if err := s.ObjectMeta.AfterFind(tx); err != nil {
		return errors.Errorf("failed to run `AfterFind` hook: %v", err)
	}

	return nil
}

type IdpSAMLMetadata struct {
	Endpoint          string `json:"endpoint"`
	Key               string `json:"key"`
	Cert              string `json:"cert"`
	AuthnNameIDFormat string `json:"authn_name_id_format"`
}

/* 当前服务作为sso idp时，记录信任的sp*/

// sso sp(单点登录服务提供商)
type ServiceProvider struct {
	imachinery.ObjectMeta
	Enable   bool   `json:"enable"`
	Protocol string `json:"protocol" binding:"required,oneof=saml oauth2"`
	Endpoint string `json:"endpoint" binding:"required"`

	SAML *SAMLSp `json:"saml" gorm:"-"`
	// Oauth2
}

func (ServiceProvider) TableName() string {
	return "service_provider"
}

func (s ServiceProvider) Validate() error {

	switch s.Protocol {
	case SSOProtocolSAML:
		if s.SAML == nil {
			return errors.Errorf("saml is empty when protocol is saml")
		}

		if s.SAML.Metadata == "" {
			return errors.Errorf("sp metadata is empty")
		}

		if _, err := samlsp.ParseMetadata([]byte(s.SAML.Metadata)); err != nil {
			return errors.WithStack(err)
		}

	case SSOProtocolOAUTH2:
	}
	return errors.Errorf("invalid sso protocol")
}

type SAMLSp struct {
	Key      string `json:"key"`
	Metadata string `json:"metadata"`
}

func (s *ServiceProvider) BeforeCreate(tx *gorm.DB) error {
	if s.Extend == nil {
		s.Extend = make(map[string]any)
	}
	if s.SAML != nil {
		s.Extend["saml"] = json.ToString(s.SAML)
	}

	if err := s.ObjectMeta.BeforeCreate(tx); err != nil {
		return errors.Errorf("failed to run `BeforeCreate` hook: %w", err)
	}
	return nil
}

// AfterCreate run after create database record.
func (s *ServiceProvider) AfterCreate(tx *gorm.DB) error {
	return tx.Save(s).Error
}

// BeforeUpdate run before update database record.
func (s *ServiceProvider) BeforeUpdate(tx *gorm.DB) error {
	if s.Extend == nil {
		s.Extend = make(map[string]any)
	}
	if s.SAML != nil {
		s.Extend["saml"] = s.SAML
	}
	if err := s.ObjectMeta.BeforeUpdate(tx); err != nil {
		return fmt.Errorf("failed to run `BeforeUpdate` hook: %w", err)
	}
	return nil
}

func (s *ServiceProvider) AfterFind(tx *gorm.DB) error {
	if err := s.ObjectMeta.AfterFind(tx); err != nil {
		return errors.Errorf("failed to run `AfterFind` hook: %v", err)
	}

	return nil
}
