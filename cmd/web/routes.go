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
		// noSurf,
	)

	// Public routes
	mux.Handle("GET /{$}", dynamicMiddleware.ThenFunc(app.home))

	mux.Handle("GET /user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Handle("POST /user/signup", dynamicMiddleware.ThenFunc(app.signupUser))

	mux.Handle("GET /user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Handle("POST /user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Handle("POST /user/logout", dynamicMiddleware.ThenFunc(app.logoutUser))

	mux.Handle("GET /venue/listing", dynamicMiddleware.ThenFunc(app.venueListing)) // access you book, add, edit, delete
	mux.Handle("GET /venue/form", dynamicMiddleware.ThenFunc(app.venueForm))       // Venue form to add new venue
	mux.Handle("POST /venue/add", dynamicMiddleware.ThenFunc(app.createVenue))
	mux.Handle("GET /venue/{id}", dynamicMiddleware.ThenFunc(app.viewVenue))

	mux.Handle("GET /venue/{id}/edit", dynamicMiddleware.ThenFunc(app.showUpdateVenueForm))
	mux.Handle("POST /venue/{id}/edit", dynamicMiddleware.ThenFunc(app.updateVenue))
	mux.Handle("POST /venue/{id}/delete", dynamicMiddleware.ThenFunc(app.deleteVenue))

	mux.Handle("POST /reservation/{id}/create", dynamicMiddleware.ThenFunc(app.createReservation))

	mux.Handle("GET /reservations", dynamicMiddleware.ThenFunc(app.showAllReservations))
	mux.Handle("GET /reservations/cancelled", dynamicMiddleware.ThenFunc(app.showCancelledReservations))

	mux.Handle("GET /reservations/update/{id}", dynamicMiddleware.ThenFunc(app.showUpdateReservationForm))
	mux.Handle("POST /reservations/update/{id}", dynamicMiddleware.ThenFunc(app.updateReservation))

	mux.Handle("POST /reservations/cancel/{id}", dynamicMiddleware.ThenFunc(app.cancelReservation))

	mux.Handle("POST /venue/{id}/review", dynamicMiddleware.ThenFunc(app.submitReview))

	// Example of protected route (uncomment when requireAuthentication is implemented)
	// protected := dynamicMiddleware.Append(app.requireAuthentication)
	// mux.Handle("GET /dashboard", protected.ThenFunc(app.dashboard))

	// Final handler with outermost middleware
	return standardMiddleware.Then(mux)
}
