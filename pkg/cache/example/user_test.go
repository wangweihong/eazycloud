package example_test

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/wangweihong/eazycloud/pkg/cache/example"
	"github.com/wangweihong/eazycloud/pkg/sets"
)

func TestUserList(t *testing.T) {
	Convey("用户列表", t, func() {
		_, err := example.GetUMInstance().Add(&example.User{
			Name:       "aaa",
			UUID:       "user1",
			Tenant:     "tenant1",
			Group:      []string{"group1", "group2"},
			Roles:      []string{"role1", "role2"},
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		})
		So(err, ShouldBeNil)

		_, err2 := example.GetUMInstance().Add(&example.User{
			Name:       "aaa",
			UUID:       "user2",
			Tenant:     "tenant2",
			Group:      []string{"group1", "group3"},
			Roles:      []string{"role1", "role2"},
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		})
		So(err2, ShouldBeNil)

		Convey("查询某个组用户列表", func() {
			Convey("不存在的用户组", func() {
				So(len(example.GetUMInstance().ListInGroup("notexist")), ShouldEqual, 0)
			})
			Convey("存在的用户组", func() {
				g1 := example.GetUMInstance().ListInGroup("group1")
				So(len(g1), ShouldEqual, 2)
				for _, v := range g1 {
					So(sets.NewString(v.Group...).Has("group1"), ShouldBeTrue)
				}
			})
			Convey("ListInGroup应和ListInGroupIndex一致", func() {
				us := example.GetUMInstance().ListInGroup("group1")
				So(len(us), ShouldNotEqual, 0)

				usIndex := example.GetUMInstance().ListInGroupIndex("group1")
				So(len(usIndex), ShouldNotEqual, 0)
			})
		})
		Convey("查询某个租户用户列表键", func() {
			t1 := example.GetUMInstance().ListInTenantIndex("tenant1")
			So(len(t1), ShouldEqual, 1)
		})
	})

}

func TestUserManager_Add(t *testing.T) {
	Convey("用户添加", t, func() {
		u1 := &example.User{
			Name:       "aaa",
			Tenant:     "tenant1",
			Group:      []string{"group1", "group2"},
			Roles:      []string{"role1", "role2"},
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		nu1, err := example.GetUMInstance().Add(u1)
		So(err, ShouldBeNil)
		So(nu1.UUID, ShouldNotBeEmpty)

		Convey("入参数据改变不影响缓存数据", func() {
			u1.Name = "bbb"

			nu2, exists, err := example.GetUMInstance().Get(nu1)
			So(err, ShouldBeNil)
			So(exists, ShouldBeTrue)
			So(nu2.Name, ShouldNotEqual, u1.Name)
		})

		Convey("添加返回数据改变不影响缓存数据", func() {
			nu1.Name = "bbb"

			nu2, exists, err := example.GetUMInstance().Get(nu1)
			So(err, ShouldBeNil)
			So(exists, ShouldBeTrue)
			So(nu2.Name, ShouldNotEqual, nu1.Name)
		})

	})

}

func TestUserManager_CleanGroup(t *testing.T) {
	Convey("清除用户中某个组的索引", t, func() {
		u1 := &example.User{
			Name:       "aaa",
			Tenant:     "tenant1",
			Group:      []string{"group1", "group2"},
			Roles:      []string{"role1", "role2"},
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		nu1, err := example.GetUMInstance().Add(u1)
		So(err, ShouldBeNil)
		So(nu1.UUID, ShouldNotBeEmpty)

		Convey("清除组group1", func() {
			us := example.GetUMInstance().ListInGroup("group1")
			So(len(us), ShouldNotEqual, 0)

			usIndex := example.GetUMInstance().ListInGroupIndex("group1")
			So(len(usIndex), ShouldNotEqual, 0)

			So(len(us), ShouldEqual, len(usIndex))

			err = example.GetUMInstance().CleanGroup("group1")
			So(err, ShouldBeNil)

			us = example.GetUMInstance().ListInGroup("group1")
			So(len(us), ShouldEqual, 0)

		})

	})
}

func TestUserManager_CleanRole(t *testing.T) {
	Convey("清除用户中某个角色的索引", t, func() {
		u1 := &example.User{
			Name:       "aaa",
			Tenant:     "tenant1",
			Group:      []string{"group1", "group2"},
			Roles:      []string{"role1", "role2"},
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		nu1, err := example.GetUMInstance().Add(u1)
		So(err, ShouldBeNil)
		So(nu1.UUID, ShouldNotBeEmpty)

		Convey("清除角色role1", func() {
			us := example.GetUMInstance().ListInRole("role1")
			So(len(us), ShouldNotEqual, 0)

			usIndex := example.GetUMInstance().ListInRoleIndex("role1")
			So(len(usIndex), ShouldNotEqual, 0)

			So(len(us), ShouldEqual, len(usIndex))

			err = example.GetUMInstance().CleanRole("role1")
			So(err, ShouldBeNil)

			us = example.GetUMInstance().ListInRole("role1")
			So(len(us), ShouldEqual, 0)

		})
	})
}
