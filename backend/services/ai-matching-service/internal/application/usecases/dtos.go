package usecases

import (
	"time"
	"github.com/google/uuid"
	"github.com/southern-martin/zride/backend/services/ai-matching-service/internal/domain"
)

// CreateMatchRequestDTO represents the request to create a match
type CreateMatchRequestDTO struct {
	PassengerID      uuid.UUID           `json:"passenger_id" validate:"required"`
	PickupLocation   domain.Location     `json:"pickup_location" validate:"required"`
	DropoffLocation  domain.Location     `json:"dropoff_location" validate:"required"`
	MaxWaitTime      time.Duration       `json:"max_wait_time"`
	PreferredCarType string              `json:"preferred_car_type"`
	MaxDistance      float64             `json:"max_distance"`
	PriceRange       domain.PriceRange   `json:"price_range"`
}

// MatchRequestResponse represents the response for match request
type MatchRequestResponse struct {
	ID               uuid.UUID       `json:"id"`
	PassengerID      uuid.UUID       `json:"passenger_id"`
	PickupLocation   domain.Location `json:"pickup_location"`
	DropoffLocation  domain.Location `json:"dropoff_location"`
	Status           string          `json:"status"`
	CreatedAt        time.Time       `json:"created_at"`
	EstimatedMatches int             `json:"estimated_matches"`
}

// AcceptMatchResponse represents the response for accepting a match
type AcceptMatchResponse struct {
	MatchResultID  uuid.UUID     `json:"match_result_id"`
	TripID         uuid.UUID     `json:"trip_id"`
	PassengerID    uuid.UUID     `json:"passenger_id"`
	DriverID       uuid.UUID     `json:"driver_id"`
	EstimatedTime  time.Duration `json:"estimated_time"`
	EstimatedPrice float64       `json:"estimated_price"`
	Status         string        `json:"status"`
}

// MatchResultResponse represents the response for match results
type MatchResultResponse struct {
	ID                uuid.UUID       `json:"id"`
	MatchRequestID    uuid.UUID       `json:"match_request_id"`
	DriverID          uuid.UUID       `json:"driver_id"`
	PassengerID       uuid.UUID       `json:"passenger_id"`
	PickupLocation    domain.Location `json:"pickup_location"`
	DropoffLocation   domain.Location `json:"dropoff_location"`
	Score             float64         `json:"score"`
	EstimatedDistance float64         `json:"estimated_distance"`
	EstimatedTime     time.Duration   `json:"estimated_time"`
	EstimatedPrice    float64         `json:"estimated_price"`
	Status            string          `json:"status"`
}

// UpdateDriverLocationDTO represents the request to update driver location
type UpdateDriverLocationDTO struct {
	DriverID uuid.UUID       `json:"driver_id" validate:"required"`
	Location domain.Location `json:"location" validate:"required"`
}

// UpdateDriverAvailabilityDTO represents the request to update driver availability
type UpdateDriverAvailabilityDTO struct {
	DriverID  uuid.UUID `json:"driver_id" validate:"required"`
	Available bool      `json:"available"`
}