package aliyun

import (
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

// 创建弹性网卡
func (c *Ecs) CreateNetworkInterface(
	clusterUUID string,
	req *ecs20140526.CreateNetworkInterfaceRequest,
) (*ecs20140526.CreateNetworkInterfaceResponse, error) {
	ret, err := c.CreateNetworkInterfaceWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) AttachNetworkInterface(
	clusterUUID string,
	req *ecs20140526.AttachNetworkInterfaceRequest,
) (*ecs20140526.AttachNetworkInterfaceResponse, error) {
	ret, err := c.AttachNetworkInterfaceWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) DetachNetworkInterface(
	clusterUUID string,
	req *ecs20140526.DetachNetworkInterfaceRequest,
) (*ecs20140526.DetachNetworkInterfaceResponse, error) {
	ret, err := c.DetachNetworkInterfaceWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) DeleteNetworkInterface(
	clusterUUID string,
	req *ecs20140526.DeleteNetworkInterfaceRequest,
) (*ecs20140526.DeleteNetworkInterfaceResponse, error) {
	ret, err := c.DeleteNetworkInterfaceWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)

}

func (c *Ecs) DescribeNetworkInterfaces(
	clusterUUID string,
	req *ecs20140526.DescribeNetworkInterfacesRequest,
) (*ecs20140526.DescribeNetworkInterfacesResponse, error) {
	ret, err := c.DescribeNetworkInterfacesWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)

}

func (c *Ecs) DescribeNetworkInterfaceAttribute(
	clusterUUID string,
	req *ecs20140526.DescribeNetworkInterfaceAttributeRequest,
) (*ecs20140526.DescribeNetworkInterfaceAttributeResponse, error) {
	ret, err := c.DescribeNetworkInterfaceAttributeWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)

}

func (c *Ecs) ModifyNetworkInterfaceAttribute(
	clusterUUID string,
	req *ecs20140526.ModifyNetworkInterfaceAttributeRequest,
) (*ecs20140526.ModifyNetworkInterfaceAttributeResponse, error) {
	ret, err := c.ModifyNetworkInterfaceAttributeWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)

}

func (c *Ecs) AssignPrivateIpAddresses(
	clusterUUID string,
	req *ecs20140526.AssignPrivateIpAddressesRequest,
) (*ecs20140526.AssignPrivateIpAddressesResponse, error) {
	ret, err := c.AssignPrivateIpAddressesWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)

}

func (c *Ecs) UnassignPrivateIpAddresses(
	clusterUUID string,
	req *ecs20140526.UnassignPrivateIpAddressesRequest,
) (*ecs20140526.UnassignPrivateIpAddressesResponse, error) {
	ret, err := c.UnassignPrivateIpAddressesWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}
