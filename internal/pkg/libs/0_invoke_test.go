package libs_test

import (
	"io/ioutil"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

type RegistryOption struct {
	Endpoint string `json:"endpoint" mapstructure:"endpoint"`
	Repo     string `json:"repo"  mapstructure:"repo"`
	Tag      string `json:"tag"  mapstructure:"tag"`
}

type Option struct {
	Registry *RegistryOption `json:"registry" mapstructure:"registry"`
}

func load() *Option {
	errors.UpdateModuleInfo(errors.NewModuleGetter("github.com/wangweihong/eazycloud", "127.0.01", 12345))

	yamlConfig, err := ioutil.ReadFile("./testdata/config.yaml")
	if err != nil {
		panic(err)
	}
	v := viper.New()
	v.SetConfigType("yaml")
	// 读取 YAML 配置
	if err := v.ReadConfig(strings.NewReader(string(yamlConfig))); err != nil {
		panic(err)
	}

	var opt Option
	if err := v.Unmarshal(&opt); err != nil {
		panic(err)
	}
	return &opt
}

func TestLoad(t *testing.T) {
	Convey("TestLoad", t, func() {
		So(load(), ShouldNotBeNil)
	})

}
