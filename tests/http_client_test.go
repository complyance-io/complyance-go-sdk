package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg/config"
	httpClient "github.com/complyance-io/complyance-go-sdk/v3/pkg/http"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
)

func TestHTTPClient(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request headers
		if r.Header.Get("X-API-Key") == "" {
			t.Error("API key header not set")
		}
		if r.Header.Get("X-Request-Timestamp") == "" {
			t.Error("Timestamp header not set")
		}
		if r.Header.Get("X-Request-Signature") == "" {
			t.Error("Signature header not set")
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"status":  "success",
			"message": "Request processed successfully",
			"data": map[string]interface{}{
				"submission_id": "test_123",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create configuration
	cfg := config.New(
		config.WithAPIKey("test_api_key"),
		config.WithEnvironment(models.EnvironmentSandbox),
		config.WithBaseURL(server.URL), // Use test server URL
	)

	// Create HTTP client
	client := httpClient.NewClient(cfg)

	// Test GET request
	t.Run("GET Request", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		resp, err := client.Get(ctx, "/test", nil)
		if err != nil {
			t.Fatalf("GET request failed: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(resp.Body, &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if status, ok := response["status"].(string); !ok || status != "success" {
			t.Errorf("Expected status 'success', got %v", response["status"])
		}
	})

	// Test POST request
	t.Run("POST Request", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		body := map[string]interface{}{
			"test_key": "test_value",
		}

		resp, err := client.Post(ctx, "/test", body, nil)
		if err != nil {
			t.Fatalf("POST request failed: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(resp.Body, &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if status, ok := response["status"].(string); !ok || status != "success" {
			t.Errorf("Expected status 'success', got %v", response["status"])
		}
	})

	// Test context cancellation
	t.Run("Context Cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := client.Get(ctx, "/test", nil)
		if err == nil {
			t.Fatal("Expected error due to cancelled context, got nil")
		}
	})
}

func TestHTTPService(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the path
		switch r.URL.Path {
		case "/unify":
			// Handle PushToUnify request
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := models.UnifyResponse{
				Status:  "success",
				Message: "Request processed successfully",
				Data: map[string]interface{}{
					"submission_id": "test_123",
				},
			}
			json.NewEncoder(w).Encode(response)

		case "/status/test_123":
			// Handle GetStatus request
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := models.UnifyResponse{
				Status:  "success",
				Message: "Status retrieved successfully",
				Data: map[string]interface{}{
					"submission_id": "test_123",
					"status":        "PROCESSED",
				},
			}
			json.NewEncoder(w).Encode(response)

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create configuration
	cfg := config.New(
		config.WithAPIKey("test_api_key"),
		config.WithEnvironment(models.EnvironmentSandbox),
		config.WithBaseURL(server.URL), // Use test server URL
	)

	// Create HTTP service
	service := httpClient.NewService(cfg)

	// Test PushToUnify
	t.Run("PushToUnify", func(t *testing.T) {
		ctx := context.Background()
		source := &models.Source{
			ID:   "test_source",
			Type: models.SourceTypeFirstParty,
			Name: "Test Source",
		}
		request := models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, "US")
		request.WithPayload(map[string]interface{}{
			"test_key": "test_value",
		})

		resp, err := service.PushToUnify(ctx, request)
		if err != nil {
			t.Fatalf("PushToUnify failed: %v", err)
		}

		if resp.Status != "success" {
			t.Errorf("Expected status 'success', got %s", resp.Status)
		}

		submissionID, ok := resp.Data["submission_id"].(string)
		if !ok || submissionID != "test_123" {
			t.Errorf("Expected submission_id 'test_123', got %v", resp.Data["submission_id"])
		}
	})

	// Test GetStatus
	t.Run("GetStatus", func(t *testing.T) {
		ctx := context.Background()
		resp, err := service.GetStatus(ctx, "test_123")
		if err != nil {
			t.Fatalf("GetStatus failed: %v", err)
		}

		if resp.Status != "success" {
			t.Errorf("Expected status 'success', got %s", resp.Status)
		}

		status, ok := resp.Data["status"].(string)
		if !ok || status != "PROCESSED" {
			t.Errorf("Expected status 'PROCESSED', got %v", resp.Data["status"])
		}
	})
}