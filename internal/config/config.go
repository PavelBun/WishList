// Package config provides configuration loading from environment variables.
package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration parameters for the application.
type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	AppPort        string
	JWTSecret      string
	JWTExpiryHours int
}

// Load reads configuration from .env file and environment variables.
func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET environment variable is required")
	}

	expiryHours := 24
	if val := os.Getenv("JWT_EXPIRY_HOURS"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 {
			expiryHours = parsed
		}
	}

	return &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "wishlist"),
		AppPort:        getEnv("APP_PORT", "8080"),
		JWTSecret:      jwtSecret,
		JWTExpiryHours: expiryHours,
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
