package gateway

import (
	"context"
	"net/http"

	"github.com/georgelopez7/grpc-project/api/proto/paymentpb"
	"github.com/georgelopez7/grpc-project/pkg/connections"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("gateway")

func PaymentRequestHandler(c *gin.Context) {
	// --- PARSE BODY ---
	var requestBody PaymentRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": err.Error()})
		return
	}

	// --- START TRACE ---
	ctx := c.Request.Context()
	ctx, span := tracer.Start(ctx, "payment_request_handler", trace.WithAttributes(
		attribute.String("payment.id", requestBody.ID),
		attribute.Float64("payment.amount", float64(requestBody.Amount)),
	))
	defer span.End()

	// --- CREATE PAYMENT REQUEST ---
	paymentRequest := &paymentpb.PaymentRequest{
		Id:       requestBody.ID,
		Amount:   requestBody.Amount,
		Sender:   requestBody.Sender,
		Receiver: requestBody.Receiver,
	}

	// --- CHECK PAYMENT FOR FRAUD ---
	fraudResponse, err := connections.FraudClient.FraudCheck(ctx, paymentRequest)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to check for fraud: " + err.Error()})
		return
	}

	// --- CHECK PAYMENT FOR VALIDATION ---
	validationResponse, err := connections.ValidationClient.ValidatePayment(context.TODO(), paymentRequest)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to validate payment: " + err.Error()})
		return
	}

	// --- SEND RESPONSE ---
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
