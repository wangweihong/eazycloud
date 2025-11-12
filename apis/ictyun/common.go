package ictyun

import "fmt"

type PagingParam struct {
	PageNum  *int `json:"pageNo"`   //页码，取值范围：正整数（≥1），注：默认值为
	PageSize *int `json:"pageSize"` //每页记录数目，取值范围：[1, 50]，注：默认值为1
}

type ResponseResult struct {
	ErrorCode   string      `json:"errorCode"`
	Error       string      `json:"error"`
	Details     string      `json:"details"`
	Message     string      `json:"message"`
	Description string      `json:"description"`
	StatusCode  interface{} `json:"statusCode"`
	MsgDesc     string      `json:"msgDesc"`
}

func (r ResponseResult) String() string {
	return fmt.Sprintf("errorCode:%s,error:%s,details:%s,message:%s,description:%s,statusCode:%v,msgDesc:%v",
		r.ErrorCode, r.Error, r.Details, r.Message, r.Description, r.StatusCode, r.MsgDesc)
}
