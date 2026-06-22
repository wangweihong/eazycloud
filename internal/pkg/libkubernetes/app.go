package libkubernetes

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/wangweihong/eazycloud/apis/ikubernetes"
)

func ReplicaSetCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	replicaSet *appsv1.ReplicaSet,
	opts metav1.CreateOptions,
) (*appsv1.ReplicaSet, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.ReplicaSet, error) {
		resource, err := c.client.AppsV1().ReplicaSets(namespace).Create(ctx, replicaSet, opts)
		return resource, err
	})
	return resp, err
}

func ReplicaSetDelete(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.DeleteOptions,
) error {

	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.ReplicaSet, error) {
		err := c.client.AppsV1().ReplicaSets(namespace).Delete(ctx, name, opts)
		return nil, err
	})
	return err
}

func ReplicaSetUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	replicaSet *appsv1.ReplicaSet,
	opts metav1.UpdateOptions,
) (*appsv1.ReplicaSet, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.ReplicaSet, error) {
		resource, err := c.client.AppsV1().ReplicaSets(namespace).Update(ctx, replicaSet, opts)
		return resource, err
	})
	return resp, err
}

func ReplicaSetGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*appsv1.ReplicaSet, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.ReplicaSet, error) {
		resource, err := c.client.AppsV1().ReplicaSets(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func ReplicaSetList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*appsv1.ReplicaSetList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.ReplicaSetList, error) {
		resource, err := c.client.AppsV1().ReplicaSets(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err
}

func StatefulSetCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	deployment *appsv1.StatefulSet,
	opts metav1.CreateOptions,
) (*appsv1.StatefulSet, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.StatefulSet, error) {
		resource, err := c.client.AppsV1().StatefulSets(namespace).Create(ctx, deployment, opts)
		return resource, err
	})
	return resp, err
}

func StatefulSetPatch(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	pt types.PatchType,
	data []byte,
	opts metav1.PatchOptions,
	subresources ...string,
) (*appsv1.StatefulSet, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.StatefulSet, error) {
		resource, err := c.client.AppsV1().StatefulSets(namespace).Patch(ctx, name, pt, data, opts, subresources...)
		return resource, err
	})
	return resp, err
}

func StatefulSetDelete(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.DeleteOptions,
) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.StatefulSet, error) {
		err := c.client.AppsV1().StatefulSets(namespace).Delete(ctx, name, opts)
		return nil, err
	})
	return err
}

func StatefulSetUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	deployment *appsv1.StatefulSet,
	opts metav1.UpdateOptions,
) (*appsv1.StatefulSet, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.StatefulSet, error) {
		resource, err := c.client.AppsV1().StatefulSets(namespace).Update(ctx, deployment, opts)
		return resource, err
	})
	return resp, err
}

func RevisionList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opt metav1.ListOptions,
) (*appsv1.ControllerRevisionList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.ControllerRevisionList, error) {
		resource, err := c.client.AppsV1().ControllerRevisions(namespace).List(ctx, opt)
		return resource, err
	})
	return resp, err
}

func StatefulSetGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*appsv1.StatefulSet, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.StatefulSet, error) {
		resource, err := c.client.AppsV1().StatefulSets(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func StatefulSetList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*appsv1.StatefulSetList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.StatefulSetList, error) {
		resource, err := c.client.AppsV1().StatefulSets(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err
}

func DaemonSetCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	daemonset *appsv1.DaemonSet,
	opts metav1.CreateOptions,
) (*appsv1.DaemonSet, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.DaemonSet, error) {
		resource, err := c.client.AppsV1().DaemonSets(namespace).Create(ctx, daemonset, opts)
		return resource, err
	})
	return resp, err
}

func DaemonSetPatch(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	pt types.PatchType,
	data []byte,
	opts metav1.PatchOptions,
	subresources ...string,
) (*appsv1.DaemonSet, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.DaemonSet, error) {
		resource, err := c.client.AppsV1().DaemonSets(namespace).Patch(ctx, name, pt, data, opts, subresources...)
		return resource, err
	})
	return resp, err
}

func DaemonSetUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	daemonset *appsv1.DaemonSet,
	opts metav1.UpdateOptions,
) (*appsv1.DaemonSet, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.DaemonSet, error) {
		resource, err := c.client.AppsV1().DaemonSets(namespace).Update(ctx, daemonset, opts)
		return resource, err
	})
	return resp, err
}

func DaemonSetDelete(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.DaemonSet, error) {
		err := c.client.AppsV1().DaemonSets(namespace).Delete(ctx, name, opts)
		return nil, err
	})
	return err
}

func DaemonSetGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*appsv1.DaemonSet, error) {

	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.DaemonSet, error) {
		resource, err := c.client.AppsV1().DaemonSets(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func DaemonSetList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*appsv1.DaemonSetList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.DaemonSetList, error) {
		resource, err := c.client.AppsV1().DaemonSets(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err
}

func DeploymentCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	deployment *appsv1.Deployment,
	opts metav1.CreateOptions,
) (*appsv1.Deployment, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.Deployment, error) {
		resource, err := c.client.AppsV1().Deployments(namespace).Create(ctx, deployment, opts)
		return resource, err
	})
	return resp, err
}

func DeploymentPatch(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	pt types.PatchType,
	data []byte,
	opts metav1.PatchOptions,
	subresources ...string,
) (*appsv1.Deployment, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.Deployment, error) {
		resource, err := c.client.AppsV1().Deployments(namespace).Patch(ctx, name, pt, data, opts, subresources...)
		return resource, err
	})
	return resp, err
}

func DeploymentUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	deployment *appsv1.Deployment,
	opts metav1.UpdateOptions,
) (*appsv1.Deployment, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.Deployment, error) {
		resource, err := c.client.AppsV1().Deployments(namespace).Update(ctx, deployment, opts)
		return resource, err
	})
	return resp, err
}

func DeploymentList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*appsv1.DeploymentList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.DeploymentList, error) {
		resource, err := c.client.AppsV1().Deployments(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err
}

func DeploymentGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*appsv1.Deployment, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*appsv1.Deployment, error) {
		resource, err := c.client.AppsV1().Deployments(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func DeploymentGetScale(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*autoscalingv1.Scale, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*autoscalingv1.Scale, error) {
		resource, err := c.client.AppsV1().Deployments(namespace).GetScale(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func DeploymentDelete(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.DeleteOptions,
) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*autoscalingv1.Scale, error) {
		err := c.client.AppsV1().Deployments(namespace).Delete(ctx, name, opts)
		return nil, err
	})
	return err
}
