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
	Name             string                         `json:"name"`
	Log              *log.Options                   `json:"log"    mapstructure:"log"`
	TCP              *grpcoptions.TCPOptions        `json:"tcp"    mapstructure:"tcp"`
	UnixSocket       *grpcoptions.UnixSocketOptions `json:"unix"   mapstructure:"unix"`
	ServerRunOptions *grpcoptions.ServerRunOptions  `json:"server" mapstructure:"server"`
}

// NewOptions creates a new Options object with default parameters.
func NewOptions() *Options {
	s := Options{
		Name: "example-gRPC",

		Log:              log.NewOptions(),
		TCP:              grpcoptions.NewTCPOptions(),
		UnixSocket:       grpcoptions.NewUnixSocketOptions(),
		ServerRunOptions: grpcoptions.NewServerRunOptions(),
	}

	return &s
}

// Flags returns flags for a specific server by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.Log.AddFlags(fss.FlagSet("logs"))
	o.TCP.AddFlags(fss.FlagSet("tcp"))
	o.UnixSocket.AddFlags(fss.FlagSet("unix"))
	o.ServerRunOptions.AddFlags(fss.FlagSet("server"))
	return fss
}

func (o *Options) String() string {
	// hide annoying cert data in log
	cert := o.TCP.ServerCert.CopyAndHide()
	data, _ := json.Marshal(o)
	o.TCP.ServerCert = *cert

	return string(data)
}

// Complete fills in any fields not set that are required to have valid data.
// 补全指定的选项.
func (o *Options) Complete() error {
	if err := o.TCP.Complete(); err != nil {
		return err
	}
	return nil
}
