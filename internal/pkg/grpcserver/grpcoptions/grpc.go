package grpcoptions

import (
	"fmt"
	"path"

	"github.com/spf13/pflag"

	"github.com/wangweihong/eazycloud/internal/pkg/genericoptions"
)

// GRPCOptions are for creating an unauthenticated, unauthorized, insecure port.
// No one should be using these anymore.
type GRPCOptions struct {
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port"    mapstructure:"bind-port"`
	MaxMsgSize  int    `json:"max-msg-size" mapstructure:"max-msg-size"`
	TlsEnable   bool   `json:"tls-enable"   mapstructure:"tls-enable"`
	// ServerCert is the TLS cert info for serving secure traffic
	ServerCert genericoptions.GeneratableKeyCert `json:"tls"          mapstructure:"tls"`
}

// NewGRPCOptions is for creating an unauthenticated, unauthorized, insecure port.
// No one should be using these anymore.
func NewGRPCOptions() *GRPCOptions {
	return &GRPCOptions{
		BindAddress: "0.0.0.0",
		BindPort:    8081,
		MaxMsgSize:  4 * 1024 * 1024,
		TlsEnable:   false,
	}
}

// Validate is used to parse and validate the parameters entered by the user at
// the command line when the program starts.
func (s *GRPCOptions) Validate() []error {
	var errors []error

	if s.BindPort < 1 || s.BindPort > 65535 {
		errors = append(
			errors,
			fmt.Errorf(
				"--grpc.bind-port %v must be between 1 and 65535",
				s.BindPort,
			),
		)
	}

	if s.TlsEnable {
		if s.ServerCert.CertKey.KeyFile != "" || s.ServerCert.CertKey.CertFile != "" {
			if s.ServerCert.CertKey.KeyFile == "" || s.ServerCert.CertKey.CertFile == "" {
				errors = append(
					errors,
					fmt.Errorf(
						" --grpc.tls.cert-key.cert-file and --grpc.tls.cert-key.private-key-file must provided together",
					),
				)
			}
			return errors
		}

		if s.ServerCert.CertDirectory != "" || s.ServerCert.PairName != "" {
			if s.ServerCert.CertDirectory == "" || s.ServerCert.PairName == "" {
				errors = append(
					errors,
					fmt.Errorf(" --grpc.tls.cert-dir and --grpc.tls.pair-name must provided together"),
				)
				return errors
			}
			return errors
		}

		errors = append(
			errors, fmt.Errorf(
				"if required tls server, you should set --grpc.tls.cert-key.cert-file "+
					"and --grpc.tls.cert-key.key-file to real Tls Certs or set --grpc.tls.cert-dir and --grpc.tls.pair-name "),
		)
	}

	return errors
}

// AddFlags adds flags related to features for a specific api server to the
// specified FlagSet.
func (s *GRPCOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.BindAddress, "grpc.bind-address", s.BindAddress, ""+
		"The IP address on which to serve the --grpc.bind-port(set to 0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")

	fs.IntVar(&s.BindPort, "grpc.bind-port", s.BindPort, ""+
		"The port on which gRPC server to serve unsecured.")

	fs.IntVar(&s.MaxMsgSize, "grpc.max-msg-size", s.MaxMsgSize, "gRPC max message size.")

	fs.BoolVar(&s.TlsEnable, "grpc.tls-enable", s.TlsEnable,
		"Whether enabled gRPC tls verified.",
	)

	fs.StringVar(&s.ServerCert.CertDirectory, "grpc.tls.cert-dir", s.ServerCert.CertDirectory, ""+
		"The directory where the TLS certs are located. "+
		"If --grpc.tls.cert-key.cert-file and --grpc.tls.cert-key.private-key-file are provided, "+
		"this flag will be ignored.")

	fs.StringVar(&s.ServerCert.PairName, "grpc.tls.pair-name", s.ServerCert.PairName, ""+
		"The name which will be used with --grpc.tls.cert-dir to make a cert and key filenames. "+
		"It becomes <cert-dir>/<pair-name>.crt and <cert-dir>/<pair-name>.key")

	fs.StringVar(&s.ServerCert.CertKey.CertFile, "grpc.tls.cert-file", s.ServerCert.CertKey.CertFile, ""+
		"File containing the default x509 Certificate for gRPC server. (CA cert, if any, concatenated "+
		"after server cert).")

	fs.StringVar(&s.ServerCert.CertKey.KeyFile, "grpc.tls.cert-key.",
		s.ServerCert.CertKey.KeyFile, ""+
			"File containing the default x509 private key matching --grpc.tls.cert-file.")
}

// Complete fills in any fields not set that are required to have valid data.
func (s *GRPCOptions) Complete() error {
	if s == nil || s.BindPort == 0 {
		return nil
	}

	keyCert := &s.ServerCert.CertKey
	if len(keyCert.CertFile) != 0 || len(keyCert.KeyFile) != 0 {
		return nil
	}

	if len(s.ServerCert.CertDirectory) > 0 {
		if len(s.ServerCert.PairName) == 0 {
			return fmt.Errorf("--grpc.tls.pair-name is required if --grpc.tls.cert-dir is set")
		}
		keyCert.CertFile = path.Join(s.ServerCert.CertDirectory, s.ServerCert.PairName+".crt")
		keyCert.KeyFile = path.Join(s.ServerCert.CertDirectory, s.ServerCert.PairName+".key")
	}

	return nil
}
