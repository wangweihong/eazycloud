package kubernetes

// import (
// 	"errors"
// 	"fmt"

// 	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
// 	gapi "github.com/grafana/grafana-apis-golang-client"
// 	"github.com/wangweihong/eazycloud/apis/iapiserver"
// 	corev1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// )

// var (
// 	prometheusServiceName = "prometheus-k8s"
// 	grafanaServiceName    = "grafana"
// 	statefulset           = "statefulset"
// 	daemonset             = "daemonset"
// 	deployment            = "deployment"
// 	pod                   = "pod"
// 	nodePort              = "NodePort"
// 	componentEtcd         = "etcd"
// )

// func (k *kubernetesService) getPrometheusUrlFromPrometheusService(clusterInfo *iapiserver.Cluster) (string, error) {
// 	service, err := clientset.ServiceGet(ctx,clusterInfo, monitorNamespace, prometheusServiceName, metav1.GetOptions{})
// 	if err != nil {
// 		return "", err
// 	}

// 	if service.Spec.Type != corev1.ServiceType(nodePort) {
// 		return "", errors.Errorf("namespace:monitoring,service:prometheus-k8s,type:NodePort")
// 	}

// 	if len(service.Spec.Ports) == 0 {
// 		return "", errors.Errorf("namespace:monitoring,service:prometheus-k8s,type:NodePort")
// 	}

// 	url := "http://" + utils.ParseAddrFromURLNoError(clusterInfo.Config.Host) + ":" +strconv.Itoa(int(service.Spec.Ports[0].NodePort))
// 	return url, nil
// }

// func (k *kubernetesService) getPrometheusConfigProqlOption(clusterInfo *iapiserver.Cluster, req *topke.Monitoring, namespace string) (*topke.PrometheusConfig, *topke.ProqlOption, error) {
// 	conf := &topke.PrometheusConfig{}
// 	opt := &topke.ProqlOption{}
// 	var err error

// 	if !clientset.IsMonitorServiceReady(clusterInfo.UUID) {
// 		return nil, nil, errors.Errorf("cluster '%v' not ready",clusterInfo.UUID)
// 	}

// 	conf.Address, err = tm.getPrometheusUrlFromPrometheusService(clusterInfo)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	opt.Namespace = namespace
// 	if req != nil {
// 		opt.Level = req.Level
// 		opt.Node = req.NodeName
// 		opt.PodName = req.Name
// 		opt.Time = time.Unix(req.Time, 0)
// 		opt.StartTime = time.Unix(req.StartTime, 0)
// 		opt.EndTime = time.Unix(req.EndTime, 0)
// 		opt.Duration = time.Duration(req.Duration) * time.Second
// 	}

// 	return conf, opt, nil
// }

// func (k *kubernetesService) checkMonitoringListRequest(req *topke.MonitoringListRequest) error {
// 	if req == nil {
// 		return errors.Errorf("metrics is empty")
// 	}
// 	if req.StartTime == 0 && req.EndTime == 0 {
// 		d, _ := time.ParseDuration("-30m")
// 		req.EndTime = time.Now().Unix()
// 		req.StartTime = time.Now().Add(d).Unix()
// 	}

// 	if req.Duration == 0 {
// 		req.Duration = 1
// 	}

// 	switch req.MonitoringParam.Level {
// 	case libprometheus.NamespaceLevel,
// 		libprometheus.NodeLevel,
// 		libprometheus.ClusterLevel,
// 		libprometheus.PodLevel,
// 		libprometheus.AlertLevel:
// 	default:
// 		log.Errorf("monitoring.level[%v] not support", req.MonitoringParam.Level)
// 	}

// 	return nil
// }

// func (k *kubernetesService) MonitorServiceStatusGet(ctx context.Context, req *topke.MonitoringGetRequest)*topke.MonitoringResponse {
// 	resp := &topke.MonitoringResponse{Status: status.SuccessStatus, MonitorStatus:&topke.MonitorStatus{PrometheusEnable: true, GrafanaEnable: true}}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	resp.MonitorStatus.PrometheusUrl, err = tm.getPrometheusUrlFromPrometheusService(clusterInfo)
// 	if err != nil {
// 		resp.MonitorStatus.PrometheusEnable = false
// 		resp.MonitorStatus.PrometheusError = err.Error()
// 	}

// 	resp.MonitorStatus.GrafanaUrl, err = tm.getGrafanaUrlFromGrafanaService(clusterInfo)
// 	if err != nil {
// 		resp.MonitorStatus.GrafanaEnable = false
// 		resp.MonitorStatus.GrafanaError = err.Error()
// 	}

// 	return resp
// }

// func (k *kubernetesService) MonitorMetricsHistoryList(ctx context.Context, req *topke.MonitoringListRequest)*topke.MonitoringMetricsHistoryListResponse {
// 	resp := &topke.MonitoringMetricsHistoryListResponse{Status: status.SuccessStatus}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := tm.checkMonitoringListRequest(req); err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	proConf, opt, err := tm.getPrometheusConfigProqlOption(clusterInfo, req.ToMonitoring(), req.Namespace)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	timeRange := req.ToRange()
// 	resp.TimeRange = &timeRange
// 	resp.Metrics, err = libprometheus.GetNamedMetricsHistory(ctx, proConf, req.ToMonitoring().Metrics, opt)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	sort.SliceStable(resp.Metrics, func(i, j int) bool {
// 		return resp.Metrics[i].MetricName < resp.Metrics[j].MetricName
// 	})

// 	return resp
// }

// func (k *kubernetesService) MonitorClusterHardwareResourceHistoryList(ctx context.Context, req *topke.MonistoringHardwareHistoryListRequest) *topke.MonistoringHardwareHistoryListResponse {
// 	resp := &topke.MonistoringHardwareHistoryListResponse{Status: status.SuccessStatus}

// 	metrics := []KV{cluster_cpu_ratio, cluster_cpu_load_1, cluster_cpu_load_5, cluster_cpu_load_15,
// cluster_memory_ratio,
// 		cluster_disk_used_total, cluster_disk_io_rate, cluster_network_receive_rate, cluster_network_send_rate,
// 		cluster_pod_unknown, cluster_pod_failed, cluster_pod_pending, cluster_pod_succeeded, cluster_pod_running}
// 	rets, err := tm.autoMetricHistoryList(ctx, req.ClusterUUID, req.MonitoringTimeRange.ToRange(), metrics...)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	timeRange := req.ToRange()
// 	resp.TimeRange = &timeRange
// 	resp.Metrics = make([]topke.Metric, 0, len(rets))
// 	for _, v := range rets {
// 		resp.Metrics = append(resp.Metrics, *v)
// 	}

// 	sort.SliceStable(resp.Metrics, func(i, j int) bool {
// 		return resp.Metrics[i].MetricName < resp.Metrics[j].MetricName
// 	})

// 	return resp
// }

// func (k *kubernetesService) MonitorClusterNodeResourceHistoryGet(ctx context.Context, req
// *topke.MonitoringNodeGetRequest) *topke.MonitoringNodeGetResponse {
// 	resp := &topke.MonitoringNodeGetResponse{Status: status.SuccessStatus}

// 	if req.NodeName == "" {
// 		resp.Status = status.NewStatusDesc(scode.TopECParameterEmpty, "node_name is empty")
// 		return resp
// 	}

// 	timeRange := req.ToRange()
// 	history, err := tm.autoMetricHistoryList(ctx, req.ClusterUUID, timeRange, nodePhysicalResourceMetrics...)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.TimeRange = &timeRange
// 	for _, v := range history {
// 		resp.Metrics = append(resp.Metrics, *v)
// 	}
// 	return resp
// }

// // func (k *kubernetesService) MonitorNamespaceList(ctx context.Context, req *topke.MonitoringNamespaceListRequest)*topke.MonitoringNamespaceListResponse {
// // 	resp := &topke.MonitoringNamespaceListResponse{Status: status.SuccessStatus}

// // 	clusterList, err := GetTopkeManager().GetVisitScope(req.ResourceListRequest)
// // 	if err != nil {
// // 		resp.Status = err
// // 		return resp
// // 	}

// // 	wg := utils.NewWaitGroup(ctx)
// // 	for _, cluster := range clusterList {
// // 		cluster := cluster
// // 		extra := getObjectExtraInfos(cluster)
// // 		wg.Start(utils.NewWaitGroupHandleFunc("", ctx, func() utils.WaitGroupResult {
// // 			clusterListOne := topke.NewEachResourceRangeListState(cluster.UUID, cluster.Name)
// // 			ret, err := tm.MonitorNamespaceStateList(ctx, cluster)
// // 			if err != nil {
// // 				return utils.NewWaitGroupResult(clusterListOne, err)
// // 			}
// // 			nsList, err := clientset.NamespaceList(cluster, metav1.ListOptions{})
// // 			if err != nil {
// // 				return utils.NewWaitGroupResult(clusterListOne, err)
// // 			}

// // 			rets := make([]*topke.MonitorData, 0, len(ret))
// // 			for _, ns := range nsList.Items {
// // 				if req.Fuzzy != "" && NewFieldFilter(ns.Name, ns.Namespace, cluster.Name).Filter(req.Fuzzy) {
// // 					continue
// // 				}

// // 				nsInfo, ok := ret[ns.Name]
// // 				if !ok {
// // 					nsInfo = &topke.MonitorData{Name: ns.Name, Extra: extra}
// // 					ret[ns.Name] = nsInfo
// // 				}
// // 				nsInfo.Cpu = nsInfo.Cpu / 1000
// // 				rets = append(rets, nsInfo)
// // 			}
// // 			clusterListOne.TotalCount = len(rets)
// // 			clusterListOne.List = rets
// // 			return utils.NewWaitGroupResult(clusterListOne, nil)
// // 		}))
// // 	}
// // 	wg.Wait()

// // 	if req.SortBy == "" {
// // 		req.SortBy = "Cpu"
// // 	}
// // 	CutPagingSliceFromWgResultsV2(wg, &resp.EachRangeListState, &resp.MonitorDatas, req.PageNumber, req.PageSize,
// // &resp.TotalCount, req.SortBy, req.SortDesc)
// // 	return resp
// // }

// // get all metrics name in promtheus
// func (k *kubernetesService) MonitorMetricNameList(ctx context.Context, req *topke.MonitoringMetricsNameListRequest)*topke.MonitoringMetricsNameListResponse {
// 	resp := &topke.MonitoringMetricsNameListResponse{Status: status.SuccessStatus}
// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if !clientset.IsMonitorServiceReady(clusterInfo.UUID) {
// 		resp.Status = status.NewStatusDesc(scode.TopECTopkeKubeMonitorNotHealth, fmt.Sprintf("cluster_uuid[%v]",
// clusterInfo.UUID))
// 		return resp
// 	}

// 	prometheusAddress, err := tm.getPrometheusUrlFromPrometheusService(clusterInfo)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	//query all metrics name in promtheus
// 	ret, err := libprometheus.LabelValues(ctx,
// 		&topke.PrometheusConfig{Address: prometheusAddress},
// 		"__name__",
// 		15)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.TotalCount = len(ret)
// 	resp.MetricList = make([]string, 0, ret.Len())
// 	for _, v := range ret {
// 		resp.MetricList = append(resp.MetricList, string(v))
// 	}

// 	return resp
// }

// func (k *kubernetesService) MonitorRuleAutoList(ctx context.Context, req *topke.MonitoringListRequest)
// *topke.MonitoringListResponse {
// 	resp := &topke.MonitoringListResponse{Status: status.SuccessStatus}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := tm.checkMonitoringListRequest(req); err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	proConf, opt, err := tm.getPrometheusConfigProqlOption(clusterInfo, req.ToMonitoring(), req.Namespace)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.Metrics, err = libprometheus.GetAutoMetricsHistory(ctx, proConf, req.ToMonitoring().Metrics, opt)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}
// 	resp.Metrics = tm.metricCombine(resp.Metrics)
// 	if len(resp.Metrics) != 0 {
// 		resp.Metric = &resp.Metrics[0]
// 	}
// 	resp.Metrics = nil
// 	return resp
// }

