package amadeus

// tokenResponse contains either a valid access token
// or an error that occurred while fetching the token
type tokenResponse struct {
	Token string
	Err   error
}

// Client is a client for the Amadeus API.
// It takes care of refreshing the access token regularly
// in the background while serving the currently valid
// token to clients
type Client struct {
	baseURL     string
	accessToken chan tokenResponse
}

// Create a new client and start the token refreshing goroutine.
func New() *Client {
	c := &Client{
		baseURL:     "https://test.api.amadeus.com/v1",
		accessToken: make(chan tokenResponse),
	}
	go c.refreshToken()
	return c
}

// AuthResponse contains the unmarshaled response from the Amadeus
// authorization API.
// It is a blend of the success and error response, so that we can
// unmarshal the response first and check for success or error later.
type AuthResponse struct {
	AuthSuccessResponse
	AuthErrorResponse
}

type AuthSuccessResponse struct {
	Type            string `json:"type"`
	Username        string `json:"username"`
	ApplicationName string `json:"application_name"`
	ClientID        string `json:"client_id"`
	TokenType       string `json:"token_type"`
	AccessToken     string `json:"access_token"`
	ExpiresIn       int    `json:"expires_in"`
	State           string `json:"state"`
	Scope           string `json:"scope"`
}

type AuthErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	Code             int    `json:"code"`
	Title            string `json:"title"`
}

type SearchParameters struct {
	StartLocationCode string `json:"startLocationCode,omitempty"`
	EndAddressLine    string `json:"endAddressLine,omitempty"`
	EndCityName       string `json:"endCityName,omitempty"`
	EndZipCode        string `json:"endZipCode,omitempty"`
	EndCountryCode    string `json:"endCountryCode,omitempty"`
	EndName           string `json:"endName,omitempty"`
	EndGeoCode        string `json:"endGeoCode,omitempty"`
	EndLocationCode   string `json:"endLocationCode,omitempty"`
	StartAddressLine  string `json:"startAddressLine,omitempty"`
	StartCityName     string `json:"startCityName,omitempty"`
	StartZipCode      string `json:"startZipCode,omitempty"`
	StartCountryCode  string `json:"startCountryCode,omitempty"`
	StartName         string `json:"startName,omitempty"`
	StartGeoCode      string `json:"startGeoCode,omitempty"`
	TransferType      string `json:"transferType,omitempty"`
	StartDateTime     string `json:"startDateTime,omitempty"`
	ProviderCodes     string `json:"providerCodes,omitempty"`
	Passengers        int    `json:"passengers,omitempty"`
	StopOvers         []struct {
		Duration       string `json:"duration,omitempty"`
		SequenceNumber int    `json:"sequenceNumber,omitempty"`
		AddressLine    string `json:"addressLine,omitempty"`
		CountryCode    string `json:"countryCode,omitempty"`
		CityName       string `json:"cityName,omitempty"`
		ZipCode        string `json:"zipCode,omitempty"`
		Name           string `json:"name,omitempty"`
		GeoCode        string `json:"geoCode,omitempty"`
		StateCode      string `json:"stateCode,omitempty"`
	} `json:"stopOvers,omitempty"`
	StartConnectedSegment struct {
		TransportationType   string `json:"transportationType,omitempty"`
		TransportationNumber string `json:"transportationNumber,omitempty"`
		Departure            struct {
			LocalDateTime string `json:"localDateTime,omitempty"`
			IataCode      string `json:"iataCode,omitempty"`
		} `json:"departure,omitempty"`
		Arrival struct {
			LocalDateTime string `json:"localDateTime,omitempty"`
			IataCode      string `json:"iataCode,omitempty"`
		} `json:"arrival,omitempty"`
	} `json:"startConnectedSegment,omitempty"`
	PassengerCharacteristics []struct {
		PassengerTypeCode string `json:"passengerTypeCode,omitempty"`
		Age               int    `json:"age,omitempty"`
	} `json:"passengerCharacteristics,omitempty"`
}

