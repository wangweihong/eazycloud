package sequential_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wangweihong/eazycloud/pkg/util/sequential"
)

func TestNewSequentialList(t *testing.T) {
	Convey("TestNewSequentialList", t, func() {
		s := sequential.NewSequentialList("a", "a", "b")
		So(s.Len(), ShouldEqual, 3)
		So(s.Has("a"), ShouldBeTrue)
		So(len(s.Indices("a")), ShouldEqual, 2)
		So(s.List(), ShouldResemble, []interface{}{"a", "a", "b"})

		var b []interface{}
		err := s.ForEach(func(f interface{}) error {
			b = append(b, f)
			return nil
		})
		So(err, ShouldBeNil)
		So(b, ShouldResemble, s.List())

		s.Inject("a")
		So(len(s.Indices("a")), ShouldEqual, 3)
		So(s.List(), ShouldResemble, []interface{}{"a", "a", "b", "a"})

		s.DeleteAtIndex(0)
		So(len(s.Indices("a")), ShouldEqual, 2)
		So(s.Indices("a"), ShouldResemble, []int{0, 2})
		So(s.Indices("b"), ShouldResemble, []int{1})
		So(s.List(), ShouldResemble, []interface{}{"a", "b", "a"})

		s.Delete("a")
		So(len(s.Indices("a")), ShouldEqual, 0)
		So(s.Has("a"), ShouldBeFalse)
		So(s.List(), ShouldResemble, []interface{}{"b"})
	})
}

func TestNewLimitSequentialList(t *testing.T) {
	Convey("TestNewLimitSequentialMap", t, func() {
		s := sequential.NewLimitSequentialList(2, "a", "b", "c")
		So(s.Len(), ShouldEqual, 2)
		So(s.Has("a"), ShouldBeFalse)
		So(s.Has("b"), ShouldBeTrue)
		So(s.Has("c"), ShouldBeTrue)
		So(s.Get(1), ShouldResemble, "c")
		So(s.Get(2), ShouldBeNil)

		s.Clear()
		s.Inject("a")
		s.Inject("a")
		So(s.Len(), ShouldEqual, 2)
		So(s.List(), ShouldResemble, []interface{}{"a", "a"})

		s.Delete("a")
		So(s.Len(), ShouldEqual, 0)

		s.Clear()
		So(s.Inject("a"), ShouldEqual, 0)
		So(s.Inject("a"), ShouldEqual, 1)
		So(s.Inject("a"), ShouldEqual, 1)
		So(s.Indices("a"), ShouldResemble, []int{0, 1})
	})
}
