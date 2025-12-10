package models

// ErrorCode represents error classification codes
type ErrorCode string

const (
	// ErrorCodeConfigurationError represents configuration-related errors
	ErrorCodeConfigurationError ErrorCode = "CONFIGURATION_ERROR"

	// ErrorCodeValidationError represents input validation errors
	ErrorCodeValidationError ErrorCode = "VALIDATION_ERROR"

	// ErrorCodeNetworkError represents network and connectivity errors
	ErrorCodeNetworkError ErrorCode = "NETWORK_ERROR"

	// ErrorCodeAPIError represents API response errors
	ErrorCodeAPIError ErrorCode = "API_ERROR"

	// ErrorCodeAuthenticationError represents authentication errors
	ErrorCodeAuthenticationError ErrorCode = "AUTHENTICATION_ERROR"

	// ErrorCodeRateLimitError represents rate limiting errors
	ErrorCodeRateLimitError ErrorCode = "RATE_LIMIT_ERROR"

	// ErrorCodeServerError represents server-side errors
	ErrorCodeServerError ErrorCode = "SERVER_ERROR"

	// ErrorCodeUnknownError represents unclassified errors
	ErrorCodeUnknownError ErrorCode = "UNKNOWN_ERROR"
)

// String returns the string representation of the error code
func (ec ErrorCode) String() string {
	return string(ec)
}