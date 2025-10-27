package application

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/log"
	"helm.sh/helm/v3/pkg/kube"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	scheme2 "k8s.io/client-go/kubernetes/scheme"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/internal/apiserver/store"
	"github.com/wangweihong/eazycloud/internal/pkg/libkubernetes"

	"k8s.io/apimachinery/pkg/runtime"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/resource"

	diskcached "k8s.io/client-go/discovery/cached/disk"

	"k8s.io/client-go/restmapper"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/ghodss/yaml"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/releaseutil"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	"k8s.io/kubectl/pkg/scheme"
)

func getCapabilities(ctx context.Context, dc *discovery.DiscoveryClient) (*chartutil.Capabilities, error) {
	// force a discovery cache invalidation to always fetch the latest server version/capabilities.
	kubeVersion, err := dc.ServerVersion()
	if err != nil {
		return nil, errors.Wrap(err, "could not get server version from Kubernetes")
	}
	// Issue #6361:
	// Client-Go emits an error when an API service is registered but unimplemented.
	// We trap that error here and print a warning. But since the discovery client continues
	// building the API object, it is correctly populated with all valid APIs.
	// See https://github.com/kubernetes/kubernetes/issues/72051#issuecomment-521157642
	apiVersions, err := GetVersionSet(dc)
	if err != nil {
		if discovery.IsGroupDiscoveryFailedError(err) {
			log.Infof("WARNING: The Kubernetes server has an orphaned API service. Server reports: ", err)
			log.Infof("WARNING: To fix this, kubectl delete apiservice <service-name>")
		} else {
			return nil, errors.Wrap(err, "could not get apiVersions from Kubernetes")
		}
	}

	return &chartutil.Capabilities{
		APIVersions: apiVersions,
		KubeVersion: chartutil.KubeVersion{
			Version: kubeVersion.GitVersion,
			Major:   kubeVersion.Major,
			Minor:   kubeVersion.Minor,
		},
	}, nil
}

func GetVersionSet(client discovery.ServerResourcesInterface) ([]string, error) {
	groups, resources, err := client.ServerGroupsAndResources()
	if err != nil && !discovery.IsGroupDiscoveryFailedError(err) {
		return allKnownVersions(), errors.Errorf("could not get apiVersions from Kubernetes:%v", err)
	}

	if len(groups) == 0 && len(resources) == 0 {
		return allKnownVersions(), nil
	}

	versionMap := make(map[string]any)
	versions := []string{}
	// Extract the groups
	for _, g := range groups {
		for _, gv := range g.Versions {
			versionMap[gv.GroupVersion] = struct{}{}
		}
	}

	// Extract the resources
	var id string
	var ok bool
	for _, r := range resources {
		for _, rl := range r.APIResources {
			// A Kind at a GroupVersion can show up more than once. We only want
			// it displayed once in the final output.
			id = path.Join(r.GroupVersion, rl.Kind)
			if _, ok = versionMap[id]; !ok {
				versionMap[id] = struct{}{}
			}
		}
	}

	// Convert to a form that NewVersionSet can use
	for k := range versionMap {
		versions = append(versions, k)
	}

	return versions, nil
}
func allKnownVersions() []string {
	// We should register the built in extension APIs as well so CRDs are
	// supported in the default version set. This has caused problems with `helm
	// template` in the past, so let's be safe
	apiextensionsv1beta1.AddToScheme(scheme.Scheme)
	apiextensionsv1.AddToScheme(scheme.Scheme)

	groups := scheme.Scheme.PrioritizedVersionsAllGroups()
	vs := make([]string, 0, len(groups))
	for _, gv := range groups {
		vs = append(vs, gv.String())
	}
	return vs
}

type Capabilities struct {
	// KubeVersion is the Kubernetes version.
	KubeVersion KubeVersion
	// APIversions are supported Kubernetes API versions.
	APIVersions []string
}

// KubeVersion is the Kubernetes version.
type KubeVersion struct {
	Version string // Kubernetes version
	Major   string // Kubernetes major version
	Minor   string // Kubernetes minor version
}

const notesFileSuffix = "NOTES.txt"

