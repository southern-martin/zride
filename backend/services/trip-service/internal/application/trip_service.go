package application

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/southern-martin/zride/backend/services/trip-service/internal/domain"
	"github.com/southern-martin/zride/backend/shared"
)

// TripService handles trip business logic
type TripService struct {
	tripRepo domain.TripRepository
}

// NewTripService creates a new trip service
func NewTripService(tripRepo domain.TripRepository) *TripService {
	return &TripService{
		tripRepo: tripRepo,
	}
}

// CreateTripRequest represents a request to create a trip
type CreateTripRequest struct {
	PassengerID     uuid.UUID               `json:"passenger_id" binding:"required"`
	PickupLocation  domain.Location         `json:"pickup_location" binding:"required"`
	DropoffLocation domain.Location         `json:"dropoff_location" binding:"required"`
	ScheduledTime   *time.Time              `json:"scheduled_time,omitempty"`
	Notes           string                  `json:"notes,omitempty"`
}

// UpdateTripRequest represents a request to update a trip
type UpdateTripRequest struct {
	PickupLocation   *domain.Location `json:"pickup_location,omitempty"`
	DropoffLocation  *domain.Location `json:"dropoff_location,omitempty"`
	ScheduledTime    *time.Time       `json:"scheduled_time,omitempty"`
	PassengerNotes   string           `json:"passenger_notes,omitempty"`
}

// AcceptTripRequest represents a request to accept a trip
type AcceptTripRequest struct {
	DriverID  uuid.UUID `json:"driver_id" binding:"required"`
	VehicleID uuid.UUID `json:"vehicle_id" binding:"required"`
}

// TripSearchRequest represents a trip search request
type TripSearchRequest struct {
	PassengerID     *uuid.UUID            `json:"passenger_id,omitempty"`
	DriverID        *uuid.UUID            `json:"driver_id,omitempty"`
	Status          *domain.TripStatus    `json:"status,omitempty"`
	FromTime        *time.Time            `json:"from_time,omitempty"`
	ToTime          *time.Time            `json:"to_time,omitempty"`
	PickupRadius    *float64              `json:"pickup_radius,omitempty"`
	DropoffRadius   *float64              `json:"dropoff_radius,omitempty"`
	PickupLocation  *domain.Location      `json:"pickup_location,omitempty"`
	DropoffLocation *domain.Location      `json:"dropoff_location,omitempty"`
	MaxPrice        *float64              `json:"max_price,omitempty"`
	MinPrice        *float64              `json:"min_price,omitempty"`
	Page            int                   `json:"page,omitempty"`
	PageSize        int                   `json:"page_size,omitempty"`
	SortBy          string                `json:"sort_by,omitempty"`
	SortOrder       string                `json:"sort_order,omitempty"`
}

// TripResponse represents a trip response
type TripResponse struct {
	*domain.Trip
	EstimatedFare *float64 `json:"estimated_fare,omitempty"`
}

