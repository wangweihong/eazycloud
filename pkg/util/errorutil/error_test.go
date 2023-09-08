package errorutil_test

import (
	"fmt"
	"testing"

	"github.com/wangweihong/eazycloud/pkg/util/errorutil"
)

func TestErrorMsg(t *testing.T) {
	if errorutil.ErrorMsg(nil) != "" {
		t.Fatalf("no match")
	}

	if errorutil.ErrorMsg(fmt.Errorf("err")) != "err" {
		t.Fatalf("no match")
	}
}
