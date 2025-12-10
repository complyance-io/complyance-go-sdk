/*
Exceptions for the Complyance SDK matching Python SDK exactly.
*/
package complyancesdk

// SDKError Main SDK error matching Python SDK
type SDKError struct {
	ErrorDetail *ErrorDetail
}

// NewSDKError creates a new SDK error
func NewSDKError(errorDetail *ErrorDetail) *SDKError {
	return &SDKError{
		ErrorDetail: errorDetail,
	}
}

// GetErrorDetail getter for error detail
func (s *SDKError) GetErrorDetail() *ErrorDetail {
	return s.ErrorDetail
}

// Error implements the error interface
func (s *SDKError) Error() string {
	if s.ErrorDetail != nil {
		return s.ErrorDetail.String()
	}
	return "Unknown SDK error"
}

// String string representation
func (s *SDKError) String() string {
	if s.ErrorDetail != nil {
		return s.ErrorDetail.String()
	}
	return "Unknown SDK error"
}

// ValidationError Validation error exception
type ValidationError struct {
	*SDKError
}

// NewValidationError creates a new validation error
func NewValidationError(message string, suggestion *string) *ValidationError {
	errorDetail := NewErrorDetailWithCode(ErrorCodeValidationFailed, message)
	if suggestion != nil {
		errorDetail.Suggestion = suggestion
	}
	return &ValidationError{
		SDKError: NewSDKError(errorDetail),
	}
}

// NetworkError Network error exception
type NetworkError struct {
	*SDKError
}

// NewNetworkError creates a new network error
func NewNetworkError(message string, suggestion *string) *NetworkError {
	errorDetail := NewErrorDetailWithCode(ErrorCodeNetworkError, message)
	if suggestion != nil {
		errorDetail.Suggestion = suggestion
	}
	errorDetail.Retryable = true
	return &NetworkError{
		SDKError: NewSDKError(errorDetail),
	}
}

// APIError API error exception
type APIError struct {
	*SDKError
}

// NewAPIError creates a new API error
func NewAPIError(message string, suggestion *string) *APIError {
	errorDetail := NewErrorDetailWithCode(ErrorCodeAPIError, message)
	if suggestion != nil {
		errorDetail.Suggestion = suggestion
	}
	return &APIError{
		SDKError: NewSDKError(errorDetail),
	}
}

// ConfigurationError Configuration error exception
type ConfigurationError struct {
	*SDKError
}

// NewConfigurationError creates a new configuration error
func NewConfigurationError(message string, suggestion *string) *ConfigurationError {
	errorDetail := NewErrorDetailWithCode(ErrorCodeMissingField, message)
	if suggestion != nil {
		errorDetail.Suggestion = suggestion
	}
	return &ConfigurationError{
		SDKError: NewSDKError(errorDetail),
	}
}
