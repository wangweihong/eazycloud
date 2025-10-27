package kubernetes

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/internal/pkg/libkubernetes"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/paging"
	"github.com/wangweihong/gotoolbox/pkg/sets"
)

func (k *kubernetesService) PrometheusRuleList(ctx context.Context, req *iapiserver.PrometheusRuleListRequest) (*iapiserver.PrometheusRuleListResponse, error) {
	resp := &iapiserver.PrometheusRuleListResponse{}

	if req.Namespace == "" {
		req.Namespace = iapiserver.NamespaceMonitoring
	}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ret, err := libkubernetes.PrometheusRuleList(ctx, cluster.Config, req.Namespace, req.ToListOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var resInfos []*iapiserver.PrometheusRuleInfo
	for i, dm := range ret.Items {
		if !sets.NewString(dm.Name).ContainAny(req.FuzzyFields()...) {
			continue
		}
		di := iapiserver.NewPrometheusRuleInfo(&ret.Items[i], cluster)
		resInfos = append(resInfos, di)
	}
	resp.List = resInfos
	s, e := paging.Index(len(resp.List), req.PageNum, req.PageSize)
	resp.List = resp.List[s:e]
	return resp, nil
}

func (k *kubernetesService) PrometheusRuleGet(ctx context.Context, req *iapiserver.PrometheusRuleGetRequest) (*iapiserver.PrometheusRuleInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.PrometheusRuleGet(ctx, cluster.Config, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewPrometheusRuleInfo(meta, cluster), nil
}

func (k *kubernetesService) PrometheusRuleCreate(ctx context.Context, req *iapiserver.PrometheusRuleRequest) (*iapiserver.PrometheusRuleInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.PrometheusRuleCreate(ctx, cluster.Config, req.PrometheusRule.Namespace, req.PrometheusRule, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewPrometheusRuleInfo(meta, cluster), nil
}

func (k *kubernetesService) PrometheusRuleUpdate(ctx context.Context, req *iapiserver.PrometheusRuleRequest) (*iapiserver.PrometheusRuleInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.PrometheusRuleUpdate(ctx, cluster.Config, req.PrometheusRule.Namespace, req.PrometheusRule, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewPrometheusRuleInfo(meta, cluster), nil
}

func (k *kubernetesService) PrometheusRuleDelete(ctx context.Context, req *iapiserver.PrometheusRuleRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = libkubernetes.PrometheusRuleDelete(ctx, cluster.Config, req.PrometheusRule.Namespace, req.PrometheusRule.Name, req.DeleteOpts); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