func renderResource(ch *chart.Chart, values chartutil.Values, releaseName string, includeCrds bool) ([]*release.Hook, *bytes.Buffer, string, error) {
	hs := []*release.Hook{}
	b := bytes.NewBuffer(nil)
	files, err := engine.Render(ch, values)
	if err != nil {
		return hs, b, "", err
	}
	// NOTES.txt gets rendered like all the other files, but because it's not a hook nor a resource,
	// pull it out of here into a separate file so that we can actually use the output of the rendered
	// text file. We have to spin through this map because the file contains path information, so we
	// look for terminating NOTES.txt. We also remove it from the files so that we don't have to skip
	// it in the sortHooks.
	var notesBuffer bytes.Buffer
	for k, v := range files {
		if strings.HasSuffix(k, notesFileSuffix) {
			if k == path.Join(ch.Name(), "templates", notesFileSuffix) {
				// If buffer contains data, add newline before adding more
				if notesBuffer.Len() > 0 {
					notesBuffer.WriteString("\n")
				}
				notesBuffer.WriteString(v)
			}
			delete(files, k)
		}
	}
	notes := notesBuffer.String()

	// Sort hooks, manifests, and partials. Only hooks and manifests are returned,
	// as partials are not used after renderer.Render. Empty manifests are also
	// removed here.
	hs, manifests, err := releaseutil.SortManifests(files, nil, releaseutil.InstallOrder)
	if err != nil {
		// By catching parse errors here, we can prevent bogus releases from going
		// to Kubernetes.
		//
		// We return the files as a big blob of data to help the user debug parser
		// errors.
		for name, content := range files {
			if strings.TrimSpace(content) == "" {
				continue
			}
			fmt.Fprintf(b, "---\n# Source: %s\n%s\n", name, content)
		}
		return hs, b, "", err
	}

	if includeCrds {
		for _, crd := range ch.CRDObjects() {
			fmt.Fprintf(b, "---\n# Source: %s\n%s\n", crd.Name, string(crd.File.Data[:]))
		}
	}
	for _, m := range manifests {
		fmt.Fprintf(b, "---\n# Source: %s\n%s\n", m.Name, m.Content)
	}

	return hs, b, notes, nil
}

var drivePathPattern = regexp.MustCompile(`^[a-zA-Z]:/`)
var utf8bom = []byte{0xEF, 0xBB, 0xBF}

func loadArchiveFilesHelm(data []byte) ([]*loader.BufferedFile, error) {
	files := []*loader.BufferedFile{}
	tr := tar.NewReader(bytes.NewReader(data))
	for {
		b := bytes.NewBuffer(nil)
		hd, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if hd.FileInfo().IsDir() {
			// Use this instead of hd.Typeflag because we don't have to do any
			// inference chasing.
			continue
		}

		switch hd.Typeflag {
		// We don't want to process these extension header files.
		case tar.TypeXGlobalHeader, tar.TypeXHeader:
			continue
		}

		// Archive could contain \ if generated on Windows
		delimiter := "/"
		if strings.ContainsRune(hd.Name, '\\') {
			delimiter = "\\"
		}

		parts := strings.Split(hd.Name, delimiter)
		n := strings.Join(parts[1:], delimiter)

		// Normalize the path to the / delimiter
		n = strings.Replace(n, delimiter, "/", -1)

		if path.IsAbs(n) {
			return nil, errors.New("chart illegally contains absolute paths")
		}

		n = path.Clean(n)
		if n == "." {
			// In this case, the original path was relative when it should have been absolute.
			return nil, errors.Errorf("chart illegally contains content outside the base directory: %q", hd.Name)
		}
		if strings.HasPrefix(n, "..") {
			return nil, errors.New("chart illegally references parent directory")
		}

		// In some particularly arcane acts of path creativity, it is possible to intermix
		// UNIX and Windows style paths in such a way that you produce a result of the form
		// c:/foo even after all the built-in absolute path checks. So we explicitly check
		// for this condition.
		if drivePathPattern.MatchString(n) {
			return nil, errors.New("chart contains illegally named files")
		}

		if parts[0] == "Chart.yaml" {
			return nil, errors.New("chart yaml not in base directory")
		}

		if _, err := io.Copy(b, tr); err != nil {
			return nil, err
		}

		data := bytes.TrimPrefix(b.Bytes(), utf8bom)

		files = append(files, &loader.BufferedFile{Name: n, Data: data})
		b.Reset()
	}

	if len(files) == 0 {
		return nil, errors.New("no files in chart archive")
	}
	return files, nil
}

