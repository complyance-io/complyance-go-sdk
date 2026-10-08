package complyancesdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// RawUnifyBulkResponse retains every indexed item without projecting domain fields.
type RawUnifyBulkResponse struct {
	Summary RawUnifyBulkSummary `json:"summary"`
	Results []RawUnifyBulkItem  `json:"results"`
}
type RawUnifyBulkSummary struct {
	Total     int `json:"total"`
	Succeeded int `json:"succeeded"`
	Failed    int `json:"failed"`
}
type RawUnifyBulkItem struct {
	Index   int                    `json:"index"`
	Success bool                   `json:"success"`
	Result  map[string]interface{} `json:"result,omitempty"`
	Error   map[string]interface{} `json:"error,omitempty"`
}

// SubmitJSONSingle sends an already-formed JSON body once with an explicit contract selector.
func (a *APIClient) SubmitJSONSingle(body []byte, origin string, newAPI bool) (map[string]interface{}, error) {
	request, err := rawJSONObject(body)
	if err != nil {
		return nil, err
	}
	if _, bulk := request["invoices"]; newAPI && bulk {
		return nil, fmt.Errorf("Single JSON request cannot contain invoices")
	}
	return a.rawSingle(body, origin, "application/json", &newAPI, newAPI)
}

// SubmitLegacyJSONSingle sends a caller-formed legacy body once without the new-api header.
func (a *APIClient) SubmitLegacyJSONSingle(body []byte, origin string) (map[string]interface{}, error) {
	if _, err := rawJSONObject(body); err != nil {
		return nil, err
	}
	return a.rawSingle(body, origin, "application/json", nil, false)
}

// SubmitJSONBulk sends the revamped JSON bulk envelope once.
func (a *APIClient) SubmitJSONBulk(body []byte, origin string) (*RawUnifyBulkResponse, error) {
	request, err := rawJSONObject(body)
	if err != nil {
		return nil, err
	}
	invoices, ok := request["invoices"].([]interface{})
	if !ok || len(invoices) == 0 {
		return nil, fmt.Errorf("Bulk JSON request requires nonempty invoices array")
	}
	for _, invoice := range invoices {
		if _, ok := invoice.(map[string]interface{}); !ok {
			return nil, fmt.Errorf("Bulk invoices must be objects")
		}
	}
	newAPI := true
	return a.rawBulk(body, origin, "application/json", &newAPI)
}
func rawJSONObject(body []byte) (map[string]interface{}, error) {
	var request map[string]interface{}
	if err := json.Unmarshal(body, &request); err != nil {
		return nil, fmt.Errorf("Request must be a JSON object: %w", err)
	}
	if request == nil {
		return nil, fmt.Errorf("Request must be a JSON object")
	}
	return request, nil
}
func (a *APIClient) SubmitXmlSingle(body []byte, origin string) (map[string]interface{}, error) {
	return a.rawSingle(body, origin, "application/xml", nil, true)
}
func (a *APIClient) SubmitXmlBulk(body []byte, origin string) (*RawUnifyBulkResponse, error) {
	return a.rawBulk(body, origin, "application/xml", nil)
}
func readUnifyXML(path string) ([]byte, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("XML path must be a regular file")
	}
	return os.ReadFile(path)
}
func (a *APIClient) SubmitXmlSingleFile(path, origin string) (map[string]interface{}, error) {
	body, err := readUnifyXML(path)
	if err != nil {
		return nil, err
	}
	return a.SubmitXmlSingle(body, origin)
}
func (a *APIClient) SubmitXmlBulkFile(path, origin string) (*RawUnifyBulkResponse, error) {
	body, err := readUnifyXML(path)
	if err != nil {
		return nil, err
	}
	return a.SubmitXmlBulk(body, origin)
}
func rawEndpoint(origin string) (string, error) {
	parsed, err := url.Parse(origin)
	if err != nil {
		return "", fmt.Errorf("Invalid public API origin")
	}
	if (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" || parsed.User != nil || (parsed.Path != "" && parsed.Path != "/") || parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", fmt.Errorf("Public API origin must contain only HTTP scheme and host")
	}
	parsed.Path = "/api/v3/unify"
	return parsed.String(), nil
}
func (a *APIClient) rawSingle(body []byte, origin, contentType string, newAPI *bool, acceptValidation bool) (map[string]interface{}, error) {
	response, err := a.sendRawUnify(body, origin, contentType, newAPI, acceptValidation)
	if err != nil {
		return nil, err
	}
	if _, bulk := response["results"]; bulk {
		return nil, fmt.Errorf("Expected a single Unify result")
	}
	if _, bulk := response["summary"]; bulk {
		return nil, fmt.Errorf("Expected a single Unify result")
	}
	return response, nil
}
func (a *APIClient) sendRawUnify(body []byte, origin, contentType string, newAPI *bool, acceptValidation bool) (map[string]interface{}, error) {
	if len(body) == 0 {
		return nil, fmt.Errorf("Request body is required")
	}
	if a.apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}
	endpoint, err := rawEndpoint(origin)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+a.apiKey)
	if newAPI != nil {
		request.Header.Set("new-api", fmt.Sprint(*newAPI))
	}
	response, err := a.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	encoded, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	parseErr := json.Unmarshal(encoded, &result)
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		if response.StatusCode != 422 || !acceptValidation || parseErr != nil || !rawValidationResult(result) {
			detail := NewErrorDetailWithCode(ErrorCodeAPIError, fmt.Sprintf("Unify request failed with HTTP %d", response.StatusCode))
			detail.AddContextValue("httpStatus", response.StatusCode)
			detail.AddContextValue("responseBody", string(encoded))
			return nil, NewSDKError(detail)
		}
	}
	if parseErr != nil {
		return nil, fmt.Errorf("Invalid Unify JSON response: %w", parseErr)
	}
	if result == nil {
		return nil, fmt.Errorf("Expected a Unify response object")
	}
	return result, nil
}
func rawValidationResult(result map[string]interface{}) bool {
	if _, bulk := result["results"]; bulk {
		return false
	}
	if _, bulk := result["summary"]; bulk {
		return false
	}
	_, document := result["documentId"].(string)
	_, stage := result["validationStage"].(string)
	_, errors := result["errors"].([]interface{})
	return document && stage && errors
}
func nonnegativeCount(value interface{}) (int, bool) {
	number, ok := value.(float64)
	if !ok || number < 0 || number > float64(int(^uint(0)>>1)) {
		return 0, false
	}
	count := int(number)
	return count, count >= 0 && float64(count) == number
}
func (a *APIClient) rawBulk(body []byte, origin, contentType string, newAPI *bool) (*RawUnifyBulkResponse, error) {
	result, err := a.sendRawUnify(body, origin, contentType, newAPI, false)
	if err != nil {
		return nil, err
	}
	summary, ok := result["summary"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Bulk response requires summary object")
	}
	items, ok := result["results"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Bulk response requires results array")
	}
	total, validTotal := nonnegativeCount(summary["total"])
	succeeded, validSucceeded := nonnegativeCount(summary["succeeded"])
	failed, validFailed := nonnegativeCount(summary["failed"])
	if !validTotal || !validSucceeded || !validFailed || total != len(items) || succeeded > total || failed != total-succeeded {
		return nil, fmt.Errorf("Invalid bulk summary counts")
	}
	bulk := &RawUnifyBulkResponse{Summary: RawUnifyBulkSummary{total, succeeded, failed}, Results: make([]RawUnifyBulkItem, len(items))}
	successes := 0
	for position, value := range items {
		item, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Bulk result must be object")
		}
		index, valid := nonnegativeCount(item["index"])
		success, validSuccess := item["success"].(bool)
		if !valid || index != position || !validSuccess {
			return nil, fmt.Errorf("Invalid bulk index or success")
		}
		entry := RawUnifyBulkItem{Index: index, Success: success}
		resultValue, hasResult := item["result"]
		errorValue, hasError := item["error"]
		if hasResult == hasError {
			return nil, fmt.Errorf("Bulk item requires exactly one of result or error")
		}
		if hasResult {
			entry.Result, ok = resultValue.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("Bulk result must be an object")
			}
		} else {
			entry.Error, ok = errorValue.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("Bulk error must be an object")
			}
		}
		if success {
			if !hasResult {
				return nil, fmt.Errorf("Successful bulk item requires result object")
			}
			successes++
		}
		bulk.Results[position] = entry
	}
	if successes != succeeded {
		return nil, fmt.Errorf("Bulk success count mismatch")
	}
	return bulk, nil
}

