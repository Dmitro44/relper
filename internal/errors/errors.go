package errors

import (
	"fmt"
)

type ErrorCode string

const (
	ErrorCodeConflict     ErrorCode = "CONFLICT"
	ErrorCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrorCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrorCodeInternal     ErrorCode = "INTERNAL_ERROR"
	ErrorCodeBadRequest   ErrorCode = "BAD_REQUEST"
	ErrorCodeUnauthorized ErrorCode = "UNAUTHORIZED"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func NewInternalError(message string, err error) *AppError {
	return Wrap(ErrorCodeInternal, message, err)
}

func NewNotFoundError(message string) *AppError {
	return New(ErrorCodeNotFound, message)
}

func NewForbiddenError(message string) *AppError {
	return New(ErrorCodeForbidden, message)
}

func NewBadRequestError(message string) *AppError {
	return New(ErrorCodeBadRequest, message)
}

func NewUnauthorizedError(message string) *AppError {
	return New(ErrorCodeUnauthorized, message)
}

func NewConflictError(message string) *AppError {
	return New(ErrorCodeConflict, message)
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func Wrap(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
