package main

import (
	"log"
	"net/http"
)

func startServer(a *app) {
	mux := http.NewServeMux()

	// Route for the search form page
	mux.HandleFunc("/", a.HomeHandler)

	// Route for submitting the search
	mux.HandleFunc("/search", a.SearchHandler)

	// Route for the booking handler
	mux.HandleFunc("/booking", a.BookingHandler)

	// Start the server
	go func() {
		log.Println("Listening on http://localhost:8020")
		log.Fatal(http.ListenAndServe("localhost:8020", mux))
	}()
}
