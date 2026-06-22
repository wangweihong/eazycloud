package ictyun

type EcsListRequest struct {
	PagingParam

	RegionID        string  `json:"regionID"`        //资源池ID
	InstanceName    *string `json:"instanceName"`    //云主机名称，精准匹配
	AzName          *string `json:"azName"`          //可用区名称
	InstanceIDList  *string `json:"instanceIDList"`  //云主机ID列表，多台使用英文逗号分割
	VpcID           *string `json:"vpcID"`           //虚拟私有云ID
	SecurityGroupID *string `json:"securityGroupID"` //安全组ID，模糊匹配
	State           *string `json:"state"`           //云主机状态，取值范围：active（开机），shutoff（关机），expired（已到期）
	Keyword         *string `json:"keyword"`         //关键字，对部分参数进行模糊查询，包含：instanceName、displayName、instanceID、privateIP
	ProjectID       *string `json:"projectID"`       //企业项目ID
	ResourceID      *string `json:"resourceID"`      //资源ID
}

type EcsListResponse struct {
	ReturnObj struct {
		CurrentCount int       `json:"currentCount"`
		TotalCount   int       `json:"totalCount"`
		TotalPage    int       `json:"totalPage"`
		Results      []EcsInfo `json:"results"`
	} `json:"returnObj"`
}

type EcsInfo struct {
	// 云主机组信息
	AffinityGroup struct {
		AffinityGroupName string `json:"affinityGroupName"`
		AffinityGroupID   string `json:"affinityGroupID"`
		Policy            string `json:"policy"`
	} `json:"affinityGroup"`
	AvailableDay int `json:"availableDay"`
	Addresses    []struct {
		VpcName     string `json:"vpcName"`
		AddressList []struct {
			// IP地址
			Addr string `json:"addr"`
			// IP版本
			Version interface{} `json:"version"`
			// 网络类型，取值范围:fixed（内网）,floating（弹性公网）
			Type string `json:"type"`
		} `json:"addressList"`
	} `json:"addresses"`
	UpdatedTime string `json:"updatedTime"`
	ProjectID   string `json:"projectID"`
	// 镜像信息
	Image struct {
		ImageID   string `json:"imageID"`
		ImageName string `json:"imageName"`
	} `json:"image"`
	SecGroupList []struct {
		SecurityGroupName string `json:"securityGroupName"`
		SecurityGroupID   string `json:"securityGroupID"`
	} `json:"secGroupList"`
	VpcName    string `json:"vpcName"`
	AzName     string `json:"azName"`
	ZabbixName string `json:"zabbixName"`
	// 操作系统类型，取值范围:1（linux）,2（windows）,3（redhat）,4（ubuntu）,5（centos）,6（oracle）
	OsType      interface{} `json:"osType"`
	DisplayName string      `json:"displayName"`
	//云主机规格信息
	Flavor struct {
		FlavorID     string      `json:"flavorID"`
		FlavorName   string      `json:"flavorName"`
		FlavorCPU    int         `json:"flavorCPU"`
		FlavorRAM    int         `json:"flavorRAM"`
		GpuType      interface{} `json:"gpuType"`
		GpuCount     interface{} `json:"gpuCount"`
		GpuVendor    interface{} `json:"gpuVendor"`
		VideoMemSize interface{} `json:"videoMemSize"`
	} `json:"flavor"`
	NetworkCardList []struct {
		IPv4Address   string   `json:"IPv4Address"`
		IPv6Address   []string `json:"IPv6Address"`
		IsMaster      bool     `json:"isMaster"`
		SubnetCidr    string   `json:"subnetCidr"`
		NetworkCardID string   `json:"networkCardID"`
		Gateway       string   `json:"gateway"`
		SecurityGroup []string `json:"securityGroup"`
		SubnetID      string   `json:"subnetID"`
	} `json:"networkCardList"`
	// 内网IP列表
	FixedIPList []string `json:"fixedIPList"`
	VipCount    int      `json:"vipCount"`
	KeypairName string   `json:"keypairName"`
	// 云主机状态：取值范围
	// backuping: 备份中，creating: 创建中，expired: 已到期，freezing: 冻结中，rebuild: 重装，restarting: 重启中，running: 运行中，
	// starting: 开机中，stopped: 已关机，stopping: 关机中，error: 错误，snapshotting: 快照创建中
	InstanceStatus string `json:"instanceStatus"`
	//	虚拟私有云ID
	VpcID          string   `json:"vpcID"`
	AttachedVolume []string `json:"attachedVolume"`
	OnDemand       bool     `json:"onDemand"`
	// 公网IP
	FloatingIP  string `json:"floatingIP"`
	InstanceID  string `json:"instanceID"`
	ResourceID  string `json:"resourceID"`
	VipInfoList []struct {
		VipID          string `json:"vipID"`
		VipAddress     string `json:"vipAddress"`
		VipBindNicIP   string `json:"vipBindNicIP"`
		VipBindNicIPv6 string `json:"vipBindNicIPv6"`
		NicID          string `json:"nicID"`
	} `json:"vipInfoList"`
	PrivateIP    string   `json:"privateIP"`
	PrivateIPv6  string   `json:"privateIPv6"`
	ExpiredTime  string   `json:"expiredTime"`
	SubnetIDList []string `json:"subnetIDList"`
	CreatedTime  string   `json:"createdTime"`
	InstanceName string   `json:"instanceName"`
}