func rawConfiguredClient() (*APIClient, error) {
	if globalSDK == nil || globalSDK.apiClient == nil {
		return nil, fmt.Errorf("SDK must be configured")
	}
	return globalSDK.apiClient, nil
}
func SubmitJSONSingle(body []byte, origin string, newAPI bool) (map[string]interface{}, error) {
	client, err := rawConfiguredClient()
	if err != nil {
		return nil, err
	}
	return client.SubmitJSONSingle(body, origin, newAPI)
}

func SubmitLegacyJSONSingle(body []byte, origin string) (map[string]interface{}, error) {
	client, err := rawConfiguredClient()
	if err != nil {
		return nil, err
	}
	return client.SubmitLegacyJSONSingle(body, origin)
}
func SubmitJSONBulk(body []byte, origin string) (*RawUnifyBulkResponse, error) {
	client, err := rawConfiguredClient()
	if err != nil {
		return nil, err
	}
	return client.SubmitJSONBulk(body, origin)
}
func SubmitXmlSingle(body []byte, origin string) (map[string]interface{}, error) {
	client, err := rawConfiguredClient()
	if err != nil {
		return nil, err
	}
	return client.SubmitXmlSingle(body, origin)
}
func SubmitXmlBulk(body []byte, origin string) (*RawUnifyBulkResponse, error) {
	client, err := rawConfiguredClient()
	if err != nil {
		return nil, err
	}
	return client.SubmitXmlBulk(body, origin)
}
func SubmitXmlSingleFile(path, origin string) (map[string]interface{}, error) {
	client, err := rawConfiguredClient()
	if err != nil {
		return nil, err
	}
	return client.SubmitXmlSingleFile(path, origin)
}
func SubmitXmlBulkFile(path, origin string) (*RawUnifyBulkResponse, error) {
	client, err := rawConfiguredClient()
	if err != nil {
		return nil, err
	}
	return client.SubmitXmlBulkFile(path, origin)
}
