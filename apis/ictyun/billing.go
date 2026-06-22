package ictyun

type BillCycleFeeListRequest struct {
	PagingParam
	ProductCode   *string `json:"productCode"`    //产品编码
	BllingCycleId string  `json:"billingCycleId"` //账期如202312
	BillType      *string `json:"billType"`       //账单类型
	ContractId    *string `json:"contractId"`     //合同标识
	ProjectId     *string `json:"projectId"`      //项目id
	MasterOrderId *string `json:"masterOrderId"`  //主订单标识
}

type BillCycleFeeListResponse struct {
	ReturnObj struct {
		Result []struct {
			OrderType       string `json:"orderType"`       //订单类型
			ResourceID      string `json:"resourceId"`      //资源id
			OrderID         string `json:"orderId"`         //订单标识
			BillMode        string `json:"billMode"`        //计费模式
			DiscountAmount  string `json:"discountAmount"`  //优惠金额
			ConsumeDate     string `json:"consumeDate"`     //消费时间
			PayableAmount   string `json:"payableAmount"`   //应付金额
			ProductName     string `json:"productName"`     //主资源产品名称
			ServID          string `json:"servId"`          //资源实例id
			PayMethod       string `json:"payMethod"`       //付款方式
			Price           string `json:"price"`           //官网价（原价）
			ServiceTag      string `json:"serviceTag"`      //服务标识
			ContractName    string `json:"contractName"`    //合同名称
			Cash            string `json:"cash"`            //余额支付
			IsAgencyAdvance string `json:"isAgencyAdvance"` //
			Amount          string `json:"amount"`          //金额
			OrderNo         string `json:"orderNo"`         //订单编码
			Coupon          string `json:"coupon"`          //代金券抵扣
			ChannleAmount   string `json:"channleAmount"`   //现金金额
			BillType        string `json:"billType"`        //账单类型
			ChannelAmount   string `json:"channelAmount"`   //
			ProductCode     string `json:"productCode"`     //产品名称
			RegionID        string `json:"regionId"`        //资源池名称，多个用|隔开
			BillingCycle    string `json:"billingCycle"`    //账期
			ContractID      string `json:"contractId"`      //合同id
			ProjectName     string `json:"projectName"`     //项目名称
			ContractCode    string `json:"contractCode"`    //合同编码
			PayStatus       string `json:"payStatus"`       //支付状态
			ProjectID       string `json:"projectId"`       //项目标识
			ResourceType    string `json:"resourceType"`    //资源类型
		} `json:"result"`
		AccountID      string `json:"accountId"`
		BillingCycleID string `json:"billingCycleId"`
		PageNo         int    `json:"pageNo"`
		PageSize       int64  `json:"pageSize"`
		TotalCount     int    `json:"totalCount"`
	} `json:"returnObj"`
}

type BillOnDemandFeeListRequest struct {
	PagingParam
	ProductCode    *string `json:"productCode"`    //产品编码
	BillingCycleId string  `json:"billingCycleId"` //账期如202312
	HasTotal       *bool   `json:"hasTotal"`       //是否包含汇总
}

