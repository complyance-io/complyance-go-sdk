package tests

import (
	"encoding/json"
	"testing"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestSourceValidation(t *testing.T) {
	// Test valid source
	source := models.NewSource("test-source", models.SourceTypeFirstParty, "Test Source")
	err := source.Validate()
	assert.NoError(t, err)

	// Test missing ID
	invalidSource := models.NewSource("", models.SourceTypeFirstParty, "Test Source")
	err = invalidSource.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "source ID is required")

	// Test missing type
	invalidSource = models.NewSource("test-source", "", "Test Source")
	err = invalidSource.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "source type is required")

	// Test missing name
	invalidSource = models.NewSource("test-source", models.SourceTypeFirstParty, "")
	err = invalidSource.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "source name is required")

	// Test invalid type
	invalidSource = models.NewSource("test-source", "INVALID_TYPE", "Test Source")
	err = invalidSource.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid source type")
}

func TestSourceMetadata(t *testing.T) {
	// Create source with metadata
	source := models.NewSource("test-source", models.SourceTypeFirstParty, "Test Source")
	source.WithVersion("1.0.0")
	source.AddMetadata("tenant_id", "tenant-123")
	source.AddMetadata("region", "us-east-1")

	// Verify metadata
	assert.Equal(t, "1.0.0", source.Version)
	assert.Equal(t, "tenant-123", source.Metadata["tenant_id"])
	assert.Equal(t, "us-east-1", source.Metadata["region"])

	// Test WithMetadata
	metadata := map[string]interface{}{
		"new_key": "new_value",
		"number":  42,
	}
	source.WithMetadata(metadata)
	assert.Equal(t, "new_value", source.Metadata["new_key"])
	assert.Equal(t, 42, source.Metadata["number"])
}

func TestUnifyRequestValidation(t *testing.T) {
	// Create valid source
	source := models.NewSource("test-source", models.SourceTypeFirstParty, "Test Source")

	// Test valid request
	req := models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, "SA")
	req.WithPayload(map[string]interface{}{
		"invoice_number": "INV-001",
	})

	err := req.Validate()
	assert.NoError(t, err)

	// Test missing source
	invalidReq := models.NewUnifyRequest(nil, models.DocumentTypeTaxInvoice, "SA")
	err = invalidReq.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "source is required")

	// Test invalid source
	invalidSource := models.NewSource("", models.SourceTypeFirstParty, "Test Source")
	invalidReq = models.NewUnifyRequest(invalidSource, models.DocumentTypeTaxInvoice, "SA")
	err = invalidReq.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid source")

	// Test missing document type
	invalidReq = models.NewUnifyRequest(source, "", "SA")
	err = invalidReq.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "document type is required")

	// Test invalid document type
	invalidReq = models.NewUnifyRequest(source, "INVALID_TYPE", "SA")
	err = invalidReq.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid document type")

	// Test missing country
	invalidReq = models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, "")
	err = invalidReq.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "country is required")

	// Test invalid country format
	invalidReq = models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, "USA")
	err = invalidReq.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "country must be a 2-letter ISO code")

	// Test missing payload
	invalidReq = models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, "SA")
	invalidReq.Payload = nil
	err = invalidReq.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "payload is required")
}

func TestUnifyRequestBuilder(t *testing.T) {
	// Create source
	source := models.NewSource("test-source", models.SourceTypeFirstParty, "Test Source")

	// Create request using builder pattern
	req := models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, "SA")
	req.WithOperation(models.OperationSingle)
	req.WithMode(models.ModeDocuments)
	req.WithPurpose(models.PurposeInvoicing)
	req.WithPayload(map[string]interface{}{
		"invoice_number": "INV-001",
		"issue_date":     "2023-01-01",
	})
	req.AddPayloadField("customer_name", "Test Customer")
	req.AddDestination("EMAIL", map[string]interface{}{
		"email": "test@example.com",
	})

	// Verify request fields
	assert.Equal(t, source, req.Source)
	assert.Equal(t, models.DocumentTypeTaxInvoice, req.DocumentType)
	assert.Equal(t, "SA", req.Country)
	assert.Equal(t, models.OperationSingle, req.Operation)
	assert.Equal(t, models.ModeDocuments, req.Mode)
	assert.Equal(t, models.PurposeInvoicing, req.Purpose)
	assert.Equal(t, "INV-001", req.Payload["invoice_number"])
	assert.Equal(t, "2023-01-01", req.Payload["issue_date"])
	assert.Equal(t, "Test Customer", req.Payload["customer_name"])
	assert.Equal(t, 1, len(req.Destinations))
	assert.Equal(t, "EMAIL", req.Destinations[0].Type)
	assert.Equal(t, "test@example.com", req.Destinations[0].Config["email"])
}

