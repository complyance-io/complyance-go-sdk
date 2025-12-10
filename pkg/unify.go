/*
Unify API implementation for the GETS Unify Go SDK.

This contains the logical document type processing functionality.
*/
package complyancesdk

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)


// PushToUnify Push to Unify API with logical document types but full control over operation, mode, and purpose
func PushToUnify(
	sourceName string,
	sourceVersion string,
	logicalType LogicalDocType,
	country Country,
	operation Operation,
	mode Mode,
	purpose Purpose,
	payload map[string]interface{},
	destinations []*Destination,
) (*UnifyResponse, error) {
	if globalSDK == nil || globalSDK.config == nil {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeMissingField,
			"SDK not configured",
		))
	}

	// Process queued submissions first before handling new requests
	ProcessQueuedSubmissionsFirst()

	// Validate required parameters
	// Handle sourceName and sourceVersion based on purpose
	var finalSourceName, finalSourceVersion string
	if purpose == PurposeMapping {
		// For MAPPING purpose, sourceName and sourceVersion are optional - set to empty string if None
		if sourceName == "" {
			finalSourceName = ""
		} else {
			finalSourceName = sourceName
		}
		if sourceVersion == "" {
			finalSourceVersion = ""
		} else {
			finalSourceVersion = sourceVersion
		}
	} else {
		// For all other purposes, sourceName and sourceVersion are mandatory
		if strings.TrimSpace(sourceName) == "" {
			return nil, NewSDKError(NewErrorDetailWithCode(
				ErrorCodeMissingField,
				"Source name is required",
			))
		}
		if strings.TrimSpace(sourceVersion) == "" {
			return nil, NewSDKError(NewErrorDetailWithCode(
				ErrorCodeMissingField,
				"Source version is required",
			))
		}
		finalSourceName = sourceName
		finalSourceVersion = sourceVersion
	}

	if logicalType == "" {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeMissingField,
			"Logical document type is required",
		))
	}

	if country == "" {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeMissingField,
			"Country is required",
		))
	}

	if operation == "" {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeMissingField,
			"Operation is required",
		))
	}

	if mode == "" {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeMissingField,
			"Mode is required",
		))
	}

	if purpose == "" {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeMissingField,
			"Purpose is required",
		))
	}

	if payload == nil {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeMissingField,
			"Payload is required",
		))
	}

	// Validate country restrictions for current environment
	if err := validateCountryForEnvironment(country, globalSDK.config.Environment); err != nil {
		return nil, err
	}

	// Evaluate country policy to get base document type and meta.config flags
	policy := CountryPolicyRegistryInstance.Evaluate(country, logicalType)

	// Merge meta.config flags into payload
	mergedPayload := deepMergeIntoMetaConfig(payload, policy.GetMetaConfigFlags())

	// Auto-set invoice_data.document_type based on LogicalDocType
	setInvoiceDataDocumentType(mergedPayload, logicalType)

	// Create source reference
	sourceRef := NewSourceRef(finalSourceName, finalSourceVersion)

	// Auto-generate destinations if none provided and auto-generation is enabled
	var finalDestinations []*Destination
	if destinations == nil && globalSDK.config.AutoGenerateTaxDestination {
		finalDestinations = generateDefaultDestinations(string(country), getMetaConfigDocumentType(logicalType))
	} else {
		finalDestinations = destinations
		if finalDestinations == nil {
			finalDestinations = []*Destination{}
		}
	}

	// Build and send request using the resolved base document type
	return pushToUnifyInternalWithDocumentType(
		sourceRef, policy.GetBaseType(),
		getMetaConfigDocumentType(logicalType),
		country, operation, mode, purpose, mergedPayload, finalDestinations,
	)
}

