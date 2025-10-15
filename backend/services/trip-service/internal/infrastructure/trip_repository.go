package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/southern-martin/zride/backend/services/trip-service/internal/domain"
	"github.com/southern-martin/zride/backend/shared"
)

// TripRepository implements the domain.TripRepository interface
type TripRepository struct {
	db *sql.DB
}

// NewTripRepository creates a new trip repository
func NewTripRepository(db *sql.DB) *TripRepository {
	return &TripRepository{
		db: db,
	}
}

// CreateTrip creates a new trip
func (r *TripRepository) CreateTrip(ctx context.Context, trip *domain.Trip) error {
	query := `
		INSERT INTO trips (
			id, passenger_id, driver_id, vehicle_id, pickup_location, dropoff_location,
			scheduled_time, requested_at, accepted_at, pickup_time, dropoff_time,
			status, route_info, pricing_info, passenger_notes, driver_notes,
			cancellation_info, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
		)`

	_, err := r.db.ExecContext(ctx, query,
		trip.ID, trip.PassengerID, trip.DriverID, trip.VehicleID,
		trip.PickupLocation, trip.DropoffLocation, trip.ScheduledTime,
		trip.RequestedAt, trip.AcceptedAt, trip.PickupTime, trip.DropoffTime,
		trip.Status, trip.RouteInfo, trip.PricingInfo,
		trip.PassengerNotes, trip.DriverNotes, trip.CancellationInfo,
		trip.CreatedAt, trip.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				return shared.NewError(shared.ErrCodeConflict, "Trip already exists", err.Error())
			case "23503": // foreign_key_violation
				return shared.NewError(shared.ErrCodeValidation, "Invalid foreign key reference", err.Error())
			}
		}
		return shared.NewError(shared.ErrCodeInternal, "Failed to create trip", err.Error())
	}

	return nil
}

// GetTripByID retrieves a trip by its ID
func (r *TripRepository) GetTripByID(ctx context.Context, id uuid.UUID) (*domain.Trip, error) {
	query := `
		SELECT id, passenger_id, driver_id, vehicle_id, pickup_location, dropoff_location,
			   scheduled_time, requested_at, accepted_at, pickup_time, dropoff_time,
			   status, route_info, pricing_info, passenger_notes, driver_notes,
			   cancellation_info, created_at, updated_at, deleted_at
		FROM trips
		WHERE id = $1 AND deleted_at IS NULL`

	trip := &domain.Trip{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&trip.ID, &trip.PassengerID, &trip.DriverID, &trip.VehicleID,
		&trip.PickupLocation, &trip.DropoffLocation, &trip.ScheduledTime,
		&trip.RequestedAt, &trip.AcceptedAt, &trip.PickupTime, &trip.DropoffTime,
		&trip.Status, &trip.RouteInfo, &trip.PricingInfo,
		&trip.PassengerNotes, &trip.DriverNotes, &trip.CancellationInfo,
		&trip.CreatedAt, &trip.UpdatedAt, &trip.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.NewError(shared.ErrCodeNotFound, "Trip not found", "")
		}
		return nil, shared.NewError(shared.ErrCodeInternal, "Failed to get trip", err.Error())
	}

	return trip, nil
}

