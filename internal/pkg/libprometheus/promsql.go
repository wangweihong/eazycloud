package libprometheus

import (
	"fmt"

	"github.com/wangweihong/eazycloud/apis/iprometheus"
	"github.com/wangweihong/gotoolbox/pkg/log"
)

var (
	ClusterLevel   = "cluster"
	NodeLevel      = "node"
	NamespaceLevel = "namespace"
	WorkloadLevel  = "workload"
	ApiserverLevel = "apiserver"
	PodLevel       = "pod"
	AlertLevel     = "alert"
	proql          map[string]map[string]string
	proql2         map[string]map[string][]string
)

func init() {
	proql2 = map[string]map[string][]string{
		AlertLevel: {
			"alertInfo": {`ALERTS{}`, `ALERTS{namespace=~"%v"}`},
		},
	}
	proql = map[string]map[string]string{
		ClusterLevel: {
			"cluster_cpu_ratio":            `sum(instance:node_cpu:ratio)/count(instance:node_cpu:ratio)`,
			"cluster_cpu_load_1":           `sum(node_load1)/count(node_load1)`,
			"cluster_cpu_load_5":           `sum(node_load5)/count(node_load5)`,
			"cluster_cpu_load_15":          `sum(node_load15)/count(node_load15)`,
			"cluster_memory_ratio":         `(sum(node_memory_MemTotal_bytes)-sum(node_memory_MemAvailable_bytes))/sum(node_memory_MemTotal_bytes)`,
			"cluster_disk_used_total":      `sum(max(node_filesystem_size_bytes{device=~"/dev/.*", device!~"/dev/loop\\d+", job="node-exporter"} - node_filesystem_avail_bytes{device=~"/dev/.*", device!~"/dev/loop\\d+", job="node-exporter"}) by (device, instance))`,
			"cluster_disk_io_rate":         `sum(rate(node_disk_io_time_seconds_total[1m]))`,
			"cluster_disk_read_rate":       `sum(rate(node_disk_read_time_seconds_total{}[1m]))`,
			"cluster_disk_write_rate":      `sum(rate(node_disk_write_time_seconds_total{}[1m]))`,
			"cluster_network_receive_rate": `sum(instance:node_network_receive_bytes:rate:sum)`,
			"cluster_network_send_rate":    `instance:node_network_transmit_bytes:rate:sum`,
			"cluster_pod_unknown":          `sum(kube_pod_status_phase{phase="Unknown"})`,
			"cluster_pod_failed":           `sum(kube_pod_status_phase{phase="Failed"})`,
			"cluster_pod_pending":          `sum(kube_pod_status_phase{phase="Pending"})`,
			"cluster_pod_succeeded":        `sum(kube_pod_status_phase{phase="Succeeded"})`,
			"cluster_pod_running":          `sum(kube_pod_status_phase{phase="Running"})`,
			"cluster_pod_count":            `sum(kube_pod_info{})`,
			"cluster_deployment_count":     `sum(kube_deployment_labels{})`,
			"cluster_statefulset_count":    `sum(kube_statefulset_labels{})`,
			"cluster_daemonset_count":      `sum(kube_daemonset_labels{})`,
			"cluster_cronjob_count":        `sum(kube_cronjob_labels{})`,
			"cluster_job_count":            `sum(kube_job_info{})`,
			"cluster_ingress_count":        `sum(kube_ingress_labels{})`,
			"cluster_service_count":        `sum(kube_service_info{})`,
			"cluster_pvc_count":            `sum(kube_persistentvolumeclaim_labels{})`,
			"cluster_secret_count":         `sum(kube_secret_info{})`,
			"cluster_configmap_count":      `sum(kube_configmap_info{})`,
			"cluster_cpu_rate":             `namespace:container_cpu_usage_seconds_total:sum_rate{}`,
			"cluster_memory_rate":          `namespace:container_memory_usage_bytes:sum{}`,
			"cluster_namespace_count":      `count(kube_namespace_labels)`,
		},
		ApiserverLevel: {
			"apiserver_request_duration_seconds_create": `sum(cluster:apiserver_request_duration_seconds:mean5m{verb="CREATE"}>0)/count(cluster:apiserver_request_duration_seconds:mean5m{verb="CREATE"}>0)`,                                     // 请求平均耗时
			"apiserver_request_duration_seconds_delete": `sum(cluster:apiserver_request_duration_seconds:mean5m{verb=~"DELETE|DELETECOLLECTION"}>0)/count(cluster:apiserver_request_duration_seconds:mean5m{verb=~"DELETE|DELETECOLLECTION"}>0)`, // 请求平均耗时
			"apiserver_request_duration_seconds_update": `sum(cluster:apiserver_request_duration_seconds:mean5m{verb=~"POST|PATCH|PUT|UPDATE"}>0)/count(cluster:apiserver_request_duration_seconds:mean5m{verb=~"POST|PATCH|PUT|UPDATE"}>0)`,     // 请求平均耗时
			"apiserver_request_duration_seconds_get":    `sum(cluster:apiserver_request_duration_seconds:mean5m{verb=~"GET|LIST"}>0)/count(cluster:apiserver_request_duration_seconds:mean5m{verb=~"GET|LIST"}>0)`,                               // 请求平均耗时
			"apiserver_request_duration_seconds":        `sum(cluster:apiserver_request_duration_seconds:mean5m>0)/count(cluster:apiserver_request_duration_seconds:mean5m>0)`,                                                                   // 请求平均耗时
			"apiserver_request_time_seconds":            `sum(code_resource:apiserver_request_total:rate5m)`,                                                                                                                                     // 每秒请求次数
		},
		NodeLevel: {
			"node_cpu_ratio":          `instance:node_cpu:ratio{instance=~"%v"}`,
			"node_cpu_load_1":         `node_load1{instance=~"%v"}`,
			"node_cpu_load_5":         `node_load5{instance=~"%v"}`,
			"node_cpu_load_15":        `node_load15{instance=~"%v"}`,
			"node_cpu_used":           `instance:node_cpu:rate:sum{instance="%v"}`,
			"node_cpu_total":          `instance:node_num_cpu:sum{instance="%v"}`,
			"node_memory_ratio":       `instance:node_memory_utilisation:ratio{instance="%v"}`,
			"node_memory_used":        `node_memory_MemTotal_bytes{instance="%v"} - node_memory_MemAvailable_bytes{instance="%v"}`,
			"node_memory_total":       `node_memory_MemTotal_bytes{instance="%v"}`,
			"node_disk_io_rate":       `sum(rate(node_disk_io_time_seconds_total{instance="friday-10-30-100-166"}[1m]))/count(rate(node_disk_io_time_seconds_total{instance="friday-10-30-100-166"}[1m]))`,
			"node_disk_ratio":         `1-max(  max by (device) (    node_filesystem_avail_bytes{job="node-exporter", instance="%v", fstype!=""}  ))/max(  max by (device) (    node_filesystem_size_bytes{job="node-exporter", instance="%v", fstype!=""}  ))`,
			"node_disk_used":          `max(  max by (device) (    node_filesystem_size_bytes{job="node-exporter", instance="%v", fstype!=""}  )) - max(  max by (device) (    node_filesystem_avail_bytes{job="node-exporter", instance="%v", fstype!=""}  ))`,
			"node_disk_total":         `max(  max by (device) (    node_filesystem_size_bytes{job="node-exporter", instance="%v", fstype!=""}  ))`,
			"node_disk_read_rate":     `rate(node_disk_read_time_seconds_total{instance="%v"}[1m])`,
			"node_disk_write_rate":    `rate(node_disk_write_time_seconds_total{instance="%v"}[1m])`,
			"node_pod_ratio":          `sum(kube_pod_info{node="%v"})/sum(kube_node_status_capacity_pods{node="%v"})`,
			"node_pod_used":           `sum(kube_pod_info{node="%v"})`,
			"node_pod_total":          `sum(kube_node_status_capacity_pods{node="%v"})`,
			"node_pod_running":        `kubelet_running_pod_count{node=~"%v"}`,
			"node_network_read_rate":  `instance:node_network_receive_bytes_excluding_lo:rate1m{instance="%v"}`,
			"node_network_write_rate": `instance:node_network_transmit_bytes_excluding_lo:rate1m{instance="%v"}`,
		},
		NamespaceLevel: {
			"namespace_pod_count":           `sum(kube_pod_info{namespace=~"%v"})`,
			"namespace_deployment_count":    `sum(kube_deployment_labels{namespace=~"%v"})`,
			"namespace_statefulset_count":   `sum(kube_statefulset_labels{namespace=~"%v"})`,
			"namespace_daemonset_count":     `sum(kube_daemonset_labels{namespace=~"%v"})`,
			"namespace_cronjob_count":       `sum(kube_cronjob_labels{namespace=~"%v"})`,
			"namespace_job_count":           `sum(kube_job_info{namespace=~"%v"})`,
			"namespace_ingress_count":       `sum(kube_ingress_labels{namespace=~"%v"})`,
			"namespace_service_count":       `sum(kube_service_info{namespace=~"%v"})`,
			"namespace_pvc_count":           `sum(kube_persistentvolumeclaim_labels{namespace="%v"})`,
			"namespace_secret_count":        `sum(kube_secret_info{namespace=~"%v"})`,
			"namespace_configmap_count":     `sum(kube_configmap_info{namespace=~"%v"})`,
			"namespace_cpu_rate":            `namespace:container_cpu_usage_seconds_total:sum_rate{namespace=~"%v"}`,
			"namespce_memory_rate":          `namespace:container_memory_usage_bytes:sum{namespace=~"%v"}`,
			"namespce_network_send_size":    `sum(rate(container_network_transmit_bytes_total{namespace=~"%v"}[1m]))`,
			"namespce_network_receive_size": `sum(rate(container_network_receive_bytes_total{namespace=~"%v"}[1m]))`,
		},
		WorkloadLevel: {
			"namespace_workload_cpu_rate":             `sum(rate(container_cpu_usage_seconds_total{namespace="%v",pod=~"%v"}[1m]))`,
			"namespace_workload_memory_size":          `sum(container_memory_working_set_bytes{ namespace="%v", pod=~"%v", container!="POD", container!="", image!=""})`,
			"namespace_workload_network_receive_size": `sum(rate(container_network_receive_bytes_total{namespace="%v",pod=~"%v"}[1m]))`,
			"namespace_workload_network_send_size":    `sum(rate(container_network_transmit_bytes_total{namespace="%v",pod=~"%v"}[1m]))`,
		},
		PodLevel: {
			"namespace_pod_cpu_rate":             `sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate{ namespace="%v",pod="%v"})`,
			"namespace_pod_memory_size":          `sum(container_memory_working_set_bytes{cluster="", namespace="%v", pod="%v", container!="POD", container!="", image!=""})`,
			"namespace_pod_network_receive_size": `sum(rate(container_network_receive_bytes_total{namespace="%v",pod="%v"}[1m]))`,
			"namespace_pod_network_send_size":    `sum(rate(container_network_transmit_bytes_total{namespace="%v",pod="%v"}[1m]))`,
		},
	}
}

