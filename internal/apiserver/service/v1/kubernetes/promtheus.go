package kubernetes

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	//gapi "github.com/grafana/grafana-openapi-client-go"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	prometheusv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/apis/ikubernetes"
	"github.com/wangweihong/eazycloud/apis/iprometheus"
	"github.com/wangweihong/eazycloud/internal/pkg/clientset"
	"github.com/wangweihong/eazycloud/internal/pkg/libkubernetes"
	"github.com/wangweihong/eazycloud/internal/pkg/libprometheus"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/log"
	"github.com/wangweihong/gotoolbox/pkg/maputil"
	"github.com/wangweihong/gotoolbox/pkg/mathutil"
	"github.com/wangweihong/gotoolbox/pkg/paging"
	"github.com/wangweihong/gotoolbox/pkg/sets"
	"github.com/wangweihong/gotoolbox/pkg/waitgroup"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (k *kubernetesService) AlertManagerCreate(ctx context.Context, req *iapiserver.AlertManagerRequest) (*iapiserver.AlertManagerInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.AlertManagerCreate(ctx, cluster.Config, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewAlertManagerInfo(meta, cluster), nil
}

func (k *kubernetesService) AlertManagerDelete(ctx context.Context, req *iapiserver.AlertManagerRequest) error {
	req.Resource.ResourceVersion = ""
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := libkubernetes.AlertManagerDelete(ctx, cluster.Config, req.Resource.Namespace, req.Resource.Name, req.DeleteOpts); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) AlertManagerUpdate(ctx context.Context, req *iapiserver.AlertManagerRequest) (*iapiserver.AlertManagerInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.AlertManagerUpdate(ctx, cluster.Config, req.Resource.Namespace, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return iapiserver.NewAlertManagerInfo(meta, cluster), nil
}

func (k *kubernetesService) AlertManagerGet(ctx context.Context, req *iapiserver.AlertManagerGetRequest) (*iapiserver.AlertManagerInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.AlertManagerGet(ctx, cluster.Config, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewAlertManagerInfo(meta, cluster), nil
}

func (k *kubernetesService) AlertManagerList(ctx context.Context, req *iapiserver.AlertManagerListRequest) (*iapiserver.AlertManagerListResponse, error) {
	resp := &iapiserver.AlertManagerListResponse{}

	if req.Namespace == "" {
		req.Namespace = ikubernetes.MonitorNamespace
	}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ret, err := libkubernetes.AlertManagerList(ctx, cluster.Config, req.Namespace, req.ToListOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var resources []*iapiserver.AlertManagerInfo
	for _, dm := range ret.Items {
		if req.Fuzzy != "" && !sets.NewString(dm.Name).ContainAny(strings.Split(req.Fuzzy, ",")...) {
			continue
		}

		resources = append(resources, iapiserver.NewAlertManagerInfo(&dm, cluster))
	}

	resp.List = resources
	s, e := paging.Index(len(resources), req.PageNum, req.PageSize)
	resp.List = resp.List[s:e]
	return resp, nil
}

func (k *kubernetesService) ServiceMonitorCreate(ctx context.Context, req *iapiserver.ServiceMonitorRequest) (*iapiserver.ServiceMonitorInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.ServiceMonitorCreate(ctx, cluster.Config, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewServiceMonitorInfo(meta, cluster), nil
}

func (k *kubernetesService) ServiceMonitorDelete(ctx context.Context, req *iapiserver.ServiceMonitorRequest) error {
	req.Resource.ResourceVersion = ""
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := libkubernetes.ServiceMonitorDelete(ctx, cluster.Config, req.Resource.Namespace, req.Resource.Name, req.DeleteOpts); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) ServiceMonitorUpdate(ctx context.Context, req *iapiserver.ServiceMonitorRequest) (*iapiserver.ServiceMonitorInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.ServiceMonitorUpdate(ctx, cluster.Config, req.Resource.Namespace, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return iapiserver.NewServiceMonitorInfo(meta, cluster), nil
}

func (k *kubernetesService) ServiceMonitorGet(ctx context.Context, req *iapiserver.ServiceMonitorGetRequest) (*iapiserver.ServiceMonitorInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.ServiceMonitorGet(ctx, cluster.Config, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewServiceMonitorInfo(meta, cluster), nil
}

func (k *kubernetesService) ServiceMonitorList(ctx context.Context, req *iapiserver.ServiceMonitorListRequest) (*iapiserver.ServiceMonitorListResponse, error) {
	resp := &iapiserver.ServiceMonitorListResponse{}

	if req.Namespace == "" {
		req.Namespace = ikubernetes.MonitorNamespace
	}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ret, err := libkubernetes.ServiceMonitorList(ctx, cluster.Config, req.Namespace, req.ToListOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var resources []*iapiserver.ServiceMonitorInfo
	for _, dm := range ret.Items {
		if req.Fuzzy != "" && !sets.NewString(dm.Name).ContainAny(strings.Split(req.Fuzzy, ",")...) {
			continue
		}

		resources = append(resources, iapiserver.NewServiceMonitorInfo(&dm, cluster))
	}

	resp.List = resources
	s, e := paging.Index(len(resources), req.PageNum, req.PageSize)
	resp.List = resp.List[s:e]
	return resp, nil
}

func (k *kubernetesService) PodMonitorCreate(ctx context.Context, req *iapiserver.PodMonitorRequest) (*iapiserver.PodMonitorInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.PodMonitorCreate(ctx, cluster.Config, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewPodMonitorInfo(meta, cluster), nil
}

func (k *kubernetesService) PodMonitorDelete(ctx context.Context, req *iapiserver.PodMonitorRequest) error {
	req.Resource.ResourceVersion = ""
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := libkubernetes.PodMonitorDelete(ctx, cluster.Config, req.Resource.Namespace, req.Resource.Name, req.DeleteOpts); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) PodMonitorUpdate(ctx context.Context, req *iapiserver.PodMonitorRequest) (*iapiserver.PodMonitorInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.PodMonitorUpdate(ctx, cluster.Config, req.Resource.Namespace, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return iapiserver.NewPodMonitorInfo(meta, cluster), nil
}

func (k *kubernetesService) PodMonitorGet(ctx context.Context, req *iapiserver.PodMonitorGetRequest) (*iapiserver.PodMonitorInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.PodMonitorGet(ctx, cluster.Config, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewPodMonitorInfo(meta, cluster), nil
}

func (k *kubernetesService) PodMonitorList(ctx context.Context, req *iapiserver.PodMonitorListRequest) (*iapiserver.PodMonitorListResponse, error) {
	resp := &iapiserver.PodMonitorListResponse{}

	if req.Namespace == "" {
		req.Namespace = ikubernetes.MonitorNamespace
	}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ret, err := libkubernetes.PodMonitorList(ctx, cluster.Config, req.Namespace, req.ToListOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var resources []*iapiserver.PodMonitorInfo
	for _, dm := range ret.Items {
		if req.Fuzzy != "" && !sets.NewString(dm.Name).ContainAny(strings.Split(req.Fuzzy, ",")...) {
			continue
		}

		resources = append(resources, iapiserver.NewPodMonitorInfo(&dm, cluster))
	}

	resp.List = resources
	s, e := paging.Index(len(resources), req.PageNum, req.PageSize)
	resp.List = resp.List[s:e]
	return resp, nil
}

var (
	prometheusServiceName = "prometheus-k8s"
	grafanaServiceName    = "grafana"
	statefulset           = "statefulset"
	daemonset             = "daemonset"
	deployment            = "deployment"
	pod                   = "pod"
	nodePort              = "NodePort"
	componentEtcd         = "etcd"
)

func (k *kubernetesService) getPrometheusUrlFromPrometheusService(ctx context.Context, clusterInfo *iapiserver.Cluster) (string, error) {
	service, err := clientset.ServiceGet(ctx, clusterInfo, ikubernetes.MonitorNamespace, prometheusServiceName, metav1.GetOptions{})
	if err != nil {
		return "", errors.WithStack(err)
	}

	if service.Spec.Type != corev1.ServiceType(nodePort) {
		return "", errors.Errorf("namespace:monitoring,service:prometheus-k8s,type:NodePort")
	}

	if len(service.Spec.Ports) == 0 {
		return "", errors.Errorf("namespace:monitoring,service:prometheus-k8s,type:NodePort")
	}

	url := "http://" + ParseAddrFromURLNoError(clusterInfo.Config.Host) + ":" + strconv.Itoa(int(service.Spec.Ports[0].NodePort))
	return url, nil
}

func ParseAddrFromURLNoError(rawurl string) string {
	addr, _ := ParseAddrFromURL(rawurl)
	return addr
}

func ParseAddrFromURL(rawurl string) (string, error) {
	if !strings.HasPrefix(rawurl, "http://") && !strings.HasPrefix(rawurl, "https://") {
		rawurl = "http://" + rawurl
	}

	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}

	if !strings.Contains(u.Host, ":") {
		return u.Host, nil
	}

	ip, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		return "", err
	}

	return ip, nil
}

func (k *kubernetesService) getPrometheusConfigProqlOption(ctx context.Context, clusterInfo *iapiserver.Cluster, req *iapiserver.Monitoring, namespace string) (*iprometheus.Config, *iprometheus.ProqlOption, error) {
	conf := &iprometheus.Config{}
	opt := &iprometheus.ProqlOption{}
	var err error

	if !clientset.IsMonitorServiceReady(clusterInfo.ID) {
		return nil, nil, errors.Errorf("cluster '%v' not ready", clusterInfo.ID)
	}

	conf.Address, err = k.getPrometheusUrlFromPrometheusService(ctx, clusterInfo)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	opt.Namespace = namespace
	if req != nil {
		opt.Level = req.Level
		opt.Node = req.NodeName
		opt.PodName = req.Name
		opt.Time = time.Unix(req.Time, 0)
		opt.StartTime = time.Unix(req.StartTime, 0)
		opt.EndTime = time.Unix(req.EndTime, 0)
		opt.Duration = time.Duration(req.Duration) * time.Second
	}

	return conf, opt, nil
}

func (k *kubernetesService) MonitorServiceStatusGet(ctx context.Context, req *iapiserver.MonitoringGetRequest) (*iapiserver.MonitoringResponse, error) {
	resp := &iapiserver.MonitoringResponse{MonitorStatus: &iapiserver.MonitorStatus{PrometheusEnable: true, GrafanaEnable: true}}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.MonitorStatus.PrometheusUrl, err = k.getPrometheusUrlFromPrometheusService(ctx, cluster)
	if err != nil {
		resp.MonitorStatus.PrometheusEnable = false
		resp.MonitorStatus.PrometheusError = err.Error()
	}

	resp.MonitorStatus.GrafanaUrl, err = k.getGrafanaUrlFromGrafanaService(ctx, cluster)
	if err != nil {
		resp.MonitorStatus.GrafanaEnable = false
		resp.MonitorStatus.GrafanaError = err.Error()
	}

	return resp, nil
}

func (k *kubernetesService) MonitorMetricsHistoryList(ctx context.Context, req *iapiserver.MonitoringListRequest) (*iapiserver.MonitoringMetricsHistoryListResponse, error) {
	resp := &iapiserver.MonitoringMetricsHistoryListResponse{}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	proConf, opt, err := k.getPrometheusConfigProqlOption(ctx, cluster, req.ToMonitoring(), req.Namespace)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	timeRange := req.ToRange()
	resp.TimeRange = &timeRange
	resp.Metrics, err = libprometheus.GetNamedMetricsHistory(ctx, proConf, req.ToMonitoring().Metrics, opt)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sort.SliceStable(resp.Metrics, func(i, j int) bool {
		return resp.Metrics[i].MetricName < resp.Metrics[j].MetricName
	})

	return nil, errors.WithStack(err)
}

func (k *kubernetesService) MonitorClusterHardwareResourceHistoryList(ctx context.Context, req *iapiserver.MonistoringHardwareHistoryListRequest) (*iapiserver.MonistoringHardwareHistoryListResponse, error) {
	resp := &iapiserver.MonistoringHardwareHistoryListResponse{}

	metrics := []KV{cluster_cpu_ratio, cluster_cpu_load_1, cluster_cpu_load_5, cluster_cpu_load_15,
		cluster_memory_ratio,
		cluster_disk_used_total, cluster_disk_io_rate, cluster_network_receive_rate, cluster_network_send_rate,
		cluster_pod_unknown, cluster_pod_failed, cluster_pod_pending, cluster_pod_succeeded, cluster_pod_running}
	rets, err := k.autoMetricHistoryList(ctx, req.Cluster, req.MonitoringTimeRange.ToRange(), metrics...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	timeRange := req.ToRange()
	resp.TimeRange = &timeRange
	resp.Metrics = make([]iprometheus.Metric, 0, len(rets))
	for _, v := range rets {
		resp.Metrics = append(resp.Metrics, *v)
	}

	sort.SliceStable(resp.Metrics, func(i, j int) bool {
		return resp.Metrics[i].MetricName < resp.Metrics[j].MetricName
	})

	return nil, errors.WithStack(err)
}

func (k *kubernetesService) MonitorClusterNodeResourceHistoryGet(ctx context.Context, req *iapiserver.MonitoringNodeGetRequest) (*iapiserver.MonitoringNodeGetResponse, error) {
	resp := &iapiserver.MonitoringNodeGetResponse{}

	timeRange := req.ToRange()
	history, err := k.autoMetricHistoryList(ctx, req.Cluster, timeRange, nodePhysicalResourceMetrics...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.TimeRange = &timeRange
	for _, v := range history {
		resp.Metrics = append(resp.Metrics, *v)
	}
	return nil, errors.WithStack(err)
}

func (k *kubernetesService) MonitorNamespaceList(ctx context.Context, req *iapiserver.MonitoringNamespaceListRequest) (*iapiserver.MonitoringNamespaceListResponse, error) {
	resp := &iapiserver.MonitoringNamespaceListResponse{}
	if req.SortBy == "" {
		req.SortBy = "cpu"
	}
	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList2[*iapiserver.MonitorData](ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.MonitorData]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.MonitorData](cluster.ID, cluster.Name)

			ret, err := k.MonitorNamespaceStateList(ctx, cluster)
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, nil)
			}
			nsList, err := clientset.NamespaceList(ctx, cluster, metav1.ListOptions{})
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, nil)
			}

			rets := make([]*iapiserver.MonitorData, 0, len(ret))
			for _, ns := range nsList.Items {
				if req.Fuzzy != "" && NewFieldFilter(ns.Name, ns.Namespace, cluster.Name).Filter(req.Fuzzy) {
					continue
				}

				nsInfo, ok := ret[ns.Name]
				if !ok {
					nsInfo = &iapiserver.MonitorData{Name: ns.Name}
					ret[ns.Name] = nsInfo
				}
				nsInfo.Cpu = nsInfo.Cpu / 1000
				rets = append(rets, nsInfo)
			}
			clusterListOne.TotalCount = len(rets)
			clusterListOne.List = rets
			return waitgroup.NewGenericResult(clusterListOne, nil)
		}, req.SortBy, !req.SortDesc, 10*time.Second)
	return resp, err

}

// get all metrics name in promtheus
func (k *kubernetesService) MonitorMetricNameList(ctx context.Context, req *iapiserver.MonitoringMetricsNameListRequest) (*iapiserver.MonitoringMetricsNameListResponse, error) {
	resp := &iapiserver.MonitoringMetricsNameListResponse{}
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	if !clientset.IsMonitorServiceReady(cluster.ID) {
		return nil, errors.Errorf("monitor service not ready")
	}

	prometheusAddress, err := k.getPrometheusUrlFromPrometheusService(ctx, cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	//query all metrics name in promtheus
	ret, err := libprometheus.LabelValues(ctx,
		&iprometheus.Config{Address: prometheusAddress},
		"__name__",
		15)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.TotalCount = len(ret)
	resp.MetricList = make([]string, 0, ret.Len())
	for _, v := range ret {
		resp.MetricList = append(resp.MetricList, string(v))
	}

	return nil, errors.WithStack(err)
}

func (k *kubernetesService) MonitorRuleAutoList(ctx context.Context, req *iapiserver.MonitoringListRequest) (*iapiserver.MonitoringListResponse, error) {
	resp := &iapiserver.MonitoringListResponse{}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	proConf, opt, err := k.getPrometheusConfigProqlOption(ctx, cluster, req.ToMonitoring(), req.Namespace)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.Metrics, err = libprometheus.GetAutoMetricsHistory(ctx, proConf, req.ToMonitoring().Metrics, opt)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp.Metrics = k.metricCombine(resp.Metrics)
	if len(resp.Metrics) != 0 {
		resp.Metric = &resp.Metrics[0]
	}
	resp.Metrics = nil
	return resp, nil
}

type Combiner struct {
	big  bool
	meta map[int64]float64
}

func NewCombinerSmall() *Combiner {
	return &Combiner{big: false, meta: map[int64]float64{}}
}

func NewCombinerBig() *Combiner {
	return &Combiner{big: true, meta: map[int64]float64{}}
}

func (m *Combiner) Push(time int64, value float64) {
	switch m.big {
	case true:
		oldValue, ok := m.meta[time]
		if ok {
			if value > oldValue {
				m.meta[time] = value
			}
		} else {
			m.meta[time] = value
		}

	case false:
		oldValue, ok := m.meta[time]
		if ok {
			if value < oldValue {
				m.meta[time] = value
			}
		} else {
			m.meta[time] = value
		}

	}
}

func (m *Combiner) GetTypeString() string {
	if m.big {
		return "max"
	}
	return "least"
}

func (m *Combiner) GetInsert() []iprometheus.Point {
	resp := make([]iprometheus.Point, len(m.meta), len(m.meta))

	var i, j int
	for k, v := range m.meta {
		resp[i][0], resp[i][1], j = float64(k), v, i
		for j > 0 && resp[j][0] < resp[j-1][0] {
			resp[j], resp[j-1] = resp[j-1], resp[j]
			j = j - 1
		}
		i++
	}

	return resp
}

func (m *Combiner) GetMax() float64 {
	if len(m.meta) == 0 {
		return 0
	}

	isFirst := true
	max := float64(0)
	for _, v := range m.meta {
		if isFirst {
			max = v
			isFirst = false
			continue
		}
		if v > max {
			max = v
		}
	}

	return max
}

func (m *Combiner) GetQuick() []iapiserver.Point {
	resp := make([]iapiserver.Point, len(m.meta), len(m.meta))

	i := 0
	for k, v := range m.meta {
		p := iapiserver.Point{}
		p[0] = float64(k)
		p[1] = v
		resp[i] = p
	}

	sort.Slice(resp, func(i, j int) bool {
		return resp[i][0] < resp[j][0]
	})

	return resp
}

func (k *kubernetesService) metricCombine(ms []iprometheus.Metric) []iprometheus.Metric {
	resp := make([]iprometheus.Metric, len(ms), len(ms))

	for k, metricv := range ms {
		tmp := iprometheus.Metric{MetricName: metricv.MetricName}
		tmp.MetricValues = make([]iprometheus.MetricValue, 2, 2)
		tmp.MetricType = metricv.MetricType
		tmp.Error = metricv.Error

		cb := NewCombinerBig()
		cs := NewCombinerSmall()
		for _, mcv := range metricv.MetricValues {
			for _, pv := range mcv.Series {
				cb.Push(int64(pv.Timestamp()), pv.Value())
				cs.Push(int64(pv.Timestamp()), pv.Value())
			}
		}

		tmp.MetricValues[0].Series = cs.GetInsert()
		tmp.MetricValues[0].Name = cs.GetTypeString()
		tmp.MetricValues[0].Max = cs.GetMax()
		tmp.MetricValues[1].Series = cb.GetInsert()
		tmp.MetricValues[1].Name = cb.GetTypeString()
		tmp.MetricValues[1].Max = cb.GetMax()
		resp[k] = tmp
	}

	return resp
}

func (k *kubernetesService) MonitorAlertList(ctx context.Context, req *iapiserver.MonitoringListRequest) (*iapiserver.MonitoringAlertListResponse, error) {
	resp := &iapiserver.MonitoringAlertListResponse{}

	if req.SortBy == "" {
		req.SortBy = "ClusterName Namespace ActiveTime"
	}

	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList2[*iapiserver.Alert](ctx, k.store, &resp.AlertList, req.ResourceListRequest,
		func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.Alert]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.Alert](cluster.ID, cluster.Name)

			proConf, _, err := k.getPrometheusConfigProqlOption(ctx, cluster, req.ToMonitoring(), req.Namespace)
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}
			ret, err := libprometheus.AlertList(ctx, proConf, 15)
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			alertList := resp.AlertList
			for _, v := range ret.Alerts {
				v := v

				alert := iapiserver.Alert{}
				alert.Namespace = string(v.Labels["namespace"])
				alert.Name = string(v.Labels["alertname"])
				alert.Message = string(v.Annotations["message"])
				alert.ActiveTime = v.ActiveAt.Unix()
				alert.Pod = string(v.Labels["pod"])
				alert.Container = string(v.Labels["container"])
				alert.Severity = alertLevelValueDesc[string(v.Labels["severity"])]

				v2, _ := strconv.ParseFloat(v.Value, 64)
				alert.Value = fmt.Sprintf("%.6f", v2)
				alert.Alertstate = string(v.State)
				alert.Metadata = &v
				if req.Namespace != "" && req.Namespace != alert.Namespace {
					continue
				}

				if req.Fuzzy != "" && !sets.NewString(alert.Namespace, alert.Name, alert.Pod, alert.Container,
					alert.Alertstate, alert.Severity, alert.ClusterName,
					alert.Value, alert.Message).ContainAny(req.FuzzyFields()...) {
					continue
				}
				alertList = append(alertList, &alert)
			}

			clusterListOne.TotalCount = len(alertList)
			clusterListOne.List = alertList
			return waitgroup.NewGenericResult(clusterListOne, nil)
		}, req.SortBy, !req.SortDesc, 10*time.Second)
	return resp, err

}

func (k *kubernetesService) MonitorAlertDelete(ctx context.Context, req *iapiserver.MonitoringRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	proConf, _, err := k.getPrometheusConfigProqlOption(ctx, cluster, req.Monitoring, req.Namespace)
	if err != nil {
		return errors.WithStack(err)
	}

	ret, err := libprometheus.AlertList(ctx, proConf, 15)
	if err != nil {
		return errors.WithStack(err)
	}
	deleteAlert := req.AlertDeleteList[0]

	for _, v := range ret.Alerts {
		if v.Labels.Equal(deleteAlert.Labels) {
			if err := libprometheus.DeleteSeries(ctx, proConf, []string{"ALERTS" + deleteAlert.Labels.String()}, v.ActiveAt); err != nil {
				return errors.WithStack(err)
			}
		}
	}
	return nil
}

func (k *kubernetesService) MonitorPodList(ctx context.Context, req *iapiserver.MonitoringPodListRequest) (*iapiserver.MonitoringPodListResponse, error) {
	resp := &iapiserver.MonitoringPodListResponse{}
	if req.SortBy == "" {
		req.SortBy = "cluster_name namespace name cpu.used"
	}
	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList2[*iapiserver.PodOverview](ctx, k.store, &resp.PodList, req.ResourceListRequest,
		func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.PodOverview]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.PodOverview](cluster.ID, cluster.Name)

			nodeList, err := clientset.NodeList(ctx, cluster, metav1.ListOptions{})
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, nil)
			}
			podList, err := clientset.PodList(ctx, cluster, req.Namespace, metav1.ListOptions{}) // local cache
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, nil)
			}
			nsPods, err := k.MonitorClusterPodStateList(ctx, cluster) // k8s
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, nil)
			}

			nodes := map[string]*corev1.Node{} // node info
			for i := range nodeList.Items {
				nodes[nodeList.Items[i].Name] = &nodeList.Items[i]
			}
			rets := make([]*iapiserver.PodOverview, 0, len(podList.Items))
			for _, v := range podList.Items { // all pod
				info := nsPods[v.Namespace][v.Name]
				if info == nil {
					info = &iapiserver.PodOverview{Namespace: v.Namespace, Name: v.Name}
					info.NodeName = v.Status.HostIP
				}
				info.NodeAddr = getNodeAddr(nodes[info.NodeName])
				info.CreateTime = v.CreationTimestamp.Unix()

				if !sets.NewString(info.Name, info.Namespace, info.ClusterName, info.NodeName, info.NodeAddr).ContainAny(req.FuzzyFields()...) {
					continue
				}

				rets = append(rets, info)
			}
			clusterListOne.TotalCount = len(rets)
			clusterListOne.List = rets
			return waitgroup.NewGenericResult(clusterListOne, nil)
		}, req.SortBy, !req.SortDesc, 10*time.Second)
	return resp, err
}

