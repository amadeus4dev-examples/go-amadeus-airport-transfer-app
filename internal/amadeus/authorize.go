package amadeus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// token returns the current access token. If none exists yet, or if the existing one has expired, it fetches a new one from the Amadeus authorization API. If fetching fails, token returns an error.
func (c *Client) token() (string, error) {
	t := <-c.accessToken
	return t.Token, t.Err
}

// startTokenFetcher starts a goroutine that fetches a new access token from the Amadeus authorization API if there is none yet, or if the current one expires. It returns channels for returning the current token, or an error if the token could not be fetched.
func (c *Client) refreshToken() {
	var token string
	var expiration time.Duration
	var err error

	// Set the initial token, before any client can request it.
	token, expiration, err = authorize(c.baseURL)

	// Set a new timer to fire when 90% of the expiration duration has passed.
	// We want a new token *before* the current one expires.
	expired := time.After(expiration * 90 / 100)

	for {
		select {
		// The expiration timer has fired and wrote the current time to `expired`.
		case <-expired:
			token, expiration, err = authorize(c.baseURL)
			// Set a new timer to fire when 90% of the expiration duration has passed.
			expired = time.After(expiration * 90 / 100)

		case c.accessToken <- tokenResponse{Token: token, Err: err}:
			// Someone has read the token, nothing to do.
			// The next iteration will send the token to the channel again.
		}
	}
}

// authorize reads client ID and secret from the environment variables and updates the access token and its lifespan (in seconds) from the Amadeus authorization API.
func authorize(baseURL string) (token string, lifespan time.Duration, err error) {

	url := baseURL + "/security/oauth2/token"
	method := "POST"

	id := os.Getenv("AMADEUS_CLIENT_ID")
	secret := os.Getenv("AMADEUS_CLIENT_SECRET")

	if id == "" || secret == "" {
		return "", 0, fmt.Errorf("authorize: missing client ID or secret (check the environment variables AMADEUS_CLIENT_ID and AMADEUS_CLIENT_SECRET)")
	}

	payload := strings.NewReader("client_id=" + id +
		"&client_secret=" + secret +
		"&grant_type=client_credentials")

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return "", 0, fmt.Errorf("authorize: http.NewRequest: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("authorize: client.Do: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", 0, fmt.Errorf("authorize: io.ReadAll: %w", err)
	}

	// Unmarshal the response. AuthResponse is a struct that combines
	// the responses for the successful case and for the error case.
	// Unmarshal() does not complain if the JSON does not fill the
	// struct completely, which we use here to simplify the unmarshalling.
	var authResponse AuthResponse
	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		return "", 0, fmt.Errorf("authorize: json.Unmarshal: %w", err)
	}
	if authResponse.Error != "" {
		return "", 0, fmt.Errorf("authorize: %w (%s: %s (error: %s, code: %d)",
			err,
			authResponse.Title,
			authResponse.ErrorDescription,
			authResponse.Error,
			authResponse.Code)
	}

	return authResponse.AccessToken,
		time.Duration(authResponse.ExpiresIn) * time.Second,
		nil

}
