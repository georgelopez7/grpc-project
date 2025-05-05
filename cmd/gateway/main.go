package main

import (
	"os"

	"github.com/georgelopez7/grpc-project/internal/gateway"
	"github.com/georgelopez7/grpc-project/pkg/logging"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	// --- METRICS ---
	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "2114"
	}

	prometheus.MustRegister(gateway.PaymentRequestsCount)
	logging.InitMetricsEndpoint(metricsPort)
}

func main() {
	gateway.StartGateway()
}
