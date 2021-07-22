package udp

import (
	"errors"
	"net"
	"time"

	"github.com/foursking/ztgo/addressing"
	"github.com/foursking/ztgo/log"
)

// ServerConfig .
type ServerConfig struct {
	Addr         string
	WriteTimeout time.Duration // 写入超时时间
	RetryTimes   int           // 失败重试次数
}

// Server udp server
type Server struct {
	c       *ServerConfig
	handler func(*Conn)
}

// NewServer 创建一个udp server对象，并且可以选择指定一个单例名字
func NewServer(c *ServerConfig, handler func(*Conn)) *Server {
	fixServerConfig(c)
	return &Server{c, handler}
}

func fixServerConfig(c *ServerConfig) {
	if c.WriteTimeout == 0 {
		c.WriteTimeout = defaultWriteTimeout
	}
}

// Run .
func (s *Server) Run() (err error) {
	if s.handler == nil {
		err = errors.New("start running failed: socket handler not defined")
		log.Errorf("udp: %v", err)
		return
	}
	a, err := addressing.Address(s.c.Addr)
	if err != nil {
		log.Errorf("udp server addressing addr(%s) error(%v)", s.c.Addr, err)
		return
	}
	defer a.ReportL5(&err)
	addr, err := net.ResolveUDPAddr("udp", a.Addr)
	if err != nil {
		log.Errorf("udp: resolveUDPAddr(%s) error(%v)", a.Addr, err)
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Errorf("udp: listenUDP error(%v)", err)
		return
	}
	for {
		s.handler(&Conn{
			conn:         conn,
			writeTimeout: s.c.WriteTimeout,
			retryTimes:   s.c.RetryTimes,
		})
	}
}
