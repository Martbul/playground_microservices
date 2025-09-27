package config

import (
	"os"
)

type Config struct {
	Port           string
	AuthService    string
	ProductService string
	AllowedOrigins []string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		AuthService:    getEnv("AUTH_SERVICE", "localhost:8081"),
		ProductService: getEnv("PRODUCT_SERVICE", "localhost:8082"),
		AllowedOrigins: []string{
			getEnv("CLIENT_URL", "http://localhost:8083"),
			"http://localhost:3000", // For development
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}