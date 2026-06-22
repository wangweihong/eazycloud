package iapiserver

import (
	"fmt"
	"time"

	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/json"
	"gorm.io/gorm"
	"helm.sh/helm/v3/pkg/chart"

	"github.com/wangweihong/eazycloud/apis/ikubernetes"
	"github.com/wangweihong/eazycloud/apis/imachinery"
)

type Cluster struct {
	imachinery.ObjectMeta
	Type      string                     `json:"type"`
	IsStop    bool                       `json:"is_stop"`
	Config    *ikubernetes.ClusterConfig `json:"config"`
	State     string                     `json:"state"`
	IsManaged bool                       //纳管集群
	StateDesc string                     `json:"state_desc"`
	// KubeDeployConfig *KubeDeployConfig      `json:"DeployConfig,omitempty"`
	// KubeDeployStage  *DeployStageController `json:"DeployStage,omitempty"`
	// KubeVersion      *version.Info          `json:"kube_version"`
	IsHighAvailable  *bool
	ClusterNodeScale string `json:"cluster_node_scale"` //集群工作节点规模，部署前检测
	ServicePodCidr   string `json:"ServicePodCidr"`
	PodCidr          string `json:"PodCidr"`
}

func (s Cluster) FuzzyFields() []string {
	return []string{s.Name, s.Type, s.Description}
}

const (
	AppSourceLocal = iota
	AppSourceExternal
)

const (
	AppStoreRemoteAuthModeNone  = "NONE"
	AppStoreRemoteAuthModeBASIC = "BASIC"
	AppStoreRemoteAuthModeOAUTH = "OAUTH"
)

// +k8s:deepcopy-gen=true
type AppStoreRemoteConfig struct {
	IndexGeneratedAt imachinery.Time   `json:"index_generated_at"`
	URL              string            `json:"url"`
	AuthMode         string            `json:"auth_mode"`
	AuthInfo         *AppStoreAuthInfo `json:"auth_info"`
	TlsConfig        *ClientTlsConfig  `json:"auth_tls_config"`
}

// +k8s:deepcopy-gen=true
type AppStore struct {
	imachinery.ObjectMeta

	State        string                `json:"state"`
	StateMessage string                `json:"state_message"`
	AppSource    int                   `json:"app_source"    binding:"oneof=0 1"`
	RemoteConfig *AppStoreRemoteConfig `json:"remote_config" gorm:"-"`

	ApplicationTemplates  []ApplicationTemplate `json:"templates"`
	ApplicationCategories []ApplicationCategory `json:"categories"`
}

func (s *AppStore) Validate() error {
	if s.AppSource == AppSourceExternal {
		if s.RemoteConfig == nil {
			return errors.Errorf("external app store must set remote config")
		}
		//FIXME: url check
		if s.RemoteConfig.URL == "" {
			return errors.Errorf("external app store url is empty")
		}
	}
	return nil
}

func (s *AppStore) FuzzyFields() []string {
	return []string{s.Name, s.State, s.Description}
}

// BeforeCreate run before create database record.
func (s *AppStore) BeforeCreate(tx *gorm.DB) error {
	if err := s.ObjectMeta.BeforeCreate(tx); err != nil {
		return errors.Errorf("failed to run `BeforeCreate` hook: %v", err)
	}

	if s.Extend == nil {
		s.Extend = make(map[string]any)
	}

	s.State = HealthStatusHealthy
	s.StateMessage = ""
	if s.AppSource == AppSourceExternal {
		s.Extend["remote_config"] = s.RemoteConfig
	}
	s.ExtendShadow = json.ToString(s.Extend)

	return nil
}

// AfterCreate run after create database record.
func (s *AppStore) AfterCreate(tx *gorm.DB) error {
	return tx.Save(s).Error
}

// BeforeUpdate run before update database record.
func (s *AppStore) BeforeUpdate(tx *gorm.DB) error {
	if err := s.ObjectMeta.BeforeUpdate(tx); err != nil {
		return fmt.Errorf("failed to run `BeforeUpdate` hook: %v", err)
	}

	if s.Extend == nil {
		s.Extend = make(map[string]any)
	}

	if s.AppSource == AppSourceExternal {
		s.Extend["remote_config"] = s.RemoteConfig
	}
	s.ExtendShadow = json.ToString(s.Extend)

	return nil
}