func (k *kubernetesService) MonitorNamespaceTopGet(ctx context.Context, req *iapiserver.MonitoringTopGetRequest) (*iapiserver.MonitoringTopGetResponse, error) {
	resp := &iapiserver.MonitoringTopGetResponse{}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	stateList, err := k.namespaceOverview(ctx, cluster, req.Namespace)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.TotalCount = len(stateList)
	s, e := paging.Index(len(stateList), req.PageNum, req.PageSize)
	resp.MonitorDatas = stateList[s:e]
	return nil, errors.WithStack(err)
}

func (k *kubernetesService) MonitorStatefulSetState(ctx context.Context, clusterInfo *iapiserver.Cluster, namespace string) ([]*iapiserver.MonitorData, []string, error) {
	workLoadState := []*iapiserver.MonitorData{}
	seletedPod := []string{}
	var err error

	podList := &iapiserver.PodListResponse{}
	workloadList := &iapiserver.StatefulSetListResponse{}

	{
		wg := waitgroup.NewWaitGroup(nil)
		wg.Start(waitgroup.NewWaitGroupHandleFunc("", func(ctx context.Context) waitgroup.Result {
			var err error
			podList, err = k.PodList(ctx, &iapiserver.PodListRequest{
				ResourceListRequest: iapiserver.ResourceListRequest{
					Cluster:   clusterInfo.ID,
					Namespace: namespace,
				},
			})
			return waitgroup.NewResult(nil, err)
		}))

		wg.Start(waitgroup.NewWaitGroupHandleFunc("", func(ctx context.Context) waitgroup.Result {
			var err error
			workloadList, err = k.StatefulSetList(ctx, &iapiserver.StatefulSetListRequest{
				ResourceListRequest: iapiserver.ResourceListRequest{
					Cluster:   clusterInfo.ID,
					Namespace: namespace,
				}})
			return waitgroup.NewResult(nil, err)
		}))
		wg.Wait()

		for _, v := range wg.GetResults() {
			if v.Error != nil {
				return nil, nil, err
			}
		}
	}
	lock := sync.RWMutex{}

	waitgroup.RunConcurrently(ctx, workloadList.List, func(ctx context.Context, one *iapiserver.StatefulSetInfo) waitgroup.Result {
		objs := []metav1.ObjectMeta{}
		for _, pod := range podList.List {
			if pod != nil && pod.Resource != nil && pod.Resource.Labels != nil {
				if maputil.Equal(pod.Resource.Labels, one.Resource.Spec.Selector.MatchLabels) {
					objs = append(objs, pod.Resource.ObjectMeta)
					lock.Lock()
					seletedPod = append(seletedPod, pod.Resource.Name)
					lock.Unlock()
				}
			}
		}

		ret, err := k.MonitorPodState(ctx, clusterInfo, objs)
		if err != nil {
			return waitgroup.NewResult(nil, err)
		}

		ret.Type = daemonset
		ret.Namespace = one.Resource.Namespace
		ret.Name = one.Resource.Name
		lock.Lock()
		workLoadState = append(workLoadState, ret)
		lock.Unlock()
		return waitgroup.NewResult(nil, nil)
	})

	return workLoadState, seletedPod, nil
}

