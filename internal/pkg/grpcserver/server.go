package grpcserver

import (
	"fmt"
	"net"

	"github.com/wangweihong/eazycloud/internal/pkg/grpcserver/interceptor"

	"google.golang.org/grpc/reflection"

	"github.com/wangweihong/eazycloud/internal/pkg/grpcserver/service/debugservice"
	"github.com/wangweihong/eazycloud/internal/pkg/grpcserver/service/versionservice"

	"golang.org/x/sync/errgroup"

	"github.com/wangweihong/eazycloud/pkg/log"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	*grpc.Server
	address string
	// install services
	Version bool
	Reflect bool
	Debug   bool
	// install interceptors
	UnaryInterceptors  []string
	StreamInterceptors []string
}

func (s *GRPCServer) Run() {
	var eg errgroup.Group
	eg.Go(func() error {
		listen, err := net.Listen("tcp", s.address)
		if err != nil {
			return fmt.Errorf("failed to listen: %w", err)
		}
		log.Infof("start gRPC server at %s", s.address)

		if err := s.Serve(listen); err != nil {
			return fmt.Errorf("failed to start grpc server: %w", err)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		log.Fatal(err.Error())
	}
}

func (s *GRPCServer) Close() {
	s.GracefulStop()
	log.Infof("gRPC server on %s stopped", s.address)
}

// 安装通用服务的中间件和api
// 1. 这里安装的api仅会被提前安装的插件所影响
// 2. 这里安装的中间件会影响后续所有的接口。如果不希望这里有影响, 可以将中间件和通用路由特性等选项关闭。
func initGenericGRPCServer(s *GRPCServer) {
	s.InstallAPIs()
}

func (s *GRPCServer) InstallAPIs() {
	if s.Reflect {
		reflection.Register(s)
	}

	if s.Version {
		versionservice.RegisterVersionService(s.Server)
	}

	if s.Debug {
		debugservice.RegisterDebugServer(s.Server)
	}

	log.Info(
		"gRPC run with service",
		log.Bool("reflect", s.Reflect),
		log.Bool("version", s.Version),
		log.Bool("debug", s.Debug),
	)
}

func installInterceptors(interceptors []string, opt []grpc.ServerOption) []grpc.ServerOption {
	// panic recovery option
	chainUnaryInterceptor := []grpc.UnaryServerInterceptor{}
	// streamInterceptor := []grpc.StreamServerInterceptor{}

	//// 定制panic recover interceptor
	//recoveryOptions := []recovery.Option{
	//	recovery.WithRecoveryHandlerContext(recovery.CustomPanicHandler),
	//}

	// install custom interceptors
	for _, m := range interceptors {
		mw, ok := interceptor.UnaryServerInterceptorList[m]
		if !ok {
			log.Warnf("can not find  unary server interceptor: %s", m)
			continue
		}

		log.Infof("install unary server interceptors: %s", m)
		chainUnaryInterceptor = append(chainUnaryInterceptor, mw)
	}
	opt = append(opt, grpc.ChainUnaryInterceptor(chainUnaryInterceptor...))
	//opt = append(opt,grpc.ChainUnaryInterceptor(chainUnaryInterceptor...))
	//chainUnaryInterceptor = append(chainUnaryInterceptor,
	//	// otelgrpc.UnaryServerInterceptor(),
	//	requestid.UnaryServerInterceptor(),
	//	context.UnaryServerInterceptor(),
	//	logging.UnaryServerInterceptor(),
	//	recovery.UnaryServerInterceptor(recoveryOptions...),
	//)

	//streamUnaryInterceptor := []grpc.StreamServerInterceptor{}
	//streamUnaryInterceptor = append(streamUnaryInterceptor,
	//	recovery.StreamServerInterceptor(recoveryOptions...),
	//)
	return opt
}
