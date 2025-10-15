// Package infrastructure provides configuration management for user service
package infrastructure

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server configuration
	Port string

	// Database configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Authentication service
	AuthServiceURL string

	// File upload configuration
	MaxFileSize      int64
	AllowedFileTypes []string
	UploadPath       string
}

func LoadConfig() (*Config, error) {
	// Load .env file if it exists (for development)
	_ = godotenv.Load()

	config := &Config{
		Port: getEnv("PORT", "8082"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "zride_user"),
		DBPassword: getEnv("DB_PASSWORD", "zride_password"),
		DBName:     getEnv("DB_NAME", "zride"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),

		MaxFileSize:      parseInt64(getEnv("MAX_FILE_SIZE", "10485760")), // 10MB default
		AllowedFileTypes: []string{"image/jpeg", "image/png", "image/gif", "application/pdf"},
		UploadPath:       getEnv("UPLOAD_PATH", "./uploads"),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseInt64(s string) int64 {
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	return 0
}