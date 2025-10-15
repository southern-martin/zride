package postgres


import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/southern-martin/zride/backend/services/ai-matching-service/internal/domain"
	"github.com/southern-martin/zride/backend/shared/errors"
)

type MatchRequestRepository struct {
	db *sql.DB
}

func NewMatchRequestRepository(db *sql.DB) domain.MatchRequestRepository {
	return &MatchRequestRepository{db: db}
}

func (r *MatchRequestRepository) Create(ctx context.Context, request *domain.MatchRequest) error {
	query := `
		INSERT INTO match_requests (
			id, passenger_id, pickup_location, dropoff_location, request_time,
			max_wait_time, preferred_car_type, max_distance, price_range,
			status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	pickupLocationJSON, _ := json.Marshal(request.PickupLocation)
	dropoffLocationJSON, _ := json.Marshal(request.DropoffLocation)
	priceRangeJSON, _ := json.Marshal(request.PriceRange)

	_, err := r.db.ExecContext(ctx, query,
		request.ID,
		request.PassengerID,
		pickupLocationJSON,
		dropoffLocationJSON,
		request.RequestTime,
		request.MaxWaitTime.Nanoseconds(),
		request.PreferredCarType,
		request.MaxDistance,
		priceRangeJSON,
		request.Status,
		request.CreatedAt,
		request.UpdatedAt,
	)

	if err != nil {
		return errors.NewAppError(errors.CodeInternalError, "failed to create match request", err)
	}

	return nil
}

func (r *MatchRequestRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.MatchRequest, error) {
	query := `
		SELECT id, passenger_id, pickup_location, dropoff_location, request_time,
			   max_wait_time, preferred_car_type, max_distance, price_range,
			   status, created_at, updated_at
		FROM match_requests WHERE id = $1
	`

	var request domain.MatchRequest
	var pickupLocationJSON, dropoffLocationJSON, priceRangeJSON []byte
	var maxWaitTimeNanos int64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&request.ID,
		&request.PassengerID,
		&pickupLocationJSON,
		&dropoffLocationJSON,
		&request.RequestTime,
		&maxWaitTimeNanos,
		&request.PreferredCarType,
		&request.MaxDistance,
		&priceRangeJSON,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewAppError(errors.CodeNotFound, "match request not found", err)
		}
		return nil, errors.NewAppError(errors.CodeInternalError, "failed to get match request", err)
	}

	// Parse JSON fields
	json.Unmarshal(pickupLocationJSON, &request.PickupLocation)
	json.Unmarshal(dropoffLocationJSON, &request.DropoffLocation)
	json.Unmarshal(priceRangeJSON, &request.PriceRange)
	request.MaxWaitTime = time.Duration(maxWaitTimeNanos)

	return &request, nil
}

func (r *MatchRequestRepository) GetByPassengerID(ctx context.Context, passengerID uuid.UUID) ([]*domain.MatchRequest, error) {
	query := `
		SELECT id, passenger_id, pickup_location, dropoff_location, request_time,
			   max_wait_time, preferred_car_type, max_distance, price_range,
			   status, created_at, updated_at
		FROM match_requests 
		WHERE passenger_id = $1 
		ORDER BY created_at DESC
		LIMIT 50
	`

	rows, err := r.db.QueryContext(ctx, query, passengerID)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeInternalError, "failed to query match requests", err)
	}
	defer rows.Close()

	var requests []*domain.MatchRequest
	for rows.Next() {
		var request domain.MatchRequest
		var pickupLocationJSON, dropoffLocationJSON, priceRangeJSON []byte
		var maxWaitTimeNanos int64

		err := rows.Scan(
			&request.ID,
			&request.PassengerID,
			&pickupLocationJSON,
			&dropoffLocationJSON,
			&request.RequestTime,
			&maxWaitTimeNanos,
			&request.PreferredCarType,
			&request.MaxDistance,
			&priceRangeJSON,
			&request.Status,
			&request.CreatedAt,
			&request.UpdatedAt,
		)

		if err != nil {
			continue
		}

		// Parse JSON fields
		json.Unmarshal(pickupLocationJSON, &request.PickupLocation)
		json.Unmarshal(dropoffLocationJSON, &request.DropoffLocation)
		json.Unmarshal(priceRangeJSON, &request.PriceRange)
		request.MaxWaitTime = time.Duration(maxWaitTimeNanos)

		requests = append(requests, &request)
	}

	return requests, nil
}

func (r *MatchRequestRepository) Update(ctx context.Context, request *domain.MatchRequest) error {
	query := `
		UPDATE match_requests SET
			pickup_location = $2, dropoff_location = $3, request_time = $4,
			max_wait_time = $5, preferred_car_type = $6, max_distance = $7,
			price_range = $8, status = $9, updated_at = $10
		WHERE id = $1
	`

	pickupLocationJSON, _ := json.Marshal(request.PickupLocation)
	dropoffLocationJSON, _ := json.Marshal(request.DropoffLocation)
	priceRangeJSON, _ := json.Marshal(request.PriceRange)
	request.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		request.ID,
		pickupLocationJSON,
		dropoffLocationJSON,
		request.RequestTime,
		request.MaxWaitTime.Nanoseconds(),
		request.PreferredCarType,
		request.MaxDistance,
		priceRangeJSON,
		request.Status,
		request.UpdatedAt,
	)

	if err != nil {
		return errors.NewAppError(errors.CodeInternalError, "failed to update match request", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.NewAppError(errors.CodeNotFound, "match request not found", nil)
	}

	return nil
}

func (r *MatchRequestRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM match_requests WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.NewAppError(errors.CodeInternalError, "failed to delete match request", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.NewAppError(errors.CodeNotFound, "match request not found", nil)
	}

	return nil
}

func (r *MatchRequestRepository) GetPendingRequests(ctx context.Context, limit int) ([]*domain.MatchRequest, error) {
	query := `
		SELECT id, passenger_id, pickup_location, dropoff_location, request_time,
			   max_wait_time, preferred_car_type, max_distance, price_range,
			   status, created_at, updated_at
		FROM match_requests 
		WHERE status = $1 
		ORDER BY created_at ASC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, domain.MatchStatusPending, limit)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeInternalError, "failed to query pending requests", err)
	}
	defer rows.Close()

	var requests []*domain.MatchRequest
	for rows.Next() {
		var request domain.MatchRequest
		var pickupLocationJSON, dropoffLocationJSON, priceRangeJSON []byte
		var maxWaitTimeNanos int64

		err := rows.Scan(
			&request.ID,
			&request.PassengerID,
			&pickupLocationJSON,
			&dropoffLocationJSON,
			&request.RequestTime,
			&maxWaitTimeNanos,
			&request.PreferredCarType,
			&request.MaxDistance,
			&priceRangeJSON,
			&request.Status,
			&request.CreatedAt,
			&request.UpdatedAt,
		)

		if err != nil {
			continue
		}

		// Parse JSON fields
		json.Unmarshal(pickupLocationJSON, &request.PickupLocation)
		json.Unmarshal(dropoffLocationJSON, &request.DropoffLocation)
		json.Unmarshal(priceRangeJSON, &request.PriceRange)
		request.MaxWaitTime = time.Duration(maxWaitTimeNanos)

		requests = append(requests, &request)
	}

	return requests, nil
}

