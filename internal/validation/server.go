package validation

import (
	"log"
	"net"

	"github.com/georgelopez7/grpc-project/api/proto/paymentpb"
	"google.golang.org/grpc"
)

type ValidationServer struct {
	paymentpb.UnimplementedValidationServiceServer
}

func InitValidationServer() {
	// --- LISTEN ---
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// --- gRPC SERVER ---
	s := grpc.NewServer()
	paymentpb.RegisterValidationServiceServer(s, &ValidationServer{})

	log.Println("Validator service listening on :50053")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
