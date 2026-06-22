package libs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"
	"github.com/wangweihong/gotoolbox/pkg/httpcli/interceptorcli"
	"github.com/wangweihong/gotoolbox/pkg/rate"
	"github.com/wangweihong/gotoolbox/pkg/validation"
)

func DefaultCallOptions() []httpcli.Option {
	return []httpcli.Option{
		httpcli.WithTimeout(30 * time.Second),
		httpcli.WithIntercepts(
			interceptorcli.DecodeResponseInterceptor("decode"),
			interceptorcli.StatusCodeInterceptor("httpstatus"),
		),
	}
}

func Invoke[R any, T any](ctx context.Context, err error, c *httpcli.Client, httpReq *httpcli.HttpRequest, req R, resp T, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	return invoke(ctx, err, c, httpReq, req, resp, opts...)
}

func invoke[R any, T any](ctx context.Context, err error, c *httpcli.Client, httpReq *httpcli.HttpRequest, req R, resp T, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if any(req) != nil {
		if val, ok := any(req).(validation.Validator); ok {
			if err := val.Validate(); err != nil {
				return nil, errors.WithStack(err)
			}
		}
	}

	httpResp, err := c.Invoke(ctx, httpReq, req, resp, opts...)
	return httpResp, errors.WithStack(err)
}

func InvokeApi[T any, R any](ctx context.Context, err error, req T, call func(T) (R, error), opts ...RateLimitOption) (R, error) {
	var zero R
	if err != nil {
		return zero, errors.WithStack(err)
	}
	// 获取请求类型作为限流键
	reqType := getRequestType(req)
	// 应用限流选项
	limiter := getOrCreateLimiter(reqType, opts...)
	// 等待令牌可用（支持上下文取消）
	if err := limiter.Wait(ctx); err != nil {
		return zero, errors.WithStack(err)
	}

	return call(req)
}

// 全局限流器映射，按请求类型存储限流器
var (
	rateLimiters     = make(map[string]*rate.Limiter)
	rateLimitersLock sync.RWMutex
)

// 获取或创建限流器
func getOrCreateLimiter(key string, opts ...RateLimitOption) *rate.Limiter {
	rateLimitersLock.RLock()
	limiter, exists := rateLimiters[key]
	rateLimitersLock.RUnlock()

	if exists {
		return limiter
	}

	// 默认配置
	config := &rateLimiterConfig{
		rps:   10, // 默认10 QPS
		burst: 20, // 默认突发容量20
	}

	// 应用选项配置
	for _, opt := range opts {
		opt(config)
	}

	// 创建新限流器
	newLimiter := rate.NewLimiter(rate.Limit(config.rps), config.burst)

	rateLimitersLock.Lock()
	defer rateLimitersLock.Unlock()

	// 双重检查防止重复创建
	if limiter, exists := rateLimiters[key]; exists {
		return limiter
	}

	rateLimiters[key] = newLimiter
	return newLimiter
}

// 获取请求类型标识
func getRequestType(req any) string {
	// 使用反射获取类型名称作为唯一标识
	// 实际应用中可根据需求使用更复杂的标识逻辑
	return fmt.Sprintf("%T", req)
}

// 限流配置选项
type RateLimitOption func(*rateLimiterConfig)

// 可选配置函数
func WithRateLimit(rps, burst int) RateLimitOption {
	return func(c *rateLimiterConfig) {
		c.rps = rps
		c.burst = burst
	}
}

type rateLimiterConfig struct {
	rps   int
	burst int
}

// 动态更新特定请求类型的限流配置
func UpdateRateLimit(reqType string, rps, burst int) {
	rateLimitersLock.Lock()
	defer rateLimitersLock.Unlock()

	if limiter, exists := rateLimiters[reqType]; exists {
		limiter.SetLimit(rate.Limit(rps))
		limiter.SetBurst(burst)
	} else {
		rateLimiters[reqType] = rate.NewLimiter(rate.Limit(rps), burst)
	}
}

// 获取当前限流统计
func GetRateLimitStats(reqType string) (currentRPS float64, currentBurst int) {
	rateLimitersLock.RLock()
	defer rateLimitersLock.RUnlock()

	if limiter, exists := rateLimiters[reqType]; exists {
		return float64(limiter.Limit()), limiter.Burst()
	}
	return 0, 0
}
