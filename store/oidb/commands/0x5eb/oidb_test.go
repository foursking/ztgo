package oidb_0x5eb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"git.code.oa.com/qdgo/core/net/udp"
	"git.code.oa.com/qdgo/core/store/oidb/pbv2"

	"github.com/gogo/protobuf/proto"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	cfg = &pbv2.Config{
		Cmd:         "0x5eb",
		ServiceType: 254,
		UDPClient: &udp.ClientConfig{
			RemoteAddr:   "l5://737473:917504",
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
		},
	}
)

func TestOidb_Do(t *testing.T) {
	Convey("ecode oidb do", t, func() {
		o := pbv2.New(cfg)
		o.SetClientIP("1.2.3.4")
		req := &ReqBody{
			RptUint64Uins:   []uint64{79461183, 251184321},
			Uint32Appid:     proto.Uint32(715021401),
			Uint32ReqNick:   proto.Uint32(1),
			Uint32ReqFaceId: proto.Uint32(1),
			Uint32ReqAllow:  proto.Uint32(1),
			Uint32ReqGender: proto.Uint32(1),
			Uint32ReqPerson: proto.Uint32(1),
		}
		o.SetReq(req).SetRsp(&RspBody{})
		o.SetSig(79461183, 717048201, 1, []byte("@oEJ1H5jLy"))
		err := o.Do(context.Background())
		So(err, ShouldBeNil)
		fmt.Printf("res(%+v)", o.RspBody)
	})
}
