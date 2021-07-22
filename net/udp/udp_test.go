package udp

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClient(t *testing.T) {
	Convey("ecode udp client", t, func() {
		// Server
		sc := &ServerConfig{
			Addr:         "ip://127.0.0.1:9999",
			WriteTimeout: time.Second,
		}
		go NewServer(sc, func(conn *Conn) {
			defer conn.Close()
			for {
				data, err := conn.Receive(context.Background())
				fmt.Printf("server receive data(%s) error(%v)", data, err)
				t.Logf("receive data (%s)", data)
				if len(data) > 0 {
					if err := conn.Send(context.Background(), append([]byte("> "), data...)); err != nil {
						t.Error(err)
					}
				}
				if err != nil {
					t.Error(err)
				}
			}
		}).Run()

		time.Sleep(time.Second)

		// Client
		cc := &ClientConfig{
			RemoteAddr:   "ip://127.0.0.1:9999",
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
		}
		for {
			conn, err := NewClient(cc)
			if err != nil {
				t.Error(err)
			}
			fmt.Printf("\nsend a request \n")
			if err = conn.Send(context.Background(), []byte(time.Now().String())); err != nil {
				t.Error(err)
			}
			if b, err := conn.SendAndReceive(context.Background(), []byte(time.Now().String())); err == nil {
				fmt.Println("wwwwww", string(b), conn.LocalAddr(), conn.RemoteAddr())
			} else {
				t.Error(err)
			}
			conn.Close()
			time.Sleep(time.Second)
		}
	})
}
