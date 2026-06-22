package iapiserver

import (
	gvalidator "github.com/go-playground/validator/v10"
	"github.com/wangweihong/eazycloud/pkg/validator"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	validator.RegisterValidatorNoTrans("namespaced", ValidateNamespaceScopeResource)
	validator.RegisterValidatorNoTrans("clusterd", ValidateClusterScopeResource)
}

// ValidateNamespaceScopeResource 检验命名空间级资源是否合法
func ValidateNamespaceScopeResource(fl gvalidator.FieldLevel) bool {
	if fl.Field().Interface() == nil {
		return false
	}

	// 如果是iapiserver.ResourceGetRequest结构则直接判断
	if gr, ok := fl.Field().Interface().(ResourceGetRequest); ok {
		return !(gr.Namespace == "" || gr.Name == "")
	}
	// 这里是因为之前没法找到获取检测的结构体字段临时想到的方法
	// fieldVal := fl.Field()
	// if fieldVal.Kind() != reflect.Ptr {
	// 	fieldVal = reflectutil.CreatePointerToValue(fl.Field())
	// }
	if !fl.Field().CanAddr() {
		return false
	}
	// 即使结构体中定义的指针,fl.Field()获得解引用的类型。因此需要通过fl.Field().Addr()
	cm, ok := fl.Field().Addr().Interface().(metav1.Object)
	if !ok || cm == nil {
		return false
	}

	if cm.GetNamespace() == "" || cm.GetName() == "" {
		return false
	}

	return true
}

// ValidateClusterScopeResource 检验集群级资源是否合法
func ValidateClusterScopeResource(fl gvalidator.FieldLevel) bool {
	if fl.Field().Interface() == nil {
		return false
	}

	if gr, ok := fl.Field().Interface().(ResourceGetRequest); ok {
		return !(gr.Name == "")
	}

	if !fl.Field().CanAddr() {
		return false
	}

	cm, ok := fl.Field().Addr().Interface().(metav1.Object)
	if !ok || cm == nil {
		return false
	}

	if cm.GetName() == "" {
		return false
	}

	return true
}
