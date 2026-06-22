package clientset

import (
	"github.com/wangweihong/eazycloud/apis/iapiserver"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	appsinformers "k8s.io/client-go/informers/apps/v1"
	autoscalinginformers "k8s.io/client-go/informers/autoscaling/v1"
	batchinformers "k8s.io/client-go/informers/batch/v1"
	coreinformers "k8s.io/client-go/informers/core/v1"
	networkinginformers "k8s.io/client-go/informers/networking/v1"
	policyinformers "k8s.io/client-go/informers/policy/v1"
	rbacinformers "k8s.io/client-go/informers/rbac/v1"
	storageinformers "k8s.io/client-go/informers/storage/v1"

	"k8s.io/client-go/tools/cache"

	"fmt"

	"k8s.io/client-go/informers"
)

// ResourceCacheController logs the name and namespace of pods that are added,
// deleted, or updated
// use informer policy for we can use local cache data
type ResourceCacheController struct {
	cluster         *iapiserver.Cluster
	stopChan        chan struct{}
	informerFactory informers.SharedInformerFactory

	podInformer                   coreinformers.PodInformer
	serviceInformer               coreinformers.ServiceInformer
	configMapInformer             coreinformers.ConfigMapInformer
	secretInformer                coreinformers.SecretInformer
	serviceAccountInformer        coreinformers.ServiceAccountInformer
	replicaControllerInformer     coreinformers.ReplicationControllerInformer
	namespaceInformer             coreinformers.NamespaceInformer
	limitRangeInformer            coreinformers.LimitRangeInformer
	resourceQuotaInformer         coreinformers.ResourceQuotaInformer
	persistentVolumeClaimInformer coreinformers.PersistentVolumeClaimInformer
	persistentVolumeInformer      coreinformers.PersistentVolumeInformer
	nodeInformer                  coreinformers.NodeInformer
	endpointsInformer             coreinformers.EndpointsInformer
	eventInformer                 coreinformers.EventInformer
	//batch
	jobInformer     batchinformers.JobInformer
	cronJobInformer batchinformers.CronJobInformer

	//apps
	deploymentInformer  appsinformers.DeploymentInformer
	daemonSetInformer   appsinformers.DaemonSetInformer
	statefulSetInformer appsinformers.StatefulSetInformer
	replicaSetInformer  appsinformers.ReplicaSetInformer

	//rbac
	roleInformer               rbacinformers.RoleInformer
	rolebindingInformer        rbacinformers.RoleBindingInformer
	clusterroleInformer        rbacinformers.ClusterRoleInformer
	clusterrolebindingInformer rbacinformers.ClusterRoleBindingInformer
	//storage
	storageClassInformer storageinformers.StorageClassInformer
	//application
	//applicationInformer applicationv1beta1Informers.ApplicationInformer

	//networkPolicy
	ingressInformer              networkinginformers.IngressInformer
	networkingInformer           networkinginformers.NetworkPolicyInformer
	podDisruptionBudegetInformer policyinformers.PodDisruptionBudgetInformer

	//hap
	hpaInformer autoscalinginformers.HorizontalPodAutoscalerInformer
}

// Run starts shared informers and waits for the shared informer cache to
// synchronize.
func (c *ResourceCacheController) Run(stopCh chan struct{}) error {
	// Starts all the shared informers that have been created by the factory so
	// far.
	c.informerFactory.Start(stopCh)
	//c.applicationInformerFactory.Start(stopCh)
	// wait for the initial synchronization of the local cache.
	// poll every 100 millisecond. should add timeout to stop too long wait?
	// this only reason WaitForNamedCacheSync error is call stopCh signal that will trigger timeout error.
	if !cache.WaitForNamedCacheSync("ResourceCacheController", stopCh,
		c.podInformer.Informer().HasSynced,
		c.serviceInformer.Informer().HasSynced,
		c.serviceAccountInformer.Informer().HasSynced,
		c.secretInformer.Informer().HasSynced,
		c.configMapInformer.Informer().HasSynced,
		c.replicaControllerInformer.Informer().HasSynced,
		c.namespaceInformer.Informer().HasSynced,
		c.limitRangeInformer.Informer().HasSynced,
		c.resourceQuotaInformer.Informer().HasSynced,
		c.persistentVolumeClaimInformer.Informer().HasSynced,
		c.persistentVolumeInformer.Informer().HasSynced,
		c.nodeInformer.Informer().HasSynced,
		c.endpointsInformer.Informer().HasSynced,
		c.eventInformer.Informer().HasSynced,

		//batch
		c.cronJobInformer.Informer().HasSynced,
		c.jobInformer.Informer().HasSynced,

		//apps
		c.deploymentInformer.Informer().HasSynced,
		c.daemonSetInformer.Informer().HasSynced,
		c.statefulSetInformer.Informer().HasSynced,
		c.replicaControllerInformer.Informer().HasSynced,

		// storage
		c.storageClassInformer.Informer().HasSynced,

		// rbac
		c.roleInformer.Informer().HasSynced,
		c.rolebindingInformer.Informer().HasSynced,
		c.clusterroleInformer.Informer().HasSynced,
		c.clusterrolebindingInformer.Informer().HasSynced,
		// networking
		c.networkingInformer.Informer().HasSynced,
		c.ingressInformer.Informer().HasSynced,
		//autoscaling
		c.hpaInformer.Informer().HasSynced,
	) {
		return fmt.Errorf("Failed to sync informer")
	}
	return nil
}

