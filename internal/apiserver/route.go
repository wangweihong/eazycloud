package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/wangweihong/gotoolbox/pkg/errors"

	"github.com/wangweihong/eazycloud/internal/apiserver/controller/v1/application"
	"github.com/wangweihong/eazycloud/internal/apiserver/controller/v1/authentication"
	"github.com/wangweihong/eazycloud/internal/apiserver/controller/v1/kubernetes"
	"github.com/wangweihong/eazycloud/internal/apiserver/controller/v1/registry"
	"github.com/wangweihong/eazycloud/internal/apiserver/controller/v1/setting"
	"github.com/wangweihong/eazycloud/internal/apiserver/store"
	"github.com/wangweihong/eazycloud/internal/apiserver/store/postgresql"
	"github.com/wangweihong/eazycloud/internal/pkg/code"
	"github.com/wangweihong/eazycloud/pkg/core"
	"github.com/wangweihong/eazycloud/pkg/httpsvr/genericmiddleware"
)

func initRouter(g *gin.Engine) {
	InstallMiddleware(g)
	InstallApis(g)
}

func InstallMiddleware(g *gin.Engine) {
	g.Use(genericmiddleware.RequestID())
	g.Use(genericmiddleware.Context())
	g.Use(genericmiddleware.LoggerMiddleware())
}

func InstallApis(g *gin.Engine) *gin.Engine {
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errors.NewStatusF(code.ErrPageNotFound, "Page not found."), nil)
	})
	storeIns, _ := postgresql.GetPostgresSQLFactoryOr(nil)
	v1 := g.Group("/v1/eazycloud")
	{
		installRegistryApis(v1, storeIns)
		InstallApplicationApis(v1, storeIns)
	}

	return g
}

func installAuthApis(rg *gin.RouterGroup, storeIns store.Factory) {
	authv1 := rg.Group("/auth")
	{
		authController := authentication.NewController(storeIns)
		otp := authv1.Group("/otp")
		{
			otp.GET("qrcode", authController.OTPGenerateOrGet)
			otp.POST("validate", authController.OTPValidate)
		}
		// 修改以下路由需要同步修改iapiserver.SsoURL相关的常量
		sso := authv1.Group("/sso")
		{
			sp := sso.Group("/sp")
			{
				sp.GET("/saml/metadata", authController.SpSsoSamlInitiator)
				sp.POST("/saml/initiator", authController.SpSsoSamlInitiator)
				sp.POST("/saml/acs", authController.SpSsoSamlAcs)
				sp.POST("/saml/slo", authController.SpSsoSamlSLO)
				//oauth2
				// sp.POST("/oauth2/initiator", authController.SpSsoInitiator)
				// sp.POST("/oauth2/acs", authController.SpSsoInitiator)

			}

			idp := sso.Group("/idp")
			{
				// //saml
				idp.POST("/saml/answer", authController.IdpServeSAMLProtocolSSO)
				// sp.GET("/saml/metadata", authController.SpSsoInitiator)
				// //oauth2
				// idp.POST("/oauth2/answer", authController.SpSsoInitiator)
			}
		}
	}
}

