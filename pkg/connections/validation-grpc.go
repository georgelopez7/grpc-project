package connections

import (
	"log"
	"os"

	"github.com/georgelopez7/grpc-project/api/proto/paymentpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ValidationClient paymentpb.ValidationServiceClient

func InitValidationServiceClient() {
	// --- VALIDATION SERVICE ADDRESS ---
	validationServiceAddress := os.Getenv("VALIDATION_SERVICE_ADDR")
	if validationServiceAddress == "" {
		validationServiceAddress = "localhost:50052"
	}

	// --- GRPC CONNECTION OPTIONS ---
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// --- CREATE CONNECTION ---
	clientConn, err := grpc.NewClient(validationServiceAddress, opts...)
	if err != nil {
		log.Fatalf("⛔ Failed to connect to validation service: %v", err)
	}

	// --- CREATE VALIDATION CLIENT ---
	ValidationClient = paymentpb.NewValidationServiceClient(clientConn)
	log.Printf("🟢 Connected to validation service at %s", validationServiceAddress)
}
