package options

// BackendOptions contains configuration items related to server features.
type BackendOptions struct {
	Address string `json:"address" mapstructure:"address"`
}
