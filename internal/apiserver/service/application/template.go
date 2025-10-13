package application

import (
	"context"
	"encoding/base64"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/apis/imachinery"

	"github.com/wangweihong/gotoolbox/pkg/errors"
	"helm.sh/helm/v3/pkg/repo"
)

func (s *applicationService) ApplicationTemplateValidate(ctx context.Context, req *iapiserver.ApplicationTemplateValidateRequest) error {
	_, err := renderPackageDataToManifests(&RenderParam{packageTarData: base64.StdEncoding.EncodeToString(req.PackageData)})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *applicationService) ApplicationTemplateList(ctx context.Context, req *iapiserver.ApplicationTemplateListRequest) (*iapiserver.ApplicationTemplateListResponse, error) {
	metas, total, err := s.store.ApplicationTemplates().List(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp := &iapiserver.ApplicationTemplateListResponse{}
	resp.Total = total
	resp.List = metas
	return resp, nil
}

// func (s *applicationService) ApplicationTemplateList(ctx context.Context, req *iapiserver.ApplicationTemplateListRequest) (*iapiserver.ApplicationTemplateListResponse, error) {
// 	resp := &iapiserver.ApplicationTemplateListResponse{}

// 	appStore, err := s.store.AppStores().Get(ctx, req.FilterStore)
// 	if err != nil {
// 		return nil, errors.WithStack(err)
// 	}
// 	if appStore.AppSource != iapiserver.AppSourceExternal {
// 		return nil, errors.Errorf("target app store is not external")
// 	}

// 	indexFile, err := loadRepoIndex(appStore)
// 	if err != nil {
// 		return nil, errors.WithStack(err)
// 	}

// 	resp.List = make([]*iapiserver.ApplicationTemplateEntry, 0, len(indexFile.Entries))
// 	for chartName, chartVersions := range indexFile.Entries {
// 		if req.Keyword != "" && !strings.Contains(chartName, req.Keyword) {
// 			continue
// 		}

// 		apiData := convertHelmChartToApplicationTemplateEntry(appStore, chartName, chartVersions)
// 		// for _, v := range instances {
// 		// 	if v.AppStoreUUID == apiData.AppStoreUUID && v.Template == apiData.Name {
// 		// 		apiData.InstanceNum += 1
// 		// 	}
// 		// }

// 		resp.List = append(resp.List, apiData)
// 	}

// 	metas, err := s.store.ApplicationTemplates().List(ctx)
// 	if err != nil {
// 		return nil, errors.WithStack(err)
// 	}
// 	var filterMeta []*iapiserver.ApplicationTemplateEntry
// 	// TODO: 返回实例数
// 	// 案分类过滤
// 	for _, meta := range metas {
// 		if req.Keyword != "" && sets.NewString(meta.FuzzyFields()...).Contain(req.Fuzzy) {
// 			continue
// 		}

// 		if !stringutil.MatchIfNotEmptry(meta.AppStoreID, req.FilterStore) {
// 			continue
// 		}

// 		filterMeta = append(filterMeta, &iapiserver.ApplicationTemplateEntry{ApplicationTemplate: *meta})
// 	}

// 	i, e := paging.Index(len(filterMeta), req.PageNum, req.PageSize)
// 	resp.Total = int64(len(filterMeta))
// 	resp.List = filterMeta[i:e]
// 	return resp, nil
// }

func (s *applicationService) ApplicationTemplateGet(ctx context.Context, req *iapiserver.ApplicationTemplateGetRequest) (*iapiserver.ApplicationTemplateGetResponse, error) {
	meta, err := s.store.ApplicationTemplates().Get(ctx, req.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &iapiserver.ApplicationTemplateGetResponse{ApplicationTemplate: *meta}, nil
}

func (s *applicationService) ApplicationTemplateAdd(ctx context.Context, req *iapiserver.ApplicationTemplateAddRequest) (*iapiserver.ApplicationTemplateAddResponse, error) {
	req.PackageData = []byte(base64.StdEncoding.EncodeToString(req.PackageData))
	renderData, err := renderPackageDataToManifests(&RenderParam{packageTarData: string(req.PackageData)})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.ApplicationTemplate.Versions = []iapiserver.ApplicationTemplateVersion{
		{
			ObjectMeta: imachinery.ObjectMeta{
				Name: renderData.Chart.Metadata.Version,
			},
			PackageData: string(req.PackageData),
			Metadata:    renderData.Chart.Metadata,
		},
	}

	if req.ApplicationTemplate.Categories == nil {
		catetory, err := s.store.ApplicationCategories().GetByName(ctx, req.AppStoreID, iapiserver.ApplicationCategoryUncategoriedName)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		req.ApplicationTemplate.Categories = append(req.ApplicationTemplate.Categories, *catetory)
	}

	meta, err := s.store.ApplicationTemplates().Add(ctx, &req.ApplicationTemplate)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &iapiserver.ApplicationTemplateAddResponse{ApplicationTemplate: *meta}, nil
}

func (s *applicationService) ApplicationTemplateUpdate(ctx context.Context, req *iapiserver.ApplicationTemplateUpdateRequest) error {
	// meta, err := s.store.ApplicationTemplates().Get(ctx, req.ID)
	// if err != nil {
	// 	return errors.WithStack(err)
	// }

	// if err := s.store.ApplicationTemplates().Update(ctx, &req.ApplicationTemplate); err != nil {
	// 	return errors.WithStack(err)
	// }

	// if req.URL != "" && req.URL != meta.URL {
	// 	repoIndex, err := fetchRemoteRepoIndexYaml(req.URL, req.AuthMode, req.AuthInfo, req.TlsConfig)
	// 	if err != nil {
	// 		return errors.WithStack(err)
	// 	}
	// 	async.GoRoutine(ctx, func(n context.Context) {
	// 		if err := writeRepoIndexFile(repoIndex, meta); err != nil {
	// 			log.Errorf("write repo index file %v fail:%v", generateLocalApplicationTemplateIndexPath(meta.ID, meta.Name))
	// 		}
	// 	})
	// }

	return nil
}

func (s *applicationService) ApplicationTemplateDelete(ctx context.Context, req *iapiserver.ApplicationTemplateDeleteRequest) error {
	meta, err := s.store.ApplicationTemplates().Get(ctx, req.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	store, err := s.store.AppStores().Get(ctx, meta.AppStoreID)
	if err != nil {
		return errors.WithStack(err)
	}

	if store.AppSource == iapiserver.AppSourceExternal {
		return errors.Errorf("external app store doesn't support delete template")
	}

	return nil
}

func convertHelmChartToApplicationTemplateEntry(appStore *iapiserver.AppStore, templateName string, meta []*repo.ChartVersion) *iapiserver.ApplicationTemplateEntry {
	//var versionInfo *iapiserver.ApplicationTemplateVersion

	// if meta is empty, it must be helm bug!
	// if len(meta) != 0 {
	// 	versionInfo = &iapiserver.ApplicationTemplateVersion{
	// 		Metadata:      meta[0].Metadata,
	// 		Name:          meta[0].Metadata.Version,
	// 		PackageData:   "",
	// 		Operator:      "",
	// 		StateOperator: "",
	// 		State:         "",
	// 		StateTime:     0,
	// 		UpdateLog:     "",
	// 	}
	// }

	info := &iapiserver.ApplicationTemplateEntry{
		ApplicationTemplate: iapiserver.ApplicationTemplate{
			ObjectMeta: imachinery.ObjectMeta{
				Name:        templateName,
				Description: "",
				CreatedAt:   imachinery.Time{},
				UpdatedAt:   imachinery.Time{},
			},

			AppStoreID: appStore.ID,
			Icon:       "",

			State:     "",
			StateTime: imachinery.Now(),
			//	CategorySet: nil,
			//		AppSource:   appStore.AppSource,
			Versions: nil,
		},
		//	CategorySet:   nil,
		//LatestVersion: versionInfo,
		//	VersionNum:    len(meta),
		//	InstanceNum:   0,
		//AppStoreName:  appStore.Name,
	}
	return info
}

func (s *applicationService) ApplicationTemplateVersionAdd(ctx context.Context, req *iapiserver.ApplicationTemplateVersionAddRequest) (*iapiserver.ApplicationTemplateVersionAddResponse, error) {
	return nil, nil
}

func (s *applicationService) ApplicationTemplateVersionDelete(ctx context.Context, req *iapiserver.ApplicationTemplateVersionDeleteRequest) error {
	return nil
}

func (s *applicationService) ApplicationTemplateVersionList(ctx context.Context, req *iapiserver.ApplicationTemplateVersionListRequest) (*iapiserver.ApplicationTemplateVersionListResponse, error) {
	return nil, nil
}

func (s *applicationService) ApplicationTemplateVersionGet(ctx context.Context, req *iapiserver.ApplicationTemplateVersionGetRequest) (*iapiserver.ApplicationTemplateVersionGetResponse, error) {
	return nil, nil
}

func (s *applicationService) ApplicationTemplateVersionUpdate(ctx context.Context, req *iapiserver.ApplicationTemplateVersionUpdateRequest) error {
	return nil
}
