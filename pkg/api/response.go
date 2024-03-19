package api

// Response is a basic structure for API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Create a successful response
func NewSuccessResponse(data interface{}) Response {
	return Response{
		Success: true,
		Message: "success",
		Data:    data,
	}
}

// Create an error response
func NewErrorResponse(message string) Response {
	return Response{
		Success: false,
		Message: message,
	}
}
