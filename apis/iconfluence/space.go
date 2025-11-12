package iconfluence

type Space struct {
	ID    interface{} `json:"id"`
	Key   string      `json:"key"`
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Links struct {
		Webui string `json:"webui"`
		Self  string `json:"self"`
	} `json:"_links"`
	Expandable struct {
		Metadata        string `json:"metadata"`
		Icon            string `json:"icon"`
		Description     string `json:"description"`
		RetentionPolicy string `json:"retentionPolicy"`
		Homepage        string `json:"homepage"`
	} `json:"_expandable"`
}
type SpaceListResponse struct {
	Results []Space `json:"results"`
	Start   int     `json:"start"`
	Limit   int     `json:"limit"`
	Size    int     `json:"size"`
	Links   struct {
		Self    string `json:"self"`
		Base    string `json:"base"`
		Context string `json:"context"`
	} `json:"_links"`
}
