package libprometheus

import (
	"time"

	"github.com/prometheus/client_golang/api"
	apiv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/wangweihong/eazycloud/apis/iprometheus"
	"github.com/wangweihong/gotoolbox/pkg/errors"
)

const (
	defaultTimeout    = 30
	DefaultTimeoutGet = 10
)

type PrometheusClient struct {
	client apiv1.API
}

func (m *PrometheusClient) Close() {}

func newPrometheusClient(conf *iprometheus.Config) (*PrometheusClient, error) {
	if conf == nil {
		return nil, errors.Errorf("prometheus config is nil")
	}

	cfg := api.Config{
		Address: conf.Address,
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &PrometheusClient{client: apiv1.NewAPI(client)}, nil
}

func getTimeOut(i ...int) time.Duration {
	if len(i) <= 0 {
		return 5
	}
	if i[0] < 1 {
		return time.Duration(5)
	}

	return time.Duration(i[0])
}
