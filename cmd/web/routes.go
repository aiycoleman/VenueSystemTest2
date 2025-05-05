// filename: routes.go
// Description: Maps specific URL paths to their corresponding handler functions

package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {

	// Create a new ServeMux (multiplexer) to manage request routing.
	mux := http.NewServeMux()

	// File server for serving static files (CSS, JS, images)
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Serve static files from the "/static" path.
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Route for the all pages, mapping the root URL to the their respective handler.
	mux.HandleFunc("GET /{$}", app.home)

	// User Authentication
	// mux

	// Wrap the router with a logging middleware to track requests.
	return app.session.Enable(app.loggingMiddleware(mux))
}
