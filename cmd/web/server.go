// filename: server.go
// Description: Starting and configuring the HTTP server

package main

import (
	"log/slog"
	"net/http"
	"time"
)

func (app *application) serve() error {
	// srv configures and starts the HTTP server with defined settings,
	// including request handling, timeouts, and error logging.
	srv := &http.Server{
		Addr:         *app.addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError), // Logs server errors
	}
	app.logger.Info("starting server", "addr", srv.Addr)
	return srv.ListenAndServe() // begins handling http requests
}
