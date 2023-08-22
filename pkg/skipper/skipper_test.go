package skipper_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/wangweihong/eazycloud/pkg/skipper"
)

func TestSkipper(t *testing.T) {
	Convey("skipper", t, func() {
		s := skipper.AllowPathPrefixSkipper("/version", "/debug")
		So(skipper.Skip("/version/debug", s), ShouldBeTrue)
		So(skipper.Skip("/debug", s), ShouldBeTrue)
		So(skipper.Skip("/v1/test", s), ShouldBeFalse)

		ns := skipper.AllowPathPrefixNoSkipper("/version", "/debug")
		So(skipper.Skip("/version/debug", ns), ShouldBeFalse)
		So(skipper.Skip("/debug", ns), ShouldBeFalse)
		So(skipper.Skip("/v1/test", ns), ShouldBeTrue)
	})
}
