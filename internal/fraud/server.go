package fraud

import (
	"context"
	"log"
	"net"

	"github.com/georgelopez7/grpc-project/api/proto/paymentpb"
	"github.com/georgelopez7/grpc-project/pkg/logging"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type FraudServer struct {
	paymentpb.UnimplementedFraudServiceServer
}

func InitFraudServer() {
	// TRACING
	tracerName := "fraud-service"
	shutdown := logging.InitTracer(tracerName)
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Fatalf("Failed to shutdown tracer: %v", err)
		}
	}()

	// LISTEN
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// GRPC SERVER
	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()), // IMPORTANT: PROPAGATE TRACE CONTEXT
	)
	paymentpb.RegisterFraudServiceServer(s, &FraudServer{})

	log.Println("💀 Fraud service listening on :50052")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
