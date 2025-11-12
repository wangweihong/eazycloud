package kubernetes

import (
	"github.com/gin-gonic/gin"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

func (rc *KubernetesController) PodList(c *gin.Context) {
	runKubernetes(c, &iapiserver.PodListRequest{}, func(r *iapiserver.PodListRequest) (any, error) {
		return rc.srv.Kubernetes().PodList(c, r)
	})
}

func (rc *KubernetesController) PodGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.ResourceGetRequest{}, func(r *iapiserver.ResourceGetRequest) (any, error) {
		return rc.srv.Kubernetes().PodGet(c, r)
	})
}

func (rc *KubernetesController) PodAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.PodRequest{}, func(r *iapiserver.PodRequest) (any, error) {
		return rc.srv.Kubernetes().PodAdd(c, r)
	})
}

func (rc *KubernetesController) PodDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.PodRequest{}, func(r *iapiserver.PodRequest) (any, error) {
		err := rc.srv.Kubernetes().PodDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) PodBatchDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.PodBatchRequest{}, func(r *iapiserver.PodBatchRequest) (any, error) {
		ret := rc.srv.Kubernetes().PodBatchDelete(c, r)
		return ret, nil
	})
}

func (rc *KubernetesController) PodLogList(c *gin.Context) {
	runKubernetes(c, &iapiserver.PodLogRequest{}, func(r *iapiserver.PodLogRequest) (any, error) {
		return rc.srv.Kubernetes().PodLogList(c, r)
	})
}

func (rc *KubernetesController) PodLogStream(c *gin.Context) {
	runKubernetes(c, &iapiserver.PodLogRequest{}, func(r *iapiserver.PodLogRequest) (any, error) {
		return rc.srv.Kubernetes().PodLogStream(c, r)
	})
}

func (rc *KubernetesController) PodPatch(c *gin.Context) {
	runKubernetes(c, &iapiserver.PodRequest{}, func(r *iapiserver.PodRequest) (any, error) {
		return rc.srv.Kubernetes().PodPatch(c, r)
	})
}

func (rc *KubernetesController) PodEvict(c *gin.Context) {
	runKubernetes(c, &iapiserver.PodRequest{}, func(r *iapiserver.PodRequest) (any, error) {
		err := rc.srv.Kubernetes().PodEvict(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) GetComponentPod(c *gin.Context) {
	runKubernetes(c, &iapiserver.PodListRequest{}, func(r *iapiserver.PodListRequest) (any, error) {
		return rc.srv.Kubernetes().GetComponentPod(c, r)
	})
}

func (rc *KubernetesController) ServiceGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.ServiceGetRequest{}, func(r *iapiserver.ServiceGetRequest) (any, error) {
		return rc.srv.Kubernetes().ServiceGet(c, r)
	})
}

func (rc *KubernetesController) ServiceList(c *gin.Context) {
	runKubernetes(c, &iapiserver.ServiceListRequest{}, func(r *iapiserver.ServiceListRequest) (any, error) {
		return rc.srv.Kubernetes().ServiceList(c, r)
	})
}

func (rc *KubernetesController) ServiceAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.ServiceRequest{}, func(r *iapiserver.ServiceRequest) (any, error) {
		return rc.srv.Kubernetes().ServiceAdd(c, r)
	})
}

func (rc *KubernetesController) ServiceDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.ServiceRequest{}, func(r *iapiserver.ServiceRequest) (any, error) {
		err := rc.srv.Kubernetes().ServiceDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) ServiceBatchDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.ServiceBatchRequest{}, func(r *iapiserver.ServiceBatchRequest) (any, error) {
		ret := rc.srv.Kubernetes().ServiceBatchDelete(c, r)
		return ret, nil
	})
}

func (rc *KubernetesController) ServiceUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.ServiceRequest{}, func(r *iapiserver.ServiceRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().ServiceUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) ConfigMapGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.ConfigMapGetRequest{}, func(r *iapiserver.ConfigMapGetRequest) (any, error) {
		return rc.srv.Kubernetes().ConfigMapGet(c, r)
	})
}

func (rc *KubernetesController) ConfigMapList(c *gin.Context) {
	runKubernetes(c, &iapiserver.ConfigMapListRequest{}, func(r *iapiserver.ConfigMapListRequest) (any, error) {
		return rc.srv.Kubernetes().ConfigMapList(c, r)
	})
}

func (rc *KubernetesController) ConfigMapAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.ConfigMapRequest{}, func(r *iapiserver.ConfigMapRequest) (any, error) {
		return rc.srv.Kubernetes().ConfigMapAdd(c, r)
	})
}

