package shared

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound          = errors.New("record not found")
	ErrUnauthorized      = errors.New("unauthorized access")
	ErrForbidden         = errors.New("forbidden resource")
	ErrConflict          = errors.New("record conflict")
	ErrBadRequest        = errors.New("bad request")
	ErrInsufficientStock = errors.New("insufficient product stock")
	ErrSingleSellerOnly  = errors.New("all items in order must be from the same seller")
	ErrInvalidStatus     = errors.New("invalid status transition")
)

type AppError struct {
	Code       string
	Message    string
	StatusCode int
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewAppError(code string, message string, statusCode int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

