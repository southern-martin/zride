package interfaces

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the trip service
func SetupRoutes(tripHandler *TripHandler) *gin.Engine {
	// Create Gin engine
	r := gin.New()

	// Add middleware
	r.Use(LoggerMiddleware())
	r.Use(RecoveryMiddleware())
	r.Use(CORSMiddleware())
	r.Use(RequestIDMiddleware())

	// Health check endpoint (no auth required)
	r.GET("/health", tripHandler.HealthCheck)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Public routes (no auth required)
		public := v1.Group("/trips")
		{
			public.GET("/nearby", tripHandler.GetNearbyTrips)
			public.GET("/search", tripHandler.SearchTrips)
		}

		// Protected routes (auth required)
		protected := v1.Group("/trips")
		protected.Use(AuthMiddleware())
		{
			// Trip CRUD operations
			protected.POST("", tripHandler.CreateTrip)
			protected.GET("/:id", tripHandler.GetTrip)
			protected.PUT("/:id", tripHandler.UpdateTrip)
			protected.DELETE("/:id", tripHandler.DeleteTrip)

			// Trip actions
			protected.POST("/:id/accept", tripHandler.AcceptTrip)
			protected.POST("/:id/start", tripHandler.StartTrip)
			protected.POST("/:id/complete", tripHandler.CompleteTrip)
			protected.POST("/:id/cancel", tripHandler.CancelTrip)

			// User-specific trips
			protected.GET("/my", tripHandler.GetMyTrips)
		}
	}

	return r
}

// SetupRoutesWithMiddleware configures routes with custom middleware
func SetupRoutesWithMiddleware(tripHandler *TripHandler, middleware ...gin.HandlerFunc) *gin.Engine {
	r := SetupRoutes(tripHandler)

	// Add custom middleware
	for _, mw := range middleware {
		r.Use(mw)
	}

	return r
}

// RouteInfo represents route information for documentation
type RouteInfo struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
	AuthRequired bool  `json:"auth_required"`
}

// GetRoutes returns all available routes
func GetRoutes() []RouteInfo {
	return []RouteInfo{
		{
			Method:       "GET",
			Path:         "/health",
			Description:  "Health check endpoint",
			AuthRequired: false,
		},
		{
			Method:       "GET",
			Path:         "/api/v1/trips/nearby",
			Description:  "Find trips near a location",
			AuthRequired: false,
		},
		{
			Method:       "GET",
			Path:         "/api/v1/trips/search",
			Description:  "Search trips with filters",
			AuthRequired: false,
		},
		{
			Method:       "POST",
			Path:         "/api/v1/trips",
			Description:  "Create a new trip",
			AuthRequired: true,
		},
		{
			Method:       "GET",
			Path:         "/api/v1/trips/:id",
			Description:  "Get trip by ID",
			AuthRequired: true,
		},
		{
			Method:       "PUT",
			Path:         "/api/v1/trips/:id",
			Description:  "Update trip",
			AuthRequired: true,
		},
		{
			Method:       "DELETE",
			Path:         "/api/v1/trips/:id",
			Description:  "Delete trip",
			AuthRequired: true,
		},
		{
			Method:       "POST",
			Path:         "/api/v1/trips/:id/accept",
			Description:  "Accept a trip as driver",
			AuthRequired: true,
		},
		{
			Method:       "POST",
			Path:         "/api/v1/trips/:id/start",
			Description:  "Start a trip",
			AuthRequired: true,
		},
		{
			Method:       "POST",
			Path:         "/api/v1/trips/:id/complete",
			Description:  "Complete a trip",
			AuthRequired: true,
		},
		{
			Method:       "POST",
			Path:         "/api/v1/trips/:id/cancel",
			Description:  "Cancel a trip",
			AuthRequired: true,
		},
		{
			Method:       "GET",
			Path:         "/api/v1/trips/my",
			Description:  "Get current user's trips",
			AuthRequired: true,
		},
	}
}