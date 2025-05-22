package main

import (
	"os"

	"github.com/georgelopez7/grpc-project/internal/validation"
	"github.com/georgelopez7/grpc-project/pkg/logging"
)

func init() {
	// METRICS
	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "2113"
	}
	logging.InitMetricsEndpoint(metricsPort)
}

func main() {
	validation.InitValidationServer()
}