// PushToUnifyFromJSON Push to Unify API with logical document types using JSON string payload
func PushToUnifyFromJSON(
	sourceName string,
	sourceVersion string,
	logicalType LogicalDocType,
	country Country,
	operation Operation,
	mode Mode,
	purpose Purpose,
	jsonPayload string,
	destinations []*Destination,
) (*UnifyResponse, error) {
	if strings.TrimSpace(jsonPayload) == "" {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeEmptyPayload,
			"Payload is required but was null or empty",
		).WithSuggestion(`Provide a non-empty JSON payload string. Example: '{"invoiceNumber":"INV-123","amount":1000}'`))
	}

	var payloadMap map[string]interface{}
	if err := json.Unmarshal([]byte(jsonPayload), &payloadMap); err != nil {
		errorDetail := NewErrorDetailWithCode(
			ErrorCodeMalformedJSON,
			fmt.Sprintf("Failed to parse JSON payload: %s", err.Error()),
		).WithSuggestion(`Ensure the payload is valid JSON. Example: '{"invoiceNumber":"INV-123","amount":1000}'`)
		
		// Add context for debugging
		payloadSnippet := jsonPayload
		if len(jsonPayload) > 100 {
			payloadSnippet = jsonPayload[:100] + "..."
		}
		errorDetail.AddContextValue("payloadSnippet", payloadSnippet)
		errorDetail.AddContextValue("parseError", err.Error())
		
		return nil, NewSDKError(errorDetail)
	}

	if payloadMap == nil {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeMalformedJSON,
			"Failed to parse JSON payload: parsed result is nil",
		).WithSuggestion(`Ensure the payload is valid JSON and represents an object structure. Example: '{"invoiceNumber":"INV-123"}'`))
	}

	return PushToUnify(
		sourceName, sourceVersion, logicalType, country,
		operation, mode, purpose, payloadMap, destinations,
	)
}

// PushToUnifyFromStruct Push to Unify API with logical document types using struct payload
func PushToUnifyFromStruct(
	sourceName string,
	sourceVersion string,
	logicalType LogicalDocType,
	country Country,
	operation Operation,
	mode Mode,
	purpose Purpose,
	payloadStruct interface{},
	destinations []*Destination,
) (*UnifyResponse, error) {
	if payloadStruct == nil {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeMissingField,
			"Payload struct is required but was nil",
		).WithSuggestion("Provide a valid payload struct. Example: struct{InvoiceNumber string `json:\"invoiceNumber\"`; Amount int `json:\"amount\"`}"))
	}

	// Convert struct to map[string]interface{} via JSON marshaling/unmarshaling
	jsonBytes, err := json.Marshal(payloadStruct)
	if err != nil {
		errorDetail := NewErrorDetailWithCode(
			ErrorCodeInvalidPayloadFormat,
			fmt.Sprintf("Failed to convert payload struct to JSON: %s", err.Error()),
		).WithSuggestion("Ensure the struct structure is compatible with the SDK payload format. " +
			"The struct should be JSON serializable. " +
			"Example: struct{InvoiceNumber string `json:\"invoiceNumber\"`; Amount int `json:\"amount\"`}")
		
		// Add context for debugging
		errorDetail.AddContextValue("structType", fmt.Sprintf("%T", payloadStruct))
		errorDetail.AddContextValue("conversionError", err.Error())
		
		return nil, NewSDKError(errorDetail)
	}

	var payloadMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &payloadMap); err != nil {
		errorDetail := NewErrorDetailWithCode(
			ErrorCodeInvalidPayloadFormat,
			fmt.Sprintf("Failed to convert payload struct to map: %s", err.Error()),
		).WithSuggestion("Ensure the struct structure is compatible with the SDK payload format. " +
			"The struct should be JSON serializable and deserializable to a map structure.")
		
		// Add context for debugging
		errorDetail.AddContextValue("structType", fmt.Sprintf("%T", payloadStruct))
		errorDetail.AddContextValue("conversionError", err.Error())
		
		return nil, NewSDKError(errorDetail)
	}

	if payloadMap == nil {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeInvalidPayloadFormat,
			"Failed to convert payload struct to map: conversion returned nil result",
		).WithSuggestion("Ensure the struct structure is compatible with the SDK payload format. " +
			"The struct should be convertible to a map structure."))
	}

	return PushToUnify(
		sourceName, sourceVersion, logicalType, country,
		operation, mode, purpose, payloadMap, destinations,
	)
}

// getMetaConfigDocumentType Get meta config document type
func getMetaConfigDocumentType(logicalType LogicalDocType) string {
	logicalName := string(logicalType)
	if strings.Contains(logicalName, "CREDIT_NOTE") {
		return "credit_note"
	} else if strings.Contains(logicalName, "DEBIT_NOTE") {
		return "debit_note"
	} else {
		return "tax_invoice"
	}
}

// setInvoiceDataDocumentType Automatically sets the invoice_data.document_type field based on LogicalDocType
func setInvoiceDataDocumentType(payload map[string]interface{}, logicalType LogicalDocType) {
	if payload == nil {
		return
	}

	invoiceDataRaw, exists := payload["invoice_data"]
	if !exists {
		return
	}

	invoiceData, ok := invoiceDataRaw.(map[string]interface{})
	if !ok {
		return
	}

	// Determine document type string based on LogicalDocType
	var documentType string
	logicalName := string(logicalType)
	if strings.Contains(logicalName, "CREDIT_NOTE") {
		documentType = "credit_note"
	} else if strings.Contains(logicalName, "DEBIT_NOTE") {
		documentType = "debit_note"
	} else {
		documentType = "tax_invoice" // Default for TAX_INVOICE and SIMPLIFIED_TAX_INVOICE
	}

	// Set the document_type field
	invoiceData["document_type"] = documentType
}


