package sql

import (
	"github.com/foursking/ztgo/stat/metric"
)

var (
	_metricReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: "ztgo",
		Subsystem: "sql",
		Name:      "client_duration_ms",
		Help:      "sql client requests duration(ms).",
		Labels:    []string{"name", "addr", "command"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500},
	})
	_metricReqErr = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: "ztgo",
		Subsystem: "sql",
		Name:      "client_error_total",
		Help:      "sql client requests error count.",
		Labels:    []string{"name", "addr", "command", "error"},
	})
	_metricConnTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: "ztgo",
		Subsystem: "sql",
		Name:      "client_connection_total",
		Help:      "sql client connections total count.",
		Labels:    []string{"name", "addr", "state"},
	})
	_metricConnCurrent = metric.NewGaugeVec(&metric.GaugeVecOpts{
		Namespace: "ztgo",
		Subsystem: "sql",
		Name:      "client_connection_current",
		Help:      "sql client connections current.",
		Labels:    []string{"name", "addr", "state"},
	})
)
