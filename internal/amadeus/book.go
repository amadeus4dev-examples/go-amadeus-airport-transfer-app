package amadeus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Book receives an offer ID that the user selects from the transfer
// offers page, and returns a BookingResponse struct containing a booking
// confirmation, or an error.
func (c *Client) Book(offerId string) (BookingResponse, error) {

	url := c.baseURL + "/ordering/transfer-orders?offerId=" + offerId
	method := "POST"

	// This data is typically collected beforehand, e.g. when the
	// client registers with the service and when they book a flight.
	// For brevity, we use mock data instead.
	payload := strings.NewReader(`{
  "data": {
    "note": "Note to driver",
    "passengers": [
      {
        "firstName": "John",
        "lastName": "Doe",
        "title": "MR",
        "contacts": {
          "phoneNumber": "+33123456789",
          "email": "user@email.com"
        },
        "billingAddress": {
          "line": "Avenue de la Bourdonnais, 19",
          "zip": "75007",
          "countryCode": "FR",
          "cityName": "Paris"
        }
      }
    ],
    "payment": {
      "methodOfPayment": "CREDIT_CARD",
      "creditCard": {
        "number": "4111111111111111",
        "holderName": "JOHN DOE",
        "vendorCode": "VI",
        "expiryDate": "0928",
        "cvv": "111"
      }
    },
    "startConnectedSegment": {
      "transportationType": "FLIGHT",
      "transportationNumber": "AF380",
      "departure": {
        "uicCode": "7400001",
        "iataCode": "CDG",
        "localDateTime": "` + time.Now().Format("2006-01-02T15:04:05") + `"
      },
      "arrival": {
        "uicCode": "7400001",
        "iataCode": "CDG",
        "localDateTime": "` + time.Now().Format("2006-01-02T15:04:05") + `"
      }
    },
    "endConnectedSegment": {
      "transportationType": "FLIGHT",
      "transportationNumber": "AF380",
      "departure": {
        "uicCode": "7400001",
        "iataCode": "CDG",
        "localDateTime": "` + time.Now().Format("2006-01-02T15:04:05") + `"
      },
      "arrival": {
        "uicCode": "7400001",
        "iataCode": "CDG",
        "localDateTime": "` + time.Now().Format("2006-01-02T15:04:05") + `"
      }
    }
  }`)

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return BookingResponse{}, fmt.Errorf("book: http.NewRequest: %w", err)
	}

	token, err := c.token()
	if err != nil {
		return BookingResponse{}, fmt.Errorf("book: c.token: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		return BookingResponse{}, fmt.Errorf("book: client.Do: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return BookingResponse{}, fmt.Errorf("book: io.ReadAll: %w", err)
	}

	// Check for API errors.
	// HTTP status is 200 even if the booking fails,
	// because technically, the call succeeded.
	// Hence, we check for the occurrence of "errors" in the response body.
	if bytes.Contains(body, []byte("errors")) {
		errorResult := BookingErrorResponse{}
		err = json.Unmarshal(body, &errorResult)
		if err != nil {
			return BookingResponse{}, fmt.Errorf("book: json.Unmarshal: %w", err)
		}
		return BookingResponse{}, fmt.Errorf("Booking failed: %s (code %d)", errorResult.Errors[0].Detail, errorResult.Errors[0].Code)
	}

	result := BookingResponse{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return BookingResponse{}, fmt.Errorf("book: json.Unmarshal: %w", err)
	}
	return result, nil
}
