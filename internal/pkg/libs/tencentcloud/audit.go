package tecentcloud

import (
	"context"

	audit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"
	tcommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

type AuditClient struct {
	c   *audit.Client
	err error
}

func NewAuditClient(securityId, securityKey string, region string) *AuditClient {
	credential := tcommon.NewCredential(securityId, securityKey)
	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = HttpProfileUserEndpoint
	client, err := audit.NewClient(credential, region, clientProfile)
	return &AuditClient{
		c:   client,
		err: err,
	}
}

func (c *AuditClient) LookUpEvents(ctx context.Context, clusterUUID string, region string, req *audit.LookUpEventsRequest) (*audit.LookUpEventsResponse, error) {
	resp, err := invoke(ctx, c.err, req, c.c.LookUpEvents)
	return resp, errors.WithStack(err)
}
