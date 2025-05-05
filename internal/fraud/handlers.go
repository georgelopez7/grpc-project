package fraud

import (
	"context"

	"github.com/georgelopez7/grpc-project/api/proto/paymentpb"
	"github.com/georgelopez7/grpc-project/pkg/utils"
)

func (s *FraudServer) FraudCheck(ctx context.Context, req *paymentpb.PaymentRequest) (*paymentpb.FraudResponse, error) {
	// --- CHECK IF NUMBER IS A FIBONACCI NUMBER ---
	isFibonacci := utils.IsFibonacci(int(req.Amount))
	if isFibonacci {
		return &paymentpb.FraudResponse{
			IsFraudulent: true,
			Message:      "Payment is fraudulent (Fibonacci number)!",
		}, nil
	}

	return &paymentpb.FraudResponse{
		IsFraudulent: false,
		Message:      "Payment is not fraudulent!",
	}, nil
}
