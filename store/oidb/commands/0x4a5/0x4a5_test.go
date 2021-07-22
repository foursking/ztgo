package oidb_0x4a5

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/foursking/ztgo/net/udp"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	cfg = &Config{
		ServiceType: 1,
		UDPClient: &udp.ClientConfig{
			RemoteAddr:   "l5://1002241:65536",
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
		},
	}
	email   = "17621627786@qidian.qq.com"
	uin     = uint64(0)
	appid   = uint32(3000401) // 联合登录态的APPID
	keyType = uint32(36)      // 联合登录态类型
	sKey    = []byte("@KE9elzKOA")
)

func TestOidb0x4a5(t *testing.T) {
	Convey("ecode 0x4a5", t, func() {
		o, err := New(cfg)
		So(err, ShouldBeNil)
		defer o.Close()
		o.SetEmail(email).SetSig(uin, appid, keyType, sKey)
		rsp, err := o.Do(context.Background())
		So(err, ShouldBeNil)
		fmt.Printf("rsp(%+v)", rsp)
	})
}
