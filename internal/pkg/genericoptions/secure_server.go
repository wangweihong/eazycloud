package genericoptions

import (
	"fmt"
	"path"

	"github.com/wangweihong/eazycloud/pkg/app"

	"github.com/wangweihong/eazycloud/internal/pkg/genericserver"

	"github.com/spf13/pflag"
)

var _ app.CompleteableOptions = &SecureServingOptions{}

// SecureServingOptions contains configuration items related to HTTPS server startup.
type SecureServingOptions struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	// BindPort is ignored when Listener is set, will serve HTTPS even with 0.
	BindPort int `json:"bind-port"    mapstructure:"bind-port"`
	// Required set to true means that BindPort cannot be zero.
	Required bool `json:"required"     mapstructure:"required"`
	// ServerCert is the TLS cert info for serving secure t`raffic
	ServerCert GeneratableKeyCert `json:"tls"          mapstructure:"tls"`
}

// CertKey contains configuration items related to certificate.
type CertKey struct {
	// CertFile is a file containing a PEM-encoded certificate, and possibly the complete certificate chain
	CertFile string `json:"cert-file"        mapstructure:"cert-file"`
	// KeyFile is a file containing a PEM-encoded private key for the certificate specified by CertFile
	KeyFile string `json:"private-key-file" mapstructure:"private-key-file"`
}

// GeneratableKeyCert contains configuration items related to certificate.
type GeneratableKeyCert struct {
	// CertKey allows setting an explicit cert/key file to use.
	CertKey CertKey `json:"cert-key" mapstructure:"cert-key"`

	// CertDirectory specifies a directory to write generated certificates to if CertFile/KeyFile aren't explicitly set.
	// PairName is used to determine the filenames within CertDirectory.
	// If CertDirectory and PairName are not set, an in-memory certificate will be generated.
	CertDirectory string `json:"cert-dir"  mapstructure:"cert-dir"`
	// PairName is the name which will be used with CertDirectory to make a cert and key filenames.
	// It becomes CertDirectory/PairName.crt and CertDirectory/PairName.key
	PairName string `json:"pair-name" mapstructure:"pair-name"`
}

// NewSecureServingOptions creates a SecureServingOptions object with default parameters.
func NewSecureServingOptions() *SecureServingOptions {
	return &SecureServingOptions{
		BindAddress: "0.0.0.0",
		BindPort:    8443,
		Required:    false,
	}
}

// ApplyTo applies the run options to the method receiver and returns self.
func (s *SecureServingOptions) ApplyTo(c *genericserver.Config) error {
	// SecureServing is required to serve https
	c.SecureServing = &genericserver.SecureServingInfo{
		BindAddress: s.BindAddress,
		BindPort:    s.BindPort,
		CertKey: genericserver.CertKey{
			CertFile: s.ServerCert.CertKey.CertFile,
			KeyFile:  s.ServerCert.CertKey.KeyFile,
		},
		Required: s.Required,
	}

	return nil
}

// Validate is used to parse and validate the parameters entered by the user at
// the command line when the program starts.
func (s *SecureServingOptions) Validate() []error {
	if s == nil {
		return nil
	}

	errors := []error{}

	if s.Required {
		if s.BindPort < 1 || s.BindPort > 65535 {
			errors = append(
				errors,
				fmt.Errorf(
					"--secure.bind-port %v must be between 1 and 65535",
					s.BindPort,
				),
			)
		}
		if s.ServerCert.CertKey.KeyFile != "" || s.ServerCert.CertKey.CertFile != "" {
			if s.ServerCert.CertKey.KeyFile == "" || s.ServerCert.CertKey.CertFile == "" {
				errors = append(
					errors,
					fmt.Errorf(
						" --secure.tls.cert-key.cert-file and --secure.tls.cert-key.private-key-file must provided together",
					),
				)
			}
			return errors
		}

		if s.ServerCert.CertDirectory != "" || s.ServerCert.PairName != "" {
			if s.ServerCert.CertDirectory == "" || s.ServerCert.PairName == "" {
				errors = append(
					errors,
					fmt.Errorf("  --secure.tls.cert-dir and --secure.tls.pair-name must provided together"),
				)
				return errors
			}
			return errors
		}

		errors = append(
			errors, fmt.Errorf(
				"if required tls server, you should set --secure.tls.cert-key.cert-file "+
					"and --secure.tls.cert-key.key-file to real Tls Certs or set --secure.tls.cert-dir and --secure.tls.pair-name "),
		)
	}

	return errors
}

// AddFlags adds flags related to HTTPS server for a specific APIServer to the
// specified FlagSet.
func (s *SecureServingOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&s.Required, "secure.required", s.Required,
		"Whether require secure server, if not require, turning off secure (HTTPs) port",
	)
	fs.StringVar(&s.BindAddress, "secure.bind-address", s.BindAddress, ""+
		"The IP address on which to listen for the --secure.bind-port port. The "+
		"associated interface(s) must be reachable by the rest of the engine, and by CLI/web "+
		"clients. If blank, all interfaces will be used (0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")
	desc := "The port on which to serve HTTPS with authentication and authorization."

	fs.IntVar(&s.BindPort, "secure.bind-port", s.BindPort, desc)

	fs.StringVar(&s.ServerCert.CertDirectory, "secure.tls.cert-dir", s.ServerCert.CertDirectory, ""+
		"The directory where the TLS certs are located. "+
		"If --secure.tls.cert-key.cert-file and --secure.tls.cert-key.private-key-file are provided, "+
		"this flag will be ignored.")

	fs.StringVar(&s.ServerCert.PairName, "secure.tls.pair-name", s.ServerCert.PairName, ""+
		"The name which will be used with --secure.tls.cert-dir to make a cert and key filenames. "+
		"It becomes <cert-dir>/<pair-name>.crt and <cert-dir>/<pair-name>.key")

	fs.StringVar(&s.ServerCert.CertKey.CertFile, "secure.tls.cert-key.cert-file", s.ServerCert.CertKey.CertFile, ""+
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
		"after server cert).")

	fs.StringVar(&s.ServerCert.CertKey.KeyFile, "secure.tls.cert-key.private-key-file",
		s.ServerCert.CertKey.KeyFile, ""+
			"File containing the default x509 private key matching --secure.tls.cert-key.cert-file.")
}

// Complete fills in any fields not set that are required to have valid data.
func (s *SecureServingOptions) Complete() error {
	if s == nil || !s.Required {
		return nil
	}

	keyCert := &s.ServerCert.CertKey
	if len(keyCert.CertFile) != 0 || len(keyCert.KeyFile) != 0 {
		return nil
	}

	if len(s.ServerCert.CertDirectory) > 0 {
		if len(s.ServerCert.PairName) == 0 {
			return fmt.Errorf("--secure.tls.pair-name is required if --secure.tls.cert-dir is set")
		}

		keyCert.CertFile = path.Join(s.ServerCert.CertDirectory, s.ServerCert.PairName+".crt")
		keyCert.KeyFile = path.Join(s.ServerCert.CertDirectory, s.ServerCert.PairName+".key")
	}

	return nil
}
