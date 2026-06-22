package kubeagent

import (
	"context"

	"github.com/wangweihong/gotoolbox/pkg/log"
	"github.com/wangweihong/gotoolbox/pkg/shutdown"
	"github.com/wangweihong/gotoolbox/pkg/shutdown/managers/posixsignal"

	"github.com/wangweihong/eazycloud/apis/ikubeagent"
	"github.com/wangweihong/eazycloud/internal/kubeagent/config"
	"github.com/wangweihong/eazycloud/internal/kubeagent/store/postgresql"
	"github.com/wangweihong/eazycloud/pkg/httpsvr"
)

type server struct {
	// api服务,提供http和tls
	httpServer *httpsvr.GenericHTTPServer
	// 控制服务关闭时处理动作, 如捕捉到信号后如何处理
	gracefulShutdown *shutdown.GracefulShutdown
}

// preparedServer is a private wrapper that enforces a call of PrepareRun() before Run can be invoked.
type preparedServer struct {
	*server
}

// 创建服务器实例.
func createServer(cfg *config.Config) (*server, error) {
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	// 构建通用的http(s) server服务配置
	genericConfig, err := buildGenericHTTPServerConfig(cfg)
	if err != nil {
		return nil, err
	}

	// 补全通用服务器配置, 并生成通用服务实例
	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}

	server := &server{
		httpServer:       genericServer,
		gracefulShutdown: gs,
	}

	return server, nil
}

// 根据服务器配置应用到通用服务器配置上.
func buildGenericHTTPServerConfig(cfg *config.Config) (genericConfig *httpsvr.Config, lastErr error) {
	genericConfig = httpsvr.NewConfig()
	if lastErr = cfg.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.FeatureOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.InsecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.SecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	return
}

// PrepareRun prepares the server to run, by setting up the server instance.
func (s *server) PrepareRun() preparedServer {
	initRouter(s.httpServer.Engine)

	storeIns, _ := postgresql.GetPostgresSQLFactoryOr(nil)
	d, _ := storeIns.InstallStateStores().GetByName(context.Background(), ikubeagent.KubernetesInstallStateUniqueName)
	// 上一次正在部署过程中,程序异常退出, 数据库状态没有更新
	if d != nil && d.State == string(ikubeagent.KubernetesDeployStateDeploying) {
		d.State = string(ikubeagent.KubernetesDeployStateError)
		d.ErrorMessage = "program restart"
		if _, err := storeIns.InstallStateStores().Upsert(context.Background(), d); err != nil {
			log.Errorf("last time program restart when deploying, try to set error deploy fail:%v ", err)
		}
	}

	// 设置服务优雅退出回调处理
	s.gracefulShutdown.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		s.httpServer.Close()
		return nil
	}))

	return preparedServer{s}
}

func (s preparedServer) Run(stopCh <-chan struct{}) error {
	// start shutdown managers
	if err := s.gracefulShutdown.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}
	return s.httpServer.Run()
}
