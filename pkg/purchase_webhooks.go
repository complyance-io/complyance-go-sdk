package complyancesdk

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ListPurchaseInvoices fetches purchase invoices from the documents API.
func ListPurchaseInvoices(filters map[string]string) (map[string]interface{}, error) {
	if globalSDK == nil || globalSDK.config == nil {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeMissingField,
			"SDK not configured",
		))
	}

	query := url.Values{}
	query.Set("type", "purchases")
	for key, value := range filters {
		if strings.TrimSpace(value) != "" {
			query.Set(key, value)
		}
	}

	return getJSON(fmt.Sprintf("/documents?%s", query.Encode()))
}

// GetPurchaseInvoice fetches a single purchase invoice from the documents API.
func GetPurchaseInvoice(id string) (map[string]interface{}, error) {
	if strings.TrimSpace(id) == "" {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeMissingField,
			"Purchase invoice id is required",
		))
	}

	return getJSON(fmt.Sprintf("/documents/%s?type=purchases", url.PathEscape(id)))
}

// VerifyWebhookSignature verifies an inbound webhook signature using constant-time comparison.
func VerifyWebhookSignature(payload string, signature string, secret string, algorithm string) (bool, error) {
	normalizedAlgorithm := strings.ToLower(strings.TrimSpace(algorithm))
	if normalizedAlgorithm == "" {
		normalizedAlgorithm = "sha256"
	}

	var expected []byte
	switch normalizedAlgorithm {
	case "sha256":
		hash := hmac.New(sha256.New, []byte(secret))
		hash.Write([]byte(payload))
		expected = hash.Sum(nil)
	case "sha512":
		hash := hmac.New(sha512.New, []byte(secret))
		hash.Write([]byte(payload))
		expected = hash.Sum(nil)
	default:
		return false, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeInvalidArgument,
			fmt.Sprintf("Unsupported webhook HMAC algorithm: %s", algorithm),
		))
	}

	provided, err := hex.DecodeString(strings.TrimSpace(signature))
	if err != nil {
		return false, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeAuthenticationFailed,
			fmt.Sprintf("Invalid webhook signature encoding: %v", err),
		))
	}

	if !hmac.Equal(expected, provided) {
		return false, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeAuthenticationFailed,
			"Webhook signature verification failed",
		))
	}

	return true, nil
}

func getJSON(path string) (map[string]interface{}, error) {
	if globalSDK == nil || globalSDK.config == nil {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeMissingField,
			"SDK not configured",
		))
	}

	request, err := http.NewRequest("GET", resolveServiceURL(path), nil)
	if err != nil {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeNetworkError,
			fmt.Sprintf("Failed to build purchase invoice request: %v", err),
		))
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", globalSDK.config.APIKey))
	request.Header.Set("X-API-Key", globalSDK.config.APIKey)

	response, err := globalSDK.apiClient.httpClient.Do(request)
	if err != nil {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeNetworkError,
			fmt.Sprintf("Network error: %v", err),
		))
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeAPIError,
			fmt.Sprintf("Failed to read purchase invoice response: %v", err),
		))
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeAPIError,
			fmt.Sprintf("Purchase invoice request failed with status %d", response.StatusCode),
		))
	}

	if len(responseBody) == 0 {
		return map[string]interface{}{}, nil
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeAPIError,
			fmt.Sprintf("Failed to parse purchase invoice response: %v", err),
		))
	}

	return parsed, nil
}

func resolveServiceURL(path string) string {
	baseURL := globalSDK.config.Environment.GetBaseURL()
	normalizedBase := strings.TrimSuffix(baseURL, "/unify")
	if strings.HasPrefix(path, "/") {
		return normalizedBase + path
	}
	return normalizedBase + "/" + path
}
