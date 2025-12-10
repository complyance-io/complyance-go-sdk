package errors

import (
	"errors"
	"fmt"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
)

// Common error types
var (
	ErrInvalidConfig      = errors.New("invalid configuration")
	ErrInvalidRequest     = errors.New("invalid request")
	ErrNetworkFailure     = errors.New("network failure")
	ErrAPIError           = errors.New("API error")
	ErrAuthenticationFail = errors.New("authentication failed")
	ErrRateLimitExceeded  = errors.New("rate limit exceeded")
	ErrServerError        = errors.New("server error")
	ErrCircuitOpen        = errors.New("circuit breaker is open")
	ErrContextCanceled    = errors.New("context canceled")
	ErrTimeout            = errors.New("request timed out")
)

// SDKError represents an error from the SDK
type SDKError struct {
	// Code is the error classification code
	Code models.ErrorCode

	// Message is the human-readable error message
	Message string

	// Suggestion is a recommended action to resolve the error
	Suggestion string

	// Context contains additional error-specific information
	Context map[string]interface{}

	// Err is the underlying error
	Err error
}

// Error returns the error message
func (e *SDKError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *SDKError) Unwrap() error {
	return e.Err
}

// Is reports whether the target error matches this error
func (e *SDKError) Is(target error) bool {
	if target == nil {
		return false
	}
	
	t, ok := target.(*SDKError)
	if !ok {
		return errors.Is(e.Err, target)
	}
	
	return e.Code == t.Code
}

// NewConfigError creates a new configuration error
func NewConfigError(message string, err error) *SDKError {
	return &SDKError{
		Code:    models.ErrorCodeConfigurationError,
		Message: message,
		Err:     err,
	}
}

// NewValidationError creates a new validation error
func NewValidationError(message string, err error) *SDKError {
	return &SDKError{
		Code:    models.ErrorCodeValidationError,
		Message: message,
		Err:     err,
	}
}

// NewNetworkError creates a new network error
func NewNetworkError(message string, err error) *SDKError {
	return &SDKError{
		Code:    models.ErrorCodeNetworkError,
		Message: message,
		Err:     err,
	}
}

// NewAPIError creates a new API error
func NewAPIError(message string, err error) *SDKError {
	return &SDKError{
		Code:    models.ErrorCodeAPIError,
		Message: message,
		Err:     err,
	}
}

// NewAuthError creates a new authentication error
func NewAuthError(message string, err error) *SDKError {
	return &SDKError{
		Code:    models.ErrorCodeAuthenticationError,
		Message: message,
		Err:     err,
	}
}

// NewRateLimitError creates a new rate limit error
func NewRateLimitError(message string, err error) *SDKError {
	return &SDKError{
		Code:    models.ErrorCodeRateLimitError,
		Message: message,
		Err:     err,
	}
}

// NewServerError creates a new server error
func NewServerError(message string, err error) *SDKError {
	return &SDKError{
		Code:    models.ErrorCodeServerError,
		Message: message,
		Err:     err,
	}
}

// NewUnknownError creates a new unknown error
func NewUnknownError(message string, err error) *SDKError {
	return &SDKError{
		Code:    models.ErrorCodeUnknownError,
		Message: message,
		Err:     err,
	}
}

// WithSuggestion adds a suggestion to the error
func (e *SDKError) WithSuggestion(suggestion string) *SDKError {
	e.Suggestion = suggestion
	return e
}

// WithContext adds context to the error
func (e *SDKError) WithContext(context map[string]interface{}) *SDKError {
	e.Context = context
	return e
}

// AddContext adds a single context key-value pair
func (e *SDKError) AddContext(key string, value interface{}) *SDKError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// IsConfigError returns true if the error is a configuration error
func IsConfigError(err error) bool {
	var sdkErr *SDKError
	if errors.As(err, &sdkErr) {
		return sdkErr.Code == models.ErrorCodeConfigurationError
	}
	return false
}

// IsValidationError returns true if the error is a validation error
func IsValidationError(err error) bool {
	var sdkErr *SDKError
	if errors.As(err, &sdkErr) {
		return sdkErr.Code == models.ErrorCodeValidationError
	}
	return false
}

// IsNetworkError returns true if the error is a network error
func IsNetworkError(err error) bool {
	var sdkErr *SDKError
	if errors.As(err, &sdkErr) {
		return sdkErr.Code == models.ErrorCodeNetworkError
	}
	return false
}

// IsAPIError returns true if the error is an API error
func IsAPIError(err error) bool {
	var sdkErr *SDKError
	if errors.As(err, &sdkErr) {
		return sdkErr.Code == models.ErrorCodeAPIError
	}
	return false
}

// IsAuthError returns true if the error is an authentication error
func IsAuthError(err error) bool {
	var sdkErr *SDKError
	if errors.As(err, &sdkErr) {
		return sdkErr.Code == models.ErrorCodeAuthenticationError
	}
	return false
}

// IsRateLimitError returns true if the error is a rate limit error
func IsRateLimitError(err error) bool {
	var sdkErr *SDKError
	if errors.As(err, &sdkErr) {
		return sdkErr.Code == models.ErrorCodeRateLimitError
	}
	return false
}

// IsServerError returns true if the error is a server error
func IsServerError(err error) bool {
	var sdkErr *SDKError
	if errors.As(err, &sdkErr) {
		return sdkErr.Code == models.ErrorCodeServerError
	}
	return false
}

// IsRetryableError returns true if the error is retryable
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, ErrNetworkFailure) || 
	   errors.Is(err, ErrTimeout) || 
	   errors.Is(err, ErrServerError) || 
	   errors.Is(err, ErrRateLimitExceeded) {
		return true
	}

	var sdkErr *SDKError
	if errors.As(err, &sdkErr) {
		switch sdkErr.Code {
		case models.ErrorCodeNetworkError, 
			 models.ErrorCodeServerError, 
			 models.ErrorCodeRateLimitError:
			return true
		}
	}

	return false
}