// func (k *kubernetesService) metricCombine(ms []topke.Metric) []topke.Metric {
// 	resp := make([]topke.Metric, len(ms), len(ms))

// 	for k, metricv := range ms {
// 		tmp := topke.Metric{MetricName: metricv.MetricName}
// 		tmp.MetricValues = make([]topke.MetricValue, 2, 2)
// 		tmp.MetricType = metricv.MetricType
// 		tmp.Error = metricv.Error

// 		cb := NewCombinerBig()
// 		cs := NewCombinerSmall()
// 		for _, mcv := range metricv.MetricValues {
// 			for _, pv := range mcv.Series {
// 				cb.Push(int64(pv.Timestamp()), pv.Value())
// 				cs.Push(int64(pv.Timestamp()), pv.Value())
// 			}
// 		}

// 		tmp.MetricValues[0].Series = cs.GetInsert()
// 		tmp.MetricValues[0].Name = cs.GetTypeString()
// 		tmp.MetricValues[0].Max = cs.GetMax()
// 		tmp.MetricValues[1].Series = cb.GetInsert()
// 		tmp.MetricValues[1].Name = cb.GetTypeString()
// 		tmp.MetricValues[1].Max = cb.GetMax()
// 		resp[k] = tmp
// 	}

// 	return resp
// }

// // func (k *kubernetesService) MonitorAlertList(ctx context.Context, req *topke.MonitoringListRequest)
// // *topke.MonitoringListResponse {
// // 	resp := &topke.MonitoringListResponse{Status: status.SuccessStatus}

// // 	clusterList, err := GetTopkeManager().GetVisitScope(req.ResourceListRequest)
// // 	if err != nil {
// // 		resp.Status = status.UpdateStatus(resp.Status)
// // 		return resp
// // 	}
// // 	wg := utils.NewWaitGroup(nil)
// // 	for _, cluster := range clusterList {
// // 		cluster := cluster
// // 		wg.Start(utils.NewWaitGroupHandleFunc("", nil, func() utils.WaitGroupResult {
// // 			clusterListOne := topke.NewEachResourceRangeListState(cluster.UUID, cluster.Name)
// // 			proConf, _, err := tm.getPrometheusConfigProqlOption(cluster, req.ToMonitoring(), req.Namespace)
// // 			if err != nil {
// // 				return utils.NewWaitGroupResult(clusterListOne, err)
// // 			}
// // 			ret, err := libprometheus.AlertList(ctx, proConf, 15)
// // 			if err != nil {
// // 				return utils.NewWaitGroupResult(clusterListOne, err)
// // 			}
// // 			extra := getObjectExtraInfos(cluster)

// // 			alertList := resp.AlertList
// // 			for _, v := range ret.Alerts {
// // 				v := v

// // 				alert := topke.Alert{}
// // 				alert.TenantUuid = extra[topke.TopkeAnnnotationTenantUUIDKey]
// // 				alert.Tenant = extra[topke.TopkeAnnnotationTenantNameKey]
// // 				alert.ClusterUuid = extra[topke.TopkeAnnnotationClusterUUIDKey]
// // 				alert.ClusterName = extra[topke.TopkeAnnnotationClusterNameKey]
// // 				alert.Namespace = string(v.Labels["namespace"])
// // 				alert.Name = string(v.Labels["alertname"])
// // 				alert.Message = string(v.Annotations["message"])
// // 				alert.ActiveTime = v.ActiveAt.Unix()
// // 				alert.Pod = string(v.Labels["pod"])
// // 				alert.Container = string(v.Labels["container"])
// // 				alert.Severity = alertLevelValueDesc[string(v.Labels["severity"])]

// // 				v2, _ := strconv.ParseFloat(v.Value, 64)
// // 				alert.Value = fmt.Sprintf("%.6f", v2)
// // 				alert.Alertstate = string(v.State)
// // 				alert.Metadata = &v
// // 				if req.Namespace != "" && req.Namespace != alert.Namespace {
// // 					continue
// // 				}

// // 				if req.Fuzzy != "" && !utils.IsContain([]string{alert.Namespace, alert.Name, alert.Pod, alert.Container,
// // alert.Alertstate, 					alert.Severity, alert.TenantUuid, alert.Tenant, alert.ClusterUuid, alert.ClusterName,
// // alert.Value, alert.Message}, req.Fuzzy) {
// // 					continue
// // 				}
// // 				alertList = append(alertList, &alert)
// // 			}

// // 			clusterListOne.TotalCount = len(alertList)
// // 			clusterListOne.List = alertList
// // 			return utils.NewWaitGroupResult(clusterListOne, nil)
// // 		}))
// // 	}
// // 	wg.Wait()

// // 	if req.SortBy == "" {
// // 		req.SortBy = "ClusterName Namespace ActiveTime"
// // 	}
// // 	CutPagingSliceFromWgResultsV2(wg, &resp.EachRangeListState, &resp.AlertList, req.PageNumber, req.PageSize,
// // &resp.TotalCount, req.SortBy, req.SortDesc)

// // 	return resp
// // }

// func (k *kubernetesService) MonitorAlertDelete(ctx context.Context, req *topke.MonitoringRequest)*topke.MonitoringResponse {
// 	resp := &topke.MonitoringResponse{Status: status.SuccessStatus}

// 	if len(req.AlertDeleteList) == 0 {
// 		resp.Status = errors.Errorf("alert_delete_list is empty")
// 		return resp
// 	}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	proConf, _, err := tm.getPrometheusConfigProqlOption(clusterInfo, req.Monitoring, req.Namespace)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	ret, err := libprometheus.AlertList(ctx, proConf, 15)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}
// 	deleteAlert := req.AlertDeleteList[0]

// 	for _, v := range ret.Alerts {
// 		if v.Labels.Equal(deleteAlert.Labels) {
// 			if err := libprometheus.DeleteSeries(ctx, proConf, []string{"ALERTS" + deleteAlert.Labels.String()}, v.ActiveAt);
// err != nil {
// 				resp.Status = err
// 				return resp
// 			}
// 		}
// 	}
// 	return resp
// }

// // func (k *kubernetesService) MonitorPodList(ctx context.Context, req *topke.MonitoringPodListRequest)
// // *topke.MonitoringPodListResponse {
// // 	resp := &topke.MonitoringPodListResponse{Status: status.SuccessStatus}

// // 	clusterList, err := GetTopkeManager().GetVisitScope(req.ResourceListRequest)
// // 	if err != nil {
// // 		resp.Status = err
// // 		return resp
// // 	}

// // 	wg := utils.NewWaitGroup(nil)
// // 	for _, cluster := range clusterList {
// // 		cluster := cluster
// // 		wg.Start(utils.NewWaitGroupHandleFunc("", nil, func() utils.WaitGroupResult {
// // 			clusterListOne := topke.NewEachResourceRangeListState(cluster.UUID, cluster.Name)
// // 			nodeList, err := clientset.NodeList(cluster, metav1.ListOptions{}) // local cache
// // 			if err != nil {
// // 				return utils.NewWaitGroupResult(clusterListOne, err)
// // 			}
// // 			podList, err := clientset.PodList(cluster, req.Namespace, metav1.ListOptions{}) // local cache
// // 			if err != nil {
// // 				return utils.NewWaitGroupResult(clusterListOne, err)
// // 			}
// // 			nsPods, err := tm.MonitorClusterPodStateList(ctx, cluster) // k8s
// // 			if err != nil {
// // 				return utils.NewWaitGroupResult(clusterListOne, err)
// // 			}
// // 			nodes := map[string]*corev1.Node{} // node info
// // 			for i := range nodeList.Items {
// // 				nodes[nodeList.Items[i].Name] = &nodeList.Items[i]
// // 			}
// // 			rets := make([]*topke.PodOverview, 0, len(podList.Items))
// // 			for _, v := range podList.Items { // all pod
// // 				info := nsPods[v.Namespace][v.Name]
// // 				if info == nil {
// // 					info = &topke.PodOverview{Namespace: v.Namespace, Name: v.Name, Extra: getObjectExtraInfos(cluster)}
// // 					info.NodeName = v.Status.HostIP
// // 				}
// // 				info.NodeAddr = getNodeAddr(nodes[info.NodeName])
// // 				info.CreateTime = v.CreationTimestamp.Unix()

// // 				if !utils.IsContain([]string{info.Name, info.Namespace, info.ClusterName, info.NodeName, info.NodeAddr},
// // req.Fuzzy) {
// // 					continue
// // 				}

// // 				rets = append(rets, info)
// // 			}

// // 			clusterListOne.TotalCount = len(rets)
// // 			clusterListOne.List = rets
// // 			return utils.NewWaitGroupResult(clusterListOne, nil)
// // 		}))
// // 	}
// // 	wg.Wait()

// // 	if req.SortBy == "" {
// // 		req.SortBy = "Extra/topke.cluster.name/ Namespace Name Cpu/Used"
// // 	}
// // 	CutPagingSliceFromWgResultsV2(wg, &resp.EachRangeListState, &resp.PodList, req.PageNumber, req.PageSize,
// // &resp.TotalCount, req.SortBy, req.SortDesc)
// // 	return resp
// // }

// func (k *kubernetesService) MonitorNamespaceTopGet(ctx context.Context, req *topke.MonitoringTopGetRequest)*topke.MonitoringTopGetResponse {
// 	resp := &topke.MonitoringTopGetResponse{Status: status.SuccessStatus}

// 	if req.ClusterUUID == "" || req.Namespace == "" {
// 		resp.Status = errors.Errorf("topke_cluster_uuid or namespace is empty")
// 		return resp
// 	}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	stateList, err := tm.namespaceOverview(ctx, clusterInfo, req.Namespace)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.TotalCount = len(stateList)
// 	s, e := utils.PagingIndex(len(stateList), req.PageNumber, req.PageSize)
// 	resp.MonitorDatas = stateList[s:e]
// 	return resp
// }

// // func (k *kubernetesService) MonitorStatefulSetState(ctx context.Context, clusterInfo *iapiserver.Cluster, namespace string)([]*topke.MonitorData, []string, error) {
// // 	workLoadState := []*topke.MonitorData{}
// // 	seletedPod := []string{}
// // 	var err error

// // 	podList := &topke.PodListResponse{}
// // 	workloadList := &topke.StatefulSetListResponse{}

// // 	wg := utils.NewWaitGroup(nil)
// // 	wg.Start(utils.WaitGroupRoutineFunc{Call: func() utils.WaitGroupResult {
// // 		podList = tm.PodListAll(ctx, &topke.PodListRequest{
// // 			ResourceListRequest: topke.ResourceListRequest{
// // 				ClusterUUID: clusterInfo.UUID,
// // 				Namespace:   namespace,
// // 			},
// // 		})
// // 		if podList.Status != status.SuccessStatus {
// // 			return utils.NewWaitGroupResult(nil, status.UpdateStatus(podList.Status))
// // 		}
// // 		return utils.NewWaitGroupResult(nil, nil)
// // 	}})
// // 	wg.Start(utils.WaitGroupRoutineFunc{Call: func() utils.WaitGroupResult {
// // 		workloadList = GetTopkeManager().StatefulSetListAll(ctx, &topke.StatefulSetListRequest{
// // 			ResourceListRequest: topke.ResourceListRequest{
// // 				ClusterUUID: clusterInfo.UUID,
// // 				Namespace:   namespace,
// // 			},
// // 		})
// // 		if workloadList.Status != status.SuccessStatus {
// // 			return utils.NewWaitGroupResult(nil, status.UpdateStatus(workloadList.Status))
// // 		}
// // 		return utils.NewWaitGroupResult(nil, nil)
// // 	}})
// // 	wg.Wait()

