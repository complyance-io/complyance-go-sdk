package complyancesdk

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type rawRoundTrip func(*http.Request) (*http.Response, error)

func (r rawRoundTrip) RoundTrip(request *http.Request) (*http.Response, error) { return r(request) }

func rawTestClient(t *testing.T, contentType, selector string, input []byte, status int, output string) (*APIClient, *int) {
	t.Helper()
	calls := 0
	client := NewAPIClient("local-test-key", EnvironmentLocal, NewNoRetryConfig())
	client.httpClient = &http.Client{Transport: rawRoundTrip(func(request *http.Request) (*http.Response, error) {
		calls++
		body, err := io.ReadAll(request.Body)
		if err != nil {
			t.Fatal(err)
		}
		if request.Method != "POST" || request.URL.String() != "https://dev.gets.complyance.io/api/v3/unify" {
			t.Fatalf("Unexpected endpoint: %v", request.URL)
		}
		if request.Header.Get("Content-Type") != contentType || request.Header.Get("new-api") != selector || request.Header.Get("Authorization") != "Bearer local-test-key" {
			t.Fatal("Unexpected headers")
		}
		if !bytes.Equal(body, input) {
			t.Fatal("Request bytes changed")
		}
		return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(output)), Header: make(http.Header)}, nil
	})}
	return client, &calls
}

func TestRawSingleTransportPreservesBytes(t *testing.T) {
	xml := []byte("<?xml version=\"1.0\"?><UnifyRequest><payload>مرحبا</payload></UnifyRequest>")
	json := []byte(" {\"payload\": {\"name\":\"مرحبا\"}} \n")
	for _, mode := range []string{"xml", "file", "json-true", "json-false"} {
		t.Run(mode, func(t *testing.T) {
			input, contentType, selector := xml, "application/xml", ""
			if strings.HasPrefix(mode, "json") {
				input = json
				contentType = "application/json"
				selector = strings.TrimPrefix(mode, "json-")
			}
			client, calls := rawTestClient(t, contentType, selector, input, 200, `{"documentId":"mock-doc","unknown":{"keep":true}}`)
			var result map[string]interface{}
			var err error
			if mode == "xml" {
				result, err = client.SubmitXmlSingle(input, "https://dev.gets.complyance.io")
			}
			if mode == "file" {
				path := filepath.Join(t.TempDir(), "request.xml")
				if err = os.WriteFile(path, input, 0600); err != nil {
					t.Fatal(err)
				}
				result, err = client.SubmitXmlSingleFile(path, "https://dev.gets.complyance.io")
			}
			if strings.HasPrefix(mode, "json") {
				result, err = client.SubmitJSONSingle(input, "https://dev.gets.complyance.io", mode == "json-true")
			}
			if err != nil || result["documentId"] != "mock-doc" || result["unknown"] == nil || *calls != 1 {
				t.Fatalf("Unexpected response: %v %v calls=%d", result, err, *calls)
			}
		})
	}
}

func TestRawSingleValidationAndHTTPFailures(t *testing.T) {
	input := []byte("<UnifyRequest/>")
	validation := `{"documentId":"persisted-doc","validationStage":"business","errors":[]}`
	for _, test := range []struct {
		status int
		body   string
		accept bool
	}{
		{422, validation, true}, {422, `{"message":"invalid"}`, false}, {403, validation, false}, {500, `{"message":"failed"}`, false}, {200, `[]`, false}, {200, `null`, false}, {200, `{"results":[]}`, false},
	} {
		t.Run(fmt.Sprint(test.status, test.body), func(t *testing.T) {
			client, calls := rawTestClient(t, "application/xml", "", input, test.status, test.body)
			result, err := client.SubmitXmlSingle(input, "https://dev.gets.complyance.io")
			if (err == nil) != test.accept || *calls != 1 {
				t.Fatalf("Unexpected %v %v calls=%d", result, err, *calls)
			}
		})
	}
	for _, newAPI := range []bool{false, true} {
		jsonInput := []byte(`{}`)
		client, calls := rawTestClient(t, "application/json", fmt.Sprint(newAPI), jsonInput, 422, validation)
		_, err := client.SubmitJSONSingle(jsonInput, "https://dev.gets.complyance.io", newAPI)
		if (err == nil) != newAPI || *calls != 1 {
			t.Fatalf("422 contract mismatch: %v", err)
		}
	}
}

func TestRawLegacyEnvelopeIsPreserved(t *testing.T) {
	input := []byte(` {"action":"submit","env":"sandbox","defaults":{"country":"AE","purpose":"invoicing","source":"example:1","documentType":{"base":"tax_invoice","modifiers":[]}},"invoices":[{"payload":{"invoice_data":{"document_number":"MOCK-1"}}}]} `)
	envelope := `{"status":"success","data":{"results":[{"document":{"documentId":"legacy-doc"},"unknown":true}]}}`
	client, calls := rawTestClient(t, "application/json", "false", input, 200, envelope)
	result, err := client.SubmitJSONSingle(input, "https://dev.gets.complyance.io", false)
	if err != nil || result["data"] == nil || result["status"] != "success" || *calls != 1 {
		t.Fatalf("Legacy envelope changed: %v %v", result, err)
	}
}

