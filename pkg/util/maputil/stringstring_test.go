package maputil_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/wangweihong/eazycloud/pkg/util/maputil"
	"testing"
)

func TestStringStringMap_Init(t *testing.T) {
	Convey("TestStringStringMap_Init", t, func() {
		var nilMap map[string]string

		Convey("not nil", func() {
			So(nilMap, ShouldBeNil)

			maputil.StringStringMap(nilMap).Init()
			So(nilMap, ShouldBeNil)

			nilMap = maputil.StringStringMap(nilMap).Init()
			So(nilMap, ShouldNotBeNil)
		})
	})
}

func TestStringStringMap_Set(t *testing.T) {
	Convey("TestStringStringMap_Set", t, func() {
		var nilMap map[string]string

		Convey("not nil", func() {
			So(nilMap, ShouldBeNil)

			maputil.StringStringMap(nilMap).Set("1", "2")
			So(nilMap, ShouldBeNil)

			nilMap = maputil.StringStringMap(nilMap).Set("a", "b").Set("c", "D")
			So(nilMap, ShouldNotBeNil)
			So(len(nilMap), ShouldEqual, 2)
		})
	})
}

func TestStringStringMap_DeepCopy(t *testing.T) {
	Convey("TestStringStringMap_Set", t, func() {
		var nilMap map[string]string

		Convey("not nil", func() {
			So(nilMap, ShouldBeNil)

			nilMap = maputil.StringStringMap(nilMap).DeepCopy()
			So(nilMap, ShouldNotBeNil)

			nilMap = maputil.StringStringMap(nilMap).Set("a", "b")
			So(nilMap, ShouldNotBeNil)
			So(len(nilMap), ShouldEqual, 1)
			So(maputil.StringStringMap(nilMap).Has("a"), ShouldBeTrue)
		})
	})
}

func TestStringStringMap_Delete(t *testing.T) {
	Convey("TestStringStringMap_Delete", t, func() {
		Convey("nil", func() {
			var nilMap map[string]string
			maputil.StringStringMap(nilMap).Delete("a")
		})
		Convey("not nil", func() {
			d := make(map[string]string)
			d["a"] = "b"

			maputil.StringStringMap(d).Delete("a")
			So(maputil.StringStringMap(d).Has("a"), ShouldBeFalse)
		})
	})
}