func (rc *KubernetesController) ConfigMapDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.ConfigMapRequest{}, func(r *iapiserver.ConfigMapRequest) (any, error) {
		err := rc.srv.Kubernetes().ConfigMapDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) ConfigMapBatchDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.ConfigMapBatchRequest{}, func(r *iapiserver.ConfigMapBatchRequest) (any, error) {
		ret := rc.srv.Kubernetes().ConfigMapBatchDelete(c, r)
		return ret, nil
	})
}

func (rc *KubernetesController) ConfigMapUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.ConfigMapRequest{}, func(r *iapiserver.ConfigMapRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().ConfigMapUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) LimitRangeGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.LimitRangeGetRequest{}, func(r *iapiserver.LimitRangeGetRequest) (any, error) {
		return rc.srv.Kubernetes().LimitRangeGet(c, r)
	})
}

func (rc *KubernetesController) LimitRangeList(c *gin.Context) {
	runKubernetes(c, &iapiserver.LimitRangeListRequest{}, func(r *iapiserver.LimitRangeListRequest) (any, error) {
		return rc.srv.Kubernetes().LimitRangeList(c, r)
	})
}

func (rc *KubernetesController) LimitRangeAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.LimitRangeRequest{}, func(r *iapiserver.LimitRangeRequest) (any, error) {
		return rc.srv.Kubernetes().LimitRangeAdd(c, r)
	})
}

func (rc *KubernetesController) LimitRangeDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.LimitRangeRequest{}, func(r *iapiserver.LimitRangeRequest) (any, error) {
		err := rc.srv.Kubernetes().LimitRangeDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) LimitRangeUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.LimitRangeRequest{}, func(r *iapiserver.LimitRangeRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().LimitRangeUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) ResourceQuotaGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.ResourceQuotaGetRequest{}, func(r *iapiserver.ResourceQuotaGetRequest) (any, error) {
		return rc.srv.Kubernetes().ResourceQuotaGet(c, r)
	})
}

func (rc *KubernetesController) ResourceQuotaList(c *gin.Context) {
	runKubernetes(c, &iapiserver.ResourceQuotaListRequest{}, func(r *iapiserver.ResourceQuotaListRequest) (any, error) {
		return rc.srv.Kubernetes().ResourceQuotaList(c, r)
	})
}

func (rc *KubernetesController) ResourceQuotaAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.ResourceQuotaRequest{}, func(r *iapiserver.ResourceQuotaRequest) (any, error) {
		return rc.srv.Kubernetes().ResourceQuotaAdd(c, r)
	})
}

func (rc *KubernetesController) ResourceQuotaDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.ResourceQuotaRequest{}, func(r *iapiserver.ResourceQuotaRequest) (any, error) {
		err := rc.srv.Kubernetes().ResourceQuotaDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) ResourceQuotaUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.ResourceQuotaRequest{}, func(r *iapiserver.ResourceQuotaRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().ResourceQuotaUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) SecretGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.SecretGetRequest{}, func(r *iapiserver.SecretGetRequest) (any, error) {
		return rc.srv.Kubernetes().SecretGet(c, r)
	})
}

func (rc *KubernetesController) SecretList(c *gin.Context) {
	runKubernetes(c, &iapiserver.SecretListRequest{}, func(r *iapiserver.SecretListRequest) (any, error) {
		return rc.srv.Kubernetes().SecretList(c, r)
	})
}

func (rc *KubernetesController) SecretAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.SecretRequest{}, func(r *iapiserver.SecretRequest) (any, error) {
		return rc.srv.Kubernetes().SecretAdd(c, r)
	})
}

func (rc *KubernetesController) SecretDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.SecretRequest{}, func(r *iapiserver.SecretRequest) (any, error) {
		err := rc.srv.Kubernetes().SecretDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) SecretBatchDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.SecretBatchRequest{}, func(r *iapiserver.SecretBatchRequest) (any, error) {
		ret := rc.srv.Kubernetes().SecretBatchDelete(c, r)
		return ret, nil
	})
}

func (rc *KubernetesController) SecretUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.SecretRequest{}, func(r *iapiserver.SecretRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().SecretUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) NodeGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.NodeGetRequest{}, func(r *iapiserver.NodeGetRequest) (any, error) {
		return rc.srv.Kubernetes().NodeGet(c, r)
	})
}

func (rc *KubernetesController) NodeList(c *gin.Context) {
	runKubernetes(c, &iapiserver.NodeListRequest{}, func(r *iapiserver.NodeListRequest) (any, error) {
		return rc.srv.Kubernetes().NodeList(c, r)
	})
}