// // 	for _, v := range wg.GetResults() {
// // 		if v.Error != nil {
// // 			return nil, nil, err
// // 		}
// // 	}

// // 	gLock := sync.RWMutex{}
// // 	for _, one := range workloadList.List {
// // 		if one != nil && one.StatefulSet != nil && one.StatefulSet.Spec.Selector != nil {
// // 			one := one
// // 			wg.Start(utils.WaitGroupRoutineFunc{Call: func() utils.WaitGroupResult {
// // 				var objs []metav1.ObjectMeta
// // 				for _, pod := range podList.List {
// // 					if pod != nil && pod.Pod != nil && pod.Pod.Labels != nil {
// // 						if IsMapContainMap(pod.Pod.Labels, one.StatefulSet.Spec.Selector.MatchLabels) {
// // 							objs = append(objs, pod.Pod.ObjectMeta)
// // 							gLock.Lock()
// // 							seletedPod = append(seletedPod, pod.Pod.Name)
// // 							gLock.Unlock()
// // 						}
// // 					}
// // 				}

// // 				ret, err := tm.MonitorPodState(ctx, clusterInfo, objs)
// // 				if err != nil {
// // 					logrus.Errorf("%s", err)
// // 					return utils.WaitGroupResult{}
// // 				}
// // 				ret.Type = statefulset
// // 				ret.Namespace = one.StatefulSet.Namespace
// // 				ret.Name = one.StatefulSet.Name
// // 				gLock.Lock()
// // 				workLoadState = append(workLoadState, ret)
// // 				gLock.Unlock()

// // 				return utils.WaitGroupResult{}
// // 			}})
// // 		}
// // 	}
// // 	wg.Wait()

// // 	return workLoadState, seletedPod, nil
// // }

// // func (k *kubernetesService) MonitorDaemonsetState(ctx context.Context, clusterInfo *iapiserver.Cluster, namespace string)([]*topke.MonitorData, []string, error) {
// // 	workLoadState := []*topke.MonitorData{}
// // 	seletedPod := []string{}
// // 	var err error

// // 	podListResp := &topke.PodListResponse{}
// // 	workloadList := &topke.DaemonSetListResponse{}

// // 	wg := utils.NewWaitGroup(nil)
// // 	wg.Start(utils.WaitGroupRoutineFunc{Call: func() utils.WaitGroupResult {
// // 		podListResp = tm.PodListAll(ctx, &topke.PodListRequest{
// // 			ResourceListRequest: topke.ResourceListRequest{
// // 				ClusterUUID: clusterInfo.UUID,
// // 				Namespace:   namespace,
// // 			},
// // 		})
// // 		if podListResp.Status != status.SuccessStatus {
// // 			return utils.NewWaitGroupResult(nil, status.UpdateStatus(podListResp.Status))
// // 		}
// // 		return utils.NewWaitGroupResult(nil, nil)
// // 	}})
// // 	wg.Start(utils.WaitGroupRoutineFunc{Call: func() utils.WaitGroupResult {
// // 		workloadList = GetTopkeManager().DaemonSetListAll(ctx, &topke.DaemonSetListRequest{
// // 			ResourceListRequest: topke.ResourceListRequest{
// // 				ClusterUUID: clusterInfo.UUID,
// // 				Namespace:   namespace,
// // 			}})
// // 		if workloadList.Status != status.SuccessStatus {
// // 			return utils.NewWaitGroupResult(nil, status.UpdateStatus(workloadList.Status))
// // 		}
// // 		return utils.NewWaitGroupResult(nil, nil)
// // 	}})
// // 	wg.Wait()

// // 	for _, v := range wg.GetResults() {
// // 		if v.Error != nil {
// // 			return nil, nil, err
// // 		}
// // 	}

// // 	lock := sync.RWMutex{}
// // 	for _, one := range workloadList.List {
// // 		if one != nil && one.DaemonSet != nil && one.DaemonSet.Spec.Selector != nil {
// // 			one := one
// // 			wg.Start(utils.WaitGroupRoutineFunc{Call: func() utils.WaitGroupResult {
// // 				objs := []metav1.ObjectMeta{}
// // 				for _, pod := range podListResp.List {
// // 					if pod != nil && pod.Pod != nil && pod.Pod.Labels != nil {
// // 						if IsMapContainMap(pod.Pod.Labels, one.DaemonSet.Spec.Selector.MatchLabels) {
// // 							objs = append(objs, pod.Pod.ObjectMeta)
// // 							lock.Lock()
// // 							seletedPod = append(seletedPod, pod.Pod.Name)
// // 							lock.Unlock()
// // 						}
// // 					}
// // 				}

// // 				ret, err := tm.MonitorPodState(ctx, clusterInfo, objs)
// // 				if err != nil {
// // 					logrus.Errorf("%v", err)
// // 					return utils.WaitGroupResult{}
// // 				}

// // 				ret.Type = daemonset
// // 				ret.Namespace = one.DaemonSet.Namespace
// // 				ret.Name = one.DaemonSet.Name
// // 				lock.Lock()
// // 				workLoadState = append(workLoadState, ret)
// // 				lock.Unlock()
// // 				return utils.WaitGroupResult{}
// // 			}})
// // 		}
// // 	}
// // 	wg.Wait()

// // 	return workLoadState, seletedPod, nil
// // }

// // func (k *kubernetesService) MonitorNodeList(ctx context.Context, req *topke.MonitoringListRequest)*topke.MonitoringListResponse {
// // 	resp := &topke.MonitoringListResponse{Status: status.SuccessStatus}

// // 	clusterList, err := GetTopkeManager().GetVisitScope(req.ResourceListRequest)
// // 	if err != nil {
// // 		resp.Status = err
// // 		return resp
// // 	}

// // 	wg := utils.NewWaitGroup(nil)
// // 	for _, cluster := range clusterList {
// // 		cluster := cluster
// // 		wg.Start(utils.NewWaitGroupHandleFunc("", nil, func() utils.WaitGroupResult {
// // 			clusterListOne := topke.NewEachResourceRangeListState(cluster.UUID, cluster.Name)

// // 			noMap, err := getClusterNodesOverview(cluster)
// // 			if err != nil {
// // 				return utils.NewWaitGroupResult(clusterListOne, err)
// // 			}

// // 			rets := make([]*topke.NodeOverview, 0, len(noMap))
// // 			for _, v := range noMap {
// // 				if !utils.IsContain([]string{v.NodeName, v.ClusterName, v.NodeAddr,
// // v.Extra[topke.TopkeAnnnotationClusterNameKey],
// // 					v.Extra[topke.TopkeAnnnotationTenantNameKey]}, req.Fuzzy) {
// // 					continue
// // 				}
// // 				rets = append(rets, v)
// // 			}

// // 			clusterListOne.TotalCount = len(rets)
// // 			clusterListOne.List = rets
// // 			return utils.NewWaitGroupResult(clusterListOne, nil)
// // 		}))
// // 	}
// // 	wg.Wait()

// // 	if req.SortBy == "" {
// // 		req.SortBy = "CpuUsed"
// // 	}
// // 	CutPagingSliceFromWgResultsV2(wg, &resp.EachRangeListState, &resp.NodeList, req.PageNumber, req.PageSize,
// // &resp.TotalCount, req.SortBy, req.SortDesc)
// // 	return resp
// // }

// /*
// *
// 返回整个集群的pod的cpu memory netread netwrite
// 返回值对应关系 ： namespace podname podinfo
// 没有namespace或者podname的container不会被返回
// */
// func (k *kubernetesService) podRealtimeInfo(ctx context.Context, cluster *iapiserver.Cluster)(map[string]map[string]*topke.PodOverview, error) {
// 	cpuMetrics := &topke.Metric{}
// 	memoryMetrics := &topke.Metric{}
// 	netReadMetrics := &topke.Metric{}
// 	netWriteMetrics := &topke.Metric{}

// 	promConf, _, err := tm.getPrometheusConfigProqlOption(cluster, nil, "")
// 	if err != nil {
// 		return nil, err
// 	}

// 	wg := utils.NewWaitGroup(nil)
// 	wg.Start(utils.WaitGroupRoutineFunc{Call: func() utils.WaitGroupResult { // cpu
// 		cpuMetrics, err = libprometheus.GetAutoMetricsRealTime(ctx, promConf,
// "rate(container_cpu_usage_seconds_total{}[1m])")
// 		if err != nil {
// 			logrus.Errorf("%v", err)
// 			return utils.WaitGroupResult{}
// 		}
// 		return utils.WaitGroupResult{}
// 	}})
// 	wg.Start(utils.WaitGroupRoutineFunc{Call: func() utils.WaitGroupResult { // memory
// 		memoryMetrics, err = libprometheus.GetAutoMetricsRealTime(ctx, promConf, "container_memory_usage_bytes{}")
// 		if err != nil {
// 			logrus.Errorf("%v", err)
// 			return utils.WaitGroupResult{}
// 		}
// 		return utils.WaitGroupResult{}
// 	}})
// 	wg.Start(utils.WaitGroupRoutineFunc{Call: func() utils.WaitGroupResult { // net read
// 		netReadMetrics, err = libprometheus.GetAutoMetricsRealTime(ctx, promConf,
// "rate(container_network_receive_bytes_total{}[1m])")
// 		if err != nil {
// 			logrus.Errorf("%v", err)
// 			return utils.WaitGroupResult{}
// 		}
// 		return utils.WaitGroupResult{}
// 	}})
// 	wg.Start(utils.WaitGroupRoutineFunc{Call: func() utils.WaitGroupResult { // net write
// 		netWriteMetrics, err = libprometheus.GetAutoMetricsRealTime(ctx, promConf,
// "rate(container_network_transmit_bytes_total{}[1m])")
// 		if err != nil {
// 			logrus.Errorf("%v", err)
// 			return utils.WaitGroupResult{}
// 		}
// 		return utils.WaitGroupResult{}
// 	}})
// 	wg.Wait()

// 	var namespace string
// 	var podname string
// 	nsNamepod := map[string]map[string]*topke.PodOverview{} //
// 	callBack := func(v *topke.MetricValue) (*topke.PodOverview, float64) {
// 		var value float64
// 		if v == nil {
// 			return nil, 0
// 		}

// 		namespace = v.Metadata["namespace"]
// 		podname = v.Metadata["pod"]
// 		if podname == "" || namespace == "" {
// 			return nil, 0
// 		}

// 		namePod, ok := nsNamepod[namespace]
// 		if !ok {
// 			namePod = map[string]*topke.PodOverview{}
// 			nsNamepod[namespace] = namePod
// 		}
// 		pod, ok := namePod[podname]
// 		if !ok {
// 			pod = &topke.PodOverview{}
// 			namePod[podname] = pod
// 		}
// 		pod.Namespace = namespace
// 		pod.Name = podname
// 		if v.Sample != nil {
// 			value = v.Sample.Value()
// 		}

// 		return pod, value
// 	}

// 	if cpuMetrics != nil {
// 		for _, v := range cpuMetrics.MetricValues { // cpu
// 			podOver, value := callBack(&v)
// 			if podOver != nil {
// 				podOver.Cpu.Used += value
// 			}
// 		}
// 	}

// 	if memoryMetrics != nil {
// 		for _, v := range memoryMetrics.MetricValues { // memory
// 			podOver, value := callBack(&v)
// 			if podOver != nil {
// 				podOver.Memory.Used += value
// 			}
// 		}
// 	}

// 	if netReadMetrics != nil {
// 		for _, v := range netReadMetrics.MetricValues { // net read
// 			podOver, value := callBack(&v)
// 			if podOver != nil {
// 				podOver.NetRead.Used += value
// 			}
// 		}
// 	}

// 	if netWriteMetrics != nil {
// 		for _, v := range netWriteMetrics.MetricValues { // net write
// 			podOver, value := callBack(&v)
// 			if podOver != nil {
// 				podOver.NetWrite.Used += value
// 			}
// 		}
// 	}

