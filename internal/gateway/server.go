package gateway

import (
	"context"
	"log"
	"os"

	"github.com/georgelopez7/grpc-project/pkg/logging"
	"github.com/gin-gonic/gin"
)

func StartGateway() {
	// TRACING
	tracerName := "gateway"
	shutdown := logging.InitTracer(tracerName)
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Fatalf("Failed to shut down tracer: %v", err)
		}
	}()

	// SERVER
	router := gin.Default()

	// ROUTES
	router.POST("/api/v1/payments", AddMetrics, PaymentRequestHandler)

	// PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// RUN SERVER
	router.Run(":" + port)
}
