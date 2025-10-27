package libkubernetes

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes"

	"github.com/wangweihong/eazycloud/apis/ikubernetes"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func PersistentVolumeClaimCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	pvc *v1.PersistentVolumeClaim,
	opts metav1.CreateOptions,
) (*v1.PersistentVolumeClaim, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.PersistentVolumeClaim, error) {
		resource, err := c.client.CoreV1().PersistentVolumeClaims(namespace).Create(ctx, pvc, opts)
		return resource, err
	})
	return resp, err
}

func PersistentVolumeClaimDelete(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.DeleteOptions,
) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.PersistentVolumeClaim, error) {
		err := c.client.CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, name, opts)
		return nil, err
	})
	return err
}

func PersistentVolumeClaimUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	pvc *v1.PersistentVolumeClaim,
	opts metav1.UpdateOptions,
) (*v1.PersistentVolumeClaim, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.PersistentVolumeClaim, error) {
		resource, err := c.client.CoreV1().PersistentVolumeClaims(namespace).Update(ctx, pvc, opts)
		return resource, err
	})
	return resp, err
}

func PersistentVolumeClaimGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*v1.PersistentVolumeClaim, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.PersistentVolumeClaim, error) {
		resource, err := c.client.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func PersistentVolumeClaimList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*v1.PersistentVolumeClaimList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.PersistentVolumeClaimList, error) {
		resource, err := c.client.CoreV1().PersistentVolumeClaims(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err
}

func ConfigMapCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	configMap *v1.ConfigMap,
	opts metav1.CreateOptions,
) (*v1.ConfigMap, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ConfigMap, error) {
		resource, err := c.client.CoreV1().ConfigMaps(namespace).Create(ctx, configMap, opts)
		return resource, err
	})
	return resp, err
}

func ConfigMapUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	configMap *v1.ConfigMap,
	opts metav1.UpdateOptions,
) (*v1.ConfigMap, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ConfigMap, error) {
		resource, err := c.client.CoreV1().ConfigMaps(namespace).Update(ctx, configMap, opts)
		return resource, err
	})
	return resp, err
}

func ConfigMapDelete(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ConfigMap, error) {
		err := c.client.CoreV1().ConfigMaps(namespace).Delete(ctx, name, opts)
		return nil, err
	})
	return err
}

func ConfigMapGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*v1.ConfigMap, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ConfigMap, error) {
		resource, err := c.client.CoreV1().ConfigMaps(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func ConfigMapList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*v1.ConfigMapList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ConfigMapList, error) {
		resource, err := c.client.CoreV1().ConfigMaps(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err
}

func SecretCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	secret *v1.Secret,
	opts metav1.CreateOptions,
) (*v1.Secret, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Secret, error) {
		resource, err := c.client.CoreV1().Secrets(namespace).Create(ctx, secret, opts)
		return resource, err
	})
	return resp, err
}

func SecretUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	secret *v1.Secret,
	opts metav1.UpdateOptions,
) (*v1.Secret, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Secret, error) {
		resource, err := c.client.CoreV1().Secrets(namespace).Update(ctx, secret, opts)
		return resource, err
	})
	return resp, err
}

func SecretDelete(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Secret, error) {
		err := c.client.CoreV1().Secrets(namespace).Delete(ctx, name, opts)
		return nil, err
	})
	return err
}

func SecretGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*v1.Secret, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Secret, error) {
		resource, err := c.client.CoreV1().Secrets(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func SecretList(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, opts metav1.ListOptions) (*v1.SecretList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.SecretList, error) {
		resource, err := c.client.CoreV1().Secrets(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err

}

func EventList(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, opts metav1.ListOptions) (*v1.EventList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.EventList, error) {
		resource, err := c.client.CoreV1().Events(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err

}

func EventGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*v1.Event, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Event, error) {
		resource, err := c.client.CoreV1().Events(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err

}

func EventSearch(pctx context.Context, config *ikubernetes.ClusterConfig, scheme *runtime.Scheme, obj runtime.Object) (*v1.EventList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.EventList, error) {
		resource, err := c.client.CoreV1().Events("").Search(scheme, obj)
		return resource, err
	})
	return resp, err

}

func EventCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	event *v1.Event,
	opts metav1.CreateOptions,
) (*v1.Event, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Event, error) {
		resource, err := c.client.CoreV1().Events(namespace).Create(ctx, event, opts)

		return resource, err
	})
	return resp, err
}

func ServiceCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	service *v1.Service,
	opts metav1.CreateOptions,
) (*v1.Service, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Service, error) {
		resource, err := c.client.CoreV1().Services(namespace).Create(ctx, service, opts)

		return resource, err
	})
	return resp, err

}

func ServiceDelete(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	service string,
	opts metav1.DeleteOptions,
) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Service, error) {
		err := c.client.CoreV1().Services(namespace).Delete(ctx, service, opts)

		return nil, err
	})
	return err

}

func ServiceGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	service string,
	opts metav1.GetOptions,
) (*v1.Service, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Service, error) {
		resource, err := c.client.CoreV1().Services(namespace).Get(ctx, service, opts)

		return resource, err
	})
	return resp, err
}

func ServiceList(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, opts metav1.ListOptions) (*v1.ServiceList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ServiceList, error) {
		resource, err := c.client.CoreV1().Services(namespace).List(ctx, opts)

		return resource, err
	})
	return resp, err
}

func ServiceUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	service *v1.Service,
	opts metav1.UpdateOptions,
) (*v1.Service, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Service, error) {
		resource, err := c.client.CoreV1().Services(namespace).Update(ctx, service, opts)

		return resource, err
	})
	return resp, err
}

func ServicePatch(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	service string,
	pt types.PatchType,
	data []byte,
	opts metav1.PatchOptions,
	subResources []string,
) (*v1.Service, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Service, error) {
		resource, err := c.client.CoreV1().Services(namespace).Patch(ctx, service, pt, data, opts, subResources...)

		return resource, err
	})
	return resp, err
}

// node management
func NodeAdd(pctx context.Context, config *ikubernetes.ClusterConfig, node *v1.Node, opts metav1.CreateOptions) (*v1.Node, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Node, error) {
		resource, err := c.client.CoreV1().Nodes().Create(ctx, node, opts)

		return resource, err
	})
	return resp, err

}

func NodeUpdate(pctx context.Context, config *ikubernetes.ClusterConfig, node *v1.Node, opts metav1.UpdateOptions) (*v1.Node, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Node, error) {
		resource, err := c.client.CoreV1().Nodes().Update(ctx, node, opts)

		return resource, err
	})
	return resp, err

}

func NodeDelete(pctx context.Context, config *ikubernetes.ClusterConfig, name string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Node, error) {
		err := c.client.CoreV1().Nodes().Delete(ctx, name, opts)
		return nil, err
	})
	return err

}

func NodeGet(pctx context.Context, config *ikubernetes.ClusterConfig, name string, opts metav1.GetOptions) (*v1.Node, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Node, error) {
		resource, err := c.client.CoreV1().Nodes().Get(ctx, name, opts)

		return resource, err
	})
	return resp, err

}

func NodeList(pctx context.Context, config *ikubernetes.ClusterConfig, opts metav1.ListOptions) (*v1.NodeList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.NodeList, error) {
		resource, err := c.client.CoreV1().Nodes().List(ctx, opts)
		return resource, err
	})
	return resp, err

}

func NodeCordon(pctx context.Context, config *ikubernetes.ClusterConfig, name string) (*v1.Node, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Node, error) {
		resource, err := c.client.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		if resource.Spec.Unschedulable {
			return resource, nil
		}
		if err, _ := nodePatchOrReplace(c.client, resource, true); err != nil {
			return resource, err
		}
		return resource, nil
	})
	return resp, err

}

func nodePatchOrReplace(clientset kubernetes.Interface, node *corev1.Node, drain bool) (error, error) {
	client := clientset.CoreV1().Nodes()

	oldData, err := json.Marshal(node)
	if err != nil {
		return err, nil
	}

	node.Spec.Unschedulable = drain

	newData, err := json.Marshal(node)
	if err != nil {
		return err, nil
	}

	patchBytes, patchErr := strategicpatch.CreateTwoWayMergePatch(oldData, newData, node)
	if patchErr == nil {
		patchOptions := metav1.PatchOptions{}
		_, err = client.Patch(context.TODO(), node.GetName(), types.StrategicMergePatchType, patchBytes, patchOptions)
	} else {
		updateOptions := metav1.UpdateOptions{}
		_, err = client.Update(context.TODO(), node, updateOptions)
	}
	return err, patchErr
}

