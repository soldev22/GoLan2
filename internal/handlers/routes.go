package handlers

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/weather", h.Weather)
	mux.HandleFunc("/forecast", h.Forecast)
	mux.HandleFunc("/about", h.About)
}
