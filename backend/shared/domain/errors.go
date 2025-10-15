// Package domain contains common error definitions
package domain

import "fmt"

// DomainError represents domain-specific errors
type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewDomainError creates a new domain error
func NewDomainError(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// WithDetails adds details to the domain error
func (e *DomainError) WithDetails(key string, value interface{}) *DomainError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// Common domain errors
var (
	ErrNotFound          = NewDomainError("NOT_FOUND", "Resource not found")
	ErrUnauthorized      = NewDomainError("UNAUTHORIZED", "Unauthorized access")
	ErrForbidden         = NewDomainError("FORBIDDEN", "Access forbidden")
	ErrValidation        = NewDomainError("VALIDATION_ERROR", "Validation failed")
	ErrConflict          = NewDomainError("CONFLICT", "Resource conflict")
	ErrInternalError     = NewDomainError("INTERNAL_ERROR", "Internal server error")
	ErrBadRequest        = NewDomainError("BAD_REQUEST", "Bad request")
	ErrServiceUnavailable = NewDomainError("SERVICE_UNAVAILABLE", "Service unavailable")
)