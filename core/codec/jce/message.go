package jce

import "bytes"

// Message jce消息体接口， 类似protobuf
type Message interface {
	Encode(buf *bytes.Buffer) error
	Decode(buf *bytes.Buffer) error
	ClassName() string
}

// Marshal jce 打包函数 与标准包 json xml proto 保持一致
func Marshal(msg Message) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := msg.Encode(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal jce 解包
func Unmarshal(data []byte, msg Message) error {
	buf := bytes.NewBuffer(data)
	return msg.Decode(buf)
}
