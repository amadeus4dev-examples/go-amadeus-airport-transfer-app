package main

import (
	"html/template"
	"net/http"

	"airport-transfer-app/internal/amadeus"
)

// SearchHandler receives a query URL containing start address and airport code, queries the Amadeus Transfer Search API, and renders a new page with a list of offers, or a message if there are no offers available.
func (a *app) SearchHandler(w http.ResponseWriter, r *http.Request) {

	// Parse the query parameters from the request URL
	queryParams := r.URL.Query()

	// Retrieve the query parameters and save them to local variables
	searchParams := amadeus.SearchParameters{
		StartAddressLine: queryParams.Get("streetAddress") + " " + queryParams.Get("houseNumber"),
		StartCityName:    queryParams.Get("city"),
		StartZipCode:     queryParams.Get("zipCode"),
		StartCountryCode: queryParams.Get("countryCode"),
		StartGeoCode:     queryParams.Get("latitude") + "," + queryParams.Get("longitude"),
		EndLocationCode:  queryParams.Get("endLocationCode"),
		StartDateTime:    string(queryParams.Get("startDateTime")),
	}

	// Check if any parameter (except houseNumber) is empty
	if searchParams.EndLocationCode == "" ||
		searchParams.StartAddressLine == " " ||
		searchParams.StartCityName == "" ||
		searchParams.StartZipCode == "" ||
		searchParams.StartCountryCode == "" ||
		searchParams.StartGeoCode == "" {
		template.Must(template.New("incompleteAddress").Parse(incompleteAddressTemplate)).Execute(w, searchParams)
		return
	}

	// Call the Amadeus Transfer Search API
	// (see internal/amadeus/search.go)

	response, err := a.amadeusClient.Search(searchParams)
	if err != nil {
		template.Must(template.New("searchError").Parse(searchErrorTemplate)).Execute(w, struct {
			Search amadeus.SearchParameters
			Error  string
		}{searchParams, err.Error()})
		return
	}

	// Parse the offer list template
	tmpl, err := template.New("offerList").Parse(offerListTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template to the ResponseWriter
	err = tmpl.Execute(w, response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

const offerListTemplate = `
    <!DOCTYPE html>
    <html>
        <head>
            <title>Available Transfers</title>
        </head>
        <body>
			{{if ne (len .Data) 0}}
			<table>
				{{range .Data}}
				<tr>
					<td>Transfer Type</td>
					<td>{{.TransferType}}</td>
				</tr>
				<tr>
					<td>Start Time</td>
					<td>{{.Start.DateTime}}</td>
				</tr>
				<tr>
					<td>Arrival Time</td>
					<td>{{.End.DateTime}}</td>
				</tr>
				<tr>
					<td>Service Provider</td>
					<td>{{.ServiceProvider.Name}}</td>
				</tr>
				<tr>
					<td>Estimated Cost</td>
					<td>{{.Quotation.CurrencyCode}} {{.Quotation.MonetaryAmount}}</td>
				</tr>
					<td><button class="book" onclick="bookOffer('{{.ID}}')">Book this transfer</button></td>
					<td></td>
				{{end}}
			</table>
			{{else}}
				<p>Sorry, there are no transfers available.</p>
			{{end}}
			<p><a href="/">New search</a></p>
			<script>
				function bookOffer(offerId) {
					document.querySelectorAll(".book").forEach(function(bookButton) {
						bookButton.disabled = true
					})
					var queryString = "/booking?offerId=" + offerId;
					window.location.href = queryString;
				}
			</script>

        </body>
    </html>

`

const incompleteAddressTemplate = `<html>
<body>
  <h1>Address data is incomplete</h1>
  <p>Street address: {{.StartAddressLine}}</p>
  <p>City: {{.StartCityName}}</p>
  <p>Zip code: {{.StartZipCode}}</p>
  <p>Country code: {{.StartCountryCode}}</p>
  <p><a href="/">New search</a></p>
</body>
</html>`

const searchErrorTemplate = `<html>
<body>
  <h1>Search failed</h1>
  <p>We're sorry, but there was an error with your search.</p>
  <p><strong>{{.Error}}</strong></p>
  <p>Start address: {{.Search.StartAddressLine}}<br/>
  City: {{.Search.StartCityName}}<br/>
  Zip code: {{.Search.StartZipCode}}<br/>
  Country code: {{.Search.StartCountryCode}}</p>
  <p><a href="/">New search</a></p>
</body>
</html>`
