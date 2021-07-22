package stat

import (
	"context"
	"time"

	"github.com/foursking/ztgo/log"

	"github.com/micro/go-micro/v2/client"
)

var DefaultLatencyFn = func(ctx context.Context, req client.Request, latency time.Duration) {
	log.Infof("%s->%s latency: %s", req.Service(), req.Method(), latency)
}

type latencyKey struct{}

type LatencyFn func(ctx context.Context, req client.Request, latency time.Duration)

type latencyWrapper struct {
	client.Client
	fn LatencyFn
}

func NewDefaultLatencyWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &latencyWrapper{Client: c, fn: DefaultLatencyFn}
	}
}

func NewLatencyWrapper(fn LatencyFn) client.Wrapper {
	return func(c client.Client) client.Client {
		return &latencyWrapper{Client: c, fn: fn}
	}
}

func (w *latencyWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	st := time.Now()
	defer func() {
		d := time.Now().Sub(st)
		w.fn(ctx, req, d)
		// inject to context
		ctx = context.WithValue(ctx, latencyKey{}, d)
	}()

	return w.Client.Call(ctx, req, rsp, opts...)
}
