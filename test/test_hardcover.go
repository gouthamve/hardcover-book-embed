package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gouthamve/hardcover-book-embed/internal/hardcover"
)

func main() {
	apiToken := os.Getenv("HARDCOVER_API_TOKEN")
	if apiToken == "" {
		log.Fatal("HARDCOVER_API_TOKEN environment variable is required")
	}

	if len(os.Args) < 2 {
		log.Fatal("Usage: test_hardcover <username> [book-type]")
	}

	username := os.Args[1]
	bookType := "currently-reading"
	if len(os.Args) >= 3 {
		bookType = os.Args[2]
	}

	client := hardcover.NewClient(apiToken)

	fmt.Printf("Fetching %s books for user: %s\n", bookType, username)
	fmt.Println("API Token:", apiToken[:10]+"...")
	fmt.Println()

	var books *hardcover.UserBooksResponse
	var err error

	if bookType == "last-read" {
		books, err = client.GetUserLastReadBooksByUsername(username)
	} else {
		books, err = client.GetUserBooksByUsername(username)
	}
	if err != nil {
		log.Fatalf("Error fetching books: %v", err)
	}

	fmt.Printf("Successfully fetched %d %s books:\n", books.Count, bookType)
	fmt.Println("========================================")

	if books.Count == 0 {
		if bookType == "last-read" {
			fmt.Println("No recently read books.")
		} else {
			fmt.Println("No books currently being read.")
		}
		return
	}

	for i, userBook := range books.Books {
		fmt.Printf("\n%d. %s\n", i+1, userBook.Book.Title)

		if userBook.Rating != nil {
			fmt.Printf("   Rating: %.1f/5 stars\n", *userBook.Rating)
		}

		if userBook.Book.Image != nil && userBook.Book.Image.URL != "" {
			if strings.Contains(userBook.Book.Image.URL, "static/covers/cover") {
				fmt.Printf("   Cover: %s (fallback)\n", userBook.Book.Image.URL)
			} else {
				fmt.Printf("   Cover: %s\n", userBook.Book.Image.URL)
			}
		}

		fmt.Printf("   URL: https://hardcover.app/books/%s\n", userBook.Book.Slug)
		fmt.Printf("   Updated: %s\n", userBook.UpdatedAt.Format("2006-01-02 15:04:05"))

		if bookType == "last-read" && userBook.LastReadDate != nil && !userBook.LastReadDate.IsZero() {
			fmt.Printf("   Last Read: %s\n", userBook.LastReadDate.Format("2006-01-02"))
		}
	}

	fmt.Println("\n========================================")
	fmt.Println("Raw JSON Response:")
	jsonData, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
	} else {
		fmt.Println(string(jsonData))
	}
}
