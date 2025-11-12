package ictyun

type EbsListRequest struct {
	PagingParam
	RegionID string `json:"regionID"`
}

type EBSDiskInfo struct {
	InstanceStatus string `json:"instanceStatusByID"`
	DiskName       string `json:"diskName"`
	DiskFreeze     bool   `json:"diskFreeze"`
	Attachments    []struct {
		InstanceID   string `json:"instanceID"`
		Device       string `json:"device"`
		AttachmentID string `json:"attachmentID"`
	} `json:"attachmentsByID"`
	DiskMode       string `json:"diskMode"`
	MultiAttach    bool   `json:"multiAttach"`
	InstanceID     string `json:"instanceIDByID"`
	ProjectID      string `json:"projectID"`
	RegionID       string `json:"regionID"`
	DiskType       string `json:"diskType"`
	ExpireTime     int64  `json:"expireTime"`
	IsEncrypt      bool   `json:"isEncrypt"`
	DiskSize       int    `json:"diskSize"`
	IsPackaged     bool   `json:"isPackaged"`
	DiskStatus     string `json:"diskStatus"`
	AzName         string `json:"azName"`
	IsSystemVolume bool   `json:"isSystemVolumeByID"`
	InstanceName   string `json:"instanceNameByID"`
	CreateTime     int64  `json:"createTime"`
	DiskID         string `json:"diskID"`
}

type EbsListResponse struct {
	ReturnObj struct {
		CurrentCount int           `json:"currentCount"`
		TotalCount   int           `json:"totalCount"`
		TotalPage    int           `json:"totalPage"`
		DiskList     []EBSDiskInfo `json:"diskList"`
		DiskTotal    int           `json:"diskTotal"`
	} `json:"returnObj"`
}

type EbsGetByNameRequest struct {
	PagingParam
	RegionID string `json:"regionID"` //资源池ID
	DiskName string `json:"diskName"` //云硬盘名称
}

type EbsGetByIDRequest struct {
	RegionID string `json:"regionID"` //资源池ID
	DskID    string `json:"diskID"`   //云硬盘ID
}

type EbsGetResponse struct {
	ReturnObj EBSDiskInfo `json:"returnObj"`
}

type EbsCreateRequest struct {
	RegionID        string  `json:"regionID"`        //客户端存根
	ClientToken     string  `json:"clientToken"`     //资源池ID
	MultiAttach     string  `json:"multiAttach"`     //是否多云主机挂载，默认false，XSSD类型盘不支持多挂载
	IsEncrypt       bool    `json:"isEncrypt"`       //是否加密盘，默认false，XSSD类型盘不支持加密
	LmsUUID         string  `json:"kmsUUID"`         //如果是加密盘，需要提供kms的uuid
	ProjectID       *string `json:"projectID"`       //企业项目ID,默认为0
	DiskMode        string  `json:"diskMode"`        //磁盘模式，VBD/ISCSI/FCSAN，XSSD类型盘不支持ISCSI和FCSAN
	DiskType        string  `json:"diskType"`        //磁盘类型，SATA/SAS/SSD-genric/SSD/FAST-SSD，极速ssd类型盘（FAST-SSD）不支持iscsi、只有高io类型（SAS）支持fcsan，XSSD类型盘不支持多挂载，加密，ISCSI和FCSAN
	DiskName        string  `json:"diskName"`        //磁盘命名，单账户单资源池下，命名需唯一
	DiskSize        int     `json:"diskSize"`        //磁盘大小，单位GB
	OnDemand        bool    `json:"onDemand"`        //是否按需下单。默认为true
	CycleType       string  `json:"cycleType"`       //包周期类型，year/month。onDemand为false时，必须指定
	CycleCount      int     `json:"cycleCount"`      //包周期数。onDemand为false时必须指定。周期最大长度不能超过5年
	ImageID         string  `json:"imageID"`         //镜像ID
	AzName          string  `json:"azName"`          //多可用区资源池下，必须指定可用区
	ProvisionedIops int     `json:"provisionedIops"` //XSSD类型盘的预配置iops值，最小值为1，其他类型的盘不支持设置
}

