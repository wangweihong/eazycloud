package libkubernetes

import (
	"context"

	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/wangweihong/eazycloud/apis/ikubernetes"
)

func PodEvict(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, podName string, opts metav1.DeleteOptions) error {
	eviction := policyv1.Eviction{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      podName,
		},
		DeleteOptions: &opts,
	}
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*policyv1.PodDisruptionBudget, error) {
		err := c.client.PolicyV1().Evictions(namespace).Evict(ctx, &eviction)
		return nil, err
	})
	return err
}

func PodDisruptionBudgetCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	configMap *policyv1.PodDisruptionBudget,
	opts metav1.CreateOptions,
) (*policyv1.PodDisruptionBudget, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*policyv1.PodDisruptionBudget, error) {
		resource, err := c.client.PolicyV1().PodDisruptionBudgets(namespace).Create(ctx, configMap, opts)
		return resource, err
	})
	return resp, err
}

func PodDisruptionBudgetUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	configMap *policyv1.PodDisruptionBudget,
	opts metav1.UpdateOptions,
) (*policyv1.PodDisruptionBudget, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*policyv1.PodDisruptionBudget, error) {
		resource, err := c.client.PolicyV1().PodDisruptionBudgets(namespace).Update(ctx, configMap, opts)
		return resource, err
	})
	return resp, err
}

func PodDisruptionBudgetDelete(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.DeleteOptions,
) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*policyv1.PodDisruptionBudget, error) {
		err := c.client.PolicyV1().PodDisruptionBudgets(namespace).Delete(ctx, name, opts)
		return nil, err
	})
	return err
}

func PodDisruptionBudgetGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*policyv1.PodDisruptionBudget, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*policyv1.PodDisruptionBudget, error) {
		resource, err := c.client.PolicyV1().PodDisruptionBudgets(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func PodDisruptionBudgetList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*policyv1.PodDisruptionBudgetList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*policyv1.PodDisruptionBudgetList, error) {
		resource, err := c.client.PolicyV1().PodDisruptionBudgets(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err
}