type BillOnDemandFeeListResponse struct {
	ReturnObj struct {
		Result []struct {
			OrderType       string `json:"orderType"`       //订单类型
			ResourceID      string `json:"resourceId"`      //资源实例
			OrderID         string `json:"orderId"`         //订单ID
			BillMode        string `json:"billMode"`        //计费模式 1：包周期 2：按需
			DiscountAmount  string `json:"discountAmount"`  //优惠金额
			ConsumeDate     string `json:"consumeDate"`     //消费时间
			PayableAmount   string `json:"payableAmount"`   //应付金额
			ProductName     string `json:"productName"`     //资源产品名称
			ServID          string `json:"servId"`          //资源实例ID
			PayMethod       string `json:"payMethod"`       //支付方式
			Price           string `json:"price"`           //官网价（原价）
			ServiceTag      string `json:"serviceTag"`      //服务类型
			ContractName    string `json:"contractName"`    //合同名称
			Cash            string `json:"cash"`            //余额支付
			IsAgencyAdvance string `json:"isAgencyAdvance"` //
			Amount          string `json:"amount"`          //金额
			OrderNo         string `json:"orderNo"`         //订单NO
			Coupon          string `json:"coupon"`          //代金券抵扣
			ChannleAmount   string `json:"channleAmount"`   //现金支付
			BillType        string `json:"billType"`        //账单类型
			ChannelAmount   string `json:"channelAmount"`   //
			ProductCode     string `json:"productCode"`     //资源产品代码
			RegionID        string `json:"regionId"`        //资源池名称
			BillingOnDemand string `json:"billingOnDemand"` //
			ContractID      string `json:"contractId"`      //
			ProjectName     string `json:"projectName"`     //项目名称
			ContractCode    string `json:"contractCode"`    //合同编码
			PayStatus       string `json:"payStatus"`       //支付状态
			ProjectID       string `json:"projectId"`       //projectId
			ResourceType    string `json:"resourceType"`    //资源类型
		} `json:"result"`
		AccountID         string `json:"accountId"`
		BillingOnDemandID string `json:"billingCycleId"`
		PageNo            int    `json:"pageNo"`
		PageSize          int64  `json:"pageSize"`
		TotalCount        int    `json:"totalCount"`
	} `json:"returnObj"`
}

type BillCycleBillDetailProdCycleIdListRequest struct {
	PagingParam
	ProductCode   string  `json:"productCode"`    //产品编码
	BllingCycleId *string `json:"billingCycleId"` //账期如202312
	BillType      string  `json:"billType"`       //账单类型
	ContractId    string  `json:"contractId"`     //合同标识
	ProjectId     string  `json:"projectId"`      //项目id
	MasterOrderId string  `json:"masterOrderId"`  //主订单标识
}

type BillCycleBillDetailProdCycleIdListResponse struct {
	ReturnObj struct {
		Result []struct {
			Amount           string `json:"amount"`           //金额
			Coupon           string `json:"coupon"`           //代金券抵扣
			BillType         string `json:"billType"`         //账单类型
			BillMode         string `json:"billMode"`         //计费模式 1：包周期 2：按需
			DiscountAmount   string `json:"discountAmount"`   //优惠金额
			ConsumeDate      string `json:"consumeDate"`      //消费时间
			PayableAmount    string `json:"payableAmount"`    //应付金额
			ProductName      string `json:"productName"`      //资源产品名称
			ProductCode      string `json:"productCode"`      //
			BillingCycleID   string `json:"billingCycleId"`   //
			PayMethod        string `json:"payMethod"`        //支付方式
			Price            string `json:"price"`            //官网价（原价）
			ServiceTag       string `json:"serviceTag"`       //服务类型
			ServiceTagName   string `json:"serviceTagName"`   //
			ContractName     string `json:"contractName"`     //合同名称
			ContractCode     string `json:"contractCode"`     //合同编码
			ResourceTypeName string `json:"resourceTypeName"` //
			ResourceType     string `json:"resourceType"`     //资源类型
		} `json:"result"`
		PageNo     int `json:"pageNo"`
		PageSize   int `json:"pageSize"`
		TotalCount int `json:"totalCount"`
	} `json:"returnObj"`
}

type BillOnDemandBillDetailResDetailListRequest struct {
	PagingParam
	ProductCode    *string `json:"productCode"`    //产品编码
	BillingCycleId string  `json:"billingCycleId"` //账期如202312
	BillType       *string `json:"billType"`       //账单类型
	ContractId     *string `json:"contractId"`     //合同标识
	ProjectId      *string `json:"projectId"`      //项目id
	HasTotal       *string `json:"hasTotal"`       //是否包括汇总
}

