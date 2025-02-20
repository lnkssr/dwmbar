package weather_state

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Stats struct {
	Temperature string
}

var (
	cache      *Stats
	cacheMutex sync.Mutex
	lastUpdate time.Time
)

func Get() (*Stats, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if time.Since(lastUpdate) < 30*time.Minute && cache != nil {
		return cache, nil
	}

	weather, err := fetchWeather()
	if err != nil {
		return nil, err
	}

	cache = weather
	lastUpdate = time.Now()
	return cache, nil
}

func fetchWeather() (*Stats, error) {
	ip, err := getPublicIP()
	if err != nil {
		log.Println("Error getting IP:", err)
		return nil, err
	}

	loc, err := getLocationByIP(ip)
	if err != nil {
		log.Println("Error determining location:", err)
		return nil, err
	}

	temp, err := getTemperature(loc.Latitude, loc.Longitude)
	if err != nil {
		log.Println("Error getting temperature:", err)
		return nil, err
	}

	return &Stats{
		Temperature: fmt.Sprintf("%.1f°C", temp),
	}, nil
}

type Location struct {
	City      string
	Latitude  float64
	Longitude float64
}

func getPublicIP() (string, error) {
	resp, err := http.Get("https://api64.ipify.org?format=text")
	if err != nil {
		return "", fmt.Errorf("error getting public IP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get public IP. Code: %d", resp.StatusCode)
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response from api64.ipify.org: %v", err)
	}

	log.Println("Public IP:", string(ip))
	return string(ip), nil
}

func getLocationByIP(ip string) (*Location, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=city,lat,lon", ip)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error getting location by IP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get location data. Code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	log.Println("Location API response:", string(body))

	var result struct {
		City string  `json:"city"`
		Lat  float64 `json:"lat"`
		Lon  float64 `json:"lon"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decoding location data: %v", err)
	}

	if result.City == "" || result.Lat == 0 || result.Lon == 0 {
		return nil, errors.New("invalid city and coordinates data")
	}

	return &Location{
		City:      result.City,
		Latitude:  result.Lat,
		Longitude: result.Lon,
	}, nil
}

func getTemperature(lat, lon float64) (float64, error) {
	weatherURL := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current_weather=true", lat, lon)
	resp, err := http.Get(weatherURL)
	if err != nil {
		return 0, fmt.Errorf("error getting weather data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get weather data. Code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	log.Println("Weather API response:", string(body))

	var data struct {
		CurrentWeather struct {
			Temperature float64 `json:"temperature"`
		} `json:"current_weather"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return 0, fmt.Errorf("error decoding weather data: %v", err)
	}

	return data.CurrentWeather.Temperature, nil
}
