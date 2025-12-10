package models

// UnifyResponse represents a response from the Complyance Unified API
type UnifyResponse struct {
	// Status is the response status (success, error)
	Status string `json:"status"`

	// Message is a human-readable response message
	Message string `json:"message"`

	// Data contains the response payload
	Data map[string]interface{} `json:"data,omitempty"`

	// Error contains error details if status is "error"
	Error *ErrorDetail `json:"error,omitempty"`

	// Metadata contains additional response information
	Metadata *ResponseMetadata `json:"metadata,omitempty"`
}

// IsSuccess returns true if the response status is "success"
func (r *UnifyResponse) IsSuccess() bool {
	return r.Status == "success"
}

// IsError returns true if the response status is "error"
func (r *UnifyResponse) IsError() bool {
	return r.Status == "error"
}

// GetErrorCode returns the error code if the response is an error
func (r *UnifyResponse) GetErrorCode() ErrorCode {
	if r.Error != nil {
		return r.Error.Code
	}
	return ""
}

// GetErrorMessage returns the error message if the response is an error
func (r *UnifyResponse) GetErrorMessage() string {
	if r.Error != nil {
		return r.Error.Message
	}
	return ""
}

// NewSuccessResponse creates a new successful response
func NewSuccessResponse(message string, data map[string]interface{}) *UnifyResponse {
	return &UnifyResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(errorDetail *ErrorDetail) *UnifyResponse {
	return &UnifyResponse{
		Status:  "error",
		Message: errorDetail.Message,
		Error:   errorDetail,
	}
}

// WithMetadata adds metadata to the response
func (r *UnifyResponse) WithMetadata(metadata *ResponseMetadata) *UnifyResponse {
	r.Metadata = metadata
	return r
}