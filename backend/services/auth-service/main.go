package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/southern-martin/zride/backend/services/auth-service/internal/application"
	"github.com/southern-martin/zride/backend/services/auth-service/internal/infrastructure"
	"github.com/southern-martin/zride/backend/services/auth-service/internal/interfaces"
	"github.com/southern-martin/zride/backend/shared"
)

func main() {
	// Load configuration
	config, err := infrastructure.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := initDatabase(config)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize services
	userRepo := infrastructure.NewUserRepository(db)
	jwtService := infrastructure.NewJWTService(
		config.JWTSecret,
		config.JWTAccessExpiry,
		config.JWTRefreshExpiry,
	)
	zaloOAuthService := infrastructure.NewZaloOAuthService(
		config.ZaloAppID,
		config.ZaloAppSecret,
	)
	
	authService := application.NewAuthService(userRepo, jwtService, zaloOAuthService)

	// Initialize HTTP router
	router := interfaces.NewRouter(authService, jwtService)
	
	// Setup routes
	r := router.SetupRoutes()

	// Start server
	port := ":" + config.Port
	log.Printf("Auth service starting on port %s", port)
	
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDatabase(config *infrastructure.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.DBSSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, shared.NewDatabaseError("failed to open database connection", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, shared.NewDatabaseError("failed to ping database", err)
	}

	log.Println("Database connection established successfully")
	return db, nil
}