package interfaces

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/southern-martin/zride/backend/shared"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": shared.ErrorResponse{
					Code:    shared.ErrCodeUnauthorized,
					Message: "Authorization header required",
					Details: "Missing Authorization header",
				},
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": shared.ErrorResponse{
					Code:    shared.ErrCodeUnauthorized,
					Message: "Invalid authorization header format",
					Details: "Expected format: Bearer <token>",
				},
			})
			c.Abort()
			return
		}

		token := parts[1]
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": shared.ErrorResponse{
					Code:    shared.ErrCodeUnauthorized,
					Message: "Token is required",
					Details: "Empty token provided",
				},
			})
			c.Abort()
			return
		}

		// TODO: Validate token with auth service
		// For now, we'll extract user ID from a mock implementation
		// In real implementation, this would call the auth service to validate the token
		userID := extractUserIDFromToken(token)
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": shared.ErrorResponse{
					Code:    shared.ErrCodeUnauthorized,
					Message: "Invalid token",
					Details: "Token validation failed",
				},
			})
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("user_id", userID)
		c.Next()
	})
}

// CORSMiddleware handles CORS
func CORSMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// LoggerMiddleware logs HTTP requests
func LoggerMiddleware() gin.HandlerFunc {
	return gin.Logger()
}

// RecoveryMiddleware handles panics
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeInternal,
				Message: "Internal server error",
				Details: "An unexpected error occurred",
			},
		})
	})
}

// RequestIDMiddleware adds a request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	})
}

// RateLimitMiddleware implements basic rate limiting
func RateLimitMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// TODO: Implement proper rate limiting with Redis
		// This is a placeholder implementation
		c.Next()
	})
}

// extractUserIDFromToken extracts user ID from JWT token
// TODO: Replace with actual JWT validation
func extractUserIDFromToken(token string) string {
	// Mock implementation - in real scenario, this would:
	// 1. Call auth service to validate token
	// 2. Extract claims from validated token
	// 3. Return user ID from claims
	
	// For development, we'll use a simple check
	if len(token) > 10 { // Basic validation
		return "550e8400-e29b-41d4-a716-446655440001" // Mock user ID
	}
	return ""
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	// Simple implementation - in production, use UUID or similar
	return "req_" + string(rune(1000000 + (1000000 * 1))) // Mock ID
}