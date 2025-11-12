package version130

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/executil"
	"github.com/wangweihong/gotoolbox/pkg/log"
	"github.com/wangweihong/gotoolbox/pkg/netutil"
	"github.com/wangweihong/gotoolbox/pkg/template"
	"github.com/wangweihong/gotoolbox/pkg/typeutil"

	"github.com/wangweihong/eazycloud/apis/ikubeagent"
	"github.com/wangweihong/eazycloud/internal/pkg/run"
)

// func runKubeadmInit(config *ikubeagent.KubernetesMasterDeployConfig, isMaster0 bool) error {
// 	var stdout string
// 	var stderr string
// 	var err error
// 	defer func() {
// 		recordKubeadmLog(stdout, stderr)
// 	}()
// 	if err := generateKubeadmConfigYaml(config); err != nil {
// 		return err
// 	}
// 	args := []string{"init"}
// 	args = append(args, "--config", ikubeagent.KubeadmConfigYamlPath)
// 	if config.ControlPlaneConfig.HAConfig != nil {
// 		args = append(args, "--upload-certs")
// 	}

// 	args = append(args, "-v=9")
// 	if config.DataBaseConfig != nil && config.DataBaseConfig.Local != nil && (config.DataBaseConfig.Local.PeerPort != nil || config.DataBaseConfig.Local.ListenPort != nil) {
// 		// kubeadm preflight阶段硬编码了2379和2380的检测, 即使kubeadm.yaml更改了etcd端口仍然会检测这两个端口
// 		args = append(args, "--ignore-preflight-errors=Port-2379,Port-2380")
// 	}

// 	log.Debugf("run command:%v,args:%v", ikubeagent.KubeadmBinary, args)

// 	f, err := os.Create(ikubeagent.KubeadmLogPath)
// 	if err != nil {
// 		log.Errorf("deploy k8s control-plane [%v:%v] fail:%v", ikubeagent.KubeadmBinary, args, run.TrimError(err))
// 		return run.TrimError(err)
// 	}
// 	defer f.Close()

// 	stdout, stderr, err = executil.ExecuteCmdSplitStdoutStderr(ikubeagent.KubeadmBinary, args, 600)
// 	if err != nil {
// 		log.Errorf("deploy k8s control-plane [%v:%v] fail:%v", ikubeagent.KubeadmBinary, args, run.TrimError(err))
// 		return run.TrimError(err)
// 	}

// 	log.Infof("run kubeadm init success")
// 	return nil
// }

func generateKubeadmConfigYaml(config *ikubeagent.KubernetesDeployConfig) error {
	if config.ControlPlaneConfig == nil {
		return fmt.Errorf("control plane config is empty")
	}

	var repository string
	if config.RegistryConfig != nil {
		repository = config.RegistryConfig.Address
		if config.RegistryConfig.Project != "" {
			repository = repository + "/" + config.RegistryConfig.Project
		}
	}

	localIP, err := netutil.GetLocalIP()
	if err != nil {
		return err
	}

	ctxs := map[string]any{
		"LocalAdvertiseAddress": localIP,
		//kubeadm will use image repository replace registry.k8s.io, but we need to keep registry.k8s.io in project
		//for better distinguish image source.
		"ImageRepository": repository + "/" + "registry.k8s.io",
	}

	if config.DataBaseConfig.Local != nil {
		ctxs["EtcdListenPort"] = typeutil.GenericIndirectPrefined(config.DataBaseConfig.Local.ListenPort, 2379)
		ctxs["EtcdPeerPort"] = typeutil.GenericIndirectPrefined(config.DataBaseConfig.Local.PeerPort, 2380)
	}

	if config.ControlPlaneConfig.HAConfig != nil {
		ctxs["ControlPlaneEndpoint"] = config.ControlPlaneConfig.HAConfig.VIP + ":" + strconv.Itoa(
			int(config.ControlPlaneConfig.HAConfig.ExternalPort),
		)
	}

	if config.ControlPlaneConfig.ServiceCIDR != "" {
		_, svcSubnetCIDR, err := net.ParseCIDR(config.ControlPlaneConfig.ServiceCIDR)
		if err != nil {
			return err
		}
		//calculate new dns server ip, otherwise nslookup in pod will fail because cannot connect to dns server.
		dnsIP, err := netutil.GetIndexedIP(svcSubnetCIDR, 10)
		if err != nil {
			return errors.Wrap(err, "unable to get internal Kubernetes Service IP from the given service CIDR")
		}

		ctxs["ServiceSubnet"] = config.ControlPlaneConfig.ServiceCIDR
		ctxs["ClusterDNS"] = dnsIP.String()
	}
	// pod CIDR 必须和网络插件保持一致
	ctxs["PodSubnet"] = config.ControlPlaneConfig.PodCIDR

	if config.NodeConfig != nil {
		ctxs["NodeRegistrationName"] = config.NodeConfig.NodeName
	}

	if config.ControlPlaneConfig.KubeProxy != nil {
		if config.ControlPlaneConfig.KubeProxy.Mode != "" {
			ctxs["KubeProxyMode"] = config.ControlPlaneConfig.KubeProxy.Mode
		}
	}

	t := KubeadmConfigTemplate
	if err := t.SetContexts(ctxs).LocateToDisk().Error(); err != nil {
		log.Errorf("generate template %v fail:%v", t.Name(), err)
		return err
	}
	return nil
}

