package oidb_0x88d

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/foursking/ztgo/net/udp"
	"github.com/foursking/ztgo/store/oidb/pbv2"

	"github.com/gogo/protobuf/proto"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	cfg = &pbv2.Config{
		Cmd:         "0x88d",
		ServiceType: 0,
		UDPClient: &udp.ClientConfig{
			RemoteAddr:   "l5://1052993:65536",
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
		},
	}
)

func TestOidb_Do(t *testing.T) {
	Convey("ecode oidb do", t, func() {
		o := pbv2.New(cfg)
		req := &ReqBody{
			Uint32Appid: proto.Uint32(715021401),
			Stzreqgroupinfo: []*ReqGroupInfo{
				&ReqGroupInfo{
					Uint64GroupCode: proto.Uint64(562357819),
					Stgroupinfo: &GroupInfo{
						StringGroupName:      []byte(""),
						Uint32GroupLevel:     proto.Uint32(0),
						StringGroupMemo:      []byte(""),
						Uint32GroupMemberNum: proto.Uint32(0),
					},
				},
			},
		}
		o.SetReq(req).SetRsp(&RspBody{})
		o.SetSig(251184321, 717048201, 1, []byte("@KzUNlVWtr"))
		err := o.Do(context.Background())
		So(err, ShouldBeNil)
		fmt.Printf("res(%+v) type(%v)", o.RspBody, o.RspBody.String())
	})
}
