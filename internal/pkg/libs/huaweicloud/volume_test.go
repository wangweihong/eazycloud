package huaweicloud_test

import (
	"context"
	"testing"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/evs/v2/model"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/wangweihong/eazycloud/internal/pkg/libs/huaweicloud"
	"github.com/wangweihong/gotoolbox/pkg/json"
)

func Test11(t *testing.T) {
	Convey("", t, func() {
		data := `
	{
		"cluster_uuid": "{{global.cluster}}",
		"region": "{{global.region}}",
		"body": {
			"volume": {
			    "name": "aaaatest2",
				"availability_zone": "cn-south-1g",
				"size": 10,
				"volume_type": "SSD"
			},
			"bssParam": {
				"chargingMode": "postPaid"
			}
		}
	}`

		type CreateVolumeRequest struct {
			model.CreateVolumeRequest
			ClusterUUID string `json:"cluster_uuid"`
			Region      string `json:"region"`
		}
		var arg CreateVolumeRequest
		//arg := model.CreateVolumeRequest{}
		err := json.Unmarshal([]byte(data), &arg)
		if err != nil {
			t.Fatal(err)
		}

		//requestDef := evs.GenReqDefForCreateVolume()
		//_, err = buildRequest(&arg, requestDef)
		//json.PrintStructObject(arg)
		//fmt.Println(arg.Body.BssParam.ChargingMode)
		//fmt.Println(arg)
		resp, err := huaweicloud.NewVolume(ac).VolumeCreate(context.Background(), &arg.CreateVolumeRequest)
		So(err, ShouldBeNil)
		json.PrintStructObject(resp)
	})

}
