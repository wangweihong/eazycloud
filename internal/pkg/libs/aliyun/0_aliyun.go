package aliyun

import (
	"fmt"
	"reflect"

	util "github.com/alibabacloud-go/tea-utils/v2/service"
)

const (
	defaultRegion string = "cn-shenzhen"
	futureRegion  string = "cn-beijing"
)

type client struct {
	address   string
	accessKey string
	secretKey string
}

type Invoker func(opt util.RuntimeOptions) (interface{}, error)

type InvokeOption struct {
	opt util.RuntimeOptions
}

func invoke(invoker Invoker, resp interface{}, opt ...InvokeOption) error {
	ropt := util.RuntimeOptions{}
	if opt != nil {
		ropt = opt[0].opt
	}

	ret, err := invoker(ropt)
	if err != nil {
		return err
	}

	// 当必需未传时, 返回值可能为nil
	if ret == nil {
		return fmt.Errorf("ret is nil")
	}

	if resp != nil {
		rt := reflect.TypeOf(resp)
		if rt.Kind() == reflect.Ptr && rt.Elem().Kind() == reflect.Ptr {
			if rt.Elem() == reflect.TypeOf(ret) {
				reflect.Indirect(reflect.ValueOf(resp)).Set(reflect.ValueOf(ret))
				return nil
			}
			return fmt.Errorf("response type assertion fail:%v", rt.String())
		} else {
			return fmt.Errorf("resp is not pointer")
		}
	}

	return nil
}
