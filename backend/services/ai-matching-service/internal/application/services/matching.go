package services


import (
	"context"
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/southern-martin/zride/backend/services/ai-matching-service/internal/domain"
	"github.com/southern-martin/zride/backend/shared/errors"
	"github.com/southern-martin/zride/backend/shared/logger"
)

const (
	// Earth's radius in kilometers
	EarthRadiusKM = 6371.0
	
	// Base pricing per km in VND
	BasePricePerKM = 8000.0
	
	// Car type multipliers
	CarTypeMotorbike = 1.0
	CarTypeCar4Seat  = 1.5
	CarTypeCar7Seat  = 2.0
	
	// Average speed in km/h for time estimation
	AverageSpeedKMH = 25.0
)

// AIMatchingService implements domain.MatchingService with AI algorithms
type AIMatchingService struct {
	driverRepo domain.DriverRepository
	logger     logger.Logger
}

// NewAIMatchingService creates a new AI matching service
func NewAIMatchingService(driverRepo domain.DriverRepository, logger logger.Logger) domain.MatchingService {
	return &AIMatchingService{
		driverRepo: driverRepo,
		logger:     logger,
	}
}

// FindMatches finds suitable drivers for a match request using AI algorithms
func (s *AIMatchingService) FindMatches(ctx context.Context, request *domain.MatchRequest, config domain.MatchingConfig) ([]*domain.MatchResult, error) {
	// Get available drivers in the search radius
	drivers, err := s.driverRepo.GetAvailableDriversInRadius(
		ctx, 
		request.PickupLocation, 
		config.MaxSearchRadius,
	)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeInternalError, "failed to get available drivers", err)
	}

	// Filter drivers based on criteria
	filteredDrivers := s.filterDrivers(drivers, request, config)
	
	if len(filteredDrivers) == 0 {
		return []*domain.MatchResult{}, nil
	}

	// Apply matching algorithm
	var matches []*domain.MatchResult
	
	switch config.Algorithm {
	case domain.AlgorithmNearest:
		matches = s.nearestAlgorithm(ctx, request, filteredDrivers, config)
	case domain.AlgorithmWeighted:
		matches = s.weightedAlgorithm(ctx, request, filteredDrivers, config)
	case domain.AlgorithmML:
		matches = s.mlAlgorithm(ctx, request, filteredDrivers, config)
	case domain.AlgorithmHybrid:
		matches = s.hybridAlgorithm(ctx, request, filteredDrivers, config)
	default:
		matches = s.weightedAlgorithm(ctx, request, filteredDrivers, config)
	}

	// Limit results
	if len(matches) > config.MaxDrivers {
		matches = matches[:config.MaxDrivers]
	}

	s.logger.Info(ctx, "Matches found", map[string]interface{}{
		"request_id":     request.ID,
		"total_drivers":  len(drivers),
		"filtered_drivers": len(filteredDrivers),
		"matches_found":  len(matches),
		"algorithm":      config.Algorithm,
	})

	return matches, nil
}

// filterDrivers filters drivers based on basic criteria
func (s *AIMatchingService) filterDrivers(drivers []*domain.Driver, request *domain.MatchRequest, config domain.MatchingConfig) []*domain.Driver {
	var filtered []*domain.Driver
	
	for _, driver := range drivers {
		// Check if driver meets minimum rating
		if driver.Rating < config.MinDriverRating {
			continue
		}
		
		// Check if driver accepts the car type preference
		if request.PreferredCarType != "" && driver.CarType != request.PreferredCarType {
			continue
		}
		
		// Check if driver is within maximum distance
		distance := s.CalculateDistance(driver.CurrentLocation, request.PickupLocation)
		if distance > request.MaxDistance {
			continue
		}
		
		// Check if driver was active recently (within last 5 minutes)
		if time.Since(driver.LastActiveTime) > 5*time.Minute {
			continue
		}
		
		filtered = append(filtered, driver)
	}
	
	return filtered
}