func loadArchiveFilesHelmV2(data []byte) ([]*loader.BufferedFile, error) {
	return loader.LoadArchiveFiles(bytes.NewReader(data))
}

var _ genericclioptions.RESTClientGetter = &TopCmdConfig{}
var _ clientcmd.ClientConfig = &TopCmdConfig{}
var _ clientcmd.ConfigAccess = &TopCmdConfig{}

func NewTopCmdConfig(cluster *iapiserver.Cluster) (*TopCmdConfig, error) {
	return &TopCmdConfig{config: cluster.Config.ToRestConfig()}, nil
}

type TopCmdConfig struct {
	config *rest.Config
}

// ToRESTConfig returns restconfig
func (f *TopCmdConfig) ToRESTConfig() (*rest.Config, error) {
	return f.ToRawKubeConfigLoader().ClientConfig()
}

// ToDiscoveryClient returns discovery client
func (f *TopCmdConfig) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	config, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	// The more groups you have, the more discovery requests you need to make.
	// given 25 groups (our groups + a few custom resources) with one-ish version each, discovery needs to make 50 requests
	// double it just so we don't end up here again for a while.  This config is only used for discovery.
	config.Burst = 100

	// retrieve a user-provided value for the "cache-dir"
	// defaulting to ~/.kube/http-cache if no user-value is given.
	httpCacheDir := "/tmp/.kube/http-cache"
	//if f.CacheDir != "" {
	//	httpCacheDir = f.CacheDir
	//}

	discoveryCacheDir := computeDiscoverCacheDir(filepath.Join("/tmp", ".kube", "cache", "discovery"), config.Host)
	return diskcached.NewCachedDiscoveryClientForConfig(config, discoveryCacheDir, httpCacheDir, time.Duration(10*time.Minute))
}

// overlyCautiousIllegalFileCharacters matches characters that *might* not be supported.  Windows is really restrictive, so this is really restrictive
var overlyCautiousIllegalFileCharacters = regexp.MustCompile(`[^(\w/\.)]`)

// computeDiscoverCacheDir takes the parentDir and the host and comes up with a "usually non-colliding" name.
func computeDiscoverCacheDir(parentDir, host string) string {
	// strip the optional scheme from host if its there:
	schemelessHost := strings.Replace(strings.Replace(host, "https://", "", 1), "http://", "", 1)
	// now do a simple collapse of non-AZ09 characters.  Collisions are possible but unlikely.  Even if we do collide the problem is short lived
	safeHost := overlyCautiousIllegalFileCharacters.ReplaceAllString(schemelessHost, "_")
	return filepath.Join(parentDir, safeHost)
}

// ToRESTMapper returns a restmapper
func (f *TopCmdConfig) ToRESTMapper() (meta.RESTMapper, error) {
	discoveryClient, err := f.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}
	warningHandler := func(string) {}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	expander := restmapper.NewShortcutExpander(mapper, discoveryClient, warningHandler)
	return expander, nil
}

// ToRawKubeConfigLoader return kubeconfig loader as-is
func (f *TopCmdConfig) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return f
}

// RawConfig returns the merged result of all overrides
func (f *TopCmdConfig) RawConfig() (clientcmdapi.Config, error) {
	return *generateKubeClusterApiConfigFromRestConfig(f.config), nil
}

// ClientConfig returns a complete client config
func (f *TopCmdConfig) ClientConfig() (*rest.Config, error) {
	return f.config, nil
}

// Namespace returns the namespace resulting from the merged
// result of all overrides and a boolean indicating if it was
// overridden
func (f *TopCmdConfig) Namespace() (string, bool, error) {
	return "default", false, nil
}

// ConfigAccess returns the rules for loading/persisting the config.
func (f *TopCmdConfig) ConfigAccess() clientcmd.ConfigAccess {
	return f
}

func (g *TopCmdConfig) Load() (*clientcmdapi.Config, error) {
	return generateKubeClusterApiConfigFromRestConfig(g.config), nil
}

func (g *TopCmdConfig) GetLoadingPrecedence() []string {
	return nil
}

