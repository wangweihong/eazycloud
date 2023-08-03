package grpcoptions

import (
	"fmt"
	"path"

	"github.com/spf13/pflag"

	"github.com/wangweihong/eazycloud/pkg/util/stringutil"

	"github.com/wangweihong/eazycloud/internal/pkg/genericoptions"
)

// TCPOptions are for creating an generic gRPC server.
type TCPOptions struct {
	Required    bool   `json:"required"     mapstructure:"required"`
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port"    mapstructure:"bind-port"`
	TlsEnable   bool   `json:"tls-enable"   mapstructure:"tls-enable"`
	// ServerCert is the TLS cert info for serving secure traffic
	ServerCert genericoptions.GeneratableKeyCert `json:"tls"          mapstructure:"tls"`
}

// NewTCPOptions is for creating an generic tcp listen gRPC server.
func NewTCPOptions() *TCPOptions {
	return &TCPOptions{
		BindAddress: "0.0.0.0",
		BindPort:    8081,
		TlsEnable:   false,
		Required:    true,
	}
}

// Validate is used to parse and validate the parameters entered by the user at
// the command line when the program starts.
func (s *TCPOptions) Validate() []error {
	var errors []error

	if s.BindAddress == "" {
		errors = append(errors, fmt.Errorf("gRPC bind address `--tcp.bind-address` is empty"))
	}

	// BindPort = 0 means random port. maybe should support it ??
	if s.BindPort < 1 || s.BindPort > 65535 {
		errors = append(
			errors,
			fmt.Errorf(
				"--tcp.bind-port %v must be between 1 and 65535",
				s.BindPort,
			),
		)
	}

	if s.TlsEnable {
		if !stringutil.BothEmptyOrNone(s.ServerCert.CertKey.KeyFile, s.ServerCert.CertKey.CertFile) {
			errors = append(
				errors,
				fmt.Errorf(
					" --tcp.tls.cert-key.cert-file and --tcp.tls.cert-key.private-key-file must provided together",
				),
			)
			return errors
		}

		if !stringutil.BothEmptyOrNone(s.ServerCert.CertDirectory, s.ServerCert.PairName) {
			errors = append(
				errors,
				fmt.Errorf(" --grpc.tls.cert-dir and --grpc.tls.pair-name must provided together"),
			)
			return errors
		}

		errors = append(
			errors, fmt.Errorf(
				"if required tls server, you should set --tcp.tls.cert-key.cert-file "+
					"and --tcp.tls.cert-key.key-file to real Tls Certs or set --tcp.tls.cert-dir and --tcp.tls.pair-name "),
		)
	}

	return errors
}

// AddFlags adds flags related to features for a specific api server to the
// specified FlagSet.
func (s *TCPOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&s.Required, "tcp.required", s.Required,
		"Whether require tcp server, if not require, turning off tcp server")

	fs.StringVar(&s.BindAddress, "tcp.bind-address", s.BindAddress, ""+
		"The IP address on which to serve the --grpc.bind-port(set to 0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")

	fs.IntVar(&s.BindPort, "tcp.bind-port", s.BindPort, ""+
		"The port on which gRPC server to serve.  0 for turning off insecure (HTTP) port.")

	fs.BoolVar(&s.TlsEnable, "tcp.tls-enable", s.TlsEnable,
		"Whether enabled gRPC tls verified.",
	)

	fs.StringVar(&s.ServerCert.CertDirectory, "tcp.tls.cert-dir", s.ServerCert.CertDirectory, ""+
		"The directory where the TLS certs are located. "+
		"If --tcp.tls.cert-key.cert-file and --tcp.tls.cert-key.private-key-file are provided, "+
		"this flag will be ignored.")

	fs.StringVar(&s.ServerCert.PairName, "tcp.tls.pair-name", s.ServerCert.PairName, ""+
		"The name which will be used with --tcp.tls.cert-dir to make a cert and key filenames. "+
		"It becomes <cert-dir>/<pair-name>.crt and <cert-dir>/<pair-name>.key")

	fs.StringVar(&s.ServerCert.CertKey.CertFile, "tcp.tls.cert-file", s.ServerCert.CertKey.CertFile, ""+
		"File containing the default x509 Certificate for gRPC server. (CA cert, if any, concatenated "+
		"after server cert).")

	fs.StringVar(&s.ServerCert.CertKey.KeyFile, "tcp.tls.cert-key.",
		s.ServerCert.CertKey.KeyFile, ""+
			"File containing the default x509 private key matching --tcp.tls.cert-file.")
}

// Complete fills in any fields not set that are required to have valid data.
func (s *TCPOptions) Complete() error {
	if s == nil {
		return nil
	}

	keyCert := &s.ServerCert.CertKey
	if len(keyCert.CertFile) != 0 || len(keyCert.KeyFile) != 0 {
		return nil
	}

	if len(s.ServerCert.CertDirectory) > 0 {
		if len(s.ServerCert.PairName) == 0 {
			return fmt.Errorf("--tcp.tls.pair-name is required if --tcp.tls.cert-dir is set")
		}
		keyCert.CertFile = path.Join(s.ServerCert.CertDirectory, s.ServerCert.PairName+".crt")
		keyCert.KeyFile = path.Join(s.ServerCert.CertDirectory, s.ServerCert.PairName+".key")
	}

	return nil
}
