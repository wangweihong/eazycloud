package setting

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/internal/apiserver/store"
)

type settingService struct {
	store store.Factory
}

type SettingSrv interface {
	IdentityProviderList(ctx context.Context, req *iapiserver.IdentityProviderListRequest) (*iapiserver.IdentityProviderListResponse, error)
	IdentityProviderGet(ctx context.Context, req *iapiserver.IdentityProviderGetRequest) (*iapiserver.IdentityProviderGetResponse, error)
	IdentityProviderAdd(ctx context.Context, req *iapiserver.IdentityProviderAddRequest) (*iapiserver.IdentityProviderAddResponse, error)
	IdentityProviderDelete(ctx context.Context, req *iapiserver.IdentityProviderDeleteRequest) error
	IdentityProviderUpdate(ctx context.Context, req *iapiserver.IdentityProviderUpdateRequest) error

	ServiceProviderList(ctx context.Context, req *iapiserver.ServiceProviderListRequest) (*iapiserver.ServiceProviderListResponse, error)
	ServiceProviderGet(ctx context.Context, req *iapiserver.ServiceProviderGetRequest) (*iapiserver.ServiceProviderGetResponse, error)
	ServiceProviderAdd(ctx context.Context, req *iapiserver.ServiceProviderAddRequest) (*iapiserver.ServiceProviderAddResponse, error)
	ServiceProviderDelete(ctx context.Context, req *iapiserver.ServiceProviderDeleteRequest) error
	ServiceProviderUpdate(ctx context.Context, req *iapiserver.ServiceProviderUpdateRequest) error

	IdentityProviderSAMLMetadataUpset(ctx context.Context, req *iapiserver.IdentityProviderMetadataUpsetRequest) error
	IdentityProviderSAMLMetadataGet(ctx context.Context) (*iapiserver.IdentityProviderMetadataGetResponse, error)
	ServiceProviderSAMLMetadataUpset(ctx context.Context, req *iapiserver.ServiceProviderMetadataUpsetRequest) error
}

func NewService(str store.Factory) *settingService {
	return &settingService{store: str}
}
