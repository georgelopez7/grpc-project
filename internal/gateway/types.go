package gateway

type PaymentRequest struct {
	ID       string `json:"id" binding:"required"`
	Amount   int32  `json:"amount" binding:"required"`
	Sender   string `json:"sender" binding:"required"`
	Receiver string `json:"receiver" binding:"required"`
}

type FraudulentResponse struct {
	IsFraudulent bool   `json:"is_fraudulent"`
	Message      string `json:"message"`
}

type ValidResponse struct {
	IsValid bool   `json:"is_valid"`
	Message string `json:"message"`
}
