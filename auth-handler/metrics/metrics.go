package metrics

import (
	"net/http"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	apiRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "Number of API requests",
		},
		[]string{"service", "endpoint"},
	)
	apiRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_request_duration_seconds",
			Help:    "Duration of API requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "endpoint"},
	)
	memoryUsage = prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "memory_usage_bytes",
			Help: "Current memory usage in bytes",
		},
		func() float64 {
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)
			return float64(mem.Alloc)
		},
	)
)

func init() {
	prometheus.MustRegister(apiRequests, apiRequestDuration, memoryUsage)
}

func RecordMetrics(service, endpoint string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(apiRequestDuration.WithLabelValues(service, endpoint))
		defer timer.ObserveDuration()

		apiRequests.WithLabelValues(service, endpoint).Inc()

		next.ServeHTTP(w, r)
	})
}
