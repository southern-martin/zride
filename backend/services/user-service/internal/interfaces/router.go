// Package interfaces provides routing for user service
package interfaces

import (
	"github.com/gin-gonic/gin"
	"github.com/southern-martin/zride/backend/services/user-service/internal/application"
)

type Router struct {
	userHandler     *UserHandler
	authServiceURL  string
}

func NewRouter(userService *application.UserService, authServiceURL string) *Router {
	return &Router{
		userHandler:    NewUserHandler(userService),
		authServiceURL: authServiceURL,
	}
}

func (r *Router) SetupRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	
	router := gin.New()
	router.Use(LoggingMiddleware())
	router.Use(RecoveryMiddleware())
	router.Use(CORSMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "user-service"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (limited public access)
		users := v1.Group("/users")
		{
			users.POST("/", r.userHandler.CreateUserProfile)                    // Create profile
			users.GET("/zalo/:zalo_id", r.userHandler.GetUserProfileByZaloID)   // Get by Zalo ID
		}

		// Protected routes (authentication required)
		protected := v1.Group("/")
		protected.Use(AuthMiddleware(r.authServiceURL))
		{
			// User profile management
			userProfile := protected.Group("/users")
			{
				userProfile.GET("/:id", r.userHandler.GetUserProfile)               // Get profile
				userProfile.PUT("/:id", r.userHandler.UpdateUserProfile)            // Update profile
				userProfile.PUT("/:id/preferences", r.userHandler.UpdateUserPreferences) // Update preferences
				userProfile.POST("/:id/switch-type", r.userHandler.SwitchUserType)  // Switch user type
			}

			// Vehicle management
			vehicles := protected.Group("/users/:id/vehicles")
			{
				vehicles.POST("/", r.userHandler.CreateVehicle)                     // Create vehicle
				vehicles.GET("/", r.userHandler.GetUserVehicles)                    // Get user vehicles
			}
			
			// Individual vehicle operations
			vehicle := protected.Group("/vehicles")
			{
				vehicle.PUT("/:vehicle_id", r.userHandler.UpdateVehicle)            // Update vehicle
			}

			// Rating management
			ratings := protected.Group("/")
			{
				ratings.POST("/ratings", r.userHandler.CreateRating)                // Create rating
				ratings.GET("/users/:id/ratings", r.userHandler.GetUserRatings)     // Get user ratings
			}
		}
	}

	return router
}