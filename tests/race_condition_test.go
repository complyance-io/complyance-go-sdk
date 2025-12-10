package tests

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/config"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/retry"
	"github.com/stretchr/testify/assert"
)

// TestConcurrentSDKUsage tests concurrent usage of the SDK
// This test is designed to be run with the race detector enabled:
// go test -race ./tests/race_condition_test.go
func TestConcurrentSDKUsage(t *testing.T) {
	// Create configuration
	cfg := config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
	)

	// Configure the SDK
	err := pkg.Configure(cfg)
	assert.NoError(t, err)

	// Create a valid source
	source := models.NewSource("test-source", models.SourceTypeFirstParty, "Test Source")

	// Create a valid payload
	payload := map[string]interface{}{
		"invoice_number": "INV-001",
		"issue_date":     "2023-01-01",
	}

	// Test concurrent access to static methods
	t.Run("ConcurrentStaticMethods", func(t *testing.T) {
		var wg sync.WaitGroup
		const goroutines = 10

		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func(id int) {
				defer wg.Done()
				ctx := context.Background()

				// Alternate between different methods to increase contention
				switch id % 3 {
				case 0:
					_, err := pkg.SubmitInvoice(ctx, source, "SA", payload)
					assert.NoError(t, err)
				case 1:
					_, err := pkg.CreateMapping(ctx, source, "SA", payload)
					assert.NoError(t, err)
				case 2:
					_, err := pkg.GetStatus(ctx, "test_123")
					assert.NoError(t, err)
				}
			}(i)
		}

		wg.Wait()
	})

	// Test concurrent client creation and usage
	t.Run("ConcurrentClientUsage", func(t *testing.T) {
		var wg sync.WaitGroup
		const goroutines = 10

		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				
				// Create a new client in each goroutine
				client, err := pkg.NewClient(cfg)
				assert.NoError(t, err)
				
				ctx := context.Background()
				_, err = client.SubmitInvoice(ctx, source, "SA", payload)
				assert.NoError(t, err)
			}()
		}

		wg.Wait()
	})

	// Test concurrent batch processing
	t.Run("ConcurrentBatchProcessing", func(t *testing.T) {
		var wg sync.WaitGroup
		const goroutines = 5
		const requestsPerBatch = 10

		// Create client
		client, err := pkg.NewClient(cfg)
		assert.NoError(t, err)

		// Create requests
		requests := make([]*models.UnifyRequest, requestsPerBatch)
		for i := 0; i < requestsPerBatch; i++ {
			req := models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, "SA")
			req.WithPayload(map[string]interface{}{
				"invoice_number": "INV-001",
				"issue_date":     "2023-01-01",
			})
			requests[i] = req
		}

		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()
				
				ctx := context.Background()
				responses, errs := client.BatchProcess(ctx, requests)
				
				assert.Equal(t, requestsPerBatch, len(responses))
				for _, err := range errs {
					assert.NoError(t, err)
				}
			}()
		}

		wg.Wait()
	})
}

// TestConcurrentRetryOperations tests concurrent operations on retry components
func TestConcurrentRetryOperations(t *testing.T) {
	// Create retry configuration
	retryConfig := &config.RetryConfig{
		MaxRetries:           3,
		BaseDelay:            10 * time.Millisecond,
		MaxDelay:             100 * time.Millisecond,
		JitterFactor:         0.1,
		CircuitBreakerEnabled: true,
		FailureThreshold:     5,
		CircuitBreakerTimeout: 100 * time.Millisecond,
		RetryableHTTPCodes:   []int{408, 429, 500, 502, 503, 504},
	}

	// Test concurrent access to retry strategy
	t.Run("ConcurrentRetryStrategy", func(t *testing.T) {
		strategy := retry.NewStrategy(retryConfig)
		var wg sync.WaitGroup
		const goroutines = 10

		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func(id int) {
				defer wg.Done()
				ctx := context.Background()

				// Simulate different retry scenarios
				err := strategy.Do(ctx, func(ctx context.Context) error {
					if id%3 == 0 {
						return nil // Success
					} else {
						// Fail once then succeed
						if strategy.Metrics.GetAttempts() == 0 {
							return &retry.RetryableError{
								Message: "Temporary error",
							}
						}
						return nil
					}
				})
				assert.NoError(t, err)
			}(i)
		}

		wg.Wait()
	})

	// Test concurrent access to circuit breaker
	t.Run("ConcurrentCircuitBreaker", func(t *testing.T) {
		cb := retry.NewCircuitBreaker(5, 100*time.Millisecond)
		var wg sync.WaitGroup
		const goroutines = 20
		const operationsPerGoroutine = 50

		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func(id int) {
				defer wg.Done()
				
				for j := 0; j < operationsPerGoroutine; j++ {
					// Mix operations to increase contention
					switch (id + j) % 4 {
					case 0:
						cb.RecordSuccess()
					case 1:
						cb.RecordFailure()
					case 2:
						_ = cb.IsOpen()
					case 3:
						_ = cb.GetState()
					}
				}
			}(i)
		}

		wg.Wait()
	})

	// Test concurrent access to retry metrics
	t.Run("ConcurrentRetryMetrics", func(t *testing.T) {
		metrics := retry.NewMetrics()
		var wg sync.WaitGroup
		const goroutines = 10
		const operationsPerGoroutine = 100

		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func(id int) {
				defer wg.Done()
				
				for j := 0; j < operationsPerGoroutine; j++ {
					// Mix operations to increase contention
					switch (id + j) % 5 {
					case 0:
						metrics.RecordAttempt()
					case 1:
						metrics.RecordSuccess()
					case 2:
						metrics.RecordFailure()
					case 3:
						metrics.RecordCircuitOpen()
					case 4:
						_ = metrics.GetAttempts()
						_ = metrics.GetSuccesses()
						_ = metrics.GetFailures()
					}
				}
			}(i)
		}

		wg.Wait()
	})
}

// TestConcurrentConfigModification tests concurrent modification of configuration
func TestConcurrentConfigModification(t *testing.T) {
	// Create base configuration
	cfg := config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
	)

	// Test concurrent modification
	t.Run("ConcurrentModification", func(t *testing.T) {
		var wg sync.WaitGroup
		const goroutines = 10

		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func(id int) {
				defer wg.Done()
				
				// Create a copy of the config to modify
				localCfg := config.New(
					config.WithAPIKey(cfg.APIKey),
					config.WithEnvironment(cfg.Environment),
				)
				
				// Modify different aspects of the config
				switch id % 5 {
				case 0:
					localCfg.WithTimeout(time.Duration(100+id) * time.Millisecond)
				case 1:
					localCfg.WithBaseURL("https://api.example.com/v" + string(id%10+48))
				case 2:
					source := models.NewSource("source-"+string(id%10+48), models.SourceTypeFirstParty, "Test Source")
					localCfg.WithSources([]*models.Source{source})
				case 3:
					retryConfig := &config.RetryConfig{
						MaxRetries: id%5 + 1,
						BaseDelay:  time.Duration(id%10+1) * time.Millisecond,
					}
					localCfg.WithRetryConfig(retryConfig)
				case 4:
					// Just validate
					err := localCfg.Validate()
					assert.NoError(t, err)
				}
				
				// Create a client with the modified config
				client, err := pkg.NewClient(localCfg)
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}(i)
		}

		wg.Wait()
	})
}