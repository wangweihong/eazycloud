package grpcserver_test

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/wangweihong/eazycloud/internal/pkg/grpcserver"
	"github.com/wangweihong/eazycloud/internal/pkg/grpcserver/apis/debug"
	"github.com/wangweihong/eazycloud/internal/pkg/grpcserver/apis/version"
	pkgversion "github.com/wangweihong/eazycloud/pkg/version"
)

func testInstallApi(conf *grpcserver.GRPCConfig, addr string) {
	s, err := conf.Complete().New()
	So(err, ShouldBeNil)
	go func() {
		s.Run()
	}()
	// Wait for the server to start (you can use a more sophisticated wait mechanism)
	time.Sleep(3 * time.Second)

	// Set up a gRPC connection to the server
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	So(err, ShouldBeNil)
	defer conn.Close()

	_, err = debug.NewDebugServiceClient(conn).
		Sleep(context.Background(), &debug.SleepRequest{Duration: durationpb.New(50 * time.Millisecond)})
	So(err, ShouldBeNil)

	versionResp, err := version.NewVersionServiceClient(conn).Version(context.Background(), &version.VersionRequest{})
	So(err, ShouldBeNil)
	So(versionResp.GitCommit, ShouldEqual, pkgversion.Get().GitCommit)

	s.Close()
}

func TestGRPCServer_InstallAPI(t *testing.T) {
	Convey("grpc通用服务测试", t, func() {
		conf := grpcserver.NewConfig()
		conf.Debug = true
		conf.Version = true
		conf.Reflect = true
		// 必须设置. 不设置将会遇到rpc error: code = ResourceExhausted desc = grpc: received message larger than max (7 vs. 0)
		conf.MaxMsgSize = 4 * 1024 * 1024

		Convey("测试version,debug", func() {
			Convey("tcp连接", func() {
				// 随机端口
				conf.Addr = "0.0.0.0:56218"
				testInstallApi(conf, conf.Addr)
			})

			// 注意, windows不支持unix domain socket
			Convey("unix socket", func() {
				conf.UnixSocket = "/tmp/test.socket"
				testUnixSocket(conf, "unix://"+conf.UnixSocket)
			})
		})
	})
}