func makeProql(metric string, sql *iprometheus.ProqlOption) (query string) {
	log.Debugf("level[%v],metric[%v],namespace[%v]", sql.Level, metric, sql.Namespace)

	switch sql.Level {
	case ClusterLevel:
		query = makeClusterSql(metric, sql)
	case NodeLevel:
		query = makeNodeSql(metric, sql)
	case ApiserverLevel:
		query = makeApiserverSql(metric, sql)
	case NamespaceLevel:
		query = makeNamespaceSql(metric, sql)
	case WorkloadLevel:
		query = makeWorkloadSql(metric, sql)
	case PodLevel:
		query = makePodSql(metric, sql)
	case AlertLevel:
		query = makeAlertSql(metric, sql)
	default:
		query = metric
	}
	log.Debugf("query[%v]", query)
	return query
}

func makeAlertSql(metric string, sql *iprometheus.ProqlOption) string {
	scopeMetricPair, _ := proql2[sql.Level]
	if len(scopeMetricPair) == 0 {
		log.Errorf("%s", sql.Level)
		return metric
	}

	exp, ok := scopeMetricPair[metric]
	if !ok {
		log.Errorf("%s", metric)
		return metric
	}
	return fullTemplate(exp, sql.Paras)
}

func fullTemplate(temp []string, paras []string) string {
	if len(temp) <= len(paras) {
		log.Errorf("%s", len(temp), len(paras), temp, paras)
		return "foosql"
	}
	switch len(paras) {
	case 0:
		return temp[0]
	case 1:
		return fmt.Sprintf(temp[1], paras[0])
	case 2:
		return fmt.Sprintf(temp[2], paras[0], paras[1])
	case 3:
		return fmt.Sprintf(temp[3], paras[0], paras[1], paras[2])
	default:
		return fmt.Sprintf(temp[4], paras[0], paras[1], paras[2], paras[3])
	}
}

