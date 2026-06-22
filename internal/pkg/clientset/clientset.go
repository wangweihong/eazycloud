package clientset

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/wangweihong/gotoolbox/pkg/errors"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/internal/pkg/api/constant"

	"sync"
	"time"
)

var (
	informerEnabled = true
	//informerEnabled        = true //enable informer. if informer has problem, we can disable just set to false.
	defaultResyncPeriod    = 24 * time.Hour
	informerMapLock        = sync.Mutex{}
	informerControllerMaps = make(map[string]*ResourceCacheController) // record cluster and resource cache controllers
)

func RegisterInformerController(cluster *iapiserver.Cluster) error {
	if !informerEnabled { // do nothing if informer disable
		return nil
	}
	informerMapLock.Lock()
	defer informerMapLock.Unlock()

	if _, ok := informerControllerMaps[cluster.ID]; ok {
		return fmt.Errorf("cluster %v has register in informer controller map", cluster.ID)
	}
	stopCh := make(chan struct{})
	controller, err := NewResourceCacheControllerWithCluster(cluster, stopCh)
	if err != nil {
		return err
	}
	informerControllerMaps[cluster.ID] = controller

	logrus.Debugf("start to run cluster %v/%v's informer controller", cluster.ID, cluster.Name)
	go func() {
		if err := controller.Run(controller.stopChan); err != nil {
			logrus.Errorf("fail to run cluster %v/%v's informer controller %v", cluster.ID, cluster.Name, err.Error())
		}
	}()
	logrus.Debugf("finish to run cluster %v/%v's informer controller", cluster.ID, cluster.Name)
	return nil
}

func UnregisterInformerController(cluster *iapiserver.Cluster) {
	if !informerEnabled { // do nothing if informer disable
		return
	}
	informerMapLock.Lock()
	defer informerMapLock.Unlock()

	controller, ok := informerControllerMaps[cluster.ID]
	if !ok {
		//  not in informerControllerMap is not a big deal
		logrus.Debugf("cluster %v/%v does not register in informer controller map", cluster.ID, cluster.Name)
		return
	}
	delete(informerControllerMaps, cluster.ID)
	cleanClusterMonitorService(cluster.ID) // remove monitor
	cleanClusterNodePorts(cluster.ID)      // remove cluster node Port
	if controller.stopChan != nil {
		close(controller.stopChan)
	}
}

func getInformerController(cluster *iapiserver.Cluster) (*ResourceCacheController, error) {
	if !informerEnabled {
		return nil, errors.New("informer not Enabled")
	}

	if cluster.State != constant.ClusterStateRunning {
		return nil, errors.Errorf("cluster[%v] is unhealthy", cluster.ID)
	}
	informerMapLock.Lock()
	defer informerMapLock.Unlock()

	controller, ok := informerControllerMaps[cluster.ID]
	if !ok {
		return nil, errors.New("informer controller is not start")
	}
	return controller, nil
}