// 	return nsNamepod, nil
// }

// func (k *kubernetesService) MonitorPodContainerHistoryList(ctx context.Context, req *topke.MonistoringPodContainerHistoryListRequest) *topke.MonistoringPodContainerHistoryListResponse {
// 	resp := &topke.MonistoringPodContainerHistoryListResponse{
// 		Status: status.SuccessStatus,
// 	}

// 	//podName
// 	if req.Name == "" {
// 		resp.Status = status.NewStatusDesc(scode.TopECParameterEmpty, "monitoring.name is empty")
// 		return resp
// 	}

// 	timeRange := req.MonitoringTimeRange.ToRange()
// 	metricValues, err := tm.autoMetricHistoryList(ctx, req.ClusterUUID, timeRange, combineSql(req.Name,podContainerMetrics...)...)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.TimeRange = &timeRange
// 	resp.Metrics = make([]topke.Metric, 0, len(metricValues))
// 	for _, v := range metricValues {
// 		resp.Metrics = append(resp.Metrics, *v)
// 	}

// 	return resp
// }

// func (k *kubernetesService) MonitorControllerManagerGet(ctx context.Context, req *topke.MonitoringGetRequest)*topke.MonitoringControllerManagerResponse {
// 	resp := &topke.MonitoringControllerManagerResponse{Status: status.SuccessStatus}
// 	timeRange := req.MonitoringTimeRange.ToRange()
// 	metricValues, err := tm.autoMetricHistoryList(ctx, req.ClusterUUID, timeRange, controllerManagerQueryMetrics...)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.TimeRange = &timeRange
// 	resp.UpCount = int(tm.getMetricsLastSum(metricValues[cm_up_count.Name]))
// 	resp.WorkQueueAddRate = *tm.getMetricRangeAvg(metricValues[cm_work_queue_add_rate.Name])
// 	resp.CpuUsed = *tm.getMetricRangeAvg(metricValues[cm_cpu_used.Name])
// 	resp.MemoryUsed = *tm.getMetricRangeAvg(metricValues[cm_memory_used.Name])
// 	resp.Goroutine = *tm.getMetricRangeAvg(metricValues[cm_goroutines.Name])
// 	return resp
// }

// func (k *kubernetesService) MonitorSchedulerGet(ctx context.Context, req *topke.MonitoringGetRequest)*topke.MonitoringSchedulerResponse {
// 	resp := &topke.MonitoringSchedulerResponse{Status: status.SuccessStatus}
// 	timeRange := req.MonitoringTimeRange.ToRange()
// 	metricValues, err := tm.autoMetricHistoryList(ctx, req.ClusterUUID, timeRange, schedulerQueryMetrics...)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.TimeRange = &timeRange
// 	resp.UpCount = int(tm.getMetricsLastSum(getMetricsValueFromMetrics(metricValues, scheduler_up_count.Name)))
// 	resp.ScheduleCount = *tm.getMetricRangeAvg(getMetricsValueFromMetrics(metricValues, scheduler_schedule_count.Name))
// 	resp.CpuUsed = *tm.getMetricRangeAvg(getMetricsValueFromMetrics(metricValues, scheduler_cpu_used.Name))
// 	resp.MemoryUsed = *tm.getMetricRangeAvg(getMetricsValueFromMetrics(metricValues, scheduler_memory_used.Name))
// 	resp.Goroutine = *tm.getMetricRangeAvg(getMetricsValueFromMetrics(metricValues, scheduler_goroutines.Name))
// 	return resp
// }

// func (k *kubernetesService) MonitorApiserverGet(ctx context.Context, req *topke.MonitoringGetRequest)*topke.MonitoringApiserverResponse {
// 	resp := &topke.MonitoringApiserverResponse{Status: status.SuccessStatus}

// 	timeRange := req.MonitoringTimeRange.ToRange()
// 	metricValues, err := tm.autoMetricHistoryList(ctx, req.ClusterUUID, timeRange, apiserverQueryMetrics...)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.TimeRange = &timeRange
// 	resp.CpuUsed = *tm.getMetricRangeAvg(metricValues[apiserver_cpu_used.Name])
// 	resp.MemoryUsed = *tm.getMetricRangeAvg(metricValues[apiserver_memory_used.Name])
// 	resp.Goroutine = *tm.getMetricRangeAvg(metricValues[apiserver_goroutines.Name])
// 	return resp
// }

// func (k *kubernetesService) MonitorKubeletGet(ctx context.Context, req *topke.MonitoringGetRequest)*topke.MonitoringKubeletResponse {
// 	resp := &topke.MonitoringKubeletResponse{Status: status.SuccessStatus}

// 	//当前场景为主机名
// 	if req.Name == "" {
// 		resp.Status = errors.Errorf("monitoring.name is empty")
// 		return resp
// 	}

// 	timeRange := req.MonitoringTimeRange.ToRange()
// 	metricValues, err := tm.autoMetricHistoryList(ctx, req.ClusterUUID, timeRange, combineSql(req.Name,kubeletQueryMetrics...)...)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.TimeRange = &timeRange
// 	resp.UpCount = int(tm.getMetricsLastSum(metricValues[kubelet_up_count.Name]))
// 	resp.PodCount = int(tm.getMetricsLastSum(metricValues[kubelet_running_pod_count.Name]))
// 	resp.RpcRate = *tm.getMetricRangeAvg(metricValues[kubelet_rpc_rate.Name])
// 	return resp
// }

// func (k *kubernetesService) autoMetricHistoryList(ctx context.Context, clusterUUID string, timeRange prometheusv1.Range,metrics ...KV) (map[string]*topke.Metric, error) {
// 	resp := map[string]*topke.Metric{}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if !clientset.IsMonitorServiceReady(clusterInfo.UUID) {
// 		return nil, status.NewStatusDesc(scode.TopECTopkeKubeMonitorNotHealth, fmt.Sprintf("cluster_uuid[%v]",clusterInfo.UUID))
// 	}
// 	prometheusAddress, err := tm.getPrometheusUrlFromPrometheusService(clusterInfo)
// 	if err != nil {
// 		return nil, err
// 	}

// 	wg := utils.NewWaitGroup(ctx)
// 	for _, metric := range metrics {
// 		metric := metric
// 		metricValue := &topke.Metric{MetricName: metric.Name}
// 		resp[metric.Name] = metricValue
// 		wg.Start(utils.WaitGroupRoutineFunc{Call: func() utils.WaitGroupResult {
// 			metrics, err := libprometheus.GetAutoMetricsHistory2(ctx, &topke.PrometheusConfig{Address: prometheusAddress},metric.Value, timeRange)
// 			if err != nil {
// 				metricValue.Error = err.Error()
// 				logrus.Errorf("%s", err)
// 				return utils.WaitGroupResult{}
// 			}
// 			metricValue.MetricValues = metrics.MetricValues
// 			for k := range metricValue.MetricValues {
// 				metricValue.MetricValues[k].Name = metricValue.MetricName
// 			}
// 			return utils.WaitGroupResult{}
// 		}})
// 	}
// 	wg.Wait()

// 	return resp, nil
// }

// // func (k *kubernetesService) autoMetricRealtimeList(ctx context.Context, clusterInfo *iapiserver.Cluster, metrics ...KV)(map[string]*topke.Metric, error) {
// // 	resp := map[string]*topke.Metric{}

// // 	proConf, _, err := tm.getPrometheusConfigProqlOption(clusterInfo, nil, "")
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	wg := utils.NewWaitGroup(nil)
// // 	for _, metric := range metrics {
// // 		metric := metric
// // 		metricValue := &topke.Metric{MetricName: metric.Name}
// // 		resp[metric.Name] = metricValue
// // 		wg.Start(utils.WaitGroupRoutineFunc{Call: func() utils.WaitGroupResult {
// // 			metrics, err := libprometheus.GetAutoMetricsRealTime(ctx, proConf, metric.Value, 5)
// // 			if err != nil {
// // 				metricValue.Error = err.Error()
// // 				logrus.Errorf("sql[%v],err[%v]", metric.Value, err.Error())
// // 				return utils.WaitGroupResult{}
// // 			}
// // 			metricValue.MetricValues = metrics.MetricValues
// // 			return utils.WaitGroupResult{}
// // 		}})
// // 	}
// // 	wg.Wait()

// // 	return resp, nil
// // }

// // func (k *kubernetesService) MonitorEtcdGet(ctx context.Context, req *topke.MonitoringGetRequest)*topke.MonitoringEtcdResponse {
// // 	resp := &topke.MonitoringEtcdResponse{Status: status.SuccessStatus}
// // 	timeRange := req.MonitoringTimeRange.ToRange()
// // 	metricValues, err := tm.autoMetricHistoryList(ctx, req.ClusterUUID, timeRange, etcd_has_leader, etcd_change_total,
// // 		etcd_net_read_byte, etcd_net_write_byte, etcd_db_size, etcd_raft_proposal_applied, etcd_raft_proposal_commited,
// // 		etcd_raft_proposal_failed, etcd_raft_proposal_pending)
// // 	if err != nil {
// // 		resp.Status = err
// // 		return resp
// // 	}

// // 	resp.TimeRange = &timeRange
// // 	if tm.getMetricsLastSum(metricValues[etcd_has_leader.Name]) > 0 {
// // 		resp.HasLeader = true
// // 	}
// // 	resp.LeaderChangeCount = int(tm.getMetricsLastSum(getMetricsValueFromMetrics(metricValues, etcd_change_total.Name)))
// // 	resp.NetReadByte = *tm.getMetricRangeAvg(getMetricsValueFromMetrics(metricValues, etcd_net_read_byte.Name))
// // 	resp.NetWriteByte = *tm.getMetricRangeAvg(getMetricsValueFromMetrics(metricValues, etcd_net_write_byte.Name))
// // 	resp.MemorySize = *tm.getMetricRangeAvg(getMetricsValueFromMetrics(metricValues, etcd_db_size.Name))
// // 	resp.RaftProposalApplied = *tm.getMetricRangeAvg(getMetricsValueFromMetrics(metricValues,
// // etcd_raft_proposal_applied.Name)) 	resp.RaftProposalCommited =
// // *tm.getMetricRangeAvg(getMetricsValueFromMetrics(metricValues, etcd_raft_proposal_commited.Name))
// // 	resp.RaftProposalFailed = *tm.getMetricRangeAvg(getMetricsValueFromMetrics(metricValues,
// // etcd_raft_proposal_failed.Name)) 	resp.RaftProposalPending =
// // *tm.getMetricRangeAvg(getMetricsValueFromMetrics(metricValues, etcd_raft_proposal_pending.Name))

// // 	return resp
// // }

// func getMetricsValueFromMetrics(metrics map[string]*topke.Metric, metricName string) *topke.Metric {
// 	metric, _ := metrics[metricName]
// 	return metric
// }

// func (k *kubernetesService) getMetricsLastSum(metric *topke.Metric) float64 {
// 	resp := 0.0
// 	if metric != nil {
// 		for _, v := range metric.MetricValues {
// 			if len(v.Series) != 0 {
// 				resp += v.Series[len(v.Series)-1].Value()
// 			}
// 		}
// 	}

// 	return resp
// }

// func (k *kubernetesService) getMetricLastAvg(metric *topke.Metric) float64 {
// 	resp := 0.0
// 	if metric != nil {
// 		for _, v := range metric.MetricValues {
// 			if len(v.Series) != 0 {
// 				resp += v.Series[len(v.Series)-1].Value()
// 			}
// 		}
// 		resp = resp / float64(len(metric.MetricValues))
// 	}
// 	return resp
// }

// func (k *kubernetesService) getMetricRangeSum(metric *topke.Metric) float64 {
// 	resp := 0.0
// 	if metric != nil {
// 		for _, v := range metric.MetricValues {
// 			if len(v.Series) != 0 {
// 				resp += v.Series[len(v.Series)-1].Value()
// 			}
// 		}
// 		resp = resp / float64(len(metric.MetricValues))
// 	}
// 	return resp
// }

