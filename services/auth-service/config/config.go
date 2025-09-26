package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port                string
	DatabaseURL         string
	JWTSecret           string
	JWTExpirationHours  int
	RefreshTokenExpDays int
}

//!TODO: Addreal postgress DB
func Load() *Config {
	return &Config{
		Port:                getEnv("PORT", "8081"),
		DatabaseURL:         getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/auth_db?sslmode=disable"),
		JWTSecret:           getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
		JWTExpirationHours:  getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		RefreshTokenExpDays: getEnvAsInt("REFRESH_TOKEN_EXP_DAYS", 30),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}