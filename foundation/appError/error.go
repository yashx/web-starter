package appError

import (
	"net/http"

	shakTypes "github.com/yashx/shak/types"
)

type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Cause      error  `json:"-"`
	HttpStatus int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Code + " " + e.Message
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

const (
	_internalServerErrorCode = "E000"
	_badRequestErrorCode     = "E001"
	_invalidStateErrorCode   = "E002"
)

func InternalServerError() *AppError {
	return &AppError{
		Code:       _internalServerErrorCode,
		Message:    "Internal server error",
		HttpStatus: http.StatusInternalServerError,
	}
}

func InternalServerErrorWithCause(cause error) *AppError {
	return &AppError{
		Code:       _internalServerErrorCode,
		Message:    "Internal server error",
		HttpStatus: http.StatusInternalServerError,
		Cause:      cause,
	}
}

func BadRequestError(message string) *AppError {
	return &AppError{
		Code:       _badRequestErrorCode,
		Message:    message,
		HttpStatus: http.StatusBadRequest,
	}
}

func BadRequestErrorWithCause(message string, cause error) *AppError {
	return &AppError{
		Code:       _badRequestErrorCode,
		Message:    message,
		HttpStatus: http.StatusBadRequest,
		Cause:      cause,
	}
}

func BadRequestErrorFromValidationError(err *shakTypes.ValidationError) *AppError {
	return &AppError{
		Code:       _badRequestErrorCode,
		Message:    err.Error(),
		HttpStatus: http.StatusBadRequest,
		Cause:      err,
	}
}

func InvalidStateError(message string) *AppError {
	return &AppError{
		Code:       _invalidStateErrorCode,
		Message:    message,
		HttpStatus: http.StatusBadRequest,
	}
}

func InvalidStateErrorWithCause(message string, cause error) *AppError {
	return &AppError{
		Code:       _invalidStateErrorCode,
		Message:    message,
		HttpStatus: http.StatusBadRequest,
		Cause:      cause,
	}
}

func IsInternalServerError(err *AppError) bool {
	if err == nil {
		return false
	}
	return err.Code == _internalServerErrorCode
}