// func (k *kubernetesService) getMetricRange(metric *topke.Metric) float64 {
// 	if metric == nil || len(metric.MetricValues) == 0 || len(metric.MetricValues[0].Series) == 0 {
// 		return 0
// 	}
// 	return metric.MetricValues[0].Series[len(metric.MetricValues[0].Series)-1].Value()
// }

// func (k *kubernetesService) getMetricRangeAvg(metric *topke.Metric) *topke.MetricOne {
// 	resp := &topke.MetricOne{}
// 	if metric != nil {
// 		resp.MetricName = metric.MetricName
// 		resp.Error = metric.Error
// 		if len(metric.MetricValues) != 0 {
// 			resp.Series = metric.MetricValues[0].Series
// 			resp.Sample = metric.MetricValues[0].Sample
// 		}
// 	}
// 	return resp
// }

// func (k *kubernetesService) metricsToMetricOnes(metric *topke.Metric) []*topke.MetricOne {
// 	var resp []*topke.MetricOne

// 	if metric != nil {
// 		resp = make([]*topke.MetricOne, 0, len(metric.MetricValues))
// 		for _, v := range metric.MetricValues {
// 			resp = append(resp, &topke.MetricOne{
// 				MetricName: metric.MetricName,
// 				Series:     v.Series,
// 				Sample:     v.Sample,
// 				Error:      metric.Error,
// 			})
// 		}
// 	}

// 	return resp
// }

// func (k *kubernetesService) MonitorNodeMasterGet(ctx context.Context, req *topke.MonitoringGetRequest)*topke.MonitoringResponse {
// 	resp := &topke.MonitoringResponse{Status: status.SuccessStatus, NodeUsedTotal: &topke.NodeOverview{}}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	rets, err := getClusterNodesOverview(clusterInfo)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	for _, ret := range rets {
// 		if ret.Role != nodeRoleMaster {
// 			continue
// 		}
// 		ret.CpuUsedRatio = Divide(ret.CpuUsedRatio, 100)
// 		ret.MemoryUsedRatio = Divide(ret.MemoryUsedRatio, 100)
// 		ret.DiskUsedRatio = Divide(ret.DiskUsedRatio, 100)
// 		ret.PodUsedRatio = Divide(ret.PodUsedRatio, 100)
// 		resp.NodeList = append(resp.NodeList, ret)
// 		resp.NodeUsedTotal.CpuCapacity += ret.CpuCapacity
// 		resp.NodeUsedTotal.CpuUsed += ret.CpuUsed
// 		resp.NodeUsedTotal.MemoryCapacity += ret.MemoryCapacity
// 		resp.NodeUsedTotal.MemoryUsed += ret.MemoryUsed
// 		resp.NodeUsedTotal.DiskCapacity += ret.DiskCapacity
// 		resp.NodeUsedTotal.DiskUsed += ret.DiskUsed
// 		resp.NodeUsedTotal.PodCapacity += ret.PodCapacity
// 		resp.NodeUsedTotal.PodUsed += ret.PodUsed
// 	}

// 	resp.TotalCount = len(resp.NodeList)
// 	resp.NodeUsedTotal.CpuUsedRatio = Divide(resp.NodeUsedTotal.CpuUsed, resp.NodeUsedTotal.CpuCapacity)
// 	resp.NodeUsedTotal.MemoryUsedRatio = Divide(resp.NodeUsedTotal.MemoryUsed, resp.NodeUsedTotal.MemoryCapacity)
// 	resp.NodeUsedTotal.DiskUsedRatio = Divide(resp.NodeUsedTotal.DiskUsed, resp.NodeUsedTotal.DiskCapacity)
// 	resp.NodeUsedTotal.PodUsedRatio = Divide(float64(resp.NodeUsedTotal.PodUsed),
// float64(resp.NodeUsedTotal.PodCapacity))

// 	return resp
// }

// func (k *kubernetesService) NodeState(clusterInfo *iapiserver.Cluster, nodeName string, timeout int64) (*topke.MonitorData,error) {
// 	resp := &topke.MonitorData{Name: nodeName}

// 	rets, err := getClusterNodesOverview(clusterInfo)
// 	if err != nil {
// 		return nil, err
// 	}
// 	nodeOverview, ok := rets[nodeName]
// 	if !ok {
// 		return resp, nil
// 	}

// 	resp.Cpu = nodeOverview.CpuUsedRatio
// 	resp.CpuUsed = nodeOverview.CpuUsed
// 	resp.CpuTotal = nodeOverview.CpuCapacity
// 	resp.Memory = nodeOverview.MemoryUsedRatio
// 	resp.MemoryUsed = nodeOverview.MemoryUsed
// 	resp.MemoryTotal = nodeOverview.MemoryCapacity
// 	resp.DiskRatio = nodeOverview.DiskUsedRatio
// 	resp.DiskUsed = nodeOverview.DiskUsed
// 	resp.DiskTotal = nodeOverview.DiskCapacity
// 	resp.PodRatio = nodeOverview.PodUsedRatio
// 	resp.PodCount = float64(nodeOverview.PodUsed)
// 	resp.PodTotal = float64(nodeOverview.PodCapacity)

// 	return resp, nil
// }

// func (k *kubernetesService) MonitorWorkloadStateList(ctx context.Context, req *topke.MonitoringListRequest)*topke.MonitoringListResponse {
// 	resp := &topke.MonitoringListResponse{Status: status.SuccessStatus}

// 	if err := tm.checkMonitoringListRequest(req); err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var ls labels.Selector
// 	switch req.Type {
// 	case deployment:
// 		wl, err := clientset.DeploymentGet(clusterInfo, req.Namespace, req.Name, metav1.GetOptions{})
// 		if err != nil {
// 			resp.Status = err
// 			return resp
// 		}
// 		ls, _ = metav1.LabelSelectorAsSelector(wl.Spec.Selector)
// 	case statefulset:
// 		wl, err := clientset.StatefulSetGet(clusterInfo, req.Namespace, req.Name, metav1.GetOptions{})
// 		if err != nil {
// 			resp.Status = err
// 			return resp
// 		}

// 		ls, _ = metav1.LabelSelectorAsSelector(wl.Spec.Selector)
// 	case daemonset:
// 		wl, err := clientset.DaemonSetGet(clusterInfo, req.Namespace, req.Name, metav1.GetOptions{})
// 		if err != nil {
// 			resp.Status = err
// 			return resp
// 		}
// 		ls, _ = metav1.LabelSelectorAsSelector(wl.Spec.Selector)
// 	default:
// 		resp.Status = errors.Errorf("monitor.type not support")
// 		return resp
// 	}

// 	logrus.Infof("workload label selectors:%v", ls.String())
// 	podList, err := clientset.PodList(clusterInfo, req.Namespace, metav1.ListOptions{LabelSelector: ls.String()})
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	if len(podList.Items) == 0 {
// 		logrus.Info("podList Items is empty")
// 		return resp
// 	}

// 	podStateList, err := tm.MonitorPodStateList(ctx, clusterInfo, req.ToMonitoring())
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}
// 	pods, ok := podStateList[req.Namespace]
// 	if !ok {
// 		return resp
// 	}
// 	for _, v := range podList.Items {
// 		if podMetrics, ok := pods[v.Name]; ok {
// 			resp.Metrics = append(resp.Metrics, podMetrics)

// 		}
// 	}

// 	return resp
// }

// // func (k *kubernetesService) MonitorDeploymentState(ctx context.Context, clusterInfo *iapiserver.Cluster, namespace string)([]*topke.MonitorData, []string, error) {
// // 	workLoadState := []*topke.MonitorData{}
// // 	seletedPod := []string{}
// // 	if namespace == "" {
// // 		return nil, nil, errors.Errorf("namespace is empty")
// // 	}

// // 	podList := &topke.PodListResponse{}
// // 	workloadList := &topke.DeploymentListResponse{}

// // 	wg := utils.NewWaitGroup(nil)
// // 	wg.Start(utils.WaitGroupRoutineFunc{Call: func() utils.WaitGroupResult {
// // 		podList = tm.PodListAll(ctx, &topke.PodListRequest{
// // 			ResourceListRequest: topke.ResourceListRequest{
// // 				ClusterUUID: clusterInfo.UUID,
// // 				Namespace:   namespace,
// // 			},
// // 		})
// // 		if podList.Status != status.SuccessStatus {
// // 			return utils.NewWaitGroupResult(nil, status.UpdateStatus(podList.Status))
// // 		}
// // 		return utils.NewWaitGroupResult(nil, nil)
// // 	}})
// // 	wg.Start(utils.WaitGroupRoutineFunc{Call: func() utils.WaitGroupResult {
// // 		in := &topke.DeploymentListRequest{}
// // 		in.ClusterUUID = clusterInfo.UUID
// // 		in.Namespace = namespace
// // 		workloadList = GetTopkeManager().DeploymentListAll(ctx, in)
// // 		if workloadList.Status != status.SuccessStatus {
// // 			return utils.NewWaitGroupResult(nil, status.UpdateStatus(workloadList.Status))
// // 		}
// // 		return utils.NewWaitGroupResult(nil, nil)
// // 	}})
// // 	wg.Wait()

// // 	for _, v := range wg.GetResults() {
// // 		if v.Error != nil {
// // 			return nil, nil, v.Error
// // 		}
// // 	}

// // 	wg = utils.NewWaitGroup(nil)
// // 	lock := sync.RWMutex{}
// // 	for _, one := range workloadList.List {
// // 		one := one
// // 		wg.Start(utils.WaitGroupRoutineFunc{Call: func() utils.WaitGroupResult {
// // 			var objs []metav1.ObjectMeta
// // 			for _, pod := range podList.List {
// // 				if pod.Pod == nil || one.Deployment == nil || one.Deployment.Spec.Selector == nil {
// // 					continue
// // 				}

// // 				if IsMapContainMap(pod.Pod.Labels, one.Deployment.Spec.Selector.MatchLabels) {
// // 					objs = append(objs, pod.Pod.ObjectMeta)
// // 					lock.Lock()
// // 					seletedPod = append(seletedPod, pod.Pod.Name)
// // 					lock.Unlock()
// // 				}
// // 			}

// // 			ret, err := tm.MonitorPodState(ctx, clusterInfo, objs)
// // 			if err != nil {
// // 				logrus.Errorf("%v", err)
// // 				return utils.WaitGroupResult{}
// // 			}
// // 			ret.Type = deployment
// // 			ret.Namespace = one.Deployment.Namespace
// // 			ret.Name = one.Deployment.Name

// // 			lock.Lock()
// // 			workLoadState = append(workLoadState, ret)
// // 			lock.Unlock()
// // 			return utils.WaitGroupResult{}
// // 		}})
// // 	}
// // 	wg.Wait()

// // 	return workLoadState, seletedPod, nil
// // }

// func (k *kubernetesService) MonitorPodStateList(ctx context.Context, clusterInfo *iapiserver.Cluster, req *topke.Monitoring)(map[string]map[string]topke.Metric, error) {
// 	resp := map[string]map[string]topke.Metric{}

// 	metrics := []KV{
// 		namespace_pod_cpu_rate,
// 		namespace_pod_memory_size,
// 		namespace_pod_network_receive_size,
// 		namespace_pod_network_send_size,
// 	}
// 	rets, err := tm.autoMetricHistoryList(ctx, clusterInfo.UUID, req.MonitoringTimeRange.ToRange(), metrics...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, v := range rets {
// 		for _, podvalue := range v.MetricValues {
// 			namespace := podvalue.Metadata["namespace"]
// 			podname := podvalue.Metadata["pod"]

// 			if namespace == "" || podname == "" {
// 				continue
// 			}