// UpdateTrip updates an existing trip
func (r *TripRepository) UpdateTrip(ctx context.Context, trip *domain.Trip) error {
	query := `
		UPDATE trips SET
			driver_id = $2, vehicle_id = $3, pickup_location = $4, dropoff_location = $5,
			scheduled_time = $6, accepted_at = $7, pickup_time = $8, dropoff_time = $9,
			status = $10, route_info = $11, pricing_info = $12, passenger_notes = $13,
			driver_notes = $14, cancellation_info = $15, updated_at = $16
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query,
		trip.ID, trip.DriverID, trip.VehicleID, trip.PickupLocation, trip.DropoffLocation,
		trip.ScheduledTime, trip.AcceptedAt, trip.PickupTime, trip.DropoffTime,
		trip.Status, trip.RouteInfo, trip.PricingInfo, trip.PassengerNotes,
		trip.DriverNotes, trip.CancellationInfo, time.Now(),
	)

	if err != nil {
		return shared.NewError(shared.ErrCodeInternal, "Failed to update trip", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return shared.NewError(shared.ErrCodeInternal, "Failed to get update result", err.Error())
	}

	if rowsAffected == 0 {
		return shared.NewError(shared.ErrCodeNotFound, "Trip not found", "")
	}

	return nil
}

// DeleteTrip soft deletes a trip
func (r *TripRepository) DeleteTrip(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE trips SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return shared.NewError(shared.ErrCodeInternal, "Failed to delete trip", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return shared.NewError(shared.ErrCodeInternal, "Failed to get delete result", err.Error())
	}

	if rowsAffected == 0 {
		return shared.NewError(shared.ErrCodeNotFound, "Trip not found", "")
	}

	return nil
}

// SearchTrips searches for trips based on criteria
func (r *TripRepository) SearchTrips(ctx context.Context, criteria domain.TripSearchCriteria) ([]*domain.Trip, int64, error) {
	whereClauses := []string{"deleted_at IS NULL"}
	args := []interface{}{}
	argIndex := 1

	// Build WHERE clauses
	if criteria.PassengerID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("passenger_id = $%d", argIndex))
		args = append(args, *criteria.PassengerID)
		argIndex++
	}

	if criteria.DriverID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("driver_id = $%d", argIndex))
		args = append(args, *criteria.DriverID)
		argIndex++
	}

	if criteria.Status != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *criteria.Status)
		argIndex++
	}

	if criteria.FromTime != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *criteria.FromTime)
		argIndex++
	}

	if criteria.ToTime != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *criteria.ToTime)
		argIndex++
	}

	// Add location-based searches if provided
	if criteria.PickupLocation != nil && criteria.PickupRadius != nil {
		whereClauses = append(whereClauses, fmt.Sprintf(`
			ST_DWithin(
				ST_Point((pickup_location->>'lng')::float, (pickup_location->>'lat')::float)::geography,
				ST_Point($%d, $%d)::geography,
				$%d
			)`, argIndex, argIndex+1, argIndex+2))
		args = append(args, criteria.PickupLocation.Longitude, criteria.PickupLocation.Latitude, *criteria.PickupRadius*1000)
		argIndex += 3
	}

	if criteria.MinPrice != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("(pricing_info->>'total_fare')::float >= $%d", argIndex))
		args = append(args, *criteria.MinPrice)
		argIndex++
	}

	if criteria.MaxPrice != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("(pricing_info->>'total_fare')::float <= $%d", argIndex))
		args = append(args, *criteria.MaxPrice)
		argIndex++
	}

	whereClause := strings.Join(whereClauses, " AND ")

	// Count total results
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM trips WHERE %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, shared.NewError(shared.ErrCodeInternal, "Failed to count trips", err.Error())
	}

	// Build ORDER BY clause
	orderBy := "created_at DESC"
	if criteria.SortBy != "" {
		direction := "DESC"
		if criteria.SortOrder == "asc" {
			direction = "ASC"
		}
		
		switch criteria.SortBy {
		case "created_at", "price", "distance":
			if criteria.SortBy == "price" {
				orderBy = fmt.Sprintf("(pricing_info->>'total_fare')::float %s", direction)
			} else if criteria.SortBy == "distance" {
				orderBy = fmt.Sprintf("(route_info->>'distance_km')::float %s", direction)
			} else {
				orderBy = fmt.Sprintf("%s %s", criteria.SortBy, direction)
			}
		}
	}

	// Main query with pagination
	query := fmt.Sprintf(`
		SELECT id, passenger_id, driver_id, vehicle_id, pickup_location, dropoff_location,
			   scheduled_time, requested_at, accepted_at, pickup_time, dropoff_time,
			   status, route_info, pricing_info, passenger_notes, driver_notes,
			   cancellation_info, created_at, updated_at, deleted_at
		FROM trips
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, argIndex, argIndex+1)

	args = append(args, criteria.Limit, criteria.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, shared.NewError(shared.ErrCodeInternal, "Failed to search trips", err.Error())
	}
	defer rows.Close()

	trips := []*domain.Trip{}
	for rows.Next() {
		trip := &domain.Trip{}
		err := rows.Scan(
			&trip.ID, &trip.PassengerID, &trip.DriverID, &trip.VehicleID,
			&trip.PickupLocation, &trip.DropoffLocation, &trip.ScheduledTime,
			&trip.RequestedAt, &trip.AcceptedAt, &trip.PickupTime, &trip.DropoffTime,
			&trip.Status, &trip.RouteInfo, &trip.PricingInfo,
			&trip.PassengerNotes, &trip.DriverNotes, &trip.CancellationInfo,
			&trip.CreatedAt, &trip.UpdatedAt, &trip.DeletedAt,
		)
		if err != nil {
			return nil, 0, shared.NewError(shared.ErrCodeInternal, "Failed to scan trip", err.Error())
		}
		trips = append(trips, trip)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, shared.NewError(shared.ErrCodeInternal, "Failed to iterate trips", err.Error())
	}

	return trips, total, nil
}

