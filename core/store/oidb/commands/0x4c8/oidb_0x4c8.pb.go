// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: library/database/oidb/commands/0x4c8/oidb_0x4c8.proto

package oidb_0x4c8

import (
	fmt "fmt"
	github_com_gogo_protobuf_proto "github.com/gogo/protobuf/proto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type ReqUsrInfo struct {
	DstUin               *uint64  `protobuf:"varint,1,req,name=dstUin" json:"dstUin,omitempty"`
	Timestamp            *uint32  `protobuf:"varint,2,req,name=timestamp" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReqUsrInfo) Reset()         { *m = ReqUsrInfo{} }
func (m *ReqUsrInfo) String() string { return proto.CompactTextString(m) }
func (*ReqUsrInfo) ProtoMessage()    {}
func (*ReqUsrInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_31fba6237500278c, []int{0}
}
func (m *ReqUsrInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ReqUsrInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ReqUsrInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ReqUsrInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReqUsrInfo.Merge(m, src)
}
func (m *ReqUsrInfo) XXX_Size() int {
	return m.Size()
}
func (m *ReqUsrInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ReqUsrInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ReqUsrInfo proto.InternalMessageInfo

func (m *ReqUsrInfo) GetDstUin() uint64 {
	if m != nil && m.DstUin != nil {
		return *m.DstUin
	}
	return 0
}

func (m *ReqUsrInfo) GetTimestamp() uint32 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

type QQHeadUrlReq struct {
	SrcUsrType           *uint32       `protobuf:"varint,1,req,name=srcUsrType" json:"srcUsrType,omitempty"`
	SrcUin               *uint64       `protobuf:"varint,2,req,name=srcUin" json:"srcUin,omitempty"`
	DstUsrType           *uint32       `protobuf:"varint,3,req,name=dstUsrType" json:"dstUsrType,omitempty"`
	DstUsrInfos          []*ReqUsrInfo `protobuf:"bytes,4,rep,name=dstUsrInfos" json:"dstUsrInfos,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *QQHeadUrlReq) Reset()         { *m = QQHeadUrlReq{} }
func (m *QQHeadUrlReq) String() string { return proto.CompactTextString(m) }
func (*QQHeadUrlReq) ProtoMessage()    {}
func (*QQHeadUrlReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_31fba6237500278c, []int{1}
}
func (m *QQHeadUrlReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QQHeadUrlReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QQHeadUrlReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QQHeadUrlReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QQHeadUrlReq.Merge(m, src)
}
func (m *QQHeadUrlReq) XXX_Size() int {
	return m.Size()
}
func (m *QQHeadUrlReq) XXX_DiscardUnknown() {
	xxx_messageInfo_QQHeadUrlReq.DiscardUnknown(m)
}

var xxx_messageInfo_QQHeadUrlReq proto.InternalMessageInfo

func (m *QQHeadUrlReq) GetSrcUsrType() uint32 {
	if m != nil && m.SrcUsrType != nil {
		return *m.SrcUsrType
	}
	return 0
}

func (m *QQHeadUrlReq) GetSrcUin() uint64 {
	if m != nil && m.SrcUin != nil {
		return *m.SrcUin
	}
	return 0
}

func (m *QQHeadUrlReq) GetDstUsrType() uint32 {
	if m != nil && m.DstUsrType != nil {
		return *m.DstUsrType
	}
	return 0
}

func (m *QQHeadUrlReq) GetDstUsrInfos() []*ReqUsrInfo {
	if m != nil {
		return m.DstUsrInfos
	}
	return nil
}

type RspHeadInfo struct {
	DstUin               *uint64  `protobuf:"varint,1,req,name=dstUin" json:"dstUin,omitempty"`
	FaceType             *uint32  `protobuf:"varint,2,req,name=faceType" json:"faceType,omitempty"`
	Timestamp            *uint32  `protobuf:"varint,3,req,name=timestamp" json:"timestamp,omitempty"`
	FaceFlag             *uint32  `protobuf:"varint,4,req,name=faceFlag" json:"faceFlag,omitempty"`
	Url                  *string  `protobuf:"bytes,5,req,name=url" json:"url,omitempty"`
	Sysid                *uint32  `protobuf:"varint,6,opt,name=sysid" json:"sysid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RspHeadInfo) Reset()         { *m = RspHeadInfo{} }
func (m *RspHeadInfo) String() string { return proto.CompactTextString(m) }
func (*RspHeadInfo) ProtoMessage()    {}
func (*RspHeadInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_31fba6237500278c, []int{2}
}
func (m *RspHeadInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RspHeadInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RspHeadInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RspHeadInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RspHeadInfo.Merge(m, src)
}
func (m *RspHeadInfo) XXX_Size() int {
	return m.Size()
}
func (m *RspHeadInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_RspHeadInfo.DiscardUnknown(m)
}

var xxx_messageInfo_RspHeadInfo proto.InternalMessageInfo

func (m *RspHeadInfo) GetDstUin() uint64 {
	if m != nil && m.DstUin != nil {
		return *m.DstUin
	}
	return 0
}

func (m *RspHeadInfo) GetFaceType() uint32 {
	if m != nil && m.FaceType != nil {
		return *m.FaceType
	}
	return 0
}

func (m *RspHeadInfo) GetTimestamp() uint32 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

func (m *RspHeadInfo) GetFaceFlag() uint32 {
	if m != nil && m.FaceFlag != nil {
		return *m.FaceFlag
	}
	return 0
}

func (m *RspHeadInfo) GetUrl() string {
	if m != nil && m.Url != nil {
		return *m.Url
	}
	return ""
}

func (m *RspHeadInfo) GetSysid() uint32 {
	if m != nil && m.Sysid != nil {
		return *m.Sysid
	}
	return 0
}

type QQHeadUrlRsp struct {
	SrcUsrType           *uint32        `protobuf:"varint,1,req,name=srcUsrType" json:"srcUsrType,omitempty"`
	SrcUin               *uint64        `protobuf:"varint,2,req,name=srcUin" json:"srcUin,omitempty"`
	Result               *int32         `protobuf:"varint,3,req,name=result" json:"result,omitempty"`
	DstUsrType           *uint32        `protobuf:"varint,4,req,name=dstUsrType" json:"dstUsrType,omitempty"`
	DstHeadInfos         []*RspHeadInfo `protobuf:"bytes,5,rep,name=dstHeadInfos" json:"dstHeadInfos,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *QQHeadUrlRsp) Reset()         { *m = QQHeadUrlRsp{} }
func (m *QQHeadUrlRsp) String() string { return proto.CompactTextString(m) }
func (*QQHeadUrlRsp) ProtoMessage()    {}
func (*QQHeadUrlRsp) Descriptor() ([]byte, []int) {
	return fileDescriptor_31fba6237500278c, []int{3}
}
func (m *QQHeadUrlRsp) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QQHeadUrlRsp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QQHeadUrlRsp.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QQHeadUrlRsp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QQHeadUrlRsp.Merge(m, src)
}
func (m *QQHeadUrlRsp) XXX_Size() int {
	return m.Size()
}
func (m *QQHeadUrlRsp) XXX_DiscardUnknown() {
	xxx_messageInfo_QQHeadUrlRsp.DiscardUnknown(m)
}

var xxx_messageInfo_QQHeadUrlRsp proto.InternalMessageInfo

func (m *QQHeadUrlRsp) GetSrcUsrType() uint32 {
	if m != nil && m.SrcUsrType != nil {
		return *m.SrcUsrType
	}
	return 0
}

func (m *QQHeadUrlRsp) GetSrcUin() uint64 {
	if m != nil && m.SrcUin != nil {
		return *m.SrcUin
	}
	return 0
}

func (m *QQHeadUrlRsp) GetResult() int32 {
	if m != nil && m.Result != nil {
		return *m.Result
	}
	return 0
}

func (m *QQHeadUrlRsp) GetDstUsrType() uint32 {
	if m != nil && m.DstUsrType != nil {
		return *m.DstUsrType
	}
	return 0
}

func (m *QQHeadUrlRsp) GetDstHeadInfos() []*RspHeadInfo {
	if m != nil {
		return m.DstHeadInfos
	}
	return nil
}

func init() {
	proto.RegisterType((*ReqUsrInfo)(nil), "oidb_0x4c8.ReqUsrInfo")
	proto.RegisterType((*QQHeadUrlReq)(nil), "oidb_0x4c8.QQHeadUrlReq")
	proto.RegisterType((*RspHeadInfo)(nil), "oidb_0x4c8.RspHeadInfo")
	proto.RegisterType((*QQHeadUrlRsp)(nil), "oidb_0x4c8.QQHeadUrlRsp")
}

func init() {
	proto.RegisterFile("library/database/oidb/commands/0x4c8/oidb_0x4c8.proto", fileDescriptor_31fba6237500278c)
}

var fileDescriptor_31fba6237500278c = []byte{
	// 350 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x90, 0xcf, 0x4a, 0xeb, 0x40,
	0x14, 0xc6, 0xc9, 0xbf, 0x72, 0x7b, 0xda, 0xc2, 0x25, 0x5c, 0x7a, 0xc3, 0xe5, 0x12, 0x42, 0x56,
	0x59, 0x35, 0x22, 0x0a, 0x05, 0x77, 0x5d, 0x88, 0x2e, 0x3b, 0x98, 0xb5, 0x4c, 0x33, 0x53, 0x09,
	0xe4, 0x5f, 0xe7, 0x4c, 0xc1, 0xbe, 0x8d, 0xe0, 0x63, 0xf8, 0x02, 0x2e, 0x7d, 0x04, 0xe9, 0x93,
	0xc8, 0x4c, 0x52, 0xd3, 0x54, 0x70, 0xe1, 0x6e, 0xbe, 0x8f, 0xf9, 0xce, 0x39, 0xdf, 0x0f, 0x2e,
	0xf3, 0x6c, 0x25, 0xa8, 0xd8, 0xc5, 0x8c, 0x4a, 0xba, 0xa2, 0xc8, 0xe3, 0x2a, 0x63, 0xab, 0x38,
	0xad, 0x8a, 0x82, 0x96, 0x0c, 0xe3, 0xb3, 0xc7, 0x8b, 0x74, 0xae, 0xbd, 0x7b, 0xfd, 0x9c, 0xd5,
	0xa2, 0x92, 0x95, 0x0b, 0x9d, 0x13, 0x2e, 0x00, 0x08, 0xdf, 0x24, 0x28, 0x6e, 0xcb, 0x75, 0xe5,
	0x4e, 0x61, 0xc0, 0x50, 0x26, 0x59, 0xe9, 0x19, 0x81, 0x19, 0xd9, 0xa4, 0x55, 0xee, 0x7f, 0x18,
	0xca, 0xac, 0xe0, 0x28, 0x69, 0x51, 0x7b, 0x66, 0x60, 0x46, 0x13, 0xd2, 0x19, 0xe1, 0x93, 0x01,
	0xe3, 0xe5, 0xf2, 0x86, 0x53, 0x96, 0x88, 0x9c, 0xf0, 0x8d, 0xeb, 0x03, 0xa0, 0x48, 0x13, 0x14,
	0x77, 0xbb, 0x9a, 0xeb, 0x51, 0x13, 0x72, 0xe4, 0xa8, 0x35, 0x4a, 0x65, 0xa5, 0x9e, 0x65, 0x93,
	0x56, 0xa9, 0x9c, 0x5a, 0xd8, 0xe6, 0xac, 0x26, 0xd7, 0x39, 0xee, 0x1c, 0x46, 0x8d, 0x52, 0xc7,
	0xa2, 0x67, 0x07, 0x56, 0x34, 0x3a, 0x9f, 0xce, 0x8e, 0x0a, 0x76, 0x5d, 0xc8, 0xf1, 0xd7, 0xf0,
	0xd9, 0x80, 0x11, 0xc1, 0x5a, 0xdd, 0xf8, 0x6d, 0xd1, 0x7f, 0xf0, 0x6b, 0x4d, 0x53, 0xae, 0xf7,
	0x37, 0x3d, 0x3f, 0x75, 0x1f, 0x82, 0x75, 0x02, 0xe1, 0x90, 0xbc, 0xce, 0xe9, 0x83, 0x67, 0x77,
	0x49, 0xa5, 0xdd, 0xdf, 0x60, 0x6d, 0x45, 0xee, 0x39, 0x81, 0x19, 0x0d, 0x89, 0x7a, 0xba, 0x7f,
	0xc0, 0xc1, 0x1d, 0x66, 0xcc, 0x1b, 0x04, 0x46, 0x34, 0x21, 0x8d, 0x08, 0x5f, 0x7a, 0x20, 0xb1,
	0xfe, 0x31, 0xc8, 0x29, 0x0c, 0x04, 0xc7, 0x6d, 0x2e, 0xf5, 0x9d, 0x0e, 0x69, 0xd5, 0x09, 0x60,
	0xfb, 0x0b, 0xe0, 0x2b, 0x18, 0x33, 0x94, 0x07, 0x4a, 0xe8, 0x39, 0x9a, 0xf0, 0xdf, 0x1e, 0xe1,
	0x8e, 0x22, 0xe9, 0x7d, 0x5e, 0x8c, 0x5f, 0xf7, 0xbe, 0xf1, 0xb6, 0xf7, 0x8d, 0xf7, 0xbd, 0x6f,
	0x7c, 0x04, 0x00, 0x00, 0xff, 0xff, 0x8a, 0x3f, 0x8d, 0x93, 0x9c, 0x02, 0x00, 0x00,
}

func (m *ReqUsrInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ReqUsrInfo) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.DstUin == nil {
		return 0, github_com_gogo_protobuf_proto.NewRequiredNotSetError("dstUin")
	} else {
		dAtA[i] = 0x8
		i++
		i = encodeVarintOidb_0X4C8(dAtA, i, uint64(*m.DstUin))
	}
	if m.Timestamp == nil {
		return 0, github_com_gogo_protobuf_proto.NewRequiredNotSetError("timestamp")
	} else {
		dAtA[i] = 0x10
		i++
		i = encodeVarintOidb_0X4C8(dAtA, i, uint64(*m.Timestamp))
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *QQHeadUrlReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QQHeadUrlReq) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.SrcUsrType == nil {
		return 0, github_com_gogo_protobuf_proto.NewRequiredNotSetError("srcUsrType")
	} else {
		dAtA[i] = 0x8
		i++
		i = encodeVarintOidb_0X4C8(dAtA, i, uint64(*m.SrcUsrType))
	}
	if m.SrcUin == nil {
		return 0, github_com_gogo_protobuf_proto.NewRequiredNotSetError("srcUin")
	} else {
		dAtA[i] = 0x10
		i++
		i = encodeVarintOidb_0X4C8(dAtA, i, uint64(*m.SrcUin))
	}
	if m.DstUsrType == nil {
		return 0, github_com_gogo_protobuf_proto.NewRequiredNotSetError("dstUsrType")
	} else {
		dAtA[i] = 0x18
		i++
		i = encodeVarintOidb_0X4C8(dAtA, i, uint64(*m.DstUsrType))
	}
	if len(m.DstUsrInfos) > 0 {
		for _, msg := range m.DstUsrInfos {
			dAtA[i] = 0x22
			i++
			i = encodeVarintOidb_0X4C8(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *RspHeadInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RspHeadInfo) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.DstUin == nil {
		return 0, github_com_gogo_protobuf_proto.NewRequiredNotSetError("dstUin")
	} else {
		dAtA[i] = 0x8
		i++
		i = encodeVarintOidb_0X4C8(dAtA, i, uint64(*m.DstUin))
	}
	if m.FaceType == nil {
		return 0, github_com_gogo_protobuf_proto.NewRequiredNotSetError("faceType")
	} else {
		dAtA[i] = 0x10
		i++
		i = encodeVarintOidb_0X4C8(dAtA, i, uint64(*m.FaceType))
	}
	if m.Timestamp == nil {
		return 0, github_com_gogo_protobuf_proto.NewRequiredNotSetError("timestamp")
	} else {
		dAtA[i] = 0x18
		i++
		i = encodeVarintOidb_0X4C8(dAtA, i, uint64(*m.Timestamp))
	}
	if m.FaceFlag == nil {
		return 0, github_com_gogo_protobuf_proto.NewRequiredNotSetError("faceFlag")
	} else {
		dAtA[i] = 0x20
		i++
		i = encodeVarintOidb_0X4C8(dAtA, i, uint64(*m.FaceFlag))
	}
	if m.Url == nil {
		return 0, github_com_gogo_protobuf_proto.NewRequiredNotSetError("url")
	} else {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintOidb_0X4C8(dAtA, i, uint64(len(*m.Url)))
		i += copy(dAtA[i:], *m.Url)
	}
	if m.Sysid != nil {
		dAtA[i] = 0x30
		i++
		i = encodeVarintOidb_0X4C8(dAtA, i, uint64(*m.Sysid))
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *QQHeadUrlRsp) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QQHeadUrlRsp) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.SrcUsrType == nil {
		return 0, github_com_gogo_protobuf_proto.NewRequiredNotSetError("srcUsrType")
	} else {
		dAtA[i] = 0x8
		i++
		i = encodeVarintOidb_0X4C8(dAtA, i, uint64(*m.SrcUsrType))
	}
	if m.SrcUin == nil {
		return 0, github_com_gogo_protobuf_proto.NewRequiredNotSetError("srcUin")
	} else {
		dAtA[i] = 0x10
		i++
		i = encodeVarintOidb_0X4C8(dAtA, i, uint64(*m.SrcUin))
	}
	if m.Result == nil {
		return 0, github_com_gogo_protobuf_proto.NewRequiredNotSetError("result")
	} else {
		dAtA[i] = 0x18
		i++
		i = encodeVarintOidb_0X4C8(dAtA, i, uint64(*m.Result))
	}
	if m.DstUsrType == nil {
		return 0, github_com_gogo_protobuf_proto.NewRequiredNotSetError("dstUsrType")
	} else {
		dAtA[i] = 0x20
		i++
		i = encodeVarintOidb_0X4C8(dAtA, i, uint64(*m.DstUsrType))
	}
	if len(m.DstHeadInfos) > 0 {
		for _, msg := range m.DstHeadInfos {
			dAtA[i] = 0x2a
			i++
			i = encodeVarintOidb_0X4C8(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintOidb_0X4C8(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *ReqUsrInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.DstUin != nil {
		n += 1 + sovOidb_0X4C8(uint64(*m.DstUin))
	}
	if m.Timestamp != nil {
		n += 1 + sovOidb_0X4C8(uint64(*m.Timestamp))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *QQHeadUrlReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SrcUsrType != nil {
		n += 1 + sovOidb_0X4C8(uint64(*m.SrcUsrType))
	}
	if m.SrcUin != nil {
		n += 1 + sovOidb_0X4C8(uint64(*m.SrcUin))
	}
	if m.DstUsrType != nil {
		n += 1 + sovOidb_0X4C8(uint64(*m.DstUsrType))
	}
	if len(m.DstUsrInfos) > 0 {
		for _, e := range m.DstUsrInfos {
			l = e.Size()
			n += 1 + l + sovOidb_0X4C8(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *RspHeadInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.DstUin != nil {
		n += 1 + sovOidb_0X4C8(uint64(*m.DstUin))
	}
	if m.FaceType != nil {
		n += 1 + sovOidb_0X4C8(uint64(*m.FaceType))
	}
	if m.Timestamp != nil {
		n += 1 + sovOidb_0X4C8(uint64(*m.Timestamp))
	}
	if m.FaceFlag != nil {
		n += 1 + sovOidb_0X4C8(uint64(*m.FaceFlag))
	}
	if m.Url != nil {
		l = len(*m.Url)
		n += 1 + l + sovOidb_0X4C8(uint64(l))
	}
	if m.Sysid != nil {
		n += 1 + sovOidb_0X4C8(uint64(*m.Sysid))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *QQHeadUrlRsp) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SrcUsrType != nil {
		n += 1 + sovOidb_0X4C8(uint64(*m.SrcUsrType))
	}
	if m.SrcUin != nil {
		n += 1 + sovOidb_0X4C8(uint64(*m.SrcUin))
	}
	if m.Result != nil {
		n += 1 + sovOidb_0X4C8(uint64(*m.Result))
	}
	if m.DstUsrType != nil {
		n += 1 + sovOidb_0X4C8(uint64(*m.DstUsrType))
	}
	if len(m.DstHeadInfos) > 0 {
		for _, e := range m.DstHeadInfos {
			l = e.Size()
			n += 1 + l + sovOidb_0X4C8(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovOidb_0X4C8(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozOidb_0X4C8(x uint64) (n int) {
	return sovOidb_0X4C8(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ReqUsrInfo) Unmarshal(dAtA []byte) error {
	var hasFields [1]uint64
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOidb_0X4C8
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ReqUsrInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ReqUsrInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DstUin", wireType)
			}
			var v uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOidb_0X4C8
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.DstUin = &v
			hasFields[0] |= uint64(0x00000001)
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var v uint32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOidb_0X4C8
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Timestamp = &v
			hasFields[0] |= uint64(0x00000002)
		default:
			iNdEx = preIndex
			skippy, err := skipOidb_0X4C8(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthOidb_0X4C8
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthOidb_0X4C8
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}
	if hasFields[0]&uint64(0x00000001) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("dstUin")
	}
	if hasFields[0]&uint64(0x00000002) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("timestamp")
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QQHeadUrlReq) Unmarshal(dAtA []byte) error {
	var hasFields [1]uint64
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOidb_0X4C8
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QQHeadUrlReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QQHeadUrlReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SrcUsrType", wireType)
			}
			var v uint32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOidb_0X4C8
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.SrcUsrType = &v
			hasFields[0] |= uint64(0x00000001)
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SrcUin", wireType)
			}
			var v uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOidb_0X4C8
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.SrcUin = &v
			hasFields[0] |= uint64(0x00000002)
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DstUsrType", wireType)
			}
			var v uint32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOidb_0X4C8
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.DstUsrType = &v
			hasFields[0] |= uint64(0x00000004)
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DstUsrInfos", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOidb_0X4C8
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthOidb_0X4C8
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthOidb_0X4C8
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DstUsrInfos = append(m.DstUsrInfos, &ReqUsrInfo{})
			if err := m.DstUsrInfos[len(m.DstUsrInfos)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipOidb_0X4C8(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthOidb_0X4C8
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthOidb_0X4C8
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}
	if hasFields[0]&uint64(0x00000001) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("srcUsrType")
	}
	if hasFields[0]&uint64(0x00000002) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("srcUin")
	}
	if hasFields[0]&uint64(0x00000004) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("dstUsrType")
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RspHeadInfo) Unmarshal(dAtA []byte) error {
	var hasFields [1]uint64
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOidb_0X4C8
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RspHeadInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RspHeadInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DstUin", wireType)
			}
			var v uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOidb_0X4C8
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.DstUin = &v
			hasFields[0] |= uint64(0x00000001)
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FaceType", wireType)
			}
			var v uint32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOidb_0X4C8
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.FaceType = &v
			hasFields[0] |= uint64(0x00000002)
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var v uint32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOidb_0X4C8
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Timestamp = &v
			hasFields[0] |= uint64(0x00000004)
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FaceFlag", wireType)
			}
			var v uint32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOidb_0X4C8
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.FaceFlag = &v
			hasFields[0] |= uint64(0x00000008)
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Url", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOidb_0X4C8
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthOidb_0X4C8
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOidb_0X4C8
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Url = &s
			iNdEx = postIndex
			hasFields[0] |= uint64(0x00000010)
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sysid", wireType)
			}
			var v uint32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOidb_0X4C8
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Sysid = &v
		default:
			iNdEx = preIndex
			skippy, err := skipOidb_0X4C8(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthOidb_0X4C8
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthOidb_0X4C8
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}
	if hasFields[0]&uint64(0x00000001) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("dstUin")
	}
	if hasFields[0]&uint64(0x00000002) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("faceType")
	}
	if hasFields[0]&uint64(0x00000004) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("timestamp")
	}
	if hasFields[0]&uint64(0x00000008) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("faceFlag")
	}
	if hasFields[0]&uint64(0x00000010) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("url")
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QQHeadUrlRsp) Unmarshal(dAtA []byte) error {
	var hasFields [1]uint64
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOidb_0X4C8
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QQHeadUrlRsp: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QQHeadUrlRsp: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SrcUsrType", wireType)
			}
			var v uint32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOidb_0X4C8
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.SrcUsrType = &v
			hasFields[0] |= uint64(0x00000001)
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SrcUin", wireType)
			}
			var v uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOidb_0X4C8
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.SrcUin = &v
			hasFields[0] |= uint64(0x00000002)
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Result", wireType)
			}
			var v int32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOidb_0X4C8
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Result = &v
			hasFields[0] |= uint64(0x00000004)
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DstUsrType", wireType)
			}
			var v uint32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOidb_0X4C8
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.DstUsrType = &v
			hasFields[0] |= uint64(0x00000008)
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DstHeadInfos", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOidb_0X4C8
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthOidb_0X4C8
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthOidb_0X4C8
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DstHeadInfos = append(m.DstHeadInfos, &RspHeadInfo{})
			if err := m.DstHeadInfos[len(m.DstHeadInfos)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipOidb_0X4C8(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthOidb_0X4C8
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthOidb_0X4C8
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}
	if hasFields[0]&uint64(0x00000001) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("srcUsrType")
	}
	if hasFields[0]&uint64(0x00000002) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("srcUin")
	}
	if hasFields[0]&uint64(0x00000004) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("result")
	}
	if hasFields[0]&uint64(0x00000008) == 0 {
		return github_com_gogo_protobuf_proto.NewRequiredNotSetError("dstUsrType")
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipOidb_0X4C8(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowOidb_0X4C8
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowOidb_0X4C8
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowOidb_0X4C8
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthOidb_0X4C8
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthOidb_0X4C8
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowOidb_0X4C8
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipOidb_0X4C8(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthOidb_0X4C8
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthOidb_0X4C8 = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowOidb_0X4C8   = fmt.Errorf("proto: integer overflow")
)
