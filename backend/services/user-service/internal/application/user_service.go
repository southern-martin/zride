// Package application contains user service business logic
package application

import (
	"github.com/google/uuid"
	"github.com/southern-martin/zride/backend/services/user-service/internal/domain"
	"github.com/southern-martin/zride/backend/shared"
)

type UserService struct {
	userRepo   domain.UserProfileRepository
	vehicleRepo domain.VehicleRepository
	ratingRepo domain.RatingRepository
}

// DTOs for the user service
type CreateUserProfileRequest struct {
	ZaloID   string `json:"zalo_id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Picture  string `json:"picture"`
	UserType string `json:"user_type"`
}

type UpdateProfileRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
	Bio     string `json:"bio"`
}

type UpdatePreferencesRequest struct {
	Preferences map[string]string `json:"preferences"`
	Languages   []string          `json:"languages"`
}

type CreateVehicleRequest struct {
	LicensePlate string `json:"license_plate"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Color        string `json:"color"`
	VehicleType  string `json:"vehicle_type"`
	Year         int    `json:"year"`
	Capacity     int    `json:"capacity"`
}

type UpdateVehicleRequest struct {
	Brand    string `json:"brand"`
	Model    string `json:"model"`
	Color    string `json:"color"`
	Year     int    `json:"year"`
	Capacity int    `json:"capacity"`
}

type CreateRatingRequest struct {
	RatedUserID uuid.UUID `json:"rated_user_id"`
	TripID      uuid.UUID `json:"trip_id"`
	Score       float64   `json:"score"`
	Comment     string    `json:"comment"`
	Tags        []string  `json:"tags"`
}

type UserProfileResponse struct {
	ID              uuid.UUID         `json:"id"`
	ZaloID          string            `json:"zalo_id"`
	Name            string            `json:"name"`
	Phone           string            `json:"phone"`
	Email           string            `json:"email"`
	Picture         string            `json:"picture"`
	UserType        string            `json:"user_type"`
	IsActive        bool              `json:"is_active"`
	IsVerified      bool              `json:"is_verified"`
	Rating          float64           `json:"rating"`
	TotalTrips      int               `json:"total_trips"`
	TotalRatings    int               `json:"total_ratings"`
	Preferences     map[string]string `json:"preferences"`
	Languages       []string          `json:"languages"`
	Bio             string            `json:"bio"`
	HasActiveVehicle bool             `json:"has_active_vehicle"`
}

func NewUserService(userRepo domain.UserProfileRepository, vehicleRepo domain.VehicleRepository, ratingRepo domain.RatingRepository) *UserService {
	return &UserService{
		userRepo:   userRepo,
		vehicleRepo: vehicleRepo,
		ratingRepo: ratingRepo,
	}
}

// CreateUserProfile creates a new user profile
func (s *UserService) CreateUserProfile(req *CreateUserProfileRequest) (*UserProfileResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByZaloID(req.ZaloID)
	if existingUser != nil {
		return nil, shared.NewConflictError("user with this Zalo ID already exists", nil)
	}

	// Create new user profile
	profile, err := domain.NewUserProfile(req.ZaloID, req.Name, req.Phone, req.Email, req.Picture, req.UserType)
	if err != nil {
		return nil, shared.NewValidationError(err.Error(), err)
	}

	if err := s.userRepo.Create(profile); err != nil {
		return nil, shared.NewDatabaseError("failed to create user profile", err)
	}

	return s.buildUserProfileResponse(profile), nil
}

// GetUserProfile retrieves a user profile by ID
func (s *UserService) GetUserProfile(userID uuid.UUID) (*UserProfileResponse, error) {
	profile, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return s.buildUserProfileResponse(profile), nil
}

// GetUserProfileByZaloID retrieves a user profile by Zalo ID
func (s *UserService) GetUserProfileByZaloID(zaloID string) (*UserProfileResponse, error) {
	profile, err := s.userRepo.FindByZaloID(zaloID)
	if err != nil {
		return nil, err
	}

	return s.buildUserProfileResponse(profile), nil
}

// UpdateUserProfile updates a user profile
func (s *UserService) UpdateUserProfile(userID uuid.UUID, req *UpdateProfileRequest) (*UserProfileResponse, error) {
	profile, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if err := profile.UpdateProfile(req.Name, req.Phone, req.Email, req.Picture, req.Bio); err != nil {
		return nil, shared.NewValidationError(err.Error(), err)
	}

	if err := s.userRepo.Update(profile); err != nil {
		return nil, shared.NewDatabaseError("failed to update user profile", err)
	}

	return s.buildUserProfileResponse(profile), nil
}

// UpdateUserPreferences updates user preferences
func (s *UserService) UpdateUserPreferences(userID uuid.UUID, req *UpdatePreferencesRequest) (*UserProfileResponse, error) {
	profile, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	profile.UpdatePreferences(req.Preferences)
	profile.UpdateLanguages(req.Languages)

	if err := s.userRepo.Update(profile); err != nil {
		return nil, shared.NewDatabaseError("failed to update user preferences", err)
	}

	return s.buildUserProfileResponse(profile), nil
}