// DriverRepository implementation
type DriverRepository struct {
	db *sql.DB
}

func NewDriverRepository(db *sql.DB) domain.DriverRepository {
	return &DriverRepository{db: db}
}

func (r *DriverRepository) GetAvailableDriversInRadius(ctx context.Context, location domain.Location, radius float64) ([]*domain.Driver, error) {
	query := `
		SELECT d.id, d.user_id, d.current_location, d.is_available, d.car_type,
			   d.rating, d.completed_trips, d.last_active_time, d.max_distance,
			   d.preferred_areas,
			   ST_Distance(
				   ST_SetSRID(ST_MakePoint($2, $1), 4326)::geography,
				   ST_SetSRID(ST_MakePoint(
					   (d.current_location->>'longitude')::float,
					   (d.current_location->>'latitude')::float
				   ), 4326)::geography
			   ) / 1000 as distance_km
		FROM drivers d
		WHERE d.is_available = true
		  AND ST_DWithin(
			  ST_SetSRID(ST_MakePoint($2, $1), 4326)::geography,
			  ST_SetSRID(ST_MakePoint(
				  (d.current_location->>'longitude')::float,
				  (d.current_location->>'latitude')::float
			  ), 4326)::geography,
			  $3 * 1000
		  )
		ORDER BY distance_km
		LIMIT 50
	`

	rows, err := r.db.QueryContext(ctx, query, location.Latitude, location.Longitude, radius)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeInternalError, "failed to query drivers", err)
	}
	defer rows.Close()

	var drivers []*domain.Driver
	for rows.Next() {
		var driver domain.Driver
		var currentLocationJSON, preferredAreasJSON []byte
		var distanceKM float64

		err := rows.Scan(
			&driver.ID,
			&driver.UserID,
			&currentLocationJSON,
			&driver.IsAvailable,
			&driver.CarType,
			&driver.Rating,
			&driver.CompletedTrips,
			&driver.LastActiveTime,
			&driver.MaxDistance,
			&preferredAreasJSON,
			&distanceKM,
		)

		if err != nil {
			continue
		}

		// Parse JSON fields
		json.Unmarshal(currentLocationJSON, &driver.CurrentLocation)
		json.Unmarshal(preferredAreasJSON, &driver.PreferredAreas)

		drivers = append(drivers, &driver)
	}

	return drivers, nil
}

