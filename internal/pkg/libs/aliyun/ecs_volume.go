package aliyun

import (
	"context"

	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

// EcsList 查询弹性云主机列表
func (c *Ecs) CreateDisk(
	ctx context.Context,
	req *ecs20140526.CreateDiskRequest,
	opt ...util.RuntimeOptions,
) (*ecs20140526.CreateDiskResponse, error) {
	ret, err := c.CreateDiskWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) DeleteDisk(

	RegionId *string,
	req *ecs20140526.DeleteDiskRequest,
) (*ecs20140526.DeleteDiskResponse, error) {
	ret, err := c.DeleteDiskWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) ListDisks(

	req *ecs20140526.DescribeDisksRequest,
	opts ...*util.RuntimeOptions,
) (*ecs20140526.DescribeDisksResponse, error) {
	ret, err := c.DescribeDisksWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) AttachDisk(

	RegionId *string,
	req *ecs20140526.AttachDiskRequest,
) (*ecs20140526.AttachDiskResponse, error) {
	ret, err := c.AttachDiskWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) DetachDisk(

	RegionId *string,
	req *ecs20140526.DetachDiskRequest,
) (*ecs20140526.DetachDiskResponse, error) {

	ret, err := c.DetachDiskWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) ModifyDiskAttribute(

	req *ecs20140526.ModifyDiskAttributeRequest,
) (*ecs20140526.ModifyDiskAttributeResponse, error) {

	ret, err := c.ModifyDiskAttributeWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) ReplaceSystemDisk(

	RegionId *string,
	req *ecs20140526.ReplaceSystemDiskRequest,
) (*ecs20140526.ReplaceSystemDiskResponse, error) {

	ret, err := c.ReplaceSystemDiskWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)

}

func (c *Ecs) ReInitDisk(
	RegionId *string,
	req *ecs20140526.ReInitDiskRequest,
) (*ecs20140526.ReInitDiskResponse, error) {

	ret, err := c.ReInitDiskWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) ResetDisk(
	RegionId *string,
	req *ecs20140526.ResetDiskRequest,
) (*ecs20140526.ResetDiskResponse, error) {
	ret, err := c.ResetDiskWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) ResizeDisk(

	RegionId *string,
	req *ecs20140526.ResizeDiskRequest,
) (*ecs20140526.ResizeDiskResponse, error) {
	ret, err := c.ResizeDiskWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) ResetDisks(

	req *ecs20140526.ResetDisksRequest,
) (*ecs20140526.ResetDisksResponse, error) {
	ret, err := c.ResetDisksWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) ModifyDiskSpec(

	RegionId *string,
	req *ecs20140526.ModifyDiskSpecRequest,
) (*ecs20140526.ModifyDiskSpecResponse, error) {
	ret, err := c.ModifyDiskSpecWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}