func (rc *KubernetesController) NodeAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.NodeRequest{}, func(r *iapiserver.NodeRequest) (any, error) {
		return rc.srv.Kubernetes().NodeAdd(c, r)
	})
}

func (rc *KubernetesController) NodeDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.NodeRequest{}, func(r *iapiserver.NodeRequest) (any, error) {
		err := rc.srv.Kubernetes().NodeDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) NodeUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.NodeRequest{}, func(r *iapiserver.NodeRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().NodeUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) NodeCordon(c *gin.Context) {
	runKubernetes(c, &iapiserver.NodeRequest{}, func(r *iapiserver.NodeRequest) (any, error) {
		err := rc.srv.Kubernetes().NodeCordon(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) NodeUncordon(c *gin.Context) {
	runKubernetes(c, &iapiserver.NodeRequest{}, func(r *iapiserver.NodeRequest) (any, error) {
		err := rc.srv.Kubernetes().NodeUncordon(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) NodeDrain(c *gin.Context) {
	runKubernetes(c, &iapiserver.NodeDrainRequest{}, func(r *iapiserver.NodeDrainRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().NodeDrain(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) NodeGatewayUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.NodeRequest{}, func(r *iapiserver.NodeRequest) (any, error) {
		err := rc.srv.Kubernetes().NodeGatewayUpdate(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) NodeTaintUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.NodeRequest{}, func(r *iapiserver.NodeRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().NodeTaintUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) NodeUpdateLabel(c *gin.Context) {
	runKubernetes(c, &iapiserver.NodeRequest{}, func(r *iapiserver.NodeRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().NodeUpdateLabel(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) NodeEvent(c *gin.Context) {
	runKubernetes(c, &iapiserver.NodeEventRequest{}, func(r *iapiserver.NodeEventRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().NodeEvent(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) EventGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.EventGetRequest{}, func(r *iapiserver.EventGetRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().EventGet(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) EventList(c *gin.Context) {
	runKubernetes(c, &iapiserver.EventListRequest{}, func(r *iapiserver.EventListRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().EventList(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) JobGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.JobGetRequest{}, func(r *iapiserver.JobGetRequest) (any, error) {
		return rc.srv.Kubernetes().JobGet(c, r)
	})
}

func (rc *KubernetesController) JobList(c *gin.Context) {
	runKubernetes(c, &iapiserver.JobListRequest{}, func(r *iapiserver.JobListRequest) (any, error) {
		return rc.srv.Kubernetes().JobList(c, r)
	})
}

func (rc *KubernetesController) JobAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.JobRequest{}, func(r *iapiserver.JobRequest) (any, error) {
		return rc.srv.Kubernetes().JobAdd(c, r)
	})
}

func (rc *KubernetesController) JobDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.JobRequest{}, func(r *iapiserver.JobRequest) (any, error) {
		err := rc.srv.Kubernetes().JobDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) JobBatchDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.JobBatchRequest{}, func(r *iapiserver.JobBatchRequest) (any, error) {
		ret := rc.srv.Kubernetes().JobBatchDelete(c, r)
		return ret, nil
	})
}

func (rc *KubernetesController) CronJobGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.CronJobGetRequest{}, func(r *iapiserver.CronJobGetRequest) (any, error) {
		return rc.srv.Kubernetes().CronJobGet(c, r)
	})
}

func (rc *KubernetesController) CronJobList(c *gin.Context) {
	runKubernetes(c, &iapiserver.CronJobListRequest{}, func(r *iapiserver.CronJobListRequest) (any, error) {
		return rc.srv.Kubernetes().CronJobList(c, r)
	})
}

func (rc *KubernetesController) CronJobAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.CronJobRequest{}, func(r *iapiserver.CronJobRequest) (any, error) {
		return rc.srv.Kubernetes().CronJobAdd(c, r)
	})
}

func (rc *KubernetesController) CronJobDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.CronJobRequest{}, func(r *iapiserver.CronJobRequest) (any, error) {
		err := rc.srv.Kubernetes().CronJobDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) CronJobBatchDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.CronJobBatchRequest{}, func(r *iapiserver.CronJobBatchRequest) (any, error) {
		ret := rc.srv.Kubernetes().CronJobBatchDelete(c, r)
		return ret, nil
	})
}

func (rc *KubernetesController) CronJobUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.CronJobRequest{}, func(r *iapiserver.CronJobRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().CronJobUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) CronJobUpdateSuspend(c *gin.Context) {
	runKubernetes(c, &iapiserver.CronJobRequest{}, func(r *iapiserver.CronJobRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().CronJobUpdateSuspend(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) CronJobUpdateSchedule(c *gin.Context) {
	runKubernetes(c, &iapiserver.CronJobRequest{}, func(r *iapiserver.CronJobRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().CronJobUpdateSchedule(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) CronJobTrigger(c *gin.Context) {
	runKubernetes(c, &iapiserver.CronJobRequest{}, func(r *iapiserver.CronJobRequest) (any, error) {
		return rc.srv.Kubernetes().CronJobTrigger(c, r)
	})
}

func (rc *KubernetesController) DeploymentGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.DeploymentGetRequest{}, func(r *iapiserver.DeploymentGetRequest) (any, error) {
		return rc.srv.Kubernetes().DeploymentGet(c, r)
	})
}

func (rc *KubernetesController) DeploymentList(c *gin.Context) {
	runKubernetes(c, &iapiserver.DeploymentListRequest{}, func(r *iapiserver.DeploymentListRequest) (any, error) {
		return rc.srv.Kubernetes().DeploymentList(c, r)
	})
}

func (rc *KubernetesController) DeploymentAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.DeploymentRequest{}, func(r *iapiserver.DeploymentRequest) (any, error) {
		return rc.srv.Kubernetes().DeploymentAdd(c, r)
	})
}

func (rc *KubernetesController) DeploymentRecreate(c *gin.Context) {
	runKubernetes(c, &iapiserver.DeploymentRequest{}, func(r *iapiserver.DeploymentRequest) (any, error) {
		return rc.srv.Kubernetes().DeploymentRecreate(c, r)
	})
}

func (rc *KubernetesController) DeploymentDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.DeploymentRequest{}, func(r *iapiserver.DeploymentRequest) (any, error) {
		err := rc.srv.Kubernetes().DeploymentDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) DeploymentBatchDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.DeploymentBatchRequest{}, func(r *iapiserver.DeploymentBatchRequest) (any, error) {
		ret := rc.srv.Kubernetes().DeploymentBatchDelete(c, r)
		return ret, nil
	})
}

func (rc *KubernetesController) DeploymentUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.DeploymentRequest{}, func(r *iapiserver.DeploymentRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().DeploymentUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) DeploymentVersionList(c *gin.Context) {
	runKubernetes(c, &iapiserver.DeploymentVersionListRequest{}, func(r *iapiserver.DeploymentVersionListRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().DeploymentVersionList(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) DeploymentVersionUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.DeploymentRequest{}, func(r *iapiserver.DeploymentRequest) (any, error) {
		err := rc.srv.Kubernetes().DeploymentVersionUpdate(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) DaemonSetGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.DaemonSetGetRequest{}, func(r *iapiserver.DaemonSetGetRequest) (any, error) {
		return rc.srv.Kubernetes().DaemonSetGet(c, r)
	})
}

func (rc *KubernetesController) DaemonSetList(c *gin.Context) {
	runKubernetes(c, &iapiserver.DaemonSetListRequest{}, func(r *iapiserver.DaemonSetListRequest) (any, error) {
		return rc.srv.Kubernetes().DaemonSetList(c, r)
	})
}

func (rc *KubernetesController) DaemonSetAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.DaemonSetRequest{}, func(r *iapiserver.DaemonSetRequest) (any, error) {
		return rc.srv.Kubernetes().DaemonSetAdd(c, r)
	})
}

func (rc *KubernetesController) DaemonSetRecreate(c *gin.Context) {
	runKubernetes(c, &iapiserver.DaemonSetRequest{}, func(r *iapiserver.DaemonSetRequest) (any, error) {
		return rc.srv.Kubernetes().DaemonSetRecreate(c, r)
	})
}

func (rc *KubernetesController) DaemonSetDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.DaemonSetRequest{}, func(r *iapiserver.DaemonSetRequest) (any, error) {
		err := rc.srv.Kubernetes().DaemonSetDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) DaemonSetBatchDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.DaemonSetBatchRequest{}, func(r *iapiserver.DaemonSetBatchRequest) (any, error) {
		ret := rc.srv.Kubernetes().DaemonSetBatchDelete(c, r)
		return ret, nil
	})
}

func (rc *KubernetesController) DaemonSetUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.DaemonSetRequest{}, func(r *iapiserver.DaemonSetRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().DaemonSetUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) DaemonSetVersionList(c *gin.Context) {
	runKubernetes(c, &iapiserver.DaemonSetVersionListRequest{}, func(r *iapiserver.DaemonSetVersionListRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().DaemonSetVersionList(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) DaemonSetVersionUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.DaemonSetRequest{}, func(r *iapiserver.DaemonSetRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().DaemonSetVersionUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) StatefulSetGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.StatefulSetGetRequest{}, func(r *iapiserver.StatefulSetGetRequest) (any, error) {
		return rc.srv.Kubernetes().StatefulSetGet(c, r)
	})
}

func (rc *KubernetesController) StatefulSetList(c *gin.Context) {
	runKubernetes(c, &iapiserver.StatefulSetListRequest{}, func(r *iapiserver.StatefulSetListRequest) (any, error) {
		return rc.srv.Kubernetes().StatefulSetList(c, r)
	})
}

func (rc *KubernetesController) StatefulSetAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.StatefulSetRequest{}, func(r *iapiserver.StatefulSetRequest) (any, error) {
		return rc.srv.Kubernetes().StatefulSetAdd(c, r)
	})
}

func (rc *KubernetesController) StatefulSetRecreate(c *gin.Context) {
	runKubernetes(c, &iapiserver.StatefulSetRequest{}, func(r *iapiserver.StatefulSetRequest) (any, error) {
		return rc.srv.Kubernetes().StatefulSetRecreate(c, r)
	})
}

func (rc *KubernetesController) StatefulSetDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.StatefulSetRequest{}, func(r *iapiserver.StatefulSetRequest) (any, error) {
		err := rc.srv.Kubernetes().StatefulSetDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) StatefulSetBatchDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.StatefulSetBatchRequest{}, func(r *iapiserver.StatefulSetBatchRequest) (any, error) {
		ret := rc.srv.Kubernetes().StatefulSetBatchDelete(c, r)
		return ret, nil
	})
}

func (rc *KubernetesController) StatefulSetUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.StatefulSetRequest{}, func(r *iapiserver.StatefulSetRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().StatefulSetUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) StatefulSetVersionList(c *gin.Context) {
	runKubernetes(c, &iapiserver.StatefulSetVersionListRequest{}, func(r *iapiserver.StatefulSetVersionListRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().StatefulSetVersionList(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) StatefulSetVersionUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.StatefulSetRequest{}, func(r *iapiserver.StatefulSetRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().StatefulSetVersionUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) ReplicaSetGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.ReplicaSetGetRequest{}, func(r *iapiserver.ReplicaSetGetRequest) (any, error) {
		return rc.srv.Kubernetes().ReplicaSetGet(c, r)
	})
}

func (rc *KubernetesController) ReplicaSetList(c *gin.Context) {
	runKubernetes(c, &iapiserver.ReplicaSetListRequest{}, func(r *iapiserver.ReplicaSetListRequest) (any, error) {
		return rc.srv.Kubernetes().ReplicaSetList(c, r)
	})
}

func (rc *KubernetesController) ReplicaSetAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.ReplicaSetRequest{}, func(r *iapiserver.ReplicaSetRequest) (any, error) {
		return rc.srv.Kubernetes().ReplicaSetAdd(c, r)
	})
}

func (rc *KubernetesController) ReplicaSetDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.ReplicaSetRequest{}, func(r *iapiserver.ReplicaSetRequest) (any, error) {
		err := rc.srv.Kubernetes().ReplicaSetDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) ReplicaSetUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.ReplicaSetRequest{}, func(r *iapiserver.ReplicaSetRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().ReplicaSetUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) HpaGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.HpaGetRequest{}, func(r *iapiserver.HpaGetRequest) (any, error) {
		return rc.srv.Kubernetes().HpaGet(c, r)
	})
}

