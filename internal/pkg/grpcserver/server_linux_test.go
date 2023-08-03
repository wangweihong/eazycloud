package grpcserver_test

import "github.com/wangweihong/eazycloud/internal/pkg/grpcserver"

func testUnixSocket(conf *grpcserver.GRPCConfig, addr string) {
	testInstallApi(conf, "unix://"+conf.UnixSocket)
}
