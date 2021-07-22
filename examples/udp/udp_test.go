package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/foursking/ztgo/net/udp"
)

func TestUdpClient(t *testing.T) {
	cc := &udp.ClientConfig{
		RemoteAddr:   "ip://127.0.0.1:8888",
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}
	for {
		conn, err := udp.NewClient(cc)
		if err != nil {
			t.Error(err)
		}
		sendb := []byte(time.Now().String())
		b, err := conn.SendAndReceive(context.Background(), sendb)
		fmt.Printf("send(%s) recv(%s) laddr(%v) raddr(%v) error(%v) \n", sendb, b, conn.LocalAddr(), conn.RemoteAddr(), err)
		conn.Close()
		time.Sleep(time.Microsecond)
		//break
	}
}
