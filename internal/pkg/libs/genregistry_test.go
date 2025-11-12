package libs_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wangweihong/eazycloud/apis/igenregistry"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"
	"github.com/wangweihong/gotoolbox/pkg/httpcli/interceptorcli"
	"github.com/wangweihong/gotoolbox/pkg/sliceutil"
)

func TestGenRegistry_ListImages(t *testing.T) {
	Convey("TestGenRegistry_ListRepo", t, func() {
		opt := load()
		So(opt, ShouldNotBeNil)
		So(opt.Registry, ShouldNotBeNil)

		resp := &igenregistry.ListImagesResponse{}
		opts := sliceutil.Append(libs.DefaultCallOptions(), httpcli.WithInsecure())
		httpResp, err := libs.NewGenenralRegistryClient(opt.Registry.Endpoint, opts...).ListImages(context.Background(), nil, resp)
		//	log.Errorf("%#+v", err)
		So(err, ShouldBeNil)
		t.Log(httpResp.GetBody())
		t.Log(resp)
	})
}

func TestGenRegistry_GetImageTags(t *testing.T) {
	Convey("TestGenRegistry_GetImageTags", t, func() {
		opt := load()
		So(opt, ShouldNotBeNil)
		So(opt.Registry, ShouldNotBeNil)

		resp := &igenregistry.GetImageTagsResponse{}
		opts := sliceutil.Append(libs.DefaultCallOptions(), httpcli.WithInsecure())
		httpResp, err := libs.NewGenenralRegistryClient(opt.Registry.Endpoint, opts...).GetImageTags(context.Background(), &igenregistry.GetImageTagsRequest{Repo: opt.Registry.Repo}, resp)
		//	log.Errorf("%#+v", err)
		So(err, ShouldBeNil)
		t.Log(httpResp.GetBody())
		t.Log(resp)
	})
}

func TestGenRegistry_GetImageManifests(t *testing.T) {
	Convey("TestGenRegistry_GetImageManifests", t, func() {
		opt := load()
		So(opt, ShouldNotBeNil)
		So(opt.Registry, ShouldNotBeNil)

		resp := &igenregistry.GetImageManifestsResponse{}
		opts := []httpcli.Option{httpcli.WithIntercepts(
			interceptorcli.DecodeResponseInterceptor("decode"),
			interceptorcli.StatusCodeInterceptor("httpstatus"),
			interceptorcli.LoggingInterceptor("logging"),
		), httpcli.WithInsecure()}

		httpResp, err := libs.NewGenenralRegistryClient(opt.Registry.Endpoint, opts...).GetImageManifests(context.Background(), &igenregistry.GetImageManifestsRequest{Repo: opt.Registry.Repo, Tag: opt.Registry.Tag}, resp)
		//	log.Errorf("%#+v", err)
		So(err, ShouldBeNil)
		t.Log(httpResp.GetBody())
		t.Log(resp)
	})
}
