// filename: routes.go
// Description: Maps specific URL paths to their corresponding handler functions

package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Static file server
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Base middleware chain: Logging + CSRF + Sessions
	standardMiddleware := alice.New(
		app.recoverPanic,
		app.logRequest,
		app.secureHeaders,
	)
	dynamicMiddleware := alice.New(
		app.session.Enable,
		app.loggingMiddleware,
		noSurf,
	)

	// Public routes
	mux.Handle("GET /{$}", dynamicMiddleware.ThenFunc(app.home))

	mux.Handle("GET /user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Handle("POST /user/signup", dynamicMiddleware.ThenFunc(app.signupUser))

	mux.Handle("GET /user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Handle("POST /user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Handle("POST /user/logout", dynamicMiddleware.ThenFunc(app.logoutUser))

	// Example of protected route (uncomment when requireAuthentication is implemented)
	// protected := dynamicMiddleware.Append(app.requireAuthentication)
	// mux.Handle("GET /dashboard", protected.ThenFunc(app.dashboard))

	// Final handler with outermost middleware
	return standardMiddleware.Then(mux)
}
