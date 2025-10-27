package validator

import (
	"fmt"
	"reflect"

	"github.com/wangweihong/gotoolbox/pkg/validation"
)

type Validator interface {
	Validate() error
}

// ValidateAll 先基于结构体字段进行binding tag的go-validator检测, 如果字段实现了Validator, 再次进行验证
func ValidateAll(s interface{}) error {
	cval := NewCustomValidator(validation.LangEN)
	cval.Engine()

	if err := cval.Validate(s); err != nil {
		return err
	}

	return validateCustom(s)
}

// 自定义验证入口
func validateCustom(s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil // 跳过nil指针
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil // 非结构体直接返回
	}
	return validateRecursive(v)
}

// 递归验证函数
func validateRecursive(v reflect.Value) error {
	// 处理指针
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	// 仅处理结构体
	if v.Kind() != reflect.Struct {
		return nil
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		// 跳过不可导出字段
		if !field.CanInterface() {
			continue
		}
		// 优先检查并执行自定义验证
		if validator, ok := getCustomValidator(field); ok {
			if err := validator.Validate(); err != nil {
				return fmt.Errorf("%s: %w", fieldType.Name, err)
			}
			continue // 跳过递归
		}
		// 递归处理嵌套结构
		switch field.Kind() {
		case reflect.Struct:
			if err := validateRecursive(field); err != nil {
				return err
			}
		case reflect.Ptr:
			if field.Type().Elem().Kind() == reflect.Struct {
				if err := validateRecursive(field); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func getCustomValidator(field reflect.Value) (Validator, bool) {
	if validator, ok := field.Interface().(Validator); ok {
		return validator, true
	}
	if field.CanAddr() {
		if validator, ok := field.Addr().Interface().(Validator); ok {
			return validator, true
		}
	}
	return nil, false
}
