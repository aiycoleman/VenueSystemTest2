// filename: render.go
// Description: Rendering HTML templates efficiently using a template cache and a buffer bool.

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
)

// bufferPool is a pool of reusable byte buffers to optimize memory usage when rendering templates.
var bufferPool = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}

// render retrieves a template from the cache, executes it with the provided data, and writes the result to the response.
func (app *application) render(w http.ResponseWriter, status int, page string, data *TemplateData) error {
	// Get a buffer from the pool and reset it.
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf) // Return the buffer to the pool after use.

	// Retrieve the requested template from the cache.
	ts, ok := app.templateCache[page]
	if !ok {
		// Log and return an error if the template does not exist.
		err := fmt.Errorf("template %s does not exist", page)
		app.logger.Error("template does not exist", "template", page, "error", err)
		return err
	}
	// Execute the template, writing the output to the buffer.
	err := ts.Execute(buf, data)
	if err != nil {
		// Log and return an error if the template rendering fails.
		err = fmt.Errorf("failed to render template %s: %w", page, err)
		app.logger.Error("failed to render template", "template", page, "error", err)
		return err
	}

	// Write the HTTP status code before sending the response.
	w.WriteHeader(status)

	// Write the rendered template from the buffer to the response writer.
	_, err = buf.WriteTo(w)
	if err != nil {
		// Log and return an error if writing to the response fails.
		err = fmt.Errorf("failed to write template to response: %w", err)
		app.logger.Error("failed to write template to response", "error", err)
		return err
	}

	return nil
}
