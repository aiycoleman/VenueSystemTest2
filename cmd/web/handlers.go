// filename: handlers.go
// Description: Handling HTTP requests

package main

import (
	"net/http"
)

// ------------------------------- Home Handler --------------------------------
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Welcome to Venue Verge!"
	data.HeaderText = "Find Your Perfect Venue!"

	err := app.render(w, http.StatusOK, "home.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render home page", "template", "home.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
