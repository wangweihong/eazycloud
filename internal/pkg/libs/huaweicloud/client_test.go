package huaweicloud_test

import (
	"os"

	"github.com/wangweihong/eazycloud/internal/pkg/libs/huaweicloud"
)

var (
	ak = os.Getenv("HUAWEI_AK")
	sk = os.Getenv("HUAWEI_SK")
	//region = "cn-south-1"
	region = "cn-east-3"

	ac, err = huaweicloud.NewAccessConfig(ak, sk, region)
)