type BillOnDemandBillDetailResDetailListResponse struct {
	ReturnObj struct {
		Result []struct {
			ResourceID        string `json:"resourceId"`        //资源ID
			LabelInfo         string `json:"labelInfo"`         //资源标签
			MasterOrderNo     string `json:"masterOrderNo"`     //主订单编号
			BillMode          string `json:"billMode"`          //计费模式 1：包周期 2：按需
			DiscountAmount    string `json:"discountAmount"`    //优惠金额
			ConsumeDate       string `json:"consumeDate"`       //消费时间
			PayableAmount     string `json:"payableAmount"`     //应付金额
			ProductName       string `json:"productName"`       //产品名称
			RegionCode        string `json:"regionCode"`        //资源池编码
			ServID            string `json:"servId"`            //资源实例ID
			Price             string `json:"price"`             //官网价
			ServiceTag        string `json:"serviceTag"`        //服务标识
			ContractName      string `json:"contractName"`      //
			RealResourceID    string `json:"realResourceId"`    //真实资源ID
			SalesAttribute    string `json:"salesAttribute"`    //产销品规格
			KeySalesAttribute string `json:"keySalesAttribute"` //关键产销品规格
			Amount            string `json:"amount"`            //金额
			OfferName         string `json:"offerName"`         //
			Coupon            string `json:"coupon"`            //代金券金额
			BillType          string `json:"billType"`          //账单类型
			MasterOrderID     string `json:"masterOrderId"`     //主订单ID
			ProductCode       string `json:"productCode"`       //
			BillingCycleID    string `json:"billingCycleId"`    //
			RegionID          string `json:"regionId"`          //资源池ID
			ProjectName       string `json:"projectName"`       //项目名称
			ContractCode      string `json:"contractCode"`      //
			ProjectID         string `json:"projectId"`         //项目ID
			ResourceType      string `json:"resourceType"`      //资源类型
		} `json:"result"`
		PageNo     int `json:"pageNo"`
		PageSize   int `json:"pageSize"`
		TotalCount int `json:"totalCount"`
	} `json:"returnObj"`
}

type BillOnDemandBillDetailProductDetailListRequest struct {
	PagingParam
	ProductCode    *string `json:"productCode" `   //产品编码
	BillingCycleId string  `json:"billingCycleId"` //账期如202312
	BillType       *string `json:"billType" `      //账单类型
	ContractId     *string `json:"contractId"`     //合同标识
	ProjectId      *string `json:"projectId"`      //项目id
	HasTotal       *bool   `json:"hasTotal"`       //是否包括汇总
}

type BillOnDemandBillDetailProductDetailListResponse struct {
	ReturnObj struct {
		Result []struct {
			Amount         string `json:"amount"`         //金额
			Coupon         string `json:"coupon"`         //代金券抵扣
			BillType       string `json:"billType"`       //账单类型
			BillMode       string `json:"billMode"`       //计费模式 1：包周期 2：按需
			DiscountAmount string `json:"discountAmount"` //优惠金额
			ConsumeDate    string `json:"consumeDate"`    //消费时间
			PayableAmount  string `json:"payableAmount"`  //应付金额
			ProductName    string `json:"productName"`    //资源产品名称
			ProductCode    string `json:"productCode"`    //资源产品代码
			BillingCycleID string `json:"billingCycleId"` //
			PayMethod      string `json:"payMethod"`      //支付方式
			Price          string `json:"price"`          //官网价（原价）
			ServiceTag     string `json:"serviceTag"`     //服务类型
			ContractName   string `json:"contractName"`   //合同名称
			ContractCode   string `json:"contractCode"`   //合同编码
			ResourceType   string `json:"resourceType"`   //资源类型
		} `json:"result"`
		PageNo     int `json:"pageNo"`
		PageSize   int `json:"pageSize"`
		TotalCount int `json:"totalCount"`
	} `json:"returnObj"`
}

type BillOnDemandBillDetailResCycleIdListRequest struct {
	PagingParam
	ProductCode    *string `json:"productCode"`    //产品编码
	BillingCycleId string  `json:"billingCycleId"` //账期如202312
	BillType       *string `json:"billType"`       //账单类型
	ContractId     *string `json:"contractId"`     //合同标识
	ProjectId      *string `json:"projectId"`      //项目id
	HasTotal       *string `json:"hasTotal"`       //是否包括汇总
}

