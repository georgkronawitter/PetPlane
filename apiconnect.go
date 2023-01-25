package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var formOrigin string
var formDestination string
var formDepartureDate string

const (
	prompt    = "Tell me if a five kilo dog can travel with airline "
	model     = "text-davinci-003"
	openaiurl = "https://api.openai.com/v1/completions"
	maxChars  = 200
)

func getBearerToken() string {
	data := url.Values{"grant_type": {"client_credentials"},
		"client_id":     {"kxKK7o0FAgdqBjpURgy47i1KShqqwNfh"},
		"client_secret": {"3p3lgOAQxJzNA6Av"}}

	var authToken struct {
		Type             string
		Username         string
		Application_name string
		Client_id        string
		Token_type       string
		Access_token     string
		Expires_in       int32
		State            string
		Scope            string
	}

	resp, err := http.PostForm("https://test.api.amadeus.com/v1/security/oauth2/token", data)
	if err != nil {
		log.Fatal(err)
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(responseBody, &authToken)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(authToken)
	token := authToken.Access_token
	return token
}

//define Response struct
var Response struct {
	Meta struct {
		Count int `json:"count"`
	} `json:"meta"`

	Data []struct {
		Type                     string `json:"type"`
		ID                       string `json:"id"`
		Source                   string `json:"source"`
		InstantTicketingRequired bool   `json:"instantTicketingRequired"`
		NonHomogeneous           bool   `json:"nonHomogeneous"`
		OneWay                   bool   `json:"oneWay"`
		LastTicketingDate        string `json:"lastTicketingDate"`
		NumberOfBookableSeats    int    `json:"numberOfBookableSeats"`
		Itineraries              []struct {
			Duration string `json:"duration"`
			Segments []struct {
				Departure struct {
					IataCode string `json:"iataCode"`
					Terminal string `json:"terminal"`
					At       string `json:"at"`
				} `json:"departure"`
				Arrival struct {
					IataCode string `json:"iataCode"`
					Terminal string `json:"terminal"`
					At       string `json:"at"`
				} `json:"arrival,omitempty"`
				CarrierCode string `json:"carrierCode"`
				Number      string `json:"number"`
				Aircraft    struct {
					Code string `json:"code"`
				} `json:"aircraft"`
				Operating struct {
					CarrierCode string `json:"carrierCode"`
				} `json:"operating"`
				Duration        string `json:"duration"`
				ID              string `json:"id"`
				NumberOfStops   int    `json:"numberOfStops"`
				BlacklistedInEU bool   `json:"blacklistedInEU"`
			} `json:"segments"`
		} `json:"itineraries"`
		Price struct {
			Currency string `json:"currency"`
			Total    string `json:"total"`
			Base     string `json:"base"`
			Fees     []struct {
				Amount string `json:"amount"`
				Type   string `json:"type"`
			} `json:"fees"`
			GrandTotal string `json:"grandTotal"`
		} `json:"price"`
		PricingOptions struct {
			FareType                []string `json:"fareType"`
			IncludedCheckedBagsOnly bool     `json:"includedCheckedBagsOnly"`
		} `json:"pricingOptions"`
		ValidatingAirlineCodes []string `json:"validatingAirlineCodes"`
		TravelerPricings       []struct {
			TravelerID   string `json:"travelerId"`
			FareOption   string `json:"fareOption"`
			TravelerType string `json:"travelerType"`
			Price        struct {
				Currency string `json:"currency"`
				Total    string `json:"total"`
				Base     string `json:"base"`
			} `json:"price"`
			FareDetailsBySegment []struct {
				SegmentID           string `json:"segmentId"`
				Cabin               string `json:"cabin"`
				FareBasis           string `json:"fareBasis"`
				Class               string `json:"class"`
				IncludedCheckedBags struct {
					Weight     int    `json:"weight"`
					WeightUnit string `json:"weightUnit"`
				} `json:"includedCheckedBags"`
			} `json:"fareDetailsBySegment"`
		} `json:"travelerPricings"`
	} `json:"data"`
	Dictionaries struct {
		Locations struct {
			Bkk struct {
				CityCode    string `json:"cityCode"`
				CountryCode string `json:"countryCode"`
			} `json:"BKK"`
			Mnl struct {
				CityCode    string `json:"cityCode"`
				CountryCode string `json:"countryCode"`
			} `json:"MNL"`
			Syd struct {
				CityCode    string `json:"cityCode"`
				CountryCode string `json:"countryCode"`
			} `json:"SYD"`
		} `json:"locations"`
		Aircraft struct {
			Num321 string `json:"321"`
			Num333 string `json:"333"`
		} `json:"aircraft"`
		Currencies struct {
			Eur string `json:"EUR"`
		} `json:"currencies"`
		Carriers struct {
			Pr string `json:"PR"`
		} `json:"carriers"`
	} `json:"dictionaries"`
}

//define Request struct
type Request struct {
	CurrencyCode       string               `json:"currencyCode"`
	OriginDestinations []OriginDestinations `json:"originDestinations"`
	Travelers          []Travelers          `json:"travelers"`
	Sources            []string             `json:"sources"`
}

type OriginDestinations struct {
	ID                      string                 `json:"id"`
	OriginLocationCode      string                 `json:"originLocationCode"`
	DestinationLocationCode string                 `json:"destinationLocationCode"`
	DepartureDateTimeRange  DepartureDateTimeRange `json:"departureDateTimeRange"`
}

type DepartureDateTimeRange struct {
	Date string `json:"date"`
	Time string `json:"time"`
}

type Travelers struct {
	ID           string `json:"id"`
	TravelerType string `json:"travelerType"`
}

/*  type    SearchCriteria struct {
        MaxFlightOffers int `json:"maxFlightOffers"`
        FlightFilters   struct {
            CabinRestrictions []struct {
                Cabin                string   `json:"cabin"`
                Coverage             string   `json:"coverage"`
                OriginDestinationIds []string `json:"originDestinationIds"`
            } `json:"cabinRestrictions"`
            CarrierRestrictions struct {
                ExcludedCarrierCodes []string `json:"excludedCarrierCodes"`
            } `json:"carrierRestrictions"`
        } `json:"flightFilters"`
    }
*/

type request struct {
	Prompt           string  `json:"prompt"`
	Model            string  `json:"model"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float32 `json:"temperature"`
	TopP             int     `json:"top_p"`
	N                int     `json:"n"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
}
type choices struct {
	Text string `json:"text"`
}

type OpenAIResponse struct {
	ID        string    `json:"id"`
	Model     string    `json:"model"`
	Prompt    string    `json:"prompt"`
	Completed string    `json:"completed"`
	Choices   []choices `json:"choices"`
}

func sendAmadeusPriceRequest(itineraries string) []byte {
	//retrieve Bearer Token
	token := getBearerToken()
	fmt.Println("Authentication token received: ", token)
	//fill Request instance wth Request data
	requestData := "'data': { 'type': 'flight-offers-pricing', 'flightOffers': [ { 'type': 'flight-offer', 'id': '1', 'source': 'GDS', 'instantTicketingRequired': false, 'nonHomogeneous': false, 'oneWay': false, 'lastTicketingDate': '2020-08-04', 'numberOfBookableSeats': 9, 'itineraries': [" + itineraries + "]}]}"
	url := "https://test.api.amadeus.com/v1/shopping/flight-offers/pricing?forceClass=false"
	//marshall request Data into b, then turn it into bytes
	b, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println(err)
	}
	b = []byte(b)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	fmt.Printf(requestData)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("\nPrice search request sent, response:", string(body))
	return body
}

