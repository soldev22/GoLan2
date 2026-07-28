package services

import (
	"context"
	"strings"
	"time"

	"weather-app/internal/models"
	"weather-app/internal/weather"
)

type PageService interface {
	Home() models.PageData
	Weather(city string) models.PageData
	Forecast(city string) models.PageData
	About() models.PageData
	Error() models.PageData
}

type pageService struct {
	provider weather.Provider
}

func NewPageService(provider weather.Provider) PageService {
	if provider == nil {
		provider = weather.NewClient()
	}
	return &pageService{provider: provider}
}

func (s *pageService) Home() models.PageData {
	return models.PageData{
		Title:           "Home",
		Heading:         "Welcome",
		Message:         "Search for a city to view live weather.",
		ContentTemplate: "index_content",
		SearchQuery:     "",
		NavItems:        navItems("/"),
		Year:            time.Now().Year(),
	}
}

func (s *pageService) Weather(city string) models.PageData {
	return s.renderWeatherPage("Weather", "Current Weather", "Current weather data is shown below.", city, "/weather")
}

func (s *pageService) Forecast(city string) models.PageData {
	return s.renderWeatherPage("Forecast", "Forecast", "A short forecast is shown below.", city, "/forecast")
}

func (s *pageService) About() models.PageData {
	return models.PageData{
		Title:           "About",
		Heading:         "About",
		Message:         "About this weather application.",
		ContentTemplate: "index_content",
		NavItems:        navItems("/about"),
		Year:            time.Now().Year(),
	}
}

func (s *pageService) Error() models.PageData {
	return models.PageData{
		Title:           "Error",
		Heading:         "Something went wrong",
		Message:         "Please try again later.",
		ContentTemplate: "error_content",
		NavItems:        navItems(""),
		Year:            time.Now().Year(),
	}
}

func (s *pageService) renderWeatherPage(title, heading, message, city, activePath string) models.PageData {
	trimmedCity := strings.TrimSpace(city)
	page := models.PageData{
		Title:           title,
		Heading:         heading,
		Message:         message,
		ContentTemplate: "weather_content",
		SearchQuery:     trimmedCity,
		NavItems:        navItems(activePath),
		Year:            time.Now().Year(),
	}
	if trimmedCity == "" {
		page.ErrorMessage = "Enter a city name to get weather details."
		page.HasWeather = false
		return page
	}

	current, err := s.provider.Current(context.Background(), trimmedCity)
	if err != nil {
		page.ErrorMessage = err.Error()
		page.HasWeather = false
		return page
	}

	forecast, err := s.provider.Forecast(context.Background(), trimmedCity)
	if err != nil {
		page.ErrorMessage = err.Error()
		page.HasWeather = false
		return page
	}

	page.Weather = current
	page.Forecast = forecast
	page.HasWeather = true
	return page
}

func navItems(activePath string) []models.NavItem {
	return []models.NavItem{
		{Label: "Home", Path: "/", Active: activePath == "/"},
		{Label: "Weather", Path: "/weather", Active: activePath == "/weather"},
		{Label: "Forecast", Path: "/forecast", Active: activePath == "/forecast"},
		{Label: "About", Path: "/about", Active: activePath == "/about"},
	}
}