func (rc *KubernetesController) HpaList(c *gin.Context) {
	runKubernetes(c, &iapiserver.HpaListRequest{}, func(r *iapiserver.HpaListRequest) (any, error) {
		return rc.srv.Kubernetes().HpaList(c, r)
	})
}

func (rc *KubernetesController) HpaAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.HpaRequest{}, func(r *iapiserver.HpaRequest) (any, error) {
		return rc.srv.Kubernetes().HpaAdd(c, r)
	})
}

func (rc *KubernetesController) HpaDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.HpaRequest{}, func(r *iapiserver.HpaRequest) (any, error) {
		err := rc.srv.Kubernetes().HpaDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) HpaUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.HpaRequest{}, func(r *iapiserver.HpaRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().HpaUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) NamespaceGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.NamespaceGetRequest{}, func(r *iapiserver.NamespaceGetRequest) (any, error) {
		return rc.srv.Kubernetes().NamespaceGet(c, r)
	})
}

func (rc *KubernetesController) NamespaceGatewayGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.PodListRequest{}, func(r *iapiserver.PodListRequest) (any, error) {
		return rc.srv.Kubernetes().NamespaceGatewayGet(c, r)
	})
}

func (rc *KubernetesController) NamespaceList(c *gin.Context) {
	runKubernetes(c, &iapiserver.NamespaceListRequest{}, func(r *iapiserver.NamespaceListRequest) (any, error) {
		return rc.srv.Kubernetes().NamespaceList(c, r)
	})
}