func TestRawLegacyOmittedSelectorPreservesBodyAndAuth(t *testing.T) {
	input := []byte(` {"action":"submit","env":"sandbox","defaults":{"country":"AE","purpose":"invoicing","source":{"name":"example","version":"1"},"logicalDocumentType":{"base":"tax_invoice","modifiers":[]}},"invoices":[{"payload":{"invoice_data":{"document_number":"MOCK-1"}}}]} `)
	envelope := `{"status":"success","data":{"results":[{"documentId":"legacy-omitted-doc"}]}}`
	client, calls := rawTestClient(t, "application/json", "", input, 200, envelope)
	transport := client.httpClient.Transport
	client.httpClient.Transport = rawRoundTrip(func(request *http.Request) (*http.Response, error) {
		if _, present := request.Header[http.CanonicalHeaderKey("new-api")]; present {
			t.Fatal("Legacy omitted request included new-api header")
		}
		return transport.RoundTrip(request)
	})
	result, err := client.SubmitLegacyJSONSingle(input, "https://dev.gets.complyance.io")
	if err != nil || result["data"] == nil || *calls != 1 {
		t.Fatalf("Unexpected legacy response %v %v", result, err)
	}
	if _, err = client.SubmitLegacyJSONSingle([]byte(`null`), "https://dev.gets.complyance.io"); err == nil || *calls != 1 {
		t.Fatal("Invalid omitted-selector JSON sent")
	}
}

func TestRawNetworkFailureDoesNotRetry(t *testing.T) {
	calls := 0
	client := NewAPIClient("key", EnvironmentLocal, NewDefaultRetryConfig())
	client.httpClient = &http.Client{Transport: rawRoundTrip(func(request *http.Request) (*http.Response, error) {
		calls++
		return nil, fmt.Errorf("network unavailable")
	})}
	if _, err := client.SubmitJSONSingle([]byte(`{}`), "https://dev.gets.complyance.io", true); err == nil || calls != 1 {
		t.Fatalf("Unexpected retry count %d: %v", calls, err)
	}
}

func TestRawBulkPreservesIndexedResults(t *testing.T) {
	input := []byte("<UnifyRequest><invoices/></UnifyRequest>")
	output := `{"summary":{"total":2,"succeeded":1,"failed":1},"results":[{"index":0,"success":true,"result":{"payloadId":"p-1","unknown":42}},{"index":1,"success":false,"error":{"code":"invalid_argument","message":"missing"}}]}`
	for _, mode := range []string{"xml", "file", "json"} {
		t.Run(mode, func(t *testing.T) {
			request := input
			contentType, selector := "application/xml", ""
			if mode == "json" {
				request = []byte(`{"invoices":[{}]}`)
				contentType = "application/json"
				selector = "true"
			}
			client, calls := rawTestClient(t, contentType, selector, request, 200, output)
			var result *RawUnifyBulkResponse
			var err error
			switch mode {
			case "xml":
				result, err = client.SubmitXmlBulk(input, "https://dev.gets.complyance.io")
			case "file":
				path := filepath.Join(t.TempDir(), "bulk.xml")
				if err = os.WriteFile(path, input, 0600); err != nil {
					t.Fatal(err)
				}
				result, err = client.SubmitXmlBulkFile(path, "https://dev.gets.complyance.io")
			case "json":
				result, err = client.SubmitJSONBulk(request, "https://dev.gets.complyance.io")
			}
			if err != nil || *calls != 1 {
				t.Fatal(err)
			}
			if result.Summary.Total != 2 || result.Results[0].Result["payloadId"] != "p-1" || result.Results[1].Error["code"] != "invalid_argument" {
				t.Fatalf("Lost result: %+v", result)
			}
		})
	}
}

func TestRawBulkRejectsMalformedShape(t *testing.T) {
	input := []byte("xml")
	for _, body := range []string{`{}`, `{"summary":{"total":0,"succeeded":0,"failed":0}}`, `{"summary":{"total":1,"succeeded":0,"failed":1},"results":[]}`, `{"summary":{"total":1,"succeeded":1,"failed":0},"results":[{"success":true,"result":{}}]}`, `{"summary":{"total":1,"succeeded":1,"failed":0},"results":[{"index":0,"success":false,"error":{}}]}`, `{"summary":{"total":1,"succeeded":1,"failed":0},"results":[{"index":0,"success":true}]}`, `{"summary":{"total":-1,"succeeded":0,"failed":0},"results":[]}`} {
		client, _ := rawTestClient(t, "application/xml", "", input, 200, body)
		if _, err := client.SubmitXmlBulk(input, "https://dev.gets.complyance.io"); err == nil {
			t.Fatalf("Accepted malformed response %s", body)
		}
	}
}

