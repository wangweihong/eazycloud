package ctyun_test

//
//import (
//	"context"
//	"io"
//	"net/http"
//	"os"
//	"path/filepath"
//	"strings"
//	"testing"
//
//	"github.com/aws/aws-sdk-go/aws"
//	"github.com/aws/aws-sdk-go/aws/credentials"
//	"github.com/aws/aws-sdk-go/aws/session"
//	"github.com/aws/aws-sdk-go/service/s3"
//	. "github.com/smartystreets/goconvey/convey"
//	"github.com/wangweihong/eazycloud/internal/pkg/multicloud/ctyun"
//	"github.com/wangweihong/eazycloud/internal/pkg/multicloud/ctyun/ictyun"
//	"github.com/wangweihong/gotoolbox/pkg/json"
//	"github.com/wangweihong/gotoolbox/sets"
//)
//
//const (
//	bucket = "bucket-6369"
//	object = "dir1/key1" //如果指定某个目录下, 则key为相对路径
//)
//
//func TestOss_OssBucketList(t *testing.T) {
//	Convey("TestOss_OssBucketList", t, func() {
//		c, err := ctyun.NewClient(endpoint, ak, sk)
//		So(err, ShouldBeNil)
//		resp, err := ctyun.NewOss(c).OssBucketList(context.Background(), &ictyun.OssBucketListRequest{
//			RegionID: region,
//		})
//		So(err, ShouldBeNil)
//		json.PrintStructObject(resp)
//	})
//}
//
//func TestOss_OssBucketGetACL(t *testing.T) {
//	Convey("TestOss_OssBucketGetACL", t, func() {
//		c, err := ctyun.NewClient(endpoint, ak, sk)
//		So(err, ShouldBeNil)
//		resp, err := ctyun.NewOss(c).OssBucketGetACL(context.Background(), &ictyun.OssBucketGetACLRequest{
//			RegionID: region,
//			Bucket:   bucket,
//		})
//		So(err, ShouldBeNil)
//		//json.PrintStructObject(resp)
//
//		acl := "private"
//		for _, v := range resp.ReturnObj.Grants {
//			if strings.Contains(v.Grantee.URI, "AllUsers") {
//				if v.Permission == "READ" {
//					acl = "public-read"
//				}
//				if sets.NewString("WRITE", "WRITE_ACP", "FULL_CONTROL").Has(v.Permission) {
//					acl = "public-read-write"
//				}
//			}
//		}
//
//		So(acl, ShouldNotBeEmpty)
//	})
//}
//
//func TestOss_OssBucketObjectDownload(t *testing.T) {
//	Convey("下载存储桶对象", t, func() {
//		c, err := ctyun.NewClient(endpoint, ak, sk)
//		So(err, ShouldBeNil)
//		//1. 提取下载链接
//		resp, err := ctyun.NewOss(c).OssObjectDownloadLink(context.Background(), &ictyun.OssObjectDownloadLinkRequest{
//			RegionID: region,
//			Bucket:   bucket,
//			Key:      object,
//		})
//		So(err, ShouldBeNil)
//
//		//2. 通过链接下载文件
//		httpResp, err := http.Get(resp.ReturnObj)
//		So(err, ShouldBeNil)
//		So(httpResp.StatusCode, ShouldEqual, http.StatusOK)
//		defer httpResp.Body.Close()
//
//		file, err := os.Create(filepath.Base(object))
//		So(err, ShouldBeNil)
//		defer file.Close()
//
//		_, err = io.Copy(file, httpResp.Body)
//		So(err, ShouldBeNil)
//	})
//}
//
//func TestOss_OssBucketObjectUpload(t *testing.T) {
//	Convey("上传存储桶对象", t, func() {
//		// 生成的存储对象的acl
//		var acl string = "default"
//		fp, err := os.Open("example1")
//		So(err, ShouldBeNil)
//		defer fp.Close()
//
//		c, err := ctyun.NewClient(endpoint, ak, sk)
//		So(err, ShouldBeNil)
//
//		// 获取aws访问端点
//		endpointRet, err := ctyun.NewOss(c).OssEndpointList(context.Background(), &ictyun.OssEndpointRequest{
//			RegionID: region,
//		})
//		So(err, ShouldBeNil)
//		So(len(endpointRet.ReturnObj.InternetEndpoint), ShouldNotEqual, 0)
//
//		// 获取aws访问密钥
//		keyRet, err := ctyun.NewOss(c).OssKeyList(context.Background(), &ictyun.OssKeyRequest{
//			RegionID: region,
//		})
//		So(err, ShouldBeNil)
//		So(len(keyRet.ReturnObj), ShouldNotEqual, 0)
//
//		//继承存储桶的ACL
//		if acl == "default" {
//			aclRet, err := ctyun.NewOss(c).OssBucketGetACL(context.Background(), &ictyun.OssBucketGetACLRequest{
//				RegionID: region,
//				Bucket:   bucket,
//			})
//			So(err, ShouldBeNil)
//			acl = "private"
//			for _, v := range aclRet.ReturnObj.Grants {
//				if strings.Contains(v.Grantee.URI, "AllUsers") {
//					if v.Permission == "READ" {
//						acl = "public-read"
//					}
//					if sets.NewString("WRITE", "WRITE_ACP", "FULL_CONTROL").Has(v.Permission) {
//						acl = "public-read-write"
//					}
//				}
//			}
//		}
//
//		sess, err := session.NewSession(&aws.Config{
//			Credentials:      credentials.NewStaticCredentials(keyRet.ReturnObj[0].AccessKey, keyRet.ReturnObj[0].SecretKey, ""),
//			Endpoint:         aws.String(endpointRet.ReturnObj.InternetEndpoint[0]),
//			Region:           aws.String("default"),
//			DisableSSL:       aws.Bool(true),
//			S3ForcePathStyle: aws.Bool(true),
//		})
//		So(err, ShouldBeNil)
//		_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
//			Bucket: aws.String(bucket),
//			Key:    aws.String(object), //如果要上传到指定的目录, 则加上目录名+"/"
//			ACL:    aws.String(acl),
//			Body:   fp,
//		})
//		So(err, ShouldBeNil)
//	})
//}
//
//func TestOss_OssBucketObjectDelete(t *testing.T) {
//	Convey("TestOss_OssBucketObjectDelete", t, func() {
//		c, err := ctyun.NewClient(endpoint, ak, sk)
//		So(err, ShouldBeNil)
//
//		Convey("删除根目录下文件", func() {
//			_, err = ctyun.NewOss(c).OssBucketDelete(context.Background(), &ictyun.OssBucketDeleteRequest{
//				RegionID: region,
//				Bucket:   "example1",
//			})
//			So(err, ShouldBeNil)
//		})
//
//		Convey("删除指定目录下文件", func() {
//			_, err = ctyun.NewOss(c).OssBucketDelete(context.Background(), &ictyun.OssBucketDeleteRequest{
//				RegionID: region,
//				Bucket:   "dir/example1",
//			})
//			So(err, ShouldBeNil)
//		})
//
//	})
//}
