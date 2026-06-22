package tecentcloud

import (
	"context"

	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	tcommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

const (
	//云盘状态
	DiskStateUNATTACHED  = "UNATTACHED"  //表示未挂载
	DiskStateATTACHING   = "ATTACHING"   //表示挂载中
	DiskStateATTACHED    = "ATTACHED"    //表示已挂载
	DiskStateDETACHING   = "DETACHING"   //表示解挂中
	DiskStateEXPANDING   = "EXPANDING"   //表示扩容中
	DiskStateROLLBACKING = "ROLLBACKING" //表示回滚中
	DiskStateTORECYCLE   = "TORECYCLE"   //表示待回收
	DiskStateDUMPING     = "DUMPING"     //表示拷贝硬盘中
	//快照状态
	SnapshotStateNORMAL            = "NORMAL"              //表示正常
	SnapshotStateCREATING          = "CREATING"            //表示创建中
	SnapshotStateROLLBACKING       = "ROLLBACKING"         //表示回滚中
	SnapshotStateCOPYINGFROMREMOTE = "COPYING_FROM_REMOTE" //表示跨地域复制快照拷贝中
)

type EbsClient struct {
	c   *cbs.Client
	err error
}

func NewEbsClient(securityId, securityKey string, region string) *EbsClient {
	credential := tcommon.NewCredential(securityId, securityKey)
	clientProfile := profile.NewClientProfile()
	//clientProfile.HttpProfile.Endpoint = endpoint
	client, err := cbs.NewClient(credential, region, clientProfile)
	return &EbsClient{
		c:   client,
		err: err,
	}
}

func (c *EbsClient) DescribeDisks(ctx context.Context, req *cbs.DescribeDisksRequest) (*cbs.DescribeDisksResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeDisks, libs.WithRateLimit(5, 10))
	return resp, errors.WithStack(err)
}

func (c *EbsClient) CreateDisks(ctx context.Context, req *cbs.CreateDisksRequest) (*cbs.CreateDisksResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.CreateDisks)
	return resp, errors.WithStack(err)
}

func (c *EbsClient) ModifyDiskAttributes(ctx context.Context, req *cbs.ModifyDiskAttributesRequest) (*cbs.ModifyDiskAttributesResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.ModifyDiskAttributes)
	return resp, errors.WithStack(err)
}

func (c *EbsClient) ResizeDisk(ctx context.Context, req *cbs.ResizeDiskRequest) (*cbs.ResizeDiskResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.ResizeDisk)
	return resp, errors.WithStack(err)
}

// func (c *EbsClient) DescribeDiskOperationLogs(req *cbs.DescribeDiskOperationLogsRequest) (*cbs.DescribeDiskOperationLogsResponse, error) {
// 	resp, err := invoke(ctx,c.err, req, c.c.DescribeDiskOperationLogs)
// 	return resp, errors.WithStack(err)
// }

func (c *EbsClient) TerminateDisks(ctx context.Context, req *cbs.TerminateDisksRequest) (*cbs.TerminateDisksResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.TerminateDisks)
	return resp, errors.WithStack(err)
}

func (c *EbsClient) DescribeDiskConfigQuota(ctx context.Context, req *cbs.DescribeDiskConfigQuotaRequest) (*cbs.DescribeDiskConfigQuotaResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeDiskConfigQuota)
	return resp, errors.WithStack(err)
}

func (c *EbsClient) DescribeInstancesDiskNum(ctx context.Context, req *cbs.DescribeInstancesDiskNumRequest) (*cbs.DescribeInstancesDiskNumResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeInstancesDiskNum)
	return resp, errors.WithStack(err)
}

func (c *EbsClient) AttachDisks(ctx context.Context, req *cbs.AttachDisksRequest) (*cbs.AttachDisksResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.AttachDisks)
	return resp, errors.WithStack(err)
}

func (c *EbsClient) DetachDisks(ctx context.Context, req *cbs.DetachDisksRequest) (*cbs.DetachDisksResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DetachDisks)
	return resp, errors.WithStack(err)
}

func (c *EbsClient) DescribeSnapshots(ctx context.Context, req *cbs.DescribeSnapshotsRequest) (*cbs.DescribeSnapshotsResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeSnapshots)
	return resp, errors.WithStack(err)
}

func (c *EbsClient) CreateSnapshot(ctx context.Context, req *cbs.CreateSnapshotRequest) (*cbs.CreateSnapshotResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.CreateSnapshot)
	return resp, errors.WithStack(err)
}

func (c *EbsClient) DeleteSnapshots(ctx context.Context, req *cbs.DeleteSnapshotsRequest) (*cbs.DeleteSnapshotsResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DeleteSnapshots)
	return resp, errors.WithStack(err)
}

func (c *EbsClient) ModifySnapshotAttribute(ctx context.Context, req *cbs.ModifySnapshotAttributeRequest) (*cbs.ModifySnapshotAttributeResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.ModifySnapshotAttribute)
	return resp, errors.WithStack(err)
}

func (c *EbsClient) ApplySnapshot(ctx context.Context, req *cbs.ApplySnapshotRequest) (*cbs.ApplySnapshotResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.ApplySnapshot)
	return resp, errors.WithStack(err)
}

func (c *EbsClient) CreateAutoSnapshotPolicy(ctx context.Context, req *cbs.CreateAutoSnapshotPolicyRequest) (*cbs.CreateAutoSnapshotPolicyResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.CreateAutoSnapshotPolicy)
	return resp, errors.WithStack(err)
}

func (c *EbsClient) DeleteAutoSnapshotPolicies(ctx context.Context, req *cbs.DeleteAutoSnapshotPoliciesRequest) (*cbs.DeleteAutoSnapshotPoliciesResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DeleteAutoSnapshotPolicies)
	return resp, errors.WithStack(err)
}

func (c *EbsClient) ModifyAutoSnapshotPolicyAttribute(ctx context.Context, req *cbs.ModifyAutoSnapshotPolicyAttributeRequest) (*cbs.ModifyAutoSnapshotPolicyAttributeResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.ModifyAutoSnapshotPolicyAttribute)
	return resp, errors.WithStack(err)
}

func (c *EbsClient) DescribeAutoSnapshotPolicies(ctx context.Context, req *cbs.DescribeAutoSnapshotPoliciesRequest) (*cbs.DescribeAutoSnapshotPoliciesResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeAutoSnapshotPolicies)
	return resp, errors.WithStack(err)
}

func (c *EbsClient) BindAutoSnapshotPolicy(ctx context.Context, req *cbs.BindAutoSnapshotPolicyRequest) (*cbs.BindAutoSnapshotPolicyResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.BindAutoSnapshotPolicy)
	return resp, errors.WithStack(err)
}

func (c *EbsClient) UnbindAutoSnapshotPolicy(ctx context.Context, req *cbs.UnbindAutoSnapshotPolicyRequest) (*cbs.UnbindAutoSnapshotPolicyResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.UnbindAutoSnapshotPolicy)
	return resp, errors.WithStack(err)
}

func (c *EbsClient) DescribeDiskAssociatedAutoSnapshotPolicy(ctx context.Context, req *cbs.DescribeDiskAssociatedAutoSnapshotPolicyRequest) (*cbs.DescribeDiskAssociatedAutoSnapshotPolicyResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeDiskAssociatedAutoSnapshotPolicy)
	return resp, errors.WithStack(err)
}
