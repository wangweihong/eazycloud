package errors_test

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/wangweihong/eazycloud/pkg/errors"
)

func TestAggregateError(t *testing.T) {
	Convey("Aggregates", t, func() {
		Convey("e", func() {
			e1 := errors.Wrap(101, "error1")
			e2 := errors.Wrap(101, "error2")

			fmt.Println(errors.NewAggregate(e1, e2).Error())
		})
	})
}
