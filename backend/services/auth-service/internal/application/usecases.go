// Package application contains auth service use cases
package application

import (
	"context"
	"time"

	"github.com/southern-martin/zride/backend/services/auth-service/internal/domain"
	"github.com/southern-martin/zride/backend/shared/application"
	sharedDomain "github.com/southern-martin/zride/backend/shared/domain"
)

// LoginUseCase handles user login
type LoginUseCase struct {
	userRepo        domain.UserRepository
	sessionRepo     domain.AuthSessionRepository
	zaloService     domain.ZaloService
	tokenService    domain.TokenService
}

// NewLoginUseCase creates new login use case
func NewLoginUseCase(
	userRepo domain.UserRepository,
	sessionRepo domain.AuthSessionRepository,
	zaloService domain.ZaloService,
	tokenService domain.TokenService,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		zaloService:  zaloService,
		tokenService: tokenService,
	}
}

// Execute executes login use case
func (uc *LoginUseCase) Execute(ctx context.Context, cmd *LoginCommand) (*LoginResponseDTO, error) {
	// Verify Zalo access token
	zaloUser, err := uc.zaloService.VerifyAccessToken(ctx, cmd.ZaloAccessToken)
	if err != nil {
		return nil, err
	}

	// Check if user exists
	user, err := uc.userRepo.FindByZaloID(ctx, zaloUser.ID)
	if err != nil {
		// Create new user if not exists
		user, err = domain.NewUser(zaloUser.ID, zaloUser.Name, zaloUser.Phone, zaloUser.Email, zaloUser.Avatar)
		if err != nil {
			return nil, err
		}

		if err := uc.userRepo.Save(ctx, user); err != nil {
			return nil, err
		}
	}

	// Update last login
	user.UpdateLastLogin()
	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := uc.tokenService.GenerateAccessToken(user.GetID())
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.tokenService.GenerateRefreshToken(user.GetID())
	if err != nil {
		return nil, err
	}

	// Save session
	expiresAt := time.Now().Add(24 * time.Hour) // 24 hours
	session := domain.NewAuthSession(
		user.GetID(),
		accessToken,
		refreshToken,
		cmd.DeviceInfo,
		cmd.IPAddress,
		expiresAt,
	)

	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		return nil, err
	}

	// Update user refresh token
	user.SetRefreshToken(refreshToken)
	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	return &LoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    86400, // 24 hours in seconds
		User:         mapUserToDTO(user),
	}, nil
}

// RefreshTokenUseCase handles token refresh
type RefreshTokenUseCase struct {
	userRepo     domain.UserRepository
	sessionRepo  domain.AuthSessionRepository
	tokenService domain.TokenService
}

// NewRefreshTokenUseCase creates new refresh token use case
func NewRefreshTokenUseCase(
	userRepo domain.UserRepository,
	sessionRepo domain.AuthSessionRepository,
	tokenService domain.TokenService,
) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		tokenService: tokenService,
	}
}

// Execute executes refresh token use case
func (uc *RefreshTokenUseCase) Execute(ctx context.Context, cmd *RefreshTokenCommand) (*RefreshTokenResponseDTO, error) {
	// Validate refresh token
	claims, err := uc.tokenService.ValidateRefreshToken(cmd.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Find session by refresh token
	session, err := uc.sessionRepo.FindByRefreshToken(ctx, cmd.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Check if session is active and not expired
	if !session.IsActive || session.IsExpired() {
		return nil, sharedDomain.ErrUnauthorized
	}

	// Find user
	user, err := uc.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new tokens
	accessToken, err := uc.tokenService.GenerateAccessToken(user.GetID())
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := uc.tokenService.GenerateRefreshToken(user.GetID())
	if err != nil {
		return nil, err
	}

	// Update session
	session.AccessToken = accessToken
	session.RefreshToken = newRefreshToken
	session.ExpiresAt = time.Now().Add(24 * time.Hour)
	session.UpdateTimestamp()

	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		return nil, err
	}

	// Update user refresh token
	user.SetRefreshToken(newRefreshToken)
	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	return &RefreshTokenResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    86400, // 24 hours in seconds
		User:         mapUserToDTO(user),
	}, nil
}

