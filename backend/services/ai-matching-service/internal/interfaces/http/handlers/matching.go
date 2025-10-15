package handlers


import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/southern-martin/zride/backend/services/ai-matching-service/internal/application/usecases"
	"github.com/southern-martin/zride/backend/shared/errors"
	"github.com/southern-martin/zride/backend/shared/logger"
)

type MatchingHandler struct {
	matchingUseCase *usecases.MatchingUseCase
	logger          logger.Logger
}

func NewMatchingHandler(matchingUseCase *usecases.MatchingUseCase, logger logger.Logger) *MatchingHandler {
	return &MatchingHandler{
		matchingUseCase: matchingUseCase,
		logger:          logger,
	}
}

// CreateMatchRequest creates a new match request
// @Summary Create match request
// @Description Create a new match request for passenger
// @Tags matching
// @Accept json
// @Produce json
// @Param request body usecases.CreateMatchRequestDTO true "Match request"
// @Success 201 {object} usecases.MatchRequestResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/matching/requests [post]
func (h *MatchingHandler) CreateMatchRequest(c *gin.Context) {
	var req usecases.CreateMatchRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, errors.NewAppError(errors.CodeValidationError, "invalid request body", err))
		return
	}

	// Set defaults
	if req.MaxWaitTime == 0 {
		req.MaxWaitTime = 10 * time.Minute
	}
	if req.MaxDistance == 0 {
		req.MaxDistance = 15.0 // 15km
	}

	response, err := h.matchingUseCase.CreateMatchRequest(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.logger.Info(c.Request.Context(), "Match request created successfully", map[string]interface{}{
		"match_request_id": response.ID,
		"passenger_id":     response.PassengerID,
	})

	c.JSON(http.StatusCreated, response)
}

// AcceptMatch handles driver accepting a match
// @Summary Accept match
// @Description Driver accepts a match request
// @Tags matching
// @Accept json
// @Produce json
// @Param driverID path string true "Driver ID"
// @Param matchResultID path string true "Match Result ID"
// @Success 200 {object} usecases.AcceptMatchResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/matching/drivers/{driverID}/matches/{matchResultID}/accept [post]
func (h *MatchingHandler) AcceptMatch(c *gin.Context) {
	driverIDStr := c.Param("driverID")
	matchResultIDStr := c.Param("matchResultID")

	driverID, err := uuid.Parse(driverIDStr)
	if err != nil {
		h.handleError(c, errors.NewAppError(errors.CodeValidationError, "invalid driver ID", err))
		return
	}

	matchResultID, err := uuid.Parse(matchResultIDStr)
	if err != nil {
		h.handleError(c, errors.NewAppError(errors.CodeValidationError, "invalid match result ID", err))
		return
	}

	response, err := h.matchingUseCase.AcceptMatch(c.Request.Context(), driverID, matchResultID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.logger.Info(c.Request.Context(), "Match accepted successfully", map[string]interface{}{
		"match_result_id": matchResultID,
		"driver_id":       driverID,
		"trip_id":         response.TripID,
	})

	c.JSON(http.StatusOK, response)
}

// GetMatchRequests gets match requests for a passenger
// @Summary Get match requests
// @Description Get match requests for a passenger
// @Tags matching
// @Produce json
// @Param passengerID path string true "Passenger ID"
// @Success 200 {array} usecases.MatchRequestResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/matching/passengers/{passengerID}/requests [get]
func (h *MatchingHandler) GetMatchRequests(c *gin.Context) {
	passengerIDStr := c.Param("passengerID")

	passengerID, err := uuid.Parse(passengerIDStr)
	if err != nil {
		h.handleError(c, errors.NewAppError(errors.CodeValidationError, "invalid passenger ID", err))
		return
	}

	responses, err := h.matchingUseCase.GetMatchRequests(c.Request.Context(), passengerID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, responses)
}

// GetAvailableMatches gets available matches for a driver
// @Summary Get available matches
// @Description Get available matches for a driver
// @Tags matching
// @Produce json
// @Param driverID path string true "Driver ID"
// @Param limit query int false "Limit results" default(10)
// @Success 200 {array} usecases.MatchResultResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/matching/drivers/{driverID}/matches [get]
func (h *MatchingHandler) GetAvailableMatches(c *gin.Context) {
	driverIDStr := c.Param("driverID")

	driverID, err := uuid.Parse(driverIDStr)
	if err != nil {
		h.handleError(c, errors.NewAppError(errors.CodeValidationError, "invalid driver ID", err))
		return
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	responses, err := h.matchingUseCase.GetAvailableMatches(c.Request.Context(), driverID, limit)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, responses)
}

// HealthCheck checks the health of the service
// @Summary Health check
// @Description Check the health of the AI matching service
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/health [get]
func (h *MatchingHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "ai-matching-service",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// handleError handles errors and returns appropriate HTTP responses
func (h *MatchingHandler) handleError(c *gin.Context, err error) {
	if appErr := errors.GetAppError(err); appErr != nil {
		h.logger.Error(c.Request.Context(), "Application error", err, map[string]interface{}{
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
			"code":   appErr.Code(),
		})

		c.JSON(appErr.HTTPStatusCode(), ErrorResponse{
			Error:   appErr.Error(),
			Code:    appErr.Code(),
			Message: appErr.Message(),
		})
		return
	}

	// Generic error
	h.logger.Error(c.Request.Context(), "Unexpected error", err, map[string]interface{}{
		"path":   c.Request.URL.Path,
		"method": c.Request.Method,
	})

	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: "internal server error",
		Code:  errors.CodeInternalError,
	})
}