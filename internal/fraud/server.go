package fraud

import (
	"log"
	"net"

	"github.com/georgelopez7/grpc-project/api/proto/paymentpb"
	"google.golang.org/grpc"
)

type FraudServer struct {
	paymentpb.UnimplementedFraudServiceServer
}

func InitFraudServer() {
	// --- LISTEN ---
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// --- gRPC SERVER ---
	s := grpc.NewServer()
	paymentpb.RegisterFraudServiceServer(s, &FraudServer{})

	log.Println("Fraud service listening on :50052")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