func (k *kubernetesService) MonitorDaemonsetState(ctx context.Context, clusterInfo *iapiserver.Cluster, namespace string) ([]*iapiserver.MonitorData, []string, error) {
	workLoadState := []*iapiserver.MonitorData{}
	seletedPod := []string{}
	var err error

	podListResp := &iapiserver.PodListResponse{}
	workloadList := &iapiserver.DaemonSetListResponse{}
	{
		wg := waitgroup.NewWaitGroup(nil)
		wg.Start(waitgroup.NewWaitGroupHandleFunc("", func(ctx context.Context) waitgroup.Result {
			var err error
			podListResp, err = k.PodList(ctx, &iapiserver.PodListRequest{
				ResourceListRequest: iapiserver.ResourceListRequest{
					Cluster:   clusterInfo.ID,
					Namespace: namespace,
				},
			})
			return waitgroup.NewResult(nil, err)
		}))

		wg.Start(waitgroup.NewWaitGroupHandleFunc("", func(ctx context.Context) waitgroup.Result {
			var err error
			workloadList, err = k.DaemonSetList(ctx, &iapiserver.DaemonSetListRequest{
				ResourceListRequest: iapiserver.ResourceListRequest{
					Cluster:   clusterInfo.ID,
					Namespace: namespace,
				}})
			return waitgroup.NewResult(nil, err)
		}))
		wg.Wait()

		for _, v := range wg.GetResults() {
			if v.Error != nil {
				return nil, nil, err
			}
		}
	}
	lock := sync.RWMutex{}

	waitgroup.RunConcurrently(ctx, workloadList.List, func(ctx context.Context, one *iapiserver.DaemonSetInfo) waitgroup.Result {
		objs := []metav1.ObjectMeta{}
		for _, pod := range podListResp.List {
			if pod != nil && pod.Resource != nil && pod.Resource.Labels != nil {
				if maputil.Equal(pod.Resource.Labels, one.Resource.Spec.Selector.MatchLabels) {
					objs = append(objs, pod.Resource.ObjectMeta)
					lock.Lock()
					seletedPod = append(seletedPod, pod.Resource.Name)
					lock.Unlock()
				}
			}
		}

		ret, err := k.MonitorPodState(ctx, clusterInfo, objs)
		if err != nil {
			return waitgroup.NewResult(nil, err)
		}

		ret.Type = daemonset
		ret.Namespace = one.Resource.Namespace
		ret.Name = one.Resource.Name
		lock.Lock()
		workLoadState = append(workLoadState, ret)
		lock.Unlock()
		return waitgroup.NewResult(nil, nil)
	})

	return workLoadState, seletedPod, nil
}

