package validation

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/georgelopez7/grpc-project/api/proto/paymentpb"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("validation/validation_check")

func (s *ValidationServer) ValidatePayment(ctx context.Context, req *paymentpb.PaymentRequest) (*paymentpb.ValidationResponse, error) {
	// START SPAN
	_, span := tracer.Start(ctx, "validate_payment_request")
	defer span.End()

	// ADD ATTRIBUTES FOR OBSERVABILITY
	span.SetAttributes(
		attribute.String("payment.id", req.Id),
		attribute.Float64("payment.amount", float64(req.Amount)),
	)

	// CHECK AMOUNT IS GREATER THAN ZERO
	if req.Amount <= 0 {
		slog.WarnContext(ctx, "🚨 Payment is below the minimum allowed!")
		return &paymentpb.ValidationResponse{
			IsValid: false,
			Message: "Amount must be greater than zero!",
		}, nil
	}

	// CHECK AMOUNT IS LESS THAN MAX
	MAX_AMOUNT := ValidationConfig.MaxAmount
	if req.Amount > int32(MAX_AMOUNT) {
		slog.WarnContext(ctx, "🚨 Payment is above the maximum allowed!")
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
