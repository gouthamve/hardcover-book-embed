package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gouthamve/hardcover-book-embed/internal/api"
	"github.com/gouthamve/hardcover-book-embed/internal/cache"
	"github.com/gouthamve/hardcover-book-embed/internal/hardcover"
	"github.com/gouthamve/hardcover-book-embed/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	apiToken := os.Getenv("HARDCOVER_API_TOKEN")
	if apiToken == "" {
		log.Fatal("HARDCOVER_API_TOKEN environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "9090"
	}

	cacheTTLStr := os.Getenv("CACHE_TTL_MINUTES")
	cacheTTL := 30 * time.Minute
	if cacheTTLStr != "" {
		if minutes, err := strconv.Atoi(cacheTTLStr); err == nil {
			cacheTTL = time.Duration(minutes) * time.Minute
		}
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}

	client := hardcover.NewClient(apiToken)
	memCache := cache.NewMemoryCache(cacheTTL)
	server := api.NewServer(client, memCache, allowedOrigins)

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register routes with patterns and metrics middleware
	mux.HandleFunc("GET /api/books/currently-reading/{username}",
		api.MetricsMiddleware("currently-reading")(server.HandleUserCurrentlyReading))
	mux.HandleFunc("GET /api/books/last-read/{username}",
		api.MetricsMiddleware("last-read")(server.HandleUserLastRead))

	// Handle OPTIONS for CORS
	mux.HandleFunc("OPTIONS /api/books/currently-reading/{username}", server.HandleUserCurrentlyReading)
	mux.HandleFunc("OPTIONS /api/books/last-read/{username}", server.HandleUserLastRead)

	mux.HandleFunc("GET /test-widget.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/test-widget.html")
	})

	mux.HandleFunc("GET /embed.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/embed.html")
	})

	// Static file handler with caching
	staticHandler := api.NewStaticHandler("./web/static")
	mux.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	// Initialize metrics
	metrics.Init()

	// Start metrics server
	go func() {
		metricsMux := http.NewServeMux()
		metricsMux.Handle("/metrics", promhttp.Handler())
		log.Printf("Metrics server starting on port %s", metricsPort)
		if err := http.ListenAndServe(":"+metricsPort, metricsMux); err != nil {
			log.Fatal("Metrics server failed to start:", err)
		}
	}()

	log.Printf("Server starting on port %s", port)
	log.Printf("Metrics available on port %s/metrics", metricsPort)
	log.Printf("Cache TTL: %v", cacheTTL)
	log.Printf("Allowed origins: %s", allowedOrigins)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