type BillOnDemandBillDetailResCycleIdListResponse struct {
	ReturnObj struct {
		Result []struct {
			ResourceID        string `json:"resourceId"`        //资源Id
			LabelInfo         string `json:"labelInfo"`         //资源标签
			BillMode          string `json:"billMode"`          //计费模式 1：包周期 2：按需
			DiscountAmount    string `json:"discountAmount"`    //优惠金额
			ConsumeDate       string `json:"consumeDate"`       //消费时间
			PayableAmount     string `json:"payableAmount"`     //应付金额
			ProductName       string `json:"productName"`       //项目名称
			RegionCode        string `json:"regionCode"`        //区域编码
			ServID            string `json:"servId"`            //资源实例id
			Price             string `json:"price"`             //官网价（原价）
			ServiceTag        string `json:"serviceTag"`        //服务内类型
			ContractName      string `json:"contractName"`      //合同名称
			RealResourceID    string `json:"realResourceId"`    //真实资源ID
			SalesAttribute    string `json:"salesAttribute"`    //销售品规格
			KeySalesAttribute string `json:"keySalesAttribute"` //关键规格
			Amount            string `json:"amount"`            //金额
			OfferName         string `json:"offerName"`         //销售品名称
			Coupon            string `json:"coupon"`            //代金券抵扣
			BillType          string `json:"billType"`          //计费模式 1：包周期 2：按需
			ProductCode       string `json:"productCode"`       //资源产品代码
			BillingCycleID    string `json:"billingCycleId"`    //账期
			RegionID          string `json:"regionId"`          //区域id
			ProjectName       string `json:"projectName"`       //项目
			ContractCode      string `json:"contractCode"`      //合同编码
			ResourceType      string `json:"resourceType"`      //资源类型
		} `json:"result"`
		PageNo     int `json:"pageNo"`
		PageSize   int `json:"pageSize"`
		TotalCount int `json:"totalCount"`
	} `json:"returnObj"`
}

type BillOnDemandBillDetailUsageCycleIdListRequest struct {
	PagingParam
	ProductCode    *string `json:"productCode"`    //产品编码
	BillingCycleId string  `json:"billingCycleId"` //账期如202312
	BillType       *string `json:"billType"`       //账单类型
	ContractId     *string `json:"contractId"`     //合同标识
	ProjectId      *string `json:"projectId"`      //项目id
	HasTotal       *string `json:"hasTotal"`       //是否包括汇总
}

type BillOnDemandBillDetailUsageCycleIdListResponse struct {
	ReturnObj struct {
		Result []struct {
			ResourceID        string `json:"resourceId"`        //资源实例ID
			LabelInfo         string `json:"labelInfo"`         //资源标签
			Usage             string `json:"usage"`             //使用量
			BillMode          string `json:"billMode"`          //计费模式 1：包周期 2：按需
			DiscountAmount    string `json:"discountAmount"`    //优惠金额
			ConsumeDate       string `json:"consumeDate"`       //消费时间
			DeductUsage       string `json:"deductUsage"`       //
			PayableAmount     string `json:"payableAmount"`     //应付金额
			PricefactorValue  string `json:"pricefactorValue"`  //
			ProductName       string `json:"productName"`       //资源产品名称
			RegionCode        string `json:"regionCode"`        //区域编码
			ServID            string `json:"servId"`            //资源实例id
			Price             string `json:"price"`             //官网价（原价）
			ServiceTag        string `json:"serviceTag"`        //服务类型
			ContractName      string `json:"contractName"`      //
			RealResourceID    string `json:"realResourceId"`    //真实资源ID
			SalesAttribute    string `json:"salesAttribute"`    //销售品规格
			KeySalesAttribute string `json:"keySalesAttribute"` //关键规格
			UsageType         string `json:"usageType"`         //使用量类型
			Amount            string `json:"amount"`            //金额
			UsageTypeID       string `json:"usageTypeId"`       //使用量类型ID
			OfferName         string `json:"offerName"`         //销售品名称
			Coupon            string `json:"coupon"`            //代金券抵扣
			BillType          string `json:"billType"`          //账单类型
			ProductCode       string `json:"productCode"`       //资源产品代码
			BillingCycleID    string `json:"billingCycleId"`    //账期
			RegionID          string `json:"regionId"`          //区域ID
			ProjectName       string `json:"projectName"`       //项目名称
			ContractCode      string `json:"contractCode"`      //合同编码
			ProjectID         string `json:"projectId"`         //项目id
			ResourceType      string `json:"resourceType"`      //资源类型
		} `json:"result"`
		PageNo     int `json:"pageNo"`
		PageSize   int `json:"pageSize"`
		TotalCount int `json:"totalCount"`
	} `json:"returnObj"`
}

