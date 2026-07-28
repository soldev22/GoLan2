package services

import (
	"context"
	"testing"

	"weather-app/internal/models"
)

type stubProvider struct {
	current  models.Weather
	forecast []models.ForecastItem
	err      error
}

func (s *stubProvider) Current(context.Context, string) (models.Weather, error) {
	return s.current, s.err
}

func (s *stubProvider) Forecast(context.Context, string) ([]models.ForecastItem, error) {
	return s.forecast, s.err
}

func TestWeatherPageUsesProviderData(t *testing.T) {
	provider := &stubProvider{
		current: models.Weather{
			City:         "London",
			TemperatureC: 18,
			Condition:    "Cloudy",
		},
		forecast: []models.ForecastItem{{
			Day:       "Mon",
			HighC:     20,
			LowC:      12,
			Condition: "Cloudy",
		}},
	}

	svc := NewPageService(provider)
	page := svc.Weather("London")

	if page.Weather.City != "London" {
		t.Fatalf("expected city London, got %s", page.Weather.City)
	}
	if !page.HasWeather {
		t.Fatalf("expected weather data to be present")
	}
	if len(page.Forecast) != 1 {
		t.Fatalf("expected one forecast item, got %d", len(page.Forecast))
	}
}