func (k *kubernetesService) MonitorNodeList(ctx context.Context, req *iapiserver.MonitoringListRequest) (*iapiserver.MonitoringNodeListResponse, error) {
	resp := &iapiserver.MonitoringNodeListResponse{}
	if req.SortBy == "" {
		req.SortBy = "CpuUsed"
	}
	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList2[*iapiserver.NodeOverview](ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.NodeOverview]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.NodeOverview](cluster.ID, cluster.Name)

			noMap, err := k.getClusterNodesOverview(ctx, cluster)
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, nil)
			}
			rets := make([]*iapiserver.NodeOverview, 0, len(noMap))
			for _, v := range noMap {
				if !sets.NewString(v.NodeName, v.ClusterName, v.NodeAddr).ContainAny(req.FuzzyFields()...) {
					continue
				}
				rets = append(rets, v)
			}
			clusterListOne.TotalCount = len(rets)
			clusterListOne.List = rets
			return waitgroup.NewGenericResult(clusterListOne, nil)
		}, req.SortBy, !req.SortDesc, 10*time.Second)
	return resp, err
}

/*
*
返回整个集群的pod的cpu memory netread netwrite
返回值对应关系 ： namespace podname podinfo
没有namespace或者podname的container不会被返回
*/
func (k *kubernetesService) podRealtimeInfo(ctx context.Context, cluster *iapiserver.Cluster) (map[string]map[string]*iapiserver.PodOverview, error) {
	cpuMetrics := &iprometheus.Metric{}
	memoryMetrics := &iprometheus.Metric{}
	netReadMetrics := &iprometheus.Metric{}
	netWriteMetrics := &iprometheus.Metric{}

	promConf, _, err := k.getPrometheusConfigProqlOption(ctx, cluster, nil, "")
	if err != nil {
		return nil, err
	}

	wg := waitgroup.NewWaitGroup(nil)
	wg.Start(waitgroup.NewWaitGroupHandleFunc("", func(ctx context.Context) waitgroup.Result {
		var err error
		cpuMetrics, err = libprometheus.GetAutoMetricsRealTime(ctx, promConf, "rate(container_cpu_usage_seconds_total{}[1m])")
		return waitgroup.NewResult(nil, err)
	}))
	wg.Start(waitgroup.NewWaitGroupHandleFunc("", func(ctx context.Context) waitgroup.Result {
		var err error
		memoryMetrics, err = libprometheus.GetAutoMetricsRealTime(ctx, promConf, "container_memory_usage_bytes{}")
		return waitgroup.NewResult(nil, err)
	}))
	wg.Start(waitgroup.NewWaitGroupHandleFunc("", func(ctx context.Context) waitgroup.Result {
		var err error
		netReadMetrics, err = libprometheus.GetAutoMetricsRealTime(ctx, promConf, "rate(container_network_receive_bytes_total{}[1m])")
		return waitgroup.NewResult(nil, err)
	}))
	wg.Start(waitgroup.NewWaitGroupHandleFunc("", func(ctx context.Context) waitgroup.Result {
		var err error
		netWriteMetrics, err = libprometheus.GetAutoMetricsRealTime(ctx, promConf, "rate(container_network_transmit_bytes_total{}[1m])")
		return waitgroup.NewResult(nil, err)
	}))

	wg.Wait()

	var namespace string
	var podname string
	nsNamepod := map[string]map[string]*iapiserver.PodOverview{}
	callBack := func(v *iprometheus.MetricValue) (*iapiserver.PodOverview, float64) {
		var value float64
		if v == nil {
			return nil, 0
		}

		namespace = v.Metadata["namespace"]
		podname = v.Metadata["pod"]
		if podname == "" || namespace == "" {
			return nil, 0
		}

		namePod, ok := nsNamepod[namespace]
		if !ok {
			namePod = map[string]*iapiserver.PodOverview{}
			nsNamepod[namespace] = namePod
		}
		pod, ok := namePod[podname]
		if !ok {
			pod = &iapiserver.PodOverview{}
			namePod[podname] = pod
		}
		pod.Namespace = namespace
		pod.Name = podname
		if v.Sample != nil {
			value = v.Sample.Value()
		}

		return pod, value
	}

	if cpuMetrics != nil {
		for _, v := range cpuMetrics.MetricValues { // cpu
			podOver, value := callBack(&v)
			if podOver != nil {
				podOver.Cpu.Used += value
			}
		}
	}

	if memoryMetrics != nil {
		for _, v := range memoryMetrics.MetricValues { // memory
			podOver, value := callBack(&v)
			if podOver != nil {
				podOver.Memory.Used += value
			}
		}
	}

	if netReadMetrics != nil {
		for _, v := range netReadMetrics.MetricValues { // net read
			podOver, value := callBack(&v)
			if podOver != nil {
				podOver.NetRead.Used += value
			}
		}
	}

	if netWriteMetrics != nil {
		for _, v := range netWriteMetrics.MetricValues { // net write
			podOver, value := callBack(&v)
			if podOver != nil {
				podOver.NetWrite.Used += value
			}
		}
	}

	return nsNamepod, nil
}

