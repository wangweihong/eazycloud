package httpcli

import (
	"net/http"
	"net/url"
	"time"
)

type Option func(*Client)

// WithTimeout 设置所有连接超时操作.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = &timeout
	}
}

// WithReport 是否打印返参.
func WithReport() Option {
	return func(c *Client) {
		c.report = true
	}
}

// WithInsecure 是否跳过服务端证书检测.
func WithInsecure() Option {
	return func(c *Client) {
		c.tlsEnabled = true
		c.skipTlsVerified = true
	}
}

// WithServerCA 设置服务端CA证书数据.
func WithServerCA(serverCAData string) Option {
	return func(c *Client) {
		c.tlsEnabled = true
		c.serverCA = serverCAData
	}
}

// WithMTLS 是否开启双向认证.
func WithMTLS(serverCAData string, clientCertData string, clientKeyData string) Option {
	return func(c *Client) {
		c.mtlsEnabled = true
		c.tlsEnabled = true
		c.clientCertData = clientCertData
		c.clientKeyData = clientKeyData
		c.serverCA = serverCAData
	}
}

// WithIntercepts 插入拦截器
// 注意顺序:序号0的拦截器第一个执行调用前的处理, 最后一个执行调用后的处理.
func WithIntercepts(inters ...Interceptor) Option {
	return func(c *Client) {
		c.chainInterceptors = inters
	}
}

// WithCallOption 通用请求选项.
// 用于对单独请求选项设置。
func WithCallOption(copt ...CallOption) Option {
	return func(c *Client) {
		c.callOpts = copt
	}
}

// WithTransport 通用请求选项.
func WithTransport(tp *http.Transport) Option {
	return func(c *Client) {
		c.transport = tp
	}
}

// WithProxy 请求代理选项.
// WithProxy(http.ProxyFromEnvironment)
func WithProxy(proxy func(*http.Request) (*url.URL, error)) Option {
	return func(c *Client) {
		c.proxy = proxy
	}
}