func (g *TopCmdConfig) GetStartingConfig() (*clientcmdapi.Config, error) {
	return generateKubeClusterApiConfigFromRestConfig(g.config), nil
}

func (g *TopCmdConfig) GetDefaultFilename() string {
	return ""
}

func (g *TopCmdConfig) IsExplicitFile() bool {
	return false
}

func (g *TopCmdConfig) GetExplicitFile() string {
	return ""
}

func (g *TopCmdConfig) IsDefaultConfig(config *rest.Config) bool {
	return false
}

func generateKubeClusterApiConfigFromRestConfig(config *rest.Config) *clientcmdapi.Config {
	clusterName := "kubernetes"
	userName := "kubernetes-admin"

	contextNmae := fmt.Sprintf("%s@%s", userName, clusterName)
	apiConfig := &clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			clusterName: {
				Server:                   config.Host,
				CertificateAuthorityData: config.CAData,
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			contextNmae: {
				Cluster:  clusterName,
				AuthInfo: userName,
			},
		},
		AuthInfos:      map[string]*clientcmdapi.AuthInfo{},
		CurrentContext: contextNmae,
	}
	apiConfig.AuthInfos[userName] = &clientcmdapi.AuthInfo{
		ClientKeyData:         config.KeyData,
		ClientCertificateData: config.CertData,
	}
	return apiConfig
}

func createResource(info *resource.Info) error {
	obj, err := resource.NewHelper(info.Client, info.Mapping).Create(info.Namespace, true, info.Object)
	if err != nil {
		return err
	}
	return info.Refresh(obj, true)
}

func deleteResource(info *resource.Info) error {
	policy := metav1.DeletePropagationBackground
	opts := &metav1.DeleteOptions{PropagationPolicy: &policy}
	_, err := resource.NewHelper(info.Client, info.Mapping).DeleteWithOptions(info.Namespace, info.Name, opts)
	return err
}

const (
	appManagedByLabel              = "app.kubernetes.io/managed-by"
	appManagedByHelm               = "Helm"
	helmReleaseNameAnnotation      = "meta.helm.sh/release-name"
	helmReleaseNamespaceAnnotation = "meta.helm.sh/release-namespace"
)

// setMetadataVisitor adds release tracking metadata to all resources. If force is enabled, existing
// ownership metadata will be overwritten. Otherwise an error will be returned if any resource has an
// existing and conflicting value for the managed by label or Helm release/namespace annotations.
func setMetadataVisitor(releaseName, releaseNamespace string, force bool) resource.VisitorFunc {
	return func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}

		//if !force {
		//	if err := checkOwnership(info.Object, releaseName, releaseNamespace); err != nil {
		//		return fmt.Errorf("%s cannot be owned: %s", resourceString(info), err)
		//	}
		//}

		if err := mergeLabels(info.Object, map[string]string{
			appManagedByLabel: appManagedByHelm,
		}); err != nil {
			return errors.Errorf("%s labels could not be updated: %s", resourceString(info), err)
		}

		if err := mergeAnnotations(info.Object, map[string]string{
			helmReleaseNameAnnotation:      releaseName,
			helmReleaseNamespaceAnnotation: releaseNamespace,
		}); err != nil {
			return errors.Errorf("%s annotations could not be updated: %s", resourceString(info), err)
		}

		return nil
	}
}

var accessor = meta.NewAccessor()

func resourceString(info *resource.Info) string {
	_, k := info.Mapping.GroupVersionKind.ToAPIVersionAndKind()
	return fmt.Sprintf("%s %q in namespace %q", k, info.Name, info.Namespace)
}

func mergeLabels(obj runtime.Object, labels map[string]string) error {
	current, err := accessor.Labels(obj)
	if err != nil {
		return err
	}
	return accessor.SetLabels(obj, mergeStrStrMaps(current, labels))
}

func mergeAnnotations(obj runtime.Object, annotations map[string]string) error {
	current, err := accessor.Annotations(obj)
	if err != nil {
		return err
	}
	return accessor.SetAnnotations(obj, mergeStrStrMaps(current, annotations))
}

// merge two maps, always taking the value on the right
func mergeStrStrMaps(current, desired map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range current {
		result[k] = v
	}
	for k, desiredVal := range desired {
		result[k] = desiredVal
	}
	return result
}

