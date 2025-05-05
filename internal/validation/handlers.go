package validation

import (
	"context"
	"fmt"

	"github.com/georgelopez7/grpc-project/api/proto/paymentpb"
)

func (s *ValidationServer) ValidatePayment(ctx context.Context, req *paymentpb.PaymentRequest) (*paymentpb.ValidationResponse, error) {
	// --- CHECK AMOUNT IS GREATER THAN ZERO ---
	if req.Amount <= 0 {
		return &paymentpb.ValidationResponse{
			IsValid: false,
			Message: "Amount must be greater than zero!",
		}, nil
	}

	// --- CHECK AMOUNT IS LESS THAN MAX ---
	MAX_AMOUNT := ValidationConfig.MaxAmount
	if req.Amount > int32(MAX_AMOUNT) {
		return &paymentpb.ValidationResponse{
			IsValid: false,
			Message: fmt.Sprintf("Amount is above the maximum allowed (%d)!", MAX_AMOUNT),
		}, nil
	}

	return &paymentpb.ValidationResponse{
		IsValid: true,
		Message: "Payment passed validation check!",
	}, nil
}
