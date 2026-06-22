package application

import (
	"context"

	"github.com/wangweihong/gotoolbox/pkg/async"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/log"
	"github.com/wangweihong/gotoolbox/pkg/paging"
	"github.com/wangweihong/gotoolbox/pkg/sets"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/apis/imachinery"
)

func (s *applicationService) AppStoreQuery(ctx context.Context, req *iapiserver.AppStoreListRequest) (*iapiserver.AppStoreListResponse, error) {
	metas, total, err := s.store.AppStores().Query(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp := &iapiserver.AppStoreListResponse{}
	resp.Total = total
	resp.List = metas
	return resp, nil
}

func (s *applicationService) AppStoreList(ctx context.Context, req *iapiserver.AppStoreListRequest) (*iapiserver.AppStoreListResponse, error) {
	metas, err := s.store.AppStores().List(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var filterMeta []*iapiserver.AppStore
	for _, meta := range metas {
		if req.AppSource != nil && *req.AppSource != meta.AppSource {
			continue
		}

		if req.Keyword != "" && sets.NewString(meta.State, meta.Name, meta.Description, meta.StateMessage).Contain(req.Keyword) {
			continue
		}

		filterMeta = append(filterMeta, meta)
	}

	resp := &iapiserver.AppStoreListResponse{}
	i, e := paging.Index(len(filterMeta), req.PageNum, req.PageSize)
	resp.Total = int64(len(filterMeta))
	resp.List = filterMeta[i:e]
	return resp, nil
}

func (s *applicationService) AppStoreGet(ctx context.Context, req *iapiserver.AppStoreGetRequest) (*iapiserver.AppStoreGetResponse, error) {
	meta, err := s.store.AppStores().Get(ctx, req.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &iapiserver.AppStoreGetResponse{AppStore: *meta}, nil
}

func (s *applicationService) AppStoreAdd(ctx context.Context, req *iapiserver.AppStoreAddRequest) (*iapiserver.AppStoreAddResponse, error) {
	var repoIndex *repo.IndexFile
	var err error

	switch req.AppSource {
	case iapiserver.AppSourceLocal:
		req.ApplicationCategories = []iapiserver.ApplicationCategory{
			{
				ObjectMeta: imachinery.ObjectMeta{
					Name: iapiserver.ApplicationCategoryUncategoriedName,
				},
			},
		}
	case iapiserver.AppSourceExternal:
		req.ApplicationCategories = nil
		repoIndex, err = fetchRemoteRepoIndexYaml(req.RemoteConfig.URL, req.RemoteConfig.AuthMode, req.RemoteConfig.AuthInfo, req.RemoteConfig.TlsConfig)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	default:
		return nil, errors.Errorf("invalid app store source")
	}

	req.ApplicationTemplates = nil
	meta, err := s.store.AppStores().Add(ctx, &req.AppStore)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if repoIndex != nil {
		async.GoRoutine(ctx, func(ctx context.Context) {
			if err := writeRepoIndexFile(repoIndex, meta); err != nil {
				log.Errorf("write repo index %v error:%v", generateLocalAppStoreIndexPath(meta.ID, meta.Name), err)
			}
			if err := writeStoreInfoFile(meta); err != nil {
				log.Errorf("write store index %v error:%v", generateLocalAppStoreInfoPath(meta.ID, meta.Name), err)
			}
		})
	}
	return &iapiserver.AppStoreAddResponse{AppStore: *meta}, nil
}

func (s *applicationService) AppStoreUpdate(ctx context.Context, req *iapiserver.AppStoreUpdateRequest) error {
	meta, err := s.store.AppStores().Get(ctx, req.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := s.store.AppStores().Update(ctx, &req.AppStore); err != nil {
		return errors.WithStack(err)
	}

	if req.AppSource == iapiserver.AppSourceExternal {
		if req.RemoteConfig.URL != "" && req.RemoteConfig.URL != meta.RemoteConfig.URL {
			repoIndex, err := fetchRemoteRepoIndexYaml(req.RemoteConfig.URL, req.RemoteConfig.AuthMode, req.RemoteConfig.AuthInfo, req.RemoteConfig.TlsConfig)
			if err != nil {
				return errors.WithStack(err)
			}
			async.GoRoutine(ctx, func(n context.Context) {
				if err := writeRepoIndexFile(repoIndex, meta); err != nil {
					log.Errorf("write repo index file %v fail:%v", generateLocalAppStoreIndexPath(meta.ID, meta.Name))
				}
			})
		}
	}

	return nil
}

func (s *applicationService) AppStoreSync(ctx context.Context, req *iapiserver.AppStoreSyncRequest) error {
	meta, err := s.store.AppStores().Get(ctx, req.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	if meta.AppSource == iapiserver.AppSourceLocal {
		return errors.Errorf("local app store doesn't support sync")
	}

	repoIndex, err := fetchRemoteRepoIndexYaml(meta.RemoteConfig.URL, meta.RemoteConfig.AuthMode, meta.RemoteConfig.AuthInfo, meta.RemoteConfig.TlsConfig)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := writeRepoIndexFile(repoIndex, meta); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *applicationService) AppStoreDelete(ctx context.Context, req *iapiserver.AppStoreDeleteRequest) error {
	err := s.store.AppStores().Delete(ctx, req.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *applicationService) AppStoreProbe(ctx context.Context, req *iapiserver.AppStoreProbeRequest) *iapiserver.AppStoreProbeResponse {
	resp := &iapiserver.AppStoreProbeResponse{Healthy: true}
	_, err := fetchRemoteRepoIndexYaml(req.RemoteConfig.URL, req.RemoteConfig.AuthMode, req.RemoteConfig.AuthInfo, req.RemoteConfig.TlsConfig)
	if err != nil {
		resp.Healthy = false
		resp.Error = err.Error()
	}

	return resp
}
