package udp

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"git.code.oa.com/qdgo/core/event"
	"git.code.oa.com/qdgo/core/log"

	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/util/addr"
)

const (
	EventHandleUDPError = "udp.handle_error"
)

func init() {
	cmd.DefaultServers["udp"] = NewMicroServer
}

type MicroServer struct {
	sync.Mutex
	opts         server.Options
	handler      server.Handler
	exit         chan chan error
	registerOnce sync.Once
	registered   bool
	rb           []byte // read buffer

	// chans to control read/handle/write logic
	handleCh chan *message
	writeCh  chan *message

	// waitGroups to take charge of the workers life cycle
	readwg   sync.WaitGroup
	handlewg sync.WaitGroup
	writewg  sync.WaitGroup

	// configs
	readBufSize    int // buffer size for ReadFromUDP
	readerNum      int // readLoop number
	handlerNum     int // handleLoop number
	writerNum      int // writeLoop number
	handleChanSize int // handler message chan size
	writeChanSize  int // writer message chan size
}

type message struct {
	data  []byte
	raddr *net.UDPAddr
}

func NewMicroServer(opts ...server.Option) server.Server {
	srv := MicroServer{
		opts: newOptions(opts...),
		exit: make(chan chan error),

		// default configs
		readBufSize:    65536,
		readerNum:      runtime.NumCPU(),
		handlerNum:     runtime.NumCPU(),
		writerNum:      runtime.NumCPU(),
		handleChanSize: 4096,
		writeChanSize:  1024,
	}
	srv.init()
	return &srv
}

func (s *MicroServer) init() {
	if v := s.opts.Context.Value(readBufKey{}); v != nil {
		if i, ok := v.(int); ok && i > 0 {
			s.readBufSize = i
		}
	}
	if v := s.opts.Context.Value(readerNumKey{}); v != nil {
		if i, ok := v.(int); ok && i > 0 {
			s.readerNum = i
		}
	}
	if v := s.opts.Context.Value(handlerNumKey{}); v != nil {
		if i, ok := v.(int); ok && i > 0 {
			s.handlerNum = i
		}
	}
	if v := s.opts.Context.Value(writerNumKey{}); v != nil {
		if i, ok := v.(int); ok && i > 0 {
			s.writerNum = i
		}
	}
	if v := s.opts.Context.Value(handleChanSizeKey{}); v != nil {
		if i, ok := v.(int); ok && i > 0 {
			s.handleChanSize = i
		}
	}
	if v := s.opts.Context.Value(writeChanSizeKey{}); v != nil {
		if i, ok := v.(int); ok && i > 0 {
			s.writeChanSize = i
		}
	}
	if v := s.opts.Context.Value(handleFuncKey{}); v != nil {
		if fn, ok := v.(UDPHandleFunc); ok {
			if err := s.Handle(s.NewHandler(fn)); err != nil {
				log.Fatalf("UDP server: new handler error(%v)", err)
			}
		}
	}
	s.rb = make([]byte, s.readBufSize)
	s.handleCh = make(chan *message, s.handleChanSize)
	s.writeCh = make(chan *message, s.writeChanSize)
}

func (s *MicroServer) Options() server.Options {
	s.Lock()
	opts := s.opts
	s.Unlock()
	return opts
}

func (s *MicroServer) Init(opts ...server.Option) error {
	s.Lock()
	for _, o := range opts {
		o(&s.opts)
	}
	s.Unlock()
	return nil
}

func (s *MicroServer) Handle(handler server.Handler) error {
	if _, ok := handler.Handler().(UDPHandleFunc); !ok {
		return nil
	}
	s.Lock()
	s.handler = handler
	s.Unlock()
	return nil
}

func (s *MicroServer) NewHandler(handler interface{}, opts ...server.HandlerOption) server.Handler {
	options := server.HandlerOptions{
		Metadata: make(map[string]map[string]string),
	}
	for _, o := range opts {
		o(&options)
	}
	var eps []*registry.Endpoint
	if !options.Internal {
		for name, metadata := range options.Metadata {
			eps = append(eps, &registry.Endpoint{
				Name:     name,
				Metadata: metadata,
			})
		}
	}
	return &udpHandler{
		eps:  eps,
		hd:   handler,
		opts: options,
	}
}

func (s *MicroServer) NewSubscriber(string, interface{}, ...server.SubscriberOption) server.Subscriber {
	return nil
}

func (s *MicroServer) Subscribe(server.Subscriber) error {
	return nil
}

