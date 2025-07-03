package main

import (
	// "encoding/json"
	"fmt"
	"os"
	"github.com/carlmjohnson/requests"
	"context"
	"time"
	"strings"
)

const (
	IPAPI = "http://ip-api.com/json"
	NWSPOINTSAPI = "https://api.weather.gov/points/"
	LAYOUT =  "2 January 2006 at 03:04 PM MST"
)

type Period struct {
    Name                       string `json:"name"`
    Number                     int    `json:"number"`
    StartTime                  time.Time `json:"startTime"`   // Updated for time parsing
    EndTime                    time.Time `json:"endTime"`     // Updated for time parsing
    Icon                       string `json:"icon"`
    IsDaytime                  bool   `json:"isDaytime"`
    Temperature                int    `json:"temperature"`
    TemperatureUnit            string `json:"temperatureUnit"`
    ProbabilityOfPrecipitation struct {
        UnitCode string `json:"unitCode"`
        Value    *int   `json:"value"` // Nullable field example
    } `json:"probabilityOfPrecipitation"`
    ShortForecast    string `json:"shortForecast"`
    DetailedForecast string `json:"detailedForecast"`
    WindDirection    string `json:"windDirection"`
    WindSpeed        string `json:"windSpeed"`
}

type WeatherResponse struct {
    Properties struct {
        Periods    []Period  `json:"periods"`
        GeneratedAt time.Time `json:"generatedAt"`
        UpdateTime  time.Time `json:"updateTime"`
    } `json:"properties"`
}

type location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	TimeZone string `json:"timezone"`
}

type forecastAPI struct {
	Properties struct {
		API string `json:"forecast"`
	} `json:"properties"`
}

func main() {
	if err := getWeather(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// Requests a geo-ip api to get the lat long for the current system based on systems public ip address to get local weather
func getNWSApi(ctx context.Context, loc *location) (string, error) {
	// var loc location
	if err := requests.URL(IPAPI).ToJSON(&loc).Fetch(ctx); err != nil {
		return "", err
	}

	// print the lat long for testing
	// fmt.Printf("Lat: %f; Lon: %f\n", loc.Lat, loc.Lon)

	nwsApi := fmt.Sprintf("%s%0.4f,%0.4f", NWSPOINTSAPI, loc.Lat, loc.Lon)
	return nwsApi, nil
}

// Request the NSW points api with the coordinates to get the forecast url for the NWS office code responsible for the given coordinates
func getForecastApi(ctx context.Context, nwsApi string) (string, error) {
	var fapi forecastAPI
	if err := requests.URL(nwsApi).ToJSON(&fapi).Fetch(ctx); err != nil {
		return "", err
	}

	// Prints the forecast api resource for testing
	// fmt.Printf("ForecastApi: %s\n", fapi.Properties.API)
	return fapi.Properties.API, nil
}

// Requests forecast API for the NWS office code and parses it into a map[string]interface{} for now
func getForecast(ctx context.Context, forecast *WeatherResponse, foreCastAPI string) error {
	if err := requests.URL(foreCastAPI).ToJSON(&forecast).Fetch(ctx); err != nil {
		return err
	}
	return nil
}

func getWeather() error {
	ctx := context.TODO()

	var loc location
	nwsApi, err := getNWSApi(ctx, &loc)
	if err != nil {
		return err
	}
	
	foreCastAPI, err := getForecastApi(ctx, nwsApi)
	if err != nil {
		return err
	}

	var forecast WeatherResponse
	if err = getForecast(ctx, &forecast, foreCastAPI); err != nil {
		return err
	}

	var sb strings.Builder
	sb.WriteString("\nLocal Forecast:\n\n")
	for _, p := range forecast.Properties.Periods {
		tz, err := time.LoadLocation(loc.TimeZone)
		if err != nil {
			return err
		}
		adjustedStartTime := p.StartTime.In(tz)
		adjustedEndTime:= p.EndTime.In(tz)
		startTime := adjustedStartTime.Format(LAYOUT)
		endTime := adjustedEndTime.Format(LAYOUT)

		sb.WriteString(fmt.Sprintf("%s %s to %s\n", p.Name, startTime, endTime))
		sb.WriteString(fmt.Sprintf("Temp: %s | Wind: %s | Precip: %s\n\n", p.formatTemp(), p.formatWind(), p.formatPrecipitation()))
	}
	fmt.Println(sb.String())
	return nil
}

var wmoUnit = map[string]string{
	"percent": "%",
	"inches": "in",
	"centimeters": "cm",
}

func (p *Period) formatPrecipitation() string {
	unit := strings.Split(p.ProbabilityOfPrecipitation.UnitCode, ":")[1]
	unit = wmoUnit[unit]
	return fmt.Sprintf("%d%s", *p.ProbabilityOfPrecipitation.Value, unit)
}

func (p *Period) formatWind() string {
	return fmt.Sprintf("%s at %s", p.WindDirection, p.WindSpeed)
}

func (p *Period) formatTemp() string {
	return fmt.Sprintf("%d%s", p.Temperature, p.TemperatureUnit)
}

