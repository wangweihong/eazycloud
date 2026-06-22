package ictyun

type LogListRequest struct {
	PagingParam
	RegionID     string  `json:"regionID" `    //资源池ID
	Product      *string `json:"product"`      //查询的产品,默认compute,（compute主机，volume磁盘，network网络，monitor监控）
	ActionType   *string `json:"actionType" `  //操作动作类型
	Action       *string `json:"action"`       //操作动作
	DelegateType *string `json:"delegateType"` //委托类型，1表示获取所有用户操作,非1指获取用户自己的操作
	QueryContent *string `json:"queryContent"` //模糊查询关联资源id，或者资源name或者操作日志id
}

type LogListResponse struct {
	ReturnObj struct {
		CurrentCount int         `json:"currentCount"`
		TotalPage    int         `json:"totalPage"`
		PageNo       interface{} `json:"pageNo"`
		PageSize     interface{} `json:"pageSize"`
		DataList     []struct {
			Username     string `json:"username"`
			FinishTime   string `json:"finishTime"`
			RegionID     string `json:"regionID"`
			IsSubUser    bool   `json:"isSubUser"`
			ResourceType string `json:"resourceType"`
			SubUserID    string `json:"subUserID"`
			ExtraInfo    struct {
			} `json:"extraInfo"`
			Action       string `json:"action"`
			ActionType   string `json:"actionType"`
			StartTime    string `json:"startTime"`
			ResourceUUID string `json:"resourceUUID"`
			ResourceName string `json:"resourceName"`
			AccountID    string `json:"accountID"`
			RegionName   string `json:"regionName"`
		} `json:"dataList"`
		TotalCount int `json:"totalCount"`
	} `json:"returnObj"`
}
