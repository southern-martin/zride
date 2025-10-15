// Package infrastructure provides PostgreSQL user profile repository implementation
package infrastructure

import (
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/southern-martin/zride/backend/services/user-service/internal/domain"
	"github.com/southern-martin/zride/backend/shared"
)

type UserProfileRepository struct {
	db *sql.DB
}

func NewUserProfileRepository(db *sql.DB) domain.UserProfileRepository {
	return &UserProfileRepository{db: db}
}

func (r *UserProfileRepository) FindByID(id uuid.UUID) (*domain.UserProfile, error) {
	query := `
		SELECT id, zalo_id, name, phone, email, picture, user_type, is_active, is_verified,
		       verification_doc, rating, total_trips, total_ratings, preferences, languages,
		       bio, created_at, updated_at, last_active_at, version
		FROM user_profiles WHERE id = $1 AND deleted_at IS NULL`

	profile := &domain.UserProfile{}
	var lastActiveAt sql.NullTime
	var preferencesJSON, languagesJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&profile.ID, &profile.ZaloID, &profile.Name, &profile.Phone, &profile.Email,
		&profile.Picture, &profile.UserType, &profile.IsActive, &profile.IsVerified,
		&profile.VerificationDoc, &profile.Rating, &profile.TotalTrips, &profile.TotalRatings,
		&preferencesJSON, &languagesJSON, &profile.Bio, &profile.CreatedAt, &profile.UpdatedAt,
		&lastActiveAt, &profile.Version,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.NewNotFoundError("user profile not found", err)
		}
		return nil, shared.NewDatabaseError("failed to find user profile", err)
	}

	if lastActiveAt.Valid {
		profile.LastActiveAt = &lastActiveAt.Time
	}

	// Unmarshal JSON fields
	if err := json.Unmarshal(preferencesJSON, &profile.Preferences); err != nil {
		profile.Preferences = make(map[string]string)
	}
	if err := json.Unmarshal(languagesJSON, &profile.Languages); err != nil {
		profile.Languages = []string{"vi"}
	}

	return profile, nil
}

func (r *UserProfileRepository) FindByZaloID(zaloID string) (*domain.UserProfile, error) {
	query := `
		SELECT id, zalo_id, name, phone, email, picture, user_type, is_active, is_verified,
		       verification_doc, rating, total_trips, total_ratings, preferences, languages,
		       bio, created_at, updated_at, last_active_at, version
		FROM user_profiles WHERE zalo_id = $1 AND deleted_at IS NULL`

	profile := &domain.UserProfile{}
	var lastActiveAt sql.NullTime
	var preferencesJSON, languagesJSON []byte

	err := r.db.QueryRow(query, zaloID).Scan(
		&profile.ID, &profile.ZaloID, &profile.Name, &profile.Phone, &profile.Email,
		&profile.Picture, &profile.UserType, &profile.IsActive, &profile.IsVerified,
		&profile.VerificationDoc, &profile.Rating, &profile.TotalTrips, &profile.TotalRatings,
		&preferencesJSON, &languagesJSON, &profile.Bio, &profile.CreatedAt, &profile.UpdatedAt,
		&lastActiveAt, &profile.Version,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.NewNotFoundError("user profile not found", err)
		}
		return nil, shared.NewDatabaseError("failed to find user profile", err)
	}

	if lastActiveAt.Valid {
		profile.LastActiveAt = &lastActiveAt.Time
	}

	// Unmarshal JSON fields
	if err := json.Unmarshal(preferencesJSON, &profile.Preferences); err != nil {
		profile.Preferences = make(map[string]string)
	}
	if err := json.Unmarshal(languagesJSON, &profile.Languages); err != nil {
		profile.Languages = []string{"vi"}
	}

	return profile, nil
}

func (r *UserProfileRepository) FindByEmail(email string) (*domain.UserProfile, error) {
	// Implementation similar to FindByZaloID
	return r.findByField("email", email)
}

func (r *UserProfileRepository) FindByPhone(phone string) (*domain.UserProfile, error) {
	// Implementation similar to FindByZaloID
	return r.findByField("phone", phone)
}

