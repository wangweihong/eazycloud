package kubernetes

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/internal/pkg/clientset"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/paging"
	"github.com/wangweihong/gotoolbox/pkg/sets"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *kubernetesService) NodePortList(ctx context.Context, req *iapiserver.NodePortListRequest) (*iapiserver.NodePortListResponse, error) {
	resp := &iapiserver.NodePortListResponse{}

	clusters, err := getVisitScope(ctx, k.store, req.ResourceListRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, v := range clusters {
		nodePorts := clientset.GetClusterNodePortServices(v.ID)
		for _, np := range nodePorts {
			if req.Namespace != "" && np.Service.Namespace != req.Namespace {
				continue
			}

			if req.Fuzzy != "" && !strings.Contains(np.Service.Name, req.Fuzzy) {
				continue
			}

			resp.List = append(resp.List, &iapiserver.NodePort{
				Port:     np.Port,
				Service:  iapiserver.NewServiceInfo(np.Service, v),
				Endpoint: strings.Split(strings.TrimPrefix(v.Config.Host, "https://"), ":")[0] + ":" + strconv.Itoa(int(np.Port)),
			})
		}
	}

	sort.Slice(resp.List, func(i, j int) bool {
		if resp.List[i].Port == resp.List[j].Port {
			return resp.List[i].Service.Resource.Name < resp.List[j].Service.Resource.Name
		}
		return resp.List[i].Port < resp.List[j].Port
	})

	resp.TotalCount = len(resp.List)
	s, e := paging.Index(len(resp.List), req.PageNum, req.PageSize)
	resp.List = resp.List[s:e]
	return resp, nil
}

func (k *kubernetesService) IsNodePortUsed(ctx context.Context, req *iapiserver.NodePortRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	svrList, err := clientset.ServiceList(ctx, cluster, "", metav1.ListOptions{})
	if err != nil {
		return errors.WithStack(err)
	}

	ports := sets.NewInt(req.Ports...)

	for _, r := range svrList.Items {
		if r.Spec.Type == v1.ServiceTypeNodePort {
			for _, v := range r.Spec.Ports {
				if v.NodePort != 0 {
					if ports.Has(int(v.NodePort)) {
						return errors.Errorf("port %v has been used", v.NodePort)
					}
				}
			}
		}
	}
	//nodePorts := clientset.GetClusterNodePortServices(req.ClusterUUID)
	//for _, v := range req.Ports {
	//	if _, ok := nodePorts[v]; ok {
	//		resp.Status = status.NewStatusDesc(scode.TopECTopkePortHasBeenUsed, strconv.Itoa(int(v)))
	//		return resp
	//	}
	//}
	return nil
}

