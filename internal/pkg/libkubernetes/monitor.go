package libkubernetes

//
//import (
//	"context"
//	"github.com/wangweihong/eazycloud/internal/pkg/apis/topke"
//
//	"time"
//
//	//monitorv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
//	monitorv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring"
//
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//)
//
// func ServiceMonitorCreate(config *ikubernetes.ClusterConfig, namespace string, sm *monitorv1.ServiceMonitor, opts
// metav1.CreateOptions) (*monitorv1.ServiceMonitor, error) {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return nil, err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Second)
//	defer cancel()
//	sm, err = c.monitorClient.ServiceMonitors(namespace).Create(ctx, sm, opts)
//	if err != nil {
//		return nil, err
//	}
//
//	return sm, nil
//}
//
// func ServiceMonitorDelete(config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.DeleteOptions)
// error {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Second)
//	defer cancel()
//	err = c.monitorClient.ServiceMonitors(namespace).Delete(ctx, name, opts)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
// func ServiceMonitorUpdate(config *ikubernetes.ClusterConfig, namespace string, sm *monitorv1.ServiceMonitor, opts
// metav1.UpdateOptions) (*monitorv1.ServiceMonitor, error) {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return nil, err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Second)
//	defer cancel()
//	sm, err = c.monitorClient.ServiceMonitors(namespace).Update(ctx, sm, opts)
//	if err != nil {
//		return nil, err
//	}
//
//	return sm, nil
//}
//
// func ServiceMonitorGet(config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.GetOptions)
// (*monitorv1.ServiceMonitor, error) {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return nil, err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Second)
//	defer cancel()
//	sm, err := c.monitorClient.ServiceMonitors(namespace).Get(ctx, name, opts)
//	if err != nil {
//		return nil, err
//	}
//	sm.TypeMeta = getResourceTypeMeta(topke.KubernetesResourceKindServiceMonitor)
//
//	return sm, nil
//}
//
// func ServiceMonitorList(config *ikubernetes.ClusterConfig, namespace string, opts metav1.ListOptions, timeout ...int)
// (*monitorv1.ServiceMonitorList, error) {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return nil, err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(getTimeOut(timeout...))*time.Second)
//	defer cancel()
//	sm, err := c.monitorClient.ServiceMonitors(namespace).List(ctx, opts)
//	if err != nil {
//		return nil, err
//	}
//
//	return sm, nil
//}
//
// func PodMonitorCreate(config *ikubernetes.ClusterConfig, namespace string, pm *monitorv1.PodMonitor, opts
// metav1.CreateOptions) (*monitorv1.PodMonitor, error) {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return nil, err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Second)
//	defer cancel()
//	sm, err := c.monitorClient.PodMonitors(namespace).Create(ctx, pm, opts)
//	if err != nil {
//		return nil, err
//	}
//
//	return sm, nil
//}
//
// func PodMonitorDelete(config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.DeleteOptions)
// error {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Second)
//	defer cancel()
//	err = c.monitorClient.PodMonitors(namespace).Delete(ctx, name, opts)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
// func PodMonitorUpdate(config *ikubernetes.ClusterConfig, namespace string, pm *monitorv1.PodMonitor, opts
// metav1.UpdateOptions) (*monitorv1.PodMonitor, error) {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return nil, err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Second)
//	defer cancel()
//	pm, err = c.monitorClient.PodMonitors(namespace).Update(ctx, pm, opts)
//	if err != nil {
//		return nil, err
//	}
//
//	return pm, nil
//}
//
// func PodMonitorGet(config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.GetOptions)
// (*monitorv1.PodMonitor, error) {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return nil, err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Second)
//	defer cancel()
//	pm, err := c.monitorClient.PodMonitors(namespace).Get(ctx, name, opts)
//	if err != nil {
//		return nil, err
//	}
//	pm.TypeMeta = getResourceTypeMeta(topke.KubernetesResourceKindPodMoniotr)
//
//	return pm, nil
//}
//
// func PodMonitorList(config *ikubernetes.ClusterConfig, namespace string, opts metav1.ListOptions, timeout ...int)
// (*monitorv1.PodMonitorList, error) {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return nil, err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(getTimeOut(timeout...))*time.Second)
//	defer cancel()
//	pm, err := c.monitorClient.PodMonitors(namespace).List(ctx, opts)
//	if err != nil {
//		return nil, err
//	}
//
//	return pm, nil
//}
//
// func PrometheusCreate(config *ikubernetes.ClusterConfig, namespace string, prom *monitorv1.Prometheus, opts
// metav1.CreateOptions) (*monitorv1.Prometheus, error) {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return nil, err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Second)
//	defer cancel()
//	pm, err := c.monitorClient.Prometheuses(namespace).Create(ctx, prom, opts)
//	if err != nil {
//		return nil, err
//	}
//
//	return pm, nil
//}
//
// func PrometheusDelete(config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.DeleteOptions)
// error {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Second)
//	defer cancel()
//	err = c.monitorClient.Prometheuses(namespace).Delete(ctx, name, opts)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
// func PrometheusUpdate(config *ikubernetes.ClusterConfig, namespace string, prom *monitorv1.Prometheus, opts
// metav1.UpdateOptions) (*monitorv1.Prometheus, error) {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return nil, err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Second)
//	defer cancel()
//	prom, err = c.monitorClient.Prometheuses(namespace).Update(ctx, prom, opts)
//	if err != nil {
//		return nil, err
//	}
//
//	return prom, nil
//}
//
// func PrometheusGet(config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.GetOptions)
// (*monitorv1.Prometheus, error) {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return nil, err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Second)
//	defer cancel()
//	prom, err := c.monitorClient.Prometheuses(namespace).Get(ctx, name, opts)
//	if err != nil {
//		return nil, err
//	}
//	prom.TypeMeta = getResourceTypeMeta(topke.KubernetesResourceKindPrometheus)
//
//	return prom, nil
//}
//
// func PrometheusList(config *ikubernetes.ClusterConfig, namespace string, opts metav1.ListOptions, timeout ...int)
// (*monitorv1.PrometheusList, error) {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return nil, err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(getTimeOut(timeout...))*time.Second)
//	defer cancel()
//	prom, err := c.monitorClient.Prometheuses(namespace).List(ctx, opts)
//	if err != nil {
//		return nil, err
//	}
//
//	return prom, nil
//}
//
// func PrometheusRuleCreate(config *ikubernetes.ClusterConfig, namespace string, pr *monitorv1.PrometheusRule, opts
// metav1.CreateOptions) (*monitorv1.PrometheusRule, error) {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return nil, err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Second)
//	defer cancel()
//	prom, err := c.monitorClient.PrometheusRules(namespace).Create(ctx, pr, opts)
//	if err != nil {
//		return nil, err
//	}
//
//	return prom, nil
//}
//
// func PrometheusRuleDelete(config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.DeleteOptions)
// error {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Second)
//	defer cancel()
//	err = c.monitorClient.PrometheusRules(namespace).Delete(ctx, name, opts)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
// func PrometheusRuleUpdate(config *ikubernetes.ClusterConfig, namespace string, pr *monitorv1.PrometheusRule, opts
// metav1.UpdateOptions) (*monitorv1.PrometheusRule, error) {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return nil, err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Second)
//	defer cancel()
//	pr, err = c.monitorClient.PrometheusRules(namespace).Update(ctx, pr, opts)
//	if err != nil {
//		return nil, err
//	}
//
//	return pr, nil
//}
//
// func PrometheusRuleGet(config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.GetOptions)
// (*monitorv1.PrometheusRule, error) {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return nil, err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Second)
//	defer cancel()
//	pr, err := c.monitorClient.PrometheusRules(namespace).Get(ctx, name, opts)
//	if err != nil {
//		return nil, err
//	}
//	pr.TypeMeta = getResourceTypeMeta(topke.KubernetesResourceKindPrometheusRule)
//
//	return pr, nil
//}
//
// func PrometheusRuleList(config *ikubernetes.ClusterConfig, namespace string, opts metav1.ListOptions, timeout ...int)
// (*monitorv1.PrometheusRuleList, error) {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return nil, err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(getTimeOut(timeout...))*time.Second)
//	defer cancel()
//	pr, err := c.monitorClient.PrometheusRules(namespace).List(ctx, opts)
//	if err != nil {
//		return nil, err
//	}
//
//	return pr, nil
//}
//
// func AlertManagerCreate(config *ikubernetes.ClusterConfig, namespace string, am *monitorv1.Alertmanager, opts
// metav1.CreateOptions) (*monitorv1.Alertmanager, error) {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return nil, err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Second)
//	defer cancel()
//	pr, err := c.monitorClient.Alertmanagers(namespace).Create(ctx, am, opts)
//	if err != nil {
//		return nil, err
//	}
//
//	return pr, nil
//}
//
// func AlertManagerDelete(config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.DeleteOptions)
// error {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Second)
//	defer cancel()
//	err = c.monitorClient.Alertmanagers(namespace).Delete(ctx, name, opts)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
// func AlertManagerUpdate(config *ikubernetes.ClusterConfig, namespace string, am *monitorv1.Alertmanager, opts
// metav1.UpdateOptions) (*monitorv1.Alertmanager, error) {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return nil, err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Second)
//	defer cancel()
//	pr, err := c.monitorClient.Alertmanagers(namespace).Update(ctx, am, opts)
//	if err != nil {
//		return nil, err
//	}
//
//	return pr, nil
//}
//
// func AlertManagerGet(config *ikubernetes.ClusterConfig, namespace string, name string, opts metav1.GetOptions)
// (*monitorv1.Alertmanager, error) {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return nil, err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultTimeout)*time.Second)
//	defer cancel()
//	pr, err := c.monitorClient.Alertmanagers(namespace).Get(ctx, name, opts)
//	if err != nil {
//		return nil, err
//	}
//	pr.TypeMeta = getResourceTypeMeta(topke.KubernetesResourceKindServiceMonitor)
//
//	return pr, nil
//}
//
// func AlertManagerList(config *ikubernetes.ClusterConfig, namespace string, opts metav1.ListOptions, timeout ...int)
// (*monitorv1.AlertmanagerList, error) {
//	c, err := NewMonitorClient(config)
//	if err != nil {
//		return nil, err
//	}
//	defer c.Close()
//	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(getTimeOut(timeout...))*time.Second)
//	defer cancel()
//	pr, err := c.monitorClient.Alertmanagers(namespace).List(ctx, opts)
//	if err != nil {
//		return nil, err
//	}
//
//	return pr, nil
//}
