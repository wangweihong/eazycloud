package ikubernetes

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	KubernetesResourceKindService                 = "service"
	KubernetesResourceKindConfigMap               = "configmap"
	KubernetesResourceKindSecret                  = "secret"
	KubernetesResourceKindEvent                   = "event"
	KubernetesResourceKindNode                    = "node"
	KubernetesResourceKindLimitRange              = "limitrange"
	KubernetesResourceKindNamespace               = "namespace"
	KubernetesResourceKindEndpoints               = "endpoints"
	KubernetesResourceKindPersistentVolume        = "persistentvolume"
	KubernetesResourceKindPersistentVolumeClaim   = "persistentvolumeclaim"
	KubernetesResourceKindPod                     = "pod"
	KubernetesResourceKindServiceAccount          = "serviceaccount"
	KubernetesResourceKindResourceQuota           = "resourcequota"
	KubernetesResourceKindDeployment              = "deployment"
	KubernetesResourceKindReplicaSet              = "replicaset"
	KubernetesResourceKindStatefulSet             = "statefulset"
	KubernetesResourceKindDaemonSet               = "daemonset"
	KubernetesResourceKindIngress                 = "ingress"
	KubernetesResourceKindHorizontalPodAutoscaler = "hpa"
	KubernetesResourceKindJob                     = "job"
	KubernetesResourceKindCronJob                 = "cronjob"
	KubernetesResourceKindRole                    = "role"
	KubernetesResourceKindRoleBinding             = "rolebinding"
	KubernetesResourceKindClusterRole             = "clusterrole"
	KubernetesResourceKindClusterRoleBinding      = "clusterrolebinding"
	KubernetesResourceKindStorageClass            = "storageclass"
	KubernetesResourceKindApplication             = "application"
	KubernetesResourceKindNetworkPolicy           = "NetworkPolicy"
	//
	KubernetesResourceKindServiceMonitor = "serviceMonitor"
	KubernetesResourceKindPodMoniotr     = "podMonitor"
	KubernetesResourceKindAlertManager   = "alertManager"
	KubernetesResourceKindPrometheus     = "prometheus"
	KubernetesResourceKindPrometheusRule = "prometheusRule"
)

type ResourceMeta struct {
	metav1.TypeMeta
	Config any // 指针
}
