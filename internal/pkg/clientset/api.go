package clientset

import (
	"fmt"
	"io"
	"reflect"

	"github.com/wangweihong/gotoolbox/pkg/errors"

	"github.com/wangweihong/eazycloud/internal/pkg/libkubernetes"

	policyv1 "k8s.io/api/policy/v1"

	rbacv1 "k8s.io/api/rbac/v1"

	autoscalingv1 "k8s.io/api/autoscaling/v1"
	networkingv1 "k8s.io/api/networking/v1"

	storagev1 "k8s.io/api/storage/v1"

	k8stypes "k8s.io/apimachinery/pkg/types"

	appsv1 "k8s.io/api/apps/v1"

	batchv1 "k8s.io/api/batch/v1"

	"k8s.io/apimachinery/pkg/runtime"

	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/apis/ikubernetes"
)

type KubernetesListAction func(c *ResourceCacheController, labelSet labels.Set) (any, error)
type KubernetesGetAction func(c *ResourceCacheController) (any, error)

func informerAction[T any](ctx context.Context, cluster *iapiserver.Cluster, action func(context.Context, *ResourceCacheController, labels.Set) (T, error), opt *metav1.ListOptions) (T, error) {
	var zero T
	informerController, err := getInformerController(cluster)
	if err != nil {
		return zero, err
	}

	var labelSet labels.Set
	if opt != nil {
		labelSet, err = labels.ConvertSelectorToLabelsMap(opt.LabelSelector)
		if err != nil {
			return zero, err
		}
	}
	result, err := action(ctx, informerController, labelSet)
	if err != nil {
		return zero, err
	}

	libkubernetes.HideManagerField(result)
	libkubernetes.InjectTypeMeta(result)
	return result, nil
}

// type listFunc func(namespace string,resp any)error
func informerList(
	cluster *iapiserver.Cluster,
	opts metav1.ListOptions,
	action KubernetesListAction,
	resp any,
) (any, error) {
	informerController, err := getInformerController(cluster)
	if err != nil {
		return nil, err
	}
	labelSet, err := labels.ConvertSelectorToLabelsMap(opts.LabelSelector)
	if err != nil {
		return nil, err
	}
	result, err := action(informerController, labelSet)
	if err != nil {
		return nil, err
	}

	if result != nil {
		libkubernetes.HideManagerField(result)
		libkubernetes.InjectTypeMeta(result)
		if resp != nil {
			rt := reflect.TypeOf(resp)
			if rt.Kind() == reflect.Ptr && rt.Elem().Kind() == reflect.Ptr {
				if rt.Elem() == reflect.TypeOf(result) {
					reflect.Indirect(reflect.ValueOf(resp)).Set(reflect.ValueOf(result))
					return result, nil
				}
				return result, fmt.Errorf("type assertion fail")
			} else {
				return result, fmt.Errorf("callKubernetes resp is not pointer")
			}
		}
	}
	return result, nil
}

func informerGet(ctx context.Context, cluster *iapiserver.Cluster, action KubernetesGetAction, resp any) (any, error) {
	informerController, err := getInformerController(cluster)
	if err != nil {
		return nil, err
	}

	result, err := action(informerController)
	if err != nil {
		return nil, err
	}

	if result != nil {
		libkubernetes.HideManagerField(result)
		libkubernetes.InjectTypeMeta(result)
		if resp != nil {
			rt := reflect.TypeOf(resp)
			if rt.Kind() == reflect.Ptr && rt.Elem().Kind() == reflect.Ptr {
				if rt.Elem() == reflect.TypeOf(result) {
					reflect.Indirect(reflect.ValueOf(resp)).Set(reflect.ValueOf(result))
					return result, nil
				}
				return result, errors.Errorf("type assertion fail")
			} else {
				return result, errors.Errorf("callKubernetes resp is not pointer")
			}
		}
	}
	return result, nil
}

func PodCreate(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res *v1.Pod, opts metav1.CreateOptions) (*v1.Pod, error) {
	return libkubernetes.PodCreate(ctx, cluster.Config, namespace, res, opts)
}

func PodDelete(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res *v1.Pod, opts metav1.DeleteOptions) error {
	return libkubernetes.PodDelete(ctx, cluster.Config, namespace, res.Name, opts)
}

func PodDeleteCollection(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	delOpts metav1.DeleteOptions,
	listOpts metav1.ListOptions,
) error {
	return libkubernetes.PodDeleteCollection(ctx, cluster.Config, namespace, delOpts, listOpts)
}

func PodUpdate(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res *v1.Pod, opts metav1.UpdateOptions) (*v1.Pod, error) {
	return libkubernetes.PodUpdate(ctx, cluster.Config, namespace, res, opts)
}

func PodLogList(ctx context.Context, cluster *iapiserver.Cluster, namespace string, name string, opts v1.PodLogOptions) ([]ikubernetes.PodLog, error) {
	return libkubernetes.PodLogList(ctx, cluster.Config, namespace, name, opts)
}
func PodLogStream(ctx context.Context, cluster *iapiserver.Cluster, namespace string, name string, opts v1.PodLogOptions) (io.ReadCloser, error) {
	return libkubernetes.PodLogStream(ctx, cluster.Config, namespace, name, opts)
}

