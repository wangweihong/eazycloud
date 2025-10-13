package cache

import (
	"context"
	"reflect"
	"sort"
	"sync"

	"github.com/google/uuid"
	"github.com/wangweihong/gotoolbox/pkg/cache"
	"github.com/wangweihong/gotoolbox/pkg/errors"

	"github.com/wangweihong/eazycloud/apis/iregistry"
)

func registryKeyFunc(obj any) (string, error) {
	if obj == nil {
		return "", errors.Errorf("object is nil")
	}
	data, ok := obj.(*iregistry.Registry)
	if !ok {
		return "", errors.Errorf("object is %v,not %v type", reflect.TypeOf(obj), reflect.TypeOf(&iregistry.Registry{}))
	}
	return data.ID, nil
}

func registryIndexer() cache.Indexer {
	indexers := make(map[string]cache.IndexFunc)
	return cache.NewIndexer(registryKeyFunc, indexers)
}

type registry struct {
	lock    sync.RWMutex
	indexer cache.Indexer
}

func newRegistry(ds *datastore) *registry {
	return &registry{
		indexer: ds.registries,
	}
}

func (s *registry) List(ctx context.Context) ([]*iregistry.Registry, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	list := s.indexer.List()
	result := make([]*iregistry.Registry, 0, len(list))
	for _, v := range list {
		if v, ok := v.(*iregistry.Registry); ok {
			result = append(result, v)
		}
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result, nil
}

func (s *registry) Get(ctx context.Context, id string) (*iregistry.Registry, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	tmp := &iregistry.Registry{}
	tmp.ID = id

	meta, exist, err := s.indexer.Get(tmp)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.Errorf("not exist")
	}

	data, ok := meta.(*iregistry.Registry)
	if !ok {
		return nil, errors.Errorf("data type not iregistry.Registry")
	}

	return data, nil
}

func (s *registry) Create(ctx context.Context, data *iregistry.Registry) (*iregistry.Registry, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	data.ID = uuid.New().String()
	err := s.indexer.Add(data)
	return data, err
}

func (s *registry) Delete(ctx context.Context, id string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	tmp := &iregistry.Registry{}
	tmp.ID = id
	return s.indexer.Delete(tmp)
}

func (s *registry) Update(ctx context.Context, data *iregistry.Registry) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.indexer.Update(data)
}

func (s *registry) Sync(ctx context.Context, pages []*iregistry.Registry) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for i := range pages {
		if err := s.indexer.Add(pages[i]); err != nil {
			return err
		}
	}

	return nil
}