func (s *MicroServer) Start() error {
	s.Lock()
	opts := s.opts
	s.Unlock()
	handler, ok := s.handler.Handler().(UDPHandleFunc)
	if !ok || handler == nil {
		log.Fatalf("UDP server: invalid handler(%+v)", handler)
	}
	address, err := net.ResolveUDPAddr("udp", opts.Address)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		return err
	}
	if err = s.Register(); err != nil {
		return err
	}
	log.Infof("UDP Server listen on %s", opts.Address)
	s.readwg.Add(s.readerNum)
	for i := 0; i < s.readerNum; i++ {
		go s.readLoop(conn)
	}
	s.handlewg.Add(s.handlerNum)
	for i := 0; i < s.handlerNum; i++ {
		go s.handleLoop(handler)
	}
	s.writewg.Add(s.writerNum)
	for i := 0; i < s.writerNum; i++ {
		go s.writeLoop(conn)
	}
	go func() {
		t := new(time.Ticker)
		if opts.RegisterInterval > time.Duration(0) {
			t = time.NewTicker(opts.RegisterInterval)
		}
		var ch chan error
	Loop:
		for {
			select {
			// register self on interval
			case <-t.C:
				if err := s.Register(); err != nil {
					log.Errorf("UDP server register error(%v)", err)
				}
			// wait for exit
			case ch = <-s.exit:
				break Loop
			}
		}
		if err = conn.Close(); err != nil {
			log.Errorf("UDP server: close connection error(%v)", err)
		}
		s.readwg.Wait()
		close(s.handleCh)
		s.handlewg.Wait()
		close(s.writeCh)
		s.writewg.Wait()
		if err = s.Deregister(); err != nil {
			log.Errorf("UDP server: deregister error(%v)", err)
		}
		if err = opts.Broker.Disconnect(); err != nil {
			log.Errorf("UDP server: broker disconnect error(%v)", err)
		}
		ch <- err
	}()
	return nil
}

func (s *MicroServer) readLoop(conn *net.UDPConn) {
	defer s.readwg.Done()
	var tempDelay time.Duration
	for {
		size, raddr, err := conn.ReadFromUDP(s.rb)
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			// 连接被关闭
			if e, ok := err.(*net.OpError); ok {
				if strings.Contains(e.Error(), "use of closed network connection") {
					log.Debugf("UDP server: connection has been closed")
					return
				}
			}
			log.Errorf("UDP server: read error(%v)", err)
			return
		}
		tempDelay = 0
		data := make([]byte, size)
		copy(data, s.rb[:size])
		s.handleCh <- &message{
			data:  data,
			raddr: raddr,
		}
	}
}

func (s *MicroServer) handleLoop(handler UDPHandleFunc) {
	defer s.handlewg.Done()
	for msg := range s.handleCh {
		rsp, err := handler(msg.raddr, msg.data)
		if err != nil {
			event.Emit(context.TODO(), EventHandleUDPError, msg.raddr, msg.data, err)
			continue
		}
		data := make([]byte, len(rsp))
		copy(data, rsp)
		s.writeCh <- &message{
			data:  data,
			raddr: msg.raddr,
		}
	}
	log.Infof("UDP server: handleCh has been closed")
	return
}

func (s *MicroServer) writeLoop(conn *net.UDPConn) {
	defer s.writewg.Done()
	var err error
	for msg := range s.writeCh {
		if _, err = conn.WriteToUDP(msg.data, msg.raddr); err != nil {
			log.Errorf("UDP server: write to (%s) error(%v)", msg.raddr, err)
		}
	}
	log.Infof("UDP server: writeCh has been closed")
	return
}

func (s *MicroServer) Register() error {
	s.Lock()
	opts := s.opts
	eps := s.handler.Endpoints()
	s.Unlock()
	service := serviceDef(opts)
	service.Endpoints = eps
	rOpts := []registry.RegisterOption{
		registry.RegisterTTL(opts.RegisterTTL),
	}
	s.registerOnce.Do(func() {
		log.Infof("registering node: %s", opts.Name+"-"+opts.Id)
	})
	if err := opts.Registry.Register(service, rOpts...); err != nil {
		return err
	}
	s.Lock()
	if s.registered {
		s.Unlock()
		return nil
	}
	s.registered = true
	s.Unlock()
	return nil
}

func (s *MicroServer) Deregister() error {
	s.Lock()
	opts := s.opts
	s.Unlock()

	log.Infof("deregistering node: %s", opts.Name+"-"+opts.Id)

	service := serviceDef(opts)
	if err := opts.Registry.Deregister(service); err != nil {
		return err
	}
	s.Lock()
	if !s.registered {
		s.Unlock()
		return nil
	}
	s.registered = false
	s.Unlock()
	return nil
}

func (s *MicroServer) Stop() error {
	ch := make(chan error)
	s.exit <- ch
	return <-ch
}

func (s *MicroServer) String() string {
	return "udp"
}

func serviceDef(opts server.Options) *registry.Service {
	var advt, host string
	var port int
	if len(opts.Advertise) > 0 {
		advt = opts.Advertise
	} else {
		advt = opts.Address
	}
	parts := strings.Split(advt, ":")
	if len(parts) > 1 {
		host = strings.Join(parts[:len(parts)-1], ":")
		port, _ = strconv.Atoi(parts[len(parts)-1])
	} else {
		host = parts[0]
	}
	address, err := addr.Extract(host)
	if err != nil {
		address = host
	}
	node := &registry.Node{
		Id:       opts.Name + "-" + opts.Id,
		Address:  fmt.Sprintf("%s:%d", address, port),
		Metadata: opts.Metadata,
	}
	node.Metadata["server"] = "udp"
	node.Metadata["broker"] = opts.Broker.String()
	node.Metadata["registry"] = opts.Registry.String()
	node.Metadata["protocol"] = "udp"
	return &registry.Service{
		Name:    opts.Name,
		Version: opts.Version,
		Nodes:   []*registry.Node{node},
	}
}