func (r *DriverRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Driver, error) {
	query := `
		SELECT id, user_id, current_location, is_available, car_type,
			   rating, completed_trips, last_active_time, max_distance,
			   preferred_areas
		FROM drivers WHERE id = $1
	`

	var driver domain.Driver
	var currentLocationJSON, preferredAreasJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&driver.ID,
		&driver.UserID,
		&currentLocationJSON,
		&driver.IsAvailable,
		&driver.CarType,
		&driver.Rating,
		&driver.CompletedTrips,
		&driver.LastActiveTime,
		&driver.MaxDistance,
		&preferredAreasJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewAppError(errors.CodeNotFound, "driver not found", err)
		}
		return nil, errors.NewAppError(errors.CodeInternalError, "failed to get driver", err)
	}

	// Parse JSON fields
	json.Unmarshal(currentLocationJSON, &driver.CurrentLocation)
	json.Unmarshal(preferredAreasJSON, &driver.PreferredAreas)

	return &driver, nil
}

func (r *DriverRepository) UpdateLocation(ctx context.Context, driverID uuid.UUID, location domain.Location) error {
	locationJSON, _ := json.Marshal(location)
	query := `
		UPDATE drivers SET 
			current_location = $2, 
			last_active_time = $3 
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, driverID, locationJSON, time.Now())
	if err != nil {
		return errors.NewAppError(errors.CodeInternalError, "failed to update driver location", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.NewAppError(errors.CodeNotFound, "driver not found", nil)
	}

	return nil
}

func (r *DriverRepository) UpdateAvailability(ctx context.Context, driverID uuid.UUID, available bool) error {
	query := `UPDATE drivers SET is_available = $2, last_active_time = $3 WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, driverID, available, time.Now())
	if err != nil {
		return errors.NewAppError(errors.CodeInternalError, "failed to update driver availability", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.NewAppError(errors.CodeNotFound, "driver not found", nil)
	}

	return nil
}

func (r *DriverRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Driver, error) {
	query := `
		SELECT id, user_id, current_location, is_available, car_type,
			   rating, completed_trips, last_active_time, max_distance,
			   preferred_areas
		FROM drivers WHERE user_id = $1
	`

	var driver domain.Driver
	var currentLocationJSON, preferredAreasJSON []byte

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&driver.ID,
		&driver.UserID,
		&currentLocationJSON,
		&driver.IsAvailable,
		&driver.CarType,
		&driver.Rating,
		&driver.CompletedTrips,
		&driver.LastActiveTime,
		&driver.MaxDistance,
		&preferredAreasJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewAppError(errors.CodeNotFound, "driver not found", err)
		}
		return nil, errors.NewAppError(errors.CodeInternalError, "failed to get driver", err)
	}

	// Parse JSON fields
	json.Unmarshal(currentLocationJSON, &driver.CurrentLocation)
	json.Unmarshal(preferredAreasJSON, &driver.PreferredAreas)

	return &driver, nil
}

