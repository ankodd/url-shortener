package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics is used to store metrics
type Metrics struct {
	ErrorsCount  prometheus.Counter
	RequestCount prometheus.Counter
}

// NewMetrics returns new metrics
func NewMetrics() *Metrics {
	return &Metrics{
		ErrorsCount: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "url_shortener",
			Subsystem: "http",
			Name:      "errors_count",
			Help:      "The total number of HTTP errors",
		}),
		RequestCount: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "url_shortener",
			Subsystem: "http",
			Name:      "request_count",
			Help:      "The total number of HTTP requests",
		}),
	}
}

// IncError is used to increment errors
func (m *Metrics) IncError() {
	m.ErrorsCount.Inc()
}

// IncRequest is used to increment requests
func (m *Metrics) IncRequest() {
	m.RequestCount.Inc()
}
