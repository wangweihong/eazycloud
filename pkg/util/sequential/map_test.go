package sequential_test

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wangweihong/eazycloud/pkg/util/sequential"
)

func TestNewSequentialMap(t *testing.T) {
	Convey("TestStringIntMap_Init", t, func() {
		nonexist := "nonexist"
		s := sequential.NewSequentialMap()
		s.Inject("a", 123)
		s.Inject("b", 321)
		s.Inject("c", 456)

		So(s.Len(), ShouldEqual, 3)
		So(s.Keys(), ShouldResemble, []interface{}{"a", "b", "c"})
		So(s.Values(), ShouldResemble, []interface{}{123, 321, 456})
		So(s.Has("a"), ShouldBeTrue)
		So(s.Has(nonexist), ShouldBeFalse)
		So(s.Get("a"), ShouldEqual, 123)
		So(s.Get(nonexist), ShouldEqual, nil)

		s.Delete(nonexist)
		So(s.Len(), ShouldEqual, 3)
		So(s.Keys(), ShouldResemble, []interface{}{"a", "b", "c"})
		So(s.Values(), ShouldResemble, []interface{}{123, 321, 456})

		s.Delete("b")
		So(s.Len(), ShouldEqual, 2)
		So(s.Keys(), ShouldResemble, []interface{}{"a", "c"})
		So(s.Values(), ShouldResemble, []interface{}{123, 456})

		err := s.ForEach(func(value interface{}) error {
			fmt.Println(value)
			return nil
		})
		So(err, ShouldBeNil)

		s.Inject("b", 321)
		So(s.Len(), ShouldEqual, 3)
		So(s.Keys(), ShouldResemble, []interface{}{"a", "c", "b"})
		So(s.Values(), ShouldResemble, []interface{}{123, 456, 321})

		s.Inject("b", 789)
		So(s.Len(), ShouldEqual, 3)
		So(s.Keys(), ShouldResemble, []interface{}{"a", "c", "b"})
		So(s.Values(), ShouldResemble, []interface{}{123, 456, 789})
	})
}
