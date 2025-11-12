package ictyun

type VpcListRequest struct {
	PagingParam
	RegionID string  `json:"regionID"`
	VpcID    *string `json:"vpcID"`
}

type VpcListResponse struct {
	ReturnObj struct {
		Vpcs []struct {
			VpcID          string   `json:"vpcID"`
			Name           string   `json:"name"`
			Description    string   `json:"description"`
			Cidr           string   `json:"CIDR"`
			Ipv6Enabled    bool     `json:"ipv6Enabled"`
			EnableIpv6     bool     `json:"enableIpv6"`
			Ipv6CIDRS      []string `json:"ipv6CIDRS"`
			SubnetIDs      []string `json:"subnetIDs"`
			NatGatewayIDs  []string `json:"natGatewayIDs"`
			SecondaryCIDRS []string `json:"secondaryCIDRS"`
			ProjectID      string   `json:"projectID"`
		} `json:"vpcs"`
		TotalCount   int `json:"totalCount"`
		TotalPage    int `json:"totalPage"`
		CurrentCount int `json:"currentCount"`
	} `json:"returnObj"`
}

type VpcCreateRequest struct {
	RegionID string `json:"regionID"`
	// 客户端存根。用于保证订单幂等性, 长度 1 - 64。可以随机生成一个
	ClientToken string `json:"clientToken"`
	// 支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32
	Name string `json:"name"`
	// VPC的网段。建议使用 192.168.0.0/16、172.16.0.0/12、10.0.0.0/
	CIDR string `json:"CIDR"`
	// 长度限制128
	Description *string `json:"description"`
	// 是否开启 IPv6 网段
	EnableIPV6 *bool   `json:"enableIpv6"`
	ProjectID  *string `json:"projectID"`
}

type VpcCreateResponse struct {
	ReturnObj struct {
		VpcID string `json:"vpcID"`
	} `json:"returnObj"`
}

type VpcDeleteRequest struct {
	RegionID  string  `json:"regionID"`
	VpcID     string  `json:"vpcID"`
	ProjectID *string `json:"projectID"`
}

type VpcDeleteResponse struct {
}

type VpcSubnetListRequest struct {
	PagingParam
	// 资源池
	RegionID string `json:"regionID"`
	// 所属vpc id
	VpcID string `json:"vpcID"`
	// 多个subnet用","隔开
	SubnetID *string `json:"subnetID"`
}

type VpcSubnetListResponse struct {
	ReturnObj struct {
		Subnets []struct {
			SubnetID          string   `json:"subnetID"`
			Name              string   `json:"name"`
			Description       string   `json:"description"`
			VpcID             string   `json:"vpcID"`
			Cidr              string   `json:"CIDR"`
			AvailableIPCount  int      `json:"availableIPCount"`
			GatewayIP         string   `json:"gatewayIP"`
			AvailabilityZones []string `json:"availabilityZones"`
			RouteTableID      string   `json:"routeTableID"`
			NetworkACLID      string   `json:"networkAclID"`
			Start             string   `json:"start"`
			End               string   `json:"end"`
			Ipv6Enabled       int      `json:"ipv6Enabled"`
			Ipv6CIDR          string   `json:"ipv6CIDR"`
			Ipv6Start         string   `json:"ipv6Start"`
			Ipv6End           string   `json:"ipv6End"`
			Ipv6GatewayIP     string   `json:"ipv6GatewayIP"`
			DNSList           []string `json:"dnsList"`
			NtpList           []string `json:"ntpList"`
			Type              int      `json:"type"`
			CreateAt          string   `json:"createAt"`
			UpdateAt          string   `json:"updateAt"`
			ProjectID         string
		} `json:"subnets"`
		TotalCount   int `json:"totalCount"`
		TotalPage    int `json:"totalPage"`
		CurrentCount int `json:"currentCount"`
	} `json:"returnObj"`
}

type VpcSubnetCreateRequest struct {
	RegionID string `json:"regionID"`
	// 客户端存根。用于保证订单幂等性, 长度 1 - 64。可以随机生成一个
	ClientToken string `json:"clientToken"`
	VpcID       string `json:"vpcID"`
	// 支持拉丁字母、中文、数字，下划线，连字符，中文 / 英文字母开头，不能以 http: / https: 开头，长度 2 - 32
	Name string `json:"name"`
	// 子网网段。建议使用 192.168.0.0/16、172.16.0.0/12、10.0.0.0/
	// 一个 Subnet 只能指定一个网段，创建后无法修改网段。
	CIDR string `json:"CIDR"`
	// 长度限制128
	Description *string `json:"description"`
	// 是否开启 IPv6 网段
	EnableIPV6 *bool `json:"enableIpv6"`
	// 子网 dns 列表, 最多同时支持 4 个 dns 地址
	DnsList []string `json:"dnsList"`
	// 子网类型：common（普通子网）/ cbm（裸金属子网），默认为普通子网
	SubnetType *string `json:"subnetType"`
	// 子网网关 IP
	SubnetGatewayIP *string `json:"subnetGatewayIP"`
}

type VpcSubnetCreateResponse struct {
	ReturnObj struct {
		SubnetID string
	} `json:"returnObj"`
}

type VpcSubnetDeleteRequest struct {
	RegionID string `json:"regionID"`
	SubnetID string `json:"subnetID"`
}

type VpcSubnetDeleteResponse struct {
}
