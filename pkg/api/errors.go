package api

import "fmt"

// APIError is a structure for storing API error information
type APIError struct {
	StatusCode int    // HTTP status code
	ErrorCode  string // Custom error code
	Message    string // Error message
}

// Error method implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("[%d] %s: %s", e.StatusCode, e.ErrorCode, e.Message)
}

// NewAPIError creates an instance of APIError
func NewAPIError(statusCode int, errorCode, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
	}
}

// BadRequestError creates an APIError instance representing a BadRequest
func BadRequestError(message string) *APIError {
	return NewAPIError(400, "BadRequest", message)
}

// UnauthorizedError creates an APIError instance representing an Unauthorized
func UnauthorizedError(message string) *APIError {
	return NewAPIError(401, "Unauthorized", message)
}

// NotFoundError creates an APIError instance representing a NotFound
func NotFoundError(message string) *APIError {
	return NewAPIError(404, "NotFound", message)
}

// InternalServerError creates an APIError instance representing an InternalServer error
func InternalServerError(message string) *APIError {
	return NewAPIError(500, "InternalServerError", message)
}
