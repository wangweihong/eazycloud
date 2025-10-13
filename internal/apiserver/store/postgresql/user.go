package postgresql

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/imanager"
)

type user struct {
	ds *datastore
}

func newUser(ds *datastore) *user {
	return &user{ds}
}

func (s *user) List(ctx context.Context) ([]*imanager.User, error) {
	return nil, nil
}

func (s *user) Get(ctx context.Context, id string) (*imanager.User, error) {
	return nil, nil
}

func (s *user) GetByName(ctx context.Context, name string) (*imanager.User, error) {
	return nil, nil
}

func (s *user) Add(ctx context.Context, data *imanager.User) (*imanager.User, error) {
	return data, nil
}

func (s *user) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *user) Update(ctx context.Context, data *imanager.User) error {
	return nil
}

func (s *user) Sync(ctx context.Context, pages []*imanager.User) error {
	if len(pages) == 0 {
		return nil
	}

	return nil
}
