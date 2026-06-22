package aliyun

import (
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

// 删除密钥对
func (c *Ecs) DeleteKeyPairs(
	clusterUUID string,
	req *ecs20140526.DeleteKeyPairsRequest,
) (*ecs20140526.DeleteKeyPairsResponse, error) {
	ret, err := c.DeleteKeyPairsWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) DescribeKeyPairs(
	clusterUUID string,
	req *ecs20140526.DescribeKeyPairsRequest,
) (*ecs20140526.DescribeKeyPairsResponse, error) {
	ret, err := c.DescribeKeyPairsWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// 删除密钥对
func (c *Ecs) CreateKeyPair(
	clusterUUID string,
	req *ecs20140526.CreateKeyPairRequest,
) (*ecs20140526.CreateKeyPairResponse, error) {
	ret, err := c.CreateKeyPairWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) ImportKeyPair(
	clusterUUID string,
	req *ecs20140526.ImportKeyPairRequest,
) (*ecs20140526.ImportKeyPairResponse, error) {
	ret, err := c.ImportKeyPairWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) AttachKeyPair(
	clusterUUID string,
	req *ecs20140526.AttachKeyPairRequest,
) (*ecs20140526.AttachKeyPairResponse, error) {
	ret, err := c.AttachKeyPairWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) DetachKeyPair(
	clusterUUID string,
	req *ecs20140526.DetachKeyPairRequest,
) (*ecs20140526.DetachKeyPairResponse, error) {
	ret, err := c.DetachKeyPairWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}
