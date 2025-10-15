// Package application contains Data Transfer Objects (DTOs)
package application

import "time"

// BaseDTO represents the base structure for DTOs
type BaseDTO struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PaginationRequestDTO represents pagination request
type PaginationRequestDTO struct {
	Page     int    `json:"page" form:"page" binding:"min=1"`
	PageSize int    `json:"page_size" form:"page_size" binding:"min=1,max=100"`
	SortBy   string `json:"sort_by" form:"sort_by"`
	SortDir  string `json:"sort_dir" form:"sort_dir" binding:"oneof=asc desc"`
}

// PaginationResponseDTO represents pagination response
type PaginationResponseDTO[T any] struct {
	Items      []T `json:"items"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// LocationDTO represents geographical location
type LocationDTO struct {
	Latitude  float64 `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" binding:"required,min=-180,max=180"`
	Address   string  `json:"address" binding:"required,min=1"`
}

// ErrorResponseDTO represents error response
type ErrorResponseDTO struct {
	Error ErrorDetailDTO `json:"error"`
}

// ErrorDetailDTO represents error details
type ErrorDetailDTO struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// HealthCheckDTO represents health check response
type HealthCheckDTO struct {
	Status      string            `json:"status"`
	Timestamp   time.Time         `json:"timestamp"`
	Service     string            `json:"service"`
	Version     string            `json:"version"`
	Uptime      string            `json:"uptime"`
	Dependencies map[string]string `json:"dependencies"`
}

// NewPaginationRequestDTO creates pagination request with defaults
func NewPaginationRequestDTO() PaginationRequestDTO {
	return PaginationRequestDTO{
		Page:     1,
		PageSize: 20,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

// NewPaginationResponseDTO creates pagination response
func NewPaginationResponseDTO[T any](items []T, totalItems, page, pageSize int) PaginationResponseDTO[T] {
	totalPages := (totalItems + pageSize - 1) / pageSize
	if totalPages < 0 {
		totalPages = 0
	}
	
	return PaginationResponseDTO[T]{
		Items:      items,
		TotalItems: totalItems,
		TotalPages: totalPages,
		Page:       page,
		PageSize:   pageSize,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// NewErrorResponseDTO creates error response DTO
func NewErrorResponseDTO(code, message string, details map[string]interface{}) ErrorResponseDTO {
	return ErrorResponseDTO{
		Error: ErrorDetailDTO{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}