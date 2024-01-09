package httpcli

import (
	"context"
	"net/http"

	"fmt"
	"github.com/wangweihong/eazycloud/pkg/json"
)

type Request struct {
	Method   string
	URI      string
	Query    map[string]interface{}
	BodyData interface{}
	Header   http.Header
	err      error
}

func NewRequest() *Request {
	return &Request{}
}

func (r *Request) POST() *Request {
	r.Method = "POST"
	return r
}

func (r *Request) GET() *Request {
	r.Method = "GET"
	return r
}

func (r *Request) PUT() *Request {
	r.Method = "PUT"
	return r
}

func (r *Request) DELETE() *Request {
	r.Method = "DELETE"
	return r
}

func (r *Request) SetURI(uri string) *Request {
	r.URI = uri
	return r
}

func (r *Request) SetParam(v interface{}) *Request {
	r.BodyData = v
	return r
}

func (r *Request) SetQueryRaw(v map[string]interface{}) *Request {
	r.Query = v
	return r
}

func (r *Request) SetHeaders(v http.Header) *Request {
	r.Header = v
	return r
}

func (r *Request) SetHeader(key string, value ...string) *Request {
	if r.Header == nil {
		r.Header = make(map[string][]string)
	}

	if key == "" {
		return r
	}
	r.Header.Del(key)
	for _, v := range value {
		r.Header.Add(key, v)
	}
	return r
}

func (r *Request) AddHeader(key string, value ...string) *Request {
	if r.Header == nil {
		r.Header = make(map[string][]string)
	}

	if key == "" {
		return r
	}

	for _, v := range value {
		r.Header.Add(key, v)
	}
	return r
}

func (r *Request) SetQuery(input interface{}) *Request {
	if input == nil {
		return r
	}
	jsonData, err := json.Marshal(input)
	if err != nil {
		r.err = err
		return r
	}

	// 将 JSON 字符串解码为 map
	var resultMap map[string]interface{}
	err = json.Unmarshal(jsonData, &resultMap)
	if err != nil {
		r.err = err
		return r
	}
	r.Query = resultMap
	return r
}

func (r *Request) HttpRequest(ctx context.Context, c *Client) (*http.Request, error) {
	if r.err != nil {
		return nil, fmt.Errorf("request error:%v", r.err)
	}

	var opts []CallOption
	opts = append(opts, QueryCallOption(r.Query))
	opts = append(opts, AddHeaderCallOption(r.Header))

	return NewHttpRequest(ctx, c.GetAddr(), r.Method, r.URI, r.BodyData, opts...)
}

func (r *Request) Invoke(ctx context.Context, c *Client, resp interface{}, opts ...CallOption) (*RawResponse, error) {
	if c == nil {
		return nil, fmt.Errorf("client is nil")
	}

	if r.err != nil {
		return nil, fmt.Errorf("request error:%v", r.err)
	}

	opts = append(opts, QueryCallOption(r.Query))
	opts = append(opts, AddHeaderCallOption(r.Header))

	return c.Invoke(
		ctx,
		r.Method,
		r.URI,
		r.BodyData,
		resp,
		opts...)
}
