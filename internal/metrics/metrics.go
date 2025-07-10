package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP Metrics
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hardcoverembed_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"endpoint", "method", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "hardcoverembed_http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint", "method"},
	)

	HTTPRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "hardcoverembed_http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed",
		},
	)

	// Cache Metrics
	CacheHitsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hardcoverembed_cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"endpoint", "username"},
	)

	CacheMissesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hardcoverembed_cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"endpoint", "username"},
	)

	CacheSize = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "hardcoverembed_cache_size",
			Help: "Current number of items in cache",
		},
	)

	CacheEvictionsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "hardcoverembed_cache_evictions_total",
			Help: "Total number of cache evictions",
		},
	)

	// Hardcover API Metrics
	HardcoverAPIRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hardcoverembed_hardcover_api_requests_total",
			Help: "Total number of Hardcover API requests",
		},
		[]string{"endpoint", "status", "username"},
	)

	HardcoverAPIRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "hardcoverembed_hardcover_api_request_duration_seconds",
			Help:    "Hardcover API request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint", "username"},
	)
)

// Init initializes the metrics (placeholder for any future initialization needs)
func Init() {
	// Currently metrics are auto-registered with promauto
	// This function exists for future expansion if needed
}
