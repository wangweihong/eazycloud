package application

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/internal/apiserver/store"
)

type ApplicationSrv interface {
	AppStoreQuery(ctx context.Context, req *iapiserver.AppStoreListRequest) (*iapiserver.AppStoreListResponse, error)
	AppStoreList(ctx context.Context, req *iapiserver.AppStoreListRequest) (*iapiserver.AppStoreListResponse, error)
	AppStoreGet(ctx context.Context, req *iapiserver.AppStoreGetRequest) (*iapiserver.AppStoreGetResponse, error)
	AppStoreAdd(ctx context.Context, req *iapiserver.AppStoreAddRequest) (*iapiserver.AppStoreAddResponse, error)
	AppStoreDelete(ctx context.Context, req *iapiserver.AppStoreDeleteRequest) error
	AppStoreUpdate(ctx context.Context, req *iapiserver.AppStoreUpdateRequest) error
	AppStoreSync(ctx context.Context, req *iapiserver.AppStoreSyncRequest) error

	ApplicationTemplateValidate(ctx context.Context, req *iapiserver.ApplicationTemplateValidateRequest) error
	ApplicationTemplateList(ctx context.Context, req *iapiserver.ApplicationTemplateListRequest) (*iapiserver.ApplicationTemplateListResponse, error)
	ApplicationTemplateGet(ctx context.Context, req *iapiserver.ApplicationTemplateGetRequest) (*iapiserver.ApplicationTemplateGetResponse, error)
	ApplicationTemplateAdd(ctx context.Context, req *iapiserver.ApplicationTemplateAddRequest) (*iapiserver.ApplicationTemplateAddResponse, error)
	ApplicationTemplateDelete(ctx context.Context, req *iapiserver.ApplicationTemplateDeleteRequest) error
	ApplicationTemplateUpdate(ctx context.Context, req *iapiserver.ApplicationTemplateUpdateRequest) error

	ApplicationTemplateVersionList(ctx context.Context, req *iapiserver.ApplicationTemplateVersionListRequest) (*iapiserver.ApplicationTemplateVersionListResponse, error)
	ApplicationTemplateVersionGet(ctx context.Context, req *iapiserver.ApplicationTemplateVersionGetRequest) (*iapiserver.ApplicationTemplateVersionGetResponse, error)
	ApplicationTemplateVersionAdd(ctx context.Context, req *iapiserver.ApplicationTemplateVersionAddRequest) (*iapiserver.ApplicationTemplateVersionAddResponse, error)
	ApplicationTemplateVersionDelete(ctx context.Context, req *iapiserver.ApplicationTemplateVersionDeleteRequest) error
	ApplicationTemplateVersionUpdate(ctx context.Context, req *iapiserver.ApplicationTemplateVersionUpdateRequest) error

	ApplicationInstanceList(ctx context.Context, req *iapiserver.ApplicationInstanceListRequest) (*iapiserver.ApplicationInstanceListResponse, error)
	ApplicationInstanceGet(ctx context.Context, req *iapiserver.ApplicationInstanceGetRequest) (*iapiserver.ApplicationInstanceGetResponse, error)
	ApplicationInstanceAdd(ctx context.Context, req *iapiserver.ApplicationInstanceAddRequest) (*iapiserver.ApplicationInstanceAddResponse, error)
	ApplicationInstanceDelete(ctx context.Context, req *iapiserver.ApplicationInstanceDeleteRequest) error
	ApplicationInstanceUpdate(ctx context.Context, req *iapiserver.ApplicationInstanceUpdateRequest) error

	ApplicationCategoryList(ctx context.Context, req *iapiserver.ApplicationCategoryListRequest) (*iapiserver.ApplicationCategoryListResponse, error)
	ApplicationCategoryGet(ctx context.Context, req *iapiserver.ApplicationCategoryGetRequest) (*iapiserver.ApplicationCategoryGetResponse, error)
	ApplicationCategoryAdd(ctx context.Context, req *iapiserver.ApplicationCategoryAddRequest) (*iapiserver.ApplicationCategoryAddResponse, error)
	ApplicationCategoryDelete(ctx context.Context, req *iapiserver.ApplicationCategoryDeleteRequest) error
	ApplicationCategoryUpdate(ctx context.Context, req *iapiserver.ApplicationCategoryUpdateRequest) error
}

type applicationService struct {
	store store.Factory
}

func NewService(str store.Factory) *applicationService {
	return &applicationService{store: str}
}
