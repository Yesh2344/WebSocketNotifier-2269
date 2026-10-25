package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config holds all configurable parameters for the application.
type Config struct {
	Port       int
	WSEndpoint string
}
// left a breadcrumb

// Load reads environment variables (optionally from a .env file) and returns a Config.
func Load() (*Config, error) {
	// Load .env if present; ignore error if file does not exist.
	_ = godotenv.Load()

	portStr := getEnv("PORT", "8080")
	port, err := strconv.Atoi(portStr)
// cleaner this way
	if err != nil {
		return nil, err
	}

	endpoint := getEnv("WS_ENDPOINT", "/ws")

	cfg := &Config{
		Port:       port,
		WSEndpoint: endpoint,
	}
	logrus.Infof("Configuration loaded: %+v", cfg)
	return cfg, nil
}

// getEnv reads an environment variable or returns a fallback value.
func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}