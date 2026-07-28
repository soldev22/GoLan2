package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"weather-app/internal/models"
)

func TestTemplatesParseAndRender(t *testing.T) {
	tmpl, err := template.ParseFiles(
		"../../templates/layouts/base.html",
		"../../templates/partials/nav.html",
		"../../templates/index.html",
		"../../templates/weather.html",
		"../../templates/error.html",
	)
	if err != nil {
		t.Fatalf("parse templates: %v", err)
	}

	var buf bytes.Buffer
	data := models.PageData{
		Title:           "Home",
		Heading:         "Welcome",
		Message:         "Weather App starter is running.",
		ContentTemplate: "index_content",
		NavItems: []models.NavItem{{
			Label:  "Home",
			Path:   "/",
			Active: true,
		}},
		Year: 2026,
	}

	if err := tmpl.ExecuteTemplate(&buf, "base", data); err != nil {
		t.Fatalf("execute template: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Weather App") {
		t.Fatalf("rendered output missing expected content: %s", output)
	}
}

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	Health(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != "ok" {
		t.Fatalf("expected body ok, got %q", got)
	}
}
