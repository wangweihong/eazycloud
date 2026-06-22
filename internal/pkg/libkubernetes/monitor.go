package libkubernetes

import (
	"context"

	monitorv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/wangweihong/eazycloud/apis/ikubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ServiceMonitorCreate(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, res *monitorv1.ServiceMonitor, opts metav1.CreateOptions) (*monitorv1.ServiceMonitor, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.ServiceMonitor, error) {
		resource, err := c.monitorClient.MonitoringV1().ServiceMonitors(namespace).Create(ctx, res, opts)
		return resource, err
	})
	return resp, err
}

func ServiceMonitorDelete(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.ServiceMonitor, error) {
		err := c.monitorClient.MonitoringV1().ServiceMonitors(namespace).Delete(ctx, name, opts)
		return nil, err
	})
	return err

}

func ServiceMonitorUpdate(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, res *monitorv1.ServiceMonitor, opts metav1.UpdateOptions) (*monitorv1.ServiceMonitor, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.ServiceMonitor, error) {
		resource, err := c.monitorClient.MonitoringV1().ServiceMonitors(namespace).Update(ctx, res, opts)
		return resource, err
	})

	return resp, err
}

func ServiceMonitorGet(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.GetOptions) (*monitorv1.ServiceMonitor, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.ServiceMonitor, error) {
		resource, err := c.monitorClient.MonitoringV1().ServiceMonitors(namespace).Get(ctx, name, opts)
		if err != nil {
			return nil, err
		}
		resource.TypeMeta = getResourceTypeMeta(ikubernetes.KubernetesResourceKindServiceMonitor)
		return resource, err
	})
	return resp, err

}

func ServiceMonitorList(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, opts metav1.ListOptions, timeout ...int) (*monitorv1.ServiceMonitorList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.ServiceMonitorList, error) {
		resource, err := c.monitorClient.MonitoringV1().ServiceMonitors(namespace).List(ctx, opts)
		return resource, err
	})

	return resp, err
}

func PodMonitorCreate(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, res *monitorv1.PodMonitor, opts metav1.CreateOptions) (*monitorv1.PodMonitor, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.PodMonitor, error) {
		resource, err := c.monitorClient.MonitoringV1().PodMonitors(namespace).Create(ctx, res, opts)
		return resource, err
	})
	return resp, err
}

func PodMonitorDelete(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.PodMonitor, error) {
		err := c.monitorClient.MonitoringV1().PodMonitors(namespace).Delete(ctx, name, opts)
		return nil, err
	})
	return err

}

func PodMonitorUpdate(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, res *monitorv1.PodMonitor, opts metav1.UpdateOptions) (*monitorv1.PodMonitor, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.PodMonitor, error) {
		resource, err := c.monitorClient.MonitoringV1().PodMonitors(namespace).Update(ctx, res, opts)
		return resource, err
	})

	return resp, err
}

func PodMonitorGet(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.GetOptions) (*monitorv1.PodMonitor, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.PodMonitor, error) {
		resource, err := c.monitorClient.MonitoringV1().PodMonitors(namespace).Get(ctx, name, opts)
		if err != nil {
			return nil, err
		}
		resource.TypeMeta = getResourceTypeMeta(ikubernetes.KubernetesResourceKindPodMonitor)
		return resource, err
	})
	return resp, err

}

func PodMonitorList(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, opts metav1.ListOptions, timeout ...int) (*monitorv1.PodMonitorList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.PodMonitorList, error) {
		resource, err := c.monitorClient.MonitoringV1().PodMonitors(namespace).List(ctx, opts)
		return resource, err
	})

	return resp, err
}

func PrometheusCreate(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, res *monitorv1.Prometheus, opts metav1.CreateOptions) (*monitorv1.Prometheus, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.Prometheus, error) {
		resource, err := c.monitorClient.MonitoringV1().Prometheuses(namespace).Create(ctx, res, opts)
		return resource, err
	})

	return resp, err
}

func PrometheusDelete(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.Prometheus, error) {
		err := c.monitorClient.MonitoringV1().Prometheuses(namespace).Delete(ctx, name, opts)
		return nil, err
	})

	return err
}

