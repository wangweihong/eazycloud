package example

import (
	"context"
	"sort"
	"strings"
	"sync"

	"github.com/wangweihong/eazycloud/internal/pkg/libs/aliyun"
	"github.com/wangweihong/gotoolbox/pkg/paging"

	"github.com/wangweihong/gotoolbox/pkg/log"

	"github.com/wangweihong/gotoolbox/pkg/waitgroup"

	"github.com/wangweihong/gotoolbox/pkg/typeutil"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type PagingParam struct {
	PageNumber *int32
	PageSize   *int32
}

type BucketListRequest struct {
	PagingParam
	Prefix    *string //prefix
	Delimiter *string
	Fuzzy     string
}

type BucketListUnit struct {
	Property  oss.BucketProperties
	ACL       *oss.GetBucketACLResult
	Lifecycle *oss.GetBucketLifecycleResult
	Info      *oss.GetBucketInfoResult
	ObjectNum int
}

type BucketListResponse struct {
	Info  []BucketListUnit //存储桶列表
	Total int              //总数
}

// BucketList 查询存储桶列表, 支持模糊搜索，分页查询
func BucketList(ctx context.Context, c *aliyun.Oss, req *BucketListRequest) (*BucketListResponse, error) {
	resp := &BucketListResponse{}
	ret, err := c.OssBucketList(ctx, req.Prefix, req.Delimiter)
	if err != nil {
		return nil, err
	}

	filterRet := make([]oss.BucketProperties, 0, len(ret))
	//模糊搜索
	if req.Fuzzy != "" {
		for _, v := range ret {
			if !strings.Contains(v.Name, req.Fuzzy) {
				continue
			}
			filterRet = append(filterRet, v)
		}
	} else {
		filterRet = ret
	}

	resp.Total = len(filterRet)
	sort.SliceStable(filterRet, func(i, j int) bool {
		return filterRet[i].Name < filterRet[j].Name
	})

	s, e := paging.Index(resp.Total, typeutil.Int32ValueToInt(req.PageNumber)-1, typeutil.Int32ValueToInt(req.PageSize))
	buckets := filterRet[s:e]
	lock := sync.Mutex{}
	waitgroup.RunConcurrently(ctx, buckets, func(ctx context.Context, v oss.BucketProperties) waitgroup.Result {
		acl, err := c.OssBucketGetACL(ctx, v.Name)
		if err != nil {
			log.L(ctx).Debugf("OssBucketGetACL err:%v", err)
		}
		lifecycle, err := c.OssBucketGetLifeCycle(ctx, v.Name)
		if err != nil {
			log.L(ctx).Debugf("OssBucketGetLifeCycle err:%v", err)
		}
		info, err := c.OssBucketGetInfo(ctx, v.Name)
		if err != nil {
			log.L(ctx).Debugf("OssBucketGetInfo err:%v", err)
		}
		objects, _, err := c.OssBucketObjectList(ctx, v.Name, nil, nil)
		if err != nil {
			log.L(ctx).Debugf("OssObjectList err:%v", err)
		}

		lock.Lock()
		defer lock.Unlock()
		entry := BucketListUnit{
			Property:  v,
			ACL:       acl,
			Lifecycle: lifecycle,
			Info:      info,
			ObjectNum: len(objects),
		}
		resp.Info = append(resp.Info, entry)
		return waitgroup.NewResult(nil, nil)
	})

	return resp, nil
}

type BucketCreateRequest struct {
	BucketName         string //存储桶名
	ACL                string //访问控制
	StorageClass       string //存储类型
	DataRedundancyType string //数据冗余类型
}

type BucketCreateResponse struct {
}

// 创建存储桶
func BucketCreate(ctx context.Context, c *aliyun.Oss, req *BucketCreateRequest) (*BucketCreateResponse, error) {
	resp := &BucketCreateResponse{}

	if err := c.OssBucketCreate(ctx,
		req.BucketName,
		req.ACL,
		req.StorageClass,
		req.DataRedundancyType); err != nil {
		return nil, err
	}

	return resp, nil
}

type BucketDeleteRequest struct {
	BucketName string //存储桶名
}

type BucketDeleteResponse struct {
}

func BucketDelete(ctx context.Context, c *aliyun.Oss, req *BucketDeleteRequest) (*BucketDeleteResponse, error) {
	resp := &BucketDeleteResponse{}

	if err := c.OssBucketDelete(ctx, req.BucketName); err != nil {
		return nil, err
	}

	return resp, nil
}

type BucketObjectListRequest struct {
	PagingParam
	BucketName string  //存储桶名
	Prefix     *string //前缀
	Delimiter  *string // 路径分隔字符
	Fuzzy      string  // 模糊搜索
}

type BucketObjectListResponse struct {
	Contents     []oss.ObjectProperties //对象列表
	CommonPrefix []string               //文件夹列表
	Total        int                    //所有对象总和(包含文件和对象)
}

// BucketObjectList 查询存储桶对象列表, 支持模糊搜索，文件夹/文件组合分页查询
func BucketObjectList(ctx context.Context, c *aliyun.Oss, req *BucketObjectListRequest) (*BucketObjectListResponse, error) {
	resp := &BucketObjectListResponse{}

	ret, commonPrefix, err := c.OssBucketObjectList(
		ctx,
		req.BucketName,
		req.Prefix,
		req.Delimiter)
	if err != nil {
		return nil, err
	}

	// 组合分页
	type CombinedItem struct {
		Directory *string               // 目录
		File      *oss.ObjectProperties // 文件
	}
	var combinedItems []CombinedItem

	filterRet := make([]oss.ObjectProperties, 0, len(ret))
	for _, v := range ret {
		if req.Fuzzy != "" && !strings.Contains(v.Key, req.Fuzzy) {
			continue
		}

		if typeutil.StringValue(req.Prefix) != "" {
			if v.Key == typeutil.StringValue(req.Prefix) {
				continue
			}
			v.Key = strings.TrimPrefix(v.Key, typeutil.StringValue(req.Prefix))
		}

		filterRet = append(filterRet, v)
	}
	filterCommonPrefix := make([]string, 0, len(commonPrefix))
	for _, v := range commonPrefix {
		if req.Fuzzy != "" && !strings.Contains(v, req.Fuzzy) {
			continue
		}
		if typeutil.StringValue(req.Prefix) != "" {
			v = strings.TrimPrefix(v, typeutil.StringValue(req.Prefix))
		}

		filterCommonPrefix = append(filterCommonPrefix, v)
	}
	sort.SliceStable(filterRet, func(i, j int) bool {
		return filterRet[i].Key < filterRet[j].Key
	})
	sort.SliceStable(filterCommonPrefix, func(i, j int) bool {
		return filterCommonPrefix[i] < filterCommonPrefix[j]
	})
	for i := range filterCommonPrefix {
		combinedItems = append(combinedItems, CombinedItem{Directory: typeutil.String(filterCommonPrefix[i])})
	}
	for i := range filterRet {
		combinedItems = append(combinedItems, CombinedItem{File: &filterRet[i]})
	}

	resp.Total = len(combinedItems)
	s, e := paging.Index(resp.Total, typeutil.Int32ValueToInt(req.PageNumber)-1, typeutil.Int32ValueToInt(req.PageSize))
	pageItems := combinedItems[s:e]
	for _, v := range pageItems {
		if v.Directory != nil {
			resp.CommonPrefix = append(resp.CommonPrefix, *v.Directory)
		}
		if v.File != nil {
			resp.Contents = append(resp.Contents, *v.File)
		}
	}

	return resp, nil
}
