package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TripStatus represents the status of a trip
type TripStatus string

const (
	TripStatusRequested  TripStatus = "requested"
	TripStatusAccepted   TripStatus = "accepted"
	TripStatusInProgress TripStatus = "in_progress"
	TripStatusCompleted  TripStatus = "completed"
	TripStatusCancelled  TripStatus = "cancelled"
)

// Location represents a geographic location
type Location struct {
	Address   string  `json:"address"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	PlaceID   string  `json:"place_id,omitempty"`
}

// Value implements the driver.Valuer interface
func (l Location) Value() (driver.Value, error) {
	return json.Marshal(l)
}

// Scan implements the sql.Scanner interface
func (l *Location) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, l)
	case string:
		return json.Unmarshal([]byte(v), l)
	default:
		return errors.New("cannot scan into Location")
	}
}

// RouteInfo contains route calculation details
type RouteInfo struct {
	DistanceKM      float64 `json:"distance_km"`
	DurationMinutes int     `json:"duration_minutes"`
	PolylineEncoded string  `json:"polyline_encoded,omitempty"`
	Tolls           float64 `json:"tolls,omitempty"`
}

// Value implements the driver.Valuer interface
func (r RouteInfo) Value() (driver.Value, error) {
	return json.Marshal(r)
}

// Scan implements the sql.Scanner interface
func (r *RouteInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, r)
	case string:
		return json.Unmarshal([]byte(v), r)
	default:
		return errors.New("cannot scan into RouteInfo")
	}
}

// PricingInfo contains pricing details
type PricingInfo struct {
	BaseFare      float64 `json:"base_fare"`
	DistanceFare  float64 `json:"distance_fare"`
	TimeFare      float64 `json:"time_fare"`
	SurgeRate     float64 `json:"surge_rate"`
	TotalFare     float64 `json:"total_fare"`
	Currency      string  `json:"currency"`
	PaymentMethod string  `json:"payment_method,omitempty"`
}

// Value implements the driver.Valuer interface
func (p PricingInfo) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan implements the sql.Scanner interface
func (p *PricingInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, p)
	case string:
		return json.Unmarshal([]byte(v), p)
	default:
		return errors.New("cannot scan into PricingInfo")
	}
}

// Trip represents a trip in the system
type Trip struct {
	ID               uuid.UUID    `json:"id" db:"id"`
	PassengerID      uuid.UUID    `json:"passenger_id" db:"passenger_id"`
	DriverID         *uuid.UUID   `json:"driver_id,omitempty" db:"driver_id"`
	VehicleID        *uuid.UUID   `json:"vehicle_id,omitempty" db:"vehicle_id"`
	PickupLocation   Location     `json:"pickup_location" db:"pickup_location"`
	DropoffLocation  Location     `json:"dropoff_location" db:"dropoff_location"`
	ScheduledTime    *time.Time   `json:"scheduled_time,omitempty" db:"scheduled_time"`
	RequestedAt      time.Time    `json:"requested_at" db:"requested_at"`
	AcceptedAt       *time.Time   `json:"accepted_at,omitempty" db:"accepted_at"`
	PickupTime       *time.Time   `json:"pickup_time,omitempty" db:"pickup_time"`
	DropoffTime      *time.Time   `json:"dropoff_time,omitempty" db:"dropoff_time"`
	Status           TripStatus   `json:"status" db:"status"`
	RouteInfo        *RouteInfo   `json:"route_info,omitempty" db:"route_info"`
	PricingInfo      *PricingInfo `json:"pricing_info,omitempty" db:"pricing_info"`
	PassengerNotes   string       `json:"passenger_notes,omitempty" db:"passenger_notes"`
	DriverNotes      string       `json:"driver_notes,omitempty" db:"driver_notes"`
	CancellationInfo string       `json:"cancellation_info,omitempty" db:"cancellation_info"`
	CreatedAt        time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time   `json:"deleted_at,omitempty" db:"deleted_at"`
}

// NewTrip creates a new trip
func NewTrip(passengerID uuid.UUID, pickup, dropoff Location, scheduledTime *time.Time, notes string) *Trip {
	return &Trip{
		ID:              uuid.New(),
		PassengerID:     passengerID,
		PickupLocation:  pickup,
		DropoffLocation: dropoff,
		ScheduledTime:   scheduledTime,
		RequestedAt:     time.Now(),
		Status:          TripStatusRequested,
		PassengerNotes:  notes,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

// AcceptTrip accepts a trip by a driver
func (t *Trip) AcceptTrip(driverID, vehicleID uuid.UUID) error {
	if t.Status != TripStatusRequested {
		return fmt.Errorf("cannot accept trip with status %s", t.Status)
	}

	now := time.Now()
	t.DriverID = &driverID
	t.VehicleID = &vehicleID
	t.AcceptedAt = &now
	t.Status = TripStatusAccepted
	t.UpdatedAt = now

	return nil
}

// StartTrip marks the trip as in progress
func (t *Trip) StartTrip() error {
	if t.Status != TripStatusAccepted {
		return fmt.Errorf("cannot start trip with status %s", t.Status)
	}

	now := time.Now()
	t.PickupTime = &now
	t.Status = TripStatusInProgress
	t.UpdatedAt = now

	return nil
}

// CompleteTrip marks the trip as completed
func (t *Trip) CompleteTrip() error {
	if t.Status != TripStatusInProgress {
		return fmt.Errorf("cannot complete trip with status %s", t.Status)
	}

	now := time.Now()
	t.DropoffTime = &now
	t.Status = TripStatusCompleted
	t.UpdatedAt = now

	return nil
}

// CancelTrip cancels the trip with a reason
func (t *Trip) CancelTrip(reason string) error {
	if t.Status == TripStatusCompleted {
		return errors.New("cannot cancel completed trip")
	}

	now := time.Now()
	t.Status = TripStatusCancelled
	t.CancellationInfo = reason
	t.UpdatedAt = now

	return nil
}

// UpdateRouteInfo updates the route information
func (t *Trip) UpdateRouteInfo(routeInfo RouteInfo) {
	t.RouteInfo = &routeInfo
	t.UpdatedAt = time.Now()
}

// UpdatePricing updates the pricing information
func (t *Trip) UpdatePricing(pricingInfo PricingInfo) {
	t.PricingInfo = &pricingInfo
	t.UpdatedAt = time.Now()
}

// IsActive checks if the trip is in an active state
func (t *Trip) IsActive() bool {
	return t.Status == TripStatusRequested ||
		t.Status == TripStatusAccepted ||
		t.Status == TripStatusInProgress
}

// CanBeModified checks if the trip can be modified
func (t *Trip) CanBeModified() bool {
	return t.Status == TripStatusRequested || t.Status == TripStatusAccepted
}

// GetDuration returns the trip duration if completed
func (t *Trip) GetDuration() *time.Duration {
	if t.PickupTime != nil && t.DropoffTime != nil {
		duration := t.DropoffTime.Sub(*t.PickupTime)
		return &duration
	}
	return nil
}

// Validate validates the trip data
func (t *Trip) Validate() error {
	if t.PassengerID == uuid.Nil {
		return errors.New("passenger ID is required")
	}

	if t.PickupLocation.Address == "" {
		return errors.New("pickup location address is required")
	}

	if t.DropoffLocation.Address == "" {
		return errors.New("dropoff location address is required")
	}

	if t.PickupLocation.Latitude < -90 || t.PickupLocation.Latitude > 90 {
		return errors.New("pickup latitude must be between -90 and 90")
	}

	if t.PickupLocation.Longitude < -180 || t.PickupLocation.Longitude > 180 {
		return errors.New("pickup longitude must be between -180 and 180")
	}

	if t.DropoffLocation.Latitude < -90 || t.DropoffLocation.Latitude > 90 {
		return errors.New("dropoff latitude must be between -90 and 90")
	}

	if t.DropoffLocation.Longitude < -180 || t.DropoffLocation.Longitude > 180 {
		return errors.New("dropoff longitude must be between -180 and 180")
	}

	if t.ScheduledTime != nil && t.ScheduledTime.Before(time.Now()) {
		return errors.New("scheduled time cannot be in the past")
	}

	return nil
}

// TripSearchCriteria represents search criteria for finding trips
type TripSearchCriteria struct {
	PassengerID      *uuid.UUID  `json:"passenger_id,omitempty"`
	DriverID         *uuid.UUID  `json:"driver_id,omitempty"`
	Status           *TripStatus `json:"status,omitempty"`
	FromTime         *time.Time  `json:"from_time,omitempty"`
	ToTime           *time.Time  `json:"to_time,omitempty"`
	PickupRadius     *float64    `json:"pickup_radius,omitempty"`     // in kilometers
	DropoffRadius    *float64    `json:"dropoff_radius,omitempty"`    // in kilometers
	PickupLocation   *Location   `json:"pickup_location,omitempty"`
	DropoffLocation  *Location   `json:"dropoff_location,omitempty"`
	MaxPrice         *float64    `json:"max_price,omitempty"`
	MinPrice         *float64    `json:"min_price,omitempty"`
	Limit            int         `json:"limit,omitempty"`
	Offset           int         `json:"offset,omitempty"`
	SortBy           string      `json:"sort_by,omitempty"` // "created_at", "price", "distance"
	SortOrder        string      `json:"sort_order,omitempty"` // "asc", "desc"
}