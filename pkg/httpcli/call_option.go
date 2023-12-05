package httpcli

import (
	"net/http"
	"time"
)

type callInfo struct {
	timeout            *time.Duration
	header             map[string]string
	query              map[string]interface{}
	responseNotParse   bool
	httpRequestProcess func(req *http.Request) (*http.Request, error)
}

type CallOption func(*callInfo)

// TimeoutCallOption 设置某个连接超时操作.
func TimeoutCallOption(timeout time.Duration) CallOption {
	return func(c *callInfo) {
		c.timeout = &timeout
	}
}

// HeaderCallOption 设置某个请求的头部.
func HeaderCallOption(header map[string]string) CallOption {
	return func(c *callInfo) {
		c.header = header
	}
}

// QueryCallOption 设置某个连接查询参数.
func QueryCallOption(query map[string]interface{}) CallOption {
	return func(c *callInfo) {
		if c.query == nil {
			c.query = make(map[string]interface{})
		}
		for k, v := range query {
			c.query[k] = v
		}
	}
}

// OneQueryCallOption 设置某个连接查询参数.
func OneQueryCallOption(key string, value interface{}) CallOption {
	return func(c *callInfo) {
		if key == "" {
			return
		}
		if c.query == nil {
			c.query = make(map[string]interface{})
		}
		c.query[key] = value
	}
}

// ResponseNotParseCallOption 在invoke时不对http请求体数据进行解析.
func ResponseNotParseCallOption() CallOption {
	return func(c *callInfo) {
		c.responseNotParse = true
	}
}

type ProcessRequestFunc func(req *http.Request) (*http.Request, error)

// HttpRequestProcessOption 在http请求发起调用前，对http请求进行处理. 如根据url/请求头进行加密,并写入httpReq.
func HttpRequestProcessOption(fun ProcessRequestFunc) CallOption {
	return func(c *callInfo) {
		c.httpRequestProcess = fun
	}
}
