package aliyun

import (
	"context"
	"fmt"

	"github.com/wangweihong/gotoolbox/pkg/errors"
	"github.com/wangweihong/gotoolbox/pkg/typeutil"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
)

type Ecs struct {
	*ecs20140526.Client
}

func NewEcs(ak, sk string, region string) (*Ecs, error) {
	config := &openapi.Config{
		AccessKeyId:     typeutil.String(ak),
		AccessKeySecret: typeutil.String(sk),
	}

	config.Endpoint = typeutil.String(fmt.Sprintf("ecs.%s.aliyuncs.com", region))
	c, err := ecs20140526.NewClient(config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Ecs{Client: c}, nil
}

// EcsList 查询弹性云主机列表
func (c *Ecs) EcsList(
	ctx context.Context,
	req *ecs20140526.DescribeInstancesRequest,
	opt ...util.RuntimeOptions,
) (*ecs20140526.DescribeInstancesResponse, error) {
	ret, err := c.DescribeInstancesWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) ListVmStatus(

	req *ecs20140526.DescribeInstanceStatusRequest,
) (*ecs20140526.DescribeInstanceStatusResponse, error) {
	ret, err := c.DescribeInstanceStatusWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) GetVmAttributes(

	RegionId *string,
	req *ecs20140526.DescribeInstanceAttributeRequest,
) (*ecs20140526.DescribeInstanceAttributeResponse, error) {
	ret, err := c.DescribeInstanceAttributeWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) StartVm(

	RegionId *string,
	req *ecs20140526.StartInstanceRequest,
) (*ecs20140526.StartInstanceResponse, error) {
	ret, err := c.StartInstanceWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) BatchStartVm(

	req *ecs20140526.StartInstancesRequest,
) (*ecs20140526.StartInstancesResponse, error) {
	ret, err := c.StartInstancesWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) StopVm(

	RegionId *string,
	req *ecs20140526.StopInstanceRequest,
) (*ecs20140526.StopInstanceResponse, error) {
	ret, err := c.StopInstanceWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) BatchStopVm(

	req *ecs20140526.StopInstancesRequest,
) (*ecs20140526.StopInstancesResponse, error) {
	ret, err := c.StopInstancesWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) RebootVm(

	RegionId *string,
	req *ecs20140526.RebootInstanceRequest,
) (*ecs20140526.RebootInstanceResponse, error) {
	ret, err := c.RebootInstanceWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) BatchRebootVm(

	req *ecs20140526.RebootInstancesRequest,
) (*ecs20140526.RebootInstancesResponse, error) {
	ret, err := c.RebootInstancesWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) CreateVm(

	req *ecs20140526.CreateInstanceRequest,
) (*ecs20140526.CreateInstanceResponse, error) {
	ret, err := c.CreateInstanceWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) DeleteVm(

	RegionId *string,
	req *ecs20140526.DeleteInstanceRequest,
) (*ecs20140526.DeleteInstanceResponse, error) {
	ret, err := c.DeleteInstanceWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) BatchDeleteVm(

	req *ecs20140526.DeleteInstancesRequest,
) (*ecs20140526.DeleteInstancesResponse, error) {
	ret, err := c.DeleteInstancesWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// 修改元数据
func (c *Ecs) ModifyVmMetadata(

	req *ecs20140526.ModifyInstanceMetadataOptionsRequest,
) (*ecs20140526.ModifyInstanceMetadataOptionsResponse, error) {
	ret, err := c.ModifyInstanceMetadataOptionsWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// 修改元数据
func (c *Ecs) ModifyVmAttribute(

	RegionId *string,
	req *ecs20140526.ModifyInstanceAttributeRequest,
) (*ecs20140526.ModifyInstanceAttributeResponse, error) {
	ret, err := c.ModifyInstanceAttributeWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// 修改元数据
func (c *Ecs) ModifyVncPassword(

	req *ecs20140526.ModifyInstanceVncPasswdRequest,
) (*ecs20140526.ModifyInstanceVncPasswdResponse, error) {
	ret, err := c.ModifyInstanceVncPasswdWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// 查看实例规格族列表
func (c *Ecs) ListVmSpecFamilyType(

	req *ecs20140526.DescribeInstanceTypeFamiliesRequest,
) (*ecs20140526.DescribeInstanceTypeFamiliesResponse, error) {
	ret, err := c.DescribeInstanceTypeFamiliesWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// 查看实例规格族列表
func (c *Ecs) GetVmSpecFamilyAttribute(

	RegionId *string,
	req *ecs20140526.DescribeInstanceTypesRequest,
) (*ecs20140526.DescribeInstanceTypesResponse, error) {
	ret, err := c.DescribeInstanceTypesWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) GetVmVNCURL(

	RegionId *string,
	req *ecs20140526.DescribeInstanceVncUrlRequest,
) (*ecs20140526.DescribeInstanceVncUrlResponse, error) {
	ret, err := c.DescribeInstanceVncUrlWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}
