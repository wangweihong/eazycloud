package ictyun

type KeyPairCreateRequest struct {
	RegionID    string  `json:"regionID"`    //资源池ID
	KeyPairName string  `json:"keyPairName"` //密钥对名称
	ProjectID   *string `json:"projectID"`   //企业项目ID
}

type KeyPairCreateResponse struct {
	ReturnObj struct {
		PublicKey   string `json:"publicKey"`
		PrivateKey  string `json:"privateKey"`
		KeyPairName string `json:"keyPairName"`
		FingerPrint string `json:"fingerPrint"`
		KeyPairID   string `json:"keyPairID"`
	} `json:"returnObj"`
}

type KeyPairAttachInstanceRequest struct {
	RegionID    string `json:"regionID"`    //资源池ID
	KeyPairName string `json:"keyPairName"` //密钥对名称
	InstanceID  string `json:"instanceID"`  //云服务器ID
}

type KeyPairAttachInstanceResponse struct {
	ReturnObj struct {
		InstanceID string `json:"instanceID"`
	} `json:"returnObj"`
}

type KeyPairDetachInstanceRequest struct {
	RegionID    string `json:"regionID"`    //资源池ID
	KeyPairName string `json:"keyPairName"` //密钥对名称
	InstanceID  string `json:"instanceID"`  //云服务器ID
}

type KeyPairDetachInstanceResponse struct {
	ReturnObj struct {
		InstanceID string `json:"instanceID"`
	} `json:"returnObj"`
}

type KeyPairListRequest struct {
	PagingParam
	RegionID     string  `json:"regionID"`     //资源池ID
	ProjectID    *string `json:"projectID"`    //项目ID
	KeyPairName  *string `json:"keyPairName"`  //密钥对名称
	QueryContent *string `json:"queryContent"` //模糊匹配查询内容（匹配字段：keyPairName、keyPairID）
}

type KeyPairListResponse struct {
	ReturnObj struct {
		CurrentCount int `json:"currentCount"`
		TotalCount   int `json:"totalCount"`
		Results      []struct {
			PublicKey   string `json:"publicKey"`
			FingerPrint string `json:"fingerPrint"`
			KeyPairName string `json:"keyPairName"`
			KeyPairID   string `json:"keyPairID"`
			ProjectID   string `json:"projectID"`
		} `json:"results"`
	} `json:"returnObj"`
}

type KeyPairImportRequest struct {
	RegionID    string  `json:"regionID"`     //资源池ID
	KeyPairName string  `json:"snapshotName"` //密钥对名称
	PublicKey   string  `json:"publicKey"`    //导入的公钥信息
	ProjectID   *string `json:"projectID"`    //企业项目ID
}

type KeyPairImportResponse struct {
	ReturnObj struct {
		PublicKey   string `json:"publicKey"`
		KeyPairName string `json:"keyPairName"`
		FingerPrint string `json:"fingerPrint"`
	} `json:"returnObj"`
}

type KeyPairDeleteRequest struct {
	RegionID    string `json:"regionID"`    //资源池ID
	KeyPairName string `json:"keyPairName"` //密钥对名称
}

type KeyPairDeleteResponse struct {
	ReturnObj struct {
		KeyPairName string `json:"keyPairName"`
	} `json:"returnObj"`
}