func NodeUncordon(pctx context.Context, config *ikubernetes.ClusterConfig, name string) (*v1.Node, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Node, error) {
		resource, err := c.client.CoreV1().Nodes().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		if !resource.Spec.Unschedulable {
			return resource, nil
		}
		if err, _ := nodePatchOrReplace(c.client, resource, false); err != nil {
			return resource, err

		}
		return resource, nil
	})
	return resp, err
}

func NodePatch(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	name string,
	pt types.PatchType,
	data []byte,
	opts metav1.PatchOptions,
	subResources ...string,
) (*v1.Node, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Node, error) {
		resource, err := c.client.CoreV1().Nodes().Patch(ctx, name, pt, data, opts, subResources...)
		return resource, err
	})
	return resp, err
}

func LimitRangeCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	data *v1.LimitRange,
	opts metav1.CreateOptions,
) (*v1.LimitRange, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.LimitRange, error) {
		resource, err := c.client.CoreV1().LimitRanges(namespace).Create(ctx, data, opts)
		return resource, err
	})
	return resp, err

}

func LimitRangeDelete(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.DeleteOptions,
) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Node, error) {
		err := c.client.CoreV1().LimitRanges(namespace).Delete(ctx, name, opts)
		return nil, err
	})
	return err

}

func LimitRangeUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	data *v1.LimitRange,
	opts metav1.UpdateOptions,
) (*v1.LimitRange, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.LimitRange, error) {
		resource, err := c.client.CoreV1().LimitRanges(namespace).Update(ctx, data, opts)
		return resource, err
	})
	return resp, err

}

func LimitRangeGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*v1.LimitRange, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.LimitRange, error) {
		resource, err := c.client.CoreV1().LimitRanges(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func LimitRangeList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*v1.LimitRangeList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.LimitRangeList, error) {
		resource, err := c.client.CoreV1().LimitRanges(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err
}

// namespace management, same as the topke project
func NamespaceCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace *v1.Namespace,
	opts metav1.CreateOptions,
) (*v1.Namespace, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Namespace, error) {
		resource, err := c.client.CoreV1().Namespaces().Create(ctx, namespace, opts)
		return resource, err
	})
	return resp, err
}

func NamespaceDelete(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Node, error) {
		err := c.client.CoreV1().Namespaces().Delete(ctx, namespace, opts)
		return nil, err
	})
	return err

}

func NamespaceGet(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, opts metav1.GetOptions) (*v1.Namespace, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Namespace, error) {
		resource, err := c.client.CoreV1().Namespaces().Get(ctx, namespace, opts)
		return resource, err
	})
	return resp, err

}

func NamespaceList(pctx context.Context, config *ikubernetes.ClusterConfig, opts metav1.ListOptions) (*v1.NamespaceList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.NamespaceList, error) {
		resource, err := c.client.CoreV1().Namespaces().List(ctx, opts)
		return resource, err
	})
	return resp, err

}

func NamespaceUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace *v1.Namespace,
	opts metav1.UpdateOptions,
) (*v1.Namespace, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Namespace, error) {
		resource, err := c.client.CoreV1().Namespaces().Update(ctx, namespace, opts)
		return resource, err
	})
	return resp, err
}

func NamespacePatch(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	pt types.PatchType,
	data []byte,
	opts metav1.PatchOptions,
	subResources []string,
) (*v1.Namespace, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Namespace, error) {
		resource, err := c.client.CoreV1().Namespaces().Patch(ctx, namespace, pt, data, opts, subResources...)
		return resource, err
	})
	return resp, err
}

// endpoints management, to manage the specify namespace/service endpoints, the endpoints include the pods of the
// service.
func EndpointCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	endpoints *v1.Endpoints,
	opts metav1.CreateOptions,
) (*v1.Endpoints, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Endpoints, error) {
		resource, err := c.client.CoreV1().Endpoints(namespace).Create(ctx, endpoints, opts)
		return resource, err
	})
	return resp, err
}

func EndpointDelete(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	endpoints string,
	opts metav1.DeleteOptions,
) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Node, error) {
		err := c.client.CoreV1().Endpoints(namespace).Delete(ctx, endpoints, opts)
		return nil, err
	})
	return err
}

func EndpointGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	endpoints string,
	opts metav1.GetOptions,
) (*v1.Endpoints, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Endpoints, error) {
		resource, err := c.client.CoreV1().Endpoints(namespace).Get(ctx, endpoints, opts)
		return resource, err
	})
	return resp, err
}