type postInstall struct {
	err error
}

func (pi postInstall) createSystemNamespace() error {
	if pi.err != nil {
		return pi.err
	}
	return run.RunTemplate(&NamespacesManifests, nil)

}

func (pi postInstall) deployNetworkPlugin(config *ikubeagent.KubernetesDeployConfig) error {
	if pi.err != nil {
		return pi.err
	}

	plugin := config.ControlPlaneConfig.NetworkPlugin
	if plugin == nil {
		return fmt.Errorf("network plugin config is empty")
	}

	if plugin.PluginType == ikubeagent.NetworkPluginCalico {
		if plugin.CalicoConfig == nil {
			return fmt.Errorf("network plugin calico config is nil")
		}

		calicoTemplate := CalicoManifest
		calicoCtx := map[string]any{}

		if config.ControlPlaneConfig.PodCIDR != "" {
			if _, _, err := net.ParseCIDR(config.ControlPlaneConfig.PodCIDR); err != nil {
				return fmt.Errorf("invalid pod subnet %v: %v", config.ControlPlaneConfig.PodCIDR, err)
			}
			calicoCtx["PodSubnet"] = config.ControlPlaneConfig.PodCIDR
		}

		log.Infof("calico context:%v", calicoCtx)
		if err := run.RunTemplateError(nil, &calicoTemplate, calicoCtx); err != nil {
			return err
		}

		return nil

	}
	err := fmt.Errorf("network plugintype %v not support ", plugin.PluginType)
	log.Error(err.Error())
	return err
}

func runTemplateOperation(c *setupContext, template template.ProcessorInterface, ctx map[string]any, fileMode os.FileMode) setupFunc {
	return func(c *setupContext) error {
		return run.RunTemplate(template, ctx)
	}
}

func installMonitorPlugin(c *setupContext) error {
	if c.config.ControlPlaneConfig == nil || c.config.ControlPlaneConfig.MonitorPlugin == nil {
		log.Infof("ignore monitor plugin deploy phase for non-monitor plugin config")
		return nil
	}

	ctx := map[string]any{
		"ImageRepository": getImageRepository(c.config, "monitor"),
	}

	Ops := []NamedOperation{
		{Name: "install metrics server", Op: runTemplateOperation(c, KubeadmConfigTemplate, ctx, 0644)},
		{Name: "install prometheus Operator Template", Op: runTemplateOperation(c, MonitorOperatorsTemplate, ctx, 0644)},
		{Name: "install promethues Stack Template", Op: runTemplateOperation(c, MonitorStackTemplate, ctx, 0644)},
		{Name: "generate etcd secret", Op: createSecretOfEtcdForMonitoring},
	}

	for _, op := range Ops {
		setupDep := LogOperation(op.Name, op.Op)
		if err := setupDep(c); err != nil {
			return err
		}
	}
	return nil
}

func installDisplayCardPlugin(c *setupContext) error {
	Ops := []NamedOperation{}
	ctx := map[string]any{}

	if c.config.ControlPlaneConfig != nil && c.config.ControlPlaneConfig.Gpu != nil {
		Ops = append(Ops, NamedOperation{Name: "install gpu plugin", Op: runTemplateOperation(c, Gpu, ctx, 0644)})
	}

	if c.config.ControlPlaneConfig != nil && c.config.ControlPlaneConfig.Npu != nil {
		Ops = append(Ops, NamedOperation{Name: "install npu plugin", Op: runTemplateOperation(c, Npu, ctx, 0644)})
	}

	for _, op := range Ops {
		setupDep := LogOperation(op.Name, op.Op)
		if err := setupDep(c); err != nil {
			return err
		}
	}
	return nil
}

