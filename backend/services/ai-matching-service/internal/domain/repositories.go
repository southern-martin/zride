package domain

import (
	"context"
	"github.com/google/uuid"
)

// MatchRequestRepository defines the repository interface for match requests
type MatchRequestRepository interface {
	Create(ctx context.Context, request *MatchRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*MatchRequest, error)
	GetByPassengerID(ctx context.Context, passengerID uuid.UUID) ([]*MatchRequest, error)
	Update(ctx context.Context, request *MatchRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetPendingRequests(ctx context.Context, limit int) ([]*MatchRequest, error)
}

// DriverRepository defines the repository interface for drivers
type DriverRepository interface {
	GetAvailableDriversInRadius(ctx context.Context, location Location, radius float64) ([]*Driver, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Driver, error)
	UpdateLocation(ctx context.Context, driverID uuid.UUID, location Location) error
	UpdateAvailability(ctx context.Context, driverID uuid.UUID, available bool) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Driver, error)
}

// MatchResultRepository defines the repository interface for match results
type MatchResultRepository interface {
	Create(ctx context.Context, result *MatchResult) error
	GetByID(ctx context.Context, id uuid.UUID) (*MatchResult, error)
	GetByMatchRequestID(ctx context.Context, matchRequestID uuid.UUID) ([]*MatchResult, error)
	Update(ctx context.Context, result *MatchResult) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// CacheRepository defines the repository interface for caching
type CacheRepository interface {
	SetDriverLocation(ctx context.Context, driverID uuid.UUID, location Location) error
	GetDriverLocation(ctx context.Context, driverID uuid.UUID) (*Location, error)
	SetMatchResult(ctx context.Context, key string, result *MatchResult) error
	GetMatchResult(ctx context.Context, key string) (*MatchResult, error)
	DeleteMatchResult(ctx context.Context, key string) error
}