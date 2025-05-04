// filename: tempates.go
// Description: Loading and caching HTML templates

package main

import (
	"html/template"
	"path/filepath"
)

// newTemplateCache loads all HTML templates from the `./ui/html/` directory
// and stores them in a map.
func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize an empty map to store parsed templates.
	cache := map[string]*template.Template{}

	// Find all template files in the `./ui/html/` directory with a `.tmpl` extension.
	pages, err := filepath.Glob("./ui/html/*.tmpl")
	if err != nil {
		// Return an error if there's an issue retrieving the template files.
		return nil, err
	}

	// Iterate over the found template files and parse them.
	for _, page := range pages {
		fileName := filepath.Base(page)

		// Parse the template file and store it in a variable.
		ts, err := template.ParseFiles(page)
		if err != nil {
			// Return an error if template parsing fails.
			return nil, err
		}
		// Store the parsed template in the cache using the file name as the key.
		cache[fileName] = ts
	}

	// Return the fully populated template cache.
	return cache, nil
}