func PrometheusUpdate(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, res *monitorv1.Prometheus, opts metav1.UpdateOptions) (*monitorv1.Prometheus, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.Prometheus, error) {
		resource, err := c.monitorClient.MonitoringV1().Prometheuses(namespace).Update(ctx, res, opts)
		return resource, err
	})

	return resp, err
}

func PrometheusGet(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.GetOptions) (*monitorv1.Prometheus, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.Prometheus, error) {
		resource, err := c.monitorClient.MonitoringV1().Prometheuses(namespace).Get(ctx, name, opts)
		if err != nil {
			return nil, err
		}
		resource.TypeMeta = getResourceTypeMeta(ikubernetes.KubernetesResourceKindPrometheus)
		return resource, err
	})

	return resp, err
}

func PrometheusList(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, opts metav1.ListOptions, timeout ...int) (*monitorv1.PrometheusList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.PrometheusList, error) {
		resource, err := c.monitorClient.MonitoringV1().Prometheuses(namespace).List(ctx, opts)
		return resource, err
	})

	return resp, err
}

func PrometheusRuleCreate(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, res *monitorv1.PrometheusRule, opts metav1.CreateOptions) (*monitorv1.PrometheusRule, error) {

	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.PrometheusRule, error) {
		resource, err := c.monitorClient.MonitoringV1().PrometheusRules(namespace).Create(ctx, res, opts)
		return resource, err
	})

	return resp, err
}

func PrometheusRuleDelete(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.PrometheusRule, error) {
		err := c.monitorClient.MonitoringV1().PrometheusRules(namespace).Delete(ctx, name, opts)
		return nil, err
	})

	return err

}

func PrometheusRuleUpdate(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, res *monitorv1.PrometheusRule, opts metav1.UpdateOptions) (*monitorv1.PrometheusRule, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.PrometheusRule, error) {
		resource, err := c.monitorClient.MonitoringV1().PrometheusRules(namespace).Update(ctx, res, opts)
		return resource, err
	})
	return resp, err
}

func PrometheusRuleGet(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.GetOptions) (*monitorv1.PrometheusRule, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.PrometheusRule, error) {
		resource, err := c.monitorClient.MonitoringV1().PrometheusRules(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func PrometheusRuleList(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, opts metav1.ListOptions, timeout ...int) (*monitorv1.PrometheusRuleList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.PrometheusRuleList, error) {
		resource, err := c.monitorClient.MonitoringV1().PrometheusRules(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err
}

func AlertManagerCreate(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, res *monitorv1.Alertmanager, opts metav1.CreateOptions) (*monitorv1.Alertmanager, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.Alertmanager, error) {
		resource, err := c.monitorClient.MonitoringV1().Alertmanagers(namespace).Create(ctx, res, opts)
		return resource, err
	})
	return resp, err
}

func AlertManagerDelete(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.Alertmanager, error) {
		err := c.monitorClient.MonitoringV1().Alertmanagers(namespace).Delete(ctx, name, opts)
		return nil, err
	})
	return err
}

func AlertManagerUpdate(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, res *monitorv1.Alertmanager, opts metav1.UpdateOptions) (*monitorv1.Alertmanager, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.Alertmanager, error) {
		resource, err := c.monitorClient.MonitoringV1().Alertmanagers(namespace).Update(ctx, res, opts)
		return resource, err
	})
	return resp, err
}

func AlertManagerGet(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.GetOptions) (*monitorv1.Alertmanager, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.Alertmanager, error) {
		resource, err := c.monitorClient.MonitoringV1().Alertmanagers(namespace).Get(ctx, name, opts)
		if err != nil {
			return resource, err
		}
		resource.TypeMeta = getResourceTypeMeta(ikubernetes.KubernetesResourceKindServiceMonitor)
		return resource, err
	})

	return resp, err
}

func AlertManagerList(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, opts metav1.ListOptions, timeout ...int) (*monitorv1.AlertmanagerList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*monitorv1.AlertmanagerList, error) {
		resource, err := c.monitorClient.MonitoringV1().Alertmanagers(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err
}
