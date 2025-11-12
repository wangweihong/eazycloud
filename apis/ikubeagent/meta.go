package ikubeagent

import (
	"github.com/wangweihong/eazycloud/apis/imachinery"
	"gorm.io/gorm"
)

const (
	KubernetesInstallStateUniqueName = "kubernetes"
)

type InstallState struct {
	imachinery.ObjectMeta
	State         string          `json:"state"`
	ControlPlane  bool            `json:"control_plane"`
	HighAvailable bool            `json:"high_available"`
	StartTime     imachinery.Time `json:"start_time"`
	EndTime       imachinery.Time `json:"end_time"`
	ErrorMessage  string          `json:"error_message"`
	Owner         string          `json:"owner"`
}

func (obj *InstallState) BeforeCreate(tx *gorm.DB) error {
	obj.State = string(KubernetesDeployStateUninitialized)
	return obj.ObjectMeta.BeforeCreate(tx)
}
