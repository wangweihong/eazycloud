package ictyun

type EcsSnapshotListRequest struct {
	PagingParam
	RegionID       string `json:"regionID,omitempty" description:"资源池ID" required:"true"`
	ProjectID      string `json:"projectID,omitempty" description:"企业项目ID" required:"false"`
	InstanceID     string `json:"instanceID,omitempty" description:"云主机ID" required:"false"`
	SnapshotStatus string `json:"snapshotStatus,omitempty" description:"	云主机快照状态，取值范围：pending（创建中），available（可用），restoring（恢复中），error（错误）" required:"false"`
	SnapshotID     string `json:"snapshotID,omitempty" description:"云主机快照ID，" required:"false"`
	QueryContent   string `json:"queryContent,omitempty" description:"	模糊查询查询内容,（匹配字段：instanceID、snapshotID、snapshotName" required:"false"`
	SnapshotName   string `json:"snapshotName,omitempty" description:"云主机快照名称" required:"false"`
}

type EcsSnapshotListResponse struct {
	ReturnObj struct {
		CurrentCount int `json:"currentCount"`
		TotalCount   int `json:"totalCount"`
		TotalPage    int `json:"totalPage"`
		Results      []struct {
			SnapshotID          string `json:"snapshotID"`
			InstanceID          string `json:"instanceID"`
			InstanceName        string `json:"instanceName"`
			SnapshotName        string `json:"snapshotName"`
			AzName              string `json:"azName"`
			InstanceStatus      string `json:"instanceStatus"`
			SnapshotStatus      string `json:"snapshotStatus"`
			SnapshotDescription string `json:"snapshotDescription"`
			ProjectID           string `json:"projectID"`
			CreatedTime         string `json:"createdTime"`
			UpdatedTime         string `json:"updatedTime"`
			Members             []struct {
				DiskType           string `json:"diskType"`
				DiskID             string `json:"diskID"`
				DiskName           string `json:"diskName"`
				IsBootable         bool   `json:"isBootable"`
				IsEncrypt          bool   `json:"isEncrypt"`
				DiskSize           *int   `json:"diskSize"`
				DiskSnapshotID     string `json:"diskSnapshotID"`
				DiskSnapshotStatus string `json:"diskSnapshotStatus"`
			} `json:"members "`
		} `json:"results"`
	} `json:"returnObj"`
}

type EcsSnapshotUpdateRequest struct {
	RegionID            string `json:"regionID,omitempty" description:"资源池ID" required:"true"`
	SnapshotID          string `json:"snapshotID,omitempty" description:"云主机快照ID，长度限制2~63字符" required:"true"`
	SnapshotName        string `json:"snapshotName,omitempty" description:"云主机快照名称，长度限制2~63字符" required:"true"`
	SnapshotDescription string `json:"snapshotDescription,omitempty" description:"云主机快照描述" required:"false"`
}

type EcsSnapshotUpdateResponse struct {
	ReturnObj struct {
		SnapshotID string `json:"snapshotID"`
	} `json:"returnObj"`
}

type EcsSnapshotStatusRequest struct {
	RegionID   string `json:"regionID,omitempty" description:"资源池ID" required:"true"`
	SnapshotID string `json:"snapshotID,omitempty" description:"快照ID" required:"true"`
}

type EcsSnapshotStatusResponse struct {
	ReturnObj struct {
		SnapshotStatus string `json:"snapshotStatus"`
	} `json:"returnObj"`
}

type EcsSnapshotDeleteRequest struct {
	RegionID   string `json:"regionID,omitempty" description:"资源池ID" required:"true"`
	SnapshotID string `json:"snapshotID,omitempty" description:"云主机快照ID" required:"true"`
}

type EcsSnapshotDeleteResponse struct {
	ReturnObj struct {
		SnapshotID string `json:"snapshotID"`
	} `json:"returnObj"`
}

type EcsSnapshotRestoreRequest struct {
	RegionID   string `json:"regionID,omitempty" description:"资源池ID" required:"true"`
	SnapshotID string `json:"snapshotID,omitempty" description:"云主机快照ID" required:"true"`
}

