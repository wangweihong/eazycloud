package kubernetes

import (
	"sort"
	"strings"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/gotoolbox/pkg/compareutil"
	"github.com/wangweihong/gotoolbox/pkg/paging"
	"github.com/wangweihong/gotoolbox/pkg/sets"
	"github.com/wangweihong/gotoolbox/pkg/sortutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// import (
// 	"context"
// 	"fmt"
// 	"sync"

// 	"github.com/wangweihong/eazycloud/apis/iapiserver"
// 	"github.com/wangweihong/eazycloud/internal/pkg/clientset"
// 	"github.com/wangweihong/gotoolbox/pkg/mathutil"
// 	"github.com/wangweihong/gotoolbox/pkg/sets"
// 	"github.com/wangweihong/gotoolbox/pkg/waitgroup"
// 	v1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/kubectl/pkg/util/resource"
// 	"sigs.k8s.io/yaml"
// )

// func convertResourceLimitToPersistentUnit(containerName string, requests, limits v1.ResourceList) *iapiserver.ResourceConvert {
// 	ResourceConvert := &iapiserver.ResourceConvert{
// 		Container: containerName,
// 		Limits:    make(map[string]int64),
// 		Request:   make(map[string]int64),
// 	}
// 	if requests != nil {
// 		for i, j := range requests {
// 			if string(i) == "memory" {
// 				ResourceConvert.Request[string(i)] = j.Value() / 1024 / 1024 // convert to memory
// 				continue
// 			}
// 			if string(i) == "cpu" { // if use Value(), 0.1 cpu/100m cpu will convert to 1 cpu
// 				ResourceConvert.Request[string(i)] = j.MilliValue()
// 				continue
// 			}
// 			ResourceConvert.Request[string(i)] = j.Value()
// 		}
// 	}

// 	if limits != nil {
// 		for i, j := range limits {
// 			if string(i) == "memory" {
// 				ResourceConvert.Limits[string(i)] = j.Value() / 1024 / 1024 // convert to memory
// 				continue
// 			}
// 			if string(i) == "cpu" { // if use Value(), 0.1 cpu/100m cpu will convert to 1 cpu
// 				ResourceConvert.Limits[string(i)] = j.MilliValue()
// 				continue
// 			}
// 			ResourceConvert.Limits[string(i)] = j.Value()
// 		}
// 	}
// 	return ResourceConvert
// }

// func getResourceController(uuid string, childParent map[string]map[string]struct{}, resourceInfo map[string]*iapiserver.ObjectTypeMeta) *iapiserver.ObjectTypeMeta {
// 	if uuid == "" || childParent == nil || resourceInfo == nil {
// 		return nil
// 	}

// 	parent := ""
// 	parents := childParent[uuid]
// 	for len(parents) != 0 {
// 		for parent = range parents {
// 			break
// 		}
// 		parents = childParent[parent]
// 	}

// 	return resourceInfo[parent]
// }

// func getAllControllers(ctx context.Context, clusters []*iapiserver.Cluster, req iapiserver.ListRequest) (map[string]map[string]struct{}, map[string]*iapiserver.ObjectTypeMeta, error) {
// 	childParent := make(map[string]map[string]struct{})
// 	resourceInfo := make(map[string]*iapiserver.ObjectTypeMeta)

// 	kinds := map[string]bool{"Pod": true, "Deployment": true, "StatefulSet": true, "DaemonSet": true, "Job": true, "CronJob": true, "ReplicaSet": true}

// 	wg := waitgroup.NewWaitGroup(nil)
// 	glock := sync.Mutex{}
// 	metaList := []iapiserver.ObjectTypeMeta{}
// 	for _, cluster := range clusters {
// 		cluster := cluster
// 		wg.Start(waitgroup.NewWaitGroupHandleFunc(ctx, "", func() waitgroup.Result {
// 			resourceList, err := clientset.DeploymentList(ctx, cluster, req.Namespace, metav1.ListOptions{})
// 			if err != nil {
// 				return waitgroup.Result{}
// 			}
// 			glock.Lock()
// 			defer glock.Unlock()
// 			for _, resource := range resourceList.Items {
// 				metaList = append(metaList, iapiserver.ObjectTypeMeta{ObjectMeta: resource.ObjectMeta, TypeMeta: resource.TypeMeta})
// 			}
// 			return waitgroup.Result{}
// 		}))

// 		wg.Start(waitgroup.NewWaitGroupHandleFunc(ctx, "", func() waitgroup.Result {
// 			resourceList, err := clientset.StatefulSetList(ctx, cluster, req.Namespace, metav1.ListOptions{})
// 			if err != nil {
// 				return waitgroup.Result{}
// 			}
// 			glock.Lock()
// 			defer glock.Unlock()
// 			for _, resource := range resourceList.Items {
// 				metaList = append(metaList, iapiserver.ObjectTypeMeta{ObjectMeta: resource.ObjectMeta, TypeMeta: resource.TypeMeta})
// 			}
// 			return waitgroup.Result{}
// 		}))
// 		wg.Start(waitgroup.NewWaitGroupHandleFunc(ctx, "", func() waitgroup.Result {
// 			resourceList, err := clientset.ReplicaSetList(ctx, cluster, req.Namespace, metav1.ListOptions{})
// 			if err != nil {
// 				return waitgroup.Result{}
// 			}
// 			glock.Lock()
// 			defer glock.Unlock()
// 			for _, resource := range resourceList.Items {
// 				metaList = append(metaList, iapiserver.ObjectTypeMeta{ObjectMeta: resource.ObjectMeta, TypeMeta: resource.TypeMeta})
// 			}
// 			return waitgroup.Result{}
// 		}))

// 		wg.Start(waitgroup.NewWaitGroupHandleFunc(ctx, "", func() waitgroup.Result {
// 			resourceList, err := clientset.DaemonSetList(ctx, cluster, req.Namespace, metav1.ListOptions{})
// 			if err != nil {
// 				return waitgroup.Result{}
// 			}
// 			glock.Lock()
// 			defer glock.Unlock()
// 			for _, resource := range resourceList.Items {
// 				metaList = append(metaList, iapiserver.ObjectTypeMeta{ObjectMeta: resource.ObjectMeta, TypeMeta: resource.TypeMeta})
// 			}
// 			return waitgroup.Result{}
// 		}))
// 		wg.Start(waitgroup.NewWaitGroupHandleFunc(ctx, "", func() waitgroup.Result {
// 			resourceList, err := clientset.CronJobList(ctx, cluster, req.Namespace, metav1.ListOptions{})
// 			if err != nil {
// 				return waitgroup.Result{}
// 			}
// 			glock.Lock()
// 			defer glock.Unlock()
// 			for _, resource := range resourceList.Items {
// 				metaList = append(metaList, iapiserver.ObjectTypeMeta{ObjectMeta: resource.ObjectMeta, TypeMeta: resource.TypeMeta})
// 			}
// 			return waitgroup.Result{}
// 		}))
// 		wg.Start(waitgroup.NewWaitGroupHandleFunc(ctx, "", func() waitgroup.Result {
// 			resourceList, err := clientset.JobList(ctx, cluster, req.Namespace, metav1.ListOptions{})
// 			if err != nil {
// 				return waitgroup.Result{}
// 			}
// 			glock.Lock()
// 			defer glock.Unlock()
// 			for _, resource := range resourceList.Items {
// 				metaList = append(metaList, iapiserver.ObjectTypeMeta{ObjectMeta: resource.ObjectMeta, TypeMeta: resource.TypeMeta})
// 			}
// 			return waitgroup.Result{}
// 		}))
// 		wg.Start(waitgroup.NewWaitGroupHandleFunc(ctx, "", func() waitgroup.Result {
// 			resourceList, err := clientset.PodList(ctx, cluster, req.Namespace, metav1.ListOptions{})
// 			if err != nil {
// 				return waitgroup.Result{}
// 			}
// 			glock.Lock()
// 			defer glock.Unlock()
// 			for _, resource := range resourceList.Items {
// 				metaList = append(metaList, iapiserver.ObjectTypeMeta{ObjectMeta: resource.ObjectMeta, TypeMeta: resource.TypeMeta})
// 			}
// 			return waitgroup.Result{}
// 		}))
// 	}
// 	wg.Wait()

// 	for _, meta := range metaList {
// 		meta := meta
// 		resourceInfo[string(meta.ObjectMeta.UID)] = &meta
// 		for _, parent := range meta.ObjectMeta.OwnerReferences {
// 			if !kinds[parent.Kind] {
// 				continue
// 			}
// 			parents, ok := childParent[string(meta.ObjectMeta.UID)]
// 			if !ok {
// 				parents = map[string]struct{}{}
// 				childParent[string(meta.ObjectMeta.UID)] = parents
// 			}
// 			parents[string(parent.UID)] = struct{}{}
// 		}
// 	}

// 	return childParent, resourceInfo, nil
// }

// func filterPod(podInfo *iapiserver.PodInfo, fuzzy string, FilterVolumeName string, FilterConfigMap string, FilterSecret string) bool {
// 	if FilterVolumeName != "" {
// 		isFind := false
// 		for _, v1 := range podInfo.Spec.Volumes {
// 			if v1.PersistentVolumeClaim != nil && v1.PersistentVolumeClaim.ClaimName == FilterVolumeName {
// 				isFind = true
// 				break
// 			}
// 		}
// 		if !isFind {
// 			return true
// 		}
// 	}

// 	if FilterConfigMap != "" {
// 		isFind := false
// 		for _, v1 := range podInfo.Spec.Volumes {
// 			if v1.ConfigMap != nil && v1.ConfigMap.Name == FilterConfigMap {
// 				isFind = true
// 				break
// 			}
// 		}
// 		if !isFind {
// 			return true
// 		}
// 	}

// 	if FilterSecret != "" {
// 		isFind := false
// 		for _, v1 := range podInfo.Spec.Volumes {
// 			if v1.Secret != nil && v1.Secret.SecretName == FilterSecret {
// 				isFind = true
// 				break
// 			}
// 			for _, pullSecret := range podInfo.Spec.ImagePullSecrets {
// 				if pullSecret.Name == FilterSecret {
// 					isFind = true
// 					break
// 				}

// 			}
// 		}
// 		if !isFind {
// 			return true
// 		}
// 	}

// 	return NewObjectCommonFieldFilter(podInfo).AddField(podInfo.Status.HostIP).AddField(podInfo.PodStatus.Status).AddField(podInfo.Status.PodIP).Filter(fuzzy)
// }

// func NewObjectCommonFieldFilter(object metav1.Object, fields ...string) *FieldFilter {
// 	return NewFieldFilter(fields...).AddObjectField(object)
// }

// type FieldFilter struct {
// 	fields []string
// }

// func NewFieldFilter(fields ...string) *FieldFilter {
// 	f := &FieldFilter{
// 		fields: make([]string, 0),
// 	}
// 	f.fields = append(f.fields, fields...)
// 	return f
// }

// func (f *FieldFilter) AddObjectField(obj metav1.Object) *FieldFilter {
// 	if obj == nil || obj.GetAnnotations() == nil {
// 		return f
// 	}

// 	if obj.GetName() != "" {
// 		f.fields = append(f.fields, obj.GetName())
// 	}

// 	if obj.GetNamespace() != "" {
// 		f.fields = append(f.fields, obj.GetNamespace())
// 	}

// 	return f
// }

// func (f *FieldFilter) AddField(field ...string) *FieldFilter {
// 	f.fields = append(f.fields, field...)
// 	return f
// }

// func (f *FieldFilter) Filter(fuzzy string) bool {
// 	if fuzzy == "" {
// 		return false
// 	}

// 	if sets.NewString(f.fields...).ContainAny(fuzzy) {
// 		return false
// 	}
// 	return true
// }

// // func Compare(a, b interface{}, Asc bool) bool {
// // 	ret := compareutil.Compare(a, b)
// // 	if ret ==0 {
// // 		return false
// // 	}

// // 	if !Asc{
// // 		return true
// // 	}
// // }

// // func CutPagingSliceFromWgResults(wg *utils.Group, listStateSlicePtr *[]iapiserver.EachResourceRangeListState, resourceSlicePtr interface{}, pageNum, pageSize int, totalCount *int, lessFunc func(i, j int) bool) {
// // 	if wg != nil && listStateSlicePtr != nil && resourceSlicePtr != nil {
// // 		for _, ret := range wg.GetResults() {
// // 			oneClusterResourceList, ok := ret.Data.(iapiserver.EachResourceRangeListState)
// // 			if !ok {
// // 				log.Errorf("CutPagingSliceFromWgResults ignore for non EachResourceRangeListState type:%v",reflect.TypeOf(ret.Data).Kind())
// // 				return
// // 			}
// // 			oneClusterResourceList.Result = ret.Error
// // 			resList := oneClusterResourceList.List
// // 			oneClusterResourceList.List = nil
// // 			*listStateSlicePtr = append(*listStateSlicePtr, oneClusterResourceList)
// // 			if ret.Error != nil {
// // 				continue
// // 			}
// // 			if resList != nil {
// // 				Append(resourceSlicePtr, resList)
// // 			}
// // 		}
// // 		total := Len(resourceSlicePtr)
// // 		if totalCount != nil {
// // 			*totalCount = total
// // 		}

// // 		sort.SliceStable(reflect.Indirect(reflect.ValueOf(resourceSlicePtr)).Interface(), lessFunc)
// // 		s, index := utils.PagingIndex(total, pageNum, pageSize)
// // 		Slice(resourceSlicePtr, s, index)
// // 		return
// // 	}
// // 	logrus.Error("CutPagingSliceFromWgResults do nothing because wg or istStateSlicePtr or resourceSlicePtr is nil")
// // }

func yamlToJson(in string) (string, error) {
	data, err := yaml.YAMLToJSON([]byte(in))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

type FieldFilter struct {
	fields []string
}

func NewObjectCommonFieldFilter(object metav1.Object, fields ...string) *FieldFilter {
	return NewFieldFilter(fields...).AddObjectField(object)
}

func NewFieldFilter(fields ...string) *FieldFilter {
	f := &FieldFilter{
		fields: make([]string, 0),
	}
	f.fields = append(f.fields, fields...)
	return f
}

func (f *FieldFilter) AddObjectField(obj metav1.Object) *FieldFilter {
	if obj == nil || obj.GetAnnotations() == nil {
		return f
	}

	if obj.GetName() != "" {
		f.fields = append(f.fields, obj.GetName())
	}

	if obj.GetNamespace() != "" {
		f.fields = append(f.fields, obj.GetNamespace())
	}
	return f
}

func (f *FieldFilter) AddField(field ...string) *FieldFilter {
	f.fields = append(f.fields, field...)
	return f
}

// true表示过滤该数据
func (f *FieldFilter) Filter(fuzzy string) bool {
	if fuzzy == "" {
		return false
	}
	fuzzyList := strings.Split(fuzzy, ",")
	if sets.NewString(f.fields...).ContainAny(fuzzyList...) {
		return false
	}
	return true
}

func filterPod(podInfo *iapiserver.PodInfo, fuzzy string, FilterVolumeName string, FilterConfigMap string, FilterSecret string) bool {
	if FilterVolumeName != "" {
		isFind := false
		for _, v1 := range podInfo.Resource.Spec.Volumes {
			if v1.PersistentVolumeClaim != nil && v1.PersistentVolumeClaim.ClaimName == FilterVolumeName {
				isFind = true
				break
			}
		}
		if !isFind {
			return true
		}
	}

	if FilterConfigMap != "" {
		isFind := false
		for _, v1 := range podInfo.Resource.Spec.Volumes {
			if v1.ConfigMap != nil && v1.ConfigMap.Name == FilterConfigMap {
				isFind = true
				break
			}
		}
		if !isFind {
			return true
		}
	}

	if FilterSecret != "" {
		isFind := false
		for _, v1 := range podInfo.Resource.Spec.Volumes {
			if v1.Secret != nil && v1.Secret.SecretName == FilterSecret {
				isFind = true
				break
			}
			for _, pullSecret := range podInfo.Resource.Spec.ImagePullSecrets {
				if pullSecret.Name == FilterSecret {
					isFind = true
					break
				}

			}
		}
		if !isFind {
			return true
		}
	}

	return NewObjectCommonFieldFilter(podInfo.Resource).AddField(podInfo.Resource.Status.HostIP).AddField(podInfo.PodStatus.Status).AddField(podInfo.Resource.Status.PodIP).Filter(fuzzy)
}

func convertResourceLimitToPersistentUnit(containerName string, requests, limits v1.ResourceList) *iapiserver.ResourceConvert {
	ResourceConvert := &iapiserver.ResourceConvert{
		Container: containerName,
		Limits:    make(map[string]int64),
		Request:   make(map[string]int64),
	}

	for i, j := range requests {
		if string(i) == "memory" {
			ResourceConvert.Request[string(i)] = j.Value() / 1024 / 1024 // convert to memory
			continue
		}
		if string(i) == "cpu" { // if use Value(), 0.1 cpu/100m cpu will convert to 1 cpu
			ResourceConvert.Request[string(i)] = j.MilliValue()
			continue
		}
		ResourceConvert.Request[string(i)] = j.Value()
	}

	for i, j := range limits {
		if string(i) == "memory" {
			ResourceConvert.Limits[string(i)] = j.Value() / 1024 / 1024 // convert to memory
			continue
		}
		if string(i) == "cpu" { // if use Value(), 0.1 cpu/100m cpu will convert to 1 cpu
			ResourceConvert.Limits[string(i)] = j.MilliValue()
			continue
		}
		ResourceConvert.Limits[string(i)] = j.Value()
	}

	return ResourceConvert
}

// func CutPagingSliceFromWgResults(wg *utils.Group, listStateSlicePtr *[]iapiserver.EachResourceRangeListState, resourceSlicePtr interface{}, pageNum, pageSize int, totalCount *int, lessFunc func(i, j int) bool) {
// 	if wg != nil && listStateSlicePtr != nil && resourceSlicePtr != nil {
// 		for _, ret := range wg.GetResults() {
// 			oneClusterResourceList, ok := ret.Data.(iapiserver.EachResourceRangeListState)
// 			if !ok {
// 				logrus.Errorf("CutPagingSliceFromWgResults ignore for non iapiserver.EachResourceRangeListState type:%v",
// 					reflect.TypeOf(ret.Data).Kind())
// 				return
// 			}
// 			oneClusterResourceList.Result = iapiserver.SetOutput(nil, ret.Error)
// 			resList := oneClusterResourceList.List
// 			oneClusterResourceList.List = nil
// 			*listStateSlicePtr = append(*listStateSlicePtr, oneClusterResourceList)
// 			if ret.Error != nil {
// 				continue
// 			}
// 			if resList != nil {
// 				Append(resourceSlicePtr, resList)
// 			}
// 		}
// 		total := Len(resourceSlicePtr)
// 		if totalCount != nil {
// 			*totalCount = total
// 		}

// 		sort.SliceStable(reflect.Indirect(reflect.ValueOf(resourceSlicePtr)).Interface(), lessFunc)
// 		s, index := utils.PagingIndex(total, pageNum, pageSize)
// 		Slice(resourceSlicePtr, s, index)
// 		return
// 	}
// 	logrus.Error("CutPagingSliceFromWgResults do nothing because wg or istStateSlicePtr or resourceSlicePtr is nil")
// }

// func Append(slicePtr interface{}, data interface{}) {
// 	if slicePtr == nil || data == nil {
// 		return
// 	}

// 	st := reflect.TypeOf(slicePtr)
// 	dt := reflect.TypeOf(data)
// 	if st.Kind() != reflect.Ptr || st.Elem().Kind() != reflect.Slice {
// 		panic("no slice pointer")
// 	}
// 	sv := reflect.ValueOf(slicePtr)
// 	indirect := reflect.Indirect(sv)
// 	if st.Elem() == dt {
// 		indirect.Set(reflect.AppendSlice(indirect, reflect.ValueOf(data)))
// 		return
// 	}

// 	if st.Elem().Elem() == dt {
// 		indirect.Set(reflect.Append(indirect, reflect.ValueOf(data)))
// 		return
// 	}
// 	panic(fmt.Sprintf("type not match,[%v] [%v]", st.String(), dt.String()))
// }

// func Len(data interface{}) int {
// 	if data != nil {
// 		dt := reflect.TypeOf(data)
// 		if dt.Kind() == reflect.Slice {
// 			return reflect.ValueOf(data).Len()
// 		}

// 		if dt.Kind() == reflect.Ptr && dt.Elem().Kind() == reflect.Slice {
// 			return reflect.Indirect(reflect.ValueOf(data)).Len()
// 		}
// 	}
// 	return 0
// }

// func Slice(dataSlicePtr interface{}, i, j int) {
// 	if dataSlicePtr != nil {
// 		dt := reflect.TypeOf(dataSlicePtr)
// 		if dt.Kind() == reflect.Ptr && dt.Elem().Kind() == reflect.Slice {
// 			indirect := reflect.Indirect(reflect.ValueOf(dataSlicePtr))
// 			if j > indirect.Len() {
// 				j = indirect.Len()
// 			}

// 			if i < 0 || j < 0 || i > j {
// 				return
// 			}

//				indirect.Set(indirect.Slice(i, j))
//			}
//		}
//	}
// func GetMultiClusterResources[T any](

// )

func CutPagingSliceResourceList[T any](eachClusterResources []iapiserver.EachResourceRangeListState[T], list *[]T,
	pageNum, pageSize int,
	lessFunc func(i, j int) bool,
) int {
	for _, ret := range eachClusterResources {
		*list = append(*list, ret.List...)
	}

	total := len(*list)
	if total > 0 {
		sort.SliceStable(*list, lessFunc)
		s, index := paging.Index(total, pageNum, pageSize)
		plist := *list
		plist = plist[s:index]
		*list = plist
	}
	return total
}

func CutPagingSliceResourceList2[T any](eachClusterResources []iapiserver.EachResourceRangeListState[T], list *[]T,
	pageNum, pageSize int,
	sortBy string, asc bool,
) int {
	for _, ret := range eachClusterResources {
		*list = append(*list, ret.List...)
	}

	total := len(*list)
	if total > 0 {
		sortutil.StructSliceSort(*list, sortBy, asc)
		//sort.SliceStable(*list, lessFunc)
		s, index := paging.Index(total, pageNum, pageSize)
		plist := *list
		plist = plist[s:index]
		*list = plist
	}
	return total
}

func sortWithCommonObjectParam(obj1, obj2 metav1.Object, sortType string, IsDesc bool) bool {
	switch sortType {
	case iapiserver.KubernetesResourceSortByName:
		if obj1.GetName() != obj2.GetName() { // if equal, compare create time instead
			return compareutil.Compare(obj1.GetName(), obj2.GetName(), IsDesc)
		}
	case iapiserver.KubernetesResourceSortByNamespace:
		if obj1.GetNamespace() != obj2.GetNamespace() { // if equal, compare create time instead
			return compareutil.Compare(obj1.GetNamespace(), obj2.GetNamespace(), IsDesc)
		}
	}
	//fall back to create time if sort field equal or sort type not supported
	if obj1.GetCreationTimestamp().Unix() == obj2.GetCreationTimestamp().Unix() {
		if obj1.GetName() == obj2.GetName() {
			return compareutil.Compare(obj1.GetNamespace(), obj2.GetNamespace(), IsDesc)
		}
		return compareutil.Compare(obj1.GetName(), obj2.GetName(), IsDesc)
	}
	return compareutil.Compare(obj1.GetCreationTimestamp().Unix(), obj2.GetCreationTimestamp().Unix(), IsDesc)
}
