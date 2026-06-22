package kubernetes

import (
	"bytes"
	"context"
	"net"
	"sort"
	"strconv"
	"time"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/apis/ikubeagent"
	"github.com/wangweihong/eazycloud/apis/imachinery"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/netutil"
	"github.com/wangweihong/gotoolbox/pkg/paging"
	"github.com/wangweihong/gotoolbox/pkg/waitgroup"
)

const (
	ClusterNodeScaleMedium  = "medium"
	ClusterNodeScaleLarge   = "large"
	ClusterNodeScaleXlarge  = "xlarge"
	ClusterNodeScale2XLarge = "2xlarge"
	ClusterNodeScale4Xlarge = "4xlarge"
	ClusterNodeScale8Xlarge = "8xlarge"
)

func (k *kubernetesService) ScanHosts(ctx context.Context, req *iapiserver.HostScanRequest) (*iapiserver.HostScanResponse, error) {
	ipRanges, err := netutil.GetIPRange(req.Begin, req.End)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(ipRanges) == 0 {
		return nil, errors.Errorf("no ip in range")
	}

	wg := waitgroup.RunGenericConcurrently(ctx, ipRanges, func(ctx context.Context, ip string) waitgroup.GenericResult[*iapiserver.HostScanInfo] {
		c := libs.NewKubeagentClient(ip)
		if _, err := c.CheckDependency(ctx, &ikubeagent.CheckDependencyRequest{}, nil); err != nil {
			return waitgroup.NewGenericResult[*iapiserver.HostScanInfo](nil, err)
		}

		var installResp ikubeagent.InstallStateResponse
		if _, err := c.GetInstallState(ctx, nil, &installResp); err != nil {
			return waitgroup.NewGenericResult[*iapiserver.HostScanInfo](nil, err)
		}

		if installResp.State.State == string(ikubeagent.KubernetesDeployStateUninitialized) {
			ret := &iapiserver.HostScanInfo{
				IP:       ip,
				HostInfo: installResp.HostStat,
			}
			return waitgroup.NewGenericResult[*iapiserver.HostScanInfo](ret, nil)
		}
		return waitgroup.NewGenericResult[*iapiserver.HostScanInfo](nil, errors.Errorf("node %v has installed, current state:%v", ip, installResp.State.State))
	}, 10*time.Second)

	iplist := wg.GetSuccessResultList()
	sort.SliceStable(iplist, func(i, j int) bool {
		return bytes.Compare(net.ParseIP(iplist[i].IP), net.ParseIP(iplist[j].IP)) < 0
	})
	total := len(iplist)
	return &iapiserver.HostScanResponse{
		HostsInfo:  paging.Cut(len(iplist), req.PageNum, req.PageSize, iplist),
		TotalCount: total,
	}, nil
}

func (k *kubernetesService) VipUsedTest(ctx context.Context, req *iapiserver.VipUsedTestRequest) error {
	if netutil.PingV2(req.Vip) {
		return errors.Errorf("vip %v has been used", req.Vip)
	}

	if req.LocalAdvisePort == req.ExternalPort {
		return errors.Errorf("external port %v must not equal to local advertise port %v", req.ExternalPort, req.LocalAdvisePort)
	}

	wg := waitgroup.RunGenericConcurrently(ctx, req.NodeList, func(ctx context.Context, ip string) waitgroup.GenericResult[*imachinery.Empty] {
		_, err := net.DialTimeout("tcp", net.JoinHostPort(ip, strconv.Itoa(req.LocalAdvisePort)), 5*time.Second)
		if err == nil {
			return waitgroup.NewGenericResult[*imachinery.Empty](nil, errors.Errorf("port %v has beed used in ip %v", req.LocalAdvisePort, ip))
		}
		return waitgroup.NewGenericResult[*imachinery.Empty](nil, nil)
	}, 5*time.Second)

	return wg.GetErrorList()
}

func (k *kubernetesService) SetupCluster(ctx context.Context, req *iapiserver.VipUsedTestRequest) error {
	if netutil.PingV2(req.Vip) {
		return errors.Errorf("vip %v has been used", req.Vip)
	}

	if req.LocalAdvisePort == req.ExternalPort {
		return errors.Errorf("external port %v must not equal to local advertise port %v", req.ExternalPort, req.LocalAdvisePort)
	}

	wg := waitgroup.RunGenericConcurrently(ctx, req.NodeList, func(ctx context.Context, ip string) waitgroup.GenericResult[*imachinery.Empty] {
		_, err := net.DialTimeout("tcp", net.JoinHostPort(ip, strconv.Itoa(req.LocalAdvisePort)), 5*time.Second)
		if err == nil {
			return waitgroup.NewGenericResult[*imachinery.Empty](nil, errors.Errorf("port %v has beed used in ip %v", req.LocalAdvisePort, ip))
		}
		return waitgroup.NewGenericResult[*imachinery.Empty](nil, nil)
	}, 5*time.Second)

	return wg.GetErrorList()
}

