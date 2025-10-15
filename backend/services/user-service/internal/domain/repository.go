// Package domain contains repository interfaces for user service
package domain

import (
	"github.com/google/uuid"
)

// UserProfileRepository interface for user profile data access
type UserProfileRepository interface {
	FindByID(id uuid.UUID) (*UserProfile, error)
	FindByZaloID(zaloID string) (*UserProfile, error)
	FindByEmail(email string) (*UserProfile, error)
	FindByPhone(phone string) (*UserProfile, error)
	Create(profile *UserProfile) error
	Update(profile *UserProfile) error
	Delete(id uuid.UUID) error
	FindByUserType(userType string, limit, offset int) ([]*UserProfile, error)
	FindVerifiedDrivers(limit, offset int) ([]*UserProfile, error)
	Search(query string, limit, offset int) ([]*UserProfile, error)
}

// VehicleRepository interface for vehicle data access
type VehicleRepository interface {
	FindByID(id uuid.UUID) (*Vehicle, error)
	FindByOwnerID(ownerID uuid.UUID) ([]*Vehicle, error)
	FindByLicensePlate(licensePlate string) (*Vehicle, error)
	Create(vehicle *Vehicle) error
	Update(vehicle *Vehicle) error
	Delete(id uuid.UUID) error
	FindActiveByOwnerID(ownerID uuid.UUID) ([]*Vehicle, error)
	FindPendingVerification(limit, offset int) ([]*Vehicle, error)
}

// RatingRepository interface for user rating data access
type RatingRepository interface {
	Create(rating *Rating) error
	FindByUserID(userID uuid.UUID, limit, offset int) ([]*Rating, error)
	FindByRaterID(raterID uuid.UUID, limit, offset int) ([]*Rating, error)
	GetAverageRating(userID uuid.UUID) (float64, int, error)
	FindByTripID(tripID uuid.UUID) (*Rating, error)
	Update(rating *Rating) error
	Delete(id uuid.UUID) error
}