
package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "bytes"
    "net/http"
    "net/url"
    "encoding/json"
)

var formOrigin string
var formDestination string
var formDepartureDate string

func getBearerToken () string {  
    data := url.Values{"grant_type":   {"client_credentials"},
    "client_id":    {"Mhg6yvAQn6hHARz5Us5ecR0mBRIlnAGb"},
    "client_secret":    {"kjvpKAuSCzLdIx6r"}}
    
    var authToken struct {
        Type string
        Username    string
        Application_name    string
        Client_id   string
        Token_type  string
        Access_token    string
        Expires_in  int32
        State   string
        Scope   string
    }
    
    resp, err := http.PostForm("https://test.api.amadeus.com/v1/security/oauth2/token", data)
    if err != nil { log.Fatal(err)}
    responseBody, err := ioutil.ReadAll(resp.Body)    
    if err != nil { log.Fatal(err)}
    err = json.Unmarshal(responseBody, &authToken)
    if err != nil { log.Fatal(err)}
    token := authToken.Access_token
    return token
}

//defining Response struct
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
    CurrencyCode       string `json:"currencyCode"`
    OriginDestinations []OriginDestinations `json:"originDestinations"`
    Travelers   []Travelers `json:"travelers"`
    Sources        []string `json:"sources"`
}

type OriginDestinations struct {
        ID                      string `json:"id"`
        OriginLocationCode      string `json:"originLocationCode"`
        DestinationLocationCode string `json:"destinationLocationCode"`
        DepartureDateTimeRange DepartureDateTimeRange `json:"departureDateTimeRange"`
    }

type    DepartureDateTimeRange  struct{
            Date string `json:"date"`
            Time string `json:"time"`
        }

type    Travelers struct {
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

func sendRequest(origin string, destination string, departuredate string) []byte {

    //retrieve Bearer Token
    token := getBearerToken()
    fmt.Println("Authentication token received")
    
    //fill Request instance wth Request data
    requestData := Request{CurrencyCode: "USD", OriginDestinations: []OriginDestinations{{ID: "1", OriginLocationCode: origin, DestinationLocationCode: destination, DepartureDateTimeRange: DepartureDateTimeRange{Date: departuredate, Time: "00:10:10"}}}, Travelers: []Travelers{{ID: "1", TravelerType: "ADULT"}},Sources: []string{"GDS"}/* SearchCriteria:  struct {MaxFlightOffers: 20,FlightFilters: struct {CabinRestrictions: []struct {}, CarrierRestrictions: struct {ExcludedCarrierCodes: []}}}}*/}
    url := "https://test.api.amadeus.com/v2/shopping/flight-offers"

    //Marshall Request Data into b, then turn it into bytes
    b, err := json.Marshal(requestData)
    if(err!=nil){fmt.Println(err)}

    b = []byte(b)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
    req.Header.Add("Authorization", "Bearer "+token)
    req.Header.Add("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil { fmt.Println(err) }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    //header := string(resp.StatusCode)+string(http.StatusText(resp.StatusCode))
    fmt.Println("Flight search request sent")
    fmt.Println(string(body))
    return body
}

//retrieve airline policy data
var airlinedata = GetAirlineData("airlinedata.json")

//main function
func main() {  

    //retrieve flights for given origin/destination combo
    server()
}