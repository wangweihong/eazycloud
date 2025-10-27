package kubernetes

import (
	"context"

	v1 "k8s.io/api/core/v1"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/internal/pkg/clientset"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/waitgroup"
)

func (k *kubernetesService) ConfigMapCreate(ctx context.Context, req *iapiserver.ConfigMapRequest) (*iapiserver.ConfigMapInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.ConfigMapCreate(ctx, cluster, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sConfigMapToApiConfigMap(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) ConfigMapUpdate(ctx context.Context, req *iapiserver.ConfigMapRequest) (*iapiserver.ConfigMapInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.ConfigMapUpdate(ctx, cluster, req.Resource.Namespace, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sConfigMapToApiConfigMap(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) ConfigMapDelete(ctx context.Context, req *iapiserver.ConfigMapRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := clientset.ConfigMapDelete(ctx, cluster, req.Resource.Namespace, req.Resource, req.DeleteOpts); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) ConfigMapBatchDelete(ctx context.Context, req *iapiserver.ConfigMapBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.ConfigMapRequest] {
	wg := waitgroup.RunGenericConcurrently(ctx, req.Resources, func(ctx context.Context, res *iapiserver.ConfigMapRequest) waitgroup.GenericResult[*iapiserver.ConfigMapRequest] {
		cluster, err := k.store.Clusters().Get(ctx, res.Cluster)
		if err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		if err := clientset.ConfigMapDelete(ctx, cluster, res.Resource.Namespace, res.Resource, res.DeleteOpts); err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		return waitgroup.NewGenericResult(res, nil)
	})
	return wg.BatchGenericOutput()
}

func (k *kubernetesService) ConfigMapGet(ctx context.Context, req *iapiserver.ConfigMapGetRequest) (*iapiserver.ConfigMapInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.ConfigMapGet(ctx, cluster, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sConfigMapToApiConfigMap(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) ConfigMapList(ctx context.Context, req *iapiserver.ConfigMapListRequest) (*iapiserver.ConfigMapListResponse, error) {
	resp := &iapiserver.ConfigMapListResponse{}

	clusters, err := getVisitScope(ctx, k.store, req.ResourceListRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	wg := waitgroup.RunGenericConcurrently[*iapiserver.Cluster, iapiserver.EachResourceRangeListState[*iapiserver.ConfigMapInfo]](ctx, clusters, func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.ConfigMapInfo]] {
		clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.ConfigMapInfo](cluster.ID, cluster.Name)
		resList, err := clientset.ConfigMapList(ctx, cluster, req.Namespace, req.ToListOpts())
		if err != nil {
			return waitgroup.NewGenericResult(clusterListOne, err)
		}

		var resInfos []*iapiserver.ConfigMapInfo
		for i := range resList.Items {
			resInfo := &resList.Items[i]
			if NewObjectCommonFieldFilter(resInfo).Filter(req.Fuzzy) {
				continue
			}
			resInfos = append(resInfos, convertK8sConfigMapToApiConfigMap(resInfo, cluster, req.Yaml))
		}
		clusterListOne.TotalCount = len(resInfos)
		clusterListOne.List = resInfos
		return waitgroup.NewGenericResult(clusterListOne, err)
	})

	resp.EachRangeListState = wg.GetSuccessResultList()
	resp.TotalCount = CutPagingSliceResourceList[*iapiserver.ConfigMapInfo](resp.EachRangeListState, &resp.List, req.PageNum, req.PageSize, func(i, j int) bool {
		return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
	})
	return resp, nil
}

func convertK8sConfigMapToApiConfigMap(meta *v1.ConfigMap, cluster *iapiserver.Cluster, yaml bool) *iapiserver.ConfigMapInfo {
	resp := &iapiserver.ConfigMapInfo{
		Resource: meta,
	}
	if !yaml {

	}
	return resp
}
