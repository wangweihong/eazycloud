package setting

import (
	"context"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

// 当前服务作为sp时，添加第三方idp应用作为单点登录提供商
// 在服务器使用第三方idp进行单点登陆前，需要先执行以下的步骤：
//  1. 通过idp的公共端点下载idp的元数据文件
//  2. 生成当前服务器的sp元数据文件，
//  3. 将sp元数据文件注册到idp服务

func (s *settingService) ServiceProviderList(ctx context.Context, req *iapiserver.ServiceProviderListRequest) (*iapiserver.ServiceProviderListResponse, error) {
	metas, total, err := s.store.ServiceProviders().List(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp := &iapiserver.ServiceProviderListResponse{}
	resp.Total = total
	resp.List = metas
	return resp, nil
}

func (s *settingService) ServiceProviderGet(ctx context.Context, req *iapiserver.ServiceProviderGetRequest) (*iapiserver.ServiceProviderGetResponse, error) {
	meta, err := s.store.ServiceProviders().Get(ctx, req.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &iapiserver.ServiceProviderGetResponse{ServiceProvider: *meta}, nil
}

func (s *settingService) ServiceProviderAdd(ctx context.Context, req *iapiserver.ServiceProviderAddRequest) (*iapiserver.ServiceProviderAddResponse, error) {
	meta, err := s.store.ServiceProviders().Add(ctx, &req.ServiceProvider)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &iapiserver.ServiceProviderAddResponse{ServiceProvider: *meta}, nil
}

func (s *settingService) ServiceProviderUpdate(ctx context.Context, req *iapiserver.ServiceProviderUpdateRequest) error {
	if err := s.store.ServiceProviders().Update(ctx, &req.ServiceProvider); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *settingService) ServiceProviderDelete(ctx context.Context, req *iapiserver.ServiceProviderDeleteRequest) error {
	err := s.store.ServiceProviders().Delete(ctx, req.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