func (r *UserProfileRepository) findByField(field, value string) (*domain.UserProfile, error) {
	query := `
		SELECT id, zalo_id, name, phone, email, picture, user_type, is_active, is_verified,
		       verification_doc, rating, total_trips, total_ratings, preferences, languages,
		       bio, created_at, updated_at, last_active_at, version
		FROM user_profiles WHERE ` + field + ` = $1 AND deleted_at IS NULL`

	profile := &domain.UserProfile{}
	var lastActiveAt sql.NullTime
	var preferencesJSON, languagesJSON []byte

	err := r.db.QueryRow(query, value).Scan(
		&profile.ID, &profile.ZaloID, &profile.Name, &profile.Phone, &profile.Email,
		&profile.Picture, &profile.UserType, &profile.IsActive, &profile.IsVerified,
		&profile.VerificationDoc, &profile.Rating, &profile.TotalTrips, &profile.TotalRatings,
		&preferencesJSON, &languagesJSON, &profile.Bio, &profile.CreatedAt, &profile.UpdatedAt,
		&lastActiveAt, &profile.Version,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.NewNotFoundError("user profile not found", err)
		}
		return nil, shared.NewDatabaseError("failed to find user profile", err)
	}

	if lastActiveAt.Valid {
		profile.LastActiveAt = &lastActiveAt.Time
	}

	// Unmarshal JSON fields
	if err := json.Unmarshal(preferencesJSON, &profile.Preferences); err != nil {
		profile.Preferences = make(map[string]string)
	}
	if err := json.Unmarshal(languagesJSON, &profile.Languages); err != nil {
		profile.Languages = []string{"vi"}
	}

	return profile, nil
}

func (r *UserProfileRepository) Create(profile *domain.UserProfile) error {
	preferencesJSON, _ := json.Marshal(profile.Preferences)
	languagesJSON, _ := json.Marshal(profile.Languages)

	query := `
		INSERT INTO user_profiles (id, zalo_id, name, phone, email, picture, user_type,
		                          is_active, is_verified, verification_doc, rating,
		                          total_trips, total_ratings, preferences, languages,
		                          bio, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`

	_, err := r.db.Exec(query,
		profile.ID, profile.ZaloID, profile.Name, profile.Phone, profile.Email,
		profile.Picture, profile.UserType, profile.IsActive, profile.IsVerified,
		profile.VerificationDoc, profile.Rating, profile.TotalTrips, profile.TotalRatings,
		preferencesJSON, languagesJSON, profile.Bio, profile.CreatedAt, profile.UpdatedAt,
		profile.Version,
	)

	if err != nil {
		return shared.NewDatabaseError("failed to create user profile", err)
	}

	return nil
}

func (r *UserProfileRepository) Update(profile *domain.UserProfile) error {
	preferencesJSON, _ := json.Marshal(profile.Preferences)
	languagesJSON, _ := json.Marshal(profile.Languages)

	query := `
		UPDATE user_profiles 
		SET name = $2, phone = $3, email = $4, picture = $5, user_type = $6,
		    is_active = $7, is_verified = $8, verification_doc = $9, rating = $10,
		    total_trips = $11, total_ratings = $12, preferences = $13, languages = $14,
		    bio = $15, updated_at = $16, last_active_at = $17, version = version + 1
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Exec(query,
		profile.ID, profile.Name, profile.Phone, profile.Email, profile.Picture,
		profile.UserType, profile.IsActive, profile.IsVerified, profile.VerificationDoc,
		profile.Rating, profile.TotalTrips, profile.TotalRatings, preferencesJSON,
		languagesJSON, profile.Bio, profile.UpdatedAt, profile.LastActiveAt,
	)

	if err != nil {
		return shared.NewDatabaseError("failed to update user profile", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return shared.NewDatabaseError("failed to check affected rows", err)
	}

	if rowsAffected == 0 {
		return shared.NewNotFoundError("user profile not found for update", nil)
	}

	return nil
}

func (r *UserProfileRepository) Delete(id uuid.UUID) error {
	query := `UPDATE user_profiles SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return shared.NewDatabaseError("failed to delete user profile", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return shared.NewDatabaseError("failed to check affected rows", err)
	}

	if rowsAffected == 0 {
		return shared.NewNotFoundError("user profile not found for deletion", nil)
	}

	return nil
}