func (rc *KubernetesController) NamespaceAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.NamespaceRequest{}, func(r *iapiserver.NamespaceRequest) (any, error) {
		return rc.srv.Kubernetes().NamespaceAdd(c, r)
	})
}

func (rc *KubernetesController) NamespaceDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.NamespaceRequest{}, func(r *iapiserver.NamespaceRequest) (any, error) {
		err := rc.srv.Kubernetes().NamespaceDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) NamespaceUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.NamespaceRequest{}, func(r *iapiserver.NamespaceRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().NamespaceUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) StorageClassGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.StorageClassGetRequest{}, func(r *iapiserver.StorageClassGetRequest) (any, error) {
		return rc.srv.Kubernetes().StorageClassGet(c, r)
	})
}

func (rc *KubernetesController) StorageClassList(c *gin.Context) {
	runKubernetes(c, &iapiserver.StorageClassListRequest{}, func(r *iapiserver.StorageClassListRequest) (any, error) {
		return rc.srv.Kubernetes().StorageClassList(c, r)
	})
}

func (rc *KubernetesController) StorageClassAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.StorageClassRequest{}, func(r *iapiserver.StorageClassRequest) (any, error) {
		return rc.srv.Kubernetes().StorageClassAdd(c, r)
	})
}

func (rc *KubernetesController) StorageClassDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.StorageClassRequest{}, func(r *iapiserver.StorageClassRequest) (any, error) {
		err := rc.srv.Kubernetes().StorageClassDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) StorageClassBatchDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.StorageClassBatchRequest{}, func(r *iapiserver.StorageClassBatchRequest) (any, error) {
		ret := rc.srv.Kubernetes().StorageClassBatchDelete(c, r)
		return ret, nil
	})
}

