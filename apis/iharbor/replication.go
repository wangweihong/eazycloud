package iharbor

import (
	"time"

	"github.com/wangweihong/gotoolbox/pkg/errors"
)

type AdapterNamesResponse struct {
	Names []string `json:"names"`
}

type AdapterInfoResponse struct {
	List map[string]Adapter `json:"list"`
}

type Adapter struct {
	EndpointPattern   AdapterEndpointPattern    `json:"endpoint_pattern"`
	CredentialPattern *AdapterCredentialPattern `json:"credential_pattern"`
}

type AdapterEndpointPattern struct {
	EndpointType      string                    `json:"endpoint_type"`
	Endpoints         []AdapterEndpoint         `json:"endpoints"`
	CredentialPattern *AdapterCredentialPattern `json:"credential_pattern"`
}
type AdapterCredentialPattern struct {
}

type AdapterEndpoint struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// copy from https://github.com/goharbor/harbor/blob/v1.10.3/src/replication/model/policy.go
// Policy defines the structure of a replication policy
type Policy struct {
	ID          int64  `json:"id" form:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Creator     string `json:"creator"` //创建时不传
	// source, 列表中src_regsitry.name == "Local" 则是push-based,否则是pull-based.
	SrcRegistry *Registry `json:"src_registry,omitempty"` //Pull-based, 想从远端仓库拉镜像时传。 SRc/Dest二选一, 只需要传registryID
	// destination
	DestRegistry *Registry `json:"dest_registry,omitempty"` //Push-based, 想推镜像到远端仓库时传。 SRc/Dest二选一，只需要传registryID
	// Only support two dest namespace modes:
	// Put all the src resources to the one single dest namespace
	// or keep namespaces same with the source ones (under this case,
	// the DestNamespace should be set to empty)
	DestNamespace string `json:"dest_namespace"`
	// Filters
	Filters []*Filter `json:"filters"` //TODO?
	// Trigger
	Trigger *Trigger `json:"trigger"` //触发复制的方式
	// Settings
	// TODO: rename the property name
	Deletion bool `json:"deletion"` //删除本地镜像时，同时删除远端镜像（事件触发时可选)
	// If override the image tag
	Override bool `json:"override"` //已存在时是否覆盖tag
	// Operations
	Enabled      bool      `json:"enabled"` //该策略是否启动
	CreationTime time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
}

// FilterType represents the type info of the filter.
type FilterType string

// Filter holds the info of the filter
type Filter struct {
	Type  FilterType  `json:"type"`
	Value interface{} `json:"value"` // FIXME
}

// TriggerType represents the type of trigger.
type TriggerType string

// Trigger holds info for a trigger
type Trigger struct {
	Type TriggerType `json:"type"`
	//manual/scheduled/event_based ，scheduled模式必须要传TriggerSettings的Cron.
	//事件驱动仅作用于push-based. 当指定仓库更新后，会触发事件驱动,事件驱动只能推送本次更新，而不是以前未推送版本
	Settings *TriggerSettings `json:"trigger_settings"`
}

// TriggerSettings is the setting about the trigger
type TriggerSettings struct {
	Cron string `json:"cron"` // cron tag字符串 每周："0 0 0 * * 0"（每周末午夜开始） 每天:"0 0 0 * * *"(每天午夜运行一次) 每小时："0 0 * * * *"
}

type PolicyListRequest struct {
	PagingParam
}

func (r *PolicyListRequest) Validate() error {
	r.PagingParam.Valiate()
	return nil
}

type PolicyListResponse struct {
	List       []Policy `json:"list"`
	TotalCount int      `json:"total_count"`
}
type PolicyExecuteIdentifier struct {
	PolicyID int64 `json:"policy_id"`
}

type PolicyExecuteListRequest struct {
	PolicyID int64  `json:"policy_id" form:"policy_id"`
	Status   string `json:"status" form:"status"`
	PagingParam
}


func (r *PolicyExecuteListRequest)Validate()error{
	if r.PolicyID == 0 {
		return errors.Errorf("invalid policy id")
	}
	r.PagingParam.Valiate()
	return nil
}


type PolicyExecuteListResponse struct {
	List       []PolicyExecuteResult `json:"list"`
	TotalCount int                   `json:"total_count"`
}

type PolicyExecuteResult struct {
	ID         int       `json:"id"`
	PolicyID   int       `json:"policy_id"`
	Status     string    `json:"status"`
	StatusText string    `json:"status_text"`
	Total      int       `json:"total"`
	Failed     int       `json:"failed"`
	Succeed    int       `json:"succeed"`
	InProgress int       `json:"in_progress"`
	Stopped    int       `json:"stopped"`
	Trigger    string    `json:"trigger"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
}
