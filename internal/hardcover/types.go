package hardcover

import "time"

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
	Rating    *int      `json:"rating,omitempty"`
	Book      Book      `json:"book"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GraphQLResponse struct {
	Data struct {
		Me []struct {
			UserBooks []UserBook `json:"user_books"`
		} `json:"me"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

type CurrentlyReadingResponse struct {
	Books     []UserBook `json:"books"`
	Count     int        `json:"count"`
	UpdatedAt time.Time  `json:"updated_at"`
}
