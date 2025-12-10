package tests

import (
	"encoding/json"
	"testing"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestDestinationValidation(t *testing.T) {
	// Test valid destination
	dest := models.NewDestination("EMAIL", map[string]interface{}{
		"email": "test@example.com",
	})
	err := dest.Validate()
	assert.NoError(t, err)

	// Test missing type
	invalidDest := models.NewDestination("", map[string]interface{}{
		"email": "test@example.com",
	})
	err = invalidDest.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "destination type is required")

	// Test missing config
	invalidDest = models.NewDestination("EMAIL", nil)
	err = invalidDest.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "destination config is required")

	// Test builder pattern
	dest = models.NewDestination("SMS", nil)
	dest.WithConfig(map[string]interface{}{
		"phone": "+1234567890",
	})
	err = dest.Validate()
	assert.NoError(t, err)
	assert.Equal(t, "+1234567890", dest.Config["phone"])

	// Test AddConfigField
	dest = models.NewDestination("WEBHOOK", nil)
	dest.AddConfigField("url", "https://example.com/webhook")
	dest.AddConfigField("method", "POST")
	err = dest.Validate()
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com/webhook", dest.Config["url"])
	assert.Equal(t, "POST", dest.Config["method"])
}

func TestCountryValidation(t *testing.T) {
	// Test valid country
	country := models.NewCountry(models.CountryCodeSA, "Saudi Arabia")
	err := country.Validate()
	assert.NoError(t, err)

	// Test missing code
	invalidCountry := models.NewCountry("", "Saudi Arabia")
	err = invalidCountry.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "country code is required")

	// Test invalid code length
	invalidCountry = models.NewCountry("SAU", "Saudi Arabia")
	err = invalidCountry.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "country code must be a 2-letter ISO code")

	// Test missing name
	invalidCountry = models.NewCountry(models.CountryCodeSA, "")
	err = invalidCountry.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "country name is required")

	// Test extensions
	country = models.NewCountry(models.CountryCodeSA, "Saudi Arabia")
	country.AddExtension("vat_rate", 15.0)
	country.AddExtension("currency", "SAR")
	assert.Equal(t, 15.0, country.Extensions["vat_rate"])
	assert.Equal(t, "SAR", country.Extensions["currency"])

	// Test WithExtensions
	extensions := map[string]interface{}{
		"timezone": "Asia/Riyadh",
		"language": "ar",
	}
	country.WithExtensions(extensions)
	assert.Equal(t, "Asia/Riyadh", country.Extensions["timezone"])
	assert.Equal(t, "ar", country.Extensions["language"])

	// Test String method
	str := country.String()
	assert.Contains(t, str, "Code=SA")
	assert.Contains(t, str, "Name=Saudi Arabia")
	assert.Contains(t, str, "Extensions=")
}