func (k *kubernetesService) MonitorPodContainerHistoryList(ctx context.Context, req *iapiserver.MonistoringPodContainerHistoryListRequest) (*iapiserver.MonistoringPodContainerHistoryListResponse, error) {
	resp := &iapiserver.MonistoringPodContainerHistoryListResponse{}

	timeRange := req.MonitoringTimeRange.ToRange()
	metricValues, err := k.autoMetricHistoryList(ctx, req.Cluster, timeRange, combineSql(req.Name, podContainerMetrics...)...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.TimeRange = &timeRange
	resp.Metrics = make([]iprometheus.Metric, 0, len(metricValues))
	for _, v := range metricValues {
		resp.Metrics = append(resp.Metrics, *v)
	}

	return resp, nil
}

func (k *kubernetesService) MonitorControllerManagerGet(ctx context.Context, req *iapiserver.MonitoringGetRequest) (*iapiserver.MonitoringControllerManagerResponse, error) {
	resp := &iapiserver.MonitoringControllerManagerResponse{}
	timeRange := req.MonitoringTimeRange.ToRange()
	metricValues, err := k.autoMetricHistoryList(ctx, req.Cluster, timeRange, controllerManagerQueryMetrics...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.TimeRange = &timeRange
	resp.UpCount = int(metricValues[cm_up_count.Name].LastSum())
	resp.WorkQueueAddRate = *metricValues[cm_work_queue_add_rate.Name].RangeAvg()
	resp.CpuUsed = *metricValues[cm_cpu_used.Name].RangeAvg()
	resp.MemoryUsed = *metricValues[cm_memory_used.Name].RangeAvg()
	resp.Goroutine = *metricValues[cm_goroutines.Name].RangeAvg()
	return resp, nil
}

func (k *kubernetesService) MonitorSchedulerGet(ctx context.Context, req *iapiserver.MonitoringGetRequest) (*iapiserver.MonitoringSchedulerResponse, error) {
	resp := &iapiserver.MonitoringSchedulerResponse{}
	timeRange := req.MonitoringTimeRange.ToRange()
	metricValues, err := k.autoMetricHistoryList(ctx, req.Cluster, timeRange, schedulerQueryMetrics...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.TimeRange = &timeRange
	resp.UpCount = int(getMetricsValueFromMetrics(metricValues, scheduler_up_count.Name).LastSum())
	resp.ScheduleCount = *getMetricsValueFromMetrics(metricValues, scheduler_schedule_count.Name).RangeAvg()
	resp.CpuUsed = *getMetricsValueFromMetrics(metricValues, scheduler_cpu_used.Name).RangeAvg()
	resp.MemoryUsed = *getMetricsValueFromMetrics(metricValues, scheduler_memory_used.Name).RangeAvg()
	resp.Goroutine = *getMetricsValueFromMetrics(metricValues, scheduler_goroutines.Name).RangeAvg()
	return resp, nil
}

func (k *kubernetesService) MonitorApiserverGet(ctx context.Context, req *iapiserver.MonitoringGetRequest) (*iapiserver.MonitoringApiserverResponse, error) {
	resp := &iapiserver.MonitoringApiserverResponse{}

	timeRange := req.MonitoringTimeRange.ToRange()
	metricValues, err := k.autoMetricHistoryList(ctx, req.Cluster, timeRange, apiserverQueryMetrics...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.TimeRange = &timeRange
	resp.CpuUsed = *metricValues[apiserver_cpu_used.Name].RangeAvg()
	resp.MemoryUsed = *metricValues[apiserver_memory_used.Name].RangeAvg()
	resp.Goroutine = *metricValues[apiserver_goroutines.Name].RangeAvg()
	return resp, nil
}

func (k *kubernetesService) MonitorKubeletGet(ctx context.Context, req *iapiserver.MonitoringGetRequest) (*iapiserver.MonitoringKubeletResponse, error) {
	resp := &iapiserver.MonitoringKubeletResponse{}

	//当前场景为主机名
	// if req.Name == "" {
	// 	return nil, errors.Errorf("monitoring.name is empty")
	// }

	timeRange := req.MonitoringTimeRange.ToRange()
	metricValues, err := k.autoMetricHistoryList(ctx, req.Cluster, timeRange, combineSql(req.Name, kubeletQueryMetrics...)...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.TimeRange = &timeRange
	resp.UpCount = int(metricValues[kubelet_up_count.Name].LastSum())
	resp.PodCount = int(metricValues[kubelet_running_pod_count.Name].LastSum())
	resp.RpcRate = *metricValues[kubelet_rpc_rate.Name].RangeAvg()
	return resp, nil
}

func (k *kubernetesService) autoMetricHistoryList(ctx context.Context, clusterUUID string, timeRange prometheusv1.Range, metrics ...KV) (map[string]*iprometheus.Metric, error) {
	resp := map[string]*iprometheus.Metric{}

	cluster, err := k.store.Kubernetes().Get(ctx, clusterUUID)
	if err != nil {
		return nil, err
	}

	if !clientset.IsMonitorServiceReady(cluster.ID) {
		return nil, errors.Errorf("monitor service not healthy")
	}
	prometheusAddress, err := k.getPrometheusUrlFromPrometheusService(ctx, cluster)
	if err != nil {
		return nil, err
	}

	wg := waitgroup.RunGenericConcurrently(ctx, metrics, func(ctx context.Context, metric KV) waitgroup.GenericResult[*iprometheus.Metric] {
		metricValue := &iprometheus.Metric{MetricName: metric.Name}
		metrics, err := libprometheus.GetAutoMetricsHistory2(ctx, &iprometheus.Config{Address: prometheusAddress}, metric.Value, timeRange)
		if err != nil {
			metricValue.Error = err.Error()
		}
		metricValue.MetricValues = metrics.MetricValues
		return waitgroup.NewGenericResult(metrics, nil)
	})
	for _, v := range wg.GetResults() {
		resp[v.Data.MetricName] = v.Data
	}
	return resp, nil
}

func (k *kubernetesService) autoMetricRealtimeList(ctx context.Context, clusterInfo *iapiserver.Cluster, metrics ...KV) (map[string]*iprometheus.Metric, error) {
	resp := map[string]*iprometheus.Metric{}

	proConf, _, err := k.getPrometheusConfigProqlOption(ctx, clusterInfo, nil, "")
	if err != nil {
		return nil, err
	}

	wg := waitgroup.RunGenericConcurrently(ctx, metrics, func(ctx context.Context, metric KV) waitgroup.GenericResult[*iprometheus.Metric] {
		metricValue := &iprometheus.Metric{MetricName: metric.Name}
		metrics, err := libprometheus.GetAutoMetricsRealTime(ctx, proConf, metric.Value, 5)
		if err != nil {
			metricValue.Error = err.Error()
		}
		metricValue.MetricValues = metrics.MetricValues
		return waitgroup.NewGenericResult(metrics, nil)
	})
	for _, v := range wg.GetResults() {
		resp[v.Data.MetricName] = v.Data
	}

	return resp, nil
}

func (k *kubernetesService) MonitorEtcdGet(ctx context.Context, req *iapiserver.MonitoringGetRequest) (*iapiserver.MonitoringEtcdResponse, error) {
	resp := &iapiserver.MonitoringEtcdResponse{}
	timeRange := req.MonitoringTimeRange.ToRange()
	metricValues, err := k.autoMetricHistoryList(ctx, req.Cluster, timeRange, etcd_has_leader, etcd_change_total,
		etcd_net_read_byte, etcd_net_write_byte, etcd_db_size, etcd_raft_proposal_applied, etcd_raft_proposal_commited,
		etcd_raft_proposal_failed, etcd_raft_proposal_pending)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.TimeRange = &timeRange
	if metricValues[etcd_has_leader.Name].LastSum() > 0 {
		resp.HasLeader = true
	}
	resp.LeaderChangeCount = int(getMetricsValueFromMetrics(metricValues, etcd_change_total.Name).LastSum())
	resp.NetReadByte = *getMetricsValueFromMetrics(metricValues, etcd_net_read_byte.Name).RangeAvg()
	resp.NetWriteByte = *getMetricsValueFromMetrics(metricValues, etcd_net_write_byte.Name).RangeAvg()
	resp.MemorySize = *getMetricsValueFromMetrics(metricValues, etcd_db_size.Name).RangeAvg()
	resp.RaftProposalApplied = *getMetricsValueFromMetrics(metricValues, etcd_raft_proposal_applied.Name).RangeAvg()
	resp.RaftProposalCommited = *getMetricsValueFromMetrics(metricValues, etcd_raft_proposal_commited.Name).RangeAvg()
	resp.RaftProposalFailed = *getMetricsValueFromMetrics(metricValues, etcd_raft_proposal_failed.Name).RangeAvg()
	resp.RaftProposalPending = *getMetricsValueFromMetrics(metricValues, etcd_raft_proposal_pending.Name).RangeAvg()

	return resp, nil
}

func getMetricsValueFromMetrics(metrics map[string]*iprometheus.Metric, metricName string) *iprometheus.Metric {
	metric, _ := metrics[metricName]
	return metric
}

func (k *kubernetesService) getMetricRange(metric *iprometheus.Metric) float64 {
	if metric == nil || len(metric.MetricValues) == 0 || len(metric.MetricValues[0].Series) == 0 {
		return 0
	}
	return metric.MetricValues[0].Series[len(metric.MetricValues[0].Series)-1].Value()
}

func (k *kubernetesService) MonitorNodeMasterGet(ctx context.Context, req *iapiserver.MonitoringGetRequest) (*iapiserver.MonitoringResponse, error) {
	resp := &iapiserver.MonitoringResponse{NodeUsedTotal: &iapiserver.NodeOverview{}}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	rets, err := k.getClusterNodesOverview(ctx, cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, ret := range rets {
		if ret.Role != iapiserver.NodeRoleMaster {
			continue
		}
		ret.CpuUsedRatio = mathutil.Divide(ret.CpuUsedRatio, 100)
		ret.MemoryUsedRatio = mathutil.Divide(ret.MemoryUsedRatio, 100)
		ret.DiskUsedRatio = mathutil.Divide(ret.DiskUsedRatio, 100)
		ret.PodUsedRatio = mathutil.Divide(ret.PodUsedRatio, 100)
		resp.NodeList = append(resp.NodeList, ret)
		resp.NodeUsedTotal.CpuCapacity += ret.CpuCapacity
		resp.NodeUsedTotal.CpuUsed += ret.CpuUsed
		resp.NodeUsedTotal.MemoryCapacity += ret.MemoryCapacity
		resp.NodeUsedTotal.MemoryUsed += ret.MemoryUsed
		resp.NodeUsedTotal.DiskCapacity += ret.DiskCapacity
		resp.NodeUsedTotal.DiskUsed += ret.DiskUsed
		resp.NodeUsedTotal.PodCapacity += ret.PodCapacity
		resp.NodeUsedTotal.PodUsed += ret.PodUsed
	}

	resp.TotalCount = len(resp.NodeList)
	resp.NodeUsedTotal.CpuUsedRatio = mathutil.Divide(resp.NodeUsedTotal.CpuUsed, resp.NodeUsedTotal.CpuCapacity)
	resp.NodeUsedTotal.MemoryUsedRatio = mathutil.Divide(resp.NodeUsedTotal.MemoryUsed, resp.NodeUsedTotal.MemoryCapacity)
	resp.NodeUsedTotal.DiskUsedRatio = mathutil.Divide(resp.NodeUsedTotal.DiskUsed, resp.NodeUsedTotal.DiskCapacity)
	resp.NodeUsedTotal.PodUsedRatio = mathutil.Divide(float64(resp.NodeUsedTotal.PodUsed), float64(resp.NodeUsedTotal.PodCapacity))

	return resp, nil
}

func (k *kubernetesService) NodeState(ctx context.Context, clusterInfo *iapiserver.Cluster, nodeName string, timeout int64) (*iapiserver.MonitorData, error) {
	resp := &iapiserver.MonitorData{Name: nodeName}

	rets, err := k.getClusterNodesOverview(ctx, clusterInfo)
	if err != nil {
		return nil, err
	}
	nodeOverview, ok := rets[nodeName]
	if !ok {
		return resp, nil
	}

	resp.Cpu = nodeOverview.CpuUsedRatio
	resp.CpuUsed = nodeOverview.CpuUsed
	resp.CpuTotal = nodeOverview.CpuCapacity
	resp.Memory = nodeOverview.MemoryUsedRatio
	resp.MemoryUsed = nodeOverview.MemoryUsed
	resp.MemoryTotal = nodeOverview.MemoryCapacity
	resp.DiskRatio = nodeOverview.DiskUsedRatio
	resp.DiskUsed = nodeOverview.DiskUsed
	resp.DiskTotal = nodeOverview.DiskCapacity
	resp.PodRatio = nodeOverview.PodUsedRatio
	resp.PodCount = float64(nodeOverview.PodUsed)
	resp.PodTotal = float64(nodeOverview.PodCapacity)

	return resp, nil
}

func (k *kubernetesService) MonitorWorkloadStateList(ctx context.Context, req *iapiserver.MonitoringListRequest) (*iapiserver.MonitoringListResponse, error) {
	resp := &iapiserver.MonitoringListResponse{}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	var ls labels.Selector
	switch req.Type {
	case deployment:
		wl, err := clientset.DeploymentGet(ctx, cluster, req.Namespace, req.Name, metav1.GetOptions{})
		if err != nil {
			return nil, errors.WithStack(err)
		}
		ls, _ = metav1.LabelSelectorAsSelector(wl.Spec.Selector)
	case statefulset:
		wl, err := clientset.StatefulSetGet(ctx, cluster, req.Namespace, req.Name, metav1.GetOptions{})
		if err != nil {
			return nil, errors.WithStack(err)
		}

		ls, _ = metav1.LabelSelectorAsSelector(wl.Spec.Selector)
	case daemonset:
		wl, err := clientset.DaemonSetGet(ctx, cluster, req.Namespace, req.Name, metav1.GetOptions{})
		if err != nil {
			return nil, errors.WithStack(err)
		}
		ls, _ = metav1.LabelSelectorAsSelector(wl.Spec.Selector)
	default:
		return nil, errors.Errorf("monitor.type not support")
	}

	log.Infof("workload label selectors:%v", ls.String())
	podList, err := clientset.PodList(ctx, cluster, req.Namespace, metav1.ListOptions{LabelSelector: ls.String()})
	if err != nil {

		return nil, errors.WithStack(err)
	}

	if len(podList.Items) == 0 {
		return nil, errors.Errorf(("podList Items is empty"))
	}

	podStateList, err := k.MonitorPodStateList(ctx, cluster, req.ToMonitoring())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	pods, ok := podStateList[req.Namespace]
	if !ok {
		return nil, errors.Errorf(("podStateList Items is empty"))
	}
	for _, v := range podList.Items {
		if podMetrics, ok := pods[v.Name]; ok {
			resp.Metrics = append(resp.Metrics, podMetrics)
		}
	}

	return resp, nil
}

func (k *kubernetesService) MonitorDeploymentState(ctx context.Context, clusterInfo *iapiserver.Cluster, namespace string) ([]*iapiserver.MonitorData, []string, error) {
	workLoadState := []*iapiserver.MonitorData{}
	seletedPod := []string{}
	if namespace == "" {
		return nil, nil, errors.Errorf("namespace is empty")
	}

	podList := &iapiserver.PodListResponse{}
	workloadList := &iapiserver.DeploymentListResponse{}

	{
		wg := waitgroup.NewWaitGroup(nil)
		wg.Start(waitgroup.NewWaitGroupHandleFunc("", func(ctx context.Context) waitgroup.Result {
			var err error
			podList, err = k.PodList(ctx, &iapiserver.PodListRequest{
				ResourceListRequest: iapiserver.ResourceListRequest{
					Cluster:   clusterInfo.ID,
					Namespace: namespace,
				},
			})

			return waitgroup.NewResult(nil, err)
		}))

		wg.Start(waitgroup.NewWaitGroupHandleFunc("", func(ctx context.Context) waitgroup.Result {
			in := &iapiserver.DeploymentListRequest{}
			in.Cluster = clusterInfo.ID
			in.Namespace = namespace
			var err error
			workloadList, err = k.DeploymentList(ctx, in)
			return waitgroup.NewResult(nil, err)
		}))
		wg.Wait()

		for _, v := range wg.GetResults() {
			if v.Error != nil {
				return nil, nil, v.Error
			}
		}
	}
	lock := sync.RWMutex{}
	waitgroup.RunConcurrently(ctx, workloadList.List, func(ctx context.Context, one *iapiserver.DeploymentInfo) waitgroup.Result {
		var objs []metav1.ObjectMeta
		for _, pod := range podList.List {
			if pod.Resource == nil || one.Resource == nil || one.Resource.Spec.Selector == nil {
				continue
			}

			if maputil.Equal(pod.Resource.Labels, one.Resource.Spec.Selector.MatchLabels) {
				objs = append(objs, pod.Resource.ObjectMeta)
				lock.Lock()
				seletedPod = append(seletedPod, pod.Resource.Name)
				lock.Unlock()
			}
		}

		ret, err := k.MonitorPodState(ctx, clusterInfo, objs)
		if err != nil {
			return waitgroup.NewResult(nil, err)
		}
		ret.Type = deployment
		ret.Namespace = one.Resource.Namespace
		ret.Name = one.Resource.Name

		lock.Lock()
		workLoadState = append(workLoadState, ret)
		lock.Unlock()
		return waitgroup.NewResult(nil, nil)
	})

	return workLoadState, seletedPod, nil
}

func (k *kubernetesService) MonitorPodStateList(ctx context.Context, clusterInfo *iapiserver.Cluster, req *iapiserver.Monitoring) (map[string]map[string]iprometheus.Metric, error) {
	resp := map[string]map[string]iprometheus.Metric{}

	metrics := []KV{
		namespace_pod_cpu_rate,
		namespace_pod_memory_size,
		namespace_pod_network_receive_size,
		namespace_pod_network_send_size,
	}
	rets, err := k.autoMetricHistoryList(ctx, clusterInfo.ID, req.MonitoringTimeRange.ToRange(), metrics...)
	if err != nil {
		return nil, err
	}

	for _, v := range rets {
		for _, podvalue := range v.MetricValues {
			namespace := podvalue.Metadata["namespace"]
			podname := podvalue.Metadata["pod"]

			if namespace == "" || podname == "" {
				continue
			}

			pods, ok := resp[namespace]
			if !ok {
				pods = map[string]iprometheus.Metric{}
				resp[namespace] = pods
			}
			newm := pods[podname]
			newm.MetricValues = append(newm.MetricValues, podvalue)
			pods[podname] = newm
		}
	}

	return resp, nil
}

func (k *kubernetesService) MonitorNodeStateList(ctx context.Context, clusterInfo *iapiserver.Cluster) (map[string]*iapiserver.NodeOverview, error) {
	resp := map[string]*iapiserver.NodeOverview{}

	getNodeName := func(labels map[string]string) string {
		if labels == nil {
			return ""
		}

		if nodename1, _ := labels["node"]; nodename1 != "" {
			return nodename1
		}
		instance, _ := labels["instance"]
		return instance
	}

	metrics := []KV{node_cpu_used, node_cpu_total, node_cpu_ratio, node_cpu_request_used, node_cpu_limit_used,
		node_memory_used, node_memory_total, node_memory_ratio, node_memory_request_used, node_memory_limit_used,
		node_load_1, node_load_5, node_load_15, node_disk_io_rate, node_disk_ratio, node_disk_used, node_disk_total,
		node_disk_read_rate, node_disk_write_rate,
		node_pod_ratio, node_pod_used, node_pod_total, node_pod_running}
	ret, err := k.autoMetricRealtimeList(ctx, clusterInfo, metrics...)
	if err != nil {
		return nil, err
	}

	for _, v := range metrics {
		metricInfo := ret[v.Name]
		if metricInfo == nil {
			continue
		}
		for _, nsInfo := range metricInfo.MetricValues {
			nodeName := getNodeName(nsInfo.Metadata)
			if nodeName == "" {
				continue
			}
			data, ok := resp[nodeName]
			if !ok {
				data = &iapiserver.NodeOverview{NodeName: nodeName}
				resp[nodeName] = data
			}
			if nsInfo.Sample == nil {
				continue
			}
			switch v.Name {
			case node_cpu_used.Name:
				data.CpuUsed = nsInfo.Sample.Value() * 1000
			case node_cpu_total.Name:
				data.CpuCapacity = nsInfo.Sample.Value() * 1000
			case node_cpu_ratio.Name:
				data.CpuUsedRatio = nsInfo.Sample.Value() * 100
			case node_cpu_request_used.Name:
				data.CpuResourceRequest = int64(nsInfo.Sample.Value() * 1000)
			case node_cpu_limit_used.Name:
				data.CpuResourceLimit = int64(nsInfo.Sample.Value() * 1000)
			case node_memory_used.Name:
				data.MemoryUsed = nsInfo.Sample.Value()
			case node_memory_total.Name:
				data.MemoryCapacity = nsInfo.Sample.Value()
			case node_memory_ratio.Name:
				data.MemoryUsedRatio = nsInfo.Sample.Value()
			case node_memory_request_used.Name:
				data.MemoryResourceRequest = int64(nsInfo.Sample.Value())
			case node_memory_limit_used.Name:
				data.MemoryResourceLimit = int64(nsInfo.Sample.Value())
			case node_disk_ratio.Name:
				data.DiskUsedRatio = nsInfo.Sample.Value()
			case node_disk_used.Name:
				data.DiskUsed = nsInfo.Sample.Value()
			case node_disk_total.Name:
				data.DiskCapacity = nsInfo.Sample.Value()
			case node_pod_used.Name:
				data.PodUsed = int(nsInfo.Sample.Value())
			case node_pod_total.Name:
				data.PodCapacity = int(nsInfo.Sample.Value())
			case node_pod_ratio.Name:
				data.PodUsedRatio = nsInfo.Sample.Value()
			}
		}
	}

	return resp, nil
}

// 返回值map：namespace，podname，overview
func (k *kubernetesService) MonitorClusterPodStateList(ctx context.Context, clusterInfo *iapiserver.Cluster) (map[string]map[string]*iapiserver.PodOverview, error) {
	resp := map[string]map[string]*iapiserver.PodOverview{}

	metrics := []KV{cluster_pod_cpu_used, cluster_pod_mem_used, cluster_pod_net_rb, cluster_pod_net_wb}
	ret, err := k.autoMetricRealtimeList(ctx, clusterInfo, metrics...)
	if err != nil {
		return nil, err
	}

	for _, v := range metrics {
		metricInfo := ret[v.Name]
		for _, info := range metricInfo.MetricValues {
			namespace := info.Metadata["namespace"]
			podname := info.Metadata["pod"]
			nodename := info.Metadata["node"]
			if namespace == "" || podname == "" {
				continue
			}

			data, ok := resp[namespace]
			if !ok {
				data = map[string]*iapiserver.PodOverview{}
				resp[namespace] = data
			}
			podInfo, ok := data[podname]
			if !ok {
				podInfo = &iapiserver.PodOverview{Namespace: namespace, Name: podname, NodeName: nodename}
				data[podname] = podInfo
			}
			if info.Sample == nil {
				continue
			}
			switch v.Name {
			case cluster_pod_cpu_used.Name:
				podInfo.Cpu.Used = info.Sample.Value() * 1000
			case cluster_pod_mem_used.Name:
				podInfo.Memory.Used = info.Sample.Value()
			case cluster_pod_net_wb.Name:
				podInfo.NetWrite.Used = info.Sample.Value()
			case cluster_pod_net_rb.Name:
				podInfo.NetRead.Used = info.Sample.Value()
			}
		}
	}

	return resp, nil
}

func (k *kubernetesService) MonitorNamespaceStateList(ctx context.Context, clusterInfo *iapiserver.Cluster) (map[string]*iapiserver.MonitorData, error) {
	resp := map[string]*iapiserver.MonitorData{}

	metrics := []KV{namespace_cpu_used, namespace_mem_used, namespace_pod_count, namespace_net_rb, namespace_net_wb}
	ret, err := k.autoMetricRealtimeList(ctx, clusterInfo, metrics...)
	if err != nil {
		return nil, err
	}

	for _, v := range metrics {
		metricInfo := ret[v.Name]
		for _, nsInfo := range metricInfo.MetricValues {
			nsName := nsInfo.Metadata["namespace"]
			data, ok := resp[nsName]
			if !ok {
				data = &iapiserver.MonitorData{Name: nsName}
				resp[nsName] = data
			}

			if nsInfo.Sample == nil {
				continue
			}
			switch v.Name {
			case namespace_cpu_used.Name:
				data.Cpu = nsInfo.Sample.Value()
			case namespace_mem_used.Name:
				data.Memory = nsInfo.Sample.Value()
			case namespace_pod_count.Name:
				data.PodCount = nsInfo.Sample.Value()
			case namespace_net_wb.Name:
				data.NetworkWrite = nsInfo.Sample.Value()
			case namespace_net_rb.Name:
				data.NetworkRead = nsInfo.Sample.Value()
			}
		}
	}

	return resp, nil
}

// cpu 		m
// memory 		byte
// disk 		byte
// net  		byte
func (k *kubernetesService) MonitorPodState(ctx context.Context, clusterInfo *iapiserver.Cluster, objs []metav1.ObjectMeta) (*iapiserver.MonitorData, error) {
	resp := &iapiserver.MonitorData{}

	if len(objs) != 0 {
		rets, err := k.MonitorClusterPodStateList(ctx, clusterInfo)
		if err != nil {
			return nil, err
		}

		for _, v := range objs {
			if namespaceOverview, ok := rets[v.Namespace]; ok {
				if podOverview, ok := namespaceOverview[v.Name]; ok {
					resp.Cpu += podOverview.Cpu.Used
					resp.Memory += podOverview.Memory.Used
					resp.NetworkRead += podOverview.NetRead.Used
					resp.NetworkWrite += podOverview.NetWrite.Used
				}
			}
		}
	}

	return resp, nil
}

func (k *kubernetesService) namespaceOverview(ctx context.Context, clusterInfo *iapiserver.Cluster, namespace string) (resp []*iapiserver.MonitorData, err error) {
	selectPods := map[string]bool{}

	podList, err := k.PodList(ctx, &iapiserver.PodListRequest{
		ResourceListRequest: iapiserver.ResourceListRequest{
			Cluster:   clusterInfo.ID,
			Namespace: namespace,
		},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	deps, podnames, err := k.MonitorDeploymentState(ctx, clusterInfo, namespace)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp = append(resp, deps...)
	for _, v := range podnames {
		selectPods[v] = true
	}

	stats, podnames, err := k.MonitorStatefulSetState(ctx, clusterInfo, namespace)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp = append(resp, stats...)
	for _, v := range podnames {
		selectPods[v] = true
	}
	daes, podnames, err := k.MonitorDaemonsetState(ctx, clusterInfo, namespace)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp = append(resp, daes...)
	for _, v := range podnames {
		selectPods[v] = true
	}

	wg := waitgroup.RunGenericConcurrentlyCondition(ctx, podList.List, func(p *iapiserver.PodInfo) bool {
		if p.Resource != nil {
			if _, ok := selectPods[p.Resource.Name]; ok {
				return true
			}
		}
		return false
	}, func(ctx context.Context, p *iapiserver.PodInfo) waitgroup.GenericResult[*iapiserver.MonitorData] {
		ret, err := k.MonitorPodState(ctx, clusterInfo, []metav1.ObjectMeta{p.Resource.ObjectMeta})
		if err != nil {
			return waitgroup.NewGenericResult(&iapiserver.MonitorData{}, err)
		}

		ret.Namespace = p.Resource.ObjectMeta.Namespace
		ret.Name = p.Resource.ObjectMeta.Name
		ret.Type = pod
		return waitgroup.NewGenericResult(ret, nil)
	})

	for _, v := range wg.GetResults() {
		if v.Error == nil {
			resp = append(resp, v.Data)
		}
	}
	return resp, nil
}

var (
	grafanalocker       sync.Mutex
	grafanaDashIdMap    = map[string]map[string]string{}
	grafanaTimeout      = time.Duration(10)
	grafanaOrgId        = int64(1)
	grafanaDataSource   = "var-datasource=prometheus"
	grafanaDefaultPara  = "?orgId=1&refresh=10s"
	grafanaScheme       = "http://"
	grafanaDefultDashId = "kubernetes-compute-resources-cluster" // 默认参数
)

func (k *kubernetesService) getGrafanaUrlFromGrafanaService(ctx context.Context, clusterInfo *iapiserver.Cluster) (string, error) {
	service, err := libkubernetes.ServiceGet(ctx, clusterInfo.Config, ikubernetes.MonitorNamespace, grafanaServiceName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	// if service.Spec.Type != corev1.ServiceType(nodePort) {
	// 	service.Spec.Type = corev1.ServiceType(nodePort)
	// 	service, err = libkubernetes.ServiceUpdate(ctx,clusterInfo.Config, ikuberneres.MonitorNamespace, service, metav1.UpdateOptions{})
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	time.Sleep(1 * time.Second)
	// }

	if service.Spec.Type != corev1.ServiceType(nodePort) {
		return "", errors.Errorf("grafana service is not nodePort type")
	}
	if len(service.Spec.Ports) == 0 {
		return "", errors.Errorf("namespace:monitoring,service:grafana,type:NodePort portis empty")
	}

	url := grafanaScheme + ParseAddrFromURLNoError(clusterInfo.Config.Host) + ":" + strconv.Itoa(int(service.Spec.Ports[0].NodePort))
	return url, nil
}

func (k *kubernetesService) GrafanaDashUrlGet(ctx context.Context, req *iapiserver.GrafanaGetRequest) (*iapiserver.GrafanaResponse, error) {
	resp := &iapiserver.GrafanaResponse{}

	if req.SubUrl == "" {
		req.SubUrl = grafanaDefultDashId
	}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	resp.TotalUrl, err = k.createDashUrl(ctx, cluster, req.SubUrl)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.TotalUrl += "?" + createDashPara(req.Namespace, req.Workload, req.Node, req.Pod)
	return resp, nil
}

func createDashPara(namespace, workload, node, pod string) string {
	resp := grafanaDataSource

	add := func(src string, flag string, tar string) string {
		if tar != "" {
			src = src + "&" + flag + "=" + tar
		}
		return src
	}

	resp = add(resp, "var-namespace", namespace)
	resp = add(resp, "var-workload", workload)
	resp = add(resp, "var-node", node)
	resp = add(resp, "var-pod", pod)
	return resp
}

func (k *kubernetesService) createDashUrl(ctx context.Context, clusterInfo *iapiserver.Cluster, subUrl string) (string, error) {
	grafanaUrl, err := k.getGrafanaUrlFromGrafanaService(ctx, clusterInfo)
	if err != nil {
		return "", err
	}

	// find urls
	var totalUrl string
	urlMap, err := k.CacheGet(ctx, clusterInfo)
	if err != nil {
		return "", err
	}
	for k, v := range urlMap {
		if strings.HasSuffix(k, subUrl) {
			totalUrl = v
			break
		}
	}
	if totalUrl == "" {
		return "", errors.Errorf("cluster_uuid[%v] grafana dash sub_url[%v] not exist", clusterInfo.ID, subUrl)
	}

	totalUrl = grafanaUrl + totalUrl + grafanaDefaultPara
	return totalUrl, nil
}

func (k *kubernetesService) CacheGet(ctx context.Context, clusterInfo *iapiserver.Cluster) (map[string]string, error) {
	grafanalocker.Lock()
	cacheUrls := grafanaDashIdMap[clusterInfo.ID]
	urls := make(map[string]string, len(cacheUrls))
	for k, v := range cacheUrls {
		urls[k] = v
	}
	grafanalocker.Unlock()

	if len(urls) != 0 {
		go k.CacheUpdate(ctx, clusterInfo)
		return urls, nil
	}

	if err := k.CacheUpdate(ctx, clusterInfo); err != nil {
		return nil, err
	}
	grafanalocker.Lock()
	defer grafanalocker.Unlock()
	cacheUrls = grafanaDashIdMap[clusterInfo.ID] //CacheUpdate() not error，cluster_uuid must exist
	for k, v := range cacheUrls {
		urls[k] = v
	}

	return urls, nil
}

func (k *kubernetesService) CacheUpdate(ctx context.Context, clusterInfo *iapiserver.Cluster) error {
	urls, err := k.loadGrafanaUrls(ctx, clusterInfo)
	if err != nil {
		return err
	}

	grafanaDashIdMap[clusterInfo.ID] = urls
	return nil
}

func (k *kubernetesService) loadGrafanaUrls(ctx context.Context, clusterInfo *iapiserver.Cluster) (map[string]string, error) {
	// grafanaUrl, err := k.getGrafanaUrlFromGrafanaService(ctx,clusterInfo)
	// if err != nil {
	// 	return nil, err
	// }

	// config := gapi.Config{}
	// config.OrgID = grafanaOrgId
	// config.Client = &http.Client{Timeout: grafanaTimeout * time.Second}
	// client, err := gapi.New(grafanaUrl, config)
	// if err != nil {
	// 	return nil, errors.WithStack(err)
	// }

	// dashboards, err := client.Dashboards()
	// if err != nil {
	// 	return nil, errors.WithStack(err)
	// }

	// // update cache
	// dashIdMap := map[string]string{}
	// for _, v := range dashboards {
	// 	ss := strings.Split(v.URL, "/")
	// 	dashIdMap[ss[len(ss)-1]] = v.URL
	// }
	return nil, nil

	// return dashIdMap, nil
}

func (k *kubernetesService) PrometheusResourceCreate(ctx context.Context, req *iapiserver.PrometheusRequest) (*iapiserver.PrometheusResponse, error) {
	resp := &iapiserver.PrometheusResponse{}

	if req.Resource == nil {
		req.Resource = &monitoringv1.Prometheus{}
		req.Resource.Namespace = "monitoring"
	}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	meta, err := libkubernetes.PrometheusCreate(ctx, cluster.Config, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.Info = convertK8sPrometheusToApi(meta, cluster, req.Yaml)
	return nil, errors.WithStack(err)
}

func (k *kubernetesService) PrometheusResourceDelete(ctx context.Context, req *iapiserver.PrometheusRequest) error {
	req.Resource.ResourceVersion = ""
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = libkubernetes.PrometheusDelete(ctx, cluster.Config, req.Resource.Namespace, req.Resource.Name, req.DeleteOpts); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (k *kubernetesService) PrometheusResourceUpdate(ctx context.Context, req *iapiserver.PrometheusRequest) (*iapiserver.PrometheusInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	old, err := libkubernetes.PrometheusGet(ctx, cluster.Config, req.Resource.Namespace, req.Resource.Name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	new := req.Resource
	old.Spec.Image = new.Spec.Image
	old.Spec.Replicas = new.Spec.Replicas
	old.Spec.Resources = new.Spec.Resources
	old.Spec.Retention = new.Spec.Retention
	old.Spec.RetentionSize = new.Spec.RetentionSize
	old.Spec.NodeSelector = new.Spec.NodeSelector

	meta, err := libkubernetes.PrometheusUpdate(ctx, cluster.Config, req.Resource.Namespace, old, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sPrometheusToApi(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) PrometheusResourceGet(ctx context.Context, req *iapiserver.PrometheusGetRequest) (*iapiserver.PrometheusInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.PrometheusGet(ctx, cluster.Config, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	info := convertK8sPrometheusToApi(meta, cluster, req.Yaml)
	if info.GrafanaUrl, err = k.getGrafanaUrlFromGrafanaService(ctx, cluster); err != nil {
		log.Errorf("getgrafanaurl err : clusteruuid[%v],clustername[%v],err[%v]", cluster.ID, cluster.Name, err.Error())
	}
	if info.Url, err = k.getPrometheusUrlFromPrometheusService(ctx, cluster); err != nil {
		log.Errorf("get prometheusurl err : clusteruuid[%v],clustername[%v],err[%v]", cluster.ID, cluster.Name, err.Error())
	}
	return info, nil
}

func (k *kubernetesService) PrometheusResourceList(ctx context.Context, req *iapiserver.PrometheusListRequest) (*iapiserver.PrometheusListResponse, error) {
	resp := &iapiserver.PrometheusListResponse{}
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	ret, err := libkubernetes.PrometheusList(ctx, cluster.Config, req.Namespace, req.ToListOpts(), 20)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var resInfos []*iapiserver.PrometheusInfo
	for i, dm := range ret.Items {
		if !sets.NewString(dm.Name).ContainAny(strings.Split(req.Fuzzy, " ")...) {
			continue
		}

		resInfos = append(resInfos, convertK8sPrometheusToApi(&ret.Items[i], cluster, req.Yaml))
	}

	resp.List = resInfos
	s, e := paging.Index(len(resp.List), req.PageNum, req.PageSize)
	resp.List = resp.List[s:e]
	return resp, nil
}

func convertK8sPrometheusToApi(meta *monitoringv1.Prometheus, cluster *iapiserver.Cluster, yaml bool) *iapiserver.PrometheusInfo {
	pi := &iapiserver.PrometheusInfo{
		Prometheus: meta,
	}

	pi.ResourceConvert = &iapiserver.PrometheusResourceConvert{
		Limits:  make(map[string]int64),
		Request: make(map[string]int64),
	}
	if pi.Prometheus.Spec.Resources.Requests != nil {
		for i, j := range pi.Prometheus.Spec.Resources.Requests {
			if string(i) == "memory" {
				pi.ResourceConvert.Request[string(i)] = j.Value() / 1024 / 1024 // convert to memory
				continue
			}
			if string(i) == "cpu" { // if use Value(), 0.1 cpu/100m cpu will convert to 1 cpu
				pi.ResourceConvert.Request[string(i)] = j.MilliValue()
				continue
			}
			pi.ResourceConvert.Request[string(i)] = j.Value()
		}
	}

	if pi.Prometheus.Spec.Resources.Limits != nil {
		for i, j := range pi.Prometheus.Spec.Resources.Limits {
			if string(i) == "memory" {
				pi.ResourceConvert.Limits[string(i)] = j.Value() / 1024 / 1024 // convert to memory
				continue
			}
			if string(i) == "cpu" { // if use Value(), 0.1 cpu/100m cpu will convert to 1 cpu
				pi.ResourceConvert.Limits[string(i)] = j.MilliValue()
				continue
			}
			pi.ResourceConvert.Limits[string(i)] = j.Value()
		}
	}

	return pi
}
