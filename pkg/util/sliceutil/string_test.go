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

func TestStringSlice_GetRepeat(t *testing.T) {
	Convey("TestStringSlice_GetRepeat", t, func() {
		var nilS []string
		var rm map[string]int
		var repeated bool

		rm, repeated = sliceutil.StringSlice(nilS).GetRepeat()
		So(rm, ShouldBeNil)
		So(repeated, ShouldBeFalse)

		rm, repeated = sliceutil.StringSlice([]string{"a", "a", "a", "b"}).GetRepeat()
		So(rm, ShouldNotBeNil)
		So(repeated, ShouldBeTrue)
		d, _ := rm["a"]
		So(d, ShouldEqual, 3)

		rm, repeated = sliceutil.StringSlice([]string{"b", "a"}).GetRepeat()
		So(rm, ShouldBeNil)
		So(repeated, ShouldBeFalse)
	})
}

func TestStringSlice_Sort(t *testing.T) {
	Convey("TestStringSlice_Sort", t, func() {
		var nilS []string
		So(sliceutil.StringSlice(nilS).SortAsc(), ShouldBeNil)
		So(sliceutil.StringSlice([]string{"a", "c", "b"}).SortAsc(), ShouldResemble, []string{"a", "b", "c"})
		So(sliceutil.StringSlice([]string{"a", "c", "b"}).SortAsc(), ShouldNotResemble, []string{"a", "c", "b"})

		So(sliceutil.StringSlice(nilS).SortDesc(), ShouldBeNil)
		So(sliceutil.StringSlice([]string{"a", "c", "b"}).SortDesc(), ShouldResemble, []string{"c", "b", "a"})
		So(sliceutil.StringSlice([]string{"a", "c", "b"}).SortDesc(), ShouldNotResemble, []string{"a", "c", "b"})

	})
}

func TestStringSlice_HasEmpty(t *testing.T) {
	Convey("TestStringSlice_HasEmpty", t, func() {
		var nilS []string
		num, hasEmpty := sliceutil.StringSlice(nilS).HasEmpty()
		So(num, ShouldEqual, 0)
		So(hasEmpty, ShouldBeFalse)

		num, hasEmpty = sliceutil.StringSlice([]string{"", ""}).HasEmpty()
		So(num, ShouldEqual, 2)
		So(hasEmpty, ShouldBeTrue)

		num, hasEmpty = sliceutil.StringSlice([]string{"a", "b"}).HasEmpty()
		So(num, ShouldEqual, 0)
		So(hasEmpty, ShouldBeFalse)

		num, hasEmpty = sliceutil.StringSlice([]string{"", "b"}).HasEmpty()
		So(num, ShouldEqual, 1)
		So(hasEmpty, ShouldBeTrue)
	})
}
