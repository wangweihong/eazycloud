package postgresql

import (
	"fmt"
	"sync"

	"github.com/wangweihong/gotoolbox/pkg/errors"

	"github.com/wangweihong/eazycloud/internal/apiserver/store"
	"github.com/wangweihong/eazycloud/pkg/httpsvr/genericoptions"

	"gorm.io/gorm"
)

var (
	postgresqlFactory store.Factory
	once              sync.Once
)

// GetPostgresSQLFactoryOr create postgresql factory with the given config.
func GetPostgresSQLFactoryOr(opts *genericoptions.PostgresSQLOptions) (store.Factory, error) {
	if opts == nil && postgresqlFactory == nil {
		return nil, errors.Errorf("failed to get postgresql store factory")
	}

	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		dbIns, err = opts.NewClient()
		postgresqlFactory = &datastore{dbIns}
	})

	if postgresqlFactory == nil || err != nil {
		return nil, fmt.Errorf(
			"failed to get postgresql store factory, postgresqlFactory: %+v, error: %w",
			postgresqlFactory,
			err,
		)
	}

	return postgresqlFactory, nil
}

func CheckExists(db *gorm.DB, model interface{}, fields map[string]interface{}) bool {
	query := db.Model(model)
	for field, value := range fields {
		query = query.Where(field+" = ?", value)
	}

	var count int64
	query.Count(&count)
	return count > 0
}

type datastore struct {
	db *gorm.DB
	// redis ?
}

func (ds *datastore) Close() error {
	db, err := ds.db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}

	return db.Close()
}

func (ds *datastore) EnsureScheme(metaTypes ...any) error {
	if err := ds.db.AutoMigrate(metaTypes...); err != nil {
		return err
	}
	return nil
}

func (ds *datastore) Registries() store.RegistryStore {
	return newRegistry(ds)
}

func (ds *datastore) Kubernetes() store.KubernetesStore {
	return newKubernetes(ds)
}

func (ds *datastore) Clusters() store.ClusterStore {
	return newCluster(ds)
}

func (ds *datastore) Users() store.UserStore {
	return newUser(ds)
}

func (ds *datastore) AppStores() store.AppStoreStore {
	return newAppStore(ds)
}

func (ds *datastore) ApplicationInstances() store.ApplicationInstanceStore {
	return newApplicationInstance(ds)
}

func (ds *datastore) ApplicationTemplates() store.ApplicationTemplateStore {
	return newApplicationTemplate(ds)
}

func (ds *datastore) ApplicationCategories() store.ApplicationCategoryStore {
	return newApplicationCategory(ds)
}

/* ------ setting ------- */
func (ds *datastore) IdentityProviders() store.IdentityProviderStore {
	return newIdentityProvider(ds)
}

func (ds *datastore) ServiceProviders() store.ServiceProviderStore {
	return newServiceProvider(ds)
}

func (ds *datastore) Settings() store.SettingStore {
	return newSetting(ds)
}

func (ds *datastore) OneTimeTokens() store.OneTimeTokenStore {
	return newOneTimeToken(ds)
}

func (ds *datastore) UserOTPs() store.UserOTPStore {
	return newUserOTP(ds)
}
