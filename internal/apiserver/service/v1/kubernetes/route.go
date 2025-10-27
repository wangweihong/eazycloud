package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/log"
	v1 "k8s.io/api/core/v1"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/internal/pkg/clientset"
	"github.com/wangweihong/eazycloud/internal/pkg/libkubernetes"

	//metav1 "k8s.io/api/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	appsv1 "k8s.io/api/apps/v1"

	//monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

func (k *kubernetesService) RouterGet(ctx context.Context, req *iapiserver.RouterGetRequest) (*iapiserver.ServiceInfo, error) {
	if req.Namespace == "" {
		return nil, errors.New("namespace is empty")
	}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, err
	}

	service := &v1.Service{}
	service.Name = iapiserver.IngressControllerPrefix + req.Namespace
	meta, err := clientset.ServiceGet(ctx, cluster, iapiserver.IngressControllerNamespace, service.Name, req.ToGetOpts())
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			err = errors.Errorf("network gataway not enable")
		}

		return nil, err
	}

	return convertK8sServiceToApiService(meta, cluster, req.Yaml), nil

}

// func (k *kubernetesService) RouterList(ctx context.Context, req *topke.RouterListRequest) *topke.RouterListResponse {
// 	clusterList, err := tm.GetVisitScope(req.ResourceListRequest)
// 	if err != nil {
// 		resp.Status = status.UpdateStatus(err)
// 		return resp
// 	}

// 	wg := utils.NewWaitGroup(nil)
// 	for _, cluster := range clusterList {
// 		cluster := cluster
// 		wg.Start(utils.NewWaitGroupHandleFunc("", nil, func() utils.WaitGroupResult {
// 			clusterListOne := topke.NewEachResourceRangeListState(cluster.UUID, cluster.Name)

// 			rets, err := clientset.ServiceList(cluster, topke.IngressControllerNamespace, req.ToListOpts())
// 			if err != nil {
// 				return utils.NewWaitGroupResult(clusterListOne, status.UpdateStatus(err))
// 			}

// 			list := make([]*topke.ServiceInfo, len(rets.Items), len(rets.Items))
// 			for k, v := range rets.Items {
// 				v := v
// 				one := convertK8sServiceToApiService(&v, cluster, false)
// 				list[k] = one
// 			}
// 			clusterListOne.TotalCount = len(list)
// 			clusterListOne.List = list
// 			return utils.NewWaitGroupResult(clusterListOne, nil)
// 		}))
// 	}
// 	wg.Wait()

// 	CutPagingSliceFromWgResultsV2(wg, &resp.EachRangeListState, &resp.List, req.PageNumber, req.PageSize, &resp.TotalCount, req.SortBy, req.SortDesc)

// 	return resp
// }

const (
	defaultNginxIngressImage = "k8s.gcr.io/ingress-nginx/controller:v0.20.0"
)

// func (k *kubernetesService) RouterCreate(ctx context.Context, req *iapiserver.RouterRequest) error {
// 	if req.Namespace == "" {
// 		return errors.Errorf( "namespace is empty")
// 	}

// 	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
// 	if err != nil {
// 		return err
// 	}

// 	ret, err := clientset.PodList(cluster, "kube-system", metav1.ListOptions{})
// 	if err != nil {
// 		return err
// 	}

// 	req.Image = defaultNginxIngressImage
// 	for _, v := range ret.Items {
// 		if strings.Contains(v.Name, "kube-apiserver") {
// 			for _, c := range v.Spec.Containers {
// 				if strings.Contains(c.Image, "kube-apiserver") {
// 					splitImage := strings.SplitN(c.Image, "/", 3) //
// 					if len(splitImage) == 3 {
// 						req.Image = strings.Join([]string{splitImage[0], splitImage[1], req.Image}, "/")
// 					}
// 					break
// 				}
// 			}
// 			break
// 		}
// 	}

// 	// 如果网关已经开启，那么不允许再次开启
// 	if _, err = clientset.NamespaceGet(ctx,cluster, req.Namespace, req.GetOpts); err != nil {
// 		return err
// 	}

// 	if err := k.create(ctx,cluster, req); err != nil {
// 		return err
// 	}

// 	return nil
// }

