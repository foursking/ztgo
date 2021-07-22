package main

import (
	"net"

	"git.code.oa.com/qdgo/core"
	"git.code.oa.com/qdgo/core/net/udp"

	"github.com/micro/go-micro/v2"
)

func main() {
	var handleFunc udp.UDPHandleFunc = func(raddr *net.UDPAddr, bs []byte) ([]byte, error) {
		return append(bs, 'w'), nil
	}
	us := udp.NewMicroServer(udp.HandleFunc(handleFunc))
	srv := core.NewUDPService(micro.Server(us))
	if err := srv.Run(); err != nil {
		panic(err)
	}
}
