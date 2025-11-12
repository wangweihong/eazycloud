package ctyun_test

import (
	"os"

	"github.com/wangweihong/eazycloud/internal/pkg/libs/ctyun"
)

var (
	endpoint = "global.ctapi.ctyun.cn"
	ak       = os.Getenv("AK")
	sk       = os.Getenv("SK")
	region   = "aaf589124d5d11eaa04d0242ac110002" //贵州3
	tc       = ctyun.NewClient(ak, sk)
	rgName   = "贵州3"
	flavor   = "s6.small.1"
	// 随便一个值
	clientToken = "1adsfdsfdfda"
)
