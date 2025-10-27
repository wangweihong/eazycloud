package kubernetes

import (
	"context"
	"time"

	"github.com/robfig/cron"
	"github.com/wangweihong/gotoolbox/pkg/compareutil"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/randutil"
	"github.com/wangweihong/gotoolbox/pkg/waitgroup"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/internal/pkg/clientset"
)

func (k *kubernetesService) JobCreate(ctx context.Context, req *iapiserver.JobRequest) (*iapiserver.JobInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	meta, err := clientset.JobCreate(ctx, cluster, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, err
	}

	return convertK8sJobToApiJob(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) JobDelete(ctx context.Context, req *iapiserver.JobRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}
	if req.DeleteCollection {
		meta, err := clientset.JobGet(ctx, cluster, req.Resource.Namespace, req.Resource.Name, req.GetOpts)
		if err != nil {
			return errors.WithStack(err)
		}

		selector, err := metav1.LabelSelectorAsSelector(meta.Spec.Selector)
		if err != nil {
			return errors.WithStack(err)
		}

		if err := clientset.PodDeleteCollection(ctx, cluster, req.Resource.Namespace, metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: selector.String()}); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := clientset.JobDelete(ctx, cluster, req.Resource.Namespace, req.Resource, req.DeleteOpts); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) JobBatchDelete(ctx context.Context, req *iapiserver.JobBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.JobRequest] {
	wg := waitgroup.RunGenericConcurrently[*iapiserver.JobRequest, *iapiserver.JobRequest](ctx, req.Resources, func(ctx context.Context, res *iapiserver.JobRequest) waitgroup.GenericResult[*iapiserver.JobRequest] {
		cluster, err := k.store.Clusters().Get(ctx, res.Cluster)
		if err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		if err := clientset.JobDelete(ctx, cluster, res.Resource.Namespace, res.Resource, res.DeleteOpts); err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		return waitgroup.NewGenericResult(res, nil)
	})
	return wg.BatchGenericOutput()
}

func (k *kubernetesService) JobGet(ctx context.Context, req *iapiserver.JobGetRequest) (*iapiserver.JobResponse, error) {
	resp := &iapiserver.JobResponse{}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.JobGet(ctx, cluster, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp.Info = convertK8sJobToApiJob(meta, cluster, req.Yaml)
	return resp, nil
}

func (k *kubernetesService) JobList(ctx context.Context, req *iapiserver.JobListRequest) (*iapiserver.JobListResponse, error) {
	resp := &iapiserver.JobListResponse{}
	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList[*iapiserver.JobInfo](ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, c *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.JobInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.JobInfo](c.ID, c.Name)
			resList, err := clientset.JobList(ctx, c, req.Namespace, req.ToListOpts())
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.JobInfo
			for i := range resList.Items {
				resInfo := convertK8sJobToApiJob(&resList.Items[i], c, req.Yaml)
				resInfos = append(resInfos, resInfo)
			}
			clusterListOne.TotalCount = len(resInfos)
			clusterListOne.List = resInfos
			return waitgroup.NewGenericResult(clusterListOne, nil)
		}, func(i, j int) bool {

			switch req.SortBy {
			case iapiserver.KubernetesResourceJobSortByStartTime:
				var lastTime1, lastTime2 int64
				if resp.List[i].Resource.Status.StartTime != nil {
					lastTime1 = resp.List[i].Resource.Status.StartTime.Unix()
				}
				if resp.List[j].Resource.Status.StartTime != nil {
					lastTime2 = resp.List[j].Resource.Status.StartTime.Unix()
				}
				if lastTime1 != lastTime2 {
					return compareutil.Compare(lastTime1, lastTime2, req.SortDesc)
				}
			case iapiserver.KubernetesResourceJobSortByEndTime:
				var endTime1, endTime2 int64
				if resp.List[i].Resource.Status.CompletionTime != nil {
					endTime1 = resp.List[i].Resource.Status.CompletionTime.Unix()
				}
				if resp.List[j].Resource.Status.CompletionTime != nil {
					endTime2 = resp.List[j].Resource.Status.CompletionTime.Unix()
				}
				if endTime1 != endTime2 {
					return compareutil.Compare(endTime1, endTime2, req.SortDesc)
				}
			}

			return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
		}, 10*time.Second)
	return resp, err
}

func convertK8sJobToApiJob(meta *batchv1.Job, cluster *iapiserver.Cluster, yaml bool) *iapiserver.JobInfo {
	resp := &iapiserver.JobInfo{
		Resource: meta,
	}
	if !yaml {
		resp.ResourceConvert = make([]*iapiserver.ResourceConvert, 0)
		for _, container := range meta.Spec.Template.Spec.Containers {
			resp.ResourceConvert = append(resp.ResourceConvert, convertResourceLimitToPersistentUnit(container.Name, container.Resources.Requests, container.Resources.Limits))
		}
	}
	return resp
}

