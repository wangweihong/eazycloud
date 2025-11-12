package aliyun

import (
	"context"
	"fmt"
	"strings"

	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/pathutil"
	"github.com/wangweihong/gotoolbox/pkg/typeutil"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Oss struct {
	c *oss.Client
}

func NewOss(ak, sk string, region string) (*Oss, error) {
	if region == "" {
		region = futureRegion
	}
	// 阿里云当前所有的对象存储端点都必须发往cn-beijing区域, 其他区域无效
	endpoint := fmt.Sprintf("oss-%s.aliyuncs.com", region)
	c, err := oss.New(endpoint, ak, sk)
	if err != nil {
		return nil, err
	}

	return &Oss{c: c}, nil
}

func (c *Oss) OssBucketCreate(
	ctx context.Context,
	bucketName string,
	acl string,
	storageClass string,
	DataRedundancyType string,
) error {
	var opts []oss.Option
	// 这里需要注意,acl参数是存放在请求头，sc和drt是存放在请求体
	// oss.ACL调用的是setHeader
	// oss.StorageClass/oss.RedundancyType调用的是addArg

	if acl != "" {
		opts = append(opts, oss.ACL(oss.ACLType(acl)))
	}

	if storageClass != "" {
		opts = append(opts, oss.StorageClass(oss.StorageClassType(storageClass)))
	}

	if DataRedundancyType != "" {
		opts = append(opts, oss.RedundancyType(oss.DataRedundancyType(DataRedundancyType)))
	}

	if err := c.c.CreateBucket(bucketName, opts...); err != nil {
		return err
	}

	return nil
}

func (c *Oss) OssBucketDelete(
	ctx context.Context,
	bucketName string,
) error {
	err := c.c.DeleteBucket(bucketName)
	return errors.WithStack(err)
}

// OssBucketList 返回全部存储桶列表
func (c *Oss) OssBucketList(
	ctx context.Context,
	Prefix *string,
	Delimiter *string,
) ([]oss.BucketProperties, error) {
	var objects []oss.BucketProperties

	var filterOpt []oss.Option
	if Prefix != nil {
		filterOpt = append(filterOpt, oss.Prefix(*Prefix))
	}
	if Delimiter != nil {
		filterOpt = append(filterOpt, oss.Delimiter(*Delimiter))
	}

	marker := oss.Marker("")
	for {
		var listOpt []oss.Option
		listOpt = append(listOpt, filterOpt...)
		listOpt = append(listOpt, oss.MaxKeys(1000))
		listOpt = append(listOpt, marker)
		ret, err := c.c.ListBuckets(listOpt...)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		marker = oss.Marker(ret.NextMarker)
		objects = append(objects, ret.Buckets...)
		if !ret.IsTruncated {
			break
		}
	}

	return objects, nil
}

func (c *Oss) OssBucketGetACL(ctx context.Context,
	bucketName string) (*oss.GetBucketACLResult, error) {
	ret, err := c.c.GetBucketACL(bucketName)
	return &ret, errors.WithStack(err)
}

func (c *Oss) OssBucketGetLifeCycle(ctx context.Context,
	bucketName string) (*oss.GetBucketLifecycleResult, error) {
	ret, err := c.c.GetBucketLifecycle(bucketName)
	return &ret, errors.WithStack(err)
}

func (c *Oss) OssBucketGetInfo(
	ctx context.Context,
	bucketName string,
) (*oss.GetBucketInfoResult, error) {
	ret, err := c.c.GetBucketInfo(bucketName)
	return &ret, errors.WithStack(err)
}

func (c *Oss) OssBucketGetPolicy(
	ctx context.Context,
	bucketName string,
) (*string, error) {
	ret, err := c.c.GetBucketPolicy(bucketName)
	return &ret, errors.WithStack(err)
}

// 下载对象到本地文件
func (c *Oss) OssBucketObjectDownload(
	ctx context.Context,
	bucketName string,
	objectName string,
	localFileName string,
) error {
	bucket, err := c.c.Bucket(bucketName)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := bucket.GetObjectToFile(objectName, localFileName); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

/*
	假如存储桶中有数据

d1/
f1
d1/d2/
d1/f2
d1/f3
d1/d2/ff1
d1/d2/d3/
d2/dd2/

 1. 如果prefix/delimiter都不传, 则直接返回存储桶中所有对象到contents中, CommonPrefix为空
 2. 如果delimiter设置了"/", 则

contents中只有文件f1
commonPrefix包含第一级目录d1,d2
 3. 如果delimiter设置了"/", prefix设置了"d1/",则相当于进入了d1目录

contents中的值为d1/f2, d1/f3(这里需要注意,返回的值会带上前缀)
commonPrefix包含d1的子目录d1/d2(这里需要注意,返回的值会带上前缀)
4.  如果只设置prefix为"d1", 则不会返回commonPrefix, 则是返回
d1/
d1/d2/
d1/f2
d1/f3
d1/d2/ff1
d1/d2/d3/
*/
func (c *Oss) OssBucketObjectList(
	ctx context.Context,
	bucketName string,
	Prefix *string,
	Delimiter *string,
) ( /*contents*/ []oss.ObjectProperties /*CommonPrefix*/, []string, error) {

	bucket, err := c.c.Bucket(bucketName)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	var objects []oss.ObjectProperties
	var commonPrefixes []string
	var filterOpt []oss.Option
	if Prefix != nil {
		filterOpt = append(filterOpt, oss.Prefix(*Prefix))
	}
	if Delimiter != nil {
		filterOpt = append(filterOpt, oss.Delimiter(*Delimiter))
	}

	marker := oss.Marker("")
	for {
		var listOpt []oss.Option
		listOpt = append(listOpt, filterOpt...)
		listOpt = append(listOpt, oss.MaxKeys(1000))
		listOpt = append(listOpt, marker)
		ret, err := bucket.ListObjects(listOpt...)
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}
		marker = oss.Marker(ret.NextMarker)
		objects = append(objects, ret.Objects...)

		commonPrefixes = append(commonPrefixes, ret.CommonPrefixes...)
		if !ret.IsTruncated {
			break
		}
	}

	return objects, commonPrefixes, nil
}

// OssBucketObjectDelete 删除文件对象或者目录对象下的所有子对象
func (c *Oss) OssBucketObjectDelete(
	ctx context.Context,
	bucketName string,
	objectName string,
) error {
	bucket, err := c.c.Bucket(bucketName)
	if err != nil {
		return err
	}
	// 由于阿里云api不支持删除某个文件夹下的所有对象功能(接口调用成功,但仍然存在).
	// 因此先要获取该前缀下所有对象的key, 然后通过目录最后深度排序，因此删除这些对象。

	// 只传prefix,不传Delimiter,则将会带有该前缀的所有key
	objectRet, _, err := c.OssBucketObjectList(
		ctx,
		bucketName,
		typeutil.String(objectName),
		nil)
	if err != nil {
		return errors.WithStack(err)
	}

	var fps pathutil.DirLastDepthPaths
	for _, v := range objectRet {
		fp := pathutil.DirLastDepthPath{
			Value: v.Key,
			IsDir: strings.HasSuffix(v.Key, "/"),
		}
		fps = append(fps, fp)
	}
	// 按文件优先深度排序
	fps.Sort()
	// 依次删除。
	for _, v := range fps {
		if err := bucket.DeleteObject(v.Value); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// 创建文件夹或者上传文件到指定目录
func (c *Oss) OssBucketObjectCreate(
	ctx context.Context,
	bucketName string,
	key string,
	localfile string,
	acl string,
	StorageClass string,
) error {

	bucket, err := c.c.Bucket(bucketName)
	if err != nil {
		return errors.WithStack(err)
	}

	var opts []oss.Option
	if acl != "" {
		opts = append(opts, oss.ACL(oss.ACLType(acl)))
	}

	if StorageClass != "" {
		// 一定要注意, 请求参数是放在body还是header
		opts = append(opts, oss.ObjectStorageClass(oss.StorageClassType(StorageClass)))
	}

	if localfile != "" {
		// 上传文件
		// 这里需要注意，如果上传文件，key为aaa/a.txt, 则aaa/并不存在, 则上传文件后，只存在aaa/a.txt这个对象
		// 不存在aaa/这个对象
		if err := bucket.PutObjectFromFile(key, localfile, opts...); err != nil {
			return errors.WithStack(err)
		}
	} else {
		// 创建文件夹
		if err := bucket.PutObject(key, nil, opts...); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (c *Oss) OssBucketObjectSetMeta(
	ctx context.Context,
	bucketName string,
	objectName string,
	acl *string,
	StorageClass *string,
) error {
	bucket, err := c.c.Bucket(bucketName)
	if err != nil {
		return errors.WithStack(err)
	}
	var opts []oss.Option
	// 这里需要注意,acl和sc参数是存放在请求头
	if acl != nil {
		opts = append(opts, oss.ACL(oss.ACLType(*acl)))
	}

	if StorageClass != nil {
		// 一定要注意, 请求参数是放在body还是header
		opts = append(opts, oss.ObjectStorageClass(oss.StorageClassType(*StorageClass)))
	}
	if err := bucket.SetObjectMeta(objectName, opts...); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
