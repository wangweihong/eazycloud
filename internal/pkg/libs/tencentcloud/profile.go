package tecentcloud

import (
	"context"

	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	tcommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

const (
	HttpProfileUserEndpoint = "cam.tencentcloudapi.com"
)

type CamClient struct {
	c   *cam.Client
	err error
}

func NewCamClient(securityId, securityKey string, region string) *CamClient {
	credential := tcommon.NewCredential(securityId, securityKey)
	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = HttpProfileUserEndpoint
	client, err := cam.NewClient(credential, region, clientProfile)
	return &CamClient{
		c:   client,
		err: err,
	}
}

func (c *CamClient) ListUsers(ctx context.Context, clusterUUID string, req *cam.ListUsersRequest) (*cam.ListUsersResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.ListUsers)
	return resp, errors.WithStack(err)
}

func (c *CamClient) ListGroups(ctx context.Context, clusterUUID string, req *cam.ListGroupsRequest) (*cam.ListGroupsResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.ListGroups)
	return resp, errors.WithStack(err)
}
