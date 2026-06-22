package ikubernetes

import (
	"encoding/base64"

	"k8s.io/client-go/rest"
)

type TLSClientConfig struct {
	Insecure   bool     `json:"insecure,omitempty"`
	ServerName string   `json:"server_name,omitempty"`
	CertFile   string   `json:"cert_file,omitempty"`
	KeyFile    string   `json:"key_file,omitempty"`
	CAFile     string   `json:"ca_file,omitempty"`
	CertData   string   `json:"cert_data,omitempty"`
	KeyData    string   `json:"key_data,omitempty"`
	CAData     string   `json:"ca_data,omitempty"`
	NextProtos []string `json:"next_protos,omitempty"`
}

type ClusterConfig struct {
	Host            string          `json:"host,omitempty"`
	APIPath         string          `json:"api_path,omitempty"`
	Username        string          `json:"username,omitempty"`
	Password        string          `json:"password,omitempty"`
	BearerToken     string          `json:"bearer_token,omitempty"`
	BearerTokenFile string          `json:"bearer_token_file,omitempty"`
	TlsClientConfig TLSClientConfig `json:"tls_client_config,omitempty"`
}

func (tc *ClusterConfig) ToRestConfig() *rest.Config {
	certData, _ := base64.StdEncoding.DecodeString(tc.TlsClientConfig.CertData)
	keyData, _ := base64.StdEncoding.DecodeString(tc.TlsClientConfig.KeyData)
	caData, _ := base64.StdEncoding.DecodeString(tc.TlsClientConfig.CAData)

	return &rest.Config{
		Host:                tc.Host,
		APIPath:             tc.APIPath,
		ContentConfig:       rest.ContentConfig{},
		Username:            tc.Username,
		Password:            tc.Password,
		BearerToken:         tc.BearerToken,
		BearerTokenFile:     tc.BearerTokenFile,
		Impersonate:         rest.ImpersonationConfig{},
		AuthProvider:        nil,
		AuthConfigPersister: nil,
		ExecProvider:        nil,
		TLSClientConfig: rest.TLSClientConfig{
			CertData: []byte(certData),
			KeyData:  []byte(keyData),
			CAData:   []byte(caData),
		},
		UserAgent:          "",
		DisableCompression: false,
		Transport:          nil,
		WrapTransport:      nil,
		QPS:                0,
		Burst:              0,
		RateLimiter:        nil,
		WarningHandler:     nil,
		Timeout:            0,
		Dial:               nil,
		Proxy:              nil,
	}
}
