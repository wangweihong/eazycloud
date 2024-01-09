package httpcli_test

import (
	"github.com/wangweihong/eazycloud/pkg/httpcli"
	"net/http"
	"testing"
)

func TestWithProxy(t *testing.T) {
	// 从环境变量中读取代理
	httpcli.WithProxy(http.ProxyFromEnvironment)
}
