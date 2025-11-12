package kubernetes

import (
	"context"
	"time"

	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/internal/pkg/clientset"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/waitgroup"
)

func (k *kubernetesService) NetworkPolicyAdd(ctx context.Context, req *iapiserver.NetworkPolicyRequest) (*iapiserver.NetworkPolicyInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.NetworkPolicyCreate(ctx, cluster, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewNetworkPolicyInfo(meta, cluster), nil
}

func (k *kubernetesService) NetworkPolicyUpdate(ctx context.Context, req *iapiserver.NetworkPolicyRequest) (*iapiserver.NetworkPolicyInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.NetworkPolicyUpdate(ctx, cluster, req.Resource.Namespace, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewNetworkPolicyInfo(meta, cluster), nil
}

func (k *kubernetesService) NetworkPolicyDelete(ctx context.Context, req *iapiserver.NetworkPolicyRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := clientset.NetworkPolicyDelete(ctx, cluster, req.Resource.Namespace, req.Resource, req.DeleteOpts); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) NetworkPolicyBatchDelete(ctx context.Context, req *iapiserver.NetworkPolicyBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.NetworkPolicyRequest] {
	wg := waitgroup.RunGenericConcurrently[*iapiserver.NetworkPolicyRequest, *iapiserver.NetworkPolicyRequest](
		ctx, req.Resources, func(ctx context.Context, res *iapiserver.NetworkPolicyRequest) waitgroup.GenericResult[*iapiserver.NetworkPolicyRequest] {
			cluster, err := k.store.Clusters().Get(ctx, res.Cluster)
			if err != nil {
				return waitgroup.NewGenericResult(res, err)
			}
			if err := clientset.NetworkPolicyDelete(ctx, cluster, res.Resource.Namespace, res.Resource, res.DeleteOpts); err != nil {
				return waitgroup.NewGenericResult(res, err)
			}
			return waitgroup.NewGenericResult(res, nil)
		})
	return wg.BatchGenericOutput()
}

func (k *kubernetesService) NetworkPolicyGet(ctx context.Context, req *iapiserver.NetworkPolicyGetRequest) (*iapiserver.NetworkPolicyInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.NetworkPolicyGet(ctx, cluster, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewNetworkPolicyInfo(meta, cluster), nil
}

func (k *kubernetesService) NetworkPolicyList(ctx context.Context, req *iapiserver.NetworkPolicyListRequest) (*iapiserver.NetworkPolicyListResponse, error) {
	resp := &iapiserver.NetworkPolicyListResponse{}

	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList[*iapiserver.NetworkPolicyInfo](ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, c *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.NetworkPolicyInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.NetworkPolicyInfo](c.ID, c.Name)
			resList, err := clientset.NetworkPolicyList(ctx, c, req.Namespace, req.ToListOpts())
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.NetworkPolicyInfo
			for i := range resList.Items {
				resInfo := iapiserver.NewNetworkPolicyInfo(&resList.Items[i], c)
				resInfos = append(resInfos, resInfo)
			}
			clusterListOne.TotalCount = len(resInfos)
			clusterListOne.List = resInfos
			return waitgroup.NewGenericResult(clusterListOne, nil)
		}, func(i, j int) bool {
			return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
		}, 10*time.Second)
	return resp, errors.WithStack(err)
}

func (k *kubernetesService) IngressAdd(ctx context.Context, req *iapiserver.IngressRequest) (*iapiserver.IngressInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.IngressCreate(ctx, cluster, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sIngressToApi(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) IngressUpdate(ctx context.Context, req *iapiserver.IngressRequest) (*iapiserver.IngressInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.IngressUpdate(ctx, cluster, req.Resource.Namespace, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sIngressToApi(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) IngressGet(ctx context.Context, req *iapiserver.IngressGetRequest) (*iapiserver.IngressInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.IngressGet(ctx, cluster, req.Namespace, req.Name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sIngressToApi(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) IngressDelete(ctx context.Context, req *iapiserver.IngressRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = clientset.IngressDelete(ctx, cluster, req.Resource.Namespace, req.Resource, req.DeleteOpts); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (k *kubernetesService) IngressBatchDelete(ctx context.Context, req *iapiserver.IngressBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.IngressRequest] {
	wg := waitgroup.RunGenericConcurrently[*iapiserver.IngressRequest, *iapiserver.IngressRequest](
		ctx, req.Resources, func(ctx context.Context, res *iapiserver.IngressRequest) waitgroup.GenericResult[*iapiserver.IngressRequest] {
			cluster, err := k.store.Clusters().Get(ctx, res.Cluster)
			if err != nil {
				return waitgroup.NewGenericResult(res, err)
			}
			if err := clientset.IngressDelete(ctx, cluster, res.Resource.Namespace, res.Resource, res.DeleteOpts); err != nil {
				return waitgroup.NewGenericResult(res, err)
			}
			return waitgroup.NewGenericResult(res, nil)
		})
	return wg.BatchGenericOutput()
}

func (k *kubernetesService) IngressList(ctx context.Context, req *iapiserver.IngressListRequest) (*iapiserver.IngressListResponse, error) {
	resp := &iapiserver.IngressListResponse{}

	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList(ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, c *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.IngressInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.IngressInfo](c.ID, c.Name)
			resList, err := clientset.IngressList(ctx, c, req.Namespace, req.ToListOpts())
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.IngressInfo
			for i := range resList.Items {
				resInfo := convertK8sIngressToApi(&resList.Items[i], c, req.Yaml)
				resInfos = append(resInfos, resInfo)
			}
			clusterListOne.TotalCount = len(resInfos)
			clusterListOne.List = resInfos
			return waitgroup.NewGenericResult(clusterListOne, nil)
		}, func(i, j int) bool {
			return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
		}, 10*time.Second)
	return resp, errors.WithStack(err)
}

func convertK8sIngressToApi(meta *networkingv1.Ingress, cluster *iapiserver.Cluster, yaml bool) *iapiserver.IngressInfo {
	apiInfo := iapiserver.NewIngressInfo(meta, cluster)

	if !yaml {
		ref := v1.ObjectReference{
			Kind:            "Deployment",
			Namespace:       "system-router",
			Name:            "router-" + meta.Namespace,
			UID:             "",
			APIVersion:      "apps/v1",
			ResourceVersion: "",
			FieldPath:       "",
		}
		apiInfo.Controller = ref
		routerDeployment, err := clientset.DeploymentGet(context.Background(), cluster, ref.Namespace, ref.Name, metav1.GetOptions{})
		if err == nil {
			if routerDeployment.Spec.Replicas != nil && *routerDeployment.Spec.Replicas > 0 && *routerDeployment.Spec.Replicas == routerDeployment.Status.Replicas {
				apiInfo.ControllerHealthy = true
			}
		}
	}

	return apiInfo
}
