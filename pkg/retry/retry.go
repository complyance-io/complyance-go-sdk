package retry

import (
	"context"
	"math"
	"math/rand"
	"time"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg/config"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/errors"
)

// RetryableFunc is a function that can be retried
type RetryableFunc func(ctx context.Context) error

// IsRetryable determines if an error should be retried
type IsRetryable func(err error) bool

// Strategy implements retry logic with exponential backoff and jitter
type Strategy struct {
	// Config holds the retry configuration
	Config *config.RetryConfig

	// CircuitBreaker is the circuit breaker instance
	CircuitBreaker *CircuitBreaker

	// Metrics tracks retry statistics
	Metrics *Metrics

	// IsRetryable is a function that determines if an error should be retried
	IsRetryable IsRetryable
}

// NewStrategy creates a new retry strategy with the provided configuration
func NewStrategy(cfg *config.RetryConfig) *Strategy {
	if cfg == nil {
		cfg = &config.RetryConfig{
			MaxRetries:           config.DefaultMaxRetries,
			BaseDelay:            config.DefaultBaseDelay,
			MaxDelay:             config.DefaultMaxDelay,
			JitterFactor:         config.DefaultJitterFactor,
			CircuitBreakerEnabled: true,
			FailureThreshold:     5,
			CircuitBreakerTimeout: 60 * time.Second,
			RetryableHTTPCodes:   []int{408, 429, 500, 502, 503, 504},
		}
	}

	var cb *CircuitBreaker
	if cfg.CircuitBreakerEnabled {
		cb = NewCircuitBreaker(cfg.FailureThreshold, cfg.CircuitBreakerTimeout)
	}

	return &Strategy{
		Config:        cfg,
		CircuitBreaker: cb,
		Metrics:       NewMetrics(),
		IsRetryable:   errors.IsRetryableError,
	}
}

// Do executes the provided function with retry logic
func (s *Strategy) Do(ctx context.Context, fn RetryableFunc) error {
	// Check if context is already canceled
	if ctx.Err() != nil {
		return errors.NewNetworkError("context canceled", ctx.Err())
	}

	// Check circuit breaker if enabled
	if s.CircuitBreaker != nil && s.CircuitBreaker.IsOpen() {
		s.Metrics.RecordCircuitOpen()
		return errors.NewNetworkError("circuit breaker is open", errors.ErrCircuitOpen).
			WithSuggestion("Wait for the circuit breaker to close or reset it manually")
	}

	var err error
	var attempt int

	// Initialize random number generator with current time
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Execute the function with retries
	for attempt = 0; attempt <= s.Config.MaxRetries; attempt++ {
		// Record attempt in metrics
		s.Metrics.RecordAttempt()

		// Execute the function
		err = fn(ctx)

		// If successful, record success and return
		if err == nil {
			s.Metrics.RecordSuccess()
			if s.CircuitBreaker != nil {
				s.CircuitBreaker.RecordSuccess()
			}
			return nil
		}

		// Record failure in metrics
		s.Metrics.RecordFailure()

		// Record failure in circuit breaker if enabled
		if s.CircuitBreaker != nil {
			s.CircuitBreaker.RecordFailure()
		}

		// Check if error is retryable
		if !s.IsRetryable(err) {
			return err
		}

		// Check if we've reached max retries
		if attempt >= s.Config.MaxRetries {
			break
		}

		// Check if context is canceled
		if ctx.Err() != nil {
			return errors.NewNetworkError("context canceled during retry", ctx.Err())
		}

		// Calculate delay with exponential backoff and jitter
		delay := s.calculateDelay(attempt, rnd)

		// Create a timer for the delay
		timer := time.NewTimer(delay)
		defer timer.Stop()

		// Wait for the delay or context cancellation
		select {
		case <-timer.C:
			// Continue with next attempt
		case <-ctx.Done():
			return errors.NewNetworkError("context canceled during retry delay", ctx.Err())
		}
	}

	// If we've exhausted all retries, return the last error
	return errors.NewNetworkError(
		"all retry attempts failed",
		err,
	).WithSuggestion("Check network connectivity or try again later").
		AddContext("attempts", attempt).
		AddContext("max_retries", s.Config.MaxRetries)
}

// calculateDelay computes the delay for the next retry attempt
func (s *Strategy) calculateDelay(attempt int, rnd *rand.Rand) time.Duration {
	// Calculate base delay with exponential backoff: baseDelay * 2^attempt
	backoff := float64(s.Config.BaseDelay) * math.Pow(2, float64(attempt))
	
	// Apply maximum delay cap
	if backoff > float64(s.Config.MaxDelay) {
		backoff = float64(s.Config.MaxDelay)
	}
	
	// Apply jitter: delay = backoff * (1 ± jitterFactor)
	jitter := backoff * s.Config.JitterFactor
	min := backoff - jitter
	max := backoff + jitter
	
	// Generate random delay within jitter range
	delay := min + rnd.Float64()*(max-min)
	
	return time.Duration(delay)
}

// WithIsRetryable sets the function that determines if an error should be retried
func (s *Strategy) WithIsRetryable(fn IsRetryable) *Strategy {
	s.IsRetryable = fn
	return s
}

// WithCircuitBreaker sets the circuit breaker
func (s *Strategy) WithCircuitBreaker(cb *CircuitBreaker) *Strategy {
	s.CircuitBreaker = cb
	return s
}

// WithMetrics sets the metrics collector
func (s *Strategy) WithMetrics(metrics *Metrics) *Strategy {
	s.Metrics = metrics
	return s
}