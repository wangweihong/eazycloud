package maputil_test

import (
	"testing"

	"github.com/wangweihong/eazycloud/pkg/sets"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/wangweihong/eazycloud/pkg/util/maputil"
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

func TestStringStringMap_Get(t *testing.T) {
	Convey("TestStringStringMap_Get", t, func() {
		Convey("nil", func() {
			var nilMap map[string]string

			So(maputil.StringStringMap(nilMap).Get("notexist"), ShouldEqual, "")
		})
		Convey("not nil", func() {
			d := make(map[string]string)
			d["a"] = "b"

			So(maputil.StringStringMap(d).Get("a"), ShouldEqual, "b")
			So(maputil.StringStringMap(d).Get("notexist"), ShouldEqual, "")
		})
	})
}

func TestStringStringMap_Keys(t *testing.T) {
	Convey("TestStringStringMap_Keys", t, func() {
		Convey("nil", func() {
			var nilMap map[string]string
			keys := maputil.StringStringMap(nilMap).Keys()

			So(len(keys), ShouldEqual, 0)
		})
		Convey("not nil", func() {
			d := make(map[string]string)
			d["a"] = "1"
			d["b"] = "2"

			keys := maputil.StringStringMap(d).Keys()
			So(len(keys), ShouldEqual, 2)
			So(sets.NewString(keys...).Equal(sets.NewString("a", "b")), ShouldBeTrue)
		})
	})
}
