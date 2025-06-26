package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"os"
	// "io"
)


const (
	IPAPI = "http://ip-api.com/json"
	NSWPOINTSAPI = "https://api.weather.gov/points/"
)

type location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type forecastAPI struct {
	Properties struct {
		API string `json:"forecast"`
	} `json:"properties"`
}

func main() {
	client := &http.Client{
	
	}

	res, err := client.Get(IPAPI)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	var loc location
	if err = json.NewDecoder(res.Body).Decode(&loc); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Lat: %f; Lon: %f\n", loc.Lat, loc.Lon)

	nswApi := fmt.Sprintf("%s%0.4f,%0.4f", NSWPOINTSAPI, loc.Lat, loc.Lon)
	res, err = client.Get(nswApi)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	// body, err := io.ReadAll(res.Body)
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("Body: %s\n", body)

	var fapi forecastAPI
	if err := json.NewDecoder(res.Body).Decode(&fapi); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("ForecastApi: %s\n", fapi.Properties.API)
	res, err = client.Get(fapi.Properties.API)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	var forecast map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&forecast); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	prettyJSON, err := json.MarshalIndent(forecast, "", "   ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(prettyJSON))
}