package oidb_0x485

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"sync/atomic"

	"git.code.oa.com/qdgo/core/log"
	"git.code.oa.com/qdgo/core/net/udp"
	"git.code.oa.com/qdgo/core/store/oidb/binv5"

	pkgerr "github.com/pkg/errors"
)

const (
	_cmd = 0x485

	// ErrCodeOk 成功
	ErrCodeOk = uint8(0)
	// ErrCodeNoUin 号码不存在
	ErrCodeNoUin = uint8(0x12)
	// ErrCodeServiceType cServiceType异常
	ErrCodeServiceType = uint8(0x13)
)

var (
	// ErrNoUin .
	ErrNoUin = errors.New("no uin was found")
	// ErrUnknown .
	ErrServiceType = errors.New("cServiceType error")

	//ErrDefault .
	ErrDefault = errors.New("0x485 default error")
)

// Response .
type Response struct {
	Flag uint8
}

type Oidb0x485 struct {
	head     *binv5.Head
	response *Response
	udpCli   *udp.Client
	seq      uint32
}

// New .
func New(c *binv5.Config) (*Oidb0x485, error) {
	cli, err := udp.NewClient(c.UDPClient)
	if err != nil {
		return nil, pkgerr.Wrap(err, "Oidb 0x4a5 udp client l5 error")
	}
	o := &Oidb0x485{
		head:     binv5.New(_cmd, c.ServiceType),
		response: new(Response),
		udpCli:   cli,
	}
	return o, nil
}

// SetSig .
func (o *Oidb0x485) SetSig(uin uint64, appid uint32, keyType uint32, key []byte) *Oidb0x485 {
	o.head.SetSig(uin, appid, keyType, key)
	return o
}

// marshal 先 marshal body, marshal head 里用到了 body 长度
func (o *Oidb0x485) marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	o.head.SetBody(buf.Bytes())
	return o.head.Marshal()
}

func (o *Oidb0x485) unmarshal(data []byte) (err error) {
	if err = o.head.Unmarshal(data); err != nil {
		log.Errorf("oidb 0x485 unmarshal head(%s) code(%d) error(%v)", data, o.head.ResponseCode(), err)
		return
	}
	switch o.head.ResponseCode() {
	case ErrCodeOk: // success
	case ErrCodeNoUin:
		return ErrNoUin
	case ErrCodeServiceType:
		return ErrServiceType
	default:
		return ErrDefault
	}
	buf := bytes.NewBuffer(o.head.Body)
	if err = binary.Read(buf, binary.BigEndian, &o.response.Flag); err != nil {
		err = pkgerr.Wrap(err, "oidb 0x485 unmarshal read flag")
	}
	return
}

func (o *Oidb0x485) newSeq() uint32 {
	return atomic.AddUint32(&o.seq, 1)
}

// Do .
func (o *Oidb0x485) Do(ctx context.Context) (rsp *Response, err error) {
	o.head.Assistant.ServiceSeq = o.newSeq()
	req, err := o.marshal()
	if err != nil {
		log.Errorf("oidb Do() marshal() error(%v)", err)
		return
	}
	data, err := o.udpCli.SendAndReceive(ctx, req)
	if err != nil {
		log.Errorf("oidb Do() SendAndReceive() error(%v)", err)
		return
	}
	if err = o.unmarshal(data); err != nil {
		log.Errorf("oidb Unmarshal() error(%v)", err)
		return
	}
	rsp = o.response
	return
}

// Close .
func (o *Oidb0x485) Close() error {
	return o.udpCli.Close()
}