var (
	mediumClusterNodeScaleRequirement = iapiserver.ClusterNodeScaleRequirement{
		Name:             ClusterNodeScaleMedium,
		MinNodeNum:       1,
		MaxNodeNum:       5,
		LowestCpuCores:   1,
		LowestMemoryUnit: 3.75 * 1024 * 1024 * 1024,
	}
	largeClusterNodeScaleRequirement = iapiserver.ClusterNodeScaleRequirement{
		Name:             ClusterNodeScaleLarge,
		MinNodeNum:       6,
		MaxNodeNum:       10,
		LowestCpuCores:   2,
		LowestMemoryUnit: 8 * 1024 * 1024 * 1024,
	}
	xLargeClusterNodeScaleRequirement = iapiserver.ClusterNodeScaleRequirement{
		Name:             ClusterNodeScaleXlarge,
		MinNodeNum:       11,
		MaxNodeNum:       100,
		LowestCpuCores:   4,
		LowestMemoryUnit: 15 * 1024 * 1024 * 1024,
	}
	x2LargeClusterNodeScaleRequirement = iapiserver.ClusterNodeScaleRequirement{
		Name:             ClusterNodeScale2XLarge,
		MinNodeNum:       101,
		MaxNodeNum:       250,
		LowestCpuCores:   8,
		LowestMemoryUnit: 30 * 1024 * 1024 * 1024,
	}
	x4LargeClusterNodeScaleRequirement = iapiserver.ClusterNodeScaleRequirement{
		Name:             ClusterNodeScale4Xlarge,
		MinNodeNum:       251,
		MaxNodeNum:       500,
		LowestCpuCores:   16,
		LowestMemoryUnit: 60 * 1024 * 1024 * 1024,
	}
	x8LargeClusterNodeScaleRequirement = iapiserver.ClusterNodeScaleRequirement{
		Name:             ClusterNodeScale8Xlarge,
		MinNodeNum:       501,
		MaxNodeNum:       1000, //more than 500
		LowestCpuCores:   32,
		LowestMemoryUnit: 120 * 1024 * 1024 * 1024,
	}

	K8sClusterNodeScales = []iapiserver.ClusterNodeScaleRequirement{
		mediumClusterNodeScaleRequirement,
		largeClusterNodeScaleRequirement,
		xLargeClusterNodeScaleRequirement,
		x2LargeClusterNodeScaleRequirement,
		x4LargeClusterNodeScaleRequirement,
		x8LargeClusterNodeScaleRequirement,
	}
)

func (k *kubernetesService) ClusterNodeScaleRequirement(ctx context.Context) (*iapiserver.GetClusterNodeScaleRequirementResponse, error) {
	return &iapiserver.GetClusterNodeScaleRequirementResponse{
		Requirements: K8sClusterNodeScales,
	}, nil
}

func validateMasterFitClusterNodeScales(ip string, controlPlaneHostInfo *ikubeagent.HostInfo, nodeScaleNums int) error {
	var suitableRequirement *iapiserver.ClusterNodeScaleRequirement
	for i, v := range K8sClusterNodeScales {
		if nodeScaleNums >= v.MinNodeNum {
			if v.MaxNodeNum == -1 || i == len(K8sClusterNodeScales)-1 {
				suitableRequirement = &v
				break
			}
			if nodeScaleNums <= v.MaxNodeNum {
				suitableRequirement = &v
				break
			}
		}
	}

	if suitableRequirement == nil {
		return errors.Errorf("cannot find suitable node scale requirement, scales:%v", nodeScaleNums)
	}

	if controlPlaneHostInfo.Cpu.Cores < int64(suitableRequirement.LowestCpuCores) {
		return errors.Errorf("node %v,cpu cores does not fix node scale requirement, require :%v, current:%v",
			ip, suitableRequirement.LowestCpuCores, controlPlaneHostInfo.Cpu.Cores)
	}

	if controlPlaneHostInfo.Mem.Total < uint64(suitableRequirement.LowestMemoryUnit) {
		return errors.Errorf("node %v memory does not fix node scale requirement, require :%v, current:%v",
			ip, suitableRequirement.LowestMemoryUnit, controlPlaneHostInfo.Mem.Total)
	}
	return nil
}

// 检测当前控制面板节点的性能是否需要部署的集群的工作节点规模
func (k *kubernetesService) ValidateMasterFitClusterScales(ctx context.Context, req *iapiserver.KubeMasterNodeFitClusterScaleCheckRequest) error {
	wg := waitgroup.RunGenericConcurrently(ctx, req.MasterList, func(ctx context.Context, node *ikubeagent.KubernetesNodeConfig) waitgroup.GenericResult[*imachinery.Empty] {
		var stateRet ikubeagent.InstallStateResponse
		if _, err := libs.NewKubeagentClient(node.IP).GetInstallState(ctx, nil, &stateRet); err != nil {
			return waitgroup.NewGenericResult[*imachinery.Empty](nil, errors.WithStack(err))
		}
		err := validateMasterFitClusterNodeScales(node.IP, stateRet.HostStat, req.ClusterNodeScale)
		return waitgroup.NewGenericResult[*imachinery.Empty](nil, errors.WithStack(err))
	}, 5*time.Second)

	return wg.GetErrorList()
}
