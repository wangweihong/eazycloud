package randutil_test

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wangweihong/eazycloud/pkg/util/randutil"
)

func TestRandString(t *testing.T) {
	Convey("randString", t, func() {
		fmt.Println(randutil.RandString([]rune("abcdefghijklmnopqrstuvwxyz12345678"), 32))
	})
}
