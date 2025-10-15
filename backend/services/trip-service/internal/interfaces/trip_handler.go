package interfaces

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/southern-martin/zride/backend/services/trip-service/internal/application"
	"github.com/southern-martin/zride/backend/services/trip-service/internal/domain"
	"github.com/southern-martin/zride/backend/shared"
)

// TripHandler handles trip-related HTTP requests
type TripHandler struct {
	tripService *application.TripService
}

// NewTripHandler creates a new trip handler
func NewTripHandler(tripService *application.TripService) *TripHandler {
	return &TripHandler{
		tripService: tripService,
	}
}

// CreateTrip handles POST /trips
func (h *TripHandler) CreateTrip(c *gin.Context) {
	var req application.CreateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeValidation,
				Message: "Invalid request data",
				Details: err.Error(),
			},
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeUnauthorized,
				Message: "User not authenticated",
				Details: "User ID not found in context",
			},
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeValidation,
				Message: "Invalid user ID",
				Details: err.Error(),
			},
		})
		return
	}

	// Set passenger ID to the authenticated user
	req.PassengerID = userID

	trip, err := h.tripService.CreateTrip(c.Request.Context(), req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": trip,
	})
}

// GetTrip handles GET /trips/:id
func (h *TripHandler) GetTrip(c *gin.Context) {
	tripIDStr := c.Param("id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeValidation,
				Message: "Invalid trip ID",
				Details: err.Error(),
			},
		})
		return
	}

	trip, err := h.tripService.GetTrip(c.Request.Context(), tripID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": trip,
	})
}

// UpdateTrip handles PUT /trips/:id
func (h *TripHandler) UpdateTrip(c *gin.Context) {
	tripIDStr := c.Param("id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeValidation,
				Message: "Invalid trip ID",
				Details: err.Error(),
			},
		})
		return
	}

	var req application.UpdateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeValidation,
				Message: "Invalid request data",
				Details: err.Error(),
			},
		})
		return
	}

	trip, err := h.tripService.UpdateTrip(c.Request.Context(), tripID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": trip,
	})
}

// AcceptTrip handles POST /trips/:id/accept
func (h *TripHandler) AcceptTrip(c *gin.Context) {
	tripIDStr := c.Param("id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeValidation,
				Message: "Invalid trip ID",
				Details: err.Error(),
			},
		})
		return
	}

	var req application.AcceptTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeValidation,
				Message: "Invalid request data",
				Details: err.Error(),
			},
		})
		return
	}

	// Get user ID from context and ensure it matches the driver ID
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeUnauthorized,
				Message: "User not authenticated",
				Details: "User ID not found in context",
			},
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil || userID != req.DriverID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeForbidden,
				Message: "Access denied",
				Details: "Cannot accept trip for another driver",
			},
		})
		return
	}

	trip, err := h.tripService.AcceptTrip(c.Request.Context(), tripID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": trip,
	})
}

// StartTrip handles POST /trips/:id/start
func (h *TripHandler) StartTrip(c *gin.Context) {
	tripIDStr := c.Param("id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeValidation,
				Message: "Invalid trip ID",
				Details: err.Error(),
			},
		})
		return
	}

	trip, err := h.tripService.StartTrip(c.Request.Context(), tripID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": trip,
	})
}

// CompleteTrip handles POST /trips/:id/complete
func (h *TripHandler) CompleteTrip(c *gin.Context) {
	tripIDStr := c.Param("id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeValidation,
				Message: "Invalid trip ID",
				Details: err.Error(),
			},
		})
		return
	}

	trip, err := h.tripService.CompleteTrip(c.Request.Context(), tripID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": trip,
	})
}

// CancelTrip handles POST /trips/:id/cancel
func (h *TripHandler) CancelTrip(c *gin.Context) {
	tripIDStr := c.Param("id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeValidation,
				Message: "Invalid trip ID",
				Details: err.Error(),
			},
		})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeValidation,
				Message: "Invalid request data",
				Details: err.Error(),
			},
		})
		return
	}

	trip, err := h.tripService.CancelTrip(c.Request.Context(), tripID, req.Reason)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": trip,
	})
}

