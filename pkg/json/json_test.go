package json_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/wangweihong/eazycloud/pkg/json"
)

type Data struct {
	AA int `json:"aa"`
}

func TestUnmarshal(t *testing.T) {
	Convey("Unmarshal 大小写", t, func() {
		d1 := `{"aa":2324}`
		d2 := `{"Aa":2500}`
		d3 := `{"AA":2600}`
		d4 := `{"aA":2700}`

		var a Data
		So(json.Unmarshal([]byte(d1), &a), ShouldBeNil)
		So(a.AA, ShouldEqual, 2324)
		So(json.Unmarshal([]byte(d2), &a), ShouldBeNil)
		So(a.AA, ShouldEqual, 2500)
		So(json.Unmarshal([]byte(d3), &a), ShouldBeNil)
		So(a.AA, ShouldEqual, 2600)
		So(json.Unmarshal([]byte(d4), &a), ShouldBeNil)
		So(a.AA, ShouldEqual, 2700)
	})
}