func (k *kubernetesService) RouterDelete(ctx context.Context, req *iapiserver.RouterRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return err
	}
	// 删除service daemonset configmap
	configmap := v1.ConfigMap{}
	deploymentName := iapiserver.IngressControllerPrefix + req.Namespace
	serviceName := iapiserver.IngressControllerPrefix + req.Namespace
	configmap.Name = iapiserver.IngressControllerPrefix + req.Namespace

	if err := clientset.DeploymentDelete(ctx, cluster, iapiserver.IngressControllerNamespace, deploymentName, req.DeleteOpts); err != nil {
		log.Errorf("%v", err)
	}

	if err := clientset.ServiceDelete(ctx, cluster, iapiserver.IngressControllerNamespace, serviceName, req.DeleteOpts); err != nil {
		log.Errorf("%v", err)
	}

	if err := clientset.ConfigMapDelete(ctx, cluster, iapiserver.IngressControllerNamespace, &configmap, req.DeleteOpts); err != nil {
		log.Errorf("%v", err)
	}

	return nil
}

func (k *kubernetesService) create(ctx context.Context, cluster *iapiserver.Cluster, req *iapiserver.RouterRequest) error {
	if _, err := clientset.NamespaceGet(ctx, cluster, iapiserver.IngressControllerNamespace, req.GetOpts); err != nil {
		ns := v1.Namespace{}
		ns.Name = iapiserver.IngressControllerNamespace
		if _, err = clientset.NamespaceCreate(ctx, cluster, &ns, req.CreateOpts); err != nil {
			return err
		}
	}

	deplomentName := iapiserver.IngressControllerPrefix + req.Namespace
	serviceName := iapiserver.IngressControllerPrefix + req.Namespace
	configmapName := iapiserver.IngressControllerPrefix + req.Namespace
	labels := map[string]string{
		"namespace": req.Namespace,
	}

	// create serviceAccount   		 : kind = ServiceAccount      , name = nginx-ingress-serviceaccount
	// create namespace role  		 : kind = Role				  , name = nginx-ingress-role
	// create namesapce rolebanding	 : kind = RoleBinding		  , name = nginx-ingress-role-nisa-binding
	// create cluster role   		 : kind = ClusterRole   	  , name = nginx-ingress-clusterrole
	// create cluster rolebanding	 : kind = ClusterRoleBinding  , name = nginx-ingress-clusterrole-nisa-binding
	configs := strings.Split(config, "---")
	for _, v := range configs {
		jsondata, err := yamlToJson(v)
		if err != nil {
			return err
		}

		meta := runtime.TypeMeta{}
		if err := json.Unmarshal([]byte(jsondata), &meta); err != nil {
			return err
		}
		switch meta.Kind {
		case "ConfigMap":
			c := &v1.ConfigMap{}
			if err := json.Unmarshal([]byte(jsondata), c); err != nil {
				return err
			}
			c.Name = configmapName
			if _, err = clientset.ConfigMapGet(ctx, cluster, iapiserver.IngressControllerNamespace, configmapName, req.GetOpts); err != nil {
				if _, err = clientset.ConfigMapCreate(ctx, cluster, iapiserver.IngressControllerNamespace, c, req.CreateOpts); err != nil {
					return err
				}
			}
		case "Deployment":
			c := &appsv1.Deployment{}
			if err := json.Unmarshal([]byte(jsondata), c); err != nil {
				return err
			}
			_, err = clientset.DeploymentGet(ctx, cluster, iapiserver.IngressControllerNamespace, c.Name, req.GetOpts)
			if err != nil {
				//            - --configmap=$(POD_NAMESPACE)/nginx-configuration
				//            - --watch-namespace=ltsns
				c.Name = deplomentName
				c.Spec.Selector.MatchLabels = labels
				c.Spec.Template.Labels = labels
				if c.Spec.Template.Spec.NodeSelector == nil {
					c.Spec.Template.Spec.NodeSelector = map[string]string{}
				}
				c.Spec.Template.Spec.NodeSelector["type"] = "gateway"
				c.Spec.Template.Spec.Tolerations = []v1.Toleration{
					{
						Key:    "type",
						Value:  "gateway",
						Effect: "NoExecute",
					},
				}
				c.Spec.Template.Spec.Containers[0].Image = req.Image
				c.Spec.Template.Spec.Containers[0].Args = append(c.Spec.Template.Spec.Containers[0].Args, fmt.Sprintf("--configmap=$(POD_NAMESPACE)/%v", configmapName))
				c.Spec.Template.Spec.Containers[0].Args = append(c.Spec.Template.Spec.Containers[0].Args, fmt.Sprintf("--watch-namespace=%v", req.Namespace))
				if _, err = clientset.DeploymentCreate(ctx, cluster, iapiserver.IngressControllerNamespace, c, req.CreateOpts); err != nil {
					return err
				}
			}
		case "Service":
			c := &v1.Service{}
			if err := json.Unmarshal([]byte(jsondata), c); err != nil {
				return err
			}
			if _, err = clientset.ServiceGet(ctx, cluster, iapiserver.IngressControllerNamespace, c.Name, req.GetOpts); err != nil {
				c.Name = serviceName
				c.Spec.Selector = labels
				c.Labels = labels
				if _, err = clientset.ServiceCreate(ctx, cluster, iapiserver.IngressControllerNamespace, c, req.CreateOpts); err != nil {
					return err
				}
			}
		case "ServiceAccount":
			c := &v1.ServiceAccount{}
			if err := json.Unmarshal([]byte(jsondata), c); err != nil {
				return err
			}
			if _, err = clientset.ServiceAccountGet(ctx, cluster, iapiserver.IngressControllerNamespace, c, req.GetOpts); err != nil {
				if _, err = clientset.ServiceAccountCreate(ctx, cluster, iapiserver.IngressControllerNamespace, c, req.CreateOpts); err != nil {
					return err
				}
			}
		case "Role":
			c := &rbacv1.Role{}
			if err := json.Unmarshal([]byte(jsondata), c); err != nil {
				return err
			}
			if _, err = clientset.RoleGet(ctx, cluster, iapiserver.IngressControllerNamespace, c, req.GetOpts); err != nil {
				if _, err = clientset.RoleCreate(ctx, cluster, iapiserver.IngressControllerNamespace, c, req.CreateOpts); err != nil {
					return err
				}
			}
		case "RoleBinding":
			c := &rbacv1.RoleBinding{}
			if err := json.Unmarshal([]byte(jsondata), c); err != nil {
				return err
			}
			if _, err = libkubernetes.RoleBindingGet(ctx, cluster.Config, iapiserver.IngressControllerNamespace, c.Name, req.GetOpts); err != nil {
				if _, err = libkubernetes.RoleBindingCreate(ctx, cluster.Config, iapiserver.IngressControllerNamespace, c, req.CreateOpts); err != nil {
					return err
				}
			}
		case "ClusterRole":
			c := &rbacv1.ClusterRole{}
			if err := json.Unmarshal([]byte(jsondata), c); err != nil {
				return err
			}
			if _, err = libkubernetes.ClusterRoleGet(ctx, cluster.Config, c.Name, req.GetOpts); err != nil {
				if _, err = libkubernetes.ClusterRoleCreate(ctx, cluster.Config, c, req.CreateOpts); err != nil {
					return err
				}
			}
		case "ClusterRoleBinding":
			c := &rbacv1.ClusterRoleBinding{}
			if err := json.Unmarshal([]byte(jsondata), c); err != nil {
				return err
			}
			if _, err = libkubernetes.ClusterRoleBindingGet(ctx, cluster.Config, c.Name, req.GetOpts); err != nil {
				if _, err = libkubernetes.ClusterRoleBindingCreate(ctx, cluster.Config, c, req.CreateOpts); err != nil {
					return err
				}
			}
		default:
			if meta.Kind != "" {
				log.Errorf("%v", fmt.Sprintf("kind[%v]", meta.Kind))
			}
		}
	}

	return nil
}

