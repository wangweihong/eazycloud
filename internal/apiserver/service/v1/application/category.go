package application

import (
	"context"

	"github.com/wangweihong/gotoolbox/pkg/errors"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
)

func (s *applicationService) ApplicationCategoryList(ctx context.Context, req *iapiserver.ApplicationCategoryListRequest) (*iapiserver.ApplicationCategoryListResponse, error) {
	metas, total, err := s.store.ApplicationCategories().List(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp := &iapiserver.ApplicationCategoryListResponse{}
	resp.Total = total
	resp.List = metas
	return resp, nil
}

func (s *applicationService) ApplicationCategoryGet(ctx context.Context, req *iapiserver.ApplicationCategoryGetRequest) (*iapiserver.ApplicationCategoryGetResponse, error) {
	meta, err := s.store.ApplicationCategories().Get(ctx, req.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &iapiserver.ApplicationCategoryGetResponse{ApplicationCategory: *meta}, nil
}

func (s *applicationService) ApplicationCategoryAdd(ctx context.Context, req *iapiserver.ApplicationCategoryAddRequest) (*iapiserver.ApplicationCategoryAddResponse, error) {
	meta, err := s.store.ApplicationCategories().Add(ctx, &req.ApplicationCategory)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &iapiserver.ApplicationCategoryAddResponse{ApplicationCategory: *meta}, nil
}

func (s *applicationService) ApplicationCategoryUpdate(ctx context.Context, req *iapiserver.ApplicationCategoryUpdateRequest) error {
	if err := s.store.ApplicationCategories().Update(ctx, &req.ApplicationCategory); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *applicationService) ApplicationCategoryDelete(ctx context.Context, req *iapiserver.ApplicationCategoryDeleteRequest) error {
	err := s.store.ApplicationCategories().Delete(ctx, req.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