// MatchResultRepository implementation
type MatchResultRepository struct {
	db *sql.DB
}

func NewMatchResultRepository(db *sql.DB) domain.MatchResultRepository {
	return &MatchResultRepository{db: db}
}

func (r *MatchResultRepository) Create(ctx context.Context, result *domain.MatchResult) error {
	query := `
		INSERT INTO match_results (
			id, match_request_id, driver_id, score, estimated_distance,
			estimated_time, estimated_price, match_time, status, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		result.ID,
		result.MatchRequestID,
		result.DriverID,
		result.Score,
		result.EstimatedDistance,
		result.EstimatedTime.Nanoseconds(),
		result.EstimatedPrice,
		result.MatchTime,
		result.Status,
		result.CreatedAt,
	)

	if err != nil {
		return errors.NewAppError(errors.CodeInternalError, "failed to create match result", err)
	}

	return nil
}

func (r *MatchResultRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.MatchResult, error) {
	query := `
		SELECT id, match_request_id, driver_id, score, estimated_distance,
			   estimated_time, estimated_price, match_time, status, created_at
		FROM match_results WHERE id = $1
	`

	var result domain.MatchResult
	var estimatedTimeNanos int64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&result.ID,
		&result.MatchRequestID,
		&result.DriverID,
		&result.Score,
		&result.EstimatedDistance,
		&estimatedTimeNanos,
		&result.EstimatedPrice,
		&result.MatchTime,
		&result.Status,
		&result.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewAppError(errors.CodeNotFound, "match result not found", err)
		}
		return nil, errors.NewAppError(errors.CodeInternalError, "failed to get match result", err)
	}

	result.EstimatedTime = time.Duration(estimatedTimeNanos)
	return &result, nil
}

func (r *MatchResultRepository) GetByMatchRequestID(ctx context.Context, matchRequestID uuid.UUID) ([]*domain.MatchResult, error) {
	query := `
		SELECT id, match_request_id, driver_id, score, estimated_distance,
			   estimated_time, estimated_price, match_time, status, created_at
		FROM match_results 
		WHERE match_request_id = $1 
		ORDER BY score DESC
	`

	rows, err := r.db.QueryContext(ctx, query, matchRequestID)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeInternalError, "failed to query match results", err)
	}
	defer rows.Close()

	var results []*domain.MatchResult
	for rows.Next() {
		var result domain.MatchResult
		var estimatedTimeNanos int64

		err := rows.Scan(
			&result.ID,
			&result.MatchRequestID,
			&result.DriverID,
			&result.Score,
			&result.EstimatedDistance,
			&estimatedTimeNanos,
			&result.EstimatedPrice,
			&result.MatchTime,
			&result.Status,
			&result.CreatedAt,
		)

		if err != nil {
			continue
		}

		result.EstimatedTime = time.Duration(estimatedTimeNanos)
		results = append(results, &result)
	}

	return results, nil
}

func (r *MatchResultRepository) Update(ctx context.Context, result *domain.MatchResult) error {
	query := `
		UPDATE match_results SET
			score = $2, estimated_distance = $3, estimated_time = $4,
			estimated_price = $5, status = $6
		WHERE id = $1
	`

	dbResult, err := r.db.ExecContext(ctx, query,
		result.ID,
		result.Score,
		result.EstimatedDistance,
		result.EstimatedTime.Nanoseconds(),
		result.EstimatedPrice,
		result.Status,
	)

	if err != nil {
		return errors.NewAppError(errors.CodeInternalError, "failed to update match result", err)
	}

	rowsAffected, _ := dbResult.RowsAffected()
	if rowsAffected == 0 {
		return errors.NewAppError(errors.CodeNotFound, "match result not found", nil)
	}

	return nil
}

func (r *MatchResultRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM match_results WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.NewAppError(errors.CodeInternalError, "failed to delete match result", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.NewAppError(errors.CodeNotFound, "match result not found", nil)
	}

	return nil
}