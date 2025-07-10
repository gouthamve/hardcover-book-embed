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

	client := hardcover.NewClient(apiToken)

	fmt.Println("Fetching currently reading books from Hardcover API...")
	fmt.Println("API Token:", apiToken[:10]+"...")
	fmt.Println()

	books, err := client.GetCurrentlyReadingBooks()
	if err != nil {
		log.Fatalf("Error fetching books: %v", err)
	}

	fmt.Printf("Successfully fetched %d currently reading books:\n", books.Count)
	fmt.Println("========================================")

	if books.Count == 0 {
		fmt.Println("No books currently being read.")
		return
	}

	for i, userBook := range books.Books {
		fmt.Printf("\n%d. %s\n", i+1, userBook.Book.Title)

		if userBook.Rating != nil {
			fmt.Printf("   Rating: %d/5 stars\n", *userBook.Rating)
		}

		if userBook.Book.Image != nil && userBook.Book.Image.URL != "" {
			if strings.Contains(userBook.Book.Image.URL, "static/covers/cover") {
				fmt.Printf("   Cover: %s (fallback)\n", userBook.Book.Image.URL)
			} else {
				fmt.Printf("   Cover: %s\n", userBook.Book.Image.URL)
			}
		}

		fmt.Printf("   Updated: %s\n", userBook.UpdatedAt.Format("2006-01-02 15:04:05"))
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
