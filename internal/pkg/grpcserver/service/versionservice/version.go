package versionservice

import (
	"context"

	"google.golang.org/grpc"

	"github.com/wangweihong/eazycloud/internal/pkg/grpcserver/apis/version"
	pkgversion "github.com/wangweihong/eazycloud/pkg/version"
)

type versionService struct{}

func (v versionService) Version(
	ctx context.Context,
	request *version.VersionRequest,
) (*version.VersionResponse, error) {
	info := pkgversion.Get()
	return &version.VersionResponse{
		GitVersion:   info.GitVersion,
		GitCommit:    info.GitCommit,
		GitTreeState: info.GitTreeState,
		BuildDate:    info.BuildDate,
		GoVersion:    info.GoVersion,
		Compiler:     info.Compiler,
		Platform:     info.Platform,
	}, nil
}

func RegisterVersionService(s *grpc.Server) {
	version.RegisterVersionServiceServer(s, &versionService{})
}
