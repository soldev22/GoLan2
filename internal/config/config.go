package config

import (
	"os"
	"strconv"
	"time"
)

// Config keeps runtime settings loaded from environment variables.
type Config struct {
	Addr            string
	ShutdownTimeout time.Duration
}

func Load() Config {
	addr := getEnv("WEATHER_APP_ADDR", "localhost:8080")
	shutdownSeconds := getIntEnv("WEATHER_APP_SHUTDOWN_TIMEOUT_SEC", 10)

	return Config{
		Addr:            addr,
		ShutdownTimeout: time.Duration(shutdownSeconds) * time.Second,
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getIntEnv(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}
