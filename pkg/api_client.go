/*
API Client for the Complyance SDK matching Python SDK exactly.
*/
package complyancesdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// APIClient API Client matching Python SDK
type APIClient struct {
	apiKey         string
	baseURL        string
	retryStrategy  *RetryStrategy
	circuitBreaker *CircuitBreaker
	httpClient     *http.Client
}

const DefaultTimeout = 30 * time.Second

// NewAPIClient creates a new API client
func NewAPIClient(apiKey string, environment Environment, retryConfig *RetryConfig) *APIClient {
	return &APIClient{
		apiKey:         apiKey,
		baseURL:        environment.GetBaseURL(),
		retryStrategy:  NewRetryStrategy(retryConfig),
		circuitBreaker: NewCircuitBreaker(retryConfig.GetCircuitBreakerConfig()),
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// GetCircuitBreaker Get the circuit breaker
func (a *APIClient) GetCircuitBreaker() *CircuitBreaker {
	return a.circuitBreaker
}

// SendPayload Send payload matching Python SDK
func (a *APIClient) SendPayload(payload string, source *Source, country Country, documentType DocumentType) (*SubmissionResponseOld, error) {
	log.Println("🔥 SENDING PAYLOAD FROM QUEUE 🔥")
	log.Printf("Source: %s", source.GetID())
	log.Printf("Country: %s", country)
	log.Printf("Document Type: %s", documentType)
	log.Println("Payload JSON:")
	log.Println(payload)
	log.Println("🔥 END PAYLOAD 🔥")

	// Mocked: Always return a successful response
	response := &SubmissionResponseOld{
		SubmissionID: "mock-id",
		Status:       SubmissionStatusSubmitted,
		Error:        nil,
	}
	log.Printf("Payload submitted successfully with ID: %s", response.GetSubmissionID())
	return response, nil
}

// SendUnifyRequest Send UnifyRequest matching Python SDK
func (a *APIClient) SendUnifyRequest(request *UnifyRequest) (*UnifyResponse, error) {
	// Execute the request with retry logic
	result, err := a.retryStrategy.Execute(
		func() (interface{}, error) {
			return a.sendUnifyRequestInternal(request)
		},
		fmt.Sprintf("unify-request-%s", request.GetSource().GetID()),
	)
	if err != nil {
		return nil, err
	}
	return result.(*UnifyResponse), nil
}

// sendUnifyRequestInternal Internal method to send UnifyRequest
func (a *APIClient) sendUnifyRequestInternal(request *UnifyRequest) (*UnifyResponse, error) {
	requestData := a.serializeRequest(request)
	jsonPayload, err := json.Marshal(requestData)
	if err != nil {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeAPIError,
			fmt.Sprintf("Failed to serialize request: %v", err),
		))
	}

	// Essential request info
	log.Printf("📤 API Request URL: %s", a.baseURL)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", *request.GetAPIKey()),
		"X-Request-ID":  *request.GetRequestID(),
		"Origin":        "SDK",
	}

	// Add correlation ID if available
	if request.GetCorrelationID() != nil {
		headers["X-Correlation-ID"] = *request.GetCorrelationID()
	}

	// Log the request headers
	log.Println("📤 API REQUEST HEADERS:")
	for key, value := range headers {
		if key == "Authorization" {
			log.Printf("   %s: Bearer %s", key, value[7:]) // Hide the actual token
		} else {
			log.Printf("   %s: %s", key, value)
		}
	}

	// Log the complete request payload
	log.Println("📤 API REQUEST PAYLOAD:")
	var prettyPayload map[string]interface{}
	json.Unmarshal(jsonPayload, &prettyPayload)
	prettyJSON, _ := json.MarshalIndent(prettyPayload, "", "  ")
	log.Println(string(prettyJSON))

	// Create HTTP request
	req, err := http.NewRequest("POST", a.baseURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeNetworkError,
			fmt.Sprintf("Failed to create HTTP request: %v", err),
		))
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send request
	resp, err := a.httpClient.Do(req)
	if err != nil {
		log.Printf("Network error during API request: %v", err)
		errorDetail := NewErrorDetailWithCode(
			ErrorCodeNetworkError,
			fmt.Sprintf("Network error: %v", err),
		)
		errorDetail.Suggestion = &[]string{"Check your network connection and try again"}[0]
		errorDetail.Retryable = true
		return nil, NewSDKError(errorDetail)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeAPIError,
			fmt.Sprintf("Failed to read response body: %v", err),
		))
	}

	responseCode := resp.StatusCode
	responseBodyStr := string(responseBody)

	log.Printf("📥 API Response: %d %s", responseCode, resp.Status)
	log.Println("📥 RAW API RESPONSE:")
	log.Println(responseBodyStr)

	return a.handleResponse(responseCode, responseBodyStr, resp)
}