func EndpointList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*v1.EndpointsList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.EndpointsList, error) {
		resource, err := c.client.CoreV1().Endpoints(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err
}

func EndpointsUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	endpoints *v1.Endpoints,
	opts metav1.UpdateOptions,
) (*v1.Endpoints, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Endpoints, error) {
		resource, err := c.client.CoreV1().Endpoints(namespace).Update(ctx, endpoints, opts)
		return resource, err
	})
	return resp, err
}

// pod management
func PodCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	pod *v1.Pod,
	opts metav1.CreateOptions,
) (*v1.Pod, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Pod, error) {
		resource, err := c.client.CoreV1().Pods(namespace).Create(ctx, pod, opts)
		return resource, err
	})
	return resp, err
}

func PodUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	pod *v1.Pod,
	opts metav1.UpdateOptions,
) (*v1.Pod, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Pod, error) {
		resource, err := c.client.CoreV1().Pods(namespace).Update(ctx, pod, opts)
		return resource, err
	})
	return resp, err
}

func PodUpdateStatus(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	pod *v1.Pod,
	opts metav1.UpdateOptions,
) (*v1.Pod, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Pod, error) {
		resource, err := c.client.CoreV1().Pods(namespace).UpdateStatus(ctx, pod, opts)
		return resource, err
	})
	return resp, err
}

func PodDelete(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Node, error) {
		err := c.client.CoreV1().Pods(namespace).Delete(ctx, name, opts)
		return nil, err
	})
	return err
}

func PodDeleteCollection(pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	delOpts metav1.DeleteOptions,
	listOpts metav1.ListOptions,
) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Node, error) {
		err := c.client.CoreV1().Pods(namespace).DeleteCollection(ctx, delOpts, listOpts)
		return nil, err
	})
	return err
}

func PodGet(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.GetOptions) (*v1.Pod, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Pod, error) {
		resource, err := c.client.CoreV1().Pods(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func PodList(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, opts metav1.ListOptions) (*v1.PodList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.PodList, error) {
		resource, err := c.client.CoreV1().Pods(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err
}

func PodPatch(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	pt types.PatchType,
	data []byte,
	opts metav1.PatchOptions,
	subResources ...string,
) (*v1.Pod, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Pod, error) {
		resource, err := c.client.CoreV1().Pods(namespace).Patch(ctx, name, pt, data, opts, subResources...)
		return resource, err
	})
	return resp, err
}

// func PodGetEphemeralContainers(
// 	pctx context.Context,
// 	config *ikubernetes.ClusterConfig,
// 	namespace string,
// 	name string,
// 	opts metav1.GetOptions,
// ) (*v1.EphemeralContainers, error) {
// 	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Node, error) {
// 		containers, err := c.client.CoreV1().Pods(namespace).GetEphemeralContainers(ctx, name, opts)
// 		return resource, err
// 	})
// 	return resp, err
// }

// func PodUpdateEphemeralContainers(
// 	pctx context.Context,
// 	config *ikubernetes.ClusterConfig,
// 	namespace string,
// 	name string,
// 	ephemeralContainers *v1.EphemeralContainers,
// 	opts metav1.UpdateOptions,
// ) (*v1.EphemeralContainers, error) {
// 	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.EphemeralContainers, error) {
// 		resource, err := c.client.CoreV1().Pods(namespace).UpdateEphemeralContainers(ctx, name, ephemeralContainers, opts)
// 		return resource, err
// 	})
// 	return resp, err
// }

func PodLogList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts v1.PodLogOptions,
) ([]ikubernetes.PodLog, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) ([]ikubernetes.PodLog, error) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		defer cancel()
		ret, err := c.client.CoreV1().Pods(namespace).GetLogs(name, &opts).Do(ctx).Raw()
		if err != nil {
			if strings.Contains(err.Error(), "the server rejected our request for an unknown reason") ||
				strings.Contains(err.Error(), "has prevented the request from succeeding") {
				return nil, err
			}
			return nil, err
		}
		var logs []ikubernetes.PodLog
		sl := strings.Split(string(ret), "\n")
		for _, v := range sl {
			if v != "" {
				logs = append(logs, ikubernetes.PodLog{Msg: v})
			}
		}

		for i := 0; i < len(logs)/2; i++ { // 排序
			logs[i], logs[len(logs)-i-1] = logs[len(logs)-1-i], logs[i]
		}
		return logs, nil
	})
	return resp, err
}