// deepMergeIntoMetaConfig Deep merge meta.config flags into payload. User values take precedence over policy defaults
func deepMergeIntoMetaConfig(payload map[string]interface{}, configFlags map[string]interface{}) map[string]interface{} {
	// Create a deep copy of the payload
	merged := make(map[string]interface{})
	for k, v := range payload {
		merged[k] = v
	}

	metaRaw, exists := merged["meta"]
	var meta map[string]interface{}
	if !exists {
		meta = make(map[string]interface{})
	} else {
		if m, ok := metaRaw.(map[string]interface{}); ok {
			meta = m
		} else {
			meta = make(map[string]interface{})
		}
	}

	configRaw, exists := meta["config"]
	var config map[string]interface{}
	if !exists {
		config = make(map[string]interface{})
	} else {
		if c, ok := configRaw.(map[string]interface{}); ok {
			config = c
		} else {
			config = make(map[string]interface{})
		}
	}

	// Merge config flags (user values take precedence)
	mergedConfig := make(map[string]interface{})
	for k, v := range configFlags {
		mergedConfig[k] = v
	}
	for k, v := range config {
		mergedConfig[k] = v
	}

	meta["config"] = mergedConfig
	merged["meta"] = meta

	return merged
}

// generateDefaultDestinations Generate default destinations for a country and document type
func generateDefaultDestinations(country string, documentType string) []*Destination {
	destinations := []*Destination{}

	// Auto-generate tax authority destination
	authority := getDefaultTaxAuthority(country)
	if authority != "" {
		// Convert document type to lowercase with underscores (e.g., TAX_INVOICE -> tax_invoice)
		docTypeLower := strings.ToLower(documentType)
		destinations = append(destinations, NewTaxAuthorityDestination(strings.ToUpper(country), authority, docTypeLower))
	}

	return destinations
}

// getDefaultTaxAuthority Get default tax authority for a country
func getDefaultTaxAuthority(country string) string {
	countryUpper := strings.ToUpper(country)
	switch countryUpper {
	case "SA":
		return "ZATCA"
	case "MY":
		return "LHDN"
	case "AE":
		return "FTA"
	case "SG":
		return "IRAS"
	default:
		return ""
	}
}

// pushToUnifyInternalWithDocumentType Internal method to push to Unify API with custom document type string
func pushToUnifyInternalWithDocumentType(
	sourceRef *SourceRef,
	baseDocumentType DocumentType,
	documentTypeString string,
	country Country,
	operation Operation,
	mode Mode,
	purpose Purpose,
	payload map[string]interface{},
	destinations []*Destination,
) (*UnifyResponse, error) {
	// Build UnifyRequest with custom document type string
	now := time.Now().UTC().Format(time.RFC3339)
	requestID := fmt.Sprintf("req_%d_%f", time.Now().UnixNano()/int64(time.Millisecond), rand.Float64())

	request := NewUnifyRequestBuilder().
		Source(buildSourceObject(sourceRef)).
		DocumentType(baseDocumentType).
		DocumentTypeString(documentTypeString).
		Country(string(country)).
		Operation(operation).
		Mode(mode).
		Purpose(purpose).
		Payload(payload).
		Destinations(destinations).
		APIKey(globalSDK.config.APIKey).
		RequestID(requestID).
		Timestamp(now).
		Env(mapEnvironmentToAPIValue(globalSDK.config.Environment)).
		Build()

	// Handle correlation ID
	if globalSDK.config.CorrelationID != nil {
		request.SetCorrelationID(*globalSDK.config.CorrelationID)
	}

	response, err := globalSDK.apiClient.SendUnifyRequest(request)
	if err != nil {
		if sdkErr, ok := err.(*SDKError); ok {
			log.Printf("🔥 QUEUE: SDKError caught - Error: %v, ServerError: %v, QueueManager: %v",
				sdkErr, isServerError(sdkErr), globalSDK.queueManager != nil)

			// Check if the error is a 500-range server error and queue is enabled
			if isServerError(sdkErr) && globalSDK.queueManager != nil {
				// Store the complete UnifyRequest as JSON to maintain exact API format
				completeRequestJSON, jsonErr := json.Marshal(request)
				if jsonErr != nil {
					log.Printf("Failed to convert UnifyRequest to JSON, using toString(): %v", jsonErr)
					completeRequestJSON = []byte(fmt.Sprintf("%+v", request))
				} else {
					log.Printf("🔥 QUEUE: Successfully converted complete UnifyRequest to JSON with length: %d", len(completeRequestJSON))
					log.Printf("🔥 QUEUE: Complete request JSON preview: %s", string(completeRequestJSON)[:min(200, len(completeRequestJSON))])
				}

				// Create a Source object for backward compatibility with queue
				source := NewSource(sourceRef.GetName(), sourceRef.GetVersion(), nil)

				submission := NewPayloadSubmission(
					string(completeRequestJSON), // Store complete UnifyRequest as JSON to maintain exact API format
					source,
					country,
					baseDocumentType,
				)

				log.Printf("🔥 QUEUE: Created PayloadSubmission with complete request length: %d", len(submission.GetPayload()))

				// Enqueue the failed submission for background retry
				globalSDK.queueManager.Enqueue(submission)

				// Return a response indicating the submission was queued
				queuedResponse := &UnifyResponse{
					Status:  "queued",
					Message: &[]string{fmt.Sprintf("Request failed but has been queued for retry. Submission ID: %s", *request.GetRequestID())}[0],
					Data: &UnifyResponseData{
						Submission: &SubmissionResponse{
							SubmissionID: request.GetRequestID(),
						},
					},
				}

				return queuedResponse, nil
			}

			// If not a server error or queue not available, re-throw the exception
			return nil, sdkErr
		}
		return nil, err
	}

	return response, nil
}

