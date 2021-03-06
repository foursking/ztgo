package metric

import (
	"net/http"

	"github.com/foursking/ztgo/config/env"
	"github.com/foursking/ztgo/log"

	ph "github.com/prometheus/client_golang/prometheus/promhttp"
)

// BootPrometheus boots prometheus
func BootPrometheus() {
	http.Handle("/metrics", ph.Handler())
	go func() {
		//log.Infof("prometheus listening on: http://%s/metrics", env.PrometheusAddr)
		err := http.ListenAndServe(env.PrometheusAddr, nil)
		if err != nil {
			log.Fatalf("prometheus: ListenAndServe: %v", err)
		}
	}()
}