// TripListResponse represents a paginated list of trips
type TripListResponse struct {
	Trips      []*TripResponse `json:"trips"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// CreateTrip creates a new trip
func (s *TripService) CreateTrip(ctx context.Context, req CreateTripRequest) (*TripResponse, error) {
	// Create new trip
	trip := domain.NewTrip(
		req.PassengerID,
		req.PickupLocation,
		req.DropoffLocation,
		req.ScheduledTime,
		req.Notes,
	)

	// Validate trip
	if err := trip.Validate(); err != nil {
		return nil, shared.NewError(shared.ErrCodeValidation, "Invalid trip data", err.Error())
	}

	// Calculate route and pricing
	if err := s.calculateRoute(trip); err != nil {
		return nil, shared.NewError(shared.ErrCodeInternal, "Failed to calculate route", err.Error())
	}

	if err := s.calculatePricing(trip); err != nil {
		return nil, shared.NewError(shared.ErrCodeInternal, "Failed to calculate pricing", err.Error())
	}

	// Save trip
	if err := s.tripRepo.CreateTrip(ctx, trip); err != nil {
		return nil, shared.NewError(shared.ErrCodeInternal, "Failed to create trip", err.Error())
	}

	response := &TripResponse{
		Trip: trip,
	}
	
	if trip.PricingInfo != nil {
		response.EstimatedFare = &trip.PricingInfo.TotalFare
	}

	return response, nil
}

// GetTrip retrieves a trip by ID
func (s *TripService) GetTrip(ctx context.Context, id uuid.UUID) (*TripResponse, error) {
	trip, err := s.tripRepo.GetTripByID(ctx, id)
	if err != nil {
		return nil, shared.NewError(shared.ErrCodeNotFound, "Trip not found", err.Error())
	}

	response := &TripResponse{
		Trip: trip,
	}
	
	if trip.PricingInfo != nil {
		response.EstimatedFare = &trip.PricingInfo.TotalFare
	}

	return response, nil
}

// UpdateTrip updates a trip
func (s *TripService) UpdateTrip(ctx context.Context, id uuid.UUID, req UpdateTripRequest) (*TripResponse, error) {
	trip, err := s.tripRepo.GetTripByID(ctx, id)
	if err != nil {
		return nil, shared.NewError(shared.ErrCodeNotFound, "Trip not found", err.Error())
	}

	if !trip.CanBeModified() {
		return nil, shared.NewError(shared.ErrCodeValidation, "Trip cannot be modified", "Trip status does not allow modifications")
	}

	// Update trip fields
	if req.PickupLocation != nil {
		trip.PickupLocation = *req.PickupLocation
	}
	if req.DropoffLocation != nil {
		trip.DropoffLocation = *req.DropoffLocation
	}
	if req.ScheduledTime != nil {
		trip.ScheduledTime = req.ScheduledTime
	}
	if req.PassengerNotes != "" {
		trip.PassengerNotes = req.PassengerNotes
	}

	trip.UpdatedAt = time.Now()

	// Recalculate route and pricing if locations changed
	if req.PickupLocation != nil || req.DropoffLocation != nil {
		if err := s.calculateRoute(trip); err != nil {
			return nil, shared.NewError(shared.ErrCodeInternal, "Failed to recalculate route", err.Error())
		}
		if err := s.calculatePricing(trip); err != nil {
			return nil, shared.NewError(shared.ErrCodeInternal, "Failed to recalculate pricing", err.Error())
		}
	}

	// Validate updated trip
	if err := trip.Validate(); err != nil {
		return nil, shared.NewError(shared.ErrCodeValidation, "Invalid trip data", err.Error())
	}

	// Save updated trip
	if err := s.tripRepo.UpdateTrip(ctx, trip); err != nil {
		return nil, shared.NewError(shared.ErrCodeInternal, "Failed to update trip", err.Error())
	}

	response := &TripResponse{
		Trip: trip,
	}
	
	if trip.PricingInfo != nil {
		response.EstimatedFare = &trip.PricingInfo.TotalFare
	}

	return response, nil
}

// AcceptTrip accepts a trip by a driver
func (s *TripService) AcceptTrip(ctx context.Context, tripID uuid.UUID, req AcceptTripRequest) (*TripResponse, error) {
	trip, err := s.tripRepo.GetTripByID(ctx, tripID)
	if err != nil {
		return nil, shared.NewError(shared.ErrCodeNotFound, "Trip not found", err.Error())
	}

	if err := trip.AcceptTrip(req.DriverID, req.VehicleID); err != nil {
		return nil, shared.NewError(shared.ErrCodeValidation, "Cannot accept trip", err.Error())
	}

	if err := s.tripRepo.UpdateTrip(ctx, trip); err != nil {
		return nil, shared.NewError(shared.ErrCodeInternal, "Failed to accept trip", err.Error())
	}

	return &TripResponse{Trip: trip}, nil
}

// StartTrip starts a trip
func (s *TripService) StartTrip(ctx context.Context, tripID uuid.UUID) (*TripResponse, error) {
	trip, err := s.tripRepo.GetTripByID(ctx, tripID)
	if err != nil {
		return nil, shared.NewError(shared.ErrCodeNotFound, "Trip not found", err.Error())
	}

	if err := trip.StartTrip(); err != nil {
		return nil, shared.NewError(shared.ErrCodeValidation, "Cannot start trip", err.Error())
	}

	if err := s.tripRepo.UpdateTrip(ctx, trip); err != nil {
		return nil, shared.NewError(shared.ErrCodeInternal, "Failed to start trip", err.Error())
	}

	return &TripResponse{Trip: trip}, nil
}

// CompleteTrip completes a trip
func (s *TripService) CompleteTrip(ctx context.Context, tripID uuid.UUID) (*TripResponse, error) {
	trip, err := s.tripRepo.GetTripByID(ctx, tripID)
	if err != nil {
		return nil, shared.NewError(shared.ErrCodeNotFound, "Trip not found", err.Error())
	}

	if err := trip.CompleteTrip(); err != nil {
		return nil, shared.NewError(shared.ErrCodeValidation, "Cannot complete trip", err.Error())
	}

	if err := s.tripRepo.UpdateTrip(ctx, trip); err != nil {
		return nil, shared.NewError(shared.ErrCodeInternal, "Failed to complete trip", err.Error())
	}

	return &TripResponse{Trip: trip}, nil
}

// CancelTrip cancels a trip
func (s *TripService) CancelTrip(ctx context.Context, tripID uuid.UUID, reason string) (*TripResponse, error) {
	trip, err := s.tripRepo.GetTripByID(ctx, tripID)
	if err != nil {
		return nil, shared.NewError(shared.ErrCodeNotFound, "Trip not found", err.Error())
	}

	if err := trip.CancelTrip(reason); err != nil {
		return nil, shared.NewError(shared.ErrCodeValidation, "Cannot cancel trip", err.Error())
	}

	if err := s.tripRepo.UpdateTrip(ctx, trip); err != nil {
		return nil, shared.NewError(shared.ErrCodeInternal, "Failed to cancel trip", err.Error())
	}

	return &TripResponse{Trip: trip}, nil
}

// SearchTrips searches for trips
func (s *TripService) SearchTrips(ctx context.Context, req TripSearchRequest) (*TripListResponse, error) {
	// Set default pagination
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	// Convert to domain criteria
	criteria := domain.TripSearchCriteria{
		PassengerID:     req.PassengerID,
		DriverID:        req.DriverID,
		Status:          req.Status,
		FromTime:        req.FromTime,
		ToTime:          req.ToTime,
		PickupRadius:    req.PickupRadius,
		DropoffRadius:   req.DropoffRadius,
		PickupLocation:  req.PickupLocation,
		DropoffLocation: req.DropoffLocation,
		MaxPrice:        req.MaxPrice,
		MinPrice:        req.MinPrice,
		Limit:           req.PageSize,
		Offset:          (req.Page - 1) * req.PageSize,
		SortBy:          req.SortBy,
		SortOrder:       req.SortOrder,
	}

	trips, total, err := s.tripRepo.SearchTrips(ctx, criteria)
	if err != nil {
		return nil, shared.NewError(shared.ErrCodeInternal, "Failed to search trips", err.Error())
	}

	// Convert to response format
	tripResponses := make([]*TripResponse, len(trips))
	for i, trip := range trips {
		tripResponse := &TripResponse{Trip: trip}
		if trip.PricingInfo != nil {
			tripResponse.EstimatedFare = &trip.PricingInfo.TotalFare
		}
		tripResponses[i] = tripResponse
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	return &TripListResponse{
		Trips:      tripResponses,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetNearbyTrips finds trips near a location
func (s *TripService) GetNearbyTrips(ctx context.Context, location domain.Location, radiusKM float64, status *domain.TripStatus) ([]*TripResponse, error) {
	trips, err := s.tripRepo.GetNearbyTrips(ctx, location, radiusKM, status)
	if err != nil {
		return nil, shared.NewError(shared.ErrCodeInternal, "Failed to find nearby trips", err.Error())
	}

	responses := make([]*TripResponse, len(trips))
	for i, trip := range trips {
		response := &TripResponse{Trip: trip}
		if trip.PricingInfo != nil {
			response.EstimatedFare = &trip.PricingInfo.TotalFare
		}
		responses[i] = response
	}

	return responses, nil
}

// DeleteTrip deletes a trip
func (s *TripService) DeleteTrip(ctx context.Context, tripID uuid.UUID) error {
	trip, err := s.tripRepo.GetTripByID(ctx, tripID)
	if err != nil {
		return shared.NewError(shared.ErrCodeNotFound, "Trip not found", err.Error())
	}

	if trip.IsActive() {
		return shared.NewError(shared.ErrCodeValidation, "Cannot delete active trip", "Trip must be completed or cancelled before deletion")
	}

	if err := s.tripRepo.DeleteTrip(ctx, tripID); err != nil {
		return shared.NewError(shared.ErrCodeInternal, "Failed to delete trip", err.Error())
	}

	return nil
}

// calculateRoute calculates route information (mock implementation)
func (s *TripService) calculateRoute(trip *domain.Trip) error {
	// This would integrate with a real mapping service like Google Maps, Mapbox, etc.
	// For now, we'll use a simple distance calculation
	
	distance := calculateDistance(
		trip.PickupLocation.Latitude, trip.PickupLocation.Longitude,
		trip.DropoffLocation.Latitude, trip.DropoffLocation.Longitude,
	)

	// Estimate duration based on distance (assuming 30 km/h average speed in city)
	duration := int(distance / 30.0 * 60) // minutes

	routeInfo := domain.RouteInfo{
		DistanceKM:      distance,
		DurationMinutes: duration,
		Tolls:          0, // Would be calculated by routing service
	}

	trip.UpdateRouteInfo(routeInfo)
	return nil
}

// calculatePricing calculates trip pricing (mock implementation)
func (s *TripService) calculatePricing(trip *domain.Trip) error {
	if trip.RouteInfo == nil {
		return fmt.Errorf("route info required for pricing calculation")
	}

	// Vietnamese pricing structure (VND)
	baseFare := 15000.0      // Base fare: 15,000 VND
	perKmRate := 12000.0     // Per km: 12,000 VND
	perMinuteRate := 300.0   // Per minute: 300 VND
	surgeRate := 1.0         // No surge by default

	distanceFare := trip.RouteInfo.DistanceKM * perKmRate
	timeFare := float64(trip.RouteInfo.DurationMinutes) * perMinuteRate
	
	totalFare := (baseFare + distanceFare + timeFare) * surgeRate

	pricingInfo := domain.PricingInfo{
		BaseFare:     baseFare,
		DistanceFare: distanceFare,
		TimeFare:     timeFare,
		SurgeRate:    surgeRate,
		TotalFare:    totalFare,
		Currency:     "VND",
	}

	trip.UpdatePricing(pricingInfo)
	return nil
}

// calculateDistance calculates the haversine distance between two points
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth's radius in kilometers

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}