type EcsSnapshotRestoreResponse struct {
	ReturnObj struct {
		SnapshotID string `json:"snapshotID"`
	} `json:"returnObj"`
}

type EcsSnapshotCreateInstanceRequest struct {
	RegionID        string   `json:"regionID,omitempty" description:"资源池ID" required:"true"`
	ClientToken     string   `json:"clientToken,omitempty" description:"客户端存根" required:"true"`
	InstanceName    string   `json:"instanceName,omitempty" description:"云主机名称" required:"true"`
	DisplayName     string   `json:"displayName,omitempty" description:"云主机显示名称，长度为2-63字符" required:"true"`
	SnapshotID      string   `json:"snapshotID,omitempty" description:"	云主机快照ID，" required:"true"`
	NetworkCardList string   `json:"networkCardList,omitempty" description:"网卡信息列表" required:"true"`
	ExtIP           string   `json:"extIP,omitempty" description:"是否使用弹性公网IP" required:"true"`
	VpcID           string   `json:"vpcID,omitempty" description:"资源池ID" required:"true"`
	OnDemand        bool     `json:"onDemand,omitempty" description:"资源池ID" required:"true"`
	SecGroupList    []string `json:"secGroupList,omitempty" description:"购买方式，取值范围：false（按周期），true（按需" required:"true"`
	IpVersion       string   `json:"ipVersion,omitempty" description:"弹性IP版本，取值范围" required:"false"`
	Bandwidth       int      `json:"bandwidth,omitempty" description:"带宽大小单位为Mbit/s ，取值范围:[1~2000]" required:"false"`
	Ipv6AddressID   string   `json:"ipv6AddressID,omitempty" description:"弹性公网IPv6的ID" required:"false"`
	EipID           string   `json:"eipID,omitempty" description:"弹性公网IP的ID" required:"false"`
	AffinityGroupID string   `json:"affinityGroupID,omitempty" description:"云主机组ID" required:"false"`
	KeyPairID       string   `json:"keyPairID,omitempty" description:"密钥对ID" required:"false"`
	UserPassword    string   `json:"userPassword,omitempty" description:"用户密码，满足以下规则：长度在8～30个字符；必须包含大写字母、小写字母、数字以及特殊符号中的三项；" required:"false"`
	CycleCount      int      `json:"cycleCount,omitempty" description:"订购时长" required:"false"`
	CycleType       string   `json:"cycleType,omitempty" description:"订购周期类型，取值范围：MONTH：按月YEAR：按年" required:"false"`
	AutoRenewStatus int      `json:"autoRenewStatus,omitempty" description:"是否自动续订" required:"false"`
	UserData        string   `json:"userData,omitempty" description:"用户自定义数据，需要以Base64方式编码" required:"false"`
}

type EcsSnapshotCreateInstanceResponse struct {
	ReturnObj struct {
		RegionID         string `json:"regionID"`
		MasterOrderID    string `json:"masterOrderID"`
		MasterResourceID string `json:"masterResourceID"`
		MasterOrderNO    string `json:"masterOrderNO"`
	} `json:"returnObj"`
}

type EcsSnapshotStatisticRequest struct {
	RegionID       string `json:"regionID,omitempty" description:"资源池ID" required:"false"`
	InstanceIDList string `json:"instanceIDList,omitempty" description:"云主机ID列表,多台使用英文逗号分割" required:"true"`
}

type EcsSnapshotStatisticResponse struct {
	ReturnObj []struct {
		InstanceID string `json:"instanceID"`
		Count      *int   `json:"count"`
	} `json:"returnObj"`
}

type EcsSnapshotPolicyDeleteRequest struct {
	RegionID         string `json:"regionID,omitempty" description:"资源池ID" required:"true"`
	SnapshotPolicyID string `json:"snapshotPolicyID,omitempty" description:"" required:"true"`
}

type EcsSnapshotPolicyDeleteResponse struct {
	ReturnObj struct {
		SnapshotPolicyID string `json:"snapshotPolicyID"`
	} `json:"returnObj"`
}

