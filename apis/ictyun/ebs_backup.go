package ictyun

type EbsBackupListRequest struct {
	PagingParam
	RegionID     string  `json:"regionID"`     //资源池ID
	RepositoryID *string `json:"repositoryID"` //仓库ID
	VolumeID     *string `json:"volumeID"`     //云硬盘iD
	VolumeName   *string `json:"volumeName"`   //云硬盘名称，模糊过滤
	BackupName   *string `json:"backupName"`   //云硬盘备份名称，模糊过滤
}

type EbsBackupListResponse struct {
	ReturnObj struct {
		CurrentCount int `json:"currentCount"`
		TotalCount   int `json:"totalCount"`
		BackupList   []struct {
			CmkUUID           string      `json:"cmkUUID"`
			VolumeName        string      `json:"volumeName"`
			UsedSize          int         `json:"usedSize"`
			RepositoryID      string      `json:"repositoryID"`
			VMName            string      `json:"vMName"`
			Size              int         `json:"size"`
			RegionID          string      `json:"regionID"`
			Encrypted         bool        `json:"encrypted"`
			RepositoryName    string      `json:"repositoryName"`
			RestoreDate       interface{} `json:"restoreDate"`
			Status            string      `json:"status"`
			ProjectID         string      `json:"projectID"`
			Description       string      `json:"description"`
			VMID              string      `json:"vMID"`
			AzName            string      `json:"azName"`
			CreatedDate       int         `json:"createdDate"`
			Freeze            bool        `json:"freeze"`
			Paas              bool        `json:"paas"`
			BackupName        string      `json:"backupName"`
			RestoreFinishDate interface{} `json:"restoreFinishDate"`
			VolumeType        string      `json:"volumeType"`
			VolumeID          string      `json:"volumeID"`
			UpdatedDate       interface{} `json:"updatedDate"`
			FinishDate        interface{} `json:"finishDate"`
			BackupID          string      `json:"backupID"`
		} `json:"backupList"`
	} `json:"returnObj"`
	Message     string `json:"message"`
	Description string `json:"description"`
	StatusCode  int    `json:"statusCode"`
}

type EbsBackupRestoreRequest struct {
	RegionID string `json:"regionID"` //资源池ID
	BackupID string `json:"backupID"` //备份ID
	VolumeID string `json:"volumeID"` //云硬盘ID
}

type EbsBackupRestoreResponse struct {
}

type EbsBackupDeleteRequest struct {
	RegionID string `json:"regionID"` //资源池ID
	BackupID string `json:"backupID"` //备份ID
}

type EbsBackupDeleteResponse struct {
}

type EbsBackupCreateRequest struct {
	RegionID     string  `json:"regionID"`     // 资源池ID
	VolumeID     string  `json:"volumeID"`     //云硬盘ID
	RepositoryID string  `json:"repositoryID"` //备份存储库ID
	Name         string  `json:"name"`         //备份名称
	Description  *string `json:"description"`  //云硬盘备份描述
}

type EbsBackupCreateResponse struct {
	ReturnObj struct {
		CmkUUID           string `json:"cmkUUID"`
		VolumeName        string `json:"volumeName"`
		UsedSize          int    `json:"usedSize"`
		RepositoryID      string `json:"repositoryID"`
		Pass              bool   `json:"pass"`
		Size              int    `json:"size"`
		RegionID          string `json:"regionID"`
		Encrypted         bool   `json:"encrypted"`
		RepositoryName    string `json:"repositoryName"`
		RestoreDate       int    `json:"restoreDate"`
		Status            string `json:"status"`
		Description       string `json:"description"`
		VMID              string `json:"vMID"`
		AzName            string `json:"azName"`
		CreatedDate       int    `json:"createdDate"`
		Freeze            bool   `json:"freeze"`
		BackupName        string `json:"backupName"`
		RestoreFinishDate int    `json:"restoreFinishDate"`
		VolumeType        string `json:"volumeType"`
		VolumeID          string `json:"volumeID"`
		UpdatedDate       int    `json:"updatedDate"`
		FinishDate        int    `json:"finishDate"`
		BackupID          string `json:"backupID"`
		VMName            string `json:"vMName"`
		ProjectID         string `json:"projectID"`
	} `json:"returnObj"`
}

type EbsBackupPolicyListRequest struct {
	PagingParam
	RegionID   string  `json:"regionID"`   //资源池ID
	PolicyID   *string `json:"policyID"`   //备份策略ID
	PolicyName *string `json:"policyName"` //备份策略名
}

