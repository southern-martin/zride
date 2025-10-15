package routes


import (
	"github.com/gin-gonic/gin"
	"github.com/southern-martin/zride/backend/services/ai-matching-service/internal/interfaces/http/handlers"
	"github.com/southern-martin/zride/backend/services/ai-matching-service/internal/interfaces/http/middleware"
)

func SetupRoutes(router *gin.Engine, matchingHandler *handlers.MatchingHandler) {
	// Apply middleware
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLogger())
	router.Use(gin.Recovery())

	// API v1 group
	v1 := router.Group("/api/v1")
	
	// Health check
	v1.GET("/health", matchingHandler.HealthCheck)
	
	// Matching routes
	matching := v1.Group("/matching")
	{
		// Match requests
		matching.POST("/requests", matchingHandler.CreateMatchRequest)
		matching.GET("/passengers/:passengerID/requests", matchingHandler.GetMatchRequests)
		
		// Driver matches
		matching.GET("/drivers/:driverID/matches", matchingHandler.GetAvailableMatches)
		matching.POST("/drivers/:driverID/matches/:matchResultID/accept", matchingHandler.AcceptMatch)
	}
}