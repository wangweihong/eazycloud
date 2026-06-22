package application

import (
	"context"

	"github.com/wangweihong/gotoolbox/pkg/errors"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
)

func (s *applicationService) ApplicationInstanceList(ctx context.Context, req *iapiserver.ApplicationInstanceListRequest) (*iapiserver.ApplicationInstanceListResponse, error) {
	metas, total, err := s.store.ApplicationInstances().List(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	resp := &iapiserver.ApplicationInstanceListResponse{}
	resp.Total = total
	resp.List = metas
	return resp, nil
}

func (s *applicationService) ApplicationInstanceGet(ctx context.Context, req *iapiserver.ApplicationInstanceGetRequest) (*iapiserver.ApplicationInstanceGetResponse, error) {
	meta, err := s.store.ApplicationInstances().Get(ctx, req.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &iapiserver.ApplicationInstanceGetResponse{ApplicationInstance: *meta}, nil
}

func (s *applicationService) ApplicationInstanceWorkload(ctx context.Context, req *iapiserver.ApplicationInstanceWorkloadRequest) (*iapiserver.ApplicationInstanceWorkloadResponse, error) {
	return nil, errors.Errorf("not support yet")
}

func (s *applicationService) ApplicationInstanceAdd(ctx context.Context, req *iapiserver.ApplicationInstanceAddRequest) (*iapiserver.ApplicationInstanceAddResponse, error) {
	return nil, errors.Errorf("not support yet")
}

func (s *applicationService) ApplicationInstanceUpdate(ctx context.Context, req *iapiserver.ApplicationInstanceUpdateRequest) error {
	if err := s.store.ApplicationInstances().Update(ctx, &req.ApplicationInstance); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *applicationService) ApplicationInstanceDelete(ctx context.Context, req *iapiserver.ApplicationInstanceDeleteRequest) error {
	err := s.store.ApplicationInstances().Delete(ctx, req.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