type BackupPolicy struct {
	Status         int    `json:"status"`        //状态，0-停用，1-启用
	RegionID       string `json:"regionID"`      //资源池ID
	ResourceCount  int    `json:"resourceCount"` //策略绑定的云硬盘数量
	ProjectID      string `json:"projectID"`
	ResourceIDs    string `json:"resourceIDs"`   //绑定的云硬盘ID，以逗号分隔
	CycleType      string `json:"cycleType"`     //备份周期类型，day-按天备份，week-按星期备份
	PolicyID       string `json:"policyID"`      //策略ID
	CreatedDate    int    `json:"createdDate"`   //创建时间
	PolicyName     string `json:"policyName"`    //策略名
	RetentionType  string `json:"retentionType"` //备份保留类型，num-按数量保留，date-按时间保留，all-全部保留
	RetentionDay   int    `json:"retentionDay"`  //保留天数，只有retentionType为date时返回
	RetentionNum   int    `json:"retentionNum"`  //保留数量，只有retentionType为num时返回
	CycleDay       int    `json:"cycleDay"`      //备份周期，只有cycleType为day时返回
	CycleWeek      string `json:"cycleWeek"`     //备份周期，只有cycleType为week时返回，则取值范围0-6代表星期日-星期六，如果一周有多天备份，以逗号隔开
	Time           string `json:"time"`          //备份整点时间，取值范围0-23，如果一天内多个时间节点备份，以逗号隔开
	RepositoryList []struct {
		RepositoryName string `json:"repositoryName"`
		RepositoryID   string `json:"repositoryID"`
	} `json:"repositoryList"` //策略绑定的云硬盘备份存储库列表
	AccountID string `json:"accountID"` //账户ID
}

type EbsBackupPolicyListResponse struct {
	ReturnObj struct {
		PolicyList   []BackupPolicy `json:"policyList"`
		TotalCount   int            `json:"totalCount"`
		CurrentCount int            `json:"currentCount"`
	} `json:"returnObj"`
}

type EbsBackupPolicyDiskListRequest struct {
	PagingParam
	RegionID string  `json:"regionID"` //资源池ID
	PolicyID *string `json:"policyID"` //备份策略ID
	DiskID   *string `json:"diskID"`   //云硬盘ID
	DiskName *string `json:"diskName"` //云硬盘名称
}

type EbsBackupPolicyDiskListResponse struct {
	ReturnObj struct {
		CurrentCount int           `json:"currentCount"`
		TotalCount   int           `json:"totalCount"`
		DiskList     []EBSDiskInfo `json:"diskList"`
	} `json:"returnObj"`
}

type EbsBackupPolicyExecuteRequest struct {
	RegionID string `json:"regionID"` //资源池ID
	PolicyID string `json:"policyID"` //备份策略ID
}

type EbsBackupPolicyExecuteResponse struct {
}

type EbsBackupPolicyEnableRequest struct {
	RegionID string `json:"regionID"` //资源池ID
	PolicyID string `json:"policyID"` //备份策略ID
}

type EbsBackupPolicyEnableResponse struct {
}

type EbsBackupPolicyDisableRequest struct {
	RegionID string `json:"regionID"` //资源池ID
	PolicyID string `json:"policyID"` //备份策略ID
}

type EbsBackupPolicyDisableResponse struct {
}

type EbsBackupPolicyRestoreRequest struct {
	RegionID string `json:"regionID"` //资源池ID
	PolicyID string `json:"policyID"` //备份策略ID
}

type EbsBackupPolicyRestoreResponse struct {
}

type EbsBackupPolicyUpdateRequest struct {
	PolicyID              string  `json:"policyID"`              //备份策略ID
	RegionID              string  `json:"regionID"`              //资源池ID
	PolicyName            *string `json:"policyName"`            //备份策略ID
	Status                *int    `json:"status"`                //是否启用策略，0-停用，1-启用，默认0
	CycleType             *string `json:"cycleType"`             //备份周期类型，day-按天备份，week-按星期备
	CycleWeek             *string `json:"cycleWeek"`             //备份周期，只有cycleType为week时需设置，则取值范围0-6代表星期日-星期六，如果一周有多天备份，以逗号隔开
	CycleDay              *int    `json:"cycleDay"`              //备份周期，只有cycleType为day时需设置
	Time                  *string `json:"time"`                  //备份整点时间，取值范围0-23，如果一天内多个时间节点备份，以逗号隔开
	RetentionType         *string `json:"retentionType"`         //备份保留类型，num-按数量保留，date-按时间保留，all-全部保留
	RetentionNum          *int    `json:"retentionNum"`          //保留数量，只有retentionType为num时需设置,取值范围1-99999
	RetentionDay          *int    `json:"retentionDay"`          //保留天数，只有retentionType为date时需设置，取值1-99999
	RemainFirstOfCurMonth *bool   `json:"remainFirstOfCurMonth"` //是否保留每个月第一个备份，在retentionType为num时可设置，默认false
}

type EbsBackupPolicyUpdateResponse struct {
}

