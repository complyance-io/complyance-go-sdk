package tests

import (
	"context"
	"testing"
	"time"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/config"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/errors"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestSDKClientComprehensive(t *testing.T) {
	// Create mock server
	mockServer := NewMockServer()
	defer mockServer.Close()

	// Create configuration
	cfg := mockServer.CreateTestConfig()

	// Test cases for SDK configuration
	configTestCases := []struct {
		name        string
		config      *config.Config
		expectError bool
		errorType   string
	}{
		{
			name:        "Valid Configuration",
			config:      cfg,
			expectError: false,
		},
		{
			name:        "Nil Configuration",
			config:      nil,
			expectError: true,
			errorType:   "ConfigError",
		},
		{
			name: "Missing API Key",
			config: config.New(
				config.WithEnvironment(models.EnvironmentSandbox),
				config.WithBaseURL(mockServer.URL()),
			),
			expectError: true,
			errorType:   "ConfigError",
		},
		{
			name: "Missing Environment",
			config: config.New(
				config.WithAPIKey("test-api-key"),
				config.WithBaseURL(mockServer.URL()),
			),
			expectError: true,
			errorType:   "ConfigError",
		},
		{
			name: "Invalid Environment",
			config: config.New(
				config.WithAPIKey("test-api-key"),
				config.WithEnvironment("invalid"),
				config.WithBaseURL(mockServer.URL()),
			),
			expectError: true,
			errorType:   "ConfigError",
		},
		{
			name: "With Custom Timeout",
			config: config.New(
				config.WithAPIKey("test-api-key"),
				config.WithEnvironment(models.EnvironmentSandbox),
				config.WithBaseURL(mockServer.URL()),
				config.WithTimeout(30*time.Second),
			),
			expectError: false,
		},
		{
			name: "With Retry Config",
			config: config.New(
				config.WithAPIKey("test-api-key"),
				config.WithEnvironment(models.EnvironmentSandbox),
				config.WithBaseURL(mockServer.URL()),
				config.WithRetryConfig(&config.RetryConfig{
					MaxRetries:  3,
					BaseDelay:   100 * time.Millisecond,
					MaxDelay:    1 * time.Second,
					JitterFactor: 0.1,
				}),
			),
			expectError: false,
		},
	}

	// Run configuration test cases
	for _, tc := range configTestCases {
		t.Run("Configure_"+tc.name, func(t *testing.T) {
			err := pkg.Configure(tc.config)
			
			if tc.expectError {
				assert.Error(t, err)
				if tc.errorType != "" {
					switch tc.errorType {
					case "ConfigError":
						assert.True(t, errors.IsConfigError(err))
					case "ValidationError":
						assert.True(t, errors.IsValidationError(err))
					}
				}
			} else {
				assert.NoError(t, err)
			}
		})

		t.Run("NewClient_"+tc.name, func(t *testing.T) {
			client, err := pkg.NewClient(tc.config)
			
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
				if tc.errorType != "" {
					switch tc.errorType {
					case "ConfigError":
						assert.True(t, errors.IsConfigError(err))
					case "ValidationError":
						assert.True(t, errors.IsValidationError(err))
					}
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}

	// Create client for API tests
	client, err := pkg.NewClient(cfg)
	assert.NoError(t, err)

	// Create test source
	source := CreateTestSource()

	// Create test payload
	payload := CreateTestPayload()

	// Test cases for API methods
	apiTestCases := []struct {
		name          string
		method        string
		source        *models.Source
		country       string
		payload       map[string]interface{}
		submissionID  string
		expectError   bool
		errorType     string
		errorContains string
	}{
		{
			name:        "SubmitInvoice Valid",
			method:      "SubmitInvoice",
			source:      source,
			country:     "SA",
			payload:     payload,
			expectError: false,
		},
		{
			name:          "SubmitInvoice Nil Source",
			method:        "SubmitInvoice",
			source:        nil,
			country:       "SA",
			payload:       payload,
			expectError:   true,
			errorType:     "ValidationError",
			errorContains: "source cannot be nil",
		},
		{
			name:          "SubmitInvoice Empty Country",
			method:        "SubmitInvoice",
			source:        source,
			country:       "",
			payload:       payload,
			expectError:   true,
			errorType:     "ValidationError",
			errorContains: "country is required",
		},
		{
			name:          "SubmitInvoice Invalid Country Format",
			method:        "SubmitInvoice",
			source:        source,
			country:       "USA",
			payload:       payload,
			expectError:   true,
			errorType:     "ValidationError",
			errorContains: "country must be a 2-letter ISO code",
		},
		{
			name:          "SubmitInvoice Nil Payload",
			method:        "SubmitInvoice",
			source:        source,
			country:       "SA",
			payload:       nil,
			expectError:   true,
			errorType:     "ValidationError",
			errorContains: "payload cannot be nil",
		},
		{
			name:        "CreateMapping Valid",
			method:      "CreateMapping",
			source:      source,
			country:     "SA",
			payload:     payload,
			expectError: false,
		},
		{
			name:          "CreateMapping Nil Source",
			method:        "CreateMapping",
			source:        nil,
			country:       "SA",
			payload:       payload,
			expectError:   true,
			errorType:     "ValidationError",
			errorContains: "source cannot be nil",
		},
		{
			name:        "GetStatus Valid",
			method:      "GetStatus",
			submissionID: "test_123",
			expectError: false,
		},
		{
			name:          "GetStatus Empty ID",
			method:        "GetStatus",
			submissionID:  "",
			expectError:   true,
			errorType:     "ValidationError",
			errorContains: "submission ID is required",
		},
	}

	// Run API test cases
	for _, tc := range apiTestCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			var resp *models.UnifyResponse
			var err error

			switch tc.method {
			case "SubmitInvoice":
				resp, err = client.SubmitInvoice(ctx, tc.source, tc.country, tc.payload)
			case "CreateMapping":
				resp, err = client.CreateMapping(ctx, tc.source, tc.country, tc.payload)
			case "GetStatus":
				resp, err = client.GetStatus(ctx, tc.submissionID)
			}

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
				if tc.errorType != "" {
					switch tc.errorType {
					case "ConfigError":
						assert.True(t, errors.IsConfigError(err))
					case "ValidationError":
						assert.True(t, errors.IsValidationError(err))
					}
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, "success", resp.Status)
			}
		})
	}

	// Test batch processing
	t.Run("BatchProcess", func(t *testing.T) {
		// Create batch requests
		requests := make([]*models.UnifyRequest, 3)
		for i := 0; i < 3; i++ {
			req := models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, "SA")
			req.WithPayload(payload)
			requests[i] = req
		}

		ctx := context.Background()
		responses, errs := client.BatchProcess(ctx, requests)

		assert.Equal(t, 3, len(responses))
		assert.Equal(t, 3, len(errs))

		for i, resp := range responses {
			assert.NotNil(t, resp)
			assert.Equal(t, "success", resp.Status)
			assert.NoError(t, errs[i])
		}
	})

	// Test batch processing with errors
	t.Run("BatchProcess_WithErrors", func(t *testing.T) {
		// Create batch requests with one invalid request
		requests := make([]*models.UnifyRequest, 3)
		
		// Valid requests
		requests[0] = models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, "SA")
		requests[0].WithPayload(payload)
		
		// Invalid request (nil source)
		requests[1] = models.NewUnifyRequest(nil, models.DocumentTypeTaxInvoice, "SA")
		requests[1].WithPayload(payload)
		
		// Valid request
		requests[2] = models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, "SA")
		requests[2].WithPayload(payload)

		ctx := context.Background()
		responses, errs := client.BatchProcess(ctx, requests)

		assert.Equal(t, 3, len(responses))
		assert.Equal(t, 3, len(errs))

		// Check first request (valid)
		assert.NotNil(t, responses[0])
		assert.Equal(t, "success", responses[0].Status)
		assert.NoError(t, errs[0])

		// Check second request (invalid)
		assert.Nil(t, responses[1])
		assert.Error(t, errs[1])
		assert.True(t, errors.IsValidationError(errs[1]))

		// Check third request (valid)
		assert.NotNil(t, responses[2])
		assert.Equal(t, "success", responses[2].Status)
		assert.NoError(t, errs[2])
	})

	// Test server middleware
	t.Run("ServerMiddleware", func(t *testing.T) {
		middleware, err := pkg.NewServerMiddleware()
		assert.NoError(t, err)
		assert.NotNil(t, middleware)
	})
}