func (rc *KubernetesController) StorageClassUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.StorageClassRequest{}, func(r *iapiserver.StorageClassRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().StorageClassUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) PersistentVolumeGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.PersistentVolumeGetRequest{}, func(r *iapiserver.PersistentVolumeGetRequest) (any, error) {
		return rc.srv.Kubernetes().PersistentVolumeGet(c, r)
	})
}

func (rc *KubernetesController) PersistentVolumeList(c *gin.Context) {
	runKubernetes(c, &iapiserver.PersistentVolumeListRequest{}, func(r *iapiserver.PersistentVolumeListRequest) (any, error) {
		return rc.srv.Kubernetes().PersistentVolumeList(c, r)
	})
}

func (rc *KubernetesController) PersistentVolumeAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.PersistentVolumeRequest{}, func(r *iapiserver.PersistentVolumeRequest) (any, error) {
		return rc.srv.Kubernetes().PersistentVolumeAdd(c, r)
	})
}

func (rc *KubernetesController) PersistentVolumeDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.PersistentVolumeRequest{}, func(r *iapiserver.PersistentVolumeRequest) (any, error) {
		err := rc.srv.Kubernetes().PersistentVolumeDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) PersistentVolumeBatchDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.PersistentVolumeBatchRequest{}, func(r *iapiserver.PersistentVolumeBatchRequest) (any, error) {
		ret := rc.srv.Kubernetes().PersistentVolumeBatchDelete(c, r)
		return ret, nil
	})
}

func (rc *KubernetesController) PersistentVolumeUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.PersistentVolumeRequest{}, func(r *iapiserver.PersistentVolumeRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().PersistentVolumeUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) PersistentVolumeClaimGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.PersistentVolumeClaimGetRequest{}, func(r *iapiserver.PersistentVolumeClaimGetRequest) (any, error) {
		return rc.srv.Kubernetes().PersistentVolumeClaimGet(c, r)
	})
}

func (rc *KubernetesController) PersistentVolumeClaimList(c *gin.Context) {
	runKubernetes(c, &iapiserver.PersistentVolumeClaimListRequest{}, func(r *iapiserver.PersistentVolumeClaimListRequest) (any, error) {
		return rc.srv.Kubernetes().PersistentVolumeClaimList(c, r)
	})
}

