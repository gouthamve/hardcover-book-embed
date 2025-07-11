package hardcover

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gouthamve/hardcover-book-embed/internal/metrics"
	"golang.org/x/time/rate"
)

const (
	HardcoverAPIURL = "https://api.hardcover.app/v1/graphql"
	UserAgent       = "hardcover-book-embed/1.0"
)

// Client is the interface for interacting with the Hardcover API
type Client interface {
	GetUserBooksByUsername(username string) (*UserBooksResponse, error)
	GetUserLastReadBooksByUsername(username string) (*UserBooksResponse, error)
}

// HTTPClient interface allows for mocking HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// client is the concrete implementation of the Client interface
type client struct {
	apiToken    string
	httpClient  HTTPClient
	rateLimiter *rate.Limiter
}

// NewClient creates a new Hardcover API client
func NewClient(apiToken string) Client {
	// Hardcover API allows 60 requests per minute
	// We'll be conservative and limit to 50 requests per minute (0.83 per second)
	// with a burst of 5 to handle short spikes
	limiter := rate.NewLimiter(rate.Limit(0.83), 5)

	return &client{
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		rateLimiter: limiter,
	}
}

// NewClientWithHTTPClient creates a new Hardcover API client with a custom HTTP client
func NewClientWithHTTPClient(apiToken string, httpClient HTTPClient) Client {
	limiter := rate.NewLimiter(rate.Limit(0.83), 5)

	return &client{
		apiToken:    apiToken,
		httpClient:  httpClient,
		rateLimiter: limiter,
	}
}

func (c *client) GetUserBooksByUsername(username string) (*UserBooksResponse, error) {
	// Wait for rate limiter
	ctx := context.Background()
	waitStart := time.Now()
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}
	waitDuration := time.Since(waitStart).Seconds()
	if waitDuration > 0.001 { // Only record if we actually waited
		metrics.RateLimitWaitDuration.WithLabelValues("currently-reading").Observe(waitDuration)
	}

	query := fmt.Sprintf(`{
		user_books(
			where: {user: {username: {_eq: "%s"}}, status_id: {_eq: 2}},
			order_by: {updated_at: desc},
			limit: 5
		) {
			rating
			updated_at
			book {
				id
				title
				image {
					url
				}
				slug
			}
		}
	}`, username)

	reqBody := map[string]string{
		"query": query,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", HardcoverAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("User-Agent", UserAgent)

	// Track API request duration
	start := time.Now()
	resp, err := c.httpClient.Do(req)
	duration := time.Since(start).Seconds()
	metrics.HardcoverAPIRequestDuration.WithLabelValues("currently-reading", username).Observe(duration)

	if err != nil {
		metrics.HardcoverAPIRequestsTotal.WithLabelValues("currently-reading", "error", username).Inc()
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log but don't fail on close error
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		metrics.HardcoverAPIRequestsTotal.WithLabelValues("currently-reading", fmt.Sprintf("%d", resp.StatusCode), username).Inc()
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	metrics.HardcoverAPIRequestsTotal.WithLabelValues("currently-reading", "200", username).Inc()

	var graphqlResp UserBooksAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&graphqlResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(graphqlResp.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %v", graphqlResp.Errors)
	}

	// Process books and add fallback images
	books := graphqlResp.Data.UserBooks
	for i := range books {
		if books[i].Book.Image == nil {
			// Generate a fallback image based on book ID
			coverNum := (books[i].Book.ID % 9) + 1
			books[i].Book.Image = &Image{
				URL: fmt.Sprintf("https://assets.hardcover.app/static/covers/cover%d.webp", coverNum),
			}
		}
	}

	return &UserBooksResponse{
		Books:     books,
		Count:     len(books),
		UpdatedAt: time.Now(),
	}, nil
}

func (c *client) GetUserLastReadBooksByUsername(username string) (*UserBooksResponse, error) {
	// Wait for rate limiter
	ctx := context.Background()
	waitStart := time.Now()
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}
	waitDuration := time.Since(waitStart).Seconds()
	if waitDuration > 0.001 { // Only record if we actually waited
		metrics.RateLimitWaitDuration.WithLabelValues("last-read").Observe(waitDuration)
	}

	query := fmt.Sprintf(`{
		user_books(
			where: {user: {username: {_eq: "%s"}}, status_id: {_eq: 3}},
			order_by: {last_read_date: desc_nulls_last},
			limit: 5
		) {
			rating
			updated_at
			last_read_date
			book {
				id
				title
				image {
					url
				}
				slug
			}
		}
	}`, username)

	reqBody := map[string]string{
		"query": query,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", HardcoverAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("User-Agent", UserAgent)

	// Track API request duration
	start := time.Now()
	resp, err := c.httpClient.Do(req)
	duration := time.Since(start).Seconds()
	metrics.HardcoverAPIRequestDuration.WithLabelValues("last-read", username).Observe(duration)

	if err != nil {
		metrics.HardcoverAPIRequestsTotal.WithLabelValues("last-read", "error", username).Inc()
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log but don't fail on close error
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		metrics.HardcoverAPIRequestsTotal.WithLabelValues("last-read", fmt.Sprintf("%d", resp.StatusCode), username).Inc()
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	metrics.HardcoverAPIRequestsTotal.WithLabelValues("last-read", "200", username).Inc()

	var graphqlResp UserBooksAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&graphqlResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(graphqlResp.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %v", graphqlResp.Errors)
	}

	// Process books and add fallback images
	books := graphqlResp.Data.UserBooks
	for i := range books {
		if books[i].Book.Image == nil {
			// Generate a fallback image based on book ID
			coverNum := (books[i].Book.ID % 9) + 1
			books[i].Book.Image = &Image{
				URL: fmt.Sprintf("https://assets.hardcover.app/static/covers/cover%d.webp", coverNum),
			}
		}
	}

	return &UserBooksResponse{
		Books:     books,
		Count:     len(books),
		UpdatedAt: time.Now(),
	}, nil
}
