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
func TestClientParsesCurrentlyReadingAPIResponseCorrectly(t *testing.T) {
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
	response, err := client.GetUserCurrentlyReadingBooksByUsername("testuser")
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
			_, err := client.GetUserCurrentlyReadingBooksByUsername("testuser")

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
	response, err := client.GetUserCurrentlyReadingBooksByUsername("testuser")

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

func TestClientParsesReviewsAPIResponse(t *testing.T) {
	// Actual API response for user reviews
	apiResponse := `{
  "data": {
    "user_books": [
      {
        "review_length": 45,
        "review_raw": "Neeeeeeeew achievement! Read my first LitRPG!",
        "reviewed_at": "2025-04-23T00:00:00",
        "has_review": true,
        "review_has_spoilers": false,
        "review_html": null,
        "rating": 5,
        "book": {
          "id": 446681,
          "title": "Dungeon Crawler Carl",
          "image": {
            "url": "https://assets.hardcover.app/edition/31601422/2fd7c35a3c3ea037dea2e38a1b48964f1d6b7218.jpeg"
          },
          "slug": "dungeon-crawler-carl",
          "contributions": [
            {
              "author": {
                "name": "Matt Dinniman",
                "links": [],
                "slug": "matt-dinniman"
              }
            }
          ]
        },
        "review_object": [],
        "review_slate": {
          "document": {
            "object": "document",
            "children": [
              {
                "data": {},
                "type": "paragraph",
                "object": "block",
                "children": [
                  {
                    "text": "Neeeeeeeew achievement! Read my first LitRPG!",
                    "object": "text"
                  }
                ]
              }
            ]
          }
        },
        "url": null
      },
      {
        "review_length": 436,
        "review_raw": "Sometimes a book comes along with exactly what you need when you need it. This one helped me reflect to on my own mindset about productivity and make adjustments to be happier and healthier.The main premise of the book is stated early (on page 8):A philosophy for organizing knowledge work efforts in a sustainable and meaningful manner based on the following three principles.Do fewer things.Work at a natural pace.Obsess over quality.",
        "reviewed_at": "2024-04-02T00:00:00",
        "has_review": true,
        "review_has_spoilers": false,
        "review_html": null,
        "rating": 4.5,
        "book": {
          "id": 898371,
          "title": "Slow Productivity: The Lost Art of Accomplishment Without Burnout",
          "image": {
            "url": "https://assets.hardcover.app/external_data/60702703/983bb1e60d94ccec2e4f80c90715b93ffaccdc04.jpeg"
          },
          "slug": "slow-productivity",
          "contributions": [
            {
              "author": {
                "name": "Cal Newport",
                "links": [],
                "slug": "cal-newport"
              }
            }
          ]
        },
        "review_object": [],
        "review_slate": {
          "document": {
            "object": "document",
            "children": [
              {
                "data": {},
                "type": "paragraph",
                "object": "block",
                "children": [
                  {
                    "text": "Sometimes a book comes along with exactly what you need when you need it. This one helped me reflect to on my own mindset about productivity and make adjustments to be happier and healthier.",
                    "object": "text"
                  }
                ]
              }
            ]
          }
        },
        "url": null
      },
      {
        "review_length": 560,
        "review_raw": "After reading a few other magic school books this year (The Will of the Many, The Scholomance), I wasn't sure this one would live up to the hype of being the #1 trending book on Hardcover. Turns out it did.Fourth Wing takes place in a cutthroat school for dragon riders. Students learn the skills needed to defend their homeland from invading forces and protect society.At times it reminded me of The Hunger Games, LOTR and others in the dark-academia genre while still managing to be original enough to keep me wondering. Sign me up for the next in the series.",
        "reviewed_at": "2023-10-18T17:29:19.393223",
        "has_review": true,
        "review_has_spoilers": false,
        "review_html": null,
        "rating": 5,
        "book": {
          "id": 714600,
          "title": "Fourth Wing",
          "image": {
            "url": "https://assets.hardcover.app/editions/30707731/3559167047761380.jpeg"
          },
          "slug": "fourth-wing",
          "contributions": [
            {
              "author": {
                "name": "Rebecca Yarros",
                "links": [],
                "slug": "rebecca-yarros"
              }
            }
          ]
        },
        "review_object": [],
        "review_slate": {
          "document": {
            "object": "document",
            "children": []
          }
        },
        "url": null
      }
    ]
  }
}`

	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Verify request details
			if req.Method != "POST" {
				t.Errorf("expected POST request, got %s", req.Method)
			}
			if req.URL.String() != HardcoverAPIURL {
				t.Errorf("expected URL %s, got %s", HardcoverAPIURL, req.URL.String())
			}

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(apiResponse))),
			}, nil
		},
	}

	client := NewClientWithHTTPClient("test-token", mockHTTP)
	response, err := client.GetUserReviewsByUsername("testuser")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify response count
	if response.Count != 3 {
		t.Errorf("expected 3 reviews, got %d", response.Count)
	}

	// Expected data
	expectedReviews := []struct {
		BookID            int
		Title             string
		Slug              string
		AuthorName        string
		AuthorSlug        string
		Rating            float64
		ReviewLength      int
		ReviewRaw         string
		ReviewedAt        string
		HasReview         bool
		ReviewHasSpoilers bool
		HasReviewHTML     bool
	}{
		{
			BookID:            446681,
			Title:             "Dungeon Crawler Carl",
			Slug:              "dungeon-crawler-carl",
			AuthorName:        "Matt Dinniman",
			AuthorSlug:        "matt-dinniman",
			Rating:            5,
			ReviewLength:      45,
			ReviewRaw:         "Neeeeeeeew achievement! Read my first LitRPG!",
			ReviewedAt:        "2025-04-23T00:00:00",
			HasReview:         true,
			ReviewHasSpoilers: false,
			HasReviewHTML:     false,
		},
		{
			BookID:            898371,
			Title:             "Slow Productivity: The Lost Art of Accomplishment Without Burnout",
			Slug:              "slow-productivity",
			AuthorName:        "Cal Newport",
			AuthorSlug:        "cal-newport",
			Rating:            4.5,
			ReviewLength:      436,
			ReviewRaw:         "Sometimes a book comes along with exactly what you need when you need it. This one helped me reflect to on my own mindset about productivity and make adjustments to be happier and healthier.The main premise of the book is stated early (on page 8):A philosophy for organizing knowledge work efforts in a sustainable and meaningful manner based on the following three principles.Do fewer things.Work at a natural pace.Obsess over quality.",
			ReviewedAt:        "2024-04-02T00:00:00",
			HasReview:         true,
			ReviewHasSpoilers: false,
			HasReviewHTML:     false,
		},
		{
			BookID:            714600,
			Title:             "Fourth Wing",
			Slug:              "fourth-wing",
			AuthorName:        "Rebecca Yarros",
			AuthorSlug:        "rebecca-yarros",
			Rating:            5,
			ReviewLength:      560,
			ReviewRaw:         "After reading a few other magic school books this year (The Will of the Many, The Scholomance), I wasn't sure this one would live up to the hype of being the #1 trending book on Hardcover. Turns out it did.Fourth Wing takes place in a cutthroat school for dragon riders. Students learn the skills needed to defend their homeland from invading forces and protect society.At times it reminded me of The Hunger Games, LOTR and others in the dark-academia genre while still managing to be original enough to keep me wondering. Sign me up for the next in the series.",
			ReviewedAt:        "2023-10-18T17:29:19.393223",
			HasReview:         true,
			ReviewHasSpoilers: false,
			HasReviewHTML:     false,
		},
	}

	// Verify each review
	for i, expected := range expectedReviews {
		if i >= len(response.Books) {
			t.Errorf("missing review at index %d", i)
			continue
		}

		review := response.Books[i]
		book := review.Book

		// Verify book details
		if book.ID != expected.BookID {
			t.Errorf("review %d: expected book ID %d, got %d", i, expected.BookID, book.ID)
		}
		if book.Title != expected.Title {
			t.Errorf("review %d: expected title %q, got %q", i, expected.Title, book.Title)
		}
		if book.Slug != expected.Slug {
			t.Errorf("review %d: expected slug %q, got %q", i, expected.Slug, book.Slug)
		}

		// Verify rating
		if review.Rating == nil {
			t.Errorf("review %d: expected rating %v, got nil", i, expected.Rating)
		} else if *review.Rating != expected.Rating {
			t.Errorf("review %d: expected rating %v, got %v", i, expected.Rating, *review.Rating)
		}

		// Verify review content
		if review.ReviewLength == nil {
			t.Errorf("review %d: expected review_length %d, got nil", i, expected.ReviewLength)
		} else if *review.ReviewLength != expected.ReviewLength {
			t.Errorf("review %d: expected review_length %d, got %d", i, expected.ReviewLength, *review.ReviewLength)
		}

		if review.ReviewRaw == nil {
			t.Errorf("review %d: expected review_raw, got nil", i)
		} else if *review.ReviewRaw != expected.ReviewRaw {
			t.Errorf("review %d: expected review_raw %q, got %q", i, expected.ReviewRaw, *review.ReviewRaw)
		}

		// Verify review metadata
		if review.HasReview != expected.HasReview {
			t.Errorf("review %d: expected has_review %v, got %v", i, expected.HasReview, review.HasReview)
		}
		if review.ReviewHasSpoilers != expected.ReviewHasSpoilers {
			t.Errorf("review %d: expected review_has_spoilers %v, got %v", i, expected.ReviewHasSpoilers, review.ReviewHasSpoilers)
		}

		// Verify reviewed_at date parsing
		if review.ReviewedAt == nil {
			t.Errorf("review %d: expected reviewed_at, got nil", i)
		} else {
			// The API returns different date formats, test that our Date type handles them
			var expectedTime time.Time
			var err error

			// Try parsing with different formats based on the data
			if containsString(expected.ReviewedAt, ".") {
				// Format with microseconds
				expectedTime, err = time.Parse("2006-01-02T15:04:05.999999", expected.ReviewedAt)
			} else {
				// Format without timezone
				expectedTime, err = time.Parse("2006-01-02T15:04:05", expected.ReviewedAt)
			}

			if err != nil {
				t.Errorf("review %d: failed to parse expected date %q: %v", i, expected.ReviewedAt, err)
			} else {
				// Compare just the date/time components, not the exact time.Time object
				if !review.ReviewedAt.Equal(expectedTime) {
					// Allow for small differences in parsing
					diff := review.ReviewedAt.Sub(expectedTime)
					if diff < -time.Second || diff > time.Second {
						t.Errorf("review %d: expected reviewed_at %v, got %v (diff: %v)", i, expectedTime, review.ReviewedAt.Time, diff)
					}
				}
			}
		}

		// Verify review_html is null (as per test data)
		if expected.HasReviewHTML && review.ReviewHTML != nil {
			t.Errorf("review %d: expected review_html to be null, got %v", i, *review.ReviewHTML)
		}

		// Verify all books have images
		if book.Image == nil || book.Image.URL == "" {
			t.Errorf("review %d: missing image URL", i)
		}

		// Verify author details
		if book.Contributions == nil || len(book.Contributions) == 0 {
			t.Errorf("review %d: missing contributions", i)
		} else {
			// Verify the first author (assuming single author for these test books)
			author := book.Contributions[0].Author
			if author.Name != expected.AuthorName {
				t.Errorf("review %d: expected author name %q, got %q", i, expected.AuthorName, author.Name)
			}
			if author.Slug != expected.AuthorSlug {
				t.Errorf("review %d: expected author slug %q, got %q", i, expected.AuthorSlug, author.Slug)
			}
			// Verify links field exists (even if empty)
			if author.Links == nil {
				t.Errorf("review %d: expected author links to be non-nil", i)
			}
		}

		// Verify review_slate exists (basic check)
		if review.ReviewSlate == nil {
			t.Errorf("review %d: expected review_slate, got nil", i)
		}

		// Verify review_object exists
		if review.ReviewObject == nil {
			t.Errorf("review %d: expected review_object, got nil", i)
		}

		// Verify URL is empty string or null
		if review.URL != "" {
			t.Errorf("review %d: expected empty URL, got %q", i, review.URL)
		}
	}
}

