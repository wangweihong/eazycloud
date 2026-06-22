package deployment

import (
	"context"
	gerrors "errors"
	"os"
	"strings"

	"github.com/wangweihong/eazycloud/apis/ikubeagent"
	"github.com/wangweihong/eazycloud/internal/kubeagent/service/v1/deployment/version130"
	"github.com/wangweihong/eazycloud/internal/kubeagent/store"
	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/executil"
	"gorm.io/gorm"
)

type DeployConfig struct {
}

type AgentService interface {
	InitMaster(ctx context.Context, param *ikubeagent.InstallMasterRequest) (err error)
	JoinCluster(ctx context.Context, param *ikubeagent.JoinClusterRequest) (err error)
	Reset(ctx context.Context) error
	GetJoinCommand(ctx context.Context, param *ikubeagent.GetJoinCommandRequest) (string, string, error)
	GetKubeConfig(ctx context.Context) (string, error)
	GetDeployLog(ctx context.Context) (string, error)
	CheckDependency(ctx context.Context, version string) error
	GetInstallState(ctx context.Context) (*ikubeagent.InstallStateResponse, error)
	SetOwner(ctx context.Context, owner string) error
	Version() string
}

func NewService(store store.Factory) AgentService {
	kubedamPath := "/usr/bin/kubeadm"

	if _, err := os.Stat(kubedamPath); err != nil && os.IsNotExist(err) {
		return NewInvalidKubeadm(err)
	}

	state, err := store.InstallStateStores().GetByName(context.Background(), ikubeagent.KubernetesInstallStateUniqueName)
	if err != nil && !gerrors.Is(err, gorm.ErrRecordNotFound) {
		return NewInvalidKubeadm(err)
	}

	if err != nil {
		state, err = store.InstallStateStores().Upsert(context.Background(), &ikubeagent.InstallState{
			State: string(ikubeagent.KubernetesDeployStateUninitialized),
		})
		if err != nil {
			return NewInvalidKubeadm(err)
		}
	}

	output, err := executil.ExecuteTimeout("kubeadm", []string{"version", "-o", "short"}, 10)
	if err != nil {
		return NewInvalidKubeadm(err)
	}
	version := strings.TrimSuffix(output, "\n")
	switch version {
	case "v1.30.0":
		return version130.NewDeployService(version, state, store)
		// case "v1.18.0":
		// 	return version118.NewDeployService(version)
	}
	err = errors.Errorf("unsupport kubeadm version:%v", version)
	return NewInvalidKubeadm(err)
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

func (p invalidKubeadm) InitMaster(ctx context.Context, config *ikubeagent.InstallMasterRequest) (err error) {
	return p.err
}

func (p invalidKubeadm) JoinCluster(
	ctx context.Context,
	config *ikubeagent.JoinClusterRequest,
) (err error) {
	return p.err
}

func (p invalidKubeadm) Reset(ctx context.Context) error {
	return p.err
}

func (p invalidKubeadm) GetJoinCommand(ctx context.Context, param *ikubeagent.GetJoinCommandRequest) (string, string, error) {
	return "", "", p.err
}

func (p invalidKubeadm) GetKubeConfig(ctx context.Context) (string, error) {
	return "", p.err
}

func (p invalidKubeadm) GetDeployLog(ctx context.Context) (string, error) {
	return "", p.err
}

func (p invalidKubeadm) CheckDependency(ctx context.Context, version string) error {
	return p.err
}

func (p invalidKubeadm) GetInstallState(ctx context.Context) (*ikubeagent.InstallStateResponse, error) {
	return nil, p.err
}

func (p invalidKubeadm) SetOwner(ctx context.Context, owner string) error {
	return p.err
}
