package clientset

import (
	"reflect"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/wangweihong/eazycloud/apis/iapiserver"

	v1 "k8s.io/api/core/v1"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/tools/cache"
)

func (c *ResourceCacheController) resourceAdd(obj any) {
	if r, ok := obj.(metav1.Object); ok {
		logrus.Debugf("Resource %s CREATED: %s/%s\n", reflect.TypeOf(obj), r.GetNamespace(), r.GetName())
	}
}

func (c *ResourceCacheController) resourceUpdate(old, new any) {
	if r, ok := old.(metav1.Object); ok {
		logrus.Debugf("Resource %s UPDATED: %s/%s\n", reflect.TypeOf(old), r.GetNamespace(), r.GetName())
	}

}

func (c *ResourceCacheController) resourceDelete(obj any) {
	if r, ok := obj.(metav1.Object); ok {
		logrus.Debugf("Resource %s DELETEd: %s/%s\n", reflect.TypeOf(obj), r.GetNamespace(), r.GetName())
	}
}

func (c *ResourceCacheController) GeneralEventHandler() cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		// Called on creationf
		AddFunc: c.resourceAdd,
		// Called on resource update and every resyncPeriod on existing resources.
		UpdateFunc: c.resourceUpdate,
		// Called on resource deletion.
		DeleteFunc: c.resourceDelete,
	}
}

type NodePort struct {
	Service *v1.Service
	Port    int32
}

var (
	monitorServiceReadyLock sync.RWMutex
	monitorServiceReadyMap  = make(map[string]bool) // prometheus-k8s statefulset至少一个副本正常.
	prometheusComponentName = "prometheus-k8s"

	nodePortLock  sync.RWMutex
	nodePortLMaps = make(map[string]map[int32]NodePort)
)

func IsMonitorServiceReady(clusterID string) bool {
	monitorServiceReadyLock.RLock()
	defer monitorServiceReadyLock.RUnlock()

	ready, _ := monitorServiceReadyMap[clusterID]
	return ready
}

func GetAllNodePortServices() map[string]map[int32]NodePort {
	nodePortLock.RLock()
	defer nodePortLock.Unlock()
	return nodePortLMaps
}

func GetClusterNodePortServices(clusterID string) map[int32]NodePort {
	nodePortLock.RLock()
	defer nodePortLock.RUnlock()
	nps := make(map[int32]NodePort)

	nodePorts, ok := nodePortLMaps[clusterID]
	if !ok {
		return nps
	}

	for k, v := range nodePorts {
		nps[k] = v
	}
	return nps
}

func setNodePortService(clusterID string, port int32, service *v1.Service) {
	nodePortLock.Lock()
	defer nodePortLock.Unlock()
	nodePorts, ok := nodePortLMaps[clusterID]
	if !ok {
		nodePorts = make(map[int32]NodePort)
	}

	nodePorts[port] = NodePort{
		Service: service,
		Port:    port,
	}
	nodePortLMaps[clusterID] = nodePorts
}

func cleanNodePortService(clusterID string, port int32) {
	nodePortLock.Lock()
	defer nodePortLock.Unlock()
	nodePorts, ok := nodePortLMaps[clusterID]
	if !ok {
		nodePorts = make(map[int32]NodePort)
		return
	}

	delete(nodePorts, port)
	nodePortLMaps[clusterID] = nodePorts
}

func cleanClusterNodePorts(clusterID string) {
	nodePortLock.Lock()
	defer nodePortLock.Unlock()
	delete(monitorServiceReadyMap, clusterID)
}

func setMonitorServiceReady(clusterID string, ready bool) {
	monitorServiceReadyLock.Lock()
	defer monitorServiceReadyLock.Unlock()
	monitorServiceReadyMap[clusterID] = ready
}

func cleanClusterMonitorService(clusterID string) {
	monitorServiceReadyLock.Lock()
	defer monitorServiceReadyLock.Unlock()
	delete(monitorServiceReadyMap, clusterID)
}

func (c *ResourceCacheController) statefulSetAdd(obj any) {
	if r, ok := obj.(*appsv1.StatefulSet); ok {
		logrus.Debugf("StatefulSet %s CREATED: %s/%s", reflect.TypeOf(obj), r.GetNamespace(), r.GetName())
		if r.Name == prometheusComponentName {
			setMonitorServiceReady(c.cluster.ID, r.Status.ReadyReplicas > 0)
		}
	}
}

func (c *ResourceCacheController) statefulSetUpdate(old, new any) {
	if r, ok := new.(*appsv1.StatefulSet); ok {
		logrus.Debugf("StatefulSet %s UPDATED: %s/%s", reflect.TypeOf(old), r.GetNamespace(), r.GetName())
		if r.Name == prometheusComponentName {
			//FIXME: should we care about if prometheus-k8s service exist?
			setMonitorServiceReady(c.cluster.ID, r.Status.ReadyReplicas > 0)
		}
	}
}

