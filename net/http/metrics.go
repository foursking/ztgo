package http

import (
	"git.code.oa.com/qdgo/core/stat/metric"
)

var (
	_metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: "qdgo",
		Subsystem: "http",
		Name:      "server_duration_ms",
		Help:      "http server requests duration(ms).",
		Labels:    []string{"path", "caller", "method"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500},
	})
	_metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: "qdgo",
		Subsystem: "http",
		Name:      "server_code_total",
		Help:      "http server requests error count.",
		Labels:    []string{"path", "caller", "method", "code"},
	})
	_metricClientReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: "qdgo",
		Subsystem: "http",
		Name:      "client_duration_ms",
		Help:      "http client requests duration(ms).",
		Labels:    []string{"path", "method"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500},
	})
	_metricClientReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: "qdgo",
		Subsystem: "http",
		Name:      "client_code_total",
		Help:      "http client requests code count.",
		Labels:    []string{"path", "method", "code"},
	})
)
