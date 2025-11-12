package igenregistry

import "github.com/wangweihong/gotoolbox/pkg/errors"

type (
	ListImagesResponse struct {
		Repositories []string `json:"repositories"`
	}
)

type (
	GetImageTagsRequest struct {
		Repo string `json:"rego" binding:"required"`
	}

	GetImageTagsResponse RepoTags
)

func (r *GetImageTagsRequest) Valdate() error {
	if r.Repo == "" {
		return errors.Errorf("repo is empty")
	}
	return nil
}

type (
	GetImageManifestsRequest struct {
		Repo string `json:"rego" binding:"required"`
		Tag  string `json:"tag" binding:"required"`
	}

	GetImageManifestsResponse Manifets
)

func (r *GetImageManifestsRequest) Valdate() error {
	if r.Repo == "" {
		return errors.Errorf("repo is empty")
	}
	if r.Tag == "" {
		return errors.Errorf("tag is empty")
	}
	return nil
}
