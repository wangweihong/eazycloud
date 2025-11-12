package ictyun

type ImageListRequest struct {
	PagingParam
	// 资源池ID
	RegionID string `json:"regionID"`
	// 镜像可见类型
	// 取值范围（值：描述）：0：私有镜像1：公共镜像（默认值）,2：共享镜像,3：安全产品镜像,4：甄选应用镜像
	Visibility int `json:"visibility"`
	// 可用区名称
	AzName *string `json:"azName"`
	// 用于查询某个规格支持的镜像列表, 如s7.small.1
	FlavorName *string `json:"flavorName"`
	// 企业项目 ID
	ProjectID *string `json:"projectID"`
	// 用于查询某个名称的镜像
	QueryContent *string `json:"queryContent"`
	// 镜像状态
	Status *string `json:"status"`
}

type ImageEntry struct {
	Architecture    string `json:"architecture"`
	AzName          string `json:"azName"`
	BootMode        string `json:"bootMode"`
	ContainerFormat string `json:"containerFormat"`
	CreatedTime     int    `json:"createdTime"`
	Description     string `json:"description"`
	DestinationUser string `json:"destinationUser"`
	DiskFormat      string `json:"diskFormat"`
	DiskID          string `json:"diskID"`
	DiskSize        int    `json:"diskSize"`
	ImageClass      string `json:"imageClass"`
	ImageID         string `json:"imageID"`
	ImageName       string `json:"imageName"`
	ImageType       string `json:"imageType"`
	// 镜像构建云实例时支持的最大内存。为0表示没有限制
	MaximumRAM int `json:"maximumRAM"`
	// 镜像构建云实例时支持的最小内存。为0表示没有限制
	MinimumRAM       int    `json:"minimumRAM"`
	OsDistro         string `json:"osDistro"`
	OsType           string `json:"osType"`
	OsVersion        string `json:"osVersion"`
	ProjectID        string `json:"projectID"`
	SharedListLength int    `json:"sharedListLength"`
	Size             int64  `json:"size"`
	SourceServerID   string `json:"sourceServerID"`
	SourceUser       string `json:"sourceUser"`
	// 镜像状态
	// accepted：已接受共享镜像
	// active：正常
	// deactivated：已弃用
	// deactivating：弃用中
	// deleted：已删除
	// deleting：删除中
	// error：错误
	// importing：导入中
	// killed：上传出错，镜像不可读
	// pending_delete：等待删除中
	// queued：排队中
	// reactivating：取消弃用中
	// rejected：已拒绝共享镜像
	// saving：保存中
	// syncing：同步中
	// uploading：上传中
	// waiting：等待接受/拒绝共享镜像
	Status      string `json:"status"`
	Tags        string `json:"tags"`
	UpdatedTime int    `json:"updatedTime"`
	Visibility  string `json:"visibility"`
}

type ImageListResponse struct {
	ReturnObj struct {
		Images       []ImageEntry `json:"images"`
		PageNo       int          `json:"pageNo"`
		CurrentPage  int          `json:"currentPage"`
		PageSize     int          `json:"pageSize"`
		CurrentCount int          `json:"currentCount"`
		TotalCount   int          `json:"totalCount"`
	} `json:"returnObj"`
}

type ImageGetRequest struct {
	RegionID string `json:"regionID"`
	ImageID  string `json:"imageID"`
}

type ImageGetResponse struct {
	ReturnObj struct {
		Images []ImageEntry `json:"images"`
	} `json:"returnObj"`
}

type ImageDataDiskCreateRequest struct {
	RegionID    string  `json:"regionID"`
	ImageName   string  `json:"imageName"` //镜像名称。注意：长度为 2~32 个字符，
	InstanceID  string  `json:"instanceID"`
	DataDiskID  string  `json:"dataDiskID"`
	Description *string `json:"description"`
	ProjectID   *string `json:"projectID"`
}

type ImageDataDiskCreateResponse struct {
	ReturnObj struct {
		ImageID string `json:"imageID"`
	} `json:"returnObj"`
}

type ImageSystemDiskCreateRequest struct {
	RegionID    string  `json:"regionID"`    //资源池ID
	ImageName   string  `json:"imageName"`   //镜像名称。注意：长度为 2~32 个字符
	InstanceID  string  `json:"instanceID"`  //云主机 ID
	Description *string `json:"description"` //镜像描述信息.长度为 1~128 个字符
	ProjectID   *string `json:"projectID"`   //项目ID
}

type ImageSystemDiskCreateResponse struct {
	ReturnObj struct {
		TotalCount int          `json:"totalCount"`
		Images     []ImageEntry `json:"images"`
	} `json:"returnObj"`
}

type ImageDeleteRequest struct {
	RegionID string `json:"regionID"`
	ImageID  string `json:"imageID"`
}

type ImageDeleteResponse struct {
}

type ImageProperty struct {
	ImageName    string  `json:"imageName"`    //镜像名称
	OsDistro     string  `json:"osDistro"`     //操作系统的发行版名
	OsVersion    string  `json:"osVersion"`    //操作系统版本
	Architecture *string `json:"architecture"` //镜像系统架构
	BootMode     string  `json:"bootMode"`     //启动方式,bios,uefi
	Description  *string `json:"description"`  //镜像描述信息
	ImageType    string  `json:"imageType"`    //镜像种类，取值范围（值：描述）：（空或空字符串）：系统盘镜像data_disk_image：数据盘镜像
	DiskSize     *int    `json:"diskSize"`     //磁盘容量，单位为 GB"
	MaximumRAM   *int    `json:"maximumRAM"`   //最大内存，单位为 GB
	MinimumRAM   *int    `json:"minimumRAM"`   //最小内存，单位为 GB
}

type ImageImportRequest struct {
	RegionID        string        `json:"regionID"`        //资源池ID
	ImageFileSource string        `json:"imageFileSource"` //镜像文件地址，格式应为 {internetEndpoint}/{bucket}/{key
	ImageProperties ImageProperty `json:"imageProperties"` //镜像属性
	ProjectID       *string       `json:"projectID"`
}

type ImageImportResponse struct {
}

type ImageExportRequest struct {
	RegionID        string  `json:"regionID"`        //资源池ID
	ImageID         string  `json:"imageID"`         //镜像ID
	Bucket          string  `json:"bucket"`          //对象存储存储桶
	Filename        string  `json:"filename"`        //文件名
	ImageFileFormat *string `json:"imageFileFormat"` //镜像文件格式
}

type ImageExportResponse struct {
}
