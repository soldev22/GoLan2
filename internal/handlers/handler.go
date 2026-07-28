package handlers

import (
	"html/template"
	"log"
	"net/http"

	"weather-app/internal/services"
)

type Handler struct {
	pages     services.PageService
	templates *template.Template
	logger    *log.Logger
}

func New(pages services.PageService, templates *template.Template, logger *log.Logger) *Handler {
	return &Handler{
		pages:     pages,
		templates: templates,
		logger:    logger,
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	h.render(w, http.StatusOK, h.pages.Home())
}

func (h *Handler) Weather(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	city := r.URL.Query().Get("city")
	h.render(w, http.StatusOK, h.pages.Weather(city))
}

func (h *Handler) Forecast(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	city := r.URL.Query().Get("city")
	h.render(w, http.StatusOK, h.pages.Forecast(city))
}

func (h *Handler) About(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	h.render(w, http.StatusOK, h.pages.About())
}

func (h *Handler) render(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	if err := h.templates.ExecuteTemplate(w, "base", data); err != nil {
		h.logger.Printf("render template: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) RenderError(w http.ResponseWriter, status int) {
	h.render(w, status, h.pages.Error())
}

func allowMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		w.Header().Set("Allow", method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}
