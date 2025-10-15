// Package interfaces provides middleware for user service
package interfaces

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens by calling the auth service
func AuthMiddleware(authServiceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// For now, we'll do a simple validation
		// In production, you'd call the auth service to validate the token
		// TODO: Implement proper auth service validation
		
		// Mock validation - extract user info from token
		// This should be replaced with actual auth service call
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set mock user information in context
		// In production, this would come from the auth service response
		c.Set("user_id", "550e8400-e29b-41d4-a716-446655440000") // Mock UUID
		c.Set("zalo_id", "mock_zalo_id")
		c.Set("user_type", "passenger")

		c.Next()
	}
}

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// LoggingMiddleware logs requests
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return ""
	})
}

// RecoveryMiddleware handles panics
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.Recovery()
}