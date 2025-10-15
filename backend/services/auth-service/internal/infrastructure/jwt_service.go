package infrastructure

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/southern-martin/zride/backend/shared"
)

type JWTService struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	ZaloID   string    `json:"zalo_id"`
	UserType string    `json:"user_type"`
	jwt.RegisteredClaims
}

func NewJWTService(secretKey string, accessExpiry, refreshExpiry time.Duration) *JWTService {
	return &JWTService{
		secretKey:     []byte(secretKey),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

func (j *JWTService) GenerateTokens(userID uuid.UUID, zaloID, userType string) (accessToken, refreshToken string, err error) {
	// Generate access token
	accessClaims := Claims{
		UserID:   userID,
		ZaloID:   zaloID,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "zride-auth-service",
			Subject:   userID.String(),
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString(j.secretKey)
	if err != nil {
		return "", "", shared.NewInternalError("failed to generate access token", err)
	}

	// Generate refresh token
	refreshClaims := Claims{
		UserID:   userID,
		ZaloID:   zaloID,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "zride-auth-service",
			Subject:   userID.String(),
		},
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString(j.secretKey)
	if err != nil {
		return "", "", shared.NewInternalError("failed to generate refresh token", err)
	}

	return accessToken, refreshToken, nil
}

func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, shared.NewValidationError("invalid token signing method", nil)
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, shared.NewValidationError("invalid token", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, shared.NewValidationError("invalid token claims", nil)
}

func (j *JWTService) RefreshAccessToken(refreshToken string) (newAccessToken string, err error) {
	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	// Check if it's close to expiry (within 5 minutes)
	if time.Until(claims.RegisteredClaims.ExpiresAt.Time) < 5*time.Minute {
		return "", shared.NewValidationError("refresh token is about to expire", nil)
	}

	// Generate new access token
	newAccessToken, _, err = j.GenerateTokens(claims.UserID, claims.ZaloID, claims.UserType)
	return newAccessToken, err
}