type RenderParam struct {
	packageTarData string //压缩包
	caps           *chartutil.Capabilities
	options        *chartutil.ReleaseOptions
	vals           string //渲染更改项
}

type RenderData struct {
	Files         []*loader.BufferedFile `json:"files,omitempty"`
	Manifests     string                 `json:"manifests,omitempty"`
	UnpackageData []byte                 `json:"package_data,omitempty"`
	Chart         *chart.Chart           `json:"chart_metadata"`
	ValToRender   chartutil.Values       `json:"val_to_render"`
	Hooks         []*release.Hook        `json:"hooks"`
}

func objectKey(r *resource.Info) string {
	gvk := r.Object.GetObjectKind().GroupVersionKind()
	return fmt.Sprintf("%s/%s/%s/%s", gvk.GroupVersion().String(), gvk.Kind, r.Namespace, r.Name)
}

func existingResourceConflict(resources kube.ResourceList, releaseName, releaseNamespace string) (kube.ResourceList, error) {
	var requireUpdate kube.ResourceList

	err := resources.Visit(func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}

		helper := resource.NewHelper(info.Client, info.Mapping)
		existing, err := helper.Get(info.Namespace, info.Name)
		//existing, err := helper.Get(info.Namespace, info.Name, info.Export)
		if err != nil {
			if apierrors.IsNotFound(err) {
				return nil
			}
			return errors.Wrap(err, "could not get information about the resource")
		}

		// Allow adoption of the resource if it is managed by Helm and is annotated with correct release name and namespace.
		if err := checkOwnership(existing, releaseName, releaseNamespace); err != nil {
			return errors.Errorf("%s exists and cannot be imported into the current release: %s", resourceString(info), err)
		}

		requireUpdate.Append(info)
		return nil
	})

	return requireUpdate, err
}

func checkOwnership(obj runtime.Object, releaseName, releaseNamespace string) error {
	lbls, err := accessor.Labels(obj)
	if err != nil {
		return err
	}
	annos, err := accessor.Annotations(obj)
	if err != nil {
		return err
	}

	var errs []error
	if err := requireValue(lbls, appManagedByLabel, appManagedByHelm); err != nil {
		errs = append(errs, fmt.Errorf("label validation error: %s", err))
	}
	if err := requireValue(annos, helmReleaseNameAnnotation, releaseName); err != nil {
		errs = append(errs, fmt.Errorf("annotation validation error: %s", err))
	}
	if err := requireValue(annos, helmReleaseNamespaceAnnotation, releaseNamespace); err != nil {
		errs = append(errs, fmt.Errorf("annotation validation error: %s", err))
	}

	if len(errs) > 0 {
		err := errors.New("invalid ownership metadata")
		for _, e := range errs {
			err = errors.Errorf("%w; %s", err, e)
		}
		return err
	}

	return nil
}

func requireValue(meta map[string]string, k, v string) error {
	actual, ok := meta[k]
	if !ok {
		return fmt.Errorf("missing key %q: must be set to %q", k, v)
	}
	if actual != v {
		return fmt.Errorf("key %q must equal %q: current value is %q", k, v, actual)
	}
	return nil
}

// Result contains the information of created, updated, and deleted resources
// for various kube API calls along with helper methods for using those
// resources
type Result struct {
	Created kube.ResourceList
	Updated kube.ResourceList
	Deleted kube.ResourceList
}

var metadataAccessor = meta.NewAccessor()

// ResourcePolicyAnno is the annotation name for a resource policy
const ResourcePolicyAnno = "helm.sh/resource-policy"

// KeepPolicy is the resource policy type for keep
//
// This resource policy type allows resources to skip being deleted
//
//	during an uninstallRelease action.
const KeepPolicy = "keep"

