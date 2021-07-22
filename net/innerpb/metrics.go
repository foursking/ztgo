package innerpb

import "github.com/foursking/ztgo/stat/metric"

var (
	_metricClientReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: "ztgo",
		Subsystem: "innerpb",
		Name:      "client_duration_ms",
		Help:      "innerpb client requests duration(ms).",
		Labels:    []string{"command"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500},
	})
	_metricClientReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: "ztgo",
		Subsystem: "innerpb",
		Name:      "client_code_total",
		Help:      "innerpb client requests code count.",
		Labels:    []string{"command", "code"},
	})
)
