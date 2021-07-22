package udp

import (
	"net"

	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/server"
)

type UDPHandleFunc func(*net.UDPAddr, []byte) ([]byte, error)

type udpHandler struct {
	opts server.HandlerOptions
	eps  []*registry.Endpoint
	hd   interface{}
}

func (h *udpHandler) Name() string {
	return "handler"
}

func (h *udpHandler) Handler() interface{} {
	return h.hd
}

func (h *udpHandler) Endpoints() []*registry.Endpoint {
	return h.eps
}

func (h *udpHandler) Options() server.HandlerOptions {
	return h.opts
}
