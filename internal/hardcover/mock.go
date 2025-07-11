package hardcover

import (
	"fmt"
	"time"
)

// MockClient is a mock implementation of the Client interface for testing
type MockClient struct {
	// GetUserBooksByUsernameFunc allows custom behavior for testing
	GetUserBooksByUsernameFunc func(username string) (*UserBooksResponse, error)
	// GetUserLastReadBooksByUsernameFunc allows custom behavior for testing
	GetUserLastReadBooksByUsernameFunc func(username string) (*UserBooksResponse, error)

	// CallCount tracks method invocations
	GetUserBooksCalls     []string
	GetLastReadBooksCalls []string
}

// NewMockClient creates a new mock client with default behavior
func NewMockClient() *MockClient {
	return &MockClient{
		GetUserBooksCalls:     []string{},
		GetLastReadBooksCalls: []string{},
	}
}

// GetUserBooksByUsername implements the Client interface
func (m *MockClient) GetUserBooksByUsername(username string) (*UserBooksResponse, error) {
	m.GetUserBooksCalls = append(m.GetUserBooksCalls, username)

	if m.GetUserBooksByUsernameFunc != nil {
		return m.GetUserBooksByUsernameFunc(username)
	}

	// Default mock response based on real JSON data
	updatedAt1, _ := time.Parse(time.RFC3339, "2025-07-09T22:22:43.459059Z")
	updatedAt2, _ := time.Parse(time.RFC3339, "2025-07-07T19:55:40.080268Z")
	updatedAt3, _ := time.Parse(time.RFC3339, "2025-07-01T09:57:55.571541Z")
	updatedAt4, _ := time.Parse(time.RFC3339, "2025-07-01T09:57:47.96016Z")
	updatedAt5, _ := time.Parse(time.RFC3339, "2025-05-28T12:46:25.915366Z")
	responseUpdatedAt, _ := time.Parse(time.RFC3339, "2025-07-11T09:12:20.384103+02:00")
	
	return &UserBooksResponse{
		Books: []UserBook{
			{
				Book: Book{
					ID:    386725,
					Title: "Shakespeare: The World as Stage",
					Slug:  "shakespeare-the-world-as-stage",
					Image: &Image{
						URL: "https://assets.hardcover.app/edition/13396527/7281320-L.jpg",
					},
				},
				UpdatedAt: updatedAt1,
			},
			{
				Book: Book{
					ID:    1946043,
					Title: "Folk Tales & Fables from Bulgaria",
					Slug:  "folk-tales-fables-from-bulgaria",
					Image: &Image{
						URL: "https://assets.hardcover.app/edition/32101711/57e32301c2253bba6dd1e0e1d07c0ea7505752b7.jpeg",
					},
				},
				UpdatedAt: updatedAt2,
			},
			{
				Book: Book{
					ID:    662039,
					Title: "The Battle for the Falklands",
					Slug:  "the-battle-for-the-falklands",
					Image: &Image{
						URL: "https://assets.hardcover.app/book_mappings/7397268/cb4e36451c993f4a73c0db4bf97c25687059b556.jpeg",
					},
				},
				UpdatedAt: updatedAt3,
			},
			{
				Book: Book{
					ID:    120339,
					Title: "Salt: A World History",
					Slug:  "salt",
					Image: &Image{
						URL: "https://assets.hardcover.app/external_data/59911635/995d64e0921d5eb2cb16f6d153aaa6576f7daf17.jpeg",
					},
				},
				UpdatedAt: updatedAt4,
			},
			{
				Book: Book{
					ID:    1896516,
					Title: "Irani Cafe",
					Slug:  "irani-cafe",
					Image: &Image{
						URL: "https://assets.hardcover.app/static/covers/cover1.webp",
					},
				},
				UpdatedAt: updatedAt5,
			},
		},
		Count:     5,
		UpdatedAt: responseUpdatedAt,
	}, nil
}

// GetUserLastReadBooksByUsername implements the Client interface
func (m *MockClient) GetUserLastReadBooksByUsername(username string) (*UserBooksResponse, error) {
	m.GetLastReadBooksCalls = append(m.GetLastReadBooksCalls, username)

	if m.GetUserLastReadBooksByUsernameFunc != nil {
		return m.GetUserLastReadBooksByUsernameFunc(username)
	}

	// Default mock response
	return &UserBooksResponse{
		Books: []UserBook{
			{
				Book: Book{
					ID:    2,
					Title: "Mock Finished Book",
					Slug:  "mock-finished-book",
					Image: &Image{
						URL: "https://example.com/cover2.jpg",
					},
				},
				UpdatedAt: time.Now(),
				Rating:    &[]float64{4.5}[0],
			},
		},
		Count:     1,
		UpdatedAt: time.Now(),
	}, nil
}

// WithError returns a mock client that always returns an error
func (m *MockClient) WithError(err error) *MockClient {
	m.GetUserBooksByUsernameFunc = func(username string) (*UserBooksResponse, error) {
		return nil, err
	}
	m.GetUserLastReadBooksByUsernameFunc = func(username string) (*UserBooksResponse, error) {
		return nil, err
	}
	return m
}

// WithEmptyResponse returns a mock client that returns empty book lists
func (m *MockClient) WithEmptyResponse() *MockClient {
	emptyResponse := &UserBooksResponse{
		Books:     []UserBook{},
		Count:     0,
		UpdatedAt: time.Now(),
	}

	m.GetUserBooksByUsernameFunc = func(username string) (*UserBooksResponse, error) {
		return emptyResponse, nil
	}
	m.GetUserLastReadBooksByUsernameFunc = func(username string) (*UserBooksResponse, error) {
		return emptyResponse, nil
	}
	return m
}

// Reset clears all recorded calls
func (m *MockClient) Reset() {
	m.GetUserBooksCalls = []string{}
	m.GetLastReadBooksCalls = []string{}
}

// AssertCalled verifies a method was called with specific username
func (m *MockClient) AssertCalled(method string, username string) error {
	switch method {
	case "GetUserBooksByUsername":
		for _, call := range m.GetUserBooksCalls {
			if call == username {
				return nil
			}
		}
		return fmt.Errorf("GetUserBooksByUsername was not called with username: %s", username)
	case "GetUserLastReadBooksByUsername":
		for _, call := range m.GetLastReadBooksCalls {
			if call == username {
				return nil
			}
		}
		return fmt.Errorf("GetUserLastReadBooksByUsername was not called with username: %s", username)
	default:
		return fmt.Errorf("unknown method: %s", method)
	}
}
