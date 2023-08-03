package stringutil_test

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/wangweihong/eazycloud/pkg/util/stringutil"

	"testing"
)

func TestBothEmptyOrNone(t *testing.T) {
	Convey("BothEmptyOrNone", t, func() {
		So(stringutil.BothEmptyOrNone("a", ""), ShouldBeFalse)
		So(stringutil.BothEmptyOrNone("a", "b"), ShouldBeTrue)
		So(stringutil.BothEmptyOrNone("", "b"), ShouldBeFalse)
	})
}

func TestF(t *testing.T) {
	now := time.Now()
	fmt.Printf("%02d%02d%02d%02d%02d", now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

func TestEncrypt(t *testing.T) {
	userID := "J10003"
	password := "111111"
	timestamp := "0803192020"
	fixedValue := "00000000"

	data := userID + fixedValue + password + timestamp
	// 创建 MD5 加密器
	hasher := md5.New()

	// 将数据写入加密器
	_, err := hasher.Write([]byte(data))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 计算 MD5 值
	hashedData := hasher.Sum(nil)

	// 将二进制的 MD5 值转换为十六进制字符串
	md5String := hex.EncodeToString(hashedData)

	fmt.Println(md5String)
}
