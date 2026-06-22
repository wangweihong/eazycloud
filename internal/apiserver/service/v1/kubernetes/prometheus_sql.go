package kubernetes

import (
	"fmt"
	"strings"
)

type KV struct {
	Name  string //自定义名
	Value string //查询公式
}

var (
	// kubelet
	kubelet_up_count          = KV{Name: "kubelet_up_count", Value: `sum(up{job="kubelet", metrics_path="/metrics"})`}
	kubelet_running_pod_count = KV{Name: "kubelet_running_pod_count", Value: `sum(kubelet_running_pod_count{job="kubelet", metrics_path="/metrics", instance=~"%v"})`}
	kubelet_rpc_rate          = KV{Name: `kubelet_rpc_rate`, Value: `sum  (rate(rest_client_requests_total{job="kubelet", metrics_path="/metrics", instance=~"%v",code=~"2..|3..|4..|5.."}[5m]))`}
	kubelet_cpu_used          = KV{Name: "kubelet_cpu_used", Value: `rate(process_cpu_seconds_total{job="kubelet"}[5m])`}
	kubelet_memory_used       = KV{Name: "kubelet_memory_used", Value: `process_resident_memory_bytes{job="kubelet"}`}
	kubelet_goroutines        = KV{Name: "kubelet_goroutines", Value: `go_goroutines{job="kubelet"}`}

	// scheduler
	scheduler_up_count       = KV{Name: "scheduler_up_count", Value: `sum(up{job="scheduler"})`}
	scheduler_schedule_count = KV{Name: "scheduler_schedule_count", Value: `sum(rate(scheduler_scheduling_algorithm_duration_seconds_count{job="scheduler"}[5m]))`}
	scheduler_cpu_used       = KV{Name: "scheduler_cpu_used", Value: `rate(process_cpu_seconds_total{job="apiserver"}[5m])`}
	scheduler_memory_used    = KV{Name: "scheduler_memory_used", Value: `process_resident_memory_bytes{job="scheduler"}`}
	scheduler_goroutines     = KV{Name: "scheduler_goroutines", Value: `go_goroutines{job="scheduler"}`}

	// apiserver
	apiserver_cpu_used    = KV{Name: "apiserver_cpu_used", Value: `rate(process_cpu_seconds_total{job="apiserver"}[5m])`}
	apiserver_memory_used = KV{Name: "apiserver_memory_used", Value: `process_resident_memory_bytes{job="apiserver"}`}
	apiserver_goroutines  = KV{Name: "apiserver_goroutines", Value: `go_goroutines{job="apiserver"}`}

	// controller manager
	cm_up_count            = KV{Name: "cm_up_count", Value: `sum(up{job="controller-manager"})`}
	cm_work_queue_add_rate = KV{Name: "cm_work_queue_add_rate", Value: `sum(rate(workqueue_adds_total{job="controller-manager"}[5m]))`}
	cm_cpu_used            = KV{Name: "cm_cpu_used", Value: `rate(process_cpu_seconds_total{job="controller-manager"}[5m])`}
	cm_memory_used         = KV{Name: "cm_memory_used", Value: `process_resident_memory_bytes{job="controller-manager"}`}
	cm_goroutines          = KV{Name: "cm_goroutines", Value: `go_goroutines{job="controller-manager"}`}

	//etcd
	etcd_change_total           = KV{Name: "etcd_change_total", Value: "sum(etcd_server_leader_changes_seen_total)"}
	etcd_has_leader             = KV{Name: "etcd_has_leader", Value: "sum(etcd_server_has_leader)"}
	etcd_net_read_byte          = KV{Name: "etcd_net_read_byte", Value: "sum(rate(etcd_network_client_grpc_received_bytes_total[5m]))"}
	etcd_net_write_byte         = KV{Name: "etcd_net_write_byte", Value: "sum(rate(etcd_network_client_grpc_sent_bytes_total[5m]))"}
	etcd_db_size                = KV{Name: "etcd_db_size", Value: "sum(rate(etcd_mvcc_db_total_size_in_bytes[5m]))"}
	etcd_raft_proposal_applied  = KV{Name: "etcd_raft_proposal_applied", Value: "sum(rate(etcd_server_proposals_applied_total[5m]))"}
	etcd_raft_proposal_commited = KV{Name: "etcd_raft_proposal_commited", Value: "sum(rate(etcd_server_proposals_committed_total[5m]))"}
	etcd_raft_proposal_failed   = KV{Name: "etcd_raft_proposal_failed", Value: "sum(rate(etcd_server_proposals_failed_total[5m]))"}
	etcd_raft_proposal_pending  = KV{Name: "etcd_raft_proposal_pending", Value: "sum(rate(etcd_server_proposals_pending[5m]))"}

	// pod:container
	pod_container_cpu_used = KV{Name: "pod_container_cpu_used", Value: `sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate{ pod="%v", container!="POD",container!=""}) by (container)`}
	pod_container_mem_used = KV{Name: "pod_container_mem_used", Value: `sum(container_memory_working_set_bytes{pod="%v", container!="POD", container!="", image!=""})  by (container)`}
	pod_container_net_rb   = KV{Name: "pod_container_net_rb", Value: `sum(irate(container_network_receive_bytes_total{ pod="%v"}[5m])) by (pod)`}
	pod_container_net_wb   = KV{Name: "pod_container_net_wb", Value: `sum(rate(container_network_transmit_bytes_total{ pod="%v"}[5m])) by (pod)`}

	// node
	node_cpu_used              = KV{Name: "node_cpu_used", Value: `sum(instance:node_cpu:ratio{}) by (instance) * sum(instance:node_num_cpu:sum{}) by (instance)`}
	node_cpu_total             = KV{Name: "node_cpu_total", Value: `sum(instance:node_num_cpu:sum{}) by (instance)`}
	node_cpu_ratio             = KV{Name: "node_cpu_ratio", Value: `sum(instance:node_cpu:ratio{}) by (instance)`}
	node_cpu_request_used      = KV{Name: "node_cpu_request_used", Value: `sum(kube_pod_container_resource_requests_cpu_cores{}) by (node)`}
	node_cpu_limit_used        = KV{Name: "node_cpu_limit_used", Value: `sum(kube_pod_container_resource_limits_cpu_cores{}) by (node)`}
	node_memory_used           = KV{Name: "node_memory_used", Value: `sum((node_memory_MemTotal_bytes{job="node-exporter"}- node_memory_MemAvailable_bytes{job="node-exporter"})) by (instance)`}
	node_memory_total          = KV{Name: "node_memory_total", Value: `sum(node_memory_MemTotal_bytes{job="node-exporter" }) by (instance)`}
	node_memory_ratio          = KV{Name: "node_memory_ratio", Value: `sum(100 -(  node_memory_MemAvailable_bytes{job="node-exporter"}/  node_memory_MemTotal_bytes{job="node-exporter"}* 100)) by (instance)`}
	node_memory_request_used   = KV{Name: "node_memory_request_used", Value: `sum(kube_pod_container_resource_requests_memory_bytes{}) by (node)`}
	node_memory_limit_used     = KV{Name: "node_memory_limit_used", Value: `sum(kube_pod_container_resource_limits_memory_bytes{}) by (node)`}
	node_load_1                = KV{Name: "node_load_1", Value: `sum(node_load1{job="node-exporter"}) by (instance)`}
	node_load_5                = KV{Name: "node_load_5", Value: `sum(node_load5{job="node-exporter"}) by (instance)`}
	node_load_15               = KV{Name: "node_load_15", Value: `sum(node_load15{job="node-exporter"}) by (instance)`}
	node_disk_io_rate          = KV{Name: "node_disk_io_rate", Value: `sum(rate(node_disk_io_time_seconds_total{}[1m])) by (instance)/count(node_disk_io_time_seconds_total{}) by (instance)`}
	node_disk_ratio            = KV{Name: "node_disk_ratio", Value: `1-max(  max by (device,instance) (node_filesystem_avail_bytes{job="node-exporter",fstype!=""})) by (instance) / max(max by (device,instance) (node_filesystem_size_bytes{job="node-exporter",fstype!=""})) by (instance)`}
	node_disk_used             = KV{Name: "node_disk_used", Value: `max(max by (device,instance) (node_filesystem_size_bytes{job="node-exporter",fstype!=""})) by (instance)-max(max by (device,instance) (node_filesystem_avail_bytes{job="node-exporter",fstype!=""})) by (instance)`}
	node_disk_total            = KV{Name: "node_disk_total", Value: `max(  max by (device,instance) (node_filesystem_size_bytes{job="node-exporter",fstype!=""}))  by (instance)`}
	node_disk_read_rate        = KV{Name: "node_disk_read_rate", Value: `sum(rate(node_disk_read_time_seconds_total{}[1m])) by (instance)`}
	node_disk_write_rate       = KV{Name: "node_disk_write_rate", Value: `sum(rate(node_disk_write_time_seconds_total{}[1m])) by (instance)`}
	node_pod_ratio             = KV{Name: "node_pod_ratio", Value: `sum(kube_pod_info{}) by (node)/sum(kube_node_status_capacity_pods{}) by (node)`}
	node_pod_used              = KV{Name: "node_pod_used", Value: `sum(kube_pod_info{}) by (node)`}
	node_pod_total             = KV{Name: "node_pod_total", Value: `sum(kube_node_status_capacity_pods{}) by (node)`}
	node_pod_running           = KV{Name: "node_pod_running", Value: `sum(kubelet_running_pod_count{}) by (node)`}
	node_network_recerive_rate = KV{Name: "node_network_receive_rate", Value: `instance:node_network_receive_bytes:rate:sum`}
	node_network_send_rate     = KV{Name: "cluster_network_send_rate", Value: `instance:node_network_transmit_bytes:rate:sum`}

	// namespace
	namespace_cpu_used  = KV{Name: "namespace_cpu_used", Value: `sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate{}) by (namespace)`}
	namespace_mem_used  = KV{Name: "namespace_mem_used", Value: `sum(container_memory_working_set_bytes{ container!="", image!=""}) by (namespace)`}
	namespace_pod_count = KV{Name: "namespace_pod_count", Value: `sum(kube_pod_info{}) by (namespace)`}
	namespace_net_rb    = KV{Name: "namespace_net_rb", Value: `sum(irate(container_network_receive_bytes_total{namespace!=""}[5m])) by (namespace)`}
	namespace_net_wb    = KV{Name: "namespace_net_wb", Value: `sum(irate(container_network_transmit_bytes_total{namespace!=""}[5m])) by (namespace)`}

	// cluster:pod
	cluster_pod_cpu_used = KV{Name: "cluster_pod_cpu_used", Value: `sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate{container!="POD"}) by (pod,namespace,node)`}
	cluster_pod_mem_used = KV{Name: "cluster_pod_mem_used", Value: `sum(container_memory_working_set_bytes{container!="POD", container!="", image!=""}) by (pod,namespace,node)`}
	cluster_pod_net_rb   = KV{Name: "cluster_pod_net_read", Value: `sum(irate(container_network_receive_bytes_total{pod!=""}[5m])) by (pod,namespace,node)`}
	cluster_pod_net_wb   = KV{Name: "cluster_pod_net_write", Value: `sum(irate(container_network_transmit_bytes_total{pod!=""}[5m])) by (pod,namespace,node)`}

	// cluster:hardware
	cluster_cpu_ratio            = KV{Name: "cluster_cpu_ratio", Value: `sum(instance:node_cpu:ratio)/count(instance:node_cpu:ratio)`}
	cluster_cpu_load_1           = KV{Name: "cluster_cpu_load_1", Value: `sum(node_load1)/count(node_load1)`}
	cluster_cpu_load_5           = KV{Name: "cluster_cpu_load_5", Value: `sum(node_load5)/count(node_load5)`}
	cluster_cpu_load_15          = KV{Name: "cluster_cpu_load_15", Value: `sum(node_load15)/count(node_load15)`}
	cluster_memory_ratio         = KV{Name: "cluster_memory_ratio", Value: `(sum(node_memory_MemTotal_bytes)-sum(node_memory_MemAvailable_bytes))/sum(node_memory_MemTotal_bytes)`}
	cluster_disk_used_total      = KV{Name: "cluster_disk_used_total", Value: `sum(max(node_filesystem_size_bytes{device=~"/dev/.*", device!~"/dev/loop\\d+", job="node-exporter"} - node_filesystem_avail_bytes{device=~"/dev/.*", device!~"/dev/loop\\d+", job="node-exporter"}) by (device, instance))`}
	cluster_disk_io_rate         = KV{Name: "cluster_disk_io_rate", Value: `sum(rate(node_disk_io_time_seconds_total[1m]))`}
	cluster_network_receive_rate = KV{Name: "cluster_network_receive_rate", Value: `sum(instance:node_network_receive_bytes:rate:sum)`}
	cluster_network_send_rate    = KV{Name: "cluster_network_send_rate", Value: `instance:node_network_transmit_bytes:rate:sum`}
	cluster_pod_unknown          = KV{Name: "cluster_pod_unknown", Value: `sum(kube_pod_status_phase{phase="Unknown"})`}
	cluster_pod_failed           = KV{Name: "cluster_pod_failed", Value: `sum(kube_pod_status_phase{phase="Failed"})`}
	cluster_pod_pending          = KV{Name: "cluster_pod_pending", Value: `sum(kube_pod_status_phase{phase="Pending"})`}
	cluster_pod_succeeded        = KV{Name: "cluster_pod_succeeded", Value: `sum(kube_pod_status_phase{phase="Succeeded"})`}
	cluster_pod_running          = KV{Name: "cluster_pod_running", Value: `sum(kube_pod_status_phase{phase="Running"})`}

	// namespace:pod
	namespace_pod_cpu_rate             = KV{Name: "namespace_workload_cpu_rate", Value: `sum(rate(container_cpu_usage_seconds_total{}[1m])) by (namespace,pod)`}
	namespace_pod_memory_size          = KV{Name: "namespace_workload_memory_size", Value: `sum(container_memory_working_set_bytes{  container!="POD", container!="", image!=""}) by (namespace,pod)`}
	namespace_pod_network_receive_size = KV{Name: "namespace_workload_network_receive_size", Value: `sum(rate(container_network_receive_bytes_total{}[1m])) by (namespace,pod)`}
	namespace_pod_network_send_size    = KV{Name: "namespace_workload_network_send_size", Value: `sum(rate(container_network_transmit_bytes_total{}[1m])) by (namespace,pod)`}
)

