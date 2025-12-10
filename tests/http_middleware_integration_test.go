package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg/config"
	httpClient "github.com/complyance-io/complyance-go-sdk/v3/pkg/http"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestServerMiddlewareIntegration(t *testing.T) {
	// Create configuration
	cfg := config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
	)

	// Create a source
	source := models.NewSource("test-source", models.SourceTypeFirstParty, "Test Source")
	cfg.WithSources([]*models.Source{source})

	// Create middleware
	middleware := httpClient.NewServerMiddleware(cfg)

	// Create a test server that simulates the Complyance API
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request headers
		if r.Header.Get("X-API-Key") == "" {
			t.Error("API key header not set")
		}

		// Return a mock response
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
	}))
	defer apiServer.Close()

	// Update config to use test server
	cfg.WithBaseURL(apiServer.URL)

	// Test cases for the middleware
	testCases := []struct {
		name           string
		method         string
		path           string
		country        string
		payload        map[string]interface{}
		expectedStatus int
		expectError    bool
	}{
		{
			name:    "Valid Invoice Processing",
			method:  "POST",
			path:    "/api/invoices",
			country: "SA",
			payload: map[string]interface{}{
				"invoice_number": "INV-001",
				"issue_date":     "2023-01-01",
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:    "Missing Country",
			method:  "POST",
			path:    "/api/invoices",
			country: "", // Missing country
			payload: map[string]interface{}{
				"invoice_number": "INV-001",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:    "Invalid JSON Payload",
			method:  "POST",
			path:    "/api/invoices",
			country: "SA",
			payload: nil, // Will send invalid JSON
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test server with the middleware
			var payloadBytes []byte
			var err error
			
			if tc.payload != nil {
				payloadBytes, err = json.Marshal(tc.payload)
				assert.NoError(t, err)
			} else {
				// Invalid JSON for testing error handling
				payloadBytes = []byte("{invalid json")
			}

			// Create request
			req := httptest.NewRequest(tc.method, tc.path, bytes.NewReader(payloadBytes))
			req.Header.Set("Content-Type", "application/json")
			if tc.country != "" {
				req.Header.Set("X-Country", tc.country)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Create a handler that uses the middleware
			handler := middleware.ProcessInvoice(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check if response is in context
				resp := httpClient.GetResponse(r)
				assert.NotNil(t, resp)
				assert.Equal(t, "success", resp.Status)
			}))

			// Process the request
			handler.ServeHTTP(w, req)

			// Check response
			resp := w.Result()
			defer resp.Body.Close()

			// Read response body
			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			// Check status code
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			// Parse response
			var responseData map[string]interface{}
			err = json.Unmarshal(body, &responseData)
			assert.NoError(t, err)

			// Check for error response
			if tc.expectError {
				_, hasError := responseData["error"]
				assert.True(t, hasError, "Expected error in response")
			} else {
				assert.Equal(t, "success", responseData["status"])
			}
		})
	}
}

// TestMiddlewareWithRealHandler tests the middleware with a real HTTP handler
func TestMiddlewareWithRealHandler(t *testing.T) {
	// Create configuration
	cfg := config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
	)

	// Create a source
	source := models.NewSource("test-source", models.SourceTypeFirstParty, "Test Source")
	cfg.WithSources([]*models.Source{source})

	// Create middleware
	middleware := httpClient.NewServerMiddleware(cfg)

	// Create a test server that simulates the Complyance API
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return a mock response
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
	}))
	defer apiServer.Close()

	// Update config to use test server
	cfg.WithBaseURL(apiServer.URL)

	// Create a test HTTP server with the middleware
	testServer := httptest.NewServer(middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is the next handler after middleware
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Handler called"))
	})))
	defer testServer.Close()

	// Create HTTP client
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Make a request to the test server
	req, err := http.NewRequest("GET", testServer.URL+"/test", nil)
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	// Check that the handler was called
	assert.Equal(t, "Handler called", string(body))
}