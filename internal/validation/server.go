package validation

import (
	"context"
	"log"
	"net"

	"github.com/georgelopez7/grpc-project/api/proto/paymentpb"
	"github.com/georgelopez7/grpc-project/pkg/logging"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type ValidationServer struct {
	paymentpb.UnimplementedValidationServiceServer
}

func InitValidationServer() {
	// TRACING
	tracerName := "validation-service"
	shutdown := logging.InitTracer(tracerName)
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Fatalf("Failed to shutdown tracer: %v", err)
		}
	}()

	// LISTEN
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// GRPC SERVER
	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()), // IMPORTANT: PROPAGATE TRACE CONTEXT

	)
	paymentpb.RegisterValidationServiceServer(s, &ValidationServer{})

	log.Println("🔍 Validator service listening on :50053")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