// Test that reviews with spoilers are handled correctly
func TestClientHandlesReviewsWithSpoilers(t *testing.T) {
	apiResponse := `{
		"data": {
			"user_books": [
				{
					"review_length": 100,
					"review_raw": "Major plot twist revealed!",
					"reviewed_at": "2024-01-01T00:00:00",
					"has_review": true,
					"review_has_spoilers": true,
					"review_html": "<p>Major plot twist revealed!</p>",
					"rating": 4,
					"book": {
						"id": 123,
						"title": "Test Book",
						"image": {
							"url": "https://example.com/cover.jpg"
						},
						"slug": "test-book",
						"contributions": [
							{
								"author": {
									"name": "Test Author",
									"links": [],
									"slug": "test-author"
								}
							}
						]
					},
					"review_object": [],
					"review_slate": {},
					"url": "https://hardcover.app/books/test-book"
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
	response, err := client.GetUserReviewsByUsername("testuser")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(response.Books) != 1 {
		t.Fatalf("expected 1 review, got %d", len(response.Books))
	}

	review := response.Books[0]
	if !review.ReviewHasSpoilers {
		t.Error("expected review_has_spoilers to be true")
	}

	// Verify review_html is populated
	if review.ReviewHTML == nil {
		t.Error("expected review_html to be populated")
	} else if *review.ReviewHTML != "<p>Major plot twist revealed!</p>" {
		t.Errorf("expected review_html %q, got %q", "<p>Major plot twist revealed!</p>", *review.ReviewHTML)
	}

	// Verify URL is populated
	if review.URL != "https://hardcover.app/books/test-book" {
		t.Errorf("expected URL %q, got %q", "https://hardcover.app/books/test-book", review.URL)
	}
}

// Helper function
func containsString(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