// 			pods, ok := resp[namespace]
// 			if !ok {
// 				pods = map[string]topke.Metric{}
// 				resp[namespace] = pods
// 			}
// 			newm := pods[podname]
// 			newm.MetricValues = append(newm.MetricValues, podvalue)
// 			pods[podname] = newm
// 		}
// 	}

// 	return resp, nil
// }

// func (k *kubernetesService) MonitorNodeStateList(ctx context.Context, clusterInfo *iapiserver.Cluster)(map[string]*topke.NodeOverview, error) {
// 	resp := map[string]*topke.NodeOverview{}

// 	getNodeName := func(labels map[string]string) string {
// 		if labels == nil {
// 			return ""
// 		}

// 		if nodename1, _ := labels["node"]; nodename1 != "" {
// 			return nodename1
// 		}
// 		instance, _ := labels["instance"]
// 		return instance
// 	}

// 	metrics := []KV{node_cpu_used, node_cpu_total, node_cpu_ratio, node_cpu_request_used, node_cpu_limit_used,
// 		node_memory_used, node_memory_total, node_memory_ratio, node_memory_request_used, node_memory_limit_used,
// 		node_load_1, node_load_5, node_load_15, node_disk_io_rate, node_disk_ratio, node_disk_used, node_disk_total,
// node_disk_read_rate, node_disk_write_rate,
// 		node_pod_ratio, node_pod_used, node_pod_total, node_pod_running}
// 	ret, err := tm.autoMetricRealtimeList(ctx, clusterInfo, metrics...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, v := range metrics {
// 		metricInfo := ret[v.Name]
// 		if metricInfo == nil {
// 			continue
// 		}
// 		for _, nsInfo := range metricInfo.MetricValues {
// 			nodeName := getNodeName(nsInfo.Metadata)
// 			if nodeName == "" {
// 				continue
// 			}
// 			data, ok := resp[nodeName]
// 			if !ok {
// 				data = &topke.NodeOverview{NodeName: nodeName}
// 				resp[nodeName] = data
// 			}
// 			if nsInfo.Sample == nil {
// 				continue
// 			}
// 			switch v.Name {
// 			case node_cpu_used.Name:
// 				data.CpuUsed = nsInfo.Sample.Value() * 1000
// 			case node_cpu_total.Name:
// 				data.CpuCapacity = nsInfo.Sample.Value() * 1000
// 			case node_cpu_ratio.Name:
// 				data.CpuUsedRatio = nsInfo.Sample.Value() * 100
// 			case node_cpu_request_used.Name:
// 				data.CpuResourceRequest = int64(nsInfo.Sample.Value() * 1000)
// 			case node_cpu_limit_used.Name:
// 				data.CpuResourceLimit = int64(nsInfo.Sample.Value() * 1000)
// 			case node_memory_used.Name:
// 				data.MemoryUsed = nsInfo.Sample.Value()
// 			case node_memory_total.Name:
// 				data.MemoryCapacity = nsInfo.Sample.Value()
// 			case node_memory_ratio.Name:
// 				data.MemoryUsedRatio = nsInfo.Sample.Value()
// 			case node_memory_request_used.Name:
// 				data.MemoryResourceRequest = int64(nsInfo.Sample.Value())
// 			case node_memory_limit_used.Name:
// 				data.MemoryResourceLimit = int64(nsInfo.Sample.Value())
// 			case node_disk_ratio.Name:
// 				data.DiskUsedRatio = nsInfo.Sample.Value()
// 			case node_disk_used.Name:
// 				data.DiskUsed = nsInfo.Sample.Value()
// 			case node_disk_total.Name:
// 				data.DiskCapacity = nsInfo.Sample.Value()
// 			case node_pod_used.Name:
// 				data.PodUsed = int(nsInfo.Sample.Value())
// 			case node_pod_total.Name:
// 				data.PodCapacity = int(nsInfo.Sample.Value())
// 			case node_pod_ratio.Name:
// 				data.PodUsedRatio = nsInfo.Sample.Value()
// 			}
// 		}
// 	}

// 	return resp, nil
// }

// // 返回值map：namespace，podname，overview
// func (k *kubernetesService) MonitorClusterPodStateList(ctx context.Context, clusterInfo *iapiserver.Cluster)(map[string]map[string]*topke.PodOverview, error) {
// 	resp := map[string]map[string]*topke.PodOverview{}

// 	extra := getObjectExtraInfos(clusterInfo)
// 	metrics := []KV{cluster_pod_cpu_used, cluster_pod_mem_used, cluster_pod_net_rb, cluster_pod_net_wb}
// 	ret, err := tm.autoMetricRealtimeList(ctx, clusterInfo, metrics...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, v := range metrics {
// 		metricInfo := ret[v.Name]
// 		for _, info := range metricInfo.MetricValues {
// 			namespace := info.Metadata["namespace"]
// 			podname := info.Metadata["pod"]
// 			nodename := info.Metadata["node"]
// 			if namespace == "" || podname == "" {
// 				continue
// 			}

// 			data, ok := resp[namespace]
// 			if !ok {
// 				data = map[string]*topke.PodOverview{}
// 				resp[namespace] = data
// 			}
// 			podInfo, ok := data[podname]
// 			if !ok {
// 				podInfo = &topke.PodOverview{Extra: extra, Namespace: namespace, Name: podname, NodeName: nodename}
// 				data[podname] = podInfo
// 			}
// 			if info.Sample == nil {
// 				continue
// 			}
// 			switch v.Name {
// 			case cluster_pod_cpu_used.Name:
// 				podInfo.Cpu.Used = info.Sample.Value() * 1000
// 			case cluster_pod_mem_used.Name:
// 				podInfo.Memory.Used = info.Sample.Value()
// 			case cluster_pod_net_wb.Name:
// 				podInfo.NetWrite.Used = info.Sample.Value()
// 			case cluster_pod_net_rb.Name:
// 				podInfo.NetRead.Used = info.Sample.Value()
// 			}
// 		}
// 	}

// 	return resp, nil
// }

// func (k *kubernetesService) MonitorNamespaceStateList(ctx context.Context, clusterInfo *iapiserver.Cluster)(map[string]*topke.MonitorData, error) {
// 	resp := map[string]*topke.MonitorData{}

// 	extra := getObjectExtraInfos(clusterInfo)
// 	metrics := []KV{namespace_cpu_used, namespace_mem_used, namespace_pod_count, namespace_net_rb, namespace_net_wb}
// 	ret, err := tm.autoMetricRealtimeList(ctx, clusterInfo, metrics...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, v := range metrics {
// 		metricInfo := ret[v.Name]
// 		for _, nsInfo := range metricInfo.MetricValues {
// 			nsName := nsInfo.Metadata["namespace"]
// 			data, ok := resp[nsName]
// 			if !ok {
// 				data = &topke.MonitorData{Name: nsName}
// 				resp[nsName] = data
// 			}
// 			data.Extra = extra
// 			if nsInfo.Sample == nil {
// 				continue
// 			}
// 			switch v.Name {
// 			case namespace_cpu_used.Name:
// 				data.Cpu = nsInfo.Sample.Value()
// 			case namespace_mem_used.Name:
// 				data.Memory = nsInfo.Sample.Value()
// 			case namespace_pod_count.Name:
// 				data.PodCount = nsInfo.Sample.Value()
// 			case namespace_net_wb.Name:
// 				data.NetworkWrite = nsInfo.Sample.Value()
// 			case namespace_net_rb.Name:
// 				data.NetworkRead = nsInfo.Sample.Value()
// 			}
// 		}
// 	}

// 	return resp, nil
// }

// /*
// *
// cpu 		m
// memory 		byte
// disk 		byte
// net  		byte
// */
// func (k *kubernetesService) MonitorPodState(ctx context.Context, clusterInfo *iapiserver.Cluster, objs []metav1.ObjectMeta)(*topke.MonitorData, error) {
// 	resp := &topke.MonitorData{}

// 	if len(objs) != 0 {
// 		rets, err := tm.MonitorClusterPodStateList(ctx, clusterInfo)
// 		if err != nil {
// 			return nil, err
// 		}

// 		for _, v := range objs {
// 			if namespaceOverview, ok := rets[v.Namespace]; ok {
// 				if podOverview, ok := namespaceOverview[v.Name]; ok {
// 					resp.Cpu += podOverview.Cpu.Used
// 					resp.Memory += podOverview.Memory.Used
// 					resp.NetworkRead += podOverview.NetRead.Used
// 					resp.NetworkWrite += podOverview.NetWrite.Used
// 				}
// 			}
// 		}
// 	}

// 	return resp, nil
// }

// func (k *kubernetesService) namespaceOverview(ctx context.Context, clusterInfo *iapiserver.Cluster, namespace string) (resp[]*topke.MonitorData, err error) {
// 	selectPods := map[string]bool{}

// 	podList := tm.PodListAll(ctx, &topke.PodListRequest{
// 		ResourceListRequest: topke.ResourceListRequest{
// 			ClusterUUID: clusterInfo.UUID,
// 			Namespace:   namespace,
// 		},
// 	})
// 	if podList.Status != status.SuccessStatus {
// 		return nil, status.UpdateStatus(podList.Status)
// 	}

// 	deps, podnames, e := tm.MonitorDeploymentState(ctx, clusterInfo, namespace)
// 	if e != nil {
// 		return nil, status.UpdateStatus(e)
// 	}
// 	resp = append(resp, deps...)
// 	for _, v := range podnames {
// 		selectPods[v] = true
// 	}

// 	stats, podnames, e := tm.MonitorStatefulSetState(ctx, clusterInfo, namespace)
// 	if e != nil {
// 		return nil, status.UpdateStatus(e)
// 	}
// 	resp = append(resp, stats...)
// 	for _, v := range podnames {
// 		selectPods[v] = true
// 	}
// 	daes, podnames, e := tm.MonitorDaemonsetState(ctx, clusterInfo, namespace)
// 	if e != nil {
// 		return nil, status.UpdateStatus(e)
// 	}
// 	resp = append(resp, daes...)
// 	for _, v := range podnames {
// 		selectPods[v] = true
// 	}

// 	wg := utils.NewWaitGroup(nil)
// 	gLock := sync.RWMutex{}
// 	for _, v := range podList.List {
// 		if v.Pod == nil {
// 			continue
// 		}
// 		if _, ok := selectPods[v.Pod.Name]; ok {
// 			continue
// 		}

// 		v := v
// 		wg.Start(utils.WaitGroupRoutineFunc{Call: func() utils.WaitGroupResult {
// 			ret, err := tm.MonitorPodState(ctx, clusterInfo, []metav1.ObjectMeta{v.Pod.ObjectMeta})
// 			if err != nil {
// 				logrus.Errorf("%s", err)
// 				return utils.WaitGroupResult{}
// 			}

// 			ret.Namespace = v.Pod.ObjectMeta.Namespace
// 			ret.Name = v.Pod.ObjectMeta.Name
// 			ret.Type = pod

// 			gLock.Lock()
// 			resp = append(resp, ret)
// 			gLock.Unlock()
// 			return utils.WaitGroupResult{}
// 		}})
// 	}
// 	wg.Wait()

// 	return
// }

// var (
// 	grafanalocker       sync.Mutex
// 	grafanaDashIdMap    = map[string]map[string]string{}
// 	grafanaTimeout      = time.Duration(10)
// 	grafanaOrgId        = int64(1)
// 	grafanaDataSource   = "var-datasource=prometheus"
// 	grafanaDefaultPara  = "?orgId=1&refresh=10s"
// 	grafanaScheme       = "http://"
// 	grafanaDefultDashId = "kubernetes-compute-resources-cluster" // 默认参数
// )

// func (k *kubernetesService) getGrafanaUrlFromGrafanaService(clusterInfo *iapiserver.Cluster) (string, error) {
// 	if clusterInfo == nil {
// 		return "", status.NewStatusDesc(scode.TopECParameterEmpty, "clusterInfo is empty")
// 	}

