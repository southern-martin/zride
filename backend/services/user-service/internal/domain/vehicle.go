// Package domain contains vehicle-related entities
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Vehicle represents a driver's vehicle
type Vehicle struct {
	ID           uuid.UUID `json:"id" db:"id"`
	OwnerID      uuid.UUID `json:"owner_id" db:"owner_id"`
	LicensePlate string    `json:"license_plate" db:"license_plate"`
	Brand        string    `json:"brand" db:"brand"`
	Model        string    `json:"model" db:"model"`
	Year         int       `json:"year" db:"year"`
	Color        string    `json:"color" db:"color"`
	VehicleType  string    `json:"vehicle_type" db:"vehicle_type"` // motorbike, car, van
	Capacity     int       `json:"capacity" db:"capacity"`         // number of passengers
	IsActive     bool      `json:"is_active" db:"is_active"`
	IsVerified   bool      `json:"is_verified" db:"is_verified"`
	PhotoURLs    []string  `json:"photo_urls" db:"photo_urls"`
	Documents    []string  `json:"documents" db:"documents"`
	Features     []string  `json:"features" db:"features"` // air_con, wifi, etc.
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	Version      int       `json:"version" db:"version"`
}

// NewVehicle creates a new vehicle
func NewVehicle(ownerID uuid.UUID, licensePlate, brand, model, color, vehicleType string, year, capacity int) (*Vehicle, error) {
	if ownerID == uuid.Nil {
		return nil, errors.New("owner ID is required")
	}
	if licensePlate == "" {
		return nil, errors.New("license plate is required")
	}
	if brand == "" {
		return nil, errors.New("brand is required")
	}
	if model == "" {
		return nil, errors.New("model is required")
	}
	if vehicleType == "" {
		return nil, errors.New("vehicle type is required")
	}
	if year < 1900 || year > time.Now().Year()+1 {
		return nil, errors.New("invalid year")
	}
	if capacity < 1 || capacity > 8 {
		return nil, errors.New("capacity must be between 1 and 8")
	}

	now := time.Now()
	return &Vehicle{
		ID:           uuid.New(),
		OwnerID:      ownerID,
		LicensePlate: licensePlate,
		Brand:        brand,
		Model:        model,
		Year:         year,
		Color:        color,
		VehicleType:  vehicleType,
		Capacity:     capacity,
		IsActive:     true,
		IsVerified:   false,
		PhotoURLs:    make([]string, 0),
		Documents:    make([]string, 0),
		Features:     make([]string, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
		Version:      1,
	}, nil
}

// UpdateVehicleInfo updates basic vehicle information
func (v *Vehicle) UpdateVehicleInfo(brand, model, color string, year, capacity int) error {
	if brand == "" {
		return errors.New("brand is required")
	}
	if model == "" {
		return errors.New("model is required")
	}
	if year < 1900 || year > time.Now().Year()+1 {
		return errors.New("invalid year")
	}
	if capacity < 1 || capacity > 8 {
		return errors.New("capacity must be between 1 and 8")
	}

	v.Brand = brand
	v.Model = model
	v.Color = color
	v.Year = year
	v.Capacity = capacity
	v.UpdatedAt = time.Now()
	return nil
}

// AddPhoto adds a photo URL to the vehicle
func (v *Vehicle) AddPhoto(photoURL string) error {
	if photoURL == "" {
		return errors.New("photo URL is required")
	}
	
	// Check if photo already exists
	for _, url := range v.PhotoURLs {
		if url == photoURL {
			return errors.New("photo already exists")
		}
	}
	
	v.PhotoURLs = append(v.PhotoURLs, photoURL)
	v.UpdatedAt = time.Now()
	return nil
}

// RemovePhoto removes a photo URL from the vehicle
func (v *Vehicle) RemovePhoto(photoURL string) {
	for i, url := range v.PhotoURLs {
		if url == photoURL {
			v.PhotoURLs = append(v.PhotoURLs[:i], v.PhotoURLs[i+1:]...)
			v.UpdatedAt = time.Now()
			break
		}
	}
}

// AddDocument adds a document URL to the vehicle
func (v *Vehicle) AddDocument(documentURL string) error {
	if documentURL == "" {
		return errors.New("document URL is required")
	}
	
	// Check if document already exists
	for _, url := range v.Documents {
		if url == documentURL {
			return errors.New("document already exists")
		}
	}
	
	v.Documents = append(v.Documents, documentURL)
	v.UpdatedAt = time.Now()
	return nil
}

// RemoveDocument removes a document URL from the vehicle
func (v *Vehicle) RemoveDocument(documentURL string) {
	for i, url := range v.Documents {
		if url == documentURL {
			v.Documents = append(v.Documents[:i], v.Documents[i+1:]...)
			v.UpdatedAt = time.Now()
			break
		}
	}
}

// AddFeature adds a feature to the vehicle
func (v *Vehicle) AddFeature(feature string) error {
	if feature == "" {
		return errors.New("feature is required")
	}
	
	// Check if feature already exists
	for _, f := range v.Features {
		if f == feature {
			return errors.New("feature already exists")
		}
	}
	
	v.Features = append(v.Features, feature)
	v.UpdatedAt = time.Now()
	return nil
}

// RemoveFeature removes a feature from the vehicle
func (v *Vehicle) RemoveFeature(feature string) {
	for i, f := range v.Features {
		if f == feature {
			v.Features = append(v.Features[:i], v.Features[i+1:]...)
			v.UpdatedAt = time.Now()
			break
		}
	}
}

// Activate activates the vehicle
func (v *Vehicle) Activate() {
	v.IsActive = true
	v.UpdatedAt = time.Now()
}

// Deactivate deactivates the vehicle
func (v *Vehicle) Deactivate() {
	v.IsActive = false
	v.UpdatedAt = time.Now()
}

// SubmitForVerification submits vehicle for verification
func (v *Vehicle) SubmitForVerification() error {
	if len(v.Documents) == 0 {
		return errors.New("at least one document is required for verification")
	}
	if len(v.PhotoURLs) == 0 {
		return errors.New("at least one photo is required for verification")
	}
	
	v.UpdatedAt = time.Now()
	return nil
}

// ApproveVerification approves vehicle verification
func (v *Vehicle) ApproveVerification() {
	v.IsVerified = true
	v.UpdatedAt = time.Now()
}

// RejectVerification rejects vehicle verification
func (v *Vehicle) RejectVerification() {
	v.IsVerified = false
	v.UpdatedAt = time.Now()
}

// Validate validates the vehicle data
func (v *Vehicle) Validate() error {
	if v.OwnerID == uuid.Nil {
		return errors.New("owner ID is required")
	}
	if v.LicensePlate == "" {
		return errors.New("license plate is required")
	}
	if v.Brand == "" {
		return errors.New("brand is required")
	}
	if v.Model == "" {
		return errors.New("model is required")
	}
	if v.VehicleType == "" {
		return errors.New("vehicle type is required")
	}
	if v.Year < 1900 || v.Year > time.Now().Year()+1 {
		return errors.New("invalid year")
	}
	if v.Capacity < 1 || v.Capacity > 8 {
		return errors.New("capacity must be between 1 and 8")
	}
	return nil
}

// GetID returns the vehicle ID as string
func (v *Vehicle) GetID() string {
	return v.ID.String()
}