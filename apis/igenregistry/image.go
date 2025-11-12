package igenregistry

type RepoTags struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type Manifets struct {
	SchemaVersion int          `json:"schemaVersion"`
	MediaType     string       `json:"mediaType"`
	Config        ImageConfig  `json:"config"`
	Layers        []ImageLayer `json:"layers"`
}

type ImageLayer struct {
	MediaType string `json:"mediaType"`
	Size      int    `json:"size"`
	Digest    string `json:"digest"`
}

type ImageConfig struct {
	MediaType string `json:"mediaType"`
	Size      int    `json:"size"`
	Digest    string `json:"digest"`
}
