package maputil_test

import (
	"github.com/wangweihong/eazycloud/pkg/sets"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/wangweihong/eazycloud/pkg/util/maputil"
)

func TestStringIntMap_Init(t *testing.T) {
	Convey("TestStringIntMap_Init", t, func() {
		var nilMap map[string]int

		Convey("not nil", func() {
			So(nilMap, ShouldBeNil)

			maputil.StringIntMap(nilMap).Init()
			So(nilMap, ShouldBeNil)

			nilMap = maputil.StringIntMap(nilMap).Init()
			So(nilMap, ShouldNotBeNil)
		})
	})
}

func TestStringIntMap_Set(t *testing.T) {
	Convey("TestStringIntMap_Set", t, func() {
		var nilMap map[string]int

		Convey("not nil", func() {
			So(nilMap, ShouldBeNil)

			maputil.StringIntMap(nilMap).Set("1", 2)
			So(nilMap, ShouldBeNil)

			nilMap = maputil.StringIntMap(nilMap).Set("a", 3).Set("c", 3)
			So(nilMap, ShouldNotBeNil)
			So(len(nilMap), ShouldEqual, 2)
		})
	})
}

func TestStringIntMap_DeepCopy(t *testing.T) {
	Convey("TestStringIntMap_Set", t, func() {
		var nilMap map[string]int

		Convey("not nil", func() {
			So(nilMap, ShouldBeNil)

			nilMap = maputil.StringIntMap(nilMap).DeepCopy()
			So(nilMap, ShouldNotBeNil)

			nilMap = maputil.StringIntMap(nilMap).Set("a", 4)
			So(nilMap, ShouldNotBeNil)
			So(len(nilMap), ShouldEqual, 1)
			So(maputil.StringIntMap(nilMap).Has("a"), ShouldBeTrue)
		})
	})
}

func TestStringIntMap_Delete(t *testing.T) {
	Convey("TestStringIntMap_Delete", t, func() {
		Convey("nil", func() {
			var nilMap map[string]int
			maputil.StringIntMap(nilMap).Delete("a")
		})
		Convey("not nil", func() {
			d := make(map[string]int)
			d["a"] = 3

			maputil.StringIntMap(d).Delete("a")
			So(maputil.StringIntMap(d).Has("a"), ShouldBeFalse)
		})
	})
}

func TestStringIntMap_Get(t *testing.T) {
	Convey("TestStringIntMap_Get", t, func() {
		Convey("nil", func() {
			var nilMap map[string]int

			So(maputil.StringIntMap(nilMap).Get("notexist"), ShouldEqual, 0)
		})
		Convey("not nil", func() {
			d := make(map[string]int)
			d["a"] = 2

			So(maputil.StringIntMap(d).Get("a"), ShouldEqual, 2)
			So(maputil.StringIntMap(d).Get("notexist"), ShouldEqual, 0)
		})
	})
}

func TestStringIntMap_Keys(t *testing.T) {
	Convey("TestStringIntMap_Keys", t, func() {
		Convey("nil", func() {
			var nilMap map[string]int
			keys := maputil.StringIntMap(nilMap).Keys()

			So(len(keys), ShouldEqual, 0)
		})
		Convey("not nil", func() {
			d := make(map[string]int)
			d["a"] = 1
			d["b"] = 2

			keys := maputil.StringIntMap(d).Keys()
			So(len(keys), ShouldEqual, 2)
			So(sets.NewString(keys...).Equal(sets.NewString("a", "b")), ShouldBeTrue)
		})
	})
}