// isServerError determines if an SDK error represents a server error (500-range HTTP status codes).
// Only 500-range errors (500-599) should trigger queue access.
func isServerError(sdkErr *SDKError) bool {
	if sdkErr.ErrorDetail == nil {
		return false
	}

	// Check HTTP status code in context
	httpStatusObj := sdkErr.ErrorDetail.GetContextValue("httpStatus")
	if httpStatusObj != nil {
		if statusCode, ok := httpStatusObj.(int); ok {
			// Only 500-range errors (500-599) should trigger queue access
			isServerStatus := statusCode >= 500 && statusCode < 600
			if !isServerStatus {
				log.Printf("HTTP status %d detected (non 500-range) - skipping queue", statusCode)
			} else {
				log.Printf("Server error detected from HTTP status: %d", statusCode)
			}
			return isServerStatus
		} else {
			log.Printf("Invalid HTTP status format: %v", httpStatusObj)
		}
	} else {
		log.Printf("No httpStatus in ErrorDetail context, not counting as server error")
	}

	// Fallback: use error codes only when HTTP status is unavailable
	return sdkErr.ErrorDetail.Code != nil && 
		   (*sdkErr.ErrorDetail.Code == ErrorCodeInternalServerError ||
			*sdkErr.ErrorDetail.Code == ErrorCodeServiceUnavailable)
}

// buildSourceObject Build source object from SourceRef for the request
func buildSourceObject(sourceRef *SourceRef) *Source {
	source := NewSource(sourceRef.GetName(), sourceRef.GetVersion(), nil)

	// Add type if available from registry
	sourceType := getSourceTypeFromRegistry(sourceRef.GetName(), sourceRef.GetVersion())
	if sourceType != nil {
		source = NewSource(sourceRef.GetName(), sourceRef.GetVersion(), sourceType)
	}

	return source
}

// getSourceTypeFromRegistry Get source type from registry by name and version
func getSourceTypeFromRegistry(name, version string) *SourceType {
	if globalSDK != nil && globalSDK.config != nil && globalSDK.config.Sources != nil {
		for _, s := range globalSDK.config.Sources {
			if s.GetName() == name && s.GetVersion() == version {
				return s.GetSourceTypeEnum()
			}
		}
	}
	return nil
}

// mapEnvironmentToAPIValue Map Environment enum to API-expected string values
func mapEnvironmentToAPIValue(environment Environment) string {
	switch environment {
	case EnvironmentLocal, EnvironmentTest, EnvironmentStage:
		return "sandbox"
	case EnvironmentDev, EnvironmentSandbox:
		return "sandbox"
	case EnvironmentSimulation:
		return "simulation"
	case EnvironmentProduction:
		return "prod"
	default:
		return "sandbox" // Default to sandbox for safety
	}
}