var (
	schedulerQueryMetrics = []KV{
		scheduler_up_count,
		scheduler_schedule_count,
		scheduler_cpu_used,
		scheduler_memory_used,
		scheduler_goroutines,
	}
	controllerManagerQueryMetrics = []KV{
		cm_up_count,
		cm_work_queue_add_rate,
		cm_cpu_used,
		cm_memory_used,
		cm_goroutines,
	}
	apiserverQueryMetrics = []KV{
		apiserver_cpu_used,
		apiserver_memory_used,
		apiserver_goroutines,
	}

	kubeletQueryMetrics = []KV{
		kubelet_running_pod_count,
		kubelet_up_count,
		kubelet_running_pod_count,
		kubelet_memory_used,
		kubelet_goroutines,
		kubelet_rpc_rate,
	}

	podContainerMetrics = []KV{
		pod_container_cpu_used,
		pod_container_mem_used,
		pod_container_net_rb,
		pod_container_net_wb,
	}

	nodePhysicalResourceMetrics = []KV{
		node_cpu_used,
		node_cpu_total,
		node_cpu_ratio,
		node_cpu_request_used,
		node_cpu_limit_used,
		node_memory_used,
		node_memory_total,
		node_memory_ratio,
		node_memory_request_used,
		node_memory_limit_used,
		node_load_1,
		node_load_5,
		node_load_15,
		node_disk_io_rate,
		node_disk_ratio,
		node_disk_used,
		node_disk_total,
		node_disk_read_rate,
		node_disk_write_rate,
		node_pod_ratio,
		node_pod_used,
		node_pod_total,
		node_pod_running,
		node_network_recerive_rate,
		node_network_send_rate,
	}
)

// some query expression depends on var
func combineSql(para string, kvs ...KV) []KV {
	for k, v := range kvs {
		if !strings.Contains(v.Value, "%v") {
			continue
		}

		v.Value = fmt.Sprintf(v.Value, para)
		kvs[k] = v
	}

	return kvs
}
