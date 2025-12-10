package tests

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg/config"
	httpClient "github.com/complyance-io/complyance-go-sdk/v3/pkg/http"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/retry"
)

func TestRetryStrategy(t *testing.T) {
	// Create retry configuration
	retryConfig := &config.RetryConfig{
		MaxRetries:           3,
		BaseDelay:            10 * time.Millisecond, // Use small values for testing
		MaxDelay:             100 * time.Millisecond,
		JitterFactor:         0.1,
		CircuitBreakerEnabled: true,
		FailureThreshold:     5,
		CircuitBreakerTimeout: 100 * time.Millisecond,
		RetryableHTTPCodes:   []int{408, 429, 500, 502, 503, 504},
	}

	// Create retry strategy
	strategy := retry.NewStrategy(retryConfig)

	// Test successful execution
	t.Run("Successful Execution", func(t *testing.T) {
		var attempts int
		err := strategy.Do(context.Background(), func(ctx context.Context) error {
			attempts++
			return nil
		})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if attempts != 1 {
			t.Errorf("Expected 1 attempt, got %d", attempts)
		}

		if strategy.Metrics.GetAttempts() != 1 {
			t.Errorf("Expected 1 recorded attempt, got %d", strategy.Metrics.GetAttempts())
		}

		if strategy.Metrics.GetSuccesses() != 1 {
			t.Errorf("Expected 1 recorded success, got %d", strategy.Metrics.GetSuccesses())
		}
	})

	// Test retryable error
	t.Run("Retryable Error", func(t *testing.T) {
		var attempts int
		retryableErr := errors.New("retryable error")

		// Override IsRetryable function for this test
		originalIsRetryable := strategy.IsRetryable
		strategy.IsRetryable = func(err error) bool {
			return err == retryableErr
		}
		defer func() {
			strategy.IsRetryable = originalIsRetryable
		}()

		err := strategy.Do(context.Background(), func(ctx context.Context) error {
			attempts++
			if attempts <= 3 {
				return retryableErr
			}
			return nil
		})

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if attempts != 4 {
			t.Errorf("Expected 4 attempts, got %d", attempts)
		}

		if strategy.Metrics.GetFailures() < 3 {
			t.Errorf("Expected at least 3 recorded failures, got %d", strategy.Metrics.GetFailures())
		}
	})

	// Test max retries exceeded
	t.Run("Max Retries Exceeded", func(t *testing.T) {
		var attempts int
		retryableErr := errors.New("retryable error")

		// Override IsRetryable function for this test
		originalIsRetryable := strategy.IsRetryable
		strategy.IsRetryable = func(err error) bool {
			return err == retryableErr
		}
		defer func() {
			strategy.IsRetryable = originalIsRetryable
		}()

		err := strategy.Do(context.Background(), func(ctx context.Context) error {
			attempts++
			return retryableErr
		})

		if err == nil {
			t.Error("Expected error, got nil")
		}

		if attempts != retryConfig.MaxRetries+1 {
			t.Errorf("Expected %d attempts, got %d", retryConfig.MaxRetries+1, attempts)
		}
	})

	// Test context cancellation
	t.Run("Context Cancellation", func(t *testing.T) {
		var attempts int
		ctx, cancel := context.WithCancel(context.Background())
		
		err := strategy.Do(ctx, func(ctx context.Context) error {
			attempts++
			if attempts == 2 {
				cancel()
			}
			return errors.New("retryable error")
		})

		if err == nil {
			t.Error("Expected error, got nil")
		}

		if attempts != 2 {
			t.Errorf("Expected 2 attempts, got %d", attempts)
		}
	})
}

