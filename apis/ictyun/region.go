package ictyun

type RegionListRequest struct {
	RegionName *string
}

type RegionListEntry struct {
	RegionParent string `json:"regionParent"`
	RegionID     string `json:"regionID"`
	RegionType   string `json:"regionType"`
	RegionName   string `json:"regionName"`
	// 是否多区域资源池
	IsMultiZones bool `json:"isMultiZones"`
	// 只有多区域资源池才会返回
	ZoneList []string `json:"zoneList"`
}

type RegionListResponse struct {
	ReturnObj struct {
		RegionList []RegionListEntry
	} `json:"returnObj"`
}

type RegionInfoRequest struct {
	RegionName *string
}

type RegionInfoEntry struct {
	Name string
	ID   string
	UUID string
}

type RegionInfoResponse struct {
	ReturnObj struct {
		RegionList []RegionInfoEntry `json:"regions"`
	} `json:"returnObj"`
}

type RegionProductRequest struct {
	RegionID string
}

type RegionProductResponse struct {
	ReturnObj struct {
		Ch struct {
			HasCh              string `json:"hasCh"`
			HasChMonitor       string `json:"hasChMonitor"`
			HasChReconsitution string `json:"hasChReconsitution"`
		} `json:"ch"`
		Cda struct {
			HasCda           string `json:"hasCda"`
			HasCdaMonitor    string `json:"hasCdaMonitor"`
			HasCdaPermission string `json:"hasCdaPermission"`
			HasCda2          string `json:"hasCda2"`
		} `json:"cda"`
		Oss struct {
			HasOssFc          string `json:"hasOssFc"`
			HasOssVpc         string `json:"hasOssVpc"`
			HasOssIpv6        string `json:"hasOssIpv6"`
			HasOss            string `json:"hasOss"`
			HasOssIa          string `json:"hasOssIa"`
			HasOssServerOpen  string `json:"hasOssServerOpen"`
			HasOssResource    string `json:"hasOssResource"`
			HasOSSWebSite     string `json:"hasOSSWebSite"`
			HasMakeImageByOss string `json:"hasMakeImageByOss"`
		} `json:"oss"`
		Ebs struct {
			HasVolumeSpeedLimit string `json:"hasVolumeSpeedLimit"`
			StorageType         []struct {
				Type string `json:"type"`
				Name string `json:"name"`
			} `json:"storageType"`
			HasBks           string `json:"hasBks"`
			HasEbsFromBks    string `json:"hasEbsFromBks"`
			ShareEbs         string `json:"shareEbs"`
			HasMakeDiskImage string `json:"hasMakeDiskImage"`
			HasVolumePool    string `json:"HasVolumePool"`
			UpgradeSysVolume string `json:"upgradeSysVolume"`
		} `json:"ebs"`
		Lb struct {
			HasElb                         string `json:"hasElb"`
			HasListenerOnHTTPS             string `json:"hasListenerOnHttps"`
			HasLBSupportSession            string `json:"hasLBSupportSession"`
			HasElbCert                     string `json:"hasElbCert"`
			HasListenerOnHTTP              string `json:"hasListenerOnHttp"`
			HasUDPLB                       string `json:"hasUdpLB"`
			HasElbMonitor                  string `json:"hasElbMonitor"`
			HasElbSubnetIpv6               string `json:"hasElbSubnetIpv6"`
			HasLBURLTransmit               string `json:"hasLBUrlTransmit"`
			SubnetOfferingIDWithInternalLb string `json:"subnetOfferingIdWithInternalLb"`
			HasLBWhiteBlack                string `json:"hasLBWhiteBlack"`
		} `json:"lb"`
		Monitor struct {
			HasMonitor   string `json:"hasMonitor"`
			HasFGMonitor string `json:"hasFGMonitor"`
			HasAuditLog  string `json:"hasAuditLog"`
		} `json:"monitor"`
		AzList []interface{} `json:"azList"`
		Efs    struct {
			HasEfsPay        string `json:"hasEfsPay"`
			HasEfsPermission string `json:"hasEfsPermission"`
			HasEfsKms        string `json:"hasEfsKms"`
		} `json:"efs"`
		Saas struct {
			MiddlewareSaasServicestageURL string `json:"middlewareSaasServicestageUrl"`
			MiddlewareSaasServicestage    string `json:"middlewareSaasServicestage"`
		} `json:"saas"`
		SecuritySafe struct {
			SecuritySafeURL string `json:"securitySafeUrl"`
			SecuritySafe    string `json:"securitySafe"`
		} `json:"securitySafe"`
		Other struct {
			SecurityWaf                  string `json:"securityWaf"`
			HasAccountAuthVPC            string `json:"hasAccountAuthVPC"`
			HasChCreate                  string `json:"hasChCreate"`
			SecurityBgp                  string `json:"securityBgp"`
			HasRoutingTable              string `json:"hasRoutingTable"`
			HasEcmBatchReinstallByRAM    string `json:"hasEcmBatchReinstallByRam"`
			IsHT                         string `json:"isHT"`
			HasLbSession                 string `json:"hasLbSession"`
			HasHPFSPOSIX                 string `json:"hasHPFSPOSIX"`
			HasHPFSPerf                  string `json:"hasHPFSPerf"`
			RegionCode                   string `json:"regionCode"`
			HasCreateLiteEcmByApplyIms   string `json:"hasCreateLiteEcmByApplyIms"`
			HasSecRuleCopy               string `json:"hasSecRuleCopy"`
			NoVmblockDevice              string `json:"NoVmblockDevice"`
			ResourceSoldout              string `json:"resourceSoldout"`
			SecurityWafURL               string `json:"securityWafUrl"`
			HasImsIntegrityCheck         string `json:"hasImsIntegrityCheck"`
			HasCtwt                      string `json:"hasCtwt"`
			QosPolicyHeader              string `json:"qosPolicyHeader"`
			SecurityWebfirelwall         string `json:"securityWebfirelwall"`
			HasSdwanAccount              string `json:"hasSdwanAccount"`
			SecurityBgpURL               string `json:"securityBgpUrl"`
			HasRemote                    string `json:"hasRemote"`
			Gpu                          string `json:"GPU"`
			HasSdwanSafetyTube           string `json:"hasSdwanSafetyTube"`
			CNPHtList                    string `json:"CNPHtList"`
			HasIAAS                      string `json:"hasIAAS"`
			HasSnapshotPolicy            string `json:"hasSnapshotPolicy"`
			HasHPFSSoldOut               string `json:"hasHPFSSoldOut"`
			SecurityDNS                  string `json:"securityDns"`
			CtyunID                      string `json:"ctyunId"`
			HostGroupAdjustable          string `json:"hostGroupAdjustable"`
			SecurityCmsURL               string `json:"securityCmsUrl"`
			HasCdaClientRouterBfd        string `json:"hasCdaClientRouterBfd"`
			SecuritySslvpnURL            string `json:"securitySslvpnUrl"`
			SecurityCtwaf                string `json:"securityCtwaf"`
			HasEcsMonitorWithDiskMertic  string `json:"hasEcsMonitorWithDiskMertic"`
			HbaseAPIURL                  string `json:"hbaseApiUrl"`
			HasCrossZoneTabs             string `json:"hasCrossZoneTabs"`
			Lng                          string `json:"lng"`
			HasSdwanMac                  string `json:"hasSdwanMac"`
			RegionName                   string `json:"regionName"`
			HasIPBasicBillByBandWidthEIP string `json:"hasIPBasicBillByBandWidthEIP"`
			Cutover                      string `json:"cutover"`
			HasCms                       string `json:"hasCms"`
			State                        int    `json:"state"`
			HasKeyOperationValid         string `json:"hasKeyOperationValid"`
			SecurityCtwt                 string `json:"securityCtwt"`
			SecurityCloudfirewallURL     string `json:"securityCloudfirewallUrl"`
			BksExcludeXssd               string `json:"bksExcludeXssd"`
			HasLITE1Band                 string `json:"hasLITE1Band"`
			SecurityCloudParseURL        string `json:"securityCloudParseUrl"`
			PaaSFlavorVersions           string `json:"PaaSFlavorVersions"`
			SecurityWebfirelwallURL      string `json:"securityWebfirelwallUrl"`
			ZabbixPwd                    string `json:"zabbixPwd"`
			ComputedServerless           string `json:"computedServerless"`
			AzName                       []struct {
				AvailabilityZone string `json:"availabilityZone"`
				Name             string `json:"name"`
			} `json:"azName"`
			ComputedCloudDesktop        string      `json:"computedCloudDesktop"`
			HasUpdateSubnet             string      `json:"hasUpdateSubnet"`
			HasGpuQuantityPay           string      `json:"hasGpuQuantityPay"`
			ZabbixURL                   string      `json:"zabbixUrl"`
			HasCertificateEdit          string      `json:"hasCertificateEdit"`
			SecurityNgfwURL             string      `json:"securityNgfwUrl"`
			Region                      string      `json:"region"`
			HasBandwidthTabs            string      `json:"hasBandwidthTabs"`
			CdrExcludeXssd              string      `json:"cdrExcludeXssd"`
			HasElbTwoWay                string      `json:"hasElbTwoWay"`
			DataPassLubanBigdata        string      `json:"dataPassLubanBigdata"`
			SecurityKmsURL              string      `json:"securityKmsUrl"`
			HasDateTimePickLog          string      `json:"hasDateTimePickLog"`
			HasEfsIpv6                  string      `json:"hasEfsIpv6"`
			ComputedCci                 string      `json:"computedCci"`
			HasIPSpeedLimit             string      `json:"hasIpSpeedLimit"`
			HasHPFS                     string      `json:"hasHPFS"`
			SecuritySitesafe            string      `json:"securitySitesafe"`
			HasEcmBaseSG                string      `json:"hasEcmBaseSG"`
			HasSdwanEdgeUpdateVersion   string      `json:"hasSdwanEdgeUpdateVersion"`
			SecurityBgpEdge             string      `json:"securityBgpEdge"`
			SecurityCloudParse          string      `json:"securityCloudParse"`
			HasNewUpgradeRules          string      `json:"hasNewUpgradeRules"`
			HasLITE1                    string      `json:"hasLITE1"`
			SecurityCtwtURL             string      `json:"securityCtwtUrl"`
			SecurityCloudfirewall       string      `json:"securityCloudfirewall"`
			SecuritySitesafemonitorURL  string      `json:"securitySitesafemonitorUrl"`
			HasNatMonitor               string      `json:"hasNatMonitor"`
			NasSawanForClound           string      `json:"nasSawanForClound"`
			OsProject                   string      `json:"osProject"`
			CaasShareNetworkID          string      `json:"caasShareNetworkId"`
			SecurityBgpEdgeURL          string      `json:"securityBgpEdgeUrl"`
			HasSdwanlpv6                string      `json:"hasSdwanlpv6"`
			SecurityCms                 string      `json:"securityCms"`
			HasACLRulesImport           string      `json:"hasAclRulesImport"`
			HasCreateSdwanExample       string      `json:"hasCreateSdwanExample"`
			HasIstack                   string      `json:"hasIstack"`
			ServiceURL                  string      `json:"serviceUrl"`
			ZabbixUser                  string      `json:"zabbixUser"`
			ComputedCloudDesktopURL     string      `json:"computedCloudDesktopUrl"`
			DataPaasLubanBigdataURL     string      `json:"dataPaasLubanBigdataUrl"`
			HasSdwanQlinksSafeLog       string      `json:"hasSdwanQlinksSafeLog"`
			Project                     string      `json:"project"`
			HasSMSAlarm                 string      `json:"hasSMSAlarm"`
			HasLiteEcmEbsCreateSnapshot string      `json:"hasLiteEcmEbsCreateSnapshot"`
			SecuritySslvpn              string      `json:"securitySslvpn"`
			SecurityMif                 string      `json:"securityMif"`
			HasSecurityRulesImport      string      `json:"hasSecurityRulesImport"`
			HasEfsV6Tip                 string      `json:"hasEfsV6Tip"`
			DataPassLubanBigdataURL     string      `json:"dataPassLubanBigdataUrl"`
			SecurityKms                 string      `json:"securityKms"`
			HasEIPBindRestrictByVCPU    string      `json:"hasEIPBindRestrictByVCPU"`
			Liteecm                     string      `json:"liteecm"`
			Az                          string      `json:"az"`
			SecurityCloudaudit          string      `json:"securityCloudaudit"`
			HasSecAddRules              string      `json:"hasSecAddRules"`
			CustomerPlatformAttributes  interface{} `json:"customerPlatformAttributes"`
			HasSdwanCross               string      `json:"hasSdwanCross"`
			HasCertificateCa            string      `json:"hasCertificateCa"`
			SecurityLogsauditv2URL      string      `json:"securityLogsauditv2Url"`
			HTResourceType              string      `json:"HTResourceType"`
			ZoneID                      string      `json:"zoneId"`
			HasCreateImsByActive        string      `json:"hasCreateImsByActive"`
			HasTurnCycle                string      `json:"hasTurnCycle"`
			SecurityCtwafURL            string      `json:"securityCtwafUrl"`
			Online                      string      `json:"online"`
			HasCdaVirtualGatewayMonitor string      `json:"hasCdaVirtualGatewayMonitor"`
			SecurityNgfw                string      `json:"securityNgfw"`
			SecuritySitesafemonitor     string      `json:"securitySitesafemonitor"`
			HasEfsMonitor               string      `json:"hasEfsMonitor"`
			SecurityLogsauditv2         string      `json:"securityLogsauditv2"`
			DataPaasLubanBigdata        string      `json:"dataPaasLubanBigdata"`
			HasLiteEcmEbsUnsubscribe    string      `json:"hasLiteEcmEbsUnsubscribe"`
			SecurityMifURL              string      `json:"securityMifUrl"`
			SecurityDNSURL              string      `json:"securityDnsUrl"`
			LbVersion                   string      `json:"lbVersion"`
			SecurityCloudauditURL       string      `json:"securityCloudauditUrl"`
			HasSdwanCreateEdgeNat       string      `json:"hasSdwanCreateEdgeNat"`
			Lat                         string      `json:"lat"`
			HasSdwanMonitor             string      `json:"hasSdwanMonitor"`
			HasSdwanQlinksSafeFw        string      `json:"hasSdwanQlinksSafeFw"`
			HasBatchReinstall           string      `json:"hasBatchReinstall"`
			Paas                        struct {
			} `json:"paas"`
			NoVMSysStart                string `json:"NoVmSysStart"`
			HasSdwanCreateCloundDesktop string `json:"hasSdwanCreateCloundDesktop"`
			HasDiskBackupMonitor        string `json:"hasDiskBackupMonitor"`
			MonitorPT                   string `json:"monitorPT"`
		} `json:"other"`
		Kms struct {
			HasKms       string `json:"hasKms"`
			HasKmsServer string `json:"hasKmsServer"`
		} `json:"kms"`
		Image struct {
			SafeImage             string `json:"safeImage"`
			HasPrivateImageExport string `json:"hasPrivateImageExport"`
			ShareImage            string `json:"shareImage"`
			SelectionImage        string `json:"selectionImage"`
		} `json:"image"`
		Pm struct {
			PmVnc           string `json:"PmVnc"`
			HasPmGpu        string `json:"hasPmGpu"`
			PmVncDomainName string `json:"pmVncDomainName"`
			HasPmChangeVpc  string `json:"hasPmChangeVpc"`
			HasPMMultiNic   string `json:"hasPMMultiNic"`
			Pm              string `json:"pm"`
		} `json:"pm"`
		Eip struct {
			HasEipMonitor string `json:"hasEipMonitor"`
		} `json:"eip"`
		Sdwan struct {
			HasSdwanCreateEdge  string `json:"hasSdwanCreateEdge"`
			HasSdwanIpv6        string `json:"hasSdwanIpv6"`
			HasSdwanSingleArm   string `json:"hasSdwanSingleArm"`
			HasSdwan            string `json:"hasSdwan"`
			HasSdwanOspf        string `json:"hasSdwanOspf"`
			HasSdwanACLEdit     string `json:"hasSdwanAclEdit"`
			HasSdwanQuickOffice string `json:"hasSdwanQuickOffice"`
			HasSdwanApply       string `json:"hasSdwanApply"`
			HasQos              string `json:"hasQos"`
			HasSdwanEdgeApp     string `json:"hasSdwanEdgeApp"`
		} `json:"sdwan"`
		Sfs struct {
			FileSystem string `json:"fileSystem"`
		} `json:"sfs"`
		Scaling struct {
			HasEss string `json:"hasEss"`
		} `json:"scaling"`
		Ecs struct {
			HasS6              string `json:"hasS6"`
			HasChangeV6Network string `json:"hasChangeV6Network"`
			FlavorTypes        struct {
				P []string `json:"p"`
				S []string `json:"s"`
				M []string `json:"m"`
				C []string `json:"c"`
				G []string `json:"g"`
			} `json:"flavorTypes"`
			HasS3            string `json:"hasS3"`
			HasC6M6          string `json:"hasC6M6"`
			HasChangeNetwork string `json:"hasChangeNetwork"`
			HasCloneVM       string `json:"hasCloneVM"`
			HasSnapshot      string `json:"hasSnapshot"`
			HasGroup         string `json:"hasGroup"`
			HasM3            string `json:"hasM3"`
			HasAddNic        string `json:"hasAddNic"`
			HasVnc           string `json:"hasVnc"`
		} `json:"ecs"`
		Vpc struct {
			ShareNetworkID   string `json:"shareNetworkId"`
			SubnetPoolID     string `json:"subnetPoolId"`
			IsCnp            string `json:"isCnp"`
			DNS2             string `json:"dns2"`
			DNS1             string `json:"dns1"`
			HasNat           string `json:"hasNat"`
			HasVPCRouter     string `json:"hasVPCRouter"`
			HasVirtualIP     string `json:"hasVirtualIp"`
			HasVPN           string `json:"hasVPN"`
			Ipv6             string `json:"ipv6"`
			HasP2P           string `json:"hasP2P"`
			HasSharedBw      string `json:"hasSharedBw"`
			HasMultiCidr4VPC string `json:"hasMultiCidr4VPC"`
			VpcOfferingID    string `json:"vpcOfferingId"`
		} `json:"vpc"`
		Caas struct {
			CaasAPIURL     string `json:"caasApiUrl"`
			ComputedCce    string `json:"computedCce"`
			ComputedCcr    string `json:"computedCcr"`
			ComputedCceURL string `json:"computedCceUrl"`
		} `json:"caas"`
		Paas struct {
			DbPaasDrdsURL           string `json:"dbPaasDrdsUrl"`
			DbPaasHbase             string `json:"dbPaasHbase"`
			DbPaasMongodbURL        string `json:"dbPaasMongodbUrl"`
			MiddlewarePaasMq        string `json:"middlewarePaasMq"`
			MiddlewarePaasKafka     string `json:"middlewarePaasKafka"`
			PaasCidr                string `json:"paasCidr"`
			DbPaasRedis             string `json:"dbPaasRedis"`
			DbPaasTsdb              string `json:"dbPaasTsdb"`
			DbPaasRedisURL          string `json:"dbPaasRedisUrl"`
			DbPaasRdsURL            string `json:"dbPaasRdsUrl"`
			PaaSPayAsYouGoProducts  string `json:"PaaSPayAsYouGoProducts"`
			MiddlewarePaasKafkaURL  string `json:"middlewarePaasKafkaUrl"`
			DbPaasMemcacheURL       string `json:"dbPaasMemcacheUrl"`
			DbPaasTsdbURL           string `json:"dbPaasTsdbUrl"`
			DbPaasHbaseURL          string `json:"dbPaasHbaseUrl"`
			DbPaasMongodb           string `json:"dbPaasMongodb"`
			DbPaasRds               string `json:"dbPaasRds"`
			DbPaasDrds              string `json:"dbPaasDrds"`
			PaasAutoRenewProducts   string `json:"PaasAutoRenewProducts"`
			MiddlewarePaasRabbitURL string `json:"middlewarePaasRabbitUrl"`
			PaasUpgradeProducts     string `json:"PaasUpgradeProducts"`
			MiddlewarePaasRabbit    string `json:"middlewarePaasRabbit"`
			MiddlewarePaasMqtt      string `json:"middlewarePaasMqtt"`
			MiddlewarePaasMqURL     string `json:"middlewarePaasMqUrl"`
			DbPaasMemcache          string `json:"dbPaasMemcache"`
		} `json:"paas"`
		Rds struct {
			RdsStorageType string `json:"rdsStorageType"`
			RdsStorage     string `json:"rdsStorage"`
		} `json:"rds"`
		ACL struct {
			HasACL string `json:"hasAcl"`
		} `json:"acl"`
		Order struct {
			CtcsVMBind          string `json:"ctcsVmBind"`
			HasRecoveredOrder   string `json:"hasRecoveredOrder"`
			HasAutoRevoke       string `json:"hasAutoRevoke"`
			HasAutoRenew        string `json:"hasAutoRenew"`
			HasChBandwidthOrder string `json:"hasChBandwidthOrder"`
			SupportDownConfig   string `json:"supportDownConfig"`
			BillingByQuantity   string `json:"billingByQuantity"`
			HasCycleTerminate   string `json:"hasCycleTerminate"`
		} `json:"order"`
	} `json:"returnObj"`
}