func installRegistryApis(rg *gin.RouterGroup, storeIns store.Factory) {
	registryv1 := rg.Group("/registry")
	{
		registryController := registry.NewRegistryController(storeIns)

		registryv1.GET("/list", registryController.List)
		registryv1.GET("/get/:id", registryController.Get)
		registryv1.POST("/add", registryController.Add)
		registryv1.POST("/delete/:id", registryController.Delete)
		registryv1.POST("/update/:id", registryController.Update)
		registryv1.GET("/search", registryController.Search)

		project := registryv1.Group("/:registryID/project")
		{
			project.GET("/list", registryController.ProjectList)
			project.GET("/get", registryController.ProjectGet)
			project.GET("/summary", registryController.ProjectSummary)
			project.POST("/create", registryController.ProjectCreate)
			project.POST("/delete", registryController.ProjectDelete)
			project.POST("/update", registryController.ProjectUpdate)
			project.GET("/quota/get", registryController.ProjectQuotaGet)
			project.GET("/quota/update", registryController.ProjectQuotaUpdate)
		}

		repository := registryv1.Group("/repository")
		{
			repository.GET("/list", registryController.RepositoryList)
			repository.POST("/delete", registryController.RepositoryDelete)

			tag := repository.Group("/tag")
			{
				tag.GET("/tag/list", registryController.TagList)
				tag.GET("/tag/get", registryController.TagGet)
				tag.GET("/tag/manifests", registryController.TagManifests)
				tag.POST("/tag/scan", registryController.TagScan)
				tag.GET("/tag/scan-result", registryController.TagScanResult)
				tag.GET("/tag/scan-log", registryController.TagScanLogs)
				tag.GET("/tag/is-pullable", registryController.TagCanPull)
			}

		}

		remoteRegistry := registryv1.Group("/remote-registry")
		{
			remoteRegistry.GET("/list", registryController.RemoteRegistryList)
			remoteRegistry.GET("/get", registryController.RemoteRegistryGet)
			remoteRegistry.POST("/create", registryController.RemoteRegistryCreate)
			remoteRegistry.POST("/delete", registryController.RemoteRegistryDelete)
			remoteRegistry.POST("/update", registryController.RemoteRegistryUpdate)
			remoteRegistry.GET("/ping", registryController.RemoteRegistryPing)
		}

		replication := registryv1.Group("/replication")
		{
			replication.GET("/policy/list", registryController.ReplicationPolicyList)
			replication.POST("/policy/create", registryController.ReplicationPolicyCreate)
			replication.POST("/policy/delete", registryController.ReplicationPolicyDelete)
			replication.POST("/policy/update", registryController.ReplicationPolicyUpdate)
			replication.POST("/policy/execute", registryController.ReplicationPolicyExecute)
			replication.GET("/policy/execute-result", registryController.ReplicationPolicyExecuteResult)

			replication.GET("/adapter/names", registryController.ReplicationAdapterNames)
			replication.GET("/adapter/info", registryController.ReplicationAdapterList)

		}

		scanner := registryv1.Group("/scanner")
		{
			scanner.GET("/list", registryController.ScannerList)
			scanner.GET("/get", registryController.ScannerGet)
			scanner.POST("/create", registryController.ScannerCreate)
			scanner.POST("/delete", registryController.ScannerDelete)
			scanner.POST("/update", registryController.ScannerUpdate)
			scanner.GET("/ping", registryController.ScannerPing)
			scanner.GET("/set-default", registryController.ScannerSetDefault)
		}

		system := registryv1.Group("/system")
		{
			system.GET("/gc/schedule/list", registryController.GarbageCollectScheduleList)
			system.GET("/gc/schedule/get", registryController.GarbageCollectScheduleGet)
			system.POST("/gc/schedule/create", registryController.GarbageCollectScheduleCreate)
			system.POST("/gc/schedule/update", registryController.GarbageCollectScheduleUpdate)
			system.GET("/gc/log", registryController.GarbageCollectLog)
			system.POST("/gc/execute/manual", registryController.GarbageCollectManualExecute)

			system.POST("/scanall/execute/manual", registryController.SystemScanAllManualExecute)
			system.GET("/scanall/execute/result", registryController.SystemScanAllExecuteResult)
			system.GET("/scanall/schedule/list", registryController.SystemScanAllExecuteScheduleList)
			system.GET("/scanall/schedule/get", registryController.SystemScanAllExecuteScheduleGet)
			system.POST("/scanall/schedule/create", registryController.SystemScanAllScheduleCreate)
			system.POST("/scanall/schedule/delete", registryController.SystemScanAllScheduleDelete)
			system.POST("/scanall/schedule/update", registryController.SystemScanAllScheduleUpdate)

			system.GET("/cve-whitelist", registryController.SystemCVEWhiteList)

			system.GET("/log", registryController.SystemLog)

		}

		preheat := registryv1.Group("/preheat")
		{
			provider := preheat.Group("/provider")
			{
				provider.GET("/list", registryController.PreheatProviderList)
			}

			instance := preheat.Group("/instance")
			{
				instance.GET("/list", registryController.PreheatInstanceList)
				instance.POST("/create", registryController.PreheatInstanceCreate)
				instance.POST("/delete", registryController.PreheatInstanceDelete)
				instance.POST("/update", registryController.PreheatInstanceUpdate)
				instance.GET("/ping", registryController.PreheatInstancePing)
			}

		}

		helmChart := registryv1.Group("/helm-chart")
		{
			helmChart.GET("/list", registryController.HelmChartList)
			helmChart.POST("/create", registryController.HelmChartCreate)
			helmChart.POST("/delete", registryController.HelmChartDelete)

			version := helmChart.Group("/version")
			{
				version.GET("/list", registryController.HelmChartVersionList)
				version.GET("/get", registryController.HelmChartVersionGet)
				version.POST("/delete", registryController.HelmChartVersionDelete)
				version.POST("/package", registryController.HelmChartVersionPackage)
			}
		}

		artifact := registryv1.Group("/artifact")
		{
			artifact.GET("/list", registryController.ArtifactList)
			artifact.GET("/history", registryController.ArtifactBuildHistory)
			artifact.POST("/delete", registryController.ArtifactDelete)
			artifact.POST("/scan", registryController.ArtifactScan)
			artifact.GET("/scan-log", registryController.ArtifactScanLog)
		}
	}
}

