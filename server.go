package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	formOrigin = strings.ToUpper(r.FormValue("origin"))
	formDestination = strings.ToUpper(r.FormValue("destination"))
	formDepartureDate = r.FormValue("departuredate")
	fmt.Fprintf(w, "<p>Origin = %s", formOrigin)
	fmt.Fprintf(w, "<p>Destination = %s\n", formDestination)
	fmt.Fprintf(w, "<p>Date of Departure = %s\n", formDepartureDate)
	flightdata := sendRequest(formOrigin, formDestination, formDepartureDate)
	err := json.Unmarshal(flightdata, &Response)

	type flights struct {
		Price       string
		Airlinename string
		Time        string
		FaqURL      string
	}

	if err != nil {
		fmt.Println(err)
	}
	var output string
	output = "<!DOCTYPE html><html><style>body {font-family: Arial; margin: 10px; padding: 10px;} p    {color: black;}</style><head><meta charset='UTF-8' /><title>PetPlane</title></head><body>"
	//output = ""
	//iterate over flight search Response.data.itineraries.segments (=every single flight), policy for each carrier
	for _, data := range Response.Data {
		for _, itinerary := range data.Itineraries {
			if len(itinerary.Segments) < 2 { //only show direct flights
				output += "<p><b>" + data.Price.Total + " " + data.Price.Currency + "</b>&#9;"
				for _, segment := range itinerary.Segments {
					output += segment.Departure.IataCode + " to " + segment.Arrival.IataCode + ", At:&#9;" + segment.Departure.At + ".&#9;Carrier " + segment.CarrierCode
					for _, airline := range airlinedata {
						if airline.Iata == segment.CarrierCode {
							output += "<a href='" + airline.FaqURL + "'>" + airline.Passages + "</a>"
						} //end of iteration over airlines
					}
				}
			}
		}
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
	flightdata := sendRequest(formOrigin, formDestination, formDepartureDate)
	err := json.Unmarshal(flightdata, &Response)
	if err != nil {
		fmt.Println(err)
	}
	//iterate over flight search Response.data.itineraries.segments (=every single flights), show policy for each carrier
	for _, data := range Response.Data {
		for _, itinerary := range data.Itineraries {
			for _, segment := range itinerary.Segments {
				for _, airline := range airlinedata {
					if airline.Iata == segment.CarrierCode {
						airline.Iata += "\t" + airline.FaqURL + "\t" + airline.Passages
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
	http.Handle("/", fileServer)
	http.HandleFunc("/searchrequest", searchRequestHandler)
	http.HandleFunc("/form", formHandler)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
