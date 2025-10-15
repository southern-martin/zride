// Package infrastructure provides PostgreSQL vehicle repository implementation
package infrastructure

import (
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/southern-martin/zride/backend/services/user-service/internal/domain"
	"github.com/southern-martin/zride/backend/shared"
)

type VehicleRepository struct {
	db *sql.DB
}

func NewVehicleRepository(db *sql.DB) domain.VehicleRepository {
	return &VehicleRepository{db: db}
}

func (r *VehicleRepository) FindByID(id uuid.UUID) (*domain.Vehicle, error) {
	query := `
		SELECT id, owner_id, license_plate, brand, model, year, color, vehicle_type,
		       capacity, is_active, is_verified, photo_urls, documents, features,
		       created_at, updated_at, version
		FROM vehicles WHERE id = $1 AND deleted_at IS NULL`

	vehicle := &domain.Vehicle{}
	var photoURLsJSON, documentsJSON, featuresJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&vehicle.ID, &vehicle.OwnerID, &vehicle.LicensePlate, &vehicle.Brand,
		&vehicle.Model, &vehicle.Year, &vehicle.Color, &vehicle.VehicleType,
		&vehicle.Capacity, &vehicle.IsActive, &vehicle.IsVerified,
		&photoURLsJSON, &documentsJSON, &featuresJSON,
		&vehicle.CreatedAt, &vehicle.UpdatedAt, &vehicle.Version,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.NewNotFoundError("vehicle not found", err)
		}
		return nil, shared.NewDatabaseError("failed to find vehicle", err)
	}

	// Unmarshal JSON arrays
	if err := json.Unmarshal(photoURLsJSON, &vehicle.PhotoURLs); err != nil {
		vehicle.PhotoURLs = make([]string, 0)
	}
	if err := json.Unmarshal(documentsJSON, &vehicle.Documents); err != nil {
		vehicle.Documents = make([]string, 0)
	}
	if err := json.Unmarshal(featuresJSON, &vehicle.Features); err != nil {
		vehicle.Features = make([]string, 0)
	}

	return vehicle, nil
}

