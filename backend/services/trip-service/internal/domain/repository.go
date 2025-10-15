package domain

import (
	"context"
	"time"
	"github.com/google/uuid"
)

// TripRepository defines the interface for trip data access
type TripRepository interface {
	// CreateTrip creates a new trip
	CreateTrip(ctx context.Context, trip *Trip) error
	
	// GetTripByID retrieves a trip by its ID
	GetTripByID(ctx context.Context, id uuid.UUID) (*Trip, error)
	
	// UpdateTrip updates an existing trip
	UpdateTrip(ctx context.Context, trip *Trip) error
	
	// DeleteTrip soft deletes a trip
	DeleteTrip(ctx context.Context, id uuid.UUID) error
	
	// SearchTrips searches for trips based on criteria
	SearchTrips(ctx context.Context, criteria TripSearchCriteria) ([]*Trip, int64, error)
	
	// GetActiveTrips retrieves all active trips
	GetActiveTrips(ctx context.Context) ([]*Trip, error)
	
	// GetTripsByPassenger retrieves trips for a specific passenger
	GetTripsByPassenger(ctx context.Context, passengerID uuid.UUID, limit, offset int) ([]*Trip, int64, error)
	
	// GetTripsByDriver retrieves trips for a specific driver
	GetTripsByDriver(ctx context.Context, driverID uuid.UUID, limit, offset int) ([]*Trip, int64, error)
	
	// GetTripsByStatus retrieves trips by status
	GetTripsByStatus(ctx context.Context, status TripStatus, limit, offset int) ([]*Trip, int64, error)
	
	// GetNearbyTrips finds trips near a specific location
	GetNearbyTrips(ctx context.Context, location Location, radiusKM float64, status *TripStatus) ([]*Trip, error)
	
	// UpdateTripStatus updates only the status of a trip
	UpdateTripStatus(ctx context.Context, tripID uuid.UUID, status TripStatus) error
	
	// GetTripStatistics gets trip statistics for analytics
	GetTripStatistics(ctx context.Context, from, to *time.Time) (*TripStatistics, error)
}

// TripStatistics represents trip statistics
type TripStatistics struct {
	TotalTrips      int64   `json:"total_trips"`
	CompletedTrips  int64   `json:"completed_trips"`
	CancelledTrips  int64   `json:"cancelled_trips"`
	ActiveTrips     int64   `json:"active_trips"`
	TotalRevenue    float64 `json:"total_revenue"`
	AverageDistance float64 `json:"average_distance"`
	AverageFare     float64 `json:"average_fare"`
}