type EbsBackupPolicyCreateRequest struct {
	RegionID              string  `json:"regionID"`              //资源池ID
	PolicyName            string  `json:"policyName"`            //备份策略ID
	Status                *int    `json:"status"`                //是否启用策略，0-停用，1-启用，默认0
	CycleType             *string `json:"cycleType"`             //备份周期类型，day-按天备份，week-按星期备
	CycleWeek             *string `json:"cycleWeek"`             //备份周期，只有cycleType为week时需设置，则取值范围0-6代表星期日-星期六，如果一周有多天备份，以逗号隔开
	CycleDay              *int    `json:"cycleDay"`              //备份周期，只有cycleType为day时需设置
	Time                  *string `json:"time"`                  //备份整点时间，取值范围0-23，如果一天内多个时间节点备份，以逗号隔开
	RetentionType         *string `json:"retentionType"`         //备份保留类型，num-按数量保留，date-按时间保留，all-全部保留
	RetentionNum          *int    `json:"retentionNum"`          //保留数量，只有retentionType为num时需设置,取值范围1-99999
	RetentionDay          *int    `json:"retentionDay"`          //保留天数，只有retentionType为date时需设置，取值1-99999
	RemainFirstOfCurMonth *bool   `json:"remainFirstOfCurMonth"` //是否保留每个月第一个备份，在retentionType为num时可设置，默认false
}

type EbsBackupPolicyCreateResponse struct {
	ReturnObj struct {
		Status                int    `json:"status"`
		PolicyName            string `json:"policyName"`
		RetentionType         string `json:"retentionType"`
		RemainFirstOfCurMonth bool   `json:"remainFirstOfCurMonth"`
		RegionID              string `json:"regionID"`
		CycleDay              int    `json:"cycleDay"`
		RetentionNum          int    `json:"retentionNum"`
		CycleType             string `json:"cycleType"`
		Time                  string `json:"time"`
		AccountID             string `json:"accountID"`
		ProjectID             string `json:"projectID"`
	} `json:"returnObj"`
}

type EbsBackupPolicyDeleteRequest struct {
	RegionID  string `json:"regionID"`  //资源池ID
	PolicyIDs string `json:"policyIDs"` //备份策略ID,如果删除多个,请使用逗号隔开
}

type EbsBackupPolicyDeleteResponse struct {
}

type EbsBackupPolicyBindDiskRequest struct {
	RegionID  string `json:"regionID"`  //资源池ID
	PolicyID  string `json:"policyID"`  //备份策略ID
	VolumeIDs string `json:"volumeIDs"` //云硬盘ID,如果绑定多个,请使用逗号隔开
}

type EbsBackupPolicyBindDiskResponse struct {
}

type EbsBackupPolicyUnbindDiskRequest struct {
	RegionID  string `json:"regionID"`  //资源池ID
	PolicyID  string `json:"policyID"`  //备份策略ID
	VolumeIDs string `json:"volumeIDs"` //云硬盘ID,如果绑定多个,请使用逗号隔开
}

type EbsBackupPolicyUnbindDiskResponse struct {
}

type EbsBackupPolicyBindRepoRequest struct {
	RegionID     string `json:"regionID"`     //资源池ID
	PolicyIDs    string `json:"policyIDs"`    //备份策略ID
	RepositoryID string `json:"repositoryID"` //云硬盘备份存储库ID
}

type EbsBackupPolicyBindRepoResponse struct {
}

type EbsBackupPolicyUnbindRepoRequest struct {
	RegionID     string `json:"regionID"`     //资源池ID
	PolicyIDs    string `json:"policyIDs"`    //备份策略ID
	RepositoryID string `json:"repositoryID"` //云硬盘备份存储库ID
}

type EbsBackupPolicyUnbindRepoResponse struct {
}

type EbsBackupPolicyTaskListRequest struct {
	PagingParam
	RegionID   string  `json:"regionID"`   //资源池ID
	PolicyID   *string `json:"policyID"`   //备份策略ID
	TaskStatus *string `json:"taskStatus"` //备份任务状态，-1-失败，0-执行中，1-成功
}

type EbsBackupPolicyTaskListResponse struct {
	ReturnObj struct {
		TotalCount   int `json:"totalCount"`
		CurrentCount int `json:"currentCount"`
		TaskList     []struct {
			DiskName      string `json:"diskName"`
			BackupName    string `json:"backupName"`
			CompletedTime int    `json:"completedTime"`
			TaskStatus    int    `json:"taskStatus"`
			TaskID        string `json:"taskID"`
			CreatedTime   int    `json:"createdTime"`
			DiskID        string `json:"diskID"`
		} `json:"taskList"`
	} `json:"returnObj"`
}