type EcsSnapshotPolicyExecuteRequest struct {
	RegionID         string `json:"regionID,omitempty" description:"资源池ID" required:"true"`
	SnapshotPolicyID string `json:"snapshotPolicyID,omitempty" description:"" required:"true"`
}

type EcsSnapshotPolicyExecuteResponse struct {
	ReturnObj struct {
		SnapshotPolicyID string `json:"snapshotPolicyID"`
	} `json:"returnObj"`
}

type EcsSnapshotPolicyListRequest struct {
	PagingParam
	RegionID             string `json:"regionID,omitempty" description:"资源池ID" required:"true"`
	SnapshotPolicyStatus *int   `json:"snapshotPolicyStatus,omitempty" description:"快照策略状态，是否启用，取值范围：0（不启用），1（启用）" required:"false"`
	QueryContent         string `json:"queryContent,omitempty" description:"模糊匹配查询内容（匹配字段：snapshotPolicyID、snapshotPolicyName）" required:"false"`
}

type EcsSnapshotPolicyListResponse struct {
	ReturnObj struct {
		CurrentCount       int `json:"currentCount"`
		TotalCount         int `json:"totalCount"`
		TotalPage          int `json:"totalPage"`
		SnapshotPolicyList []struct {
			SnapshotPolicyID     string `json:"snapshotPolicyID"`
			SnapshotPolicyStatus *int   `json:"snapshotPolicyStatus"`
			SnapshotPolicyName   string `json:"snapshotPolicyName"`
			SnapshotTime         string `json:"snapshotTime"`
			RetentionType        string `json:"retentionType"`
			RetentionDay         string `json:"retentionDay"`
			RetentionNum         *int   `json:"retentionNum"`
			CycleType            string `json:"cycleType"`
			CycleDay             string `json:"cycleDay"`
			CycleWeek            string `json:"cycleWeek"`
			ResourceCount        *int   `json:"resourceCount"`
		} `json:"snapshotPolicyList"`
	} `json:"returnObj"`
}

type EcsSnapshotPolicyGetRequest struct {
	RegionID         string `json:"regionID,omitempty" description:"资源池ID" required:"true"`
	SnapshotPolicyID string `json:"snapshotPolicyID,omitempty" description:"" required:"true"`
}

type EcsSnapshotPolicyGetResponse struct {
	ReturnObj struct {
		RetentionType        string `json:"retentionType"`
		SnapshotPolicyStatus *int   `json:"snapshotPolicyStatus"`
		RetentionDay         string `json:"retentionDay"`
		SnapshotTime         string `json:"snapshotTime"`
		CycleDay             string `json:"cycleDay"`
		SnapshotPolicyName   string `json:"snapshotPolicyName"`
		CycleWeek            string `json:"cycleWeek"`
		SnapshotPolicyID     string `json:"snapshotPolicyID"`
		RetentionNum         *int   `json:"retentionNum"`
		CycleType            string `json:"cycleType"`
		ResourceCount        *int   `json:"resourceCount"`
	} `json:"returnObj"`
}

type EcsSnapshotPolicyInstanceListRequest struct {
	PagingParam
	RegionID         string `json:"regionID,omitempty" description:"资源池ID" required:"true"`
	SnapshotPolicyID string `json:"snapshotPolicyID,omitempty" description:"" required:"true"`
}

type EcsSnapshotPolicyInstanceListResponse struct {
	ReturnObj struct {
		CurrentCount int `json:"currentCount"`
		TotalCount   int `json:"totalCount"`
		TotalPage    int `json:"totalPage"`
		InstanceList []struct {
			InstanceID     string `json:"instanceID"`
			InstanceName   string `json:"instanceName"`
			DisplayName    string `json:"displayName"`
			InstanceStatus string `json:"instanceStatus"`
			VolumeCount    *int   `json:"volumeCount"`
		} `json:"instanceList"`
	} `json:"returnObj"`
}

type EcsSnapshotPolicyInstanceUnbindRequest struct {
	RegionID         string `json:"regionID,omitempty" description:"资源池ID" required:"true"`
	SnapshotPolicyID string `json:"snapshotPolicyID,omitempty" description:"云主机快照策略ID" required:"true"`
	InstanceIDs      string `json:"instanceIDs,omitempty" description:"云主机ID列表" required:"true"`
}