// SwitchUserType switches between passenger and driver
func (s *UserService) SwitchUserType(userID uuid.UUID, newType string) (*UserProfileResponse, error) {
	profile, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if err := profile.SwitchUserType(newType); err != nil {
		return nil, shared.NewValidationError(err.Error(), err)
	}

	if err := s.userRepo.Update(profile); err != nil {
		return nil, shared.NewDatabaseError("failed to switch user type", err)
	}

	return s.buildUserProfileResponse(profile), nil
}

// CreateVehicle creates a new vehicle for a driver
func (s *UserService) CreateVehicle(ownerID uuid.UUID, req *CreateVehicleRequest) (*domain.Vehicle, error) {
	// Check if user is a driver
	profile, err := s.userRepo.FindByID(ownerID)
	if err != nil {
		return nil, err
	}
	
	if profile.UserType != "driver" {
		return nil, shared.NewValidationError("only drivers can add vehicles", nil)
	}

	// Check if license plate already exists
	existing, _ := s.vehicleRepo.FindByLicensePlate(req.LicensePlate)
	if existing != nil {
		return nil, shared.NewConflictError("vehicle with this license plate already exists", nil)
	}

	vehicle, err := domain.NewVehicle(ownerID, req.LicensePlate, req.Brand, req.Model, req.Color, req.VehicleType, req.Year, req.Capacity)
	if err != nil {
		return nil, shared.NewValidationError(err.Error(), err)
	}

	if err := s.vehicleRepo.Create(vehicle); err != nil {
		return nil, shared.NewDatabaseError("failed to create vehicle", err)
	}

	return vehicle, nil
}

// GetUserVehicles retrieves all vehicles for a user
func (s *UserService) GetUserVehicles(ownerID uuid.UUID) ([]*domain.Vehicle, error) {
	vehicles, err := s.vehicleRepo.FindByOwnerID(ownerID)
	if err != nil {
		return nil, err
	}

	return vehicles, nil
}

// UpdateVehicle updates vehicle information
func (s *UserService) UpdateVehicle(vehicleID uuid.UUID, req *UpdateVehicleRequest) (*domain.Vehicle, error) {
	vehicle, err := s.vehicleRepo.FindByID(vehicleID)
	if err != nil {
		return nil, err
	}

	if err := vehicle.UpdateVehicleInfo(req.Brand, req.Model, req.Color, req.Year, req.Capacity); err != nil {
		return nil, shared.NewValidationError(err.Error(), err)
	}

	if err := s.vehicleRepo.Update(vehicle); err != nil {
		return nil, shared.NewDatabaseError("failed to update vehicle", err)
	}

	return vehicle, nil
}

// CreateRating creates a new rating
func (s *UserService) CreateRating(raterID uuid.UUID, req *CreateRatingRequest) (*domain.Rating, error) {
	// Verify both users exist
	_, err := s.userRepo.FindByID(raterID)
	if err != nil {
		return nil, shared.NewNotFoundError("rater not found", err)
	}

	ratedUser, err := s.userRepo.FindByID(req.RatedUserID)
	if err != nil {
		return nil, shared.NewNotFoundError("rated user not found", err)
	}

	// Create rating
	rating, err := domain.NewRating(raterID, req.RatedUserID, req.TripID, req.Score, req.Comment)
	if err != nil {
		return nil, shared.NewValidationError(err.Error(), err)
	}

	if len(req.Tags) > 0 {
		rating.SetTags(req.Tags)
	}

	if err := s.ratingRepo.Create(rating); err != nil {
		return nil, shared.NewDatabaseError("failed to create rating", err)
	}

	// Update user's rating
	if err := ratedUser.AddRating(req.Score); err != nil {
		return nil, shared.NewValidationError(err.Error(), err)
	}

	if err := s.userRepo.Update(ratedUser); err != nil {
		return nil, shared.NewDatabaseError("failed to update user rating", err)
	}

	return rating, nil
}

// GetUserRatings retrieves ratings for a user
func (s *UserService) GetUserRatings(userID uuid.UUID, limit, offset int) ([]*domain.Rating, error) {
	return s.ratingRepo.FindByUserID(userID, limit, offset)
}

// Helper function to build user profile response
func (s *UserService) buildUserProfileResponse(profile *domain.UserProfile) *UserProfileResponse {
	// Check if user has active vehicles
	hasActiveVehicle := false
	if profile.UserType == "driver" {
		vehicles, _ := s.vehicleRepo.FindActiveByOwnerID(profile.ID)
		hasActiveVehicle = len(vehicles) > 0
	}

	return &UserProfileResponse{
		ID:               profile.ID,
		ZaloID:           profile.ZaloID,
		Name:             profile.Name,
		Phone:            profile.Phone,
		Email:            profile.Email,
		Picture:          profile.Picture,
		UserType:         profile.UserType,
		IsActive:         profile.IsActive,
		IsVerified:       profile.IsVerified,
		Rating:           profile.Rating,
		TotalTrips:       profile.TotalTrips,
		TotalRatings:     profile.TotalRatings,
		Preferences:      profile.Preferences,
		Languages:        profile.Languages,
		Bio:              profile.Bio,
		HasActiveVehicle: hasActiveVehicle,
	}
}