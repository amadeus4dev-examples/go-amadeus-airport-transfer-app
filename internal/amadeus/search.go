package amadeus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Search receives search parameters from the user and calls the
// Search Tranfer API to receive a list of transfer offers.
// It returns a SearchResponse struct containing the list of offers,
// or an error if the search fails.
func (c *Client) Search(p SearchParameters) (SearchResponse, error) {
	url := c.baseURL + "/shopping/transfer-offers"
	method := "POST"

	params, err := json.Marshal(p)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("Search: json.Marshal: %w", err)
	}

	payload := bytes.NewReader(params)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return SearchResponse{}, fmt.Errorf("Search: http.NewRequest: %w", err)
	}

	token, err := c.token()
	if err != nil {
		return SearchResponse{}, fmt.Errorf("Search: c.token: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("Search: client.Do: %w", err)
	}
	defer res.Body.Close()

	// Read the response body and unmarshal it into a SearchResponse variable

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("Search: io.ReadAll: %w", err)
	}

	// Check for API errors.
	// HTTP status is 200 even if the booking fails,
	// because technically, the call succeeded.
	// Hence, we check for the occurrence of "errors" in the response body.
	if bytes.Contains(body, []byte("errors")) {
		errorResult := SearchErrorResponse{}
		err = json.Unmarshal(body, &errorResult)
		if err != nil {
			return SearchResponse{}, fmt.Errorf("Search: json.Unmarshal: %w", err)
		}
		return SearchResponse{}, fmt.Errorf("Search failed: %s: %s (code %d)", errorResult.Errors[0].Title, errorResult.Errors[0].Detail, errorResult.Errors[0].Code)
	}

	// Unmarshal the response into a SearchResponse struct
	result := SearchResponse{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("Search: json.Unmarshal: %w", err)
	}
	return result, nil
}