type SearchResponse struct {
	Data []struct {
		ID           string `json:"id"`
		Type         string `json:"type"`
		TransferType string `json:"transferType"`
		Start        struct {
			DateTime     string `json:"dateTime"`
			LocationCode string `json:"locationCode"`
		} `json:"start"`
		End struct {
			DateTime string `json:"dateTime"`
			Address  struct {
				Line        string  `json:"line"`
				Zip         string  `json:"zip"`
				CountryCode string  `json:"countryCode"`
				CityName    string  `json:"cityName"`
				Latitude    float64 `json:"latitude"`
				Longitude   float64 `json:"longitude"`
			} `json:"address"`
			Name string `json:"name"`
		} `json:"end"`
		Vehicle struct {
			Code        string `json:"code"`
			Category    string `json:"category"`
			Description string `json:"description"`
			ImageURL    string `json:"imageURL"`
			Baggages    []struct {
				Count int    `json:"count"`
				Size  string `json:"size"`
			} `json:"baggages"`
			Seats []struct {
				Count int `json:"count"`
			} `json:"seats"`
		} `json:"vehicle"`
		ServiceProvider struct {
			Code     string   `json:"code"`
			Name     string   `json:"name"`
			TermsURL string   `json:"termsUrl"`
			LogoURL  string   `json:"logoUrl"`
			Settings []string `json:"settings"`
		} `json:"serviceProvider"`
		Quotation struct {
			MonetaryAmount string `json:"monetaryAmount"`
			CurrencyCode   string `json:"currencyCode"`
			Taxes          []struct {
				MonetaryAmount string `json:"monetaryAmount"`
			} `json:"taxes"`
			TotalTaxes struct {
				MonetaryAmount string `json:"monetaryAmount"`
			} `json:"totalTaxes"`
			Base struct {
				MonetaryAmount string `json:"monetaryAmount"`
			} `json:"base"`
			Discount struct {
				MonetaryAmount string `json:"monetaryAmount"`
			} `json:"discount"`
			TotalFees struct {
				MonetaryAmount string `json:"monetaryAmount"`
			} `json:"totalFees"`
		} `json:"quotation"`
		CancellationRules []struct {
			FeeType         string `json:"feeType"`
			FeeValue        string `json:"feeValue"`
			CurrencyCode    string `json:"currencyCode"`
			MetricType      string `json:"metricType"`
			MetricMin       string `json:"metricMin"`
			MetricMax       string `json:"metricMax"`
			RuleDescription string `json:"ruleDescription"`
		} `json:"cancellationRules"`
		MethodsOfPaymentAccepted []string `json:"methodsOfPaymentAccepted"`
		PassengerCharacteristics []struct {
			PassengerTypeCode string `json:"passengerTypeCode"`
			Age               int    `json:"age"`
		} `json:"passengerCharacteristics"`
		Converted struct {
			MonetaryAmount string `json:"monetaryAmount"`
			CurrencyCode   string `json:"currencyCode"`
			Taxes          []struct {
				MonetaryAmount string `json:"monetaryAmount"`
			} `json:"taxes"`
			TotalTaxes struct {
				MonetaryAmount string `json:"monetaryAmount"`
			} `json:"totalTaxes"`
			Base struct {
				MonetaryAmount string `json:"monetaryAmount"`
			} `json:"base"`
			Discount struct {
				MonetaryAmount string `json:"monetaryAmount"`
			} `json:"discount"`
			TotalFees struct {
				MonetaryAmount string `json:"monetaryAmount"`
			} `json:"totalFees"`
		} `json:"converted"`
	} `json:"data"`
}

