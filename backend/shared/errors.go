package shared

import (
	"fmt"
)

// Base error interface
type AppError interface {
	error
	Code() string
	Message() string
}

// ValidationError represents validation errors
type ValidationError struct {
	message string
	cause   error
}

func (e *ValidationError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("validation error: %s - %v", e.message, e.cause)
	}
	return fmt.Sprintf("validation error: %s", e.message)
}

func (e *ValidationError) Code() string {
	return "VALIDATION_ERROR"
}

func (e *ValidationError) Message() string {
	return e.message
}

func NewValidationError(message string, cause error) *ValidationError {
	return &ValidationError{message: message, cause: cause}
}

// NotFoundError represents not found errors
type NotFoundError struct {
	message string
	cause   error
}

func (e *NotFoundError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("not found: %s - %v", e.message, e.cause)
	}
	return fmt.Sprintf("not found: %s", e.message)
}

func (e *NotFoundError) Code() string {
	return "NOT_FOUND"
}

func (e *NotFoundError) Message() string {
	return e.message
}

func NewNotFoundError(message string, cause error) *NotFoundError {
	return &NotFoundError{message: message, cause: cause}
}

// ConflictError represents conflict errors
type ConflictError struct {
	message string
	cause   error
}

func (e *ConflictError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("conflict: %s - %v", e.message, e.cause)
	}
	return fmt.Sprintf("conflict: %s", e.message)
}

func (e *ConflictError) Code() string {
	return "CONFLICT"
}

func (e *ConflictError) Message() string {
	return e.message
}

func NewConflictError(message string, cause error) *ConflictError {
	return &ConflictError{message: message, cause: cause}
}

// DatabaseError represents database errors
type DatabaseError struct {
	message string
	cause   error
}

func (e *DatabaseError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("database error: %s - %v", e.message, e.cause)
	}
	return fmt.Sprintf("database error: %s", e.message)
}

func (e *DatabaseError) Code() string {
	return "DATABASE_ERROR"
}

func (e *DatabaseError) Message() string {
	return e.message
}

func NewDatabaseError(message string, cause error) *DatabaseError {
	return &DatabaseError{message: message, cause: cause}
}

// InternalError represents internal server errors
type InternalError struct {
	message string
	cause   error
}

func (e *InternalError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("internal error: %s - %v", e.message, e.cause)
	}
	return fmt.Sprintf("internal error: %s", e.message)
}

func (e *InternalError) Code() string {
	return "INTERNAL_ERROR"
}

func (e *InternalError) Message() string {
	return e.message
}

func NewInternalError(message string, cause error) *InternalError {
	return &InternalError{message: message, cause: cause}
}

// ExternalServiceError represents external service errors
type ExternalServiceError struct {
	message string
	cause   error
}

func (e *ExternalServiceError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("external service error: %s - %v", e.message, e.cause)
	}
	return fmt.Sprintf("external service error: %s", e.message)
}

func (e *ExternalServiceError) Code() string {
	return "EXTERNAL_SERVICE_ERROR"
}

func (e *ExternalServiceError) Message() string {
	return e.message
}

func NewExternalServiceError(message string, cause error) *ExternalServiceError {
	return &ExternalServiceError{message: message, cause: cause}
}