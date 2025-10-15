// Package domain contains user service domain entitiespackage domain

package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// UserProfile represents a comprehensive user profile
type UserProfile struct {
	ID              uuid.UUID          `json:"id" db:"id"`
	ZaloID          string             `json:"zalo_id" db:"zalo_id"`
	Name            string             `json:"name" db:"name"`
	Phone           string             `json:"phone" db:"phone"`
	Email           string             `json:"email" db:"email"`
	Picture         string             `json:"picture" db:"picture"`
	UserType        string             `json:"user_type" db:"user_type"` // passenger, driver
	IsActive        bool               `json:"is_active" db:"is_active"`
	IsVerified      bool               `json:"is_verified" db:"is_verified"`
	VerificationDoc string             `json:"verification_doc" db:"verification_doc"`
	Rating          float64            `json:"rating" db:"rating"`
	TotalTrips      int                `json:"total_trips" db:"total_trips"`
	TotalRatings    int                `json:"total_ratings" db:"total_ratings"`
	Preferences     map[string]string  `json:"preferences" db:"preferences"`
	Languages       []string           `json:"languages" db:"languages"`
	Bio             string             `json:"bio" db:"bio"`
	CreatedAt       time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" db:"updated_at"`
	LastActiveAt    *time.Time         `json:"last_active_at" db:"last_active_at"`
	Version         int                `json:"version" db:"version"`
}

// NewUserProfile creates a new user profile from auth service data
func NewUserProfile(zaloID, name, phone, email, picture, userType string) (*UserProfile, error) {
	if zaloID == "" {
		return nil, errors.New("zalo ID is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}
	if userType != "passenger" && userType != "driver" {
		return nil, errors.New("user type must be either 'passenger' or 'driver'")
	}

	now := time.Now()
	return &UserProfile{
		ID:           uuid.New(),
		ZaloID:       zaloID,
		Name:         name,
		Phone:        phone,
		Email:        email,
		Picture:      picture,
		UserType:     userType,
		IsActive:     true,
		IsVerified:   false,
		Rating:       0.0,
		TotalTrips:   0,
		TotalRatings: 0,
		Preferences:  make(map[string]string),
		Languages:    []string{"vi"}, // Default Vietnamese
		CreatedAt:    now,
		UpdatedAt:    now,
		Version:      1,
	}, nil
}

// UpdateProfile updates basic profile information
func (u *UserProfile) UpdateProfile(name, phone, email, picture, bio string) error {
	if name == "" {
		return errors.New("name is required")
	}

	u.Name = name
	u.Phone = phone
	u.Email = email
	u.Picture = picture
	u.Bio = bio
	u.UpdatedAt = time.Now()
	return nil
}

// SwitchUserType changes between passenger and driver
func (u *UserProfile) SwitchUserType(newType string) error {
	if newType != "passenger" && newType != "driver" {
		return errors.New("user type must be either 'passenger' or 'driver'")
	}
	
	u.UserType = newType
	u.UpdatedAt = time.Now()
	
	// If switching to driver, reset verification
	if newType == "driver" {
		u.IsVerified = false
		u.VerificationDoc = ""
	}
	
	return nil
}

// UpdatePreferences updates user preferences
func (u *UserProfile) UpdatePreferences(preferences map[string]string) {
	if u.Preferences == nil {
		u.Preferences = make(map[string]string)
	}
	
	for key, value := range preferences {
		u.Preferences[key] = value
	}
	u.UpdatedAt = time.Now()
}

// UpdateLanguages updates user language preferences
func (u *UserProfile) UpdateLanguages(languages []string) {
	if len(languages) == 0 {
		languages = []string{"vi"} // Default Vietnamese
	}
	
	u.Languages = languages
	u.UpdatedAt = time.Now()
}

// UpdateActivity updates last active timestamp
func (u *UserProfile) UpdateActivity() {
	now := time.Now()
	u.LastActiveAt = &now
	u.UpdatedAt = now
}

// Deactivate deactivates the user profile
func (u *UserProfile) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Activate reactivates the user profile
func (u *UserProfile) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// SubmitVerification submits verification documents
func (u *UserProfile) SubmitVerification(documentURL string) error {
	if documentURL == "" {
		return errors.New("verification document URL is required")
	}
	
	u.VerificationDoc = documentURL
	u.UpdatedAt = time.Now()
	return nil
}

// ApproveVerification approves user verification
func (u *UserProfile) ApproveVerification() {
	u.IsVerified = true
	u.VerificationDoc = "" // Clear document after approval
	u.UpdatedAt = time.Now()
}

// RejectVerification rejects user verification
func (u *UserProfile) RejectVerification() {
	u.IsVerified = false
	u.VerificationDoc = "" // Clear document after rejection
	u.UpdatedAt = time.Now()
}

// AddRating adds a new rating to the user
func (u *UserProfile) AddRating(rating float64) error {
	if rating < 1.0 || rating > 5.0 {
		return errors.New("rating must be between 1.0 and 5.0")
	}
	
	// Calculate new average rating
	totalRatingPoints := u.Rating * float64(u.TotalRatings)
	totalRatingPoints += rating
	u.TotalRatings++
	u.Rating = totalRatingPoints / float64(u.TotalRatings)
	
	u.UpdatedAt = time.Now()
	return nil
}

// IncrementTripCount increments the total trip count
func (u *UserProfile) IncrementTripCount() {
	u.TotalTrips++
	u.UpdatedAt = time.Now()
}

// Validate validates the user profile data
func (u *UserProfile) Validate() error {
	if u.ZaloID == "" {
		return errors.New("zalo ID is required")
	}
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.UserType != "passenger" && u.UserType != "driver" {
		return errors.New("user type must be either 'passenger' or 'driver'")
	}
	return nil
}

// GetID returns the user profile ID as string
func (u *UserProfile) GetID() string {
	return u.ID.String()
}