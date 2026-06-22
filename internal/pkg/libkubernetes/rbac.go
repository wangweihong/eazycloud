package libkubernetes

import (
	"context"

	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/wangweihong/eazycloud/apis/ikubernetes"
)

func RoleCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	req *v1.Role,
	opts metav1.CreateOptions,
) (*v1.Role, error) {

	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Role, error) {
		resource, err := c.client.RbacV1().Roles(namespace).Create(ctx, req, opts)
		return resource, err
	})
	return resp, err
}

func RoleGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*v1.Role, error) {

	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Role, error) {
		resource, err := c.client.RbacV1().Roles(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err
}

func RoleDelete(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, object string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Role, error) {
		err := c.client.RbacV1().Roles(namespace).Delete(ctx, object, opts)
		return nil, err
	})
	return err
}

func RoleList(pctx context.Context, config *ikubernetes.ClusterConfig, namespace string, opts metav1.ListOptions) (*v1.RoleList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.RoleList, error) {
		resource, err := c.client.RbacV1().Roles(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err

}

func RoleUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	object *v1.Role,
	opts metav1.UpdateOptions,
) (*v1.Role, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.Role, error) {
		resource, err := c.client.RbacV1().Roles(namespace).Update(ctx, object, opts)
		return resource, err
	})
	return resp, err
}

func RoleBindingGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	name string,
	opts metav1.GetOptions,
) (*v1.RoleBinding, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.RoleBinding, error) {
		resource, err := c.client.RbacV1().RoleBindings(namespace).Get(ctx, name, opts)
		return resource, err
	})
	return resp, err

}

func RoleBindingCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	object *v1.RoleBinding,
	opts metav1.CreateOptions,
) (*v1.RoleBinding, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.RoleBinding, error) {
		resource, err := c.client.RbacV1().RoleBindings(namespace).Create(ctx, object, opts)
		return resource, err
	})
	return resp, err
}

func RoleBindingDelete(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	object string,
	opts metav1.DeleteOptions,
) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.RoleBinding, error) {
		err := c.client.RbacV1().RoleBindings(namespace).Delete(ctx, object, opts)
		return nil, err
	})
	return err
}

func RoleBindingList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	opts metav1.ListOptions,
) (*v1.RoleBindingList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.RoleBindingList, error) {
		resource, err := c.client.RbacV1().RoleBindings(namespace).List(ctx, opts)
		return resource, err
	})
	return resp, err

}

func RoleBindingUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	namespace string,
	object *v1.RoleBinding,
	opts metav1.UpdateOptions,
) (*v1.RoleBinding, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.RoleBinding, error) {
		resource, err := c.client.RbacV1().RoleBindings(namespace).Update(ctx, object, opts)
		return resource, err
	})
	return resp, err
}

func ClusterRoleCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	object *v1.ClusterRole,
	opts metav1.CreateOptions,
) (*v1.ClusterRole, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ClusterRole, error) {
		resource, err := c.client.RbacV1().ClusterRoles().Create(ctx, object, opts)
		return resource, err
	})
	return resp, err
}

func ClusterRoleDelete(pctx context.Context, config *ikubernetes.ClusterConfig, object string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ClusterRole, error) {
		err := c.client.RbacV1().ClusterRoles().Delete(ctx, object, opts)
		return nil, err
	})
	return err
}

func ClusterRoleGet(pctx context.Context, config *ikubernetes.ClusterConfig, object string, opts metav1.GetOptions) (*v1.ClusterRole, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ClusterRole, error) {
		resource, err := c.client.RbacV1().ClusterRoles().Get(ctx, object, opts)
		return resource, err
	})
	return resp, err
}

func ClusterRoleList(pctx context.Context, config *ikubernetes.ClusterConfig, opts metav1.ListOptions) (*v1.ClusterRoleList, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ClusterRoleList, error) {
		resource, err := c.client.RbacV1().ClusterRoles().List(ctx, opts)
		return resource, err
	})
	return resp, err
}

func ClusterRoleUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	object *v1.ClusterRole,
	opts metav1.UpdateOptions,
) (*v1.ClusterRole, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ClusterRole, error) {
		resource, err := c.client.RbacV1().ClusterRoles().Update(ctx, object, opts)
		return resource, err
	})
	return resp, err
}

func ClusterRoleBindingCreate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	object *v1.ClusterRoleBinding,
	opts metav1.CreateOptions,
) (*v1.ClusterRoleBinding, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ClusterRoleBinding, error) {
		resource, err := c.client.RbacV1().ClusterRoleBindings().Create(ctx, object, opts)
		return resource, err
	})
	return resp, err
}

func ClusterRoleBindingDelete(pctx context.Context, config *ikubernetes.ClusterConfig, object string, opts metav1.DeleteOptions) error {
	_, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ClusterRole, error) {
		err := c.client.RbacV1().ClusterRoleBindings().Delete(ctx, object, opts)
		return nil, err
	})
	return err
}

func ClusterRoleBindingGet(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	object string,
	opts metav1.GetOptions,
) (*v1.ClusterRoleBinding, error) {

	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ClusterRoleBinding, error) {
		resource, err := c.client.RbacV1().ClusterRoleBindings().Get(ctx, object, opts)
		return resource, err
	})
	return resp, err
}

func ClusterRoleBindingList(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	opts metav1.ListOptions,
) (*v1.ClusterRoleBindingList, error) {

	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ClusterRoleBindingList, error) {
		resource, err := c.client.RbacV1().ClusterRoleBindings().List(ctx, opts)
		return resource, err
	})
	return resp, err

}

func ClusterRoleBindingUpdate(
	pctx context.Context,
	config *ikubernetes.ClusterConfig,
	object *v1.ClusterRoleBinding,
	opts metav1.UpdateOptions,
) (*v1.ClusterRoleBinding, error) {
	resp, err := run(pctx, config, func(ctx context.Context, c *Client) (*v1.ClusterRoleBinding, error) {
		resource, err := c.client.RbacV1().ClusterRoleBindings().Update(ctx, object, opts)
		return resource, err
	})
	return resp, err
}