func installGpuPlugin(c *setupContext) error {
	if c.config.ControlPlaneConfig == nil || c.config.ControlPlaneConfig.Gpu == nil {
		log.Infof("ignore deploy phase for non gpu plugin config")
		return nil
	}
	ctx := map[string]any{
		"ImageRepository": getImageRepository(c.config, "gpu"),
	}

	Ops := []NamedOperation{
		{Name: "install gpu plugin", Op: runTemplateOperation(c, Gpu, ctx, 0644)},
	}
	for _, op := range Ops {
		setupDep := LogOperation(op.Name, op.Op)
		if err := setupDep(c); err != nil {
			return err
		}
	}
	log.Info("deploy gpu plugin success")
	return nil
}

func installStoragePlugin(c *setupContext) error {
	if c.config.ControlPlaneConfig == nil || c.config.ControlPlaneConfig.StoragePlugin == nil {
		log.Infof("ignore storage plugin deploy phase for non-storage plugin config")
		return nil
	}
	ctx := map[string]any{
		"ImageRepository": getImageRepository(c.config, "storage"),
	}

	Ops := []NamedOperation{}
	if c.config.ControlPlaneConfig.StoragePlugin.Local != nil {
		Ops = append(Ops, NamedOperation{Name: "install local storage plugin", Op: runTemplateOperation(c, LocalStorageAllInOneTemplate, ctx, 0644)})
	}

	if c.config.ControlPlaneConfig.StoragePlugin.Nfs != nil {
		Ops = append(Ops, NamedOperation{Name: "install nfs storage plugin", Op: runTemplateOperation(c, NFSStorageClassTemplate, ctx, 0644)})
	}

	if c.config.ControlPlaneConfig.StoragePlugin.NFSVolume != nil {
		log.Infof("deploy nfs pv storage plugin.")
		ctx := map[string]any{
			"NFSIP": c.config.ControlPlaneConfig.StoragePlugin.NFSVolume.ServerIP,
		}

		if c.config.ControlPlaneConfig.StoragePlugin.NFSVolume.MountPath != "" {
			ctx["NFSHostPath"] = c.config.ControlPlaneConfig.StoragePlugin.NFSVolume.MountPath
		}
		Ops = append(Ops, NamedOperation{Name: "install nfs volume plugin", Op: runTemplateOperation(c, NFSStorageClassTemplate, ctx, 0644)})

	}

	for _, op := range Ops {
		setupDep := LogOperation(op.Name, op.Op)
		if err := setupDep(c); err != nil {
			return err
		}
	}

	return nil
}

// func (pi postInstall) deployNamespaceControllerPlugin(config *ikubeagent.KubernetesDeployConfig) error {
// 	if pi.err != nil {
// 		return pi.err
// 	}

// 	if config.ControlPlaneConfig == nil || config.ControlPlaneConfig.NamespaceConfig == nil {
// 		log.Infof("ignore storage plugin deploy phase for non-namespace-controller plugin config")
// 		return nil
// 	}
// 	//install metrics server
// 	repo := getImageRepository(config, "namespace-controller")
// 	ctx := map[string]any{
// 		"ImageRepository": repo,
// 	}

// 	log.Infof("deploy namespace controller plugin.")
// 	if err := run.RunTemplateError(nil, &NamespacesControllerManifests, ctx); err != nil {
// 		return err
// 	}

// 	return nil
// }

func installKubectlPlugin(c *setupContext) error {
	if c.config.ControlPlaneConfig == nil || c.config.ControlPlaneConfig.KubectlPlugin == nil {
		log.Infof("ignore kubectl plugin deploy phase for non-kubectl plugin config")
		return nil
	}

	ctx := map[string]any{
		"ImageRepository": getImageRepository(c.config, "kubectl"),
	}
	Ops := []NamedOperation{
		{Name: "install kubectl plugin", Op: runTemplateOperation(c, KubectlManifest, ctx, 0644)},
	}
	for _, op := range Ops {
		setupDep := LogOperation(op.Name, op.Op)
		if err := setupDep(c); err != nil {
			return err
		}
	}

	return nil
}

