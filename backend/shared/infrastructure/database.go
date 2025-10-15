// Package infrastructure provides database utilities and configurationspackage infrastructure

package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/southern-martin/zride/backend/shared/domain"
	_ "github.com/lib/pq"
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	SSLMode  string
	MaxConns int
	MaxIdle  int
	ConnTTL  time.Duration
}

// NewDatabaseConfig creates database config with defaults
func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "zride",
		Username: "zride_user",
		Password: "zride_password",
		SSLMode:  "disable",
		MaxConns: 25,
		MaxIdle:  5,
		ConnTTL:  5 * time.Minute,
	}
}

// DSN returns database connection string
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.SSLMode,
	)
}

// Database represents database connection wrapper
type Database struct {
	db     *sql.DB
	config *DatabaseConfig
}

// NewDatabase creates new database connection
func NewDatabase(config *DatabaseConfig) (*Database, error) {
	db, err := sql.Open("postgres", config.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxConns)
	db.SetMaxIdleConns(config.MaxIdle)
	db.SetConnMaxLifetime(config.ConnTTL)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{
		db:     db,
		config: config,
	}, nil
}

// GetDB returns underlying sql.DB instance
func (d *Database) GetDB() *sql.DB {
	return d.db
}

// Close closes database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// Health checks database health
func (d *Database) Health(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// BaseRepository provides base repository implementation
type BaseRepository struct {
	db *Database
}

// NewBaseRepository creates base repository
func NewBaseRepository(db *Database) *BaseRepository {
	return &BaseRepository{db: db}
}

// GetDB returns the underlying database connection
func (r *BaseRepository) GetDB() *sql.DB {
	return r.db.GetDB()
}

// ExecuteInTransaction executes function within database transaction
func (r *BaseRepository) ExecuteInTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := r.db.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %w (original error: %v)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// BuildPaginationQuery builds pagination SQL query
func BuildPaginationQuery(baseQuery string, params *domain.PaginationParams) string {
	query := baseQuery
	
	if params.SortBy != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", params.SortBy, params.SortDir)
	}
	
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", params.PageSize, params.GetOffset())
	
	return query
}

// BuildCountQuery builds count query for pagination
func BuildCountQuery(baseQuery string) string {
	return fmt.Sprintf("SELECT COUNT(*) FROM (%s) as count_query", baseQuery)
}