// GetActiveTrips retrieves all active trips
func (r *TripRepository) GetActiveTrips(ctx context.Context) ([]*domain.Trip, error) {
	query := `
		SELECT id, passenger_id, driver_id, vehicle_id, pickup_location, dropoff_location,
			   scheduled_time, requested_at, accepted_at, pickup_time, dropoff_time,
			   status, route_info, pricing_info, passenger_notes, driver_notes,
			   cancellation_info, created_at, updated_at, deleted_at
		FROM trips
		WHERE status IN ('requested', 'accepted', 'in_progress') AND deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, shared.NewError(shared.ErrCodeInternal, "Failed to get active trips", err.Error())
	}
	defer rows.Close()

	trips := []*domain.Trip{}
	for rows.Next() {
		trip := &domain.Trip{}
		err := rows.Scan(
			&trip.ID, &trip.PassengerID, &trip.DriverID, &trip.VehicleID,
			&trip.PickupLocation, &trip.DropoffLocation, &trip.ScheduledTime,
			&trip.RequestedAt, &trip.AcceptedAt, &trip.PickupTime, &trip.DropoffTime,
			&trip.Status, &trip.RouteInfo, &trip.PricingInfo,
			&trip.PassengerNotes, &trip.DriverNotes, &trip.CancellationInfo,
			&trip.CreatedAt, &trip.UpdatedAt, &trip.DeletedAt,
		)
		if err != nil {
			return nil, shared.NewError(shared.ErrCodeInternal, "Failed to scan trip", err.Error())
		}
		trips = append(trips, trip)
	}

	return trips, rows.Err()
}

// GetTripsByPassenger retrieves trips for a specific passenger
func (r *TripRepository) GetTripsByPassenger(ctx context.Context, passengerID uuid.UUID, limit, offset int) ([]*domain.Trip, int64, error) {
	criteria := domain.TripSearchCriteria{
		PassengerID: &passengerID,
		Limit:       limit,
		Offset:      offset,
		SortBy:      "created_at",
		SortOrder:   "desc",
	}
	return r.SearchTrips(ctx, criteria)
}

// GetTripsByDriver retrieves trips for a specific driver
func (r *TripRepository) GetTripsByDriver(ctx context.Context, driverID uuid.UUID, limit, offset int) ([]*domain.Trip, int64, error) {
	criteria := domain.TripSearchCriteria{
		DriverID:  &driverID,
		Limit:     limit,
		Offset:    offset,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
	return r.SearchTrips(ctx, criteria)
}

// GetTripsByStatus retrieves trips by status
func (r *TripRepository) GetTripsByStatus(ctx context.Context, status domain.TripStatus, limit, offset int) ([]*domain.Trip, int64, error) {
	criteria := domain.TripSearchCriteria{
		Status:    &status,
		Limit:     limit,
		Offset:    offset,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
	return r.SearchTrips(ctx, criteria)
}

// GetNearbyTrips finds trips near a specific location
func (r *TripRepository) GetNearbyTrips(ctx context.Context, location domain.Location, radiusKM float64, status *domain.TripStatus) ([]*domain.Trip, error) {
	criteria := domain.TripSearchCriteria{
		PickupLocation: &location,
		PickupRadius:   &radiusKM,
		Status:         status,
		Limit:          50, // Default limit for nearby searches
		SortBy:         "created_at",
		SortOrder:      "desc",
	}
	
	trips, _, err := r.SearchTrips(ctx, criteria)
	return trips, err
}

// UpdateTripStatus updates only the status of a trip
func (r *TripRepository) UpdateTripStatus(ctx context.Context, tripID uuid.UUID, status domain.TripStatus) error {
	query := `UPDATE trips SET status = $1, updated_at = $2 WHERE id = $3 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, status, time.Now(), tripID)
	if err != nil {
		return shared.NewError(shared.ErrCodeInternal, "Failed to update trip status", err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return shared.NewError(shared.ErrCodeInternal, "Failed to get update result", err.Error())
	}

	if rowsAffected == 0 {
		return shared.NewError(shared.ErrCodeNotFound, "Trip not found", "")
	}

	return nil
}

// GetTripStatistics gets trip statistics for analytics
func (r *TripRepository) GetTripStatistics(ctx context.Context, from, to *time.Time) (*domain.TripStatistics, error) {
	whereClauses := []string{"deleted_at IS NULL"}
	args := []interface{}{}
	argIndex := 1

	if from != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *from)
		argIndex++
	}

	if to != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *to)
		argIndex++
	}

	whereClause := strings.Join(whereClauses, " AND ")

	query := fmt.Sprintf(`
		SELECT 
			COUNT(*) as total_trips,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_trips,
			COUNT(CASE WHEN status = 'cancelled' THEN 1 END) as cancelled_trips,
			COUNT(CASE WHEN status IN ('requested', 'accepted', 'in_progress') THEN 1 END) as active_trips,
			COALESCE(SUM(CASE WHEN status = 'completed' AND pricing_info IS NOT NULL THEN (pricing_info->>'total_fare')::float ELSE 0 END), 0) as total_revenue,
			COALESCE(AVG(CASE WHEN status = 'completed' AND route_info IS NOT NULL THEN (route_info->>'distance_km')::float END), 0) as average_distance,
			COALESCE(AVG(CASE WHEN status = 'completed' AND pricing_info IS NOT NULL THEN (pricing_info->>'total_fare')::float END), 0) as average_fare
		FROM trips
		WHERE %s`, whereClause)

	stats := &domain.TripStatistics{}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&stats.TotalTrips, &stats.CompletedTrips, &stats.CancelledTrips,
		&stats.ActiveTrips, &stats.TotalRevenue, &stats.AverageDistance, &stats.AverageFare,
	)

	if err != nil {
		return nil, shared.NewError(shared.ErrCodeInternal, "Failed to get trip statistics", err.Error())
	}

	return stats, nil
}