package oidb_0x4c8

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
		Cmd:         "0x4c8",
		ServiceType: 0,
		UDPClient: &udp.ClientConfig{
			RemoteAddr:   "l5://737473:917504",
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
		},
	}
)

func TestOidb_Do(t *testing.T) {
	Convey("test oidb do", t, func() {
		o := pbv2.New(cfg)
		o.SetClientIP("1.2.3.4")
		dstUserInfo := &ReqUsrInfo{
			DstUin:    proto.Uint64(623009470),
			Timestamp: proto.Uint32(0),
		}
		req := &QQHeadUrlReq{
			SrcUsrType:  proto.Uint32(1),
			SrcUin:      proto.Uint64(0),
			DstUsrType:  proto.Uint32(1),
			DstUsrInfos: []*ReqUsrInfo{dstUserInfo},
		}

		o.SetReq(req).SetRsp(&QQHeadUrlRsp{})
		o.SetSig(623009470, 717048201, 1, []byte("@2sSGBLnFJ"))
		err := o.Do(context.Background())
		So(err, ShouldBeNil)
		fmt.Printf("res(%+v)", o.RspBody)
	})
}