func (rc *KubernetesController) PersistentVolumeClaimAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.PersistentVolumeClaimRequest{}, func(r *iapiserver.PersistentVolumeClaimRequest) (any, error) {
		return rc.srv.Kubernetes().PersistentVolumeClaimAdd(c, r)
	})
}

func (rc *KubernetesController) PersistentVolumeClaimDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.PersistentVolumeClaimRequest{}, func(r *iapiserver.PersistentVolumeClaimRequest) (any, error) {
		err := rc.srv.Kubernetes().PersistentVolumeClaimDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) PersistentVolumeClaimBatchDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.PersistentVolumeClaimBatchRequest{}, func(r *iapiserver.PersistentVolumeClaimBatchRequest) (any, error) {
		ret := rc.srv.Kubernetes().PersistentVolumeClaimBatchDelete(c, r)
		return ret, nil
	})
}

func (rc *KubernetesController) PersistentVolumeClaimUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.PersistentVolumeClaimRequest{}, func(r *iapiserver.PersistentVolumeClaimRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().PersistentVolumeClaimUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) VolumeSnapshotClassGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.VolumeSnapshotClassGetRequest{}, func(r *iapiserver.VolumeSnapshotClassGetRequest) (any, error) {
		return rc.srv.Kubernetes().VolumeSnapshotClassGet(c, r)
	})
}

func (rc *KubernetesController) VolumeSnapshotClassList(c *gin.Context) {
	runKubernetes(c, &iapiserver.VolumeSnapshotClassListRequest{}, func(r *iapiserver.VolumeSnapshotClassListRequest) (any, error) {
		return rc.srv.Kubernetes().VolumeSnapshotClassList(c, r)
	})
}

func (rc *KubernetesController) VolumeSnapshotClassAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.VolumeSnapshotClassRequest{}, func(r *iapiserver.VolumeSnapshotClassRequest) (any, error) {
		return rc.srv.Kubernetes().VolumeSnapshotClassAdd(c, r)
	})
}

func (rc *KubernetesController) VolumeSnapshotClassDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.VolumeSnapshotClassRequest{}, func(r *iapiserver.VolumeSnapshotClassRequest) (any, error) {
		err := rc.srv.Kubernetes().VolumeSnapshotClassDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) VolumeSnapshotClassBatchDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.VolumeSnapshotClassBatchRequest{}, func(r *iapiserver.VolumeSnapshotClassBatchRequest) (any, error) {
		ret := rc.srv.Kubernetes().VolumeSnapshotClassBatchDelete(c, r)
		return ret, nil
	})
}

func (rc *KubernetesController) VolumeSnapshotClassUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.VolumeSnapshotClassRequest{}, func(r *iapiserver.VolumeSnapshotClassRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().VolumeSnapshotClassUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) VolumeSnapshotGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.VolumeSnapshotGetRequest{}, func(r *iapiserver.VolumeSnapshotGetRequest) (any, error) {
		return rc.srv.Kubernetes().VolumeSnapshotGet(c, r)
	})
}

func (rc *KubernetesController) VolumeSnapshotList(c *gin.Context) {
	runKubernetes(c, &iapiserver.VolumeSnapshotListRequest{}, func(r *iapiserver.VolumeSnapshotListRequest) (any, error) {
		return rc.srv.Kubernetes().VolumeSnapshotList(c, r)
	})
}

func (rc *KubernetesController) VolumeSnapshotAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.VolumeSnapshotRequest{}, func(r *iapiserver.VolumeSnapshotRequest) (any, error) {
		return rc.srv.Kubernetes().VolumeSnapshotAdd(c, r)
	})
}

func (rc *KubernetesController) VolumeSnapshotDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.VolumeSnapshotRequest{}, func(r *iapiserver.VolumeSnapshotRequest) (any, error) {
		err := rc.srv.Kubernetes().VolumeSnapshotDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) VolumeSnapshotBatchDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.VolumeSnapshotBatchRequest{}, func(r *iapiserver.VolumeSnapshotBatchRequest) (any, error) {
		ret := rc.srv.Kubernetes().VolumeSnapshotBatchDelete(c, r)
		return ret, nil
	})
}

func (rc *KubernetesController) VolumeSnapshotUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.VolumeSnapshotRequest{}, func(r *iapiserver.VolumeSnapshotRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().VolumeSnapshotUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) VolumeSnapshotContentGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.VolumeSnapshotContentGetRequest{}, func(r *iapiserver.VolumeSnapshotContentGetRequest) (any, error) {
		return rc.srv.Kubernetes().VolumeSnapshotContentGet(c, r)
	})
}

