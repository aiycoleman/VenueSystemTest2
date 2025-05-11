package data

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrRecordNotFound     = errors.New("models: no matching recod found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type Users struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Role           int64     `json:"role"`
	Email          string    `json:"email"`
	HashedPassword []byte    `json:"hashedpassword"`
	Active         bool      `json:"active"`
	CreatedAt      time.Time `json:"created_at"`
}

// // ValidateUsers validates the input from the signupform
// func ValidateUsers(v *validator.Validator, users *Users) {
// 	// Customer name
// 	v.Check(validator.NotBlank(users.Name), "name", "must be provided")
// 	v.Check(validator.MaxLength(users.Name, 50), "name", "must not be more than 50 characters long")

// 	// Email
// 	v.Check(validator.NotBlank(users.Email), "email", "must be provided")
// 	v.Check(validator.IsValidEmail(users.Email), "email", "invalid email address")
// 	v.Check(validator.MaxLength(users.Email, 100), "email", "must not be more than 100 characters long")

// 	v.Check(validator.IsValidChoice(users.Role), "role", "must be provided")

// 	v.Check(validator.NotBlank(string(users.HashedPassword)), "hashedpassword", "must be provided")
// }

// UsersModel holds the database connection and methods for handling users
type UsersModel struct {
	DB *sql.DB
}

func (m *UsersModel) Insert(users *Users) error {
	query := `
		INSERT INTO users (name, email, role, password_hash, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Set the current time as created_at
	users.CreatedAt = time.Now()

	err := m.DB.QueryRowContext(
		ctx,
		query,
		users.Name,
		users.Email,
		users.Role,
		users.HashedPassword,
		users.CreatedAt,
	).Scan(&users.ID, &users.CreatedAt)

	if err != nil {
		if strings.Contains(err.Error(), `duplicate key value violates unique constraint "users_email_key"`) {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (m *UsersModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	query := `
		SELECT id, password_hash
		FROM users
		WHERE email = $1
		AND activated = TRUE`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(&id, &hashedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	// Check the password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

func (m *UsersModel) Get(id int) (*Users, error) {
	return nil, nil
}
