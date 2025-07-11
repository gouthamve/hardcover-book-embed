package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gouthamve/hardcover-book-embed/internal/metrics"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// MetricsMiddleware wraps HTTP handlers to collect metrics
func MetricsMiddleware(endpoint string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Track in-flight requests
			metrics.HTTPRequestsInFlight.Inc()
			defer metrics.HTTPRequestsInFlight.Dec()

			// Wrap response writer to capture status code
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Call the actual handler
			next(rw, r)

			// Record metrics
			duration := time.Since(start).Seconds()
			status := strconv.Itoa(rw.statusCode)

			metrics.HTTPRequestsTotal.WithLabelValues(endpoint, r.Method, status).Inc()
			metrics.HTTPRequestDuration.WithLabelValues(endpoint, r.Method).Observe(duration)
		}
	}
}
