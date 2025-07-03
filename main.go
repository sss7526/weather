package main

import (
    "fmt"
    "os"
    "github.com/carlmjohnson/requests"
    "context"
    "time"
    "strings"
)

const (
    // IPAPI is the endpoint for obtaining the approximate geolocation of the user based on their public IP address.
    IPAPI = "http://ip-api.com/json"

    // NWSPOINTSAPI is the base URL for requesting information from the National Weather Service points API.
    NWSPOINTSAPI = "https://api.weather.gov/points/"

    // LAYOUT defines the output format for formatting times into a human-readable format.
    LAYOUT = "2 January 2006 at 03:04 PM MST"
)

// wmoUnit maps World Meteorological Organization (WMO) unit codes to more user-friendly units.
var wmoUnit = map[string]string{
    "percent":     "%",
    "inches":      "in",
    "centimeters": "cm",
}

// Period represents a single forecast period from the NWS API, including temperature, time, and other weather data.
type Period struct {
    Name                       string `json:"name"` 			// The name of the period (e.g., "Today", "Tonight").
    Number                     int    `json:"number"`
    StartTime                  time.Time `json:"startTime"`  	// The start time of this forecast period.
    EndTime                    time.Time `json:"endTime"`    	// The end time of this forecast period.
    Icon                       string `json:"icon"`          	// URL of the weather icon representing this period.
    IsDaytime                  bool   `json:"isDaytime"`     	// Whether this period represents daytime hours.
    Temperature                int    `json:"temperature"`   	// The temperature during this period.
    TemperatureUnit            string `json:"temperatureUnit"` 	// The unit of the temperature (e.g., "F" for Fahrenheit).
    ProbabilityOfPrecipitation struct {
        UnitCode string `json:"unitCode"` 						// The unit of the precipitation probability (usually "percent").
        Value    *int   `json:"value"`    						// The probability of precipitation as a percentage.
    } `json:"probabilityOfPrecipitation"`
    ShortForecast    string `json:"shortForecast"`    			// A brief description of the forecast (e.g., "Sunny").
    DetailedForecast string `json:"detailedForecast"` 			// A detailed description of the forecast.
    WindDirection    string `json:"windDirection"`    			// The direction from which the wind is coming.
    WindSpeed        string `json:"windSpeed"`        			// The speed of the wind.
}

// WeatherResponse represents the top-level response from the NWS API containing metadata and forecast periods.
type WeatherResponse struct {
    Properties struct {
        Periods    []Period  `json:"periods"`      // The forecast periods.
        GeneratedAt time.Time `json:"generatedAt"` // The time at which the forecast was generated.
        UpdateTime  time.Time `json:"updateTime"`  // The time when the forecast was last updated.
    } `json:"properties"`
}

// location represents the geographical location returned by the IP geolocation API.
type location struct {
    Lat          float64 `json:"lat"`    	// Latitude of the user's location.
    Lon          float64 `json:"lon"`    	// Longitude of the user's location.
    TimeZone     string  `json:"timezone"` 	// The time zone corresponding to the location.
    NWSPointsAPI string  					// The URL for the NWS points API for this location.
    ForecastAPI  string  					// The URL for the specific forecast API.
}

// forecastAPI represents the response from the points API used to retrieve the forecast URL.
type forecastAPI struct {
    Properties struct {
        API string `json:"forecast"` // The URL for the forecast data.
    } `json:"properties"`
}

// formatPrecipitation formats the probability of precipitation into a user-friendly string.
//
// If the probability field is nil or invalid, this may panic.
func (p *Period) formatPrecipitation() string {
    unit := strings.Split(p.ProbabilityOfPrecipitation.UnitCode, ":")[1]
    unit = wmoUnit[unit]
    return fmt.Sprintf("%d%s", *p.ProbabilityOfPrecipitation.Value, unit)
}

// formatWind formats the wind direction and speed into a readable string.
func (p *Period) formatWind() string {
    return fmt.Sprintf("%s at %s", p.WindDirection, p.WindSpeed)
}

// formatTemp formats the temperature and its unit into a readable string.
func (p *Period) formatTemp() string {
    return fmt.Sprintf("%d%s", p.Temperature, p.TemperatureUnit)
}

// main is the entry point for the program.
// It retrieves and prints the weather forecast for the local area based on the current public IP address.
func main() {
    if err := getWeather(); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}

// getNWSApi retrieves the NWS Points API URL for a given latitude and longitude.
//
// It updates the `NWSPointsAPI` field of the provided location struct.
// Returns an error if the IP geolocation API request fails or the response is invalid.
func getNWSApi(ctx context.Context, loc *location) error {
    if err := requests.URL(IPAPI).ToJSON(&loc).Fetch(ctx); err != nil {
        return err
    }

    loc.NWSPointsAPI = fmt.Sprintf("%s%0.4f,%0.4f", NWSPOINTSAPI, loc.Lat, loc.Lon)
    return nil
}

// getForecastApi retrieves the forecast URL for the location from the NWS Points API.
//
// It updates the `ForecastAPI` field of the location struct. Returns an error if the API request fails.
func getForecastApi(ctx context.Context, loc *location) error {
    var fapi forecastAPI
    if err := requests.URL(loc.NWSPointsAPI).ToJSON(&fapi).Fetch(ctx); err != nil {
        return err
    }
    loc.ForecastAPI = fapi.Properties.API
    return nil
}

// getForecast retrieves the actual forecast data from the NWS forecast API URL.
//
// The weather data is stored in the provided WeatherResponse pointer.
// Returns an error if the API request fails or the response is invalid.
func getForecast(ctx context.Context, forecast *WeatherResponse, foreCastAPI string) error {
    if err := requests.URL(foreCastAPI).ToJSON(&forecast).Fetch(ctx); err != nil {
        return err
    }
    return nil
}

// getWeather is the main logic for retrieving and displaying the local weather.
//
// It uses geolocation to determine the user's coordinates, retrieves the forecast data, 
// and then formats and prints it to the terminal.
//
// Returns an error if any step of the process fails.
func getWeather() error {
    ctx := context.TODO()

    var loc location
    if err := getNWSApi(ctx, &loc); err != nil {
        return err
    }
    if err := getForecastApi(ctx, &loc); err != nil {
        return err
    }

    var forecast WeatherResponse
    if err := getForecast(ctx, &forecast, loc.ForecastAPI); err != nil {
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
        adjustedEndTime := p.EndTime.In(tz)
        startTime := adjustedStartTime.Format(LAYOUT)
        endTime := adjustedEndTime.Format(LAYOUT)

        sb.WriteString(fmt.Sprintf("%s %s to %s\n", p.Name, startTime, endTime))
        sb.WriteString(fmt.Sprintf("Temp: %s | Wind: %s | Precip: %s\n", p.formatTemp(), p.formatWind(), p.formatPrecipitation()))
		sb.WriteString(fmt.Sprintf("%s\n\n", p.DetailedForecast))
    }
    fmt.Println(sb.String())
    return nil
}