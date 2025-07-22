package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gouthamve/hardcover-book-embed/internal/cache"
	"github.com/gouthamve/hardcover-book-embed/internal/hardcover"
	"github.com/gouthamve/hardcover-book-embed/internal/metrics"
)

type Server struct {
	client         hardcover.Client
	cache          *cache.MemoryCache
	allowedOrigins string
}

func NewServer(client hardcover.Client, cache *cache.MemoryCache, allowedOrigins string) *Server {
	return &Server{
		client:         client,
		cache:          cache,
		allowedOrigins: allowedOrigins,
	}
}

func (s *Server) enableCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	if s.allowedOrigins == "*" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	} else if origin != "" {
		allowedOrigins := strings.Split(s.allowedOrigins, ",")
		for _, allowed := range allowedOrigins {
			if strings.TrimSpace(allowed) == origin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}
	}

	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Max-Age", "86400")

	// Security headers for API responses
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none';")
}

func (s *Server) HandleUserCurrentlyReading(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w, r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Extract username from path parameter
	username := r.PathValue("username")

	// Validate username (alphanumeric, hyphens, underscores)
	if username == "" || !isValidUsername(username) {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("currently_reading_%s", username)

	if cached, found := s.cache.Get(cacheKey); found {
		metrics.CacheHitsTotal.WithLabelValues("currently-reading", username).Inc()
		log.Printf("Serving cached currently reading books for user: %s", username)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(cached); err != nil {
			log.Printf("Error encoding cached response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	metrics.CacheMissesTotal.WithLabelValues("currently-reading", username).Inc()

	log.Printf("Fetching currently reading books for user: %s", username)
	books, err := s.client.GetUserCurrentlyReadingBooksByUsername(username)
	if err != nil {
		log.Printf("Error fetching books for user %s: %v", username, err)
		http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
		return
	}

	s.cache.Set(cacheKey, books)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(books); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// usernameRegex validates usernames containing only alphanumeric characters, hyphens, and underscores
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// isValidUsername checks if the username contains only alphanumeric characters, hyphens, and underscores
func isValidUsername(username string) bool {
	return usernameRegex.MatchString(username)
}

func (s *Server) HandleUserLastRead(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w, r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Extract username from path parameter
	username := r.PathValue("username")

	// Validate username (alphanumeric, hyphens, underscores)
	if username == "" || !isValidUsername(username) {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("last_read_%s", username)

	if cached, found := s.cache.Get(cacheKey); found {
		metrics.CacheHitsTotal.WithLabelValues("last-read", username).Inc()
		log.Printf("Serving cached last read books for user: %s", username)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(cached); err != nil {
			log.Printf("Error encoding cached response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	metrics.CacheMissesTotal.WithLabelValues("last-read", username).Inc()

	log.Printf("Fetching last read books for user: %s", username)
	books, err := s.client.GetUserLastReadBooksByUsername(username)
	if err != nil {
		log.Printf("Error fetching last read books for user %s: %v", username, err)
		http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
		return
	}

	s.cache.Set(cacheKey, books)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(books); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) HandleUserReviews(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w, r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Extract username from path parameter
	username := r.PathValue("username")

	// Validate username (alphanumeric, hyphens, underscores)
	if username == "" || !isValidUsername(username) {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("reviews_%s", username)

	if cached, found := s.cache.Get(cacheKey); found {
		metrics.CacheHitsTotal.WithLabelValues("reviews", username).Inc()
		log.Printf("Serving cached reviews for user: %s", username)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(cached); err != nil {
			log.Printf("Error encoding cached response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	metrics.CacheMissesTotal.WithLabelValues("reviews", username).Inc()

	log.Printf("Fetching reviews for user: %s", username)
	books, err := s.client.GetUserReviewsByUsername(username)
	if err != nil {
		log.Printf("Error fetching reviews for user %s: %v", username, err)
		http.Error(w, "Failed to fetch reviews", http.StatusInternalServerError)
		return
	}

	s.cache.Set(cacheKey, books)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(books); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
