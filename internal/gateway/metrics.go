package gateway

import (
	"github.com/prometheus/client_golang/prometheus"
)

// METRICS
var PaymentRequestsCount = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "gateway_payment_requests_total",
		Help: "Total number of payment requests received by the gateway!",
	},
)
