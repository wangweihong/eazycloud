package iregistry

import "github.com/wangweihong/eazycloud/apis/imachinery"

type ListRegistryParam struct {
	imachinery.PagingParams
}

type ListRegistryResult struct {
}

type DeleteRegistryParam struct {
	ID string `json:"id" query:"id" binding:"required"`
}

type GetRegistryParam struct {
	ID string `json:"id" query:"id" binding:"required"`
}

type SearchRegistryParam struct {
	ID    string `json:"id"    query:"id"    binding:"required"`
	Query string `json:"query" query:"query" binding:"required"`
}
