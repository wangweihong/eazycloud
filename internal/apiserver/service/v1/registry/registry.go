package registry

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/apis/imachinery"
	"github.com/wangweihong/eazycloud/apis/iregistry"
)

func (r *registryService) Add(ctx context.Context, registry *iregistry.Registry, opts imachinery.CreateOptions) (*iregistry.Registry, error) {
	return r.store.Registries().Add(ctx, registry)
}

func (r *registryService) Update(ctx context.Context, registry *iregistry.Registry, opts imachinery.UpdateOptions) (*iregistry.Registry, error) {
	return nil, r.store.Registries().Update(ctx, registry)
}

func (r *registryService) Delete(ctx context.Context, registryID string, opts imachinery.DeleteOptions) error {
	return r.store.Registries().Delete(ctx, registryID)
}

func (r *registryService) Get(ctx context.Context, registryID string, opts imachinery.GetOptions) (*iregistry.Registry, error) {
	return r.store.Registries().Get(ctx, registryID)
}

func (r *registryService) List(ctx context.Context, opts *iapiserver.RegistryListRequest) (*iregistry.RegistryList, error) {
	rl := &iregistry.RegistryList{}
	var err error
	rl.Items, rl.Total, err = r.store.Registries().List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return rl, nil
}
