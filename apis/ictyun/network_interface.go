package ictyun

type NetworkInterfaceListRequest struct {
	PagingParam
	RegionID string `json:"regionID"` //资源池ID
	VpcID    string `json:"vpcID"`    //所属vpc id
	DeviceID string `json:"deviceID"` //关联设备id
	SubnetID string `json:"subnetID"` //所属子网id
}

type NetworkInterfaceListResponse struct {
	ReturnObj []struct {
		NetworkInterfaceName string      `json:"networkInterfaceName"`
		NetworkInterfaceID   string      `json:"networkInterfaceID"`
		VpcID                string      `json:"vpcID"`
		SubnetID             string      `json:"subnetID"`
		Role                 int         `json:"role"`
		MacAddress           string      `json:"macAddress"`
		PrimaryPrivateIP     string      `json:"primaryPrivateIp"`
		Ipv6Addresses        []string    `json:"ipv6Addresses"`
		InstanceID           string      `json:"instanceID"`
		InstanceType         string      `json:"instanceType"`
		InstanceOwnerID      string      `json:"instanceOwnerID"`
		Description          string      `json:"description"`
		SecurityGroupIds     []string    `json:"securityGroupIds"`
		SecondaryPrivateIps  []string    `json:"secondaryPrivateIps"`
		AdminStatus          string      `json:"adminStatus"`
		AssociatedEip        interface{} `json:"associatedEip"`
	} `json:"returnObj"`
	CurrentCount int `json:"currentCount"`
	TotalCount   int `json:"totalCount"`
	TotalPage    int `json:"totalPage"`
}

type NetworkInterfaceGetRequest struct {
	RegionID           string `json:"regionID"`           //资源池ID
	NetworkInterfaceID string `json:"NetworkInterfaceID"` //弹性网卡id
}

type NetworkInterfaceGetResponse struct {
	ReturnObj struct {
		NetworkInterfaceID string `json:"networkInterfaceID"`
		NetworkInterName   string `json:"networkInterName"`
		Description        string `json:"description"`
		Role               string `json:"role"`
		VpcID              string `json:"vpcID"`
		SubnetID           string `json:"subnetID"`
		MacAddress         string `json:"macAddress"`
		SanityCheck        string `json:"sanityCheck"`
		PrimaryPrivateIP   string `json:"primaryPrivateIp"`
		DeviceID           string `json:"deviceId"`
		DeviceType         string `json:"deviceType"`
	} `json:"returnObj"`
}

type NetworkInterfaceCreateRequest struct {
	RegionID                string   `json:"regionID"`                //资源池ID
	ClientToken             string   `json:"clientToken"`             //客户端存根
	SubnetID                string   `json:"subnetID"`                //子网ID
	PrimaryPrivateIp        *string  `json:"primaryPrivateIp"`        //弹性网卡的主私网IPv4地址
	Ipv6Addresses           []string `json:"ipv6Addresses"`           //弹性网卡的主私网IPv6地址
	SecurityGroupIds        []string `json:"securityGroupIds"`        //加入一个或多个安全组
	SecondaryPrivateIpCount *string  `json:"secondaryPrivateIpCount"` //辅助私网IP地址数量
	SecondaryPrivateIps     []string `json:"secondaryPrivateIps"`     //辅助私网IP地址
	Name                    *string  `json:"name"`                    //网卡名称
	Description             *string  `json:"description"`             //网卡的描述
}

type NetworkInterfaceCreateResponse struct {
	ReturnObj struct {
		VpcID                string   `json:"vpcID"`
		SubnetID             string   `json:"subnetID"`
		NetworkInterfaceName string   `json:"networkInterfaceName"`
		MacAddress           string   `json:"macAddress"`
		NetworkInterfaceID   string   `json:"networkInterfaceID"`
		Description          string   `json:"description"`
		Ipv6Address          []string `json:"ipv6Address"`
		SecurityGroupIds     []string `json:"securityGroupIds"`
		SecondaryPrivateIps  []string `json:"secondaryPrivateIps"`
		PrivateIPAddress     string   `json:"privateIpAddress"`
		InstanceOwnerID      string   `json:"instanceOwnerID"`
		InstanceType         string   `json:"instanceType"`
		InstanceID           string   `json:"instanceID"`
	} `json:"returnObj"`
}

type NetworkInterfaceDeleteRequest struct {
	RegionID           string `json:"regionID"`           //资源池ID
	ClientToken        string `json:"clientToken"`        //客户端存根
	NetworkInterfaceID string `json:"networkInterfaceID"` //网卡ID
}

type NetworkInterfaceDeleteResponse struct {
}

type NetworkInterfaceAttachRequest struct {
	RegionID           string  `json:"regionID"`           //资源池ID
	ClientToken        string  `json:"clientToken"`        //客户端存根
	AzName             string  `json:"azName"`             //可用区名称
	ProjectID          *string `json:"projectID"`          //企业项目ID
	NetworkInterfaceID string  `json:"networkInterfaceID"` //弹性网卡ID
	InstanceID         string  `json:"instanceID"`         //云主机ID
	InstanceType       *int    `json:"instanceType"`       //实例类型：3-虚拟机，4-物理机
}

type NetworkInterfaceAttachResponse struct {
}

type NetworkInterfaceDetachRequest struct {
	RegionID           string `json:"regionID"`           //资源池ID
	NetworkInterfaceID string `json:"networkInterfaceID"` //网卡ID
	ClientToken        string `json:"clientToken"`        //客户端存根
}

type NetworkInterfaceDetachResponse struct {
}

type NetworkInterfaceUpdateRequest struct {
	RegionID           string   `json:"regionID"`           //资源池ID
	ClientToken        string   `json:"clientToken"`        //客户端存根
	NetworkInterfaceID string   `json:"networkInterfaceID"` //网卡ID
	SecurityGroupIds   []string `json:"securityGroupIds"`   //加入一个或多个安全组
	Name               *string  `json:"name"`               //网卡名称
	Description        *string  `json:"description"`        //网卡的描述
}

type NetworkInterfaceUpdateResponse struct {
}