func InstallSettingApis(rg *gin.RouterGroup, storeIns store.Factory) {
	settingv1 := rg.Group("/setting")
	{
		settingController := setting.NewController(storeIns)
		sso := settingv1.Group("/sso")
		{
			saml := sso.Group("/saml")
			{
				saml.POST("/idp/metadata/upsert", settingController.IdentityProviderSAMLMetadataUpsert)
				saml.GET("/idp/metadata/get", settingController.IdentityProviderSAMLMetadataGet)
				saml.GET("/idp/metadata/download", settingController.IdentityProviderSAMLMetadataDownload)

				saml.POST("/sp/metadata/upsert", settingController.ServiceProviderSAMLMetadataUpsert)
				saml.GET("/sp/metadata/get", settingController.ServiceProviderSAMLMetadataGet)
				saml.GET("/sp/metadata/download", settingController.ServiceProviderSAMLMetadataDownload)

			}

			ssoapp := sso.Group("/app")
			{
				ssoapp.POST("/idp/add", settingController.IdentityProviderAdd)
				ssoapp.POST("/idp/delete", settingController.IdentityProviderDelete)
				ssoapp.POST("/idp/update", settingController.IdentityProviderUpdate)
				ssoapp.GET("/idp/get", settingController.IdentityProviderGet)
				ssoapp.GET("/idp/list", settingController.IdentityProviderList)

				ssoapp.POST("/sp/add", settingController.ServiceProviderAdd)
				ssoapp.POST("/sp/delete", settingController.ServiceProviderDelete)
				ssoapp.POST("/sp/update", settingController.ServiceProviderUpdate)
				ssoapp.GET("/sp/get", settingController.ServiceProviderGet)
				ssoapp.GET("/sp/redirect_url", settingController.ServiceProviderRedirectURL)
				ssoapp.GET("/sp/list", settingController.ServiceProviderList)
			}
		}
	}
}

func InstallApplicationApis(rg *gin.RouterGroup, storeIns store.Factory) {
	appv1 := rg.Group("/application")
	{
		appController := application.NewApplicationController(storeIns)

		store := appv1.Group("/store")
		{
			store.GET("/list", appController.AppStoreList)
			store.GET("/get", appController.AppStoreGet)
			store.POST("/add", appController.AppStoreAdd)
			store.POST("/delete", appController.AppStoreDelete)
			store.POST("/update", appController.AppStoreUpdate)
			store.POST("/sync", appController.AppStoreSync)
		}

		template := appv1.Group("/template")
		{
			template.POST("/validate", appController.ApplicationTemplateValidate)
			template.POST("/add", appController.ApplicationTemplateAdd)
			template.GET("/list", appController.ApplicationTemplateList)
			template.GET("/get", appController.ApplicationTemplateGet)
			template.POST("/delete", appController.ApplicationTemplateDelete)
			template.POST("/update", appController.ApplicationTemplateUpdate)

			version := appv1.Group("/version")
			{
				version.POST("/list", appController.ApplicationTemplateVersionAdd)
				version.POST("/delete", appController.ApplicationTemplateVersionDelete)
				version.POST("/update", appController.ApplicationTemplateVersionUpdate)
			}
		}
	}
}

