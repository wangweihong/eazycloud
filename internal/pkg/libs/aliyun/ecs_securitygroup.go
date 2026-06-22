package aliyun

import (
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

// 查看安全组
func (c *Ecs) DescribeSecurityGroups(
	clusterUUID string,
	req *ecs20140526.DescribeSecurityGroupsRequest,
) (*ecs20140526.DescribeSecurityGroupsResponse, error) {
	ret, err := c.DescribeSecurityGroupsWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)

}

func (c *Ecs) CreateSecurityGroup(
	clusterUUID string,
	req *ecs20140526.CreateSecurityGroupRequest,
) (*ecs20140526.CreateSecurityGroupResponse, error) {
	ret, err := c.CreateSecurityGroupWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)

}

// 增加安全组入方向规则
func (c *Ecs) AuthorizeSecurityGroup(
	clusterUUID string,
	req *ecs20140526.AuthorizeSecurityGroupRequest,
) (*ecs20140526.AuthorizeSecurityGroupResponse, error) {
	ret, err := c.AuthorizeSecurityGroupWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// 增加安全组出方向规则
func (c *Ecs) AuthorizeSecurityGroupEgress(
	clusterUUID string,
	req *ecs20140526.AuthorizeSecurityGroupEgressRequest,
) (*ecs20140526.AuthorizeSecurityGroupEgressResponse, error) {
	ret, err := c.AuthorizeSecurityGroupEgressWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// 移除安全组入方向规则
func (c *Ecs) RevokeSecurityGroup(
	clusterUUID string,
	req *ecs20140526.RevokeSecurityGroupRequest,
) (*ecs20140526.RevokeSecurityGroupResponse, error) {
	ret, err := c.RevokeSecurityGroupWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// 移除安全组出方向规则
func (c *Ecs) RevokeSecurityGroupEgress(
	clusterUUID string,
	req *ecs20140526.RevokeSecurityGroupEgressRequest,
) (*ecs20140526.RevokeSecurityGroupEgressResponse, error) {
	ret, err := c.RevokeSecurityGroupEgressWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// 将实例加入到指定安全组
func (c *Ecs) JoinSecurityGroup(
	clusterUUID string,
	req *ecs20140526.JoinSecurityGroupRequest,
) (*ecs20140526.JoinSecurityGroupResponse, error) {
	ret, err := c.JoinSecurityGroupWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)

}

// 将实例移出指定安全组
func (c *Ecs) LeaveSecurityGroup(
	clusterUUID string,
	req *ecs20140526.LeaveSecurityGroupRequest,
) (*ecs20140526.LeaveSecurityGroupResponse, error) {

	ret, err := c.LeaveSecurityGroupWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)

}

// 删除指定安全组
func (c *Ecs) DeleteSecurityGroup(
	clusterUUID string,
	req *ecs20140526.DeleteSecurityGroupRequest,
) (*ecs20140526.DeleteSecurityGroupResponse, error) {
	ret, err := c.DeleteSecurityGroupWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)

}

// 查询安全组的规则
func (c *Ecs) DescribeSecurityGroupAttribute(
	clusterUUID string,
	req *ecs20140526.DescribeSecurityGroupAttributeRequest,
) (*ecs20140526.DescribeSecurityGroupAttributeResponse, error) {
	ret, err := c.DescribeSecurityGroupAttributeWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// 查询安全组的规则
func (c *Ecs) ModifySecurityGroupAttribute(
	clusterUUID string,
	req *ecs20140526.ModifySecurityGroupAttributeRequest,
) (*ecs20140526.ModifySecurityGroupAttributeResponse, error) {
	ret, err := c.ModifySecurityGroupAttributeWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)

}

// 查询安全组的联通策略
func (c *Ecs) ModifySecurityGroupPolicy(
	clusterUUID string,
	req *ecs20140526.ModifySecurityGroupPolicyRequest,
) (*ecs20140526.ModifySecurityGroupPolicyResponse, error) {
	ret, err := c.ModifySecurityGroupPolicyWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// 修改安全组的入方向规则的描述信息
func (c *Ecs) ModifySecurityGroupRule(
	clusterUUID string,
	req *ecs20140526.ModifySecurityGroupRuleRequest,
) (*ecs20140526.ModifySecurityGroupRuleResponse, error) {
	ret, err := c.ModifySecurityGroupRuleWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

// 修改安全组的出方向规则的描述信息
func (c *Ecs) ModifySecurityGroupEgressRule(
	clusterUUID string,
	req *ecs20140526.ModifySecurityGroupEgressRuleRequest,
) (*ecs20140526.ModifySecurityGroupEgressRuleResponse, error) {
	ret, err := c.ModifySecurityGroupEgressRuleWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}
