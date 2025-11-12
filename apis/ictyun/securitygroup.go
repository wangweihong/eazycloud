package ictyun

type SecurityGroupDeleteRequest struct {
	RegionID        string  `json:"regionID"`        //资源池ID
	SecurityGroupID string  `json:"securityGroupID"` //安全组ID
	ProjectID       *string `json:"projectID"`       //企业项目ID
	ClientToken     string  `json:"clientToken"`     //客户端存根
}

type SecurityGroupDeleteResponse struct {
}

type SecurityGroupRule struct {
	Direction   string  `json:"direction"`   //规则方向，入方向则填写ingress
	Action      string  `json:"action"`      //授权策略，取值范围：accept（允许），drop（拒绝）
	Priority    *int    `json:"priority"`    //优先级，取值范围：[1, 100]
	Protocol    string  `json:"protocol"`    //网络协议，取值范围：ANY（任意）、TCP、UDP、ICMP
	Ethertype   string  `json:"ethertype"`   //IP类型，取值范围：IPv4、IPv6
	DestCidrIp  string  `json:"destCidrIp"`  //远端地址
	Description *string `json:"description"` //描述
	Range       string  `json:"range"`       //安全组开放的传输层协议相关的源端端口范围
}

type SecurityGroupIngressCreateRequest struct {
	RegionID           string              `json:"regionID"`           //资源池ID
	SecurityGroupID    string              `json:"securityGroupID"`    //安全组ID
	SecurityGroupRules []SecurityGroupRule `json:"securityGroupRules"` //规则信息列表
	ClientToken        string              `json:"clientToken"`        //客户端存根
}

type SecurityGroupIngressCreateResponse struct {
}

type SecurityGroupIngressDeleteRequest struct {
	RegionID            string `json:"regionID"`            //资源池ID
	SecurityGroupID     string `json:"securityGroupID"`     //安全组ID
	SecurityGroupRuleID string `json:"securityGroupRuleID"` //安全组入向规则ID
	ClientToken         string `json:"clientToken"`         //客户端存根
}

type SecurityGroupIngressDeleteResponse struct {
}

type SecurityGroupIngressUpdateRequest struct {
	RegionID            string `json:"regionID"`            //资源池ID
	SecurityGroupID     string `json:"securityGroupID"`     //安全组ID
	SecurityGroupRuleID string `json:"securityGroupRuleID"` //安全组入向规则ID
	Description         string `json:"description"`         //安全组规则描述信息
	ClientToken         string `json:"clientToken"`         //客户端存根
}

type SecurityGroupIngressUpdateResponse struct {
}

type SecurityGroupEgressDeleteRequest struct {
	RegionID            string `json:"regionID"`            //资源池ID
	SecurityGroupID     string `json:"securityGroupID"`     //安全组ID
	SecurityGroupRuleID string `json:"securityGroupRuleID"` //安全组入向规则ID
	ClientToken         string `json:"clientToken"`         //客户端存根
}

type SecurityGroupEgressDeleteResponse struct {
}

type SecurityGroupEgressCreateRequest struct {
	RegionID           string              `json:"regionID"`           //资源池ID
	SecurityGroupID    string              `json:"securityGroupID"`    //安全组ID
	SecurityGroupRules []SecurityGroupRule `json:"securityGroupRules"` //规则信息列表
	ClientToken        string              `json:"clientToken"`        //客户端存根
}

type SecurityGroupEgressCreateResponse struct {
}

type SecurityGroupEgressUpdateRequest struct {
	RegionID            string  `json:"regionID"`            //资源池ID
	SecurityGroupID     string  `json:"securityGroupID"`     //安全组ID
	SecurityGroupRuleID string  `json:"securityGroupRuleID"` //安全组入向规则ID
	Description         *string `json:"description"`         //安全组规则描述信息
	ClientToken         string  `json:"clientToken"`         //客户端存根
}

type SecurityGroupEgressUpdateResponse struct {
}

type SecurityGroupListRequest struct {
	PagingParam
	RegionID     string  `json:"regionID"`     //资源池ID
	VpcID        *string `json:"vpcID"`        //虚拟私有云ID
	QueryContent *string `json:"queryContent"` //模糊匹配查询内容（匹配字段：id、name）
	ProjectID    *string `json:"projectID"`    //企业项目ID
	InstanceID   *string `json:"instanceID"`   //云主机ID
}

