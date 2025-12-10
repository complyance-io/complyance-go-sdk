package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg/config"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
)

// MockServer creates a test server that returns predefined responses
type MockServer struct {
	Server *httptest.Server
}

// NewMockServer creates a new mock server with default handlers
func NewMockServer() *MockServer {
	mock := &MockServer{}
	
	// Create a test server
	mock.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Default handler for all requests
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		// Different responses based on path
		switch r.URL.Path {
		case "/unify":
			response := models.UnifyResponse{
				Status:  "success",
				Message: "Request processed successfully",
				Data: map[string]interface{}{
					"submission_id": "test_123",
				},
			}
			json.NewEncoder(w).Encode(response)
			
		case "/status/test_123":
			response := models.UnifyResponse{
				Status:  "success",
				Message: "Status retrieved successfully",
				Data: map[string]interface{}{
					"submission_id": "test_123",
					"status":        "PROCESSED",
				},
			}
			json.NewEncoder(w).Encode(response)
			
		case "/mapping":
			response := models.UnifyResponse{
				Status:  "success",
				Message: "Mapping created successfully",
				Data: map[string]interface{}{
					"mapping": map[string]interface{}{
						"invoice_number": "invoice.number",
						"issue_date":     "invoice.date",
					},
				},
			}
			json.NewEncoder(w).Encode(response)
			
		default:
			// Generic success response
			response := map[string]interface{}{
				"status":  "success",
				"message": "Request processed successfully",
			}
			json.NewEncoder(w).Encode(response)
		}
	}))
	
	return mock
}

// Close closes the mock server
func (m *MockServer) Close() {
	if m.Server != nil {
		m.Server.Close()
	}
}

// URL returns the URL of the mock server
func (m *MockServer) URL() string {
	return m.Server.URL
}

// CreateTestConfig creates a test configuration with the mock server URL
func (m *MockServer) CreateTestConfig() *config.Config {
	return config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
		config.WithBaseURL(m.URL()),
	)
}

// CreateTestSource creates a test source
func CreateTestSource() *models.Source {
	return models.NewSource("test-source", models.SourceTypeFirstParty, "Test Source")
}

// CreateTestPayload creates a test payload
func CreateTestPayload() map[string]interface{} {
	return map[string]interface{}{
		"invoice_number": "INV-001",
		"issue_date":     "2023-01-01",
		"buyer": map[string]interface{}{
			"name": "Test Buyer",
			"address": map[string]interface{}{
				"street":  "123 Test St",
				"city":    "Test City",
				"country": "SA",
			},
		},
		"seller": map[string]interface{}{
			"name": "Test Seller",
			"address": map[string]interface{}{
				"street":  "456 Test Ave",
				"city":    "Seller City",
				"country": "SA",
			},
		},
		"line_items": []map[string]interface{}{
			{
				"name":     "Item 1",
				"quantity": 1,
				"price":    100.00,
			},
		},
		"totals": map[string]interface{}{
			"subtotal": 100.00,
			"tax":      15.00,
			"total":    115.00,
		},
	}
}

// SkipIfShort skips the test if testing in short mode
func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}
}