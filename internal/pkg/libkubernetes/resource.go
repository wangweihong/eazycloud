package libkubernetes

import (
	"context"
	"reflect"

	"github.com/wangweihong/gotoolbox/pkg/errors"

	"github.com/wangweihong/eazycloud/apis/ikubernetes"

	corev1 "k8s.io/api/core/v1"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	batchv1 "k8s.io/api/batch/v1"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
)

func ResourceCreate(
	config *ikubernetes.ClusterConfig,
	namespace string,
	metas []*ikubernetes.ResourceMeta,
	opts metav1.CreateOptions,
) error {
	var cancels []context.CancelFunc
	defer func() {
		for _, v := range cancels {
			v()
		}
	}()

	c, err := NewClient(config)
	if err != nil {
		return err
	}

	//cm, err := NewMonitorClient(config)
	//if err != nil {
	//	return err
	//}
	//defer cm.Close()

	for _, v := range metas {
		ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
		cancels = append(cancels, cancel)

		switch v.Kind {
		case "Pod":
			_, err = c.client.CoreV1().Pods(namespace).Create(ctx, (v.Config).(*corev1.Pod), opts)
		case "NetworkPolicy":
			_, err = c.client.NetworkingV1().NetworkPolicies(namespace).Create(ctx, (v.Config).(*networkingv1.NetworkPolicy), opts)
		case "Secret":
			_, err = c.client.CoreV1().Secrets(namespace).Create(ctx, (v.Config).(*corev1.Secret), opts)
		case "ResourceQuota":
			_, err = c.client.CoreV1().ResourceQuotas(namespace).Create(ctx, (v.Config).(*corev1.ResourceQuota), opts)
		case "ConfigMap":
			_, err = c.client.CoreV1().ConfigMaps(namespace).Create(ctx, (v.Config).(*corev1.ConfigMap), opts)
		case "Endpoints":
			_, err = c.client.CoreV1().Endpoints(namespace).Create(ctx, (v.Config).(*corev1.Endpoints), opts)
		case "Event":
			_, err = c.client.CoreV1().Events(namespace).Create(ctx, (v.Config).(*corev1.Event), opts)
		case "LimitRange":
			_, err = c.client.CoreV1().LimitRanges(namespace).Create(ctx, (v.Config).(*corev1.LimitRange), opts)
		case "PersistentVolumeClaim":
			_, err = c.client.CoreV1().PersistentVolumeClaims(namespace).Create(ctx, (v.Config).(*corev1.PersistentVolumeClaim), opts)
		case "Ingress":
			_, err = c.client.NetworkingV1().Ingresses(namespace).Create(ctx, (v.Config).(*networkingv1.Ingress), opts)
		case "PodTemplate":
			_, err = c.client.CoreV1().PodTemplates(namespace).Create(ctx, (v.Config).(*corev1.PodTemplate), opts)
		case "HorizontalPodAutoscaler":
			_, err = c.client.AutoscalingV1().HorizontalPodAutoscalers(namespace).Create(ctx, (v.Config).(*autoscalingv1.HorizontalPodAutoscaler), opts)
		case "ReplicationController":
			_, err = c.client.CoreV1().ReplicationControllers(namespace).Create(ctx, (v.Config).(*corev1.ReplicationController), opts)
		case "ServiceAccount":
			_, err = c.client.CoreV1().ServiceAccounts(namespace).Create(ctx, (v.Config).(*corev1.ServiceAccount), opts)
		case "Service":
			_, err = c.client.CoreV1().Services(namespace).Create(ctx, (v.Config).(*corev1.Service), opts)
		case "Deployment":
			_, err = c.client.AppsV1().Deployments(namespace).Create(ctx, (v.Config).(*appsv1.Deployment), opts)
		case "DaemonSet":
			_, err = c.client.AppsV1().DaemonSets(namespace).Create(ctx, (v.Config).(*appsv1.DaemonSet), opts)
		case "ReplicaSet":
			_, err = c.client.AppsV1().ReplicaSets(namespace).Create(ctx, (v.Config).(*appsv1.ReplicaSet), opts)
		case "StatefulSet":
			_, err = c.client.AppsV1().StatefulSets(namespace).Create(ctx, (v.Config).(*appsv1.StatefulSet), opts)
		case "Job":
			_, err = c.client.BatchV1().Jobs(namespace).Create(ctx, (v.Config).(*batchv1.Job), opts)
		case "CronJob":
			_, err = c.client.BatchV1().CronJobs(namespace).Create(ctx, (v.Config).(*batchv1.CronJob), opts)
		case "ClusterRole":
			_, err = c.client.RbacV1().ClusterRoles().Create(ctx, (v.Config).(*rbacv1.ClusterRole), opts)
		case "ClusterRoleBinding":
			_, err = c.client.RbacV1().ClusterRoleBindings().Create(ctx, (v.Config).(*rbacv1.ClusterRoleBinding), opts)
		case "Role":
			_, err = c.client.RbacV1().Roles(namespace).Create(ctx, (v.Config).(*rbacv1.Role), opts)
		case "RoleBinding":
			_, err = c.client.RbacV1().RoleBindings(namespace).Create(ctx, (v.Config).(*rbacv1.RoleBinding), opts)
		case "StorageClass":
			// avoid multiple default storageclass exists
			storageClass := v.Config.(*storagev1.StorageClass)
			delete(storageClass.Annotations, "storageclass.kubernetes.io/is-default-class")

			_, err = c.client.StorageV1().StorageClasses().Create(ctx, (v.Config).(*storagev1.StorageClass), opts)
		//case "ServiceMonitor":
		// 	_, err = cm.monitorClient.ServiceMonitors(monitoringNamespace).Create(ctx,
		// (v.Config).(*monitoringv1.ServiceMonitor), opts)
		//case "PodMonitor":
		// 	_, err = cm.monitorClient.PodMonitors(monitoringNamespace).Create(ctx,
		// (v.Config).(*monitoringv1.PodMonitor), opts)
		//case "Prometheus":
		// 	_, err = cm.monitorClient.Prometheuses(monitoringNamespace).Create(ctx,
		// (v.Config).(*monitoringv1.Prometheus), opts)
		//case "PrometheusRule":
		// 	_, err = cm.monitorClient.PrometheusRules(monitoringNamespace).Create(ctx,
		// (v.Config).(*monitoringv1.PrometheusRule), opts)
		//case "Alertmanager":
		// 	_, err = cm.monitorClient.Alertmanagers(monitoringNamespace).Create(ctx,
		// (v.Config).(*monitoringv1.Alertmanager), opts)
		case "PersistentVolume":
			_, err = c.client.CoreV1().PersistentVolumes().Create(ctx, (v.Config).(*corev1.PersistentVolume), opts)
		default:
			return errors.Errorf("resource not support,kind[%v]apiversion[%v]", v.Kind, v.APIVersion)

		}
		if err != nil {
			return errors.Errorf("kind[%v]apiversion[%v]err[%v]", v.Kind, v.APIVersion, err.Error())
		}
	}

	return nil
}

func ValidateClusterScopeParameters(name string) error {
	if name == "" {
		return errors.New("name is empty")
	}
	return nil
}

func ValidateNamespacedScopeParameters(namespace, name string) error {
	if name == "" {
		return errors.New("name is empty")
	}

	if namespace == "" {
		return errors.New("namespace is empty")
	}
	return nil
}

func ValidateNamespacedObjectParameters(req metav1.Object) error {
	// validate metav1.Object's implement if a nil Pointer. otherwise it will panic when get field
	if req == nil || reflect.ValueOf(req).IsNil() {
		return errors.Errorf("object is empty")
	}

	if req.GetName() == "" {
		return errors.Errorf("object name is empty")
	}

	if req.GetNamespace() == "" {
		return errors.Errorf("object namespace is empty")
	}
	return nil
}

func ValidateClusterScopedObjectParameters(req metav1.Object) error {
	if req == nil || reflect.ValueOf(req).IsNil() {
		return errors.Errorf("object is empty")
	}

	if req.GetName() == "" {
		return errors.Errorf("object name is empty")
	}

	return nil
}
