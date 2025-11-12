package aliyun

import (
	"context"
	"fmt"

	"github.com/wangweihong/gotoolbox/pkg/errors"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	vpc20160428 "github.com/alibabacloud-go/vpc-20160428/v6/client"
)

type Vpc struct {
	*vpc20160428.Client
}

func NewVpc(ak, sk string, region string) (*Vpc, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(ak),
		AccessKeySecret: tea.String(sk),
	}

	config.Endpoint = tea.String(fmt.Sprintf("vpc.%s.aliyuncs.com", region))
	c, err := vpc20160428.NewClient(config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Vpc{Client: c}, nil
}

// VSwitchList 查询虚拟交换机列表
func (c *Vpc) VSwitchList(
	ctx context.Context,
	req *vpc20160428.DescribeVSwitchesRequest,
) (*vpc20160428.DescribeVSwitchesResponse, error) {
	ret, err := c.DescribeVSwitchesWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}
