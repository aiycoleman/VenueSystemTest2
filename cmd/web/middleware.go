// filename: middleware.go
// Description: Logs incoming HTTP requests and their details before passing them to the next handler

package main

import (
	"context"
	"net/http"

	"github.com/aiycoleman/VenueSystemTest2/internal/data"
	"github.com/justinas/nosurf"
)

type contextKey string

const contextKeyUser = contextKey("user")
const contextKeyIsAuthenticated = contextKey("isAuthenticated")

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

// Middleware to check if the user is authenticated
func (app *application) isAuthenticated(r *http.Request) bool {
	isAuth, ok := r.Context().Value(contextKeyIsAuthenticated).(bool)
	return ok && isAuth
}

// Middleware to enforce authentication
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

// Middleware to check the role of the authenticated user
func (app *application) requireRole(role int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := app.contextGetUser(r.Context())
			if user == nil {
				http.Redirect(w, r, "/unauthorized", http.StatusSeeOther)
				return
			}

			app.logger.Info("Checking user role",
				"userID", user.ID,
				"userRole", user.Role,
				"requiredRole", role,
			)

			// Convert role to int64 for comparison
			if user.Role != int64(role) {
				http.Redirect(w, r, "/unauthorized", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Helper function to get user from request context
// This is your context function for retrieving the user from context
func (app *application) contextGetUser(ctx context.Context) *data.Users {
	user, ok := ctx.Value(contextKeyUser).(*data.Users)
	if !ok {
		return nil
	}
	return user
}

// CSRF protection middleware
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}

// Middleware to authenticate and set user in context
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := app.session.Get(r, "authenticatedUserID").(int)
		if !ok {
			next.ServeHTTP(w, r) // Not authenticated — let the next handler decide
			return
		}

		user, err := app.users.Get(id)
		if err != nil {
			// Optionally clear session on error
			app.session.Remove(r, "authenticatedUserID")
			next.ServeHTTP(w, r)
			return
		}

		// Store user in the request context
		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		ctx = context.WithValue(ctx, contextKeyIsAuthenticated, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
