package ctyun_test

//
//import (
//	"context"
//	"testing"
//	"time"
//
//	. "github.com/smartystreets/goconvey/convey"
//	"github.com/wangweihong/eazycloud/internal/pkg/multicloud/ctyun"
//	"github.com/wangweihong/eazycloud/internal/pkg/multicloud/ctyun/ictyun"
//	"github.com/wangweihong/gotoolbox/pkg/json"
//	"github.com/wangweihong/gotoolbox/log"
//	"github.com/wangweihong/gotoolbox/paging"
//	"github.com/wangweihong/gotoolbox/pkg/typeutil"
//	"github.com/wangweihong/gotoolbox/waitgroup"
//)
//
//func TestEbs_EbsList(t *testing.T) {
//	SkipConvey("EbsList", t, func() {
//		c, err := ctyun.NewClient(endpoint, ak, sk)
//		So(err, ShouldBeNil)
//		resp, err := ctyun.NewEbs(c).EbsList(context.Background(), &ictyun.EbsListRequest{
//			RegionID: region,
//		})
//		So(err, ShouldBeNil)
//		json.PrintStructObject(resp)
//	})
//
//	SkipConvey("查询挂载某云实例的主机", t, func() {
//		c, err := ctyun.NewClient(endpoint, ak, sk)
//		So(err, ShouldBeNil)
//		var allDisks []ictyun.EBSDiskInfo
//		var i = 1
//		var size = 50
//		for {
//			ret, err := ctyun.NewEbs(c).EbsList(context.Background(), &ictyun.EbsListRequest{
//				RegionID:    region,
//				PagingParam: ictyun.PagingParam{PageNum: &i, PageSize: &size},
//			})
//			if err != nil {
//				break
//			}
//			allDisks = append(allDisks, ret.ReturnObj.DiskList...)
//		}
//
//		So(err, ShouldBeNil)
//		json.PrintStructObject(allDisks)
//	})
//
//	SkipConvey("查找全部资源池的全部云硬盘", t, func() {
//
//		c, err := ctyun.NewClient(endpoint, ak, sk)
//		So(err, ShouldBeNil)
//		rgResp, err := ctyun.NewRegion(c).RegionList(context.Background(), &ictyun.RegionListRequest{})
//		So(err, ShouldBeNil)
//
//		start := time.Now()
//		wg := waitgroup.NewWaitGroup(context.Background())
//
//		for _, r := range rgResp.ReturnObj.RegionList {
//			rid := r.RegionID
//			// 港澳及海外耗时非常长
//			if r.RegionParent != "港澳及海外" {
//				wg.Start(waitgroup.NewWaitGroupHandleFunc(context.Background(), r.RegionID+"/"+r.RegionName, func() waitgroup.WaitGroupResult {
//					var pageNum int = 1
//					var ecs []ictyun.EBSDiskInfo
//					for {
//						ecsResp, err := ctyun.NewEbs(c).EbsList(context.Background(), &ictyun.EbsListRequest{
//							RegionID: rid,
//							PagingParam: ictyun.PagingParam{
//								PageNum:  &pageNum,
//								PageSize: typeutil.Int(50),
//							},
//						})
//						if err != nil {
//							log.Error(err.Error())
//							break
//						}
//
//						ecs = append(ecs, ecsResp.ReturnObj.DiskList...)
//						if ecsResp.ReturnObj.TotalPage == pageNum {
//							break
//						}
//						pageNum++
//					}
//					return waitgroup.NewWaitGroupResult(ecs, nil)
//				}))
//			}
//		}
//		wg.Wait()
//		wg.Debug().PrintResults()
//
//		var allvms []ictyun.EBSDiskInfo
//		for _, v := range wg.GetResults() {
//			if v.Data != nil {
//				vms := v.Data.([]ictyun.EBSDiskInfo)
//				allvms = append(allvms, vms...)
//			}
//		}
//
//		s, e := paging.Index(len(allvms), 0, 10)
//		pageVms := allvms[s:e]
//		log.Infof("-----------------------%v", len(pageVms))
//		log.Infof("cost:%v", time.Now().Sub(start).Seconds())
//
//	})
//}
