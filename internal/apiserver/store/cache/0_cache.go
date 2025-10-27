package cache

import (
	"github.com/wangweihong/gotoolbox/pkg/cache"

	"github.com/wangweihong/eazycloud/internal/apiserver/store"
)

// func GetCacheFactoryOr() (store.Factory, error) {
// 	return &datastore{
// 		registries: registryIndexer(),
// 	}, nil
// }

type datastore struct {
	registries cache.Indexer
}

func (ds *datastore) Registries() store.RegistryStore {
	return newRegistry(ds)
}
func (ds *datastore) Close() error {
	return nil
}

func (ds *datastore) EnsureScheme(metaTypes ...any) error {
	return nil
}
