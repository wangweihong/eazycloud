package kubernetes

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	snapshotv1beta1 "github.com/kubernetes-csi/external-snapshotter/client/v3/apis/volumesnapshot/v1beta1"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/log"
	"github.com/wangweihong/gotoolbox/pkg/waitgroup"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/internal/pkg/clientset"
	"github.com/wangweihong/eazycloud/internal/pkg/libkubernetes"
)

const (
	NFSDriverPrefix = "nfs-pod-provisioner"
)

func (k *kubernetesService) StorageClassListAll(ctx context.Context, req *iapiserver.StorageClassListRequest) (*iapiserver.StorageClassListResponse, error) {
	resp := &iapiserver.StorageClassListResponse{}

	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList(ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.StorageClassInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.StorageClassInfo](cluster.ID, cluster.Name)
			resList, err := clientset.StorageClassList(ctx, cluster, req.Namespace, req.ToListOpts())
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.StorageClassInfo
			for i := range resList.Items {

				resInfo := iapiserver.NewStorageClassInfo(&resList.Items[i], cluster)
				if filterStorageClass(resInfo, req.Fuzzy) {
					continue
				}
				resInfos = append(resInfos, resInfo)
			}
			clusterListOne.TotalCount = len(resInfos)
			clusterListOne.List = resInfos
			return waitgroup.NewGenericResult(clusterListOne, nil)
		}, func(i, j int) bool {
			return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
		}, 10*time.Second)
	return resp, err
}

