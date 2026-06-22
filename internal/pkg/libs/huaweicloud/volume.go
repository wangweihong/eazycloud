package huaweicloud

import (
	"context"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/evs/v2/model"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
	"github.com/wangweihong/gotoolbox/pkg/errors"

	evs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/evs/v2"
)

type Volume struct {
	c   *evs.EvsClient
	err error
}

func NewVolume(cfg *AccessConfig) *Volume {
	c, err := cfg.HcEvsClient(cfg.Region)
	return &Volume{c: c, err: errors.WithStack(err)}
}

// VolumeList 查看镜像列表
func (p *Volume) VolumeList(
	ctx context.Context,
	req *model.ListVolumesRequest,
) (*model.ListVolumesResponse, error) {
	resp, err := libs.InvokeApi(ctx, p.err, req, p.c.ListVolumes)
	return resp, errors.WithStack(err)

}

// VolumeCreate 查看卷列表
func (p *Volume) VolumeCreate(
	ctx context.Context,
	req *model.CreateVolumeRequest,
) (*model.CreateVolumeResponse, error) {
	resp, err := libs.InvokeApi(ctx, p.err, req, p.c.CreateVolume)
	return resp, errors.WithStack(err)
}

// VolumeSnapshotList 查看卷快照列表
func (p *Volume) VolumeSnapshotList(
	ctx context.Context,
	req *model.ListSnapshotsRequest,
) (*model.ListSnapshotsResponse, error) {
	resp, err := libs.InvokeApi(ctx, p.err, req, p.c.ListSnapshots)
	return resp, errors.WithStack(err)
}

// VolumeSnapshotCreate 创建卷快照
func (p *Volume) VolumeSnapshotCreate(
	ctx context.Context,
	req *model.CreateSnapshotRequest,
) (*model.CreateSnapshotResponse, error) {
	resp, err := libs.InvokeApi(ctx, p.err, req, p.c.CreateSnapshot)
	return resp, errors.WithStack(err)
}
