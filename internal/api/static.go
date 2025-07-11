package api

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gouthamve/hardcover-book-embed/internal/metrics"
)

// StaticFile holds metadata about a static file
type StaticFile struct {
	Path         string
	ETag         string
	LastModified time.Time
	Size         int64
}

// StaticHandler serves static files with proper caching headers
type StaticHandler struct {
	root  string
	files map[string]*StaticFile
	mu    sync.RWMutex
}

// NewStaticHandler creates a new static file handler
func NewStaticHandler(root string) *StaticHandler {
	return &StaticHandler{
		root:  root,
		files: make(map[string]*StaticFile),
	}
}

// calculateETag generates an ETag for a file
func calculateETag(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := file.Close(); err != nil {
			// Log but don't fail on close error
			fmt.Printf("Error closing file %s: %v\n", path, err)
		}
	}()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf(`"%x"`, hash.Sum(nil)), nil
}

// getOrUpdateFileInfo gets cached file info or updates it if stale
func (h *StaticHandler) getOrUpdateFileInfo(urlPath string) (*StaticFile, error) {
	filePath := filepath.Join(h.root, urlPath)

	// Check file exists and get info
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	h.mu.RLock()
	cached, exists := h.files[urlPath]
	h.mu.RUnlock()

	// If cached and not modified, return cached version
	if exists && cached.LastModified.Equal(info.ModTime()) {
		return cached, nil
	}

	// Calculate new ETag
	etag, err := calculateETag(filePath)
	if err != nil {
		return nil, err
	}

	// Update cache
	fileInfo := &StaticFile{
		Path:         filePath,
		ETag:         etag,
		LastModified: info.ModTime(),
		Size:         info.Size(),
	}

	h.mu.Lock()
	h.files[urlPath] = fileInfo
	h.mu.Unlock()

	return fileInfo, nil
}

// ServeHTTP handles static file requests
func (h *StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Start timing the request
	start := time.Now()

	// Only handle GET and HEAD
	if r.Method != "GET" && r.Method != "HEAD" {
		metrics.StaticFileRequestsTotal.WithLabelValues("", "405").Inc()
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Clean the path
	urlPath := strings.TrimPrefix(r.URL.Path, "/static/")
	if urlPath == "" || strings.Contains(urlPath, "..") {
		metrics.StaticFileRequestsTotal.WithLabelValues(urlPath, "404").Inc()
		http.NotFound(w, r)
		return
	}

	// Get file info
	fileInfo, err := h.getOrUpdateFileInfo(urlPath)
	if err != nil {
		if os.IsNotExist(err) {
			metrics.StaticFileRequestsTotal.WithLabelValues(urlPath, "404").Inc()
			http.NotFound(w, r)
		} else {
			metrics.StaticFileRequestsTotal.WithLabelValues(urlPath, "500").Inc()
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Set caching headers
	w.Header().Set("ETag", fileInfo.ETag)
	w.Header().Set("Last-Modified", fileInfo.LastModified.UTC().Format(http.TimeFormat))

	// Set cache control for 1 year (immutable for versioned assets)
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

	// Set content type based on file extension
	switch filepath.Ext(urlPath) {
	case ".js":
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	case ".css":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case ".html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	default:
		// Let ServeFile detect the content type
	}

	// CORS for JavaScript files
	if filepath.Ext(urlPath) == ".js" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	// Check conditional requests
	// Check If-None-Match (ETag)
	if match := r.Header.Get("If-None-Match"); match != "" {
		if match == fileInfo.ETag {
			metrics.StaticFileRequestsTotal.WithLabelValues(urlPath, "304").Inc()
			metrics.StaticFileRequestDuration.WithLabelValues(urlPath).Observe(time.Since(start).Seconds())
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	// Check If-Modified-Since
	if modifiedSince := r.Header.Get("If-Modified-Since"); modifiedSince != "" {
		t, err := time.Parse(http.TimeFormat, modifiedSince)
		if err == nil && !fileInfo.LastModified.After(t) {
			metrics.StaticFileRequestsTotal.WithLabelValues(urlPath, "304").Inc()
			metrics.StaticFileRequestDuration.WithLabelValues(urlPath).Observe(time.Since(start).Seconds())
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	// Set Content-Length
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size, 10))

	// Serve the file
	http.ServeFile(w, r, fileInfo.Path)

	// Record successful request
	metrics.StaticFileRequestsTotal.WithLabelValues(urlPath, "200").Inc()
	metrics.StaticFileRequestDuration.WithLabelValues(urlPath).Observe(time.Since(start).Seconds())
}
