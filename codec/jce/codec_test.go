package jce

// XXX: benchmark encode and decode

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"
)

type RequestPacket struct {
	IVersion     int16             `tag:"1"  required:"true"`
	CPacketType  byte              `tag:"2"  required:"true"`
	IMessageType int32             `tag:"3"  required:"true"`
	IRequestId   int32             `tag:"4"  required:"true"`
	SServantName string            `tag:"5"  required:"true"`
	SFuncName    string            `tag:"6"  required:"true"`
	SBuffer      []byte            `tag:"7"  required:"true"`
	ITimeout     int32             `tag:"8"  required:"true"`
	Context      map[string]string `tag:"9"  required:"true"`
	Status       map[string]string `tag:"10"  required:"true"`
}

func (p *RequestPacket) Encode(buf *bytes.Buffer) error {
	var err error
	err = EncodeTagInt16Value(buf, p.IVersion, 1)
	if nil != err {
		return err
	}
	err = EncodeTagByteValue(buf, p.CPacketType, 2)
	if nil != err {
		return err
	}
	err = EncodeTagInt32Value(buf, p.IMessageType, 3)
	if nil != err {
		return err
	}
	err = EncodeTagInt32Value(buf, p.IRequestId, 4)
	if nil != err {
		return err
	}
	err = EncodeTagStringValue(buf, p.SServantName, 5)
	if nil != err {
		return err
	}
	err = EncodeTagStringValue(buf, p.SFuncName, 6)
	if nil != err {
		return err
	}
	err = EncodeTagBytesValue(buf, p.SBuffer, 7)
	if nil != err {
		return err
	}
	err = EncodeTagInt32Value(buf, p.ITimeout, 8)
	if nil != err {
		return err
	}
	err = EncodeTagMapValue(buf, p.Context, 9)
	if nil != err {
		return err
	}
	err = EncodeTagMapValue(buf, p.Status, 10)
	if nil != err {
		return err
	}
	return nil
}

func (p *RequestPacket) Decode(buf *bytes.Buffer) error {
	var err error
	err = DecodeTagInt16Value(buf, &p.IVersion, 1, true)
	if nil != err {
		return err
	}
	err = DecodeTagByteValue(buf, &p.CPacketType, 2, true)
	if nil != err {
		return err
	}
	err = DecodeTagInt32Value(buf, &p.IMessageType, 3, true)
	if nil != err {
		return err
	}
	err = DecodeTagInt32Value(buf, &p.IRequestId, 4, true)
	if nil != err {
		return err
	}
	err = DecodeTagStringValue(buf, &p.SServantName, 5, true)
	if nil != err {
		return err
	}
	err = DecodeTagStringValue(buf, &p.SFuncName, 6, true)
	if nil != err {
		return err
	}
	err = DecodeTagBytesValue(buf, &p.SBuffer, 7, true)
	if nil != err {
		return err
	}
	err = DecodeTagInt32Value(buf, &p.ITimeout, 8, true)
	if nil != err {
		return err
	}
	err = DecodeTagMapValue(buf, &p.Context, 9, true)
	if nil != err {
		return err
	}
	err = DecodeTagMapValue(buf, &p.Status, 10, true)
	if nil != err {
		return err
	}
	return err
}

func TestCodec(t *testing.T) {
	v1 := &RequestPacket{}
	v1.SFuncName = "helloww"
	v1.IMessageType = 12456
	v1.ITimeout = 10101
	v1.SServantName = "343242342$$"
	v1.Context = make(map[string]string)
	v1.Context["AAA"] = "BBB"
	v1.SBuffer = []byte("#######")
	var buf bytes.Buffer
	err := v1.Encode(&buf)
	if nil != err {
		t.Fatalf("###%v", err)
	}
	t.Logf("####%v", buf.Len())

	v2 := &RequestPacket{}
	err = v2.Decode(&buf)
	if nil != err {
		t.Fatalf("###%v", err)
	}
	t.Logf("####%v", v2)
}

