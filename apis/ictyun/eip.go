package ictyun

type EipListRequest struct {
	PagingParam
	// 资源池ID
	RegionID string `json:"regionID"`
	// 客户端存根。用于保证订单幂等性, 长度 1 - 64。可以随机生成一个
	ClientToken string `json:"clientToken"`
	// eip类型 normal / cn2
	EipType *string `json:"eipType"`
	// 弹性 IP 的 ip 地址
	IP  *string  `json:"ip"`
	IDs []string `json:"ids"`
	// ip类型 ipv4 / ipv6
	IPType *string `json:"ipType"`
	// eip状态
	// ACTIVE（已绑定）/ DOWN（未绑定）/ FREEZING（已冻结）/ EXPIRED（已过期），不传是查询所有状态的 EIP
	Status *string `json:"status"`
}

type EipListResponse struct {
	ReturnObj struct {
		Eips []struct {
			ID         string `json:"ID"`
			Name       string `json:"name"`
			EipAddress string `json:"eipAddress"`
			// 当前绑定的实例的 ID
			AssociationID string `json:"associationID"`
			// 当前绑定的实例类型
			AssociationType string `json:"associationType"`
			// 交换机网段内的一个 IP 地址
			PrivateIPAddress string `json:"privateIpAddress"`
			// 带宽峰值大小，单位 Mb
			Bandwidth int    `json:"bandwidth"`
			Status    string `json:"status"`
			Tags      string `json:"tags"`
			CreatedAt string `json:"createdAt"`
			UpdatedAt string `json:"updatedAt"`
			ExpiredAt string `json:"expiredAt"`
		} `json:"eips"`
	} `json:"returnObj"`
	CurrentCount int `json:"currentCount"`
	TotalCount   int `json:"totalCount"`
	TotalPage    int `json:"totalPage"`
}

type EipCreateRequest struct {
	// 资源池ID
	RegionID string `json:"regionID"`
	// 客户端存根。用于保证订单幂等性, 长度 1 - 64。可以随机生成一个
	ClientToken string `json:"clientToken"`
	// 订购类型
	// month（包月/year（包年）/ on_demand（按需）
	CycleType string `json:"cycleType"`
	Name      string `json:"name"`
	// 订购时长
	// 当cycleType=month, 支持续订1-11个月;
	// 当cycleType=year, 支持续订1-3年,
	// 当cycleType = on_demand 时，可以不传
	CycleCount *int `json:"cycleCount"`
	// 弹性IP的带宽峰值，
	// 默认为1Mbps
	Bandwidth *int `json:"bandwidth"`
	// 当cycleType 为on_demand时，可以使用bandwidthID，将弹性IP加入到共享带宽中
	BandwidthID *string `json:"bandwidthID"`
	// 按需计费类型
	// 当cycleType为 on_demand时生效，支持bandwidth（按带宽/upflowc（按流量）
	DemandBillingType *string `json:"demandBillingType"`
	// 不填默认为默认企业项目，如果需要指定企业项目，则需要填写
	ProjectID *string `json:"projectID"`
}

type EipCreateResponse struct {
	ReturnObj struct {
		MasterOrderID        string      `json:"masterOrderID"`
		MasterOrderNO        interface{} `json:"masterOrderNO"`
		MasterResourceID     string      `json:"masterResourceID"`
		MasterResourceStatus string      `json:"masterResourceStatus"`
		RegionID             string      `json:"regionID"`
	} `json:"returnObj"`
}

type EipDeleteRequest struct {
	// 资源池ID
	RegionID string `json:"regionID"`
	// 客户端存根。用于保证订单幂等性, 长度 1 - 64。可以随机生成一个
	ClientToken string `json:"clientToken"`
	EipID       string `json:"eipID"`
	// 不填默认为默认企业项目，如果需要指定企业项目，则需要填写
	ProjectID *string `json:"projectID"`
}

type EipDeleteResponse struct {
	ReturnObj struct {
		MasterOrderID        string      `json:"masterOrderID"`
		MasterOrderNO        interface{} `json:"masterOrderNO"`
		MasterResourceID     string      `json:"masterResourceID"`
		MasterResourceStatus string      `json:"masterResourceStatus"`
		RegionID             string      `json:"regionID"`
	} `json:"returnObj"`
}
