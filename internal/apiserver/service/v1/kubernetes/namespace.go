package kubernetes

import (
	"context"
	"strings"
	"time"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/internal/pkg/clientset"
	"github.com/wangweihong/gotoolbox/pkg/async"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/log"
	"github.com/wangweihong/gotoolbox/pkg/waitgroup"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *kubernetesService) NamespaceAdd(ctx context.Context, req *iapiserver.NamespaceRequest) (*iapiserver.NamespaceInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.NamespaceCreate(ctx, cluster, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewNamespaceInfo(meta, cluster), nil
}

func (k *kubernetesService) NamespaceUpdate(ctx context.Context, req *iapiserver.NamespaceRequest) (*iapiserver.NamespaceInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.NamespaceUpdate(ctx, cluster, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewNamespaceInfo(meta, cluster), nil
}

func (k *kubernetesService) NamespaceDelete(ctx context.Context, req *iapiserver.NamespaceRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	// 限制内置命名空间删除
	if iapiserver.KubernetesBuiltInNamespace.Has(req.Resource.Name) {
		return errors.WithStack(err)
	}

	if err := clientset.NamespaceDelete(ctx, cluster, req.Resource, req.DeleteOpts); err != nil {
		return errors.WithStack(err)
	}
	async.Run(ctx, func(ctx context.Context) {
		tmpReq := &iapiserver.RouterRequest{}
		tmpReq.Cluster = req.Cluster
		tmpReq.Namespace = req.Resource.Name

		if err := k.RouterDelete(ctx, tmpReq); err != nil {
			log.Errorf("delete namespace router when delete namespace %v fail:%v", req.Resource.Name, err)
		}
	})

	return nil
}

func (k *kubernetesService) NamespaceGet(ctx context.Context, req *iapiserver.NamespaceGetRequest) (*iapiserver.NamespaceInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.NamespaceGet(ctx, cluster, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return iapiserver.NewNamespaceInfo(meta, cluster), nil
}

func (k *kubernetesService) NamespaceGatewayGet(ctx context.Context, req *iapiserver.PodListRequest) (*iapiserver.NamespaceGatewayResponse, error) {
	resp := &iapiserver.NamespaceGatewayResponse{}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ret, err := clientset.PodList(ctx, cluster, iapiserver.IngressControllerNamespace, req.ToListOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, v := range ret.Items {
		if strings.Contains(v.Name, "router-"+req.Namespace+"-") {
			resp.Gateway = v.Status.HostIP
		}
	}

	return resp, nil
}

func (k *kubernetesService) isNamespaceGatewayEnable(infos *iapiserver.RouterListResponse, clusterUuid string, namespace string) bool {
	if infos != nil || clusterUuid != "" || namespace != "" {
		for _, v := range infos.List {
			if v.Resource == nil || v.Cluster.ID != clusterUuid || v.Resource.Name != iapiserver.IngressControllerPrefix+namespace {
				continue
			}
			return true
		}
	}

	return false
}

func (k *kubernetesService) NamespaceList(ctx context.Context, req *iapiserver.NamespaceListRequest) (*iapiserver.NamespaceListResponse, error) {
	resp := &iapiserver.NamespaceListResponse{}
	var err error
	routerList, err := k.RouterList(ctx, &iapiserver.RouterListRequest{ResourceListRequest: req.ResourceListRequest})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList[*iapiserver.NamespaceInfo](ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.NamespaceInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.NamespaceInfo](cluster.ID, cluster.Name)
			resList, err := clientset.NamespaceList(ctx, cluster, metav1.ListOptions{})
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.NamespaceInfo
			for i := range resList.Items {
				resInfo := iapiserver.NewNamespaceInfo(&resList.Items[i], cluster)
				if NewObjectCommonFieldFilter(resInfo.Resource).AddField(string(resInfo.Resource.Status.Phase)).Filter(req.Fuzzy) {
					continue
				}

				gateWayEnable := k.isNamespaceGatewayEnable(routerList, cluster.ID, resInfo.Resource.Name)
				resInfo.GatewayEnable = &gateWayEnable
			}
			clusterListOne.TotalCount = len(resInfos)
			clusterListOne.List = resInfos
			return waitgroup.NewGenericResult(clusterListOne, nil)
		}, func(i, j int) bool {
			return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
		}, 10*time.Second)
	return resp, errors.WithStack(err)
}

// func (k *kubernetesService)  NamespaceTree(ctx context.Context, req *iapiserver.NamespaceListRequest) *iapiserver.NamespaceTreeResponse {
// 	resp := &iapiserver.NamespaceTreeResponse{}
// 	clusters, err := tm.GetVisitScope(req.ResourceListRequest)
// 	if err != nil {

// 		return resp
// 	}

// 	wg := utils.NewWaitGroup(nil)
// 	for _, cluster := range clusters {
// 		cluster := cluster
// 		wg.Start(utils.NewWaitGroupHandleFunc("", nil, func() utils.WaitGroupResult {
// 			clusterListOne := iapiserver.EachResourceRangeListState{}
// 			clusterListOne.ClusterName = cluster.Name
// 			clusterListOne.ClusterUUID = cluster.UUID
// 			resList, err := clientset.NamespaceList(cluster, req.ToListOpts())
// 			if err != nil {
// 				return utils.NewWaitGroupResult(clusterListOne, status.UpdateStatus(err))
// 			}

// 			var resInfos []*iapiserver.NamespaceInfo
// 			for i := range resList.Items {
// 				resInfo := convertK8sNamespaceToApiNamespace(&resList.Items[i], cluster, req.Yaml)
// 				if NewObjectCommonFieldFilter(resInfo).AddField(string(resInfo.Status.Phase)).Filter(req.Fuzzy) {
// 					continue
// 				}
// 				resInfos = append(resInfos, resInfo)
// 			}

// 			clusterNamespaceInfo := &iapiserver.ClusterNamespaceInfo{}
// 			clusterNamespaceInfo.Cluster = convertClusterToClusterInfo(cluster)
// 			clusterNamespaceInfo.Namespaces = resInfos

// 			clusterListOne.TotalCount = 1
// 			clusterListOne.List = clusterNamespaceInfo
// 			return utils.NewWaitGroupResult(clusterListOne, nil)
// 		}))
// 	}
// 	wg.Wait()

// 	resp.EachRangeListState = make([]iapiserver.EachResourceRangeListState, 0)
// 	for _, ret := range wg.GetResults() {
// 		oneClusterResourceList, ok := ret.Data.(iapiserver.EachResourceRangeListState)
// 		if !ok {
// 			logrus.Error("ignore for non iapiserver.EachResourceRangeListState type")
// 			return resp
// 		}
// 		oneClusterResourceList.Result = iapiserver.SetOutput(nil, ret.Error)
// 		resList := oneClusterResourceList.List
// 		oneClusterResourceList.List = nil

// 		if resList != nil {
// 			resp.List = append(resp.List, resList.(*iapiserver.ClusterNamespaceInfo))
// 		}
// 		resp.EachRangeListState = append(resp.EachRangeListState, oneClusterResourceList)
// 	}
// 	return resp
// }
