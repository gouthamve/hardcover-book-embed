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

	// Rate Limiting Metrics
	RateLimitWaitDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "hardcoverembed_rate_limit_wait_duration_seconds",
			Help:    "Time spent waiting for rate limiter",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"endpoint"},
	)

	// Static File Metrics
	StaticFileRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hardcoverembed_static_file_requests_total",
			Help: "Total number of static file requests",
		},
		[]string{"file", "status"},
	)

	StaticFileRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "hardcoverembed_static_file_request_duration_seconds",
			Help:    "Static file request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"file"},
	)
)

// Init initializes the metrics (placeholder for any future initialization needs)
func Init() {
	// Currently metrics are auto-registered with promauto
	// This function exists for future expansion if needed
}
