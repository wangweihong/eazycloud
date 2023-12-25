package sortutil_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/wangweihong/eazycloud/pkg/util/sortutil"
	"sort"
	"testing"
)

func TestSort_Pinyin(t *testing.T) {
	Convey("拼音排序", t, func() {
		data := []string{"上海", "北京", "广州"}
		sort.Sort(sortutil.ByPinyin(data))
		So(data, ShouldResemble, []string{"北京", "广州", "上海"})
	})
}
