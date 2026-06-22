package aliyun

import (
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

func (c *Ecs) ListSnapshot(
	clusterUUID string,
	req *ecs20140526.DescribeSnapshotsRequest,
	opts ...*util.RuntimeOptions,
) (*ecs20140526.DescribeSnapshotsResponse, error) {
	ret, err := c.DescribeSnapshotsWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) ListSnapshotLinks(
	clusterUUID string,
	req *ecs20140526.DescribeSnapshotLinksRequest,
	opts ...*util.RuntimeOptions,
) (*ecs20140526.DescribeSnapshotLinksResponse, error) {
	ret, err := c.DescribeSnapshotLinksWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) ListSnapshotAutoPolicy(
	clusterUUID string,
	req *ecs20140526.DescribeAutoSnapshotPolicyExRequest,
	opts ...*util.RuntimeOptions,
) (*ecs20140526.DescribeAutoSnapshotPolicyExResponse, error) {
	ret, err := c.DescribeAutoSnapshotPolicyExWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) ListSnapshotGroups(
	clusterUUID string,
	req *ecs20140526.DescribeSnapshotGroupsRequest,
) (*ecs20140526.DescribeSnapshotGroupsResponse, error) {
	ret, err := c.DescribeSnapshotGroupsWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)

}

func (c *Ecs) ListSnapshotUsage(
	clusterUUID string,
	req *ecs20140526.DescribeSnapshotsUsageRequest,
) (*ecs20140526.DescribeSnapshotsUsageResponse, error) {
	ret, err := c.DescribeSnapshotsUsageWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) DeleteSnapshot(
	clusterUUID string,
	RegionId *string,
	req *ecs20140526.DeleteSnapshotRequest,
) (*ecs20140526.DeleteSnapshotResponse, error) {
	ret, err := c.DeleteSnapshotWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) CopySnapshot(
	clusterUUID string,
	req *ecs20140526.CopySnapshotRequest,
) (*ecs20140526.CopySnapshotResponse, error) {
	ret, err := c.CopySnapshotWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) CreateSnapshot(
	clusterUUID string,
	RegionId *string,
	req *ecs20140526.CreateSnapshotRequest,
) (*ecs20140526.CreateSnapshotResponse, error) {
	ret, err := c.CreateSnapshotWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) CreateSnapshotAutoPolicy(
	clusterUUID string,
	req *ecs20140526.CreateAutoSnapshotPolicyRequest,
) (*ecs20140526.CreateAutoSnapshotPolicyResponse, error) {
	ret, err := c.CreateAutoSnapshotPolicyWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) ApplySnapshotAutoPolicy(
	clusterUUID string,
	req *ecs20140526.ApplyAutoSnapshotPolicyRequest,
) (*ecs20140526.ApplyAutoSnapshotPolicyResponse, error) {
	ret, err := c.ApplyAutoSnapshotPolicyWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) CancelSnapshotAutoPolicy(
	clusterUUID string,
	req *ecs20140526.CancelAutoSnapshotPolicyRequest,
) (*ecs20140526.CancelAutoSnapshotPolicyResponse, error) {
	ret, err := c.CancelAutoSnapshotPolicyWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) DeleteSnapshotAutoPolicy(
	clusterUUID string,
	req *ecs20140526.DeleteAutoSnapshotPolicyRequest,
) (*ecs20140526.DeleteAutoSnapshotPolicyResponse, error) {
	ret, err := c.DeleteAutoSnapshotPolicyWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) ModifySnapshotAutoPolicy(
	clusterUUID string,
	req *ecs20140526.ModifyAutoSnapshotPolicyExRequest,
) (*ecs20140526.ModifyAutoSnapshotPolicyExResponse, error) {
	ret, err := c.ModifyAutoSnapshotPolicyExWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) ModifySnapshotAttribute(
	clusterUUID string,
	RegionId *string,
	req *ecs20140526.ModifySnapshotAttributeRequest,
) (*ecs20140526.ModifySnapshotAttributeResponse, error) {
	ret, err := c.ModifySnapshotAttributeWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}

func (c *Ecs) ModifySnapshotGroup(
	clusterUUID string,
	req *ecs20140526.ModifySnapshotGroupRequest,
) (*ecs20140526.ModifySnapshotGroupResponse, error) {
	ret, err := c.ModifySnapshotGroupWithOptions(req, &util.RuntimeOptions{})
	return ret, errors.WithStack(err)
}
