package tests

import (
	"errors"
	"testing"

	sdkerrors "github.com/complyance-io/complyance-go-sdk/v3/pkg/errors"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestSDKError(t *testing.T) {
	// Create a validation error
	originalErr := errors.New("field validation failed")
	validationErr := sdkerrors.NewValidationError("Invalid invoice number", originalErr)

	// Test error message
	assert.Contains(t, validationErr.Error(), "Invalid invoice number")
	assert.Contains(t, validationErr.Error(), "field validation failed")

	// Test error code
	assert.Equal(t, models.ErrorCodeValidationError, validationErr.Code)

	// Test unwrap
	assert.Equal(t, originalErr, errors.Unwrap(validationErr))

	// Test error type checking
	assert.True(t, sdkerrors.IsValidationError(validationErr))
	assert.False(t, sdkerrors.IsNetworkError(validationErr))
	assert.False(t, sdkerrors.IsAPIError(validationErr))

	// Test adding suggestion and context
	validationErr.WithSuggestion("Provide a valid invoice number")
	validationErr.AddContext("field", "invoice_number")
	validationErr.AddContext("value", "")

	assert.Equal(t, "Provide a valid invoice number", validationErr.Suggestion)
	assert.Equal(t, "invoice_number", validationErr.Context["field"])
	assert.Equal(t, "", validationErr.Context["value"])

	// Test with context
	contextMap := map[string]interface{}{
		"request_id": "req_123456",
		"timestamp":  "2023-01-01T12:00:00Z",
	}
	validationErr.WithContext(contextMap)
	assert.Equal(t, "req_123456", validationErr.Context["request_id"])
	assert.Equal(t, "2023-01-01T12:00:00Z", validationErr.Context["timestamp"])
}

func TestErrorTypes(t *testing.T) {
	// Create different error types
	configErr := sdkerrors.NewConfigError("Missing API key", nil)
	validationErr := sdkerrors.NewValidationError("Invalid country code", nil)
	networkErr := sdkerrors.NewNetworkError("Connection failed", nil)
	apiErr := sdkerrors.NewAPIError("API returned error", nil)
	authErr := sdkerrors.NewAuthError("Invalid API key", nil)
	rateLimitErr := sdkerrors.NewRateLimitError("Too many requests", nil)
	serverErr := sdkerrors.NewServerError("Internal server error", nil)
	unknownErr := sdkerrors.NewUnknownError("Unknown error", nil)

	// Test error codes
	assert.Equal(t, models.ErrorCodeConfigurationError, configErr.Code)
	assert.Equal(t, models.ErrorCodeValidationError, validationErr.Code)
	assert.Equal(t, models.ErrorCodeNetworkError, networkErr.Code)
	assert.Equal(t, models.ErrorCodeAPIError, apiErr.Code)
	assert.Equal(t, models.ErrorCodeAuthenticationError, authErr.Code)
	assert.Equal(t, models.ErrorCodeRateLimitError, rateLimitErr.Code)
	assert.Equal(t, models.ErrorCodeServerError, serverErr.Code)
	assert.Equal(t, models.ErrorCodeUnknownError, unknownErr.Code)

	// Test error type checking
	assert.True(t, sdkerrors.IsConfigError(configErr))
	assert.True(t, sdkerrors.IsValidationError(validationErr))
	assert.True(t, sdkerrors.IsNetworkError(networkErr))
	assert.True(t, sdkerrors.IsAPIError(apiErr))
	assert.True(t, sdkerrors.IsAuthError(authErr))
	assert.True(t, sdkerrors.IsRateLimitError(rateLimitErr))
	assert.True(t, sdkerrors.IsServerError(serverErr))

	// Test error wrapping
	wrappedErr := sdkerrors.NewNetworkError("Connection timeout", sdkerrors.ErrTimeout)
	assert.True(t, errors.Is(wrappedErr, sdkerrors.ErrTimeout))
}

func TestRetryableErrors(t *testing.T) {
	// Test retryable errors
	assert.True(t, sdkerrors.IsRetryableError(sdkerrors.ErrNetworkFailure))
	assert.True(t, sdkerrors.IsRetryableError(sdkerrors.ErrTimeout))
	assert.True(t, sdkerrors.IsRetryableError(sdkerrors.ErrServerError))
	assert.True(t, sdkerrors.IsRetryableError(sdkerrors.ErrRateLimitExceeded))

	// Test non-retryable errors
	assert.False(t, sdkerrors.IsRetryableError(sdkerrors.ErrInvalidConfig))
	assert.False(t, sdkerrors.IsRetryableError(sdkerrors.ErrInvalidRequest))
	assert.False(t, sdkerrors.IsRetryableError(sdkerrors.ErrAuthenticationFail))
	assert.False(t, sdkerrors.IsRetryableError(sdkerrors.ErrCircuitOpen))
	assert.False(t, sdkerrors.IsRetryableError(nil))

	// Test SDK error types
	networkErr := sdkerrors.NewNetworkError("Connection failed", nil)
	serverErr := sdkerrors.NewServerError("Internal server error", nil)
	rateLimitErr := sdkerrors.NewRateLimitError("Too many requests", nil)
	validationErr := sdkerrors.NewValidationError("Invalid input", nil)
	configErr := sdkerrors.NewConfigError("Missing config", nil)

	assert.True(t, sdkerrors.IsRetryableError(networkErr))
	assert.True(t, sdkerrors.IsRetryableError(serverErr))
	assert.True(t, sdkerrors.IsRetryableError(rateLimitErr))
	assert.False(t, sdkerrors.IsRetryableError(validationErr))
	assert.False(t, sdkerrors.IsRetryableError(configErr))
}