func InstallKubernetesApis(rg *gin.RouterGroup, storeIns store.Factory) {
	appv1 := rg.Group("/kubernetes")
	{
		k8sController := kubernetes.NewKubernetesController(storeIns)

		pod := appv1.Group("/pod")
		{
			pod.GET("/list", k8sController.PodList)
			pod.GET("/get", k8sController.PodGet)
			pod.POST("/add", k8sController.PodAdd)
			pod.POST("/delete", k8sController.PodDelete)
			pod.POST("/evict", k8sController.PodEvict)
			pod.POST("/batch_delete", k8sController.PodBatchDelete)
			pod.POST("/log/list", k8sController.PodLogList)
			pod.POST("/patch", k8sController.PodPatch)
			pod.POST("/component/pod", k8sController.GetComponentPod)
		}

		service := appv1.Group("/service")
		{
			service.GET("/list", k8sController.ServiceList)
			service.GET("/get", k8sController.ServiceGet)
			service.POST("/add", k8sController.ServiceAdd)
			service.POST("/delete", k8sController.ServiceDelete)
			service.POST("/batch_delete", k8sController.ServiceBatchDelete)
			service.POST("/update", k8sController.ServiceUpdate)
		}

		secret := appv1.Group("/secret")
		{
			secret.GET("/list", k8sController.SecretList)
			secret.GET("/get", k8sController.SecretGet)
			secret.POST("/add", k8sController.SecretAdd)
			secret.POST("/delete", k8sController.SecretDelete)
			secret.POST("/batch_delete", k8sController.SecretBatchDelete)
			secret.POST("/update", k8sController.SecretUpdate)
		}

		configmap := appv1.Group("/configmap")
		{
			configmap.GET("/list", k8sController.ConfigMapList)
			configmap.GET("/get", k8sController.ConfigMapGet)
			configmap.POST("/add", k8sController.ConfigMapAdd)
			configmap.POST("/delete", k8sController.ConfigMapDelete)
			configmap.POST("/batch_delete", k8sController.ConfigMapBatchDelete)
			configmap.POST("/update", k8sController.ConfigMapUpdate)
		}

		limitrange := appv1.Group("/limitrange")
		{
			limitrange.GET("/list", k8sController.LimitRangeList)
			limitrange.GET("/get", k8sController.LimitRangeGet)
			limitrange.POST("/add", k8sController.LimitRangeAdd)
			limitrange.POST("/delete", k8sController.LimitRangeDelete)
			limitrange.POST("/update", k8sController.LimitRangeUpdate)
		}

		resourcequota := appv1.Group("/resourcequota")
		{
			resourcequota.GET("/list", k8sController.ResourceQuotaList)
			resourcequota.GET("/get", k8sController.ResourceQuotaGet)
			resourcequota.POST("/add", k8sController.ResourceQuotaAdd)
			resourcequota.POST("/delete", k8sController.ResourceQuotaDelete)
			resourcequota.POST("/update", k8sController.ResourceQuotaUpdate)
		}

		namespace := appv1.Group("/namespace")
		{
			namespace.GET("/list", k8sController.NamespaceList)
			namespace.GET("/get", k8sController.NamespaceGet)
			namespace.GET("/gateway/get", k8sController.NamespaceGatewayGet)
			namespace.POST("/add", k8sController.NamespaceAdd)
			namespace.POST("/delete", k8sController.NamespaceDelete)
			namespace.POST("/update", k8sController.NamespaceUpdate)
		}

		hpa := appv1.Group("/hpa")
		{
			hpa.GET("/list", k8sController.HpaList)
			hpa.GET("/get", k8sController.HpaGet)
			hpa.POST("/add", k8sController.HpaAdd)
			hpa.POST("/delete", k8sController.HpaDelete)
			hpa.POST("/update", k8sController.HpaUpdate)
		}

		deployment := appv1.Group("/deployment")
		{
			deployment.GET("/list", k8sController.DeploymentList)
			deployment.GET("/get", k8sController.DeploymentGet)
			deployment.POST("/add", k8sController.DeploymentAdd)
			deployment.POST("/delete", k8sController.DeploymentDelete)
			deployment.POST("/batch_delete", k8sController.DeploymentBatchDelete)
			deployment.POST("/update", k8sController.DeploymentUpdate)
			deployment.POST("/recreate", k8sController.DeploymentRecreate)
			deployment.GET("/version/list", k8sController.DeploymentVersionList)
			deployment.POST("/version/update", k8sController.DeploymentVersionUpdate)
		}

		daemonset := appv1.Group("/daemonset")
		{
			daemonset.GET("/list", k8sController.DaemonSetList)
			daemonset.GET("/get", k8sController.DaemonSetGet)
			daemonset.POST("/add", k8sController.DaemonSetAdd)
			daemonset.POST("/delete", k8sController.DaemonSetDelete)
			daemonset.POST("/batch_delete", k8sController.DaemonSetBatchDelete)
			daemonset.POST("/update", k8sController.DaemonSetUpdate)
			daemonset.POST("/recreate", k8sController.DaemonSetRecreate)
			daemonset.GET("/version/list", k8sController.DaemonSetVersionList)
			daemonset.POST("/version/update", k8sController.DaemonSetVersionUpdate)
		}

		statefulset := appv1.Group("/statefulset")
		{
			statefulset.GET("/list", k8sController.StatefulSetList)
			statefulset.GET("/get", k8sController.StatefulSetGet)
			statefulset.POST("/add", k8sController.StatefulSetAdd)
			statefulset.POST("/delete", k8sController.StatefulSetDelete)
			statefulset.POST("/batch_delete", k8sController.StatefulSetBatchDelete)
			statefulset.POST("/update", k8sController.StatefulSetUpdate)
			statefulset.POST("/recreate", k8sController.StatefulSetRecreate)
			statefulset.GET("/version/list", k8sController.StatefulSetVersionList)
			statefulset.POST("/version/update", k8sController.StatefulSetVersionUpdate)
		}

		replicaset := appv1.Group("/replicaset")
		{
			replicaset.GET("/list", k8sController.ReplicaSetList)
			replicaset.GET("/get", k8sController.ReplicaSetGet)
			replicaset.POST("/add", k8sController.ReplicaSetAdd)
			replicaset.POST("/delete", k8sController.ReplicaSetDelete)
			replicaset.POST("/update", k8sController.ReplicaSetUpdate)
		}

		event := appv1.Group("/event")
		{
			event.GET("/list", k8sController.EventList)
			event.GET("/get", k8sController.EventGet)
		}

		job := appv1.Group("/job")
		{
			job.GET("/list", k8sController.JobList)
			job.GET("/get", k8sController.JobGet)
			job.POST("/add", k8sController.JobAdd)
			job.POST("/delete", k8sController.JobDelete)
			job.POST("/batch_delete", k8sController.JobBatchDelete)
		}

		cronjob := appv1.Group("/cronjob")
		{
			cronjob.GET("/list", k8sController.CronJobList)
			cronjob.GET("/get", k8sController.CronJobGet)
			cronjob.POST("/add", k8sController.CronJobAdd)
			cronjob.POST("/delete", k8sController.CronJobDelete)
			cronjob.POST("/update", k8sController.CronJobUpdate)
			cronjob.POST("/suspend/update", k8sController.CronJobUpdateSuspend)
			cronjob.POST("/schedule/update", k8sController.CronJobUpdateSchedule)
			cronjob.POST("/batch_delete", k8sController.CronJobBatchDelete)
			cronjob.POST("/trigger", k8sController.CronJobTrigger)
		}

		storageclass := appv1.Group("/storageclass")
		{
			storageclass.GET("/list", k8sController.StorageClassList)
			storageclass.GET("/get", k8sController.StorageClassGet)
			storageclass.POST("/add", k8sController.StorageClassAdd)
			storageclass.POST("/delete", k8sController.StorageClassDelete)
			storageclass.POST("/batch_delete", k8sController.StorageClassBatchDelete)
			storageclass.POST("/update", k8sController.StorageClassUpdate)
		}

		persistentvolume := appv1.Group("/persistentvolume")
		{
			persistentvolume.GET("/list", k8sController.PersistentVolumeList)
			persistentvolume.GET("/get", k8sController.PersistentVolumeGet)
			persistentvolume.POST("/add", k8sController.PersistentVolumeAdd)
			persistentvolume.POST("/delete", k8sController.PersistentVolumeDelete)
			persistentvolume.POST("/batch_delete", k8sController.PersistentVolumeBatchDelete)
			persistentvolume.POST("/update", k8sController.PersistentVolumeUpdate)
		}

		persistentvolumeclaim := appv1.Group("/persistentvolumeclaim")
		{
			persistentvolumeclaim.GET("/list", k8sController.PersistentVolumeClaimList)
			persistentvolumeclaim.GET("/get", k8sController.PersistentVolumeClaimGet)
			persistentvolumeclaim.POST("/add", k8sController.PersistentVolumeClaimAdd)
			persistentvolumeclaim.POST("/delete", k8sController.PersistentVolumeClaimDelete)
			persistentvolumeclaim.POST("/batch_delete", k8sController.PersistentVolumeClaimBatchDelete)
			persistentvolumeclaim.POST("/update", k8sController.PersistentVolumeClaimUpdate)
		}

		volumsnapshotclass := appv1.Group("/volumsnapshotclass")
		{
			volumsnapshotclass.GET("/list", k8sController.VolumeSnapshotClassList)
			volumsnapshotclass.GET("/get", k8sController.VolumeSnapshotClassGet)
			volumsnapshotclass.POST("/add", k8sController.VolumeSnapshotClassAdd)
			volumsnapshotclass.POST("/delete", k8sController.VolumeSnapshotClassDelete)
			volumsnapshotclass.POST("/batch_delete", k8sController.VolumeSnapshotClassBatchDelete)
			volumsnapshotclass.POST("/update", k8sController.VolumeSnapshotClassUpdate)
		}

		volumesnapshotcontent := appv1.Group("/volumesnapshotcontent")
		{
			volumesnapshotcontent.GET("/list", k8sController.VolumeSnapshotContentList)
			volumesnapshotcontent.GET("/get", k8sController.VolumeSnapshotContentGet)
			volumesnapshotcontent.POST("/add", k8sController.VolumeSnapshotContentAdd)
			volumesnapshotcontent.POST("/delete", k8sController.VolumeSnapshotContentDelete)
			volumesnapshotcontent.POST("/batch_delete", k8sController.VolumeSnapshotContentBatchDelete)
			volumesnapshotcontent.POST("/update", k8sController.VolumeSnapshotContentUpdate)
		}

		volumesnapshot := appv1.Group("/volumesnapshot")
		{
			volumesnapshot.GET("/list", k8sController.VolumeSnapshotList)
			volumesnapshot.GET("/get", k8sController.VolumeSnapshotGet)
			volumesnapshot.POST("/add", k8sController.VolumeSnapshotAdd)
			volumesnapshot.POST("/delete", k8sController.VolumeSnapshotDelete)
			volumesnapshot.POST("/batch_delete", k8sController.VolumeSnapshotBatchDelete)
			volumesnapshot.POST("/update", k8sController.VolumeSnapshotUpdate)
		}

		ingress := appv1.Group("/ingress")
		{
			ingress.GET("/list", k8sController.IngressList)
			ingress.GET("/get", k8sController.IngressGet)
			ingress.POST("/add", k8sController.IngressAdd)
			ingress.POST("/delete", k8sController.IngressDelete)
			ingress.POST("/batch_delete", k8sController.IngressBatchDelete)
			ingress.POST("/update", k8sController.IngressUpdate)
		}

		networkpolicy := appv1.Group("/networkpolicy")
		{
			networkpolicy.GET("/list", k8sController.NetworkPolicyList)
			networkpolicy.GET("/get", k8sController.NetworkPolicyGet)
			networkpolicy.POST("/add", k8sController.NetworkPolicyAdd)
			networkpolicy.POST("/delete", k8sController.NetworkPolicyDelete)
			networkpolicy.POST("/batch_delete", k8sController.NetworkPolicyBatchDelete)
			networkpolicy.POST("/update", k8sController.NetworkPolicyUpdate)
		}

	}
}
