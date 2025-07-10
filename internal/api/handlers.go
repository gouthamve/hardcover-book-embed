package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gouthamve/hardcover-book-embed/internal/cache"
	"github.com/gouthamve/hardcover-book-embed/internal/hardcover"
)

type Server struct {
	client         *hardcover.Client
	cache          *cache.MemoryCache
	allowedOrigins string
}

func NewServer(client *hardcover.Client, cache *cache.MemoryCache, allowedOrigins string) *Server {
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
}

func (s *Server) HandleCurrentlyReading(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w, r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	const cacheKey = "currently_reading"

	if cached, found := s.cache.Get(cacheKey); found {
		log.Println("Serving cached currently reading books")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cached)
		return
	}

	log.Println("Fetching currently reading books from Hardcover API")
	books, err := s.client.GetCurrentlyReadingBooks()
	if err != nil {
		log.Printf("Error fetching books: %v", err)
		http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
		return
	}

	s.cache.Set(cacheKey, books)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (s *Server) HandleHealth(w http.ResponseWriter, r *http.Request) {
	s.enableCORS(w, r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"status":  "healthy",
		"service": "hardcover-book-embed",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
