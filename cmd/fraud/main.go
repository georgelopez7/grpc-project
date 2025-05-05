package main

import (
	"os"

	"github.com/georgelopez7/grpc-project/internal/fraud"
	"github.com/georgelopez7/grpc-project/pkg/logging"
)

func init() {
	// --- METRICS ---
	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "2112"
	}

	logging.InitMetricsEndpoint(metricsPort)
}

func main() {
	fraud.InitFraudServer()
}
