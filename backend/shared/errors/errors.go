package errors


import (
	"fmt"
	"net/http"
)

// Error codes
const (
	CodeInvalidInput     = "INVALID_INPUT"
	CodeNotFound         = "NOT_FOUND"
	CodeUserNotFound     = "USER_NOT_FOUND"
	CodeUnauthorized     = "UNAUTHORIZED"
	CodeForbidden        = "FORBIDDEN"
	CodeConflict         = "CONFLICT"
	CodeInternalError    = "INTERNAL_ERROR"
	CodeExternalError    = "EXTERNAL_ERROR"
	CodeValidationError  = "VALIDATION_ERROR"
	CodeInvalidOperation = "INVALID_OPERATION"
	CodeTimeout          = "TIMEOUT"
)

// AppError represents application errors
type AppError struct {
	code    string
	message string
	cause   error
}

// NewAppError creates a new application error
func NewAppError(code, message string, cause error) *AppError {
	return &AppError{
		code:    code,
		message: message,
		cause:   cause,
	}
}

func (e *AppError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s (%v)", e.code, e.message, e.cause)
	}
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

func (e *AppError) Code() string {
	return e.code
}

func (e *AppError) Message() string {
	return e.message
}

func (e *AppError) Cause() error {
	return e.cause
}

// HTTPStatusCode returns the HTTP status code for the error
func (e *AppError) HTTPStatusCode() int {
	switch e.code {
	case CodeInvalidInput, CodeValidationError:
		return http.StatusBadRequest
	case CodeNotFound, CodeUserNotFound:
		return http.StatusNotFound
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeConflict:
		return http.StatusConflict
	case CodeTimeout:
		return http.StatusRequestTimeout
	case CodeInvalidOperation:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// ValidationError represents validation errors with field details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "validation failed"
	}
	return fmt.Sprintf("validation failed: %s", ve[0].Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string, value interface{}) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	}
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError extracts AppError from an error
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return nil
}