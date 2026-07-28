package models

// Weather represents a current weather snapshot.
type Weather struct {
	City         string
	TemperatureC int
	Condition    string
	FeelsLikeC   int
	HumidityPct  int
	WindKph      int
	Description  string
}

// ForecastItem represents a single forecast entry.
type ForecastItem struct {
	Day       string // ISO date (YYYY-MM-DD), for machine-readable attributes
	DayLabel  string // Formatted label for display (e.g. "Mon, Jul 28")
	HighC     int
	LowC      int
	Condition string
}
