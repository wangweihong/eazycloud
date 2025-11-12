package kubernetes

import (
	"context"
	"strconv"
	"strings"
	"time"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/internal/pkg/libkubernetes"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/json"
	"github.com/wangweihong/gotoolbox/pkg/paging"
	"github.com/wangweihong/gotoolbox/pkg/sets"
	"github.com/wangweihong/gotoolbox/pkg/sortutil"
	"github.com/wangweihong/gotoolbox/pkg/typeutil"
	"github.com/wangweihong/gotoolbox/pkg/waitgroup"

	"k8s.io/apimachinery/pkg/util/intstr"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ruleGroupPre              = "iapiserver.rulegroup.extra."
	DefaultPrometheusRuleName = "prometheus-k8s-rules"
	defaultPrometheusRuleKind = "PrometheusRule"
	defaultPrometheusRuleAv   = "monitoring.coreos.com/v1"
	defaultGroupRule          = "rulegroup_default_rule"
	ruleDescKeys              = []string{"message"}

	// alert
	alertLevelKey       = "severity"
	alertLevelValueDesc = map[string]string{
		"info":     "一般",
		"warning":  "告警",
		"critical": "严重",
	}
	proSqlOperator = map[string]string{
		"!=": "不等于",
		">=": "大于或等于",
		"<=": "小于或等于",
		"=":  "等于",
		">":  "大于",
		"<":  "小于",
	}
)

