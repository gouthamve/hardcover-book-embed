package hardcover

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"
)

// MockHTTPClient is a mock implementation of HTTPClient for testing
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	// Default response
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
	}, nil
}

// These two tests need to be changed whenever the API response changes
func TestClientParsesAPIResponseCorrectly(t *testing.T) {
	// The actual API response JSON
	apiResponse := `{
  "data": {
    "user_books": [
      {
        "rating": null,
        "updated_at": "2025-07-09T22:22:43.459059+00:00",
        "book": {
          "id": 386725,
          "title": "Shakespeare: The World as Stage",
          "image": {
            "url": "https://assets.hardcover.app/edition/13396527/7281320-L.jpg"
          },
          "slug": "shakespeare-the-world-as-stage"
        }
      },
      {
        "rating": null,
        "updated_at": "2025-07-07T19:55:40.080268+00:00",
        "book": {
          "id": 1946043,
          "title": "Folk Tales & Fables from Bulgaria",
          "image": {
            "url": "https://assets.hardcover.app/edition/32101711/57e32301c2253bba6dd1e0e1d07c0ea7505752b7.jpeg"
          },
          "slug": "folk-tales-fables-from-bulgaria"
        }
      },
      {
        "rating": null,
        "updated_at": "2025-07-01T09:57:55.571541+00:00",
        "book": {
          "id": 662039,
          "title": "The Battle for the Falklands",
          "image": {
            "url": "https://assets.hardcover.app/book_mappings/7397268/cb4e36451c993f4a73c0db4bf97c25687059b556.jpeg"
          },
          "slug": "the-battle-for-the-falklands"
        }
      },
      {
        "rating": null,
        "updated_at": "2025-07-01T09:57:47.96016+00:00",
        "book": {
          "id": 120339,
          "title": "Salt: A World History",
          "image": {
            "url": "https://assets.hardcover.app/external_data/59911635/995d64e0921d5eb2cb16f6d153aaa6576f7daf17.jpeg"
          },
          "slug": "salt"
        }
      },
      {
        "rating": null,
        "updated_at": "2025-05-28T12:46:25.915366+00:00",
        "book": {
          "id": 1896516,
          "title": "Irani Cafe",
          "image": null,
          "slug": "irani-cafe"
        }
      }
    ]
  }
}`

	// Create mock HTTP client
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Verify request details
			if req.Method != "POST" {
				t.Errorf("expected POST request, got %s", req.Method)
			}
			if req.URL.String() != HardcoverAPIURL {
				t.Errorf("expected URL %s, got %s", HardcoverAPIURL, req.URL.String())
			}

			// Return mock response
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(apiResponse))),
			}, nil
		},
	}

	// Create client with mock HTTP client
	client := NewClientWithHTTPClient("test-token", mockHTTP)

	// Call the method
	response, err := client.GetUserBooksByUsername("testuser")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify response
	if response.Count != 5 {
		t.Errorf("expected 5 books, got %d", response.Count)
	}

	// Verify book details
	expectedBooks := []struct {
		ID       int
		Title    string
		Slug     string
		HasImage bool
	}{
		{386725, "Shakespeare: The World as Stage", "shakespeare-the-world-as-stage", true},
		{1946043, "Folk Tales & Fables from Bulgaria", "folk-tales-fables-from-bulgaria", true},
		{662039, "The Battle for the Falklands", "the-battle-for-the-falklands", true},
		{120339, "Salt: A World History", "salt", true},
		{1896516, "Irani Cafe", "irani-cafe", false}, // This one has null image
	}

	for i, expected := range expectedBooks {
		if i >= len(response.Books) {
			t.Errorf("missing book at index %d", i)
			continue
		}

		book := response.Books[i].Book

		if book.ID != expected.ID {
			t.Errorf("book %d: expected ID %d, got %d", i, expected.ID, book.ID)
		}
		if book.Title != expected.Title {
			t.Errorf("book %d: expected title %q, got %q", i, expected.Title, book.Title)
		}
		if book.Slug != expected.Slug {
			t.Errorf("book %d: expected slug %q, got %q", i, expected.Slug, book.Slug)
		}

		// Check image handling
		if expected.HasImage {
			if book.Image == nil || book.Image.URL == "" {
				t.Errorf("book %d: expected image URL, got nil", i)
			}
		} else {
			// Should have fallback image
			if book.Image == nil || book.Image.URL == "" {
				t.Errorf("book %d: expected fallback image, got nil", i)
			} else if book.Image.URL != "https://assets.hardcover.app/static/covers/cover1.webp" {
				// ID 1896516 % 9 + 1 = 1
				t.Errorf("book %d: expected fallback image cover1.webp, got %s", i, book.Image.URL)
			}
		}

		// Verify rating is nil
		if response.Books[i].Rating != nil {
			t.Errorf("book %d: expected nil rating, got %v", i, *response.Books[i].Rating)
		}
	}
}

