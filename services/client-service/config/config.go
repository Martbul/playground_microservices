package config

import (
	"os"
)

type Config struct {
	Port           string
	APIGatewayURL  string
	SessionSecret  string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8083"),
		APIGatewayURL:  getEnv("API_GATEWAY_URL", "http://localhost:8080"),
		SessionSecret:  getEnv("SESSION_SECRET", "your-super-secret-session-key-change-in-production"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}