type EbsCreateResponse struct {
	ReturnObj struct {
		MasterResourceStatus string `json:"masterResourceStatus"`
		RegionID             string `json:"regionID"`
		MasterOrderID        string `json:"masterOrderID"`
		MasterResourceID     string `json:"masterResourceID"`
		MasterOrderNO        string `json:"masterOrderNO"`
		Resources            []struct {
			OrderID          string `json:"orderID"`
			Status           int    `json:"status"`
			IsMaster         bool   `json:"isMaster"`
			DiskName         string `json:"diskName"`
			ResourceType     string `json:"resourceType"`
			MasterOrderID    string `json:"masterOrderID"`
			UpdateTime       int64  `json:"updateTime"`
			MasterResourceID string `json:"masterResourceID"`
			ItemValue        int    `json:"itemValue"`
			StartTime        int64  `json:"startTime"`
			CreateTime       int64  `json:"createTime"`
			DiskID           string `json:"diskID"`
		} `json:"resources"`
	} `json:"returnObj"`
}

type EbsRefundRequest struct {
	RegionID    string `json:"regionID"`    //资源池ID
	DiskID      string `json:"diskID"`      //云硬盘ID
	ClientToken string `json:"clientToken"` //客户端存根
}

type EbsRefundResponse struct {
	ReturnObj struct {
		MasterOrderID string `json:"masterOrderID"`
		RegionID      string `json:"regionID"`
		MasterOrderNO string `json:"masterOrderNO"`
	} `json:"returnObj"`
}

type EbsResizeRequest struct {
	RegionID    string `json:"regionID"`    //资源池ID
	DiskID      string `json:"diskID"`      //磁盘ID
	ClientToken string `json:"clientToken"` //客户端存根
	DiskSize    int    `json:"diskSize"`    //变配后的磁盘大小。当前仅支持变更磁盘大小,单位为GB
}

type EbsResizeResponse struct {
	ReturnObj struct {
		MasterOrderID string `json:"masterOrderID"`
		MasterOrderNO string `json:"masterOrderNO"`
	} `json:"returnObj"`
}

type EbsRenewRequest struct {
	RegionID    string `json:"regionID"`    //资源池ID
	DiskID      string `json:"diskID"`      //云硬盘ID
	ClientToken string `json:"clientToken"` //客户端存根
	CycleType   string `json:"cycleType"`   //包周期类型，year/month
	CycleCount  int    `json:"cycleCount"`  //包周期数
}

type EbsRenewResponse struct {
	ReturnObj struct {
		MasterOrderID string `json:"masterOrderID"`
		MasterOrderNO string `json:"masterOrderNO"`
	} `json:"returnObj"`
}

type EbsUpdateRequest struct {
	RegionID string `json:"regionID"` //资源池
	DiskID   string `json:"diskID"`   //云硬盘ID
	DiskName string `json:"diskName"` //云硬盘名字
}

type EbsUpdateResponse struct{}

type EbsAttachRequest struct {
	RegionID   string `json:"regionID"`   //资源池ID
	DiskID     string `json:"diskID"`     //磁盘ID
	InstanceID string `json:"instanceID"` //云主机ID
}

type EbsAttachResponse struct {
	ReturnObj struct {
		DiskRequestID string `json:"diskRequestID"`
		DiskJobID     []struct {
			Jobid interface{} `json:"jobid"`
		} `json:"diskJobID"`
	} `json:"returnObj"`
}

type EbsDetachRequest struct {
	RegionID   string `json:"regionID"`   //资源池ID
	DiskID     string `json:"diskID"`     //磁盘ID
	InstanceID string `json:"instanceID"` //云主机ID
}

type EbsDetachResponse struct {
	ReturnObj struct {
		DiskRequestID string `json:"diskRequestID"`
		DiskJobID     []struct {
			Jobid interface{} `json:"jobid"`
		} `json:"diskJobID"`
	} `json:"returnObj"`
}