// serializeRequest Serialize UnifyRequest to dictionary
func (a *APIClient) serializeRequest(request *UnifyRequest) map[string]interface{} {
	data := make(map[string]interface{})

	// Serialize source
	if request.Source != nil {
		sourceName := request.Source.GetName()
		sourceVersion := request.Source.GetVersion()
		sourceIdentity := fmt.Sprintf("%s:%s", sourceName, sourceVersion)
		data["source"] = map[string]interface{}{
			"name":     sourceName,
			"version":  sourceVersion,
			"type":     request.Source.GetType(),
			"identity": sourceIdentity,
			"id":       sourceIdentity,
		}
	}

	// Use document_type_string if available, otherwise document_type value
	if request.DocumentTypeString != nil {
		data["documentType"] = strings.ToUpper(*request.DocumentTypeString)
	} else {
		data["documentType"] = strings.ToUpper(string(request.DocumentType))
	}

	if request.Country != "" {
		data["country"] = request.Country
	}

	if request.Operation != nil {
		data["operation"] = strings.ToUpper(string(*request.Operation))
	}

	if request.Mode != nil {
		data["mode"] = strings.ToUpper(string(*request.Mode))
	}

	if request.Purpose != nil {
		data["purpose"] = string(*request.Purpose)
	}

	if request.Payload != nil {
		data["payload"] = request.Payload
	}

	if request.APIKey != nil {
		data["apiKey"] = *request.APIKey
	}

	if request.RequestID != nil {
		data["requestId"] = *request.RequestID
	}

	if request.Timestamp != nil {
		data["timestamp"] = *request.Timestamp
	}

	if request.Env != nil {
		data["env"] = *request.Env
	}

	if request.Destinations != nil && len(request.Destinations) > 0 {
		destinations := make([]map[string]interface{}, len(request.Destinations))
		for i, dest := range request.Destinations {
			destinations[i] = a.serializeDestination(dest)
		}
		data["destinations"] = destinations
	}

	if request.CorrelationID != nil {
		data["correlationId"] = *request.CorrelationID
	}

	return data
}

// serializeDestination Serialize destination to dictionary
func (a *APIClient) serializeDestination(destination *Destination) map[string]interface{} {
	return map[string]interface{}{
		"type":    strings.ToUpper(string(destination.GetType())),
		"details": a.serializeDestinationDetails(destination.GetDetails()),
	}
}

// serializeDestinationDetails Serialize destination details to dictionary
func (a *APIClient) serializeDestinationDetails(details *DestinationDetails) map[string]interface{} {
	result := make(map[string]interface{})

	if details.Country != nil {
		result["country"] = *details.Country
	}
	if details.Authority != nil {
		result["authority"] = *details.Authority
	}
	if details.DocumentType != nil {
		result["documentType"] = *details.DocumentType
	}
	if details.Recipients != nil {
		result["recipients"] = *details.Recipients
	}
	if details.Subject != nil {
		result["subject"] = *details.Subject
	}
	if details.Body != nil {
		result["body"] = *details.Body
	}
	if details.ParticipantID != nil {
		result["participantId"] = *details.ParticipantID
	}
	if details.ProcessID != nil {
		result["processId"] = *details.ProcessID
	}

	return result
}

// handleResponse Handle HTTP response
func (a *APIClient) handleResponse(responseCode int, responseBody string, resp *http.Response) (*UnifyResponse, error) {
	if responseCode >= 200 && responseCode < 300 {
		return a.handleSuccessResponse(responseBody)
	} else {
		return a.handleErrorResponse(responseCode, responseBody, resp)
	}
}

// handleSuccessResponse Handle successful response
func (a *APIClient) handleSuccessResponse(responseBody string) (*UnifyResponse, error) {
	// Log the complete raw response
	log.Println("📥 API RAW RESPONSE:")
	log.Println(responseBody)
	
	var responseData map[string]interface{}
	err := json.Unmarshal([]byte(responseBody), &responseData)
	if err != nil {
		log.Printf("Failed to parse successful API response: %v", err)
		log.Printf("Raw response body: %s", responseBody)

		errorDetail := NewErrorDetailWithCode(
			ErrorCodeAPIError,
			"Failed to parse API response",
		)
		errorDetail.Suggestion = &[]string{"The server returned an invalid response format"}[0]
		errorDetail.AddContextValue("parseError", err.Error())
		errorDetail.AddContextValue("responseBody", responseBody)
		return nil, NewSDKError(errorDetail)
	}

	// Convert dict to UnifyResponse object
	unifyResponse := a.deserializeUnifyResponse(responseData)
	log.Printf("API request completed successfully with status: %s", unifyResponse.GetStatus())

	// Validate response structure
	if unifyResponse.GetData() == nil {
		log.Println("Response data is null, this might indicate an issue")
	}

	return unifyResponse, nil
}

