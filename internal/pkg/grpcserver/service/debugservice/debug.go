package debugservice

import (
	"context"
	"time"

	"github.com/wangweihong/eazycloud/pkg/log"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	"github.com/wangweihong/eazycloud/internal/pkg/grpcserver/apis/debug"
)

type debugService struct{}

// Panic trigger panic for test.
func (s *debugService) Panic(context.Context, *empty.Empty) (*empty.Empty, error) {
	panic("panic")
}

// Sleep sleep from specific duration.
func (s *debugService) Sleep(ctx context.Context, req *debug.SleepRequest) (*empty.Empty, error) {
	d := 30 * time.Second
	if req != nil && req.Duration != nil {
		d = req.Duration.AsDuration()
	}

	start := time.Now()
	log.F(ctx).Info("sleep")
	time.Sleep(d)
	log.F(ctx).Infof("awake,cost:%s", time.Since(start))

	return &empty.Empty{}, nil
}

// RegisterDebugServer  register debug service to gRPC.
func RegisterDebugServer(s *grpc.Server) {
	debug.RegisterDebugServiceServer(s, &debugService{})
}
