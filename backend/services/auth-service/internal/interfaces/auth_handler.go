package interfaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/southern-martin/zride/backend/services/auth-service/internal/application"
	"github.com/southern-martin/zride/backend/shared"
)

type AuthHandler struct {
	authService *application.AuthService
}

type LoginRequest struct {
	Code string `json:"code" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string                   `json:"access_token"`
	RefreshToken string                   `json:"refresh_token"`
	User         *application.UserProfile `json:"user"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

func NewAuthHandler(authService *application.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	result, err := h.authService.LoginWithZalo(req.Code)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response := LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User:         result.User,
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	result, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response := RefreshResponse{
		AccessToken: result.AccessToken,
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// For stateless JWT, logout is handled client-side
	// In a production system, you might want to maintain a blacklist of tokens
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (h *AuthHandler) ValidateToken(c *gin.Context) {
	// This endpoint can be used by other services to validate tokens
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Extract token from "Bearer <token>" format
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		return
	}

	user, err := h.authService.ValidateToken(token)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"user":  user,
	})
}

func (h *AuthHandler) handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *shared.ValidationError:
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Message})
	case *shared.NotFoundError:
		c.JSON(http.StatusNotFound, gin.H{"error": e.Message})
	case *shared.ConflictError:
		c.JSON(http.StatusConflict, gin.H{"error": e.Message})
	case *shared.ExternalServiceError:
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "External service error"})
	case *shared.InternalError:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unknown error"})
	}
}