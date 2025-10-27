package libkubernetes

import (
	"context"

	"github.com/wangweihong/gotoolbox/pkg/errors"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"

	"github.com/wangweihong/eazycloud/apis/ikubernetes"
)

func IsServerReachable(config *ikubernetes.ClusterConfig) error {
	c, err := NewClient(config)
	if err != nil {
		return err
	}
	if _, err = c.client.ServerVersion(); err != nil {
		return errors.Errorf("cluster not reachable")
	}

	return nil
}

func ToDiscoveryClient(config *ikubernetes.ClusterConfig) (*discovery.DiscoveryClient, error) {
	c, err := NewClient(config)
	if err != nil {
		return nil, err
	}
	return c.client.DiscoveryClient, nil
}

func ServerVersion(ctx context.Context, config *ikubernetes.ClusterConfig) (*version.Info, error) {
	c, err := NewClient(config)
	if err != nil {
		return nil, err
	}

	versionInfo, err := c.client.ServerVersion()
	if err != nil {
		return nil, err
	}

	return versionInfo, nil
}
