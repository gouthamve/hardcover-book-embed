package hardcover

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	HardcoverAPIURL = "https://api.hardcover.app/v1/graphql"
	UserAgent       = "hardcover-book-embed/1.0"
)

type Client struct {
	apiToken   string
	httpClient *http.Client
}

func NewClient(apiToken string) *Client {
	return &Client{
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) GetCurrentlyReadingBooks() (*CurrentlyReadingResponse, error) {
	query := `{
		me {
			user_books(where: {status_id: {_eq: 2}}) {
				rating
				updated_at
				book {
					id
					title
					image {
						url
					}
				}
			}
		}
	}`

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

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var graphqlResp GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&graphqlResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(graphqlResp.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %v", graphqlResp.Errors)
	}

	// Handle empty Me array
	if len(graphqlResp.Data.Me) == 0 {
		return &CurrentlyReadingResponse{
			Books:     []UserBook{},
			Count:     0,
			UpdatedAt: time.Now(),
		}, nil
	}

	// Process books and add fallback images
	books := graphqlResp.Data.Me[0].UserBooks
	for i := range books {
		if books[i].Book.Image == nil {
			// Generate a fallback image based on book ID
			coverNum := (books[i].Book.ID % 9) + 1
			books[i].Book.Image = &Image{
				URL: fmt.Sprintf("https://assets.hardcover.app/static/covers/cover%d.webp", coverNum),
			}
		}
	}

	return &CurrentlyReadingResponse{
		Books:     books,
		Count:     len(books),
		UpdatedAt: time.Now(),
	}, nil
}
