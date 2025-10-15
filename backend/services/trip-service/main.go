package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/southern-martin/zride/backend/services/trip-service/internal/application"
	"github.com/southern-martin/zride/backend/services/trip-service/internal/infrastructure"
	"github.com/southern-martin/zride/backend/services/trip-service/internal/interfaces"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Load configuration
	config := infrastructure.LoadConfig()
	if err := config.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Connect to database
	db, err := connectDatabase(config.GetDatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	tripRepo := infrastructure.NewTripRepository(db)

	// Initialize services
	tripService := application.NewTripService(tripRepo)

	// Initialize handlers
	tripHandler := interfaces.NewTripHandler(tripService)

	// Setup routes
	router := interfaces.SetupRoutes(tripHandler)

	// Start server
	port := config.Port
	log.Printf("Starting trip service on port %s", port)
	log.Printf("Trip service is running at http://localhost:%s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// connectDatabase establishes a connection to PostgreSQL
func connectDatabase(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Successfully connected to database")
	return db, nil
}

// setupDatabase creates tables if they don't exist
func setupDatabase(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS trips (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		passenger_id UUID NOT NULL,
		driver_id UUID,
		vehicle_id UUID,
		pickup_location JSONB NOT NULL,
		dropoff_location JSONB NOT NULL,
		scheduled_time TIMESTAMP WITH TIME ZONE,
		requested_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		accepted_at TIMESTAMP WITH TIME ZONE,
		pickup_time TIMESTAMP WITH TIME ZONE,
		dropoff_time TIMESTAMP WITH TIME ZONE,
		status VARCHAR(20) NOT NULL DEFAULT 'requested' CHECK (status IN ('requested', 'accepted', 'in_progress', 'completed', 'cancelled')),
		route_info JSONB,
		pricing_info JSONB,
		passenger_notes TEXT,
		driver_notes TEXT,
		cancellation_info TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		deleted_at TIMESTAMP WITH TIME ZONE
	);

	-- Create indexes for better performance
	CREATE INDEX IF NOT EXISTS idx_trips_passenger_id ON trips(passenger_id);
	CREATE INDEX IF NOT EXISTS idx_trips_driver_id ON trips(driver_id);
	CREATE INDEX IF NOT EXISTS idx_trips_status ON trips(status);
	CREATE INDEX IF NOT EXISTS idx_trips_created_at ON trips(created_at);
	CREATE INDEX IF NOT EXISTS idx_trips_requested_at ON trips(requested_at);
	CREATE INDEX IF NOT EXISTS idx_trips_deleted_at ON trips(deleted_at);

	-- Spatial index for location-based queries (requires PostGIS)
	-- CREATE INDEX IF NOT EXISTS idx_trips_pickup_location ON trips USING GIST ((pickup_location->>'lat')::float, (pickup_location->>'lng')::float);
	`

	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("failed to create database schema: %w", err)
	}

	log.Println("Database schema initialized successfully")
	return nil
}