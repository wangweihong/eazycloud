package identity

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/internal/apiserver/store"
)

type identityService struct {
	store store.Factory
}

type IdentitySrv interface {
	UserList(ctx context.Context, req *iapiserver.UserListRequest) (*iapiserver.UserListResponse, error)
	UserGet(ctx context.Context, req *iapiserver.UserGetRequest) (*iapiserver.UserGetResponse, error)
	UserAdd(ctx context.Context, req *iapiserver.UserAddRequest) (*iapiserver.UserAddResponse, error)
	UserDelete(ctx context.Context, req *iapiserver.UserDeleteRequest) error
	UserUpdate(ctx context.Context, req *iapiserver.UserUpdateRequest) error
}

func NewService(str store.Factory) *identityService {
	return &identityService{store: str}
}
