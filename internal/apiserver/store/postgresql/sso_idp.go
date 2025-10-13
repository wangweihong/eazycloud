package postgresql

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"gorm.io/gorm"
)

type ssoIdp struct {
	ds *datastore
}

func newIdentityProvider(ds *datastore) *ssoIdp {
	return &ssoIdp{ds}
}

func (s *ssoIdp) List(ctx context.Context, param *iapiserver.IdentityProviderListRequest) ([]*iapiserver.IdentityProvider, int64, error) {
	var meta []*iapiserver.IdentityProvider
	var total int64

	resourceSpecificFilter := func(q *gorm.DB) *gorm.DB {
		if param.Protocol != "" {
			q = q.Where("protocol = ?", param.Protocol)
		}
		return q
	}

	err := param.ToQuery(ctx, s.ds.db, resourceSpecificFilter).
		Find(&meta).Count(&total).Error

	return meta, total, err
}

func (s *ssoIdp) Get(ctx context.Context, id string) (*iapiserver.IdentityProvider, error) {
	var meta iapiserver.IdentityProvider
	meta.ID = id
	err := s.ds.db.
		Model(&iapiserver.IdentityProvider{}).
		Find(&meta).Error
	return &meta, err
}

func (s *ssoIdp) GetByName(ctx context.Context, name string) (*iapiserver.IdentityProvider, error) {
	return nil, nil
}

func (s *ssoIdp) Add(ctx context.Context, data *iapiserver.IdentityProvider) (*iapiserver.IdentityProvider, error) {
	return data, nil
}

func (s *ssoIdp) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *ssoIdp) Update(ctx context.Context, data *iapiserver.IdentityProvider) error {
	return nil
}

func (s *ssoIdp) Sync(ctx context.Context, pages []*iapiserver.IdentityProvider) error {
	if len(pages) == 0 {
		return nil
	}

	return nil
}
