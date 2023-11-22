package main

import (
	"html/template"
	"net/http"
)

// BookingHandler receives a query URL containing offer ID, queries the Amadeus Transfer Booking API, and renders a new page with a booking confirmation
func (a *app) BookingHandler(w http.ResponseWriter, r *http.Request) {
	// Get the offer ID from the query string
	offerID := r.URL.Query().Get("offerId")

	// Call the Amadeus Transfer Booking API
	// (see internal/amadeus/book.go)
	response, err := a.amadeusClient.Book(offerID)
	if err != nil {
		// Render the erorr nicely
		template.Must(template.New("bookingError").Parse(bookingErrorTemplate)).Execute(w, err)
		return
	}

	// Render the booking receipt template
	tmpl, err := template.New("bookingReceipt").Parse(bookingConfirmationTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// booking confirmation template
// include detail from the BookingResponse
var bookingConfirmationTemplate = `<html>
<body>
	<h1>Booking Confirmation</h1>
	<p>Reference: {{.Data.Reference}}</p>
	<p>Booking ID: {{.Data.ID}}</p>
	<p>Thank you for travelling with us!</p>
	<p><a href="/">New search</a></p>
</body>
</html>`

var bookingErrorTemplate = `<html>
<body>
	<h1>Booking Error</h1>
	<p>We're sorry, but there was an error with your booking.</p>
	<p>{{.}}</p>
	<p><a href="/">New search</a></p>
</body>
</html>`
