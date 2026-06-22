package iharbor

type SystemInfo struct {
	WithNotary                  bool   `json:"with_notary"`
	WithClair                   bool   `json:"with_clair"` //harbor是否部署clair(镜像扫描)
	WithAdmiral                 bool   `json:"with_admiral"`
	AdmiralEndpoint             string `json:"admiral_endpoint"`
	AuthMode                    string `json:"auth_mode"` //db_auth
	RegistryURL                 string `json:"registry_url"`
	ProjectCreationRestriction  string `json:"project_creation_restriction"`
	SelfRegistration            bool   `json:"self_registration"`
	HasCaRoot                   bool   `json:"has_ca_root"`
	HarborVersion               string `json:"harbor_version"` //版本
	NextScanAll                 int    `json:"next_scan_all"`
	RegistryStorageProviderName string `json:"registry_storage_provider_name"` //filesystem
	ReadOnly                    bool   `json:"read_only"`                      //harbor只读，用户在此模式下无法对镜像进行操作
	WithChartmuseum             bool   `json:"with_chartmuseum"`
}

type Volume struct {
	Storage struct {
		Total int64 `json:"total"` //系统容量,单位bytes
		Free  int64 `json:"free"`  //剩余容量,单位bytes
	} `json:"storage"`
}

// GC 任务
type Job struct {
	ID           int      `json:"id"`
	JobName      string   `json:"job_name"`
	JobKind      string   `json:"job_kind"`   //值:Generic（手动触发），Periodic（周期调度触发）
	Schedule     Schedule `json:"schedule"`   //周期调度
	JobStatus    string   `json:"job_status"` //值: finished
	Deleted      bool     `json:"deleted"`
	CreationTime string   `json:"creation_time"`
	UpdateTime   string   `json:"update_time"`
}

type Schedule struct {
	Type string `json:"type"` //值: Manual(立即执行，不需要带cron)，Weekly(每周)，Daily（每天），Hourly（每小时）Custom（自定义），None（无调度）
	Cron string `json:"cron"` // 每周："0 0 0 * * 0"（每周末午夜开始） 每天:"0 0 0 * * *"(每天午夜运行一次) 每小时："0 0 * * * *"
}

type ScheduleParameter struct {
	DeleteUntagged bool `json:"delete_untagged"`
	DryRun         bool `json:"dry_run"`
}

type ScheduleRequest struct {
	Schedule   *Schedule          `json:"schedule"`
	Parameters *ScheduleParameter `json:"parameters"`
}

type ScheduleGetResponse struct {
	Schedule *Schedule `json:"schedule"`
}

const (
	ScheduleTypeManual     = "Manual" //手动/立即触发
	ScheduleTypeCronWeekly = "Weekly"
	ScheduleTypeCronDaily  = "Daily"
	ScheduleTypeCronHourly = "Hourly"
	ScheduleTypeCronCustom = "Custom" //自定义
	ScheduleTypeCronNone   = "None"
)

type SystemGcIdentifier struct {
	ID int `json:"id" form:"id"`
}

type StatisticRequest struct {
}

type StatisticResponse struct {
	PrivateProjectCount int `json:"private_project_count"` //私有项目数
	PrivateRepoCount    int `json:"private_repo_count"`    //私有仓库数
	PublicProjectCount  int `json:"public_project_count"`  //公有项目数
	PublicRepoCount     int `json:"public_repo_count"`     //公有仓库数
	TotalProjectCount   int `json:"total_project_count"`   //总项目数
	TotalRepoCount      int `json:"total_repo_count"`      //总仓库数
}

type ScanAllMetrics struct {
	Total     int          `json:"total"`     //总扫描
	Completed int          `json:"completed"` //已完成
	Metrics   *ScanMetrics `json:"metrics,omitempty"`
	Requester string       `json:"requester"`
	Ongoing   bool         `json:"ongoing"` //正在扫描中
}

type ScanMetrics struct {
	Pending int `json:"Pending"`
	Running int `json:"Running"`
	Success int `json:"Success"` //已成功
}
