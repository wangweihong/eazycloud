package ictyun

type OssBucketListRequest struct {
	PagingParam
	RegionID string `json:"regionID"` //资源池ID
}

type OSSBucketInfo struct {
	CmkUUID      string `json:"cmkUUID"`
	StorageType  string `json:"storageType"`
	ProjectID    string `json:"projectID"`
	RegionID     string `json:"regionID"`
	Bucket       string `json:"bucket"`
	IsEncrypted  bool   `json:"isEncrypted"`
	AZPolicy     string `json:"AZPolicy"`
	CreationDate string `json:"creationDate"`
	RegionName   string `json:"regionName"`
}

type OssBucketListResponse struct {
	ReturnObj struct {
		BucketTotal  int             `json:"bucketTotal"`
		BucketList   []OSSBucketInfo `json:"bucketList"`
		PageSize     int             `json:"pageSize"`
		PageNo       int             `json:"pageNo"`
		TotalCount   int             `json:"totalCount"`
		CurrentCount int             `json:"currentCount"`
	} `json:"returnObj"`
}

type OssKeyRequest struct {
	RegionID string `json:"regionID"` //资源池ID
}

type OssKeyResponse struct {
	ReturnObj []struct {
		AccessKey  string `json:"accessKey"`
		SecretKey  string `json:"secretKey"`
		RegionName string `json:"regionName"`
		RegionID   string `json:"regionID"`
	} `json:"returnObj"`
}

type OssEndpointRequest struct {
	RegionID string `json:"regionID"` //资源池ID
}

type OssEndpointResponse struct {
	ReturnObj struct {
		Ipv6Endpoint     []string `json:"ipv6Endpoint"`
		IntranetEndpoint []string `json:"intranetEndpoint"`
		InternetEndpoint []string `json:"internetEndpoint"`
	} `json:"returnObj"`
}

type OssBucketGetRequest struct {
	RegionID string `json:"regionID"` //资源池ID
	Bucket   string `json:"bucket"`   //桶名称
}

type OssBucketGetResponse struct {
	ReturnObj struct {
		CmkUUID     string `json:"cmkUUID"`
		BucketQuota struct {
			Enabled    bool `json:"enabled"`
			MaxSize    int  `json:"maxSize"`
			MaxObjects int  `json:"maxObjects"`
			CheckOnRaw bool `json:"checkOnRaw"`
			MaxSizeKb  int  `json:"maxSizeKb"`
		} `json:"bucketQuota"`
		Tenant      string `json:"tenant"`
		Ctime       string `json:"ctime"`
		StorageType string `json:"storageType"`
		ProjectID   string `json:"projectID"`
		Mtime       string `json:"mtime"`
		Bucket      string `json:"bucket"`
		Owner       string `json:"owner"`
		Usage       struct {
		} `json:"usage"`
		NumShards         int    `json:"numShards"`
		BucketPreviewFlag int    `json:"bucketPreviewFlag"`
		AZPolicy          string `json:"AZPolicy"`
		ExplicitPlacement struct {
			DataExtraPool string `json:"dataExtraPool"`
			DataPool      string `json:"dataPool"`
			IndexPool     string `json:"indexPool"`
		} `json:"explicitPlacement"`
		Zonegroup     string `json:"zonegroup"`
		PlacementRule string `json:"placementRule"`
		IndexType     string `json:"indexType"`
	} `json:"returnObj"`
}

type OssBucketGetACLRequest struct {
	RegionID string `json:"regionID"` //资源池ID
	Bucket   string `json:"bucket"`   //桶名称
}

type OssBucketGetACLResponse struct {
	ReturnObj struct {
		Owner struct {
			DisplayName string `json:"displayName"`
			ID          string `json:"ID"`
		} `json:"owner"`
		Grants []struct {
			Grantee struct {
				Type         string `json:"type"`
				URI          string `json:"URI"`
				EmailAddress string `json:"emailAddress"`
				DisplayName  string `json:"displayName"`
				ID           string `json:"ID"`
			} `json:"grantee"`
			Permission string `json:"permission"`
		} `json:"grants"`
	} `json:"returnObj"`
}

