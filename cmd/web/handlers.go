// filename: handlers.go
// Description: Handling HTTP requests

package main

import (
	"errors"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/aiycoleman/VenueSystemTest2/internal/data"
	"github.com/aiycoleman/VenueSystemTest2/internal/validator"
)

// ------------------------------- Home Handler --------------------------------
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData(r)
	data.Title = "Welcome to Venue Verge!"
	data.HeaderText = "Find Your Perfect Venue!"

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
			td.Title = "Nice to See You Again!"
			td.HeaderText = "Login"

			err = app.render(w, http.StatusUnprocessableEntity, "sigin.tmpl", td)
			if err != nil {
				app.logger.Error("failed to render signup form", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}
		return
	}
	app.session.Put(r, "authenicatedUserID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenicatedUserID")
	app.session.Put(r, "flash", "You have logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
