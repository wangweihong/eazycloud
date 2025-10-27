package deployment

import (
	"github.com/wangweihong/eazycloud/apis/ikubeagent"
)

type DeployConfig struct {
}

type AgentService interface {
	InitMaster(config *ikubeagent.KubernetesDeployConfig, isMaster0 bool) (err error)
	JoinCluster(config *ikubeagent.KubernetesDeployConfig, joinCmd string, isControlPlane bool) (err error)
	Reset(isHa bool) error
	GetJoinCommand(showControlPlane bool) (string, string, error)
	GetKubeConfig() (string, error)
	GetDeployLog() (string, error)
	CheckDependency(version string) error
	Version() string
}

func NewInvalidKubeadm(err error) AgentService {
	return &invalidKubeadm{
		err: err,
	}
}

type invalidKubeadm struct {
	err error
}

func (p invalidKubeadm) Version() string {
	return ""
}

func (p invalidKubeadm) InitMaster(config *ikubeagent.KubernetesDeployConfig, isMaster0 bool) (err error) {
	return p.err
}

func (p invalidKubeadm) JoinCluster(
	config *ikubeagent.KubernetesDeployConfig,
	joinCmd string,
	isControlPlane bool,
) (err error) {
	return p.err
}

func (p invalidKubeadm) Reset(isHa bool) error {
	return p.err
}

func (p invalidKubeadm) GetJoinCommand(showControlPlane bool) (string, string, error) {
	return "", "", p.err
}

func (p invalidKubeadm) GetKubeConfig() (string, error) {
	return "", p.err
}

func (p invalidKubeadm) GetDeployLog() (string, error) {
	return "", p.err
}

func (p invalidKubeadm) CheckDependency(version string) error {
	return p.err
}