var config = `
kind: ConfigMap
apiVersion: v1
metadata:
  namespace: system-router
---
apiVersion: v1
kind: Service
metadata:
  namespace: system-router
spec:
  type: NodePort
  ports:
  - name: http
    port: 80
    targetPort: 80
    protocol: TCP
  selector:
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-ingress-controller
  namespace: system-router
spec:
  replicas: 1
  selector:
    matchLabels:
  template:
    metadata:
      labels:
      annotations:
        prometheus.io/port: "10254"
        prometheus.io/scrape: "true"
    spec:
      serviceAccountName: nginx-ingress-serviceaccount
      containers:
        - name: nginx-ingress-controller
          image: k8s.gcr.io/nginx-ingress-controller:0.20.0
          args:
            - /nginx-ingress-controller
          securityContext:
            capabilities:
              drop:
                - ALL
              add:
                - NET_BIND_SERVICE
            # www-data -> 33
            runAsUser: 33
          env:
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
          ports:
            - name: http
              containerPort: 80
            - name: https
              containerPort: 443
          livenessProbe:
            failureThreshold: 3
            httpGet:
              path: /healthz
              port: 10254
              scheme: HTTP
            initialDelaySeconds: 10
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 1
          readinessProbe:
            failureThreshold: 3
            httpGet:
              path: /healthz
              port: 10254
              scheme: HTTP
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 1
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: nginx-ingress-serviceaccount
  namespace: system-router
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: nginx-ingress-clusterrole
rules:
  - apiGroups:
    - '*'
    resources:
    - '*'
    verbs:
    - '*'
  - nonResourceURLs:
    - '*'
    verbs:
    - '*'
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: nginx-ingress-clusterrole-nisa-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: nginx-ingress-clusterrole
subjects:
  - kind: ServiceAccount
    name: nginx-ingress-serviceaccount
    namespace: system-router
---
`
