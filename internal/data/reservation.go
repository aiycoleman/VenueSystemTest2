package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/aiycoleman/VenueSystemTest2/internal/validator"
)

type Reservation struct {
	ID           int64     `json:"id"`
	VenueID      int64     `json:"venue_id"`
	CustomerID   int64     `json:"customer_id"`
	CustomerName string    `json:"customer"` // <- this line is required
	StartDate    time.Time `json:"start_date"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	VenueName    string    `json:"venue_name"`
}

// ValidateReservation validates the input from the reservation form
func ValidateReservation(v *validator.Validator, reservation *Reservation) {

	v.Check(validator.IsDateSelected(reservation.StartDate), "start_date", "must be provided")
	v.Check(validator.IsDateInFuture(reservation.StartDate), "start_date", "must be a future date")

	// Start time
	v.Check(validator.IsDateSelected(reservation.StartTime), "start_time", "must be provided")
	v.Check(validator.IsTimeInFuture(reservation.StartTime), "start_time", "must be a future time")

	// End time
	v.Check(validator.IsDateSelected(reservation.EndTime), "end_time", "must be provided")
	v.Check(reservation.EndTime.After(reservation.StartTime), "end_time", "must be after the start time")
}

// ReservationModel holds the database connection and methods for handling reservations
type ReservationModel struct {
	DB *sql.DB
}

// Insert adds a new reservation record to the database
func (m *ReservationModel) Insert(reservation *Reservation) error {
	// Set creation time before insert
	reservation.CreatedAt = time.Now()

	query := `
		INSERT INTO reservation (venue, customer, start_date, start_time, end_time, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use QueryRowContext to assign the returned id and created_at
	err := m.DB.QueryRowContext(
		ctx,
		query,
		reservation.VenueID,
		reservation.CustomerID,
		reservation.StartDate,
		reservation.StartTime,
		reservation.EndTime,
		reservation.Status, // make sure you're passing this
		reservation.CreatedAt,
	).Scan(&reservation.ID, &reservation.CreatedAt)

	return err
}

// FetchAllConfirmedReservations retrieves all confirmed reservations from the database
func (m *ReservationModel) FetchAllConfirmedReservations() ([]*Reservation, error) {
	query := `
        SELECT r.id, r.venue, r.start_date, r.start_time, r.end_time, r.status, r.created_at, venue.name
		FROM reservation r
		JOIN venue ON r.venue = venue.id
		WHERE r.status = 1
		ORDER BY r.created_at DESC`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*Reservation
	for rows.Next() {
		r := &Reservation{}
		err := rows.Scan(&r.ID, &r.VenueID, &r.StartDate, &r.StartTime, &r.EndTime, &r.Status, &r.CreatedAt, &r.VenueName)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}

// FetchAllCancelledReservations retrieves all cancelled reservations from the database
func (m *ReservationModel) FetchAllCancelledReservations() ([]*Reservation, error) {
	query := `
		SELECT r.id, r.venue, r.start_date, r.start_time, r.end_time, r.status, r.created_at, venue.name
		FROM reservation r
		JOIN venue ON r.venue = venue.id
		WHERE r.status = 2
		ORDER BY r.created_at DESC`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*Reservation
	for rows.Next() {
		r := &Reservation{}
		err := rows.Scan(&r.ID, &r.VenueID, &r.StartDate, &r.StartTime, &r.EndTime, &r.Status, &r.CreatedAt, &r.VenueName)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}

// Update updates an existing reservation record in the database
func (m *ReservationModel) Update(reservation *Reservation) error {
	query := `
		UPDATE reservation
		SET start_date = $1, start_time = $2, end_time = $3, status = $4
		WHERE id = $5
		RETURNING id`

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query and return the result
	return m.DB.QueryRowContext(
		ctx,
		query,
		reservation.StartDate,
		reservation.StartTime,
		reservation.EndTime,
		reservation.Status,
		reservation.ID,
	).Scan(&reservation.ID)
}

// Cancel updates the status of an existing reservation to 'cancelled'
func (m *ReservationModel) Cancel(reservationID int64) error {
	query := `
		UPDATE reservation
		SET status = 2
		WHERE id = $1`

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the update query
	_, err := m.DB.ExecContext(ctx, query, reservationID)
	return err // Return any error that occurred during the update operation
}

// Query for 1 Reservation data by
func (m *ReservationModel) FetchByID(id int) (*Reservation, error) {
	query := `
	SELECT 
		r.id, r.venue, r.customer, c.name,
		r.start_date, r.start_time, r.end_time, r.status, r.created_at,
		v.name AS venue_name
	FROM reservation r
	JOIN venue v ON r.venue = v.id
	JOIN users c ON r.customer = c.id
	WHERE r.id = $1;`

	row := m.DB.QueryRow(query, id)

	var res Reservation
	err := row.Scan(
		&res.ID,           // r.id
		&res.VenueID,      // r.venue
		&res.CustomerID,   // r.customer
		&res.CustomerName, // c.name
		&res.StartDate,    // r.start_date
		&res.StartTime,    // r.start_time
		&res.EndTime,      // r.end_time
		&res.Status,       // r.status
		&res.CreatedAt,    // r.created_at
		&res.VenueName,    // v.name
	)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