// nearestAlgorithm implements nearest driver matching
func (s *AIMatchingService) nearestAlgorithm(ctx context.Context, request *domain.MatchRequest, drivers []*domain.Driver, config domain.MatchingConfig) []*domain.MatchResult {
	type driverDistance struct {
		driver   *domain.Driver
		distance float64
	}
	
	var driverDistances []driverDistance
	for _, driver := range drivers {
		distance := s.CalculateDistance(driver.CurrentLocation, request.PickupLocation)
		driverDistances = append(driverDistances, driverDistance{
			driver:   driver,
			distance: distance,
		})
	}
	
	// Sort by distance
	sort.Slice(driverDistances, func(i, j int) bool {
		return driverDistances[i].distance < driverDistances[j].distance
	})
	
	var results []*domain.MatchResult
	for _, dd := range driverDistances {
		estimatedTime, _ := s.EstimateTime(dd.driver.CurrentLocation, request.PickupLocation)
		estimatedPrice, _ := s.EstimatePrice(request.PickupLocation, request.DropoffLocation, dd.driver.CarType)
		
		result := &domain.MatchResult{
			ID:                uuid.New(),
			MatchRequestID:    request.ID,
			DriverID:          dd.driver.ID,
			Score:             1.0 / (1.0 + dd.distance), // Higher score for shorter distance
			EstimatedDistance: dd.distance,
			EstimatedTime:     estimatedTime,
			EstimatedPrice:    estimatedPrice,
			MatchTime:         time.Now(),
			Status:            domain.MatchResultStatusPending,
			CreatedAt:         time.Now(),
		}
		
		results = append(results, result)
	}
	
	return results
}

// weightedAlgorithm implements weighted scoring algorithm
func (s *AIMatchingService) weightedAlgorithm(ctx context.Context, request *domain.MatchRequest, drivers []*domain.Driver, config domain.MatchingConfig) []*domain.MatchResult {
	type driverScore struct {
		driver *domain.Driver
		score  float64
		result *domain.MatchResult
	}
	
	var driverScores []driverScore
	for _, driver := range drivers {
		score, err := s.ScoreMatch(ctx, request, driver, config.Criteria)
		if err != nil {
			continue
		}
		
		distance := s.CalculateDistance(driver.CurrentLocation, request.PickupLocation)
		estimatedTime, _ := s.EstimateTime(driver.CurrentLocation, request.PickupLocation)
		estimatedPrice, _ := s.EstimatePrice(request.PickupLocation, request.DropoffLocation, driver.CarType)
		
		result := &domain.MatchResult{
			ID:                uuid.New(),
			MatchRequestID:    request.ID,
			DriverID:          driver.ID,
			Score:             score,
			EstimatedDistance: distance,
			EstimatedTime:     estimatedTime,
			EstimatedPrice:    estimatedPrice,
			MatchTime:         time.Now(),
			Status:            domain.MatchResultStatusPending,
			CreatedAt:         time.Now(),
		}
		
		driverScores = append(driverScores, driverScore{
			driver: driver,
			score:  score,
			result: result,
		})
	}
	
	// Sort by score (highest first)
	sort.Slice(driverScores, func(i, j int) bool {
		return driverScores[i].score > driverScores[j].score
	})
	
	var results []*domain.MatchResult
	for _, ds := range driverScores {
		results = append(results, ds.result)
	}
	
	return results
}

// driverFeatures represents features extracted for ML algorithm
type driverFeatures struct {
	driver              *domain.Driver
	distance           float64
	normalizedRating   float64
	experienceScore    float64
	availabilityScore  float64
	locationSimilarity float64
	timePreference     float64
	mlScore           float64
	result            *domain.MatchResult
}

// mlAlgorithm implements machine learning based matching
func (s *AIMatchingService) mlAlgorithm(ctx context.Context, request *domain.MatchRequest, drivers []*domain.Driver, config domain.MatchingConfig) []*domain.MatchResult {
	// For now, implement as enhanced weighted algorithm with additional features
	// In production, this would use a trained ML model
	
	var driverFeaturesList []driverFeatures
	
	for _, driver := range drivers {
		distance := s.CalculateDistance(driver.CurrentLocation, request.PickupLocation)
		estimatedTime, _ := s.EstimateTime(driver.CurrentLocation, request.PickupLocation)
		estimatedPrice, _ := s.EstimatePrice(request.PickupLocation, request.DropoffLocation, driver.CarType)
		
		// Extract features
		features := driverFeatures{
			driver:              driver,
			distance:           distance,
			normalizedRating:   driver.Rating / 5.0, // Normalize rating to 0-1
			experienceScore:    math.Min(float64(driver.CompletedTrips)/100.0, 1.0), // Cap at 100 trips
			availabilityScore:  1.0 - (time.Since(driver.LastActiveTime).Minutes() / 60.0), // Recent activity score
			locationSimilarity: s.calculateLocationSimilarity(driver, request),
			timePreference:     s.calculateTimePreference(request),
		}
		
		// Simple ML scoring (in production, use trained model)
		features.mlScore = s.calculateMLScore(features, config.Criteria)
		
		features.result = &domain.MatchResult{
			ID:                uuid.New(),
			MatchRequestID:    request.ID,
			DriverID:          driver.ID,
			Score:             features.mlScore,
			EstimatedDistance: distance,
			EstimatedTime:     estimatedTime,
			EstimatedPrice:    estimatedPrice,
			MatchTime:         time.Now(),
			Status:            domain.MatchResultStatusPending,
			CreatedAt:         time.Now(),
		}
		
		driverFeaturesList = append(driverFeaturesList, features)
	}
	
	// Sort by ML score
	sort.Slice(driverFeaturesList, func(i, j int) bool {
		return driverFeaturesList[i].mlScore > driverFeaturesList[j].mlScore
	})
	
	var results []*domain.MatchResult
	for _, df := range driverFeaturesList {
		results = append(results, df.result)
	}
	
	return results
}