// SearchTrips handles GET /trips/search
func (h *TripHandler) SearchTrips(c *gin.Context) {
	var req application.TripSearchRequest

	// Parse query parameters
	if passengerID := c.Query("passenger_id"); passengerID != "" {
		if id, err := uuid.Parse(passengerID); err == nil {
			req.PassengerID = &id
		}
	}

	if driverID := c.Query("driver_id"); driverID != "" {
		if id, err := uuid.Parse(driverID); err == nil {
			req.DriverID = &id
		}
	}

	if status := c.Query("status"); status != "" {
		tripStatus := domain.TripStatus(status)
		req.Status = &tripStatus
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			req.Page = page
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			req.PageSize = pageSize
		}
	}

	req.SortBy = c.Query("sort_by")
	req.SortOrder = c.Query("sort_order")

	// Parse location and radius parameters
	if lat := c.Query("pickup_lat"); lat != "" {
		if lng := c.Query("pickup_lng"); lng != "" {
			if latFloat, err := strconv.ParseFloat(lat, 64); err == nil {
				if lngFloat, err := strconv.ParseFloat(lng, 64); err == nil {
					req.PickupLocation = &domain.Location{
						Latitude:  latFloat,
						Longitude: lngFloat,
						Address:   c.Query("pickup_address"),
					}
					if radiusStr := c.Query("pickup_radius"); radiusStr != "" {
						if radius, err := strconv.ParseFloat(radiusStr, 64); err == nil {
							req.PickupRadius = &radius
						}
					}
				}
			}
		}
	}

	response, err := h.tripService.SearchTrips(c.Request.Context(), req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// GetMyTrips handles GET /trips/my
func (h *TripHandler) GetMyTrips(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeUnauthorized,
				Message: "User not authenticated",
				Details: "User ID not found in context",
			},
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeValidation,
				Message: "Invalid user ID",
				Details: err.Error(),
			},
		})
		return
	}

	var req application.TripSearchRequest
	req.PassengerID = &userID

	// Parse pagination parameters
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			req.Page = page
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			req.PageSize = pageSize
		}
	}

	response, err := h.tripService.SearchTrips(c.Request.Context(), req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// GetNearbyTrips handles GET /trips/nearby
func (h *TripHandler) GetNearbyTrips(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.Query("radius")

	if latStr == "" || lngStr == "" || radiusStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeValidation,
				Message: "Missing required parameters",
				Details: "lat, lng, and radius are required",
			},
		})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeValidation,
				Message: "Invalid latitude",
				Details: err.Error(),
			},
		})
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeValidation,
				Message: "Invalid longitude",
				Details: err.Error(),
			},
		})
		return
	}

	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeValidation,
				Message: "Invalid radius",
				Details: err.Error(),
			},
		})
		return
	}

	location := domain.Location{
		Latitude:  lat,
		Longitude: lng,
		Address:   c.Query("address"),
	}

	var status *domain.TripStatus
	if statusStr := c.Query("status"); statusStr != "" {
		tripStatus := domain.TripStatus(statusStr)
		status = &tripStatus
	}

	trips, err := h.tripService.GetNearbyTrips(c.Request.Context(), location, radius, status)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": trips,
	})
}

// DeleteTrip handles DELETE /trips/:id
func (h *TripHandler) DeleteTrip(c *gin.Context) {
	tripIDStr := c.Param("id")
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeValidation,
				Message: "Invalid trip ID",
				Details: err.Error(),
			},
		})
		return
	}

	err = h.tripService.DeleteTrip(c.Request.Context(), tripID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// HealthCheck handles GET /health
func (h *TripHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "trip-service",
		"version":   "1.0.0",
		"timestamp": "2024-01-01T00:00:00Z",
	})
}

// handleServiceError converts service errors to HTTP responses
func handleServiceError(c *gin.Context, err error) {
	if appError, ok := err.(*shared.Error); ok {
		var statusCode int
		switch appError.Code {
		case shared.ErrCodeNotFound:
			statusCode = http.StatusNotFound
		case shared.ErrCodeValidation:
			statusCode = http.StatusBadRequest
		case shared.ErrCodeUnauthorized:
			statusCode = http.StatusUnauthorized
		case shared.ErrCodeForbidden:
			statusCode = http.StatusForbidden
		case shared.ErrCodeConflict:
			statusCode = http.StatusConflict
		default:
			statusCode = http.StatusInternalServerError
		}

		c.JSON(statusCode, gin.H{
			"error": shared.ErrorResponse{
				Code:    appError.Code,
				Message: appError.Message,
				Details: appError.Details,
			},
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": shared.ErrorResponse{
				Code:    shared.ErrCodeInternal,
				Message: "Internal server error",
				Details: err.Error(),
			},
		})
	}
}