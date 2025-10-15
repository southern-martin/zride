package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/southern-martin/zride/backend/services/ai-matching-service/internal/domain"
	"github.com/southern-martin/zride/backend/shared/errors"
	"github.com/southern-martin/zride/backend/shared/logger"
)

type MatchingUseCase struct {
	matchRequestRepo  domain.MatchRequestRepository
	driverRepo       domain.DriverRepository
	matchResultRepo  domain.MatchResultRepository
	cacheRepo        domain.CacheRepository
	matchingService  domain.MatchingService
	notificationSvc  domain.NotificationService
	externalSvc      domain.ExternalService
	config          domain.MatchingConfig
	logger          logger.Logger
}

func NewMatchingUseCase(
	matchRequestRepo domain.MatchRequestRepository,
	driverRepo domain.DriverRepository,
	matchResultRepo domain.MatchResultRepository,
	cacheRepo domain.CacheRepository,
	matchingService domain.MatchingService,
	notificationSvc domain.NotificationService,
	externalSvc domain.ExternalService,
	config domain.MatchingConfig,
	logger logger.Logger,
) *MatchingUseCase {
	return &MatchingUseCase{
		matchRequestRepo: matchRequestRepo,
		driverRepo:       driverRepo,
		matchResultRepo:  matchResultRepo,
		cacheRepo:        cacheRepo,
		matchingService:  matchingService,
		notificationSvc:  notificationSvc,
		externalSvc:      externalSvc,
		config:           config,
		logger:           logger,
	}
}

// CreateMatchRequest creates a new match request and finds matches
func (uc *MatchingUseCase) CreateMatchRequest(ctx context.Context, req CreateMatchRequestDTO) (*MatchRequestResponse, error) {
	// Validate user exists
	_, err := uc.externalSvc.GetUserByID(ctx, req.PassengerID)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeUserNotFound, "passenger not found", err)
	}

	// Create match request
	matchRequest := &domain.MatchRequest{
		ID:               uuid.New(),
		PassengerID:      req.PassengerID,
		PickupLocation:   req.PickupLocation,
		DropoffLocation:  req.DropoffLocation,
		RequestTime:      time.Now(),
		MaxWaitTime:      req.MaxWaitTime,
		PreferredCarType: req.PreferredCarType,
		MaxDistance:      req.MaxDistance,
		PriceRange:       req.PriceRange,
		Status:           domain.MatchStatusPending,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := uc.matchRequestRepo.Create(ctx, matchRequest); err != nil {
		return nil, errors.NewAppError(errors.CodeInternalError, "failed to create match request", err)
	}

	uc.logger.Info(ctx, "Match request created", map[string]interface{}{
		"match_request_id": matchRequest.ID,
		"passenger_id":     matchRequest.PassengerID,
	})

	// Find matches asynchronously
	go uc.processMatchRequest(context.Background(), matchRequest)

	return &MatchRequestResponse{
		ID:               matchRequest.ID,
		PassengerID:      matchRequest.PassengerID,
		PickupLocation:   matchRequest.PickupLocation,
		DropoffLocation:  matchRequest.DropoffLocation,
		Status:           string(matchRequest.Status),
		CreatedAt:        matchRequest.CreatedAt,
		EstimatedMatches: 0, // Will be updated when matches are found
	}, nil
}

// processMatchRequest processes the match request and finds suitable drivers
func (uc *MatchingUseCase) processMatchRequest(ctx context.Context, request *domain.MatchRequest) {
	// Find matches using the configured algorithm
	matches, err := uc.matchingService.FindMatches(ctx, request, uc.config)
	if err != nil {
		uc.logger.Error(ctx, "Failed to find matches", err, map[string]interface{}{
			"match_request_id": request.ID,
		})
		return
	}

	if len(matches) == 0 {
		// No matches found - update status and notify
		request.Status = domain.MatchStatusExpired
		if err := uc.matchRequestRepo.Update(ctx, request); err != nil {
			uc.logger.Error(ctx, "Failed to update match request status", err, map[string]interface{}{
				"match_request_id": request.ID,
			})
		}
		
		// Notify passenger about no matches
		if err := uc.notificationSvc.NotifyMatchTimeout(ctx, request.PassengerID, request); err != nil {
			uc.logger.Error(ctx, "Failed to send timeout notification", err, map[string]interface{}{
				"match_request_id": request.ID,
			})
		}
		return
	}

	// Save match results
	for _, match := range matches {
		if err := uc.matchResultRepo.Create(ctx, match); err != nil {
			uc.logger.Error(ctx, "Failed to save match result", err, map[string]interface{}{
				"match_result_id": match.ID,
			})
			continue
		}

		// Notify driver about potential match
		if err := uc.notificationSvc.NotifyDriverMatch(ctx, match.DriverID, match); err != nil {
			uc.logger.Error(ctx, "Failed to notify driver", err, map[string]interface{}{
				"driver_id": match.DriverID,
				"match_result_id": match.ID,
			})
		}
	}

	// Update match request status
	request.Status = domain.MatchStatusMatched
	request.UpdatedAt = time.Now()
	if err := uc.matchRequestRepo.Update(ctx, request); err != nil {
		uc.logger.Error(ctx, "Failed to update match request", err, map[string]interface{}{
			"match_request_id": request.ID,
		})
	}

	uc.logger.Info(ctx, "Match processing completed", map[string]interface{}{
		"match_request_id": request.ID,
		"matches_found":    len(matches),
	})
}