func updateResources(original, target kube.ResourceList, force bool) (*Result, error) {
	updateErrors := []string{}
	res := &Result{}

	err := target.Visit(func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}
		kind := info.Mapping.GroupVersionKind.Kind
		helper := resource.NewHelper(info.Client, info.Mapping)
		if _, err := helper.Get(info.Namespace, info.Name); err != nil {
			if !apierrors.IsNotFound(err) {
				log.Infof("could not get information about the resource %s called %q in %s", kind, info.Name, info.Namespace)
				return errors.Wrap(err, "could not get information about the resource")
			}

			// Append the created resource to the results, even if something fails
			res.Created = append(res.Created, info)

			// Since the resource does not exist, create it.
			if err := createResource(info); err != nil {
				log.Errorf("failed to create resource %s called %q in %s", kind, info.Name, info.Namespace)
				return errors.Wrap(err, "failed to create resource")
			}

			log.Infof("Created a new %s called %q in %s success", kind, info.Name, info.Namespace)
			return nil
		}

		originalInfo := original.Get(info)
		if originalInfo == nil {
			log.Infof("no %s with the name %q found", kind, info.Name)
			return errors.Errorf("no %s with the name %q found", kind, info.Name)
		}

		if err := updateResource(info, originalInfo.Object, force); err != nil {
			updateErrors = append(updateErrors, err.Error())
		}
		// Because we check for errors later, append the info regardless
		res.Updated = append(res.Updated, info)

		return nil
	})

	switch {
	case err != nil:
		return res, err
	case len(updateErrors) != 0:
		return res, errors.Errorf(strings.Join(updateErrors, " && "))
	}

	for _, info := range original.Difference(target) {
		log.Infof("Deleting %q in %s...", info.Name, info.Namespace)

		if err := info.Get(); err != nil {
			log.Infof("Unable to get obj %q, err: %s", info.Name, err.Error())
			continue
		}
		annotations, err := metadataAccessor.Annotations(info.Object)
		if err != nil {
			log.Infof("Unable to get annotations on %q, err: %s", info.Name, err.Error())
		}
		if annotations != nil && annotations[ResourcePolicyAnno] == KeepPolicy {
			log.Infof("Skipping delete of %q due to annotation [%s=%s]", info.Name, ResourcePolicyAnno, KeepPolicy)
			continue
		}
		if err := deleteResource(info); err != nil {
			log.Infof("Failed to delete %q, err: %s", info.ObjectName(), err.Error())
			continue
		}
		res.Deleted = append(res.Deleted, info)
	}
	return res, nil
}

func updateResource(target *resource.Info, currentObj runtime.Object, force bool) error {
	var (
		obj    runtime.Object
		helper = resource.NewHelper(target.Client, target.Mapping)
		kind   = target.Mapping.GroupVersionKind.Kind
	)

	// if --force is applied, attempt to replace the existing resource with the new object.
	if force {
		var err error
		obj, err = helper.Replace(target.Namespace, target.Name, true, target.Object)
		if err != nil {
			log.Infof("failed to replace resource %s called %q in %s", kind, target.Name, target.Namespace)
			return errors.Wrap(err, "failed to replace resource")
		}
		log.Infof("Replaced %q with kind %s for kind %s", target.Name, currentObj.GetObjectKind().GroupVersionKind().Kind, kind)
	} else {
		patch, patchType, err := createPatch(target, currentObj)
		if err != nil {
			return errors.Wrap(err, "failed to create patch")
		}

		if patch == nil || string(patch) == "{}" {
			log.Infof("Looks like there are no changes for %s %q", target.Mapping.GroupVersionKind.Kind, target.Name)
			// This needs to happen to make sure that Helm has the latest info from the API
			// Otherwise there will be no labels and other functions that use labels will panic
			if err := target.Get(); err != nil {
				log.Infof("failed to refresh resource information")
				return errors.Wrap(err, "failed to refresh resource information")
			}
			return nil
		}
		// send patch to server
		obj, err = helper.Patch(target.Namespace, target.Name, patchType, patch, nil)
		if err != nil {
			log.Infof("cannot patch %q with kind %s", target.Name, kind)
			return errors.Wrapf(err, "cannot patch %q with kind %s", target.Name, kind)
		}
	}

	target.Refresh(obj, true)
	return nil
}

