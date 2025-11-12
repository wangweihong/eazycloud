package iconfluence

import "time"

type ContentPage struct {
	ID      string          `json:"id" pg:",pk"`
	Type    string          `json:"type"`
	Status  string          `json:"status"`
	Title   string          `json:"title"`
	History *ContentHistory `json:"history"`
	// 当前页的祖先节点信息
	Ancestors []ContentAncestors `json:"ancestors"`
	Body      *ContentBody       `json:"body"`
	Meta      *ContentMeta       `pg:"rel:has-one"`
}

// ContentMeta record message
type ContentMeta struct {
	ID string `json:"id" pg:",pk"`
	// 访问次数
	VisitCount int `json:"visit_count"`
	// 是否共享
	Share bool `json:"share"`
}

type ContentAncestors struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Title  string `json:"title"`
}

type ContentLastUpdated struct {
	By struct {
		Type           string `json:"type"`
		Username       string `json:"username"`
		UserKey        string `json:"userKey"`
		ProfilePicture struct {
			Path      string `json:"path"`
			Width     int    `json:"width"`
			Height    int    `json:"height"`
			IsDefault bool   `json:"isDefault"`
		} `json:"profilePicture"`
		DisplayName string `json:"displayName"`
		Links       struct {
			Self string `json:"self"`
		} `json:"_links"`
		Expandable struct {
			Status string `json:"status"`
		} `json:"_expandable"`
	} `json:"by"`
	When      time.Time `json:"when"`
	Message   string    `json:"message"`
	Number    int       `json:"number"`
	MinorEdit bool      `json:"minorEdit"`
	Hidden    bool      `json:"hidden"`
	Links     struct {
		Self string `json:"self"`
	} `json:"_links"`
	Expandable struct {
		Content string `json:"content"`
	} `json:"_expandable"`
}

type ContentHistory struct {
	Latest      bool               `json:"latest"`
	LastUpdated ContentLastUpdated `json:"lastUpdated"`
	CreatedBy   struct {
		Type     string `json:"type"`
		Username string `json:"username"`
		UserKey  string `json:"userKey"`
		//ProfilePicture struct {
		//	Path      string `json:"path"`
		//	Width     int    `json:"width"`
		//	Height    int    `json:"height"`
		//	IsDefault bool   `json:"isDefault"`
		//} `json:"profilePicture"`
		DisplayName string `json:"displayName"`
	} `json:"createdBy"`
	CreatedDate time.Time `json:"createdDate"`
}

type ContentBody struct {
	View struct {
		Value          string `json:"value"`
		Representation string `json:"representation"`
	} `json:"view"`
}

type ContentDescendants struct {
	Attachment struct {
		Results []struct {
			ID       string `json:"id"`
			Type     string `json:"type"`
			Status   string `json:"status"`
			Title    string `json:"title"`
			Metadata struct {
				MediaType string `json:"mediaType"`
			} `json:"metadata"`
			Extensions struct {
				MediaType string `json:"mediaType"`
				FileSize  int    `json:"fileSize"`
				Comment   string `json:"comment"`
			} `json:"extensions"`
		} `json:"results"`
		Start int `json:"start"`
		Limit int `json:"limit"`
		Size  int `json:"size"`
	} `json:"attachment"`
}

type ContentAllPagesResponse struct {
	Results []*ContentPage `json:"results"`
	Start   int            `json:"start"`
	Limit   int            `json:"limit"`
	Size    int            `json:"size"`
}

type ContentTreeNode struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Title  string `json:"title"`

	Children []*ContentTreeNode `json:"children"`
	Share    bool               `json:"share"`
}
