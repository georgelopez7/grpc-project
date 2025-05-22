package gateway

import (
	"log/slog"
	"net/http"

	"github.com/georgelopez7/grpc-project/api/proto/paymentpb"
	"github.com/georgelopez7/grpc-project/pkg/connections"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("gateway/payment_handler")

func PaymentRequestHandler(c *gin.Context) {
	ctx := c.Request.Context()

	slog.InfoContext(ctx, "✅ Received payment request", c.Request.Method, c.Request.URL.Path)

	// PARSE BODY
	var requestBody PaymentRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": err.Error()})
		return
	}

	// // THIS WILL ONLY SHOW UP IN DEVELOPMENT MODE
	// slog.DebugContext(ctx, "🧠 Parsed request body", "payment_id", requestBody.ID, "amount", requestBody.Amount)

	// START TRACE
	ctx, span := tracer.Start(ctx, "handle_payment_request", trace.WithAttributes(
		attribute.String("payment.id", requestBody.ID),
		attribute.Float64("payment.amount", float64(requestBody.Amount)),
	))
	defer span.End()

	// CREATE PAYMENT REQUEST
	paymentRequest := &paymentpb.PaymentRequest{
		Id:       requestBody.ID,
		Amount:   requestBody.Amount,
		Sender:   requestBody.Sender,
		Receiver: requestBody.Receiver,
	}

	// CHECK PAYMENT FOR FRAUD
	fraudResponse, err := connections.FraudClient.FraudCheck(ctx, paymentRequest)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to check for fraud: " + err.Error()})
		return
	}
	slog.InfoContext(ctx, "🔎 Fraud check result", "fraud_status", fraudResponse.IsFraudulent, "message", fraudResponse.Message)

	// CHECK PAYMENT FOR VALIDATION
	validationResponse, err := connections.ValidationClient.ValidatePayment(ctx, paymentRequest)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to validate payment: " + err.Error()})
		return
	}
	slog.InfoContext(ctx, "👓 Validation result", "validation_status", validationResponse.IsValid, "message", validationResponse.Message)

	// ADDITIONAL SPAN
	_, endSpan := tracer.Start(ctx, "finalize_payment_response")
	defer endSpan.End()

	endSpan.SetAttributes(
		attribute.String("response.status", "success"),
		attribute.String("response.fraud", fraudResponse.Message),
		attribute.String("response.validation", validationResponse.Message),
	)

	slog.InfoContext(ctx, "🚀 Payment request processed successfully")

	// SEND RESPONSE
	c.JSON(200, gin.H{
		"fraud_status": FraudulentResponse{
			IsFraudulent: fraudResponse.IsFraudulent,
			Message:      fraudResponse.Message,
		},
		"validation_status": ValidResponse{
			IsValid: validationResponse.IsValid,
			Message: validationResponse.Message,
		},
	})
}