type OssBucketCreateRequest struct {
	PagingParam
	RegionID    string `json:"regionID" description:"资源池ID" required:"true"`
	ACL         string `json:"ACL" description:"桶权限，默认为private，可选值为 private,public-read,public-read-write" required:"false"`
	Bucket      string `json:"bucket" description:"桶名称，不可为空。长度 3-63 个字符内（含）字符只能有大小写字母、数字以及英文句号" required:"true"`
	ProjectID   string `json:"projectID" description:"企业项目ID" required:"false"`
	CmkUUID     string `json:"cmkUUID" description:"cmkUUID，若 isEncrypted 为 true，则此参数必传" required:"false"`
	IsEncrypted bool   `json:"isEncrypted" description:"加密状态" required:"false"`
	StorageType string `json:"storageType" description:"存储类型，可选的值为 STANDARD, STANDARD_IA, GLACIER，分别表示标准、低频、归档" required:"false"`
	AZPolicy    string `json:"AZPolicy" description:"可选值为single-az，multi-az，默认为single-az" required:"false"`
}

type OssBucketCreateResponse struct {
}

type OssBucketDeleteRequest struct {
	PagingParam
	RegionID string `json:"regionID"` //资源池ID
	Bucket   string `json:"bucket"`   //桶名称
}

type OssBucketDeleteResponse struct {
}

type OssObjectListRequest struct {
	RegionID  string `json:"regionID" description:"资源池ID" required:"true"`
	Bucket    string `json:"bucket" description:"桶名称" required:"true"`
	Delimiter string `json:"delimiter" description:"定界符" required:"false"`
	Prefix    string `json:"prefix" description:"返回的key的前缀" required:"false"`
	MaxKeys   int    `json:"maxKeys" description:"一次返回keys的最大数目" required:"false"`
	Marker    string `json:"marker" description:"指示从哪里开始列出" required:"false"`
}

type OssObjectListResponse struct {
	ReturnObj struct {
		Name         string          `json:"name"`
		MaxKeys      int             `json:"maxKeys"`
		Prefix       string          `json:"prefix"`
		Marker       string          `json:"marker"`
		EncodingType string          `json:"encodingType"`
		IsTruncated  bool            `json:"isTruncated"`
		Contents     []OssObjectInfo `json:"contents"`
	} `json:"returnObj"`
}

type OssObjectInfo struct {
	LastModified string `json:"lastModified"`
	ETag         string `json:"ETag"`
	StorageClass string `json:"storageClass"`
	Key          string `json:"key"`
	Owner        struct {
		DisplayName string `json:"displayName"`
		ID          string `json:"ID"`
	} `json:"owner"`
	Type string `json:"type"`
	Size int    `json:"size"`
}

type OssObjectNumRequest struct {
	RegionID string `json:"regionID"` //资源池ID
	Bucket   string `json:"bucket"`   //桶名称
}

type OssObjectNumResponse struct {
	ReturnObj struct {
		ObjectsNum int `json:"objectsNum"`
	} `json:"returnObj"`
}

type OssObjectDeleteRequest struct {
	RegionID string `json:"regionID"` //资源池ID
	Bucket   string `json:"bucket"`   //桶名称
	Key      string `json:"key"`      //对象名
}

type OssObjectDeleteResponse struct {
}

type OssObjectUploadLinkRequest struct {
	RegionID string `json:"regionID"` //资源池ID
	Bucket   string `json:"bucket"`   //桶名称
	Key      string `json:"key"`      //对象名
}

type OssObjectUploadLinkResponse struct {
	ReturnObj struct {
		URL    string `json:"url"`
		Fields struct {
			Policy         string `json:"policy"`
			AWSAccessKeyID string `json:"AWSAccessKeyId"`
			Key            string `json:"key"`
			Signature      string `json:"signature"`
		} `json:"fields"`
	} `json:"returnObj"`
}

type OssObjectDownloadLinkRequest struct {
	RegionID string `json:"regionID"` //资源池ID
	Bucket   string `json:"bucket"`   //桶名称
	Key      string `json:"key"`      //对象名
}

type OssObjectDownloadLinkResponse struct {
	ReturnObj string `json:"returnObj"`
}

type OssDirectoryCreateRequest struct {
	RegionID string `json:"regionID"` //资源池ID
	Bucket   string `json:"bucket"`   //桶名称
	Key      string `json:"key"`      //对象名
}

type OssDirectoryCreateResponse struct {
}

type OssDirectoryDeleteRequest struct {
	RegionID string `json:"regionID"` //资源池ID
	Bucket   string `json:"bucket"`   //桶名称
	Key      string `json:"key"`      //对象名
}

type OssDirectoryDeleteResponse struct {
}
