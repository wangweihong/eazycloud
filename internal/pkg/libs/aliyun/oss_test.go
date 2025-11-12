package aliyun_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func TestAAAAA(t *testing.T) {
	// 从环境变量中获取访问凭证。运行本代码示例之前，请确保已设置环境变量OSS_ACCESS_KEY_ID和OSS_ACCESS_KEY_SECRET。
	//provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	os.Exit(-1)
	//}

	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	client, err := oss.New("https://oss-cn-hangzhou.aliyuncs.com", os.Getenv("OSS_ACCESS_KEY_ID"), os.Getenv("OSS_ACCESS_KEY_SECRET"))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	// 列举当前账号所有地域下的存储空间，限定此次列举存储空间的最大个数为500。MaxKeys默认值为100，最大值为1000。
	lsRes, err := client.ListBuckets(oss.MaxKeys(500))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 打印存储空间列表。
	fmt.Println("My buckets max num:", lsRes.Buckets)
	for _, bucket := range lsRes.Buckets {
		fmt.Println("Bucket with maxKeys: ", bucket.Name)
	}
}
