// filename: middleware.go
// Description: Logs incoming HTTP requests and their details before passing them to the next handler

package main

import (
	"net/http"

	"github.com/justinas/nosurf"
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

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				app.logger.Error("panic recovered", "error", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Info("incoming request", "method", r.Method, "url", r.URL.String(), "remote", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// is authenticated is not initialized

// func (app *application) requireAuthentication(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if !app.isAuthenticated(r) {
// 			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
// 			return
// 		}
// 		w.Header().Add("Cache-Control", "no-store")
// 		next.ServeHTTP(w, r)
// 	})
// }

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}