type SecurityGroupListResponse struct {
	ReturnObj []struct {
		SecurityGroupName     string `json:"securityGroupName"`
		ID                    string `json:"id"`
		VMNum                 int    `json:"vmNum"`
		Origin                string `json:"origin"`
		VpcName               string `json:"vpcName"`
		VpcID                 string `json:"vpcID"`
		CreationTime          string `json:"creationTime"`
		Description           string `json:"description"`
		SecurityGroupRuleList []struct {
			Direction       string `json:"direction"`
			Priority        int    `json:"priority"`
			Ethertype       string `json:"ethertype"`
			Protocol        string `json:"protocol"`
			Range           string `json:"range"`
			DestCidrIP      string `json:"destCidrIp"`
			Description     string `json:"description"`
			Origin          string `json:"origin"`
			CreateTime      string `json:"createTime"`
			ID              string `json:"id"`
			Action          string `json:"action"`
			SecurityGroupID string `json:"securityGroupID"`
		} `json:"securityGroupRuleList"`
		ProjectID string `json:"projectID"`
	} `json:"returnObj"`
	CurrentCount int `json:"currentCount"`
	TotalCount   int `json:"totalCount"`
	TotalPage    int `json:"totalPage"`
}

type SecurityGroupJoinRequest struct {
	RegionID           string  `json:"regionID"`           //资源池ID
	SecurityGroupID    string  `json:"securityGroupID"`    //安全组ID
	InstanceID         string  `json:"instanceID"`         //云主机ID
	Action             string  `json:"action"`             //系统规定参数，绑定安全组填写joinSecurityGroup
	NetworkInterfaceID *string `json:"networkInterfaceID"` //弹性网卡ID
}

type SecurityGroupJoinResponse struct {
}

type SecurityGroupCreateRequest struct {
	RegionID        string  `json:"regionID"`        //资源池ID
	SecurityGroupID string  `json:"securityGroupID"` //安全组ID
	ClientToken     string  `json:"clientToken"`     //客户端存根
	ProjectID       *string `json:"projectID"`       //企业项目ID
	VpcID           string  `json:"vpcID"`           //虚拟私有云ID
	Name            string  `json:"name"`            //安全组名称
	Description     *string `json:"description"`     //安全组描述信息
}

type SecurityGroupCreateResponse struct {
	ReturnObj struct {
		SecurityGroupID string `json:"securityGroupID"`
	} `json:"returnObj"`
}

type SecurityGroupLeaveRequest struct {
	RegionID        string `json:"regionID"`        //资源池ID
	SecurityGroupID string `json:"securityGroupID"` //安全组ID
	InstanceID      string `json:"instanceID"`      //云主机ID
}

type SecurityGroupLeaveResponse struct{}

type SecurityGroupGetRequest struct {
	RegionID        string `json:"regionID"`        //资源池ID
	SecurityGroupID string `json:"securityGroupID"` //安全组ID
	Direction       string `json:"direction"`       //安全组规则授权方向，取值范围：egress（出方向），ingress（入方向），all（不区分方向
}

type SecurityGroupGetResponse struct {
	ReturnObj struct {
		SecurityGroupName     string `json:"securityGroupName"`
		ID                    string `json:"id"`
		VMNum                 int    `json:"vmNum"`
		Origin                string `json:"origin"`
		VpcName               string `json:"vpcName"`
		VpcID                 string `json:"vpcID"`
		CreationTime          string `json:"creationTime"`
		Description           string `json:"description"`
		SecurityGroupRuleList []struct {
			Direction       string `json:"direction"`
			Priority        int    `json:"priority"`
			Ethertype       string `json:"ethertype"`
			Protocol        string `json:"protocol"`
			Range           string `json:"range"`
			DestCidrIP      string `json:"destCidrIp"`
			Description     string `json:"description"`
			Origin          string `json:"origin"`
			CreateTime      string `json:"createTime"`
			ID              string `json:"id"`
			Action          string `json:"action"`
			SecurityGroupID string `json:"securityGroupID"`
		} `json:"securityGroupRuleList"`
	} `json:"returnObj"`
}