func (r *UserProfileRepository) FindByUserType(userType string, limit, offset int) ([]*domain.UserProfile, error) {
	query := `
		SELECT id, zalo_id, name, phone, email, picture, user_type, is_active, is_verified,
		       verification_doc, rating, total_trips, total_ratings, preferences, languages,
		       bio, created_at, updated_at, last_active_at, version
		FROM user_profiles 
		WHERE user_type = $1 AND deleted_at IS NULL 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userType, limit, offset)
	if err != nil {
		return nil, shared.NewDatabaseError("failed to find users by type", err)
	}
	defer rows.Close()

	var profiles []*domain.UserProfile
	for rows.Next() {
		profile := &domain.UserProfile{}
		var lastActiveAt sql.NullTime
		var preferencesJSON, languagesJSON []byte

		err := rows.Scan(
			&profile.ID, &profile.ZaloID, &profile.Name, &profile.Phone, &profile.Email,
			&profile.Picture, &profile.UserType, &profile.IsActive, &profile.IsVerified,
			&profile.VerificationDoc, &profile.Rating, &profile.TotalTrips, &profile.TotalRatings,
			&preferencesJSON, &languagesJSON, &profile.Bio, &profile.CreatedAt, &profile.UpdatedAt,
			&lastActiveAt, &profile.Version,
		)

		if err != nil {
			return nil, shared.NewDatabaseError("failed to scan user profile", err)
		}

		if lastActiveAt.Valid {
			profile.LastActiveAt = &lastActiveAt.Time
		}

		// Unmarshal JSON fields
		if err := json.Unmarshal(preferencesJSON, &profile.Preferences); err != nil {
			profile.Preferences = make(map[string]string)
		}
		if err := json.Unmarshal(languagesJSON, &profile.Languages); err != nil {
			profile.Languages = []string{"vi"}
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func (r *UserProfileRepository) FindVerifiedDrivers(limit, offset int) ([]*domain.UserProfile, error) {
	query := `
		SELECT id, zalo_id, name, phone, email, picture, user_type, is_active, is_verified,
		       verification_doc, rating, total_trips, total_ratings, preferences, languages,
		       bio, created_at, updated_at, last_active_at, version
		FROM user_profiles 
		WHERE user_type = 'driver' AND is_verified = true AND is_active = true AND deleted_at IS NULL 
		ORDER BY rating DESC, total_trips DESC 
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, shared.NewDatabaseError("failed to find verified drivers", err)
	}
	defer rows.Close()

	var profiles []*domain.UserProfile
	for rows.Next() {
		profile := &domain.UserProfile{}
		var lastActiveAt sql.NullTime
		var preferencesJSON, languagesJSON []byte

		err := rows.Scan(
			&profile.ID, &profile.ZaloID, &profile.Name, &profile.Phone, &profile.Email,
			&profile.Picture, &profile.UserType, &profile.IsActive, &profile.IsVerified,
			&profile.VerificationDoc, &profile.Rating, &profile.TotalTrips, &profile.TotalRatings,
			&preferencesJSON, &languagesJSON, &profile.Bio, &profile.CreatedAt, &profile.UpdatedAt,
			&lastActiveAt, &profile.Version,
		)

		if err != nil {
			return nil, shared.NewDatabaseError("failed to scan user profile", err)
		}

		if lastActiveAt.Valid {
			profile.LastActiveAt = &lastActiveAt.Time
		}

		// Unmarshal JSON fields
		if err := json.Unmarshal(preferencesJSON, &profile.Preferences); err != nil {
			profile.Preferences = make(map[string]string)
		}
		if err := json.Unmarshal(languagesJSON, &profile.Languages); err != nil {
			profile.Languages = []string{"vi"}
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func (r *UserProfileRepository) Search(query string, limit, offset int) ([]*domain.UserProfile, error) {
	searchQuery := `
		SELECT id, zalo_id, name, phone, email, picture, user_type, is_active, is_verified,
		       verification_doc, rating, total_trips, total_ratings, preferences, languages,
		       bio, created_at, updated_at, last_active_at, version
		FROM user_profiles 
		WHERE (name ILIKE $1 OR phone ILIKE $1 OR email ILIKE $1) AND deleted_at IS NULL 
		ORDER BY rating DESC, total_trips DESC 
		LIMIT $2 OFFSET $3`

	searchPattern := "%" + query + "%"
	rows, err := r.db.Query(searchQuery, searchPattern, limit, offset)
	if err != nil {
		return nil, shared.NewDatabaseError("failed to search user profiles", err)
	}
	defer rows.Close()

	var profiles []*domain.UserProfile
	for rows.Next() {
		profile := &domain.UserProfile{}
		var lastActiveAt sql.NullTime
		var preferencesJSON, languagesJSON []byte

		err := rows.Scan(
			&profile.ID, &profile.ZaloID, &profile.Name, &profile.Phone, &profile.Email,
			&profile.Picture, &profile.UserType, &profile.IsActive, &profile.IsVerified,
			&profile.VerificationDoc, &profile.Rating, &profile.TotalTrips, &profile.TotalRatings,
			&preferencesJSON, &languagesJSON, &profile.Bio, &profile.CreatedAt, &profile.UpdatedAt,
			&lastActiveAt, &profile.Version,
		)

		if err != nil {
			return nil, shared.NewDatabaseError("failed to scan user profile", err)
		}

		if lastActiveAt.Valid {
			profile.LastActiveAt = &lastActiveAt.Time
		}

		// Unmarshal JSON fields
		if err := json.Unmarshal(preferencesJSON, &profile.Preferences); err != nil {
			profile.Preferences = make(map[string]string)
		}
		if err := json.Unmarshal(languagesJSON, &profile.Languages); err != nil {
			profile.Languages = []string{"vi"}
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}