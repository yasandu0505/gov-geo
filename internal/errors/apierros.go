package apierrors

import "net/http"

// APIError defines a standard API error structure
type APIError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}

// Predefined API errors
var (
	ErrDepartmentNotFound = &APIError{Code: http.StatusNotFound, Message: "Department not found"}
	ErrMinistryNotFound   = &APIError{Code: http.StatusNotFound, Message: "Ministry not found"}
	ErrMissingField       = &APIError{Code: http.StatusBadRequest, Message: "Missing required field"}
	ErrInvalidInput       = &APIError{Code: http.StatusBadRequest, Message: "Invalid input"}
	ErrInternal           = &APIError{Code: http.StatusInternalServerError, Message: "Internal server error"}
)
