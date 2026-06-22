package libkubernetes_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wangweihong/gotoolbox/pkg/json"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/wangweihong/eazycloud/apis/ikubernetes"
	"github.com/wangweihong/eazycloud/internal/pkg/libkubernetes"
)

var (
	config *ikubernetes.ClusterConfig
)

func init() {
	f, err := ioutil.ReadFile("./testdata/kubeconfig")
	if err != nil {
		panic(err)
	}

	config, err = libkubernetes.NewClusterConfigFromKubeconfig(string(f))
	if err != nil {
		panic(err)
	}

}
func TestServerVersion(t *testing.T) {
	Convey("TestServerVersion", t, func() {
		v, err := libkubernetes.ServerVersion(context.Background(), config)
		So(err, ShouldBeNil)

		fmt.Println(v)
	})
}

func TestPodList(t *testing.T) {
	Convey("TestPodList", t, func() {
		v, err := libkubernetes.PodList(context.Background(), config, "kube-system", v1.ListOptions{})
		So(err, ShouldBeNil)
		So(len(v.Items), ShouldNotEqual, 0)
		json.PrintObject(v.Items[0])

		So(v.Items[0].Kind, ShouldEqual, "Pod")
		So(v.Items[0].APIVersion, ShouldEqual, "v1")
	})
}
