package example_server

import (
	"github.com/wangweihong/eazycloud/internal/exampleserver/config"
	"github.com/wangweihong/eazycloud/pkg/httpsvr"
	"github.com/wangweihong/eazycloud/pkg/log"
	"github.com/wangweihong/eazycloud/pkg/shutdown"
	"github.com/wangweihong/eazycloud/pkg/shutdown/managers/posixsignal"
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
