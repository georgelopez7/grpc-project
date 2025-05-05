package gateway

import (
	"github.com/gin-gonic/gin"
)

func SendMetrics(c *gin.Context) {
	// --- INCREMENT PAYMENT REQUESTS COUNTER ---
	PaymentRequestsCount.Inc()
	c.Next()
}
