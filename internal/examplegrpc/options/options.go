package options

import (
	"github.com/wangweihong/eazycloud/internal/pkg/grpcserver/grpcoptions"
	"github.com/wangweihong/eazycloud/pkg/app"
	cliflag "github.com/wangweihong/eazycloud/pkg/cli/flag"
	"github.com/wangweihong/eazycloud/pkg/json"
	"github.com/wangweihong/eazycloud/pkg/log"
)

var (
	_ app.PrintableOptions    = &Options{}
	_ app.CompleteableOptions = &Options{}
)

// Options runs a http server.
type Options struct {
	Name             string                        `json:"name"`
	Log              *log.Options                  `json:"log"    mapstructure:"log"`
	GRPC             *grpcoptions.GRPCOptions      `json:"grpc"   mapstructure:"grpc"`
	ServerRunOptions *grpcoptions.ServerRunOptions `json:"server" mapstructure:"server"`
}

// NewOptions creates a new Options object with default parameters.
func NewOptions() *Options {
	s := Options{
		Name: "example-gRPC",

		Log:              log.NewOptions(),
		GRPC:             grpcoptions.NewGRPCOptions(),
		ServerRunOptions: grpcoptions.NewServerRunOptions(),
	}

	return &s
}

// Flags returns flags for a specific server by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.Log.AddFlags(fss.FlagSet("logs"))
	o.GRPC.AddFlags(fss.FlagSet("grpc"))
	o.ServerRunOptions.AddFlags(fss.FlagSet("server"))
	return fss
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

// Complete fills in any fields not set that are required to have valid data.
// 补全指定的选项.
func (o *Options) Complete() error {
	if err := o.GRPC.Complete(); err != nil {
		return err
	}
	return nil
}
