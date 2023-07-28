package grpcserver

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/wangweihong/eazycloud/internal/pkg/grpcserver/interceptor"

	"github.com/wangweihong/eazycloud/internal/pkg/genericoptions"
	"github.com/wangweihong/eazycloud/pkg/log"
	//"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc".
)

// GRPCConfig defines  configuration for the grpc server.
type GRPCConfig struct {
	TlsEnable          bool
	Addr               string
	MaxMsgSize         int
	ServerCert         genericoptions.GeneratableKeyCert
	Version            bool
	Reflect            bool
	Debug              bool
	UnaryInterceptors  []string
	StreamInterceptors []string
}

// NewConfig returns a Config struct with the default values.
func NewConfig() *GRPCConfig {
	return &GRPCConfig{
		Version: true,
		Reflect: false,
		Debug:   false,
		UnaryInterceptors: []string{
			interceptor.InterceptorNameRequestID,
			interceptor.InterceptorNameContext,
			interceptor.InterceptorNameLogger,
			interceptor.InterceptorNameRecovery,
		},
	}
}

type CompletedGRPCConfig struct {
	*GRPCConfig
}

// Complete fills in any fields not set that are required to have valid data and can be derived from other fields.
func (c *GRPCConfig) Complete() *CompletedGRPCConfig {
	if c.Addr == "" {
		c.Addr = "127.0.0.1:8081"
	}

	return &CompletedGRPCConfig{c}
}

// New create a grpc Server instance.
func (c *CompletedGRPCConfig) New() (*GRPCServer, error) {
	opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(c.MaxMsgSize)}

	if c.TlsEnable {
		creds, err := credentials.NewServerTLSFromFile(c.ServerCert.CertKey.CertFile, c.ServerCert.CertKey.KeyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %s", err.Error())
		}
		log.Info("gRPC service run with TLS")
		opts = append(opts, grpc.Creds(creds))
	}

	opts = installInterceptors(c.UnaryInterceptors, opts)
	// opts = append(opts, grpc.ChainStreamInterceptor(streamUnaryInterceptor...))

	gRPCServer := &GRPCServer{
		Server:             grpc.NewServer(opts...),
		address:            c.Addr,
		Version:            c.Version,
		Reflect:            c.Reflect,
		Debug:              c.Debug,
		UnaryInterceptors:  c.UnaryInterceptors,
		StreamInterceptors: c.StreamInterceptors,
	}

	initGenericGRPCServer(gRPCServer)

	return gRPCServer, nil
}
