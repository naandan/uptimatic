package utils

import "net/http"

type AppError struct {
	Code    string                      `json:"code"`
	Message string                      `json:"message"`
	Fields  map[string][]map[string]any `json:"fields,omitempty"`
	Status  int                         `json:"-"`
	Err     error                       `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func NewAppError(status int, code, message string, err error) *AppError {
	return &AppError{
		Status:  status,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func ValidationErrorErr(fields map[string][]map[string]any) *AppError {
	return &AppError{
		Code:    ValidationError,
		Message: "Payload validation failed",
		Fields:  fields,
		Status:  http.StatusUnprocessableEntity,
	}
}

func UniqueFieldError(field string) *AppError {
	return ValidationErrorErr(map[string][]map[string]any{
		field: {{"code": Unique}},
	})
}

func InternalServerError(message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, InternalError, message, err)
}

func ConflictError(message string, err error) *AppError {
	return NewAppError(http.StatusConflict, Conflict, message, err)
}