// LogoutUseCase handles user logout
type LogoutUseCase struct {
	sessionRepo  domain.AuthSessionRepository
	tokenService domain.TokenService
}

// NewLogoutUseCase creates new logout use case
func NewLogoutUseCase(
	sessionRepo domain.AuthSessionRepository,
	tokenService domain.TokenService,
) *LogoutUseCase {
	return &LogoutUseCase{
		sessionRepo:  sessionRepo,
		tokenService: tokenService,
	}
}

// Execute executes logout use case
func (uc *LogoutUseCase) Execute(ctx context.Context, cmd *LogoutCommand) error {
	// Validate access token
	_, err := uc.tokenService.ValidateAccessToken(cmd.AccessToken)
	if err != nil {
		return err
	}

	// Find and revoke session
	session, err := uc.sessionRepo.FindByAccessToken(ctx, cmd.AccessToken)
	if err != nil {
		return err
	}

	session.Revoke()
	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		return err
	}

	// Revoke token
	return uc.tokenService.RevokeToken(cmd.AccessToken)
}

// GetUserUseCase handles get user profile
type GetUserUseCase struct {
	userRepo domain.UserRepository
}

// NewGetUserUseCase creates new get user use case
func NewGetUserUseCase(userRepo domain.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{userRepo: userRepo}
}

// Execute executes get user use case
func (uc *GetUserUseCase) Execute(ctx context.Context, query *GetUserQuery) (*UserDTO, error) {
	user, err := uc.userRepo.FindByID(ctx, query.UserID)
	if err != nil {
		return nil, err
	}

	userDTO := mapUserToDTO(user)
	return &userDTO, nil
}

// ValidateTokenUseCase handles token validation
type ValidateTokenUseCase struct {
	userRepo     domain.UserRepository
	sessionRepo  domain.AuthSessionRepository
	tokenService domain.TokenService
}

// NewValidateTokenUseCase creates new validate token use case
func NewValidateTokenUseCase(
	userRepo domain.UserRepository,
	sessionRepo domain.AuthSessionRepository,
	tokenService domain.TokenService,
) *ValidateTokenUseCase {
	return &ValidateTokenUseCase{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		tokenService: tokenService,
	}
}

// Execute executes validate token use case
func (uc *ValidateTokenUseCase) Execute(ctx context.Context, query *ValidateTokenQuery) (*TokenValidationResponseDTO, error) {
	// Validate token
	claims, err := uc.tokenService.ValidateAccessToken(query.AccessToken)
	if err != nil {
		return &TokenValidationResponseDTO{Valid: false}, nil
	}

	// Find session
	session, err := uc.sessionRepo.FindByAccessToken(ctx, query.AccessToken)
	if err != nil {
		return &TokenValidationResponseDTO{Valid: false}, nil
	}

	// Check if session is active and not expired
	if !session.IsActive || session.IsExpired() {
		return &TokenValidationResponseDTO{Valid: false}, nil
	}

	// Find user
	user, err := uc.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return &TokenValidationResponseDTO{Valid: false}, nil
	}

	userDTO := mapUserToDTO(user)
	return &TokenValidationResponseDTO{
		Valid:  true,
		UserID: user.GetID(),
		ZaloID: user.ZaloID,
		User:   &userDTO,
	}, nil
}

// Helper function to map domain user to DTO
func mapUserToDTO(user *domain.User) UserDTO {
	dto := UserDTO{
		BaseDTO: application.BaseDTO{
			ID:        user.GetID(),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		ZaloID:   user.ZaloID,
		Name:     user.Name,
		Phone:    user.Phone,
		Email:    user.Email,
		Avatar:   user.Avatar,
		IsActive: user.IsActive,
	}

	if user.LastLoginAt != nil {
		dto.LastLoginAt = user.LastLoginAt.Format(time.RFC3339)
	}

	return dto
}