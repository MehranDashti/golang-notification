package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int
	Message string
	Details interface{}
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code=%d message=%s: %v",
			e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code=%d message=%s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}

	return nil, false
}

func NotFound(message string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: message}
}

func BadRequest(message string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message}
}

func BadRequestWithDetails(message string, details interface{}) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		Details: details,
	}
}

func Internal(message string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}

func InternalWithDetails(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Err:     err,
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func Gone(message string) *AppError {
	return &AppError{
		Code:    http.StatusGone,
		Message: message,
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: message,
	}
}
