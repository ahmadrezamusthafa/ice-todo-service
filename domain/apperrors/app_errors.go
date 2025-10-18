package apperrors

import (
	"fmt"
	"net/http"
)

type ErrorType string

const (
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeValidation   ErrorType = "VALIDATION"
	ErrorTypeDatabase     ErrorType = "DATABASE"
	ErrorTypeInternal     ErrorType = "INTERNAL"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden    ErrorType = "FORBIDDEN"
	ErrorTypeConflict     ErrorType = "CONFLICT"
	ErrorTypeStorage      ErrorType = "STORAGE"
)

type AppError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Err     error     `json:"error,omitempty"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) StatusCode() int {
	switch e.Type {
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeValidation:
		return http.StatusBadRequest
	case ErrorTypeDatabase:
		return http.StatusInternalServerError
	case ErrorTypeInternal:
		return http.StatusInternalServerError
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeConflict:
		return http.StatusConflict
	case ErrorTypeStorage:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func NewNotFoundError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: message,
		Err:     err,
	}
}

func NewValidationError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Err:     err,
	}
}

func NewDatabaseError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeDatabase,
		Message: message,
		Err:     err,
	}
}

func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: message,
		Err:     err,
	}
}

func NewUnauthorizedError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeUnauthorized,
		Message: message,
		Err:     err,
	}
}

func NewForbiddenError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeForbidden,
		Message: message,
		Err:     err,
	}
}

func NewConflictError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeConflict,
		Message: message,
		Err:     err,
	}
}

func NewStorageError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeStorage,
		Message: message,
		Err:     err,
	}
}