func (k *kubernetesService) StorageClassDelete(ctx context.Context, req *iapiserver.StorageClassRequest) (*iapiserver.StorageClassInfo, error) {

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resList, err := clientset.PersistentVolumeClaimList(ctx, cluster, "", metav1.ListOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range resList.Items {
		if resList.Items[i].Spec.StorageClassName != nil && *resList.Items[i].Spec.StorageClassName == req.Resource.Name {
			return nil, errors.Errorf("cannot delete storageclass if has persistent volumes")
		}
	}

	if err := clientset.StorageClassDelete(ctx, cluster, req.Resource, req.DeleteOpts); err != nil {
		return nil, errors.WithStack(err)
	}

	return nil, errors.WithStack(err)
}

func (k *kubernetesService) StorageClassBatchDelete(ctx context.Context, req *iapiserver.StorageClassBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.StorageClassRequest] {
	wg := waitgroup.RunGenericConcurrently(
		ctx, req.Resources, func(ctx context.Context, res *iapiserver.StorageClassRequest) waitgroup.GenericResult[*iapiserver.StorageClassRequest] {
			cluster, err := k.store.Clusters().Get(ctx, res.Cluster)
			if err != nil {
				return waitgroup.NewGenericResult(res, err)
			}
			if err := clientset.StorageClassDelete(ctx, cluster, res.Resource, res.DeleteOpts); err != nil {
				return waitgroup.NewGenericResult(res, err)
			}
			return waitgroup.NewGenericResult(res, nil)
		})
	return wg.BatchGenericOutput()
}

var (
	nfsNecessaryParameterKey = []string{
		"nfs-server",
		"nfs-path",
	}

	storageClassNecessaryParameterKey = []string{
		"resturl",
		"restuser",
		"restuserkey",
	}
)

const (
	nfsControllerYaml = `
kind: Deployment
apiVersion: apps/v1
metadata:
 name: nfs-pod-provisioner  #  replace name with-storage
 namespace: nfs
spec:
 replicas: 1
 strategy:
   type: Recreate
 selector:
   matchLabels:
     app: nfs-pod-provisioner #  replace name with-storage
 template:
   metadata:
     labels:
       app: nfs-pod-provisioner #  replace name with-storage
   spec:
     serviceAccountName: nfs-pod-provisioner-sa # name of service account created in rbac.yaml
     containers:
       - name: nfs-pod-provisioner
         image: gcr.io/k8s-staging-sig-storage/nfs-subdir-external-provisioner:v4.0.0   #  replace name with-storage
         #image: nfs-subdir-external-provisioner:v4.0.0   #  replace name with-storage
         volumeMounts:
           - name: nfs-pod-provisioner
             mountPath: /persistentvolumes
         env:
           - name: PROVISIONER_NAME # do not change
             value: nfs-test # SAME AS PROVISONER NAME VALUE IN STORAGECLASS
           - name: NFS_SERVER # do not change
             value: 10.30.100.157 # Ip of the NFS SERVER
           - name: NFS_PATH # do not change
             value: /nfsshare # path to nfs directory setup
     volumes:
      - name: nfs-pod-provisioner   # same as volumemouts name
        nfs:
          server: 10.30.100.157  # nfs server ip
          path: /nfsshare
`
)

// handle topsc-provisioner storageclass
func (k *kubernetesService) StorageClassCreate(ctx context.Context, req *iapiserver.StorageClassRequest) (*iapiserver.StorageClassInfo, error) {
	if req.Resource.Provisioner == "" {
		return nil, errors.Errorf("req.StorageClass.Provisioner is empty")
	}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// glusterfs
	if req.Resource.Provisioner == iapiserver.StorageClassDriverGlusterFs {
		meta, err := k.storageClassGlusterfsCreate(ctx, req)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return iapiserver.NewStorageClassInfo(meta, cluster), nil
	}

	var nfsDeploy *appsv1.Deployment
	// nfs storage class
	if strings.HasPrefix(req.Resource.Provisioner, iapiserver.StorageClassDriverNfs) {
		if req.Resource.Parameters == nil {
			return nil, errors.Errorf("req.StorageClass.Parameter is empty")
		}

		for _, key := range nfsNecessaryParameterKey {
			value, ok := req.Resource.Parameters[key]
			if !ok {
				return nil, errors.Errorf("req.StorageClass.Parameter necessary key is missing")
			}
			if value == "" {
				return nil, errors.Errorf("req.StorageClass.Parameter necessary key/value is empty")

			}
		}

		// ret, err := clientset.PodList(ctx, cluster, topke.KubeSystemNamespace, metav1.ListOptions{})
		// if err != nil {
		// 	return nil,errors.WithStack(err)
		// }

		// image := "quay.io/external_storage/nfs-subdir-external-provisioner:v4.0.0"
		// for _, v := range ret.Items {
		// 	if strings.Contains(v.Name, "kube-apiserver") {
		// 		for _, c := range v.Spec.Containers {
		// 			if strings.Contains(c.Image, iapiserver.ApiServerImageRepository) {
		// 				splitImage := strings.SplitN(c.Image, "/", 3) //
		// 				if len(splitImage) == 3 {
		// 					image = strings.Join([]string{splitImage[0], splitImage[1], image}, "/")
		// 				}
		// 				break
		// 			}
		// 		}
		// 		break
		// 	}
		// }

		nfsServerIP := req.Resource.Parameters["nfs-server"]
		nfsPath := req.Resource.Parameters["nfs-path"]

		nfsDeploy = &appsv1.Deployment{}
		nfsDeploy.ResourceVersion = ""
		if err := yaml.Unmarshal([]byte(nfsControllerYaml), nfsDeploy); err != nil {
			return nil, errors.WithStack(err)
		}
		nfsDeploy.Namespace = "kube-system"
		nfsDeploy.Name = "nfs-pod-provisioner" + "-" + req.Resource.Name
		if nfsDeploy.Spec.Selector.MatchLabels == nil {
			nfsDeploy.Spec.Selector.MatchLabels = make(map[string]string)
		}
		nfsDeploy.Spec.Selector.MatchLabels["app"] = nfsDeploy.Name

		if nfsDeploy.Spec.Template.Labels == nil {
			nfsDeploy.Spec.Template.Labels = make(map[string]string)
		}
		nfsDeploy.Spec.Template.Labels["app"] = nfsDeploy.Name

		//handle pod template
		for i, v := range nfsDeploy.Spec.Template.Spec.Volumes {
			if v.NFS != nil && v.Name == "nfs-pod-provisioner" {
				v.NFS.Server = nfsServerIP
				v.NFS.Path = nfsPath
				nfsDeploy.Spec.Template.Spec.Volumes[i] = v

				break
			}
		}

		for i, container := range nfsDeploy.Spec.Template.Spec.Containers {
			if container.Name == "nfs-pod-provisioner" {
				for j, e := range container.Env {
					if e.Name == "PROVISIONER_NAME" {
						e.Value = req.Resource.Provisioner
					}
					if e.Name == "NFS_SERVER" {
						e.Value = nfsServerIP
					}
					if e.Name == "NFS_PATH" {
						e.Value = nfsPath
					}
					container.Env[j] = e
				}

				// //containerImage
				// if image != "" {
				// 	container.Image = image
				// }
			}
			nfsDeploy.Spec.Template.Spec.Containers[i] = container
		}
	}

	meta, err := clientset.StorageClassCreate(ctx, cluster, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if nfsDeploy != nil {
		//TODO: delete deployment when storageclass delete
		go func() {
			if _, err := clientset.DeploymentCreate(ctx, cluster, nfsDeploy.Namespace, nfsDeploy, metav1.CreateOptions{}); err != nil {
				log.Errorf("create topke nfs  driver err: %v", err.Error())
			}
		}()
	}

	return iapiserver.NewStorageClassInfo(meta, cluster), nil
}

func (k *kubernetesService) storageClassGlusterfsCreateCheck(req *iapiserver.StorageClassRequest) error {
	if req.Resource == nil || req.Resource.Parameters == nil {
		return errors.Errorf("resturl is empty")
	}

	resturl := req.Resource.Parameters["resturl"]
	if resturl == "" {
		return errors.Errorf("resturl is empty")
	}

	if req.Resource.Parameters["restauthenabled"] == "true" {
		user := req.Resource.Parameters["restuser"]
		userkey := req.Resource.Parameters["restuserkey"]
		if user == "" || userkey == "" {
			return errors.Errorf("restuser or restuserkey is empty")
		}
		decoded, err := base64.StdEncoding.DecodeString(userkey)
		if err != nil {
			return errors.WithStack(err)
		}
		decodestr := string(decoded)
		req.Resource.Parameters["restuserkey"] = decodestr
	}

	// refer volumetype
	//replicate := req.StorageClass.Parameters["volumetype"]
	//if replicate == "" {
	//	req.StorageClass.Parameters["volumetype"] = "replicate:2"
	//} else {
	//	replicateInt, err := strconv.Atoi(replicate)
	//	if err != nil {
	//		return errors.Errorf( fmt.Sprintf("replicate is err[%v]", replicateInt))
	//	}
	//	req.StorageClass.Parameters["volumetype"] = fmt.Sprintf("replicate:%v", replicateInt)
	//}

	for _, key := range storageClassNecessaryParameterKey {
		value, ok := req.Resource.Parameters[key]
		if !ok {
			return errors.Errorf("req.StorageClass.Parameter necessary key is missing:")

		}
		if value == "" {
			return errors.Errorf("req.StorageClass.Parameter necessary key/value is empty")
		}
	}

	return nil
}

func (k *kubernetesService) storageClassGlusterfsCreate(ctx context.Context, req *iapiserver.StorageClassRequest) (*storagev1.StorageClass, error) {
	if req.Resource.Parameters == nil {
		return nil, errors.Errorf("resturl is empty")
	}

	if err := k.storageClassGlusterfsCreateCheck(req); err != nil {
		return nil, errors.WithStack(err)
	}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.StorageClassCreate(ctx, cluster, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return meta, nil
}

func (k *kubernetesService) StorageClassUpdate(ctx context.Context, req *iapiserver.StorageClassRequest) (*iapiserver.StorageClassInfo, error) {
	if err := libkubernetes.ValidateClusterScopedObjectParameters(req.Resource); err != nil {
		return nil, errors.WithStack(err)
	}
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.StorageClassUpdate(ctx, cluster, req.Resource.Namespace, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewStorageClassInfo(meta, cluster), nil
}

func (k *kubernetesService) StorageClassSetDefault(ctx context.Context, req *iapiserver.StorageClassRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return err
	}

	storageList, err := clientset.StorageClassList(ctx, cluster, req.Resource.Namespace, metav1.ListOptions{})
	if err != nil {
		return err
	}
	//remove old default storageclass default flag
	wg := waitgroup.RunGenericConcurrently[storagev1.StorageClass, storagev1.StorageClass](ctx, storageList.Items,
		func(ctx context.Context, res storagev1.StorageClass) waitgroup.GenericResult[storagev1.StorageClass] {
			var needUpdate bool
			if res.Name == req.Resource.Name {
				annotations := res.GetAnnotations()
				if annotations == nil {
					annotations = make(map[string]string)
				}
				annotations["storageclass.kubernetes.io/is-default-class"] = "true"
				res.SetAnnotations(annotations)
				needUpdate = true
			} else {
				if res.Annotations != nil {
					if isDefault, ok := res.Annotations["storageclass.kubernetes.io/is-default-class"]; ok && isDefault == "true" {
						delete(res.Annotations, "storageclass.kubernetes.io/is-default-class")
						needUpdate = true
					}
				}
			}
			if needUpdate {
				if _, err := clientset.StorageClassUpdate(ctx, cluster, "", &res, req.UpdateOpts); err != nil {
					return waitgroup.NewGenericResult(res, err)
				}
			}
			return waitgroup.NewGenericResult(res, nil)
		})

	for _, v := range wg.GetResults() {
		if v.Data.Name == req.Resource.Name && v.Error != nil {
			return err
		}
	}

	return err
}

func (k *kubernetesService) StorageClassGet(ctx context.Context, req *iapiserver.StorageClassGetRequest) (*iapiserver.StorageClassInfo, error) {
	if err := libkubernetes.ValidateClusterScopeParameters(req.Name); err != nil {
		return nil, errors.WithStack(err)
	}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.StorageClassGet(ctx, cluster, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return iapiserver.NewStorageClassInfo(meta, cluster), nil
}

func (k *kubernetesService) PersistentVolumeListAll(ctx context.Context, req *iapiserver.PersistentVolumeListRequest) (*iapiserver.PersistentVolumeListResponse, error) {
	resp := &iapiserver.PersistentVolumeListResponse{}
	clusters, err := getVisitScope(ctx, k.store, req.ResourceListRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	wg := waitgroup.RunGenericConcurrently[*iapiserver.Cluster, iapiserver.EachResourceRangeListState[*iapiserver.PersistentVolumeInfo]](
		ctx, clusters, func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.PersistentVolumeInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.PersistentVolumeInfo](cluster.ID, cluster.Name)
			resList, err := clientset.PersistentVolumeList(ctx, cluster, req.Namespace, metav1.ListOptions{FieldSelector: req.FieldSelector, LabelSelector: req.LabelSelector})
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.PersistentVolumeInfo
			for i := range resList.Items {
				resInfo := iapiserver.NewPersistentVolumeInfo(&resList.Items[i], cluster)
				if req.StorageClassFilter != "" && resInfo.Resource.Spec.StorageClassName != req.StorageClassFilter {
					continue
				}
				if filterPersistentVolume(resInfo, req.Fuzzy, req.StorageClassFilter) {
					continue
				}
				resInfos = append(resInfos, resInfo)
			}
			clusterListOne.TotalCount = len(resInfos)
			clusterListOne.List = resInfos
			return waitgroup.NewGenericResult(clusterListOne, err)
		})

	resp.EachRangeListState = wg.GetSuccessResultList()
	resp.TotalCount = CutPagingSliceResourceList[*iapiserver.PersistentVolumeInfo](resp.EachRangeListState, &resp.List, req.PageNum, req.PageSize, func(i, j int) bool {
		return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
	})
	return resp, nil
}

func (k *kubernetesService) PersistentVolumeCreate(ctx context.Context, req *iapiserver.PersistentVolumeRequest) (*iapiserver.PersistentVolumeInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.PersistentVolumeCreate(ctx, cluster, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewPersistentVolumeInfo(meta, cluster), nil
}

func (k *kubernetesService) PersistentVolumeUpdate(ctx context.Context, req *iapiserver.PersistentVolumeRequest) (*iapiserver.PersistentVolumeInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.PersistentVolumeUpdate(ctx, cluster, req.Resource.Namespace, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return iapiserver.NewPersistentVolumeInfo(meta, cluster), nil
}

func (k *kubernetesService) PersistentVolumeGet(ctx context.Context, req *iapiserver.PersistentVolumeGetRequest) (*iapiserver.PersistentVolumeInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.PersistentVolumeGet(ctx, cluster, req.Name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return iapiserver.NewPersistentVolumeInfo(meta, cluster), nil
}

func (k *kubernetesService) PersistentVolumeDelete(ctx context.Context, req *iapiserver.PersistentVolumeRequest) (*iapiserver.PersistentVolumeInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := clientset.PersistentVolumeDelete(ctx, cluster, req.Resource.Namespace, req.Resource.Name, req.DeleteOpts); err != nil {
		return nil, errors.WithStack(err)
	}

	return nil, errors.WithStack(err)
}

func (k *kubernetesService) PersistentVolumeBatchDelete(ctx context.Context, req *iapiserver.PersistentVolumeBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.PersistentVolumeRequest] {
	wg := waitgroup.RunGenericConcurrently(ctx, req.Resources, func(ctx context.Context, res *iapiserver.PersistentVolumeRequest) waitgroup.GenericResult[*iapiserver.PersistentVolumeRequest] {
		cluster, err := k.store.Clusters().Get(ctx, res.Cluster)
		if err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		if err := clientset.PersistentVolumeDelete(ctx, cluster, res.Resource.Namespace, res.Resource.Name, res.DeleteOpts); err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		return waitgroup.NewGenericResult(res, nil)
	})
	return wg.BatchGenericOutput()
}

func filterStorageClass(resInfo *iapiserver.StorageClassInfo, fuzzy string) bool {
	if fuzzy == "" {
		return false
	}

	if resInfo.Resource.AllowVolumeExpansion != nil {
		if fuzzy == "可扩展" || fuzzy == "不可扩展" {
			if fuzzy == "可扩展" && !*resInfo.Resource.AllowVolumeExpansion {
				return false
			}
			if fuzzy == "不可扩展" && *resInfo.Resource.AllowVolumeExpansion {
				return false
			}
			//do not return,check if other filter match
		}
	}
	var reclaimPolicy, bindPolicy string
	if resInfo.Resource.ReclaimPolicy != nil {
		reclaimPolicy = string(*resInfo.Resource.ReclaimPolicy)
	}
	if resInfo.Resource.VolumeBindingMode != nil {
		bindPolicy = string(*resInfo.Resource.VolumeBindingMode)
	}
	return NewObjectCommonFieldFilter(resInfo.Resource).AddField(resInfo.Resource.Provisioner).AddField(reclaimPolicy).AddField(bindPolicy).Filter(fuzzy)
}

func filterPersistentVolume(resource *iapiserver.PersistentVolumeInfo, fuzzy string, storageclasFilter string) bool {
	if fuzzy == "" && storageclasFilter == "" {
		return false
	}
	resInfo := resource.Resource

	if storageclasFilter != "" && resInfo.Spec.StorageClassName != storageclasFilter {
		return true
	}
	var accessModes []string
	for _, accessMode := range resInfo.Spec.AccessModes {
		accessModes = append(accessModes, string(accessMode))
	}

	var volumeMode string
	if resInfo.Spec.VolumeMode != nil {
		volumeMode = string(*resInfo.Spec.VolumeMode)
	}
	var capacity string
	if resInfo.Spec.Capacity.Storage() != nil {
		capacity = resInfo.Spec.Capacity.Storage().String()
	}
	var claimRef string
	if resInfo.Spec.ClaimRef != nil {
		claimRef = resInfo.Spec.ClaimRef.Namespace + "/" + resInfo.Spec.ClaimRef.Name
	}

	return NewObjectCommonFieldFilter(resInfo).
		AddField(accessModes...).
		AddField(resInfo.Spec.StorageClassName).
		AddField(string(resInfo.Spec.PersistentVolumeReclaimPolicy)).
		AddField(string(resInfo.Status.Phase)).
		AddField(volumeMode).
		AddField(capacity).
		AddField(claimRef).
		Filter(fuzzy)
}

func (k *kubernetesService) PersistentVolumeClaimCreate(ctx context.Context, req *iapiserver.PersistentVolumeClaimRequest) (*iapiserver.PersistentVolumeClaimInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.PersistentVolumeClaimCreate(ctx, cluster, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sPersistentVolumeClaimToApi(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) PersistentVolumeClaimDelete(ctx context.Context, req *iapiserver.PersistentVolumeClaimRequest) (*iapiserver.PersistentVolumeClaimInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resList, err := clientset.PodList(ctx, cluster, req.Resource.Namespace, req.ListOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range resList.Items {
		podInfo := convertK8sPodToApi(&resList.Items[i], cluster, nil)
		if !filterPod(podInfo, "", req.Resource.Name, "", "") {
			return nil, errors.Errorf("pvc still mount by pod")
		}
	}

	if err := clientset.PersistentVolumeClaimDelete(ctx, cluster, req.Resource.Namespace, req.Resource.Name, req.DeleteOpts); err != nil {
		return nil, errors.WithStack(err)
	}

	return nil, errors.WithStack(err)
}

func (k *kubernetesService) PersistentVolumeClaimClaimBatchDelete(ctx context.Context, req *iapiserver.PersistentVolumeClaimBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.PersistentVolumeClaimRequest] {
	wg := waitgroup.RunGenericConcurrently(ctx, req.Resources, func(ctx context.Context, res *iapiserver.PersistentVolumeClaimRequest) waitgroup.GenericResult[*iapiserver.PersistentVolumeClaimRequest] {
		cluster, err := k.store.Clusters().Get(ctx, res.Cluster)
		if err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		if err := clientset.PersistentVolumeClaimDelete(ctx, cluster, res.Resource.Namespace, res.Resource.Name, res.DeleteOpts); err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		return waitgroup.NewGenericResult(res, nil)
	})
	return wg.BatchGenericOutput()
}

func (k *kubernetesService) PersistentVolumeClaimUpdate(ctx context.Context, req *iapiserver.PersistentVolumeClaimRequest) (*iapiserver.PersistentVolumeClaimInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.PersistentVolumeClaimUpdate(ctx, cluster, req.Resource.Namespace, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return convertK8sPersistentVolumeClaimToApi(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) PersistentVolumeClaimExpand(ctx context.Context, req *iapiserver.PersistentVolumeClaimRequest) (*iapiserver.PersistentVolumeClaimInfo, error) {
	if req.Resource.Spec.Resources.Requests.Storage() == nil {
		return nil, errors.Errorf("unspecific storage size")
	}

	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.PersistentVolumeClaimGet(ctx, cluster, req.Resource.Namespace, req.Resource.Name, req.GetOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if meta.Spec.StorageClassName != nil {
		sc, err := clientset.StorageClassGet(ctx, cluster, *meta.Spec.StorageClassName, metav1.GetOptions{})
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if sc.AllowVolumeExpansion == nil || !*sc.AllowVolumeExpansion {
			return nil, errors.Errorf("storageclass not support expand")
		}
	}

	if req.Resource.Spec.Resources.Requests.Storage() != nil &&
		req.Resource.Spec.Resources.Requests.Storage().Size() < meta.Spec.Resources.Requests.Storage().Size() {
		return nil, errors.Errorf("request storage size smalller than current")
	}

	meta.ResourceVersion = ""
	meta.Spec.Resources.Requests = req.Resource.Spec.Resources.Requests

	meta, err = clientset.PersistentVolumeClaimUpdate(ctx, cluster, req.Resource.Namespace, meta, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sPersistentVolumeClaimToApi(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) PersistentVolumeClaimGet(ctx context.Context, req *iapiserver.PersistentVolumeClaimGetRequest) (*iapiserver.PersistentVolumeClaimInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := clientset.PersistentVolumeClaimGet(ctx, cluster, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sPersistentVolumeClaimToApi(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) PersistentVolumeClaimListAll(ctx context.Context, req *iapiserver.PersistentVolumeClaimListRequest) (*iapiserver.PersistentVolumeClaimListResponse, error) {
	resp := &iapiserver.PersistentVolumeClaimListResponse{}

	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList(ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.PersistentVolumeClaimInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.PersistentVolumeClaimInfo](cluster.ID, cluster.Name)
			resList, err := clientset.PersistentVolumeClaimList(ctx, cluster, req.Namespace, metav1.ListOptions{FieldSelector: req.FieldSelector, LabelSelector: req.LabelSelector})
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.PersistentVolumeClaimInfo
			for i := range resList.Items {
				resInfo := convertK8sPersistentVolumeClaimToApi(&resList.Items[i], cluster, req.Yaml)

				if filterPersistentVolumeClaim(resInfo, req.Fuzzy) {
					continue
				}
				resInfos = append(resInfos, resInfo)
			}
			clusterListOne.TotalCount = len(resInfos)
			clusterListOne.List = resInfos
			return waitgroup.NewGenericResult(clusterListOne, err)
		}, func(i, j int) bool {
			return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
		}, 10*time.Second)
	return resp, err
}

func convertK8sPersistentVolumeClaimToApi(meta *v1.PersistentVolumeClaim, cluster *iapiserver.Cluster, yaml bool) *iapiserver.PersistentVolumeClaimInfo {
	resp := iapiserver.NewPersistentVolumeClaimInfo(meta, cluster)

	if !yaml {
		allowExpansion, allowSnapshot := getVolumeAttribute(context.Background(), cluster, meta)
		resp.AllowExpansion = &allowExpansion
		resp.AllowSnapshot = &allowSnapshot
	}
	return resp
}

func filterPersistentVolumeClaim(pvci *iapiserver.PersistentVolumeClaimInfo, fuzzy string) bool {
	if fuzzy == "" {
		return false
	}
	resInfo := pvci.Resource
	var storageclasFilter string
	if resInfo.Spec.StorageClassName != nil {
		storageclasFilter = *resInfo.Spec.StorageClassName
	}

	var accessModes []string
	for _, accessMode := range resInfo.Spec.AccessModes {
		accessModes = append(accessModes, string(accessMode))
	}

	var volumeMode string
	if resInfo.Spec.VolumeMode != nil {
		volumeMode = string(*resInfo.Spec.VolumeMode)
	}
	var reqeustCapcity string
	if resInfo.Spec.Resources.Requests.Storage() != nil {
		reqeustCapcity = resInfo.Spec.Resources.Requests.Storage().String()
	}
	var bindCapacity string
	if resInfo.Status.Capacity.Storage() != nil {
		reqeustCapcity = resInfo.Status.Capacity.Storage().String()
	}

	return NewObjectCommonFieldFilter(resInfo).
		AddField(accessModes...).
		AddField(storageclasFilter).
		AddField(volumeMode).
		AddField(string(resInfo.Status.Phase)).
		AddField(resInfo.Spec.VolumeName).
		AddField(reqeustCapcity).
		AddField(bindCapacity).
		Filter(fuzzy)
}

func getVolumeAttribute(ctx context.Context, cluster *iapiserver.Cluster, pvc *v1.PersistentVolumeClaim) (bool, bool) {
	var allowExpansion, allowSnapshot bool
	if pvc.Spec.StorageClassName != nil {
		sc, err := clientset.StorageClassGet(ctx, cluster, *pvc.Spec.StorageClassName, metav1.GetOptions{})
		if err == nil {
			if sc.AllowVolumeExpansion != nil {
				//sc allow and pvc is bound
				if *sc.AllowVolumeExpansion && pvc.Status.Phase == v1.ClaimBound {
					allowExpansion = true
				}
			}
			return allowExpansion, allowSnapshot
		}
		log.Errorf("pvc %v/%v's storage class %v not found", pvc.Namespace, pvc.Name, *pvc.Spec.StorageClassName)
	}
	return allowExpansion, allowSnapshot
}

func (k *kubernetesService) VolumeSnapshotClassListAll(ctx context.Context, req *iapiserver.VolumeSnapshotClassListRequest) (*iapiserver.VolumeSnapshotClassListResponse, error) {
	resp := &iapiserver.VolumeSnapshotClassListResponse{}

	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList(ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.VolumeSnapshotClassInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.VolumeSnapshotClassInfo](cluster.ID, cluster.Name)
			resList, err := libkubernetes.VolumeSnapshotClassList(ctx, cluster.Config, req.Namespace, req.ToListOpts())
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.VolumeSnapshotClassInfo
			for i := range resList.Items {
				resInfo := convertK8sVolumeSnapshotClassToApiVolumeSnapshotClass(&resList.Items[i], cluster, req.Yaml)

				resInfos = append(resInfos, resInfo)
			}
			clusterListOne.TotalCount = len(resInfos)
			clusterListOne.List = resInfos
			return waitgroup.NewGenericResult(clusterListOne, err)
		}, func(i, j int) bool {
			return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
		}, 10*time.Second)
	return resp, err
}

func convertK8sVolumeSnapshotClassToApiVolumeSnapshotClass(meta *snapshotv1beta1.VolumeSnapshotClass, cluster *iapiserver.Cluster, yaml bool) *iapiserver.VolumeSnapshotClassInfo {
	return &iapiserver.VolumeSnapshotClassInfo{
		Resource: meta,
	}
}

func (k *kubernetesService) VolumeSnapshotClassGet(ctx context.Context, req *iapiserver.VolumeSnapshotClassGetRequest) (*iapiserver.VolumeSnapshotClassInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.VolumeSnapshotClassGet(ctx, cluster.Config, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return convertK8sVolumeSnapshotClassToApiVolumeSnapshotClass(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) VolumeSnapshotClassCreate(ctx context.Context, req *iapiserver.VolumeSnapshotClassRequest) (*iapiserver.VolumeSnapshotClassInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.VolumeSnapshotClassCreate(ctx, cluster.Config, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sVolumeSnapshotClassToApiVolumeSnapshotClass(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) VolumeSnapshotClassDelete(ctx context.Context, req *iapiserver.VolumeSnapshotClassRequest) (*iapiserver.VolumeSnapshotInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := libkubernetes.VolumeSnapshotClassDelete(ctx, cluster.Config, req.Resource.Name, req.DeleteOpts); err != nil {
		return nil, errors.WithStack(err)
	}
	return nil, nil
}

func (k *kubernetesService) VolumeSnapshotClassBatchDelete(ctx context.Context, req *iapiserver.VolumeSnapshotClassBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.VolumeSnapshotClassRequest] {
	wg := waitgroup.RunGenericConcurrently(ctx, req.Resources, func(ctx context.Context, res *iapiserver.VolumeSnapshotClassRequest) waitgroup.GenericResult[*iapiserver.VolumeSnapshotClassRequest] {
		cluster, err := k.store.Clusters().Get(ctx, res.Cluster)
		if err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		if err := libkubernetes.VolumeSnapshotClassDelete(ctx, cluster.Config, res.Resource.Name, res.DeleteOpts); err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		return waitgroup.NewGenericResult(res, nil)
	})
	return wg.BatchGenericOutput()
}

func (k *kubernetesService) VolumeSnapshotClassUpdate(ctx context.Context, req *iapiserver.VolumeSnapshotClassRequest) (*iapiserver.VolumeSnapshotClassInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.VolumeSnapshotClassUpdate(ctx, cluster.Config, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sVolumeSnapshotClassToApiVolumeSnapshotClass(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) VolumeSnapshotContentListAll(ctx context.Context, req *iapiserver.VolumeSnapshotContentListRequest) (*iapiserver.VolumeSnapshotContentListResponse, error) {
	resp := &iapiserver.VolumeSnapshotContentListResponse{}

	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList(ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.VolumeSnapshotContentInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.VolumeSnapshotContentInfo](cluster.ID, cluster.Name)
			resList, err := libkubernetes.VolumeSnapshotContentList(ctx, cluster.Config, req.Namespace, req.ToListOpts())
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.VolumeSnapshotContentInfo
			for i := range resList.Items {
				resInfo := convertK8sVolumeSnapshotContentToApiVolumeSnapshotContent(&resList.Items[i], cluster, req.Yaml)
				if NewObjectCommonFieldFilter(resInfo.Resource).Filter(req.Fuzzy) {
					continue
				}
				resInfos = append(resInfos, resInfo)
			}
			clusterListOne.TotalCount = len(resInfos)
			clusterListOne.List = resInfos
			return waitgroup.NewGenericResult(clusterListOne, err)
		}, func(i, j int) bool {
			return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
		}, 10*time.Second)

	return resp, err
}

func convertK8sVolumeSnapshotContentToApiVolumeSnapshotContent(meta *snapshotv1beta1.VolumeSnapshotContent, cluster *iapiserver.Cluster, yaml bool) *iapiserver.VolumeSnapshotContentInfo {
	return &iapiserver.VolumeSnapshotContentInfo{
		Resource: meta,
	}
}

func (k *kubernetesService) VolumeSnapshotContentGet(ctx context.Context, req *iapiserver.VolumeSnapshotContentGetRequest) (*iapiserver.VolumeSnapshotContentInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.VolumeSnapshotContentGet(ctx, cluster.Config, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return convertK8sVolumeSnapshotContentToApiVolumeSnapshotContent(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) VolumeSnapshotContentCreate(ctx context.Context, req *iapiserver.VolumeSnapshotContentRequest) (*iapiserver.VolumeSnapshotContentInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.VolumeSnapshotContentCreate(ctx, cluster.Config, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sVolumeSnapshotContentToApiVolumeSnapshotContent(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) VolumeSnapshotContentDelete(ctx context.Context, req *iapiserver.VolumeSnapshotContentRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := libkubernetes.VolumeSnapshotContentDelete(ctx, cluster.Config, req.Resource.Name, req.DeleteOpts); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) VolumeSnapshotContentBatchDelete(ctx context.Context, req *iapiserver.VolumeSnapshotContentBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.VolumeSnapshotContentRequest] {
	wg := waitgroup.RunGenericConcurrently(ctx, req.Resources, func(ctx context.Context, res *iapiserver.VolumeSnapshotContentRequest) waitgroup.GenericResult[*iapiserver.VolumeSnapshotContentRequest] {
		cluster, err := k.store.Clusters().Get(ctx, res.Cluster)
		if err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		if err := libkubernetes.VolumeSnapshotContentDelete(ctx, cluster.Config, res.Resource.Name, res.DeleteOpts); err != nil {
			return waitgroup.NewGenericResult(res, err)
		}
		return waitgroup.NewGenericResult(res, nil)
	})
	return wg.BatchGenericOutput()
}

func (k *kubernetesService) VolumeSnapshotContentUpdate(ctx context.Context, req *iapiserver.VolumeSnapshotContentRequest) (*iapiserver.VolumeSnapshotContentInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.VolumeSnapshotContentUpdate(ctx, cluster.Config, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sVolumeSnapshotContentToApiVolumeSnapshotContent(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) VolumeSnapshotListAll(ctx context.Context, req *iapiserver.VolumeSnapshotListRequest) (*iapiserver.VolumeSnapshotListResponse, error) {
	resp := &iapiserver.VolumeSnapshotListResponse{}

	var err error
	resp.EachRangeListState, resp.TotalCount, err = multiClusterResourceList(ctx, k.store, &resp.List, req.ResourceListRequest,
		func(ctx context.Context, cluster *iapiserver.Cluster) waitgroup.GenericResult[iapiserver.EachResourceRangeListState[*iapiserver.VolumeSnapshotInfo]] {
			clusterListOne := iapiserver.NewEachResourceRangeListState[*iapiserver.VolumeSnapshotInfo](cluster.ID, cluster.Name)
			resList, err := libkubernetes.VolumeSnapshotList(ctx, cluster.Config, req.Namespace, req.ToListOpts())
			if err != nil {
				return waitgroup.NewGenericResult(clusterListOne, err)
			}

			var resInfos []*iapiserver.VolumeSnapshotInfo
			for i := range resList.Items {
				resInfo := convertK8sVolumeSnapshotToApiVolumeSnapshot(&resList.Items[i], cluster, req.Yaml)
				if filterVolumeSnapshot(resInfo, req.Fuzzy, req.PersistentVolumeClaimName) {
					continue
				}
				resInfos = append(resInfos, resInfo)
			}
			clusterListOne.TotalCount = len(resInfos)
			clusterListOne.List = resInfos
			return waitgroup.NewGenericResult(clusterListOne, err)
		}, func(i, j int) bool {
			return sortWithCommonObjectParam(resp.List[i].Resource, resp.List[j].Resource, req.SortBy, req.SortDesc)
		}, 10*time.Second)

	return resp, err
}

func convertK8sVolumeSnapshotToApiVolumeSnapshot(meta *snapshotv1beta1.VolumeSnapshot, cluster *iapiserver.Cluster, yaml bool) *iapiserver.VolumeSnapshotInfo {
	return &iapiserver.VolumeSnapshotInfo{
		Resource: meta,
	}
}

func (k *kubernetesService) VolumeSnapshotGet(ctx context.Context, req *iapiserver.VolumeSnapshotGetRequest) (*iapiserver.VolumeSnapshotInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.VolumeSnapshotGet(ctx, cluster.Config, req.Namespace, req.Name, req.ToGetOpts())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sVolumeSnapshotToApiVolumeSnapshot(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) VolumeSnapshotCreate(ctx context.Context, req *iapiserver.VolumeSnapshotRequest) (*iapiserver.VolumeSnapshotInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.VolumeSnapshotCreate(ctx, cluster.Config, req.Resource.Namespace, req.Resource, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sVolumeSnapshotToApiVolumeSnapshot(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) VolumeSnapshotDelete(ctx context.Context, req *iapiserver.VolumeSnapshotRequest) error {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := libkubernetes.VolumeSnapshotDelete(ctx, cluster.Config, req.Resource.Namespace, req.Resource.Name, req.DeleteOpts); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (k *kubernetesService) VolumeSnapshotBatchDelete(ctx context.Context, req *iapiserver.VolumeSnapshotBatchRequest) waitgroup.BatchGenericOutput[*iapiserver.VolumeSnapshotRequest] {
	wg := waitgroup.RunGenericConcurrently[*iapiserver.VolumeSnapshotRequest, *iapiserver.VolumeSnapshotRequest](
		ctx, req.Resources, func(ctx context.Context, res *iapiserver.VolumeSnapshotRequest) waitgroup.GenericResult[*iapiserver.VolumeSnapshotRequest] {
			cluster, err := k.store.Clusters().Get(ctx, res.Cluster)
			if err != nil {
				return waitgroup.NewGenericResult(res, err)
			}
			if err := libkubernetes.VolumeSnapshotDelete(ctx, cluster.Config, res.Resource.Namespace, res.Resource.Name, res.DeleteOpts); err != nil {
				return waitgroup.NewGenericResult(res, err)
			}
			return waitgroup.NewGenericResult(res, nil)
		})
	return wg.BatchGenericOutput()
}

func (k *kubernetesService) VolumeSnapshotUpdate(ctx context.Context, req *iapiserver.VolumeSnapshotRequest) (*iapiserver.VolumeSnapshotInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.VolumeSnapshotUpdate(ctx, cluster.Config, req.Resource.Namespace, req.Resource, req.UpdateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return convertK8sVolumeSnapshotToApiVolumeSnapshot(meta, cluster, req.Yaml), nil
}

func (k *kubernetesService) VolumeSnapshotClone(ctx context.Context, req *iapiserver.VolumeSnapshotRequest) (*iapiserver.PersistentVolumeClaimInfo, error) {
	cluster, err := k.store.Kubernetes().Get(ctx, req.Cluster)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, err := libkubernetes.VolumeSnapshotGet(ctx, cluster.Config, req.Resource.Namespace, req.Resource.Name, req.GetOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if meta.Status == nil || meta.Status.ReadyToUse == nil || !*meta.Status.ReadyToUse {
		return nil, errors.Errorf("snapshot is not ready")
	}

	var sourcePvcStorageClass *string
	// if pvc not specific, we don't known which storageclass should set
	if meta.Spec.Source.PersistentVolumeClaimName != nil {
		sourcePVC, err := libkubernetes.PersistentVolumeClaimGet(ctx, cluster.Config, req.Resource.Namespace, *meta.Spec.Source.PersistentVolumeClaimName, req.GetOpts)
		if err != nil {

			return nil, errors.WithStack(err)
		}
		sourcePvcStorageClass = sourcePVC.Spec.StorageClassName
	}

	snapshotApiGroup := "snapshot.storage.k8s.io"
	newpvc := &v1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.PersistentVolumeClaimName,
			Namespace: meta.Namespace,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			DataSource: &v1.TypedLocalObjectReference{
				Name:     meta.Name,
				Kind:     "VolumeSnapshot",
				APIGroup: &snapshotApiGroup,
			},
			AccessModes: req.PersistentVolumeClaimAccessModes, // must t set access mode
			Resources: v1.VolumeResourceRequirements{
				Requests: v1.ResourceList{
					"storage": *meta.Status.RestoreSize, // must set size
				},
			},
			StorageClassName: sourcePvcStorageClass,
		},
	}

	newpvc, err = libkubernetes.PersistentVolumeClaimCreate(ctx, cluster.Config, newpvc.Namespace, newpvc, req.CreateOpts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return convertK8sPersistentVolumeClaimToApi(newpvc, cluster, req.Yaml), nil
}

func filterVolumeSnapshot(resInfo *iapiserver.VolumeSnapshotInfo, fuzzy string, persistentVolumeFilter string) bool {
	if fuzzy == "" && persistentVolumeFilter == "" {
		return false
	}

	if persistentVolumeFilter != "" && resInfo.Resource.Spec.Source.PersistentVolumeClaimName != nil && *resInfo.Resource.Spec.Source.PersistentVolumeClaimName != persistentVolumeFilter {
		return true
	}

	var ready bool
	//volume snaphot status will be nil if has error
	if resInfo.Resource.Status != nil {
		if resInfo.Resource.Status.ReadyToUse != nil {
			ready = *resInfo.Resource.Status.ReadyToUse
		}
	}
	if fuzzy == "就绪" || fuzzy == "未就绪" {
		if fuzzy == "就绪" && !ready {
			return true
		}
		if fuzzy == "未就绪" && ready {
			return true
		}
	}

	var content string
	if resInfo.Resource.Status != nil {
		if resInfo.Resource.Status.BoundVolumeSnapshotContentName != nil {
			content = *resInfo.Resource.Status.BoundVolumeSnapshotContentName
		}
	}

	var snapshotClass string
	if resInfo.Resource.Spec.VolumeSnapshotClassName != nil {
		snapshotClass = *resInfo.Resource.Spec.VolumeSnapshotClassName
	}
	return NewObjectCommonFieldFilter(resInfo.Resource).
		AddField(content).
		AddField(snapshotClass).
		Filter(fuzzy)
}
