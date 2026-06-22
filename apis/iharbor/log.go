package iharbor

import "time"

type Log struct {
	LogID      int       `json:"log_id"`
	Username   string    `json:"username"`
	ProjectID  int       `json:"project_id"`
	RepoName   string    `json:"repo_name"`
	RepoTag    string    `json:"repo_tag"`
	GUID       string    `json:"guid"`
	Operation  string    `json:"operation"`
	OpTime     time.Time `json:"op_time"`
	OpTimeUnix int64     `json:"op_time_unix"`

	//v2
	ID           int    `json:"id"`
	Resource     string `json:"resource"`
	ResourceType string `json:"resource_type"`
}

func (p *Log) Convert() *Log {
	if !p.OpTime.IsZero() {
		p.OpTimeUnix = p.OpTime.Unix()
	}
	return p
}

type LogListRequest struct {
	PagingParam
}

type LogListResponse struct {
	List       []Log `json:"list"`
	TotalCount int   `json:"total_count"`
}
