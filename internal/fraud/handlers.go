package fraud

import (
	"context"
	"log/slog"

	"github.com/georgelopez7/grpc-project/api/proto/paymentpb"
	"github.com/georgelopez7/grpc-project/pkg/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("fraud/fraud_check")

func (s *FraudServer) FraudCheck(ctx context.Context, req *paymentpb.PaymentRequest) (*paymentpb.FraudResponse, error) {
	// START SPAN
	_, span := tracer.Start(ctx, "evaluate_fibonacci_number")
	defer span.End()

	// ADD ATTRIBUTES FOR OBSERVABILITY
	span.SetAttributes(
		attribute.String("payment.id", req.Id),
		attribute.Float64("payment.amount", float64(req.Amount)),
	)

	// CHECK IF NUMBER IS A FIBONACCI NUMBER
	isFibonacci := utils.IsFibonacci(int(req.Amount))
	if isFibonacci {
		slog.WarnContext(ctx, "🚨 Payment is a Fibonacci number!")
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
