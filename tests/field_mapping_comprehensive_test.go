package tests

import (
	"context"
	"testing"
	"time"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/config"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestFieldMappingComprehensive(t *testing.T) {
	// Create configuration
	cfg := config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
	)

	client, err := pkg.NewClient(cfg)
	assert.NoError(t, err)

	// Create a valid source
	source := models.NewSource("test-source", models.SourceTypeFirstParty, "Test Source")

	// Define test cases using table-driven test pattern
	testCases := []struct {
		name          string
		source        *models.Source
		country       string
		payload       map[string]interface{}
		expectError   bool
		errorContains string
	}{
		{
			name:    "Valid Mapping Request",
			source:  source,
			country: "SA",
			payload: map[string]interface{}{
				"invoice_number": "INV-001",
				"issue_date":     "2023-01-01",
				"buyer": map[string]interface{}{
					"name": "Test Buyer",
				},
			},
			expectError: false,
		},
		{
			name:    "Missing Required Fields",
			source:  source,
			country: "SA",
			payload: map[string]interface{}{
				// Missing invoice_number
				"issue_date": "2023-01-01",
			},
			expectError:   false, // Mapping should still work even with missing fields
		},
		{
			name:    "Invalid Date Format",
			source:  source,
			country: "SA",
			payload: map[string]interface{}{
				"invoice_number": "INV-001",
				"issue_date":     "01/01/2023", // Wrong format
			},
			expectError:   false, // Mapping should identify the format issue
		},
		{
			name:    "Complex Nested Structure",
			source:  source,
			country: "SA",
			payload: map[string]interface{}{
				"invoice_number": "INV-001",
				"issue_date":     "2023-01-01",
				"buyer": map[string]interface{}{
					"name":    "Test Buyer",
					"address": map[string]interface{}{
						"street":  "123 Test St",
						"city":    "Test City",
						"country": "SA",
						"postal_code": "12345",
					},
					"tax_id": "123456789",
				},
				"seller": map[string]interface{}{
					"name":    "Test Seller",
					"address": map[string]interface{}{
						"street":  "456 Test Ave",
						"city":    "Seller City",
						"country": "SA",
						"postal_code": "54321",
					},
					"tax_id": "987654321",
				},
				"line_items": []map[string]interface{}{
					{
						"name":        "Item 1",
						"quantity":    1,
						"price":       100.00,
						"tax_rate":    15.0,
						"tax_amount":  15.0,
						"total":       115.0,
					},
					{
						"name":        "Item 2",
						"quantity":    2,
						"price":       50.00,
						"tax_rate":    15.0,
						"tax_amount":  15.0,
						"total":       115.0,
					},
				},
				"totals": map[string]interface{}{
					"subtotal":   200.00,
					"tax_total":  30.00,
					"total":      230.00,
				},
			},
			expectError: false,
		},
		{
			name:    "Nil Source",
			source:  nil,
			country: "SA",
			payload: map[string]interface{}{
				"invoice_number": "INV-001",
			},
			expectError:   true,
			errorContains: "source cannot be nil",
		},
		{
			name:    "Empty Country",
			source:  source,
			country: "",
			payload: map[string]interface{}{
				"invoice_number": "INV-001",
			},
			expectError:   true,
			errorContains: "country is required",
		},
		{
			name:    "Invalid Country Format",
			source:  source,
			country: "USA", // Should be 2 letters
			payload: map[string]interface{}{
				"invoice_number": "INV-001",
			},
			expectError:   true,
			errorContains: "country must be a 2-letter ISO code",
		},
		{
			name:    "Nil Payload",
			source:  source,
			country: "SA",
			payload: nil,
			expectError:   true,
			errorContains: "payload cannot be nil",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.CreateMapping(ctx, tc.source, tc.country, tc.payload)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, "success", resp.Status)
				
				// For mapping requests, check that the response contains mapping data
				if resp.Data != nil {
					_, hasMapping := resp.Data["mapping"]
					assert.True(t, hasMapping, "Response should contain mapping data")
				}
			}
		})
	}
}