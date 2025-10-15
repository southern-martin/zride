package domain


import (
	"time"
	"github.com/google/uuid"
)

// Location represents a geographical coordinate
type Location struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
	Address   string  `json:"address"`
}

// MatchRequest represents a request for driver-passenger matching
type MatchRequest struct {
	ID              uuid.UUID        `json:"id"`
	PassengerID     uuid.UUID        `json:"passenger_id" validate:"required"`
	PickupLocation  Location         `json:"pickup_location" validate:"required"`
	DropoffLocation Location         `json:"dropoff_location" validate:"required"`
	RequestTime     time.Time        `json:"request_time"`
	MaxWaitTime     time.Duration    `json:"max_wait_time"` // Maximum time to wait for a driver
	PreferredCarType string          `json:"preferred_car_type"`
	MaxDistance     float64          `json:"max_distance"` // Maximum distance to search for drivers (km)
	PriceRange      PriceRange       `json:"price_range"`
	Status          MatchStatus      `json:"status"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

// PriceRange represents acceptable price range for a trip
type PriceRange struct {
	MinPrice float64 `json:"min_price"`
	MaxPrice float64 `json:"max_price"`
}

// MatchStatus represents the status of a match request
type MatchStatus string

const (
	MatchStatusPending   MatchStatus = "pending"
	MatchStatusMatched   MatchStatus = "matched"
	MatchStatusExpired   MatchStatus = "expired"
	MatchStatusCancelled MatchStatus = "cancelled"
)

// Driver represents a driver available for matching
type Driver struct {
	ID               uuid.UUID      `json:"id"`
	UserID           uuid.UUID      `json:"user_id"`
	CurrentLocation  Location       `json:"current_location"`
	IsAvailable      bool           `json:"is_available"`
	CarType          string         `json:"car_type"`
	Rating           float64        `json:"rating"`
	CompletedTrips   int            `json:"completed_trips"`
	LastActiveTime   time.Time      `json:"last_active_time"`
	MaxDistance      float64        `json:"max_distance"` // Maximum distance willing to travel (km)
	PreferredAreas   []Location     `json:"preferred_areas"`
}

// MatchResult represents the result of a matching algorithm
type MatchResult struct {
	ID                uuid.UUID        `json:"id"`
	MatchRequestID    uuid.UUID        `json:"match_request_id"`
	DriverID          uuid.UUID        `json:"driver_id"`
	Score             float64          `json:"score"` // Matching score (0-1)
	EstimatedDistance float64          `json:"estimated_distance"` // Distance from driver to pickup (km)
	EstimatedTime     time.Duration    `json:"estimated_time"` // Time to reach pickup
	EstimatedPrice    float64          `json:"estimated_price"` // Estimated trip price in VND
	MatchTime         time.Time        `json:"match_time"`
	Status            MatchResultStatus `json:"status"`
	CreatedAt         time.Time        `json:"created_at"`
}

// MatchResultStatus represents the status of a match result
type MatchResultStatus string

const (
	MatchResultStatusPending  MatchResultStatus = "pending"
	MatchResultStatusAccepted MatchResultStatus = "accepted"
	MatchResultStatusRejected MatchResultStatus = "rejected"
	MatchResultStatusExpired  MatchResultStatus = "expired"
)

// MatchingCriteria represents criteria for matching algorithms
type MatchingCriteria struct {
	DistanceWeight     float64 `json:"distance_weight"`     // Weight for distance factor (0-1)
	RatingWeight       float64 `json:"rating_weight"`       // Weight for driver rating (0-1)
	TimeWeight         float64 `json:"time_weight"`         // Weight for estimated time (0-1)
	PriceWeight        float64 `json:"price_weight"`        // Weight for price factor (0-1)
	ExperienceWeight   float64 `json:"experience_weight"`   // Weight for driver experience (0-1)
}

// DefaultMatchingCriteria returns default matching criteria
func DefaultMatchingCriteria() MatchingCriteria {
	return MatchingCriteria{
		DistanceWeight:   0.4,
		RatingWeight:     0.2,
		TimeWeight:       0.2,
		PriceWeight:      0.1,
		ExperienceWeight: 0.1,
	}
}

// MatchingAlgorithm represents different matching algorithms
type MatchingAlgorithm string

const (
	AlgorithmNearest    MatchingAlgorithm = "nearest"    // Nearest driver algorithm
	AlgorithmWeighted   MatchingAlgorithm = "weighted"   // Weighted score algorithm
	AlgorithmML         MatchingAlgorithm = "ml"         // Machine learning algorithm
	AlgorithmHybrid     MatchingAlgorithm = "hybrid"     // Hybrid algorithm
)

// MatchingConfig represents configuration for matching service
type MatchingConfig struct {
	Algorithm           MatchingAlgorithm `json:"algorithm"`
	MaxDrivers          int               `json:"max_drivers"`           // Maximum drivers to consider
	MaxSearchRadius     float64           `json:"max_search_radius"`     // Maximum search radius in km
	MinDriverRating     float64           `json:"min_driver_rating"`     // Minimum driver rating
	Criteria            MatchingCriteria  `json:"criteria"`
	RealTimeEnabled     bool              `json:"real_time_enabled"`     // Enable real-time matching
	BatchSize           int               `json:"batch_size"`            // Batch size for processing
	TimeoutDuration     time.Duration     `json:"timeout_duration"`      // Matching timeout
}

// DefaultMatchingConfig returns default matching configuration
func DefaultMatchingConfig() MatchingConfig {
	return MatchingConfig{
		Algorithm:       AlgorithmWeighted,
		MaxDrivers:      10,
		MaxSearchRadius: 15.0, // 15km radius
		MinDriverRating: 3.0,
		Criteria:        DefaultMatchingCriteria(),
		RealTimeEnabled: true,
		BatchSize:       50,
		TimeoutDuration: 30 * time.Second,
	}
}