// hybridAlgorithm combines multiple algorithms
func (s *AIMatchingService) hybridAlgorithm(ctx context.Context, request *domain.MatchRequest, drivers []*domain.Driver, config domain.MatchingConfig) []*domain.MatchResult {
	// Get results from different algorithms
	nearestResults := s.nearestAlgorithm(ctx, request, drivers, config)
	weightedResults := s.weightedAlgorithm(ctx, request, drivers, config)
	mlResults := s.mlAlgorithm(ctx, request, drivers, config)
	
	// Combine scores with weights
	driverScores := make(map[uuid.UUID]float64)
	driverResults := make(map[uuid.UUID]*domain.MatchResult)
	
	// Weight: 30% nearest, 40% weighted, 30% ML
	for _, result := range nearestResults {
		driverScores[result.DriverID] += result.Score * 0.3
		driverResults[result.DriverID] = result
	}
	
	for _, result := range weightedResults {
		driverScores[result.DriverID] += result.Score * 0.4
		driverResults[result.DriverID] = result
	}
	
	for _, result := range mlResults {
		driverScores[result.DriverID] += result.Score * 0.3
		driverResults[result.DriverID] = result
	}
	
	// Sort by combined score
	type hybridScore struct {
		driverID uuid.UUID
		score    float64
	}
	
	var hybridScores []hybridScore
	for driverID, score := range driverScores {
		hybridScores = append(hybridScores, hybridScore{
			driverID: driverID,
			score:    score,
		})
	}
	
	sort.Slice(hybridScores, func(i, j int) bool {
		return hybridScores[i].score > hybridScores[j].score
	})
	
	var results []*domain.MatchResult
	for _, hs := range hybridScores {
		if result, exists := driverResults[hs.driverID]; exists {
			result.Score = hs.score // Update with combined score
			results = append(results, result)
		}
	}
	
	return results
}

// ScoreMatch calculates matching score between request and driver
func (s *AIMatchingService) ScoreMatch(ctx context.Context, request *domain.MatchRequest, driver *domain.Driver, criteria domain.MatchingCriteria) (float64, error) {
	// Calculate distance score (closer is better)
	distance := s.CalculateDistance(driver.CurrentLocation, request.PickupLocation)
	distanceScore := 1.0 / (1.0 + distance/10.0) // Normalize with 10km factor
	
	// Calculate rating score (higher is better)
	ratingScore := driver.Rating / 5.0 // Normalize to 0-1
	
	// Calculate time score (faster pickup is better)
	estimatedTime, _ := s.EstimateTime(driver.CurrentLocation, request.PickupLocation)
	timeScore := 1.0 / (1.0 + estimatedTime.Minutes()/30.0) // Normalize with 30min factor
	
	// Calculate price score (within range is better)
	estimatedPrice, _ := s.EstimatePrice(request.PickupLocation, request.DropoffLocation, driver.CarType)
	priceScore := 1.0
	if request.PriceRange.MaxPrice > 0 {
		if estimatedPrice <= request.PriceRange.MaxPrice {
			priceScore = 1.0 - (estimatedPrice-request.PriceRange.MinPrice)/(request.PriceRange.MaxPrice-request.PriceRange.MinPrice)
		} else {
			priceScore = 0.5 // Penalty for exceeding max price
		}
	}
	
	// Calculate experience score
	experienceScore := math.Min(float64(driver.CompletedTrips)/50.0, 1.0) // Normalize to 50 trips max
	
	// Weighted final score
	finalScore := (distanceScore * criteria.DistanceWeight) +
		(ratingScore * criteria.RatingWeight) +
		(timeScore * criteria.TimeWeight) +
		(priceScore * criteria.PriceWeight) +
		(experienceScore * criteria.ExperienceWeight)
	
	return finalScore, nil
}