func (r *VehicleRepository) FindByOwnerID(ownerID uuid.UUID) ([]*domain.Vehicle, error) {
	query := `
		SELECT id, owner_id, license_plate, brand, model, year, color, vehicle_type,
		       capacity, is_active, is_verified, photo_urls, documents, features,
		       created_at, updated_at, version
		FROM vehicles WHERE owner_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, ownerID)
	if err != nil {
		return nil, shared.NewDatabaseError("failed to find vehicles by owner", err)
	}
	defer rows.Close()

	var vehicles []*domain.Vehicle
	for rows.Next() {
		vehicle := &domain.Vehicle{}
		var photoURLsJSON, documentsJSON, featuresJSON []byte

		err := rows.Scan(
			&vehicle.ID, &vehicle.OwnerID, &vehicle.LicensePlate, &vehicle.Brand,
			&vehicle.Model, &vehicle.Year, &vehicle.Color, &vehicle.VehicleType,
			&vehicle.Capacity, &vehicle.IsActive, &vehicle.IsVerified,
			&photoURLsJSON, &documentsJSON, &featuresJSON,
			&vehicle.CreatedAt, &vehicle.UpdatedAt, &vehicle.Version,
		)

		if err != nil {
			return nil, shared.NewDatabaseError("failed to scan vehicle", err)
		}

		// Unmarshal JSON arrays
		if err := json.Unmarshal(photoURLsJSON, &vehicle.PhotoURLs); err != nil {
			vehicle.PhotoURLs = make([]string, 0)
		}
		if err := json.Unmarshal(documentsJSON, &vehicle.Documents); err != nil {
			vehicle.Documents = make([]string, 0)
		}
		if err := json.Unmarshal(featuresJSON, &vehicle.Features); err != nil {
			vehicle.Features = make([]string, 0)
		}

		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}

func (r *VehicleRepository) FindByLicensePlate(licensePlate string) (*domain.Vehicle, error) {
	query := `
		SELECT id, owner_id, license_plate, brand, model, year, color, vehicle_type,
		       capacity, is_active, is_verified, photo_urls, documents, features,
		       created_at, updated_at, version
		FROM vehicles WHERE license_plate = $1 AND deleted_at IS NULL`

	vehicle := &domain.Vehicle{}
	var photoURLsJSON, documentsJSON, featuresJSON []byte

	err := r.db.QueryRow(query, licensePlate).Scan(
		&vehicle.ID, &vehicle.OwnerID, &vehicle.LicensePlate, &vehicle.Brand,
		&vehicle.Model, &vehicle.Year, &vehicle.Color, &vehicle.VehicleType,
		&vehicle.Capacity, &vehicle.IsActive, &vehicle.IsVerified,
		&photoURLsJSON, &documentsJSON, &featuresJSON,
		&vehicle.CreatedAt, &vehicle.UpdatedAt, &vehicle.Version,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.NewNotFoundError("vehicle not found", err)
		}
		return nil, shared.NewDatabaseError("failed to find vehicle", err)
	}

	// Unmarshal JSON arrays
	if err := json.Unmarshal(photoURLsJSON, &vehicle.PhotoURLs); err != nil {
		vehicle.PhotoURLs = make([]string, 0)
	}
	if err := json.Unmarshal(documentsJSON, &vehicle.Documents); err != nil {
		vehicle.Documents = make([]string, 0)
	}
	if err := json.Unmarshal(featuresJSON, &vehicle.Features); err != nil {
		vehicle.Features = make([]string, 0)
	}

	return vehicle, nil
}

func (r *VehicleRepository) Create(vehicle *domain.Vehicle) error {
	photoURLsJSON, _ := json.Marshal(vehicle.PhotoURLs)
	documentsJSON, _ := json.Marshal(vehicle.Documents)
	featuresJSON, _ := json.Marshal(vehicle.Features)

	query := `
		INSERT INTO vehicles (id, owner_id, license_plate, brand, model, year, color,
		                     vehicle_type, capacity, is_active, is_verified, photo_urls,
		                     documents, features, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`

	_, err := r.db.Exec(query,
		vehicle.ID, vehicle.OwnerID, vehicle.LicensePlate, vehicle.Brand,
		vehicle.Model, vehicle.Year, vehicle.Color, vehicle.VehicleType,
		vehicle.Capacity, vehicle.IsActive, vehicle.IsVerified,
		photoURLsJSON, documentsJSON, featuresJSON,
		vehicle.CreatedAt, vehicle.UpdatedAt, vehicle.Version,
	)

	if err != nil {
		return shared.NewDatabaseError("failed to create vehicle", err)
	}

	return nil
}

func (r *VehicleRepository) Update(vehicle *domain.Vehicle) error {
	photoURLsJSON, _ := json.Marshal(vehicle.PhotoURLs)
	documentsJSON, _ := json.Marshal(vehicle.Documents)
	featuresJSON, _ := json.Marshal(vehicle.Features)

	query := `
		UPDATE vehicles 
		SET license_plate = $2, brand = $3, model = $4, year = $5, color = $6,
		    vehicle_type = $7, capacity = $8, is_active = $9, is_verified = $10,
		    photo_urls = $11, documents = $12, features = $13, updated_at = $14,
		    version = version + 1
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Exec(query,
		vehicle.ID, vehicle.LicensePlate, vehicle.Brand, vehicle.Model,
		vehicle.Year, vehicle.Color, vehicle.VehicleType, vehicle.Capacity,
		vehicle.IsActive, vehicle.IsVerified, photoURLsJSON, documentsJSON,
		featuresJSON, vehicle.UpdatedAt,
	)

	if err != nil {
		return shared.NewDatabaseError("failed to update vehicle", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return shared.NewDatabaseError("failed to check affected rows", err)
	}

	if rowsAffected == 0 {
		return shared.NewNotFoundError("vehicle not found for update", nil)
	}

	return nil
}

func (r *VehicleRepository) Delete(id uuid.UUID) error {
	query := `UPDATE vehicles SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return shared.NewDatabaseError("failed to delete vehicle", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return shared.NewDatabaseError("failed to check affected rows", err)
	}

	if rowsAffected == 0 {
		return shared.NewNotFoundError("vehicle not found for deletion", nil)
	}

	return nil
}

func (r *VehicleRepository) FindActiveByOwnerID(ownerID uuid.UUID) ([]*domain.Vehicle, error) {
	query := `
		SELECT id, owner_id, license_plate, brand, model, year, color, vehicle_type,
		       capacity, is_active, is_verified, photo_urls, documents, features,
		       created_at, updated_at, version
		FROM vehicles WHERE owner_id = $1 AND is_active = true AND deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, ownerID)
	if err != nil {
		return nil, shared.NewDatabaseError("failed to find active vehicles", err)
	}
	defer rows.Close()

	var vehicles []*domain.Vehicle
	for rows.Next() {
		vehicle := &domain.Vehicle{}
		var photoURLsJSON, documentsJSON, featuresJSON []byte

		err := rows.Scan(
			&vehicle.ID, &vehicle.OwnerID, &vehicle.LicensePlate, &vehicle.Brand,
			&vehicle.Model, &vehicle.Year, &vehicle.Color, &vehicle.VehicleType,
			&vehicle.Capacity, &vehicle.IsActive, &vehicle.IsVerified,
			&photoURLsJSON, &documentsJSON, &featuresJSON,
			&vehicle.CreatedAt, &vehicle.UpdatedAt, &vehicle.Version,
		)

		if err != nil {
			return nil, shared.NewDatabaseError("failed to scan vehicle", err)
		}

		// Unmarshal JSON arrays
		if err := json.Unmarshal(photoURLsJSON, &vehicle.PhotoURLs); err != nil {
			vehicle.PhotoURLs = make([]string, 0)
		}
		if err := json.Unmarshal(documentsJSON, &vehicle.Documents); err != nil {
			vehicle.Documents = make([]string, 0)
		}
		if err := json.Unmarshal(featuresJSON, &vehicle.Features); err != nil {
			vehicle.Features = make([]string, 0)
		}

		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}

func (r *VehicleRepository) FindPendingVerification(limit, offset int) ([]*domain.Vehicle, error) {
	query := `
		SELECT id, owner_id, license_plate, brand, model, year, color, vehicle_type,
		       capacity, is_active, is_verified, photo_urls, documents, features,
		       created_at, updated_at, version
		FROM vehicles 
		WHERE is_verified = false AND is_active = true AND deleted_at IS NULL
		  AND array_length(documents, 1) > 0 AND array_length(photo_urls, 1) > 0
		ORDER BY created_at ASC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, shared.NewDatabaseError("failed to find pending verification vehicles", err)
	}
	defer rows.Close()

	var vehicles []*domain.Vehicle
	for rows.Next() {
		vehicle := &domain.Vehicle{}
		var photoURLsJSON, documentsJSON, featuresJSON []byte

		err := rows.Scan(
			&vehicle.ID, &vehicle.OwnerID, &vehicle.LicensePlate, &vehicle.Brand,
			&vehicle.Model, &vehicle.Year, &vehicle.Color, &vehicle.VehicleType,
			&vehicle.Capacity, &vehicle.IsActive, &vehicle.IsVerified,
			&photoURLsJSON, &documentsJSON, &featuresJSON,
			&vehicle.CreatedAt, &vehicle.UpdatedAt, &vehicle.Version,
		)

		if err != nil {
			return nil, shared.NewDatabaseError("failed to scan vehicle", err)
		}

		// Unmarshal JSON arrays
		if err := json.Unmarshal(photoURLsJSON, &vehicle.PhotoURLs); err != nil {
			vehicle.PhotoURLs = make([]string, 0)
		}
		if err := json.Unmarshal(documentsJSON, &vehicle.Documents); err != nil {
			vehicle.Documents = make([]string, 0)
		}
		if err := json.Unmarshal(featuresJSON, &vehicle.Features); err != nil {
			vehicle.Features = make([]string, 0)
		}

		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}