// 	service, err := libkubernetes.ServiceGet(clusterInfo.Config, monitorNamespace, grafanaServiceName, metav1.GetOptions{})
// 	if err != nil {
// 		return "", err
// 	}
// 	if service.Spec.Type != corev1.ServiceType(nodePort) {
// 		service.Spec.Type = corev1.ServiceType(nodePort)
// 		service, err = libkubernetes.ServiceUpdate(clusterInfo.Config, monitorNamespace, service, metav1.UpdateOptions{})
// 		if err != nil {
// 			return "", err
// 		}
// 		time.Sleep(1 * time.Second)
// 	}

// 	if service.Spec.Type != corev1.ServiceType(nodePort) {
// 		return "", status.NewStatusDesc(scode.TopECTopkePrometheusNotExist,
// fmt.Sprintf("namespace:monitoring,service:grafana,type:NodePort"))
// 	}
// 	if len(service.Spec.Ports) == 0 {
// 		return "", status.NewStatusDesc(scode.TopECTopkePrometheusNotExist,
// fmt.Sprintf("namespace:monitoring,service:grafana,type:NodePort"))
// 	}

// 	url := grafanaScheme + utils.ParseAddrFromURLNoError(clusterInfo.Config.Host) + ":" +
// strconv.Itoa(int(service.Spec.Ports[0].NodePort))
// 	return url, nil
// }

// func (k *kubernetesService) GrafanaDashUrlGet(ctx context.Context, req *topke.GrafanaGetRequest) *topke.GrafanaResponse {
// 	resp := &topke.GrafanaResponse{Status: status.SuccessStatus}

// 	if req.SubUrl == "" {
// 		req.SubUrl = grafanaDefultDashId
// 	}

// 		cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	resp.TotalUrl, err = tm.createDashUrl(clusterInfo, req.SubUrl)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.TotalUrl += "?" + createDashPara(req.Namespace, req.Workload, req.Node, req.Pod)
// 	return resp
// }

// func createDashPara(namespace, workload, node, pod string) string {
// 	resp := grafanaDataSource

// 	add := func(src string, flag string, tar string) string {
// 		if tar != "" {
// 			src = src + "&" + flag + "=" + tar
// 		}
// 		return src
// 	}

// 	resp = add(resp, "var-namespace", namespace)
// 	resp = add(resp, "var-workload", workload)
// 	resp = add(resp, "var-node", node)
// 	resp = add(resp, "var-pod", pod)
// 	return resp
// }

// func (k *kubernetesService) createDashUrl(clusterInfo *iapiserver.Cluster, subUrl string) (string, error) {
// 	grafanaUrl, err := tm.getGrafanaUrlFromGrafanaService(clusterInfo)
// 	if err != nil {
// 		return "", err
// 	}

// 	// find urls
// 	var totalUrl string
// 	urlMap, err := tm.CacheGet(clusterInfo)
// 	if err != nil {
// 		return "", err
// 	}
// 	for k, v := range urlMap {
// 		if strings.HasSuffix(k, subUrl) {
// 			totalUrl = v
// 			break
// 		}
// 	}
// 	if totalUrl == "" {
// 		return "", errors.Errorf("cluster_uuid[%v] grafana dash sub_url[%v] not exist", clusterInfo.UUID, subUrl)
// 	}

// 	totalUrl = grafanaUrl + totalUrl + grafanaDefaultPara
// 	return totalUrl, nil
// }

// func (k *kubernetesService) CacheGet(clusterInfo *iapiserver.Cluster) (map[string]string, error) {
// 	grafanalocker.Lock()
// 	cacheUrls := grafanaDashIdMap[clusterInfo.UUID]
// 	urls := make(map[string]string, len(cacheUrls))
// 	for k, v := range cacheUrls {
// 		urls[k] = v
// 	}
// 	grafanalocker.Unlock()

// 	if len(urls) != 0 {
// 		go tm.CacheUpdate(clusterInfo)
// 		return urls, nil
// 	}

// 	if err := tm.CacheUpdate(clusterInfo); err != nil {
// 		return nil, err
// 	}
// 	grafanalocker.Lock()
// 	defer grafanalocker.Unlock()
// 	cacheUrls = grafanaDashIdMap[clusterInfo.UUID] //CacheUpdate() not error，cluster_uuid must exist
// 	for k, v := range cacheUrls {
// 		urls[k] = v
// 	}

// 	return urls, nil
// }

// func (k *kubernetesService) CacheUpdate(clusterInfo *iapiserver.Cluster) error {
// 	urls, err := tm.loadGrafanaUrls(clusterInfo)
// 	if err != nil {
// 		return err
// 	}

// 	grafanaDashIdMap[clusterInfo.UUID] = urls
// 	return nil
// }

// func (k *kubernetesService) loadGrafanaUrls(clusterInfo *iapiserver.Cluster) (map[string]string, error) {
// 	grafanaUrl, err := tm.getGrafanaUrlFromGrafanaService(clusterInfo)
// 	if err != nil {
// 		return nil, err
// 	}

// 	config := gapi.Config{}
// 	config.OrgID = grafanaOrgId
// 	config.Client = &http.Client{Timeout: grafanaTimeout * time.Second}
// 	client, err := gapi.New(grafanaUrl, config)
// 	if err != nil {
// 		return nil, status.NewStatusDesc(scode.TopECTopkeKubeResourceGetError, err.Error())
// 	}

// 	dashboards, err := client.Dashboards()
// 	if err != nil {
// 		return nil, status.NewStatusDesc(scode.TopECTopkeKubeResourceGetError, err.Error())
// 	}

// 	// update cache
// 	dashIdMap := map[string]string{}
// 	for _, v := range dashboards {
// 		ss := strings.Split(v.URL, "/")
// 		dashIdMap[ss[len(ss)-1]] = v.URL
// 	}

// 	return dashIdMap, nil
// }

// func (k *kubernetesService) AlertManagerCreate(ctx context.Context, req *topke.AlertManagerRequest)*topke.AlertManagerResponse {
// 	resp := &topke.AlertManagerResponse{Status: status.SuccessStatus}

// 	if req.AlertManager == nil {
// 		req.AlertManager = &monitoringv1.Alertmanager{}
// 		req.AlertManager.Namespace = monitorNamespace
// 	}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	meta, err := libkubernetes.AlertManagerCreate(cluster.Config, req.AlertManager.Namespace, req.AlertManager,req.CreateOpts)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.Info = convertAlertManagerToApi(meta)
// 	return resp
// }

// func (k *kubernetesService) AlertManagerDelete(ctx context.Context, req *topke.AlertManagerRequest)*topke.AlertManagerResponse {
// 	resp := &topke.AlertManagerResponse{
// 		Status: status.SuccessStatus,
// 	}

// 	if req.AlertManager == nil {
// 		resp.Status = status.NewStatusDesc(scode.TopECParameterEmpty, "AlertManager")
// 		return resp
// 	}

// 	req.AlertManager.ResourceVersion = ""

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := libkubernetes.AlertManagerDelete(cluster.Config, req.AlertManager.Namespace, req.AlertManager.Name,req.DeleteOpts); err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	return resp
// }

// func (k *kubernetesService) AlertManagerUpdate(ctx context.Context, req *topke.AlertManagerRequest)*topke.AlertManagerResponse {
// 	resp := &topke.AlertManagerResponse{
// 		Status: status.SuccessStatus,
// 	}

// 	if req.AlertManager == nil {
// 		resp.Status = status.NewStatusDesc(scode.TopECParameterEmpty, "AlertManager")
// 		return resp
// 	}
// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	meta, err := libkubernetes.AlertManagerUpdate(cluster.Config, req.AlertManager.Namespace, req.AlertManager,req.UpdateOpts)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.Info = convertAlertManagerToApi(meta)
// 	return resp
// }

// func (k *kubernetesService) AlertManagerGet(ctx context.Context, req *topke.AlertManagerGetRequest)*topke.AlertManagerResponse {
// 	resp := &topke.AlertManagerResponse{
// 		Status: status.SuccessStatus,
// 	}

// 	if err := validateNamespacedScopeParameters(req.Namespace, req.Name); err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	meta, err := libkubernetes.AlertManagerGet(cluster.Config, req.Namespace, req.Name, req.ToGetOpts())
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.Info = convertAlertManagerToApi(meta)
// 	return resp
// }

// func (k *kubernetesService) AlertManagerList(ctx context.Context, req *topke.AlertManagerListRequest)*topke.AlertManagerListResponse {
// 	resp := &topke.AlertManagerListResponse{Status: status.SuccessStatus}

// 	if req.Namespace == "" {
// 		req.Namespace = monitorNamespace
// 	}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ret, err := libkubernetes.AlertManagerList(cluster.Config, req.Namespace, req.ToListOpts(),libkubernetes.DefaultTimeoutGet)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	var sms []*topke.AlertManagerInfo
// 	for _, dm := range ret.Items {
// 		if !utils.IsContain([]string{dm.Name}, req.Fuzzy) {
// 			continue
// 		}

// 		tmp := convertAlertManagerToApi(&dm)
// 		sms = append(sms, tmp)
// 	}

// 	resp.List = sms
// 	s, e := utils.PagingIndex(len(sms), req.PageNumber, req.PageSize)
// 	resp.List = resp.List[s:e]
// 	return resp
// }

// func convertAlertManagerToApi(meta *monitoringv1.Alertmanager) *topke.AlertManagerInfo {
// 	ai := &topke.AlertManagerInfo{
// 		AlertManager: meta,
// 	}
// 	return ai
// }

// func (k *kubernetesService) PodMonitorCreate(ctx context.Context, req *topke.PodMonitorRequest) *topke.PodMonitorResponse{
// 	resp := &topke.PodMonitorResponse{Status: status.SuccessStatus}

// 	if req.PodMonitor == nil {
// 		resp.Status = status.NewStatusDesc(scode.TopECParameterEmpty, "name is empty")
// 		return resp
// 	}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	meta, err := libkubernetes.PodMonitorCreate(cluster.Config, req.PodMonitor.Namespace, req.PodMonitor, req.CreateOpts)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.Info = convertK8sPodMonitorToApiPodMonitor(meta, cluster, req.Yaml)
// 	return resp
// }

// func (k *kubernetesService) PodMonitorDelete(ctx context.Context, req *topke.PodMonitorRequest) *topke.PodMonitorResponse{
// 	resp := &topke.PodMonitorResponse{
// 		Status: status.SuccessStatus,
// 	}

// 	if req.PodMonitor == nil {
// 		return nil,errors.Errorf( "name is empty")
// 	}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := libkubernetes.PodMonitorDelete(cluster.Config, req.PodMonitor.Namespace, req.PodMonitor.Name, req.DeleteOpts);err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	return resp
// }

// func (k *kubernetesService) PodMonitorUpdate(ctx context.Context, req *topke.PodMonitorRequest) *topke.PodMonitorResponse{
// 	resp := &topke.PodMonitorResponse{
// 		Status: status.SuccessStatus,
// 	}

// 	if req.PodMonitor == nil {
// 		resp.Status = status.NewStatusDesc(scode.TopECParameterEmpty, "name is empty")
// 		return resp
// 	}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	meta, err := libkubernetes.PodMonitorUpdate(cluster.Config, req.PodMonitor.Namespace, req.PodMonitor, req.UpdateOpts)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.Info = convertK8sPodMonitorToApiPodMonitor(meta, cluster, req.Yaml)
// 	return resp
// }

// func (k *kubernetesService) PodMonitorGet(ctx context.Context, req *topke.PodMonitorGetRequest) *topke.PodMonitorResponse{
// 	resp := &topke.PodMonitorResponse{Status: status.SuccessStatus}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}
// 	meta, err := libkubernetes.PodMonitorGet(cluster.Config, req.Namespace, req.Name, req.ToGetOpts())
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.Info = convertK8sPodMonitorToApiPodMonitor(meta, cluster, req.Yaml)
// 	return resp
// }