type RegionZoneListRequest struct {
	RegionID string `json:"regionID"`
}

type ZoneEntry struct {
	// 如 cn-hn-cs42-hncs1A-public-ctcloud
	Name string `json:"name"`
	// 中文名, 如可用区1
	AzDisplayName string `json:"azDisplayName"`
}

type RegionZoneListResponse struct {
	ReturnObj struct {
		ZoneList []ZoneEntry `json:"zoneList"`
	} `json:"returnObj"`
}

type RegionResourceSummaryRequest struct {
	RegionID string `json:"regionID"`
}

type RegionResourceSummaryResponse struct {
	ReturnObj struct {
		Resources struct {
			VolumeSnapshot struct {
				TotalCount       int `json:"total_count"`
				DetailTotalCount int `json:"detail_total_count"`
			} `json:"VOLUME_SNAPSHOT"`
			VMGroup struct {
				TotalCount int            `json:"total_count"`
				Detail     map[string]int `json:"detail"`
			} `json:"Vm_Group"`
			Acllist struct {
				TotalCount int            `json:"total_count"`
				Detail     map[string]int `json:"detail"`
			} `json:"ACLLIST"`
			Nat struct {
				TotalCount       int            `json:"total_count"`
				DetailTotalCount int            `json:"detail_total_count"`
				Detail           map[string]int `json:"detail"`
			} `json:"NAT"`
			IPPool struct {
				TotalCount int            `json:"total_count"`
				Detail     map[string]int `json:"detail"`
			} `json:"IP_POOL"`
			PublicIP struct {
				TotalCount       int `json:"total_count"`
				DetailTotalCount int `json:"detail_total_count"`
			} `json:"Public_IP"`
			Image struct {
				TotalCount int            `json:"total_count"`
				Detail     map[string]int `json:"detail"`
			} `json:"IMAGE"`
			Volume struct {
				VoRootSize       int `json:"vo_root_size"`
				VoDiskCount      int `json:"vo_disk_count"`
				TotalSize        int `json:"total_size"`
				VoDiskSize       int `json:"vo_disk_size"`
				DetailTotalCount int `json:"detail_total_count"`
				TotalCount       int `json:"total_count"`
				VoRootCount      int `json:"vo_root_count"`
			} `json:"Volume"`
			Bms struct {
				MemoryCount        int `json:"memory_count"`
				TotalCount         int `json:"total_count"`
				DetailTotalCount   int `json:"detail_total_count"`
				CPUCount           int `json:"cpu_count"`
				BmShutdCount       int `json:"bm_shutd_count"`
				ExpireRunningCount int `json:"expire_running_count"`
				BmRunningCount     int `json:"bm_running_count"`
				ExpireCount        int `json:"expire_count"`
				ExpireShutdCount   int `json:"expire_shutd_count"`
			} `json:"BMS"`
			Snapshot struct {
				TotalCount int            `json:"total_count"`
				Detail     map[string]int `json:"detail"`
			} `json:"SNAPSHOT"`
			LbListener struct {
				TotalCount int            `json:"total_count"`
				Detail     map[string]int `json:"detail"`
			} `json:"LB_LISTENER"`
			DiskBackup struct {
				TotalCount       int            `json:"total_count"`
				DetailTotalCount int            `json:"detail_total_count"`
				Detail           map[string]int `json:"detail"`
			} `json:"Disk_Backup"`
			Loadbalancer struct {
				TotalCount int            `json:"total_count"`
				Detail     map[string]int `json:"detail"`
			} `json:"LOADBALANCER"`
			Vpc struct {
				TotalCount       int `json:"total_count"`
				DetailTotalCount int `json:"detail_total_count"`
			} `json:"VPC"`
			VM struct {
				VMShutdCount       int     `json:"vm_shutd_count"`
				MemoryCount        float64 `json:"memory_count"`
				ExpireCount        int     `json:"expire_count"`
				DetailTotalCount   int     `json:"detail_total_count"`
				CPUCount           int     `json:"cpu_count"`
				ExpireRunningCount int     `json:"expire_running_count"`
				TotalCount         int     `json:"total_count"`
				ExpireShutdCount   int     `json:"expire_shutd_count"`
				VMRunningCount     int     `json:"vm_running_count"`
			} `json:"VM"`
			OSBackup struct {
				TotalSize        int `json:"total_size"`
				DetailTotalCount int `json:"detail_total_count"`
			} `json:"OS_Backup"`
			TrafficMirrorFlow struct {
				TotalCount int            `json:"total_count"`
				Detail     map[string]int `json:"detail"`
			} `json:"TrafficMirror_Flow"`
			TrafficMirrorFilter struct {
				TotalCount int            `json:"total_count"`
				Detail     map[string]int `json:"detail"`
			} `json:"TrafficMirror_Filter"`
			Cbr struct {
				TotalSize        int            `json:"total_size"`
				TotalCount       int            `json:"total_count"`
				DetailTotalCount int            `json:"detail_total_count"`
				Detail           map[string]int `json:"detail"`
			} `json:"CBR"`
			Cert struct {
				TotalCount int            `json:"total_count"`
				Detail     map[string]int `json:"detail"`
			} `json:"CERT"`
			AzDisplayName string `json:"az_display_name"`
			CbrVbs        struct {
				TotalCount       int            `json:"total_count"`
				DetailTotalCount int            `json:"detail_total_count"`
				Detail           map[string]int `json:"detail"`
			} `json:"CBR_VBS"`
		} `json:"resources"`
	} `json:"returnObj"`
}
