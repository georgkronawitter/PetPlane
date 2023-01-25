package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var AirlinePetPolicies = make(map[string]string)
var initialHTML string = "<html><head><script src='/files/script.js'></script><link rel='stylesheet' type='text/css' href='/files/style.css' /><meta charset='UTF-8' /><title>PetPlane</title></head><body><div id=''></div><div><form method='POST' class='search' action='/form'><div class='logo'>PetPlane</div><label>&#11014;&#65039; Origin</label><input name='origin' size=6 type='text' value='" + formOrigin + "' /> <label>&#11015;&#65039; Destination</label><input name='destination' size=6 type='text' value='" + formDestination + "' /> <label>Departure Date</label><input name='departuredate' type='date' value='" + formDepartureDate + "' /> <input type='submit' class='button' value=' &#128269; Search' /> </form></div>"

func priceRequestHandler(w http.ResponseWriter, r *http.Request) {
	formItinerary := strings.ToUpper(r.FormValue("itinerary"))
	output := initialHTML
	sendAmadeusPriceRequest(formItinerary)
	fmt.Fprintf(w, output)
	return
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	formOrigin = strings.ToUpper(r.FormValue("origin"))
	formDestination = strings.ToUpper(r.FormValue("destination"))
	formDepartureDate = r.FormValue("departuredate")
	flightdata := sendAmadeusSearchRequest(formOrigin, formDestination, formDepartureDate) //send request to Amadeus API with origin, destination, date
	err := json.Unmarshal(flightdata, &Response)
	fmt.Printf(string(flightdata))
	type flights struct {
		Price       string
		Airlinename string
		Time        string
		FaqURL      string
	}

	if err != nil {
		fmt.Println(err)
	}

	hasNoDirectFlights := true
	output := initialHTML

	//iterate over flight search Response.data.itineraries.segments (=every single flight), policy for each carrier
	for _, data := range Response.Data {
		for _, itinerary := range data.Itineraries {
			if len(itinerary.Segments) < 2 { //only show direct flights
				//save itinerary data
				flightOfferJson, err := json.Marshal(data)
				if err != nil {
					fmt.Println(err)
				}

				//render price
				output += "<div class='result'><form name='getPrice' method='POST' action='/pricerequest'><input type='hidden' name='itinerary' value='" + string(flightOfferJson) + "'><a href='javascript:document.getPrice.submit()'><div class='price'>" + data.Price.Total + " " + data.Price.Currency + "</div>&#9; "
				hasNoDirectFlights = false
				for _, segment := range itinerary.Segments {
					//render flight details
					output += "<div class='flightdetails'>&nbsp;&#9992;&#65039; " + segment.Departure.IataCode + " to " + segment.Arrival.IataCode + ",&#9at " + segment.Departure.At[11:16] + ", with <b><script>document.write(IATAmapper('" + segment.CarrierCode + "'));</script></b></div></a></form>"
					for _, airline := range airlinedata {
						if airline.Iata == segment.CarrierCode {
							if AirlinePetPolicies[airline.Iata] == "" { //retrieve OpenAI assessment for airline pet policy if none present
								fmt.Println("connecting with ChatGPTfor pet policy info for " + airline.Iata)
								AirlinePetPolicies[airline.Iata] = sendOpenAIRequest(airline.Iata)
							} else {
								fmt.Println("ChatGPT policy found for " + airline.Iata)
							}
							output += "<p><div class='petpolicydb'><a href='" + airline.FaqURL + "' alt='Link to pet policy of airline" + airline.Iata + "' target='_blank'>🏷️ " + airline.Passages[0:270] + "[...]</a></div><div class='petpolicyopenai'>🤖 >> " + AirlinePetPolicies[airline.Iata] + " << </div>"
						} //end of iteration over airlines

					}

				}
				output += "</div>" //close 'results' div
			}
		}
	}
	if hasNoDirectFlights {
		output += "<p><b>No direct flights found</b></p>"
	}
	output += "</body></html>"
	fmt.Fprintf(w, output)
	return
}

func searchRequestHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	formOrigin = strings.ToUpper(r.FormValue("origin"))
	formDestination = strings.ToUpper(r.FormValue("destination"))
	formDepartureDate = r.FormValue("departuredate")
	flightdata := sendAmadeusSearchRequest(formOrigin, formDestination, formDepartureDate)
	err := json.Unmarshal(flightdata, &Response)
	if err != nil {
		fmt.Println(err)
	}
	//iterate over flight search Response.data.itineraries.segments (=every single flights), append policy for each carrier to IATA string
	for _, data := range Response.Data {
		for _, itinerary := range data.Itineraries {
			for _, segment := range itinerary.Segments {
				for _, airline := range airlinedata {
					if airline.Iata == segment.CarrierCode {
						if AirlinePetPolicies[airline.Iata] == "" {
							fmt.Println("connecting with ChatGPTfor pet policy info")
							AirlinePetPolicies[airline.Iata] = sendOpenAIRequest(airline.Iata)
						}
						airline.Iata += "\t" + airline.FaqURL + "\t" + airline.Passages + AirlinePetPolicies[airline.Iata]
					} //end of iteration over airlines
				}
			}
		}
	}
	jsonResponse, err := json.Marshal(Response)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(jsonResponse)
}

func server() {
	fileServer := http.FileServer(http.Dir("./"))
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/searchrequest", searchRequestHandler)
	http.HandleFunc("/pricerequest", priceRequestHandler)
	http.Handle("/", fileServer)
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}

//main function
func main() {
	//retrieve flights for given origin/destination combo
	server()
}
