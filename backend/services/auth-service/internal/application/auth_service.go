package application

import (
	"github.com/google/uuid"
	"github.com/southern-martin/zride/backend/services/auth-service/internal/domain"
	"github.com/southern-martin/zride/backend/services/auth-service/internal/infrastructure"
)

type AuthService struct {
	userRepo         domain.UserRepository
	jwtService       *infrastructure.JWTService
	zaloOAuthService *infrastructure.ZaloOAuthService
}

type LoginResult struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *UserProfile `json:"user"`
}

type RefreshResult struct {
	AccessToken string `json:"access_token"`
}

type UserProfile struct {
	ID       uuid.UUID `json:"id"`
	ZaloID   string    `json:"zalo_id"`
	Name     string    `json:"name"`
	Picture  string    `json:"picture"`
	UserType string    `json:"user_type"`
	IsActive bool      `json:"is_active"`
}

func NewAuthService(userRepo domain.UserRepository, jwtService *infrastructure.JWTService, zaloOAuthService *infrastructure.ZaloOAuthService) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		jwtService:       jwtService,
		zaloOAuthService: zaloOAuthService,
	}
}

func (s *AuthService) LoginWithZalo(code string) (*LoginResult, error) {
	// Exchange code for Zalo access token
	tokenResp, err := s.zaloOAuthService.ExchangeCodeForToken(code)
	if err != nil {
		return nil, err
	}

	// Get user info from Zalo
	zaloUser, err := s.zaloOAuthService.GetUserInfo(tokenResp.AccessToken)
	if err != nil {
		return nil, err
	}

	// Find or create user in our database
	user, err := s.userRepo.FindByZaloID(zaloUser.ID)
	if err != nil {
		// User doesn't exist, create new one
		user = &domain.User{
			ID:       uuid.New(),
			ZaloID:   zaloUser.ID,
			Name:     zaloUser.Name,
			Picture:  zaloUser.Picture,
			UserType: "passenger", // Default to passenger
			IsActive: true,
		}

		if err := s.userRepo.Create(user); err != nil {
			return nil, err
		}
	}

	// Generate JWT tokens
	accessToken, refreshToken, err := s.jwtService.GenerateTokens(user.ID, user.ZaloID, user.UserType)
	if err != nil {
		return nil, err
	}

	userProfile := &UserProfile{
		ID:       user.ID,
		ZaloID:   user.ZaloID,
		Name:     user.Name,
		Picture:  user.Picture,
		UserType: user.UserType,
		IsActive: user.IsActive,
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         userProfile,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*RefreshResult, error) {
	accessToken, err := s.jwtService.RefreshAccessToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return &RefreshResult{
		AccessToken: accessToken,
	}, nil
}

func (s *AuthService) ValidateToken(token string) (*UserProfile, error) {
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	return &UserProfile{
		ID:       user.ID,
		ZaloID:   user.ZaloID,
		Name:     user.Name,
		Picture:  user.Picture,
		UserType: user.UserType,
		IsActive: user.IsActive,
	}, nil
}