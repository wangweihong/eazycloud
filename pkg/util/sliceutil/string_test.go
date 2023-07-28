package sliceutil_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/wangweihong/eazycloud/pkg/util/sliceutil"
)

func TestStringSlice_DeepCopy(t *testing.T) {
	Convey("TestStringSlice_DeepCopy", t, func() {
		Convey("nil", func() {
			var d []string
			b := sliceutil.StringSlice(d).DeepCopy()

			So(d, ShouldBeNil)
			So(b, ShouldNotBeNil)
		})

		Convey("not nil", func() {
			d := []string{"aaa"}
			b := sliceutil.StringSlice(d).DeepCopy()
			d = append(d, "c")

			So(b, ShouldNotBeNil)
			So(len(b), ShouldEqual, 1)
		})
	})
}

func TestStringSlice_HasRepeat(t *testing.T) {
	Convey("TestStringSlice_HasRepeat", t, func() {
		var nilS []string
		So(sliceutil.StringSlice(nilS).HasRepeat(), ShouldBeFalse)
		So(sliceutil.StringSlice([]string{"a", "a"}).HasRepeat(), ShouldBeTrue)
		So(sliceutil.StringSlice([]string{"b", "a"}).HasRepeat(), ShouldBeFalse)
	})
}
