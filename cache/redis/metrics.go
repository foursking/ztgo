package redis

import (
	"git.code.oa.com/qdgo/core/stat/metric"
)

var (
	_metricReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: "qdgo",
		Subsystem: "redis",
		Name:      "client_duration_ms",
		Help:      "redis client requests duration(ms).",
		Labels:    []string{"command", "key"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500},
	})
	_metricReqErr = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: "qdgo",
		Subsystem: "redis",
		Name:      "client_error_total",
		Help:      "redis client requests error count.",
		Labels:    []string{"command", "key", "error"},
	})
	_metricHits = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: "qdgo",
		Subsystem: "redis",
		Name:      "client_hits_total",
		Help:      "redis client hits total.",
		Labels:    []string{"command", "key"},
	})
	_metricMisses = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: "qdgo",
		Subsystem: "redis",
		Name:      "client_misses_total",
		Help:      "redis client misses total.",
		Labels:    []string{"command", "key"},
	})
)