func sendOpenAIRequest(airline string) string {
	reqBody := request{
		Prompt:           prompt + airline,
		Model:            model,
		MaxTokens:        maxChars,
		Temperature:      0.5,
		TopP:             1,
		N:                1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", openaiurl, bytes.NewBuffer(reqBytes))
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+OpenAIApiKey)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	defer res.Body.Close()
	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	var apiRes OpenAIResponse
	if err := json.Unmarshal(resBytes, &apiRes); err != nil {
		fmt.Println(err)
		return "error"
	}
	return apiRes.Choices[0].Text
}

func sendAmadeusSearchRequest(origin string, destination string, departuredate string) []byte {
	//retrieve Bearer Token
	token := getBearerToken()
	fmt.Println("Authentication token for search received: ", token)
	//fill Request instance wth Request data
	requestData := Request{CurrencyCode: "USD", OriginDestinations: []OriginDestinations{{ID: "1", OriginLocationCode: origin, DestinationLocationCode: destination, DepartureDateTimeRange: DepartureDateTimeRange{Date: departuredate, Time: "00:10:10"}}}, Travelers: []Travelers{{ID: "1", TravelerType: "ADULT"}, {ID: "2", TravelerType: "ADULT"}}, Sources: []string{"GDS"} /* SearchCriteria:  struct {MaxFlightOffers: 20,FlightFilters: struct {CabinRestrictions: []struct {}, CarrierRestrictions: struct {ExcludedCarrierCodes: []}}}}*/}
	url := "https://test.api.amadeus.com/v2/shopping/flight-offers"
	//Marshall Request Data into b, then turn it into bytes
	b, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println(err)
	}
	b = []byte(b)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//header := string(resp.StatusCode) + string(http.StatusText(resp.StatusCode))
	fmt.Println(string(b))
	fmt.Println("Flight search request sent")
	//fmt.Println(string(body), header)
	return body
}

//retrieve airline policy data
var airlinedata = GetAirlineData("airlinedata.json")
