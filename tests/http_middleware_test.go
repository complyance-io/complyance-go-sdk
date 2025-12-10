package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/config"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
)

// MockLogger implements the Logger interface for testing
type MockLogger struct {
	DebugCalls int
	InfoCalls  int
	ErrorCalls int
	LastMsg    string
}

func (l *MockLogger) Debug(msg string, fields map[string]interface{}) {
	l.DebugCalls++
	l.LastMsg = msg
}

func (l *MockLogger) Info(msg string, fields map[string]interface{}) {
	l.InfoCalls++
	l.LastMsg = msg
}

func (l *MockLogger) Error(msg string, fields map[string]interface{}) {
	l.ErrorCalls++
	l.LastMsg = msg
}

// TestServerMiddleware tests the HTTP middleware functionality
func TestServerMiddleware(t *testing.T) {
	// Create SDK configuration
	cfg := config.New().
		WithAPIKey("test-api-key").
		WithEnvironment(config.EnvironmentSandbox).
		WithTimeout(5 * time.Second)

	// Add a source
	source := &models.Source{
		ID:   "test-source",
		Type: models.SourceTypeFirstParty,
		Name: "Test Source",
	}
	cfg.AddSource(source)

	// Configure the SDK
	if err := pkg.Configure(cfg); err != nil {
		t.Fatalf("Failed to configure SDK: %v", err)
	}

	// Create a new server middleware
	middleware, err := pkg.NewServerMiddleware()
	if err != nil {
		t.Fatalf("Failed to create middleware: %v", err)
	}

	// Add logger to middleware
	logger := &MockLogger{}
	middleware.WithLogger(logger)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"success"}`))
	})

	// Create a test server
	server := httptest.NewServer(middleware.Handler(testHandler))
	defer server.Close()

	// Test basic handler functionality
	t.Run("BasicHandler", func(t *testing.T) {
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if result["status"] != "success" {
			t.Errorf("Expected status 'success', got %v", result["status"])
		}
	})

	// Test invoice processing
	t.Run("InvoiceProcessing", func(t *testing.T) {
		// Create a test server with the invoice processing middleware
		invoiceHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if response is in context
			response := pkg.GetResponse(r)
			if response == nil {
				t.Error("Expected response in context, got nil")
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"processed"}`))
		})

		invoiceServer := httptest.NewServer(middleware.ProcessInvoice(invoiceHandler))
		defer invoiceServer.Close()

		// Create a test invoice payload
		payload := `{
			"invoice_number": "TEST-001",
			"issue_date": "2023-01-01",
			"amount": 100.00
		}`

		// Make a request with country in query parameter
		req, err := http.NewRequest("POST", invoiceServer.URL+"?country=US", strings.NewReader(payload))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Send the request
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		
		// We expect an error because the middleware will try to make a real API call
		// which will fail in the test environment. This is expected behavior.
		if err == nil {
			defer resp.Body.Close()
			// If we got a response, it should be an error
			if resp.StatusCode < 400 {
				t.Errorf("Expected error status code, got %d", resp.StatusCode)
			}
		}

		// Verify that the logger was called
		if logger.ErrorCalls == 0 {
			t.Error("Expected logger.Error to be called")
		}
	})
}

// TestBatchProcessing tests the batch processing functionality
func TestBatchProcessing(t *testing.T) {
	// Create SDK configuration
	cfg := config.New().
		WithAPIKey("test-api-key").
		WithEnvironment(config.EnvironmentSandbox).
		WithTimeout(5 * time.Second)

	// Add a source
	source := &models.Source{
		ID:   "test-source",
		Type: models.SourceTypeFirstParty,
		Name: "Test Source",
	}
	cfg.AddSource(source)

	// Create a client
	client, err := pkg.NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create test requests
	requests := []*models.UnifyRequest{
		models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, "US").
			WithPayload(map[string]interface{}{"invoice_number": "INV-001"}),
		models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, "GB").
			WithPayload(map[string]interface{}{"invoice_number": "INV-002"}),
		models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, "FR").
			WithPayload(map[string]interface{}{"invoice_number": "INV-003"}),
	}

	// Process requests concurrently
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	responses, errors := client.BatchProcess(ctx, requests)

	// Verify results
	if len(responses) != len(requests) {
		t.Errorf("Expected %d responses, got %d", len(requests), len(responses))
	}
	if len(errors) != len(requests) {
		t.Errorf("Expected %d errors, got %d", len(requests), len(errors))
	}

	// All requests should fail because we're not making real API calls
	for i, err := range errors {
		if err == nil {
			t.Errorf("Expected error for request %d, got nil", i)
		}
	}
}

// TestErrorWrapping tests the error wrapping functionality
func TestErrorWrapping(t *testing.T) {
	// Create SDK configuration
	cfg := config.New().
		WithAPIKey("test-api-key").
		WithEnvironment(config.EnvironmentSandbox).
		WithTimeout(5 * time.Second)

	// Configure the SDK
	if err := pkg.Configure(cfg); err != nil {
		t.Fatalf("Failed to configure SDK: %v", err)
	}

	// Test with invalid request (nil)
	ctx := context.Background()
	_, err := pkg.PushToUnify(ctx, nil)
	
	if err == nil {
		t.Fatal("Expected error for nil request, got nil")
	}

	// Check error type and message
	if !strings.Contains(err.Error(), "request cannot be nil") {
		t.Errorf("Expected error message to contain 'request cannot be nil', got: %s", err.Error())
	}

	// Test with invalid country
	source := &models.Source{
		ID:   "test-source",
		Type: models.SourceTypeFirstParty,
		Name: "Test Source",
	}
	
	_, err = pkg.SubmitInvoice(ctx, source, "USA", map[string]interface{}{"test": "data"})
	
	if err == nil {
		t.Fatal("Expected error for invalid country, got nil")
	}

	// Check error type and message
	if !strings.Contains(err.Error(), "country must be a 2-letter ISO code") {
		t.Errorf("Expected error message to contain 'country must be a 2-letter ISO code', got: %s", err.Error())
	}
}