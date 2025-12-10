package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/config"
	httpClient "github.com/complyance-io/complyance-go-sdk/v3/pkg/http"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/retry"
)

// BenchmarkSDKClient benchmarks the SDK client performance
func BenchmarkSDKClient(b *testing.B) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success","message":"Request processed successfully","data":{"submission_id":"test_123"}}`))
	}))
	defer server.Close()

	// Create configuration
	cfg := config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
		config.WithBaseURL(server.URL),
	)

	// Create client
	client, err := pkg.NewClient(cfg)
	if err != nil {
		b.Fatalf("Failed to create client: %v", err)
	}

	// Create source
	source := models.NewSource("test-source", models.SourceTypeFirstParty, "Test Source")

	// Create payload
	payload := map[string]interface{}{
		"invoice_number": "INV-001",
		"issue_date":     "2023-01-01",
	}

	// Benchmark SubmitInvoice
	b.Run("SubmitInvoice", func(b *testing.B) {
		ctx := context.Background()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := client.SubmitInvoice(ctx, source, "SA", payload)
			if err != nil {
				b.Fatalf("SubmitInvoice failed: %v", err)
			}
		}
	})

	// Benchmark CreateMapping
	b.Run("CreateMapping", func(b *testing.B) {
		ctx := context.Background()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := client.CreateMapping(ctx, source, "SA", payload)
			if err != nil {
				b.Fatalf("CreateMapping failed: %v", err)
			}
		}
	})

	// Benchmark GetStatus
	b.Run("GetStatus", func(b *testing.B) {
		ctx := context.Background()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := client.GetStatus(ctx, "test_123")
			if err != nil {
				b.Fatalf("GetStatus failed: %v", err)
			}
		}
	})
}

// BenchmarkRetryStrategy benchmarks the retry strategy performance
func BenchmarkRetryStrategy(b *testing.B) {
	// Create retry configuration
	retryConfig := &config.RetryConfig{
		MaxRetries:           3,
		BaseDelay:            1 * time.Millisecond, // Use small values for benchmarking
		MaxDelay:             10 * time.Millisecond,
		JitterFactor:         0.1,
		CircuitBreakerEnabled: true,
		FailureThreshold:     5,
		CircuitBreakerTimeout: 10 * time.Millisecond,
		RetryableHTTPCodes:   []int{408, 429, 500, 502, 503, 504},
	}

	// Benchmark successful execution
	b.Run("SuccessfulExecution", func(b *testing.B) {
		strategy := retry.NewStrategy(retryConfig)
		ctx := context.Background()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			err := strategy.Do(ctx, func(ctx context.Context) error {
				return nil // Always succeed
			})
			if err != nil {
				b.Fatalf("Retry strategy failed: %v", err)
			}
		}
	})

	// Benchmark with retries
	b.Run("WithRetries", func(b *testing.B) {
		strategy := retry.NewStrategy(retryConfig)
		ctx := context.Background()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			var attempts int
			err := strategy.Do(ctx, func(ctx context.Context) error {
				attempts++
				if attempts < 3 {
					return &httpClient.ResponseError{
						StatusCode: 503,
						Message:    "Service Unavailable",
					}
				}
				return nil
			})
			if err != nil {
				b.Fatalf("Retry strategy failed: %v", err)
			}
		}
	})

	// Benchmark circuit breaker
	b.Run("CircuitBreaker", func(b *testing.B) {
		cb := retry.NewCircuitBreaker(5, 10*time.Millisecond)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			// Reset circuit breaker state for each iteration
			cb.Reset()
			
			// Record some failures
			for j := 0; j < 3; j++ {
				cb.RecordFailure()
			}
			
			// Check if open
			_ = cb.IsOpen()
			
			// Record success
			cb.RecordSuccess()
		}
	})
}

// BenchmarkHTTPClient benchmarks the HTTP client performance
func BenchmarkHTTPClient(b *testing.B) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success","message":"Request processed successfully","data":{"submission_id":"test_123"}}`))
	}))
	defer server.Close()

	// Create configuration
	cfg := config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
		config.WithBaseURL(server.URL),
	)

	// Create HTTP client
	client := httpClient.NewClient(cfg)

	// Benchmark GET request
	b.Run("GETRequest", func(b *testing.B) {
		ctx := context.Background()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := client.Get(ctx, "/test", nil)
			if err != nil {
				b.Fatalf("GET request failed: %v", err)
			}
		}
	})

	// Benchmark POST request
	b.Run("POSTRequest", func(b *testing.B) {
		ctx := context.Background()
		body := map[string]interface{}{
			"test_key": "test_value",
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := client.Post(ctx, "/test", body, nil)
			if err != nil {
				b.Fatalf("POST request failed: %v", err)
			}
		}
	})
}

// BenchmarkBatchProcessing benchmarks batch processing performance
func BenchmarkBatchProcessing(b *testing.B) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success","message":"Request processed successfully","data":{"submission_id":"test_123"}}`))
	}))
	defer server.Close()

	// Create configuration
	cfg := config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
		config.WithBaseURL(server.URL),
	)

	// Create client
	client, err := pkg.NewClient(cfg)
	if err != nil {
		b.Fatalf("Failed to create client: %v", err)
	}

	// Create source
	source := models.NewSource("test-source", models.SourceTypeFirstParty, "Test Source")

	// Benchmark with different batch sizes
	batchSizes := []int{1, 5, 10, 20, 50}
	for _, size := range batchSizes {
		b.Run(fmt.Sprintf("BatchSize_%d", size), func(b *testing.B) {
			// Create batch requests
			requests := make([]*models.UnifyRequest, size)
			for i := 0; i < size; i++ {
				req := models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, "SA")
				req.WithPayload(map[string]interface{}{
					"invoice_number": fmt.Sprintf("INV-%04d", i),
					"issue_date":     "2023-01-01",
				})
				requests[i] = req
			}

			ctx := context.Background()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				responses, errs := client.BatchProcess(ctx, requests)
				if len(responses) != size {
					b.Fatalf("Expected %d responses, got %d", size, len(responses))
				}
				for j, err := range errs {
					if err != nil {
						b.Fatalf("Request %d failed: %v", j, err)
					}
				}
			}
		})
	}
}