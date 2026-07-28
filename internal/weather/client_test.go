package weather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClientCurrentAndForecast(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/geocoding":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"results":[{"latitude":55.9533,"longitude":-3.1883}]}`))
		case "/forecast":
			if r.URL.Query().Get("current") != "" {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"current":{"temperature_2m":19.2,"relative_humidity_2m":81,"apparent_temperature":18.4,"weather_code":2,"wind_speed_10m":22}}`))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"daily":{"time":["2026-07-28","2026-07-29","2026-07-30"],"temperature_2m_max":[20,21,18],"temperature_2m_min":[12,13,11],"weather_code":[2,61,0]}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient(Config{
		ForecastURL:    server.URL + "/forecast",
		GeocodingURL:   server.URL + "/geocoding",
		RequestTimeout: time.Second,
	})

	current, err := client.Current(context.Background(), "edinburgh")
	if err != nil {
		t.Fatalf("current: %v", err)
	}
	if current.City != "Edinburgh" || current.TemperatureC != 19 || current.Condition != "Mostly cloudy" {
		t.Fatalf("unexpected current weather: %+v", current)
	}

	forecast, err := client.Forecast(context.Background(), "edinburgh")
	if err != nil {
		t.Fatalf("forecast: %v", err)
	}
	if len(forecast) != 3 {
		t.Fatalf("expected 3 forecast items, got %d", len(forecast))
	}
	// Day must be the raw ISO date for machine-readable attributes
	if forecast[0].Day != "2026-07-28" {
		t.Fatalf("expected raw ISO date in Day, got %q", forecast[0].Day)
	}
	// DayLabel must be the human-readable formatted string
	if !strings.Contains(forecast[0].DayLabel, "Jul") {
		t.Fatalf("expected formatted forecast day label, got %+v", forecast[0])
	}
}

func TestClientRejectsBlankCity(t *testing.T) {
	client := NewClient()

	if _, err := client.Current(context.Background(), "   "); err != ErrInvalidCity {
		t.Fatalf("expected ErrInvalidCity, got %v", err)
	}
	if _, err := client.Forecast(context.Background(), ""); err != ErrInvalidCity {
		t.Fatalf("expected ErrInvalidCity, got %v", err)
	}
}

func TestClientPropagatesGeocodingFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(Config{
		ForecastURL:    server.URL + "/forecast",
		GeocodingURL:   server.URL + "/geocoding",
		RequestTimeout: time.Second,
	})

	_, err := client.Current(context.Background(), "edinburgh")
	if err == nil || !strings.Contains(err.Error(), "geocoding lookup failed") {
		t.Fatalf("expected geocoding failure, got %v", err)
	}
}

func TestClientHonorsTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results":[{"latitude":55.9533,"longitude":-3.1883}]}`))
	}))
	defer server.Close()

	client := NewClient(Config{
		ForecastURL:    server.URL + "/forecast",
		GeocodingURL:   server.URL + "/geocoding",
		RequestTimeout: 10 * time.Millisecond,
	})

	_, err := client.Current(context.Background(), "edinburgh")
	if err == nil {
		t.Fatal("expected timeout error")
	}
}
