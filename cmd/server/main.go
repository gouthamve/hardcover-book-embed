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

	// Register routes with patterns
	mux.HandleFunc("GET /api/books/currently-reading/{username}", server.HandleUserCurrentlyReading)
	mux.HandleFunc("GET /api/books/last-read/{username}", server.HandleUserLastRead)
	
	// Handle OPTIONS for CORS
	mux.HandleFunc("OPTIONS /api/books/currently-reading/{username}", server.HandleUserCurrentlyReading)
	mux.HandleFunc("OPTIONS /api/books/last-read/{username}", server.HandleUserLastRead)

	mux.HandleFunc("GET /test-widget.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/test-widget.html")
	})

	mux.HandleFunc("GET /embed.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/embed.html")
	})

	mux.HandleFunc("GET /widget.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, "./web/widget.js")
	})

	log.Printf("Server starting on port %s", port)
	log.Printf("Cache TTL: %v", cacheTTL)
	log.Printf("Allowed origins: %s", allowedOrigins)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
