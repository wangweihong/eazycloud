package request

import (
	"context"
	"net/http"

	"fmt"
	"github.com/wangweihong/eazycloud/pkg/httpcli"
	"github.com/wangweihong/eazycloud/pkg/json"
)

type Request struct {
	Method   string
	URI      string
	Query    map[string]interface{}
	BodyData interface{}
	Header   map[string]string
	err  error
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

func (r *Request) SetQuery(v map[string]interface{}) *Request {
	r.Query = v
	return r
}

func (r *Request) SetQueryFromStruct(input interface{}) *Request {
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
		r.err= err
		return r
	}

	return r
}

func (r *Request)HttpRequest(ctx context.Context,c *httpcli.Client)(*http.Request,error){
	if r.err != nil {
		return nil, fmt.Errorf("request error:%v",r.err)
	}

	var opts []httpcli.CallOption
	opts = append(opts, httpcli.QueryCallOption(r.Query))
	opts = append(opts,httpcli.HeaderCallOption(r.Header))

	return httpcli.NewHttpRequest(ctx,c.GetAddr(),r.Method,r.URI,r.BodyData,opts...)
}


func (r *Request)Invoke(ctx context.Context,c *httpcli.Client,resp interface{})(*httpcli.RawResponse,error){
	if c == nil {
		return nil, fmt.Errorf("client is nil")
	}

	if r.err != nil {
		return nil, fmt.Errorf("request error:%v",r.err)
	}

	var opts []httpcli.CallOption
	opts = append(opts, httpcli.QueryCallOption(r.Query))
	opts = append(opts,httpcli.HeaderCallOption(r.Header))

	return  c.Invoke(
		ctx,
		r.Method,
		r.URI,
		r.BodyData,
		resp,
		opts...)
}