func (k *kubernetesService) BuiltinPrometheusRuleSpecGroupCreate(ctx context.Context, req *iapiserver.BuiltinPrometheusRuleSpecGroupRequest) (*iapiserver.PrometheusRuleInfo, error) {
	if len(req.GroupList[0].Spec.Rules) == 0 {
		req.GroupList[0].Spec.Rules = []monitoringv1.Rule{*k.getDefaultRule()}
	}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 获取内建默认的规则
	promRule, err := k.getDefaultPrometheusRuleResource(ctx, cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := k.groupAdd(promRule, req.GroupList); err != nil {
		return nil, errors.WithStack(err)
	}

	ret, err := libkubernetes.PrometheusRuleUpdate(ctx, cluster.Config, promRule.Namespace, promRule, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return iapiserver.NewPrometheusRuleInfo(ret, cluster), nil
}

func (k *kubernetesService) BuiltinPrometheusRuleSpecGroupUpdate(ctx context.Context, req *iapiserver.BuiltinPrometheusRuleSpecGroupRequest) (*iapiserver.PrometheusRuleInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	promRule, err := k.getDefaultPrometheusRuleResource(ctx, cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := k.groupUpdate(promRule, req.GroupList); err != nil {
		return nil, errors.WithStack(err)
	}

	promRule, err = libkubernetes.PrometheusRuleUpdate(ctx, cluster.Config, promRule.Namespace, promRule, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewPrometheusRuleInfo(promRule, cluster), nil
}

func (k *kubernetesService) BuiltinPrometheusRuleSpecGroupDelete(ctx context.Context, req *iapiserver.BuiltinPrometheusRuleSpecGroupRequest) (*iapiserver.PrometheusRuleInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	promRule, err := k.getDefaultPrometheusRuleResource(ctx, cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := k.groupDelete(promRule, req.GroupList); err != nil {
		return nil, errors.WithStack(err)
	}

	promRule, err = libkubernetes.PrometheusRuleUpdate(ctx, cluster.Config, promRule.Namespace, promRule, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewPrometheusRuleInfo(promRule, cluster), nil
}

func (k *kubernetesService) BuiltinPrometheusRuleSpecGroupGet(ctx context.Context, req *iapiserver.BuiltinPrometheusRuleSpecGroupGetRequest) (*iapiserver.PrometheusRuleInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	promRule, err := k.getDefaultPrometheusRuleResource(ctx, cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	groupList := k.group2Manager(k.filterGroup(promRule, true), cluster)

	for _, v := range groupList {
		if v.Name == req.GroupName {
			return iapiserver.NewPrometheusRuleInfo(promRule, cluster), nil
		}
	}
	return nil, errors.Errorf("rule group %v not found")

}

func (k *kubernetesService) BuiltinPrometheusRuleSpecGroupList(ctx context.Context, req *iapiserver.BuiltinPrometheusRuleSpecGroupListRequest) (*iapiserver.BuiltinPrometheusRuleSpecGroupListResponse, error) {
	resp := &iapiserver.BuiltinPrometheusRuleSpecGroupListResponse{}

	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList2[*iapiserver.PrometheusRuleSpecGroup](ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.PrometheusRuleSpecGroup]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.PrometheusRuleSpecGroup](cluster.ID, cluster.Name)
			defaultPromRule, err := k.getDefaultPrometheusRuleResource(ctx, cluster)
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}
			defaultPromRule = k.filterGroup(defaultPromRule, true)
			ruleGroupList := k.group2Manager(defaultPromRule, cluster)

			var rets []*iapiserver.PrometheusRuleSpecGroup
			for _, v := range ruleGroupList {
				if sets.NewString(v.Name, string(typeutil.GenericIndirectValue(v.Spec.Interval))).ContainAny(req.FuzzyFields()...) {
					continue
				}
				rets = append(rets, v)
			}

			clusterListOne.TotalCount = len(rets)
			clusterListOne.List = rets
			return waitgroup.NewGenericResult(clusterListOne, err)
		}, req.SortBy, !req.SortDesc, 10*time.Second)
	return resp, err
}

func (k *kubernetesService) BuiltinPrometheusRuleSpecGroupRuleCreate(ctx context.Context, req *iapiserver.BuiltinPrometheusRuleSpecGroupRuleRequest) (*iapiserver.PrometheusRuleInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	promRule, err := k.getDefaultPrometheusRuleResource(ctx, cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := k.prometheusRuleResourceAddSpecRules(promRule, req.GroupName, req.RuleList); err != nil {
		return nil, errors.WithStack(err)
	}

	promeRule, err := libkubernetes.PrometheusRuleUpdate(ctx, cluster.Config, promRule.Namespace, promRule, metav1.UpdateOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewPrometheusRuleInfo(promeRule, cluster), nil
}

func (k *kubernetesService) BuiltinPrometheusRuleSpecGroupRuleDelete(ctx context.Context, req *iapiserver.BuiltinPrometheusRuleSpecGroupRuleRequest) (*iapiserver.PrometheusRuleInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	promRule, err := k.getDefaultPrometheusRuleResource(ctx, cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := k.rulesDelete(promRule, req.GroupName, req.RuleList); err != nil {
		return nil, errors.WithStack(err)
	}

	promeRule, err := libkubernetes.PrometheusRuleUpdate(ctx, cluster.Config, promRule.Namespace, promRule, metav1.UpdateOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewPrometheusRuleInfo(promeRule, cluster), nil
}

func (k *kubernetesService) BuiltinPrometheusRuleSpecGroupRuleList(ctx context.Context, req *iapiserver.BuiltinPrometheusRuleSpecGroupRuleListRequest) (*iapiserver.BuiltinPrometheusRuleSpecGroupRuleListResponse, error) {
	resp := &iapiserver.BuiltinPrometheusRuleSpecGroupRuleListResponse{}

	if req.SortBy == "" {
		req.SortBy = "Spec.Alert"
	}
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	promRule, err := k.getDefaultPrometheusRuleResource(ctx, cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	promRule = k.filterGroup(promRule, true)
	ruleList, err := k.findRule(cluster, promRule, req.GroupName, req.RuleName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, v := range ruleList {
		if !sets.NewString(v.Name, v.Spec.Expr.String(), string(typeutil.GenericIndirectValue(v.Spec.For))).ContainAny(req.FuzzyFields()...) {
			continue
		}
		resp.List = append(resp.List, v)
	}

	sortutil.StructSliceSort(resp.List, req.SortBy, req.SortDesc)

	resp.TotalCount = len(resp.List)
	s, e := paging.Index(resp.TotalCount, req.PageNum, req.PageSize)
	resp.List = resp.List[s:e]
	return resp, nil
}

func (k *kubernetesService) BuiltinPrometheusRuleSpecGroupRuleGet(ctx context.Context, req *iapiserver.BuiltinPrometheusRuleSpecGroupRuleGetRequest) (*iapiserver.BuiltinPrometheusRuleSpecGroupRuleGetResponse, error) {
	resp := &iapiserver.BuiltinPrometheusRuleSpecGroupRuleGetResponse{}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	promRule, err := k.getDefaultPrometheusRuleResource(ctx, cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	promRule = k.filterGroup(promRule, true)

	ruleList, err := k.findRule(cluster, promRule, req.GroupName, req.RuleName)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp.List = ruleList
	resp.TotalCount = len(ruleList)

	return resp, nil
}

func (k *kubernetesService) BuiltinPrometheusRuleSpecRuleUpdate(ctx context.Context, req *iapiserver.BuiltinPrometheusRuleSpecGroupRuleRequest) (*iapiserver.PrometheusRuleInfo, error) {

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	promRule, err := k.getDefaultPrometheusRuleResource(ctx, cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := k.rulesUpdate(promRule, req.GroupName, req.RuleList); err != nil {
		return nil, errors.WithStack(err)
	}
	if _, err := libkubernetes.PrometheusRuleUpdate(ctx, cluster.Config, promRule.Namespace, promRule, metav1.UpdateOptions{}); err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewPrometheusRuleInfo(promRule, cluster), nil
}

func (k *kubernetesService) BuiltinPrometheusRuleSpecAlertCreate(ctx context.Context, req *iapiserver.PrometheusRuleRequest) (*iapiserver.PrometheusRuleInfo, error) {
	if req.PrometheusRule == nil || len(req.PrometheusRule.Spec.Groups) == 0 {
		return nil, errors.Errorf("alert is empty")
	}

	oldRuleId := make(map[string]struct{}, len(req.PrometheusRule.Spec.Groups)*5)
	for _, groupv := range req.PrometheusRule.Spec.Groups {
		for _, rulev := range groupv.Rules {
			if rulev.Alert == "" || rulev.For == nil || len(rulev.Labels) == 0 || len(rulev.Annotations) == 0 {
				return nil, errors.Errorf("alert for labels or annotations is empty")

			}

			if _, ok := oldRuleId[groupv.Name+rulev.Alert]; ok {
				return nil, errors.Errorf("group[%v] alert[%v] repeated!", groupv.Name, rulev.Alert)
			}
			oldRuleId[groupv.Name+rulev.Alert] = struct{}{}
		}
	}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	promeRule, err := libkubernetes.PrometheusRuleGet(ctx, cluster.Config, req.PrometheusRule.Namespace, DefaultPrometheusRuleName, req.GetOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for groupk, groupv := range promeRule.Spec.Groups {
		for _, rulev := range groupv.Rules {
			_, ok := oldRuleId[groupv.Name+rulev.Alert]
			if ok {
				return nil, errors.Errorf("group[%v] alert[%v] already exist!", groupv.Name, rulev.Alert)

			}
		}

		for _, reqGroupv := range req.PrometheusRule.Spec.Groups {
			if reqGroupv.Name == groupv.Name {
				groupv.Rules = append(groupv.Rules, reqGroupv.Rules...)
				promeRule.Spec.Groups[groupk] = groupv
				break
			}
		}
	}

	promeRule, err = libkubernetes.PrometheusRuleUpdate(ctx, cluster.Config, req.PrometheusRule.Namespace, promeRule, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewPrometheusRuleInfo(promeRule, cluster), nil
}

func (k *kubernetesService) PrometheusRuleSpecAlertUpdate(ctx context.Context, req *iapiserver.PrometheusRuleRequest) (*iapiserver.PrometheusRuleInfo, error) {
	if req.PrometheusRule == nil || len(req.PrometheusRule.Spec.Groups) == 0 || len(req.PrometheusRule.Spec.Groups[0].Rules) == 0 {
		return nil, errors.Errorf("alert is empty")
	}

	groupName := req.PrometheusRule.Spec.Groups[0].Name
	alertName := req.PrometheusRule.Spec.Groups[0].Rules[0].Alert
	targetRule := req.PrometheusRule.Spec.Groups[0].Rules[0]

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	promeRule, err := libkubernetes.PrometheusRuleGet(ctx, cluster.Config, req.PrometheusRule.Namespace, DefaultPrometheusRuleName, req.GetOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	isFind := false
	for _, groupv := range promeRule.Spec.Groups {
		if groupv.Name == groupName {
			for rulek, rulev := range groupv.Rules {
				if rulev.Alert == alertName {
					groupv.Rules[rulek] = targetRule
					isFind = true
					break
				}
			}
			break
		}
	}
	if !isFind {
		return nil, errors.Errorf("group[%v] rule[%v]", groupName, alertName)

	}

	promeRule, err = libkubernetes.PrometheusRuleUpdate(ctx, cluster.Config, req.PrometheusRule.Namespace, promeRule, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return iapiserver.NewPrometheusRuleInfo(promeRule, cluster), nil
}

func (k *kubernetesService) PrometheusRuleSpecAlertDelete(ctx context.Context, req *iapiserver.PrometheusRuleAlertRequest) error {
	if len(req.PrometheusRule.Spec.Groups) == 0 || len(req.PrometheusRule.Spec.Groups[0].Rules) == 0 {
		return errors.Errorf("missing alert group")
	}

	groupName := req.PrometheusRule.Spec.Groups[0].Name
	alertName := req.PrometheusRule.Spec.Groups[0].Rules[0].Alert
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	promeRule, err := libkubernetes.PrometheusRuleGet(ctx, cluster.Config, req.PrometheusRule.Namespace, DefaultPrometheusRuleName, req.GetOpts)
	if err != nil {
		return errors.WithStack(err)
	}

	isFind := false
	for _, groupv := range promeRule.Spec.Groups {
		if groupv.Name == groupName {
			for rulek, rulev := range groupv.Rules {
				if rulev.Alert == alertName {
					groupv.Rules = append(groupv.Rules[:rulek], groupv.Rules[rulek+1:]...)
					isFind = true
					break
				}
			}
			break
		}
	}
	if !isFind {

	}

	promeRule, err = libkubernetes.PrometheusRuleUpdate(ctx, cluster.Config, req.PrometheusRule.Namespace, promeRule, req.UpdateOpts)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) getDefaultPrometheusRuleResource(ctx context.Context, cluster *iapiserver.Cluster) (*monitoringv1.PrometheusRule, error) {
	promRule, err := libkubernetes.PrometheusRuleGet(ctx, cluster.Config, iapiserver.NamespaceMonitoring, DefaultPrometheusRuleName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return promRule, nil
}

func (k *kubernetesService) filterGroup(in *monitoringv1.PrometheusRule, showRules bool) *monitoringv1.PrometheusRule {
	if in == nil {
		return nil
	}

	resp := *in
	resp.Spec.Groups = nil

	findRecord := false
	findDefaultRule := false
	for _, gv := range in.Spec.Groups {

		gvt := gv
		gvt.Rules = make([]monitoringv1.Rule, 0, len(gv.Rules))
		for _, v := range gv.Rules {
			if isDefaultRule(&v) {
				findDefaultRule = true
			}
			if v.Record == "" && !isDefaultRule(&v) {
				v.Expr = intstr.FromString(strings.Replace(v.Expr.String(), "\n", "", -1))
				v.Expr = intstr.FromString(strings.Replace(v.Expr.String(), " ", "", -1))
				delete(v.Annotations, "runbook_url")
				gvt.Rules = append(gvt.Rules, v)
			} else {
				findRecord = true
			}
		}
		if findRecord && len(gvt.Rules) == 0 && !findDefaultRule {
			continue
		}
		if !showRules {
			gvt.Rules = nil
		}

		resp.Spec.Groups = append(resp.Spec.Groups, gvt)
	}

	return &resp
}

func (k *kubernetesService) groupAdd(promRule *monitoringv1.PrometheusRule, groupsAdd []*iapiserver.PrometheusRuleSpecGroup) error {
	if len(groupsAdd) == 0 {
		return nil
	}
	groupName := make(map[string]struct{}, len(groupsAdd))
	for _, v := range groupsAdd {
		if v.Name == "" || v.Spec.Interval == nil {
			return errors.Errorf("group name or group interval is empty")
		}
		_, ok := groupName[v.Name]
		if ok {
			return errors.Errorf("group name[%v] already exist", v.Name)
		}
		groupName[v.Name] = struct{}{}
		if v.Annotations == nil {
			v.Annotations = map[string]string{}
		}
	}
	for _, v := range promRule.Spec.Groups {
		if _, ok := groupName[v.Name]; ok {
			return errors.Errorf("group name[%v] already exist", v.Name)
		}
	}

	if promRule.Annotations == nil {
		promRule.Annotations = map[string]string{}
	}

	for _, v := range groupsAdd {
		if promRule == nil {
			continue
		}

		v.Spec.Name = v.Name
		promRule.Spec.Groups = append(promRule.Spec.Groups, v.Spec.RuleGroup)
	}
	return nil
}

func (k *kubernetesService) groupDelete(in *monitoringv1.PrometheusRule, groupsDelete []*iapiserver.PrometheusRuleSpecGroup) error {
	if len(groupsDelete) == 0 {
		return nil
	}
	groupNames := make(map[string]struct{}, len(in.Spec.Groups))
	deleteNames := make(map[string]struct{}, len(groupsDelete))
	for _, v := range in.Spec.Groups {
		groupNames[v.Name] = struct{}{}
	}

	for _, v := range groupsDelete { // check if name exist
		_, ok := groupNames[v.Name]
		if !ok {
			return errors.Errorf("group name[%v] not exist", v.Name)
		}
		deleteNames[v.Name] = struct{}{}
	}

	// delete
	tmpOld := in.Spec.Groups
	in.Spec.Groups = make([]monitoringv1.RuleGroup, len(in.Spec.Groups)-len(deleteNames))
	index := 0
	for _, v := range tmpOld {
		if _, ok := deleteNames[v.Name]; ok {
			delete(in.Annotations, ruleGroupPre+v.Name)
			continue
		}
		in.Spec.Groups[index] = v
		index++
	}
	return nil
}

func (k *kubernetesService) groupUpdate(in *monitoringv1.PrometheusRule, groupsUpdate []*iapiserver.PrometheusRuleSpecGroup) error {
	if len(groupsUpdate) == 0 {
		return nil
	}
	if in == nil {
		return errors.Errorf("prometheusRule is nil")
	}
	if in.Annotations == nil {
		in.Annotations = map[string]string{}
	}

	var group *monitoringv1.RuleGroup
	for k, v := range in.Spec.Groups {
		if v.Name == groupsUpdate[0].Name {
			group = &in.Spec.Groups[k]
			break
		}
	}
	if group == nil {
		return errors.Errorf("group[%v] not exist", groupsUpdate[0].Name)
	}

	extraStr := in.Annotations[ruleGroupPre+group.Name]
	extraMap := map[string]string{}
	if extraStr != "" {
		if err := json.Unmarshal([]byte(extraStr), &extraMap); err != nil {
			return errors.WithStack(err)
		}
	}
	extraByte, err := json.Marshal(extraMap)
	if err != nil {
		return errors.WithStack(err)
	}
	in.Annotations[ruleGroupPre+group.Name] = string(extraByte)

	//  change  interval
	group.Interval = groupsUpdate[0].Spec.Interval

	return nil
}

func (k *kubernetesService) prometheusRuleResourceAddSpecRules(in *monitoringv1.PrometheusRule, groupName string, rulesAdd []*iapiserver.PrometheusRule) error {
	if len(rulesAdd) == 0 {
		return nil
	}
	if groupName == "" {
		return errors.Errorf("group_name is empty")
	}

	// find group
	var group *monitoringv1.RuleGroup
	for k, v := range in.Spec.Groups {
		if v.Name == groupName {
			group = &in.Spec.Groups[k]
			break
		}
	}
	if group == nil {
		return errors.Errorf("group[%v] not exist", groupName)
	}

	// check if rule already exist
	for _, v := range group.Rules {
		if v.Alert == rulesAdd[0].Name {
			return errors.Errorf("rule[%v] already exist", v.Alert)
		}
	}

	// conver
	if rulesAdd[0].Spec.Labels == nil {
		rulesAdd[0].Spec.Labels = map[string]string{}
	}
	if rulesAdd[0].Spec.Annotations == nil {
		rulesAdd[0].Spec.Annotations = map[string]string{}
	}
	rulesAdd[0].Spec.Annotations[ruleDescKeys[0]] = rulesAdd[0].Spec.Desc
	rulesAdd[0].Spec.Labels[alertLevelKey] = rulesAdd[0].Spec.Level
	rulesAdd[0].Spec.Alert = rulesAdd[0].Name
	rulesAdd[0].Spec.Expr = intstr.FromString(rulesAdd[0].Spec.Expr.String() + rulesAdd[0].Spec.Operator + rulesAdd[0].Spec.Value)

	group.Rules = append(group.Rules, rulesAdd[0].Spec.Rule)

	return nil
}

func (k *kubernetesService) rulesDelete(in *monitoringv1.PrometheusRule, groupName string, rulesDelete []*iapiserver.PrometheusRule) error {
	if len(rulesDelete) == 0 {
		return nil
	}
	if groupName == "" {
		return errors.Errorf("group_name is empty")
	}
	deleteNames := map[string]struct{}{}
	for _, v := range rulesDelete {
		deleteNames[v.Name] = struct{}{}
	}

	var group *monitoringv1.RuleGroup
	for k, v := range in.Spec.Groups {
		if v.Name == groupName {
			group = &in.Spec.Groups[k]
			break
		}
	}
	if group == nil {
		return errors.Errorf("group[%v] not exist", groupName)
	}

	deleteIndex := map[int]struct{}{}
	tempRules := group.Rules
	for k, v := range group.Rules {
		_, ok := deleteNames[v.Alert]
		if ok {
			deleteIndex[k] = struct{}{}
		}
	}
	group.Rules = make([]monitoringv1.Rule, 0, len(tempRules)-len(deleteIndex))
	for k, v := range tempRules {
		_, ok := deleteIndex[k]
		if !ok {
			group.Rules = append(group.Rules, v)
		}
	}

	return nil
}

func (k *kubernetesService) rulesUpdate(in *monitoringv1.PrometheusRule, groupName string, rulesUpdate []*iapiserver.PrometheusRule) error {
	if len(rulesUpdate) == 0 {
		return nil
	}

	if err := k.rulesDelete(in, groupName, rulesUpdate); err != nil {
		return errors.WithStack(err)
	}

	if err := k.prometheusRuleResourceAddSpecRules(in, groupName, rulesUpdate); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) addRuleGroupInPrometheusRule(promRule *monitoringv1.PrometheusRule, group *iapiserver.PrometheusRuleSpecGroup) *monitoringv1.PrometheusRule {
	if promRule == nil || group == nil {
		return promRule
	}

	group.Spec.Name = group.Name
	if promRule.Annotations == nil {
		promRule.Annotations = map[string]string{}
	}

	promRule.Spec.Groups = append(promRule.Spec.Groups, group.Spec.RuleGroup)

	return promRule
}

func parseExpr(expr string) (string, string, string) {
	for k, v := range proSqlOperator {
		index := strings.LastIndex(expr, k)
		if index <= 0 {
			continue
		}

		metric := expr[:index]
		//operator := expr[index : index+len(k)]
		value := expr[index+len(k):]

		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			continue
		}

		return metric, v, value
	}

	return "", "", ""
}

func (k *kubernetesService) group2Manager(group *monitoringv1.PrometheusRule, cluster *iapiserver.Cluster) []*iapiserver.PrometheusRuleSpecGroup {
	if group == nil {
		return nil
	}

	resp := make([]*iapiserver.PrometheusRuleSpecGroup, len(group.Spec.Groups))
	for k, v := range group.Spec.Groups {
		g := &iapiserver.PrometheusRuleSpecGroup{}
		g.Spec.RuleGroup = v
		g.Spec.RuleCount = len(v.Rules)
		g.Spec.Rules = nil
		g.ObjectMeta = metav1.ObjectMeta{}
		g.Labels = map[string]string{}
		g.Annotations = map[string]string{}

		g.Name = v.Name
		g.Namespace = group.Namespace
		if g.Spec.Interval == nil {
			g.Spec.Interval = monitoringv1.DurationPointer("0s")
		}

		resp[k] = g
	}

	return resp
}

func (k *kubernetesService) getGroupExtraInfos(rule *monitoringv1.PrometheusRule, groupName string) map[string]string {
	if rule == nil || groupName == "" {
		return nil
	}

	groupExtraInfo := rule.Annotations[ruleGroupPre+groupName]
	if groupExtraInfo != "" {
		resp := map[string]string{}
		if err := json.Unmarshal([]byte(groupExtraInfo), &resp); err != nil {
			return nil
		}
		return resp
	}

	return nil
}

func (k *kubernetesService) findRule(cluster *iapiserver.Cluster, promRule *monitoringv1.PrometheusRule, groupName string, ruleName string) ([]*iapiserver.PrometheusRule, error) {
	var resp []*iapiserver.PrometheusRule
	var group *monitoringv1.RuleGroup

	if promRule == nil {
		return nil, errors.Errorf("promethues_rule is nil")
	}
	if groupName == "" {
		return nil, errors.Errorf("group_name[%v] is empty", groupName)
	}

	// find group
	for _, v := range promRule.Spec.Groups {
		if v.Name == groupName {
			group = &v
			break
		}
	}
	if group == nil {
		return nil, errors.Errorf("group_name[%v] is not exist", groupName)
	}

	// call back
	trans := func(v monitoringv1.Rule) iapiserver.PrometheusRule {
		r := iapiserver.PrometheusRule{}
		r.Name = v.Alert
		r.Namespace = iapiserver.NamespaceMonitoring
		r.Spec.Rule = v
		expr, operator, value := parseExpr(v.Expr.String())
		r.Spec.Rule.Expr = intstr.FromString(expr)
		r.Spec.Operator = operator
		r.Spec.Value = value
		r.Spec.Level = alertLevelValueDesc[v.Labels[alertLevelKey]]
		for _, dv := range ruleDescKeys {
			r.Spec.Desc = v.Annotations[dv]
			if r.Spec.Desc != "" {
				break
			}
		}
		return r
	}

	// find rule list
	if ruleName == "" {
		for _, v := range group.Rules {
			r := trans(v)
			resp = append(resp, &r)
		}
		return resp, nil
	} else {
		for _, v := range group.Rules {
			if v.Alert == ruleName {
				r := trans(v)
				resp = append(resp, &r)
				return resp, nil
			}
		}
		return nil, errors.Errorf("group_name[%v] rule_name[%v] not exist", groupName, ruleName)
	}
}

func (k *kubernetesService) getDefaultRule() *monitoringv1.Rule {
	resp := &monitoringv1.Rule{}
	resp.Alert = defaultGroupRule
	resp.Expr = intstr.FromString(defaultGroupRule)
	resp.For = monitoringv1.DurationPointer("10000000000s")
	resp.Labels = map[string]string{alertLevelKey: "info"}
	resp.Annotations = map[string]string{"message": "this is group default alert"}
	return resp
}

func isDefaultRule(rule *monitoringv1.Rule) bool {
	if rule != nil && rule.Alert == defaultGroupRule && rule.Expr.String() == defaultGroupRule {
		return true
	}
	return false
}
