// Package infrastructure provides PostgreSQL user repository implementation
package infrastructure

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/southern-martin/zride/backend/services/auth-service/internal/domain"
	"github.com/southern-martin/zride/backend/shared"
)

// UserRepository implements domain.UserRepository interface
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates new PostgreSQL user repository
func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &UserRepository{
		db: db,
	}
}

// FindByID finds user by ID
func (r *UserRepository) FindByID(userID uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, zalo_id, name, phone, email, picture, user_type, is_active, 
		       last_login_at, refresh_token, created_at, updated_at, version
		FROM users WHERE id = $1 AND deleted_at IS NULL`

	user := &domain.User{}
	var lastLoginAt sql.NullTime

	err := r.db.QueryRow(query, userID).Scan(
		&user.ID, &user.ZaloID, &user.Name, &user.Phone, &user.Email,
		&user.Picture, &user.UserType, &user.IsActive, &lastLoginAt,
		&user.RefreshToken, &user.CreatedAt, &user.UpdatedAt, &user.Version,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.NewNotFoundError("user not found", err)
		}
		return nil, shared.NewDatabaseError("failed to find user by ID", err)
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

// FindByZaloID finds user by Zalo ID
func (r *UserRepository) FindByZaloID(zaloID string) (*domain.User, error) {
	query := `
		SELECT id, zalo_id, name, phone, email, picture, user_type, is_active, 
		       last_login_at, refresh_token, created_at, updated_at, version
		FROM users WHERE zalo_id = $1 AND deleted_at IS NULL`

	user := &domain.User{}
	var lastLoginAt sql.NullTime

	err := r.db.QueryRow(query, zaloID).Scan(
		&user.ID, &user.ZaloID, &user.Name, &user.Phone, &user.Email,
		&user.Picture, &user.UserType, &user.IsActive, &lastLoginAt,
		&user.RefreshToken, &user.CreatedAt, &user.UpdatedAt, &user.Version,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.NewNotFoundError("user not found", err)
		}
		return nil, shared.NewDatabaseError("failed to find user by Zalo ID", err)
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

// FindByEmail finds user by email
func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, zalo_id, name, phone, email, picture, user_type, is_active, 
		       last_login_at, refresh_token, created_at, updated_at, version
		FROM users WHERE email = $1 AND deleted_at IS NULL`

	user := &domain.User{}
	var lastLoginAt sql.NullTime

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.ZaloID, &user.Name, &user.Phone, &user.Email,
		&user.Picture, &user.UserType, &user.IsActive, &lastLoginAt,
		&user.RefreshToken, &user.CreatedAt, &user.UpdatedAt, &user.Version,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.NewNotFoundError("user not found", err)
		}
		return nil, shared.NewDatabaseError("failed to find user by email", err)
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

// FindByPhone finds user by phone
func (r *UserRepository) FindByPhone(phone string) (*domain.User, error) {
	query := `
		SELECT id, zalo_id, name, phone, email, picture, user_type, is_active, 
		       last_login_at, refresh_token, created_at, updated_at, version
		FROM users WHERE phone = $1 AND deleted_at IS NULL`

	user := &domain.User{}
	var lastLoginAt sql.NullTime

	err := r.db.QueryRow(query, phone).Scan(
		&user.ID, &user.ZaloID, &user.Name, &user.Phone, &user.Email,
		&user.Picture, &user.UserType, &user.IsActive, &lastLoginAt,
		&user.RefreshToken, &user.CreatedAt, &user.UpdatedAt, &user.Version,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, shared.NewNotFoundError("user not found", err)
		}
		return nil, shared.NewDatabaseError("failed to find user by phone", err)
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

// Create creates a new user
func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (id, zalo_id, name, phone, email, picture, user_type, is_active, 
		                  refresh_token, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.Exec(query,
		user.ID, user.ZaloID, user.Name, user.Phone, user.Email,
		user.Picture, user.UserType, user.IsActive, user.RefreshToken,
		user.CreatedAt, user.UpdatedAt, user.Version,
	)

	if err != nil {
		return shared.NewDatabaseError("failed to create user", err)
	}

	return nil
}

// Update updates an existing user
func (r *UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users 
		SET name = $2, phone = $3, email = $4, picture = $5, user_type = $6, 
		    is_active = $7, last_login_at = $8, refresh_token = $9, 
		    updated_at = $10, version = version + 1
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Exec(query,
		user.ID, user.Name, user.Phone, user.Email, user.Picture,
		user.UserType, user.IsActive, user.LastLoginAt, user.RefreshToken,
		user.UpdatedAt,
	)

	if err != nil {
		return shared.NewDatabaseError("failed to update user", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return shared.NewDatabaseError("failed to check affected rows", err)
	}

	if rowsAffected == 0 {
		return shared.NewNotFoundError("user not found for update", nil)
	}

	return nil
}

// Delete soft deletes a user
func (r *UserRepository) Delete(userID uuid.UUID) error {
	query := `UPDATE users SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Exec(query, userID)
	if err != nil {
		return shared.NewDatabaseError("failed to delete user", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return shared.NewDatabaseError("failed to check affected rows", err)
	}

	if rowsAffected == 0 {
		return shared.NewNotFoundError("user not found for deletion", nil)
	}

	return nil
}