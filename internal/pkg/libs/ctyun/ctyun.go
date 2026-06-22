package ctyun

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/wangweihong/eazycloud/apis/ictyun"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"
	"github.com/wangweihong/gotoolbox/pkg/httpcli/httphandler"

	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/httpcli/httpconfig"
	"github.com/wangweihong/gotoolbox/pkg/log"
	"github.com/wangweihong/gotoolbox/pkg/skipper"

	"github.com/wangweihong/gotoolbox/pkg/httpcli/interceptorcli"

	"io"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/go-basic/uuid"
)

const (
	SERVICE_ECS       = "ecs"
	SERVICE_VPC       = "vpc"
	SERVICE_IMAGE     = "image"
	SERVICE_ACCT      = "acct"
	SERVICE_EBS       = "ebs"
	SERVICE_EBSBACKUP = "ebsbackup"
	SERVICE_MONITOR   = "monitor"
	SERVICE_OSS       = "oss"
)

type Client struct {
	cc        *httpcli.Client
	accessKey string
	secretKey string
	err       error
}

func NewClient(ak, sk string, opt ...httpcli.Option) *Client {
	c := &Client{
		accessKey: ak,
		secretKey: sk,
	}

	if ak == "" || sk == "" {
		c.err = errors.Errorf("accessKey or secretKey empty")
		return c
	}

	var err error
	var cc *httpcli.Client
	if opt != nil {
		cc, err = httpcli.NewClient(nil, opt...)
	} else {
		// 建立长连接?
		HTTPTransport := &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second, // 连接超时时间
				KeepAlive: 60 * time.Second, // 保持长连接的时间
			}).DialContext, // 设置连接的参数
			MaxIdleConns:          500,              // 最大空闲连接
			IdleConnTimeout:       60 * time.Second, // 空闲连接的超时时间
			ExpectContinueTimeout: 30 * time.Second, // 等待服务第一个响应的超时时间
			MaxIdleConnsPerHost:   100,              // 每个host保持的空闲连接数
		}
		cfg := httpconfig.DefaultHttpConfig()
		cfg.HttpHandler = httphandler.NewHttpHandler().
			AddRequestHandler(genSignRequestMethod(ak, sk))

		cc, err = httpcli.NewClient(
			cfg,
			httpcli.WithTransport(HTTPTransport),
			httpcli.WithTimeout(30*time.Second),
			httpcli.WithIntercepts(
				// 注意顺序, 队列也靠后的越早执行调用后逻辑
				interceptorcli.TraceInterceptor("trace"),
				ErrorCodeInterceptor("ErrorCodeInterceptor"),
				interceptorcli.StatusCodeInterceptor("NoSuccessStatusCodeInterceptor"),
				interceptorcli.LoggingInterceptor("Logging"),
			),
		)
	}
	c.err = err
	c.cc = cc
	return c
}

// 错误码拦截.
func ErrorCodeInterceptor(name string, skipperFunc ...skipper.SkipperFunc) httpcli.Interceptor {
	return httpcli.NewInterceptor("error code", func(ctx context.Context, req *httpcli.HttpRequest, arg, reply interface{}, cc *httpcli.Client, invoker httpcli.Invoker, opts ...httpcli.CallOption) (*httpcli.HttpResponse, error) {
		if skipper.Skip(req.GetPath(), skipperFunc...) {
			log.F(ctx).Debugf("skip interceptor %s for rawrurl %s", name, req.GetPath())
			return invoker(ctx, req, arg, reply, cc, opts...)
		}

		rawResp, err := invoker(ctx, req, arg, reply, cc, opts...)
		if err != nil {
			return rawResp, errors.WithStack(err)
		}
		var ctRet ictyun.ResponseResult
		if err := rawResp.Decode(&ctRet); err != nil {
			return rawResp, errors.WithStack(err)
		}

		switch statusCode := ctRet.StatusCode.(type) {
		case string:
			if statusCode != "800" {
				return rawResp, errors.Errorf("url:%v,response result:%v", req.GetPath(), ctRet.String())
			}
		case float64:
			if statusCode != 800 {
				return rawResp, errors.Errorf("url:%v,resp:%v", rawResp.GetBody(), ctRet.String())
			}
		default:
			return rawResp, errors.Errorf("unknown status code type:%v", reflect.TypeOf(ctRet.StatusCode).String())
		}

		if reply != nil {
			decoder := json.NewDecoder(bytes.NewReader([]byte(rawResp.GetBody())))
			decoder.UseNumber()
			if err := decoder.Decode(reply); err != nil {
				return rawResp, errors.WithStack(err)
			}
		}

		return rawResp, nil
	})
}

