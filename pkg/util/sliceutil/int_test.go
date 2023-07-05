package sliceutil_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/wangweihong/eazycloud/pkg/util/sliceutil"
	"testing"
)

func TestIntSlice_DeepCopy(t *testing.T) {
	Convey("TestIntSlice_DeepCopy", t, func() {
		Convey("nil", func() {
			var d []int
			b := sliceutil.IntSlice(d).DeepCopy()

			So(d, ShouldBeNil)
			So(b, ShouldNotBeNil)
		})

		Convey("not nil", func() {
			d := []int{1}
			b := sliceutil.IntSlice(d).DeepCopy()
			d = append(d, 2)

			So(b, ShouldNotBeNil)
			So(len(b), ShouldEqual, 1)
		})
	})
}
