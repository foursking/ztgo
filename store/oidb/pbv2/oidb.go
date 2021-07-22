package pbv2

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"strconv"
	"sync/atomic"
	"time"

	"git.code.oa.com/qdgo/core/config/env"
	"git.code.oa.com/qdgo/core/log"
	"git.code.oa.com/qdgo/core/net/ip"
	"git.code.oa.com/qdgo/core/net/udp"
	"git.code.oa.com/qdgo/core/store/oidb"

	"github.com/golang/protobuf/proto"
	pkgerr "github.com/pkg/errors"
)

const (
	_stx uint8 = 0x28
	_etx uint8 = 0x29
)

// Config .
type Config struct {
	cmd         uint32 // cmd uint32 format
	Cmd         string // cmd string format
	ServiceType uint32
	UDPClient   *udp.ClientConfig
}

// Oidb .
type Oidb struct {
	c   *Config
	seq uint32 // req sequence number

	Head    OIDBHead
	ReqBody proto.Message
	RspBody proto.Message
}

// New .
func New(c *Config) *Oidb {
	fixConfig(c)
	return &Oidb{
		c: c,
		Head: OIDBHead{
			Uint32Command:     proto.Uint32(c.cmd),
			Uint32ServiceType: proto.Uint32(c.ServiceType),
			StrServiceName:    proto.String(env.AppName),
		},
	}
}

func fixConfig(c *Config) {
	if c.Cmd == "" {
		panic("invalid oidb cmd")
	}
	cmd, err := strconv.ParseInt(c.Cmd, 0, 64)
	if err != nil {
		panic("invalid oidb cmd")
	}
	c.cmd = uint32(cmd)
	if c.UDPClient == nil {
		panic("invalid oidb udp client config")
	}
	if c.UDPClient.ReadTimeout <= 0 {
		c.UDPClient.ReadTimeout = 200 * time.Millisecond
	}
	if c.UDPClient.WriteTimeout <= 0 {
		c.UDPClient.WriteTimeout = 200 * time.Millisecond
	}
}

// SetHead .
func (o *Oidb) SetClientIP(clientIP string) {
	sip := ip.Atoi(ip.InternalIP()) // server ip
	cip := ip.Atoi(clientIP)        // client ip
	o.Head.Uint32FromAddr = proto.Uint32(cip)
	o.Head.Uint32LocalAddr = proto.Uint32(sip)
}

func (o *Oidb) newSeq() uint32 {
	return atomic.AddUint32(&o.seq, 1)
}

// SetReq .
func (o *Oidb) SetReq(req proto.Message) *Oidb {
	o.ReqBody = req
	return o
}

// SetRsp .
func (o *Oidb) SetRsp(rsp proto.Message) *Oidb {
	o.RspBody = rsp
	return o
}

// SetSig 设置oidb登录态 如 SetSig(181431178, 0, 1, []byte("skey"))
func (o *Oidb) SetSig(uin uint64, appid uint32, keyType uint32, key []byte) *Oidb {
	o.Head.Uint64Uin = proto.Uint64(uin)
	o.Head.MsgLoginSig = &LoginSig{
		Uint32Type:  proto.Uint32(keyType),
		BytesSig:    key,
		Uint32Appid: proto.Uint32(appid),
	}
	return o
}

func (o *Oidb) marshal() ([]byte, error) {
	encHead, err := proto.Marshal(&o.Head)
	if err != nil {
		return nil, errors.New("oidb marshal req head error")
	}
	if o.ReqBody == nil {
		return nil, errors.New("oidb req body empty")
	}
	encBody, err := proto.Marshal(o.ReqBody)
	if err != nil {
		return nil, errors.New("oidb marshal req body error")
	}
	var headLen = uint32(len(encHead))
	var bodyLen = uint32(len(encBody))
	var packetLen = int(10 + headLen + bodyLen) // 10 == stx + headLen + bodyLen + etx
	buf := new(bytes.Buffer)
	buf.Grow(packetLen)
	if err := binary.Write(buf, binary.BigEndian, _stx); err != nil {
		return nil, pkgerr.Wrap(err, "oidb Marshal() write stx error")
	}
	if err := binary.Write(buf, binary.BigEndian, headLen); err != nil {
		return nil, pkgerr.Wrap(err, "oidb Marshal() write head len error")
	}
	if err := binary.Write(buf, binary.BigEndian, bodyLen); err != nil {
		return nil, pkgerr.Wrap(err, "oidb Marshal() write body len error")
	}
	if headLen > 0 {
		if err := binary.Write(buf, binary.BigEndian, encHead); err != nil {
			return nil, pkgerr.Wrap(err, "oidb Marshal() write enc head error")
		}
	}
	if bodyLen > 0 {
		if err := binary.Write(buf, binary.BigEndian, encBody); err != nil {
			return nil, pkgerr.Wrap(err, "oidb Marshal() write enc body error")
		}
	}
	if err := binary.Write(buf, binary.BigEndian, _etx); err != nil {
		return nil, pkgerr.Wrap(err, "oidb Marshal() write etx error")
	}
	return buf.Bytes(), nil
}

func (o *Oidb) unmarshal(data []byte) (err error) {
	dataLen := uint32(len(data))
	if dataLen < 10 || _stx != data[0] || _etx != data[dataLen-1] {
		return errors.New("oidb unmarshal bad pkg")
	}
	var (
		headLen uint32
		bodyLen uint32
		buf     = bytes.NewBuffer(data[1:9])
	)
	if err = binary.Read(buf, binary.BigEndian, &headLen); err != nil {
		return err
	}
	if 1+4+4+headLen+1 > dataLen {
		err = pkgerr.Wrapf(errors.New("oidb unmarshal head len invalid"), "data(%s)", data)
		return
	}
	if err = proto.Unmarshal(data[9:9+headLen], &o.Head); err != nil {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &bodyLen); err != nil {
		return
	}
	if 1+4+4+headLen+bodyLen+1 > dataLen {
		err = pkgerr.Wrapf(errors.New("oidb unmarshal body len invalid"), "data(%s)", data)
		return
	}
	if err = proto.Unmarshal(data[9+headLen:9+headLen+bodyLen], o.RspBody); err != nil {
		return
	}
	return
}

// Do .
func (o *Oidb) Do(ctx context.Context) (err error) {
	now := time.Now()
	o.seq = o.newSeq()
	req, err := o.marshal()
	if err != nil {
		log.Error("oidb Marshal() proto(%s) serviceType(%d) error(%v)", o.c.Cmd, o.c.ServiceType, err)
		return
	}
	udpCli, err := udp.NewClient(o.c.UDPClient)
	if err != nil {
		log.Error("new oidb proto(%s) serviceType(%d) error(%v)", o.c.Cmd, o.c.ServiceType, err)
		return
	}
	defer udpCli.Close()
	data, err := udpCli.SendAndReceive(ctx, req)
	oidb.MetricReqDur.Observe(int64(time.Since(now)/time.Millisecond), o.c.Cmd, o.c.UDPClient.RemoteAddr, "SendAndReceive")
	if err != nil {
		log.Error("oidb request proto(%s) serviceType(%d) error(%v)", o.c.Cmd, o.c.ServiceType, err)
		return
	}
	if err = o.unmarshal(data); err != nil {
		log.Error("oidb Unmarshal() proto(%s) serviceType(%d) error(%v)", o.c.Cmd, o.c.ServiceType, err)
	}
	return
}
