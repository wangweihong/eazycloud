package iharbor

import (
	"time"

	"github.com/wangweihong/gotoolbox/pkg/errors"
)

type Repository struct {
	Description    string     `json:"description"`
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	ProjectID      int        `json:"project_id"`
	ProjectName    string     `json:"project_name,omitempty"`
	RepositoryName string     `json:"repository_name,omitempty"`
	ArtifactCount  int        `json:"artifact_count"`
	PullCount      int        `json:"pull_count"`
	StarCount      int        `json:"star_count"`
	TagsCount      int        `json:"tags_count"`
	CreationTime   *time.Time `json:"creation_time,omitempty"`
	UpdateTime     *time.Time `json:"update_time,omitempty"`
	Size           int64      `json:"size"`
	Artifacts      *Artifact  `json:"artifacts"`
}

func (p *Repository) Convert() *Repository {

	if p.RepositoryName == "" {
		p.RepositoryName = p.Name
	}
	return p
}

type SearchRepository struct {
	ProjectID      int    `json:"project_id"`
	ProjectName    string `json:"project_name,omitempty"`
	ProjectPublic  bool   `json:"project_public"`
	RepositoryName string `json:"repository_name,omitempty"`
	PullCount      int    `json:"pull_count"`
	TagsCount      int    `json:"tags_count"`
}

type HarborSearchResult struct {
	Project    []Project          `json:"project"`
	Repository []SearchRepository `json:"repository"` //使用这个字段是因为search接口查询出来的repository的数据比较简单，而不是repository详情的数据
	//Repository []Repository `json:"repository"`
}

type RepositoryListRequest struct {
	ProjectIdentifier
	Q    string `json:"q" form:"q"`       //query 仓库名
	Sort string `json:"sort" form:"sort"` //Sort method, valid values include: 'name', '-name', 'creation_time', '-creation_time', 'update_time', '-update_time'. Here '-' stands for descending order.
	PagingParam
}

func (r *RepositoryListRequest) Validate() error {
	if r.ProjectName == "" {
		return errors.Errorf("project name is empty")
	}

	if r.Q != "" {
		r.Q = "name=~" + r.Q
	}
	r.PagingParam.Valiate()
	return nil
}

type RepositoryDeleteRequest struct {
	ProjectIdentifier
	RepositoryIdentifier //必须是完整的repo name
}

type RepositoryIdentifier struct {
	RepoName string `json:"repo_name" form:"repo_name"` //必须是完整的repo name。 project / repo.
}

type RepositoryDeleteResponse struct {
}

type RepositoryListResponse struct {
	List       []Repository `json:"list"`
	TotalCount int64        `json:"total_count"`
}

type RepositoryListSearchResponse struct {
	List       []SearchRepository `json:"list"`
	TotalCount int64              `json:"total_count"`
}

type RepositoryGetRequest struct {
	ProjectIdentifier
	RepositoryName string `json:"q" form:"q"`
}

type RepositoryGetResponse struct {
	Repository Repository `json:"repository"`
}
