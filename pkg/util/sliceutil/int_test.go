package sliceutil_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/wangweihong/eazycloud/pkg/util/sliceutil"
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

func TestIntSlice_HasRepeat(t *testing.T) {
	Convey("TestIntSlice_HasRepeat", t, func() {
		var nilS []int
		So(sliceutil.IntSlice(nilS).HasRepeat(), ShouldBeFalse)
		So(sliceutil.IntSlice([]int{123, 123}).HasRepeat(), ShouldBeTrue)
		So(sliceutil.IntSlice([]int{123, 245}).HasRepeat(), ShouldBeFalse)
	})
}

func TestIntSlice_GetRepeat(t *testing.T) {
	Convey("TestIntSlice_GetRepeat", t, func() {
		var nilS []int
		var rm map[int]int
		var repeated bool

		rm, repeated = sliceutil.IntSlice(nilS).GetRepeat()
		So(rm, ShouldBeNil)
		So(repeated, ShouldBeFalse)

		rm, repeated = sliceutil.IntSlice([]int{12, 12, 12, 3}).GetRepeat()
		So(rm, ShouldNotBeNil)
		So(repeated, ShouldBeTrue)
		d, _ := rm[12]
		So(d, ShouldEqual, 3)

		rm, repeated = sliceutil.IntSlice([]int{12, 3}).GetRepeat()
		So(rm, ShouldBeNil)
		So(repeated, ShouldBeFalse)
	})
}

func TestIntSlice_Sort(t *testing.T) {
	Convey("TestIntSlice_Sort", t, func() {
		var nilS []int
		So(sliceutil.IntSlice(nilS).SortAsc(), ShouldBeNil)
		So(sliceutil.IntSlice([]int{1, 3, 2}).SortAsc(), ShouldResemble, []int{1, 2, 3})
		So(sliceutil.IntSlice([]int{1, 3, 2}).SortAsc(), ShouldNotResemble, []int{1, 3, 2})

		So(sliceutil.IntSlice(nilS).SortDesc(), ShouldBeNil)
		So(sliceutil.IntSlice([]int{1, 3, 2}).SortDesc(), ShouldResemble, []int{3, 2, 1})
		So(sliceutil.IntSlice([]int{1, 3, 2}).SortDesc(), ShouldNotResemble, []int{1, 3, 2})
	})
}
