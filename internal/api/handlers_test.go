package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gouthamve/hardcover-book-embed/internal/cache"
	"github.com/gouthamve/hardcover-book-embed/internal/hardcover"
)

func TestHandleUserCurrentlyReading(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		mockResponse   *hardcover.UserBooksResponse
		mockError      error
		expectedStatus int
		expectedBooks  int
	}{
		{
			name:           "successful request with default mock",
			username:       "testuser",
			expectedStatus: http.StatusOK,
			expectedBooks:  5, // Default mock returns 5 books
		},
		{
			name:           "invalid username",
			username:       "test@user",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty username",
			username:       "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "API error",
			username:       "testuser",
			mockError:      fmt.Errorf("API error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			mockClient := hardcover.NewMockClient()
			if tt.mockResponse != nil {
				mockClient.GetUserBooksByUsernameFunc = func(username string) (*hardcover.UserBooksResponse, error) {
					return tt.mockResponse, tt.mockError
				}
			} else if tt.mockError != nil {
				mockClient.GetUserBooksByUsernameFunc = func(username string) (*hardcover.UserBooksResponse, error) {
					return nil, tt.mockError
				}
			}

			// Create server with mock
			cache := cache.NewMemoryCache(5 * time.Minute)
			server := NewServer(mockClient, cache, "*")

			// Create request
			url := fmt.Sprintf("/api/books/currently-reading/%s", tt.username)
			req := httptest.NewRequest("GET", url, nil)
			req.SetPathValue("username", tt.username)

			// Record response
			w := httptest.NewRecorder()
			server.HandleUserCurrentlyReading(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response body for successful requests
			if tt.expectedStatus == http.StatusOK {
				var response hardcover.UserBooksResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if response.Count != tt.expectedBooks {
					t.Errorf("expected %d books, got %d", tt.expectedBooks, response.Count)
				}
			}

			// Verify mock was called correctly for valid usernames
			if tt.expectedStatus != http.StatusBadRequest && tt.username != "" {
				if err := mockClient.AssertCalled("GetUserBooksByUsername", tt.username); err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestCaching(t *testing.T) {
	// Create mock client that tracks calls
	mockClient := hardcover.NewMockClient()
	callCount := 0
	mockClient.GetUserBooksByUsernameFunc = func(username string) (*hardcover.UserBooksResponse, error) {
		callCount++
		return &hardcover.UserBooksResponse{
			Books: []hardcover.UserBook{
				{
					Book: hardcover.Book{
						ID:    1,
						Title: "Cached Book",
						Slug:  "cached-book",
					},
					UpdatedAt: time.Now(),
				},
			},
			Count:     1,
			UpdatedAt: time.Now(),
		}, nil
	}

	// Create server with short cache TTL for testing
	cache := cache.NewMemoryCache(1 * time.Second)
	server := NewServer(mockClient, cache, "*")

	username := "cachetest"
	url := fmt.Sprintf("/api/books/currently-reading/%s", username)

	// First request - should hit the API
	req1 := httptest.NewRequest("GET", url, nil)
	req1.SetPathValue("username", username)
	w1 := httptest.NewRecorder()
	server.HandleUserCurrentlyReading(w1, req1)

	if callCount != 1 {
		t.Errorf("expected 1 API call, got %d", callCount)
	}

	// Second request - should hit cache
	req2 := httptest.NewRequest("GET", url, nil)
	req2.SetPathValue("username", username)
	w2 := httptest.NewRecorder()
	server.HandleUserCurrentlyReading(w2, req2)

	if callCount != 1 {
		t.Errorf("expected 1 API call (cached), got %d", callCount)
	}

	// Wait for cache to expire
	time.Sleep(2 * time.Second)

	// Third request - should hit API again
	req3 := httptest.NewRequest("GET", url, nil)
	req3.SetPathValue("username", username)
	w3 := httptest.NewRecorder()
	server.HandleUserCurrentlyReading(w3, req3)

	if callCount != 2 {
		t.Errorf("expected 2 API calls (cache expired), got %d", callCount)
	}
}