func TestValidationResults(t *testing.T) {
	// Create validation results
	results := models.NewValidationResults()
	
	// Add error
	results.AddError("invoice_number", "Invoice number is required")
	
	// Add warning
	results.AddWarning("issue_date", "Issue date is in the future")
	
	// Add info
	results.AddInfo("customer_name", "Customer name is recommended")
	
	// Test counts
	assert.Equal(t, 3, results.Count())
	assert.Equal(t, 1, results.ErrorCount())
	assert.Equal(t, 1, results.WarningCount())
	assert.Equal(t, 1, results.InfoCount())
	
	// Test HasErrors and HasWarnings
	assert.True(t, results.HasErrors())
	assert.True(t, results.HasWarnings())
	
	// Test individual results
	assert.Equal(t, "invoice_number", results.Results[0].Field)
	assert.Equal(t, "Invoice number is required", results.Results[0].Message)
	assert.Equal(t, models.ValidationSeverityError, results.Results[0].Severity)
	assert.True(t, results.Results[0].IsError())
	
	assert.Equal(t, "issue_date", results.Results[1].Field)
	assert.Equal(t, "Issue date is in the future", results.Results[1].Message)
	assert.Equal(t, models.ValidationSeverityWarning, results.Results[1].Severity)
	assert.True(t, results.Results[1].IsWarning())
	
	assert.Equal(t, "customer_name", results.Results[2].Field)
	assert.Equal(t, "Customer name is recommended", results.Results[2].Message)
	assert.Equal(t, models.ValidationSeverityInfo, results.Results[2].Severity)
	assert.True(t, results.Results[2].IsInfo())
	
	// Test validation result builder pattern
	result := models.NewValidationResult("total_amount", "Total amount is negative", models.ValidationSeverityError)
	result.WithCode("NEGATIVE_AMOUNT")
	result.WithPath("$.invoice.total_amount")
	result.WithValue(-100.0)
	result.WithExpected("positive number")
	
	assert.Equal(t, "total_amount", result.Field)
	assert.Equal(t, "Total amount is negative", result.Message)
	assert.Equal(t, models.ValidationSeverityError, result.Severity)
	assert.Equal(t, "NEGATIVE_AMOUNT", result.Code)
	assert.Equal(t, "$.invoice.total_amount", result.Path)
	assert.Equal(t, -100.0, result.Value)
	assert.Equal(t, "positive number", result.Expected)
}

func TestMetadataModels(t *testing.T) {
	// Test RequestMetadata
	reqMetadata := models.NewRequestMetadata()
	assert.NotEmpty(t, reqMetadata.RequestID)
	assert.NotEmpty(t, reqMetadata.Timestamp)
	assert.NotNil(t, reqMetadata.ClientInfo)
	
	reqMetadata.WithAPIKey("test-api-key")
	reqMetadata.WithEnvironment("sandbox")
	assert.Equal(t, "test-api-key", reqMetadata.APIKey)
	assert.Equal(t, "sandbox", reqMetadata.Environment)
	
	// Test ClientInfo
	clientInfo := models.NewClientInfo()
	clientInfo.WithSDKVersion("2.0.0")
	clientInfo.WithOSInfo("linux", "5.10.0")
	assert.Equal(t, "2.0.0", clientInfo.SDKVersion)
	assert.Equal(t, "linux", clientInfo.OSName)
	assert.Equal(t, "5.10.0", clientInfo.OSVersion)
	
	reqMetadata.WithClientInfo(clientInfo)
	assert.Equal(t, "2.0.0", reqMetadata.ClientInfo.SDKVersion)
	
	// Test ResponseMetadata
	respMetadata := models.NewResponseMetadata()
	assert.NotEmpty(t, respMetadata.Timestamp)
	assert.NotNil(t, respMetadata.ServerInfo)
	
	respMetadata.WithRequestID("req_12345")
	respMetadata.WithProcessingTime(150)
	respMetadata.WithTraceID("trace_12345")
	assert.Equal(t, "req_12345", respMetadata.RequestID)
	assert.Equal(t, int64(150), respMetadata.ProcessingTime)
	assert.Equal(t, "trace_12345", respMetadata.TraceID)
	
	// Test ServerInfo
	serverInfo := &models.ServerInfo{
		Version: "1.0.0",
		Region:  "us-east-1",
		NodeID:  "node-123",
	}
	respMetadata.WithServerInfo(serverInfo)
	assert.Equal(t, "1.0.0", respMetadata.ServerInfo.Version)
	assert.Equal(t, "us-east-1", respMetadata.ServerInfo.Region)
	assert.Equal(t, "node-123", respMetadata.ServerInfo.NodeID)
}