// AcceptMatch handles driver accepting a match
func (uc *MatchingUseCase) AcceptMatch(ctx context.Context, driverID uuid.UUID, matchResultID uuid.UUID) (*AcceptMatchResponse, error) {
	// Get match result
	matchResult, err := uc.matchResultRepo.GetByID(ctx, matchResultID)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeNotFound, "match result not found", err)
	}

	// Verify driver
	if matchResult.DriverID != driverID {
		return nil, errors.NewAppError(errors.CodeUnauthorized, "unauthorized to accept this match", nil)
	}

	// Check if match is still pending
	if matchResult.Status != domain.MatchResultStatusPending {
		return nil, errors.NewAppError(errors.CodeInvalidOperation, "match is no longer available", nil)
	}

	// Get match request
	matchRequest, err := uc.matchRequestRepo.GetByID(ctx, matchResult.MatchRequestID)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeNotFound, "match request not found", err)
	}

	// Create trip in trip service
	tripData := &domain.TripData{
		PassengerID:       matchRequest.PassengerID,
		DriverID:          driverID,
		PickupLocation:    matchRequest.PickupLocation,
		DropoffLocation:   matchRequest.DropoffLocation,
		EstimatedPrice:    matchResult.EstimatedPrice,
		EstimatedTime:     int(matchResult.EstimatedTime.Minutes()),
		EstimatedDistance: matchResult.EstimatedDistance,
	}

	trip, err := uc.externalSvc.CreateTrip(ctx, tripData)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeInternalError, "failed to create trip", err)
	}

	// Update match result status
	matchResult.Status = domain.MatchResultStatusAccepted
	if err := uc.matchResultRepo.Update(ctx, matchResult); err != nil {
		return nil, errors.NewAppError(errors.CodeInternalError, "failed to update match result", err)
	}

	// Update driver availability
	if err := uc.driverRepo.UpdateAvailability(ctx, driverID, false); err != nil {
		uc.logger.Error(ctx, "Failed to update driver availability", err, map[string]interface{}{
			"driver_id": driverID,
		})
	}

	// Notify passenger about accepted match
	if err := uc.notificationSvc.NotifyPassengerMatch(ctx, matchRequest.PassengerID, matchResult); err != nil {
		uc.logger.Error(ctx, "Failed to notify passenger", err, map[string]interface{}{
			"passenger_id": matchRequest.PassengerID,
		})
	}

	uc.logger.Info(ctx, "Match accepted successfully", map[string]interface{}{
		"match_result_id": matchResultID,
		"driver_id":       driverID,
		"trip_id":         trip.ID,
	})

	return &AcceptMatchResponse{
		MatchResultID: matchResult.ID,
		TripID:        trip.ID,
		PassengerID:   matchRequest.PassengerID,
		DriverID:      driverID,
		EstimatedTime: matchResult.EstimatedTime,
		EstimatedPrice: matchResult.EstimatedPrice,
		Status:        "accepted",
	}, nil
}

// GetMatchRequests gets match requests for a passenger
func (uc *MatchingUseCase) GetMatchRequests(ctx context.Context, passengerID uuid.UUID) ([]*MatchRequestResponse, error) {
	requests, err := uc.matchRequestRepo.GetByPassengerID(ctx, passengerID)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeInternalError, "failed to get match requests", err)
	}

	var responses []*MatchRequestResponse
	for _, req := range requests {
		// Get match results count
		matches, _ := uc.matchResultRepo.GetByMatchRequestID(ctx, req.ID)
		
		responses = append(responses, &MatchRequestResponse{
			ID:               req.ID,
			PassengerID:      req.PassengerID,
			PickupLocation:   req.PickupLocation,
			DropoffLocation:  req.DropoffLocation,
			Status:           string(req.Status),
			CreatedAt:        req.CreatedAt,
			EstimatedMatches: len(matches),
		})
	}

	return responses, nil
}

// GetAvailableMatches gets available matches for a driver
func (uc *MatchingUseCase) GetAvailableMatches(ctx context.Context, driverID uuid.UUID, limit int) ([]*MatchResultResponse, error) {
	driver, err := uc.driverRepo.GetByID(ctx, driverID)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeNotFound, "driver not found", err)
	}

	if !driver.IsAvailable {
		return []*MatchResultResponse{}, nil
	}

	// Get pending requests in driver's area
	pendingRequests, err := uc.matchRequestRepo.GetPendingRequests(ctx, limit*2)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeInternalError, "failed to get pending requests", err)
	}

	var responses []*MatchResultResponse
	for _, req := range pendingRequests {
		// Check if driver is within range
		distance := uc.matchingService.CalculateDistance(driver.CurrentLocation, req.PickupLocation)
		if distance > driver.MaxDistance {
			continue
		}

		// Calculate match score
		score, err := uc.matchingService.ScoreMatch(ctx, req, driver, uc.config.Criteria)
		if err != nil {
			continue
		}

		estimatedTime, _ := uc.matchingService.EstimateTime(driver.CurrentLocation, req.PickupLocation)
		estimatedPrice, _ := uc.matchingService.EstimatePrice(req.PickupLocation, req.DropoffLocation, driver.CarType)

		responses = append(responses, &MatchResultResponse{
			ID:                uuid.New(),
			MatchRequestID:    req.ID,
			DriverID:          driverID,
			PassengerID:       req.PassengerID,
			PickupLocation:    req.PickupLocation,
			DropoffLocation:   req.DropoffLocation,
			Score:             score,
			EstimatedDistance: distance,
			EstimatedTime:     estimatedTime,
			EstimatedPrice:    estimatedPrice,
			Status:            "pending",
		})

		if len(responses) >= limit {
			break
		}
	}

	return responses, nil
}