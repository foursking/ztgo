package event

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultEvent_Emit(t *testing.T) {
	e := newEvent()
	e.On("start", func(v ...interface{}) {
		if tmp, ok := v[0].(string); ok {
			assert.Equal(t, "hello", tmp)
		}
	})
	e.Emit(context.TODO(), "start", "hello")
}

func TestDefaultEvent_Off(t *testing.T) {
	e := newEvent()

	l := func(v ...interface{}) {}

	e.On("start", l)
	ls := e.Listeners("start")
	assert.Equal(t, 1, len(ls))

	e.Off("start", l)

	ls = e.Listeners("start")
	assert.Equal(t, 0, len(ls))
}

func TestDefaultEvent_AsyncOn(t *testing.T) {
	e := newEvent()
	e.AsyncOn("start", func(v ...interface{}) {
		time.Sleep(time.Second) // sleep 1 sec
	})

	st := time.Now()
	e.Emit(context.TODO(), "start", "hello")

	assert.Less(t, time.Now().Sub(st).Seconds(), float64(1))
}
