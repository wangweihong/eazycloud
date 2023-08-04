package grpcserver_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"testing"
	"time"

	"google.golang.org/grpc/credentials"

	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/wangweihong/eazycloud/internal/pkg/grpcserver"
	"github.com/wangweihong/eazycloud/internal/pkg/grpcserver/apis/debug"
	"github.com/wangweihong/eazycloud/internal/pkg/grpcserver/apis/version"
	pkgversion "github.com/wangweihong/eazycloud/pkg/version"
)

const (
	keyData  = "-----BEGIN PRIVATE KEY-----\n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\ne/YTDd9miWTfQ1+TMbcVszqBaaBZ\n-----END PRIVATE KEY-----\n"
	certData = "-----BEGIN CERTIFICATE-----\n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\n-----END CERTIFICATE-----\n"
	caData   = "-----BEGIN CERTIFICATE-----\nMIIFlTCCA32gAwIBAgIUOc6BpN0Oub+CYhSubB5/F9E9GaMwDQYJKoZIhvcNAQEN\nBQAwWjELMAkGA1UEBhMCQ04xEjAQBgNVBAgMCUd1YW5nZG9uZzERMA8GA1UEBwwI\nU2hlbnpoZW4xEjAQBgNVBAoMCUVhenlDbG91ZDEQMA4GA1UECwwHRGV2ZWxvcDAe\nFw0yMzA3MjYwNjQ2NTJaFw0zMzA3MjMwNjQ2NTJaMFoxCzAJBgNVBAYTAkNOMRIw\nEAYDVQQIDAlHdWFuZ2RvbmcxETAPBgNVBAcMCFNoZW56aGVuMRIwEAYDVQQKDAlF\nYXp5Q2xvdWQxEDAOBgNVBAsMB0RldmVsb3AwggIiMA0GCSqGSIb3DQEBAQUAA4IC\nDwAwggIKAoICAQCXS7sY/f2KGF4cis/tcQUyArXpyJ3MgiGpCJmv94GUSIVAzbWU\neuKdmqbh+zBvmGX8Jgan0KAC+2o5/8WYZLw9v7H1Py1DtI12MYW/QaI4+734ZsHc\nyg4IK3rmTGXWR1TLusUdJcywMBSl7BpJ8C1Vr1JomaSFuE/tP9i3fv0BT02lFlHd\n0+AvI9c9Ridhrymnn2qAFY8EKuPmu1lRV8wUB8oIL7lM3/CIFJUkGsmUxdJuTZU4\npuDn8DnTQ886jrPepxe4+j0zweZQbVgfqiQaZ+Ubl48fUtS0HZbthXJUUa260XPY\ns4nkUSGZakn2lOOBu4BUAfkwLSqrqpOX3sMujOiEXZ80BM6jQ55INurqQuMVdWTf\neJr+by5X7plUdF3Hd+7oikw6d2TvM8CoQTVNH+F+MEWJ8sncugBOhPIbJJbW6P9e\nuekQq6r61oyqg4AUGTvi6tDc6UC6FpDnP/oJlei7BCH3fr0ASws1ePyQ+reyp439\neiwMdqEHAzon7VZu7Tot9eLELSsIPb/5EkwV6qAirn3aTKol13M7+YI2UptCSIa/\n5+ruY92OueuZf9Z6VoqTt95lSrlET/femVikXCflFySWR+xbSWjdEY8faSDAIYnt\nEg/5MPcK7k3aGL3BdiU+krS7UAdbVbx8y03Zw8hO7UdfNAD2S1y8BKH2HQIDAQAB\no1MwUTAdBgNVHQ4EFgQU2YYebmvS23/VSPatO1awTf1SjkYwHwYDVR0jBBgwFoAU\n2YYebmvS23/VSPatO1awTf1SjkYwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0B\nAQ0FAAOCAgEAMMRVan8I2h7Rj1eEIrfdL8/5pLT6O+bPSYoNUgnj1EuSYAr+u+sS\nr0l7gFwg1WISQlxRbmyZKQBE44+8Pn/qHVRKC0RUVja8j70HpxY2clLLolWunETQ\n6MbplZ8w4ei1Rx5L6S4Vcz+EAJwmseTJa5B8U69coZzeuiyHAUKmnLsSudJr9bTc\n5vnMOve+eG8Y4EKcpYMwJJy8272eQFNXwKmIrfD/5qTV03aMVcANXxvGpZWBYz5w\nCE/NDsMO/BnRFm4//ml5cKiTppG9u3/94Ah2bz4dATZl8AxwfQ2vOVQqKDXm5XwD\nH2XV82FJDTfAUfYwQZhSzXXwMRYgnKfDyLxmuRrRO1NCN8ddFpH+SLbhNjQ56cJT\n2qvsXJ/n8AaKeAr2mGEJ0d8cy69IxLZSdLDmL081GdhRGTExMnIxfLL0wvOMNlMM\nQokVRGuKShrE8LFNVjTlfSLBmuVKomugCXn+VVJFhRMFkSSg1kcmusEqEWuWsx9Y\n7hMs/MVVRS0YjTjxgTFNFcevNXkY0xoHzur1ccCIArJXroF2UzUctF/dpqMJP6bk\npIT7Pu3UNj/qChMLZ8ostJyhM/24PwkLLHy1v9lU+86lYWZGLhL3QSSnctypI4fF\nGCx0CfIEfjVIKsvaSa4v+JTX/yiSnUj8ChNM7r5I2bDxB7vy2wYPn/Y=\n-----END CERTIFICATE-----\n"
)

func testInstallApi(conf *grpcserver.GRPCConfig, addr string, tlsEnable bool) {
	s, err := conf.Complete().New()
	So(err, ShouldBeNil)
	go func() {
		s.Run()
	}()
	// Wait for the server to start (you can use a more sophisticated wait mechanism)
	time.Sleep(3 * time.Second)

	var conn *grpc.ClientConn

	if !tlsEnable {
		// Set up a gRPC connection to the server
		conn, err = grpc.Dial(addr, grpc.WithInsecure())
		So(err, ShouldBeNil)
	} else {
		certPool := x509.NewCertPool()
		ok := certPool.AppendCertsFromPEM([]byte(caData))
		So(ok, ShouldBeTrue)

		// 创建 gRPC 客户端连接
		creds := credentials.NewTLS(&tls.Config{
			RootCAs: certPool,
		})

		conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(creds))
		So(err, ShouldBeNil)
	}
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
			Convey("TCP连接", func() {
				SkipConvey("tcp连接, 无TLS", func() {
					// 随机端口
					conf.Addr = "0.0.0.0:56218"
					testInstallApi(conf, conf.Addr, false)
				})
				Convey("tcp连接,有TLS", func() {
					// 随机端口
					conf.Addr = "0.0.0.0:56219"
					conf.TlsEnable = true
					conf.ServerCert.Cert = certData
					conf.ServerCert.Key = keyData
					testInstallApi(conf, conf.Addr, true)
				})
			})

			// 注意, windows不支持unix domain socket
			SkipConvey("unix socket", func() {
				conf.UnixSocket = "/tmp/test.socket"
				testUnixSocket(conf, "unix://"+conf.UnixSocket)
			})
		})
	})
}
