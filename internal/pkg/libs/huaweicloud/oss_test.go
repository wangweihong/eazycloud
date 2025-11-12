package huaweicloud_test

import (
	"bytes"
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"

	"github.com/wangweihong/eazycloud/internal/pkg/libs/huaweicloud"
	"github.com/wangweihong/gotoolbox/pkg/json"
)

var (
	oss, _ = huaweicloud.NewOss(ac)
)

// const bucketName = "asadassadadaadada"
const bucketName = "asadassadadaadada"

func TestOss_OssBucketList(t *testing.T) {
	Convey("bucket info", t, func() {
		Convey("bucket list", func() {
			resp, err := oss.OssBucketList(context.Background(), &obs.ListBucketsInput{})
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, 200)
			json.PrintStructObject(resp)
		})

		//Convey("bucket check exists", func() {
		//	resp, err := oss.OssBucketCheckExist(context.Background(), bucketName)
		//	So(err, ShouldBeNil)
		//	json.PrintStructObject(resp)
		//})
	})
}

func TestOss_OssBucketStorageInfo(t *testing.T) {
	Convey("bucket storage info", t, func() {
		Convey("1", func() {

			resp, err := oss.OssBucketGetStorageInfo(context.Background(), bucketName)
			So(err, ShouldBeNil)
			json.PrintStructObject(resp)
		})
	})
}

func TestOss_OssBucketCreate(t *testing.T) {
	Convey("bucket create", t, func() {
		Convey("1", func() {
			newB := "vsf111123ddd"
			req := &obs.CreateBucketInput{}
			req.Bucket = newB
			req.ACL = obs.AclPublicReadWrite
			req.StorageClass = obs.StorageClassWarm
			//必传参数
			//req.Location = region
			//req.AvailableZone = "3az"

			resp, err := oss.OssBucketCreate(context.Background(), req)
			So(err, ShouldBeNil)
			json.PrintStructObject(resp)
		})
	})
}

func TestOss_OssBucketDelete(t *testing.T) {
	Convey("bucket delete", t, func() {
		Convey("1", func() {
			resp, err := oss.OssBucketDelete(context.Background(), "a1231313aaa")
			So(err, ShouldBeNil)
			json.PrintStructObject(resp)
		})
	})
}

func TestOss_OssBucketObjectUpload(t *testing.T) {
	Convey("bucket object upload", t, func() {
		Convey("1", func() {
			b := bytes.NewBuffer([]byte("aaavv"))
			req := &obs.PutObjectInput{}
			req.Bucket = bucketName
			req.Key = "test2"
			req.StorageClass = "STANDARD"
			req.Body = b
			req.ACL = "private"
			resp, err := oss.OssBucketObjectUpload(context.Background(), req)
			So(err, ShouldBeNil)
			json.PrintStructObject(resp)
		})
	})
}

func TestOss_OssBucketObjectDownload(t *testing.T) {
	Convey("bucket object upload", t, func() {
		Convey("1", func() {
			req := &obs.GetObjectInput{}
			req.Bucket = bucketName
			req.Key = "cccc.txt"
			resp, err := oss.OssBucketObjectDownload(context.Background(), req)
			So(err, ShouldBeNil)
			json.PrintStructObject(resp)
		})
	})
}
