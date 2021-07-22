package addressing

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/foursking/ztgo/event"
	"github.com/foursking/ztgo/log"
)

const (
	_typeUnknown = iota
	_typeTCP
	_typeUDP
	_typeDNS
	_typeSock

	EventAddressingResult = "addressing.result"
)

var (
	ErrEmptyAddr       = errors.New("addressing: empty addr")
	ErrInvalidAddr     = errors.New("addressing: invalid addr")
	ErrUnknownAddrType = errors.New("addressing: unknown address type")
)

var (
	_lock sync.RWMutex
	_map  = make(map[string]*address)
)

type address struct {
	typ   int
	addr  string // ip:port
	proto string
}

// Addr .
type Addr struct {
	Addr  string
	Proto string

	beginTime time.Time
}

// IPPort .
func (a *Addr) IPPort() (string, int) {
	s := strings.Split(a.Addr, ":")
	if len(s) != 2 {
		return "", 0
	}
	port, _ := strconv.Atoi(s[1])
	return s[0], port
}

// Address .
func Address(addr string) (res *Addr, err error) {
	defer event.Emit(context.TODO(), EventAddressingResult, addr, res)
	if addr == "" {
		err = ErrEmptyAddr
		return
	}
	_lock.RLock()
	a, ok := _map[addr]
	_lock.RUnlock()
	if !ok {
		a, err = newAddress(addr)
		if err != nil {
			log.Errorf("addressing: new address(%s) error(%v)", addr, err)
			return
		}
	}
	res = &Addr{Addr: a.addr, Proto: a.proto}
	return
}

func newAddress(addr string) (*address, error) {
	a := new(address)
	addrs := strings.Split(addr, "://")
	if len(addrs) != 2 {
		return nil, ErrInvalidAddr
	}
	a.proto = addrs[0]
	switch a.proto {
	case "tcp", "ip":
		a.typ = _typeTCP
		a.addr = addrs[1]
		_map[addr] = a
	case "udp":
		a.typ = _typeUDP
		a.addr = addrs[1]
		_map[addr] = a
	case "dns":
		a.typ = _typeDNS
		a.addr = addrs[1]
	case "sock":
		a.typ = _typeSock
		a.addr = addrs[1]
	default:
		return nil, ErrUnknownAddrType
	}
	_lock.Lock()
	_map[addr] = a
	_lock.Unlock()
	return a, nil
}