func TestUnifyRequestValidateWithResults(t *testing.T) {
	// Create valid source
	source := models.NewSource("test-source", models.SourceTypeFirstParty, "Test Source")
	
	// Test valid request
	req := models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, "SA")
	req.WithPayload(map[string]interface{}{
		"invoice_number": "INV-001",
	})
	
	results := req.ValidateWithResults()
	assert.False(t, results.HasErrors())
	assert.Equal(t, 0, results.ErrorCount())
	
	// Test invalid request
	invalidReq := models.NewUnifyRequest(nil, "", "")
	invalidReq.Operation = ""
	invalidReq.Mode = ""
	invalidReq.Purpose = ""
	invalidReq.Payload = nil
	
	results = invalidReq.ValidateWithResults()
	assert.True(t, results.HasErrors())
	assert.GreaterOrEqual(t, results.ErrorCount(), 6) // At least 6 errors
	
	// Check specific error messages
	errorFields := make(map[string]bool)
	for _, result := range results.Results {
		if result.IsError() {
			errorFields[result.Field] = true
		}
	}
	
	assert.True(t, errorFields["source"])
	assert.True(t, errorFields["document_type"])
	assert.True(t, errorFields["country"])
	assert.True(t, errorFields["operation"])
	assert.True(t, errorFields["mode"])
	assert.True(t, errorFields["purpose"])
	assert.True(t, errorFields["payload"])
	
	// Test invalid destination
	req = models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, "SA")
	req.WithPayload(map[string]interface{}{
		"invoice_number": "INV-001",
	})
	req.AddDestination("", nil)
	
	results = req.ValidateWithResults()
	assert.True(t, results.HasErrors())
	assert.True(t, errorFields["destinations[0].type"] || errorFields["destinations[0].config"])
}

func TestValidatorInterface(t *testing.T) {
	// Create valid objects
	source := models.NewSource("test-source", models.SourceTypeFirstParty, "Test Source")
	country := models.NewCountry(models.CountryCodeSA, "Saudi Arabia")
	dest := models.NewDestination("EMAIL", map[string]interface{}{
		"email": "test@example.com",
	})
	
	// Test ValidateAll with valid objects
	err := models.ValidateAll(source, country, dest)
	assert.NoError(t, err)
	
	// Test ValidateAll with invalid object
	invalidSource := models.NewSource("", "", "")
	err = models.ValidateAll(source, country, invalidSource)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "source type is required")
}

func TestJSONSerialization(t *testing.T) {
	// Create a country with extensions
	country := models.NewCountry(models.CountryCodeSA, "Saudi Arabia")
	country.AddExtension("vat_rate", 15.0)
	country.AddExtension("currency", "SAR")
	
	// Serialize to JSON
	jsonData, err := json.Marshal(country)
	assert.NoError(t, err)
	
	// Deserialize from JSON
	var deserializedCountry models.Country
	err = json.Unmarshal(jsonData, &deserializedCountry)
	assert.NoError(t, err)
	
	// Verify fields
	assert.Equal(t, "SA", deserializedCountry.Code)
	assert.Equal(t, "Saudi Arabia", deserializedCountry.Name)
	assert.Equal(t, 15.0, deserializedCountry.Extensions["vat_rate"])
	assert.Equal(t, "SAR", deserializedCountry.Extensions["currency"])
	
	// Test validation results serialization
	results := models.NewValidationResults()
	results.AddError("invoice_number", "Invoice number is required")
	results.AddWarning("issue_date", "Issue date is in the future")
	
	// Serialize to JSON
	jsonData, err = json.Marshal(results)
	assert.NoError(t, err)
	
	// Deserialize from JSON
	var deserializedResults models.ValidationResults
	err = json.Unmarshal(jsonData, &deserializedResults)
	assert.NoError(t, err)
	
	// Verify fields
	assert.Equal(t, 2, len(deserializedResults.Results))
	assert.Equal(t, "invoice_number", deserializedResults.Results[0].Field)
	assert.Equal(t, "Invoice number is required", deserializedResults.Results[0].Message)
	assert.Equal(t, models.ValidationSeverityError, deserializedResults.Results[0].Severity)
	assert.Equal(t, "issue_date", deserializedResults.Results[1].Field)
	assert.Equal(t, "Issue date is in the future", deserializedResults.Results[1].Message)
	assert.Equal(t, models.ValidationSeverityWarning, deserializedResults.Results[1].Severity)
}