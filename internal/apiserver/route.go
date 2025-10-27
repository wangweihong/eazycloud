package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/wangweihong/gotoolbox/pkg/errors"

	"github.com/wangweihong/eazycloud/internal/apiserver/controller/v1/application"
	"github.com/wangweihong/eazycloud/internal/apiserver/controller/v1/authentication"
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
			system.POST("/scanall/schedule/create", registryController.GarbageCollectScheduleCreate)
			system.POST("/scanall/schedule/delete", registryController.GarbageCollectScheduleDelete)
			system.POST("/scanall/schedule/update", registryController.GarbageCollectScheduleDelete)

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
