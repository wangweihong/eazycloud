package ictyun

type MonitorEcsListRequest struct {
	PagingParam
	RegionID string `json:"regionID"` //资源池ID
}

type MonitorEcsListResponse struct {
	ReturnObj struct {
		Page         int `json:"page"`
		PageSize     int `json:"pageSize"`
		TotalPage    int `json:"totalPage"`
		CurrentCount int `json:"currentCount"`
		TotalCount   int `json:"totalCount"`
		List         []struct {
			DeviceUUID   string `json:"deviceUUID"`
			InstanceID   string `json:"instanceID"`
			InstanceName string `json:"instanceName"`
		} `json:"list"`
	} `json:"returnObj"`
}

type MonitorDiskListRequest struct {
	PagingParam
	RegionID string `json:"regionID"` //资源池ID
}

type MonitorDiskListResponse struct {
	ReturnObj struct {
		Page         int `json:"page"`
		PageSize     int `json:"pageSize"`
		TotalCount   int `json:"totalCount"`
		CurrentCount int `json:"currentCount"`
		TotalPage    int `json:"totalPage"`
		List         []struct {
			DeviceUUID   string `json:"deviceUUID"`
			InstanceID   string `json:"instanceID"`
			InstanceName string `json:"instanceName"`
		} `json:"list"`
	} `json:"returnObj"`
}

/*
	本参数表示设备类型。默认值为所有类型。取值范围：

vm：云主机。
bare_metal：裸金属。
disk：云磁盘。
scaling：弹性伸缩。
traffic：共享带宽。
eip：弹性IP。
elb：负载均衡。
listener：监听器。
cstor_sfs：弹性文件。
site_monitor：站点监控。
natgw：NAT网关。
zos_bucket：对象存储-存储桶。
zos_user：对象存储-用户。
vnet_monitor_endpoint_statistic：VPC终端节点。
vnet_endpoint_service_statistic：VPC终端节点服务。
根据以上范围取值。
*/
type MonitorItemListRequest struct {
	DeviceType *string `json:"deviceType"` //设备类型
}

type MonitorItemListResponse struct {
	ReturnObj struct {
		MonitorItems []struct {
			DeviceType string `json:"deviceType"`
			Items      []struct {
				Name string `json:"name"`
				Desc string `json:"desc"`
				Unit string `json:"unit"`
			} `json:"items"`
		} `json:"monitorItems"`
	} `json:"returnObj"`
}

type MonitorHistoryListRequest struct {
	RegionID       string   `json:"regionID,omitempty" description:"资源池ID" required:"true"`
	ItemNameList   []string `json:"itemNameList,omitempty" description:"待查的监控项名称" required:"true"`
	StartTime      string   `json:"startTime,omitempty" description:"查询起始时间戳" required:"true"`
	EndTime        string   `json:"endTime,omitempty" description:"查询结束时间戳" required:"true"`
	DeviceUUIDList []string `json:"deviceUUIDList,omitempty" description:"查询设备ID列表" required:"true"`
}

type MonitorHistoryListResponse struct {
	ReturnObj struct {
		Result []struct {
			RegionID          string `json:"regionID"`
			DeviceUUID        string `json:"deviceUUID"`
			ItemAggregateList []struct {
				ItemName string `json:"itemName"`
				ItemData []struct {
					Value        string `json:"value"`
					SamplingTime int    `json:"samplingTime"`
				} `json:"itemData"`
			} `json:"itemAggregateList"`
		} `json:"result"`
	} `json:"returnObj"`
}
