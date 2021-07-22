package main

import (
	"net"

	core "github.com/foursking/ztgo"
	"github.com/foursking/ztgo/net/udp"

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
