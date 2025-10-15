// Package domain contains auth service domain entities and value objects
package domain

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

// User represents the user aggregate root
type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	ZaloID       string     `json:"zalo_id" db:"zalo_id"`
	Name         string     `json:"name" db:"name"`
	Phone        string     `json:"phone" db:"phone"`
	Email        string     `json:"email" db:"email"`
	Picture      string     `json:"picture" db:"picture"`
	UserType     string     `json:"user_type" db:"user_type"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at" db:"last_login_at"`
	RefreshToken string     `json:"-" db:"refresh_token"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	Version      int        `json:"version" db:"version"`
}

// NewUser creates a new user
func NewUser(zaloID, name, phone, email, picture string) (*User, error) {
	// Validate required fields
	if zaloID == "" {
		return nil, errors.New("zalo ID is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}

	// Validate email format if provided
	if email != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(email) {
			return nil, errors.New("invalid email format")
		}
	}

	now := time.Now()
	return &User{
		ID:        uuid.New(),
		ZaloID:    zaloID,
		Name:      name,
		Phone:     phone,
		Email:     email,
		Picture:   picture,
		UserType:  "passenger", // Default to passenger
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
		Version:   1,
	}, nil
}

// GetID returns the user ID as string
func (u *User) GetID() string {
	return u.ID.String()
}

// Validate validates the user data
func (u *User) Validate() error {
	if u.ZaloID == "" {
		return errors.New("zalo ID is required")
	}
	if u.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

// UpdateLastLogin updates the last login timestamp
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
	u.UpdatedAt = now
}

// ChangeUserType changes the user type (passenger/driver)
func (u *User) ChangeUserType(userType string) error {
	if userType != "passenger" && userType != "driver" {
		return errors.New("user type must be either 'passenger' or 'driver'")
	}
	u.UserType = userType
	u.UpdatedAt = time.Now()
	return nil
}

// Deactivate deactivates the user
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Activate activates the user
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// UpdateProfile updates user profile information
func (u *User) UpdateProfile(name, phone, email, picture string) error {
	if name == "" {
		return errors.New("name is required")
	}

	// Validate email format if provided
	if email != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(email) {
			return errors.New("invalid email format")
		}
	}

	u.Name = name
	u.Phone = phone
	u.Email = email
	u.Picture = picture
	u.UpdatedAt = time.Now()
	return nil
}