func TestCircuitBreaker(t *testing.T) {
	// Create circuit breaker
	cb := retry.NewCircuitBreaker(3, 100*time.Millisecond)

	// Test initial state
	t.Run("Initial State", func(t *testing.T) {
		if cb.IsOpen() {
			t.Error("Expected circuit to be closed initially")
		}

		if cb.GetState() != retry.CircuitClosed {
			t.Errorf("Expected CircuitClosed state, got %v", cb.GetState())
		}

		if cb.GetFailureCount() != 0 {
			t.Errorf("Expected 0 failures, got %d", cb.GetFailureCount())
		}
	})

	// Test opening circuit
	t.Run("Opening Circuit", func(t *testing.T) {
		// Record failures to reach threshold
		cb.RecordFailure()
		cb.RecordFailure()
		cb.RecordFailure()

		if !cb.IsOpen() {
			t.Error("Expected circuit to be open after threshold failures")
		}

		if cb.GetState() != retry.CircuitOpen {
			t.Errorf("Expected CircuitOpen state, got %v", cb.GetState())
		}

		if cb.GetFailureCount() != 3 {
			t.Errorf("Expected 3 failures, got %d", cb.GetFailureCount())
		}
	})

	// Test half-open state
	t.Run("Half-Open State", func(t *testing.T) {
		// Wait for timeout
		time.Sleep(cb.GetTimeout() + 10*time.Millisecond)

		// Circuit should transition to half-open on next check
		if cb.IsOpen() {
			t.Error("Expected circuit to be half-open after timeout")
		}

		if cb.GetState() != retry.CircuitHalfOpen {
			t.Errorf("Expected CircuitHalfOpen state, got %v", cb.GetState())
		}
	})

	// Test closing circuit
	t.Run("Closing Circuit", func(t *testing.T) {
		// Record success to close circuit
		cb.RecordSuccess()

		if cb.IsOpen() {
			t.Error("Expected circuit to be closed after success")
		}

		if cb.GetState() != retry.CircuitClosed {
			t.Errorf("Expected CircuitClosed state, got %v", cb.GetState())
		}

		if cb.GetFailureCount() != 0 {
			t.Errorf("Expected 0 failures after success, got %d", cb.GetFailureCount())
		}
	})

	// Test reset
	t.Run("Reset", func(t *testing.T) {
		// Open the circuit
		cb.RecordFailure()
		cb.RecordFailure()
		cb.RecordFailure()

		if !cb.IsOpen() {
			t.Error("Expected circuit to be open")
		}

		// Reset the circuit
		cb.Reset()

		if cb.IsOpen() {
			t.Error("Expected circuit to be closed after reset")
		}

		if cb.GetState() != retry.CircuitClosed {
			t.Errorf("Expected CircuitClosed state, got %v", cb.GetState())
		}

		if cb.GetFailureCount() != 0 {
			t.Errorf("Expected 0 failures after reset, got %d", cb.GetFailureCount())
		}
	})
}