func deployPrepareHA(c *setupContext) error {
	// ignore non-control-plane and ha
	if c.config.ControlPlaneConfig == nil || c.config.ControlPlaneConfig.HAConfig == nil {
		log.Infof("ignore ha Prepare phase form non-ha config")
		return nil
	}

	keepalivedServiceOps := []NamedOperation{
		{Name: "generateKeepalivedCheckerScript", Op: generateKeepalivedCheckerScript},
		{Name: "generateKeepalivedService", Op: generateKeepalivedService},
		{Name: "generateKeepalivedConfig", Op: generateKeepalivedConfig},
	}

	for _, op := range keepalivedServiceOps {
		setupDep := LogOperation(op.Name, op.Op)
		if err := setupDep(c); err != nil {
			return err
		}
	}

	haproxyServiceOps := []NamedOperation{
		{Name: "generateHaproxyConfig", Op: generateHaproxyConfig},
		{Name: "generateHaproxyPodTemplate", Op: generateHaproxyPodTemplate},
	}

	for _, op := range haproxyServiceOps {
		setupDep := LogOperation(op.Name, op.Op)
		if err := setupDep(c); err != nil {
			return err
		}
	}

	log.Infof("write kubeadm config file > %v", c.KubeadmConfigYamlPath)
	if err := os.WriteFile(c.KubeadmConfigYamlPath, []byte(c.config.ControlPlaneConfig.HAConfig.KubeadmConfigYaml), 0755); err != nil {
		log.Errorf("write kubeadm config file > %v err:%v", c.KubeadmConfigYamlPath, err)
		return err
	}

	log.Infof("run prepare work for HA complete")
	return nil
}

func generateHaproxyService(config *ikubeagent.KubernetesDeployConfig) error {
	ctx := make(map[string]any)
	ctx["HAPROXY_PORT"] = strconv.Itoa(int(config.ControlPlaneConfig.HAConfig.ExternalPort))
	var serverIPs []string
	for _, v := range config.ControlPlaneConfig.HAConfig.ControlPlaneNodes {
		serverIPs = append(serverIPs, v.IP)
	}
	ctx["ServerIPs"] = serverIPs
	ctx["ServerPort"] = 6443
	t := HaproxyConfigTemplate
	if err := t.SetContexts(ctx).LocateToDisk().Error(); err != nil {
		log.Errorf("save template %v to path %v fail:%v", t.Name(), t.FilePath, err)
		return err
	}

	var repository string
	if config.RegistryConfig != nil {
		repository = config.RegistryConfig.Address
		if config.RegistryConfig.Project != "" {
			repository = repository + "/" + config.RegistryConfig.Project
		}
	}

	ctx = make(map[string]any)
	ctx["EXTERNAL_PORT"] = strconv.Itoa(int(config.ControlPlaneConfig.HAConfig.ExternalPort))
	ctx["ImageRepository"] = repository
	t = HaproxyYamlTemplate
	if err := t.SetContexts(ctx).LocateToDisk().Error(); err != nil {
		log.Errorf("save template %v to path %v fail:%v", t.Name(), t.FilePath, err)
		return err
	}

	return nil
}

func createSecretOfEtcdForMonitoring(c *setupContext) error {
	if c.config.ControlPlaneConfig == nil || c.config.ControlPlaneConfig.MonitorPlugin == nil {
		log.Infof("ignore createSecretOfEtcd deploy phase for non-monitor plugin config")
		return nil
	}
	// kubectl --kubeconfig /etc/kubernetes/admin.conf  create secret \
	// 		generic etcd-ca --from-file=/etc/kubernetes/pki/etcd/ca.key --from-file=/etc/kubernetes/pki/etcd/ca.crt -n
	// monitoring
	secretName := "etcd-ca"
	namespace := "monitoring"
	cakey := fmt.Sprintf("--from-file=%v", ikubeagent.EtcdCaPath+"ca.key")
	cacrt := fmt.Sprintf("--from-file=%v", ikubeagent.EtcdCaPath+"ca.crt")

	args := []string{
		"--kubeconfig",
		ikubeagent.KubeConfigPath,
		"create",
		"secret",
		"generic",
		secretName,
		cakey,
		cacrt,
		"-n",
		namespace,
	}
	_, stderr, err := executil.ExecuteCmdSplitStdoutStderr(ikubeagent.KubectlBinary, args, 0)
	if err != nil || stderr != "" {
		log.Errorf("run command [%v:%v] fail:%v,stderr:%v", ikubeagent.KubectlBinary, args, run.TrimError(err), stderr)
		return run.TrimError(err)
	}

	return nil
}

