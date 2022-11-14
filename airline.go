package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Airline struct {
	Id           int16  `json: "id"`
	Airlinename  string `json: Airlinename`
	Iata         string `json: Iata`
	FaqURL       string `json: FaqURL`
	Passages     string `json: Passages`
	Friendlyness string `json: Friendliness`
}

func GetAirlineData(filename string) []Airline {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	data := []Airline{}

	_ = json.Unmarshal([]byte(file), &data)
	return data
}
