package logging

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// METRICS
func InitMetricsEndpoint(port string) {
	// SET UP METRICS ENDPOINT
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":"+port, nil)
	}()
}
