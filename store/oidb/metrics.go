package oidb

import "git.code.oa.com/qdgo/core/stat/metric"

var (
	MetricReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: "qdgo",
		Subsystem: "oidb",
		Name:      "client_duration_ms",
		Help:      "oidb client requests duration(ms).",
		Labels:    []string{"name", "addr", "command"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500},
	})
	MetricReqErr = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: "qdgo",
		Subsystem: "oidb",
		Name:      "client_error_total",
		Help:      "oidb client requests error count.",
		Labels:    []string{"name", "addr", "command", "error"},
	})
)