func (c *Client) serviceEndpoint(service string) string {
	switch service {
	case SERVICE_ECS, SERVICE_VPC, SERVICE_IMAGE:
		return fmt.Sprintf("https://ct%s-global.ctapi.ctyun.cn", service)
	case SERVICE_ACCT, SERVICE_EBS, SERVICE_MONITOR, SERVICE_EBSBACKUP:
		return fmt.Sprintf("https://%s-global.ctapi.ctyun.cn", service)
	case SERVICE_OSS:
		return fmt.Sprintf("https://zos-global.ctapi.ctyun.cn")
	}
	return "https://global.ctapi.ctyun.cn"
}

func genSignRequestMethod(ak, sk string) func(req *http.Request) error {
	return func(req *http.Request) error {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("ctyun-eop-request-id", uuid.New())
		sh, _ := time.LoadLocation("Asia/Shanghai")
		req.Header.Set("eop-date", time.Now().In(sh).Format("20060102T150405Z"))

		signature, err := eopSign(ak, sk, req)
		if err != nil {
			return errors.WithStack(err)
		}
		req.Header.Set("Eop-Authorization", signature)
		return nil
	}
}

func eopSign(accessKey, secretKey string, req *http.Request) (string, error) {
	eopDate := req.Header.Get("eop-date")
	requestId := req.Header.Get("ctyun-eop-request-id")
	headerStr := fmt.Sprintf("ctyun-eop-request-id:%s\neop-date:%s\n",
		requestId,
		eopDate,
	)

	keys := make([]string, 0, len(req.URL.Query()))
	for key := range req.URL.Query() {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(req.URL.Query().Get(key))))
	}

	var body []byte
	if req.Method == "POST" {
		var err error
		body, err = io.ReadAll(req.Body)
		if err != nil {
			return "", errors.WithStack(err)
		}
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	var hmacSha256 = func(secret, data []byte) []byte {
		hasher := hmac.New(sha256.New, []byte(secret))
		hasher.Write(data)
		return hasher.Sum(nil)
	}

	hash := sha256.New()
	hash.Write(body)
	bodyHash := hex.EncodeToString(hash.Sum(nil))

	signStr := fmt.Sprintf("%s\n%s\n%s", headerStr, strings.Join(parts, "&"), bodyHash)

	kTime := hmacSha256([]byte(secretKey), []byte(eopDate))
	kAk := hmacSha256(kTime, []byte(accessKey))
	t := strings.Split(eopDate, "T")[0]
	kDate := hmacSha256(kAk, []byte(t))
	signBase64 := base64.StdEncoding.EncodeToString(hmacSha256(kDate, []byte(signStr)))

	return fmt.Sprintf("%s Headers=ctyun-eop-request-id;eop-date Signature=%s", accessKey, signBase64), nil
}

func (c *Client) EbsBackups() *EbsBackup {
	return &EbsBackup{
		c:           c,
		serviceType: SERVICE_EBSBACKUP,
	}
}

func (c *Client) Billings() *Billing {
	return &Billing{
		c:           c,
		serviceType: SERVICE_ACCT,
	}
}

func (c *Client) Ebss() *Ebs {
	return &Ebs{
		c:           c,
		serviceType: SERVICE_EBS,
	}
}

func (c *Client) Ecss() *Ecs {
	return &Ecs{
		c:           c,
		serviceType: SERVICE_ECS,
	}
}

func (c *Client) Vpcs() *Vpc {
	return &Vpc{
		c:           c,
		serviceType: SERVICE_VPC,
	}
}

func (c *Client) Images() *Image {
	return &Image{
		c:           c,
		serviceType: SERVICE_IMAGE,
	}
}

func (c *Client) Monitors() *Monitor {
	return &Monitor{
		c:           c,
		serviceType: SERVICE_MONITOR,
	}
}

func (c *Client) Osss(service string) *Oss {
	return &Oss{
		c:           c,
		serviceType: SERVICE_OSS,
	}
}