func NewResourceCacheControllerWithCluster(
	cluster *iapiserver.Cluster,
	stopChan chan struct{},
) (*ResourceCacheController, error) {
	rc, err := NewResourceCacheController(cluster.Config.ToRestConfig(), stopChan)
	if err != nil {
		return nil, err
	}
	rc.cluster = cluster
	return rc, nil
}

func NewResourceCacheController(restConfig *rest.Config, stopChan chan struct{}) (*ResourceCacheController, error) {
	kubernetesClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	informerFactory := informers.NewSharedInformerFactory(kubernetesClient, defaultResyncPeriod)
	//applicationInformerFactory := applicationInformers.NewSharedInformerFactory(appclientSet, defaultResyncPeriod)

	c := &ResourceCacheController{
		stopChan:                      stopChan,
		informerFactory:               informerFactory,
		podInformer:                   informerFactory.Core().V1().Pods(),
		serviceInformer:               informerFactory.Core().V1().Services(),
		configMapInformer:             informerFactory.Core().V1().ConfigMaps(),
		secretInformer:                informerFactory.Core().V1().Secrets(),
		serviceAccountInformer:        informerFactory.Core().V1().ServiceAccounts(),
		persistentVolumeClaimInformer: informerFactory.Core().V1().PersistentVolumeClaims(),
		persistentVolumeInformer:      informerFactory.Core().V1().PersistentVolumes(),
		replicaControllerInformer:     informerFactory.Core().V1().ReplicationControllers(),
		namespaceInformer:             informerFactory.Core().V1().Namespaces(),
		limitRangeInformer:            informerFactory.Core().V1().LimitRanges(),
		resourceQuotaInformer:         informerFactory.Core().V1().ResourceQuotas(),
		nodeInformer:                  informerFactory.Core().V1().Nodes(),
		endpointsInformer:             informerFactory.Core().V1().Endpoints(),
		eventInformer:                 informerFactory.Core().V1().Events(),
		//batch
		cronJobInformer: informerFactory.Batch().V1().CronJobs(),
		jobInformer:     informerFactory.Batch().V1().Jobs(),
		//apps
		deploymentInformer:  informerFactory.Apps().V1().Deployments(),
		daemonSetInformer:   informerFactory.Apps().V1().DaemonSets(),
		statefulSetInformer: informerFactory.Apps().V1().StatefulSets(),
		replicaSetInformer:  informerFactory.Apps().V1().ReplicaSets(),
		//extension
		//storage
		storageClassInformer: informerFactory.Storage().V1().StorageClasses(),
		//rbac
		roleInformer:               informerFactory.Rbac().V1().Roles(),
		rolebindingInformer:        informerFactory.Rbac().V1().RoleBindings(),
		clusterroleInformer:        informerFactory.Rbac().V1().ClusterRoles(),
		clusterrolebindingInformer: informerFactory.Rbac().V1().ClusterRoleBindings(),
		//network
		networkingInformer: informerFactory.Networking().V1().NetworkPolicies(),
		ingressInformer:    informerFactory.Networking().V1().Ingresses(),
		// policy
		podDisruptionBudegetInformer: informerFactory.Policy().V1().PodDisruptionBudgets(),
		// hpa
		hpaInformer: informerFactory.Autoscaling().V1().HorizontalPodAutoscalers(),
	}

	c.podInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	c.serviceInformer.Informer().AddEventHandler(c.ServiceEvent())
	c.serviceAccountInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	c.configMapInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	c.secretInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	c.persistentVolumeClaimInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	c.persistentVolumeInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	c.replicaControllerInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	c.namespaceInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	c.limitRangeInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	c.resourceQuotaInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	c.nodeInformer.Informer().AddEventHandler(c.NodeEvent())
	c.endpointsInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	c.eventInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	//batch
	c.cronJobInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	c.jobInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	//apps
	c.deploymentInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	c.daemonSetInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	c.statefulSetInformer.Informer().AddEventHandler(c.MonitorPrometheusEvent())
	c.replicaSetInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	//extension
	//storage
	c.storageClassInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	//rbac
	c.roleInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	c.rolebindingInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	c.clusterroleInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	c.clusterrolebindingInformer.Informer().AddEventHandler(c.GeneralEventHandler())

	//networking
	c.networkingInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	c.ingressInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	// policy
	c.podDisruptionBudegetInformer.Informer().AddEventHandler(c.GeneralEventHandler())
	//set watch handler when list-watch fail?
	// autoscaling
	c.hpaInformer.Informer().AddEventHandler((c.GeneralEventHandler()))

	return c, nil
}
