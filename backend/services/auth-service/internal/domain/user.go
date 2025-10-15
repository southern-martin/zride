package domain
// Package domain contains auth service domain entities and value objects
package domain

import (
	"errors"
	"regexp"
	"time"

	"github.com/southern-martin/zride/backend/shared/domain"
)

// User represents the user aggregate root
type User struct {
	domain.Entity
	ZaloID       string    `json:"zalo_id" db:"zalo_id"`
	Name         string    `json:"name" db:"name"`
	Phone        string    `json:"phone" db:"phone"`
	Email        string    `json:"email" db:"email"`
	Avatar       string    `json:"avatar" db:"avatar"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at" db:"last_login_at"`
	RefreshToken string    `json:"-" db:"refresh_token"`
	Version      int       `json:"version" db:"version"`
}

// NewUser creates a new user
func NewUser(zaloID, name, phone, email, avatar string) (*User, error) {
	// Validate required fields
	if zaloID == "" {
		return nil, errors.New("zalo ID is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}

	// Validate email format if provided
	if email != "" && !isValidEmail(email) {
		return nil, errors.New("invalid email format")
	}

	// Validate phone format if provided
	if phone != "" && !isValidPhone(phone) {
		return nil, errors.New("invalid phone format")
	}

	user := &User{
		Entity:   domain.NewEntity(),
		ZaloID:   zaloID,
		Name:     name,
		Phone:    phone,
		Email:    email,
		Avatar:   avatar,
		IsActive: true,
		Version:  1,
	}

	return user, nil
}

// GetID implements AggregateRoot interface
func (u *User) GetID() string {
	return u.ID.String()
}

// GetVersion implements AggregateRoot interface
func (u *User) GetVersion() int {
	return u.Version
}

// MarkAsModified implements AggregateRoot interface
func (u *User) MarkAsModified() {
	u.Version++
	u.UpdateTimestamp()
}

// UpdateProfile updates user profile
func (u *User) UpdateProfile(name, phone, email, avatar string) error {
	if name == "" {
		return errors.New("name is required")
	}

	// Validate email format if provided
	if email != "" && !isValidEmail(email) {
		return errors.New("invalid email format")
	}

	// Validate phone format if provided
	if phone != "" && !isValidPhone(phone) {
		return errors.New("invalid phone format")
	}

	u.Name = name
	u.Phone = phone
	u.Email = email
	u.Avatar = avatar
	u.MarkAsModified()

	return nil
}

// UpdateLastLogin updates last login timestamp
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
	u.MarkAsModified()
}

// SetRefreshToken sets refresh token
func (u *User) SetRefreshToken(token string) {
	u.RefreshToken = token
	u.MarkAsModified()
}

// ClearRefreshToken clears refresh token (logout)
func (u *User) ClearRefreshToken() {
	u.RefreshToken = ""
	u.MarkAsModified()
}

// Deactivate deactivates user account
func (u *User) Deactivate() {
	u.IsActive = false
	u.ClearRefreshToken()
	u.MarkAsModified()
}

// Activate activates user account
func (u *User) Activate() {
	u.IsActive = true
	u.MarkAsModified()
}

// AuthSession represents an authentication session
type AuthSession struct {
	domain.Entity
	UserID       string    `json:"user_id" db:"user_id"`
	AccessToken  string    `json:"access_token" db:"access_token"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	DeviceInfo   string    `json:"device_info" db:"device_info"`
	IPAddress    string    `json:"ip_address" db:"ip_address"`
}

// NewAuthSession creates a new auth session
func NewAuthSession(userID, accessToken, refreshToken, deviceInfo, ipAddress string, expiresAt time.Time) *AuthSession {
	return &AuthSession{
		Entity:       domain.NewEntity(),
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		IsActive:     true,
		DeviceInfo:   deviceInfo,
		IPAddress:    ipAddress,
	}
}

// IsExpired checks if session is expired
func (s *AuthSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// Revoke revokes the session
func (s *AuthSession) Revoke() {
	s.IsActive = false
	s.UpdateTimestamp()
}

// Utility functions for validation
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidPhone(phone string) bool {
	// Vietnamese phone number format
	phoneRegex := regexp.MustCompile(`^(\+84|84|0)[0-9]{9,10}$`)
	return phoneRegex.MatchString(phone)
}