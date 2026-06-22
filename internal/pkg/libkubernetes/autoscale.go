package libkubernetes

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/ikubernetes"

	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func HpaCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	hpa *autoscalingv1.HorizontalPodAutoscaler,
	opts metav1.CreateOptions,
) (*autoscalingv1.HorizontalPodAutoscaler, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*autoscalingv1.HorizontalPodAutoscaler, error) {
		resource, err := c.client.AutoscalingV1().HorizontalPodAutoscalers(namespace).Create(ctx, hpa, opts)
		return resource, err
	})
	return resp, err
}

func HpaDelete(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*autoscalingv1.HorizontalPodAutoscaler, error) {
		err := c.client.AutoscalingV1().HorizontalPodAutoscalers(namespace).Delete(ctx, name, opts)
		return nil, err
	})
	return err
}

func HpaUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	hpa *autoscalingv1.HorizontalPodAutoscaler,
	opts metav1.UpdateOptions,
) (*autoscalingv1.HorizontalPodAutoscaler, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*autoscalingv1.HorizontalPodAutoscaler, error) {
		resource, err := c.client.AutoscalingV1().HorizontalPodAutoscalers(namespace).Update(ctx, hpa, opts)
		return resource, err
	})
	return resp, err
}

func HpaGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*autoscalingv1.HorizontalPodAutoscaler, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*autoscalingv1.HorizontalPodAutoscaler, error) {
		resource, err := c.client.AutoscalingV1().HorizontalPodAutoscalers(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func HpaList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*autoscalingv1.HorizontalPodAutoscalerList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*autoscalingv1.HorizontalPodAutoscalerList, error) {
		resource, err := c.client.AutoscalingV1().HorizontalPodAutoscalers(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err
}