type SearchErrorResponse struct {
	Errors []struct {
		Status int    `json:"status"`
		Code   int    `json:"code"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
		Source struct {
			Parameter string `json:"parameter"`
		} `json:"source"`
	} `json:"errors"`
}

type BookingParameters struct {
	Data struct {
		Note       string `json:"note"`
		Passengers []struct {
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Title     string `json:"title"`
			Contacts  struct {
				PhoneNumber string `json:"phoneNumber"`
				Email       string `json:"email"`
			} `json:"contacts"`
			BillingAddress struct {
				Line        string `json:"line"`
				Zip         string `json:"zip"`
				CountryCode string `json:"countryCode"`
				CityName    string `json:"cityName"`
			} `json:"billingAddress"`
		} `json:"passengers"`
		Agency struct {
			Contacts []struct {
				Email struct {
					Address string `json:"address"`
				} `json:"email"`
			} `json:"contacts"`
		} `json:"agency"`
		Payment struct {
			MethodOfPayment string `json:"methodOfPayment"`
			CreditCard      struct {
				Number     string `json:"number"`
				HolderName string `json:"holderName"`
				VendorCode string `json:"vendorCode"`
				ExpiryDate string `json:"expiryDate"`
				Cvv        string `json:"cvv"`
			} `json:"creditCard"`
		} `json:"payment"`
		ExtraServices []struct {
			Code   string `json:"code"`
			ItemID string `json:"itemId"`
		} `json:"extraServices"`
		Equipment []struct {
			Code string `json:"code"`
		} `json:"equipment"`
		Corporation struct {
			Address struct {
				Line        string `json:"line"`
				Zip         string `json:"zip"`
				CountryCode string `json:"countryCode"`
				CityName    string `json:"cityName"`
			} `json:"address"`
			Info struct {
				Au string `json:"AU"`
				Ce string `json:"CE"`
			} `json:"info"`
		} `json:"corporation"`
		StartConnectedSegment struct {
			TransportationType   string `json:"transportationType"`
			TransportationNumber string `json:"transportationNumber"`
			Departure            struct {
				UicCode       string `json:"uicCode"`
				IataCode      string `json:"iataCode"`
				LocalDateTime string `json:"localDateTime"`
			} `json:"departure"`
			Arrival struct {
				UicCode       string `json:"uicCode"`
				IataCode      string `json:"iataCode"`
				LocalDateTime string `json:"localDateTime"`
			} `json:"arrival"`
		} `json:"startConnectedSegment"`
		EndConnectedSegment struct {
			TransportationType   string `json:"transportationType"`
			TransportationNumber string `json:"transportationNumber"`
			Departure            struct {
				UicCode       string `json:"uicCode"`
				IataCode      string `json:"iataCode"`
				LocalDateTime string `json:"localDateTime"`
			} `json:"departure"`
			Arrival struct {
				UicCode       string `json:"uicCode"`
				IataCode      string `json:"iataCode"`
				LocalDateTime string `json:"localDateTime"`
			} `json:"arrival"`
		} `json:"endConnectedSegment"`
	} `json:"data"`
}

type BookingResponse struct {
	Data struct {
		Type       string `json:"type"`
		Reference  string `json:"reference"`
		ID         string `json:"id"`
		Passengers []struct {
			Type      string `json:"type"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Title     string `json:"title"`
			Contacts  struct {
				Email       string `json:"email"`
				PhoneNumber string `json:"phoneNumber"`
			} `json:"contacts"`
			BillingAddress struct {
				Line        string `json:"line"`
				Zip         string `json:"zip"`
				CountryCode string `json:"countryCode"`
				CityName    string `json:"cityName"`
			} `json:"billingAddress"`
		} `json:"passengers"`
		Transfers []struct {
			Status            string `json:"status"`
			ConfirmNbr        string `json:"confirmNbr"`
			Note              string `json:"note"`
			MethodOfPayment   string `json:"methodOfPayment"`
			OfferID           string `json:"offerId"`
			TransferType      string `json:"transferType"`
			CancellationRules []struct {
				FeeType         string `json:"feeType"`
				FeeValue        string `json:"feeValue"`
				CurrencyCode    string `json:"currencyCode"`
				MetricMax       string `json:"metricMax"`
				MetricType      string `json:"metricType"`
				MetricMin       string `json:"metricMin"`
				RuleDescription string `json:"ruleDescription"`
			} `json:"cancellationRules"`
			Start struct {
				DateTime     string `json:"dateTime"`
				LocationCode string `json:"locationCode"`
			} `json:"start"`
			End struct {
				DateTime string `json:"dateTime"`
				Address  struct {
					Line        string  `json:"line"`
					Zip         string  `json:"zip"`
					CountryCode string  `json:"countryCode"`
					CityName    string  `json:"cityName"`
					Latitude    float64 `json:"latitude"`
					Longitude   float64 `json:"longitude"`
				} `json:"address"`
				Name string `json:"name"`
			} `json:"end"`
			Vehicle struct {
				Code        string `json:"code"`
				Category    string `json:"category"`
				Description string `json:"description"`
				Baggages    []struct {
					Count int    `json:"count"`
					Size  string `json:"size"`
				} `json:"baggages"`
				Seats []struct {
					Count int `json:"count"`
				} `json:"seats"`
				ImageURL string `json:"imageURL"`
			} `json:"vehicle"`
			ServiceProvider struct {
				Code     string `json:"code"`
				Name     string `json:"name"`
				TermsURL string `json:"termsUrl"`
				LogoURL  string `json:"logoUrl"`
			} `json:"serviceProvider"`
			Quotation struct {
				MonetaryAmount string `json:"monetaryAmount"`
				CurrencyCode   string `json:"currencyCode"`
				Taxes          []struct {
					MonetaryAmount string `json:"monetaryAmount"`
				} `json:"taxes"`
				IsEstimated bool `json:"isEstimated"`
				TotalFees   struct {
					MonetaryAmount string `json:"monetaryAmount"`
				} `json:"totalFees"`
				TotalTaxes struct {
					MonetaryAmount string `json:"monetaryAmount"`
				} `json:"totalTaxes"`
			} `json:"quotation"`
			Converted struct {
				MonetaryAmount string `json:"monetaryAmount"`
				CurrencyCode   string `json:"currencyCode"`
				Taxes          []struct {
					MonetaryAmount string `json:"monetaryAmount"`
				} `json:"taxes"`
				IsEstimated bool `json:"isEstimated"`
				TotalFees   struct {
					MonetaryAmount string `json:"monetaryAmount"`
				} `json:"totalFees"`
				TotalTaxes struct {
					MonetaryAmount string `json:"monetaryAmount"`
				} `json:"totalTaxes"`
			} `json:"converted"`
		} `json:"transfers"`
	} `json:"data"`
}

type BookingErrorResponse struct {
	Errors []struct {
		Code   int    `json:"code"`
		Detail string `json:"detail"`
	} `json:"errors"`
}
