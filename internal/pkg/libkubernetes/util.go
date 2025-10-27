package libkubernetes

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/wangweihong/eazycloud/apis/ikubernetes"

	snapshotv1beta1 "github.com/kubernetes-csi/external-snapshotter/client/v3/apis/volumesnapshot/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

var (
	MetaMap              = map[string]metav1.TypeMeta{}
	ObjectReflectTypeMap = map[reflect.Type]metav1.TypeMeta{}
	Scheme               = runtime.NewScheme()
)

func init() {
	_ = clientgoscheme.AddToScheme(Scheme)

	MetaMap = map[string]metav1.TypeMeta{
		//corev1
		ikubernetes.KubernetesResourceKindService:               {Kind: "Service", APIVersion: "v1"},
		ikubernetes.KubernetesResourceKindConfigMap:             {Kind: "ConfigMap", APIVersion: "v1"},
		ikubernetes.KubernetesResourceKindSecret:                {Kind: "Secret", APIVersion: "v1"},
		ikubernetes.KubernetesResourceKindEvent:                 {Kind: "Event", APIVersion: "v1"},
		ikubernetes.KubernetesResourceKindNode:                  {Kind: "Node", APIVersion: "v1"},
		ikubernetes.KubernetesResourceKindLimitRange:            {Kind: "LimitRange", APIVersion: "v1"},
		ikubernetes.KubernetesResourceKindNamespace:             {Kind: "Namespace", APIVersion: "v1"},
		ikubernetes.KubernetesResourceKindEndpoints:             {Kind: "Endpoints", APIVersion: "v1"},
		ikubernetes.KubernetesResourceKindPersistentVolume:      {Kind: "PersistentVolume", APIVersion: "v1"},
		ikubernetes.KubernetesResourceKindPersistentVolumeClaim: {Kind: "PersistentVolumeClaim", APIVersion: "v1"},
		ikubernetes.KubernetesResourceKindPod:                   {Kind: "Pod", APIVersion: "v1"},
		ikubernetes.KubernetesResourceKindServiceAccount:        {Kind: "ServiceAccount", APIVersion: "v1"},
		ikubernetes.KubernetesResourceKindResourceQuota:         {Kind: "ResourceQuota", APIVersion: "v1"},
		//apps
		ikubernetes.KubernetesResourceKindDeployment:  {Kind: "Deployment", APIVersion: "apps/v1"},
		ikubernetes.KubernetesResourceKindReplicaSet:  {Kind: "ReplicaSet", APIVersion: "apps/v1"},
		ikubernetes.KubernetesResourceKindStatefulSet: {Kind: "StatefulSet", APIVersion: "apps/v1"},
		ikubernetes.KubernetesResourceKindDaemonSet:   {Kind: "DaemonSet", APIVersion: "apps/v1"},
		// FIXME
		ikubernetes.KubernetesResourceKindIngress: {Kind: "Ingress", APIVersion: "extensions/v1beta1"},
		//autoscaling
		ikubernetes.KubernetesResourceKindHorizontalPodAutoscaler: {
			Kind:       "HorizontalPodAutoscaler",
			APIVersion: "autoscaling/v1",
		},

		// batch
		ikubernetes.KubernetesResourceKindJob:     {Kind: "Job", APIVersion: "batch/v1"},
		ikubernetes.KubernetesResourceKindCronJob: {Kind: "CronJob", APIVersion: "batch/v1beta1"},

		//rbac
		ikubernetes.KubernetesResourceKindRole: {Kind: "Role", APIVersion: "rbac.authorization.k8s.io/v1"},
		ikubernetes.KubernetesResourceKindRoleBinding: {
			Kind:       "RoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ikubernetes.KubernetesResourceKindClusterRole: {
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ikubernetes.KubernetesResourceKindClusterRoleBinding: {
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},

		//crd
		ikubernetes.KubernetesResourceKindStorageClass: {Kind: "StorageClass", APIVersion: "storage.k8s.io/v1"},

		// monitor : prometheus-operator
		ikubernetes.KubernetesResourceKindServiceMonitor: {Kind: "ServiceMonitor", APIVersion: "monitoring.coreos.com/v1"},
		ikubernetes.KubernetesResourceKindPodMonitor:     {Kind: "PodMonitor", APIVersion: "monitoring.coreos.com/v1"},
		ikubernetes.KubernetesResourceKindAlertManager:   {Kind: "Alertmanager", APIVersion: "monitoring.coreos.com/v1"},
		ikubernetes.KubernetesResourceKindPrometheus:     {Kind: "Prometheus", APIVersion: "monitoring.coreos.com/v1"},
		ikubernetes.KubernetesResourceKindPrometheusRule: {Kind: "PrometheusRule", APIVersion: "monitoring.coreos.com/v1"},
	}

	ObjectReflectTypeMap = map[reflect.Type]metav1.TypeMeta{
		//corev1
		reflect.TypeOf(&corev1.Service{}):               {Kind: "Service", APIVersion: "v1"},
		reflect.TypeOf(&corev1.ConfigMap{}):             {Kind: "ConfigMap", APIVersion: "v1"},
		reflect.TypeOf(&corev1.Secret{}):                {Kind: "Secret", APIVersion: "v1"},
		reflect.TypeOf(&corev1.Event{}):                 {Kind: "Event", APIVersion: "v1"},
		reflect.TypeOf(&corev1.Node{}):                  {Kind: "Node", APIVersion: "v1"},
		reflect.TypeOf(&corev1.LimitRange{}):            {Kind: "LimitRange", APIVersion: "v1"},
		reflect.TypeOf(&corev1.Namespace{}):             {Kind: "Namespace", APIVersion: "v1"},
		reflect.TypeOf(&corev1.Endpoints{}):             {Kind: "Endpoints", APIVersion: "v1"},
		reflect.TypeOf(&corev1.PersistentVolume{}):      {Kind: "PersistentVolume", APIVersion: "v1"},
		reflect.TypeOf(&corev1.PersistentVolumeClaim{}): {Kind: "PersistentVolumeClaim", APIVersion: "v1"},
		reflect.TypeOf(&corev1.Pod{}):                   {Kind: "Pod", APIVersion: "v1"},
		reflect.TypeOf(&corev1.ServiceAccount{}):        {Kind: "ServiceAccount", APIVersion: "v1"},
		reflect.TypeOf(&corev1.ResourceQuota{}):         {Kind: "ResourceQuota", APIVersion: "v1"},
		//apps
		reflect.TypeOf(&appsv1.Deployment{}):         {Kind: "Deployment", APIVersion: "apps/v1"},
		reflect.TypeOf(&appsv1.ReplicaSet{}):         {Kind: "ReplicaSet", APIVersion: "apps/v1"},
		reflect.TypeOf(&appsv1.StatefulSet{}):        {Kind: "StatefulSet", APIVersion: "apps/v1"},
		reflect.TypeOf(&appsv1.DaemonSet{}):          {Kind: "DaemonSet", APIVersion: "apps/v1"},
		reflect.TypeOf(&appsv1.ControllerRevision{}): {Kind: "ControllerRevision", APIVersion: "apps/v1"},
		//extensions
		reflect.TypeOf(&extensionsv1beta1.Ingress{}): {Kind: "Ingress", APIVersion: "extensions/v1beta1"},
		//autoscaling
		reflect.TypeOf(&autoscalingv1.HorizontalPodAutoscaler{}): {
			Kind:       "HorizontalPodAutoscaler",
			APIVersion: "autoscaling/v1",
		},
		// batch
		reflect.TypeOf(&batchv1.Job{}):          {Kind: "Job", APIVersion: "batch/v1"},
		reflect.TypeOf(&batchv1beta1.CronJob{}): {Kind: "CronJob", APIVersion: "batch/v1beta1"},
		//rbac
		reflect.TypeOf(&rbacv1.Role{}): {Kind: "Role", APIVersion: "rbac.authorization.k8s.io/v1"},
		reflect.TypeOf(&rbacv1.RoleBinding{}): {
			Kind:       "RoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		reflect.TypeOf(&rbacv1.ClusterRole{}): {
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		reflect.TypeOf(&rbacv1.ClusterRoleBinding{}): {
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		reflect.TypeOf(&storagev1.StorageClass{}):     {Kind: "StorageClass", APIVersion: "storage.k8s.io/v1"},
		reflect.TypeOf(&networkingv1.NetworkPolicy{}): {Kind: "NetworkPolicy", APIVersion: "networking.k8s.io/v1"},

		reflect.TypeOf(&networkingv1.NetworkPolicy{}): {Kind: "NetworkPolicy", APIVersion: "networking.k8s.io/v1"},

		reflect.TypeOf(&snapshotv1beta1.VolumeSnapshot{}): {
			Kind:       "VolumeSnapshot",
			APIVersion: "snapshot.storage.k8s.io/v1",
		},
		reflect.TypeOf(&snapshotv1beta1.VolumeSnapshotClass{}): {
			Kind:       "VolumeSnapshotClass",
			APIVersion: "snapshot.storage.k8s.io/v1",
		},
		reflect.TypeOf(&snapshotv1beta1.VolumeSnapshotContent{}): {
			Kind:       "VolumeSnapshotContent",
			APIVersion: "snapshot.storage.k8s.io/v1",
		},
	}
}

func getResourceTypeMeta(resourceName string) metav1.TypeMeta {
	meta, ok := MetaMap[resourceName]
	if !ok {
		return metav1.TypeMeta{Kind: "Unknown", APIVersion: "Unknown"}
	}
	return meta
}

const (
	defaultRequestTimeout = 1 * time.Minute
)

func run[T any](ctx context.Context, config *ikubernetes.ClusterConfig, action func(context.Context, *Client) (T, error), options ...Option) (T, error) {
	ci := &callInfo{timeout: defaultRequestTimeout}
	for _, option := range options {
		option(ci)
	}
	var zero T
	c, err := NewClient(config)
	if err != nil {
		return zero, err
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, time.Duration(ci.timeout)*time.Second)
	defer cancel()

	ret, err := action(ctx, c)
	if err != nil {
		return zero, err
	}

	HideManagerField(ret)

	InjectTypeMeta(ret)
	return ret, nil
}

type Action[T any, K any] func(context.Context, *Client, T, K) error

func HideManagerField(obj any) {
	if obj == nil {
		return
	}

	switch data := obj.(type) {
	case metav1.Object:
		data.SetManagedFields(nil)
		return
	case metav1.ListInterface:
		HideListManagerField(data)
		return
	}
}

func HideListManagerField(list metav1.ListInterface) {
	if list == nil {
		return
	}
	iterateStructField := func(rt reflect.Type, rv reflect.Value) {
		fieldRV := rv.FieldByName("Items")
		f, _ := rt.FieldByName("Items")
		fieldRt := f.Type

		if !fieldRV.IsNil() {
			if fieldRt.Kind() == reflect.Slice {
				for i := 0; i < fieldRV.Len(); i++ {
					iv := fieldRV.Index(i)
					it := reflect.TypeOf(iv)

					switch it.Kind() {
					case reflect.Struct:
						obj := iv.Addr().Interface().(metav1.Object)
						obj.SetManagedFields(nil)
					case reflect.Ptr:
						obj := iv.Interface().(metav1.Object)
						obj.SetManagedFields(nil)
					}
				}
			}
			return
		}
	}

	rt := reflect.TypeOf(list)
	rv := reflect.ValueOf(list)
	switch rt.Kind() {
	case reflect.Struct:
		iterateStructField(rt, rv)
		return
	case reflect.Ptr:
		if !rv.IsValid() || rv.IsNil() {
			return
		}

		if rt.Elem().Kind() != reflect.Struct {
			return
		}
		rv = rv.Elem()
		rt = rt.Elem()
		if !rv.IsValid() {
			return
		}

		iterateStructField(rt, rv)
		return
	}
}

// inject apiVersion and Kind into object list
func InjectTypeMeta(obj any) {
	if obj == nil {
		return
	}

	switch data := obj.(type) {
	//check first, because list is also a runtime.Object.
	case metav1.ListInterface:
		InjectListManagerField(data)
		return

	case runtime.Object:
		reflectMeta := GetResourceTypeMetaFromReflectType(reflect.TypeOf(obj))
		_ = meta.NewAccessor().SetKind(data, reflectMeta.Kind)
		_ = meta.NewAccessor().SetAPIVersion(data, reflectMeta.APIVersion)
		return

	}
}

func InjectListManagerField(list metav1.ListInterface) {
	if list == nil {
		return
	}
	iterateStructField := func(rt reflect.Type, rv reflect.Value) {
		fieldRV := rv.FieldByName("Items")
		f, _ := rt.FieldByName("Items")
		fieldRt := f.Type

		if !fieldRV.IsNil() {
			if fieldRt.Kind() == reflect.Slice {
				for i := 0; i < fieldRV.Len(); i++ {
					iv := fieldRV.Index(i)
					it := reflect.TypeOf(iv)

					switch it.Kind() {
					case reflect.Struct:
						data := iv.Addr().Interface().(runtime.Object)
						reflectMeta := GetResourceTypeMetaFromReflectType(reflect.TypeOf(data))
						_ = meta.NewAccessor().SetKind(data, reflectMeta.Kind)
						_ = meta.NewAccessor().SetAPIVersion(data, reflectMeta.APIVersion)
					case reflect.Ptr:
						data := iv.Interface().(runtime.Object)
						reflectMeta := GetResourceTypeMetaFromReflectType(reflect.TypeOf(data))
						_ = meta.NewAccessor().SetKind(data, reflectMeta.Kind)
						_ = meta.NewAccessor().SetAPIVersion(data, reflectMeta.APIVersion)
					}
				}
			}
			return
		}
	}

	rt := reflect.TypeOf(list)
	rv := reflect.ValueOf(list)
	switch rt.Kind() {
	case reflect.Struct:
		iterateStructField(rt, rv)
		return
	case reflect.Ptr:
		if !rv.IsValid() || rv.IsNil() {
			return
		}

		if rt.Elem().Kind() != reflect.Struct {
			return
		}
		rv = rv.Elem()
		rt = rt.Elem()
		if !rv.IsValid() {
			return
		}

		iterateStructField(rt, rv)
		return
	}
}

func GetResourceTypeMetaFromReflectType(rtype reflect.Type) metav1.TypeMeta {
	meta, ok := ObjectReflectTypeMap[rtype]
	if !ok {
		return metav1.TypeMeta{Kind: "", APIVersion: ""}
	}
	return meta
}

const (
	annotationtScheduleChineseKey = "schedule_translate"
)

func TranslateCronJobSchedule(job *batchv1.CronJob) {
	translateCronJobSchedule(job)
}

func translateCronJobSchedule(job *batchv1.CronJob) {
	chinese := TranslateCronToChinese(job.Spec.Schedule)
	if job.Annotations == nil {
		job.Annotations = make(map[string]string)
	}
	job.Annotations[annotationtScheduleChineseKey] = chinese
}
func TranslateCronJobListSchedule(jobList *batchv1.CronJobList) {
	translateCronJobListSchedule(jobList)
}
func translateCronJobListSchedule(jobList *batchv1.CronJobList) {
	for i := range jobList.Items {
		translateCronJobSchedule(&jobList.Items[i])
	}
}

func handleSlash(str string, unit string) (trans string) {
	trans = str + unit
	slashSlice := strings.SplitN(str, "/", 2)
	if len(slashSlice) == 1 {
		return
	}

	trans = ""
	if strings.Contains(slashSlice[0], "-") {
		trans = handleCross(slashSlice[0], unit)
	} else {
		trans = "从" + slashSlice[0] + unit + "开始"
	}

	trans += "每" + slashSlice[1] + unit
	return trans
}

// x-y 表示x至y
func handleCross(str string, unit string) string {
	slashSlice := strings.SplitN(str, "-", 2)
	if len(slashSlice) == 1 {
		return str + unit
	}
	return slashSlice[0] + "至" + slashSlice[1] + unit
}

func translateMonth(cronItem string) string {
	if cronItem != "*" {
		return handleSlash(cronItem, "月")
	} else {
		return "每月"
	}
}
func translateWeek(cronItem string) string {
	result := ""
	if cronItem != "*" && cronItem != "?" {
		for _, v := range cronItem {
			switch string(v) {
			case "1":
				result += "星期天"
			case "2":
				result += "星期一"
			case "3":
				result += "星期二"
			case "4":
				result += "星期三"
			case "5":
				result += "星期四"
			case "6":
				result += "星期五"
			case "7":
				result += "星期六"
			case "-":
				result += "至"
			default:
				result += string(v)
			}
		}
	}
	return result
}

func translateDay(cronItem string) string {
	if cronItem != "?" {
		if cronItem != "*" {
			return handleSlash(cronItem, "日")
		} else {
			return "每日"
		}
	}
	return ""
}

func translateHour(cronItem string) string {
	if cronItem != "*" {
		return handleSlash(cronItem, "时")
	} else {
		return "每时"
	}
}

func translateMinute(cronItem string) string {
	if cronItem != "*" {
		return handleSlash(cronItem, "分钟")
	} else {
		return "每分钟"
	}
}

func translateSecond(cronItem string) string {
	if cronItem != "*" {
		return handleSlash(cronItem, "秒")
	} else {
		return "每秒"
	}
}
func TranslateCronToChinese(crontab string) (result string) {
	if crontab == "" {
		//	err = fmt.Errorf("crontab is empty")
		return
	}
	cronslice := strings.Split(crontab, " ")
	switch len(cronslice) {
	case 6:
		result += translateMonth(cronslice[4])  //解析月
		result += translateWeek(cronslice[5])   //解析周
		result += translateDay(cronslice[3])    //解析天
		result += translateHour(cronslice[2])   //解析时
		result += translateMinute(cronslice[1]) //解析分
		result += translateSecond(cronslice[0]) //解析秒
	case 5:
		result += translateWeek(cronslice[4]) //解析周
		//	result += translateWeek(cronslice[4])   //解析周
		result += translateMonth(cronslice[3])  //解析月
		result += translateDay(cronslice[2])    //解析天
		result += translateHour(cronslice[1])   //解析时
		result += translateMinute(cronslice[0]) //解析分
	default:
		result = crontab
	}

	//container not translate,return origin
	if strings.Contains(result, "-") || strings.Contains(result, "/") {
		result = crontab
	}

	// 说人话
	if strings.HasPrefix(result, "每月每日每时每分钟") {
		result = strings.TrimPrefix(result, "每月每日每时每分钟")
		if strings.Contains(result, "每") {
			return result
		}
		return "每" + result
	}

	if strings.HasPrefix(result, "每月每日每时") {
		result = strings.TrimPrefix(result, "每月每日每时")
		if strings.Contains(result, "每") {
			return result
		}
		return "每" + result
	}

	if strings.HasPrefix(result, "每月每日") {
		result = strings.TrimPrefix(result, "每月每日")
		if strings.Contains(result, "每") {
			return result
		}
		return "每" + result
	}

	if strings.HasSuffix(result, "每秒") {
		result = strings.TrimSuffix(result, "每秒")
	}

	if strings.HasSuffix(result, "每分钟") {
		result = strings.TrimSuffix(result, "每分钟")
	}

	if strings.HasSuffix(result, "每时") {
		result = strings.TrimSuffix(result, "每时")
	}

	return result
}
