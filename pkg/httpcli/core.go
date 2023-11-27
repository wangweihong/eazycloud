package httpcli

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/wangweihong/eazycloud/pkg/code"
	"github.com/wangweihong/eazycloud/pkg/errors"
	"github.com/wangweihong/eazycloud/pkg/log"
	"github.com/wangweihong/eazycloud/pkg/tls/httptls"
	"github.com/wangweihong/eazycloud/pkg/util/callerutil"
)

type Client struct {
	conn            *http.Client
	transport       *http.Transport
	timeout         *time.Duration
	addr            string
	report          bool
	tlsEnabled      bool
	skipTlsVerified bool
	// gRPC服务证书
	serverCA       string
	mtlsEnabled    bool
	clientKeyData  string
	clientCertData string
	// 通用调用参数
	callOpts []CallOption
	// 拦截器列表
	chainInterceptors []Interceptor
}

func NewClient(addr string, options ...Option) (*Client, error) {
	c := &Client{
		addr: addr,
	}
	for _, o := range options {
		o(c)
	}

	if err := c.validate(); err != nil {
		return nil, errors.WrapError(code.ErrHTTPClientGenerateError, err)
	}

	return c, nil
}

func (c *Client) GetAddr() string {
	return c.addr
}

func (c *Client) validate() error {
	if c.addr == "" {
		return fmt.Errorf("client addr is empty")
	}

	if c.tlsEnabled {
		if !c.skipTlsVerified {
			if c.serverCA == "" {
				return fmt.Errorf("must set serverCA when tlsEnabled and not skipTlsVerified")
			}
		}
	}

	if c.mtlsEnabled {
		if c.clientKeyData == "" || c.clientCertData == "" {
			return fmt.Errorf("must provide clientKeyPEMData and clientCertPEMData when enable mTls")
		}

		if c.serverCA == "" {
			return fmt.Errorf("must set serverCA when mtlsEnabled enable")
		}
	}
	return nil
}

func (c *Client) getConn(ctx context.Context) (*http.Client, error) {
	// reuse conn
	// lock?
	if c.conn != nil {
		log.F(ctx).Debug("connection has exist, reuse it.")
		return c.conn, nil
	}

	var creds *tls.Config

	if c.tlsEnabled {
		var err error

		if c.skipTlsVerified {
			creds, err = httptls.NewTlsClientSkipVerifiedCredentials()
		} else {
			if c.mtlsEnabled {
				// 如果开启双向认证,需要加载服务器
				creds, err = httptls.NewMutualTlsClientCredentials([]byte(c.serverCA), []byte(c.clientCertData), []byte(c.clientKeyData))
			} else {
				creds, err = httptls.NewTlsClientCredentials([]byte(c.serverCA))
			}
		}
		if err != nil {
			log.F(ctx).Errorf("generate tls credential fail:%w ", err)
			return nil, errors.WrapError(code.ErrHTTPClientGenerateError, err)
		}
	}
	tr := &http.Transport{}
	if c.transport != nil {
		tr = c.transport
	}

	if creds != nil {
		tr.TLSClientConfig = creds
	}

	conn := http.Client{
		Transport: tr,
	}
	c.conn = &conn
	return c.conn, nil
}

type Interceptor func(ctx context.Context, method string, rawURL string, arg, reply interface{}, cc *Client, invoker Invoker, opts ...CallOption) (*RawResponse, error)

func (c *Client) Invoke(
	/*
			注意1. ctx的生命周期,如果Client在一个服务请求的异步动作,不能直接使用服务请求的ctx。否则当服务器结束后,context撤销
		会导致所有Client的请求都会立即`context cancel`失败返回
			注意2. log fieldCtx的信息覆盖问题. 如果ctx来自一个服务请求,服务的中间件可能在该请求中设置了请求相关信息。如果client
		也使用这个ctx, 在记录信息，注意字段的覆盖和冗余问题.
	*/
	ctx context.Context,
	method string,
	rawURL string,
	arg, reply interface{},
	opts ...CallOption,
) (*RawResponse, error) {
	opts = combine(c.callOpts, opts)
	file, line, fn := callerutil.CallerDepth(2)
	callerMsg := fmt.Sprintf("%s:%s:%d", file, fn, line)

	if c.chainInterceptors != nil {
		rawResp, err := c.chainInterceptors[0](
			ctx,
			method,
			rawURL,
			arg,
			reply,
			c,
			getChainUnaryInvoker(c.chainInterceptors, 0, invoke),
			opts...)

		log.F(ctx).
			Debug("Interceptor Invoked called.", log.String("caller", callerMsg), log.Err(err), log.Every("arg", arg), log.Every("reply", reply))
		return rawResp, err
	}

	rawResp, err := invoke(ctx, method, rawURL, arg, reply, c, opts...)
	log.F(ctx).
		Debug("Invoked called.", log.String("caller", callerMsg), log.Err(err), log.Every("arg", arg), log.Every("reply", reply))
	return rawResp, err
}

func combine(o1 []CallOption, o2 []CallOption) []CallOption {
	// we don't use append because o1 could have extra capacity whose
	// elements would be overwritten, which could cause inadvertent
	// sharing (and race conditions) between concurrent calls
	if len(o1) == 0 {
		return o2
	} else if len(o2) == 0 {
		return o1
	}
	ret := make([]CallOption, len(o1)+len(o2))
	copy(ret, o1)
	copy(ret[len(o1):], o2)
	return ret
}

