// filename: middleware.go
// Description: Logs incoming HTTP requests and their details before passing them to the next handler

package main

import (
	"net/http"
)

func (app *application) loggingMiddleware(next http.Handler) http.Handler {
	// Define a handler function that wraps the provided handler.
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Extract relevant request details for logging.
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		// Log the incoming request details.
		app.logger.Info("received request", "ip", ip, "protocol", proto, "method", method, "uri", uri)
		// Call the next handler in the chain to process the request.
		next.ServeHTTP(w, r)
		// Log a message after the request has been processed.
		app.logger.Info("Request processed")
	})

	// Return the wrapped handler function as an http.Handler.
	return fn

}
