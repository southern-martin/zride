package domain

import (
	"context"
	"time"
	"github.com/google/uuid"
)

// MatchingService defines the interface for matching algorithms
type MatchingService interface {
	FindMatches(ctx context.Context, request *MatchRequest, config MatchingConfig) ([]*MatchResult, error)
	ScoreMatch(ctx context.Context, request *MatchRequest, driver *Driver, criteria MatchingCriteria) (float64, error)
	CalculateDistance(from, to Location) float64
	EstimateTime(from, to Location) (time.Duration, error)
	EstimatePrice(from, to Location, carType string) (float64, error)
}

// NotificationService defines the interface for sending notifications
type NotificationService interface {
	NotifyDriverMatch(ctx context.Context, driverID uuid.UUID, matchResult *MatchResult) error
	NotifyPassengerMatch(ctx context.Context, passengerID uuid.UUID, matchResult *MatchResult) error
	NotifyMatchTimeout(ctx context.Context, passengerID uuid.UUID, matchRequest *MatchRequest) error
}

// ExternalService defines the interface for external service calls
type ExternalService interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
	CreateTrip(ctx context.Context, tripData *TripData) (*Trip, error)
	UpdateTripStatus(ctx context.Context, tripID uuid.UUID, status string) error
}

// User represents a user from external service
type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Email    string    `json:"email"`
	Rating   float64   `json:"rating"`
}

// TripData represents trip creation data
type TripData struct {
	PassengerID      uuid.UUID `json:"passenger_id"`
	DriverID         uuid.UUID `json:"driver_id"`
	PickupLocation   Location  `json:"pickup_location"`
	DropoffLocation  Location  `json:"dropoff_location"`
	EstimatedPrice   float64   `json:"estimated_price"`
	EstimatedTime    int       `json:"estimated_time"` // in minutes
	EstimatedDistance float64  `json:"estimated_distance"` // in km
}

// Trip represents a trip from external service
type Trip struct {
	ID              uuid.UUID `json:"id"`
	PassengerID     uuid.UUID `json:"passenger_id"`
	DriverID        uuid.UUID `json:"driver_id"`
	Status          string    `json:"status"`
	PickupLocation  Location  `json:"pickup_location"`
	DropoffLocation Location  `json:"dropoff_location"`
	CreatedAt       time.Time `json:"created_at"`
}