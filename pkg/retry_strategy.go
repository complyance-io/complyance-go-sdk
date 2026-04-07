/*
Retry Strategy implementation matching Python SDK exactly.
*/
package complyancesdk

import (
	"log"
	"math"
	"math/rand"
	"time"
)

// RetryStrategy Retry strategy implementation matching Python SDK
type RetryStrategy struct {
	config *RetryConfig
}

// NewRetryStrategy creates a new retry strategy
func NewRetryStrategy(config *RetryConfig) *RetryStrategy {
	return &RetryStrategy{
		config: config,
	}
}

// Execute operation with retry logic
func (r *RetryStrategy) Execute(operation func() (interface{}, error), operationName string) (interface{}, error) {
	var lastError error

	for attempt := 0; attempt < r.config.MaxAttempts; attempt++ {
		log.Printf("Executing %s, attempt %d/%d", operationName, attempt+1, r.config.MaxAttempts)

		result, err := operation()
		if err == nil {
			if attempt > 0 {
				log.Printf("Operation %s succeeded after %d attempts", operationName, attempt+1)
			}
			return result, nil
		}

		lastError = err

		// Check if this error should be retried
		shouldRetry := false
		if sdkErr, ok := err.(*SDKError); ok && sdkErr.ErrorDetail != nil && sdkErr.ErrorDetail.Code != nil {
			shouldRetry = r.config.ShouldRetry(*sdkErr.ErrorDetail.Code)
		}

		// If this is the last attempt or error is not retryable, don't retry
		if attempt == r.config.MaxAttempts-1 || !shouldRetry {
			log.Printf("Operation %s failed after %d attempts: %v", operationName, attempt+1, err)
			break
		}

		// Calculate delay for next attempt
		delayMs := r.calculateDelay(attempt + 1)
		log.Printf("Operation %s failed (attempt %d), retrying in %fms: %v", operationName, attempt+1, delayMs, err)

		// Sleep before retry
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
	}

	// If we get here, all retries failed
	if sdkErr, ok := lastError.(*SDKError); ok {
		// Create max retries exceeded error
		maxRetriesError := NewErrorDetailWithCode(
			ErrorCodeMaxRetriesExceeded,
			"Operation failed after "+string(r.config.MaxAttempts)+" retry attempts",
		)
		maxRetriesError.Suggestion = &[]string{"Maximum retry attempts exceeded. Check your network connection and try again later"}[0]
		maxRetriesError.AddContextValue("maxAttempts", r.config.MaxAttempts)
		maxRetriesError.AddContextValue("originalError", sdkErr.String())
		return nil, NewSDKError(maxRetriesError)
	} else {
		return nil, lastError
	}
}

// calculateDelay Calculate delay for retry attempt with exponential backoff and jitter
func (r *RetryStrategy) calculateDelay(attempt int) float64 {
	if attempt <= 0 {
		return 0
	}

	// Calculate exponential backoff
	delay := math.Min(
		float64(r.config.MaxDelayMs),
		float64(r.config.BaseDelayMs)*math.Pow(r.config.BackoffMultiplier, float64(attempt-1)),
	)

	// Add jitter
	if r.config.JitterFactor > 0 {
		jitter := (rand.Float64()*2 - 1) * r.config.JitterFactor // Random between -jitterFactor and +jitterFactor
		delay = delay * (1 + jitter)
	}

	return math.Max(0, delay)
}