func (rc *KubernetesController) VolumeSnapshotContentList(c *gin.Context) {
	runKubernetes(c, &iapiserver.VolumeSnapshotContentListRequest{}, func(r *iapiserver.VolumeSnapshotContentListRequest) (any, error) {
		return rc.srv.Kubernetes().VolumeSnapshotContentList(c, r)
	})
}

func (rc *KubernetesController) VolumeSnapshotContentAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.VolumeSnapshotContentRequest{}, func(r *iapiserver.VolumeSnapshotContentRequest) (any, error) {
		return rc.srv.Kubernetes().VolumeSnapshotContentAdd(c, r)
	})
}

func (rc *KubernetesController) VolumeSnapshotContentDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.VolumeSnapshotContentRequest{}, func(r *iapiserver.VolumeSnapshotContentRequest) (any, error) {
		err := rc.srv.Kubernetes().VolumeSnapshotContentDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) VolumeSnapshotContentBatchDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.VolumeSnapshotContentBatchRequest{}, func(r *iapiserver.VolumeSnapshotContentBatchRequest) (any, error) {
		ret := rc.srv.Kubernetes().VolumeSnapshotContentBatchDelete(c, r)
		return ret, nil
	})
}

func (rc *KubernetesController) VolumeSnapshotContentUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.VolumeSnapshotContentRequest{}, func(r *iapiserver.VolumeSnapshotContentRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().VolumeSnapshotContentUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) IngressGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.IngressGetRequest{}, func(r *iapiserver.IngressGetRequest) (any, error) {
		return rc.srv.Kubernetes().IngressGet(c, r)
	})
}

func (rc *KubernetesController) IngressList(c *gin.Context) {
	runKubernetes(c, &iapiserver.IngressListRequest{}, func(r *iapiserver.IngressListRequest) (any, error) {
		return rc.srv.Kubernetes().IngressList(c, r)
	})
}

func (rc *KubernetesController) IngressAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.IngressRequest{}, func(r *iapiserver.IngressRequest) (any, error) {
		return rc.srv.Kubernetes().IngressAdd(c, r)
	})
}

func (rc *KubernetesController) IngressDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.IngressRequest{}, func(r *iapiserver.IngressRequest) (any, error) {
		err := rc.srv.Kubernetes().IngressDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) IngressBatchDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.IngressBatchRequest{}, func(r *iapiserver.IngressBatchRequest) (any, error) {
		ret := rc.srv.Kubernetes().IngressBatchDelete(c, r)
		return ret, nil
	})
}

func (rc *KubernetesController) IngressUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.IngressRequest{}, func(r *iapiserver.IngressRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().IngressUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}

func (rc *KubernetesController) NetworkPolicyGet(c *gin.Context) {
	runKubernetes(c, &iapiserver.NetworkPolicyGetRequest{}, func(r *iapiserver.NetworkPolicyGetRequest) (any, error) {
		return rc.srv.Kubernetes().NetworkPolicyGet(c, r)
	})
}

func (rc *KubernetesController) NetworkPolicyList(c *gin.Context) {
	runKubernetes(c, &iapiserver.NetworkPolicyListRequest{}, func(r *iapiserver.NetworkPolicyListRequest) (any, error) {
		return rc.srv.Kubernetes().NetworkPolicyList(c, r)
	})
}

func (rc *KubernetesController) NetworkPolicyAdd(c *gin.Context) {
	runKubernetes(c, &iapiserver.NetworkPolicyRequest{}, func(r *iapiserver.NetworkPolicyRequest) (any, error) {
		return rc.srv.Kubernetes().NetworkPolicyAdd(c, r)
	})
}

func (rc *KubernetesController) NetworkPolicyDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.NetworkPolicyRequest{}, func(r *iapiserver.NetworkPolicyRequest) (any, error) {
		err := rc.srv.Kubernetes().NetworkPolicyDelete(c, r)
		return nil, errors.WithStack(err)
	})
}

func (rc *KubernetesController) NetworkPolicyBatchDelete(c *gin.Context) {
	runKubernetes(c, &iapiserver.NetworkPolicyBatchRequest{}, func(r *iapiserver.NetworkPolicyBatchRequest) (any, error) {
		ret := rc.srv.Kubernetes().NetworkPolicyBatchDelete(c, r)
		return ret, nil
	})
}

func (rc *KubernetesController) NetworkPolicyUpdate(c *gin.Context) {
	runKubernetes(c, &iapiserver.NetworkPolicyRequest{}, func(r *iapiserver.NetworkPolicyRequest) (any, error) {
		ret, err := rc.srv.Kubernetes().NetworkPolicyUpdate(c, r)
		return ret, errors.WithStack(err)
	})
}
