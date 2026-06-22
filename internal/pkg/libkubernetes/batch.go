package libkubernetes

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/ikubernetes"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func JobCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	job *batchv1.Job,
	opts metav1.CreateOptions,
) (*batchv1.Job, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*batchv1.Job, error) {
		resource, err := c.client.BatchV1().Jobs(namespace).Create(ctx, job, opts)
		return resource, err
	})
	return resp, err
}

func JobDelete(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*batchv1.Job, error) {
		err := c.client.BatchV1().Jobs(namespace).Delete(ctx, name, opts)
		return nil, err
	})
	return err
}

func JobGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*batchv1.Job, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*batchv1.Job, error) {
		resource, err := c.client.BatchV1().Jobs(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func JobUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	job *batchv1.Job,
	opts metav1.UpdateOptions,
) (*batchv1.Job, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*batchv1.Job, error) {
		resource, err := c.client.BatchV1().Jobs(namespace).Update(ctx, job, opts)
		return resource, err
	})
	return resp, err
}

func JobList(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, opts metav1.ListOptions) (*batchv1.JobList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*batchv1.JobList, error) {
		resource, err := c.client.BatchV1().Jobs(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err
}

func CronJobCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	cronjob *batchv1.CronJob,
	opts metav1.CreateOptions,
) (*batchv1.CronJob, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*batchv1.CronJob, error) {
		resource, err := c.client.BatchV1().CronJobs(namespace).Create(ctx, cronjob, opts)
		return resource, err
	})
	return resp, err
}

func CronJobUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	cronjob *batchv1.CronJob,
	opts metav1.UpdateOptions,
) (*batchv1.CronJob, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*batchv1.CronJob, error) {
		resource, err := c.client.BatchV1().CronJobs(namespace).Update(ctx, cronjob, opts)
		return resource, err
	})
	return resp, err
}

func CronJobDelete(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*batchv1.CronJob, error) {
		err := c.client.BatchV1().CronJobs(namespace).Delete(ctx, name, opts)
		return nil, err
	})
	return err
}

func CronJobGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*batchv1.CronJob, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*batchv1.CronJob, error) {
		resource, err := c.client.BatchV1().CronJobs(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func CronJobList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*batchv1.CronJobList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*batchv1.CronJobList, error) {
		resource, err := c.client.BatchV1().CronJobs(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err
}