func PodGet(ctx context.Context, cluster *iapiserver.Cluster, namespace string, name string, opts metav1.GetOptions) (*v1.Pod, error) {
	if !informerEnabled {
		return libkubernetes.PodGet(ctx, cluster.Config, namespace, name, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.Pod, error) {
		return c.podInformer.Lister().Pods(namespace).Get(name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func PodList(ctx context.Context, cluster *iapiserver.Cluster, namespace string, opts metav1.ListOptions) (*v1.PodList, error) {
	//field selector only enable in server side.
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.PodList(ctx, cluster.Config, namespace, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.PodList, error) {
		list, err := c.podInformer.Lister().Pods(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &v1.PodList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func PodEvict(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res *v1.Pod, opts metav1.DeleteOptions) error {
	return libkubernetes.PodEvict(ctx, cluster.Config, namespace, res.Name, opts)
}

func ServiceCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.Service,
	opts metav1.CreateOptions,
) (*v1.Service, error) {
	return libkubernetes.ServiceCreate(ctx, cluster.Config, namespace, res, opts)
}

func ServiceDelete(ctx context.Context, cluster *iapiserver.Cluster, namespace string, Name string, opts metav1.DeleteOptions) error {
	return libkubernetes.ServiceDelete(ctx, cluster.Config, namespace, Name, opts)
}

func ServiceUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.Service,
	opts metav1.UpdateOptions,
) (*v1.Service, error) {
	return libkubernetes.ServiceUpdate(ctx, cluster.Config, namespace, res, opts)
}

func ServiceGet(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res string, opts metav1.GetOptions) (*v1.Service, error) {
	if !informerEnabled {
		return libkubernetes.ServiceGet(ctx, cluster.Config, namespace, res, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.Service, error) {
		return c.serviceInformer.Lister().Services(namespace).Get(res)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func ServiceList(ctx context.Context, cluster *iapiserver.Cluster, namespace string, opts metav1.ListOptions) (*v1.ServiceList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.ServiceList(ctx, cluster.Config, namespace, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.ServiceList, error) {
		list, err := c.serviceInformer.Lister().Services(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &v1.ServiceList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func SecretCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.Secret,
	opts metav1.CreateOptions,
) (*v1.Secret, error) {
	return libkubernetes.SecretCreate(ctx, cluster.Config, namespace, res, opts)
}

func SecretDelete(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res *v1.Secret, opts metav1.DeleteOptions) error {
	return libkubernetes.SecretDelete(ctx, cluster.Config, namespace, res.Name, opts)
}

func SecretUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.Secret,
	opts metav1.UpdateOptions,
) (*v1.Secret, error) {
	return libkubernetes.SecretUpdate(ctx, cluster.Config, namespace, res, opts)
}

func SecretGet(ctx context.Context, cluster *iapiserver.Cluster, namespace string, name string, opts metav1.GetOptions) (*v1.Secret, error) {
	if !informerEnabled {
		return libkubernetes.SecretGet(ctx, cluster.Config, namespace, name, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.Secret, error) {
		return c.secretInformer.Lister().Secrets(namespace).Get(name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func SecretList(ctx context.Context, cluster *iapiserver.Cluster, namespace string, opts metav1.ListOptions) (*v1.SecretList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.SecretList(ctx, cluster.Config, namespace, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.SecretList, error) {
		list, err := c.secretInformer.Lister().Secrets(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &v1.SecretList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ConfigMapCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.ConfigMap,
	opts metav1.CreateOptions,
) (*v1.ConfigMap, error) {
	return libkubernetes.ConfigMapCreate(ctx, cluster.Config, namespace, res, opts)
}

func ConfigMapDelete(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res *v1.ConfigMap, opts metav1.DeleteOptions) error {
	return libkubernetes.ConfigMapDelete(ctx, cluster.Config, namespace, res.Name, opts)
}

func ConfigMapUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.ConfigMap,
	opts metav1.UpdateOptions,
) (*v1.ConfigMap, error) {
	return libkubernetes.ConfigMapUpdate(ctx, cluster.Config, namespace, res, opts)
}

func ServiceAccountCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.ServiceAccount,
	opts metav1.CreateOptions,
) (*v1.ServiceAccount, error) {
	return libkubernetes.ServiceAccountCreate(ctx, cluster.Config, namespace, res, opts)
}

func ServiceAccountDelete(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.ServiceAccount,
	opts metav1.DeleteOptions,
) error {
	return libkubernetes.ServiceAccountDelete(ctx, cluster.Config, namespace, res.Name, opts)
}

func ServiceAccountUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.ServiceAccount,
	opts metav1.UpdateOptions,
) (*v1.ServiceAccount, error) {
	return libkubernetes.ServiceAccountUpdate(ctx, cluster.Config, namespace, res, opts)
}

func ServiceAccountGet(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.ServiceAccount,
	opts metav1.GetOptions,
) (*v1.ServiceAccount, error) {
	if !informerEnabled {
		return libkubernetes.ServiceAccountGet(ctx, cluster.Config, namespace, res.Name, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.ServiceAccount, error) {
		return c.serviceAccountInformer.Lister().ServiceAccounts(namespace).Get(res.Name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func ServiceAccountList(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	opts metav1.ListOptions,
) (*v1.ServiceAccountList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.ServiceAccountList(ctx, cluster.Config, namespace, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.ServiceAccountList, error) {
		list, err := c.serviceAccountInformer.Lister().ServiceAccounts(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &v1.ServiceAccountList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ConfigMapGet(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*v1.ConfigMap, error) {
	if !informerEnabled {
		return libkubernetes.ConfigMapGet(ctx, cluster.Config, namespace, name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.ConfigMap, error) {
		return c.configMapInformer.Lister().ConfigMaps(namespace).Get(name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func ConfigMapList(ctx context.Context, cluster *iapiserver.Cluster, namespace string, opts metav1.ListOptions) (*v1.ConfigMapList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.ConfigMapList(ctx, cluster.Config, namespace, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.ConfigMapList, error) {
		list, err := c.configMapInformer.Lister().ConfigMaps(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &v1.ConfigMapList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NamespaceCreate(ctx context.Context, cluster *iapiserver.Cluster, res *v1.Namespace, opts metav1.CreateOptions) (*v1.Namespace, error) {
	return libkubernetes.NamespaceCreate(ctx, cluster.Config, res, opts)
}

func NamespaceDelete(ctx context.Context, cluster *iapiserver.Cluster, res *v1.Namespace, opts metav1.DeleteOptions) error {
	return libkubernetes.NamespaceDelete(ctx, cluster.Config, res.Name, opts)
}

func NamespaceUpdate(ctx context.Context, cluster *iapiserver.Cluster, res *v1.Namespace, opts metav1.UpdateOptions) (*v1.Namespace, error) {
	return libkubernetes.NamespaceUpdate(ctx, cluster.Config, res, opts)
}

func NamespaceGet(ctx context.Context, cluster *iapiserver.Cluster, name string, opts metav1.GetOptions) (*v1.Namespace, error) {
	if !informerEnabled {
		return libkubernetes.NamespaceGet(ctx, cluster.Config, name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.Namespace, error) {
		return c.namespaceInformer.Lister().Get(name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func NamespaceList(ctx context.Context, cluster *iapiserver.Cluster, opts metav1.ListOptions) (*v1.NamespaceList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.NamespaceList(ctx, cluster.Config, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.NamespaceList, error) {
		list, err := c.namespaceInformer.Lister().List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &v1.NamespaceList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func LimitRangeCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.LimitRange,
	opts metav1.CreateOptions,
) (*v1.LimitRange, error) {
	return libkubernetes.LimitRangeCreate(ctx, cluster.Config, namespace, res, opts)
}

func LimitRangeDelete(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res *v1.LimitRange, opts metav1.DeleteOptions) error {
	return libkubernetes.LimitRangeDelete(ctx, cluster.Config, namespace, res.Name, opts)
}

func LimitRangeUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.LimitRange,
	opts metav1.UpdateOptions,
) (*v1.LimitRange, error) {
	return libkubernetes.LimitRangeUpdate(ctx, cluster.Config, namespace, res, opts)
}

func LimitRangeGet(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*v1.LimitRange, error) {
	if !informerEnabled {
		return libkubernetes.LimitRangeGet(ctx, cluster.Config, namespace, name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.LimitRange, error) {
		return c.limitRangeInformer.Lister().LimitRanges(namespace).Get(name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func LimitRangeList(ctx context.Context, cluster *iapiserver.Cluster, namespace string, opts metav1.ListOptions) (*v1.LimitRangeList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.LimitRangeList(ctx, cluster.Config, namespace, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.LimitRangeList, error) {
		list, err := c.limitRangeInformer.Lister().LimitRanges(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &v1.LimitRangeList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ResourceQuotaCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.ResourceQuota,
	opts metav1.CreateOptions,
) (*v1.ResourceQuota, error) {
	return libkubernetes.ResourceQuotaCreate(ctx, cluster.Config, namespace, res, opts)
}

func ResourceQuotaDelete(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.ResourceQuota,
	opts metav1.DeleteOptions,
) error {
	return libkubernetes.ResourceQuotaDelete(ctx, cluster.Config, namespace, res.Name, opts)
}

func ResourceQuotaUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.ResourceQuota,
	opts metav1.UpdateOptions,
) (*v1.ResourceQuota, error) {
	return libkubernetes.ResourceQuotaUpdate(ctx, cluster.Config, namespace, res, opts)
}

func ResourceQuotaGet(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*v1.ResourceQuota, error) {
	if !informerEnabled {
		return libkubernetes.ResourceQuotaGet(ctx, cluster.Config, namespace, name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.ResourceQuota, error) {
		return c.resourceQuotaInformer.Lister().ResourceQuotas(namespace).Get(name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func ResourceQuotaList(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	opts metav1.ListOptions,
) (*v1.ResourceQuotaList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.ResourceQuotaList(ctx, cluster.Config, namespace, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.ResourceQuotaList, error) {
		list, err := c.resourceQuotaInformer.Lister().ResourceQuotas(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &v1.ResourceQuotaList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func PersistentVolumeClaimCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.PersistentVolumeClaim,
	opts metav1.CreateOptions,
) (*v1.PersistentVolumeClaim, error) {
	return libkubernetes.PersistentVolumeClaimCreate(ctx, cluster.Config, namespace, res, opts)
}

func PersistentVolumeClaimDelete(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	name string,
	opts metav1.DeleteOptions,
) error {
	return libkubernetes.PersistentVolumeClaimDelete(ctx, cluster.Config, namespace, name, opts)
}

func PersistentVolumeClaimUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.PersistentVolumeClaim,
	opts metav1.UpdateOptions,
) (*v1.PersistentVolumeClaim, error) {
	return libkubernetes.PersistentVolumeClaimUpdate(ctx, cluster.Config, namespace, res, opts)
}

func PersistentVolumeClaimGet(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*v1.PersistentVolumeClaim, error) {
	if !informerEnabled {
		return libkubernetes.PersistentVolumeClaimGet(ctx, cluster.Config, namespace, name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.PersistentVolumeClaim, error) {
		return c.persistentVolumeClaimInformer.Lister().PersistentVolumeClaims(namespace).Get(name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func PersistentVolumeClaimList(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	opts metav1.ListOptions,
) (*v1.PersistentVolumeClaimList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.PersistentVolumeClaimList(ctx, cluster.Config, namespace, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.PersistentVolumeClaimList, error) {
		list, err := c.persistentVolumeClaimInformer.Lister().PersistentVolumeClaims(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &v1.PersistentVolumeClaimList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func PersistentVolumeCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.PersistentVolume,
	opts metav1.CreateOptions,
) (*v1.PersistentVolume, error) {
	return libkubernetes.PersistentVolumeCreate(ctx, cluster.Config, namespace, res, opts)
}

func PersistentVolumeDelete(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	name string,
	opts metav1.DeleteOptions,
) error {
	return libkubernetes.PersistentVolumeDelete(ctx, cluster.Config, namespace, name, opts)
}

func PersistentVolumeUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.PersistentVolume,
	opts metav1.UpdateOptions,
) (*v1.PersistentVolume, error) {
	return libkubernetes.PersistentVolumeUpdate(ctx, cluster.Config, namespace, res, opts)
}

func PersistentVolumeGet(ctx context.Context, cluster *iapiserver.Cluster, name string, opts metav1.GetOptions) (*v1.PersistentVolume, error) {
	if !informerEnabled {
		return libkubernetes.PersistentVolumeGet(ctx, cluster.Config, name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.PersistentVolume, error) {
		return c.persistentVolumeInformer.Lister().Get(name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func PersistentVolumeList(ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	opts metav1.ListOptions,
) (*v1.PersistentVolumeList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.PersistentVolumeList(ctx, cluster.Config, namespace, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.PersistentVolumeList, error) {
		list, err := c.persistentVolumeInformer.Lister().List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &v1.PersistentVolumeList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func NodeUpdate(ctx context.Context, cluster *iapiserver.Cluster, res *v1.Node, opts metav1.UpdateOptions) (*v1.Node, error) {
	return libkubernetes.NodeUpdate(ctx, cluster.Config, res, opts)
}

func NodeCreate(ctx context.Context, cluster *iapiserver.Cluster, res *v1.Node, opts metav1.CreateOptions) (*v1.Node, error) {
	return libkubernetes.NodeAdd(ctx, cluster.Config, res, opts)
}

func NodeCordon(ctx context.Context, cluster *iapiserver.Cluster, res *v1.Node) (*v1.Node, error) {
	return libkubernetes.NodeCordon(ctx, cluster.Config, res.Name)
}

func NodeUncordon(ctx context.Context, cluster *iapiserver.Cluster, res *v1.Node) (*v1.Node, error) {
	return libkubernetes.NodeUncordon(ctx, cluster.Config, res.Name)
}

func NodeDelete(ctx context.Context, cluster *iapiserver.Cluster, res *v1.Node, opts metav1.DeleteOptions) error {
	return libkubernetes.NodeDelete(ctx, cluster.Config, res.Name, opts)
}
func NodeGet(ctx context.Context, cluster *iapiserver.Cluster, name string, opts metav1.GetOptions) (*v1.Node, error) {
	if !informerEnabled {
		return libkubernetes.NodeGet(ctx, cluster.Config, name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.Node, error) {
		return c.nodeInformer.Lister().Get(name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func NodeList(ctx context.Context, cluster *iapiserver.Cluster, opts metav1.ListOptions) (*v1.NodeList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.NodeList(ctx, cluster.Config, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.NodeList, error) {
		list, err := c.nodeInformer.Lister().List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &v1.NodeList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func EndpointsCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.Endpoints,
	opts metav1.CreateOptions,
) (*v1.Endpoints, error) {
	return libkubernetes.EndpointCreate(ctx, cluster.Config, namespace, res, opts)
}

func EndpointsDelete(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res *v1.Endpoints, opts metav1.DeleteOptions) error {
	return libkubernetes.EndpointDelete(ctx, cluster.Config, namespace, res.Name, opts)
}

func EndpointsUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.Endpoints,
	opts metav1.UpdateOptions,
) (*v1.Endpoints, error) {
	return libkubernetes.EndpointsUpdate(ctx, cluster.Config, namespace, res, opts)
}

func EndpointsGet(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*v1.Endpoints, error) {
	if !informerEnabled {
		return libkubernetes.EndpointGet(ctx, cluster.Config, namespace, name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.Endpoints, error) {
		return c.endpointsInformer.Lister().Endpoints(namespace).Get(name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func EndpointsList(ctx context.Context, cluster *iapiserver.Cluster, namespace string, opts metav1.ListOptions) (*v1.EndpointsList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.EndpointList(ctx, cluster.Config, namespace, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.EndpointsList, error) {
		list, err := c.endpointsInformer.Lister().Endpoints(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &v1.EndpointsList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func EventCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *v1.Event,
	opts metav1.CreateOptions,
) (*v1.Event, error) {
	return libkubernetes.EventCreate(ctx, cluster.Config, namespace, res, opts)
}

func EventSearch(ctx context.Context, cluster *iapiserver.Cluster, scheme *runtime.Scheme, obj runtime.Object) (*v1.EventList, error) {
	return libkubernetes.EventSearch(ctx, cluster.Config, scheme, obj)
}

func EventGet(ctx context.Context, cluster *iapiserver.Cluster, namespace string, name string, opts metav1.GetOptions) (*v1.Event, error) {
	if !informerEnabled {
		return libkubernetes.EventGet(ctx, cluster.Config, namespace, name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.Event, error) {
		return c.eventInformer.Lister().Events(namespace).Get(name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func EventList(ctx context.Context, cluster *iapiserver.Cluster, namespace string, opts metav1.ListOptions) (*v1.EventList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.EventList(ctx, cluster.Config, namespace, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*v1.EventList, error) {
		list, err := c.eventInformer.Lister().Events(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &v1.EventList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func CronJobCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *batchv1.CronJob,
	opts metav1.CreateOptions,
) (*batchv1.CronJob, error) {
	return libkubernetes.CronJobCreate(ctx, cluster.Config, namespace, res, opts)
}

func CronJobDelete(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *batchv1.CronJob,
	opts metav1.DeleteOptions,
) error {
	return libkubernetes.CronJobDelete(ctx, cluster.Config, namespace, res.Name, opts)
}

func CronJobUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *batchv1.CronJob,
	opts metav1.UpdateOptions,
) (*batchv1.CronJob, error) {
	return libkubernetes.CronJobUpdate(ctx, cluster.Config, namespace, res, opts)
}

func CronJobGet(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*batchv1.CronJob, error) {
	if !informerEnabled {
		return libkubernetes.CronJobGet(ctx, cluster.Config, namespace, name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*batchv1.CronJob, error) {
		return c.cronJobInformer.Lister().CronJobs(namespace).Get(name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func CronJobList(ctx context.Context, cluster *iapiserver.Cluster, namespace string, opts metav1.ListOptions) (*batchv1.CronJobList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.CronJobList(ctx, cluster.Config, namespace, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*batchv1.CronJobList, error) {
		list, err := c.cronJobInformer.Lister().CronJobs(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &batchv1.CronJobList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func JobCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *batchv1.Job,
	opts metav1.CreateOptions,
) (*batchv1.Job, error) {
	return libkubernetes.JobCreate(ctx, cluster.Config, namespace, res, opts)
}

func JobDelete(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res *batchv1.Job, opts metav1.DeleteOptions) error {
	return libkubernetes.JobDelete(ctx, cluster.Config, namespace, res.Name, opts)
}

func JobUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *batchv1.Job,
	opts metav1.UpdateOptions,
) (*batchv1.Job, error) {
	return libkubernetes.JobUpdate(ctx, cluster.Config, namespace, res, opts)
}

func JobGet(ctx context.Context, cluster *iapiserver.Cluster, namespace string, name string, opts metav1.GetOptions) (*batchv1.Job, error) {
	if !informerEnabled {
		return libkubernetes.JobGet(ctx, cluster.Config, namespace, name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*batchv1.Job, error) {
		return c.jobInformer.Lister().Jobs(namespace).Get(name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func JobList(ctx context.Context, cluster *iapiserver.Cluster, namespace string, opts metav1.ListOptions) (*batchv1.JobList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.JobList(ctx, cluster.Config, namespace, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*batchv1.JobList, error) {
		list, err := c.jobInformer.Lister().Jobs(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &batchv1.JobList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func DeploymentCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *appsv1.Deployment,
	opts metav1.CreateOptions,
) (*appsv1.Deployment, error) {
	return libkubernetes.DeploymentCreate(ctx, cluster.Config, namespace, res, opts)
}

func DeploymentDelete(ctx context.Context, cluster *iapiserver.Cluster, namespace string, resName string, opts metav1.DeleteOptions) error {
	return libkubernetes.DeploymentDelete(ctx, cluster.Config, namespace, resName, opts)
}

func DeploymentUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *appsv1.Deployment,
	opts metav1.UpdateOptions,
) (*appsv1.Deployment, error) {
	return libkubernetes.DeploymentUpdate(ctx, cluster.Config, namespace, res, opts)
}

func DeploymentPatch(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *appsv1.Deployment,
	pt k8stypes.PatchType,
	data []byte,
	opts metav1.PatchOptions,
	subresources ...string,
) (*appsv1.Deployment, error) {
	return libkubernetes.DeploymentPatch(ctx, cluster.Config, namespace, res.Name, pt, data, opts, subresources...)
}

func DeploymentGet(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res string,
	opts metav1.GetOptions,
) (*appsv1.Deployment, error) {
	if !informerEnabled {
		return libkubernetes.DeploymentGet(ctx, cluster.Config, namespace, res, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*appsv1.Deployment, error) {
		return c.deploymentInformer.Lister().Deployments(namespace).Get(res)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func DeploymentList(ctx context.Context, cluster *iapiserver.Cluster, namespace string, opts metav1.ListOptions) (*appsv1.DeploymentList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.DeploymentList(ctx, cluster.Config, namespace, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*appsv1.DeploymentList, error) {
		list, err := c.deploymentInformer.Lister().Deployments(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &appsv1.DeploymentList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func DaemonSetCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *appsv1.DaemonSet,
	opts metav1.CreateOptions,
) (*appsv1.DaemonSet, error) {
	return libkubernetes.DaemonSetCreate(ctx, cluster.Config, namespace, res, opts)
}

func DaemonSetPatch(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *appsv1.DaemonSet,
	pt k8stypes.PatchType,
	data []byte,
	opts metav1.PatchOptions,
	subresources ...string,
) (*appsv1.DaemonSet, error) {
	return libkubernetes.DaemonSetPatch(ctx, cluster.Config, namespace, res.Name, pt, data, opts, subresources...)
}

func DaemonSetDelete(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res *appsv1.DaemonSet, opts metav1.DeleteOptions) error {
	return libkubernetes.DaemonSetDelete(ctx, cluster.Config, namespace, res.Name, opts)
}

func DaemonSetUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *appsv1.DaemonSet,
	opts metav1.UpdateOptions,
) (*appsv1.DaemonSet, error) {
	return libkubernetes.DaemonSetUpdate(ctx, cluster.Config, namespace, res, opts)
}

func DaemonSetGet(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res string,
	opts metav1.GetOptions,
) (*appsv1.DaemonSet, error) {
	if !informerEnabled {
		return libkubernetes.DaemonSetGet(ctx, cluster.Config, namespace, res, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*appsv1.DaemonSet, error) {
		return c.daemonSetInformer.Lister().DaemonSets(namespace).Get(res)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil

}

func DaemonSetList(ctx context.Context, cluster *iapiserver.Cluster, namespace string, opts metav1.ListOptions) (*appsv1.DaemonSetList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.DaemonSetList(ctx, cluster.Config, namespace, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*appsv1.DaemonSetList, error) {
		list, err := c.daemonSetInformer.Lister().DaemonSets(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &appsv1.DaemonSetList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func StatefulSetCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *appsv1.StatefulSet,
	opts metav1.CreateOptions,
) (*appsv1.StatefulSet, error) {
	return libkubernetes.StatefulSetCreate(ctx, cluster.Config, namespace, res, opts)
}

func StatefulSetDelete(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *appsv1.StatefulSet,
	opts metav1.DeleteOptions,
) error {
	return libkubernetes.StatefulSetDelete(ctx, cluster.Config, namespace, res.Name, opts)
}

func StatefulSetUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *appsv1.StatefulSet,
	opts metav1.UpdateOptions,
) (*appsv1.StatefulSet, error) {
	return libkubernetes.StatefulSetUpdate(ctx, cluster.Config, namespace, res, opts)
}

func StatefulSetPatch(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *appsv1.StatefulSet,
	pt k8stypes.PatchType,
	data []byte,
	opts metav1.PatchOptions,
	subresources ...string,
) (*appsv1.StatefulSet, error) {
	return libkubernetes.StatefulSetPatch(ctx, cluster.Config, namespace, res.Name, pt, data, opts, subresources...)
}

func StatefulSetGet(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res string,
	opts metav1.GetOptions,
) (*appsv1.StatefulSet, error) {
	if !informerEnabled {
		return libkubernetes.StatefulSetGet(ctx, cluster.Config, namespace, res, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*appsv1.StatefulSet, error) {
		return c.statefulSetInformer.Lister().StatefulSets(namespace).Get(res)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func StatefulSetList(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	opts metav1.ListOptions,
) (*appsv1.StatefulSetList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.StatefulSetList(ctx, cluster.Config, namespace, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*appsv1.StatefulSetList, error) {
		list, err := c.statefulSetInformer.Lister().StatefulSets(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &appsv1.StatefulSetList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func ReplicaSetCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *appsv1.ReplicaSet,
	opts metav1.CreateOptions,
) (*appsv1.ReplicaSet, error) {
	return libkubernetes.ReplicaSetCreate(ctx, cluster.Config, namespace, res, opts)
}

func ReplicaSetDelete(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *appsv1.ReplicaSet,
	opts metav1.DeleteOptions,
) error {
	return libkubernetes.ReplicaSetDelete(ctx, cluster.Config, namespace, res.Name, opts)
}

func ReplicaSetUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *appsv1.ReplicaSet,
	opts metav1.UpdateOptions,
) (*appsv1.ReplicaSet, error) {
	return libkubernetes.ReplicaSetUpdate(ctx, cluster.Config, namespace, res, opts)
}

func ReplicaSetGet(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*appsv1.ReplicaSet, error) {
	if !informerEnabled {
		return libkubernetes.ReplicaSetGet(ctx, cluster.Config, namespace, name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*appsv1.ReplicaSet, error) {
		return c.replicaSetInformer.Lister().ReplicaSets(namespace).Get(name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil

}

func ReplicaSetList(ctx context.Context, cluster *iapiserver.Cluster, namespace string, opts metav1.ListOptions) (*appsv1.ReplicaSetList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.ReplicaSetList(ctx, cluster.Config, namespace, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*appsv1.ReplicaSetList, error) {
		list, err := c.replicaSetInformer.Lister().ReplicaSets(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &appsv1.ReplicaSetList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func RevisionList(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	opt metav1.ListOptions,
) (*appsv1.ControllerRevisionList, error) {
	return libkubernetes.RevisionList(ctx, cluster.Config, namespace, opt)
}

func IngressCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *networkingv1.Ingress,
	opts metav1.CreateOptions,
) (*networkingv1.Ingress, error) {
	return libkubernetes.IngressCreate(ctx, cluster.Config, namespace, res, opts)
}

func IngressDelete(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *networkingv1.Ingress,
	opts metav1.DeleteOptions,
) error {
	return libkubernetes.IngressDelete(ctx, cluster.Config, namespace, res.Name, opts)
}

func IngressUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *networkingv1.Ingress,
	opts metav1.UpdateOptions,
) (*networkingv1.Ingress, error) {
	return libkubernetes.IngressUpdate(ctx, cluster.Config, namespace, res, opts)
}

func IngressGet(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*networkingv1.Ingress, error) {
	if !informerEnabled {
		return libkubernetes.IngressGet(ctx, cluster.Config, namespace, name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*networkingv1.Ingress, error) {
		return c.ingressInformer.Lister().Ingresses(namespace).Get(name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func HpaCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *autoscalingv1.HorizontalPodAutoscaler,
	opts metav1.CreateOptions,
) (*autoscalingv1.HorizontalPodAutoscaler, error) {
	return libkubernetes.HpaCreate(ctx, cluster.Config, namespace, res, opts)
}

func HpaDelete(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *autoscalingv1.HorizontalPodAutoscaler,
	opts metav1.DeleteOptions,
) error {
	return libkubernetes.HpaDelete(ctx, cluster.Config, namespace, res.Name, opts)
}

func HpaUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *autoscalingv1.HorizontalPodAutoscaler,
	opts metav1.UpdateOptions,
) (*autoscalingv1.HorizontalPodAutoscaler, error) {
	return libkubernetes.HpaUpdate(ctx, cluster.Config, namespace, res, opts)
}

func HpaGet(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*autoscalingv1.HorizontalPodAutoscaler, error) {
	if !informerEnabled {
		return libkubernetes.HpaGet(ctx, cluster.Config, namespace, name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*autoscalingv1.HorizontalPodAutoscaler, error) {
		return c.hpaInformer.Lister().HorizontalPodAutoscalers(namespace).Get(name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func HpaList(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	opts metav1.ListOptions,
) (*autoscalingv1.HorizontalPodAutoscalerList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.HpaList(ctx, cluster.Config, namespace, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*autoscalingv1.HorizontalPodAutoscalerList, error) {
		list, err := c.hpaInformer.Lister().HorizontalPodAutoscalers(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &autoscalingv1.HorizontalPodAutoscalerList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func IngressList(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	opts metav1.ListOptions,
) (*networkingv1.IngressList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.IngressList(ctx, cluster.Config, namespace, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*networkingv1.IngressList, error) {
		list, err := c.ingressInformer.Lister().Ingresses(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &networkingv1.IngressList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func StorageClassCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	res *storagev1.StorageClass,
	opts metav1.CreateOptions,
) (*storagev1.StorageClass, error) {
	return libkubernetes.StorageClassCreate(ctx, cluster.Config, res, opts)
}

func StorageClassDelete(ctx context.Context, cluster *iapiserver.Cluster, res *storagev1.StorageClass, opts metav1.DeleteOptions) error {
	return libkubernetes.StorageClassDelete(ctx, cluster.Config, res.Name, opts)
}

func StorageClassUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *storagev1.StorageClass,
	opts metav1.UpdateOptions,
) (*storagev1.StorageClass, error) {
	return libkubernetes.StorageClassUpdate(ctx, cluster.Config, namespace, res, opts)
}

func StorageClassGet(ctx context.Context, cluster *iapiserver.Cluster, name string, opts metav1.GetOptions) (*storagev1.StorageClass, error) {
	if !informerEnabled {
		return libkubernetes.StorageClassGet(ctx, cluster.Config, name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*storagev1.StorageClass, error) {
		return c.storageClassInformer.Lister().Get(name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func StorageClassList(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	opts metav1.ListOptions,
) (*storagev1.StorageClassList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.StorageClassList(ctx, cluster.Config, namespace, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*storagev1.StorageClassList, error) {
		list, err := c.storageClassInformer.Lister().List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &storagev1.StorageClassList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func NetworkPolicyCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *networkingv1.NetworkPolicy,
	opts metav1.CreateOptions,
) (*networkingv1.NetworkPolicy, error) {
	return libkubernetes.NetworkPolicyCreate(ctx, cluster.Config, namespace, res, opts)
}

func NetworkPolicyDelete(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *networkingv1.NetworkPolicy,
	opts metav1.DeleteOptions,
) error {
	return libkubernetes.NetworkPolicyDelete(ctx, cluster.Config, namespace, res, opts)
}

func NetworkPolicyUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *networkingv1.NetworkPolicy,
	opts metav1.UpdateOptions,
) (*networkingv1.NetworkPolicy, error) {
	return libkubernetes.NetworkPolicyUpdate(ctx, cluster.Config, namespace, res, opts)
}

func NetworkPolicyGet(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*networkingv1.NetworkPolicy, error) {
	if !informerEnabled {
		return libkubernetes.NetworkPolicyGet(ctx, cluster.Config, namespace, name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*networkingv1.NetworkPolicy, error) {
		return c.networkingInformer.Lister().NetworkPolicies(namespace).Get(name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func NetworkPolicyList(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	opts metav1.ListOptions,
) (*networkingv1.NetworkPolicyList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.NetworkPolicyList(ctx, cluster.Config, namespace, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*networkingv1.NetworkPolicyList, error) {
		list, err := c.networkingInformer.Lister().NetworkPolicies(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &networkingv1.NetworkPolicyList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func RoleCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *rbacv1.Role,
	opts metav1.CreateOptions,
) (*rbacv1.Role, error) {
	return libkubernetes.RoleCreate(ctx, cluster.Config, namespace, res, opts)
}

func RoleDelete(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res *rbacv1.Role, opts metav1.DeleteOptions) error {
	return libkubernetes.RoleDelete(ctx, cluster.Config, namespace, res.Name, opts)
}

func RoleUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *rbacv1.Role,
	opts metav1.UpdateOptions,
) (*rbacv1.Role, error) {
	return libkubernetes.RoleUpdate(ctx, cluster.Config, namespace, res, opts)
}

func RoleGet(ctx context.Context, cluster *iapiserver.Cluster, namespace string, res *rbacv1.Role, opts metav1.GetOptions) (*rbacv1.Role, error) {
	if !informerEnabled {
		return libkubernetes.RoleGet(ctx, cluster.Config, namespace, res.Name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*rbacv1.Role, error) {
		return c.roleInformer.Lister().Roles(namespace).Get(res.Name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func RoleList(ctx context.Context, cluster *iapiserver.Cluster, namespace string, opts metav1.ListOptions) (*rbacv1.RoleList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.RoleList(ctx, cluster.Config, namespace, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*rbacv1.RoleList, error) {
		list, err := c.roleInformer.Lister().Roles(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &rbacv1.RoleList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func RoleBindingCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *rbacv1.RoleBinding,
	opts metav1.CreateOptions,
) (*rbacv1.RoleBinding, error) {
	return libkubernetes.RoleBindingCreate(ctx, cluster.Config, namespace, res, opts)
}

func RoleBindingDelete(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *rbacv1.RoleBinding,
	opts metav1.DeleteOptions,
) error {
	return libkubernetes.RoleBindingDelete(ctx, cluster.Config, namespace, res.Name, opts)
}

func RoleBindingUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *rbacv1.RoleBinding,
	opts metav1.UpdateOptions,
) (*rbacv1.RoleBinding, error) {
	return libkubernetes.RoleBindingUpdate(ctx, cluster.Config, namespace, res, opts)
}

func RoleBindingGet(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	res *rbacv1.RoleBinding,
	opts metav1.GetOptions,
) (*rbacv1.RoleBinding, error) {
	if !informerEnabled {
		return libkubernetes.RoleBindingGet(ctx, cluster.Config, namespace, res.Name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*rbacv1.RoleBinding, error) {
		return c.rolebindingInformer.Lister().RoleBindings(namespace).Get(res.Name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func RoleBindingList(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	opts metav1.ListOptions,
) (*rbacv1.RoleBindingList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.RoleBindingList(ctx, cluster.Config, namespace, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*rbacv1.RoleBindingList, error) {
		list, err := c.rolebindingInformer.Lister().RoleBindings(namespace).List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &rbacv1.RoleBindingList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func ClusterRoleCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	res *rbacv1.ClusterRole,
	opts metav1.CreateOptions,
) (*rbacv1.ClusterRole, error) {
	return libkubernetes.ClusterRoleCreate(ctx, cluster.Config, res, opts)
}

func ClusterRoleDelete(ctx context.Context, cluster *iapiserver.Cluster, res *rbacv1.ClusterRole, opts metav1.DeleteOptions) error {
	return libkubernetes.ClusterRoleDelete(ctx, cluster.Config, res.Name, opts)
}

func ClusterRoleUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	res *rbacv1.ClusterRole,
	opts metav1.UpdateOptions,
) (*rbacv1.ClusterRole, error) {
	return libkubernetes.ClusterRoleUpdate(ctx, cluster.Config, res, opts)
}

func ClusterRoleGet(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	res *rbacv1.ClusterRole,
	opts metav1.GetOptions,
) (*rbacv1.ClusterRole, error) {
	if !informerEnabled {
		return libkubernetes.ClusterRoleGet(ctx, cluster.Config, res.Name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*rbacv1.ClusterRole, error) {
		return c.clusterroleInformer.Lister().Get(res.Name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func ClusterRoleList(ctx context.Context, cluster *iapiserver.Cluster, opts metav1.ListOptions) (*rbacv1.ClusterRoleList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.ClusterRoleList(ctx, cluster.Config, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*rbacv1.ClusterRoleList, error) {
		list, err := c.clusterroleInformer.Lister().List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &rbacv1.ClusterRoleList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i])
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func ClusterRoleBindingCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	res *rbacv1.ClusterRoleBinding,
	opts metav1.CreateOptions,
) (*rbacv1.ClusterRoleBinding, error) {
	return libkubernetes.ClusterRoleBindingCreate(ctx, cluster.Config, res, opts)
}

func ClusterRoleBindingDelete(ctx context.Context, cluster *iapiserver.Cluster, res *rbacv1.ClusterRoleBinding, opts metav1.DeleteOptions) error {
	return libkubernetes.ClusterRoleBindingDelete(ctx, cluster.Config, res.Name, opts)
}

func ClusterRoleBindingUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	res *rbacv1.ClusterRoleBinding,
	opts metav1.UpdateOptions,
) (*rbacv1.ClusterRoleBinding, error) {
	return libkubernetes.ClusterRoleBindingUpdate(ctx, cluster.Config, res, opts)
}

func ClusterRoleBindingGet(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	res *rbacv1.ClusterRoleBinding,
	opts metav1.GetOptions,
) (*rbacv1.ClusterRoleBinding, error) {
	if !informerEnabled {
		return libkubernetes.ClusterRoleBindingGet(ctx, cluster.Config, res.Name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*rbacv1.ClusterRoleBinding, error) {
		return c.clusterrolebindingInformer.Lister().Get(res.Name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func ClusterRoleBindingList(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	opts metav1.ListOptions,
) (*rbacv1.ClusterRoleBindingList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.ClusterRoleBindingList(ctx, cluster.Config, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*rbacv1.ClusterRoleBindingList, error) {
		list, err := c.clusterrolebindingInformer.Lister().List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &rbacv1.ClusterRoleBindingList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func PodDisruptionBudgetCreate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	res *policyv1.PodDisruptionBudget,
	opts metav1.CreateOptions,
) (*policyv1.PodDisruptionBudget, error) {
	return libkubernetes.PodDisruptionBudgetCreate(ctx, cluster.Config, res.Namespace, res, opts)
}

func PodDisruptionBudgetDelete(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	res *policyv1.PodDisruptionBudget,
	opts metav1.DeleteOptions,
) error {
	return libkubernetes.PodDisruptionBudgetDelete(ctx, cluster.Config, res.Namespace, res.Name, opts)
}

func PodDisruptionBudgetUpdate(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	res *policyv1.PodDisruptionBudget,
	opts metav1.UpdateOptions,
) (*policyv1.PodDisruptionBudget, error) {
	return libkubernetes.PodDisruptionBudgetUpdate(ctx, cluster.Config, res.Namespace, res, opts)
}

func PodDisruptionBudgetGet(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	res *policyv1.PodDisruptionBudget,
	opts metav1.GetOptions,
) (*policyv1.PodDisruptionBudget, error) {
	if !informerEnabled {
		return libkubernetes.PodDisruptionBudgetGet(ctx, cluster.Config, res.Namespace, res.Name, opts)
	}
	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*policyv1.PodDisruptionBudget, error) {
		return c.podDisruptionBudegetInformer.Lister().PodDisruptionBudgets(res.Namespace).Get(res.Name)
	}, nil)
	if err != nil {
		return nil, err
	}
	return resp.DeepCopy(), nil
}

func PodDisruptionBudgetList(
	ctx context.Context,
	cluster *iapiserver.Cluster,
	namespace string,
	opts metav1.ListOptions,
) (*policyv1.PodDisruptionBudgetList, error) {
	if !informerEnabled || opts.FieldSelector != "" {
		return libkubernetes.PodDisruptionBudgetList(ctx, cluster.Config, namespace, opts)
	}

	resp, err := informerAction(ctx, cluster, func(ctx context.Context, c *ResourceCacheController, labelSet labels.Set) (*policyv1.PodDisruptionBudgetList, error) {
		list, err := c.podDisruptionBudegetInformer.Lister().List(labelSet.AsSelector())
		if err != nil {
			return nil, err
		}
		resList := &policyv1.PodDisruptionBudgetList{}
		for i := range list {
			resList.Items = append(resList.Items, *list[i].DeepCopy())
		}
		return resList, nil
	}, &opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
