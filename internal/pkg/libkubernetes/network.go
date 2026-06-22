package libkubernetes

import (
	"context"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/wangweihong/eazycloud/apis/ikubernetes"
)

func NetworkPolicyList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*networkingv1.NetworkPolicyList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*networkingv1.NetworkPolicyList, error) {
		resource, err := c.client.NetworkingV1().NetworkPolicies(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err
}

func NetworkPolicyGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*networkingv1.NetworkPolicy, error) {

	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*networkingv1.NetworkPolicy, error) {
		resource, err := c.client.NetworkingV1().NetworkPolicies(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func NetworkPolicyCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	res *networkingv1.NetworkPolicy,
	opts metav1.CreateOptions,
) (*networkingv1.NetworkPolicy, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*networkingv1.NetworkPolicy, error) {
		resource, err := c.client.NetworkingV1().NetworkPolicies(namespace).Create(ctx, res, opts)
		return resource, err
	})
	return resp, err
}

func NetworkPolicyUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	res *networkingv1.NetworkPolicy,
	opts metav1.UpdateOptions,
) (*networkingv1.NetworkPolicy, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*networkingv1.NetworkPolicy, error) {
		resource, err := c.client.NetworkingV1().NetworkPolicies(namespace).Update(ctx, res, opts)
		return resource, err
	})
	return resp, err
}

func NetworkPolicyDelete(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	res *networkingv1.NetworkPolicy,
	opts metav1.DeleteOptions,
) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*networkingv1.NetworkPolicy, error) {
		err := c.client.NetworkingV1().NetworkPolicies(namespace).Delete(ctx, res.Name, opts)
		return nil, err
	})
	return err
}

func IngressCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	ingress *networkingv1.Ingress,
	opts metav1.CreateOptions,
) (*networkingv1.Ingress, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*networkingv1.Ingress, error) {
		resource, err := c.client.NetworkingV1().Ingresses(namespace).Create(ctx, ingress, opts)
		return resource, err
	})
	return resp, err
}

func IngressDelete(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*networkingv1.NetworkPolicy, error) {
		err := c.client.NetworkingV1().Ingresses(namespace).Delete(ctx, name, opts)
		return nil, err
	})
	return err
}

func IngressUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	ingress *networkingv1.Ingress,
	opts metav1.UpdateOptions,
) (*networkingv1.Ingress, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*networkingv1.Ingress, error) {
		resource, err := c.client.NetworkingV1().Ingresses(namespace).Update(ctx, ingress, opts)
		return resource, err
	})
	return resp, err
}

func IngressGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*networkingv1.Ingress, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*networkingv1.Ingress, error) {
		resource, err := c.client.NetworkingV1().Ingresses(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func IngressList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*networkingv1.IngressList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*networkingv1.IngressList, error) {
		resource, err := c.client.NetworkingV1().Ingresses(namespace).List(ctx, opts)
		return resource, err
	})

	return resp, err
}