// func (k *kubernetesService) PodMonitorList(ctx context.Context, req *topke.PodMonitorListRequest)*topke.PodMonitorListResponse {
// 	resp := &topke.PodMonitorListResponse{Status: status.SuccessStatus}

// 	if req.Namespace == "" {
// 		req.Namespace = monitorNamespace
// 	}
// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ret, err := libkubernetes.PodMonitorList(cluster.Config, req.Namespace, req.ToListOpts(), libkubernetes.DefaultTimeoutGet)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	var resInfos []*topke.PodMonitorInfo
// 	for _, dm := range ret.Items {
// 		if !utils.IsContain([]string{dm.Name}, req.Fuzzy) {
// 			continue
// 		}

// 		resInfos = append(resInfos, convertK8sPodMonitorToApiPodMonitor(dm, cluster, req.Yaml))
// 	}

// 	resp.List = resInfos
// 	s, e := utils.PagingIndex(len(resp.List), req.PageNumber, req.PageSize)
// 	resp.List = resp.List[s:e]
// 	return resp
// }

// func convertK8sPodMonitorToApiPodMonitor(meta *monitoringv1.PodMonitor, cluster *iapiserver.Cluster, yaml bool)*topke.PodMonitorInfo {
// 	return &topke.PodMonitorInfo{
// 		PodMonitor: meta,
// }
// }

// func (k *kubernetesService) PrometheusResourceCreate(ctx context.Context, req *topke.PrometheusRequest)*topke.PrometheusResponse {
// 	resp := &topke.PrometheusResponse{
// 		Status: status.SuccessStatus,
// 	}

// 	if req.Prometheus == nil {
// 		req.Prometheus = &monitoringv1.Prometheus{}
// 		req.Prometheus.Namespace = monitorNamespace
// 	}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	meta, err := libkubernetes.PrometheusCreate(cluster.Config, req.Prometheus.Namespace, req.Prometheus, req.CreateOpts)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.Info = convertK8sPrometheusToApiPrometheus(meta, cluster, req.Yaml)
// 	return resp
// }

// func (k *kubernetesService) PrometheusResourceDelete(ctx context.Context, req *topke.PrometheusRequest)*topke.PrometheusResponse {
// 	resp := &topke.PrometheusResponse{
// 		Status: status.SuccessStatus,
// 	}

// 	if req.Prometheus == nil {
// 		resp.Status = status.NewStatusDesc(scode.ScodeManagerCommonParameterError, "name is empty")
// 		return resp
// 	}

// 	req.Prometheus.ResourceVersion = ""
// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err = libkubernetes.PrometheusDelete(cluster.Config, req.Prometheus.Namespace, req.Prometheus.Name, req.DeleteOpts); err!= nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	return resp
// }

// func (k *kubernetesService) PrometheusResourceUpdate(ctx context.Context, req *topke.PrometheusRequest)*topke.PrometheusResponse {
// 	resp := &topke.PrometheusResponse{
// 		Status: status.SuccessStatus,
// 	}

// 	if req.Prometheus == nil {
// 		resp.Status = status.NewStatusDesc(scode.ScodeManagerCommonParameterError, "name is empty")
// 		return resp
// 	}
// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	old, err := libkubernetes.PrometheusGet(cluster.Config, req.Prometheus.Namespace, req.Prometheus.Name, metav1.GetOptions{})
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	new := req.Prometheus
// 	old.Spec.Image = new.Spec.Image
// 	old.Spec.Replicas = new.Spec.Replicas
// 	old.Spec.Resources = new.Spec.Resources
// 	old.Spec.Retention = new.Spec.Retention
// 	old.Spec.RetentionSize = new.Spec.RetentionSize
// 	old.Spec.NodeSelector = new.Spec.NodeSelector

// 	meta, err := libkubernetes.PrometheusUpdate(cluster.Config, req.Prometheus.Namespace, old, req.UpdateOpts)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.Info = convertK8sPrometheusToApiPrometheus(meta, cluster, req.Yaml)
// 	return resp
// }

// func (k *kubernetesService) PrometheusResourceGet(ctx context.Context, req *topke.PrometheusGetRequest)*topke.PrometheusResponse {
// 	resp := &topke.PrometheusResponse{Status: status.SuccessStatus}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	meta, err := libkubernetes.PrometheusGet(cluster.Config, req.Namespace, req.Name, req.ToGetOpts())
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.Info = convertK8sPrometheusToApiPrometheus(meta, cluster, req.Yaml)
// 	if resp.Info.GrafanaUrl, err = tm.getGrafanaUrlFromGrafanaService(cluster); err != nil {
// 		logrus.Errorf("getgrafanaurl err : clusteruuid[%v],clustername[%v],err[%v]", cluster.UUID, cluster.Name,err.Error())
// 	}
// 	if resp.Info.Url, err = tm.getPrometheusUrlFromPrometheusService(cluster); err != nil {
// 		logrus.Errorf("get prometheusurl err : clusteruuid[%v],clustername[%v],err[%v]", cluster.UUID, cluster.Name,err.Error())
// 	}

// 	return resp
// }

// func (k *kubernetesService) PrometheusResourceList(ctx context.Context, req *topke.PrometheusListRequest)*topke.PrometheusListResponse {
// 	resp := &topke.PrometheusListResponse{
// 		Status: status.SuccessStatus,
// 	}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ret, err := libkubernetes.PrometheusList(cluster.Config, req.Namespace, req.ToListOpts(), libkubernetes.DefaultTimeoutGet)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	var resInfos []*topke.PrometheusInfo
// 	for _, dm := range ret.Items {
// 		if !utils.IsContain([]string{dm.Name}, req.Fuzzy) {
// 			continue
// 		}

// 		resInfos = append(resInfos, convertK8sPrometheusToApiPrometheus(dm, cluster, req.Yaml))
// 	}

// 	resp.List = resInfos
// 	s, e := utils.PagingIndex(len(resp.List), req.PageNumber, req.PageSize)
// 	resp.List = resp.List[s:e]
// 	return resp
// }

// func convertK8sPrometheusToApiPrometheus(meta *monitoringv1.Prometheus, cluster *iapiserver.Cluster, yaml bool)*topke.PrometheusInfo {

// 	pi := &topke.PrometheusInfo{
// 		Prometheus: meta,
// 	}

// 	pi.ResourceConvert = &topke.PrometheusResourceConvert{
// 		Limits:  make(map[string]int64),
// 		Request: make(map[string]int64),
// 	}
// 	if pi.Prometheus.Spec.Resources.Requests != nil {
// 		for i, j := range pi.Prometheus.Spec.Resources.Requests {
// 			if string(i) == "memory" {
// 				pi.ResourceConvert.Request[string(i)] = j.Value() / 1024 / 1024 // convert to memory
// 				continue
// 			}
// 			if string(i) == "cpu" { // if use Value(), 0.1 cpu/100m cpu will convert to 1 cpu
// 				pi.ResourceConvert.Request[string(i)] = j.MilliValue()
// 				continue
// 			}
// 			pi.ResourceConvert.Request[string(i)] = j.Value()
// 		}
// 	}

// 	if pi.Prometheus.Spec.Resources.Limits != nil {
// 		for i, j := range pi.Prometheus.Spec.Resources.Limits {
// 			if string(i) == "memory" {
// 				pi.ResourceConvert.Limits[string(i)] = j.Value() / 1024 / 1024 // convert to memory
// 				continue
// 			}
// 			if string(i) == "cpu" { // if use Value(), 0.1 cpu/100m cpu will convert to 1 cpu
// 				pi.ResourceConvert.Limits[string(i)] = j.MilliValue()
// 				continue
// 			}
// 			pi.ResourceConvert.Limits[string(i)] = j.Value()
// 		}
// 	}

// 	return pi
// }

// const (
// 	monitorNamespace = "monitoring"
// )

// func (k *kubernetesService) ServiceMonitorCreate(ctx context.Context, req *topke.ServiceMonitorRequest)*topke.ServiceMonitorResponse {
// 	resp := &topke.ServiceMonitorResponse{
// 		Status: status.SuccessStatus,
// 	}

// 	if req.ServiceMonitor == nil {
// 		resp.Status = errors.Errorf( "name is empty")
// 		return resp
// 	}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	meta, err := libkubernetes.ServiceMonitorCreate(cluster.Config, req.ServiceMonitor.Namespace, req.ServiceMonitor,req.CreateOpts)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.Info = convertK8sServerMonitorToTopkeServiceMonitor(meta)
// 	return resp
// }

// func (k *kubernetesService) ServiceMonitorDelete(ctx context.Context, req *topke.ServiceMonitorRequest)*topke.ServiceMonitorResponse {
// 	resp := &topke.ServiceMonitorResponse{
// 		Status: status.SuccessStatus,
// 	}

// 	if req.ServiceMonitor == nil {
// 		resp.Status = status.NewStatusDesc(scode.ScodeManagerCommonParameterError, "name is empty")
// 		return resp
// 	}

// 	req.ServiceMonitor.ResourceVersion = ""

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := libkubernetes.ServiceMonitorDelete(cluster.Config, req.ServiceMonitor.Namespace, req.ServiceMonitor.Name,req.DeleteOpts); err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	return resp
// }

// func (k *kubernetesService) ServiceMonitorUpdate(ctx context.Context, req *topke.ServiceMonitorRequest)*topke.ServiceMonitorResponse {
// 	resp := &topke.ServiceMonitorResponse{Status: status.SuccessStatus}

// 	if req.ServiceMonitor == nil {
// 		resp.Status = status.NewStatusDesc(scode.ScodeManagerCommonParameterError, "name is empty")
// 		return resp
// 	}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	meta, err := libkubernetes.ServiceMonitorUpdate(cluster.Config, req.ServiceMonitor.Namespace, req.ServiceMonitor,req.UpdateOpts)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.Info = convertK8sServerMonitorToTopkeServiceMonitor(meta)
// 	return resp
// }

// func (k *kubernetesService) ServiceMonitorGet(ctx context.Context, req *topke.ServiceMonitorGetRequest)*topke.ServiceMonitorResponse {
// 	resp := &topke.ServiceMonitorResponse{Status: status.SuccessStatus}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	meta, err := libkubernetes.ServiceMonitorGet(cluster.Config, req.Namespace, req.Name, req.ToGetOpts())
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	resp.Info = convertK8sServerMonitorToTopkeServiceMonitor(meta)
// 	return resp
// }

// func (k *kubernetesService) ServiceMonitorList(ctx context.Context, req *topke.ServiceMonitorListRequest)*topke.ServiceMonitorListResponse {
// 	resp := &topke.ServiceMonitorListResponse{
// 		Status: status.SuccessStatus,
// 	}

// 	if req.Namespace == "" {
// 		req.Namespace = monitorNamespace
// 	}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ret, err := libkubernetes.ServiceMonitorList(cluster.Config,
// 		req.Namespace,
// 		metav1.ListOptions{FieldSelector:req.FieldSelector, LabelSelector: req.LabelSelector}, libkubernetes.DefaultTimeoutGet)
// 	if err != nil {
// 		resp.Status = err
// 		return resp
// 	}

// 	var sms []*topke.ServiceMonitorInfo
// 	for _, dm := range ret.Items {
// 		if !utils.IsContain([]string{dm.Name}, req.Fuzzy) {
// 			continue
// 		}

// 		sms = append(sms, convertK8sServerMonitorToTopkeServiceMonitor(dm))
// 	}

// 	resp.List = sms
// 	s, e := utils.PagingIndex(len(sms), req.PageNumber, req.PageSize)
// 	resp.List = resp.List[s:e]
// 	return resp
// }

// func convertK8sServerMonitorToTopkeServiceMonitor(sm *monitoringv1.ServiceMonitor) *topke.ServiceMonitorInfo {
// 	si := &topke.ServiceMonitorInfo{
// 		ServiceMonitor: sm,
// 	}
// 	return si
// }
