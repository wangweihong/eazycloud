package store

var client Factory

// Factory defines the iam platform storage interface.
type Factory interface {
	Registries() RegistryStore
	Kubernetes() KubernetesStore
	Clusters() ClusterStore
	Users() UserStore

	//applicatoins
	ApplicationInstances() ApplicationInstanceStore
	ApplicationTemplates() ApplicationTemplateStore
	ApplicationCategories() ApplicationCategoryStore
	AppStores() AppStoreStore

	// settings
	IdentityProviders() IdentityProviderStore
	ServiceProviders() ServiceProviderStore
	Settings() SettingStore

	EnsureScheme(metaTypes ...any) error
	Close() error
}

// Client return the store client instance.
func Client() Factory {
	return client
}

// SetClient set the iam store client.
func SetClient(factory Factory) {
	client = factory
}