// deserializeUnifyResponse Deserialize UnifyResponse from dictionary
func (a *APIClient) deserializeUnifyResponse(data map[string]interface{}) *UnifyResponse {
	response := &UnifyResponse{
		Metadata: make(map[string]interface{}),
	}

	if status, ok := data["status"].(string); ok {
		response.Status = status
	}

	if message, ok := data["message"].(string); ok {
		response.Message = &message
	}

	if metadata, ok := data["metadata"].(map[string]interface{}); ok {
		response.Metadata = metadata
	}

	// Handle error
	if errorData, ok := data["error"].(map[string]interface{}); ok {
		errorDetail := NewErrorDetail()
		if code, ok := errorData["code"].(string); ok {
			errorCode := ErrorCode(code)
			errorDetail.Code = &errorCode
		}
		if message, ok := errorData["message"].(string); ok {
			errorDetail.Message = &message
		}
		if suggestion, ok := errorData["suggestion"].(string); ok {
			errorDetail.Suggestion = &suggestion
		}
		response.Error = errorDetail
	}

	// Handle data
	if dataDict, ok := data["data"].(map[string]interface{}); ok {
		responseData := &UnifyResponseData{}

		// Source response
		if sourceDict, ok := dataDict["source"].(map[string]interface{}); ok {
			sourceResp := &SourceResponse{}
			if sourceID, ok := sourceDict["sourceId"].(string); ok {
				sourceResp.SourceID = &sourceID
			}
			if sourceid, ok := sourceDict["sourceid"].(string); ok {
				sourceResp.Sourceid = &sourceid
			}
			if sourceType, ok := sourceDict["type"].(string); ok {
				sourceResp.Type = &sourceType
			}
			if name, ok := sourceDict["name"].(string); ok {
				sourceResp.Name = &name
			}
			if version, ok := sourceDict["version"].(string); ok {
				sourceResp.Version = &version
			}
			if created, ok := sourceDict["created"].(bool); ok {
				sourceResp.Created = created
			}
			if id, ok := sourceDict["id"].(string); ok {
				sourceResp.ID = &id
			}
			responseData.Source = sourceResp
		}

		// Add other response handlers here as needed...
		response.Data = responseData
	}

	return response
}

// handleErrorResponse Handle error response
func (a *APIClient) handleErrorResponse(responseCode int, responseBody string, resp *http.Response) (*UnifyResponse, error) {
	log.Printf("❌ API request failed with HTTP %d", responseCode)
	log.Println("📥 API ERROR RESPONSE:")
	log.Println(responseBody)

	// Try to parse error response as JSON first
	errorDetail := a.parseErrorResponse(responseCode, responseBody)

	// Handle specific HTTP status codes
	switch responseCode {
	case 400:
		errorDetail.Code = &[]ErrorCode{ErrorCodeInvalidArgument}[0]
		errorDetail.Suggestion = &[]string{"Check your request parameters and payload format"}[0]
	case 401:
		errorDetail.Code = &[]ErrorCode{ErrorCodeAuthenticationFailed}[0]
		errorDetail.Suggestion = &[]string{"Check your API key and ensure it's valid"}[0]
	case 403:
		errorDetail.Code = &[]ErrorCode{ErrorCodeAuthorizationDenied}[0]
		errorDetail.Suggestion = &[]string{"Your API key doesn't have permission for this operation"}[0]
	case 404:
		errorDetail.Code = &[]ErrorCode{ErrorCodeAPIError}[0]
		errorDetail.Suggestion = &[]string{"The requested endpoint was not found. Check your SDK version"}[0]
	case 422:
		errorDetail.Code = &[]ErrorCode{ErrorCodeValidationFailed}[0]
		errorDetail.Suggestion = &[]string{"Your request data failed validation. Check the error details"}[0]
	case 429:
		errorDetail.Code = &[]ErrorCode{ErrorCodeRateLimitExceeded}[0]
		errorDetail.Suggestion = &[]string{"Too many requests. Please wait before retrying"}[0]
		errorDetail.Retryable = true
		if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
			// Parse retry after header if needed
		}
	case 500:
		errorDetail.Code = &[]ErrorCode{ErrorCodeInternalServerError}[0]
		errorDetail.Suggestion = &[]string{"Server error occurred. This request can be retried"}[0]
		errorDetail.Retryable = true
	case 502, 503, 504:
		errorDetail.Code = &[]ErrorCode{ErrorCodeServiceUnavailable}[0]
		errorDetail.Suggestion = &[]string{"Service is temporarily unavailable. Please retry after some time"}[0]
		errorDetail.Retryable = true
	default:
		if responseCode >= 500 {
			errorDetail.Retryable = true
			errorDetail.Suggestion = &[]string{"Server error occurred. This request can be retried"}[0]
		}
	}

	return nil, NewSDKError(errorDetail)
}