func (c *ResourceCacheController) statefulSetDelete(obj any) {
	if r, ok := obj.(*appsv1.StatefulSet); ok {
		logrus.Debugf("StatefulSet %s DELETEd: %s/%s", reflect.TypeOf(obj), r.GetNamespace(), r.GetName())
		if r.Name == prometheusComponentName {
			setMonitorServiceReady(c.cluster.ID, false)
		}
	}
}

func (c *ResourceCacheController) NodeEvent() cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		// Called on creation
		AddFunc: c.nodeAdd,
		// Called on resource update and every resyncPeriod on existing resources.
		UpdateFunc: c.nodeUpdate,
		// Called on resource deletion.
		DeleteFunc: c.nodeDelete,
	}
}

func (c *ResourceCacheController) MonitorPrometheusEvent() cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		// Called on creation
		AddFunc: c.statefulSetAdd,
		// Called on resource update and every resyncPeriod on existing resources.
		UpdateFunc: c.statefulSetUpdate,
		// Called on resource deletion.
		DeleteFunc: c.statefulSetDelete,
	}
}

func (c *ResourceCacheController) ServiceEvent() cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		// Called on creation
		AddFunc: c.serviceAdd,
		// Called on resource update and every resyncPeriod on existing resources.
		UpdateFunc: c.serviceUpdate,
		// Called on resource deletion.
		DeleteFunc: c.serviceDelete,
	}
}

func (c *ResourceCacheController) serviceAdd(obj any) {
	if r, ok := obj.(*v1.Service); ok {
		logrus.Debugf("Service %s CREATED: %s/%s", reflect.TypeOf(obj), r.GetNamespace(), r.GetName())
		if r.Spec.Type == v1.ServiceTypeNodePort {
			for _, v := range r.Spec.Ports {
				if v.NodePort != 0 {
					setNodePortService(c.cluster.ID, v.NodePort, r)
				}
			}
		}
	}
}

func (c *ResourceCacheController) serviceUpdate(old, new any) {
	if r, ok := old.(*v1.Service); ok {
		if r.Spec.Type == v1.ServiceTypeNodePort {
			for _, v := range r.Spec.Ports {
				if v.NodePort != 0 {
					cleanNodePortService(c.cluster.ID, v.NodePort)
				}
			}
		}
	}

	if r, ok := new.(*v1.Service); ok {
		logrus.Debugf("Service %s UPDATED: %s/%s", reflect.TypeOf(new), r.GetNamespace(), r.GetName())
		if r.Spec.Type == v1.ServiceTypeNodePort {
			for _, v := range r.Spec.Ports {
				if v.NodePort != 0 {
					setNodePortService(c.cluster.ID, v.NodePort, r)
				}
			}
		}
	}
}

func (c *ResourceCacheController) serviceDelete(obj any) {
	if r, ok := obj.(*v1.Service); ok {
		logrus.Debugf("Service %s DELETEd: %s/%s", reflect.TypeOf(obj), r.GetNamespace(), r.GetName())
		if r.Spec.Type == v1.ServiceTypeNodePort {
			for _, v := range r.Spec.Ports {
				if v.NodePort != 0 {
					cleanNodePortService(c.cluster.ID, v.NodePort)
				}
			}
		}
	}
}

// only care about volume add/update, don't care node delete
type LabelChanged struct {
	Labels  map[string]string
	Cluster *iapiserver.Cluster
}

var (
	NodeLabelUpdateChan = make(chan LabelChanged, 50)
)

func (c *ResourceCacheController) nodeAdd(obj any) {
	if r, ok := obj.(*v1.Node); ok {
		logrus.Debugf("Node %s CREATED: %s/%s", reflect.TypeOf(obj), r.GetNamespace(), r.GetName())
		go func() {
			NodeLabelUpdateChan <- LabelChanged{
				Labels:  r.Labels,
				Cluster: c.cluster,
			}
		}()
	}
}

func (c *ResourceCacheController) nodeUpdate(old, new any) {
	if r, ok := new.(*v1.Node); ok {
		logrus.Debugf("Node %s UPDATED: %s/%s", reflect.TypeOf(new), r.GetNamespace(), r.GetName())
		go func() {
			NodeLabelUpdateChan <- LabelChanged{
				Labels:  r.Labels,
				Cluster: c.cluster,
			}
		}()
	}
}

func (c *ResourceCacheController) nodeDelete(obj any) {
	if r, ok := obj.(*v1.Node); ok {
		logrus.Debugf("Node %s DELETE: %s/%s\n", reflect.TypeOf(obj), r.GetNamespace(), r.GetName())
	}
}
