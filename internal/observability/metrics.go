package observability

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	InFlightRequests = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "tgbot",
		Subsystem: "msg_handler",
		Name:      "in_flight_requests_total",
	})

	TotalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "tgbot",
		Subsystem: "msg_handler",
		Name:      "in_total_requests_total",
	})
	SummaryResponseTime = promauto.NewSummary(prometheus.SummaryOpts{
		Namespace: "tgbot",
		Subsystem: "msg_handler",
		Name:      "summary_response_time_seconds",
		Objectives: map[float64]float64{
			0.5:  0.1,
			0.9:  0.01,
			0.99: 0.001,
		},
	})
	HistogramResponseTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "tgbot",
			Subsystem: "msg_handler",
			Name:      "histogram_response_time_seconds",
			Buckets: prometheus.ExponentialBuckets(0.0001, 2, 16),
		},
		[]string{"code"},
	)
)

func LogRequest(f func() error) {
	startTime := time.Now()
	TotalRequests.Inc()
	InFlightRequests.Inc()
	err := f()
	InFlightRequests.Dec()
	duration := time.Since(startTime)

	code := ""
	 if err != nil {
		code = "error"
	} else {
		code = "ok"
	}
	HistogramResponseTime.
		WithLabelValues(code).
		Observe(duration.Seconds())

	SummaryResponseTime.Observe(duration.Seconds())
}