type EcsSnapshotPolicyInstanceUnbindResponse struct {
	ReturnObj []struct {
		InstanceIDList []string `json:"instanceIDList"`
	} `json:"returnObj"`
}

type EcsSnapshotPolicyInstanceBindRequest struct {
	RegionID         string `json:"regionID,omitempty" description:"资源池ID" required:"true"`
	SnapshotPolicyID string `json:"snapshotPolicyID,omitempty" description:"云主机快照策略ID" required:"true"`
	InstanceIDs      string `json:"instanceIDs,omitempty" description:"云主机ID列表" required:"true"`
}

type EcsSnapshotPolicyInstanceBindResponse struct {
	ReturnObj []struct {
		InstanceIDList []string `json:"instanceIDList"`
	} `json:"returnObj"`
}

type EcsSnapshotPolicyEnableRequest struct {
	RegionID         string `json:"regionID,omitempty" description:"资源池ID" required:"true"`
	SnapshotPolicyID string `json:"snapshotPolicyID,omitempty" description:"云主机快照策略ID" required:"true"`
}

type EcsSnapshotPolicyEnableResponse struct {
	ReturnObj struct {
		SnapshotPolicyID     string `json:"snapshotPolicyID"`
		RetentionNum         string `json:"retentionNum"`
		SnapshotTime         string `json:"snapshotTime"`
		SnapshotPolicyStatus *int   `json:"snapshotPolicyStatus"`
		CycleType            string `json:"cycleType"`
		CycleDay             string `json:"cycleDay"`
		RetentionType        string `json:"retentionType"`
		RetentionDay         string `json:"retentionDay"`
		SnapshotPolicyName   string `json:"snapshotPolicyName"`
		CycleWeek            string `json:"cycleWeek"`
	} `json:"returnObj"`
}

type EcsSnapshotPolicyCreateRequest struct {
	RegionID             string `json:"regionID,omitempty" description:"资源池ID" required:"true"`
	SnapshotPolicyName   string `json:"snapshotPolicyName,omitempty" description:"云主机快照策略名称" required:"true"`
	SnapshotTime         string `json:"snapshotTime,omitempty" description:"快照整点时间，时间取值范围：0~23" required:"false"`
	CycleType            string `json:"cycleType,omitempty" description:"云主机快照周期类型，取值范围：day（天），week（周）" required:"true"`
	CycleDay             *int   `json:"cycleDay,omitempty" description:"快照周期（天），取值范围：[1, 10]" required:"false"`
	CycleWeek            string `json:"cycleWeek,omitempty" description:"快照周期（星期），星期取值范围：0~6" required:"false"`
	RetentionType        string `json:"retentionType,omitempty" description:"云主机快照保留类型，取值范围：date（按时间保存），num（按数量保存）" required:"true"`
	RetentionDay         *int   `json:"retentionDay,omitempty" description:"云主机快照保留天数，单位为天，取值范围：[1, 365" required:"false"`
	RetentionNum         *int   `json:"retentionNum,omitempty" description:"云主机快照保留数量，取值范围：[1, 30]" required:"false"`
	SnapshotPolicyStatus *int   `json:"snapshotPolicyStatus,omitempty" description:"快照策略状态，是否启用，取值范围：0（不启用），1（启用）" required:"false"`
}

type EcsSnapshotPolicyCreateResponse struct {
	ReturnObj struct {
		SnapshotPolicyID     string `json:"snapshotPolicyID"`
		RetentionNum         string `json:"retentionNum"`
		SnapshotTime         string `json:"snapshotTime"`
		SnapshotPolicyStatus *int   `json:"snapshotPolicyStatus"`
		CycleType            string `json:"cycleType"`
		CycleDay             string `json:"cycleDay"`
		RetentionType        string `json:"retentionType"`
		RetentionDay         *int   `json:"retentionDay"`
		SnapshotPolicyName   string `json:"snapshotPolicyName"`
		CycleWeek            string `json:"cycleWeek"`
	} `json:"returnObj"`
}

