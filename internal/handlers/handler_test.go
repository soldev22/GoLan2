package handlers

import (
	"bytes"
	"html/template"
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
			Label: "Home",
			Path:  "/",
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