func TestUnifyResponseJSON(t *testing.T) {
	// Create a success response
	successData := map[string]interface{}{
		"submission_id": "sub_123456",
		"status":        "PROCESSED",
	}
	successResp := models.NewSuccessResponse("Document processed successfully", successData)
	successResp.WithMetadata(&models.ResponseMetadata{
		RequestID:      "req_123456",
		ProcessingTime: 150,
		TraceID:        "trace_123456",
	})

	// Serialize to JSON
	successJSON, err := json.Marshal(successResp)
	assert.NoError(t, err)

	// Deserialize from JSON
	var deserializedSuccess models.UnifyResponse
	err = json.Unmarshal(successJSON, &deserializedSuccess)
	assert.NoError(t, err)

	// Verify fields
	assert.Equal(t, "success", deserializedSuccess.Status)
	assert.Equal(t, "Document processed successfully", deserializedSuccess.Message)
	assert.Equal(t, "sub_123456", deserializedSuccess.Data["submission_id"])
	assert.Equal(t, "PROCESSED", deserializedSuccess.Data["status"])
	assert.Equal(t, "req_123456", deserializedSuccess.Metadata.RequestID)
	assert.Equal(t, int64(150), deserializedSuccess.Metadata.ProcessingTime)
	assert.Equal(t, "trace_123456", deserializedSuccess.Metadata.TraceID)

	// Create an error response
	errorDetail := models.NewErrorDetail(models.ErrorCodeValidationError, "Invalid invoice number")
	errorDetail.WithSuggestion("Provide a valid invoice number")
	errorDetail.AddContext("field", "invoice_number")
	errorResp := models.NewErrorResponse(errorDetail)

	// Serialize to JSON
	errorJSON, err := json.Marshal(errorResp)
	assert.NoError(t, err)

	// Deserialize from JSON
	var deserializedError models.UnifyResponse
	err = json.Unmarshal(errorJSON, &deserializedError)
	assert.NoError(t, err)

	// Verify fields
	assert.Equal(t, "error", deserializedError.Status)
	assert.Equal(t, "Invalid invoice number", deserializedError.Message)
	assert.Equal(t, models.ErrorCodeValidationError, deserializedError.Error.Code)
	assert.Equal(t, "Invalid invoice number", deserializedError.Error.Message)
	assert.Equal(t, "Provide a valid invoice number", deserializedError.Error.Suggestion)
	assert.Equal(t, "invoice_number", deserializedError.Error.Context["field"])
}

func TestUnifyResponseHelpers(t *testing.T) {
	// Create a success response
	successResp := models.NewSuccessResponse("Success", nil)
	assert.True(t, successResp.IsSuccess())
	assert.False(t, successResp.IsError())
	assert.Equal(t, models.ErrorCode(""), successResp.GetErrorCode())
	assert.Equal(t, "", successResp.GetErrorMessage())

	// Create an error response
	errorDetail := models.NewErrorDetail(models.ErrorCodeValidationError, "Error message")
	errorResp := models.NewErrorResponse(errorDetail)
	assert.False(t, errorResp.IsSuccess())
	assert.True(t, errorResp.IsError())
	assert.Equal(t, models.ErrorCodeValidationError, errorResp.GetErrorCode())
	assert.Equal(t, "Error message", errorResp.GetErrorMessage())
}