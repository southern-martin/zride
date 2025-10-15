package infrastructure

import (
	"os"
	"strconv"
	"time"

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

	// JWT configuration
	JWTSecret         string
	JWTAccessExpiry   time.Duration
	JWTRefreshExpiry  time.Duration

	// Zalo configuration
	ZaloAppID     string
	ZaloAppSecret string
}

func LoadConfig() (*Config, error) {
	// Load .env file if it exists (for development)
	_ = godotenv.Load()

	config := &Config{
		Port: getEnv("PORT", "8081"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "zride_user"),
		DBPassword: getEnv("DB_PASSWORD", "zride_password"),
		DBName:     getEnv("DB_NAME", "zride"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		JWTSecret:        getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		JWTAccessExpiry:  parseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m")),
		JWTRefreshExpiry: parseDuration(getEnv("JWT_REFRESH_EXPIRY", "7d")),

		ZaloAppID:     getEnv("ZALO_APP_ID", ""),
		ZaloAppSecret: getEnv("ZALO_APP_SECRET", ""),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 15 * time.Minute // default fallback
	}
	return d
}

func parseInt(s string, defaultValue int) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return defaultValue
}