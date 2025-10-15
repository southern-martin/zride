package main


import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/southern-martin/zride/backend/services/ai-matching-service/internal/application/services"
	"github.com/southern-martin/zride/backend/services/ai-matching-service/internal/application/usecases"
	"github.com/southern-martin/zride/backend/services/ai-matching-service/internal/domain"
	redisRepo "github.com/southern-martin/zride/backend/services/ai-matching-service/internal/infrastructure/cache/redis"
	"github.com/southern-martin/zride/backend/services/ai-matching-service/internal/infrastructure/database/postgres"
	"github.com/southern-martin/zride/backend/services/ai-matching-service/internal/interfaces/http/handlers"
	"github.com/southern-martin/zride/backend/services/ai-matching-service/internal/interfaces/http/routes"
	"github.com/southern-martin/zride/backend/shared/logger"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize logger
	logLevel := logger.ParseLogLevel(getEnv("LOG_LEVEL", "info"))
	appLogger := logger.NewStandardLogger(logLevel)

	// Initialize database
	db, err := initDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize Redis
	redisClient := initRedis()
	defer redisClient.Close()

	// Initialize repositories
	matchRequestRepo := postgres.NewMatchRequestRepository(db)
	driverRepo := postgres.NewDriverRepository(db)
	matchResultRepo := postgres.NewMatchResultRepository(db)
	cacheRepo := redisRepo.NewCacheRepository(redisClient)

	// Initialize services
	matchingService := services.NewAIMatchingService(driverRepo, appLogger)
	
	// Initialize external services (mock for now)
	externalSvc := &MockExternalService{}
	notificationSvc := &MockNotificationService{}

	// Initialize use cases
	matchingConfig := domain.DefaultMatchingConfig()
	matchingUseCase := usecases.NewMatchingUseCase(
		matchRequestRepo,
		driverRepo,
		matchResultRepo,
		cacheRepo,
		matchingService,
		notificationSvc,
		externalSvc,
		matchingConfig,
		appLogger,
	)

	// Initialize handlers
	matchingHandler := handlers.NewMatchingHandler(matchingUseCase, appLogger)

	// Setup Gin router
	router := gin.New()
	routes.SetupRoutes(router, matchingHandler)

	// Start server
	port := getEnv("PORT", "8084")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	appLogger.Info(context.Background(), "AI Matching Service started", map[string]interface{}{
		"port": port,
	})

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info(context.Background(), "Shutting down server...", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	appLogger.Info(context.Background(), "Server exited", nil)
}

func initDatabase() (*sql.DB, error) {
	dbURL := getEnv("DATABASE_URL", "postgres://zride_user:zride_password@localhost:5432/zride_ai_matching?sslmode=disable")
	
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func initRedis() *redis.Client {
	redisURL := getEnv("REDIS_URL", "localhost:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: redisPassword,
		DB:       2, // Use DB 2 for matching service
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	}

	return client
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Mock services for development
type MockExternalService struct{}

func (s *MockExternalService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return &domain.User{
		ID:     userID,
		Name:   "Mock User",
		Phone:  "+84123456789",
		Email:  "mock@example.com",
		Rating: 4.5,
	}, nil
}

func (s *MockExternalService) CreateTrip(ctx context.Context, tripData *domain.TripData) (*domain.Trip, error) {
	return &domain.Trip{
		ID:              uuid.New(),
		PassengerID:     tripData.PassengerID,
		DriverID:        tripData.DriverID,
		Status:          "created",
		PickupLocation:  tripData.PickupLocation,
		DropoffLocation: tripData.DropoffLocation,
		CreatedAt:       time.Now(),
	}, nil
}

func (s *MockExternalService) UpdateTripStatus(ctx context.Context, tripID uuid.UUID, status string) error {
	return nil
}

type MockNotificationService struct{}

func (s *MockNotificationService) NotifyDriverMatch(ctx context.Context, driverID uuid.UUID, matchResult *domain.MatchResult) error {
	log.Printf("Mock notification: Driver %s has new match %s", driverID, matchResult.ID)
	return nil
}

func (s *MockNotificationService) NotifyPassengerMatch(ctx context.Context, passengerID uuid.UUID, matchResult *domain.MatchResult) error {
	log.Printf("Mock notification: Passenger %s match accepted %s", passengerID, matchResult.ID)
	return nil
}

func (s *MockNotificationService) NotifyMatchTimeout(ctx context.Context, passengerID uuid.UUID, matchRequest *domain.MatchRequest) error {
	log.Printf("Mock notification: No matches found for passenger %s request %s", passengerID, matchRequest.ID)
	return nil
}