type EcsGetStatusListRequest struct {
	PagingParam
	RegionID       string `json:"regionID"`       //资源池ID
	InstanceIDList string `json:"instanceIDList"` //云主机ID列表，多台使用英文逗号分割
	ProjectID      string `json:"projectID"`      //企业项目ID
	AzName         string `json:"azName"`         //可用区
}

type EcsGetStatusListResponse struct {
	ReturnObj struct {
		CurrentCount int `json:"currentCount"`
		TotalCount   int `json:"totalCount"`
		TotalPage    int `json:"totalPage"`
		StatusList   []struct {
			InstanceID     string `json:"instanceID"`
			InstanceStatus string `json:"instanceStatus"`
		} `json:"statusList"`
	} `json:"returnObj"`
}

type EcsGetStatisticsRequest struct {
	RegionID  string `json:"regionID"`  //资源池ID
	ProjectID string `json:"projectID"` //企业项目ID
}

type EcsGetStatisticsResponse struct {
	ReturnObj struct {
		InstanceGetStatistics struct {
			TotalCount          int `json:"totalCount"`
			CPUCount            int `json:"cpuCount"`
			ShutdownCount       int `json:"shutdownCount"`
			ExpireCount         int `json:"expireCount"`
			MemoryCount         int `json:"memoryCount"`
			ExpireRunningCount  int `json:"expireRunningCount"`
			RunningCount        int `json:"RunningCount"`
			ExpireShutdownCount int `json:"expireShutdownCount"`
		} `json:"instanceGetStatistics"`
	} `json:"returnObj"`
}

type EcsCreateRequest struct {
	RegionID        string   `json:"regionID"`        //资源池ID
	AzName          string   `json:"azName"`          //可用区名称.如果资源池不存在可用区，则传default
	InstanceName    string   `json:"instanceName"`    //云主机名称
	DisplayName     string   `json:"displayName"`     //云主机显示名称，长度为2-63字符
	FlavorID        string   `json:"flavorID" `       //云主机规格ID
	ImageType       int      `json:"imageType"`       //镜像类型，取值范围：0（私有镜像），1（公有镜像），2（共享镜像），3（安全镜像），4（甄选镜像）
	ImageID         string   `json:"imageID"`         //镜像ID
	VpcID           string   `json:"vpcID"`           //虚拟私有云ID
	BootDiskType    string   `json:"bootDiskType"`    //系统盘类型，取值范围： SATA（普通IO）， SAS（高IO）， SSD（超高IO）， SSD-genric（通用型SSD）， FAST-SSD（极速型SSD
	BootDiskSize    int      `json:"bootDiskSize"`    //系统盘大小单位为GiB，取值范围：[40, 32768]
	ExtIP           string   `json:"extIP"`           //是否使用弹性公网IP,取值范围：0（不使用），1（自动分配），2（使用已有）
	Bandwidth       *int     `json:"bandwidth"`       //带宽大小，单位为Mbit/s，取值范围：[1, 2000]
	IPVersion       *string  `json:"ipVersion"`       //弹性IP版本，取值范围：ipv4（v4地址），ipv6（v6地址），不指定默认为ipv4
	OnDemand        bool     `json:"onDemand"`        //购买方式，取值范围：false（按周期），true（按需）
	CycleCount      *int     `json:"cycleCount"`      //订购时长
	CycleType       *string  `json:"cycleType"`       //订购周期类型，取值范围：MONTH：按月，YEAR：按年
	UserPassword    *string  `json:"userPassword"`    //用户密码                                                                         //用户密码，满足以下规则：长度在8～30个字符；必须包含大写字母、小写字母、数字以及特殊符号中的三项
	AffinityGroupID *string  `json:"affinityGroupID"` //云主机组ID
	UserData        *string  `json:"userData"`        //用户自定义数据，需要以Base64方式编码
	ProjectID       string   `json:"projectID"`       //企业项目ID
	SecGroupList    []string `json:"secGroupList"`    //安全组ID列表
	NetworkCardList []struct {
		NicName  *string `json:"nicName"`   //网卡名
		FixedIP  *string `json:"fixedIP"`   //内网IPv4地址
		SubnetID string  `json:"subnetID"`  //子网ID
		IsMaster bool    `json:"isMaster" ` //是否主网卡,true（表示主网卡),false（表示扩展网卡）注：只能含有一个主网卡
	} `json:"networkCardList" description:"网卡列表"`
	DataDiskList []struct {
		DiskMode string `json:"diskMode"` //云硬盘属性，取值范围：FCSAN（光纤通道协议的SAN网络），ISCSI（小型计算机系统接口），VBD（虚拟块存储设备
		DiskName string `json:"diskName"` //云硬盘名称，长度限制2~63字符，不支持中文
		DiskType string `json:"diskType"` //云硬盘类型，取值范围：SATA（普通IO），SAS（高IO），SSD（超高IO），SSD-genric通用型SSD），FAST-SSD（极速型SSD)
		DiskSize int    `json:"diskSize"` //磁盘容量，单位为GB，取值范围：[40, 32768]
	} `json:"dataDiskList"` //数据盘列表
	ClientToken     string   `json:"clientToken"`     //客户端存根，用于保证订单幂等性.[可自行生成随机字符串代替]
	PayVoucherPrice *float64 `json:"payVoucherPrice"` //代金券
	AutoRenewStatus int      `json:"autoRenewStatus"` //是否自动续订，取值范围：0（不续费），1（自动续费）
}

