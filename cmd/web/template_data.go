// filename: tempalte_data.go
// Description: Container for dynamic content passed to HTML templates

package main

import (
	"net/http"

	"github.com/aiycoleman/VenueSystemTest2/internal/data"
	"github.com/justinas/nosurf"
)

// Holds dynamic data that can be passed to HTML templates.
type TemplateData struct {
	Title           string
	HeaderText      string
	Flash           string
	CSRFToken       string
	Venue           *data.Venue
	Venues          []data.Venue
	Reservation     []data.Reservation
	Reviews         []data.Review
	FormErrors      map[string]string
	FormData        map[string]string
	IsAuthenticated bool
}

// Initializes a new TemplateData struct with default values.
func NewTemplateData(r *http.Request) *TemplateData {
	return &TemplateData{
		Title:      "Default Title",
		HeaderText: "Default HeaderText",
		CSRFToken:  nosurf.Token(r),
		// Flash: string,
		Venues:          []data.Venue{},
		Reservation:     []data.Reservation{},
		Reviews:         []data.Review{},
		FormErrors:      map[string]string{},
		FormData:        map[string]string{},
		IsAuthenticated: false,
	}
}
