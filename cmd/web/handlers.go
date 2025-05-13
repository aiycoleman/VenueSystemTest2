// filename: handlers.go
// Description: Handling HTTP requests

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/aiycoleman/VenueSystemTest2/internal/data"
	"github.com/aiycoleman/VenueSystemTest2/internal/validator"
)

// ------------------------------- Home Handler --------------------------------
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Log the session value for debugging
	app.logger.Info("Session userID", "value", app.session.Get(r, "authenticatedUserID"))

	data := NewTemplateData(r)
	data.Title = "Welcome to Venue Verge!"
	data.HeaderText = "Find Your Perfect Venue!"
	data.IsAuthenticated = app.isAuthenticated(r)

	err := app.render(w, http.StatusOK, "home.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render home page", "template", "home.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData(r)
	data.Title = "Sign Up Today!"
	data.HeaderText = "Create Account"

	err := app.render(w, http.StatusOK, "signup.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render signup page", "template", "signup.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Extract form values
	name := r.FormValue("name")
	role := r.FormValue("role")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Create a new validator instance
	v := validator.NewValidator()

	// Validate password before hashing
	v.Check(validator.NotBlank(password), "password", "must be provided")
	v.Check(validator.MinLength(password, 10), "password", "must be at least 10 characters long")
	v.Check(validator.MaxLength(password, 50), "password", "must not be more than 50 characters long")

	// Validate other fields
	v.Check(validator.NotBlank(name), "name", "must be provided")
	v.Check(validator.MaxLength(name, 50), "name", "must not be more than 50 characters long")

	v.Check(validator.NotBlank(email), "email", "must be provided")
	v.Check(validator.IsValidEmail(email), "email", "invalid email address")
	v.Check(validator.MaxLength(email, 100), "email", "must not be more than 100 characters long")

	v.Check(validator.IsValidChoice(role), "role", "must be provided")

	// If validation fails, re-render form with errors
	if !v.ValidData() {
		formData := map[string]string{
			"name":  name,
			"role":  role,
			"email": email,
		}

		td := NewTemplateData(r)
		td.Title = "Sign Up"
		td.HeaderText = "Create Account"
		td.FormErrors = v.Errors
		td.FormData = formData

		err = app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", td)
		if err != nil {
			app.logger.Error("failed to render signup form", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Hash the password after validation passes
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		app.logger.Error("failed to hash password", "error", err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	// Convert role to integer (assuming it's valid)
	roleInt, err := strconv.ParseInt(role, 10, 64)
	if err != nil {
		app.logger.Error("invalid role value", "error", err)
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	// Create the user struct
	users := &data.Users{
		Name:           name,
		Role:           roleInt,
		Email:          email,
		HashedPassword: hashedpassword,
	}

	// Insert user into the database
	err = app.users.Insert(users)
	if err != nil {
		if errors.Is(err, data.ErrDuplicateEmail) {
			v.AddError("email", "A user with this email already exists")

			formData := map[string]string{
				"name":  name,
				"role":  role,
				"email": email,
			}

			td := NewTemplateData(r)
			td.Title = "Sign Up"
			td.HeaderText = "Create Account"
			td.FormErrors = v.Errors
			td.FormData = formData

			err = app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", td)
			if err != nil {
				app.logger.Error("failed to render signup form", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		app.logger.Error("failed to insert user", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set session flash message and redirect to login
	app.session.Put(r, "flash", "Signup was successful.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData(r)
	data.Title = "Hello, Nice to See You Again!"
	data.HeaderText = "Login"

	err := app.render(w, http.StatusOK, "signin.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render signin page", "template", "signin.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Extract form values
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Check the web form fields to validity
	errors_user := make(map[string]string)
	id, err := app.users.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, data.ErrInvalidCredentials) {
			errors_user["default"] = "Email or Password is incorrect"

			td := NewTemplateData(r)
			td.Title = "Hello, Nice to See You Again!"
			td.HeaderText = "Login"
			td.FormErrors = errors_user

			err = app.render(w, http.StatusUnprocessableEntity, "signin.tmpl", td)
			if err != nil {
				app.logger.Error("failed to render signin form", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}
		return
	}
	app.session.Put(r, "authenticatedUserID", id)
	app.logger.Info("Session userID", "value", app.session.Get(r, "authenticatedUserID"))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You have logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ---------------------------------------- Venue Handlers ------------------------------------
// Create a record of the respective table
func (app *application) createVenue(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Log the session value for debugging
	app.logger.Info("Session userID", "value", app.session.Get(r, "authenticatedUserID"))

	// Get userID
	userId, ok := app.session.Get(r, "authenticatedUserID").(int)
	if !ok {
		app.session.Put(r, "flash", "Please log in to create a venue.")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	// Extract form values
	// owner := r.FormValue("owner")
	venue_name := r.FormValue("venue_name")
	description := r.FormValue("description")
	location := r.FormValue("location")
	email := r.FormValue("email")
	priceStr := r.FormValue("price_per_hour")
	capacityStr := r.FormValue("max_capacity")
	imageLink := r.FormValue("image")

	// Convert numeric inputs
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		app.logger.Error("invalid price input", "input", priceStr, "error", err)
		price = 0
	}

	capacity, err := strconv.ParseInt(capacityStr, 10, 64)
	if err != nil {
		app.logger.Error("invalid capacity input", "input", capacityStr, "error", err)
		capacity = 0
	}

	// Create the venue struct
	venue := &data.Venue{
		OwnerID:     int64(userId),
		VenueName:   venue_name,
		Description: description,
		Location:    location,
		Email:       email,
		Price:       price,
		MaxCapacity: capacity,
		Image:       imageLink,
	}

	// Validate
	v := validator.NewValidator()
	data.ValidateVenue(v, venue)

	if !v.ValidData() {
		formData := map[string]string{
			"venue_name":     venue_name,
			"description":    description,
			"location":       location,
			"email":          email,
			"price_per_hour": priceStr,
			"max_capacity":   capacityStr,
			"image":          imageLink,
		}

		td := NewTemplateData(r)
		td.Title = "Add Venue"
		td.HeaderText = "Add New Venue Details"
		td.FormErrors = v.Errors
		td.FormData = formData
		td.IsAuthenticated = app.isAuthenticated(r)

		err = app.render(w, http.StatusUnprocessableEntity, "venueform.tmpl", td)
		if err != nil {
			app.logger.Error("failed to render venue form", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	err = app.venue.Insert(venue)
	if err != nil {
		app.logger.Error("failed to insert venue", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.session.Put(r, "flash", "Venue created successfully!")
	app.logger.Info("")
	http.Redirect(w, r, "/venue/listing", http.StatusSeeOther)
}

func (app *application) viewVenue(w http.ResponseWriter, r *http.Request) {
	// Log the session value for debugging
	app.logger.Info("Session userID", "value", app.session.Get(r, "authenticatedUserID"))

	// Set the Content-Security-Policy header to allow external images
	w.Header().Set("Content-Security-Policy", "img-src 'self' https: data:;")

	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) < 2 {
		http.NotFound(w, r)
		return
	}

	idStr := parts[1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid venue ID", http.StatusBadRequest)
		return
	}

	// Convert the id to int64 before passing to GetVenueByID
	venue, err := app.venue.GetVenueByID(id)
	if err != nil {
		app.logger.Error("failed to fetch venue", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if venue == nil {
		http.NotFound(w, r)
		return
	}

	// Fetch reviews for the venue, passing the venue ID as int64
	reviews, err := app.review.GetReviewByVenueID(int64(id))
	if err != nil {
		app.logger.Error("failed to fetch reviews", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert slice of pointers to slice of values
	var reviewList []data.Review
	for _, r := range reviews {
		reviewList = append(reviewList, *r) // Dereference each review pointer
	}

	// Initialize TemplateData
	data := NewTemplateData(r)

	// Set title and header text
	data.Title = venue.VenueName
	data.HeaderText = "Details for " + venue.VenueName
	data.Flash = app.session.PopString(r, "flash")
	data.IsAuthenticated = app.isAuthenticated(r)

	// Add the single venue to the data
	data.Venue = venue

	// Add the reviews (now as a slice of values) to the data
	data.Reviews = reviewList

	// Render the viewvenue template
	err = app.render(w, http.StatusOK, "viewvenue.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render venue view page", "template", "viewvenue.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Form page displayed to add venue
func (app *application) venueForm(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(contextKeyUser).(*data.Users)
	fmt.Println("User role is:", user.Role)

	data := NewTemplateData(r)
	data.Title = "Add Venue"
	data.HeaderText = "Establish Your New Venue!"
	data.IsAuthenticated = app.isAuthenticated(r)

	err := app.render(w, http.StatusOK, "venueform.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render venueform page", "template", "venueform.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Initial page displayed
func (app *application) venueListing(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Security-Policy header to allow external images
	w.Header().Set("Content-Security-Policy", "img-src 'self' https: data:;")

	venues, err := app.venue.FetchAllVenues()
	if err != nil {
		app.logger.Error("failed to get venues", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := NewTemplateData(r)
	data.Title = "Venue"
	data.HeaderText = "Your latest Venue Posts!"
	data.Flash = app.session.PopString(r, "flash")
	data.IsAuthenticated = app.isAuthenticated(r)

	for _, j := range venues {
		// Dereference each pointer
		data.Venues = append(data.Venues, *j)
	}

	err = app.render(w, http.StatusOK, "venue.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render venue page", "template", "venue.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) showUpdateVenueForm(w http.ResponseWriter, r *http.Request) {
	// Get the path and remove "/venue/"
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.NotFound(w, r)
		return
	}

	idStr := parts[2] // this will be "1"
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.logger.Error("invalid venue ID", "input", idStr, "error", err)
		http.Error(w, "Invalid venue ID", http.StatusBadRequest)
		return
	}

	// Fetch the venue by ID
	venue, err := app.venue.GetVenueByID(int(id))
	if err != nil {
		app.logger.Error("failed to fetch venue", "id", id, "error", err)
		http.Error(w, "Venue not found", http.StatusNotFound)
		return
	}

	// Prepare the template data
	tmplData := NewTemplateData(r)
	tmplData.Title = "Edit Venue"
	tmplData.Venue = venue // Pass the venue pointer
	tmplData.IsAuthenticated = app.isAuthenticated(r)

	// Render the template
	err = app.render(w, http.StatusOK, "editvenue.tmpl", tmplData)
	if err != nil {
		app.logger.Error("failed to render update form", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// updateVenue is the handler that processes the update from the form submission
func (app *application) updateVenue(w http.ResponseWriter, r *http.Request) {
	log.Println("received request")

	// Extract ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "invalid URL", http.StatusBadRequest)
		return
	}
	idStr := pathParts[2]

	venueID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println("invalid venue ID:", idStr, err)
		http.Error(w, "invalid venue ID", http.StatusBadRequest)
		return
	}

	// Fetch existing venue
	venue, err := app.venue.GetVenueByID(int(venueID))
	if err != nil {
		log.Println("failed to get venue:", err)
		http.Error(w, "venue not found", http.StatusNotFound)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		log.Println("form parse error:", err)
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

	// Update venue fields
	venue.VenueName = r.FormValue("venue_name")
	venue.Email = r.FormValue("email")
	venue.Description = r.FormValue("description")
	venue.Location = r.FormValue("location")
	venue.Image = r.FormValue("image")

	priceStr := r.FormValue("price")
	maxCapStr := r.FormValue("max_capacity")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err == nil {
		venue.Price = price
	}

	maxCap, err := strconv.ParseInt(maxCapStr, 10, 64)
	if err == nil {
		venue.MaxCapacity = maxCap
	}

	// Validation
	v := validator.NewValidator()
	data.ValidateVenue(v, venue)

	if !v.ValidData() {
		formData := map[string]string{
			"venue_name":   venue.VenueName,
			"email":        venue.Email,
			"description":  venue.Description,
			"location":     venue.Location,
			"price":        priceStr,
			"max_capacity": maxCapStr,
			"image":        venue.Image,
		}

		td := NewTemplateData(r)
		td.Title = "Update Venue"
		td.HeaderText = "Update Venue Details"
		td.FormErrors = v.Errors
		td.FormData = formData
		td.IsAuthenticated = app.isAuthenticated(r)

		err = app.render(w, http.StatusUnprocessableEntity, "editvenue.tmpl", td)
		if err != nil {
			log.Println("failed to render venue form:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Perform the update
	err = app.venue.Update(venue)
	if err != nil {
		log.Println("failed to update venue:", err)
		http.Error(w, "unable to update venue", http.StatusInternalServerError)
		return
	}

	// After successful update, redirect to view venue page
	app.session.Put(r, "flash", "Update Made successfully!")
	app.logger.Info("")
	http.Redirect(w, r, fmt.Sprintf("/venue/%d", venueID), http.StatusSeeOther)
}

func (app *application) deleteVenue(w http.ResponseWriter, r *http.Request) {
	// Extract the venue ID from the URL path
	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		http.NotFound(w, r)
		return
	}

	idStr := parts[1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid venue ID", http.StatusBadRequest)
		return
	}

	// Call the Delete method from the model to remove the venue
	err = app.venue.Delete(int64(id))
	if err != nil {
		app.logger.Error("failed to delete venue", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect to the venues list page or show a success message
	app.session.Put(r, "flash", "Venue Removed successfully!")
	app.logger.Info("")
	http.Redirect(w, r, "/venue/listing", http.StatusSeeOther)
}

// ------------------------------ Reviews
func (app *application) submitReview(w http.ResponseWriter, r *http.Request) {

	// Extract venue ID from URL: /venue/{id}/review
	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 3 {
		http.NotFound(w, r)
		return
	}

	venueID, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "Invalid venue ID", http.StatusBadRequest)
		return
	}

	// Get userID
	userId, ok := app.session.Get(r, "authenticatedUserID").(int)
	if !ok {
		app.session.Put(r, "flash", "Please log in to create a venue.")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	// Parse form data
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Create the review object
	review := data.Review{
		VenueID:    int64(venueID),
		CustomerID: int64(userId),
		Comment:    r.FormValue("comment"),
		CreatedAt:  time.Now(),
	}

	// Insert the review into the database
	err = app.review.Insert(&review)
	if err != nil {
		app.logger.Error("failed to insert review", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect back to venue view
	app.session.Put(r, "flash", "Review Added successfully!")
	app.logger.Info("")
	http.Redirect(w, r, fmt.Sprintf("/venue/%d", venueID), http.StatusSeeOther)
}

// ------------------------------------------- Reservation -------------------------------------------
func (app *application) createReservation(w http.ResponseWriter, r *http.Request) {
	// Log the session value for debugging
	app.logger.Info("Session userID", "value", app.session.Get(r, "authenticatedUserID"))

	// Get userID
	userId, ok := app.session.Get(r, "authenticatedUserID").(int)
	if !ok {
		app.session.Put(r, "flash", "Please log in to create a reservation.")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	// Extract the venue ID from the URL manually
	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) < 3 {
		http.NotFound(w, r)
		return
	}

	idStr := parts[1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid venue ID", http.StatusBadRequest)
		return
	}

	// Parse form input
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Log the form values to check if they're coming through
	fmt.Printf("Form Data: %v\n", r.PostForm)
	startDateStr := r.PostFormValue("start_date")
	startTimeStr := r.PostFormValue("start_time")
	endTimeStr := r.PostFormValue("end_time")
	fmt.Printf("Start Date: %s, Start Time: %s, End Time: %s\n", startDateStr, startTimeStr, endTimeStr)

	// Parse the start date
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		startDate = time.Time{}
	}

	// Combine start date with start time
	startDateTimeStr := fmt.Sprintf("%s %s", startDateStr, startTimeStr)
	startDateTime, err := time.Parse("2006-01-02 15:04", startDateTimeStr)
	if err != nil {
		startDateTime = time.Time{}
	}

	// Combine start date with end time
	endDateTimeStr := fmt.Sprintf("%s %s", startDateStr, endTimeStr)
	endDateTime, err := time.Parse("2006-01-02 15:04", endDateTimeStr)
	if err != nil {
		endDateTime = time.Time{}
	}

	// Create reservation
	reservation := &data.Reservation{
		VenueID:    int64(id),
		CustomerID: int64(userId),
		StartDate:  startDate,
		StartTime:  startDateTime,
		EndTime:    endDateTime,
		Status:     "1",
	}

	// Log the reservation data to see if it's correctly populated
	fmt.Printf("Reservation Data: %+v\n", reservation)

	// Validate reservation
	v := validator.NewValidator() // your validator setup
	data.ValidateReservation(v, reservation)

	if !v.ValidData() {
		// Manually convert url.Values to map[string]string for FormData
		formData := make(map[string]string)
		for key := range r.PostForm {
			formData[key] = r.PostFormValue(key)
		}

		tmplData := NewTemplateData(r)
		tmplData.Venue = &data.Venue{ID: int64(id)} // minimal venue data
		tmplData.FormData = formData
		tmplData.FormErrors = v.Errors
		tmplData.IsAuthenticated = app.isAuthenticated(r)

		err = app.render(w, http.StatusUnprocessableEntity, "viewvenue.tmpl", tmplData)
		if err != nil {
			app.logger.Error("failed to render view venue", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Insert into database
	err = app.reservation.Insert(reservation)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect back to venue view
	app.session.Put(r, "flash", "Reservation Made!")
	app.logger.Info("")
	http.Redirect(w, r, "/reservations", http.StatusSeeOther)
}

func (app *application) showAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := app.reservation.FetchAllConfirmedReservations()
	if err != nil {
		app.logger.Error("failed to get reservations", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := NewTemplateData(r)
	data.Title = "Confirmed Reservations"
	data.IsAuthenticated = app.isAuthenticated(r)

	for _, r := range reservations {
		// Dereference each pointer
		data.Reservation = append(data.Reservation, *r)
	}

	err = app.render(w, http.StatusOK, "reservationlist.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render venue page", "template", "reservationlist.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) showCancelledReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := app.reservation.FetchAllCancelledReservations()
	if err != nil {
		app.logger.Error("failed to get cancelled reservations", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := NewTemplateData(r)
	data.Title = "Cancelled Reservations"
	data.Flash = app.session.PopString(r, "flash")
	data.IsAuthenticated = app.isAuthenticated(r)

	for _, r := range reservations {
		data.Reservation = append(data.Reservation, *r)
	}

	err = app.render(w, http.StatusOK, "reservationlist.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render reservation page", "template", "reservationlist.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) cancelReservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the ID from the URL manually
	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 3 {
		http.NotFound(w, r)
		return
	}

	idStr := parts[2]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.logger.Error("invalid reservation ID", "input", idStr, "error", err)
		http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
		return
	}

	err = app.reservation.Cancel(id)
	if err != nil {
		app.logger.Error("failed to cancel reservation", "error", err)
		http.Error(w, "Failed to cancel reservation", http.StatusInternalServerError)
		return
	}

	app.session.Put(r, "flash", "Cancelled Reservation!")
	app.logger.Info("")
	http.Redirect(w, r, "/reservations/cancelled", http.StatusSeeOther)
}

func (app *application) showUpdateReservationForm(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) < 3 {
		http.NotFound(w, r)
		return
	}

	idStr := parts[2] // Fix: use parts[2] for ID after `/reservations/update/`
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.logger.Error("invalid reservation ID", "input", idStr, "error", err)
		http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
		return
	}

	reservation, err := app.reservation.FetchByID(id)
	if err != nil {
		app.logger.Error("failed to fetch reservation", "id", id, "error", err)
		http.Error(w, "Reservation not found", http.StatusNotFound)
		return
	}

	tmplData := NewTemplateData(r)
	tmplData.Title = "Edit Reservation"
	tmplData.Reservation = []data.Reservation{*reservation}
	tmplData.IsAuthenticated = app.isAuthenticated(r)

	err = app.render(w, http.StatusOK, "updatereservation.tmpl", tmplData)
	if err != nil {
		app.logger.Error("failed to render update form", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) updateReservation(w http.ResponseWriter, r *http.Request) {
	// Log and get user ID from session
	app.logger.Info("Session userID", "value", app.session.Get(r, "authenticatedUserID"))

	userId, ok := app.session.Get(r, "authenticatedUserID").(int)
	if !ok {
		app.session.Put(r, "flash", "Please log in to update a reservation.")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	// Parse form
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Extract reservation ID from URL
	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 3 {
		http.NotFound(w, r)
		return
	}

	idStr := parts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.logger.Error("invalid reservation ID", "input", idStr, "error", err)
		http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
		return
	}

	// Parse venue ID from form
	venueIDStr := r.PostFormValue("venue_id")
	venueID, err := strconv.ParseInt(venueIDStr, 10, 64)
	if err != nil {
		app.logger.Error("invalid venue ID", "input", venueIDStr, "error", err)
		http.Error(w, "Invalid venue ID", http.StatusBadRequest)
		return
	}

	// Parse date/time fields
	startDateStr := r.PostFormValue("start_date")
	startTimeStr := r.PostFormValue("start_time")
	endTimeStr := r.PostFormValue("end_time")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		startDate = time.Time{}
	}

	startDateTimeStr := fmt.Sprintf("%s %s", startDateStr, startTimeStr)
	startDateTime, err := time.Parse("2006-01-02 15:04", startDateTimeStr)
	if err != nil {
		startDateTime = time.Time{}
	}

	endDateTimeStr := fmt.Sprintf("%s %s", startDateStr, endTimeStr)
	endDateTime, err := time.Parse("2006-01-02 15:04", endDateTimeStr)
	if err != nil {
		endDateTime = time.Time{}
	}

	// Build updated reservation
	reservation := &data.Reservation{
		ID:         int64(id),
		VenueID:    venueID,
		CustomerID: int64(userId),
		StartDate:  startDate,
		StartTime:  startDateTime,
		EndTime:    endDateTime,
		Status:     r.PostFormValue("status"),
	}

	// Validate
	v := validator.NewValidator()
	data.ValidateReservation(v, reservation)

	if !v.ValidData() {
		formData := make(map[string]string)
		for key := range r.PostForm {
			formData[key] = r.PostFormValue(key)
		}

		tmplData := NewTemplateData(r)
		tmplData.Venue = &data.Venue{ID: venueID}
		tmplData.FormData = formData
		tmplData.FormErrors = v.Errors
		tmplData.IsAuthenticated = app.isAuthenticated(r)

		err = app.render(w, http.StatusUnprocessableEntity, "updatereservation.tmpl", tmplData)
		if err != nil {
			app.logger.Error("failed to render update form", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Perform update
	err = app.reservation.Update(reservation)
	if err != nil {
		app.logger.Error("failed to update reservation", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect on success
	app.session.Put(r, "flash", "Reservation Updated!")
	http.Redirect(w, r, "/reservations", http.StatusSeeOther)
}
