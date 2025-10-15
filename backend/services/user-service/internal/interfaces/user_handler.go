// Package interfaces provides HTTP handlers for user servicepackage interfaces

package interfaces

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/southern-martin/zride/backend/services/user-service/internal/application"
	"github.com/southern-martin/zride/backend/shared"
)

type UserHandler struct {
	userService *application.UserService
}

func NewUserHandler(userService *application.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUserProfile creates a new user profile
func (h *UserHandler) CreateUserProfile(c *gin.Context) {
	var req application.CreateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	profile, err := h.userService.CreateUserProfile(&req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, profile)
}

// GetUserProfile gets user profile by ID
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	profile, err := h.userService.GetUserProfile(userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, profile)
}

// GetUserProfileByZaloID gets user profile by Zalo ID
func (h *UserHandler) GetUserProfileByZaloID(c *gin.Context) {
	zaloID := c.Param("zalo_id")
	if zaloID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Zalo ID is required"})
		return
	}

	profile, err := h.userService.GetUserProfileByZaloID(zaloID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateUserProfile updates user profile
func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req application.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	profile, err := h.userService.UpdateUserProfile(userID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateUserPreferences updates user preferences
func (h *UserHandler) UpdateUserPreferences(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req application.UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	profile, err := h.userService.UpdateUserPreferences(userID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, profile)
}

// SwitchUserType switches between passenger and driver
func (h *UserHandler) SwitchUserType(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		UserType string `json:"user_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	profile, err := h.userService.SwitchUserType(userID, req.UserType)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, profile)
}

// CreateVehicle creates a new vehicle for a driver
func (h *UserHandler) CreateVehicle(c *gin.Context) {
	ownerIDStr := c.Param("id")
	ownerID, err := uuid.Parse(ownerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid owner ID"})
		return
	}

	var req application.CreateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	vehicle, err := h.userService.CreateVehicle(ownerID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, vehicle)
}

// GetUserVehicles gets all vehicles for a user
func (h *UserHandler) GetUserVehicles(c *gin.Context) {
	ownerIDStr := c.Param("id")
	ownerID, err := uuid.Parse(ownerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid owner ID"})
		return
	}

	vehicles, err := h.userService.GetUserVehicles(ownerID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, vehicles)
}

// UpdateVehicle updates vehicle information
func (h *UserHandler) UpdateVehicle(c *gin.Context) {
	vehicleIDStr := c.Param("vehicle_id")
	vehicleID, err := uuid.Parse(vehicleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	var req application.UpdateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	vehicle, err := h.userService.UpdateVehicle(vehicleID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, vehicle)
}

// CreateRating creates a new rating
func (h *UserHandler) CreateRating(c *gin.Context) {
	// Get rater ID from authentication context (would be set by auth middleware)
	raterIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	raterID, err := uuid.Parse(raterIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rater ID"})
		return
	}

	var req application.CreateRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	rating, err := h.userService.CreateRating(raterID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, rating)
}

// GetUserRatings gets ratings for a user
func (h *UserHandler) GetUserRatings(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100 // Cap at 100 for performance
	}

	ratings, err := h.userService.GetUserRatings(userID, limit, offset)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ratings": ratings,
		"limit":   limit,
		"offset":  offset,
	})
}

// HandleError handles different types of errors
func (h *UserHandler) handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *shared.ValidationError:
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Message()})
	case *shared.NotFoundError:
		c.JSON(http.StatusNotFound, gin.H{"error": e.Message()})
	case *shared.ConflictError:
		c.JSON(http.StatusConflict, gin.H{"error": e.Message()})
	case *shared.DatabaseError:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
	case *shared.InternalError:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unknown error"})
	}
}