func TestClientParsesLastReadAPIResponse(t *testing.T) {
	// Actual API response for last read books with last_read_date
	apiResponse := `{
  "data": {
    "user_books": [
      {
        "rating": 4,
        "updated_at": "2025-07-01T09:50:29.384902+00:00",
        "last_read_date": "2025-07-01",
        "book": {
          "id": 462038,
          "title": "System Collapse",
          "image": {
            "url": "https://assets.hardcover.app/editions/30845693/9ef45cef-6b2b-44a8-9d1c-1b625bb863fd.jpg"
          },
          "slug": "system-collapse"
        }
      },
      {
        "rating": 3.5,
        "updated_at": "2025-06-30T22:20:30.789569+00:00",
        "last_read_date": "2025-06-30",
        "book": {
          "id": 1819865,
          "title": "Edible Stories: A Novel in Sixteen Parts",
          "image": {
            "url": "https://assets.hardcover.app/edition/31943115/1d4f4b13a3a1e7ca8a2c04a019a510748ce9bffc.jpeg"
          },
          "slug": "edible-stories"
        }
      },
      {
        "rating": 5,
        "updated_at": "2025-06-30T16:53:02.392308+00:00",
        "last_read_date": "2025-06-30",
        "book": {
          "id": 277100,
          "title": "84, Charing Cross Road",
          "image": {
            "url": "https://assets.hardcover.app/external_data/30184810/181c007d085889ad012f987bdf456440386a946c.jpeg"
          },
          "slug": "84-charing-cross-road"
        }
      },
      {
        "rating": 4,
        "updated_at": "2025-06-28T23:36:41.432323+00:00",
        "last_read_date": "2025-06-28",
        "book": {
          "id": 380944,
          "title": "Medium Raw: A Bloody Valentine to the World of Food and the People Who Cook",
          "image": {
            "url": "https://assets.hardcover.app/external_data/59799176/6909db1910f811a4b9f830caaee6238ebdc7fb17.jpeg"
          },
          "slug": "medium-raw"
        }
      },
      {
        "rating": 4.5,
        "updated_at": "2025-06-26T14:16:19.284238+00:00",
        "last_read_date": "2025-06-26",
        "book": {
          "id": 435167,
          "title": "Fugitive Telemetry",
          "image": {
            "url": "https://assets.hardcover.app/book_mappings/7333079/1dd5f1cc67a36104185d7cc0c9967729fc44bf18.jpeg"
          },
          "slug": "fugitive-telemetry"
        }
      }
    ]
  }
}`

	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(apiResponse))),
			}, nil
		},
	}

	client := NewClientWithHTTPClient("test-token", mockHTTP)
	response, err := client.GetUserLastReadBooksByUsername("testuser")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify response
	if response.Count != 5 {
		t.Errorf("expected 5 books, got %d", response.Count)
	}

	// Expected data from the response
	expectedBooks := []struct {
		ID           int
		Title        string
		Slug         string
		Rating       float64
		LastReadDate string
	}{
		{462038, "System Collapse", "system-collapse", 4, "2025-07-01"},
		{1819865, "Edible Stories: A Novel in Sixteen Parts", "edible-stories", 3.5, "2025-06-30"},
		{277100, "84, Charing Cross Road", "84-charing-cross-road", 5, "2025-06-30"},
		{380944, "Medium Raw: A Bloody Valentine to the World of Food and the People Who Cook", "medium-raw", 4, "2025-06-28"},
		{435167, "Fugitive Telemetry", "fugitive-telemetry", 4.5, "2025-06-26"},
	}

	// Verify each book
	for i, expected := range expectedBooks {
		if i >= len(response.Books) {
			t.Errorf("missing book at index %d", i)
			continue
		}

		book := response.Books[i].Book
		if book.ID != expected.ID {
			t.Errorf("book %d: expected ID %d, got %d", i, expected.ID, book.ID)
		}
		if book.Title != expected.Title {
			t.Errorf("book %d: expected title %q, got %q", i, expected.Title, book.Title)
		}
		if book.Slug != expected.Slug {
			t.Errorf("book %d: expected slug %q, got %q", i, expected.Slug, book.Slug)
		}

		// Check rating (all books have ratings in this response)
		if response.Books[i].Rating == nil {
			t.Errorf("book %d: expected rating %v, got nil", i, expected.Rating)
		} else if *response.Books[i].Rating != expected.Rating {
			t.Errorf("book %d: expected rating %v, got %v", i, expected.Rating, *response.Books[i].Rating)
		}

		// Check last read date
		if response.Books[i].LastReadDate == nil {
			t.Errorf("book %d: expected last_read_date %s, got nil", i, expected.LastReadDate)
		} else {
			// The Date type embeds time.Time, we need to check it matches
			expectedDate, _ := time.Parse("2006-01-02", expected.LastReadDate)
			if !response.Books[i].LastReadDate.Equal(expectedDate) {
				t.Errorf("book %d: expected last_read_date %s, got %v", i, expected.LastReadDate, response.Books[i].LastReadDate.Time)
			}
		}

		// All books should have images
		if book.Image == nil || book.Image.URL == "" {
			t.Errorf("book %d: missing image URL", i)
		}
	}
}