type BillOnDemandBillDetailUsageDetailListRequest struct {
	PagingParam
	ProductCode    *string `json:"productCode"`    //产品编码
	BillingCycleId string  `json:"billingCycleId"` //账期如202312
	BillType       *string `json:"billType"`       //账单类型
	ContractId     *string `json:"contractId"`     //合同标识
	ProjectId      *string `json:"projectId"`      //项目id
	HasTotal       *string `json:"hasTotal"`       //是否包括汇总
}

type BillOnDemandBillDetailUsageDetailListResponse struct {
	ReturnObj struct {
		Result []struct {
			ResourceID        string `json:"resourceId"`        //资源实例id
			LabelInfo         string `json:"labelInfo"`         //资源标签
			MasterOrderNo     string `json:"masterOrderNo"`     //主订单编码
			Usage             string `json:"usage"`             //使用量
			BillMode          string `json:"billMode"`          //计费模式 1：包周期 2：按需
			DiscountAmount    string `json:"discountAmount"`    //优惠金额
			ConsumeDate       string `json:"consumeDate"`       //消费时间
			DeductUsage       string `json:"deductUsage"`       //
			PayableAmount     string `json:"payableAmount"`     //应付金额
			PricefactorValue  string `json:"pricefactorValue"`  //
			ProductName       string `json:"productName"`       //资源产品名称
			RegionCode        string `json:"regionCode"`        //区域编码
			ServID            string `json:"servId"`            //资源实例id
			Price             string `json:"price"`             //官网价（原价）
			ServiceTag        string `json:"serviceTag"`        //服务内类型
			ContractName      string `json:"contractName"`      //合同名称
			RealResourceID    string `json:"realResourceId"`    //真实资源id
			SalesAttribute    string `json:"salesAttribute"`    //销售品规格
			KeySalesAttribute string `json:"keySalesAttribute"` //关键规格
			UsageType         string `json:"usageType"`         //使用量类型
			Amount            string `json:"amount"`            //金额
			UsageTypeID       string `json:"usageTypeId"`       //使用量类型id
			OfferName         string `json:"offerName"`         //销售品名称
			Coupon            string `json:"coupon"`            //代金券抵扣
			BillType          string `json:"billType"`          //账单类型
			MasterOrderID     string `json:"masterOrderId"`     //主订单id
			ProductCode       string `json:"productCode"`       //产品代码
			BillingCycleID    string `json:"billingCycleId"`    //计费周期
			RegionID          string `json:"regionId"`          //区域id
			StateDate         string `json:"stateDate"`         //状态时间
			ProjectName       string `json:"projectName"`       //项目名
			ContractCode      string `json:"contractCode"`      //合同编码
			ProjectID         string `json:"projectId"`         //项目ID
			ResourceType      string `json:"resourceType"`      //资源类型
		} `json:"result"`
		PageNo     int `json:"pageNo"`
		PageSize   int `json:"pageSize"`
		TotalCount int `json:"totalCount"`
	} `json:"returnObj"`
}