func (k *kubernetesService) CronJobCreate(ctx context.Context, req *iapiserver.CronJobRequest) (*iapiserver.CronJobInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if req.Resource.Spec.Schedule == "" {
		return nil, errors.Errorf("missing spec.schedule")
	}

	if _, err := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow).Parse(req.Resource.Spec.Schedule); err != nil {
		return nil, errors.Errorf("invalid cron schedule :%v", req.Resource.Spec.Schedule)
	}
	meta, err := clientset.CronJobCreate(ctx, cluster, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sCronJobToApiCronJob(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) CronJobUpdate(ctx context.Context, req *iapiserver.CronJobRequest) (*iapiserver.CronJobInfo, error) {
	if req.Resource.Spec.Schedule == "" {
		return nil, errors.Errorf("missing spec.schedule")
	}

	if _, err := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow).Parse(req.Resource.Spec.Schedule); err != nil {
		return nil, errors.Errorf("invalid cron schedule :%v", req.Resource.Spec.Schedule)
	}
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.CronJobUpdate(ctx, cluster, req.Resource.Namespace, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return convertK8sCronJobToApiCronJob(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) CronJobUpdateSuspend(ctx context.Context, req *iapiserver.CronJobRequest) (*iapiserver.CronJobInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.CronJobGet(ctx, cluster, req.Resource.Namespace, req.Resource.Name, req.GetOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	meta.Spec.Suspend = req.Resource.Spec.Suspend
	meta, err = clientset.CronJobUpdate(ctx, cluster, meta.Namespace, meta, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return convertK8sCronJobToApiCronJob(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) CronJobUpdateSchedule(ctx context.Context, req *iapiserver.CronJobRequest) (*iapiserver.CronJobInfo, error) {
	if req.Resource.Spec.Schedule == "" {
		return nil, errors.Errorf("missing spec.schedule")
	}

	if _, err := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow).Parse(req.Resource.Spec.Schedule); err != nil {
		return nil, errors.Errorf("invalid cron schedule :%v", req.Resource.Spec.Schedule)
	}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.CronJobGet(ctx, cluster, req.Resource.Namespace, req.Resource.Name, req.GetOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta.Spec.Schedule = req.Resource.Spec.Schedule
	meta, err = clientset.CronJobUpdate(ctx, cluster, meta.Namespace, meta, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return convertK8sCronJobToApiCronJob(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) CronJobDelete(ctx context.Context, req *iapiserver.CronJobRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := clientset.CronJobDelete(ctx, cluster, req.Resource.Namespace, req.Resource, req.DeleteOpts); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) CronJobBatchDelete(ctx context.Context, req *iapiserver.CronJobBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.CronJobRequest] {
	wg := waitgroup.RunGenericConcurrently[*iapiserver.CronJobRequest, *iapiserver.CronJobRequest](ctx, req.Resources, func(ctx context.Context, res *iapiserver.CronJobRequest) waitgroup.GenericResult[*iapiserver.CronJobRequest] {
		cluster, err := k.store.Clusters().Get(ctx, res.Cluster)
		if err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		if err := clientset.CronJobDelete(ctx, cluster, res.Resource.Namespace, res.Resource, res.DeleteOpts); err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		return waitgroup.NewGenericResult(res, nil)
	})
	return wg.BatchGenericOutput()
}

func (k *kubernetesService) CronJobGet(ctx context.Context, req *iapiserver.CronJobGetRequest) (*iapiserver.CronJobInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	meta, err := clientset.CronJobGet(ctx, cluster, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return convertK8sCronJobToApiCronJob(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) CronJobList(ctx context.Context, req *iapiserver.CronJobListRequest) (*iapiserver.CronJobListResponse, error) {
	resp := &iapiserver.CronJobListResponse{}
	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList[*iapiserver.CronJobInfo](ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, c *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.CronJobInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.CronJobInfo](c.ID, c.Name)
			resList, err := clientset.CronJobList(ctx, c, req.Namespace, req.ToListOpts())
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.CronJobInfo
			for i := range resList.Items {
				resInfo := convertK8sCronJobToApiCronJob(&resList.Items[i], c, req.Yaml)
				resInfos = append(resInfos, resInfo)
			}
			clusterListOne.TotalCount = len(resInfos)
			clusterListOne.List = resInfos
			return waitgroup.NewGenericResult(clusterListOne, nil)
		}, func(i, j int) bool {

			switch req.SortBy {
			case iapiserver.KubernetesResourceCronJobSortByLastJobTime:
				var lastTime1, lastTime2 int64
				if resp.List[i].Resource.Status.LastScheduleTime != nil {
					lastTime1 = resp.List[i].Resource.Status.LastScheduleTime.Unix()
				}
				if resp.List[j].Resource.Status.LastScheduleTime != nil {
					lastTime2 = resp.List[j].Resource.Status.LastScheduleTime.Unix()
				}
				if lastTime1 != lastTime2 {
					return compareutil.Compare(lastTime1, lastTime2, req.SortDesc)
				}
			}
			return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
		}, 10*time.Second)
	return resp, err
}

func (k *kubernetesService) CronJobTrigger(ctx context.Context, req *iapiserver.CronJobRequest) (*iapiserver.JobInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.CronJobGet(ctx, cluster, req.Resource.Namespace, req.Resource.Name, req.GetOpts)
	if err != nil {
		return nil, err
	}

	triggerJob := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: meta.Namespace,
			Name:      meta.Name + "-trigger-" + randutil.RandNumSets(8),
		},
		Spec: meta.Spec.JobTemplate.Spec,
	}

	job, err := clientset.JobCreate(ctx, cluster, req.Resource.Namespace, triggerJob, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return convertK8sJobToApiJob(job, cluster, req.Yaml), nil
}

func convertK8sCronJobToApiCronJob(meta *batchv1.CronJob, cluster *iapiserver.Cluster, yaml bool) *iapiserver.CronJobInfo {
	resp := &iapiserver.CronJobInfo{
		Resource: meta,
	}
	if !yaml {
		resp.ResourceConvert = make([]*iapiserver.ResourceConvert, 0)
		for _, container := range meta.Spec.JobTemplate.Spec.Template.Spec.Containers {
			resp.ResourceConvert = append(resp.ResourceConvert, convertResourceLimitToPersistentUnit(container.Name, container.Resources.Requests, container.Resources.Limits))
		}
	}
	return resp
}
