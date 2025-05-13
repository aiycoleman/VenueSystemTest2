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
		app.authenticate,
		// noSurf,
	)

	// Public routes
	mux.Handle("GET /{$}", dynamicMiddleware.ThenFunc(app.home))

	mux.Handle("GET /user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Handle("POST /user/signup", dynamicMiddleware.ThenFunc(app.signupUser))

	mux.Handle("GET /user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Handle("POST /user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Handle("POST /user/logout", dynamicMiddleware.ThenFunc(app.logoutUser))

	// Protected routes - require authentication
	protected := dynamicMiddleware.Append(app.requireAuthentication)

	// Role-based access: owner role
	ownerProtected := protected.Append(app.requireRole(1))

	// Role-based access: User role
	userProtected := protected.Append(app.requireRole(2))

	// Public routes accessible by anyone
	mux.Handle("GET /venue/listing", protected.ThenFunc(app.venueListing))  // Access to book, add, edit, delete
	mux.Handle("GET /venue/form", ownerProtected.ThenFunc(app.venueForm))   // Only accessible by owner
	mux.Handle("POST /venue/add", ownerProtected.ThenFunc(app.createVenue)) // Only accessible by owner
	mux.Handle("GET /venue/{id}", protected.ThenFunc(app.viewVenue))

	mux.Handle("GET /venue/{id}/edit", ownerProtected.ThenFunc(app.showUpdateVenueForm)) // owner only
	mux.Handle("POST /venue/{id}/edit", ownerProtected.ThenFunc(app.updateVenue))        // owner only
	mux.Handle("POST /venue/{id}/delete", ownerProtected.ThenFunc(app.deleteVenue))      // owner only

	mux.Handle("POST /reservation/{id}/create", userProtected.ThenFunc(app.createReservation)) // User only

	mux.Handle("GET /reservations", userProtected.ThenFunc(app.showAllReservations))                 // User only
	mux.Handle("GET /reservations/cancelled", userProtected.ThenFunc(app.showCancelledReservations)) // User only

	mux.Handle("GET /reservations/update/{id}", userProtected.ThenFunc(app.showUpdateReservationForm)) // User only
	mux.Handle("POST /reservations/update/{id}", userProtected.ThenFunc(app.updateReservation))        // User only

	mux.Handle("POST /reservations/cancel/{id}", userProtected.ThenFunc(app.cancelReservation)) // User only

	mux.Handle("POST /venue/{id}/review", userProtected.ThenFunc(app.submitReview)) // Accessible by both user and owner

	// Final handler with outermost middleware
	return app.session.Enable(standardMiddleware.Then(mux))
}
