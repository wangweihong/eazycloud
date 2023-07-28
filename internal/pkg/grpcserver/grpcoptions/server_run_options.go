package grpcoptions

import (
	"fmt"
	"strings"

	"github.com/wangweihong/eazycloud/pkg/util/maputil"
	"github.com/wangweihong/eazycloud/pkg/util/sliceutil"

	"github.com/spf13/pflag"

	"github.com/wangweihong/eazycloud/internal/pkg/grpcserver/interceptor"
	"github.com/wangweihong/eazycloud/pkg/sets"

	"github.com/wangweihong/eazycloud/internal/pkg/grpcserver"
)

// ServerRunOptions contains the options while running a generic gRPC server.
type ServerRunOptions struct {
	Version            bool     `json:"version"             mapstructure:"version"`             // 开启版本服务
	Reflect            bool     `json:"reflect"             mapstructure:"reflect"`             // 是否开启gRPC反射服务。开启反射服务后, grpcurl工具才能获取gRPC服务接口
	Debug              bool     `json:"debug"               mapstructure:"debug"`               // 是否开启调试服务
	UnaryInterceptors  []string `json:"unary-interceptors"  mapstructure:"unary-interceptors"`  // 启动拦截器列表
	StreamInterceptors []string `json:"stream-interceptors" mapstructure:"stream-interceptors"` // 启动拦截器列表
}

// NewServerRunOptions creates a new ServerRunOptions object with default parameters.
func NewServerRunOptions() *ServerRunOptions {
	defaults := grpcserver.NewConfig()

	return &ServerRunOptions{
		Version:            defaults.Version,
		Reflect:            defaults.Reflect,
		Debug:              defaults.Debug,
		UnaryInterceptors:  defaults.UnaryInterceptors,
		StreamInterceptors: defaults.StreamInterceptors,
	}
}

// ApplyTo applies the run options to the method receiver and returns self.
func (s *ServerRunOptions) ApplyTo(c *grpcserver.GRPCConfig) error {
	c.Version = s.Version
	c.Reflect = s.Reflect
	c.Debug = s.Debug
	c.UnaryInterceptors = s.UnaryInterceptors
	c.StreamInterceptors = s.StreamInterceptors
	return nil
}

// Validate checks validation of ServerRunOptions.
func (s *ServerRunOptions) Validate() []error {
	errors := []error{}

	rm, repeated := sliceutil.StringSlice(s.UnaryInterceptors).GetRepeat()
	if repeated {
		errors = append(errors, fmt.Errorf("unary interceptors `%v` is repeated", maputil.StringIntMap(rm).Keys()))
	}

	supportedUnaryInterceptor := sets.NewString(interceptor.UnaryServerInterceptorNames...)
	if !supportedUnaryInterceptor.HasAll(s.UnaryInterceptors...) {
		invalidInterceptors := sets.NewString(s.UnaryInterceptors...).Difference(supportedUnaryInterceptor)
		errors = append(errors, fmt.Errorf("unary intercerptor `%v` is not supported", invalidInterceptors.List()))
	}
	return errors
}

// AddFlags adds flags for a specific APIServer to the specified FlagSet.
func (s *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs.BoolVar(&s.Version, "server.version", s.Version, ""+
		"Install version service.")

	fs.BoolVar(&s.Reflect, "server.reflect", s.Reflect, ""+
		"Whether enable gRPC server register reflect service. If registered, grpc client can "+
		"get gRPC service list directly. ")

	fs.BoolVar(&s.Debug, "server.debug", s.Debug, ""+
		"Install debug service.")

	fs.StringSliceVar(
		&s.UnaryInterceptors,
		"server.unary-interceptors",
		s.UnaryInterceptors,
		"List of allowed unary interceptors for server, comma separated. If this list is empty,no unary interceptors will be used."+
			"Support unary interceptors: "+strings.Join(
			interceptor.UnaryServerInterceptorNames,
			",",
		),
	)
}
