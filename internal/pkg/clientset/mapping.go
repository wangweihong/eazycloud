package clientset

import (
	"context"
	"reflect"

	"github.com/wangweihong/eazycloud/apis/iapiserver"

	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"

	"k8s.io/apimachinery/pkg/runtime"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ResourceGetFunc func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error)

// use reflect object to get data from informer?
var ObjectReflectTypeMap = map[reflect.Type]ResourceGetFunc{
	//corev1
	reflect.TypeOf(&corev1.Service{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return ServiceGet(ctx, cluster, namespace, res.(*corev1.Service).GetName(), opts)
	},
	reflect.TypeOf(&corev1.ConfigMap{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return ConfigMapGet(ctx, cluster, namespace, res.(*corev1.ConfigMap).GetName(), opts)
	},
	reflect.TypeOf(&corev1.Secret{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return SecretGet(ctx, cluster, namespace, res.(*corev1.Secret).GetName(), opts)
	},
	reflect.TypeOf(&corev1.Event{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return EventGet(ctx, cluster, namespace, res.(*corev1.Event).GetName(), opts)
	},
	reflect.TypeOf(&corev1.Node{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return NodeGet(ctx, cluster, res.(*corev1.Node).GetName(), opts)
	},
	reflect.TypeOf(&corev1.LimitRange{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return LimitRangeGet(ctx, cluster, namespace, res.(*corev1.LimitRange).GetName(), opts)
	},
	reflect.TypeOf(&corev1.Namespace{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return NamespaceGet(ctx, cluster, res.(*corev1.Namespace).GetName(), opts)
	},
	reflect.TypeOf(&corev1.Endpoints{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return EndpointsGet(ctx, cluster, namespace, res.(*corev1.Endpoints).GetName(), opts)
	},
	reflect.TypeOf(&corev1.PersistentVolume{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return PersistentVolumeGet(ctx, cluster, res.(*corev1.PersistentVolume).GetName(), opts)
	},
	reflect.TypeOf(&corev1.PersistentVolumeClaim{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return PersistentVolumeClaimGet(ctx, cluster, namespace, res.(*corev1.PersistentVolumeClaim).GetName(), opts)
	},
	reflect.TypeOf(&corev1.Pod{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return PodGet(ctx, cluster, namespace, res.(*corev1.Pod).GetName(), opts)
	},
	reflect.TypeOf(&corev1.ServiceAccount{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return ServiceAccountGet(ctx, cluster, namespace, res.(*corev1.ServiceAccount), opts)
	},
	reflect.TypeOf(&corev1.ResourceQuota{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return ResourceQuotaGet(ctx, cluster, namespace, res.(*corev1.ResourceQuota).GetName(), opts)
	},

	//apps
	reflect.TypeOf(&appsv1.Deployment{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return DeploymentGet(ctx, cluster, namespace, res.(*appsv1.Deployment).GetName(), opts)
	},
	reflect.TypeOf(&appsv1.ReplicaSet{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return ReplicaSetGet(ctx, cluster, namespace, res.(*appsv1.ReplicaSet).GetName(), opts)
	},
	reflect.TypeOf(&appsv1.StatefulSet{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return StatefulSetGet(ctx, cluster, namespace, res.(*appsv1.StatefulSet).GetName(), opts)
	},
	reflect.TypeOf(&appsv1.DaemonSet{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return DaemonSetGet(ctx, cluster, namespace, res.(*appsv1.DaemonSet).GetName(), opts)
	},
	//extensions
	reflect.TypeOf(&networkingv1.Ingress{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return IngressGet(ctx, cluster, namespace, res.(*networkingv1.Ingress).GetName(), opts)
	},
	//autoscaling
	reflect.TypeOf(&autoscalingv1.HorizontalPodAutoscaler{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return HpaGet(ctx, cluster, namespace, res.(*autoscalingv1.HorizontalPodAutoscaler).GetName(), opts)
	},

	// batch
	reflect.TypeOf(&batchv1.Job{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return JobGet(ctx, cluster, namespace, res.(*batchv1.Job).GetName(), opts)
	},
	reflect.TypeOf(&batchv1beta1.CronJob{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return CronJobGet(ctx, cluster, namespace, res.(*batchv1.CronJob).GetName(), opts)
	},

	//rbac
	reflect.TypeOf(&rbacv1.Role{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return RoleGet(ctx, cluster, namespace, res.(*rbacv1.Role), opts)
	},
	reflect.TypeOf(&rbacv1.RoleBinding{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return RoleBindingGet(ctx, cluster, namespace, res.(*rbacv1.RoleBinding), opts)
	},
	reflect.TypeOf(&rbacv1.ClusterRole{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return ClusterRoleGet(ctx, cluster, res.(*rbacv1.ClusterRole), opts)
	},
	reflect.TypeOf(&rbacv1.ClusterRoleBinding{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return ClusterRoleBindingGet(ctx, cluster, res.(*rbacv1.ClusterRoleBinding), opts)
	},
	reflect.TypeOf(&storagev1.StorageClass{}): func(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res runtime.Object, opts metav1.GetOptions) (runtime.Object, error) {
		return StorageClassGet(ctx, cluster, res.(*storagev1.StorageClass).GetName(), opts)
	},
}
