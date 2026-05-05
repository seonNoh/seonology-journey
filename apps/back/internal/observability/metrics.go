package observability

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds Prometheus metric collectors for the gRPC service.
type Metrics struct {
	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	ActiveStreams   prometheus.Gauge
}

// NewMetrics creates and registers the standard service metrics.
func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		RequestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "journey",
			Subsystem: "grpc",
			Name:      "requests_total",
			Help:      "Total number of gRPC requests by method and status.",
		}, []string{"method", "code"}),

		RequestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "journey",
			Subsystem: "grpc",
			Name:      "request_duration_seconds",
			Help:      "Duration of gRPC requests in seconds.",
			Buckets:   prometheus.DefBuckets,
		}, []string{"method"}),

		ActiveStreams: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "journey",
			Subsystem: "grpc",
			Name:      "active_streams",
			Help:      "Number of currently active gRPC streams.",
		}),
	}

	reg.MustRegister(m.RequestsTotal, m.RequestDuration, m.ActiveStreams)
	return m
}

// Handler returns an http.Handler that serves the /metrics endpoint.
func Handler() http.Handler {
	return promhttp.Handler()
}