func makeApiserverSql(metric string, sql *iprometheus.ProqlOption) string {
	scopeMetricPair, _ := proql[sql.Level]
	if len(scopeMetricPair) == 0 {
		log.Errorf("%s", sql.Level)
		return metric
	}

	exp, ok := scopeMetricPair[metric]
	if !ok {
		log.Errorf("%s", metric)
		return metric
	}

	switch metric {
	default:
		return exp
	}
}

func makeClusterSql(metric string, sql *iprometheus.ProqlOption) string {
	scopeMetricPair, _ := proql[sql.Level]
	if len(scopeMetricPair) == 0 {
		log.Errorf("%s", sql.Level)
		return metric
	}

	exp, ok := scopeMetricPair[metric]
	if !ok {
		log.Errorf("%s", metric)
		return metric
	}

	switch metric {
	default:
		return exp
	}
}

func makeNodeSql(metric string, sql *iprometheus.ProqlOption) string {
	scopeMetricPair, _ := proql[sql.Level]
	if len(scopeMetricPair) == 0 {
		log.Errorf("%s", sql.Level)
		return metric
	}

	exp, ok := scopeMetricPair[metric]
	if !ok {
		log.Errorf("%s", metric)
		return metric
	}

	switch metric {
	case "node_disk_ratio", "node_pod_ratio", "node_disk_used", "node_memory_used":
		return fmt.Sprintf(exp, sql.Node, sql.Node)
	default:
		return fmt.Sprintf(exp, sql.Node) // 前端传错了数据
	}
}

