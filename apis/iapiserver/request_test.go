package iapiserver_test

import (
	"fmt"
	"testing"

	gvalidator "github.com/go-playground/validator/v10"
	. "github.com/smartystreets/goconvey/convey"
	v1 "k8s.io/api/core/v1"

	"github.com/wangweihong/eazycloud/apis/iapiserver"
	"github.com/wangweihong/eazycloud/pkg/validator"
)

func TestRequestValidateObjectGet(t *testing.T) {
	val := gvalidator.New()
	val.SetTagName("binding")
	val.RegisterValidation("namespaced", validator.ValidateNamespaceScopeResource)
	val.RegisterValidation("clusterd", validator.ValidateClusterScopeResource)

	type NamespaceScopeGetRequest struct {
		iapiserver.ResourceGetRequest `binding:"namespaced"`
	}

	type ClusterScopeGetRequest struct {
		iapiserver.ResourceGetRequest `binding:"clusterd"`
	}

	Convey("Test Resource Get tag Validator ", t, func() {
		Convey("Test NamespaceScope name ok", func() {
			r := &NamespaceScopeGetRequest{}
			So(val.Struct(r), ShouldNotBeNil)

			r.Name = "aa"
			So(val.Struct(r), ShouldNotBeNil)

			r.Namespace = "default"
			So(val.Struct(r), ShouldBeNil)
		})

		Convey("Test ClusterScopde name ok", func() {
			r := &ClusterScopeGetRequest{}
			So(val.Struct(r), ShouldNotBeNil)

			r.Name = "aa"
			So(val.Struct(r), ShouldBeNil)
		})
	})
}

func TestRequestValidateObjectParam(t *testing.T) {
	val := gvalidator.New(gvalidator.WithRequiredStructEnabled())
	val.SetTagName("binding")
	val.RegisterValidation("namespaced", validator.ValidateNamespaceScopeResource)
	val.RegisterValidation("clusterd", validator.ValidateClusterScopeResource)

	type ConfigMapRequest struct {
		iapiserver.ResourceRequest
		ConfigMap *v1.ConfigMap `json:"configmap" binding:"required,namespaced"`
	}
	type NodeRequest struct {
		iapiserver.ResourceRequest
		Node *v1.Node `json:"node" binding:"required,clusterd"`
	}

	Convey("Test Resource Object tag Validator ", t, func() {
		Convey("Test NamespaceScope name ok", func() {
			r := &ConfigMapRequest{}
			r.Cluster = "xxx"
			So(val.Struct(r), ShouldNotBeNil)

			r.ConfigMap = &v1.ConfigMap{}
			So(val.Struct(r), ShouldNotBeNil)

			r.ConfigMap.Namespace = "aa"
			So(val.Struct(r), ShouldNotBeNil)

			r.ConfigMap.Name = "aa"

			So(val.Struct(r), ShouldBeNil)
		})

		Convey("Test ClusterScopde name ok", func() {
			fmt.Println("-----------------")
			r := &NodeRequest{}
			So(val.Struct(r), ShouldNotBeNil)

			r.Node = &v1.Node{}
			So(val.Struct(r), ShouldNotBeNil)
			r.Cluster = "xx"
			r.Node.Name = "aa"
			So(val.Struct(r), ShouldBeNil)
		})
	})
}

func TestBatchRequestValidateObjectParam(t *testing.T) {
	val := gvalidator.New(gvalidator.WithRequiredStructEnabled())
	val.SetTagName("binding")
	val.RegisterValidation("namespaced", validator.ValidateNamespaceScopeResource)
	val.RegisterValidation("clusterd", validator.ValidateClusterScopeResource)

	type ConfigMapBatchRequest struct {
		Resources []*iapiserver.ConfigMapRequest `json:"resources" binding:"dive"` // dive设置递归
	}

	Convey("Test Batach Resource Object tag Validator ", t, func() {
		Convey("Test NamespaceScope name ok", func() {
			conditionOk := iapiserver.ConfigMapRequest{}
			conditionOk.Cluster = "xx"
			conditionOk.Resource = &v1.ConfigMap{}
			conditionOk.Resource.Name = "xx"
			conditionOk.Resource.Namespace = "yy"

			So(val.Struct(&ConfigMapBatchRequest{Resources: []*iapiserver.ConfigMapRequest{&conditionOk}}), ShouldBeNil)

			conditionFail := iapiserver.ConfigMapRequest{}
			conditionFail.Cluster = "xx"
			conditionFail.Resource = &v1.ConfigMap{}
			conditionFail.Resource.Name = "xx"

			So(
				val.Struct(&ConfigMapBatchRequest{Resources: []*iapiserver.ConfigMapRequest{&conditionFail}}),
				ShouldNotBeNil,
			)

			conditionFail2 := iapiserver.ConfigMapRequest{}
			conditionFail2.Cluster = "xx"
			conditionFail2.Resource = &v1.ConfigMap{}
			conditionFail2.Resource.Namespace = "yy"

			So(
				val.Struct(&ConfigMapBatchRequest{Resources: []*iapiserver.ConfigMapRequest{&conditionFail2}}),
				ShouldNotBeNil,
			)

		})
	})
}

func TestValidatePort(t *testing.T) {
	val := gvalidator.New(gvalidator.WithRequiredStructEnabled())
	val.SetTagName("binding")
	val.RegisterValidation("ports", validator.ValidatePorts)

	type PortRequest struct {
		Ports []int `json:"resources" binding:"ports"` // dive设置递归
	}

	Convey("Test ports Validator ", t, func() {
		conditionOk := PortRequest{}
		conditionOk.Ports = append(conditionOk.Ports, 2345, 355, 123)
		So(val.Struct(&conditionOk), ShouldBeNil)

		conditionFail := PortRequest{}
		conditionFail.Ports = append(conditionFail.Ports, -1, 653355)
		So(val.Struct(&conditionFail), ShouldNotBeNil)

	})
}
