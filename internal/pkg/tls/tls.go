package tls

import (
	"fmt"
	"os"

	"github.com/wangweihong/eazycloud/pkg/util/stringutil"
)

// CertKey contains configuration items related to certificate.
type CertKey struct {
	// CertFile is a file containing a PEM-encoded certificate, and possibly the complete certificate chain
	CertFile string `json:"cert-file"        mapstructure:"cert-file"`
	// KeyFile is a file containing a PEM-encoded private key for the certificate specified by CertFile
	KeyFile string `json:"private-key-file" mapstructure:"private-key-file"`
}

// GeneratableKeyCert contains configuration items related to certificate.
// +k8s:deepcopy-gen=true
type GeneratableKeyCert struct {
	CertData CertData `json:"cert-data" mapstructure:"cert-data"`
	// CertKey allows setting an explicit cert/key file to use.
	CertKey CertKey `json:"cert-key"  mapstructure:"cert-key"`

	// CertDirectory specifies a directory to write generated certificates to if CertFile/KeyFile aren't explicitly set.
	// PairName is used to determine the filenames within CertDirectory.
	// If CertDirectory and PairName are not set, an in-memory certificate will be generated.
	CertDirectory string `json:"cert-dir"  mapstructure:"cert-dir"`
	// PairName is the name which will be used with CertDirectory to make a cert and key filenames.
	// It becomes CertDirectory/PairName.crt and CertDirectory/PairName.key
	PairName string `json:"pair-name" mapstructure:"pair-name"`
}

// CertData contains configuration items related to certificate data.
type CertData struct {
	// PEM-encoded certificate data
	Cert string `json:"cert" mapstructure:"cert"`
	// PEM-encoded private key  data
	Key string `json:"key"  mapstructure:"key"`
}

func LoadDataFromFile(certPath, keyPath string) (string, string, error) {
	certPem, err := os.ReadFile(certPath)
	if err != nil {
		return "", "", err
	}
	keyPem, err := os.ReadFile(keyPath)
	if err != nil {
		return "", "", err
	}
	return string(certPem), string(keyPem), nil
}

func (s GeneratableKeyCert) Validate() error {
	if !stringutil.BothEmptyOrNone(s.CertData.Cert, s.CertData.Key) {
		return fmt.Errorf("cert-data.cert and cert-data.key must provided together")
	}

	if s.CertData.Cert != "" {
		return nil
	}

	if !stringutil.BothEmptyOrNone(s.CertKey.KeyFile, s.CertKey.CertFile) {
		return fmt.Errorf("cert-key.cert-file and cert-key.private-key-file must provided together")
	}

	if s.CertKey.CertFile != "" {
		return nil
	}

	if !stringutil.BothEmptyOrNone(s.CertDirectory, s.PairName) {
		return fmt.Errorf(" --cert-dir and --pair-name must provided together")
	}

	if s.CertDirectory != "" {
		return nil
	}

	return fmt.Errorf(
		"if required tls server, you should choose one way to set cert and key:" +
			"cert-key.cert-file and cert-key.key-file," +
			"cert-data.cert and cert-data.key," +
			"cert-dir and pair-name ")
}

// CopyAndHide deepcopy cert and hide cert and key.
func (s *GeneratableKeyCert) CopyAndHide() *GeneratableKeyCert {
	o := s.DeepCopy()
	if s.CertData.Cert != "" {
		s.CertData.Cert = "-"
	}

	if s.CertData.Key != "" {
		s.CertData.Key = "-"
	}
	return o
}
