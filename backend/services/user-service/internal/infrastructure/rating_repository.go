// Package infrastructure provides PostgreSQL rating repository implementation
package infrastructure

import (
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/southern-martin/zride/backend/services/user-service/internal/domain"
	"github.com/southern-martin/zride/backend/shared"
)

type RatingRepository struct {
	db *sql.DB
}

func NewRatingRepository(db *sql.DB) domain.RatingRepository {
	return &RatingRepository{db: db}
}

func (r *RatingRepository) Create(rating *domain.Rating) error {
	tagsJSON, _ := json.Marshal(rating.Tags)

	query := `
		INSERT INTO ratings (id, rater_id, rated_id, trip_id, score, comment, tags,
		                    is_visible, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.Exec(query,
		rating.ID, rating.RaterID, rating.RatedID, rating.TripID,
		rating.Score, rating.Comment, tagsJSON, rating.IsVisible,
		rating.CreatedAt, rating.UpdatedAt, rating.Version,
	)

	if err != nil {
		return shared.NewDatabaseError("failed to create rating", err)
	}

	return nil
}

func (r *RatingRepository) FindByUserID(userID uuid.UUID, limit, offset int) ([]*domain.Rating, error) {
	query := `
		SELECT id, rater_id, rated_id, trip_id, score, comment, tags,
		       is_visible, created_at, updated_at, version
		FROM ratings 
		WHERE rated_id = $1 AND is_visible = true AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, shared.NewDatabaseError("failed to find ratings by user ID", err)
	}
	defer rows.Close()

	var ratings []*domain.Rating
	for rows.Next() {
		rating := &domain.Rating{}
		var tagsJSON []byte

		err := rows.Scan(
			&rating.ID, &rating.RaterID, &rating.RatedID, &rating.TripID,
			&rating.Score, &rating.Comment, &tagsJSON, &rating.IsVisible,
			&rating.CreatedAt, &rating.UpdatedAt, &rating.Version,
		)

		if err != nil {
			return nil, shared.NewDatabaseError("failed to scan rating", err)
		}

		// Unmarshal tags
		if err := json.Unmarshal(tagsJSON, &rating.Tags); err != nil {
			rating.Tags = make([]string, 0)
		}

		ratings = append(ratings, rating)
	}

	return ratings, nil
}

func (r *RatingRepository) FindByRaterID(raterID uuid.UUID, limit, offset int) ([]*domain.Rating, error) {
	query := `
		SELECT id, rater_id, rated_id, trip_id, score, comment, tags,
		       is_visible, created_at, updated_at, version
		FROM ratings 
		WHERE rater_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, raterID, limit, offset)
	if err != nil {
		return nil, shared.NewDatabaseError("failed to find ratings by rater ID", err)
	}
	defer rows.Close()

	var ratings []*domain.Rating
	for rows.Next() {
		rating := &domain.Rating{}
		var tagsJSON []byte

		err := rows.Scan(
			&rating.ID, &rating.RaterID, &rating.RatedID, &rating.TripID,
			&rating.Score, &rating.Comment, &tagsJSON, &rating.IsVisible,
			&rating.CreatedAt, &rating.UpdatedAt, &rating.Version,
		)

		if err != nil {
			return nil, shared.NewDatabaseError("failed to scan rating", err)
		}

		// Unmarshal tags
		if err := json.Unmarshal(tagsJSON, &rating.Tags); err != nil {
			rating.Tags = make([]string, 0)
		}

		ratings = append(ratings, rating)
	}

	return ratings, nil
}

func (r *RatingRepository) GetAverageRating(userID uuid.UUID) (float64, int, error) {
	query := `
		SELECT COALESCE(AVG(score), 0), COUNT(*)
		FROM ratings 
		WHERE rated_id = $1 AND is_visible = true AND deleted_at IS NULL`

	var avgRating float64
	var count int

	err := r.db.QueryRow(query, userID).Scan(&avgRating, &count)
	if err != nil {
		return 0, 0, shared.NewDatabaseError("failed to get average rating", err)
	}

	return avgRating, count, nil
}

func (r *RatingRepository) FindByTripID(tripID uuid.UUID) (*domain.Rating, error) {
	query := `
		SELECT id, rater_id, rated_id, trip_id, score, comment, tags,
		       is_visible, created_at, updated_at, version
		FROM ratings WHERE trip_id = $1 AND deleted_at IS NULL`

	rating := &domain.Rating{}
	var tagsJSON []byte

	err := r.db.QueryRow(query, tripID).Scan(
		&rating.ID, &rating.RaterID, &rating.RatedID, &rating.TripID,
		&rating.Score, &rating.Comment, &tagsJSON, &rating.IsVisible,
		&rating.CreatedAt, &rating.UpdatedAt, &rating.Version,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.NewNotFoundError("rating not found", err)
		}
		return nil, shared.NewDatabaseError("failed to find rating by trip ID", err)
	}

	// Unmarshal tags
	if err := json.Unmarshal(tagsJSON, &rating.Tags); err != nil {
		rating.Tags = make([]string, 0)
	}

	return rating, nil
}

func (r *RatingRepository) Update(rating *domain.Rating) error {
	tagsJSON, _ := json.Marshal(rating.Tags)

	query := `
		UPDATE ratings 
		SET score = $2, comment = $3, tags = $4, is_visible = $5, 
		    updated_at = $6, version = version + 1
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Exec(query,
		rating.ID, rating.Score, rating.Comment, tagsJSON,
		rating.IsVisible, rating.UpdatedAt,
	)

	if err != nil {
		return shared.NewDatabaseError("failed to update rating", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return shared.NewDatabaseError("failed to check affected rows", err)
	}

	if rowsAffected == 0 {
		return shared.NewNotFoundError("rating not found for update", nil)
	}

	return nil
}

func (r *RatingRepository) Delete(id uuid.UUID) error {
	query := `UPDATE ratings SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return shared.NewDatabaseError("failed to delete rating", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return shared.NewDatabaseError("failed to check affected rows", err)
	}

	if rowsAffected == 0 {
		return shared.NewNotFoundError("rating not found for deletion", nil)
	}

	return nil
}