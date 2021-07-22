package udp

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/foursking/ztgo/addressing"
	"github.com/foursking/ztgo/log"
	"github.com/foursking/ztgo/util/backoff"
)

const (
	defaultReadTimeout  = 200 * time.Millisecond
	defaultWriteTimeout = 200 * time.Millisecond
)

// ClientConfig .
type ClientConfig struct {
	RemoteAddr   string        // 远程地址
	LocalAddrs   []string      // 本地地址（可不指定，默认系统随机分配套接字端口）
	ReadTimeout  time.Duration // 读取超时时间
	WriteTimeout time.Duration // 写入超时时间
	RetryTimes   int           // 失败重试次数
}

// Client .
type Client struct {
	*Conn
}

// Conn 封装的链接对象
type Conn struct {
	conn         *net.UDPConn // 底层链接对象
	readTimeout  time.Duration
	writeTimeout time.Duration
	raddr        *net.UDPAddr // 远程地址
	retryTimes   int
}

const defaultReadBufferSize = 65536 // 默认数据读取缓冲区大小 64k

// NewClient .
func NewClient(c *ClientConfig) (*Client, error) {
	fixClientConfig(c)
	addr, err := addressing.Address(c.RemoteAddr)
	if err != nil {
		log.Errorf("udp client addressing addr(%s) error(%v)", c.RemoteAddr, err)
		return nil, err
	}
	//defer addr.ReportL5(&err)
	conn, err := newConn(addr.Addr, c.LocalAddrs...)
	if err != nil {
		log.Errorf("udp: NewNetConn(%s,%v) error(%v)", c.RemoteAddr, c.LocalAddrs, err)
		return nil, err
	}
	return &Client{Conn: &Conn{
		conn:         conn,
		readTimeout:  c.ReadTimeout,
		writeTimeout: c.WriteTimeout,
		retryTimes:   c.RetryTimes,
	}}, nil
}

func fixClientConfig(c *ClientConfig) {
	if c.ReadTimeout == 0 {
		c.ReadTimeout = defaultReadTimeout
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = defaultWriteTimeout
	}
}

// newConn 创建标准库UDP链接操作对象
func newConn(raddr string, laddr ...string) (*net.UDPConn, error) {
	var (
		err    error
		ra, la *net.UDPAddr
	)
	ra, err = net.ResolveUDPAddr("udp", raddr)
	if err != nil {
		return nil, err
	}
	if len(laddr) > 0 {
		la, err = net.ResolveUDPAddr("udp", laddr[0])
		if err != nil {
			return nil, err
		}
	}
	conn, err := net.DialUDP("udp", la, ra)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Send .
func (c *Conn) Send(ctx context.Context, data []byte) (err error) {
	var (
		size, length int
		doneChan     = make(chan error, 1)
		bk           = &backoff.Backoff{}
	)
	go func() {
		if c.writeTimeout > 0 {
			_ = c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
		}
		for {
			if c.raddr != nil {
				size, err = c.conn.WriteToUDP(data, c.raddr)
			} else {
				size, err = c.conn.Write(data)
			}
			if err != nil {
				// 链接已关闭，或者出错了但是不需要重试，或者重试次数超了，就退出
				if err == io.EOF || c.retryTimes == 0 || bk.Attempt() > c.retryTimes {
					doneChan <- err
					return
				}
				time.Sleep(bk.Duration())
				continue
			}
			length += size
			if length == len(data) {
				doneChan <- nil
				return
			}
		}
	}()
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-doneChan:
	}
	if err != nil {
		log.Errorf("udp: send error(%v)", err)
	}
	return
}

// Receive 接收数据
func (c *Conn) Receive(ctx context.Context) ([]byte, error) {
	var (
		err      error
		size     int          // 读取长度
		raddr    *net.UDPAddr // 当前读取的远程地址
		doneChan = make(chan error, 1)
		bk       = &backoff.Backoff{}
		buffer   = make([]byte, defaultReadBufferSize) // 读取缓冲区
	)
	go func() {
		for {
			if c.readTimeout > 0 {
				_ = c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
			}
			size, raddr, err = c.conn.ReadFromUDP(buffer)
			if err != nil {
				// 链接已关闭
				if err == io.EOF {
					doneChan <- err
					return
				}
				if bk.Attempt() >= c.retryTimes {
					doneChan <- err
					return
				}
				time.Sleep(bk.Duration())
				// 判断数据是否全部读取完毕(由于超时机制的存在，获取的数据完整性不可靠)
				if isTimeout(err) {
					_ = c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
				}
				continue
			}
			c.raddr = raddr
			doneChan <- nil
			return
		}
	}()
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-doneChan:
	}
	if err != nil {
		log.Errorf("udp: receive error(%v)", err)
	}
	return buffer[:size], err
}

// SendAndReceive 发送数据并等待接收返回数据
func (c *Conn) SendAndReceive(ctx context.Context, data []byte) ([]byte, error) {
	err := c.Send(ctx, data)
	if err != nil {
		return nil, err
	}
	return c.Receive(ctx)
}

// LocalAddr .
func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// RemoteAddr .
func (c *Conn) RemoteAddr() net.Addr {
	return c.raddr
}

// Close .
func (c *Conn) Close() error {
	return c.conn.Close()
}

// 判断是否是超时错误
func isTimeout(err error) bool {
	if err == nil {
		return false
	}
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}
	return false
}