func getChainUnaryInvoker(interceptors []Interceptor, curr int, finalInvoker Invoker) Invoker {
	if curr == len(interceptors)-1 {
		return finalInvoker
	}
	return func(ctx context.Context, method string, rawURL string, req, reply interface{}, cc *Client, opts ...CallOption) (*RawResponse, error) {
		return interceptors[curr+1](
			ctx,
			method,
			rawURL,
			req,
			reply,
			cc,
			getChainUnaryInvoker(interceptors, curr+1, finalInvoker),
			opts...)
	}
}

type RawResponse struct {
	Header     http.Header
	Body       []byte
	Cookies    []*http.Cookie
	StatusCode int
	Status     string
	ReqURL     string
	ReqAddr    string
	ReqHeader  http.Header
}

type Invoker func(ctx context.Context, method string, rawURL string, arg, reply interface{}, cc *Client, opt ...CallOption) (*RawResponse, error)

//nolint: funlen,gocognit
func invoke(
	ctx context.Context,
	method string,
	rawURL string,
	arg, reply interface{},
	cc *Client,
	opt ...CallOption,
) (*RawResponse, error) {
	log.F(ctx).Debug("invoke call.",
		log.String("method", method),
		log.String("rawURL", rawURL),
		log.Every("arg", arg))

	ci := &callInfo{}
	for _, o := range opt {
		o(ci)
	}

	reqURL := cc.addr + rawURL
	if ci.urlSetter != nil {
		var err error
		originURL := reqURL
		reqURL, err = ci.urlSetter()
		if err != nil {
			log.F(ctx).Errorf("http Do urlSetter err:%s", err.Error())
			return nil, err
		}
		log.F(ctx).Debugf("urlSetter change req url from %v to %v", originURL, reqURL)
	}
	if ci.query != nil {
		values := url.Values{}
		for k, v := range ci.query {
			value := ""
			switch v.(type) {
			case string:
				value = fmt.Sprintf("%s", v)
			case int:
				value = fmt.Sprintf("%d", v)
			default:
				value = fmt.Sprintf("%d", v)
			}
			values.Set(k, value)
		}
		if len(ci.query) > 0 {
			reqURL = fmt.Sprintf("%s?%s", reqURL, values.Encode())
		}
	}
	// refer to https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	var timeout time.Duration
	var cancel context.CancelFunc
	// 如果某个请求指定的超时时间, 则采用该超时时间
	if ci.timeout != nil && *ci.timeout >= 0 {
		timeout = *ci.timeout
	} else {
		// 如果客户设置了全局超时时间, 则采用该超时时间
		if cc.timeout != nil && *cc.timeout >= 0 {
			timeout = *cc.timeout
		}
	}

	if timeout > 0 {
		log.F(ctx).Debugf("request set timeout:%v", timeout)
		ctx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, method, reqURL, nil)
	if err != nil {
		log.F(ctx).Errorf("http.NewRequest error:%s", err.Error())
		return nil, err
	}

	for k, v := range ci.header {
		httpReq.Header.Add(k, v)
	}

	if arg != nil {
		var body []byte
		// 不对数据做任何处理
		// 如一些内含格式的字符串
		if str, ok := arg.(json.RawMessage); ok {
			body = str
		} else {
			body, err = json.Marshal(arg)
			if err != nil {
				log.F(ctx).Errorf("marshal data err:%s", err.Error())
				return nil, err
			}
		}

		if httpReq.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
			log.F(ctx).Debug("urlencoded body data")

			bf := bytes.NewBuffer(body)
			httpReq.Body = ioutil.NopCloser(bf)
			httpReq.ContentLength = int64(len(body))
		} else {
			log.F(ctx).Debug("json body data")

			httpReq.Header.Add("Content-Type", "application/json")
			httpReq.Body = ioutil.NopCloser(bytes.NewReader(body))
			httpReq.ContentLength = int64(len(body))
		}
	}

	conn, err := cc.getConn(ctx)
	if err != nil {
		log.F(ctx).Errorf("get client conn err:%s", err.Error())
		return nil, err
	}

	log.F(ctx).Debug("Before Do", log.Any("headers", httpReq.Header), log.String("requrl", httpReq.URL.String()))

	httpResp, err := conn.Do(httpReq)
	if err != nil {
		log.F(ctx).Errorf("http Do err:%s", err.Error())
		return nil, err
	}

	defer httpResp.Body.Close()

	rawResp := RawResponse{}
	rawResp.Header = httpResp.Header
	rawResp.Cookies = httpResp.Cookies()
	rawResp.StatusCode = httpResp.StatusCode
	rawResp.Status = httpResp.Status
	rawResp.ReqURL = reqURL
	rawResp.ReqAddr = cc.addr
	rawResp.ReqHeader = httpReq.Header

	log.F(ctx).Debug("After Do", log.Every("resp", rawResp))

	bodyData, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.F(ctx).Errorf("http read body err:%s", err.Error())
		return &rawResp, err
	}

	rawResp.Body = bodyData
	log.F(ctx).Debug("After Do", log.String("body", string(rawResp.Body)))

	// 由于无法覆盖服务器返回的状态码和返回请求体数据逻辑,因此提供选项允许调用者不在当前调用中进行数据解码
	// 调用者可以自定义拦截器来根据具体场景解构状态码并解析数据。 示例见NoSuccessStatusCodeInterceptor
	if !ci.responseNotParse && reply != nil {
		if err := json.Unmarshal(bodyData, reply); err != nil {
			log.F(ctx).Errorf("http decode  body err:%s", err.Error())
			return &rawResp, err
		}
		log.F(ctx).Debug("After Parse", log.Every("reply", reply))
	}
	return &rawResp, nil
}
