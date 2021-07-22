package oidb_0x485

import (
	"context"
	"fmt"
	"testing"
	"time"

	"git.code.oa.com/qdgo/core/net/udp"
	"git.code.oa.com/qdgo/core/store/oidb/binv5"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	cfg = &binv5.Config{
		ServiceType: 16,
		UDPClient: &udp.ClientConfig{
			RemoteAddr:   "l5://737473:917504",
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
		},
	}
	uin     = uint64(623009470)
	appid   = uint32(717048201) // APPID
	keyType = uint32(27)        // 登录态类型
	sKey    = []byte("2rdDmIOwj3K5Q0IyWeopSlqJX14NHfp5joGAFbpBXrk_")
)

func TestOidb0x485(t *testing.T) {
	Convey("test 0x485", t, func() {
		o, err := New(cfg)
		So(err, ShouldBeNil)
		defer o.Close()
		o.SetSig(uin, appid, keyType, sKey)
		rsp, err := o.Do(context.Background())
		So(err, ShouldBeNil)
		fmt.Printf("rsp(%+v)", rsp)
	})
}
