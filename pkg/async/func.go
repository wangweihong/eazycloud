package async

import (
	"context"
	"runtime/debug"

	"github.com/wangweihong/eazycloud/pkg/log"
)

func PanicRecover(ctx context.Context, fns ...func()) {
	if x := recover(); x != nil {
		log.Errorf("run time panic: %v %v", x, string(debug.Stack()))
		for _, fn := range fns {
			fn()
		}
	}
}

func GoRoutine(ctx context.Context, fs ...func(ctx context.Context)) {
	for i := range fs {
		f := fs[i]
		go func() {
			defer PanicRecover(ctx)
			f(ctx)
		}()
	}
}
