package tecentcloud

import (
	"context"

	"github.com/wangweihong/eazycloud/internal/pkg/libs"
)

func invoke[T any, R any](ctx context.Context, err error, req T, call func(T) (R, error), opts ...libs.RateLimitOption) (R, error) {
	return libs.InvokeApi(ctx, err, req, call, opts...)
}
