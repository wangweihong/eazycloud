package store

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/apis/iregistry"
)

type RegistryStore interface {
	List(ctx context.Context, param *iapiserver.RegistryListRequest) ([]*iregistry.Registry, int64, error)
	Get(ctx context.Context, id string) (*iregistry.Registry, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, data *iregistry.Registry) error
	Sync(ctx context.Context, datas []*iregistry.Registry) error
	Add(ctx context.Context, data *iregistry.Registry) (*iregistry.Registry, error)
}

type KubernetesStore interface {
	List(ctx context.Context) ([]*iapiserver.Cluster, error)
	Get(ctx context.Context, id string) (*iapiserver.Cluster, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, data *iapiserver.Cluster) error
	Sync(ctx context.Context, datas []*iapiserver.Cluster) error
	Add(ctx context.Context, data *iapiserver.Cluster) (*iapiserver.Cluster, error)
}

type ClusterStore interface {
	List(ctx context.Context) ([]*iapiserver.Cluster, error)
	Get(ctx context.Context, id string) (*iapiserver.Cluster, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, data *iapiserver.Cluster) error
	Sync(ctx context.Context, datas []*iapiserver.Cluster) error
	Add(ctx context.Context, data *iapiserver.Cluster) (*iapiserver.Cluster, error)
}

type SystemConfigStore interface {
	List(ctx context.Context) ([]*iapiserver.Cluster, error)
	Get(ctx context.Context, id string) (*iapiserver.Cluster, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, data *iapiserver.Cluster) error
	Sync(ctx context.Context, datas []*iapiserver.Cluster) error
	Add(ctx context.Context, data *iapiserver.Cluster) (*iapiserver.Cluster, error)
}

type AppStoreStore interface {
	List(ctx context.Context) ([]*iapiserver.AppStore, error)
	Query(ctx context.Context, req *iapiserver.AppStoreListRequest) ([]*iapiserver.AppStore, int64, error)
	Get(ctx context.Context, id string) (*iapiserver.AppStore, error)
	GetByName(ctx context.Context, id string) (*iapiserver.AppStore, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, data *iapiserver.AppStore) error
	Sync(ctx context.Context, datas []*iapiserver.AppStore) error
	Add(ctx context.Context, data *iapiserver.AppStore) (*iapiserver.AppStore, error)
}

type ApplicationInstanceStore interface {
	List(ctx context.Context, req *iapiserver.ApplicationInstanceListRequest) ([]*iapiserver.ApplicationInstance, int64, error)
	Get(ctx context.Context, id string) (*iapiserver.ApplicationInstance, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, data *iapiserver.ApplicationInstance) error
	Sync(ctx context.Context, datas []*iapiserver.ApplicationInstance) error
	Add(ctx context.Context, data *iapiserver.ApplicationInstance) (*iapiserver.ApplicationInstance, error)
}

type ApplicationTemplateStore interface {
	List(ctx context.Context, req *iapiserver.ApplicationTemplateListRequest) ([]*iapiserver.ApplicationTemplate, int64, error)
	Get(ctx context.Context, id string) (*iapiserver.ApplicationTemplate, error)
	GetByName(ctx context.Context, store string, id string) (*iapiserver.ApplicationTemplate, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, data *iapiserver.ApplicationTemplate) error
	Sync(ctx context.Context, datas []*iapiserver.ApplicationTemplate) error
	Add(ctx context.Context, data *iapiserver.ApplicationTemplate) (*iapiserver.ApplicationTemplate, error)
}

type ApplicationCategoryStore interface {
	List(ctx context.Context, req *iapiserver.ApplicationCategoryListRequest) ([]*iapiserver.ApplicationCategory, int64, error)
	Get(ctx context.Context, id string) (*iapiserver.ApplicationCategory, error)
	GetByName(ctx context.Context, store string, name string) (*iapiserver.ApplicationCategory, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, data *iapiserver.ApplicationCategory) error
	Sync(ctx context.Context, datas []*iapiserver.ApplicationCategory) error
	Add(ctx context.Context, data *iapiserver.ApplicationCategory) (*iapiserver.ApplicationCategory, error)
}

type IdentityProviderStore interface {
	List(ctx context.Context, req *iapiserver.IdentityProviderListRequest) ([]*iapiserver.IdentityProvider, int64, error)
	Get(ctx context.Context, id string) (*iapiserver.IdentityProvider, error)
	GetByName(ctx context.Context, name string) (*iapiserver.IdentityProvider, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, data *iapiserver.IdentityProvider) (*iapiserver.IdentityProvider, error)
	Sync(ctx context.Context, datas []*iapiserver.IdentityProvider) error
	Add(ctx context.Context, data *iapiserver.IdentityProvider) (*iapiserver.IdentityProvider, error)
}

type ServiceProviderStore interface {
	List(ctx context.Context, req *iapiserver.ServiceProviderListRequest) ([]*iapiserver.ServiceProvider, int64, error)
	Add(ctx context.Context, data *iapiserver.ServiceProvider) (*iapiserver.ServiceProvider, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*iapiserver.ServiceProvider, error)
	GetByKey(ctx context.Context, protocol, key string) (*iapiserver.ServiceProvider, error)
	GetByName(ctx context.Context, name string) (*iapiserver.ServiceProvider, error)
	Update(ctx context.Context, data *iapiserver.ServiceProvider) (*iapiserver.ServiceProvider, error)
	Sync(ctx context.Context, datas []*iapiserver.ServiceProvider) error
}

type SettingStore interface {
	List(ctx context.Context) ([]*iapiserver.Setting, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*iapiserver.Setting, error)
	GetByName(ctx context.Context, name string) (*iapiserver.Setting, error)
	GetMultiByNames(ctx context.Context, names ...string) ([]*iapiserver.Setting, error)
	Upsert(ctx context.Context, data *iapiserver.Setting) (*iapiserver.Setting, error)
	FirstOrCreate(ctx context.Context, data *iapiserver.Setting) (*iapiserver.Setting, error)
}

type UserStore interface {
	List(ctx context.Context, req *iapiserver.UserListRequest) ([]*iapiserver.User, int64, error)
	Get(ctx context.Context, id string) (*iapiserver.User, error)
	GetByName(ctx context.Context, name string) (*iapiserver.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, data *iapiserver.User) (*iapiserver.User, error)
	Sync(ctx context.Context, datas []*iapiserver.User) error
	Add(ctx context.Context, data *iapiserver.User) (*iapiserver.User, error)
}

type OneTimeTokenStore interface {
	GetByHash(ctx context.Context, hash string) (*iapiserver.OneTimeToken, error)
	Delete(ctx context.Context, id string) error
	Add(ctx context.Context, data *iapiserver.OneTimeToken) (*iapiserver.OneTimeToken, error)
	CleanupExpiredTokens(ctx context.Context) error
}

type UserOTPStore interface {
	List(ctx context.Context) ([]*iapiserver.UserOTP, error)
	Delete(ctx context.Context, id string) error
	GetByUser(ctx context.Context, uid string) (*iapiserver.UserOTP, error)
	Upsert(ctx context.Context, data *iapiserver.UserOTP) (*iapiserver.UserOTP, error)
	FirstOrCreate(ctx context.Context, data *iapiserver.UserOTP) (*iapiserver.UserOTP, error)
	Add(ctx context.Context, data *iapiserver.UserOTP) (*iapiserver.UserOTP, error)
}
