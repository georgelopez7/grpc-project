package connections

import (
	"log"
	"os"

	"github.com/georgelopez7/grpc-project/api/proto/paymentpb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var FraudClient paymentpb.FraudServiceClient

func ConnectFraudService() {
	// SERVICE ADDRESS
	address := os.Getenv("FRAUD_SERVICE_ADDR")
	if address == "" {
		address = "localhost:50052"
	}

	// GRPC CONNECTION OPTIONS
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()), // IMPORTANT: PROPAGATE TRACE CONTEXT
	}

	// CONNECTION
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		log.Fatalf("⛔ Failed to connect to fraud service: %v", err)
	}

	// SET FRAUD CLIENT
	FraudClient = paymentpb.NewFraudServiceClient(conn)
}
