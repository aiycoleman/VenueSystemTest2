// Filename: internal/data/review.go
// Description: Review model that initializes venue varibales and methods
package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/aiycoleman/VenueSystemTest2/internal/validator"
)

type Review struct {
	ID           int64     `json:"id"`
	CustomerID   int64     `json:"customer_id"`
	CustomerName string    `json:"customer_name"`
	VenueID      int64     `json:"venue_id"`
	Comment      string    `json:"comment"`
	CreatedAt    time.Time `json:"created_at"`
}

// ValidateReview validates input from the review form
func ValidateReview(v *validator.Validator, review *Review) {
	v.Check(validator.NotBlank(review.Comment), "comment", "must be provided")
	v.Check(validator.MaxLength(review.Comment, 500), "comment", "must not be more than 500 bytes long")
}

// ReviewModel holds the database connection and methods for handling venues
type ReviewModel struct {
	DB *sql.DB
}

// Insert adds a new review record to the database
func (m *ReviewModel) Insert(review *Review) error {
	query := `
		INSERT INTO review (customer, venue, comment, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use QueryRowContext to assign the returned id and created_at
	return m.DB.QueryRowContext(
		ctx,
		query,
		review.CustomerID,
		review.VenueID,
		review.Comment,
		review.CreatedAt,
	).Scan(&review.ID, &review.CreatedAt)
}

// GetReviewByVenueID fetches reviews by venue ID
func (m *ReviewModel) GetReviewByVenueID(venueID int64) ([]*Review, error) {
	query := `
		SELECT r.id, r.customer, u.name, r.venue, r.comment, r.created_at
		FROM review r
		JOIN users u ON r.customer = u.id
		WHERE r.venue = $1
		ORDER BY r.created_at DESC;
	`

	rows, err := m.DB.Query(query, venueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*Review
	for rows.Next() {
		var r Review
		err := rows.Scan(
			&r.ID,
			&r.CustomerID,
			&r.CustomerName,
			&r.VenueID,
			&r.Comment,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &r)
	}

	return reviews, nil
}