func TestRetryMetrics(t *testing.T) {
	// Create metrics
	metrics := retry.NewMetrics()

	// Test initial state
	t.Run("Initial State", func(t *testing.T) {
		if metrics.GetAttempts() != 0 {
			t.Errorf("Expected 0 attempts, got %d", metrics.GetAttempts())
		}

		if metrics.GetSuccesses() != 0 {
			t.Errorf("Expected 0 successes, got %d", metrics.GetSuccesses())
		}

		if metrics.GetFailures() != 0 {
			t.Errorf("Expected 0 failures, got %d", metrics.GetFailures())
		}

		if metrics.GetCircuitOpens() != 0 {
			t.Errorf("Expected 0 circuit opens, got %d", metrics.GetCircuitOpens())
		}

		if !metrics.GetLastRetryTime().IsZero() {
			t.Errorf("Expected zero last retry time, got %v", metrics.GetLastRetryTime())
		}
	})

	// Test recording metrics
	t.Run("Recording Metrics", func(t *testing.T) {
		metrics.RecordAttempt()
		metrics.RecordAttempt()
		metrics.RecordSuccess()
		metrics.RecordFailure()
		metrics.RecordCircuitOpen()

		if metrics.GetAttempts() != 2 {
			t.Errorf("Expected 2 attempts, got %d", metrics.GetAttempts())
		}

		if metrics.GetSuccesses() != 1 {
			t.Errorf("Expected 1 success, got %d", metrics.GetSuccesses())
		}

		if metrics.GetFailures() != 1 {
			t.Errorf("Expected 1 failure, got %d", metrics.GetFailures())
		}

		if metrics.GetCircuitOpens() != 1 {
			t.Errorf("Expected 1 circuit open, got %d", metrics.GetCircuitOpens())
		}

		if metrics.GetLastRetryTime().IsZero() {
			t.Error("Expected non-zero last retry time")
		}
	})

	// Test reset
	t.Run("Reset", func(t *testing.T) {
		metrics.Reset()

		if metrics.GetAttempts() != 0 {
			t.Errorf("Expected 0 attempts after reset, got %d", metrics.GetAttempts())
		}

		if metrics.GetSuccesses() != 0 {
			t.Errorf("Expected 0 successes after reset, got %d", metrics.GetSuccesses())
		}

		if metrics.GetFailures() != 0 {
			t.Errorf("Expected 0 failures after reset, got %d", metrics.GetFailures())
		}

		if metrics.GetCircuitOpens() != 0 {
			t.Errorf("Expected 0 circuit opens after reset, got %d", metrics.GetCircuitOpens())
		}

		if !metrics.GetLastRetryTime().IsZero() {
			t.Errorf("Expected zero last retry time after reset, got %v", metrics.GetLastRetryTime())
		}
	})
}

func TestRetryMiddleware(t *testing.T) {
	// Create a test server that fails a few times then succeeds
	var requestCount int
	var mu sync.Mutex
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		count := requestCount
		mu.Unlock()

		if count <= 2 {
			// Fail with a 503 Service Unavailable
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Service Unavailable"))
			return
		}

		// Succeed on the third attempt
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	}))
	defer server.Close()

	// Create configuration with retry
	cfg := config.New(
		config.WithAPIKey("test_api_key"),
		config.WithBaseURL(server.URL),
		config.WithRetryConfig(&config.RetryConfig{
			MaxRetries:           3,
			BaseDelay:            10 * time.Millisecond,
			MaxDelay:             100 * time.Millisecond,
			JitterFactor:         0.1,
			CircuitBreakerEnabled: true,
			FailureThreshold:     5,
			CircuitBreakerTimeout: 100 * time.Millisecond,
			RetryableHTTPCodes:   []int{408, 429, 500, 502, 503, 504},
		}),
	)

	// Create HTTP client
	client := httpClient.NewClient(cfg)

	// Test retry on server error
	t.Run("Retry on Server Error", func(t *testing.T) {
		ctx := context.Background()
		resp, err := client.Get(ctx, "/test", nil)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		if string(resp.Body) != "Success" {
			t.Errorf("Expected body 'Success', got '%s'", string(resp.Body))
		}

		mu.Lock()
		if requestCount != 3 {
			t.Errorf("Expected 3 requests, got %d", requestCount)
		}
		mu.Unlock()
	})
}

func TestConcurrentCircuitBreaker(t *testing.T) {
	// Create circuit breaker
	cb := retry.NewCircuitBreaker(5, 100*time.Millisecond)

	// Test concurrent access
	t.Run("Concurrent Access", func(t *testing.T) {
		var wg sync.WaitGroup
		const goroutines = 10
		const operationsPerGoroutine = 100

		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func(id int) {
				defer wg.Done()
				for j := 0; j < operationsPerGoroutine; j++ {
					if id%2 == 0 {
						cb.RecordFailure()
					} else {
						cb.RecordSuccess()
					}
					
					// Read state
					_ = cb.IsOpen()
					_ = cb.GetState()
					_ = cb.GetFailureCount()
				}
			}(i)
		}

		wg.Wait()
		
		// No assertions needed - we're testing that there are no race conditions
		// If the test completes without panicking, it passes
	})
}