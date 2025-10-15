// Package infrastructure provides HTTP utilities and middleware
package infrastructure

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/southern-martin/zride/backend/shared/application"
	"github.com/southern-martin/zride/backend/shared/domain"
)

// HTTPHandler provides common HTTP utilities
type HTTPHandler struct{}

// NewHTTPHandler creates new HTTP handler
func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{}
}

// WriteJSON writes JSON response
func (h *HTTPHandler) WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

// WriteError writes error response
func (h *HTTPHandler) WriteError(w http.ResponseWriter, statusCode int, err *domain.DomainError) {
	errorResponse := application.NewErrorResponseDTO(err.Code, err.Message, err.Details)
	h.WriteJSON(w, statusCode, errorResponse)
}

// WriteValidationError writes validation error response
func (h *HTTPHandler) WriteValidationError(w http.ResponseWriter, message string, details map[string]interface{}) {
	err := domain.ErrValidation.WithDetails("validation", details)
	err.Message = message
	h.WriteError(w, http.StatusBadRequest, err)
}

// ParsePagination parses pagination parameters from request
func (h *HTTPHandler) ParsePagination(r *http.Request) application.PaginationRequestDTO {
	pagination := application.NewPaginationRequestDTO()
	
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			pagination.Page = page
		}
	}
	
	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 && pageSize <= 100 {
			pagination.PageSize = pageSize
		}
	}
	
	if sortBy := r.URL.Query().Get("sort_by"); sortBy != "" {
		pagination.SortBy = sortBy
	}
	
	if sortDir := r.URL.Query().Get("sort_dir"); sortDir != "" {
		if strings.ToLower(sortDir) == "asc" || strings.ToLower(sortDir) == "desc" {
			pagination.SortDir = strings.ToLower(sortDir)
		}
	}
	
	return pagination
}

// ParseLocation parses location from request body
func (h *HTTPHandler) ParseLocation(r *http.Request) (*application.LocationDTO, error) {
	var location application.LocationDTO
	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		return nil, err
	}
	return &location, nil
}

// GetUserIDFromContext extracts user ID from request context
func (h *HTTPHandler) GetUserIDFromContext(r *http.Request) (string, error) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		return "", domain.ErrUnauthorized
	}
	return userID, nil
}

// SetCORSHeaders sets CORS headers
func (h *HTTPHandler) SetCORSHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
	}
}

// RequestValidator provides request validation utilities
type RequestValidator struct{}

// NewRequestValidator creates new request validator
func NewRequestValidator() *RequestValidator {
	return &RequestValidator{}
}

// ValidateRequired checks if required fields are present
func (v *RequestValidator) ValidateRequired(fields map[string]interface{}) error {
	missing := make([]string, 0)
	
	for field, value := range fields {
		if v.isEmpty(value) {
			missing = append(missing, field)
		}
	}
	
	if len(missing) > 0 {
		return domain.ErrValidation.WithDetails("missing_fields", missing)
	}
	
	return nil
}

// isEmpty checks if value is empty
func (v *RequestValidator) isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}
	
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) == ""
	case int, int32, int64:
		return v == 0
	case float32, float64:
		return v == 0
	default:
		return false
	}
}