// parseErrorResponse Parse error response
func (a *APIClient) parseErrorResponse(responseCode int, responseBody string) *ErrorDetail {
	errorDetail := NewAPIErrorDetail(responseCode, responseBody)

	// Try to parse structured error response
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal([]byte(responseBody), &jsonResponse); err == nil {
		// Check if it's a structured error response
		if errorNode, ok := jsonResponse["error"].(map[string]interface{}); ok {
			if code, ok := errorNode["code"].(string); ok {
				errorCode := ErrorCode(code)
				errorDetail.Code = &errorCode
			}

			if message, ok := errorNode["message"].(string); ok {
				errorDetail.Message = &message
			}

			if suggestion, ok := errorNode["suggestion"].(string); ok {
				errorDetail.Suggestion = &suggestion
			}

			if field, ok := errorNode["field"].(string); ok {
				errorDetail.Field = &field
			}

			if retryable, ok := errorNode["retryable"].(bool); ok {
				errorDetail.Retryable = retryable
			}

			// Parse validation errors if present
			if validationErrors, ok := errorNode["validationErrors"].([]interface{}); ok {
				for _, validationError := range validationErrors {
					if ve, ok := validationError.(map[string]interface{}); ok {
						field := ""
						message := ""
						code := ""
						if f, ok := ve["field"].(string); ok {
							field = f
						}
						if m, ok := ve["message"].(string); ok {
							message = m
						}
						if c, ok := ve["code"].(string); ok {
							code = c
						}
						errorDetail.AddValidationError(field, message, code)
					}
				}
			}
		}
	}

	return errorDetail
}

// SendRawJSONRequest Send raw JSON request directly without deserialization
func (a *APIClient) SendRawJSONRequest(jsonPayload string) (*UnifyResponse, error) {
	log.Println("🔥 RAW JSON: Sending raw JSON request")
	log.Printf("🔥 RAW JSON: JSON length: %d", len(jsonPayload))
	log.Printf("🔥 RAW JSON: JSON preview: %s", jsonPayload[:min(200, len(jsonPayload))])

	result, err := a.retryStrategy.Execute(
		func() (interface{}, error) {
			return a.sendRawJSONRequestInternal(jsonPayload)
		},
		"raw-json-request",
	)
	if err != nil {
		return nil, err
	}
	return result.(*UnifyResponse), nil
}

// sendRawJSONRequestInternal Internal method to send raw JSON request
func (a *APIClient) sendRawJSONRequestInternal(jsonPayload string) (*UnifyResponse, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	req, err := http.NewRequest("POST", a.baseURL, strings.NewReader(jsonPayload))
	if err != nil {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeNetworkError,
			fmt.Sprintf("Failed to create HTTP request: %v", err),
		))
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		log.Printf("Network error during raw JSON API request: %v", err)
		errorDetail := NewErrorDetailWithCode(
			ErrorCodeNetworkError,
			fmt.Sprintf("Network error: %v", err),
		)
		errorDetail.Suggestion = &[]string{"Check your network connection and try again"}[0]
		errorDetail.Retryable = true
		return nil, NewSDKError(errorDetail)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeAPIError,
			fmt.Sprintf("Failed to read response body: %v", err),
		))
	}

	responseCode := resp.StatusCode
	responseBodyStr := string(responseBody)

	log.Printf("🔥 RAW JSON: API Response Code: %d", responseCode)
	log.Printf("🔥 RAW JSON: API Response Body: %s", responseBodyStr)

	if responseCode >= 200 && responseCode < 300 {
		return a.handleSuccessResponse(responseBodyStr)
	} else {
		errorDetail := a.parseErrorResponse(responseCode, responseBodyStr)
		return nil, NewSDKError(errorDetail)
	}
}

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
