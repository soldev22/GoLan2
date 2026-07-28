package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"weather-app/internal/models"
)

// Provider defines weather API operations used by services.
type Provider interface {
	Current(context.Context, string) (models.Weather, error)
	Forecast(context.Context, string) ([]models.ForecastItem, error)
}

var ErrInvalidCity = errors.New("city name is required")

const (
	defaultForecastURL  = "https://api.open-meteo.com/v1/forecast"
	defaultGeocodingURL = "https://geocoding-api.open-meteo.com/v1/search"
	defaultTimeout      = 5 * time.Second
)

// Config controls the endpoints and timeout used by the weather client.
type Config struct {
	ForecastURL    string
	GeocodingURL   string
	RequestTimeout time.Duration
}

// Client wraps a simple HTTP weather provider using Open-Meteo.
type Client struct {
	forecastURL  string
	geocodingURL string
	httpClient   *http.Client
}

func NewClient(cfg ...Config) *Client {
	clientConfig := Config{}
	if len(cfg) > 0 {
		clientConfig = cfg[0]
	}
	if clientConfig.ForecastURL == "" {
		clientConfig.ForecastURL = defaultForecastURL
	}
	if clientConfig.GeocodingURL == "" {
		clientConfig.GeocodingURL = defaultGeocodingURL
	}
	if clientConfig.RequestTimeout <= 0 {
		clientConfig.RequestTimeout = defaultTimeout
	}

	return &Client{
		forecastURL:  clientConfig.ForecastURL,
		geocodingURL: clientConfig.GeocodingURL,
		httpClient:   &http.Client{Timeout: clientConfig.RequestTimeout},
	}
}

func (c *Client) Current(ctx context.Context, city string) (models.Weather, error) {
	cleanCity, err := normalizeCity(city)
	if err != nil {
		return models.Weather{}, err
	}

	coords, err := c.lookupCoordinates(ctx, cleanCity)
	if err != nil {
		return models.Weather{}, err
	}

	params := url.Values{}
	params.Set("latitude", strconv.FormatFloat(coords.Latitude, 'f', -1, 64))
	params.Set("longitude", strconv.FormatFloat(coords.Longitude, 'f', -1, 64))
	params.Set("current", "temperature_2m,relative_humidity_2m,apparent_temperature,weather_code,wind_speed_10m")
	params.Set("timezone", "auto")

	endpoint := fmt.Sprintf("%s?%s", c.forecastURL, params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return models.Weather{}, fmt.Errorf("build current weather request: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return models.Weather{}, fmt.Errorf("request current weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Weather{}, fmt.Errorf("current weather lookup failed: %s", resp.Status)
	}

	var payload currentWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return models.Weather{}, fmt.Errorf("decode current weather response: %w", err)
	}

	return models.Weather{
		City:         titleizeCity(city),
		TemperatureC: int(payload.Current.Temperature2M),
		Condition:    weatherCodeToCondition(payload.Current.WeatherCode),
		FeelsLikeC:   int(payload.Current.ApparentTemperature),
		HumidityPct:  payload.Current.RelativeHumidity2M,
		WindKph:      int(payload.Current.WindSpeed10M),
		Description:  weatherCodeToCondition(payload.Current.WeatherCode),
	}, nil
}

func (c *Client) Forecast(ctx context.Context, city string) ([]models.ForecastItem, error) {
	cleanCity, err := normalizeCity(city)
	if err != nil {
		return nil, err
	}

	coords, err := c.lookupCoordinates(ctx, cleanCity)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Set("latitude", strconv.FormatFloat(coords.Latitude, 'f', -1, 64))
	params.Set("longitude", strconv.FormatFloat(coords.Longitude, 'f', -1, 64))
	params.Set("daily", "weather_code,temperature_2m_max,temperature_2m_min")
	params.Set("forecast_days", "3")
	params.Set("timezone", "auto")

	endpoint := fmt.Sprintf("%s?%s", c.forecastURL, params.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("build forecast request: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("request forecast: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("forecast lookup failed: %s", resp.Status)
	}

	var payload forecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode forecast response: %w", err)
	}

	items := make([]models.ForecastItem, 0, len(payload.Daily.Time))
	for i := range payload.Daily.Time {
		items = append(items, models.ForecastItem{
			Day:       formatForecastDay(payload.Daily.Time[i]),
			HighC:     int(payload.Daily.MaxTemp[i]),
			LowC:      int(payload.Daily.MinTemp[i]),
			Condition: weatherCodeToCondition(payload.Daily.WeatherCode[i]),
		})
	}
	return items, nil
}

type coordinateResponse struct {
	Results []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

type currentWeatherResponse struct {
	Current struct {
		Temperature2M       float64 `json:"temperature_2m"`
		RelativeHumidity2M  int     `json:"relative_humidity_2m"`
		ApparentTemperature float64 `json:"apparent_temperature"`
		WeatherCode         int     `json:"weather_code"`
		WindSpeed10M        float64 `json:"wind_speed_10m"`
	} `json:"current"`
}

type forecastResponse struct {
	Daily struct {
		Time        []string  `json:"time"`
		MaxTemp     []float64 `json:"temperature_2m_max"`
		MinTemp     []float64 `json:"temperature_2m_min"`
		WeatherCode []int     `json:"weather_code"`
	} `json:"daily"`
}

func (c *Client) lookupCoordinates(ctx context.Context, city string) (struct {
	Latitude  float64
	Longitude float64
}, error) {
	endpoint := fmt.Sprintf("%s?name=%s&count=1&language=en&format=json", c.geocodingURL, url.QueryEscape(city))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return struct {
			Latitude  float64
			Longitude float64
		}{}, fmt.Errorf("build geocoding request: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return struct {
			Latitude  float64
			Longitude float64
		}{}, fmt.Errorf("request geocoding: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return struct {
			Latitude  float64
			Longitude float64
		}{}, fmt.Errorf("geocoding lookup failed: %s", resp.Status)
	}

	var payload coordinateResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return struct {
			Latitude  float64
			Longitude float64
		}{}, fmt.Errorf("decode geocoding response: %w", err)
	}
	if len(payload.Results) == 0 {
		return struct {
			Latitude  float64
			Longitude float64
		}{}, fmt.Errorf("city not found: %s", city)
	}

	return struct {
		Latitude  float64
		Longitude float64
	}{
		Latitude:  payload.Results[0].Latitude,
		Longitude: payload.Results[0].Longitude,
	}, nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	client := c.httpClient
	if client == nil {
		client = &http.Client{Timeout: defaultTimeout}
	}
	return client.Do(req)
}

func weatherCodeToCondition(code int) string {
	switch code {
	case 0:
		return "Clear sky"
	case 1, 2, 3:
		return "Mostly cloudy"
	case 45, 48:
		return "Fog"
	case 51, 53, 55, 61, 63, 65:
		return "Rain"
	case 71, 73, 75, 77:
		return "Snow"
	case 80, 81, 82:
		return "Showers"
	case 95, 96, 99:
		return "Thunderstorm"
	default:
		return "Unknown"
	}
}

func titleizeCity(city string) string {
	parts := strings.Fields(strings.TrimSpace(city))
	for i, part := range parts {
		if len(part) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}
	return strings.Join(parts, " ")
}

func normalizeCity(city string) (string, error) {
	cleanCity := strings.TrimSpace(city)
	if cleanCity == "" {
		return "", ErrInvalidCity
	}
	return cleanCity, nil
}

func formatForecastDay(day string) string {
	parsed, err := time.Parse("2006-01-02", day)
	if err != nil {
		return day
	}
	return parsed.Format("Mon, Jan 2")
}