func createPatch(target *resource.Info, current runtime.Object) ([]byte, types.PatchType, error) {
	oldData, err := json.Marshal(current)
	if err != nil {
		return nil, types.StrategicMergePatchType, errors.Wrap(err, "serializing current configuration")
	}
	newData, err := json.Marshal(target.Object)
	if err != nil {
		return nil, types.StrategicMergePatchType, errors.Wrap(err, "serializing target configuration")
	}

	// Fetch the current object for the three way merge
	helper := resource.NewHelper(target.Client, target.Mapping)
	currentObj, err := helper.Get(target.Namespace, target.Name)
	if err != nil && !apierrors.IsNotFound(err) {
		return nil, types.StrategicMergePatchType, errors.Wrapf(err, "unable to get data for current object %s/%s", target.Namespace, target.Name)
	}

	// Even if currentObj is nil (because it was not found), it will marshal just fine
	currentData, err := json.Marshal(currentObj)
	if err != nil {
		return nil, types.StrategicMergePatchType, errors.Wrap(err, "serializing live configuration")
	}

	// Get a versioned object
	versionedObject := AsVersioned(target)

	// Unstructured objects, such as CRDs, may not have an not registered error
	// returned from ConvertToVersion. Anything that's unstructured should
	// use the jsonpatch.CreateMergePatch. Strategic Merge Patch is not supported
	// on objects like CRDs.
	_, isUnstructured := versionedObject.(runtime.Unstructured)

	// On newer K8s versions, CRDs aren't unstructured but has this dedicated type
	_, isCRD := versionedObject.(*apiextensionsv1beta1.CustomResourceDefinition)

	if isUnstructured || isCRD {
		// fall back to generic JSON merge patch
		patch, err := jsonpatch.CreateMergePatch(oldData, newData)
		return patch, types.MergePatchType, err
	}

	patchMeta, err := strategicpatch.NewPatchMetaFromStruct(versionedObject)
	if err != nil {
		return nil, types.StrategicMergePatchType, errors.Wrap(err, "unable to create patch metadata from object")
	}

	patch, err := strategicpatch.CreateThreeWayMergePatch(oldData, newData, currentData, patchMeta, true)
	return patch, types.StrategicMergePatchType, err
}

// AsVersioned converts the given info into a runtime.Object with the correct
// group and version set
func AsVersioned(info *resource.Info) runtime.Object {
	return convertWithMapper(info.Object, info.Mapping)
}

// convertWithMapper converts the given object with the optional provided
// RESTMapping. If no mapping is provided, the default schema versioner is used
func convertWithMapper(obj runtime.Object, mapping *meta.RESTMapping) runtime.Object {
	s := kubernetesNativeScheme()
	var gv = runtime.GroupVersioner(schema.GroupVersions(s.PrioritizedVersionsAllGroups()))
	if mapping != nil {
		gv = mapping.GroupVersionKind.GroupVersion()
	}
	if obj, err := runtime.ObjectConvertor(s).ConvertToVersion(obj, gv); err == nil {
		return obj
	}
	return obj
}

var k8sNativeScheme *runtime.Scheme
var k8sNativeSchemeOnce sync.Once

// kubernetesNativeScheme returns a clean *runtime.Scheme with _only_ Kubernetes
// native resources added to it. This is required to break free of custom resources
// that may have been added to scheme.Scheme due to Helm being used as a package in
// combination with e.g. a versioned kube client. If we would not do this, the client
// may attempt to perform e.g. a 3-way-merge strategy patch for custom resources.
func kubernetesNativeScheme() *runtime.Scheme {
	k8sNativeSchemeOnce.Do(func() {
		k8sNativeScheme = runtime.NewScheme()
		scheme2.AddToScheme(k8sNativeScheme)
		// API extensions are not in the above scheme set,
		// and must thus be added separately.
		apiextensionsv1beta1.AddToScheme(k8sNativeScheme)
		apiextensionsv1.AddToScheme(k8sNativeScheme)
	})
	return k8sNativeScheme
}

// checkDependencies checks the dependencies for a chart.
func checkDependencies(ch *chart.Chart, reqs []*chart.Dependency) error {
	var missing []string

OUTER:
	for _, r := range reqs {
		for _, d := range ch.Dependencies() {
			if d.Name() == r.Name {
				continue OUTER
			}
		}
		missing = append(missing, r.Name)
	}

	if len(missing) > 0 {
		return errors.Errorf("found in Chart.yaml, but missing in charts/ directory: %s", strings.Join(missing, ", "))
	}
	return nil
}

