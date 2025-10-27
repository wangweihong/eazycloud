package registry

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/apis/imachinery"
	"github.com/wangweihong/eazycloud/apis/iregistry"
	"github.com/wangweihong/eazycloud/internal/apiserver/store"
)

type RegistrySrv interface {
	Add(ctx context.Context, registry *iregistry.Registry, opts imachinery.CreateOptions) (*iregistry.Registry, error)
	Update(ctx context.Context, registry *iregistry.Registry, opts imachinery.UpdateOptions) (*iregistry.Registry, error)
	Delete(ctx context.Context, registryID string, opts imachinery.DeleteOptions) error
	Get(ctx context.Context, registryID string, opts imachinery.GetOptions) (*iregistry.Registry, error)
	List(ctx context.Context, opts *iapiserver.RegistryListRequest) (*iregistry.RegistryList, error)
	//Search(ctx context.Context, opts imachinery.ListOptions) (*iregistry.RegistryList, error)
}

type registryService struct {
	store store.Factory
}

func NewService(str store.Factory) *registryService {
	return &registryService{store: str}
}
