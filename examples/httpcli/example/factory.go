package example

var client Factory

// Factory defines dingtalk platform client interface.
type Factory interface {
	Users() UserAPI
}

// Client return the store client instance.
func Instance() Factory {
	return client
}

// SetClient set the iam store client.
func SetInstance(factory Factory) {
	client = factory
}
