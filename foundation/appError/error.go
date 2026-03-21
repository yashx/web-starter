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

func InternalServerError(cause ...error) *AppError {
	e := &AppError{
		Code:       _internalServerErrorCode,
		Message:    "Internal server error",
		HttpStatus: http.StatusInternalServerError,
	}
	if len(cause) > 0 {
		e.Cause = cause[0]
	}
	return e
}

func BadRequestError(message string, cause ...error) *AppError {
	e := &AppError{
		Code:       _badRequestErrorCode,
		Message:    message,
		HttpStatus: http.StatusBadRequest,
	}
	if len(cause) > 0 {
		e.Cause = cause[0]
	}
	return e
}

func BadRequestErrorFromValidationError(err *shakTypes.ValidationError) *AppError {
	return &AppError{
		Code:       _badRequestErrorCode,
		Message:    err.Error(),
		HttpStatus: http.StatusBadRequest,
		Cause:      err,
	}
}

func InvalidStateError(message string, cause ...error) *AppError {
	e := &AppError{
		Code:       _invalidStateErrorCode,
		Message:    message,
		HttpStatus: http.StatusBadRequest,
	}
	if len(cause) > 0 {
		e.Cause = cause[0]
	}
	return e
}

func IsInternalServerError(err *AppError) bool {
	if err == nil {
		return false
	}
	return err.Code == _internalServerErrorCode
}
