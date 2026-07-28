package main

import (
	"context"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"weather-app/internal/config"
	"weather-app/internal/handlers"
	"weather-app/internal/middleware"
	"weather-app/internal/services"
)

func main() {
	cfg := config.Load()
	logger := log.New(os.Stdout, "", log.LstdFlags)

	tmpl, err := template.ParseFiles(
		"templates/layouts/base.html",
		"templates/partials/nav.html",
		"templates/index.html",
		"templates/weather.html",
		"templates/error.html",
	)
	if err != nil {
		logger.Fatalf("parse templates: %v", err)
	}

	pageService := services.NewPageService(nil)
	h := handlers.New(pageService, tmpl, logger)

	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux, h)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	rootHandler := middleware.Recovery(middleware.Logging(mux, logger), logger)
	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           rootHandler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Printf("server listening on http://%s", cfg.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("listen and serve: %v", err)
		}
	}()

	<-shutdownCtx.Done()
	logger.Println("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Printf("graceful shutdown failed: %v", err)
		return
	}

	logger.Println("server stopped")
}
