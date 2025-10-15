// Package domain contains auth service repository interfaces
package domain

import (
	"context"

	"github.com/southern-martin/zride/backend/shared/domain"
)

// UserRepository interface for user data access
type UserRepository interface {
	domain.Repository[*User]
	
	// Custom methods specific to user repository
	FindByZaloID(ctx context.Context, zaloID string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByPhone(ctx context.Context, phone string) (*User, error)
	UpdateLastLogin(ctx context.Context, userID string) error
	UpdateRefreshToken(ctx context.Context, userID, refreshToken string) error
	FindActiveUsers(ctx context.Context, params *domain.PaginationParams) (*domain.PaginatedResult[*User], error)
}

// AuthSessionRepository interface for auth session data access
type AuthSessionRepository interface {
	Save(ctx context.Context, session *AuthSession) error
	FindByAccessToken(ctx context.Context, token string) (*AuthSession, error)
	FindByRefreshToken(ctx context.Context, token string) (*AuthSession, error)
	FindActiveByUserID(ctx context.Context, userID string) ([]*AuthSession, error)
	RevokeSession(ctx context.Context, sessionID string) error
	RevokeAllUserSessions(ctx context.Context, userID string) error
	CleanupExpiredSessions(ctx context.Context) error
}

// ZaloService interface for Zalo integration
type ZaloService interface {
	VerifyAccessToken(ctx context.Context, accessToken string) (*ZaloUserInfo, error)
	GetUserProfile(ctx context.Context, accessToken string) (*ZaloUserInfo, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*ZaloTokenResponse, error)
}

// TokenService interface for JWT token management
type TokenService interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateAccessToken(token string) (*TokenClaims, error)
	ValidateRefreshToken(token string) (*TokenClaims, error)
	RevokeToken(token string) error
}

// ZaloUserInfo represents user info from Zalo
type ZaloUserInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

// ZaloTokenResponse represents Zalo token response
type ZaloTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	UserID    string `json:"user_id"`
	ZaloID    string `json:"zalo_id"`
	TokenType string `json:"token_type"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
}

// Events
const (
	UserRegisteredEvent = "user.registered"
	UserLoggedInEvent   = "user.logged_in"
	UserLoggedOutEvent  = "user.logged_out"
	UserProfileUpdatedEvent = "user.profile_updated"
)

// UserRegistered domain event
type UserRegistered struct {
	*domain.BaseDomainEvent
	UserID string `json:"user_id"`
	ZaloID string `json:"zalo_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

// UserLoggedIn domain event
type UserLoggedIn struct {
	*domain.BaseDomainEvent
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	IPAddress string `json:"ip_address"`
}

// UserLoggedOut domain event
type UserLoggedOut struct {
	*domain.BaseDomainEvent
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
}

// UserProfileUpdated domain event
type UserProfileUpdated struct {
	*domain.BaseDomainEvent
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
}