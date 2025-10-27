package application

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/wangweihong/gotoolbox/pkg/async"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/httpcli"
	"github.com/wangweihong/gotoolbox/pkg/log"
	"github.com/wangweihong/gotoolbox/pkg/stringutil"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
)

func fetchRemoteRepoIndexYaml(repoUrl, authMode string, authInfo *iapiserver.AppStoreAuthInfo, AppStoreTlsConfig *iapiserver.ClientTlsConfig) (*repo.IndexFile, error) {
	repoUrl = stringutil.AddSuffixIfNotHas(repoUrl, "/index.yaml")

	builder := httpcli.NewHttpRequestBuilder().
		GET().
		WithEndpoint(repoUrl)

	if authMode == iapiserver.AppStoreRemoteAuthModeBASIC {
		builder.AddBasicAuthHeaderParam(authInfo.UserName, authInfo.Password)
	}
	tr := httpcli.NewTlsClientSkipVerifiedTransport()
	if AppStoreTlsConfig != nil {
		tr = httpcli.NewTlsClientCATransport(AppStoreTlsConfig.CaData)
	}

	resp, err := builder.Build().Invoke(
		httpcli.CallOptionTransport(tr),
	)
	if err != nil {
		return nil, errors.Wrap(err, "fetch remote repo index fail")
	}

	indexFile, err := parseIndexData([]byte(resp.GetBody()))
	if err != nil {
		return nil, err
	}

	return indexFile, nil
}

// parse chart index
func parseIndexData(data []byte) (*repo.IndexFile, error) {
	i := repo.NewIndexFile()
	//Use "github.com/ghodss/yaml"
	if err := yaml.Unmarshal(data, i); err != nil {
		if strings.Contains(err.Error(), "yaml: unmarshal errors:") {
			return i, errors.Errorf("yam unmarshal error")
		}
		return i, err
	}
	i.SortEntries()
	if i.APIVersion == "" {
		return i, errors.Errorf("no API version specified")
	}
	return i, nil
}

func generateLocalAppStoreIndexPath(uuid string, name string) string {
	return filepath.Join("/var/lib/eazycloud/appstore", uuid, "index.yaml")
}

func generateLocalAppStoreInfoPath(uuid string, name string) string {
	return filepath.Join("/var/lib/eazycloud/appstore", uuid, "store.yaml")
}

func generateLocalAppStoreChartVersionPath(storeUUID string, chartVersion *repo.ChartVersion) string {
	return filepath.Join("/var/lib/eazycloud/appstore", storeUUID, "charts", generateVersionTagFileName(chartVersion))
}

func generateVersionTagFileName(chartVersion *repo.ChartVersion) string {
	return chartVersion.Name + "-" + chartVersion.Version + ".tgz"
}

func writeRepoIndexFile(repoIndex *repo.IndexFile, store *iapiserver.AppStore) error {
	indexPath := generateLocalAppStoreIndexPath(store.ID, store.Name)
	if err := os.MkdirAll(filepath.Dir(indexPath), 0755); err != nil {
		return err
	}
	if err := repoIndex.WriteFile(indexPath, 0644); err != nil {
		return err
	}
	return nil
}

func writeStoreInfoFile(store *iapiserver.AppStore) error {
	indexPath := generateLocalAppStoreInfoPath(store.ID, store.Name)
	if err := os.MkdirAll(filepath.Dir(indexPath), 0755); err != nil {
		return err
	}
	b, err := yaml.Marshal(store)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(indexPath, b, 0755)
}

func loadRepoIndex(store *iapiserver.AppStore) (*repo.IndexFile, error) {
	data, err := ioutil.ReadFile(generateLocalAppStoreIndexPath(store.ID, store.Name))
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if data == nil {
		indexFile, err := fetchRemoteRepoIndexYaml(store.RemoteConfig.URL, store.RemoteConfig.AuthMode, store.RemoteConfig.AuthInfo, store.RemoteConfig.TlsConfig)
		if err != nil {
			return nil, err
		}

		async.GoRoutine(nil, func(n context.Context) {
			if err := writeRepoIndexFile(indexFile, store); err != nil {
				log.Errorf("write repo index file %v fail:%v", generateLocalAppStoreIndexPath(store.ID, store.Name))
			}
		})
		return indexFile, nil
	}
	return parseIndexData(data)
}
