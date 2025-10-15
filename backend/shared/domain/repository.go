// Package domain contains repository interfaces following Repository pattern
package domain

import (
	"context"
)

// Repository represents the base repository interface
type Repository[T AggregateRoot] interface {
	Save(ctx context.Context, entity T) error
	FindByID(ctx context.Context, id string) (T, error)
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
}

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	SortBy   string `json:"sort_by"`
	SortDir  string `json:"sort_dir"`
}

// PaginatedResult represents paginated query result
type PaginatedResult[T any] struct {
	Items      []T `json:"items"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
}

// NewPaginationParams creates pagination parameters with defaults
func NewPaginationParams(page, pageSize int) *PaginationParams {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	return &PaginationParams{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

// GetOffset calculates the offset for database queries
func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// CalculateTotalPages calculates total pages from total items
func (p *PaginationParams) CalculateTotalPages(totalItems int) int {
	if totalItems == 0 {
		return 0
	}
	pages := totalItems / p.PageSize
	if totalItems%p.PageSize != 0 {
		pages++
	}
	return pages
}