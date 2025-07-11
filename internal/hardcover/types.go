package hardcover

import (
	"strings"
	"time"
)

// Date is a custom type that can unmarshal date-only format (YYYY-MM-DD)
type Date struct {
	time.Time
}

// UnmarshalJSON handles date-only format
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		d.Time = time.Time{}
		return nil
	}

	// Try parsing as date-only format first
	t, err := time.Parse("2006-01-02", s)
	if err == nil {
		d.Time = t
		return nil
	}

	// Fall back to full timestamp format
	t, err = time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

type Image struct {
	URL string `json:"url"`
}

type Book struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
	Image *Image `json:"image,omitempty"`
}

type Contributor struct {
	Name string `json:"name"`
}

type UserBook struct {
	Rating       *float64  `json:"rating,omitempty"`
	Book         Book      `json:"book"`
	UpdatedAt    time.Time `json:"updated_at"`
	LastReadDate *Date     `json:"last_read_date,omitempty"`
}

type UserBooksAPIResponse struct {
	Data struct {
		UserBooks []UserBook `json:"user_books"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

type UserBooksResponse struct {
	Books     []UserBook `json:"books"`
	Count     int        `json:"count"`
	UpdatedAt time.Time  `json:"updated_at"`
}
