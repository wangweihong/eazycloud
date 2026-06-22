package iregistry

const (
	RegistryStateUnknown   = "unknown"
	RegistryStateHealthy   = "healthy"
	RegistryStateUnhealthy = "unhealthy"
)

type RegistryState struct {
	State  string
	Detail string
}
