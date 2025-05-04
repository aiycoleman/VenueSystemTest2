// filename: tempalte_data.go
// Description: Container for dynamic content passed to HTML templates

package main

//"github.com/aiycoleman/VenueSystemTest2/internal/data"

// Holds dynamic data that can be passed to HTML templates.
type TemplateData struct {
	Title      string
	HeaderText string
	// User       []data.Users
	// Venue       *data.Venue
	// Venues      []data.Venue
	// Reservation []data.Reservation
	// Reviews     []data.Review
	// FormErrors  map[string]string
	// FormData    map[string]string
}

// Initializes a new TemplateData struct with default values.
func NewTemplateData() *TemplateData {
	return &TemplateData{
		Title:      "Default Title",
		HeaderText: "Default HeaderText",
		//User:       []data.Users{},
		// Venues:      []data.Venue{},
		// Reservation: []data.Reservation{},
		// Reviews:     []data.Review{},
		// FormErrors:  map[string]string{},
		// FormData:    map[string]string{},
	}
}
