package handlers

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/healthz", Health)
	mux.HandleFunc("/weather", h.Weather)
	mux.HandleFunc("/forecast", h.Forecast)
	mux.HandleFunc("/about", h.About)
}

func Health(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