type EcsCreateResponse struct {
	ReturnObj interface{} `json:"returnObj"`
}

type EcsRebootRequest struct {
	RegionID   string `json:"regionID"`   //资源池ID
	InstanceID string `json:"instanceID"` //弹性云主机ID
}
type EcsRebootResponse struct {
	ReturnObj struct {
		JobID string `json:"jobID"`
	} `json:"returnObj"`
}

type EcsStartRequest struct {
	RegionID   string `json:"regionID"`   //资源池ID
	InstanceID string `json:"instanceID"` //弹性云主机ID
}

type EcsStartResponse struct {
	ReturnObj struct {
		JobID string `json:"jobID"`
	} `json:"returnObj"`
}

type EcsStopRequest struct {
	RegionID   string `json:"regionID"`   //资源池ID
	InstanceID string `json:"instanceID"` //弹性云主机ID
}

type EcsStopResponse struct {
	ReturnObj struct {
		JobID string `json:"jobID"`
	} `json:"returnObj"`
}

type EcsDeleteRequest struct {
	RegionID   string `json:"regionID"`   //资源池ID
	InstanceID string `json:"instanceID"` //弹性云主机ID
}

type EcsDeleteResponse struct {
	ReturnObj struct {
		JobID string `json:"jobID"`
	} `json:"returnObj"`
}

type EcsUpdateRequest struct {
	RegionID    string `json:"regionID"`    //资源池ID
	InstanceID  string `json:"instanceID"`  //弹性云主机ID
	DisplayName string `json:"displayName"` //云主机显示名称，长度为2~63个字符
}

type EcsUpdateResponse struct {
	ReturnObj struct {
		InstanceID  string `json:"instanceID"`
		DisplayName string `json:"displayName"`
	} `json:"returnObj"`
}

type EcsResetPasswordRequest struct {
	EcsPwdInfo
	RegionID string `json:"regionID"` //资源池ID
}

type EcsResetPasswordResponse struct {
	ReturnObj struct {
		InstanceID string `json:"instanceID"`
	} `json:"returnObj"`
}

type EcsPwdInfo struct {
	InstanceID  string `json:"instanceID"`  //云主机ID
	NewPassword string `json:"newPassword"` //新密码
}

type EcsBatchResetPasswordRequest struct {
	RegionID      string       `json:"regionID"`      //资源池ID
	UpdatePwdInfo []EcsPwdInfo `json:"updatePwdInfo"` //批量更新密码信息列表
}

type EcsBatchResetPasswordResponse struct {
	ReturnObj struct {
		InstanceIDList string `json:"instanceIDList"`
	} `json:"returnObj"`
}

type EcsBatchOps struct {
	RegionID        string `json:"regionID"`       //资源池ID
	InstanceIDLists string `json:"instanceIDList"` //云主机ID列表，多台使用英文逗号分割
}

