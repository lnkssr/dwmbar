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

// Stats содержит температуру
type Stats struct {
	Temperature string
}

var (
	cache      *Stats
	cacheMutex sync.Mutex
	lastUpdate time.Time
)

// Get возвращает данные о погоде, обновляя их раз в 10 минут
func Get() (*Stats, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Если прошло менее 10 минут, вернуть кешированные данные
	if time.Since(lastUpdate) < 10*time.Minute && cache != nil {
		return cache, nil
	}

	// Обновляем данные о погоде
	weather, err := fetchWeather()
	if err != nil {
		return nil, err
	}

	cache = weather
	lastUpdate = time.Now()
	return cache, nil
}

// fetchWeather выполняет запрос к API и обновляет погоду
func fetchWeather() (*Stats, error) {
	ip, err := getPublicIP()
	if err != nil {
		log.Println("Ошибка при получении IP:", err)
		return nil, err
	}

	loc, err := getLocationByIP(ip)
	if err != nil {
		log.Println("Ошибка при определении местоположения:", err)
		return nil, err
	}

	temp, err := getTemperature(loc.Latitude, loc.Longitude)
	if err != nil {
		log.Println("Ошибка при получении температуры:", err)
		return nil, err
	}

	return &Stats{
		Temperature: fmt.Sprintf("%.1f°C", temp),
	}, nil
}

// Location содержит информацию о городе и координатах
type Location struct {
	City      string
	Latitude  float64
	Longitude float64
}

// getPublicIP получает публичный IP-адрес через ifconfig.me
func getPublicIP() (string, error) {
	resp, err := http.Get("https://ifconfig.me")
	if err != nil {
		return "", fmt.Errorf("ошибка при получении публичного IP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("не удалось получить публичный IP. Код: %d", resp.StatusCode)
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении ответа от ifconfig.me: %v", err)
	}

	return string(ip), nil
}

// getLocationByIP получает город и координаты по IP-адресу
func getLocationByIP(ip string) (*Location, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=city,lat,lon", ip)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении местоположения по IP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("не удалось получить данные о местоположении. Код: %d", resp.StatusCode)
	}

	var result struct {
		City string  `json:"city"`
		Lat  float64 `json:"lat"`
		Lon  float64 `json:"lon"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании данных о местоположении: %v", err)
	}

	if result.City == "" || result.Lat == 0 || result.Lon == 0 {
		return nil, errors.New("некорректные данные о городе и координатах")
	}

	return &Location{
		City:      result.City,
		Latitude:  result.Lat,
		Longitude: result.Lon,
	}, nil
}

// getTemperature получает текущую температуру через API Open-Meteo по координатам
func getTemperature(lat, lon float64) (float64, error) {
	weatherURL := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current_weather=true", lat, lon)
	resp, err := http.Get(weatherURL)
	if err != nil {
		return 0, fmt.Errorf("ошибка при получении данных о погоде: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("не удалось получить данные о погоде. Код: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("ошибка при чтении ответа от Open-Meteo: %v", err)
	}

	var data struct {
		CurrentWeather struct {
			Temperature float64 `json:"temperature"`
		} `json:"current_weather"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return 0, fmt.Errorf("ошибка при декодировании данных о погоде: %v", err)
	}

	return data.CurrentWeather.Temperature, nil
}
