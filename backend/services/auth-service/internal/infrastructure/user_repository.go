// Package infrastructure provides PostgreSQL user repository implementationpackage infrastructure

package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/southern-martin/zride/backend/services/auth-service/internal/domain"
	"github.com/southern-martin/zride/backend/shared/infrastructure"
	sharedDomain "github.com/southern-martin/zride/backend/shared/domain"
	"github.com/google/uuid"
)

// PostgreSQLUserRepository implements UserRepository interface
type PostgreSQLUserRepository struct {
	*infrastructure.BaseRepository
}

// NewPostgreSQLUserRepository creates new PostgreSQL user repository
func NewPostgreSQLUserRepository(db *infrastructure.Database) domain.UserRepository {
	return &PostgreSQLUserRepository{
		BaseRepository: infrastructure.NewBaseRepository(db),
	}
}

// Save saves user to database
func (r *PostgreSQLUserRepository) Save(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, zalo_id, name, phone, email, avatar, is_active, last_login_at, refresh_token, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			phone = EXCLUDED.phone,
			email = EXCLUDED.email,
			avatar = EXCLUDED.avatar,
			is_active = EXCLUDED.is_active,
			last_login_at = EXCLUDED.last_login_at,
			refresh_token = EXCLUDED.refresh_token,
			version = EXCLUDED.version,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.GetDB().ExecContext(ctx, query,
		user.ID,
		user.ZaloID,
		user.Name,
		user.Phone,
		user.Email,
		user.Avatar,
		user.IsActive,
		user.LastLoginAt,
		user.RefreshToken,
		user.Version,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

// FindByID finds user by ID
func (r *PostgreSQLUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, sharedDomain.ErrBadRequest.WithDetails("invalid_user_id", id)
	}

	query := `
		SELECT id, zalo_id, name, phone, email, avatar, is_active, last_login_at, refresh_token, version, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_active = true
	`

	user := &domain.User{}
	var lastLoginAt sql.NullTime

	err = r.GetDB().QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.ZaloID,
		&user.Name,
		&user.Phone,
		&user.Email,
		&user.Avatar,
		&user.IsActive,
		&lastLoginAt,
		&user.RefreshToken,
		&user.Version,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, sharedDomain.ErrNotFound.WithDetails("user_id", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

// FindByZaloID finds user by Zalo ID
func (r *PostgreSQLUserRepository) FindByZaloID(ctx context.Context, zaloID string) (*domain.User, error) {
	query := `
		SELECT id, zalo_id, name, phone, email, avatar, is_active, last_login_at, refresh_token, version, created_at, updated_at
		FROM users
		WHERE zalo_id = $1 AND is_active = true
	`

	user := &domain.User{}
	var lastLoginAt sql.NullTime

	err := r.GetDB().QueryRowContext(ctx, query, zaloID).Scan(
		&user.ID,
		&user.ZaloID,
		&user.Name,
		&user.Phone,
		&user.Email,
		&user.Avatar,
		&user.IsActive,
		&lastLoginAt,
		&user.RefreshToken,
		&user.Version,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, sharedDomain.ErrNotFound.WithDetails("zalo_id", zaloID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by zalo_id: %w", err)
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

// FindByEmail finds user by email
func (r *PostgreSQLUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, zalo_id, name, phone, email, avatar, is_active, last_login_at, refresh_token, version, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_active = true
	`

	user := &domain.User{}
	var lastLoginAt sql.NullTime

	err := r.GetDB().QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.ZaloID,
		&user.Name,
		&user.Phone,
		&user.Email,
		&user.Avatar,
		&user.IsActive,
		&lastLoginAt,
		&user.RefreshToken,
		&user.Version,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, sharedDomain.ErrNotFound.WithDetails("email", email)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

// FindByPhone finds user by phone
func (r *PostgreSQLUserRepository) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	query := `
		SELECT id, zalo_id, name, phone, email, avatar, is_active, last_login_at, refresh_token, version, created_at, updated_at
		FROM users
		WHERE phone = $1 AND is_active = true
	`

	user := &domain.User{}
	var lastLoginAt sql.NullTime

	err := r.GetDB().QueryRowContext(ctx, query, phone).Scan(
		&user.ID,
		&user.ZaloID,
		&user.Name,
		&user.Phone,
		&user.Email,
		&user.Avatar,
		&user.IsActive,
		&lastLoginAt,
		&user.RefreshToken,
		&user.Version,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, sharedDomain.ErrNotFound.WithDetails("phone", phone)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by phone: %w", err)
	}

	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

// Delete deletes user by ID
func (r *PostgreSQLUserRepository) Delete(ctx context.Context, id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return sharedDomain.ErrBadRequest.WithDetails("invalid_user_id", id)
	}

	query := `UPDATE users SET is_active = false, updated_at = $1 WHERE id = $2`
	
	result, err := r.GetDB().ExecContext(ctx, query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sharedDomain.ErrNotFound.WithDetails("user_id", id)
	}

	return nil
}

// Exists checks if user exists
func (r *PostgreSQLUserRepository) Exists(ctx context.Context, id string) (bool, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return false, sharedDomain.ErrBadRequest.WithDetails("invalid_user_id", id)
	}

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND is_active = true)`
	
	var exists bool
	err = r.GetDB().QueryRowContext(ctx, query, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

// UpdateLastLogin updates user's last login timestamp
func (r *PostgreSQLUserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return sharedDomain.ErrBadRequest.WithDetails("invalid_user_id", userID)
	}

	query := `UPDATE users SET last_login_at = $1, updated_at = $2 WHERE id = $3`
	
	now := time.Now()
	_, err = r.GetDB().ExecContext(ctx, query, now, now, id)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// UpdateRefreshToken updates user's refresh token
func (r *PostgreSQLUserRepository) UpdateRefreshToken(ctx context.Context, userID, refreshToken string) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return sharedDomain.ErrBadRequest.WithDetails("invalid_user_id", userID)
	}

	query := `UPDATE users SET refresh_token = $1, updated_at = $2 WHERE id = $3`
	
	_, err = r.GetDB().ExecContext(ctx, query, refreshToken, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update refresh token: %w", err)
	}

	return nil
}

// FindActiveUsers finds active users with pagination
func (r *PostgreSQLUserRepository) FindActiveUsers(ctx context.Context, params *sharedDomain.PaginationParams) (*sharedDomain.PaginatedResult[*domain.User], error) {
	baseQuery := "SELECT id, zalo_id, name, phone, email, avatar, is_active, last_login_at, refresh_token, version, created_at, updated_at FROM users WHERE is_active = true"
	
	// Get total count
	countQuery := infrastructure.BuildCountQuery(baseQuery)
	var totalItems int
	err := r.GetDB().QueryRowContext(ctx, countQuery).Scan(&totalItems)
	if err != nil {
		return nil, fmt.Errorf("failed to get user count: %w", err)
	}

	// Get paginated results
	paginatedQuery := infrastructure.BuildPaginationQuery(baseQuery, params)
	rows, err := r.GetDB().QueryContext(ctx, paginatedQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		var lastLoginAt sql.NullTime

		err := rows.Scan(
			&user.ID,
			&user.ZaloID,
			&user.Name,
			&user.Phone,
			&user.Email,
			&user.Avatar,
			&user.IsActive,
			&lastLoginAt,
			&user.RefreshToken,
			&user.Version,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		if lastLoginAt.Valid {
			user.LastLoginAt = &lastLoginAt.Time
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate users: %w", err)
	}

	return &sharedDomain.PaginatedResult[*domain.User]{
		Items:      users,
		TotalItems: totalItems,
		TotalPages: params.CalculateTotalPages(totalItems),
		Page:       params.Page,
		PageSize:   params.PageSize,
	}, nil
}