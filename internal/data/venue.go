// Filename: internal/data/venue.go
// Description: Venue model that initializes venue varibales and methods
package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aiycoleman/VenueSystemTest2/internal/validator"
)

type Venue struct {
	ID          int64     `json:"id"`
	OwnerID     int64     `json:"owner"`
	VenueName   string    `json:"venue_name"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Email       string    `json:"email"`
	Price       float64   `json:"price_per_hour"`
	MaxCapacity int64     `json:"max_capacity"`
	Image       string    `json:"image_link"`
	CreatedAt   time.Time `json:"created_at"`

	// Reviews []Review
}

// ValidateVenue validates input from the venue form
func ValidateVenue(v *validator.Validator, venue *Venue) {
	v.Check(validator.NotBlank(venue.VenueName), "venue_name", "must be provided")
	v.Check(validator.MaxLength(venue.VenueName, 50), "venue_name", "must not be more than 50 bytes long")

	v.Check(validator.NotBlank(venue.Description), "description", "must be provided")
	v.Check(validator.MaxLength(venue.Description, 500), "description", "must not be more than 500 bytes long")

	v.Check(validator.NotBlank(venue.Location), "location", "must be provided")
	v.Check(validator.MaxLength(venue.Location, 100), "location", "must not be more than 100 bytes long")

	v.Check(validator.NotBlank(venue.Email), "email", "must be provided")
	v.Check(validator.IsValidEmail(venue.Email), "email", "invalid email address")
	v.Check(validator.MaxLength(venue.Email, 100), "email", "must not be more than 100 bytes long")

	v.Check(validator.NotFree(venue.Price), "price_per_hour", "must be greater than 0")

	v.Check(validator.NotZeroInt(venue.MaxCapacity), "max_capacity", "must be greater than 0")

	v.Check(validator.NotBlank(venue.Image), "image_link", "must be provided")
	v.Check(validator.MinLength(venue.Image, 10), "image_link", "must be at least 10 characters")
	v.Check(validator.IsValidURL(venue.Image), "image_link", "must be a valid URL")
}

// VenueModel holds the database connection and methods for handling venues
type VenueModel struct {
	DB *sql.DB
}

// Insert adds a new venue record to the database
func (m *VenueModel) Insert(venue *Venue) error {
	query := `
		INSERT INTO venue (owner, name, description, location, email, price_per_hour, max_capacity, image_link, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use QueryRowContext to assign the returned id and created_at
	return m.DB.QueryRowContext(
		ctx,
		query,
		venue.OwnerID,
		venue.VenueName,
		venue.Description,
		venue.Location,
		venue.Email,
		venue.Price,
		venue.MaxCapacity,
		venue.Image, // Assuming Image is stored as a byte slice (you'll need to convert it)
		venue.CreatedAt,
	).Scan(&venue.ID, &venue.CreatedAt)
}

// GetVenueByID retrieves a venue by its ID from the database.
func (m *VenueModel) GetVenueByID(id int) (*Venue, error) {
	venue := &Venue{}
	query := `
		SELECT id, name, description, location, email, price_per_hour, max_capacity, image_link, created_at
		FROM venue
		WHERE id = $1`

	err := m.DB.QueryRow(query, id).Scan(
		&venue.ID,
		&venue.VenueName,
		&venue.Description,
		&venue.Location,
		&venue.Email,
		&venue.Price,
		&venue.MaxCapacity,
		&venue.Image,
		&venue.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No venue found with the given ID
		}
		return nil, err // Error fetching the venue
	}
	return venue, nil
}

// FetchAllVenues retrieves all venues from the database
func (m *VenueModel) FetchAllVenues() ([]*Venue, error) {
	query := `SELECT id, name, description, location, image_link FROM venue ORDER BY created_at DESC`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var venues []*Venue
	for rows.Next() {
		v := &Venue{}
		err := rows.Scan(&v.ID, &v.VenueName, &v.Description, &v.Location, &v.Image)
		if err != nil {
			return nil, err
		}
		fmt.Println("Venue ID from DB:", v.ID)
		venues = append(venues, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return venues, nil
}

// Update updates an existing venue record in the database
func (m *VenueModel) Update(venue *Venue) error {
	query := `
		UPDATE venue
		SET name = $1, email = $2, description = $3, location = $4, price_per_hour = $5, max_capacity = $6, image_link = $7, created_at = $8
		WHERE id = $9
		RETURNING id`

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query and return the result
	return m.DB.QueryRowContext(
		ctx,
		query,
		venue.VenueName,
		venue.Email,
		venue.Description,
		venue.Location,
		venue.Price,
		venue.MaxCapacity,
		venue.Image,
		venue.CreatedAt,
		venue.ID,
	).Scan(&venue.ID)
}

// Delete deletes a venue record from the database by its ID
func (m *VenueModel) Delete(venueID int64) error {
	query := `
		DELETE FROM venue
		WHERE id = $1`

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the delete query
	_, err := m.DB.ExecContext(ctx, query, venueID)
	return err // Return any error that occurred during the delete operation
}
