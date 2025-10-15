package interfaces

import (
	"github.com/gin-gonic/gin"
	"github.com/southern-martin/zride/backend/services/auth-service/internal/application"
	"github.com/southern-martin/zride/backend/services/auth-service/internal/infrastructure"
)

type Router struct {
	authHandler *AuthHandler
	jwtService  *infrastructure.JWTService
}

func NewRouter(authService *application.AuthService, jwtService *infrastructure.JWTService) *Router {
	return &Router{
		authHandler: NewAuthHandler(authService),
		jwtService:  jwtService,
	}
}

func (r *Router) SetupRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "auth-service"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/refresh", r.authHandler.RefreshToken)
			auth.POST("/logout", r.authHandler.Logout)
			auth.POST("/validate", r.authHandler.ValidateToken)
		}

		// Protected routes (authentication required)
		protected := v1.Group("/")
		protected.Use(AuthMiddleware(r.jwtService))
		{
			// Add protected routes here if needed
			// protected.GET("/profile", r.authHandler.GetProfile)
		}
	}

	return router
}