package tecentcloud

import (
	"context"

	tcommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

const (
	//虚拟机实例状态
	InstanceStatePENDING      = "PENDING"       //表示创建中
	InstanceStateLAUNCHFAILED = "LAUNCH_FAILED" //表示创建失败
	InstanceStateRUNNING      = "RUNNING"       //表示运行中
	InstanceStateSTOPPED      = "STOPPED"       //表示关机
	InstanceStateSTARTING     = "STARTING"      //表示开机中
	InstanceStateSTOPPING     = "STOPPING"      //表示关机中
	InstanceStateREBOOTING    = "REBOOTING"     //表示重启中
	InstanceStateSHUTDOWN     = "SHUTDOWN"      //表示停止待销毁
	InstanceStateTERMINATING  = "TERMINATING"   //表示销毁中
	//镜像类型
	ImageTypePUBLICIMAGE  = "PUBLIC_IMAGE"  //表示公共镜像
	ImageTypePRIVATEIMAGE = "PRIVATE_IMAGE" //表示私有镜像
	//镜像状态
	ImageStateCREATING     = "CREATING"     //表示创建中
	ImageStateNORMAL       = "NORMAL"       //表示正常
	ImageStateCREATEFAILED = "CREATEFAILED" //表示创建失败
	ImageStateUSING        = "USING"        //表示使用中
	ImageStateSYNCING      = "SYNCING"      //表示同步中
	ImageStateIMPORTING    = "IMPORTING"    //表示导入中
	ImageStateIMPORTFAILED = "IMPORTFAILED" //表示导入失败
)

type EcsClient struct {
	c   *cvm.Client
	err error
}

func NewEcsClient(securityId, securityKey string, region string) *EcsClient {
	credential := tcommon.NewCredential(securityId, securityKey)
	clientProfile := profile.NewClientProfile()
	//clientProfile.HttpProfile.Endpoint = endpoint
	client, err := cvm.NewClient(credential, region, clientProfile)
	return &EcsClient{
		c:   client,
		err: err,
	}
}

func (c *EcsClient) DescribeRegions(ctx context.Context, req *cvm.DescribeRegionsRequest) (*cvm.DescribeRegionsResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeRegions)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) DescribeZones(ctx context.Context, req *cvm.DescribeZonesRequest) (*cvm.DescribeZonesResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeZones)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) DescribeInstances(ctx context.Context, req *cvm.DescribeInstancesRequest) (*cvm.DescribeInstancesResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeInstances, libs.WithRateLimit(5, 10))
	return resp, errors.WithStack(err)
}

func (c *EcsClient) DescribeImages(ctx context.Context, req *cvm.DescribeImagesRequest) (*cvm.DescribeImagesResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeImages)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) CreateImage(ctx context.Context, req *cvm.CreateImageRequest) (*cvm.CreateImageResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.CreateImage)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) ImportImage(ctx context.Context, req *cvm.ImportImageRequest) (*cvm.ImportImageResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.ImportImage)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) DeleteImages(ctx context.Context, req *cvm.DeleteImagesRequest) (*cvm.DeleteImagesResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DeleteImages)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) ModifyImageAttribute(ctx context.Context, req *cvm.ModifyImageAttributeRequest) (*cvm.ModifyImageAttributeResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.ModifyImageAttribute)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) DescribeImportImageOs(ctx context.Context, req *cvm.DescribeImportImageOsRequest) (*cvm.DescribeImportImageOsResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeImportImageOs)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) DescribeImageQuota(ctx context.Context, req *cvm.DescribeImageQuotaRequest) (*cvm.DescribeImageQuotaResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeImageQuota)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) RunInstances(ctx context.Context, req *cvm.RunInstancesRequest) (*cvm.RunInstancesResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.RunInstances)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) TerminateInstances(ctx context.Context, req *cvm.TerminateInstancesRequest) (*cvm.TerminateInstancesResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.TerminateInstances)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) ModifyInstancesAttribute(ctx context.Context, req *cvm.ModifyInstancesAttributeRequest) (*cvm.ModifyInstancesAttributeResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.ModifyInstancesAttribute)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) ResetInstance(ctx context.Context, req *cvm.ResetInstanceRequest) (*cvm.ResetInstanceResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.ResetInstance)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) ResizeInstanceDisks(ctx context.Context, req *cvm.ResizeInstanceDisksRequest) (*cvm.ResizeInstanceDisksResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.ResizeInstanceDisks)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) ResetInstancesPassword(ctx context.Context, req *cvm.ResetInstancesPasswordRequest) (*cvm.ResetInstancesPasswordResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.ResetInstancesPassword)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) DescribeInstancesStatus(ctx context.Context, req *cvm.DescribeInstancesStatusRequest) (*cvm.DescribeInstancesStatusResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeInstancesStatus)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) DescribeInstanceVncUrl(ctx context.Context, req *cvm.DescribeInstanceVncUrlRequest) (*cvm.DescribeInstanceVncUrlResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeInstanceVncUrl)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) StartInstances(ctx context.Context, req *cvm.StartInstancesRequest) (*cvm.StartInstancesResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.StartInstances)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) StopInstances(ctx context.Context, req *cvm.StopInstancesRequest) (*cvm.StopInstancesResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.StopInstances)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) RebootInstances(ctx context.Context, req *cvm.RebootInstancesRequest) (*cvm.RebootInstancesResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.RebootInstances)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) DescribeInstanceFamilyConfigs(ctx context.Context, req *cvm.DescribeInstanceFamilyConfigsRequest) (*cvm.DescribeInstanceFamilyConfigsResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeInstanceFamilyConfigs)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) DescribeZoneInstanceConfigInfos(ctx context.Context, req *cvm.DescribeZoneInstanceConfigInfosRequest) (*cvm.DescribeZoneInstanceConfigInfosResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeZoneInstanceConfigInfos)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) DescribeInstanceTypeConfigs(ctx context.Context, req *cvm.DescribeInstanceTypeConfigsRequest) (*cvm.DescribeInstanceTypeConfigsResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DescribeInstanceTypeConfigs)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) ResetInstancesInternetMaxBandwidth(ctx context.Context, req *cvm.ResetInstancesInternetMaxBandwidthRequest) (*cvm.ResetInstancesInternetMaxBandwidthResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.ResetInstancesInternetMaxBandwidth)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) DisassociateSecurityGroups(ctx context.Context, req *cvm.DisassociateSecurityGroupsRequest) (*cvm.DisassociateSecurityGroupsResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.DisassociateSecurityGroups)
	return resp, errors.WithStack(err)
}

func (c *EcsClient) AssociateSecurityGroups(ctx context.Context, req *cvm.AssociateSecurityGroupsRequest) (*cvm.AssociateSecurityGroupsResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.AssociateSecurityGroups)
	return resp, errors.WithStack(err)
}
