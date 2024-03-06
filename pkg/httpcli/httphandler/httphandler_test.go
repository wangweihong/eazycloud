package httphandler_test

import (
	"net/http"
	"testing"

	"github.com/wangweihong/eazycloud/pkg/httpcli/httphandler/loghandler"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wangweihong/eazycloud/pkg/httpcli/httphandler"
	"github.com/wangweihong/eazycloud/pkg/httpcli/httphandler/multihandler"
)

func TestNewHttpHandler(t *testing.T) {
	Convey("", t, func() {
		signHandler := func(r *http.Request) error {
			if r.Header == nil {
				r.Header = make(http.Header)
			}
			r.Header.Add("abc", "123")
			return nil
		}

		reqHandlers := multihandler.NewRequestHandlers(
			signHandler,
			loghandler.RequestHandler,
		)

		r := httphandler.NewHttpHandler().AddRequestHandler(reqHandlers.RequestHandlers)
		req := &http.Request{}
		err := r.RequestHandlers(req)
		So(err, ShouldBeNil)
		So(req.Header.Get("abc"), ShouldEqual, "123")
	})
}