func setMasterAsWorker(c *setupContext) error {
	if !c.config.MasterAlsoWorker {
		log.Info("remain master, not changed to worker")
		return nil
	}
	nodeName := c.config.NodeConfig.NodeName

	args := []string{
		"--kubeconfig",
		ikubeagent.KubeConfigPath,
		"taint",
		"nodes",
		nodeName,
		"node-role.kubernetes.io/control-plane:NoSchedule-",
	}
	_, stderr, err := executil.ExecuteCmdSplitStdoutStderr(ikubeagent.KubectlBinary, args, 0)
	if err != nil || stderr != "" {
		log.Errorf("run command [%v:%v] fail:%v,stderr:%v", ikubeagent.KubectlBinary, args, run.TrimError(err), stderr)
		return run.TrimError(err)
	}

	return nil
}

func generateKeepalivedCheckerScript(c *setupContext) error {
	ctx := make(map[string]any)
	ctx["ExternalPort"] = strconv.Itoa(int(c.config.ControlPlaneConfig.HAConfig.ExternalPort))
	ctx["APISERVER_VIP"] = c.config.ControlPlaneConfig.HAConfig.VIP

	return run.WriteTemplate(KeepalivedCheckApiserverScript, ctx, 0755)
}

func generateKeepalivedService(c *setupContext) error {
	var repository string
	if c.config.RegistryConfig != nil {
		repository = c.config.RegistryConfig.Address
		if c.config.RegistryConfig.Project != "" {
			repository = repository + "/" + c.config.RegistryConfig.Project
		}
	}
	ctx := make(map[string]any)
	ctx["ImageRepository"] = repository

	return run.WriteTemplate(KeepalivedYamlTemplate, ctx, 0644)
}

func generateKeepalivedConfig(c *setupContext) error {
	log.Infof("start to generate keepalived config")

	intf, _, err := netutil.GetInterfaceAndIP()
	if err != nil {
		log.Errorf("cannot get interface and ip: %v", err)
		return err
	}
	log.Infof("get netif %v", intf)
	ctx := make(map[string]any)
	ctx["APISERVER_VIP"] = c.config.ControlPlaneConfig.HAConfig.VIP
	ctx["INTERFACE"] = intf
	ctx["STATE"] = "MASTER"
	ctx["PRIORITY"] = fmt.Sprintf("%v", c.config.NodeConfig.Priority)
	ctx["ROUTER_ID"] = fmt.Sprintf("%v", c.config.ControlPlaneConfig.HAConfig.RouteID)

	if !c.isMaster0 {
		//FIXME how to calculate other non-master keepalived priority without conflicting?
		//pass priority by caller
		ctx["STATE"] = "BACKUP"
		ctx["PRIORITY"] = fmt.Sprintf("%v", c.config.NodeConfig.Priority)
	}

	if err := run.WriteTemplate(KeepalivedConfigTempl, ctx, 0644); err != nil {
		return err
	}

	return nil
}

func generateHaproxyConfig(c *setupContext) error {
	log.Infof("start to generate haproxy config")

	var serverIPs []string
	for _, v := range c.config.ControlPlaneConfig.HAConfig.ControlPlaneNodes {
		serverIPs = append(serverIPs, v.IP)
	}
	ctx := make(map[string]any)
	ctx["HAPROXY_PORT"] = strconv.Itoa(int(c.config.ControlPlaneConfig.HAConfig.ExternalPort))
	ctx["ServerIPs"] = serverIPs
	ctx["ServerPort"] = 6443

	if err := run.WriteTemplate(HaproxyConfigTemplate, ctx, 0644); err != nil {
		return err
	}

	log.Infof("generate haproxy config success")
	return nil
}

func generateHaproxyPodTemplate(c *setupContext) error {
	log.Infof("start to generate haproxy pod template")

	ctx := make(map[string]any)
	ctx["EXTERNAL_PORT"] = strconv.Itoa(int(c.config.ControlPlaneConfig.HAConfig.ExternalPort))
	ctx["ImageRepository"] = getImageRepository(c.config, "")
	if err := run.WriteTemplate(HaproxyYamlTemplate, ctx, 0644); err != nil {
		return err
	}

	log.Infof("generate haproxy pod template success")
	return nil
}

