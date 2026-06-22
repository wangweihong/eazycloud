package v1

import (
	"github.com/wangweihong/eazycloud/internal/apiserver/service/v1/application"
	"github.com/wangweihong/eazycloud/internal/apiserver/service/v1/cluster"
	"github.com/wangweihong/eazycloud/internal/apiserver/service/v1/identity"
	"github.com/wangweihong/eazycloud/internal/apiserver/service/v1/kubernetes"
	"github.com/wangweihong/eazycloud/internal/apiserver/service/v1/registry"
	"github.com/wangweihong/eazycloud/internal/apiserver/service/v1/setting"
	"github.com/wangweihong/eazycloud/internal/apiserver/store"
)

// Service defines functions used to return resource interface.
type Service interface {
	Registries() registry.RegistrySrv
	Kubernetes() kubernetes.KubernetesSrv
	Clusters() cluster.ClusterSrv
	Applications() application.ApplicationSrv
	Settings() setting.SettingSrv
	Identities() identity.IdentitySrv
}

type service struct {
	store store.Factory
}

// NewService returns Service interface.
func NewService(store store.Factory) Service {
	return &service{
		store: store,
	}
}

func (s *service) Registries() registry.RegistrySrv {
	return registry.NewService(s.store)
}

func (s *service) Kubernetes() kubernetes.KubernetesSrv {
	return kubernetes.NewService(s.store)
}

func (s *service) Clusters() cluster.ClusterSrv {
	return cluster.NewService(s.store)
}

func (s *service) Applications() application.ApplicationSrv {
	return application.NewService(s.store)
}

func (s *service) Settings() setting.SettingSrv {
	return setting.NewService(s.store)
}

func (s *service) Identities() identity.IdentitySrv {
	return identity.NewService(s.store)
}
