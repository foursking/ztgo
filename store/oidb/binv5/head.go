package binv5

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"git.code.oa.com/qdgo/core/net/ip"
)

const (
	_stx   = 0xa
	_etx   = 0x3
	_ver   = 5
	_exVer = 600
)

// Base txt:cStx=0xa，cEtx=0x3包结构：cStx+ base + ext + Body + cEtx [命令号大于等于0x400时]
type Base struct {
	Len     uint16
	Version uint16 //接口版本号,为5
	Command uint16
	Uin     uint32
	Result  uint8
}

// Assistant 辅助信息结构
type Assistant struct {
	Version     uint16 // 辅助信息结构版本号
	Username    [11]byte
	Password    [11]byte
	ServiceIP   uint32 // 前端应用主机IP
	ServiceName [16]byte
	ServiceTime uint32 // 前端访问接口时间 time_t
	ServiceSeq  uint32
	ServiceType uint8  // 子命令字
	ClientIP    uint32 // 触发前端服务的用户IP
	ClientName  [21]byte
	ClientUin   uint32   // 用户uin
	Flag        uint32   // 传递标志
	Desc        [30]byte // 备注说明
}

// ext trans pkg head ext
type Ext struct {
	Len           uint16
	Version       uint16
	AppID         uint32
	KeyType       uint8
	KeyLen        uint16
	SessionKey    []byte //最长160
	ReserveBufLen uint16
	ReservedBuf   []byte //最长64
	ContextLen    uint16
	ContextData   []byte //最长64
}

// Head .
type Head struct {
	stx uint8
	*Base
	*Assistant
	*Ext
	Body []byte
	etx  uint8
}

// New 新建一个版本5的oidb包头
func New(command uint16, serviceType uint8) *Head {
	return &Head{
		stx: _stx,
		Base: &Base{
			Version: _ver, // 默认版本5， 版本6支持 ipv6
			Command: command,
		},
		Assistant: &Assistant{
			ServiceType: serviceType,
			ServiceIP:   ip.Atoi(ip.InternalIP()),
		},
		Ext: &Ext{
			Version: _exVer,
		},
		etx: _etx,
	}
}

// SetSig 设置登录态信息
func (h *Head) SetSig(uin uint64, appid uint32, keyType uint32, key []byte) {
	h.Base.Uin = uint32(uin)
	h.Ext.AppID = appid
	h.Ext.KeyType = uint8(keyType)
	h.Ext.SessionKey = key
}

// SetClientIP .
func (h *Head) SetClientIP(clientIP uint32) {
	h.Assistant.ClientIP = clientIP
}

// SetBody .
func (h *Head) SetBody(b []byte) {
	h.Body = b
}

// ResponseCode 请求 oidb 后的返回码
func (h *Head) ResponseCode() uint8 {
	return h.Base.Result
}

// Marshal marshal to bytes for request
func (h *Head) Marshal() ([]byte, error) {
	// 填写长度字段
	h.Ext.KeyLen = uint16(binary.Size(h.Ext.SessionKey))
	h.Ext.ReserveBufLen = uint16(binary.Size(h.Ext.ReservedBuf))
	h.Ext.ContextLen = uint16(binary.Size(h.Ext.ContextData))
	h.Ext.Len = uint16(binary.Size(h.Ext.Len) +
		binary.Size(h.Ext.Version) +
		binary.Size(h.Ext.AppID) +
		binary.Size(h.Ext.KeyType) +
		binary.Size(h.Ext.KeyLen) +
		binary.Size(h.Ext.SessionKey) +
		binary.Size(h.Ext.ReserveBufLen) +
		binary.Size(h.Ext.ReservedBuf) +
		binary.Size(h.Ext.ContextLen) +
		binary.Size(h.Ext.ContextData))
	h.Base.Len = uint16(binary.Size(h.stx) + binary.Size(h.Base) + binary.Size(h.Assistant) + int(h.Ext.Len) + binary.Size(h.Body) + binary.Size(h.etx))
	// 开始打包
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, h.stx)
	_ = binary.Write(buf, binary.BigEndian, h.Base)
	_ = binary.Write(buf, binary.BigEndian, h.Assistant)
	_ = binary.Write(buf, binary.BigEndian, h.Ext.Len)
	_ = binary.Write(buf, binary.BigEndian, h.Ext.Version)
	_ = binary.Write(buf, binary.BigEndian, h.Ext.AppID)
	_ = binary.Write(buf, binary.BigEndian, h.Ext.KeyType)
	_ = binary.Write(buf, binary.BigEndian, h.Ext.KeyLen)
	_ = binary.Write(buf, binary.BigEndian, h.Ext.SessionKey)
	_ = binary.Write(buf, binary.BigEndian, h.Ext.ReserveBufLen)
	_ = binary.Write(buf, binary.BigEndian, h.Ext.ReservedBuf)
	_ = binary.Write(buf, binary.BigEndian, h.Ext.ContextLen)
	_ = binary.Write(buf, binary.BigEndian, h.Ext.ContextData)
	_ = binary.Write(buf, binary.BigEndian, h.Body)
	_ = binary.Write(buf, binary.BigEndian, h.etx)
	return buf.Bytes(), nil
}

// Unmarshal unmarshal from response
func (h *Head) Unmarshal(data []byte) (err error) {
	if len(data) < binary.Size(h) {
		err = fmt.Errorf("oidb bin.v5 data len(%d) too short", len(data))
		return
	}
	buf := bytes.NewBuffer(data[:])
	if err = binary.Read(buf, binary.BigEndian, &h.stx); err != nil {
		return
	}
	if h.stx != _stx {
		err = fmt.Errorf("oidb bin.v5 stx error:0x%x", h.stx)
		return
	}
	if err = binary.Read(buf, binary.BigEndian, h.Base); err != nil {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, h.Assistant); err != nil {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &h.Ext.Len); err != nil {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &h.Ext.Version); err != nil {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &h.Ext.AppID); err != nil {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &h.Ext.KeyType); err != nil {
		return
	}
	var bufLen uint16
	if err = binary.Read(buf, binary.BigEndian, &bufLen); err != nil {
		return
	}
	h.Ext.SessionKey = make([]byte, bufLen)
	if err = binary.Read(buf, binary.BigEndian, &h.Ext.SessionKey); err != nil {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &bufLen); err != nil {
		return
	}
	h.Ext.ReservedBuf = make([]byte, bufLen)
	if err = binary.Read(buf, binary.BigEndian, &h.Ext.ReservedBuf); err != nil {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &bufLen); err != nil {
		return
	}
	h.Ext.ContextData = make([]byte, bufLen)
	if err = binary.Read(buf, binary.BigEndian, &h.Ext.ContextData); err != nil {
		return
	}
	bufLen = uint16(int(h.Base.Len) - binary.Size(h.stx) - binary.Size(h.Base) - binary.Size(h.Assistant) - int(h.Ext.Len) - binary.Size(h.etx))
	h.Body = make([]byte, bufLen)
	if err = binary.Read(buf, binary.BigEndian, &h.Body); err != nil {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, h.etx); err != nil {
		return
	}
	if h.etx != _etx {
		err = fmt.Errorf("oidb bin.v5 etx error:0x%x", h.etx)
	}
	return
}
