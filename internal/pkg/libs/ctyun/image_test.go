package ctyun_test

//
//func TestImage_ImageList(t *testing.T) {
//	Convey("ImageList", t, func() {
//		c, err := ctyun.NewClient(endpoint, ak, sk)
//		So(err, ShouldBeNil)
//
//		Convey("私有镜像", func() {
//			resp, err := ctyun.NewImage(c).ImageList(context.Background(), &ictyun.ImageListRequest{
//				RegionID: region,
//			})
//			So(err, ShouldBeNil)
//			log.Infof("len:%v", len(resp.ReturnObj.Images))
//		})
//
//		Convey("公共镜像", func() {
//			resp, err := ctyun.NewImage(c).ImageList(context.Background(), &ictyun.ImageListRequest{
//				RegionID:   region,
//				Visibility: 1,
//			})
//			So(err, ShouldBeNil)
//			log.Infof("len:%v", len(resp.ReturnObj.Images))
//		})
//
//	})
//}
//
//func TestImage_ImageGet(t *testing.T) {
//	Convey("ImageGet", t, func() {
//		c, err := ctyun.NewClient(endpoint, ak, sk)
//		So(err, ShouldBeNil)
//
//		Convey("公共镜像", func() {
//			resp, err := ctyun.NewImage(c).ImageGet(context.Background(), &ictyun.ImageGetRequest{
//				RegionID: region,
//				ImageID:  "ad0fc415-9455-4eea-b11c-5006eae1f9a0",
//			})
//			So(err, ShouldBeNil)
//			log.Infof("len:%v", len(resp.ReturnObj.Images))
//		})
//	})
//}
