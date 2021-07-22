package core

import (
	"git.code.oa.com/qdgo/core/addressing"
	"git.code.oa.com/qdgo/core/event"
	"git.code.oa.com/qdgo/core/log"
)

func init() {
	event.On(addressing.EventAddressingResult, logAddressing)
}

func logAddressing(v ...interface{}) {
	if len(v) < 2 {
		return
	}
	var (
		ok   bool
		str  string
		addr *addressing.Addr
	)
	if str, ok = v[0].(string); !ok {
		return
	}
	if addr, ok = v[1].(*addressing.Addr); !ok {
		return
	}
	log.Debugf("addressing addr: %s target: %#v", str, addr)
}