func TestStringCodec(t *testing.T) {
	var err error

	// string ecode
	val1 := "hello, adam!"
	var buf1 bytes.Buffer
	err = Encode(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 string
	err = Decode(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val1 != val1 {
		t.Fatal(new_val1)
	}

	var buf2 bytes.Buffer
	err = EncodeTagStringValue(&buf2, val1, 0)
	if err != nil {
		t.Fatal(err)
	}

	var new_val2 string
	err = DecodeTagStringValue(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val2 != val1 {
		t.Fatal(new_val2)
	}
}

func TestBoolCodec(t *testing.T) {
	var err error

	// string ecode
	var val1 bool = true
	var buf1 bytes.Buffer
	err = Encode(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 bool
	err = Decode(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val1 != val1 {
		t.Fatal(new_val1)
	}

	var buf2 bytes.Buffer
	err = EncodeTagBoolValue(&buf2, val1, 0)
	if err != nil {
		t.Fatal(err)
	}

	var new_val2 bool
	err = DecodeTagBoolValue(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val2 != val1 {
		t.Fatal(new_val2)
	}
}

func TestCharCodec(t *testing.T) {
	var err error

	// string ecode
	var val1 byte = 'c'
	var buf1 bytes.Buffer
	err = Encode(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 byte
	err = Decode(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val1 != val1 {
		t.Fatal(new_val1)
	}

	var buf2 bytes.Buffer
	err = EncodeTagByteValue(&buf2, val1, 0)
	if err != nil {
		t.Fatal(err)
	}

	var new_val2 byte
	err = DecodeTagByteValue(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val2 != val1 {
		t.Fatal(new_val2)
	}
}

func TestFloat32Codec(t *testing.T) {
	var err error

	// string ecode
	var val1 float32 = 134234.2342
	var buf1 bytes.Buffer
	err = Encode(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 float32
	err = Decode(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val1 != val1 {
		t.Fatal(new_val1)
	}

	var buf2 bytes.Buffer
	err = EncodeTagFloat32Value(&buf2, val1, 0)
	if err != nil {
		t.Fatal(err)
	}

	var new_val2 float32
	err = DecodeTagFloat32Value(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val2 != val1 {
		t.Fatal(new_val2)
	}
}

func TestFloat64Codec(t *testing.T) {
	var err error

	// string ecode
	var val1 float64 = 134234.2342
	var buf1 bytes.Buffer
	err = Encode(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 float64
	err = Decode(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val1 != val1 {
		t.Fatal(new_val1)
	}

	var buf2 bytes.Buffer
	err = EncodeTagFloat64Value(&buf2, val1, 0)
	if err != nil {
		t.Fatal(err)
	}

	var new_val2 float64
	err = DecodeTagFloat64Value(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val2 != val1 {
		t.Fatal(new_val2)
	}
}

func TestInt8Codec(t *testing.T) {
	var err error

	// string ecode
	var val1 int8 = 127
	var buf1 bytes.Buffer
	err = Encode(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 int8
	err = Decode(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val1 != val1 {
		t.Fatal(new_val1)
	}

	var buf2 bytes.Buffer
	err = EncodeTagInt8Value(&buf2, val1, 0)
	if err != nil {
		t.Fatal(err)
	}

	var new_val2 int8
	err = DecodeTagInt8Value(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val2 != val1 {
		t.Fatal(new_val2)
	}
}

func TestUInt8Codec(t *testing.T) {
	var err error

	// string ecode
	var val1 uint8 = 127
	var buf1 bytes.Buffer
	err = Encode(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 uint8
	err = Decode(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val1 != val1 {
		t.Fatal(new_val1)
	}

	var buf2 bytes.Buffer
	err = EncodeTagInt8Value(&buf2, int8(val1), 0)
	if err != nil {
		t.Fatal(err)
	}

	var new_val2 int8
	err = DecodeTagInt8Value(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if uint8(new_val2) != val1 {
		t.Fatal(new_val2)
	}
}

func TestInt16Codec(t *testing.T) {
	var err error

	// string ecode
	var val1 int16 = 22345
	var buf1 bytes.Buffer
	err = Encode(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 int16
	err = Decode(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val1 != val1 {
		t.Fatal(new_val1)
	}

	var buf2 bytes.Buffer
	err = EncodeTagInt16Value(&buf2, val1, 0)
	if err != nil {
		t.Fatal(err)
	}

	var new_val2 int16
	err = DecodeTagInt16Value(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val2 != val1 {
		t.Fatal(new_val2)
	}
}

func TestUInt16Codec(t *testing.T) {
	var err error

	// string ecode
	var val1 uint16 = 22345
	var buf1 bytes.Buffer
	err = Encode(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 uint16
	err = Decode(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val1 != val1 {
		t.Fatal(new_val1)
	}

	var buf2 bytes.Buffer
	err = EncodeTagInt16Value(&buf2, int16(val1), 0)
	if err != nil {
		t.Fatal(err)
	}

	var new_val2 int16
	err = DecodeTagInt16Value(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if uint16(new_val2) != val1 {
		t.Fatal(new_val2)
	}
}

func TestInt32Codec(t *testing.T) {
	var err error

	// string ecode
	var val1 int32 = -1000
	var buf1 bytes.Buffer
	err = Encode(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 int32
	err = Decode(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val1 != val1 {
		t.Fatal(new_val1)
	}

	var buf2 bytes.Buffer
	err = EncodeTagInt32Value(&buf2, val1, 0)
	if err != nil {
		t.Fatal(err)
	}

	var new_val2 int32
	err = DecodeTagInt32Value(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val2 != val1 {
		t.Fatal(new_val2)
	}
}

func TestUInt32Codec(t *testing.T) {
	var err error

	// string ecode
	var val1 uint32 = 6553621
	var buf1 bytes.Buffer
	err = Encode(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 uint32
	err = Decode(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val1 != val1 {
		t.Fatal(new_val1)
	}

	var buf2 bytes.Buffer
	err = EncodeTagInt32Value(&buf2, int32(val1), 0)
	if err != nil {
		t.Fatal(err)
	}

	var new_val2 int32
	err = DecodeTagInt32Value(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if uint32(new_val2) != val1 {
		t.Fatal(new_val2)
	}
}

func TestInt64Codec(t *testing.T) {
	var err error

	// string ecode
	var val1 int64 = 41596504421
	var buf1 bytes.Buffer
	err = Encode(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 int64
	err = Decode(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val1 != val1 {
		t.Fatal(new_val1)
	}

	var buf2 bytes.Buffer
	err = EncodeTagInt64Value(&buf2, val1, 0)
	if err != nil {
		t.Fatal(err)
	}

	var new_val2 int64
	err = DecodeTagInt64Value(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val2 != val1 {
		t.Fatal(new_val2)
	}
}

func TestUInt64Codec(t *testing.T) {
	var err error

	// string ecode
	var val1 uint64 = 41596504421
	var buf1 bytes.Buffer
	err = Encode(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 uint64
	err = Decode(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if new_val1 != val1 {
		t.Fatal(new_val1)
	}

	var buf2 bytes.Buffer
	err = EncodeTagInt64Value(&buf2, int64(val1), 0)
	if err != nil {
		t.Fatal(err)
	}

	var new_val2 int64
	err = DecodeTagInt64Value(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if uint64(new_val2) != val1 {
		t.Fatal(new_val2)
	}
}

func TestMapCodec(t *testing.T) {
	var err error

	// string ecode
	val1 := map[string]map[string][]byte{
		"req_data": {
			"struct": []byte("adam_hello"),
		},
	}
	var buf1 bytes.Buffer
	err = Encode(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 map[string]map[string][]byte
	err = Decode(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	//	if reflect.DeepEqual(val1, new_val1) {
	//		t.Fatal(new_val1)
	//	}
}

func TestBytesCodec(t *testing.T) {
	var err error

	//	// bytes ecode
	val1 := []byte("adam_hello")
	var buf1 bytes.Buffer
	err = EncodeTagBytesValue(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 []byte
	err = DecodeTagBytesValue(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(val1, new_val1) {
		t.Fatal(new_val1)
	}

	// bytes ecode
	var buf2 bytes.Buffer
	err = Encode(&buf2, &val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf2.Bytes()))

	var new_val2 []byte
	err = Decode(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(val1, new_val2) {
		t.Fatal(new_val2)
	}
}

func TestInt16ArrayCodec(t *testing.T) {
	var err error

	// int16 array ecode
	val1 := []int16{233, 3234, 23223, 15}
	var buf1 bytes.Buffer
	err = EncodeTagVectorValue(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 []int16
	err = DecodeTagVectorValue(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(val1) != len(new_val1) {
		t.Fatal(len(new_val1))
	}
	for i, _ := range val1 {
		if val1[i] != new_val1[i] {
			t.Fatal("not equal")
		}
	}

	// int16 array ecode
	var buf2 bytes.Buffer
	err = Encode(&buf2, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf2.Bytes()))

	var new_val2 []int16
	err = Decode(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	// XXX: 用反射比较
	if len(val1) != len(new_val2) {
		t.Fatal(len(new_val2))
	}
	for i, _ := range val1 {
		if val1[i] != new_val2[i] {
			t.Fatal("not equal")
		}
	}
}

func TestUInt16ArrayCodec(t *testing.T) {
	var err error

	// uint16 array ecode
	val1 := []uint16{233, 3234, 23223, 15}
	var buf1 bytes.Buffer
	err = EncodeTagVectorValue(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 []uint16
	err = DecodeTagVectorValue(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(val1) != len(new_val1) {
		t.Fatal(len(new_val1))
	}
	for i, _ := range val1 {
		if val1[i] != new_val1[i] {
			t.Fatal("not equal")
		}
	}

	// uint16 array ecode
	var buf2 bytes.Buffer
	err = Encode(&buf2, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf2.Bytes()))

	var new_val2 []uint16
	err = Decode(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	// XXX: 用反射比较
	if len(val1) != len(new_val2) {
		t.Fatal(len(new_val2))
	}
	for i, _ := range val1 {
		if val1[i] != new_val2[i] {
			t.Fatal("not equal")
		}
	}
}

func TestUInt32ArrayCodec(t *testing.T) {
	var err error

	// uint32 array ecode
	val1 := []uint32{324234233, 32342234, 23223, 15}
	var buf1 bytes.Buffer
	err = EncodeTagVectorValue(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 []uint32
	err = DecodeTagVectorValue(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(val1) != len(new_val1) {
		t.Fatal(len(new_val1))
	}
	for i, _ := range val1 {
		if val1[i] != new_val1[i] {
			t.Fatal("not equal")
		}
	}

	// uint32 array ecode
	var buf2 bytes.Buffer
	err = Encode(&buf2, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf2.Bytes()))

	var new_val2 []uint32
	err = Decode(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	// XXX: 用反射比较
	if len(val1) != len(new_val2) {
		t.Fatal(len(new_val2))
	}
	for i, _ := range val1 {
		if val1[i] != new_val2[i] {
			t.Fatal("not equal")
		}
	}
}

func TestInt8ArrayCodec(t *testing.T) {
	var err error

	// int8 array ecode
	val1 := []int8{3, 34, 23, 15}
	var buf1 bytes.Buffer
	err = EncodeTagVectorValue(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 []int8
	err = DecodeTagVectorValue(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(val1) != len(new_val1) {
		t.Fatal(len(new_val1))
	}
	for i, _ := range val1 {
		if val1[i] != new_val1[i] {
			t.Fatal("not equal")
		}
	}

	// int8 array ecode
	var buf2 bytes.Buffer
	err = Encode(&buf2, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf2.Bytes()))

	var new_val2 []int8
	err = Decode(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	// XXX: 用反射比较
	if len(val1) != len(new_val2) {
		t.Fatal(len(new_val2))
	}
	for i, _ := range val1 {
		if val1[i] != new_val2[i] {
			t.Fatal("not equal")
		}
	}
}

func TestUInt8ArrayCodec(t *testing.T) {
	var err error

	// uint8 array ecode
	val1 := []uint8{3, 34, 23, 15}
	var buf1 bytes.Buffer
	// 注意这里必须要用EncodeTagBytesValue
	err = EncodeTagBytesValue(&buf1, val1, 0)
	//	err = EncodeTagVectorValue(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 []uint8
	err = DecodeTagVectorValue(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(val1) != len(new_val1) {
		t.Fatal(len(new_val1))
	}
	for i, _ := range val1 {
		if val1[i] != new_val1[i] {
			t.Fatal("not equal")
		}
	}

	// uint8 array ecode
	var buf2 bytes.Buffer
	err = Encode(&buf2, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf2.Bytes()))

	var new_val2 []uint8
	err = Decode(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	// XXX: 用反射比较
	if len(val1) != len(new_val2) {
		t.Fatal(len(new_val2))
	}
	for i, _ := range val1 {
		if val1[i] != new_val2[i] {
			t.Fatal("not equal")
		}
	}
}

func TestInt32ArrayCodec(t *testing.T) {
	var err error

	// int32 array ecode
	val1 := []int32{24234414, 42432, 4423423, 4234}
	var buf1 bytes.Buffer
	err = EncodeTagVectorValue(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 []int32
	err = DecodeTagVectorValue(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(val1) != len(new_val1) {
		t.Fatal(len(new_val1))
	}
	for i, _ := range val1 {
		if val1[i] != new_val1[i] {
			t.Fatal("not equal")
		}
	}

	// int32 array ecode
	var buf2 bytes.Buffer
	err = Encode(&buf2, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf2.Bytes()))

	var new_val2 []int32
	err = Decode(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	// XXX: 用反射比较
	if len(val1) != len(new_val2) {
		t.Fatal(len(new_val2))
	}
	for i, _ := range val1 {
		if val1[i] != new_val2[i] {
			t.Fatal("not equal")
		}
	}
}

func TestInt64ArrayCodec(t *testing.T) {
	var err error

	// int64 array ecode
	val1 := []int64{2234234234414, 422342323432, 4423423, 4234}
	//	val1 := []int64{}
	var buf1 bytes.Buffer
	err = EncodeTagVectorValue(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 []int64
	err = DecodeTagVectorValue(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(val1) != len(new_val1) {
		t.Fatal(len(new_val1))
	}
	for i, _ := range val1 {
		if val1[i] != new_val1[i] {
			t.Fatal("not equal")
		}
	}

	// int64 array ecode
	var buf2 bytes.Buffer
	err = Encode(&buf2, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf2.Bytes()))

	var new_val2 []int64
	err = Decode(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	// XXX: 用反射比较
	if len(val1) != len(new_val2) {
		t.Fatal(len(new_val2))
	}
	for i, _ := range val1 {
		if val1[i] != new_val2[i] {
			t.Fatal("not equal")
		}
	}
}

func TestUInt64ArrayCodec(t *testing.T) {
	var err error

	// uint64 array ecode
	val1 := []uint64{2234234234414, 422342323432, 4423423, 4234}
	//	val1 := []uint64{}
	var buf1 bytes.Buffer
	err = EncodeTagVectorValue(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 []uint64
	err = DecodeTagVectorValue(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(val1) != len(new_val1) {
		t.Fatal(len(new_val1))
	}
	for i, _ := range val1 {
		if val1[i] != new_val1[i] {
			t.Fatal("not equal")
		}
	}

	// uint64 array ecode
	var buf2 bytes.Buffer
	err = Encode(&buf2, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf2.Bytes()))

	var new_val2 []uint64
	err = Decode(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	// XXX: 用反射比较
	if len(val1) != len(new_val2) {
		t.Fatal(len(new_val2))
	}
	for i, _ := range val1 {
		if val1[i] != new_val2[i] {
			t.Fatal("not equal")
		}
	}
}

func TestStructArrayCodec(t *testing.T) {
	var err error

	// uint64 array ecode
	val1 := []RequestPacket{
		RequestPacket{
			SFuncName: "hello",
			ITimeout:  10101,
		},
	}
	var buf1 bytes.Buffer
	err = EncodeTagVectorValue(&buf1, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var new_val1 []RequestPacket
	err = DecodeTagVectorValue(&buf1, &new_val1, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(val1, new_val1) {
		t.Fatal(new_val1)
	}

	// uint64 array ecode
	var buf2 bytes.Buffer
	err = Encode(&buf2, val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex.EncodeToString(buf2.Bytes()))

	var new_val2 []RequestPacket
	err = Decode(&buf2, &new_val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(val1, new_val2) {
		t.Fatal(new_val2)
	}
}