type EcsBatchResult struct {
	ReturnObj struct {
		JobIDList []struct {
			InstanceID string `json:"instanceID"`
			JobID      string `json:"jobID"`
		} `json:"jobIDList"`
	} `json:"returnObj"`
}

type EcsBatchRebootRequest EcsBatchOps
type EcsBatchRebootResponse EcsBatchResult

type EcsBatchStartRequest EcsBatchOps
type EcsBatchStartResponse EcsBatchResult

type EcsBatchStopRequest EcsBatchOps
type EcsBatchStopResponse EcsBatchResult

type EcsBatchDeleteRequest EcsBatchOps
type EcsBatchDeleteResponse EcsBatchResult

type EcsBatchCreateRequest struct {
	EcsCreateRequest
	OrderCount int `json:"orderCount"` //购买数量，取值范围：[1, 50]
}

type EcsBatchCreateResponse struct {
	EcsCreateResponse
}

type EcsGetFlavorFamilyRequest struct {
	RegionID string  `json:"regionID"` //资源池ID
	AzName   *string `json:"azName"`   //可用区列表
}
type EcsGetFlavorFamilyResponse struct {
	ReturnObj struct {
		// ["g7", "fc1", "hc1", "s2", "s7", "hm1", "fm1", "m7", "m2", "p8a", "c7", "pi7"]
		FlavorFamilyList []string `json:"flavorFamilyList"`
	} `json:"returnObj"`
}

type EcsFlavorListRequest struct {
	RegionID   string  `json:"regionID"`   //资源池ID
	AzName     *string `json:"azName"`     //可用区名称
	FlavorType *string `json:"flavorType"` //规格类型
	FlavorName *string `json:"flavorName"` //规格名称
	FlavorCPU  *int    `json:"flavorCPU"`  //VCPU个数
	FlavorRAM  *int    `json:"flavorRAM"`  //内存大小，单位为GB
	FlavorArch *string `json:"flavorArch"` //指令集架构
	// 规格系列
	// s为通用型，m为内存优化型，c为计算增强型，k代表鲲鹏，h代表海光，f代表飞腾，ip代表超高IO本地盘云主机、d代表磁盘增强云主机
	FlavorSeries *string `json:"flavorSeries"` //规格系列
	FlavorID     *string `json:"flavorID"`     //云主机规格ID

	//  自定义参数
	//是否只显示已售罄规格
	Available *bool
	Fuzzy     *string
}

type EcsFlavorInfo struct {
	GpuVendor     string  `json:"gpuVendor"`
	CPUInfo       string  `json:"cpuInfo"`
	BaseBandwidth float64 `json:"baseBandwidth"`
	FlavorName    string  `json:"flavorName"`
	VideoMemSize  float64 `json:"videoMemSize"`
	FlavorType    string  `json:"flavorType"`
	FlavorRAM     float64 `json:"flavorRAM"`
	NicMultiQueue float64 `json:"nicMultiQueue"`
	Pps           float64 `json:"pps"`
	// 规格使用的CPU数
	FlavorCPU float64 `json:"flavorCPU"`
	Bandwidth float64 `json:"bandwidth"`
	GpuType   string  `json:"gpuType"`
	FlavorID  string  `json:"flavorID"`
	GpuCount  float64 `json:"gpuCount"`
	// 用于表明该规格是否已经售罄
	// 已售罄的规格不可使用
	Available bool     `json:"available"`
	AzList    []string `json:"azList"`
	// 规格系列
	FlavorSeries     string `json:"flavorSeries"`
	FlavorSeriesName string `json:"flavorSeriesName"`
}

type EcsFlavorListResponse struct {
	ReturnObj struct {
		FlavorList []EcsFlavorInfo `json:"flavorList"`
	} `json:"returnObj"`
}

type EcsGetOrderRequest struct {
	MasterOrderID string `json:"masterOrderID"`
}

type EcsGetOrderResponse struct {
	ReturnObj struct {
		InstanceIDList []string `json:"instanceIDList"`
		// 订单状态，具体值代表含义：
		//1 待支付
		//2 已支付
		//3 完成
		//4 取消
		//5 施工失败
		//7 正在支付中
		//8 待审核
		//9 审核通过
		//10 审核未通过
		//11 撤单完成
		//12 退订中
		//13 退订完成
		//14 开通中
		//15 变更移除
		//16 自动撤单中
		//17 手动撤单中
		//18 终止中
		//22 支付失败
		//-2 待撤单
		//-1 未知
		//0 错误
		//140 已初始化
		//999 逻辑删除
		OrderStatus string `json:"orderStatus"`
	} `json:"returnObj"`
}