// 将chart包渲染成资源对象
func renderChart2Resource(ctx context.Context, s store.Factory, instance *iapiserver.ApplicationInstance, revision int, newPackageData string, newValue string) ([]*resource.Info, error) {
	revisionInfo, err := instance.GetHistory(revision)
	if err != nil {
		return nil, err
	}

	if newValue == "" {
		newValue = revisionInfo.Values
	}

	if newPackageData == "" {
		newPackageData = revisionInfo.PackageData
	}

	cluster, err := s.Clusters().Get(ctx, instance.ClusterID)
	if err != nil {
		return nil, err
	}

	cmdConfigGetter, err := NewTopCmdConfig(cluster)
	if err != nil {
		return nil, err
	}

	dc, err := libkubernetes.ToDiscoveryClient(cluster.Config)
	if err != nil {
		return nil, err
	}

	capability, err := getCapabilities(ctx, dc)
	if err != nil {
		return nil, err
	}

	var isUpgrade, isInstall bool
	if instance.Revision == 1 {
		isInstall = true
	} else {
		isUpgrade = true
	}

	renderParam := &RenderParam{
		packageTarData: newPackageData,
		vals:           newValue,
		caps:           capability,
		options: &chartutil.ReleaseOptions{
			Name:      instance.Name,
			Namespace: instance.Namespace,
			Revision:  instance.Revision,
			IsInstall: isInstall,
			IsUpgrade: isUpgrade,
		},
	}

	renderData, err := renderPackageDataToManifests(renderParam)
	if err != nil {
		return nil, err
	}
	chrt := renderData.Chart

	//FIXME：不支持Chart.yaml中配置Chart包本身不含有的Chart依赖
	if deps := chrt.Metadata.Dependencies; deps != nil {
		if err := checkDependencies(chrt, deps); err != nil {
			return nil, errors.Errorf("check chat dependencies fail:%v", err)
		}
	}

	//FIXME： install crd
	if crds := chrt.CRDObjects(); len(crds) > 0 {
	}

	factory := cmdutil.NewFactory(cmdConfigGetter)
	//scheme, err := factory.Validator(true)
	scheme, err := factory.Validator("true")
	if err != nil {
		return nil, errors.Errorf("render template %v error", err)
	}

	resources, err := factory.NewBuilder().
		ContinueOnError().
		NamespaceParam(instance.Namespace).
		DefaultNamespace().
		Flatten().
		Unstructured().
		Schema(scheme).
		Stream(bytes.NewReader([]byte(renderData.Manifests)), "").
		Do().
		Infos()
	if err != nil {
		const stopValidateMessage = "if you choose to ignore these errors, turn validation off with --validate=false"
		if strings.Contains(err.Error(), stopValidateMessage) {
			err = errors.New(strings.ReplaceAll(err.Error(), "; "+stopValidateMessage, ""))
		}
		return nil, errors.Errorf("render template %v error", err)
	}

	return resources, nil
}

func renderPackageDataToManifests(param *RenderParam) (*RenderData, error) {
	renderData := &RenderData{}
	//NOTE: 必须base64加密，不然数据无法解析
	unpackageData, err := base64.StdEncoding.DecodeString(string(param.packageTarData))
	if err != nil {
		return nil, errors.Errorf("parse package error:" + err.Error())
	}

	var vals map[string]any
	//github.com/ghodss/yaml
	if err := yaml.Unmarshal([]byte(param.vals), &vals); err != nil {
		return nil, errors.Errorf("render template %v error", err)
	}

	bufferFiles, err := loadArchiveFilesHelmV2([]byte(unpackageData))
	if err != nil {
		return nil, errors.Errorf("render template %v error", err)
	}

	chrt, err := loader.LoadFiles(bufferFiles)
	if err != nil {
		return nil, errors.Errorf("render template %v error", err)
	}

	if param.options == nil {
		param.options = &chartutil.ReleaseOptions{}
	}

	valuesToRender, err := chartutil.ToRenderValues(chrt, vals, *param.options, param.caps)
	if err != nil {
		return nil, errors.Errorf("render template %v error", err)
	}
	hooks, buffers, _, err := renderResource(chrt, valuesToRender, param.options.Name, true)
	if err != nil {
		return nil, errors.Errorf("render template %v error", err)
	}

	renderData.UnpackageData = unpackageData
	renderData.Files = bufferFiles
	renderData.Chart = chrt
	renderData.ValToRender = valuesToRender
	renderData.Manifests = buffers.String()
	renderData.Hooks = hooks

	return renderData, nil
}
