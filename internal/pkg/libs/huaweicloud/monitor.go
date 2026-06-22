package huaweicloud

import (
	"context"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1/model"
	"github.com/wangweihong/eazycloud/internal/pkg/libs"
	"github.com/wangweihong/gotoolbox/pkg/errors"

	ces "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ces/v1"
)

type Monitor struct {
	c   *ces.CesClient
	err error
}

func NewMonitor(cfg *AccessConfig) *Monitor {
	c, err := cfg.HcCesClient(cfg.Region)
	return &Monitor{c: c, err: errors.WithStack(err)}
}

// MonitorList 查询指标列表
// https://console.huaweicloud.com/apiexplorer/#/openapi/CES/sdk?version=v1&api=ListMetrics
func (p *Monitor) MonitorListMetrics(
	ctx context.Context,
	req *model.ListMetricsRequest,
) (*model.ListMetricsResponse, error) {
	resp, err := libs.InvokeApi(ctx, p.err, req, p.c.ListMetrics)
	return resp, errors.WithStack(err)
}
