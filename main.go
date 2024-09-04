package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const apiKey = ""

type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Name string `json:"name"`
}

// flag package to allow users specify city, celsius or farhrenheit
func main() {
	city := flag.String("city", "London", "City to fetch the weather for")
	unit := flag.String("unit", "C", "Temperature unit: C (Celsius) or F (Fahrenheit)")
	flag.Parse()

	weather, err := getWeather(*city)
	if err != nil {
		log.Fatal("Error fetching weather:", err)
	}

	temp := weather.Main.Temp
	if *unit == "F" {
		temp = temp*9/5 + 32
	}

	fmt.Printf("Current weather in %s: %.2f°%s\n", weather.Name, temp, *unit)

}

// fetch weather data

func getWeather(city string) (WeatherResponse, error) {
	var weather WeatherResponse
	// Encode city name to handle spaces or special characters
	city = url.QueryEscape(city)
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return weather, err
	}
	defer resp.Body.Close()

	// Check if the status code is 200 (OK)
	if resp.StatusCode != http.StatusOK {
		// Read the response body (which may be HTML or JSON with an error message)
		var bodyBytes []byte
		bodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return weather, fmt.Errorf("failed to read response body: %w", err)
		}
		return weather, fmt.Errorf("received non-200 response code: %d. Response body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Now try to decode the JSON response
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return weather, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return weather, nil
}
