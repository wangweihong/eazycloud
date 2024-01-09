package httpcli

import (
	"net/http"
	"time"
)

type callInfo struct {
	timeout            *time.Duration
	header             http.Header
	query              map[string]interface{}
	responseNotParse   bool
	httpRequestProcess func(req *http.Request) (*http.Request, error)
	urlSetter          func() (string, error)
	data               interface{}
	endpoint           string
	// 拦截器列表
	chainInterceptors []Interceptor
}

type CallOption func(*callInfo)

type CallOptions []CallOption

func (cs CallOptions) Duplicate() []CallOption {
	if cs == nil {
		return nil
	}

	n := make([]CallOption, 0, len(cs))
	for _, v := range cs {
		n = append(n, v)
	}
	return n
}

// TimeoutCallOption 设置某个连接超时操作.
func TimeoutCallOption(timeout time.Duration) CallOption {
	return func(c *callInfo) {
		c.timeout = &timeout
	}
}

// SetHeaderCallOption 设置请求头部，替换原来头部
func SetHeaderCallOption(header http.Header) CallOption {
	return func(c *callInfo) {
		if c.header == nil {
			c.header = make(map[string][]string)
		}

		c.header = header
	}
}

// SetHeaderValueCallOption 设置请求头部指定值
func SetHeaderValueCallOption(key string, value ...string) CallOption {
	return func(c *callInfo) {
		if c.header == nil {
			c.header = make(map[string][]string)
		}

		c.header.Del(key)
		for _, v := range value {
			c.header.Add(key, v)
		}
	}
}

// AddHeaderCallOption 增加某个请求的头部.
func AddHeaderCallOption(header http.Header) CallOption {
	return func(c *callInfo) {
		if c.header == nil {
			c.header = make(map[string][]string)
		}

		for k, values := range header {
			for _, v := range values {
				c.header.Add(k, v)
			}
		}
	}
}

// AddHeaderValueCallOption 增加某个请求的头部.
func AddHeaderValueCallOption(key string, value ...string) CallOption {
	return func(c *callInfo) {
		if c.header == nil {
			c.header = make(map[string][]string)
		}

		for _, v := range value {
			c.header.Add(key, v)
		}
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

type URLSetter func() (string, error)

// 有可能需要根据资源/rawURL动态更改请求URL
func URLCallOption(epf URLSetter) CallOption {
	return func(c *callInfo) {
		c.urlSetter = epf
	}
}

// 设置一些特殊处理的数据
func DataCallOption(data interface{}) CallOption {
	return func(c *callInfo) {
		c.data = data
	}
}

// 更改访问的服务器端点
func EndpointCallOption(endpoint string) CallOption {
	return func(c *callInfo) {
		c.endpoint = endpoint
	}
}

// 更改访问的拦截器列表
func InterceptorsCallOption(chainInterceptors []Interceptor) CallOption {
	return func(c *callInfo) {
		c.chainInterceptors = chainInterceptors
	}
}
