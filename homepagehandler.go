package main

import (
	_ "embed"
	"net/http"
)

// embed the home page in the binary
//
//go:embed home.html
var homeHTML []byte

// HomeHandler renders the initial search page from home.html
func (a *app) HomeHandler(w http.ResponseWriter, r *http.Request) {
	// The "/" pattern matches everything, so we need to check
	// that we're at the root here.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Set the content type to HTML
	w.Header().Set("Content-Type", "text/html")
	defer r.Body.Close()

	// Write the embedded HTML file content to the response writer
	_, err := w.Write(homeHTML)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