func makeWorkloadSql(metric string, sql *iprometheus.ProqlOption) string {
	scopeMetricPair, _ := proql[sql.Level]
	if len(scopeMetricPair) == 0 {
		log.Errorf("%s", sql.Level)
		return metric
	}

	exp, ok := scopeMetricPair[metric]
	if !ok {
		log.Errorf("%s", metric)
		return metric
	}

	switch metric {
	default:
		return fmt.Sprintf(exp, sql.Namespace, sql.PodName)
	}
}

func makePodSql(metric string, sql *iprometheus.ProqlOption) string {
	scopeMetricPair, _ := proql[sql.Level]
	if len(scopeMetricPair) == 0 {
		log.Errorf("%s", sql.Level)
		return metric
	}

	exp, ok := scopeMetricPair[metric]
	if !ok {
		log.Errorf("%s", metric)
		return metric
	}

	switch metric {
	default:
		return fmt.Sprintf(exp, sql.Namespace, sql.PodName)
	}
}

func makeNamespaceSql(metric string, sql *iprometheus.ProqlOption) string {
	scopeMetricPair, _ := proql[sql.Level]
	if len(scopeMetricPair) == 0 {
		log.Errorf("%s", sql.Level)
		return metric
	}

	exp, ok := scopeMetricPair[metric]
	if !ok {
		log.Errorf("%s", metric)
		return metric
	}

	switch metric {
	default:
		return fmt.Sprintf(exp, sql.Namespace)
	}
}