func execKubeadmInit(c *setupContext) error {
	log.Infof("start to exec kubeadm init")

	var stdout string
	var stderr string
	var err error
	defer func() {
		recordKubeadmLog(stdout, stderr)
	}()
	if err := generateKubeadmConfigYaml(c.config); err != nil {
		return err
	}
	args := []string{"init"}
	args = append(args, "--config", ikubeagent.KubeadmConfigYamlPath)
	if c.config.ControlPlaneConfig.HAConfig != nil {
		args = append(args, "--upload-certs")
	}

	log.Debugf("run command:%v,args:%v", ikubeagent.KubeadmBinary, args)

	f, err := os.Create(ikubeagent.KubeadmLogPath)
	if err != nil {
		log.Errorf("deploy k8s control-plane [%v:%v] fail:%v", ikubeagent.KubeadmBinary, args, run.TrimError(err))
		return run.TrimError(err)
	}
	defer f.Close()

	stdout, stderr, err = executil.ExecuteCmdSplitStdoutStderr(ikubeagent.KubeadmBinary, args, 600)
	if err != nil {
		log.Errorf("deploy k8s control-plane [%v:%v] fail:%v", ikubeagent.KubeadmBinary, args, run.TrimError(err))
		return run.TrimError(err)
	}

	log.Infof("run kubeadm init success")
	return nil
}

func generateKubeadmConfig(c *setupContext) error {
	log.Infof("start to generate kubeadm config")

	if c.config.ControlPlaneConfig == nil {
		return errors.New("control plane config is empty")
	}

	localIP, err := netutil.GetLocalIP()
	if err != nil {
		return err
	}

	ctxs := map[string]any{
		"LocalAdvertiseAddress": localIP,
		//kubeadm will use image repository replace registry.k8s.io, but we need to keep registry.k8s.io in project
		//for better distinguish image source.
		"ImageRepository": getImageRepository(c.config, "") + "/" + "registry.k8s.io",
	}
	// pod CIDR 必须和网络插件保持一致
	ctxs["PodSubnet"] = c.config.ControlPlaneConfig.PodCIDR

	if c.config.ControlPlaneConfig.HAConfig != nil {
		ctxs["ControlPlaneEndpoint"] = c.config.ControlPlaneConfig.HAConfig.VIP + ":" + strconv.Itoa(
			int(c.config.ControlPlaneConfig.HAConfig.ExternalPort),
		)
	}

	if c.config.ControlPlaneConfig.ServiceCIDR != "" {
		_, svcSubnetCIDR, err := net.ParseCIDR(c.config.ControlPlaneConfig.ServiceCIDR)
		if err != nil {
			return err
		}
		//calculate new dns server ip, otherwise nslookup in pod will fail for disconnect to dns server.
		dnsIP, err := netutil.GetIndexedIP(svcSubnetCIDR, 10)
		if err != nil {
			return errors.Wrap(err, "unable to get internal Kubernetes Service IP from the given service CIDR")
		}

		ctxs["ServiceSubnet"] = c.config.ControlPlaneConfig.ServiceCIDR
		ctxs["ClusterDNS"] = dnsIP.String()
	}

	if c.config.NodeConfig != nil {
		ctxs["NodeRegistrationName"] = c.config.NodeConfig.NodeName
	}

	if c.config.ControlPlaneConfig.KubeProxy != nil && c.config.ControlPlaneConfig.KubeProxy.Mode != "" {
		ctxs["KubeProxyMode"] = c.config.ControlPlaneConfig.KubeProxy.Mode
	}

	if err := run.WriteTemplate(KubeadmConfigTemplate, ctxs, 0644); err != nil {
		return err
	}

	log.Infof("generate kubeadm config success")
	return nil
}

func installSystemNamespace(c *setupContext) error {
	return run.RunTemplate(NamespacesManifests, nil)
}

func installNetworkPlugin(c *setupContext) error {
	plugin := c.config.ControlPlaneConfig.NetworkPlugin
	if plugin == nil {
		return fmt.Errorf("network plugin config is empty")
	}

	if plugin.PluginType == ikubeagent.NetworkPluginCalico {
		log.Infof("install calico network plugin")
		if plugin.CalicoConfig == nil {
			return fmt.Errorf("network plugin calico config is nil")
		}

		ctx := map[string]any{
			"PodSubnet": ikubeagent.CalicoDefaultPodSubnet,
		}

		if c.config.ControlPlaneConfig.PodCIDR != "" {
			if _, _, err := net.ParseCIDR(c.config.ControlPlaneConfig.PodCIDR); err != nil {
				return errors.Errorf("invalid pod subnet %v: %v", c.config.ControlPlaneConfig.PodCIDR, err)
			}
			ctx["PodSubnet"] = c.config.ControlPlaneConfig.PodCIDR
		}

		if err := run.RunTemplate(CalicoManifest, ctx); err != nil {
			return err
		}

		return nil

	}
	err := fmt.Errorf("network plugintype %v not support ", plugin.PluginType)
	log.Error(err.Error())
	return err
}