// CalculateDistance calculates haversine distance between two locations
func (s *AIMatchingService) CalculateDistance(from, to domain.Location) float64 {
	// Convert latitude and longitude from degrees to radians
	lat1 := from.Latitude * math.Pi / 180
	lon1 := from.Longitude * math.Pi / 180
	lat2 := to.Latitude * math.Pi / 180
	lon2 := to.Longitude * math.Pi / 180
	
	// Calculate differences
	dlat := lat2 - lat1
	dlon := lon2 - lon1
	
	// Haversine formula
	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	
	// Distance in kilometers
	return EarthRadiusKM * c
}

// EstimateTime estimates travel time between two locations
func (s *AIMatchingService) EstimateTime(from, to domain.Location) (time.Duration, error) {
	distance := s.CalculateDistance(from, to)
	
	// Simple time estimation based on average speed
	// In production, integrate with Google Maps or similar service
	hours := distance / AverageSpeedKMH
	
	return time.Duration(hours * float64(time.Hour)), nil
}

// EstimatePrice estimates trip price based on distance and car type
func (s *AIMatchingService) EstimatePrice(from, to domain.Location, carType string) (float64, error) {
	distance := s.CalculateDistance(from, to)
	
	// Base price calculation
	basePrice := distance * BasePricePerKM
	
	// Apply car type multiplier
	var multiplier float64
	switch carType {
	case "motorbike":
		multiplier = CarTypeMotorbike
	case "car_4_seat":
		multiplier = CarTypeCar4Seat
	case "car_7_seat":
		multiplier = CarTypeCar7Seat
	default:
		multiplier = CarTypeCar4Seat // Default to 4-seat car
	}
	
	finalPrice := basePrice * multiplier
	
	// Add minimum fare (20,000 VND)
	if finalPrice < 20000 {
		finalPrice = 20000
	}
	
	return finalPrice, nil
}

// calculateLocationSimilarity calculates similarity between driver's preferred areas and request location
func (s *AIMatchingService) calculateLocationSimilarity(driver *domain.Driver, request *domain.MatchRequest) float64 {
	if len(driver.PreferredAreas) == 0 {
		return 0.5 // Neutral if no preference
	}
	
	maxSimilarity := 0.0
	for _, preferredArea := range driver.PreferredAreas {
		distance := s.CalculateDistance(preferredArea, request.PickupLocation)
		similarity := 1.0 / (1.0 + distance/5.0) // 5km factor
		if similarity > maxSimilarity {
			maxSimilarity = similarity
		}
	}
	
	return maxSimilarity
}

// calculateTimePreference calculates time-based preference score
func (s *AIMatchingService) calculateTimePreference(request *domain.MatchRequest) float64 {
	hour := request.RequestTime.Hour()
	
	// Peak hours get lower preference (more demand)
	if (hour >= 7 && hour <= 9) || (hour >= 17 && hour <= 19) {
		return 0.7
	}
	
	// Off-peak hours get higher preference
	return 1.0
}

// calculateMLScore calculates ML-based score using simple feature combination
func (s *AIMatchingService) calculateMLScore(features driverFeatures, criteria domain.MatchingCriteria) float64 {
	// Simple feature combination (in production, use trained model weights)
	distanceFeature := 1.0 / (1.0 + features.distance/10.0)
	
	// Non-linear combinations for better discrimination
	ratingFeature := math.Pow(features.normalizedRating, 1.5)
	experienceFeature := math.Sqrt(features.experienceScore)
	availabilityFeature := math.Pow(features.availabilityScore, 0.8)
	
	// Feature interactions
	ratingDistanceInteraction := ratingFeature * distanceFeature
	experienceLocationInteraction := experienceFeature * features.locationSimilarity
	
	// Weighted combination with interactions
	score := (distanceFeature * criteria.DistanceWeight) +
		(ratingFeature * criteria.RatingWeight) +
		(experienceFeature * criteria.ExperienceWeight) +
		(availabilityFeature * 0.1) + // Small weight for availability
		(features.locationSimilarity * 0.1) + // Small weight for location similarity
		(features.timePreference * 0.05) + // Small weight for time preference
		(ratingDistanceInteraction * 0.1) + // Interaction term
		(experienceLocationInteraction * 0.05) // Interaction term
	
	// Apply sigmoid for normalization
	return 1.0 / (1.0 + math.Exp(-5.0*(score-0.5)))
}