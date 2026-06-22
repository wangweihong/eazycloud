package ctyun_test

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/wangweihong/eazycloud/apis/ictyun"
	"github.com/wangweihong/gotoolbox/pkg/log"
	"github.com/wangweihong/gotoolbox/pkg/paging"
	"github.com/wangweihong/gotoolbox/pkg/typeutil"
	"github.com/wangweihong/gotoolbox/pkg/waitgroup"
)

func TestEcs_EcsList(t *testing.T) {
	SkipConvey("EcsList", t, func() {

		_, err := tc.Ecss().EcsList(context.Background(), &ictyun.EcsListRequest{
			RegionID: region,
		}, nil)
		So(err, ShouldBeNil)
	})

	SkipConvey("AllEcsList", t, func() {
		rgResp := &ictyun.RegionListResponse{}
		_, err := tc.Ecss().RegionList(context.Background(), &ictyun.RegionListRequest{}, rgResp)
		So(err, ShouldBeNil)

		start := time.Now()
		wg := waitgroup.RunConcurrentlyCondition(context.Background(), rgResp.ReturnObj.RegionList, func(r ictyun.RegionListEntry) bool {
			return r.RegionParent != "港澳及海外"
		}, func(ctx context.Context, r ictyun.RegionListEntry) waitgroup.Result {
			rid := r.RegionID
			var pageNum int = 1
			var ecs []ictyun.EcsInfo
			for {
				ecsResp := &ictyun.EcsListResponse{}
				_, err := tc.Ecss().EcsList(context.Background(), &ictyun.EcsListRequest{
					RegionID: rid,
					PagingParam: ictyun.PagingParam{
						PageNum:  &pageNum,
						PageSize: typeutil.Int(50),
					},
				}, ecsResp)
				if err != nil {
					log.Error(err.Error())
					break
				}

				ecs = append(ecs, ecsResp.ReturnObj.Results...)
				if ecsResp.ReturnObj.TotalPage == pageNum {
					break
				}
				pageNum++
			}
			return waitgroup.NewResult(ecs, nil)
		})

		wg.Debug().PrintResults()

		var allvms []ictyun.EcsInfo
		for _, v := range wg.GetResults() {
			if v.Data != nil {
				vms := v.Data.([]ictyun.EcsInfo)
				allvms = append(allvms, vms...)
			}
		}

		s, e := paging.Index(len(allvms), 0, 10)
		pageVms := allvms[s:e]
		log.Infof("-----------------------%v", len(pageVms))
		log.Infof("cost:%v", time.Now().Sub(start).Seconds())

	})
}

func TestEcs_EcsCreate(t *testing.T) {
	SkipConvey("EcsCreate", t, func() {
		_, err := tc.Ecss().EcsCreate(context.Background(), &ictyun.EcsCreateRequest{
			RegionID: region,
		}, nil)
		So(err, ShouldBeNil)
	})
}

func TestEcs_EcsFlavorList(t *testing.T) {
	Convey("EcsFlavorList", t, func() {
		Convey("查询全部", func() {
			_, err := tc.Ecss().EcsFlavorList(context.Background(), &ictyun.EcsFlavorListRequest{
				RegionID: region,
			}, nil)
			So(err, ShouldBeNil)
		})
		Convey("查询通用系列", func() {

			resp := &ictyun.EcsFlavorListResponse{}
			_, err := tc.Ecss().EcsFlavorList(context.Background(), &ictyun.EcsFlavorListRequest{
				RegionID:     region,
				FlavorSeries: typeutil.String("s"),
			}, resp)
			So(err, ShouldBeNil)
			for _, v := range resp.ReturnObj.FlavorList {
				So(v.FlavorSeries, ShouldEqual, "s")
			}

		})
		Convey("查询内存为8G的规格", func() {
			resp := &ictyun.EcsFlavorListResponse{}
			_, err := tc.Ecss().EcsFlavorList(context.Background(), &ictyun.EcsFlavorListRequest{
				RegionID:  region,
				FlavorRAM: typeutil.Int(8),
			}, resp)
			So(err, ShouldBeNil)
			for _, v := range resp.ReturnObj.FlavorList {
				So(v.FlavorRAM, ShouldEqual, 8)
			}
		})
	})
}
