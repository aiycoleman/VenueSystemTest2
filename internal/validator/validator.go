// Filename: internal/validator/validator.go

package validator

import (
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Validator type is just a map
type Validator struct {
	Errors map[string]string
}

// Create a new Validator
func NewValidator() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// Check if there are any entries in the map
func (v *Validator) ValidData() bool {
	return len(v.Errors) == 0
}

// Add an error entry to the error map.
func (v *Validator) AddError(field string, message string) {
	_, exists := v.Errors[field]
	if !exists {
		v.Errors[field] = message
	}
}

// Check adds an error if the validation check fails
func (v *Validator) Check(ok bool, field string, message string) {
	if !ok {
		v.AddError(field, message)
	}
}

// returns true if data is present in the input box
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MinLength returns true if the value contains at least n characters
func MinLength(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// MaxLength returns true if the value contains no more than n characters
func MaxLength(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// IsEmail returns true if the value is a valid email address
func IsValidEmail(email string) bool {
	return EmailRX.MatchString(email)
}

// NotZeroInt checks if the integer value is not zero
func NotZeroInt(value int64) bool {
	return value > 0
}

// NotFree checks if the integer value is not zero
func NotFree(value float64) bool {
	return value > 0.0
}

// Checking if user inputed valid image link
func IsValidURL(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	// Ensure scheme is http or https and host is present
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}
	if parsedURL.Host == "" {
		return false
	}

	return true
}

// Checks if the date is in the future
func IsDateInFuture(t time.Time) bool {
	return t.After(time.Now())
}

// IsDateSelected returns true if the time is not zero (i.e., a date/time is selected)
func IsDateSelected(t time.Time) bool {
	return !t.IsZero()
}

// IsTimeInFuture checks if a time is in the future
func IsTimeInFuture(t time.Time) bool {
	return t.After(time.Now())
}

// IsValidRole checks if role is 1 (owner) or 2 (customer)
func IsValidChoice(choice string) bool {
	return choice == "1" || choice == "2"
}
