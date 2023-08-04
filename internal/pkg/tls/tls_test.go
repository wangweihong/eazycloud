package tls_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/wangweihong/eazycloud/internal/pkg/tls"
)

func TestGeneratableKeyCert_Validate(t *testing.T) {
	Convey("校验证书", t, func() {
		c := tls.GeneratableKeyCert{
			CertData:      tls.CertData{},
			CertKey:       tls.CertKey{},
			CertDirectory: "",
			PairName:      "",
		}

		So(c.Validate(), ShouldNotBeNil)
		c.CertDirectory = "./"
		So(c.Validate(), ShouldNotBeNil)
		c.PairName = "test"
		So(c.Validate(), ShouldBeNil)
		c.CertKey.CertFile = "./test.crt"
		So(c.Validate(), ShouldNotBeNil)
		c.CertKey.KeyFile = "./test.key"
		So(c.Validate(), ShouldBeNil)
		c.CertData.Cert = "xxxx"
		So(c.Validate(), ShouldNotBeNil)
		c.CertData.Key = "xxx"
		So(c.Validate(), ShouldBeNil)
	})
}