type EcsSnapshotPolicyDisableRequest struct {
	RegionID         string `json:"regionID,omitempty" description:"资源池ID" required:"true"`
	SnapshotPolicyID string `json:"snapshotPolicyID,omitempty" description:"云主机快照策略I" required:"true"`
}

type EcsSnapshotPolicyDisableResponse struct {
	ReturnObj struct {
		SnapshotPolicyID     string `json:"snapshotPolicyID"`
		RetentionNum         string `json:"retentionNum"`
		SnapshotTime         string `json:"snapshotTime"`
		SnapshotPolicyStatus *int   `json:"snapshotPolicyStatus"`
		CycleType            string `json:"cycleType"`
		CycleDay             string `json:"cycleDay"`
		RetentionType        string `json:"retentionType"`
		RetentionDay         string `json:"retentionDay"`
		SnapshotPolicyName   string `json:"snapshotPolicyName"`
		CycleWeek            string `json:"cycleWeek"`
	} `json:"returnObj"`
}

type EcsSnapshotPolicyUpdateRequest struct {
	RegionID           string `json:"regionID,omitempty" description:"资源池ID" required:"true"`
	SnapshotPolicyID   string `json:"snapshotPolicyID,omitempty" description:"" required:"true"`
	SnapshotPolicyName string `json:"snapshotPolicyName,omitempty" description:"云主机快照策略名称" required:"false"`
	SnapshotTime       string `json:"snapshotTime,omitempty" description:"快照整点时间，时间取值范围：0~23" required:"false"`
	CycleType          string `json:"cycleType,omitempty" description:"云主机快照周期类型，取值范围：day（天），week（周）" required:"false"`
	CycleDay           *int   `json:"cycleDay,omitempty" description:"快照周期（天），取值范围：[1, 10]" required:"false"`
	CycleWeek          string `json:"cycleWeek,omitempty" description:"快照周期（星期），星期取值范围：0~6" required:"false"`
	RetentionType      string `json:"retentionType,omitempty" description:"云主机快照保留类型，取值范围：date（按时间保存），num（按数量保存）" required:"false"`
	RetentionDay       *int   `json:"retentionDay,omitempty" description:"云主机快照保留天数，单位为天，取值范围：[1, 365" required:"false"`
	RetentionNum       *int   `json:"retentionNum,omitempty" description:"云主机快照保留数量，取值范围：[1, 30]" required:"false"`
}

type EcsSnapshotPolicyUpdateResponse struct {
	ReturnObj struct {
		SnapshotPolicyID     string `json:"snapshotPolicyID"`
		RetentionNum         string `json:"retentionNum"`
		SnapshotTime         string `json:"snapshotTime"`
		SnapshotPolicyStatus *int   `json:"snapshotPolicyStatus"`
		CycleType            string `json:"cycleType"`
		CycleDay             string `json:"cycleDay"`
		RetentionType        string `json:"retentionType"`
		RetentionDay         string `json:"retentionDay"`
		SnapshotPolicyName   string `json:"snapshotPolicyName"`
		CycleWeek            string `json:"cycleWeek"`
	} `json:"returnObj"`
}

type EcsSnapshotPolicyTaskListRequest struct {
	PagingParam
	RegionID         string `json:"regionID,omitempty" description:"资源池ID" required:"true"`
	SnapshotPolicyID string `json:"snapshotPolicyID,omitempty" description:"" required:"true"`
}

type EcsSnapshotPolicyTaskListResponse struct {
	ReturnObj struct {
		CurrentCount *int `json:"currentCount"`
		TotalCount   int  `json:"totalCount"`
		TotalPage    int  `json:"totalPage"`
		TaskList     []struct {
			InstanceID   string `json:"instanceID"`
			TaskID       string `json:"taskID"`
			TaskStatus   string `json:"taskStatus"`
			SnapshotName string `json:"snapshotName"`
			CreateTime   string `json:"createTime"`
			CompleteTime string `json:"completeTime"`
		} `json:"taskList"`
	} `json:"returnObj"`
}
