package binv5

import "github.com/foursking/ztgo/net/udp"

// Config .
type Config struct {
	ServiceType uint8
	UDPClient   *udp.ClientConfig
}
