package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg/config"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/errors"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/retry"
)

// RetryMiddleware implements retry logic for HTTP requests
type RetryMiddleware struct {
	strategy *retry.Strategy
	config   *config.RetryConfig
}

// NewRetryMiddleware creates a new retry middleware
func NewRetryMiddleware(cfg *config.RetryConfig) *RetryMiddleware {
	return &RetryMiddleware{
		strategy: retry.NewStrategy(cfg),
		config:   cfg,
	}
}

// Process adds retry logic to the request
func (m *RetryMiddleware) Process(req *Request) (*Request, error) {
	// Add retry middleware to the request context
	req.retryMiddleware = m
	return req, nil
}

// DoWithRetry executes an HTTP request with retry logic
func (m *RetryMiddleware) DoWithRetry(ctx context.Context, client Client, req *Request) (*Response, error) {
	var resp *Response
	var err error

	// Define the retryable function
	retryableFunc := func(ctx context.Context) error {
		resp, err = client.Do(ctx, req)
		if err != nil {
			return err
		}

		// Check if the status code is retryable
		if m.isRetryableStatusCode(resp.StatusCode) {
			return errors.NewNetworkError(
				"received retryable status code",
				nil,
			).AddContext("status_code", resp.StatusCode)
		}

		return nil
	}

	// Execute with retry
	err = m.strategy.Do(ctx, retryableFunc)
	if err != nil {
		return nil, err
	}

	// Add retry information to response headers
	if resp.Headers == nil {
		resp.Headers = make(http.Header)
	}
	resp.Headers.Set("X-Retry-Attempts", strconv.FormatInt(m.strategy.Metrics.GetAttempts(), 10))

	return resp, nil
}

// isRetryableStatusCode checks if the status code should trigger a retry
func (m *RetryMiddleware) isRetryableStatusCode(statusCode int) bool {
	for _, code := range m.config.RetryableHTTPCodes {
		if statusCode == code {
			return true
		}
	}
	return false
}

// GetMetrics returns the retry metrics
func (m *RetryMiddleware) GetMetrics() *retry.Metrics {
	return m.strategy.Metrics
}

// GetCircuitBreaker returns the circuit breaker
func (m *RetryMiddleware) GetCircuitBreaker() *retry.CircuitBreaker {
	return m.strategy.CircuitBreaker
}

// ResetMetrics resets the retry metrics
func (m *RetryMiddleware) ResetMetrics() {
	m.strategy.Metrics.Reset()
}

// ResetCircuitBreaker resets the circuit breaker
func (m *RetryMiddleware) ResetCircuitBreaker() {
	if m.strategy.CircuitBreaker != nil {
		m.strategy.CircuitBreaker.Reset()
	}
}