func (s *AppStore) AfterFind(tx *gorm.DB) error {
	if err := s.ObjectMeta.AfterFind(tx); err != nil {
		return errors.Errorf("failed to run `AfterFind` hook: %v", err)
	}

	return nil
}

// +k8s:deepcopy-gen=true
type AppStoreAuthInfo struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

const (
	ApplicationInstanceStateWaiting          = "waiting"
	ApplicationInstanceStateActive           = "active"
	ApplicationInstanceStateUnknown          = "unknown"
	ApplicationInstanceStateComponentMissing = "component missing"
)

type ApplicationInstance struct {
	imachinery.ObjectMeta
	// 部署的集群

	// 部署的集群命名空间
	Namespace string `json:"string"`
	// 应用实例的当前版本
	Revision int `json:"revision"`
	// 状态
	State string `json:"state"`
	// 状态信息
	StateMessage string `json:"state_message"`
	// 应用升级历史

	ClusterID                    string                        `json:"cluster"`
	ApplicationInstanceRevisions []ApplicationInstanceRevision `json:"histories"`
	AppStoreID                   string                        `json:"app_store_id"`
	ApplicationTemplateID        string                        `json:"application_template_id"`
}

func (s ApplicationInstance) FuzzyFields() []string {
	return []string{s.Name, s.State, s.Description}
}

func (i *ApplicationInstance) GetHistory(revision int) (*ApplicationInstanceRevision, error) {
	for _, v := range i.ApplicationInstanceRevisions {
		if v.Revision == revision {
			return &v, nil
		}
	}
	return nil, errors.Errorf("app instance revision %v not exist", revision)
}

type ApplicationInstanceRevision struct {
	imachinery.ObjectMeta
	// 实例版本
	Revision int `json:"revision"`
	// 实例模板值s
	Values string `json:"values"`
	// 实例模板版本
	TemplateVersion string `json:"template_version"`
	// 应用程序版本
	AppVersion string `json:"app_version"`
	// 应用模板包数据
	PackageData string `json:"package_data"`

	ApplicationInstanceID string `json:"application_instance_id"`
}

const (
	ApplicationTemplateTypeHelm = "helm"
	AppCategoryUncategoriedUUID = "uncategoried"
	AppCategoryUncategoriedName = "uncategoried"
	//application template state
	ApplicationTemplateStatusDeveloping = "developing"  //规划中
	ApplicationTemplateStatusOnShelves  = "on_shelves"  //上架
	ApplicationTemplateStatusOffShelves = "off_shelves" //下架
	//application template version state
	ApplicationTemplateVersionStatusDeveloping     = "developing"  // 开发中,待提交
	ApplicationTemplateVersionStatusApproving      = "approving"   // 等待审核
	ApplicationTemplateVersionStatusApproved       = "approved"    // 通过审核
	ApplicationTemplateVersionStatusApprovalReject = "rejected"    // 拒绝通过
	ApplicationTemplateVersionStatusOffShelves     = "off_shelves" // 已下架
	ApplicationTemplateVersionStatusOnShelves      = "on_shelves"  // 已上架

	ApplicationTemplateVersionActionRelease    = "release"     //发布到应用商店
	ApplicationTemplateVersionActionReject     = "reject"      //拒绝通过审核
	ApplicationTemplateVersionActionCancel     = "cancel"      //撤销审核申请
	ApplicationTemplateVersionActionSubmit     = "submit"      //提交审核
	ApplicationTemplateVersionActionPass       = "pass"        //通过审核
	ApplicationTemplateVersionActionOnShelves  = "on_shelves"  //上架
	ApplicationTemplateVersionActionOffShelves = "off_shelves" //下架
	ApplicationTemplateVersionActionCreate     = "create"      //创建
	ApplicationTemplateVersionActionDelete     = "delete"      //删除

	ApplicationTemplateSortByCreateTime        = "create_time" //sort by create_time
	ApplicationTemplateSortByCreateTimeReverse = "create_time_reverse"
	ApplicationTemplateSortByVersion           = "version"
	ApplicationTemplateSortByUpdateTime        = "update_time"
)

