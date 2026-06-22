package iregistry

import (
	"fmt"

	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/json"
	"github.com/wangweihong/gotoolbox/pkg/tls/httptls"

	"github.com/wangweihong/eazycloud/apis/imachinery"

	"gorm.io/gorm"
)

type RegistryShadow struct {
}

type Registry struct {
	// Standard object's metadata.
	imachinery.ObjectMeta `json:",inline"`

	Address          string             `json:"address"           binding:"required,url"`
	UserName         string             `json:"user_name"         binding:"required,name" comment:"用户名"`
	Password         string             `json:"password"          binding:"required"`
	VersionMajor     string             `json:"version_major"`
	State            string             `json:"state"`
	StateMessage     string             `json:"state_message"`
	TimeZoneAdaptive bool               `json:"timezone_adaptive"`
	RegistryShadow   string             `json:"-"`
	TlsConfig        *httptls.TlsConfig `json:"tls_config"        binding:"omitempty"                         gorm:"-"`
}

func (p *Registry) Validate() error {
	return nil
}

type RegistryExternalEntry struct {
	Registry
}

type RegistryList struct {
	Items []*Registry `json:"items"`

	Total int64 `json:"total"`
}

// BeforeCreate run before create database record.
func (p *Registry) BeforeCreate(tx *gorm.DB) error {
	if err := p.ObjectMeta.BeforeCreate(tx); err != nil {
		return fmt.Errorf("failed to run `BeforeCreate` hook: %w", err)
	}

	p.RegistryShadow = json.ToString(p.TlsConfig)

	return nil
}

// AfterCreate run after create database record.
func (p *Registry) AfterCreate(tx *gorm.DB) error {
	return tx.Save(p).Error
}

// BeforeUpdate run before update database record.
func (p *Registry) BeforeUpdate(tx *gorm.DB) error {
	if err := p.ObjectMeta.BeforeUpdate(tx); err != nil {
		return fmt.Errorf("failed to run `BeforeUpdate` hook: %w", err)
	}
	p.RegistryShadow = json.ToString(p.TlsConfig)

	return nil
}

func (p *Registry) AfterFind(tx *gorm.DB) error {
	if err := p.ObjectMeta.AfterFind(tx); err != nil {
		return errors.Errorf("failed to run `AfterFind` hook: %v", err)
	}

	if p.RegistryShadow != "" {
		if err := json.Unmarshal([]byte(p.RegistryShadow), &p.TlsConfig); err != nil {
			return errors.Errorf("failed to unmarshal shadow data: %v", err)
		}
	}

	return nil
}
