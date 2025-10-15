package application
// Package application contains auth service use cases and DTOs
package application

import (
	"github.com/southern-martin/zride/backend/shared/application"
)

// LoginCommand represents login command
type LoginCommand struct {
	application.BaseCommand
	ZaloAccessToken string `json:"zalo_access_token" binding:"required"`
	DeviceInfo      string `json:"device_info"`
	IPAddress       string `json:"ip_address"`
}

func NewLoginCommand(zaloAccessToken, deviceInfo, ipAddress string) *LoginCommand {
	return &LoginCommand{
		BaseCommand:     application.NewBaseCommand("auth.login"),
		ZaloAccessToken: zaloAccessToken,
		DeviceInfo:      deviceInfo,
		IPAddress:       ipAddress,
	}
}

// RefreshTokenCommand represents refresh token command
type RefreshTokenCommand struct {
	application.BaseCommand
	RefreshToken string `json:"refresh_token" binding:"required"`
	DeviceInfo   string `json:"device_info"`
	IPAddress    string `json:"ip_address"`
}

func NewRefreshTokenCommand(refreshToken, deviceInfo, ipAddress string) *RefreshTokenCommand {
	return &RefreshTokenCommand{
		BaseCommand:  application.NewBaseCommand("auth.refresh_token"),
		RefreshToken: refreshToken,
		DeviceInfo:   deviceInfo,
		IPAddress:    ipAddress,
	}
}

// LogoutCommand represents logout command
type LogoutCommand struct {
	application.BaseCommand
	AccessToken string `json:"access_token" binding:"required"`
}

func NewLogoutCommand(accessToken string) *LogoutCommand {
	return &LogoutCommand{
		BaseCommand: application.NewBaseCommand("auth.logout"),
		AccessToken: accessToken,
	}
}

// UpdateProfileCommand represents update profile command
type UpdateProfileCommand struct {
	application.BaseCommand
	UserID string `json:"user_id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Phone  string `json:"phone"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

func NewUpdateProfileCommand(userID, name, phone, email, avatar string) *UpdateProfileCommand {
	return &UpdateProfileCommand{
		BaseCommand: application.NewBaseCommand("auth.update_profile"),
		UserID:      userID,
		Name:        name,
		Phone:       phone,
		Email:       email,
		Avatar:      avatar,
	}
}

// GetUserQuery represents get user query
type GetUserQuery struct {
	application.BaseQuery
	UserID string `json:"user_id" binding:"required"`
}

func NewGetUserQuery(userID string) *GetUserQuery {
	return &GetUserQuery{
		BaseQuery: application.NewBaseQuery("auth.get_user"),
		UserID:    userID,
	}
}

// ValidateTokenQuery represents validate token query
type ValidateTokenQuery struct {
	application.BaseQuery
	AccessToken string `json:"access_token" binding:"required"`
}

func NewValidateTokenQuery(accessToken string) *ValidateTokenQuery {
	return &ValidateTokenQuery{
		BaseQuery:   application.NewBaseQuery("auth.validate_token"),
		AccessToken: accessToken,
	}
}

// Response DTOs
type LoginResponseDTO struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	TokenType    string  `json:"token_type"`
	ExpiresIn    int64   `json:"expires_in"`
	User         UserDTO `json:"user"`
}

type RefreshTokenResponseDTO struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	TokenType    string  `json:"token_type"`
	ExpiresIn    int64   `json:"expires_in"`
	User         UserDTO `json:"user"`
}

type UserDTO struct {
	application.BaseDTO
	ZaloID      string `json:"zalo_id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
	IsActive    bool   `json:"is_active"`
	LastLoginAt string `json:"last_login_at,omitempty"`
}

type TokenValidationResponseDTO struct {
	Valid  bool      `json:"valid"`
	UserID string    `json:"user_id,omitempty"`
	ZaloID string    `json:"zalo_id,omitempty"`
	User   *UserDTO  `json:"user,omitempty"`
}

// Request DTOs
type LoginRequestDTO struct {
	ZaloAccessToken string `json:"zalo_access_token" binding:"required"`
}

type RefreshTokenRequestDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdateProfileRequestDTO struct {
	Name   string `json:"name" binding:"required"`
	Phone  string `json:"phone"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}