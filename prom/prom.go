// Package prom exposes metrics used for instrumenting bq application with
// prometheus.
package prom

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Night is a wrapper around prometheus client for instrumentation in bq.
type Night struct {
	RequestCounter  *prometheus.CounterVec
	RequestInFlight prometheus.Gauge
}

// New returns a new *Night instance with the fields initialized.
func New() *Night {
	return &Night{
		RequestCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "total number of request by the service",
		}, []string{"code", "method"}),
		RequestInFlight: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "in_flight_request",
			Help: "A gauge of requests currently being served by the wrapped handler",
		}),
	}
}

// Handle returns a handler which is fully instrumented.
func Handle(n *Night, next http.Handler) http.Handler {
	return promhttp.InstrumentHandlerCounter(n.RequestCounter,
		promhttp.InstrumentHandlerInFlight(n.RequestInFlight, next),
	)
}

// Wrap instruments next handler.
func Wrap(next http.Handler) http.Handler {
	return Handle(New(), next)
}
