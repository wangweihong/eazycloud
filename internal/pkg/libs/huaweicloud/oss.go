package huaweicloud

import (
	"context"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

type Oss struct {
	c   *obs.ObsClient
	cfg *AccessConfig
}

func NewOss(cfg *AccessConfig) (*Oss, error) {
	c, err := cfg.HcOssClient(cfg.Region)
	if err != nil {
		return nil, err
	}

	return &Oss{c: c, cfg: cfg}, nil
}

// OssBucketList 查看存储桶列表
// https://support.huaweicloud.com/sdk-go-devg-obs/obs_33_0403.html
func (p *Oss) OssBucketList(
	ctx context.Context,
	req *obs.ListBucketsInput,
) (*obs.ListBucketsOutput, error) {
	resp, err := p.c.ListBuckets(req)
	return resp, errors.WithStack(err)
}

//// OssBucketCheckExist 确认存储桶是否存在
//// https://support.huaweicloud.com/sdk-go-devg-obs/obs_33_0404.html
//func (p *Oss) OssBucketCheckExist(
//	ctx context.Context,
//	bucketName string,
//) (bool, error) {
//	_, err := p.c.HeadBucket(bucketName)
//	if err != nil {
//		if strings.Contains(err.Error(), "Status=404 Not Found, Code=NoSuchBucket") {
//			return false, nil
//		}
//		return false, err
//	}
//	return true, nil
//}

// OssBucketCreate 创建存储桶
func (p *Oss) OssBucketCreate(
	ctx context.Context,
	req *obs.CreateBucketInput,
) (*obs.BaseModel, error) {
	if req.Location == "" {
		req.Location = p.cfg.Region
	}
	resp, err := p.c.CreateBucket(req)
	return resp, errors.WithStack(err)
}

// OssBucketDelete 删除存储桶
func (p *Oss) OssBucketDelete(
	ctx context.Context,
	bucketName string,
) (*obs.BaseModel, error) {
	resp, err := p.c.DeleteBucket(bucketName)
	return resp, errors.WithStack(err)
}

// OssBucketMetadata 获取桶元数据
// https://support.huaweicloud.com/sdk-go-devg-obs/obs_33_0409.html
func (p *Oss) OssBucketGetMetadata(
	ctx context.Context,
	req *obs.GetBucketMetadataInput,
) (*obs.GetBucketMetadataOutput, error) {
	resp, err := p.c.GetBucketMetadata(req)
	return resp, errors.WithStack(err)
}

// OssBucketStorageInfoGet 获取存储桶存量信息
// https://support.huaweicloud.com/sdk-go-devg-obs/obs_33_0411.html
func (p *Oss) OssBucketGetStorageInfo(
	ctx context.Context,
	bucketName string,
) (*obs.GetBucketStorageInfoOutput, error) {
	resp, err := p.c.GetBucketStorageInfo(bucketName)
	return resp, errors.WithStack(err)
}

// OssBucketGetStoragePolicy 获取存储桶类型
// https://support.huaweicloud.com/sdk-go-devg-obs/obs_33_0415.html
func (p *Oss) OssBucketGetStoragePolicy(
	ctx context.Context,
	bucketName string,
) (*obs.GetBucketStoragePolicyOutput, error) {
	resp, err := p.c.GetBucketStoragePolicy(bucketName)
	return resp, errors.WithStack(err)
}

// OssBucketSetStoragePolicy 设置存储桶类型
// https://support.huaweicloud.com/sdk-go-devg-obs/obs_33_0415.html
func (p *Oss) OssBucketSetStoragePolicy(
	ctx context.Context,
	req *obs.SetBucketStoragePolicyInput,
) (*obs.BaseModel, error) {
	resp, err := p.c.SetBucketStoragePolicy(req)
	return resp, errors.WithStack(err)
}

// OssBucketObjectList 桶对象列表
func (p *Oss) OssBucketObjectList(
	ctx context.Context,
	req *obs.ListObjectsInput,
) (*obs.ListObjectsOutput, error) {
	resp, err := p.c.ListObjects(req)
	return resp, errors.WithStack(err)
}

// OssBucketDirectoryCreate 存储桶创建文件夹
// https://support.huaweicloud.com/sdk-go-devg-obs/obs_23_0510.html
func (p *Oss) OssBucketDirectoryCreate(
	ctx context.Context,
	bucketName string,
	key string,
) (*obs.PutObjectOutput, error) {
	if !strings.HasSuffix(key, "/") {
		key = key + "/"
	}
	req := &obs.PutObjectInput{}
	req.Bucket = bucketName
	req.Key = key
	resp, err := p.c.PutObject(req)
	return resp, errors.WithStack(err)
}

// OssBucketDirectoryCreate 存储桶删除文件夹
func (p *Oss) OssBucketDirectoryDelete(
	ctx context.Context,
	bucketName string,
	key string,
) (*obs.PutObjectOutput, error) {
	if !strings.HasSuffix(key, "/") {
		key = key + "/"
	}
	req := &obs.PutObjectInput{}
	req.Bucket = bucketName
	req.Key = key
	resp, err := p.c.PutObject(req)
	return resp, errors.WithStack(err)
}

// OssBucketObjectUpload 存储桶上传文件(流式对象)
// https://support.huaweicloud.com/sdk-go-devg-obs/obs_23_0402.html
func (p *Oss) OssBucketObjectUpload(
	ctx context.Context,
	req *obs.PutObjectInput,
) (*obs.PutObjectOutput, error) {
	resp, err := p.c.PutObject(req)
	return resp, errors.WithStack(err)
}

// OssBucketObjectDownload 存储桶上传文件(流式对象)
// https://support.huaweicloud.com/sdk-go-devg-obs/obs_23_0502.html
func (p *Oss) OssBucketObjectDownload(
	ctx context.Context,
	req *obs.GetObjectInput,
) (*obs.GetObjectOutput, error) {
	resp, err := p.c.GetObject(req)
	return resp, errors.WithStack(err)
}

// OssBucketObjectDelete  删除存储桶文件(流式对象)
// https://support.huaweicloud.com/sdk-go-devg-obs/obs_23_0502.html
func (p *Oss) OssBucketObjectDelete(
	ctx context.Context,
	req *obs.DeleteObjectInput,
) (*obs.DeleteObjectOutput, error) {
	resp, err := p.c.DeleteObject(req)
	return resp, errors.WithStack(err)
}