func TestClientHandlesErrorResponses(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		responseBody  string
		expectedError string
	}{
		{
			name:          "404 not found",
			statusCode:    404,
			responseBody:  "Not Found",
			expectedError: "API request failed with status 404",
		},
		{
			name:          "500 server error",
			statusCode:    500,
			responseBody:  "Internal Server Error",
			expectedError: "API request failed with status 500",
		},
		{
			name:          "GraphQL errors",
			statusCode:    200,
			responseBody:  `{"errors": [{"message": "User not found"}], "data": {"user_books": []}}`,
			expectedError: "GraphQL errors:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTP := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: tt.statusCode,
						Body:       io.NopCloser(bytes.NewReader([]byte(tt.responseBody))),
					}, nil
				},
			}

			client := NewClientWithHTTPClient("test-token", mockHTTP)
			_, err := client.GetUserBooksByUsername("testuser")

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !containsString(err.Error(), tt.expectedError) {
				t.Errorf("expected error containing %q, got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestClientHandlesEmptyResponse(t *testing.T) {
	emptyResponse := `{
		"data": {
			"user_books": []
		}
	}`

	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(emptyResponse))),
			}, nil
		},
	}

	client := NewClientWithHTTPClient("test-token", mockHTTP)
	response, err := client.GetUserBooksByUsername("testuser")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.Count != 0 {
		t.Errorf("expected 0 books, got %d", response.Count)
	}

	if len(response.Books) != 0 {
		t.Errorf("expected empty books array, got %d books", len(response.Books))
	}
}

func TestTimeParsing(t *testing.T) {
	// Test that our time parsing handles the API's time format
	timeStr := "2025-07-09T22:22:43.459059+00:00"
	parsedTime, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		t.Fatalf("failed to parse time: %v", err)
	}

	// Verify it parsed correctly
	if parsedTime.Year() != 2025 || parsedTime.Month() != 7 || parsedTime.Day() != 9 {
		t.Errorf("time parsed incorrectly: %v", parsedTime)
	}
}

// Helper function
func containsString(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