type ApplicationTemplate struct {
	imachinery.ObjectMeta

	Icon      string          `json:"icon" form:"icon"`
	State     string          `json:"state" form:"state"`
	StateTime imachinery.Time `json:"state_time" form:"state_time"`

	AppStoreID string                       `json:"app_store_id"  form:"app_store_id" binding:"required"`
	Versions   []ApplicationTemplateVersion `json:"versions"`
	Instances  []ApplicationInstance        `json:"instances"`
	Categories []ApplicationCategory        `json:"categories" gorm:"many2many:application_template_categories;"`
}

func (obj ApplicationTemplate) FuzzyFields() []string {
	return []string{obj.Name, obj.Description, obj.State}
}

func (obj *ApplicationTemplate) BeforeCreate(tx *gorm.DB) error {
	if err := obj.ObjectMeta.BeforeCreate(tx); err != nil {
		return err
	}
	obj.Instances = nil
	obj.Versions = nil
	obj.State = ApplicationTemplateStatusDeveloping
	obj.StateTime = imachinery.Now()

	return nil
}

// BeforeUpdate run before update database record.
func (obj *ApplicationTemplate) BeforeUpdate(tx *gorm.DB) error {
	if err := obj.ObjectMeta.BeforeUpdate(tx); err != nil {
		return err
	}

	obj.StateTime = imachinery.Now()

	obj.Instances = nil
	obj.Versions = nil

	return nil
}

type ApplicationTemplateVersion struct {
	imachinery.ObjectMeta

	Metadata      *chart.Metadata `json:"metadata" gorm:"-"`
	Name          string          `json:"name"`
	PackageData   string          `json:"package_data"`
	Operator      string          `json:"operator"`
	StateOperator string          `json:"state_operator"`
	State         string          `json:"status"`
	StateTime     imachinery.Time `json:"status_time"`
	UpdateLog     string          `json:"update_log"`

	ApplicationTemplateID string `json:"application_template_id"`
}

func (obj *ApplicationTemplateVersion) BeforeCreate(tx *gorm.DB) error {
	if err := obj.ObjectMeta.BeforeCreate(tx); err != nil {
		return err
	}
	obj.State = ApplicationTemplateStatusDeveloping
	obj.StateTime = imachinery.Now()

	if obj.Extend == nil {
		obj.Extend = make(map[string]any)
	}
	obj.Extend["chart_meta"] = json.ToString(obj.Metadata)
	obj.ExtendShadow = obj.Extend.String()
	return nil
}

// BeforeUpdate run before update database record.
func (obj *ApplicationTemplateVersion) BeforeUpdate(tx *gorm.DB) error {
	obj.Extend["chart_meta"] = json.ToString(obj.Metadata)
	obj.ExtendShadow = obj.Extend.String()
	obj.UpdatedAt = imachinery.Now()

	return nil
}

// AfterFind run after find to unmarshal an extend shadow string into mExtend struct.
func (obj *ApplicationTemplateVersion) AfterFind(tx *gorm.DB) error {
	if err := json.Unmarshal([]byte(obj.ExtendShadow), &obj.Extend); err != nil {
		return err
	}

	v, _ := obj.Extend["chart_meta"]
	if err := json.Unmarshal([]byte(v.(string)), &obj.Metadata); err != nil {
		return err
	}
	return nil
}

func (s ApplicationTemplateVersion) FuzzyFields() []string {
	return []string{s.Name, s.State}
}

const (
	ApplicationCategoryUncategoriedName = "uncategoried"
)

// +k8s:deepcopy-gen=true
type ApplicationCategory struct {
	imachinery.ObjectMeta

	AppStoreID string                `json:"app_store_id" binding:"required"`
	Templates  []ApplicationTemplate `json:"templates" gorm:"many2many:application_template_categories;"`
}

type ApplicationTemplateCategory struct {
	ApplicationTemplateID string    `gorm:"primaryKey;type:varchar(36)"`
	ApplicationCategoryID string    `gorm:"primaryKey;type:varchar(36)"`
	CreatedAt             time.Time `gorm:"created_at"`
}
