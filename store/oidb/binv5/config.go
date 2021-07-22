package binv5

import "git.code.oa.com/qdgo/core/net/udp"

// Config .
type Config struct {
	ServiceType uint8
	UDPClient   *udp.ClientConfig
}
