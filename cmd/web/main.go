// filename: main.go
// Description: Entry point of the Web Journal Application
package main

import (
	"context"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"os"
	"time"

	// "github.com/aiycoleman/VenueSystemTest2/internal/data"
	_ "github.com/lib/pq"
)

// application struct holds the application's dependencies.
type application struct {
	addr *string
	// venue         *data.VenueModel
	// reservation   *data.ReservationModel
	// review        *data.ReviewModel
	logger        *slog.Logger
	templateCache map[string]*template.Template
}

// Define command-line flags for server address and database connection
func main() {
	addr := flag.String("addr", "", "HTTP network address")
	dsn := flag.String("dsn", "", "PostgreSQL DSN")

	// Parse the command-line flags
	flag.Parse()

	// Initialize a logger for structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Open a connection to the PostgreSQL database
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Log successful DB connection
	logger.Info("database connection pool established")

	// Create a cache for templates to optimize rendering performance
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Ensure the database connection is closed when the application exits
	defer db.Close()

	// Initialize the application struct with dependencies
	app := &application{
		addr: addr,
		// venue:         &data.VenueModel{DB: db},
		// review:        &data.ReviewModel{DB: db},
		// reservation:   &data.ReservationModel{DB: db},
		logger:        logger,
		templateCache: templateCache,
	}

	// Start the HTTP server
	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

// openDB establishes a connection to the PostgreSQL database.
func openDB(dsn string) (*sql.DB, error) {
	// Open the database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Set a timeout context for testing the database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping the database to ensure the connection is valid
	err = db.PingContext(ctx)
	if err != nil {
		db.Close() // Close the DB connection if the ping fails
		return nil, err
	}

	return db, nil // Return the database connection if successful
}
