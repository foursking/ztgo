package grpc

import (
	"context"
	"time"

	"github.com/foursking/ztgo/config/env"
	qdmd "github.com/foursking/ztgo/metadata"
	"github.com/foursking/ztgo/stat/metric"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/server"
)

var (
	_metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: "ztgo",
		Subsystem: "grpc",
		Name:      "server_duration_ms",
		Help:      "grpc server requests duration(ms).",
		Labels:    []string{"method", "caller"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500},
	})
	_metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: "ztgo",
		Subsystem: "grpc",
		Name:      "server_code_total",
		Help:      "grpc server requests code count.",
		Labels:    []string{"method", "caller", "code"},
	})
	_metricClientReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: "ztgo",
		Subsystem: "grpc",
		Name:      "client_duration_ms",
		Help:      "grpc client requests duration(ms).",
		Labels:    []string{"method"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500},
	})
	_metricClientReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: "ztgo",
		Subsystem: "grpc",
		Name:      "client_code_total",
		Help:      "grpc client requests code count.",
		Labels:    []string{"method", "code"},
	})
)

func NewServerMetricHandlerWrapper() server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			now := time.Now()
			method := req.Method()
			caller, ok := metadata.Get(ctx, qdmd.Caller)
			if !ok || caller == "" {
				caller = qdmd.Unknown
			}
			err := fn(ctx, req, rsp)
			if err == nil {
				_metricServerReqCodeTotal.Inc(method, caller, qdmd.Success)
			} else {
				_metricServerReqCodeTotal.Inc(method, caller, qdmd.Failure)
			}
			_metricServerReqDur.Observe(int64(time.Since(now)/time.Millisecond), method, caller)
			return err
		}
	}
}

func NewClientMetricHandlerWrapper() client.CallWrapper {
	return func(cf client.CallFunc) client.CallFunc {
		return func(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {
			now := time.Now()
			method := req.Method()
			ctx = metadata.MergeContext(ctx, metadata.Metadata{qdmd.Caller: env.AppName}, true)
			err := cf(ctx, node, req, rsp, opts)
			if err == nil {
				_metricClientReqCodeTotal.Inc(method, qdmd.Success)
			} else {
				_metricClientReqCodeTotal.Inc(method, qdmd.Failure)
			}
			_metricClientReqDur.Observe(int64(time.Since(now)/time.Millisecond), method)
			return err
		}
	}
}
