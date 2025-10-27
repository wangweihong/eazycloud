package libkubernetes

import (
	"context"

	storagev1 "k8s.io/api/storage/v1"
	v1 "k8s.io/api/storage/v1"

	"github.com/wangweihong/eazycloud/apis/ikubernetes"

	"github.com/kubernetes-csi/external-snapshotter/client/v3/apis/volumesnapshot/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func StorageClassList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*storagev1.StorageClassList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*storagev1.StorageClassList, error) {
		resource, err := c.client.StorageV1().StorageClasses().List(ctx, opts)
		return resource, err
	})
	return resp, err
}

func StorageClassDelete(pctx context.Context, config *ikubernetes.ClusterConfig, name string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*storagev1.StorageClass, error) {
		err := c.client.StorageV1().StorageClasses().Delete(ctx, name, opts)
		return nil, err
	})
	return err
}

func StorageClassCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	req *storagev1.StorageClass,
	opts metav1.CreateOptions,
) (*storagev1.StorageClass, error) {

	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*storagev1.StorageClass, error) {
		resource, err := c.client.StorageV1().StorageClasses().Create(ctx, req, opts)
		return resource, err
	})
	return resp, err
}

func StorageClassUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	req *storagev1.StorageClass,
	opts metav1.UpdateOptions,
) (*storagev1.StorageClass, error) {

	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*storagev1.StorageClass, error) {
		resource, err := c.client.StorageV1().StorageClasses().Update(ctx, req, opts)
		return resource, err
	})
	return resp, err
}

func StorageClassGet(pctx context.Context, config *ikubernetes.ClusterConfig, name string, opts metav1.GetOptions) (*v1.StorageClass, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*storagev1.StorageClass, error) {
		resource, err := c.client.StorageV1().StorageClasses().Get(ctx, name, opts)
		return resource, err
	})
	return resp, err

}

// volume snapshot
// volume snapshot content
// volume snapshot class
func VolumeSnapshotList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*v1beta1.VolumeSnapshotList, error) {

	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1beta1.VolumeSnapshotList, error) {
		resource, err := c.snapshotClient.VolumeSnapshots(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err
}

func VolumeSnapshotGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*v1beta1.VolumeSnapshot, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1beta1.VolumeSnapshot, error) {
		resource, err := c.snapshotClient.VolumeSnapshots(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func VolumeSnapshotDelete(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.DeleteOptions,
) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1beta1.VolumeSnapshot, error) {
		err := c.snapshotClient.VolumeSnapshots(namespace).Delete(ctx, name, opts)
		return nil, err
	})
	return err
}

func VolumeSnapshotCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	req *v1beta1.VolumeSnapshot,
	opts metav1.CreateOptions,
) (*v1beta1.VolumeSnapshot, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1beta1.VolumeSnapshot, error) {
		resource, err := c.snapshotClient.VolumeSnapshots(namespace).Create(ctx, req, opts)
		return resource, err
	})
	return resp, err
}

func VolumeSnapshotUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	req *v1beta1.VolumeSnapshot,
	opts metav1.UpdateOptions,
) (*v1beta1.VolumeSnapshot, error) {

	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1beta1.VolumeSnapshot, error) {
		resource, err := c.snapshotClient.VolumeSnapshots(namespace).Update(ctx, req, opts)
		return resource, err
	})
	return resp, err

}

func VolumeSnapshotContentList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*v1beta1.VolumeSnapshotContentList, error) {

	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1beta1.VolumeSnapshotContentList, error) {
		resource, err := c.snapshotClient.VolumeSnapshotContents().List(ctx, opts)
		return resource, err
	})
	return resp, err
}

func VolumeSnapshotContentGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*v1beta1.VolumeSnapshotContent, error) {

	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1beta1.VolumeSnapshotContent, error) {
		resource, err := c.snapshotClient.VolumeSnapshotContents().Get(ctx, name, opts)
		return resource, err
	})
	return resp, err

}

func VolumeSnapshotContentDelete(pctx context.Context, config *ikubernetes.ClusterConfig, name string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1beta1.VolumeSnapshotContent, error) {
		err := c.snapshotClient.VolumeSnapshotContents().Delete(ctx, name, opts)
		return nil, err
	})
	return err
}

func VolumeSnapshotContentCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	req *v1beta1.VolumeSnapshotContent,
	opts metav1.CreateOptions,
) (*v1beta1.VolumeSnapshotContent, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1beta1.VolumeSnapshotContent, error) {
		resource, err := c.snapshotClient.VolumeSnapshotContents().Create(ctx, req, opts)
		return resource, err
	})
	return resp, err
}

func VolumeSnapshotContentUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	req *v1beta1.VolumeSnapshotContent,
	opts metav1.UpdateOptions,
) (*v1beta1.VolumeSnapshotContent, error) {

	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1beta1.VolumeSnapshotContent, error) {
		resource, err := c.snapshotClient.VolumeSnapshotContents().Update(ctx, req, opts)
		return resource, err
	})
	return resp, err

}

func VolumeSnapshotClassList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*v1beta1.VolumeSnapshotClassList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1beta1.VolumeSnapshotClassList, error) {
		resource, err := c.snapshotClient.VolumeSnapshotClasses().List(ctx, opts)
		return resource, err
	})
	return resp, err

}

func VolumeSnapshotClassGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*v1beta1.VolumeSnapshotClass, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1beta1.VolumeSnapshotClass, error) {
		resource, err := c.snapshotClient.VolumeSnapshotClasses().Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func VolumeSnapshotClassDelete(pctx context.Context, config *ikubernetes.ClusterConfig, name string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1beta1.VolumeSnapshotClass, error) {
		err := c.snapshotClient.VolumeSnapshotClasses().Delete(ctx, name, opts)
		return nil, err
	})
	return err
}

func VolumeSnapshotClassCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	req *v1beta1.VolumeSnapshotClass,
	opts metav1.CreateOptions,
) (*v1beta1.VolumeSnapshotClass, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1beta1.VolumeSnapshotClass, error) {
		resource, err := c.snapshotClient.VolumeSnapshotClasses().Create(ctx, req, opts)
		return resource, err
	})
	return resp, err

}

func VolumeSnapshotClassUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	req *v1beta1.VolumeSnapshotClass,
	opts metav1.UpdateOptions,
) (*v1beta1.VolumeSnapshotClass, error) {

	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1beta1.VolumeSnapshotClass, error) {
		resource, err := c.snapshotClient.VolumeSnapshotClasses().Update(ctx, req, opts)
		return resource, err
	})
	return resp, err

}