func PodLogStream(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	podName string,
	opts v1.PodLogOptions,
) (io.ReadCloser, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (io.ReadCloser, error) {
		stream, err := c.client.CoreV1().Pods(namespace).GetLogs(podName, &opts).Stream(context.Background())
		if err != nil {
			if strings.Contains(err.Error(), "the server rejected our request for an unknown reason") ||
				strings.Contains(err.Error(), "has prevented the request from succeeding") {
				return nil, err
			}
			return nil, err
		}
		return stream, nil
	})
	return resp, err
}

func ServiceAccountCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	req *v1.ServiceAccount,
	opts metav1.CreateOptions,
) (*v1.ServiceAccount, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ServiceAccount, error) {
		resource, err := c.client.CoreV1().ServiceAccounts(namespace).Create(ctx, req, opts)
		return resource, err
	})
	return resp, err
}

func ServiceAccountUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	req *v1.ServiceAccount,
	opts metav1.UpdateOptions,
) (*v1.ServiceAccount, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ServiceAccount, error) {
		resource, err := c.client.CoreV1().ServiceAccounts(namespace).Update(ctx, req, opts)
		return resource, err
	})
	return resp, err
}

func ServiceAccountGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*v1.ServiceAccount, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ServiceAccount, error) {
		resource, err := c.client.CoreV1().ServiceAccounts(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func ServiceAccountList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*v1.ServiceAccountList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ServiceAccountList, error) {
		resource, err := c.client.CoreV1().ServiceAccounts(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err
}

func ServiceAccountDelete(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.DeleteOptions,
) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Node, error) {
		err := c.client.CoreV1().ServiceAccounts(namespace).Delete(ctx, name, opts)
		return nil, err
	})
	return err
}

func ResourceQuotaCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	rq *v1.ResourceQuota,
	opts metav1.CreateOptions,
) (*v1.ResourceQuota, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ResourceQuota, error) {
		resource, err := c.client.CoreV1().ResourceQuotas(namespace).Create(ctx, rq, opts)
		return resource, err
	})
	return resp, err
}

func ResourceQuotaUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	rq *v1.ResourceQuota,
	opts metav1.UpdateOptions,
) (*v1.ResourceQuota, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ResourceQuota, error) {
		resource, err := c.client.CoreV1().ResourceQuotas(namespace).Update(ctx, rq, opts)
		return resource, err
	})
	return resp, err
}

func ResourceQuotaDelete(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.DeleteOptions,
) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Node, error) {
		err := c.client.CoreV1().ResourceQuotas(namespace).Delete(ctx, name, opts)
		return nil, err
	})
	return err
}

func ResourceQuotaGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*v1.ResourceQuota, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ResourceQuota, error) {
		resource, err := c.client.CoreV1().ResourceQuotas(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func ResourceQuotaList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*v1.ResourceQuotaList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ResourceQuotaList, error) {
		resource, err := c.client.CoreV1().ResourceQuotas(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err
}

func PersistentVolumeCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	rq *v1.PersistentVolume,
	opts metav1.CreateOptions,
) (*v1.PersistentVolume, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.PersistentVolume, error) {
		resource, err := c.client.CoreV1().PersistentVolumes().Create(ctx, rq, opts)
		return resource, err
	})
	return resp, err
}

func PersistentVolumeDelete(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.DeleteOptions,
) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Node, error) {
		err := c.client.CoreV1().PersistentVolumes().Delete(ctx, name, opts)
		return nil, err
	})
	return err
}

func PersistentVolumeUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	rq *v1.PersistentVolume,
	opts metav1.UpdateOptions,
) (*v1.PersistentVolume, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.PersistentVolume, error) {
		resource, err := c.client.CoreV1().PersistentVolumes().Update(ctx, rq, opts)
		return resource, err
	})
	return resp, err
}

func PersistentVolumeGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	name string,
	opts metav1.GetOptions,
) (*v1.PersistentVolume, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.PersistentVolume, error) {
		resource, err := c.client.CoreV1().PersistentVolumes().Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func PersistentVolumeList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*v1.PersistentVolumeList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.PersistentVolumeList, error) {
		resource, err := c.client.CoreV1().PersistentVolumes().List(ctx, opts)
		return resource, err
	})
	return resp, err
}