type EbsBackupRepoListRequest struct {
	PagingParam
	RegionID       string  `json:"regionID"`       //资源池ID
	RepositoryName *string `json:"repositoryName"` //云硬盘备份存储库名称
	RepositoryID   *string `json:"repositoryID"`   //云硬盘备份存储库ID
	PolicyID       *string `json:"policyID"`       //策略ID
	Bind           *bool   `json:"bind"`           //绑定策略
	Status         *string `json:"status"`         //云硬盘备份存储库状态，active-可用，creating-创建中，error-失败，freezing-冻结，expired-已过期
}

type EbsBackupRepo struct {
	Status         string   `json:"status"`      //云硬盘备份存储库状态，active-可用，creating-创建中，error-失败，freezing-冻结，expired-已过期
	Paas           bool     `json:"paas"`        //是否支持PAAS
	ExpiredDate    int      `json:"expiredDate"` //到期时间
	PaymentType    string   `json:"paymentType"` //付费方式
	ProjectID      string   `json:"projectID"`
	RegionID       string   `json:"regionID"`
	UsedSize       int      `json:"usedSize"` //云硬盘备份存储库使用大小，单位Byte
	RepositoryID   string   `json:"repositoryID"`
	Freeze         bool     `json:"freeze"`         //是否冻结
	UpdatedDate    int      `json:"updatedDate"`    //更新时间
	AzName         string   `json:"azName"`         //可用域
	CreatedDate    int      `json:"createdDate"`    // 创建时间
	BackupList     []string `json:"backupList"`     //备份存储库下的可用的备份列表，元素为备份ID
	CustomerID     string   `json:"customerID"`     //用户ID
	RepositoryName string   `json:"repositoryName"` //备份存储库名称
	Expired        bool     `json:"expired"`        //备份存储库是否到期
	FreeSize       int      `json:"freeSize"`       //云硬盘备份存储库剩余大小，单位GB
	Size           int      `json:"size"`           //云硬盘备份存储库总容量，单位GB
}
type EbsBackupRepoListResponse struct {
	ReturnObj struct {
		CurrentCount   int             `json:"currentCount"`
		TotalCount     int             `json:"totalCount"`
		RepositoryList []EbsBackupRepo `json:"repositoryList"`
	} `json:"returnObj"`
}

type EbsBackupRepoCreateRequest struct {
	RegionID       string `json:"regionID"`       //资源池ID
	RepositoryName string `json:"repositoryName"` //云硬盘备份存储库名称
	CycleType      string `json:"cycleType"`      //本参数表示订购周期类型 ，取值范围：MONTH：按月YEAR：按年最长订购周期为3年
	CycleCount     int    `json:"cycleCount"`     //订购时长，与cycleType配合，cycleType为Month时，单位为月，cycleType为YEAR时，单位为年
	Size           int    `size:"size"`           //云硬盘备份存储库容量，单位GB，取值100-1024000，默认100
}

type EbsBackupRepoCreateResponse struct {
	ReturnObj struct {
		MasterOrderNO string `json:"masterOrderNO"`
		RegionID      string `json:"regionID"`
		MasterOrderID string `json:"masterOrderID"`
	} `json:"returnObj"`
}

type EbsBackupRepoResizeRequest struct {
	RegionID     string `json:"regionID"`     //资源池ID
	RepositoryID string `json:"repositoryID"` //云硬盘备份存储库ID
	Size         int    `json:"size"`         //云硬盘备份存储库容量，单位GB，取值100-1024000，默认100
}

type EbsBackupRepoResizeResponse struct {
	ReturnObj struct {
		MasterOrderNO string `json:"masterOrderNO"`
		RegionID      string `json:"regionID"`
		MasterOrderID string `json:"masterOrderID"`
	} `json:"returnObj"`
}

type EbsBackupRepoRenewRequest struct {
	RegionID     string `json:"regionID"`     //资源池ID
	RepositoryID string `json:"repositoryID"` //云硬盘备份存储库ID
	CycleType    int    `json:"cycleType"`    //本参数表示订购周期类型 ，取值范围：MONTH：按月YEAR：按年最长订购周期为3年
	CycleCount   int    `json:"cycleCount"`   //订购时长，与cycleType配合，cycleType为Month时，单位为月，cycleType为YEAR时，单位为年
}

type EbsBackupRepoRenewResponse struct {
	ReturnObj struct {
		MasterOrderNO string `json:"masterOrderNO"`
		RegionID      string `json:"regionID"`
		MasterOrderID string `json:"masterOrderID"`
	} `json:"returnObj"`
}

type EbsBackupRepoDeleteRequest struct {
	RegionID     string `json:"regionID"`     //资源池ID
	RepositoryID string `json:"RepositoryID"` //云硬盘备份存储库ID
}

type EbsBackupRepoDeleteResponse struct {
	ReturnObj struct {
		MasterOrderNO string `json:"masterOrderNO"`
		RegionID      string `json:"regionID"`
		MasterOrderID string `json:"masterOrderID"`
	} `json:"returnObj"`
}
