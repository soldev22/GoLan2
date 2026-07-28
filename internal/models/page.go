package models

// NavItem describes a single navigation link.
type NavItem struct {
	Label  string
	Path   string
	Active bool
}

// PageData carries render data for HTML templates.
type PageData struct {
	Title           string
	Heading         string
	Message         string
	ContentTemplate string
	NavItems        []NavItem
	Year            int
	Weather         Weather
	Forecast        []ForecastItem
	HasWeather      bool
	ErrorMessage    string
}
