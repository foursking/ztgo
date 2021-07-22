package event

import (
	"context"
	"reflect"
	"sync"
)

type Event interface {
	Emit(ctx context.Context, event string, v ...interface{})
	On(event string, l ListenerFunc)
	AsyncOn(event string, l ListenerFunc)
	Off(event string, l ListenerFunc)
	Listeners(event string) []Listener
	AllListeners() map[string][]Listener
}

type ListenerFunc func(v ...interface{})

type Listener struct {
	Func  ListenerFunc
	Async bool
}

type event struct {
	sync.Mutex
	listeners map[string][]Listener
}

var DefaultEvent = newEvent()

func newEvent() Event {
	return &event{listeners: make(map[string][]Listener)}
}

func (e *event) Emit(ctx context.Context, event string, v ...interface{}) {
	ls, ok := e.listeners[event]
	if !ok {
		return
	}
	for _, l := range ls {
		if l.Async {
			go l.Func(v...)
		} else {
			l.Func(v...)
		}
	}
}

func (e *event) On(event string, l ListenerFunc) {
	if l == nil {
		return
	}
	e.Lock()
	defer e.Unlock()
	nl := Listener{Func: l}
	if ls, ok := e.listeners[event]; ok {
		ls = append(ls, nl)
	} else {
		e.listeners[event] = []Listener{nl}
	}
}

func (e *event) AsyncOn(event string, l ListenerFunc) {
	if l == nil {
		return
	}
	e.Lock()
	defer e.Unlock()
	nl := Listener{Func: l, Async: true}
	if ls, ok := e.listeners[event]; ok {
		ls = append(ls, nl)
	} else {
		e.listeners[event] = []Listener{nl}
	}
}

func (e *event) Off(event string, l ListenerFunc) {
	if l == nil {
		return
	}
	e.Lock()
	defer e.Unlock()
	if ls, ok := e.listeners[event]; ok {
		lv := reflect.ValueOf(l)
		for i, v := range ls {
			if reflect.ValueOf(v.Func) == lv {
				e.listeners[event] = append(e.listeners[event][:i], e.listeners[event][i+1:]...)
			}
		}
	}
}

func (e *event) Listeners(event string) []Listener {
	if ls, ok := e.listeners[event]; ok {
		return ls
	}
	return nil
}

func (e *event) AllListeners() map[string][]Listener {
	return e.listeners
}

func Emit(ctx context.Context, event string, v ...interface{}) {
	DefaultEvent.Emit(ctx, event, v...)
}

func On(event string, l ListenerFunc) {
	DefaultEvent.On(event, l)
}

func AsyncOn(event string, l ListenerFunc) {
	DefaultEvent.AsyncOn(event, l)
}

func Off(event string, l ListenerFunc) {
	DefaultEvent.Off(event, l)
}

func Listeners(event string) []Listener {
	return DefaultEvent.Listeners(event)
}

func AllListeners() map[string][]Listener {
	return DefaultEvent.AllListeners()
}
