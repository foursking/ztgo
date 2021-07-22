package innerpb

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"git.code.oa.com/qdgo/core/codec/innerpb"
	"git.code.oa.com/qdgo/core/log"
	"git.code.oa.com/qdgo/core/metadata"
	"git.code.oa.com/qdgo/core/net/udp"

	"github.com/davecgh/go-spew/spew"
	"github.com/opentracing/opentracing-go"
)

// ClientConfig Config innerpb client configuration
type ClientConfig struct {
	Addr         string        `json:"addr"`
	ReadTimeout  time.Duration `json:"read_timeout"`  // 读取超时时间
	WriteTimeout time.Duration `json:"write_timeout"` // 写入超时时间
}

// Client innerpb client
type Client struct {
	c *ClientConfig
}

// NewClient creates innerpb client
func NewClient(c *ClientConfig) *Client {
	if c.ReadTimeout == 0 {
		c.ReadTimeout = time.Second
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = time.Second
	}
	return &Client{
		c: c,
	}
}

// Response innerpb response
type Response struct {
	Data       []byte
	RemoteAddr string
}

// Do sends request
func (cli *Client) Do(ctx context.Context, cmd uint32, uin uint64, req []byte) (rsp *Response, err error) {
	return cli.DoMessage(ctx, innerpb.NewMessageByBytes(ctx, cmd, uin, req))
}

// DoMessage send innerpb.Message
func (cli *Client) DoMessage(ctx context.Context, m *innerpb.Message) (rsp *Response, err error) {
	op := fmt.Sprintf("%s,cmd: %d", metadata.InnerPB, m.Cmd)
	kfuin := m.Head.InnerHead.UinIds.Uint64Kfuin
	span, ctx := opentracing.StartSpanFromContext(ctx, op, opentracing.Tags{"kfuin": kfuin})
	defer span.Finish()
	udpCli, err := udp.NewClient(&udp.ClientConfig{
		RemoteAddr:   cli.c.Addr,
		ReadTimeout:  cli.c.ReadTimeout,
		WriteTimeout: cli.c.WriteTimeout,
	})
	if err != nil {
		log.Errorf("new udp client kfuin(%d) req(%s) error(%v)", kfuin, m, err)
		return
	}
	defer udpCli.Close()
	reqBuf, err := m.Marshal()
	if err != nil {
		log.Errorf("innerpb marshal msg(%s) error(%v)", spew.Sdump(m), err)
		return
	}
	now := time.Now()
	rspBuf, err := udpCli.SendAndReceive(ctx, reqBuf)
	_promCmd := strconv.FormatUint(uint64(m.Cmd), 10)
	_metricClientReqDur.Observe(int64(time.Since(now)/time.Millisecond), _promCmd)
	if m.Head != nil && m.Head.InnerHead != nil && m.Head.InnerHead.Uint32Result != nil {
		_metricClientReqCodeTotal.Inc(_promCmd, strconv.Itoa(int(*m.Head.InnerHead.Uint32Result)))
	}
	if err != nil {
		log.Errorf("innerpb SendAndReceive msg(%s) error(%v)", spew.Sdump(m), err)
		return
	}
	if err = m.Unmarshal(rspBuf); err != nil {
		log.Errorf("rspMsg unmarshal msg(%s) rsp(%s) error(%v)", spew.Sdump(m), rspBuf, err)
		return
	}
	rsp = &Response{
		Data:       m.Body,
		RemoteAddr: udpCli.RemoteAddr().String(),
	}
	return
}
