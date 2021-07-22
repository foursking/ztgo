package oidb_0x4a5

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
	_cmd = 0x4a5

	// ErrCodeOk 成功
	ErrCodeOk = uint8(0)
	// ErrCodeNoUin 无对应uin
	ErrCodeNoUin = uint8(1)
	// ErrCodeInvalidEmail 账号不合法
	ErrCodeInvalidEmail = uint8(2)
	// ErrCodeUnknown 未知错误
	ErrCodeUnknown = uint8(10)
)

var (
	// ErrNoUin .
	ErrNoUin = errors.New("no uin was found")
	// ErrInvalidEmail .
	ErrInvalidEmail = errors.New("invalid email")
	// ErrUnknown .
	ErrUnknown = errors.New("unknown error")
	// ErrInvalidCode 未知的错误码
	ErrInvalidCode = errors.New("unknown error code")
)

// Config .
type Config struct {
	ServiceType uint8
	UDPClient   *udp.ClientConfig
}

// Response .
type Response struct {
	Uin  uint32
	Flag uint32
}

// Oidb0x4a5 单个email帐号查询uin
// http://oidb2.server.com/metronic/html/protocolfile/protocolList.php?appid=1189
type Oidb0x4a5 struct {
	head     *binv5.Head
	response *Response
	email    string
	udpCli   *udp.Client
	seq      uint32
}

// New .
func New(c *Config) (*Oidb0x4a5, error) {
	cli, err := udp.NewClient(c.UDPClient)
	if err != nil {
		return nil, pkgerr.Wrap(err, "Oidb 0x4a5 udp client l5 error")
	}
	o := &Oidb0x4a5{
		head:     binv5.New(_cmd, c.ServiceType),
		response: new(Response),
		udpCli:   cli,
	}
	return o, nil
}

// SetEmail .
func (o *Oidb0x4a5) SetEmail(email string) *Oidb0x4a5 {
	o.email = email
	return o
}

// SetSig .
func (o *Oidb0x4a5) SetSig(uin uint64, appid uint32, keyType uint32, key []byte) *Oidb0x4a5 {
	o.head.SetSig(uin, appid, keyType, key)
	return o
}

// marshal 先 marshal body, marshal head 里用到了 body 长度
func (o *Oidb0x4a5) marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, o.head.ServiceType)
	binary.Write(buf, binary.BigEndian, uint8(len(o.email)))
	binary.Write(buf, binary.BigEndian, []byte(o.email))
	o.head.SetBody(buf.Bytes())
	return o.head.Marshal()
}

func (o *Oidb0x4a5) unmarshal(data []byte) (err error) {
	if err = o.head.Unmarshal(data); err != nil {
		log.Errorf("oidb 0x4a5 unmarshal head(%s) code(%d) error(%v)", data, o.head.ResponseCode(), err)
		return
	}
	switch o.head.ResponseCode() {
	case ErrCodeOk: // success
	case ErrCodeNoUin:
		return ErrNoUin
	case ErrCodeInvalidEmail:
		return ErrInvalidEmail
	case ErrCodeUnknown:
		return ErrUnknown
	default:
		return pkgerr.Wrapf(ErrInvalidCode, "code is (%d)", o.head.ResponseCode())
	}
	buf := bytes.NewBuffer(o.head.Body[2:]) // 忽略前2个字节，用不到
	if err = binary.Read(buf, binary.BigEndian, &o.response.Uin); err != nil {
		err = pkgerr.Wrap(err, "oidb 0x4a5 unmarshal read uin")
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &o.response.Flag); err != nil {
		err = pkgerr.Wrap(err, "oidb 0x4a5 unmarshal read flag")
	}
	return
}

func (o *Oidb0x4a5) newSeq() uint32 {
	return atomic.AddUint32(&o.seq, 1)
}

// Do .
func (o *Oidb0x4a5) Do(ctx context.Context) (rsp *Response, err error) {
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
func (o *Oidb0x4a5) Close() error {
	return o.udpCli.Close()
}
