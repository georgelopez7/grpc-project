package connections

import (
	"log"
	"os"

	"github.com/georgelopez7/grpc-project/api/proto/paymentpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var FraudClient paymentpb.FraudServiceClient

func InitFraudServiceClient() {
	// --- FRAUD SERVICE ADDRESS ---
	fraudServiceAddr := os.Getenv("FRAUD_SERVICE_ADDR")
	if fraudServiceAddr == "" {
		fraudServiceAddr = "localhost:50052"
	}

	// --- GRPC CONNECTION OPTIONS ---
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// --- CREATE CONNECTION ---
	clientConn, err := grpc.NewClient(fraudServiceAddr, opts...)
	if err != nil {
		log.Fatalf("⛔ Failed to connect to fraud service: %v", err)
	}

	// --- CREATE FRAUD CLIENT ---
	FraudClient = paymentpb.NewFraudServiceClient(clientConn)
	log.Printf("🟢 Connected to fraud service at %s", fraudServiceAddr)
}
