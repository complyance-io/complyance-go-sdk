package models

// ErrorDetail represents detailed error information
type ErrorDetail struct {
	// Code is the error classification code
	Code ErrorCode `json:"code"`

	// Message is the human-readable error message
	Message string `json:"message"`

	// Suggestion is a recommended action to resolve the error
	Suggestion string `json:"suggestion,omitempty"`

	// Context contains additional error-specific information
	Context map[string]interface{} `json:"context,omitempty"`
}

// NewErrorDetail creates a new ErrorDetail with the provided values
func NewErrorDetail(code ErrorCode, message string) *ErrorDetail {
	return &ErrorDetail{
		Code:    code,
		Message: message,
	}
}

// WithSuggestion adds a suggestion to the error detail
func (e *ErrorDetail) WithSuggestion(suggestion string) *ErrorDetail {
	e.Suggestion = suggestion
	return e
}

// WithContext adds context to the error detail
func (e *ErrorDetail) WithContext(context map[string]interface{}) *ErrorDetail {
	e.Context = context
	return e
}

// AddContext adds a single context key-value pair
func (e *ErrorDetail) AddContext(key string, value interface{}) *ErrorDetail {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}