func TestRawInputValidationDoesNotSend(t *testing.T) {
	client, calls := rawTestClient(t, "application/xml", "", []byte("xml"), 200, `{}`)
	for _, origin := range []string{"", "file:///tmp", "https://user:key@host", "https://host/path", "https://host?q=1", "https://host#fragment"} {
		if _, err := client.SubmitXmlSingle([]byte("xml"), origin); err == nil {
			t.Fatalf("Accepted invalid origin %s", origin)
		}
	}
	if _, err := client.SubmitXmlSingle(nil, "https://host"); err == nil {
		t.Fatal("Accepted empty body")
	}
	if _, err := client.SubmitXmlSingleFile(t.TempDir(), "https://host"); err == nil {
		t.Fatal("Accepted directory")
	}
	if _, err := client.SubmitXmlSingleFile(filepath.Join(t.TempDir(), "missing"), "https://host"); err == nil {
		t.Fatal("Accepted missing file")
	}
	if *calls != 0 {
		t.Fatal("Validation sent a request")
	}
	previous := globalSDK
	globalSDK = nil
	defer func() { globalSDK = previous }()
	if _, err := SubmitXmlSingle([]byte("xml"), "https://host"); err == nil {
		t.Fatal("Accepted unconfigured SDK")
	}
}

func TestRawBulkPreservesValidationFailureResults(t *testing.T) {
	output := `{"summary":{"total":1,"succeeded":0,"failed":1},"results":[{"index":0,"success":false,"result":{"documentId":"persisted-validation-doc","validationStage":"business","errors":[{"code":"rule"}],"unknown":true}}]}`
	for _, mode := range []string{"xml", "json"} {
		t.Run(mode, func(t *testing.T) {
			input, contentType, selector := []byte("xml"), "application/xml", ""
			if mode == "json" {
				input = []byte(`{"invoices":[{}]}`)
				contentType = "application/json"
				selector = "true"
			}
			client, _ := rawTestClient(t, contentType, selector, input, 200, output)
			var response *RawUnifyBulkResponse
			var err error
			if mode == "xml" {
				response, err = client.SubmitXmlBulk(input, "https://dev.gets.complyance.io")
			} else {
				response, err = client.SubmitJSONBulk(input, "https://dev.gets.complyance.io")
			}
			if err != nil {
				t.Fatal(err)
			}
			item := response.Results[0]
			if item.Success || item.Error != nil || item.Result["documentId"] != "persisted-validation-doc" || item.Result["unknown"] != true {
				t.Fatalf("Lost validation result: %+v", item)
			}
		})
	}
}

func TestRawBulkRejectsConflictingAndNullResults(t *testing.T) {
	for _, item := range []string{`{"index":0,"success":false,"result":{},"error":{}}`, `{"index":0,"success":false,"result":null}`, `{"index":0,"success":false,"error":null}`, `{"index":0,"success":true,"error":{}}`} {
		body := `{"summary":{"total":1,"succeeded":0,"failed":1},"results":[` + item + `]}`
		client, _ := rawTestClient(t, "application/xml", "", []byte("xml"), 200, body)
		if _, err := client.SubmitXmlBulk([]byte("xml"), "https://dev.gets.complyance.io"); err == nil {
			t.Fatalf("Accepted invalid result %s", item)
		}
	}
}

func TestRawSingleRejectsBulkMarkers(t *testing.T) {
	for _, status := range []int{200, 422} {
		for _, marker := range []string{`"summary":{}`, `"results":[]`} {
			body := `{"documentId":"doc","validationStage":"business","errors":[],` + marker + `}`
			client, _ := rawTestClient(t, "application/xml", "", []byte("xml"), status, body)
			if _, err := client.SubmitXmlSingle([]byte("xml"), "https://dev.gets.complyance.io"); err == nil {
				t.Fatalf("Accepted bulk marker at %d: %s", status, marker)
			}
		}
	}
}

func TestRawJSONInvalidBodiesNeverSend(t *testing.T) {
	client, calls := rawTestClient(t, "application/json", "true", nil, 200, `{}`)
	for _, body := range []string{"", `[]`, `null`, `"value"`, `{broken`, `{"invoices":[]}`} {
		if _, err := client.SubmitJSONSingle([]byte(body), "https://dev.gets.complyance.io", true); err == nil {
			t.Fatalf("Accepted single request %s", body)
		}
	}
	for _, body := range []string{"", `[]`, `null`, `{}`, `{"invoices":[]}`, `{"invoices":null}`, `{"invoices":[null]}`, `{"invoices":[1]}`} {
		if _, err := client.SubmitJSONBulk([]byte(body), "https://dev.gets.complyance.io"); err == nil {
			t.Fatalf("Accepted bulk request %s", body)
		}
	}
	if *calls != 0 {
		t.Fatalf("Invalid JSON made %d requests", *calls)
	}
}
