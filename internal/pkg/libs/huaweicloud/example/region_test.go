package example

import (
	"fmt"
	"testing"

	"github.com/wangweihong/gotoolbox/pkg/json"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	iam "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"
)

func TestAllRegions(t *testing.T) {
	// The AK and SK used for authentication are hard-coded or stored in plaintext, which has great security risks. It is recommended that the AK and SK be stored in ciphertext in configuration files or environment variables and decrypted during use to ensure security.
	// In this example, AK and SK are stored in environment variables for authentication. Before running this example, set environment variables CLOUD_SDK_AK and CLOUD_SDK_SK in the local environment
	//ak := os.Getenv("HUAWEI_AK")
	//sk := os.Getenv("HUAWEI_SK")
	ak := "9YKG5DPKOEDXEWOXF3GJ"
	sk := "0hQw8soVAnGR2P4jXDvOVeJyfRCgB7TWX6v8pJH0"
	auth := global.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		Build()

	// https://console.huaweicloud.com/apiexplorer/#/openapi/IAM/sdk?api=KeystoneListRegions
	client := iam.NewIamClient(
		iam.IamClientBuilder().
			//https://developer.huaweicloud.com/endpoint?IAM
			WithEndpoints([]string{"iam.myhuaweicloud.com"}).
			//WithRegion(region.ValueOf("cn-south-1")).
			WithCredential(auth).
			Build())

	request := &model.KeystoneListRegionsRequest{}
	response, err := client.KeystoneListRegions(request)
	if err == nil {
		json.PrintStructObject(response)
	} else {
		fmt.Println(err)
	}
}
