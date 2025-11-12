package kubernetes

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/internal/pkg/clientset"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/waitgroup"

	policyv1 "k8s.io/api/policy/v1"
)

func (k *kubernetesService) PodDisruptionBudgetCreate(ctx context.Context, req *iapiserver.PodDisruptionBudgetRequest) (*iapiserver.PodDisruptionBudgetInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.PodDisruptionBudgetCreate(ctx, cluster, req.Resource, metav1.CreateOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewPodDisruptionBudgetInfo(meta, cluster), nil
}

func (k *kubernetesService) PodDisruptionBudgetUpdate(ctx context.Context, req *iapiserver.PodDisruptionBudgetRequest) (*iapiserver.PodDisruptionBudgetInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.PodDisruptionBudgetUpdate(ctx, cluster, req.Resource, metav1.UpdateOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewPodDisruptionBudgetInfo(meta, cluster), nil
}

func (k *kubernetesService) PodDisruptionBudgetDelete(ctx context.Context, req *iapiserver.PodDisruptionBudgetRequest) (*iapiserver.PodDisruptionBudgetInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := clientset.PodDisruptionBudgetDelete(ctx, cluster, req.Resource, metav1.DeleteOptions{}); err != nil {
		return nil, errors.WithStack(err)
	}

	return nil, nil
}

func (k *kubernetesService) PodDisruptionBudgetBatchDelete(ctx context.Context, req *iapiserver.PodDisruptionBudgetBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.PodDisruptionBudgetRequest] {
	wg := waitgroup.RunGenericConcurrently(ctx, req.Resources, func(ctx context.Context, res *iapiserver.PodDisruptionBudgetRequest) waitgroup.GenericResult[*iapiserver.PodDisruptionBudgetRequest] {
		cluster, err := k.store.Clusters().Get(ctx, res.Cluster)
		if err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		if err := clientset.PodDisruptionBudgetDelete(ctx, cluster, res.Resource, res.DeleteOpts); err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		return waitgroup.NewGenericResult(res, nil)
	})
	return wg.BatchGenericOutput()
}

func (k *kubernetesService) PodDisruptionBudgetGet(ctx context.Context, req *iapiserver.PodDisruptionBudgetGetRequest) (*iapiserver.PodDisruptionBudgetInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.PodDisruptionBudgetGet(ctx, cluster, &policyv1.PodDisruptionBudget{ObjectMeta: metav1.ObjectMeta{Namespace: req.Namespace, Name: req.Name}}, metav1.GetOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewPodDisruptionBudgetInfo(meta, cluster), nil
}

func (k *kubernetesService) PodDisruptionBudgetList(ctx context.Context, req *iapiserver.PodDisruptionBudgetListRequest) (*iapiserver.PodDisruptionBudgetListResponse, error) {
	resp := &iapiserver.PodDisruptionBudgetListResponse{}
	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList[*iapiserver.PodDisruptionBudgetInfo](ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.PodDisruptionBudgetInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.PodDisruptionBudgetInfo](cluster.ID, cluster.Name)
			resList, err := clientset.PodDisruptionBudgetList(ctx, cluster, req.Namespace, metav1.ListOptions{})
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.PodDisruptionBudgetInfo
			for i := range resList.Items {
				resInfo := iapiserver.NewPodDisruptionBudgetInfo(&resList.Items[i], cluster)
				if NewObjectCommonFieldFilter(resInfo.Resource).Filter(req.Fuzzy) {
					continue
				}
				resInfos = append(resInfos, resInfo)
			}
			return waitgroup.NewGenericResult(clusterListOne, nil)
		}, func(i, j int) bool {
			return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
		}, 10*time.Second)
	return resp, errors.WithStack(err)
}
